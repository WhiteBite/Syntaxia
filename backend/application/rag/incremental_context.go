// Package rag provides Retrieval Augmented Generation services
package rag

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"syntaxia/domain"
)

const (
	// SessionTimeout is the default session timeout
	SessionTimeout = 30 * time.Minute
	// DefaultMaxSessionTokens is the default max tokens per session
	DefaultMaxSessionTokens = 500000
	// CleanupInterval is how often to clean up expired sessions
	CleanupInterval = 5 * time.Minute
)

// ContextSession represents an incremental context loading session
type ContextSession struct {
	ID          string            `json:"id"`
	ProjectRoot string            `json:"projectRoot"`
	RepoMap     string            `json:"repoMap"`
	LoadedFiles map[string]string `json:"loadedFiles"` // path -> content
	FileTokens  map[string]int    `json:"fileTokens"`  // path -> token count
	TotalTokens int               `json:"totalTokens"`
	CreatedAt   time.Time         `json:"createdAt"`
	LastAccess  time.Time         `json:"lastAccess"`
	MaxTokens   int               `json:"maxTokens"`
}

// ContextUpdate represents the result of requesting files
type ContextUpdate struct {
	NewFiles    map[string]string `json:"newFiles"`
	TokensAdded int               `json:"tokensAdded"`
	TotalTokens int               `json:"totalTokens"`
	Skipped     []string          `json:"skipped,omitempty"` // Files skipped due to budget
}

// ContextDiff represents changes between context states
type ContextDiff struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Changed []string `json:"changed"`
}

// IncrementalContextService provides lazy loading of context
type IncrementalContextService interface {
	// StartSession starts a new incremental context session
	StartSession(ctx context.Context, projectRoot string) (*ContextSession, error)

	// RequestFiles loads specific files into the session context
	RequestFiles(ctx context.Context, sessionID string, files []string) (*ContextUpdate, error)

	// GetDiff returns changes since last request
	GetDiff(ctx context.Context, sessionID string) (*ContextDiff, error)

	// EndSession ends and cleans up a session
	EndSession(ctx context.Context, sessionID string) error

	// GetSession returns session info
	GetSession(ctx context.Context, sessionID string) (*ContextSession, error)

	// GetSessionContext returns the full context for a session
	GetSessionContext(ctx context.Context, sessionID string) (string, error)
}

// IncrementalContextServiceImpl implements IncrementalContextService
type IncrementalContextServiceImpl struct {
	log            domain.Logger
	fileReader     domain.FileContentReader
	repoMapBuilder RepoMapBuilder

	mu       sync.RWMutex
	sessions map[string]*sessionState

	stopCleanup chan struct{}
}

// sessionState holds internal session state
type sessionState struct {
	session       *ContextSession
	previousFiles map[string]string // For diff calculation
}

// NewIncrementalContextService creates a new IncrementalContextService
func NewIncrementalContextService(
	log domain.Logger,
	fileReader domain.FileContentReader,
	repoMapBuilder RepoMapBuilder,
) *IncrementalContextServiceImpl {
	svc := &IncrementalContextServiceImpl{
		log:            log,
		fileReader:     fileReader,
		repoMapBuilder: repoMapBuilder,
		sessions:       make(map[string]*sessionState),
		stopCleanup:    make(chan struct{}),
	}

	// Start cleanup goroutine
	go svc.cleanupLoop()

	return svc
}

// StartSession starts a new incremental context session
func (s *IncrementalContextServiceImpl) StartSession(ctx context.Context, projectRoot string) (*ContextSession, error) {
	s.log.Info(fmt.Sprintf("Starting incremental context session for %s", projectRoot))

	// Generate session ID
	sessionID := s.generateSessionID(projectRoot)

	// Build repo map
	var repoMapStr string
	if s.repoMapBuilder != nil {
		repoMap, err := s.repoMapBuilder.BuildMap(ctx, projectRoot, 8000)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to build repo map: %v", err))
		} else {
			repoMapStr = repoMap.Format()
		}
	}

	now := time.Now()
	session := &ContextSession{
		ID:          sessionID,
		ProjectRoot: projectRoot,
		RepoMap:     repoMapStr,
		LoadedFiles: make(map[string]string),
		FileTokens:  make(map[string]int),
		TotalTokens: estimateTokens(repoMapStr),
		CreatedAt:   now,
		LastAccess:  now,
		MaxTokens:   DefaultMaxSessionTokens,
	}

	// Store session
	s.mu.Lock()
	s.sessions[sessionID] = &sessionState{
		session:       session,
		previousFiles: make(map[string]string),
	}
	s.mu.Unlock()

	s.log.Info(fmt.Sprintf("Created session %s with repo map (%d tokens)", sessionID, session.TotalTokens))

	return session, nil
}

