package analysis

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a unit of work to be executed in parallel.
type Task[T any] struct {
	ID      string
	Execute func(ctx context.Context) (T, error)
}

// TaskResult holds the result of a task execution.
type TaskResult[T any] struct {
	ID       string
	Value    T
	Error    error
	Duration time.Duration
}

// ParallelExecutor executes tasks concurrently with configurable parallelism.
type ParallelExecutor[T any] struct {
	maxWorkers int
	timeout    time.Duration

	// Metrics
	totalTasks     int64
	completedTasks int64
	failedTasks    int64
	totalDuration  int64 // nanoseconds
}

// ExecutorConfig holds configuration for ParallelExecutor.
type ExecutorConfig struct {
	MaxWorkers int           // Maximum concurrent workers (0 = NumCPU)
	Timeout    time.Duration // Timeout per task (0 = no timeout)
}

// DefaultExecutorConfig returns default executor configuration.
func DefaultExecutorConfig() ExecutorConfig {
	return ExecutorConfig{
		MaxWorkers: runtime.NumCPU(),
		Timeout:    30 * time.Second,
	}
}

// NewParallelExecutor creates a new parallel executor.
func NewParallelExecutor[T any](cfg ExecutorConfig) *ParallelExecutor[T] {
	if cfg.MaxWorkers <= 0 {
		cfg.MaxWorkers = runtime.NumCPU()
	}
	return &ParallelExecutor[T]{
		maxWorkers: cfg.MaxWorkers,
		timeout:    cfg.Timeout,
	}
}

// Execute runs all tasks in parallel and returns results.
// Results are returned in the same order as input tasks.
func (e *ParallelExecutor[T]) Execute(ctx context.Context, tasks []Task[T]) []TaskResult[T] {
	if len(tasks) == 0 {
		return nil
	}

	results := make([]TaskResult[T], len(tasks))
	taskChan := make(chan int, len(tasks))
	var wg sync.WaitGroup

	// Determine number of workers
	numWorkers := e.maxWorkers
	if numWorkers > len(tasks) {
		numWorkers = len(tasks)
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range taskChan {
				results[idx] = e.executeTask(ctx, tasks[idx])
			}
		}()
	}

	// Send tasks to workers
	for i := range tasks {
		taskChan <- i
	}
	close(taskChan)

	// Wait for all workers to complete
	wg.Wait()

	return results
}

// ExecuteWithProgress runs tasks and reports progress via callback.
func (e *ParallelExecutor[T]) ExecuteWithProgress(
	ctx context.Context,
	tasks []Task[T],
	progress func(completed, total int),
) []TaskResult[T] {
	if len(tasks) == 0 {
		return nil
	}

	results := make([]TaskResult[T], len(tasks))
	taskChan := make(chan int, len(tasks))
	var wg sync.WaitGroup
	var completed int64

	numWorkers := e.maxWorkers
	if numWorkers > len(tasks) {
		numWorkers = len(tasks)
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range taskChan {
				results[idx] = e.executeTask(ctx, tasks[idx])
				current := atomic.AddInt64(&completed, 1)
				if progress != nil {
					progress(int(current), len(tasks))
				}
			}
		}()
	}

	// Send tasks
	for i := range tasks {
		taskChan <- i
	}
	close(taskChan)

	wg.Wait()
	return results
}

// ExecuteStream runs tasks and streams results as they complete.
func (e *ParallelExecutor[T]) ExecuteStream(
	ctx context.Context,
	tasks []Task[T],
) <-chan TaskResult[T] {
	resultChan := make(chan TaskResult[T], len(tasks))

	go func() {
		defer close(resultChan)

		if len(tasks) == 0 {
			return
		}

		taskChan := make(chan Task[T], len(tasks))
		var wg sync.WaitGroup

		numWorkers := e.maxWorkers
		if numWorkers > len(tasks) {
			numWorkers = len(tasks)
		}

		// Start workers
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range taskChan {
					select {
					case <-ctx.Done():
						return
					case resultChan <- e.executeTask(ctx, task):
					}
				}
			}()
		}

		// Send tasks
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				close(taskChan)
				wg.Wait()
				return
			case taskChan <- task:
			}
		}
		close(taskChan)

		wg.Wait()
	}()

	return resultChan
}

