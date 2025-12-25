package integration

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"
)

// BenchmarkContextGeneration benchmarks context generation performance
func BenchmarkContextGeneration(b *testing.B) {
	// Skip if not running benchmarks
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate context generation workload
		data := make([]byte, 1024*1024) // 1MB
		for j := range data {
			data[j] = byte(j % 256)
		}
		_ = data
	}
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Test buffer pool efficiency
		buf := make([]byte, 64*1024)
		_ = buf
	}
}

// BenchmarkConcurrentOperations benchmarks concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	b.ReportAllocs()

	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			// Simulate work
			time.Sleep(time.Microsecond)
		}()
	}

	wg.Wait()
}

// TestMemoryLeakDetection tests for memory leaks
func TestMemoryLeakDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping leak detection in short mode")
	}

	// Get initial memory stats
	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Run operations that might leak
	for i := 0; i < 1000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_ = ctx
		cancel() // Must cancel to prevent leak
	}

	// Force GC and get final stats
	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// Check for significant memory growth
	heapGrowth := int64(m2.HeapAlloc) - int64(m1.HeapAlloc)
	if heapGrowth > 10*1024*1024 { // 10MB threshold
		t.Errorf("Potential memory leak detected: heap grew by %d bytes", heapGrowth)
	}
}

// TestGoroutineLeakDetection tests for goroutine leaks
func TestGoroutineLeakDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping goroutine leak detection in short mode")
	}

	initialGoroutines := runtime.NumGoroutine()

	// Run operations that might leak goroutines
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond)
		}()
	}
	wg.Wait()

	// Give time for goroutines to clean up
	time.Sleep(100 * time.Millisecond)
	runtime.GC()

	finalGoroutines := runtime.NumGoroutine()
	leaked := finalGoroutines - initialGoroutines

	if leaked > 5 { // Allow small variance
		t.Errorf("Potential goroutine leak: started with %d, ended with %d (leaked %d)",
			initialGoroutines, finalGoroutines, leaked)
	}
}

// TestChannelLeakDetection tests for channel leaks
func TestChannelLeakDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping channel leak detection in short mode")
	}

	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Create and properly close channels
	for i := 0; i < 1000; i++ {
		ch := make(chan struct{}, 10)
		close(ch)
	}

	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// Channels should be garbage collected
	if m2.HeapObjects > m1.HeapObjects+1000 {
		t.Errorf("Potential channel leak: heap objects grew from %d to %d",
			m1.HeapObjects, m2.HeapObjects)
	}
}

// BenchmarkBufferPool benchmarks buffer pool performance
func BenchmarkBufferPool(b *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 64*1024)
			return &buf
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := pool.Get().(*[]byte)
		// Use buffer
		(*buf)[0] = byte(i)
		pool.Put(buf)
	}
}

// BenchmarkWithoutPool benchmarks without buffer pool for comparison
func BenchmarkWithoutPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := make([]byte, 64*1024)
		buf[0] = byte(i)
		_ = buf
	}
}

// TestContextCancellation tests proper context cancellation
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Error("Context cancellation did not propagate")
	}
}

// TestTimeoutPropagation tests timeout propagation
func TestTimeoutPropagation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got %v", ctx.Err())
		}
	case <-time.After(time.Second):
		t.Error("Timeout did not trigger")
	}
}
