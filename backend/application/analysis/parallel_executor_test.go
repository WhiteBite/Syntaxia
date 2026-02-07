package analysis

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParallelExecutor_Execute(t *testing.T) {
	tests := []struct {
		name       string
		tasks      []Task[int]
		wantLen    int
		wantErrors int
	}{
		{
			name:       "empty tasks",
			tasks:      []Task[int]{},
			wantLen:    0,
			wantErrors: 0,
		},
		{
			name: "single task",
			tasks: []Task[int]{
				{ID: "1", Execute: func(ctx context.Context) (int, error) { return 42, nil }},
			},
			wantLen:    1,
			wantErrors: 0,
		},
		{
			name: "multiple tasks",
			tasks: []Task[int]{
				{ID: "1", Execute: func(ctx context.Context) (int, error) { return 1, nil }},
				{ID: "2", Execute: func(ctx context.Context) (int, error) { return 2, nil }},
				{ID: "3", Execute: func(ctx context.Context) (int, error) { return 3, nil }},
			},
			wantLen:    3,
			wantErrors: 0,
		},
		{
			name: "tasks with errors",
			tasks: []Task[int]{
				{ID: "1", Execute: func(ctx context.Context) (int, error) { return 1, nil }},
				{ID: "2", Execute: func(ctx context.Context) (int, error) { return 0, errors.New("error") }},
				{ID: "3", Execute: func(ctx context.Context) (int, error) { return 3, nil }},
			},
			wantLen:    3,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewParallelExecutor[int](DefaultExecutorConfig())
			results := executor.Execute(context.Background(), tt.tasks)

			assert.Len(t, results, tt.wantLen)

			errorCount := 0
			for _, r := range results {
				if r.Error != nil {
					errorCount++
				}
			}
			assert.Equal(t, tt.wantErrors, errorCount)
		})
	}
}

func TestParallelExecutor_ResultOrder(t *testing.T) {
	executor := NewParallelExecutor[int](ExecutorConfig{MaxWorkers: 4})

	tasks := make([]Task[int], 10)
	for i := 0; i < 10; i++ {
		idx := i
		tasks[i] = Task[int]{
			ID: string(rune('0' + idx)),
			Execute: func(ctx context.Context) (int, error) {
				time.Sleep(time.Duration(10-idx) * time.Millisecond) // Varying delays
				return idx, nil
			},
		}
	}

	results := executor.Execute(context.Background(), tasks)

	// Results should be in same order as input tasks
	for i, r := range results {
		assert.Equal(t, i, r.Value, "Result order should match task order")
	}
}

func TestParallelExecutor_Timeout(t *testing.T) {
	executor := NewParallelExecutor[int](ExecutorConfig{
		MaxWorkers: 2,
		Timeout:    50 * time.Millisecond,
	})

	tasks := []Task[int]{
		{
			ID: "fast",
			Execute: func(ctx context.Context) (int, error) {
				return 1, nil
			},
		},
		{
			ID: "slow",
			Execute: func(ctx context.Context) (int, error) {
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				case <-time.After(200 * time.Millisecond):
					return 2, nil
				}
			},
		},
	}

	results := executor.Execute(context.Background(), tasks)

	assert.Len(t, results, 2)
	assert.NoError(t, results[0].Error)
	assert.Error(t, results[1].Error) // Should timeout
}

func TestParallelExecutor_ContextCancellation(t *testing.T) {
	executor := NewParallelExecutor[int](ExecutorConfig{MaxWorkers: 2})

	ctx, cancel := context.WithCancel(context.Background())

	var started int32
	tasks := []Task[int]{
		{
			ID: "1",
			Execute: func(ctx context.Context) (int, error) {
				atomic.AddInt32(&started, 1)
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				case <-time.After(500 * time.Millisecond):
					return 1, nil
				}
			},
		},
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	results := executor.Execute(ctx, tasks)

	assert.Len(t, results, 1)
	assert.Error(t, results[0].Error)
}

func TestParallelExecutor_ExecuteWithProgress(t *testing.T) {
	executor := NewParallelExecutor[int](ExecutorConfig{MaxWorkers: 2})

	tasks := make([]Task[int], 5)
	for i := 0; i < 5; i++ {
		idx := i
		tasks[i] = Task[int]{
			ID:      string(rune('0' + idx)),
			Execute: func(ctx context.Context) (int, error) { return idx, nil },
		}
	}

	var progressCalls []int
	var mu sync.Mutex
	results := executor.ExecuteWithProgress(context.Background(), tasks, func(completed, total int) {
		mu.Lock()
		progressCalls = append(progressCalls, completed)
		mu.Unlock()
	})

	assert.Len(t, results, 5)
	// Progress should be called at least once and at most 5 times (due to parallel execution)
	assert.GreaterOrEqual(t, len(progressCalls), 1)
	assert.LessOrEqual(t, len(progressCalls), 5)
	// Last progress call should be 5 (all tasks completed)
	mu.Lock()
	lastProgress := progressCalls[len(progressCalls)-1]
	mu.Unlock()
	assert.Equal(t, 5, lastProgress)
}

