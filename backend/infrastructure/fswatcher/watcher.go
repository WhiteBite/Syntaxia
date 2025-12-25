package fswatcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// debounceDelay is the time to wait before emitting file change events
	// This prevents multiple rapid events from triggering multiple reindexes
	debounceDelay = 500 * time.Millisecond
)

type Watcher struct {
	log           domain.Logger
	bus           domain.EventBus
	fsWatcher     *fsnotify.Watcher
	mu            sync.Mutex
	cancel        context.CancelFunc
	rootDir       string
	appCtx        context.Context
	debounceTimer *time.Timer
	pendingFiles  map[string]struct{}
	debounceMu    sync.Mutex
}

func New(ctx context.Context, bus domain.EventBus) (*Watcher, error) {
	return &Watcher{
		appCtx:       ctx,
		log:          wailsLogger{ctx: ctx},
		bus:          bus,
		pendingFiles: make(map[string]struct{}),
	}, nil
}

func (w *Watcher) shouldSkipDir(name string) bool {
	// Общий набор шумных директорий
	switch name {
	case ".git", "node_modules", ".idea", "dist", "build", ".cache", ".vite", ".wails", "out", "target", "bin", "obj", "coverage":
		return true
	default:
		return false
	}
}

func (w *Watcher) Start(path string) error {
	w.Stop()

	w.mu.Lock()
	defer w.mu.Unlock()
	w.rootDir = path

	var err error
	w.fsWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	// Рекурсивно подписываемся на директории, пропуская шумные
	err = filepath.WalkDir(w.rootDir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if w.shouldSkipDir(d.Name()) {
				w.log.Info("Watcher: skip dir " + p)
				return filepath.SkipDir
			}
			return w.fsWatcher.Add(p)
		}
		return nil
	})
	if err != nil {
		return err
	}

	go w.run(ctx)
	w.log.Info("Наблюдатель запущен для: " + path)
	return nil
}

func (w *Watcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	if w.fsWatcher != nil {
		w.fsWatcher.Close()
		w.fsWatcher = nil
		w.log.Info("Наблюдатель остановлен.")
	}
}

func (w *Watcher) run(ctx context.Context) {
	defer func() {
		w.mu.Lock()
		if w.fsWatcher != nil {
			w.fsWatcher.Close()
			w.fsWatcher = nil
		}
		w.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			// Фильтруем шумные события внутри .git (и любые подпути .git)
			// чтобы не триггерить лишние пересканы/сбросы.
			name := event.Name
			base := filepath.Base(name)
			if base == ".git" || strings.Contains(name, string(os.PathSeparator)+".git"+string(os.PathSeparator)) {
				// Пропускаем события из .git
				continue
			}

			// Debounce: collect events and emit after delay
			w.scheduleDebounce(name)
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			w.log.Error(fmt.Sprintf("Ошибка наблюдателя: %v", err))
		}
	}
}

// scheduleDebounce collects file change events and emits them after a delay
func (w *Watcher) scheduleDebounce(filePath string) {
	w.debounceMu.Lock()
	defer w.debounceMu.Unlock()

	// Add file to pending set
	w.pendingFiles[filePath] = struct{}{}

	// Reset or create timer
	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(debounceDelay, func() {
		w.flushPendingEvents()
	})
}

// flushPendingEvents emits all pending file change events
func (w *Watcher) flushPendingEvents() {
	w.debounceMu.Lock()
	files := make([]string, 0, len(w.pendingFiles))
	for f := range w.pendingFiles {
		files = append(files, f)
	}
	w.pendingFiles = make(map[string]struct{})
	w.debounceMu.Unlock()

	if len(files) > 0 {
		// Emit single event with all changed files
		w.bus.Emit("projectFilesChanged", w.rootDir)
		// Also emit individual file events for fine-grained updates
		for _, f := range files {
			w.bus.Emit("fileChanged", f)
		}
	}
}

func (w *Watcher) RefreshAndRescan() error {
	w.mu.Lock()
	path := w.rootDir
	w.mu.Unlock()
	if path == "" {
		return nil
	}
	w.log.Info("Обновление наблюдателя из-за смены правил...")
	return w.Start(path)
}

type wailsLogger struct {
	ctx context.Context
}

func (l wailsLogger) Debug(message string)   { runtime.LogDebug(l.ctx, message) }
func (l wailsLogger) Info(message string)    { runtime.LogInfo(l.ctx, message) }
func (l wailsLogger) Warning(message string) { runtime.LogWarning(l.ctx, message) }
func (l wailsLogger) Error(message string)   { runtime.LogError(l.ctx, message) }
func (l wailsLogger) Fatal(message string)   { runtime.LogFatal(l.ctx, message) }
