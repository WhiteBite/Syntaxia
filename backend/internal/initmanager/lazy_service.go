package initmanager

import (
	"context"
	"sync"
	"time"
)

// LazyService represents a service that can be lazily initialized
type LazyService[T any] struct {
	mu           sync.RWMutex
	service      T
	initialized  bool
	initFunc     func(context.Context) (T, error)
	lastAccessed time.Time
	accessCount  int64
	initTime     time.Time
}

// NewLazyService creates a new lazy service
func NewLazyService[T any](initFunc func(context.Context) (T, error)) *LazyService[T] {
	return &LazyService[T]{
		initFunc: initFunc,
	}
}

// Get returns the service instance, initializing it if necessary
func (ls *LazyService[T]) Get(ctx context.Context) (T, error) {
	ls.mu.RLock()
	if ls.initialized {
		service := ls.service
		ls.mu.RUnlock()

		// Update access tracking
		ls.mu.Lock()
		ls.lastAccessed = time.Now()
		ls.accessCount++
		ls.mu.Unlock()

		return service, nil
	}
	ls.mu.RUnlock()

	// Double-checked locking pattern
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if ls.initialized {
		ls.lastAccessed = time.Now()
		ls.accessCount++
		return ls.service, nil
	}

	startTime := time.Now()
	service, err := ls.initFunc(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	ls.service = service
	ls.initialized = true
	ls.initTime = time.Now()
	ls.lastAccessed = time.Now()
	ls.accessCount = 1

	initDuration := time.Since(startTime)
	if initDuration > 100*time.Millisecond {
		// Log slow initializations
		println("Lazy service initialized in", initDuration.String())
	}

	return service, nil
}

// IsInitialized returns true if the service has been initialized
func (ls *LazyService[T]) IsInitialized() bool {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.initialized
}

// Reset resets the service to uninitialized state
func (ls *LazyService[T]) Reset() {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	var zero T
	ls.service = zero
	ls.initialized = false
	ls.lastAccessed = time.Time{}
	ls.accessCount = 0
}

// GetStats returns statistics about the lazy service
func (ls *LazyService[T]) GetStats() map[string]interface{} {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	stats := map[string]interface{}{
		"initialized":  ls.initialized,
		"access_count": ls.accessCount,
	}

	if ls.initialized {
		stats["init_time"] = ls.initTime
		stats["last_accessed"] = ls.lastAccessed
		stats["idle_duration"] = time.Since(ls.lastAccessed).String()
	}

	return stats
}

// ShouldUnload returns true if the service should be unloaded based on idle time
func (ls *LazyService[T]) ShouldUnload(idleThreshold time.Duration) bool {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	if !ls.initialized {
		return false
	}

	return time.Since(ls.lastAccessed) > idleThreshold
}

// LazyServiceManager manages multiple lazy services
type LazyServiceManager struct {
	mu       sync.RWMutex
	services map[string]interface{} // map of service name to LazyService
}

// NewLazyServiceManager creates a new lazy service manager
func NewLazyServiceManager() *LazyServiceManager {
	return &LazyServiceManager{
		services: make(map[string]interface{}),
	}
}

// Register registers a lazy service
func (m *LazyServiceManager) Register(name string, service interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.services[name] = service
}

// GetInitializationStats returns stats for all registered services
func (m *LazyServiceManager) GetInitializationStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})
	for name, svc := range m.services {
		// Try to get stats if the service supports it
		if statsGetter, ok := svc.(interface{ GetStats() map[string]interface{} }); ok {
			stats[name] = statsGetter.GetStats()
		}
	}

	return stats
}

// UnloadUnusedServices unloads services that haven't been accessed recently
func (m *LazyServiceManager) UnloadUnusedServices(idleTime time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	unloaded := 0
	for name, svc := range m.services {
		// Try to unload if the service supports it
		if unloader, ok := svc.(interface{ ShouldUnload(time.Duration) bool }); ok {
			if unloader.ShouldUnload(idleTime) {
				if resetter, ok := svc.(interface{ Reset() }); ok {
					resetter.Reset()
					unloaded++
					println("Unloaded idle service:", name)
				}
			}
		}
	}

	return unloaded
}
