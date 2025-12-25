package memory

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite" // SQLite driver
)

// ContextMemoryImpl implements domain.ContextMemory interface
type ContextMemoryImpl struct {
	db *sql.DB
	mu sync.RWMutex
}

// Ensure ContextMemoryImpl implements domain.ContextMemory
var _ domain.ContextMemory = (*ContextMemoryImpl)(nil)

// scanContextRows scans rows into ConversationContext slice
func scanContextRows(rows *sql.Rows) ([]*domain.ConversationContext, error) {
	var contexts []*domain.ConversationContext
	for rows.Next() {
		var ctx domain.ConversationContext
		var filesJSON, symbolsJSON string
		var lastAccessed, createdAt int64

		if err := rows.Scan(&ctx.ID, &ctx.ProjectRoot, &ctx.Topic, &filesJSON, &symbolsJSON,
			&ctx.Summary, &lastAccessed, &createdAt, &ctx.MessageCount); err != nil {
			continue
		}

		_ = json.Unmarshal([]byte(filesJSON), &ctx.Files)
		_ = json.Unmarshal([]byte(symbolsJSON), &ctx.Symbols)
		ctx.LastAccessed = time.Unix(lastAccessed, 0)
		ctx.CreatedAt = time.Unix(createdAt, 0)

		contexts = append(contexts, &ctx)
	}
	return contexts, nil
}

// NewContextMemory creates a new context memory store
func NewContextMemory(cacheDir string) (*ContextMemoryImpl, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(cacheDir, "context_memory.db")
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	cm := &ContextMemoryImpl{db: db}
	if err := cm.initDB(); err != nil {
		db.Close()
		return nil, err
	}

	return cm, nil
}

func (cm *ContextMemoryImpl) initDB() error {
	schema := `
	CREATE TABLE IF NOT EXISTS contexts (
		id TEXT PRIMARY KEY,
		project_root TEXT NOT NULL,
		topic TEXT,
		files TEXT,
		symbols TEXT,
		summary TEXT,
		last_accessed INTEGER,
		created_at INTEGER,
		message_count INTEGER DEFAULT 0
	);
	
	CREATE TABLE IF NOT EXISTS preferences (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at INTEGER
	);
	
	CREATE TABLE IF NOT EXISTS task_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		context_id TEXT,
		task TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		result TEXT,
		created_at INTEGER,
		completed_at INTEGER,
		FOREIGN KEY (context_id) REFERENCES contexts(id)
	);
	
	CREATE INDEX IF NOT EXISTS idx_contexts_project ON contexts(project_root);
	CREATE INDEX IF NOT EXISTS idx_contexts_topic ON contexts(topic);
	CREATE INDEX IF NOT EXISTS idx_task_context ON task_history(context_id);
	`
	_, err := cm.db.Exec(schema)
	return err
}

