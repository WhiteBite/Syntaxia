#!/usr/bin/env pwsh
# Build script for Windows

$ErrorActionPreference = "Stop"

Write-Host "=== Building Syntaxia for Windows ===" -ForegroundColor Cyan
Write-Host ""

# Get version from git tag or use dev
$version = "dev"
$gitCommit = "unknown"
$buildDate = Get-Date -Format "yyyy-MM-dd"

try {
    $gitTag = git describe --tags --abbrev=0 2>$null
    if ($gitTag) {
        $version = $gitTag
    }
    $gitCommit = git rev-parse --short HEAD 2>$null
    if (-not $gitCommit) {
        $gitCommit = "unknown"
    }
} catch {
    Write-Host "Warning: Could not get git info, using defaults" -ForegroundColor Yellow
}

Write-Host "Version: $version" -ForegroundColor White
Write-Host "Commit: $gitCommit" -ForegroundColor White
Write-Host ""

# Build ldflags for version injection and hide console window
$ldflags = "-X syntaxia/infrastructure/version.Version=$version -X syntaxia/infrastructure/version.GitCommit=$gitCommit -X syntaxia/infrastructure/version.BuildDate=$buildDate -H windowsgui"

# Find wails (same logic as dev.ps1)
$wailsPath = $null
if (Test-Path "$env:USERPROFILE\go\bin\wails.exe") {
    $wailsPath = "$env:USERPROFILE\go\bin\wails.exe"
} elseif (Get-Command wails -ErrorAction SilentlyContinue) {
    $wailsPath = "wails"
} else {
    Write-Host "❌ Wails not found. Install with: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Red
    exit 1
}

# Ensure build directory exists
New-Item -ItemType Directory -Force -Path "build/bin" | Out-Null

Push-Location backend

try {
    # Build frontend first
    Write-Host "Building frontend..." -ForegroundColor Yellow
    Push-Location ../frontend
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "✗ Frontend build failed!" -ForegroundColor Red
        exit 1
    }
    Pop-Location
    
    # Copy frontend dist to backend/frontend/dist for embedding
    Write-Host "Preparing frontend for embedding..." -ForegroundColor Yellow
    if (Test-Path "frontend/dist") {
        Remove-Item "frontend/dist" -Recurse -Force -ErrorAction SilentlyContinue
    }
    New-Item -ItemType Directory -Force -Path "frontend" | Out-Null
    Copy-Item "../frontend/dist" "frontend/dist" -Recurse -Force
    
    Write-Host "Building Windows executable..." -ForegroundColor Green
    & $wailsPath build -clean -skipbindings -ldflags "$ldflags"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✓ Build successful!" -ForegroundColor Green
        
        # Copy to root build/bin
        if (Test-Path "build/bin/Syntaxia.exe") {
            New-Item -ItemType Directory -Force -Path "../build/bin" | Out-Null
            Copy-Item "build/bin/Syntaxia.exe" "../build/bin/Syntaxia.exe" -Force
            
            $exePath = Get-Item "../build/bin/Syntaxia.exe"
            $size = [math]::Round($exePath.Length / 1MB, 2)
            Write-Host "Output: $($exePath.FullName)" -ForegroundColor White
            Write-Host "Size: ${size}MB" -ForegroundColor White
        }
    } else {
        Write-Host ""
        Write-Host "✗ Build failed!" -ForegroundColor Red
        exit 1
    }
} finally {
    Pop-Location
}
