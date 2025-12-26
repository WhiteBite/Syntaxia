---
inclusion: fileMatch
fileMatchPattern: "frontend/src/locales/**/*|frontend/**/*.vue"
---

# Локализация (i18n) - Правила

## Структура

```
frontend/src/locales/
├── ru/                 # Русский
│   ├── common.json
│   ├── files.json
│   ├── context.json
│   ├── chat.json
│   ├── git.json
│   ├── settings.json
│   └── ...
├── en/                 # English
│   └── (аналогичная структура)
├── i18n.ts             # Конфигурация vue-i18n
└── index.ts            # Экспорт
```

## Использование

```vue
<script setup lang="ts">
import { useI18n } from "@/composables/useI18n";
const { t } = useI18n();
</script>

<template>
  <span>{{ t("files.selectAll") }}</span>
</template>
```

## Правила

### Обязательно

1. **Все тексты через `t()`** — никаких хардкод строк в UI
2. **Оба языка** — добавлять в `ru/` И `en/` одновременно
3. **Группировка по модулям** — `files.*`, `context.*`, `chat.*`

### Именование ключей

```json
// ✅ Правильно
{
  "files": {
    "selectAll": "Выбрать все",
    "deselectAll": "Снять выбор",
    "search": {
      "placeholder": "Поиск файлов..."
    }
  }
}

// ❌ Неправильно
{
  "selectAllFiles": "...",
  "files_search_placeholder": "..."
}
```

### Плейсхолдеры

```json
{
  "files": {
    "selected": "Выбрано {count} файлов",
    "tokenCount": "{tokens} токенов"
  }
}
```

```ts
t("files.selected", { count: 5 });
```

### Множественное число

```json
{
  "files": {
    "count": "нет файлов | {n} файл | {n} файла | {n} файлов"
  }
}
```

## Добавление нового текста

1. Добавить ключ в `locales/ru/{module}.json`
2. Добавить ключ в `locales/en/{module}.json`
3. Использовать через `t('module.key')`

## Проверка

При сборке проверяется что все ключи есть в обоих языках.
