// Package daemon provides a server that manages Chrome processes independently.
package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

// Server manages Chrome processes as a long-running daemon.
type Server struct {
	mu          sync.RWMutex
	manager     *Manager
	listener    net.Listener
	socketPath  string
	isRunning   bool
	stopChannel chan struct{}
}

// Request represents a client request to the daemon.
type Request struct {
	Action string `json:"action"`
	URL    string `json:"url,omitempty"`
}

// Response represents the daemon's response.
type Response struct {
	Success bool   `json:"success"`
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewServer creates a new daemon server.
func NewServer() *Server {
	socketPath := filepath.Join(os.TempDir(), "essenz-daemon.sock")
	return &Server{
		manager:     NewManager(),
		socketPath:  socketPath,
		stopChannel: make(chan struct{}),
	}
}

// Start starts the daemon server.
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("daemon already running")
	}

	// Remove existing socket file
	_ = os.Remove(s.socketPath)

	// Create Unix socket listener
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}

	s.listener = listener
	s.isRunning = true

	log.Printf("Daemon started, listening on %s", s.socketPath)

	// Start accepting connections
	go s.acceptConnections()

	return nil
}

// Stop stops the daemon server.
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	close(s.stopChannel)
	_ = s.listener.Close()
	s.manager.Shutdown()
	_ = os.Remove(s.socketPath)
	s.isRunning = false

	log.Printf("Daemon stopped")
	return nil
}

// acceptConnections handles incoming client connections.
func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.stopChannel:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				// Check if we're shutting down
				select {
				case <-s.stopChannel:
					return
				default:
					log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			go s.handleConnection(conn)
		}
	}
}

// handleConnection processes a single client connection.
func (s *Server) handleConnection(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req Request
	if err := decoder.Decode(&req); err != nil {
		s.sendError(encoder, "Invalid request format")
		return
	}

	switch req.Action {
	case "fetch":
		s.handleFetch(encoder, req.URL)
	case "ping":
		s.sendResponse(encoder, Response{Success: true})
	case "shutdown":
		s.sendResponse(encoder, Response{Success: true})
		go func() { _ = s.Stop() }()
	default:
		s.sendError(encoder, "Unknown action: "+req.Action)
	}
}

// handleFetch processes a fetch request.
func (s *Server) handleFetch(encoder *json.Encoder, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get browser context from manager
	browserCtx, browserCancel, err := s.manager.GetContext(ctx)
	if err != nil {
		s.sendError(encoder, "Failed to get browser context: "+err.Error())
		return
	}
	defer browserCancel()

	// Use chromedp directly to fetch content
	content, err := s.fetchContentWithContext(browserCtx, url)
	if err != nil {
		s.sendError(encoder, "Failed to fetch content: "+err.Error())
		return
	}

	s.sendResponse(encoder, Response{
		Success: true,
		Content: content,
	})
}

// sendResponse sends a successful response.
func (s *Server) sendResponse(encoder *json.Encoder, resp Response) {
	if err := encoder.Encode(resp); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

// sendError sends an error response.
func (s *Server) sendError(encoder *json.Encoder, errMsg string) {
	s.sendResponse(encoder, Response{
		Success: false,
		Error:   errMsg,
	})
}

// IsDaemonRunning checks if the daemon is running by attempting to connect.
func IsDaemonRunning() bool {
	socketPath := filepath.Join(os.TempDir(), "essenz-daemon.sock")
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// fetchContentWithContext fetches content using an existing browser context.
func (s *Server) fetchContentWithContext(ctx context.Context, url string) (string, error) {
	// Set timeout for the operation
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
	defer timeoutCancel()

	// Fetch page content
	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return "", fmt.Errorf("failed to fetch content from %s: %w", url, err)
	}

	return htmlContent, nil
}

// StartDaemonIfNeeded starts the daemon if it's not already running.
func StartDaemonIfNeeded() error {
	if IsDaemonRunning() {
		return nil
	}

	server := NewServer()
	return server.Start()
}
