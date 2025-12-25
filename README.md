<p align="center">
  <img src="syntaxia-icon.png" alt="Syntaxia Logo" width="128" height="128">
</p>

<h1 align="center">Syntaxia</h1>

<p align="center">
  <strong>AI-Powered Code Context Builder</strong><br>
  Build optimized context for AI assistants from your codebase
</p>

<p align="center">
  <a href="https://github.com/WhiteBite/Syntaxia/actions/workflows/build.yml">
    <img src="https://github.com/WhiteBite/Syntaxia/actions/workflows/build.yml/badge.svg" alt="Build Status">
  </a>
  <a href="https://github.com/WhiteBite/Syntaxia/releases/latest">
    <img src="https://img.shields.io/github/v/release/WhiteBite/Syntaxia?color=blue" alt="Release">
  </a>
  <a href="https://github.com/WhiteBite/Syntaxia/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/WhiteBite/Syntaxia" alt="License">
  </a>
  <img src="https://img.shields.io/github/stars/WhiteBite/Syntaxia?style=social" alt="Stars">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js&logoColor=white" alt="Vue 3">
  <img src="https://img.shields.io/badge/Wails-v2-EB4034?logo=wails&logoColor=white" alt="Wails">
  <img src="https://img.shields.io/badge/Platform-Windows%20|%20macOS%20|%20Linux-blue" alt="Platform">
</p>

<p align="center">
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-features">Features</a> •
  <a href="#-screenshots">Screenshots</a> •
  <a href="#-installation">Installation</a> •
  <a href="#-ai-providers">AI Providers</a> •
  <a href="README.ru.md">Русский</a>
</p>

---

## 🎯 What is Syntaxia?

Syntaxia helps you prepare code context for AI assistants like ChatGPT, Claude, or Gemini. Instead of manually copying files:

| Without Syntaxia | With Syntaxia |
|------------------|---------------|
| ❌ Manually copy each file | ✅ Select files in tree view |
| ❌ Exceed token limits | ✅ Auto token counting & limits |
| ❌ Include irrelevant code | ✅ AI suggests relevant files |
| ❌ Paste comments & licenses | ✅ Smart optimization removes noise |
| ❌ Switch between apps | ✅ Built-in AI chat |

---

## ⚡ Quick Start

```bash
# 1. Download latest release for your platform
# https://github.com/WhiteBite/Syntaxia/releases

# 2. Open your project folder

# 3. Select files → Build → Copy to AI chat
```

**That's it!** Your optimized context is ready for any AI assistant.

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🗂️ Smart File Browser
- Tree view with file type icons
- Filter by type (Go, TS, Vue, etc.)
- Search across all files
- Respects `.gitignore`
- Token count per file/folder

</td>
<td width="50%">

### 🔀 Git Integration
- Browse local repositories
- Switch branches
- View commit history
- Clone from GitHub/GitLab
- Select files from any commit

</td>
</tr>
<tr>
<td width="50%">

### 🤖 AI-Powered
- Context suggestions
- Built-in AI chat
- Multiple providers support
- Streaming responses
- 1M+ token context (Qwen)

</td>
<td width="50%">

### ⚡ Optimization
- Remove comments & licenses
- Collapse empty lines
- Compact JSON/YAML
- Exclude test files
- Split large contexts

</td>
</tr>
</table>

---

## 📸 Screenshots

<details>
<summary><b>🖥️ Main Interface</b> — File browser with token counting</summary>
<p align="center">
  <img src="docs/screenshots/main-interface.png" alt="Main Interface" width="800">
</p>
</details>

<details>
<summary><b>📝 Context Preview</b> — Built context ready to copy</summary>
<p align="center">
  <img src="docs/screenshots/context-built.png" alt="Context Preview" width="800">
</p>
</details>

<details>
<summary><b>💬 AI Chat</b> — Ask questions about your code</summary>
<p align="center">
  <img src="docs/screenshots/ai-chat.png" alt="AI Chat" width="800">
</p>
</details>

<details>
<summary><b>🔀 Git Panel</b> — Branch switching and history</summary>
<p align="center">
  <img src="docs/screenshots/git-panel.png" alt="Git Panel" width="800">
</p>
</details>

<details>
<summary><b>⚙️ Settings</b> — AI provider configuration</summary>
<p align="center">
  <img src="docs/screenshots/settings-ai.png" alt="Settings" width="600">
</p>
</details>

---

## 📦 Installation

### Download

| Platform | Download | Requirements |
|----------|----------|--------------|
| **Windows** | [`SyntaxiaWB.exe`](https://github.com/WhiteBite/Syntaxia/releases/latest) | Windows 10+ |
| **macOS** | [`SyntaxiaWB.app.zip`](https://github.com/WhiteBite/Syntaxia/releases/latest) | macOS 11+ |
| **Linux** | [`SyntaxiaWB`](https://github.com/WhiteBite/Syntaxia/releases/latest) | GTK3, WebKit2GTK |

### Build from Source

```bash
# Prerequisites: Go 1.24+, Node.js 20+, Wails CLI
git clone https://github.com/WhiteBite/Syntaxia.git
cd Syntaxia

# Development
./dev.ps1          # Windows
make dev           # Linux/macOS

# Production build
./build-windows.ps1
./build-macos.sh
./build-linux.sh
```

---

## 🤖 AI Providers

| Provider | Models | Max Context | Setup |
|----------|--------|-------------|-------|
| **OpenAI** | GPT-4o, GPT-4, GPT-3.5 | 128K | [Get API Key](https://platform.openai.com/api-keys) |
| **Google Gemini** | Gemini Pro, Flash | 1M | [Get API Key](https://aistudio.google.com/apikey) |
| **Qwen** | Coder Plus | 1M | [Get API Key](https://dashscope.console.aliyun.com/) |
| **OpenRouter** | 100+ models | Varies | [Get API Key](https://openrouter.ai/keys) |
| **LocalAI** | Llama, Mistral, etc. | Unlimited | Self-hosted |

---

## ⌨️ Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Ctrl+O` | Open project |
| `Ctrl+F` | Search files |
| `Ctrl+,` | Settings |
| `Ctrl+C` | Copy context |
| `Ctrl+B` | Build context |

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| **Backend** | Go 1.24, Clean Architecture |
| **Frontend** | Vue 3, TypeScript, Pinia, Tailwind CSS |
| **Desktop** | Wails v2 (Go ↔ JS bridge) |
| **CI/CD** | GitHub Actions |

---

## 🤝 Contributing

Contributions welcome! Please check:

1. Run tests before PR:
```bash
cd backend && go test ./... && golangci-lint run
cd frontend && npm run build && npm run test:run
```

2. Follow [commit conventions](https://www.conventionalcommits.org/)

---

## 📄 License

[MIT License](LICENSE) — free for personal and commercial use.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/WhiteBite">WhiteBite</a>
</p>

<p align="center">
  <a href="https://github.com/WhiteBite/Syntaxia/stargazers">⭐ Star this repo</a> if you find it useful!
</p>