// RequestFiles loads specific files into the session context
func (s *IncrementalContextServiceImpl) RequestFiles(ctx context.Context, sessionID string, files []string) (*ContextUpdate, error) {
	s.mu.Lock()
	state, ok := s.sessions[sessionID]
	if !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Update last access
	state.session.LastAccess = time.Now()

	// Store previous state for diff
	state.previousFiles = copyMap(state.session.LoadedFiles)
	s.mu.Unlock()

	update := &ContextUpdate{
		NewFiles: make(map[string]string),
		Skipped:  make([]string, 0),
	}

	// Filter out already loaded files
	filesToLoad := make([]string, 0)
	for _, f := range files {
		f = filepath.ToSlash(f)
		if _, loaded := state.session.LoadedFiles[f]; !loaded {
			filesToLoad = append(filesToLoad, f)
		}
	}

	if len(filesToLoad) == 0 {
		update.TotalTokens = state.session.TotalTokens
		return update, nil
	}

	s.log.Info(fmt.Sprintf("Loading %d files for session %s", len(filesToLoad), sessionID))

	// Read files
	contents, err := s.fileReader.ReadContents(ctx, filesToLoad, state.session.ProjectRoot, nil)
	if err != nil {
		// Try reading files individually
		contents = make(map[string]string)
		for _, f := range filesToLoad {
			fullPath := filepath.Join(state.session.ProjectRoot, f)
			data, readErr := os.ReadFile(fullPath)
			if readErr == nil {
				contents[f] = string(data)
			}
		}
	}

	// Add files within token budget
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, f := range filesToLoad {
		content, ok := contents[f]
		if !ok {
			continue
		}

		tokens := estimateTokens(content)

		// Check budget
		if state.session.TotalTokens+tokens > state.session.MaxTokens {
			// Try to truncate
			availableTokens := state.session.MaxTokens - state.session.TotalTokens
			if availableTokens < 200 {
				update.Skipped = append(update.Skipped, f)
				continue
			}
			content = truncateToTokens(content, availableTokens)
			tokens = availableTokens
		}

		state.session.LoadedFiles[f] = content
		state.session.FileTokens[f] = tokens
		state.session.TotalTokens += tokens

		update.NewFiles[f] = content
		update.TokensAdded += tokens
	}

	update.TotalTokens = state.session.TotalTokens

	s.log.Info(fmt.Sprintf("Loaded %d files, added %d tokens, total %d tokens",
		len(update.NewFiles), update.TokensAdded, update.TotalTokens))

	return update, nil
}

// GetDiff returns changes since last request
func (s *IncrementalContextServiceImpl) GetDiff(ctx context.Context, sessionID string) (*ContextDiff, error) {
	s.mu.RLock()
	state, ok := s.sessions[sessionID]
	if !ok {
		s.mu.RUnlock()
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	current := state.session.LoadedFiles
	previous := state.previousFiles
	s.mu.RUnlock()

	diff := &ContextDiff{
		Added:   make([]string, 0),
		Removed: make([]string, 0),
		Changed: make([]string, 0),
	}

	// Find added and changed files
	for path, content := range current {
		prevContent, existed := previous[path]
		if !existed {
			diff.Added = append(diff.Added, path)
		} else if prevContent != content {
			diff.Changed = append(diff.Changed, path)
		}
	}

	// Find removed files
	for path := range previous {
		if _, exists := current[path]; !exists {
			diff.Removed = append(diff.Removed, path)
		}
	}

	// Sort for consistent output
	sort.Strings(diff.Added)
	sort.Strings(diff.Removed)
	sort.Strings(diff.Changed)

	return diff, nil
}

// EndSession ends and cleans up a session
func (s *IncrementalContextServiceImpl) EndSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[sessionID]; !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	delete(s.sessions, sessionID)
	s.log.Info(fmt.Sprintf("Ended session %s", sessionID))

	return nil
}

// GetSession returns session info
func (s *IncrementalContextServiceImpl) GetSession(ctx context.Context, sessionID string) (*ContextSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Return a copy
	session := *state.session
	session.LoadedFiles = copyMap(state.session.LoadedFiles)
	session.FileTokens = copyIntMap(state.session.FileTokens)

	return &session, nil
}

