#!/usr/bin/env pwsh
# Build script for current platform only
# Note: Cross-platform builds require running on each target OS

param(
    [switch]$Clean
)

$ErrorActionPreference = "Stop"

Write-Host "=== Syntaxia - Build ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "⚠️  Note: Building for current platform only (Windows)" -ForegroundColor Yellow
Write-Host "   For macOS/Linux builds, run this script on those platforms" -ForegroundColor Yellow
Write-Host ""

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

# Clean old build artifacts
if ($Clean) {
    Write-Host "Cleaning build directory..." -ForegroundColor Yellow
    if (Test-Path "build") {
        Remove-Item -Recurse -Force "build" -ErrorAction SilentlyContinue
    }
    if (Test-Path "backend/build") {
        Remove-Item -Recurse -Force "backend/build" -ErrorAction SilentlyContinue
    }
}

# Ensure build directory exists
New-Item -ItemType Directory -Force -Path "build/bin" | Out-Null

# Change to backend directory
Push-Location backend

try {
    # Copy frontend dist to backend/frontend/dist for embedding
    Write-Host "Preparing frontend for embedding..." -ForegroundColor Yellow
    if (Test-Path "frontend/dist") {
        Remove-Item "frontend/dist" -Recurse -Force -ErrorAction SilentlyContinue
    }
    Copy-Item "../frontend/dist" "frontend/dist" -Recurse -Force
    
    Write-Host "Building for Windows..." -ForegroundColor Green
    & $wailsPath build -clean -skipbindings
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Build complete" -ForegroundColor Green
        
        # Move from backend/build/bin to root build/bin
        if (Test-Path "build/bin/Syntaxia.exe") {
            Write-Host "Moving executable to build/bin/..." -ForegroundColor Yellow
            Copy-Item "build/bin/Syntaxia.exe" "../build/bin/Syntaxia.exe" -Force
        }
    } else {
        Write-Host "✗ Build failed" -ForegroundColor Red
        exit 1
    }
    Write-Host ""

    Write-Host "=== Build Summary ===" -ForegroundColor Cyan
    Write-Host ""
    
    # List built files
    if (Test-Path "../build/bin") {
        Get-ChildItem -Path "../build/bin" -Recurse -File | ForEach-Object {
            $size = [math]::Round($_.Length / 1MB, 2)
            Write-Host "  $($_.Name) - ${size}MB" -ForegroundColor White
        }
    }
    
    Write-Host ""
    Write-Host "Build artifacts are in: build/bin/" -ForegroundColor Cyan

} finally {
    Pop-Location
}
