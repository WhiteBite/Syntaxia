---
inclusion: fileMatch
fileMatchPattern: "backend/application/**/*.go"
---

# Application Layer - Правила

## Назначение

Application — use cases и бизнес-логика. Оркестрирует domain и infrastructure.

## Структура

```
backend/application/
├── tools/              # AI tool handlers (file, git, symbol, memory)
├── ai/                 # AI сервисы (chat, providers)
├── analysis/           # Анализ кода (symbol index, call graph)
├── project/            # Проект сервисы
├── diff/               # Применение изменений, diff generation
├── tool_executor.go    # Главный executor для AI tools
└── constants.go        # Константы
```

## Правила

### 1. Импортирует domain, НЕ infrastructure напрямую

```go
// ✅ Правильно — через интерфейсы
type MyService struct {
    gitRepo domain.GitRepository  // интерфейс
    logger  domain.Logger
}

// ❌ Неправильно — прямая зависимость
import "syntaxia/infrastructure/git"
```

### 2. Dependency Injection через конструктор

```go
func NewMyService(
    gitRepo domain.GitRepository,
    logger domain.Logger,
) *MyService {
    return &MyService{
        gitRepo: gitRepo,
        logger:  logger,
    }
}
```

### 3. Use cases — методы сервисов

```go
// application/project/service.go
type ProjectService struct {
    contextBuilder domain.ContextBuilder
    gitRepo        domain.GitRepository
}

func (s *ProjectService) BuildProjectContext(ctx context.Context, path string, files []string) (*domain.ContextSummary, error) {
    // orchestration logic
}
```

## Tool Executor

Центральный компонент для AI tools:

```go
// application/tool_executor.go
type ToolExecutorImpl struct {
    handlerRegistry *tools.HandlerRegistry
    // ...
}

func (te *ToolExecutorImpl) ExecuteTool(call domain.ToolCall, projectRoot string) domain.ToolResult {
    return te.handlerRegistry.Execute(call.Name, call.Arguments, projectRoot)
}
```

## AI Chat с Tool Calling

Основной flow для редактирования файлов:

```go
// application/ai/chat_service.go
type ChatService struct {
    provider     domain.AIProvider
    toolExecutor domain.ToolExecutor
    changeManager *ChangeManager
}

// AI вызывает tools → изменения сохраняются → пользователь может откатить
```

## Change Manager

Управление изменениями файлов:

```go
// application/diff/
├── apply_service.go    # Применение изменений
├── diff_service.go     # Генерация diff
└── change_manager.go   # История изменений, rollback
```