func (e *ParallelExecutor[T]) executeTask(ctx context.Context, task Task[T]) TaskResult[T] {
	atomic.AddInt64(&e.totalTasks, 1)
	start := time.Now()

	var result TaskResult[T]
	result.ID = task.ID

	// Apply timeout if configured
	execCtx := ctx
	var cancel context.CancelFunc
	if e.timeout > 0 {
		execCtx, cancel = context.WithTimeout(ctx, e.timeout)
		defer cancel()
	}

	// Execute task
	value, err := task.Execute(execCtx)
	result.Duration = time.Since(start)
	result.Value = value
	result.Error = err

	atomic.AddInt64(&e.totalDuration, int64(result.Duration))

	if err != nil {
		atomic.AddInt64(&e.failedTasks, 1)
	} else {
		atomic.AddInt64(&e.completedTasks, 1)
	}

	return result
}

// Stats returns executor statistics.
func (e *ParallelExecutor[T]) Stats() ExecutorStats {
	total := atomic.LoadInt64(&e.totalTasks)
	completed := atomic.LoadInt64(&e.completedTasks)
	failed := atomic.LoadInt64(&e.failedTasks)
	duration := atomic.LoadInt64(&e.totalDuration)

	avgDuration := time.Duration(0)
	if total > 0 {
		avgDuration = time.Duration(duration / total)
	}

	successRate := float64(0)
	if total > 0 {
		successRate = float64(completed) / float64(total)
	}

	return ExecutorStats{
		TotalTasks:     total,
		CompletedTasks: completed,
		FailedTasks:    failed,
		AvgDuration:    avgDuration,
		SuccessRate:    successRate,
		MaxWorkers:     e.maxWorkers,
	}
}

// ResetStats resets executor statistics.
func (e *ParallelExecutor[T]) ResetStats() {
	atomic.StoreInt64(&e.totalTasks, 0)
	atomic.StoreInt64(&e.completedTasks, 0)
	atomic.StoreInt64(&e.failedTasks, 0)
	atomic.StoreInt64(&e.totalDuration, 0)
}

// ExecutorStats holds executor statistics.
type ExecutorStats struct {
	TotalTasks     int64         `json:"totalTasks"`
	CompletedTasks int64         `json:"completedTasks"`
	FailedTasks    int64         `json:"failedTasks"`
	AvgDuration    time.Duration `json:"avgDuration"`
	SuccessRate    float64       `json:"successRate"`
	MaxWorkers     int           `json:"maxWorkers"`
}

// BatchExecutor provides batch processing with automatic chunking.
type BatchExecutor[T any] struct {
	executor  *ParallelExecutor[T]
	batchSize int
}

// NewBatchExecutor creates a new batch executor.
func NewBatchExecutor[T any](cfg ExecutorConfig, batchSize int) *BatchExecutor[T] {
	if batchSize <= 0 {
		batchSize = 100
	}
	return &BatchExecutor[T]{
		executor:  NewParallelExecutor[T](cfg),
		batchSize: batchSize,
	}
}

// ExecuteBatches processes tasks in batches.
func (b *BatchExecutor[T]) ExecuteBatches(
	ctx context.Context,
	tasks []Task[T],
	batchComplete func(batchNum, totalBatches int, results []TaskResult[T]),
) []TaskResult[T] {
	if len(tasks) == 0 {
		return nil
	}

	allResults := make([]TaskResult[T], 0, len(tasks))
	totalBatches := (len(tasks) + b.batchSize - 1) / b.batchSize

	for i := 0; i < len(tasks); i += b.batchSize {
		select {
		case <-ctx.Done():
			return allResults
		default:
		}

		end := i + b.batchSize
		if end > len(tasks) {
			end = len(tasks)
		}

		batch := tasks[i:end]
		results := b.executor.Execute(ctx, batch)
		allResults = append(allResults, results...)

		if batchComplete != nil {
			batchNum := (i / b.batchSize) + 1
			batchComplete(batchNum, totalBatches, results)
		}
	}

	return allResults
}

// FileChangeTracker tracks file changes for incremental builds.
type FileChangeTracker struct {
	mu         sync.RWMutex
	timestamps map[string]time.Time
	hashes     map[string]string
}

