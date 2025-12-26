---
inclusion: fileMatch
fileMatchPattern: "frontend/src/features/**/*"
---

# Feature Modules - Архитектура

## Структура модуля

```
features/{name}/
├── model/              # Store, бизнес-логика
│   └── {name}.store.ts
├── ui/                 # Vue компоненты
│   ├── {Name}Panel.vue
│   └── {Name}Item.vue
├── composables/        # Локальные composables (опционально)
├── api/                # API вызовы (опционально)
├── lib/                # Утилиты модуля (опционально)
├── constants/          # Константы модуля (опционально)
└── index.ts            # Public API
```

## Существующие модули

| Модуль       | Назначение                             |
| ------------ | -------------------------------------- |
| `files/`     | Файловое дерево, выбор, фильтрация     |
| `context/`   | Построение и отображение контекста     |
| `ai-chat/`   | AI чат интерфейс                       |
| `git/`       | Git интеграция (статус, diff, история) |
| `task/`      | Taskflow управление                    |
| `templates/` | Шаблоны промптов                       |

## Public API (index.ts)

```ts
// features/files/index.ts
export { useFileStore } from "./model/file.store";
export { FileExplorer } from "./ui/FileExplorer.vue";
export type { FileNode, FileFilter } from "./types";
```

## Правила

### Store

```ts
// model/{name}.store.ts
export const useFileStore = defineStore('file', () => {
  // State
  const files = ref<FileNode[]>([])

  // Getters
  const selectedCount = computed(() => ...)

  // Actions
  async function loadFiles() { ... }

  return { files, selectedCount, loadFiles }
})
```

### Компоненты

- Используют store через `useXxxStore()`
- Не содержат бизнес-логику — только UI
- Эмитят события для родителя

### Зависимости между модулями

```
files → context (выбранные файлы → контекст)
context → ai-chat (контекст → чат)
git → files (статус файлов)
```

Импортировать только через `index.ts`:

```ts
// ✅ Правильно
import { useFileStore } from "@/features/files";

// ❌ Неправильно
import { useFileStore } from "@/features/files/model/file.store";
```