// SaveContext saves or updates a conversation context
func (cm *ContextMemoryImpl) SaveContext(ctx *domain.ConversationContext) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	filesJSON, _ := json.Marshal(ctx.Files)
	symbolsJSON, _ := json.Marshal(ctx.Symbols)

	_, err := cm.db.Exec(`
		INSERT OR REPLACE INTO contexts 
		(id, project_root, topic, files, symbols, summary, last_accessed, created_at, message_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, ctx.ID, ctx.ProjectRoot, ctx.Topic, string(filesJSON), string(symbolsJSON),
		ctx.Summary, ctx.LastAccessed.Unix(), ctx.CreatedAt.Unix(), ctx.MessageCount)

	return err
}

// GetContext retrieves a context by ID
func (cm *ContextMemoryImpl) GetContext(id string) (*domain.ConversationContext, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var ctx domain.ConversationContext
	var filesJSON, symbolsJSON string
	var lastAccessed, createdAt int64

	err := cm.db.QueryRow(`
		SELECT id, project_root, topic, files, symbols, summary, last_accessed, created_at, message_count
		FROM contexts WHERE id = ?
	`, id).Scan(&ctx.ID, &ctx.ProjectRoot, &ctx.Topic, &filesJSON, &symbolsJSON,
		&ctx.Summary, &lastAccessed, &createdAt, &ctx.MessageCount)

	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(filesJSON), &ctx.Files)
	_ = json.Unmarshal([]byte(symbolsJSON), &ctx.Symbols)
	ctx.LastAccessed = time.Unix(lastAccessed, 0)
	ctx.CreatedAt = time.Unix(createdAt, 0)

	// Update last accessed
	_, _ = cm.db.Exec("UPDATE contexts SET last_accessed = ? WHERE id = ?", time.Now().Unix(), id)

	return &ctx, nil
}

// FindContextByTopic finds contexts matching a topic
func (cm *ContextMemoryImpl) FindContextByTopic(projectRoot, topic string) ([]*domain.ConversationContext, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	rows, err := cm.db.Query(`
		SELECT id, project_root, topic, files, symbols, summary, last_accessed, created_at, message_count
		FROM contexts 
		WHERE project_root = ? AND (topic LIKE ? OR summary LIKE ?)
		ORDER BY last_accessed DESC
		LIMIT 10
	`, projectRoot, "%"+topic+"%", "%"+topic+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanContextRows(rows)
}

// GetRecentContexts returns recent contexts for a project
func (cm *ContextMemoryImpl) GetRecentContexts(projectRoot string, limit int) ([]*domain.ConversationContext, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if limit <= 0 {
		limit = 10
	}

	rows, err := cm.db.Query(`
		SELECT id, project_root, topic, files, symbols, summary, last_accessed, created_at, message_count
		FROM contexts 
		WHERE project_root = ?
		ORDER BY last_accessed DESC
		LIMIT ?
	`, projectRoot, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanContextRows(rows)
}

// SetPreference saves a user preference
func (cm *ContextMemoryImpl) SetPreference(key, value string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	_, err := cm.db.Exec(`
		INSERT OR REPLACE INTO preferences (key, value, updated_at)
		VALUES (?, ?, ?)
	`, key, value, time.Now().Unix())

	return err
}

// GetPreference retrieves a user preference
func (cm *ContextMemoryImpl) GetPreference(key string) (string, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var value string
	err := cm.db.QueryRow("SELECT value FROM preferences WHERE key = ?", key).Scan(&value)
	return value, err
}

// GetAllPreferences returns all user preferences
func (cm *ContextMemoryImpl) GetAllPreferences() (map[string]string, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	rows, err := cm.db.Query("SELECT key, value FROM preferences")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prefs := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}
		prefs[key] = value
	}

	return prefs, nil
}

// ExtractTopicFromMessage extracts topic keywords from a message
func ExtractTopicFromMessage(message string) string {
	// Simple extraction - look for common patterns
	message = strings.ToLower(message)

	// Common topic indicators
	patterns := []string{
		"работа над ", "работу над ", "работай над ",
		"continue with ", "work on ", "working on ",
		"fix ", "implement ", "add ", "update ", "refactor ",
		"исправь ", "добавь ", "обнови ", "рефактор ",
	}

	for _, p := range patterns {
		if idx := strings.Index(message, p); idx != -1 {
			rest := message[idx+len(p):]
			// Take first few words
			words := strings.Fields(rest)
			if len(words) > 3 {
				words = words[:3]
			}
			return strings.Join(words, " ")
		}
	}

	// Fallback: first significant words
	words := strings.Fields(message)
	var significant []string
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "и": true, "в": true, "на": true,
		"is": true, "are": true, "to": true, "for": true, "with": true,
	}

	for _, w := range words {
		if len(w) > 2 && !stopWords[w] {
			significant = append(significant, w)
			if len(significant) >= 3 {
				break
			}
		}
	}

	return strings.Join(significant, " ")
}

// Close closes the database connection
func (cm *ContextMemoryImpl) Close() error {
	if cm.db != nil {
		return cm.db.Close()
	}
	return nil
}
