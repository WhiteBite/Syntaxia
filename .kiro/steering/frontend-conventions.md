---
inclusion: fileMatch
fileMatchPattern: "frontend/**/*.{vue,ts,tsx}"
---

# Frontend Vue - Конвенции

## Структура

```
frontend/src/
├── features/         # Feature modules (основная логика)
│   └── [name]/
│       ├── model/    # Store, бизнес-логика
│       ├── ui/       # Vue компоненты
│       └── index.ts  # Public API
├── components/       # Shared компоненты
├── composables/      # Vue composables
├── stores/           # Global stores (ui, settings, project)
├── services/         # API сервисы
├── types/            # TypeScript типы
├── locales/          # i18n (ru/, en/)
└── assets/styles/    # CSS модули
```

## Именование

| Тип         | Формат             | Пример             |
| ----------- | ------------------ | ------------------ |
| Компоненты  | PascalCase.vue     | `FileTreeNode.vue` |
| Composables | useCamelCase.ts    | `useFileTree.ts`   |
| Stores      | camelCase.store.ts | `file.store.ts`    |
| CSS классы  | kebab-case         | `.panel-header`    |

## CSS

### Используй классы из `assets/styles/`:

- Кнопки: `.btn`, `.btn-primary`, `.btn-ghost`, `.action-btn`
- Инпуты: `.input`, `.search-input`
- Карточки: `.card`, `.card-hover`
- Панели: `.panel-header`, `.panel-title`
- Бейджи: `.badge`, `.badge-primary`

### Правила:

1. НЕ дублировать стили inline
2. Tailwind только для layout (flex, grid, p-_, m-_)
3. Использовать CSS переменные из `:root`

## Локализация (i18n)

```ts
// Использование
const { t } = useI18n();
t("files.selectAll");
```

Файлы: `locales/ru/*.json`, `locales/en/*.json`

**Правило:** Новые тексты добавлять в ОБА языка.

## Stores

```ts
// Feature store
export const useFileStore = defineStore("file", () => {
  const tree = useFileTree();
  const selection = useFileSelection();
  // композиция composables
});
```

## Wails API вызовы

```ts
import { App } from "@/wailsjs/go/main/App";

await App.OpenProject(path);
await App.BuildContext(options);
```

## Запрещено

- `console.log()` в продакшн коде
- `any` типы
- Магические числа (выносить в `config/constants.ts`)
- Дублирование кода

## Тестирование

```bash
npm run test:run
```

Структура: `tests/unit/`, `tests/integration/`
