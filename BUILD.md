# Build Instructions

## Prerequisites

- Go 1.21+
- Node.js 18+
- Wails CLI v2: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Platform-specific requirements:

**Windows:**

- Visual Studio Build Tools or MinGW-w64

**macOS:**

- Xcode Command Line Tools

**Linux:**

- gcc, pkg-config, gtk3, webkit2gtk

## Quick Build

### Windows

```powershell
.\build-windows.ps1
```

### Linux/macOS

```bash
./build-linux.sh    # Linux
./build-macos.sh    # macOS
```

## Cross-Platform Build

Build for all platforms:

**PowerShell (Windows/Linux/macOS):**

```powershell
.\build-all.ps1
```

**Bash (Linux/macOS):**

```bash
./build-all.sh
```

### Options

Skip specific platforms:

```powershell
.\build-all.ps1 -SkipWindows
.\build-all.ps1 -SkipMacOS
.\build-all.ps1 -SkipLinux
```

```bash
./build-all.sh --skip-windows
./build-all.sh --skip-macos
./build-all.sh --skip-linux
```

Clean build:

```powershell
.\build-all.ps1 -Clean
```

```bash
./build-all.sh --clean
```

## Development Build

For development with hot reload:

```powershell
.\dev.ps1
```

Or manually:

```bash
cd backend
wails dev
```

## Output

Build artifacts are placed in `build/bin/`:

- Windows: `Syntaxia.exe`
- macOS: `Syntaxia.app`
- Linux: `Syntaxia`

## Troubleshooting

### Windows: Missing compiler

Install Visual Studio Build Tools or MinGW-w64

### Linux: Missing dependencies

```bash
sudo apt install gcc pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
```

### macOS: Code signing issues

For development builds, you can skip signing:

```bash
wails build -skipbindings -nosyncgomod
```

## Release Build

For production releases with optimizations:

```bash
cd backend
wails build -clean -upx -webview2 embed
```

Options:

- `-clean`: Clean build cache
- `-upx`: Compress with UPX
- `-webview2 embed`: Embed WebView2 runtime (Windows)
