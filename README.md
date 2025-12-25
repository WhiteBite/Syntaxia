# Syntaxia

[Русская версия](README.ru.md)

AI Worker Factory — development automation system through AI agents. Input a task (project, feature, bug), get automated execution through configured pipelines.

## Concept

See [docs/CONCEPT.md](docs/CONCEPT.md) for detailed architecture.

```
Task Input → Context Builder → Taskflow Engine → AI Tools → Verification → Output
```

## Features

### Context Builder

- Project scanning with .gitignore support
- Token counting with model limits
- Output formats: Markdown, XML, JSON
- Smart file recommendations based on analysis

### Taskflow Engine

- Task decomposition into subtasks
- Dependency graph between tasks
- Parallel execution of independent tasks
- SLA policies (tokens, time, retries)

### AI Tool Executor

- `file_tools`: read, write, search files
- `git_tools`: status, diff, commit
- `symbol_tools`: list, search symbols
- `memory_tools`: save/restore context

### Verification Pipeline

- Static analysis (linting)
- Build verification
- Test execution
- Self-correction on failures

### AI Providers

- OpenAI, Gemini, OpenRouter, LocalAI, Qwen
- Model selection per provider
- Streaming responses

### Git Integration

- Branch switching
- File status indicators
- Repository cloning

## Installation

Download the latest release for your platform:

- Windows: `syntaxia-windows-amd64.exe`
- macOS: `syntaxia.app.zip`
- Linux: `syntaxia-linux-amd64`

Or build from source (requires Go 1.24+, Node.js 20+):

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone and build
git clone https://github.com/WhiteBite/syntaxia.git
cd syntaxia
wails build
```

## Usage

1. Open a project folder or clone a repository
2. Select files in the explorer (use filters to narrow down)
3. Review token count and adjust selection
4. Copy context to clipboard or export to PDF
5. Paste into your AI assistant

## Tech Stack

- Backend: Go, Wails
- Frontend: Vue 3, TypeScript, Pinia, Tailwind CSS
- Build: Vite, GitHub Actions

## License

MIT
