// Package daemon manages Chrome process lifecycle and connection pooling.
package daemon

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/chromedp/chromedp"
)

// Manager handles Chrome daemon lifecycle and connection management.
type Manager struct {
	mu          sync.RWMutex
	chromeCmd   *exec.Cmd
	allocCtx    context.Context
	allocCancel context.CancelFunc
	idleTimer   *time.Timer
	idleTimeout time.Duration
	isRunning   bool
	debugPort   int
	chromePID   int
}

// NewManager creates a new Chrome daemon manager.
func NewManager() *Manager {
	timeout := getIdleTimeout()
	return &Manager{
		idleTimeout: timeout,
		debugPort:   9222, // Default Chrome remote debugging port
	}
}

// GetContext returns a browser context, starting the daemon if needed.
func (m *Manager) GetContext(_ context.Context) (context.Context, context.CancelFunc, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if we need to start or reconnect
	if !m.isRunning {
		// Try to reconnect to existing Chrome process first
		if m.chromePID != 0 && m.processExists(m.chromePID) {
			if err := m.reconnect(); err != nil {
				// Reconnection failed, start new Chrome
				if err := m.start(); err != nil {
					return nil, nil, err
				}
			}
		} else {
			// Start new Chrome process
			if err := m.start(); err != nil {
				return nil, nil, err
			}
		}
	}

	// Reset idle timer
	m.resetIdleTimer()

	// Create new browser context for this operation
	browserCtx, cancel := chromedp.NewContext(m.allocCtx)
	return browserCtx, cancel, nil
}

// reconnect attempts to reconnect to an existing Chrome process.
func (m *Manager) reconnect() error {
	// Create chromedp allocator that connects to the running Chrome
	m.allocCtx, m.allocCancel = chromedp.NewRemoteAllocator(
		context.Background(),
		fmt.Sprintf("ws://localhost:%d", m.debugPort),
	)

	// Test connection
	testCtx, testCancel := chromedp.NewContext(m.allocCtx)
	defer testCancel()

	// Run a simple command to verify connection
	err := chromedp.Run(testCtx, chromedp.Navigate("about:blank"))
	if err != nil {
		return fmt.Errorf("failed to reconnect to Chrome: %w", err)
	}

	m.isRunning = true
	return nil
}

// processExists checks if a process with given PID exists.
func (m *Manager) processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Try to send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// start initializes the Chrome daemon process.
func (m *Manager) start() error {
	// Find Chrome executable
	chromePath, err := m.findChrome()
	if err != nil {
		return fmt.Errorf("failed to find Chrome: %w", err)
	}

	// Start Chrome with remote debugging
	args := []string{
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-renderer-backgrounding",
		"--disable-features=VizDisplayCompositor",
		fmt.Sprintf("--remote-debugging-port=%d", m.debugPort),
		"--user-data-dir=/tmp/essenz-chrome-profile",
		"about:blank",
	}

	m.chromeCmd = exec.Command(chromePath, args...)
	m.chromeCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Create new process group
		Setsid:  true, // Create new session
	}

	// Detach from parent process completely
	m.chromeCmd.Stdin = nil
	m.chromeCmd.Stdout = nil
	m.chromeCmd.Stderr = nil

	// Start Chrome process
	err = m.chromeCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start Chrome: %w", err)
	}

	m.chromePID = m.chromeCmd.Process.Pid

	// Detach from the process - don't wait for it
	go func() {
		_ = m.chromeCmd.Wait()
	}()

	// Wait a moment for Chrome to start
	time.Sleep(2 * time.Second)

	// Create chromedp allocator that connects to the running Chrome
	m.allocCtx, m.allocCancel = chromedp.NewRemoteAllocator(
		context.Background(),
		fmt.Sprintf("ws://localhost:%d", m.debugPort),
	)

	// Test connection
	testCtx, testCancel := chromedp.NewContext(m.allocCtx)
	defer testCancel()

	// Run a simple command to verify connection
	err = chromedp.Run(testCtx, chromedp.Navigate("about:blank"))
	if err != nil {
		_ = m.chromeCmd.Process.Kill()
		return fmt.Errorf("failed to connect to Chrome: %w", err)
	}

	m.isRunning = true
	return nil
}

// findChrome locates the Chrome executable
func (m *Manager) findChrome() (string, error) {
	// Check environment variable first
	if chromePath := os.Getenv("ESSENZ_CHROME_PATH"); chromePath != "" {
		return chromePath, nil
	}

	// Common Chrome locations
	paths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"/usr/bin/google-chrome",
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("Chrome not found in common locations")
}

// resetIdleTimer resets the idle timeout timer.
func (m *Manager) resetIdleTimer() {
	if m.idleTimer != nil {
		m.idleTimer.Stop()
	}

	m.idleTimer = time.AfterFunc(m.idleTimeout, func() {
		m.shutdownWithKill()
	})
}

// shutdown gracefully shuts down the Chrome daemon without killing process.
func (m *Manager) shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return
	}

	if m.idleTimer != nil {
		m.idleTimer.Stop()
		m.idleTimer = nil
	}

	if m.allocCancel != nil {
		m.allocCancel()
		m.allocCancel = nil
	}

	// Don't kill Chrome process - let it persist
	m.isRunning = false
	// Keep chromePID for process monitoring
}

// shutdownWithKill shuts down and kills the Chrome process.
func (m *Manager) shutdownWithKill() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return
	}

	if m.idleTimer != nil {
		m.idleTimer.Stop()
		m.idleTimer = nil
	}

	if m.allocCancel != nil {
		m.allocCancel()
		m.allocCancel = nil
	}

	// Kill Chrome process on idle timeout
	if m.chromeCmd != nil && m.chromeCmd.Process != nil {
		_ = m.chromeCmd.Process.Kill()
		m.chromeCmd = nil
	}

	m.isRunning = false
	m.chromePID = 0
}

// Shutdown manually shuts down the daemon.
func (m *Manager) Shutdown() {
	m.shutdown()
}

// IsRunning returns true if the Chrome daemon is running
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if process is still alive
	if m.chromeCmd != nil && m.chromeCmd.Process != nil {
		// Try to send signal 0 to check if process exists
		err := m.chromeCmd.Process.Signal(syscall.Signal(0))
		if err != nil {
			// Process is dead, update state
			m.isRunning = false
			return false
		}
	}

	return m.isRunning
}

// getIdleTimeout returns the idle timeout from environment or default.
func getIdleTimeout() time.Duration {
	if timeoutStr := os.Getenv("ESSENZ_DAEMON_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			return timeout
		}
	}
	return 300 * time.Second // Default 300 seconds (5 minutes)
}
