---
inclusion: always
---

# Syntaxia - Обзор проекта

## Что это

**AI-Powered Code Context Builder** — инструмент для подготовки оптимизированного контекста из кодовой базы для AI-ассистентов.

## Основной функционал

1. **Smart Context Builder** — умный сбор контекста с анализом зависимостей
2. **AI Chat с Tool Calling** — чат который может редактировать файлы
3. **Change Preview & Rollback** — просмотр изменений с возможностью отката

## Архитектура

```
File Selection → Smart Context → AI Chat → File Changes → Preview/Rollback
```

### Ключевые компоненты

| Компонент            | Назначение                                          |
| -------------------- | --------------------------------------------------- |
| Context Builder      | Сканирование проекта, токенизация, форматирование   |
| Smart Suggestions    | AI подсказки релевантных файлов на основе анализа   |
| Symbol/Import Graph  | Анализ зависимостей для умного сбора контекста      |
| AI Chat              | Чат с tool calling для редактирования файлов        |
| Change Manager       | Просмотр diff, применение/откат изменений           |

## Стек технологий

| Слой     | Технологии                                |
| -------- | ----------------------------------------- |
| Backend  | Go 1.24+, Wails v2                        |
| Frontend | Vue 3, TypeScript, Pinia, Tailwind CSS    |
| Desktop  | Wails (Go ↔ JS bridge)                    |
| AI       | OpenAI, Gemini, OpenRouter, LocalAI, Qwen |

## Структура проекта

```
syntaxia/
├── backend/                 # Go backend (Clean Architecture)
│   ├── domain/              # Бизнес-модели, интерфейсы
│   ├── application/         # Use cases, сервисы, tools
│   ├── infrastructure/      # Реализации (AI, Git, FS, DB)
│   ├── internal/            # Внутренние сервисы
│   ├── handlers/            # HTTP/Wails handlers
│   └── *_api.go             # Wails API методы
│
├── frontend/                # Vue 3 frontend
│   └── src/
│       ├── features/        # Feature modules
│       ├── components/      # Shared компоненты
│       ├── composables/     # Vue composables
│       ├── stores/          # Global Pinia stores
│       ├── locales/         # i18n (ru, en)
│       └── config/          # Константы
│
├── .kiro/                   # AI agent конфигурация
│   └── steering/            # Правила для агентов
│
└── docs/                    # Документация
```

## Команды

```bash
# Development
./dev.ps1              # Windows
make dev               # Linux/Mac

# Build
./build-windows.ps1
./build-linux.sh
./build-macos.sh

# Tests
cd backend && go test ./...
cd frontend && npm run test:run
```
