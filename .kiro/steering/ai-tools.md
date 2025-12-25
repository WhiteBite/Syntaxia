---
inclusion: fileMatch
fileMatchPattern: "backend/application/tools/**/*.go|backend/**/tool*.go"
---

# AI Tools - Конвенции

## Обзор

AI-агенты используют инструменты через `ToolExecutor`. Каждый tool handler — отдельный файл.

## Доступные инструменты

| Категория   | Tools                                                       | Файл                         |
| ----------- | ----------------------------------------------------------- | ---------------------------- |
| File        | `read_file`, `write_file`, `search_files`, `list_directory` | `file_tools.go`              |
| Git         | `git_status`, `git_diff`, `git_log`, `git_commit`           | `git_tools.go`               |
| Symbol      | `list_symbols`, `search_symbols`                            | `symbol_tools.go`            |
| Memory      | `save_context`, `find_context`                              | `memory_tools.go`            |
| CallGraph   | `get_callers`, `get_callees`                                | `callgraph_tools.go`         |
| Preferences | `get_preference`, `set_preference`                          | `preferences_tools.go`       |
| Semantic    | `semantic_search`                                           | `semantic_tools.go`          |
| Project     | `get_project_structure`                                     | `project_structure_tools.go` |

## Структура Tool Handler

```go
// handler.go - интерфейс
type ToolHandler interface {
    GetTools() []domain.Tool
    Execute(toolName string, args map[string]any, projectRoot string) (string, error)
}
```

## Создание нового Tool

1. Создать файл `backend/application/tools/{name}_tools.go`
2. Реализовать `ToolHandler` интерфейс
3. Зарегистрировать в `tool_executor.go` → `registerHandlers()`
4. Создать тесты `{name}_tools_test.go`

### Пример

```go
package tools

type MyToolsHandler struct {
    logger domain.Logger
}

func NewMyToolsHandler(logger domain.Logger) *MyToolsHandler {
    return &MyToolsHandler{logger: logger}
}

func (h *MyToolsHandler) GetTools() []domain.Tool {
    return []domain.Tool{
        {
            Name:        "my_tool",
            Description: "Does something useful",
            Parameters: domain.ToolParameters{
                Type: "object",
                Properties: map[string]domain.ToolProperty{
                    "param1": {Type: "string", Description: "First param"},
                },
                Required: []string{"param1"},
            },
        },
    }
}

func (h *MyToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
    switch toolName {
    case "my_tool":
        return h.myTool(args, projectRoot)
    default:
        return "", fmt.Errorf("unknown tool: %s", toolName)
    }
}
```

## Правила

1. Каждый handler — один файл с тестами
2. Логировать вызовы через `logger.Info()`
3. Возвращать понятные ошибки с контекстом
4. Валидировать все входные параметры
5. Ограничивать размер возвращаемых данных
