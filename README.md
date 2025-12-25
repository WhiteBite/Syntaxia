<p align="center">
  <img src="syntaxia-icon.png" alt="Syntaxia Logo" width="128" height="128">
</p>

<h1 align="center">Syntaxia</h1>

<p align="center">
  <strong>AI-Powered Code Context Builder</strong><br>
  Build optimized context for AI assistants from your codebase
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#screenshots">Screenshots</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#supported-ai">AI Providers</a> •
  <a href="#tech-stack">Tech Stack</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js" alt="Vue 3">
  <img src="https://img.shields.io/badge/Wails-v2-red" alt="Wails">
  <img src="https://img.shields.io/badge/Platform-Windows%20|%20macOS%20|%20Linux-blue" alt="Platform">
</p>

---

## What is Syntaxia?

Syntaxia helps you prepare code context for AI assistants like ChatGPT, Claude, or Gemini. Instead of manually copying files, Syntaxia:

- 📁 **Scans your project** with smart filtering (respects .gitignore)
- 🎯 **Selects relevant files** manually or with AI suggestions
- ⚡ **Optimizes output** by removing comments, licenses, empty lines
- 📋 **Exports context** ready to paste into any AI chat
- 💬 **Built-in AI Chat** for direct interaction with your codebase

---

## Features

### 🗂️ Smart File Browser
- Tree view with file type icons and sizes
- Filter by file type (Go, TypeScript, Vue, etc.)
- Search across all files
- Respects `.gitignore` and custom ignore rules
- File count badges for folders

### 🔀 Git Integration
- Browse local repositories
- Switch between branches
- View commit history
- Clone remote repositories (GitHub, GitLab)
- Select files from any branch/commit

### 🤖 AI-Powered Features
- **Context Suggestions**: AI recommends relevant files for your task
- **Built-in Chat**: Ask questions about your code directly
- **Multiple Providers**: OpenAI, Gemini, Qwen, OpenRouter, LocalAI

### ⚡ Context Optimization
- Remove comments and license headers
- Collapse empty lines
- Compact JSON/YAML files
- Exclude test files
- Token counting with limits

### 📤 Export Options
- Copy to clipboard
- Export to file (TXT, MD, JSON)
- Split large contexts into chunks
- Apply custom templates
- Include file metadata

### 🎨 Modern UI
- Dark theme with glow effects
- Resizable panels
- Keyboard shortcuts
- Multi-language (English, Russian)
- OS context menu integration

---

## Screenshots

### Welcome Screen
<p align="center">
  <img src="docs/screenshots/welcome-screen.png" alt="Welcome Screen" width="800">
</p>

### Main Interface
<p align="center">
  <img src="docs/screenshots/main-interface.png" alt="Main Interface" width="800">
</p>

### AI Chat
<p align="center">
  <img src="docs/screenshots/ai-chat.png" alt="AI Chat" width="800">
</p>

### Context Preview
<p align="center">
  <img src="docs/screenshots/context-built.png" alt="Context Preview" width="800">
</p>

### Git Integration
<p align="center">
  <img src="docs/screenshots/git-panel.png" alt="Git Panel" width="800">
</p>

### Settings
<p align="center">
  <img src="docs/screenshots/settings-ai.png" alt="AI Settings" width="600">
</p>

---

## Installation

### Download
Get the latest release for your platform:
- **Windows**: `SyntaxiaWB.exe`
- **macOS**: `SyntaxiaWB.app.zip`
- **Linux**: `SyntaxiaWB`

### Build from Source

```bash
# Prerequisites: Go 1.24+, Node.js 20+, Wails CLI

# Clone
git clone https://github.com/WhiteBite/Syntaxia.git
cd Syntaxia

# Development
./dev.ps1          # Windows
make dev           # Linux/macOS

# Build
./build-windows.ps1
./build-macos.sh
./build-linux.sh
```

---

## Usage

1. **Open Project**: Click "Open Project" or drag a folder
2. **Select Files**: Check files in the tree or use AI suggestions
3. **Configure**: Set token limit, enable optimizations
4. **Build Context**: Click "Build" to generate optimized context
5. **Export**: Copy to clipboard or save to file

### Keyboard Shortcuts
| Key | Action |
|-----|--------|
| `Ctrl+O` | Open project |
| `Ctrl+F` | Search files |
| `Ctrl+,` | Settings |
| `Ctrl+C` | Copy context |

---

<h2 id="supported-ai">Supported AI Providers</h2>

| Provider | Models | Features |
|----------|--------|----------|
| **OpenAI** | GPT-4o, GPT-4, GPT-3.5 | Chat, Streaming |
| **Google Gemini** | Gemini Pro, Flash | Chat, Large context |
| **Qwen** | Coder Plus (1M tokens) | Huge context window |
| **OpenRouter** | 100+ models | Model aggregator |
| **LocalAI** | Llama, Mistral, etc. | Self-hosted, Private |

---

<h2 id="tech-stack">Tech Stack</h2>

| Layer | Technology |
|-------|------------|
| **Backend** | Go 1.24, Clean Architecture |
| **Frontend** | Vue 3, TypeScript, Pinia, Tailwind CSS |
| **Desktop** | Wails v2 (Go ↔ JS bridge) |
| **Build** | GitHub Actions, Cross-platform |

---

## Contributing

Contributions are welcome! Please read the development guide:

```bash
# Run tests
cd backend && go test ./...
cd frontend && npm run test:run

# Lint
cd backend && golangci-lint run --config .golangci.yml
```

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/WhiteBite">WhiteBite</a>
</p>
