---
inclusion: always
---

# Syntaxia - Обзор проекта

## Что это

**AI Worker Factory** — система автоматизации разработки через AI-агентов.

Входные данные: задача (фича, баг, рефакторинг)
Выходные данные: готовый код + отчёт

## Архитектура

```
Task Input → Context Builder → Taskflow Engine → AI Tools → Verification → Output
```

### Ключевые компоненты

| Компонент       | Назначение                                                     |
| --------------- | -------------------------------------------------------------- |
| Context Builder | Сканирование проекта, токенизация, форматирование              |
| Taskflow Engine | Декомпозиция задач, граф зависимостей, параллельное выполнение |
| Tool Executor   | AI инструменты (file, git, symbol, memory)                     |
| Verification    | Linting → Build → Tests → Guardrails                           |
| Self-Correction | Автоматическое исправление ошибок                              |

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
