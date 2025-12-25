<p align="center">
  <img src="shotgun-icon.png" alt="Syntaxia Logo" width="128" height="128">
</p>

<h1 align="center">Syntaxia</h1>

<p align="center">
  <strong>AI-Powered Code Context Builder for LLMs</strong>
</p>

<p align="center">
  <a href="https://github.com/WhiteBite/Syntaxia/releases"><img src="https://img.shields.io/github/v/release/WhiteBite/Syntaxia?style=flat-square" alt="Release"></a>
  <a href="https://github.com/WhiteBite/Syntaxia/blob/main/LICENSE"><img src="https://img.shields.io/github/license/WhiteBite/Syntaxia?style=flat-square" alt="License"></a>
  <a href="https://github.com/WhiteBite/Syntaxia/stargazers"><img src="https://img.shields.io/github/stars/WhiteBite/Syntaxia?style=flat-square" alt="Stars"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#tech-stack">Tech Stack</a> •
  <a href="README.ru.md">Русский</a>
</p>

---

**Syntaxia** is a desktop application that helps developers prepare code context for AI assistants like ChatGPT, Claude, Gemini, and Copilot. Select files from your project, optimize token usage, and export perfectly formatted context for any LLM.

## Why Syntaxia?

- 🎯 **Smart File Selection** — Filter by language, size, or custom patterns
- 📊 **Token Counting** — Real-time token estimation for GPT-4, Claude, Gemini
- 🔄 **Multiple Export Formats** — Markdown, XML, JSON, PDF
- 🌳 **Git Integration** — Work with branches, view diffs, clone repos
- 🤖 **AI Provider Support** — OpenAI, Gemini, OpenRouter, LocalAI, Qwen
- 💾 **Context Memory** — Save and reuse context configurations
- ⚡ **Fast & Native** — Built with Go and Vue 3, runs on Windows, macOS, Linux

## Features

### 🗂️ Context Builder
Build optimized code context for AI assistants:
- Project scanning with `.gitignore` support
- Smart file recommendations based on code analysis
- Token counting with model-specific limits
- Chunked export for large codebases

### 🔧 AI Tools Integration
Built-in tools for AI-assisted development:
- File operations (read, write, search)
- Git operations (status, diff, commit)
- Symbol analysis (functions, classes, imports)
- Context memory (save/restore sessions)

### 📈 Verification Pipeline
Ensure code quality before AI processing:
- Static analysis integration
- Build verification
- Test execution
- Self-correction on failures

## Installation

### Download

Get the latest release for your platform:

| Platform | Download |
|----------|----------|
| Windows | [syntaxia-windows-amd64.exe](https://github.com/WhiteBite/Syntaxia/releases/latest) |
| macOS | [Syntaxia.app.zip](https://github.com/WhiteBite/Syntaxia/releases/latest) |
| Linux | [syntaxia-linux-amd64](https://github.com/WhiteBite/Syntaxia/releases/latest) |

### Build from Source

Requirements: Go 1.24+, Node.js 20+, npm

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone repository
git clone https://github.com/WhiteBite/Syntaxia.git
cd Syntaxia

# Build for your platform
./build-windows.ps1  # Windows
./build-macos.sh     # macOS
./build-linux.sh     # Linux
```

## Usage

1. **Open Project** — Select a folder or clone a Git repository
2. **Select Files** — Use filters to find relevant files
3. **Review Tokens** — Check token count against model limits
4. **Export Context** — Copy to clipboard or save as file
5. **Use with AI** — Paste into ChatGPT, Claude, or any LLM

## Tech Stack

| Layer | Technologies |
|-------|-------------|
| Backend | Go 1.24, Wails v2 |
| Frontend | Vue 3, TypeScript, Pinia, Tailwind CSS |
| Build | Vite, GitHub Actions |
| AI | OpenAI, Gemini, OpenRouter, LocalAI, Qwen |

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## License

[MIT License](LICENSE) — feel free to use in personal and commercial projects.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/WhiteBite">WhiteBite</a>
</p>
