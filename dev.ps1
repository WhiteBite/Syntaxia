# Syntaxia - Development Script
param(
    [switch]$Verbose,
    [int]$NodeMemory = 2048  # MB, увеличено для стабильной сборки
)

Write-Host "🚀 Запуск Syntaxia..." -ForegroundColor Green

# Найти wails
$wails = Get-Command wails -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source
if (-not $wails) {
    $gopath = $env:GOPATH
    if (-not $gopath) { $gopath = "$env:USERPROFILE\go" }
    $wailsPath = "$gopath\bin\wails.exe"
    if (Test-Path $wailsPath) { $wails = $wailsPath }
}

# Проверка зависимостей
$missing = @()
if (-not $wails) { $missing += "wails (go install github.com/wailsapp/wails/v2/cmd/wails@latest)" }
if (-not (Get-Command node -ErrorAction SilentlyContinue)) { $missing += "node" }
if (-not (Get-Command go -ErrorAction SilentlyContinue)) { $missing += "go" }

if ($missing.Count -gt 0) {
    Write-Host "❌ Не установлено: $($missing -join ', ')" -ForegroundColor Red
    exit 1
}

Push-Location backend
$env:GOGC = "50"
$env:NODE_OPTIONS = "--max-old-space-size=$NodeMemory"

try {
    $wailsArgs = @("dev", "-loglevel", "error")
    
    if ($Verbose) {
        & $wails dev
    } else {
        Write-Host "ℹ️  Флаги: -Verbose, -NodeMemory <MB> (default 512)" -ForegroundColor Gray
        & $wails @wailsArgs 2>&1 | Where-Object { 
            $_ -and $_ -notmatch "KnownStructs:|Not found: time\.Time|^\s*$" 
        }
    }
} finally {
    Pop-Location
    Remove-Item Env:GOGC -ErrorAction SilentlyContinue
    Remove-Item Env:NODE_OPTIONS -ErrorAction SilentlyContinue
}
