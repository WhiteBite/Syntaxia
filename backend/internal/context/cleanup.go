package context

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync/atomic"
	"time"
)

// CleanupOldStreams removes old streaming contexts
func (s *Service) CleanupOldStreams(maxAge time.Duration) error {
	s.streamsMu.Lock()
	defer s.streamsMu.Unlock()

	now := time.Now()
	for id, stream := range s.streams {
		if now.Sub(stream.CreatedAt) > maxAge {
			if err := os.Remove(stream.contextPath); err != nil && !os.IsNotExist(err) {
				s.logger.Warning(fmt.Sprintf("Failed to remove old context file %s: %v", stream.contextPath, err))
			}
			delete(s.streams, id)
			s.logger.Info(fmt.Sprintf("Cleaned up old streaming context: %s", id))
		}
	}

	// Limit active streams count
	const maxActiveStreams = 10
	if len(s.streams) > maxActiveStreams {
		type streamAge struct {
			id  string
			age time.Time
		}
		var ages []streamAge
		for id, stream := range s.streams {
			ages = append(ages, streamAge{id: id, age: stream.CreatedAt})
		}
		sort.Slice(ages, func(i, j int) bool {
			return ages[i].age.Before(ages[j].age)
		})

		// Remove oldest
		for i := 0; i < len(ages)-maxActiveStreams; i++ {
			id := ages[i].id
			if stream, exists := s.streams[id]; exists {
				os.Remove(stream.contextPath)
				delete(s.streams, id)
			}
		}
	}

	return nil
}

// periodicCleanup runs periodic cleanup of old contexts
func (s *Service) periodicCleanup() {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdownCh:
			s.logger.Info("Periodic cleanup stopped due to shutdown")
			return
		case <-ticker.C:
			if err := s.CleanupOldStreams(24 * time.Hour); err != nil {
				s.logger.Warning(fmt.Sprintf("Failed to cleanup old streams: %v", err))
			}
			s.lastCleanup = time.Now()

			// Force GC after cleanup to release memory
			runtime.GC()
		}
	}
}

// GetMemoryStats returns memory usage statistics
func (s *Service) GetMemoryStats() map[string]interface{} {
	s.streamsMu.RLock()
	defer s.streamsMu.RUnlock()

	var totalSize int64
	for _, stream := range s.streams {
		if info, err := os.Stat(stream.contextPath); err == nil {
			totalSize += info.Size()
		}
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"active_streams":     len(s.streams),
		"total_disk_size_mb": totalSize / (1024 * 1024),
		"last_cleanup":       s.lastCleanup,
		"active_operations":  atomic.LoadInt64(&s.activeOperations),
		"total_operations":   atomic.LoadInt64(&s.totalOperations),
		"total_bytes_read":   atomic.LoadInt64(&s.totalBytesRead),
		"worker_count":       s.workerCount,
		"heap_alloc_mb":      memStats.HeapAlloc / (1024 * 1024),
		"heap_sys_mb":        memStats.HeapSys / (1024 * 1024),
		"num_gc":             memStats.NumGC,
		"goroutines":         runtime.NumGoroutine(),
	}
}
