// Package daemon provides a global Chrome daemon instance.
package daemon

import (
	"context"
	"sync"
)

var (
	globalManager *Manager
	globalMutex   sync.Mutex
)

// GetGlobalManager returns the global Chrome daemon manager, creating it if necessary.
func GetGlobalManager() *Manager {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if globalManager == nil {
		globalManager = NewManager()
	}
	return globalManager
}

// GetGlobalContext returns a browser context from the global daemon.
func GetGlobalContext(ctx context.Context) (context.Context, context.CancelFunc, error) {
	manager := GetGlobalManager()
	return manager.GetContext(ctx)
}

// ShutdownGlobal shuts down the global daemon manager.
func ShutdownGlobal() {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if globalManager != nil {
		globalManager.Shutdown()
		globalManager = nil
	}
}
