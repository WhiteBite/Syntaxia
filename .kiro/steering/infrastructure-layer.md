---
inclusion: fileMatch
fileMatchPattern: "backend/infrastructure/**/*.go"
---

# Infrastructure Layer - Правила

## Назначение

Infrastructure — реализации интерфейсов из domain. Работа с внешним миром:

- Файловая система
- Git
- AI провайдеры
- База данных
- Внешние API

## Структура

```
backend/infrastructure/
├── ai/                 # AI провайдеры (OpenAI, Gemini, etc.)
├── analyzers/          # Языковые анализаторы (Go, TS, Python)
├── git/                # Git операции
├── contextbuilder/     # Построение контекста
├── memory/             # Контекстная память
├── settingsfs/         # Хранение настроек
├── filesystem/         # Файловые операции
├── fsscanner/          # Сканирование директорий
├── formatters/         # Форматирование вывода
├── embeddings/         # Векторные embeddings
├── buildpipeline/      # Build verification
├── testengine/         # Test execution
├── staticanalyzer/     # Static analysis
└── ...
```

## Правила

### 1. Реализует интерфейсы из domain

```go
// infrastructure/git/repository.go
type GitRepositoryImpl struct {
    logger domain.Logger
}

// Реализует domain.GitRepository
func (r *GitRepositoryImpl) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
    // реализация
}
```

### 2. Импортирует только domain

```go
// ✅ Правильно
import "syntaxia/domain"

// ❌ Неправильно
import "syntaxia/application/..."
```

### 3. Каждый пакет — одна ответственность

```
ai/           → только AI провайдеры
git/          → только Git операции
analyzers/    → только анализ кода
```

### 4. Конструкторы возвращают интерфейс

```go
func NewGitRepository(logger domain.Logger) domain.GitRepository {
    return &GitRepositoryImpl{logger: logger}
}
```

## AI Провайдеры

```
infrastructure/ai/
├── provider.go         # Базовый интерфейс
├── openai.go           # OpenAI
├── gemini.go           # Google Gemini
├── openrouter.go       # OpenRouter
├── localai.go          # LocalAI
└── qwen.go             # Qwen
```

Каждый провайдер реализует `domain.AIProvider`:

```go
type AIProvider interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)
    GetModels() ([]string, error)
}
```

## Анализаторы кода

```
infrastructure/analyzers/
├── registry.go         # Реестр анализаторов
├── go_analyzer.go      # Go
├── typescript_analyzer.go
├── python_analyzer.go
├── java_analyzer.go
└── ...
```

Каждый реализует `analysis.LanguageAnalyzer`:

```go
type LanguageAnalyzer interface {
    GetLanguage() string
    AnalyzeFile(content, filePath string) (*AnalysisResult, error)
    ExtractSymbols(content, filePath string) ([]Symbol, error)
}
```
