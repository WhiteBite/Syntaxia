# Syntaxia - Quick Start Script
# Быстрый запуск приложения

Write-Host "🚀 Запуск Syntaxia..." -ForegroundColor Green

# Проверяем наличие собранного приложения
$appPath = "backend\build\bin\syntaxia.exe"

if (Test-Path $appPath) {
    Write-Host "✅ Запуск приложения..." -ForegroundColor Green
    & $appPath
} else {
    Write-Host "❌ Приложение не собрано. Запустите dev.ps1 для разработки или соберите проект." -ForegroundColor Red
    exit 1
}
