---
inclusion: fileMatch
fileMatchPattern: "backend/domain/**/*.go"
---

# Domain Layer - Правила

## Назначение

Domain — ядро приложения. Содержит:

- Бизнес-модели (entities)
- Интерфейсы сервисов
- Бизнес-ошибки
- Константы конфигурации

## Структура

```
backend/domain/
├── models.go           # Основные модели (FileNode, Project, etc.)
├── models_ai.go        # AI-специфичные модели
├── models_settings.go  # Настройки
├── models_export.go    # DTO для экспорта
├── interfaces.go       # Интерфейсы сервисов
├── tools.go            # Tool/ToolCall/ToolResult
├── errors.go           # Бизнес-ошибки
├── config.go           # Константы
├── task_types.go       # Taskflow типы
└── analysis/           # Интерфейсы анализаторов
```

## Правила

### 1. Никаких внешних зависимостей

```go
// ❌ Запрещено в domain/
import "syntaxia/infrastructure/..."
import "syntaxia/application/..."

// ✅ Разрешено
import "context"
import "time"
import "fmt"
```

### 2. Интерфейсы определяются здесь

```go
// domain/interfaces.go
type ContextBuilder interface {
    BuildContext(ctx context.Context, projectPath string, ...) (*ContextSummary, error)
}

type GitRepository interface {
    GetUncommittedFiles(projectRoot string) ([]FileStatus, error)
    // ...
}
```

### 3. Модели — чистые структуры

```go
// domain/models.go
type FileNode struct {
    Name     string      `json:"name"`
    Path     string      `json:"path"`
    IsDir    bool        `json:"isDir"`
    Children []*FileNode `json:"children,omitempty"`
    Size     int64       `json:"size"`
}
```

### 4. Ошибки — типизированные

```go
// domain/errors.go
var (
    ErrProjectNotFound = errors.New("project not found")
    ErrFileNotFound    = errors.New("file not found")
    ErrInvalidPath     = errors.New("invalid path")
)
```

### 5. Константы — в config.go

```go
// domain/config.go
const (
    MaxFileSize    = 10 * 1024 * 1024  // 10MB
    DefaultTimeout = 30 * time.Second
    MaxConcurrency = 10
)
```

## Ключевые интерфейсы

| Интерфейс            | Назначение           |
| -------------------- | -------------------- |
| `ContextBuilder`     | Построение контекста |
| `GitRepository`      | Git операции         |
| `SettingsRepository` | Настройки            |
| `ToolExecutor`       | Выполнение AI tools  |
| `AIProvider`         | AI провайдеры        |
| `LanguageAnalyzer`   | Анализ кода          |