func TestParallelExecutor_ExecuteStream(t *testing.T) {
	executor := NewParallelExecutor[int](ExecutorConfig{MaxWorkers: 2})

	tasks := make([]Task[int], 3)
	for i := 0; i < 3; i++ {
		idx := i
		tasks[i] = Task[int]{
			ID:      string(rune('0' + idx)),
			Execute: func(ctx context.Context) (int, error) { return idx, nil },
		}
	}

	resultChan := executor.ExecuteStream(context.Background(), tasks)

	var results []TaskResult[int]
	for r := range resultChan {
		results = append(results, r)
	}

	assert.Len(t, results, 3)
}

func TestParallelExecutor_Stats(t *testing.T) {
	executor := NewParallelExecutor[int](ExecutorConfig{MaxWorkers: 2})

	tasks := []Task[int]{
		{ID: "1", Execute: func(ctx context.Context) (int, error) { return 1, nil }},
		{ID: "2", Execute: func(ctx context.Context) (int, error) { return 0, errors.New("error") }},
		{ID: "3", Execute: func(ctx context.Context) (int, error) { return 3, nil }},
	}

	executor.Execute(context.Background(), tasks)

	stats := executor.Stats()

	assert.Equal(t, int64(3), stats.TotalTasks)
	assert.Equal(t, int64(2), stats.CompletedTasks)
	assert.Equal(t, int64(1), stats.FailedTasks)
	assert.Greater(t, stats.SuccessRate, 0.6)
}

func TestParallelExecutor_ResetStats(t *testing.T) {
	executor := NewParallelExecutor[int](DefaultExecutorConfig())

	tasks := []Task[int]{
		{ID: "1", Execute: func(ctx context.Context) (int, error) { return 1, nil }},
	}
	executor.Execute(context.Background(), tasks)

	executor.ResetStats()
	stats := executor.Stats()

	assert.Equal(t, int64(0), stats.TotalTasks)
	assert.Equal(t, int64(0), stats.CompletedTasks)
}

func TestBatchExecutor_ExecuteBatches(t *testing.T) {
	tests := []struct {
		name       string
		taskCount  int
		batchSize  int
		wantBatches int
	}{
		{
			name:        "exact batches",
			taskCount:   10,
			batchSize:   5,
			wantBatches: 2,
		},
		{
			name:        "partial last batch",
			taskCount:   7,
			batchSize:   3,
			wantBatches: 3,
		},
		{
			name:        "single batch",
			taskCount:   3,
			batchSize:   10,
			wantBatches: 1,
		},
		{
			name:        "empty tasks",
			taskCount:   0,
			batchSize:   5,
			wantBatches: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewBatchExecutor[int](DefaultExecutorConfig(), tt.batchSize)

			tasks := make([]Task[int], tt.taskCount)
			for i := 0; i < tt.taskCount; i++ {
				idx := i
				tasks[i] = Task[int]{
					ID:      string(rune('0' + idx)),
					Execute: func(ctx context.Context) (int, error) { return idx, nil },
				}
			}

			var batchCount int
			results := executor.ExecuteBatches(context.Background(), tasks, func(batchNum, totalBatches int, results []TaskResult[int]) {
				batchCount++
			})

			if tt.taskCount > 0 {
				assert.Len(t, results, tt.taskCount)
			}
			assert.Equal(t, tt.wantBatches, batchCount)
		})
	}
}