// NewFileChangeTracker creates a new file change tracker.
func NewFileChangeTracker() *FileChangeTracker {
	return &FileChangeTracker{
		timestamps: make(map[string]time.Time),
		hashes:     make(map[string]string),
	}
}

// IsChanged checks if a file has changed since last tracked.
func (t *FileChangeTracker) IsChanged(path string, modTime time.Time, hash string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	lastTime, hasTime := t.timestamps[path]
	lastHash, hasHash := t.hashes[path]

	// New file
	if !hasTime && !hasHash {
		return true
	}

	// Check by hash if available
	if hash != "" && hasHash {
		return hash != lastHash
	}

	// Fall back to timestamp
	if hasTime {
		return modTime.After(lastTime)
	}

	return true
}

// Track records the current state of a file.
func (t *FileChangeTracker) Track(path string, modTime time.Time, hash string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.timestamps[path] = modTime
	if hash != "" {
		t.hashes[path] = hash
	}
}

// Remove removes a file from tracking.
func (t *FileChangeTracker) Remove(path string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.timestamps, path)
	delete(t.hashes, path)
}

// Clear removes all tracked files.
func (t *FileChangeTracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.timestamps = make(map[string]time.Time)
	t.hashes = make(map[string]string)
}

// GetChangedFiles returns files that have changed from the given list.
func (t *FileChangeTracker) GetChangedFiles(files []FileInfo) []string {
	changed := make([]string, 0)

	for _, f := range files {
		if t.IsChanged(f.Path, f.ModTime, f.Hash) {
			changed = append(changed, f.Path)
		}
	}

	return changed
}

// TrackFiles updates tracking for multiple files.
func (t *FileChangeTracker) TrackFiles(files []FileInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, f := range files {
		t.timestamps[f.Path] = f.ModTime
		if f.Hash != "" {
			t.hashes[f.Path] = f.Hash
		}
	}
}

// FileInfo holds file information for change tracking.
type FileInfo struct {
	Path    string
	ModTime time.Time
	Hash    string
}

// TrackerStats returns tracker statistics.
func (t *FileChangeTracker) Stats() map[string]int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return map[string]int{
		"tracked_files":       len(t.timestamps),
		"files_with_hash":     len(t.hashes),
		"files_without_hash":  len(t.timestamps) - len(t.hashes),
	}
}

// WorkerPool provides a reusable pool of workers for task execution.
type WorkerPool struct {
	workers    int
	taskQueue  chan func()
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	running    int32
	queuedJobs int64
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	if queueSize <= 0 {
		queueSize = workers * 10
	}

	ctx, cancel := context.WithCancel(context.Background())
	pool := &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func(), queueSize),
		ctx:       ctx,
		cancel:    cancel,
	}

	pool.start()
	return pool
}

func (p *WorkerPool) start() {
	atomic.StoreInt32(&p.running, 1)
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}
			task()
			atomic.AddInt64(&p.queuedJobs, -1)
		}
	}
}

// Submit adds a task to the worker pool.
// Returns error if pool is stopped or queue is full.
func (p *WorkerPool) Submit(task func()) error {
	if atomic.LoadInt32(&p.running) == 0 {
		return fmt.Errorf("worker pool is stopped")
	}

	select {
	case p.taskQueue <- task:
		atomic.AddInt64(&p.queuedJobs, 1)
		return nil
	default:
		return fmt.Errorf("task queue is full")
	}
}

// SubmitWait adds a task and waits for it to complete.
func (p *WorkerPool) SubmitWait(task func()) error {
	done := make(chan struct{})
	err := p.Submit(func() {
		defer close(done)
		task()
	})
	if err != nil {
		return err
	}
	<-done
	return nil
}

// Stop gracefully stops the worker pool.
func (p *WorkerPool) Stop() {
	if atomic.CompareAndSwapInt32(&p.running, 1, 0) {
		p.cancel()
		close(p.taskQueue)
		p.wg.Wait()
	}
}

// QueuedJobs returns the number of jobs currently in queue.
func (p *WorkerPool) QueuedJobs() int64 {
	return atomic.LoadInt64(&p.queuedJobs)
}

// IsRunning returns true if the pool is running.
func (p *WorkerPool) IsRunning() bool {
	return atomic.LoadInt32(&p.running) == 1
}
