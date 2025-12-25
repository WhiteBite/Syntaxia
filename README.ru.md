# Syntaxia

[English version](README.md)

AI Worker Factory — система автоматизации разработки через AI-агентов. На вход подаётся задача (проект, фича, баг), система автоматически выполняет её через настроенные пайплайны.

## Концепция

Подробная архитектура: [docs/CONCEPT.md](docs/CONCEPT.md)

```
Задача → Context Builder → Taskflow Engine → AI Tools → Верификация → Результат
```

## Возможности

### Context Builder

- Сканирование проекта с поддержкой .gitignore
- Подсчёт токенов с учётом лимитов модели
- Форматы вывода: Markdown, XML, JSON
- Умные рекомендации файлов на основе анализа

### Taskflow Engine

- Декомпозиция задачи на подзадачи
- Граф зависимостей между задачами
- Параллельное выполнение независимых задач
- SLA политики (токены, время, ретраи)

### AI Tool Executor

- `file_tools`: чтение, запись, поиск файлов
- `git_tools`: статус, diff, commit
- `symbol_tools`: список, поиск символов
- `memory_tools`: сохранение/восстановление контекста

### Verification Pipeline

- Статический анализ (linting)
- Проверка сборки
- Запуск тестов
- Self-correction при ошибках

### AI Провайдеры

- OpenAI, Gemini, OpenRouter, LocalAI, Qwen
- Выбор модели для каждого провайдера
- Потоковые ответы

### Git-интеграция

- Переключение веток
- Индикаторы статуса файлов
- Клонирование репозиториев

## Установка

Скачайте последний релиз для вашей платформы:

- Windows: `syntaxia-windows-amd64.exe`
- macOS: `syntaxia.app.zip`
- Linux: `syntaxia-linux-amd64`

Или соберите из исходников (требуется Go 1.24+, Node.js 20+):

```bash
# Установка Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Клонирование и сборка
git clone https://github.com/WhiteBite/syntaxia.git
cd syntaxia
wails build
```

## Использование

1. Откройте папку проекта или клонируйте репозиторий
2. Выберите файлы в проводнике (используйте фильтры)
3. Проверьте количество токенов и скорректируйте выбор
4. Скопируйте контекст в буфер или экспортируйте в PDF
5. Вставьте в ваш AI-ассистент

## Технологии

- Backend: Go, Wails
- Frontend: Vue 3, TypeScript, Pinia, Tailwind CSS
- Сборка: Vite, GitHub Actions

## Лицензия

MIT