func TestFileChangeTracker_IsChanged(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *FileChangeTracker)
		path    string
		modTime time.Time
		hash    string
		want    bool
	}{
		{
			name:    "new file",
			setup:   func(t *FileChangeTracker) {},
			path:    "new.go",
			modTime: time.Now(),
			hash:    "abc123",
			want:    true,
		},
		{
			name: "unchanged by hash",
			setup: func(t *FileChangeTracker) {
				t.Track("file.go", time.Now(), "abc123")
			},
			path:    "file.go",
			modTime: time.Now().Add(time.Hour), // Different time
			hash:    "abc123",                  // Same hash
			want:    false,
		},
		{
			name: "changed by hash",
			setup: func(t *FileChangeTracker) {
				t.Track("file.go", time.Now(), "abc123")
			},
			path:    "file.go",
			modTime: time.Now(),
			hash:    "def456", // Different hash
			want:    true,
		},
		{
			name: "unchanged by timestamp",
			setup: func(t *FileChangeTracker) {
				t.Track("file.go", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "")
			},
			path:    "file.go",
			modTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			hash:    "",
			want:    false,
		},
		{
			name: "changed by timestamp",
			setup: func(t *FileChangeTracker) {
				t.Track("file.go", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "")
			},
			path:    "file.go",
			modTime: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), // Later
			hash:    "",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewFileChangeTracker()
			tt.setup(tracker)

			got := tracker.IsChanged(tt.path, tt.modTime, tt.hash)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFileChangeTracker_GetChangedFiles(t *testing.T) {
	tracker := NewFileChangeTracker()

	// Track some files
	tracker.Track("a.go", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "hash1")
	tracker.Track("b.go", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "hash2")

	files := []FileInfo{
		{Path: "a.go", ModTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Hash: "hash1"}, // Unchanged
		{Path: "b.go", ModTime: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Hash: "hash3"}, // Changed
		{Path: "c.go", ModTime: time.Now(), Hash: "hash4"},                                  // New
	}

	changed := tracker.GetChangedFiles(files)

	assert.Len(t, changed, 2)
	assert.Contains(t, changed, "b.go")
	assert.Contains(t, changed, "c.go")
}

func TestFileChangeTracker_Remove(t *testing.T) {
	tracker := NewFileChangeTracker()

	tracker.Track("file.go", time.Now(), "hash")
	assert.False(t, tracker.IsChanged("file.go", time.Now(), "hash"))

	tracker.Remove("file.go")
	assert.True(t, tracker.IsChanged("file.go", time.Now(), "hash"))
}

func TestFileChangeTracker_Clear(t *testing.T) {
	tracker := NewFileChangeTracker()

	tracker.Track("a.go", time.Now(), "hash1")
	tracker.Track("b.go", time.Now(), "hash2")

	tracker.Clear()

	stats := tracker.Stats()
	assert.Equal(t, 0, stats["tracked_files"])
}

func TestFileChangeTracker_Stats(t *testing.T) {
	tracker := NewFileChangeTracker()

	tracker.Track("a.go", time.Now(), "hash1")
	tracker.Track("b.go", time.Now(), "")

	stats := tracker.Stats()

	assert.Equal(t, 2, stats["tracked_files"])
	assert.Equal(t, 1, stats["files_with_hash"])
	assert.Equal(t, 1, stats["files_without_hash"])
}

func TestWorkerPool_Submit(t *testing.T) {
	pool := NewWorkerPool(2, 10)
	defer pool.Stop()

	var executed int32
	err := pool.Submit(func() {
		atomic.AddInt32(&executed, 1)
	})

	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
}

func TestWorkerPool_SubmitWait(t *testing.T) {
	pool := NewWorkerPool(2, 10)
	defer pool.Stop()

	var executed bool
	err := pool.SubmitWait(func() {
		time.Sleep(10 * time.Millisecond)
		executed = true
	})

	require.NoError(t, err)
	assert.True(t, executed)
}

func TestWorkerPool_Stop(t *testing.T) {
	pool := NewWorkerPool(2, 10)

	assert.True(t, pool.IsRunning())

	pool.Stop()

	assert.False(t, pool.IsRunning())

	err := pool.Submit(func() {})
	assert.Error(t, err)
}

func TestWorkerPool_QueueFull(t *testing.T) {
	pool := NewWorkerPool(1, 2)
	defer pool.Stop()

	// Block the worker
	blocker := make(chan struct{})
	_ = pool.Submit(func() {
		<-blocker
	})

	// Fill the queue
	_ = pool.Submit(func() {})
	_ = pool.Submit(func() {})

	// Queue should be full
	err := pool.Submit(func() {})
	assert.Error(t, err)

	close(blocker)
}

func TestDefaultExecutorConfig(t *testing.T) {
	cfg := DefaultExecutorConfig()

	assert.Greater(t, cfg.MaxWorkers, 0)
	assert.Equal(t, 30*time.Second, cfg.Timeout)
}

func TestParallelExecutor_ZeroWorkers(t *testing.T) {
	// Should use NumCPU
	executor := NewParallelExecutor[int](ExecutorConfig{MaxWorkers: 0})

	tasks := []Task[int]{
		{ID: "1", Execute: func(ctx context.Context) (int, error) { return 1, nil }},
	}

	results := executor.Execute(context.Background(), tasks)
	assert.Len(t, results, 1)
}

func TestBatchExecutor_ZeroBatchSize(t *testing.T) {
	// Should use default batch size
	executor := NewBatchExecutor[int](DefaultExecutorConfig(), 0)

	tasks := make([]Task[int], 5)
	for i := 0; i < 5; i++ {
		idx := i
		tasks[i] = Task[int]{
			ID:      string(rune('0' + idx)),
			Execute: func(ctx context.Context) (int, error) { return idx, nil },
		}
	}

	results := executor.ExecuteBatches(context.Background(), tasks, nil)
	assert.Len(t, results, 5)
}
