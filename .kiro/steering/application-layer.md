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
├── tools/              # AI tool handlers
│   ├── file_tools.go
│   ├── git_tools.go
│   ├── symbol_tools.go
│   ├── memory_tools.go
│   └── ...
├── ai/                 # AI сервисы
├── analysis/           # Анализ кода
├── project/            # Проект сервисы
├── taskflow/           # Taskflow engine
├── verification/       # Verification pipeline
├── repair/             # Self-correction
├── guardrails/         # Guardrails
├── tool_executor.go    # Главный executor
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

## Taskflow Engine

Управление задачами:

```go
// application/taskflow/
├── engine.go           # Основной движок
├── decomposer.go       # Декомпозиция задач
├── scheduler.go        # Планировщик
└── executor.go         # Выполнение
```

## Verification Pipeline

```go
// application/verification/
├── pipeline.go         # Основной pipeline
├── stages.go           # Этапы верификации
└── reporter.go         # Отчёты

// Этапы: analysis → linting → building → testing → guardrails
```