// GetSessionContext returns the full context for a session
func (s *IncrementalContextServiceImpl) GetSessionContext(ctx context.Context, sessionID string) (string, error) {
	s.mu.RLock()
	state, ok := s.sessions[sessionID]
	if !ok {
		s.mu.RUnlock()
		return "", fmt.Errorf("session not found: %s", sessionID)
	}

	session := state.session
	s.mu.RUnlock()

	return FormatSessionContext(session), nil
}

// UnloadFiles removes files from the session context
func (s *IncrementalContextServiceImpl) UnloadFiles(ctx context.Context, sessionID string, files []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Store previous state for diff
	state.previousFiles = copyMap(state.session.LoadedFiles)

	for _, f := range files {
		f = filepath.ToSlash(f)
		if tokens, exists := state.session.FileTokens[f]; exists {
			state.session.TotalTokens -= tokens
			delete(state.session.LoadedFiles, f)
			delete(state.session.FileTokens, f)
		}
	}

	s.log.Info(fmt.Sprintf("Unloaded %d files from session %s", len(files), sessionID))

	return nil
}

// SetMaxTokens updates the max tokens for a session
func (s *IncrementalContextServiceImpl) SetMaxTokens(sessionID string, maxTokens int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	state.session.MaxTokens = maxTokens
	return nil
}

// GetActiveSessions returns the number of active sessions
func (s *IncrementalContextServiceImpl) GetActiveSessions() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

// Stop stops the cleanup goroutine
func (s *IncrementalContextServiceImpl) Stop() {
	close(s.stopCleanup)
}

// cleanupLoop periodically cleans up expired sessions
func (s *IncrementalContextServiceImpl) cleanupLoop() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpiredSessions()
		case <-s.stopCleanup:
			return
		}
	}
}

// cleanupExpiredSessions removes expired sessions
func (s *IncrementalContextServiceImpl) cleanupExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	expired := make([]string, 0)

	for id, state := range s.sessions {
		if now.Sub(state.session.LastAccess) > SessionTimeout {
			expired = append(expired, id)
		}
	}

	for _, id := range expired {
		delete(s.sessions, id)
		s.log.Info(fmt.Sprintf("Cleaned up expired session %s", id))
	}

	if len(expired) > 0 {
		s.log.Info(fmt.Sprintf("Cleaned up %d expired sessions", len(expired)))
	}
}

// generateSessionID generates a unique session ID
func (s *IncrementalContextServiceImpl) generateSessionID(projectRoot string) string {
	data := fmt.Sprintf("%s:%d", projectRoot, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:8])
}

// FormatSessionContext formats the session context as a string for AI
func FormatSessionContext(session *ContextSession) string {
	var sb strings.Builder

	// Add repo map first
	if session.RepoMap != "" {
		sb.WriteString(session.RepoMap)
		sb.WriteString("\n")
	}

	// Add loaded files
	if len(session.LoadedFiles) > 0 {
		sb.WriteString("# Loaded Files\n\n")

		// Sort files for consistent output
		paths := make([]string, 0, len(session.LoadedFiles))
		for path := range session.LoadedFiles {
			paths = append(paths, path)
		}
		sort.Strings(paths)

		for _, path := range paths {
			content := session.LoadedFiles[path]
			ext := filepath.Ext(path)
			if ext != "" {
				ext = ext[1:]
			}

			sb.WriteString(fmt.Sprintf("## %s\n", path))
			sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", ext, content))
		}
	}

	return sb.String()
}

// FormatContextUpdate formats a context update for AI
func FormatContextUpdate(update *ContextUpdate) string {
	var sb strings.Builder

	if len(update.NewFiles) == 0 {
		return "No new files loaded."
	}

	sb.WriteString("# Newly Loaded Files\n\n")

	// Sort files for consistent output
	paths := make([]string, 0, len(update.NewFiles))
	for path := range update.NewFiles {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		content := update.NewFiles[path]
		ext := filepath.Ext(path)
		if ext != "" {
			ext = ext[1:]
		}

		sb.WriteString(fmt.Sprintf("## %s\n", path))
		sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", ext, content))
	}

	if len(update.Skipped) > 0 {
		sb.WriteString("## Skipped (token budget exceeded)\n")
		for _, path := range update.Skipped {
			sb.WriteString(fmt.Sprintf("- %s\n", path))
		}
	}

	return sb.String()
}

// Helper functions

func copyMap(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

func copyIntMap(m map[string]int) map[string]int {
	result := make(map[string]int, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
