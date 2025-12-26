---
inclusion: fileMatch
fileMatchPattern: "backend/*_api.go|frontend/**/wailsjs/**"
---

# Wails API - Конвенции

## API файлы (backend/)

| Файл              | Методы                                         |
| ----------------- | ---------------------------------------------- |
| `project_api.go`  | OpenProject, GetFileTree, GetProjectInfo       |
| `context_api.go`  | BuildContext, GetContextContent, DeleteContext |
| `ai_api.go`       | AgenticChat, GenerateCode, StreamResponse      |
| `git_api.go`      | GetGitStatus, GetGitDiff, GetGitLog            |
| `analysis_api.go` | AnalyzeProject, GetSymbols, SearchSymbols      |
| `settings_api.go` | GetSettings, SaveSettings, GetAPIKey           |
| `taskflow_api.go` | ExecuteTaskflow, GetTaskStatus                 |
| `window_api.go`   | Minimize, Maximize, Close                      |
| `apply_api.go`    | ApplyChanges, PreviewChanges                   |

## Сигнатуры методов

```go
// С ошибкой
func (a *App) BuildContext(options ContextBuildOptions) (*ContextResult, error)

// Без ошибки (простые)
func (a *App) GetVersion() string

// JSON для frontend
func (a *App) AgenticChat(requestJSON string) (string, error)
```

## Именование методов

| Префикс                | Действие         |
| ---------------------- | ---------------- |
| `Get*`                 | Получение данных |
| `Set*` / `Save*`       | Сохранение       |
| `Delete*` / `Remove*`  | Удаление         |
| `Build*` / `Generate*` | Создание         |
| `Execute*` / `Run*`    | Выполнение       |

## Frontend вызовы

```ts
import { App } from "@/wailsjs/go/main/App";

// Вызов
const result = await App.BuildContext(options);

// Обработка ошибок
try {
  await App.OpenProject(path);
} catch (err) {
  console.error("Failed:", err);
}
```

## События (Streaming)

```go
// Backend
runtime.EventsEmit(a.ctx, "ai:stream:chunk", chunk)
```

```ts
// Frontend
import { EventsOn, EventsOff } from "@/wailsjs/runtime/runtime";

EventsOn("ai:stream:chunk", (chunk) => {
  /* handle */
});
onUnmounted(() => EventsOff("ai:stream:chunk"));
```

## Добавление нового метода

1. Добавить метод в соответствующий `*_api.go`
2. Запустить `wails generate module`
3. Использовать через `App.MethodName()`
