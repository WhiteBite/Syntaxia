---
inclusion: fileMatch
fileMatchPattern: "backend/**/*.go"
---

# Backend Go - Конвенции

## Архитектура (Clean Architecture)

```
backend/
├── domain/           # Слой 1: Бизнес-сущности, интерфейсы (НЕ импортирует другие слои)
├── application/      # Слой 2: Use cases (импортирует только domain)
├── infrastructure/   # Слой 3: Реализации (импортирует domain)
├── handlers/         # Слой 4: API handlers (импортирует application, domain)
├── internal/         # Внутренние сервисы
└── *_api.go          # Wails API методы
```

## Правила импорта

- `domain` → ничего не импортирует из проекта
- `application` → только `domain`
- `infrastructure` → только `domain`
- `handlers` → `application` + `domain`

## Именование

| Тип        | Формат          | Пример                            |
| ---------- | --------------- | --------------------------------- |
| Файлы      | snake_case      | `context_builder.go`              |
| Пакеты     | lowercase       | `analyzers`                       |
| Интерфейсы | Descriptive/-er | `LanguageAnalyzer`, `SymbolIndex` |
| Реализации | +Impl           | `ToolExecutorImpl`                |
| Публичные  | PascalCase      | `BuildContext`                    |
| Приватные  | camelCase       | `parseFile`                       |

## Структура файла

```go
package name

import (
    "стандартные"

    "внешние"

    "syntaxia/..."
)

const (...)
type MyStruct struct {...}
func NewMyStruct() *MyStruct {...}
func (s *MyStruct) Method() {...}
func helperFunc() {...}
```

## Error Handling

```go
// ✅ Правильно - wrap с контекстом
if err != nil {
    return fmt.Errorf("failed to read file %s: %w", path, err)
}

// ❌ Неправильно - потеря контекста
if err != nil {
    return err
}
```

## Tool Handlers (application/tools/)

```go
// Каждый handler реализует интерфейс
type ToolHandler interface {
    Execute(params map[string]interface{}) (interface{}, error)
}
```

Файлы: `file_tools.go`, `git_tools.go`, `memory_tools.go`, `symbol_tools.go`

## Тестирование

```go
// Table-driven tests
func TestSearchFiles(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
        wantErr  bool
    }{
        {"exact match", "main.go", 1, false},
        {"invalid", "[bad", 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test
        })
    }
}
```

## Безопасность

1. API ключи → secure_storage
2. Валидация путей (path traversal)
3. Лимиты размера файлов
4. Таймауты для внешних вызовов
