package wailsbridge

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Bridge implements domain.Logger and domain.EventBus using the Wails runtime.
// It acts as a bridge between the backend application logic and the Wails frontend/runtime.
type Bridge struct {
	ctx context.Context
}

// New creates a new Bridge instance.
func New(ctx context.Context) *Bridge {
	return &Bridge{ctx: ctx}
}

// SetWailsContext allows updating the context, which is necessary because the
// Wails context changes during the application lifecycle (e.g., OnStartup vs. OnDomReady).
func (b *Bridge) SetWailsContext(ctx context.Context) {
	b.ctx = ctx
}

// --- domain.Logger implementation ---

func (b *Bridge) Debug(message string) {
	runtime.LogDebug(b.ctx, message)
}

func (b *Bridge) Info(message string) {
	runtime.LogInfo(b.ctx, message)
}

func (b *Bridge) Warning(message string) {
	runtime.LogWarning(b.ctx, message)
}

func (b *Bridge) Error(message string) {
	runtime.LogError(b.ctx, message)
}

func (b *Bridge) Fatal(message string) {
	runtime.LogFatal(b.ctx, message)
}

// --- domain.EventBus implementation ---

func (b *Bridge) Emit(eventName string, data ...interface{}) {
	// MEMORY OPTIMIZATION: Removed all logging from Emit to prevent log accumulation
	// Logs were causing memory leaks by accumulating in Wails runtime
	if len(data) > 0 {
		runtime.EventsEmit(b.ctx, eventName, data[0])
	} else {
		runtime.EventsEmit(b.ctx, eventName)
	}
}

// --- Wails Dialogs ---

// OpenDirectoryDialog opens a native directory selection dialog.
func (b *Bridge) OpenDirectoryDialog() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(b.ctx, runtime.OpenDialogOptions{
		Title: "Select Project Directory",
	})
	if err != nil {
		return "", fmt.Errorf("failed to open directory dialog: %w", err)
	}
	return dir, nil
}

// --- Window Management ---

// WindowState represents the window position and size
type WindowState struct {
	X          int  `json:"x"`
	Y          int  `json:"y"`
	Width      int  `json:"width"`
	Height     int  `json:"height"`
	Maximized  bool `json:"maximized"`
	Fullscreen bool `json:"fullscreen"`
}

// GetWindowState returns current window position and size
// Returns zero state if window is already closed (prevents panic on shutdown)
func (b *Bridge) GetWindowState() WindowState {
	// Use recover to handle panic from Wails when window is already closed
	defer func() {
		_ = recover()
	}()

	w, h := runtime.WindowGetSize(b.ctx)
	// If size is 0, window is likely closed - return empty state
	if w == 0 || h == 0 {
		return WindowState{}
	}

	x, y := runtime.WindowGetPosition(b.ctx)
	maximized := runtime.WindowIsMaximised(b.ctx)
	fullscreen := runtime.WindowIsFullscreen(b.ctx)
	return WindowState{
		X:          x,
		Y:          y,
		Width:      w,
		Height:     h,
		Maximized:  maximized,
		Fullscreen: fullscreen,
	}
}

// SetWindowState restores window position and size
func (b *Bridge) SetWindowState(state WindowState) {
	if state.Fullscreen {
		runtime.WindowFullscreen(b.ctx)
		return
	}
	if state.Maximized {
		runtime.WindowMaximise(b.ctx)
		return
	}
	if state.Width > 0 && state.Height > 0 {
		runtime.WindowSetSize(b.ctx, state.Width, state.Height)
	}
	if state.X >= 0 && state.Y >= 0 {
		runtime.WindowSetPosition(b.ctx, state.X, state.Y)
	}
}

// WindowMaximize maximizes the window
func (b *Bridge) WindowMaximize() {
	runtime.WindowMaximise(b.ctx)
}

// WindowUnmaximize restores the window from maximized state
func (b *Bridge) WindowUnmaximize() {
	runtime.WindowUnmaximise(b.ctx)
}

// WindowToggleMaximize toggles between maximized and normal state
func (b *Bridge) WindowToggleMaximize() {
	runtime.WindowToggleMaximise(b.ctx)
}

// WindowFullscreen sets the window to fullscreen
func (b *Bridge) WindowFullscreen() {
	runtime.WindowFullscreen(b.ctx)
}

// WindowUnfullscreen exits fullscreen mode
func (b *Bridge) WindowUnfullscreen() {
	runtime.WindowUnfullscreen(b.ctx)
}
