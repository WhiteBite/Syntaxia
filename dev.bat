@echo off
chcp 65001 > nul

echo 🚀 Запуск Syntaxia в режиме разработки...

rem Проверка зависимостей
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 'go' не найден. Установите Go.
    goto:eof
)

where node >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 'node' не найден. Установите Node.js.
    goto:eof
)

where wails >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 'wails' не найден. Установите: go install github.com/wailsapp/wails/v2/cmd/wails@latest
    goto:eof
)

rem Переход в папку бэкенда и запуск
pushd backend

set "GOGC=50"
set "NODE_OPTIONS=--max-old-space-size=2048"

wails dev -loglevel error

popd

rem Очистка
set "GOGC="
set "NODE_OPTIONS="

echo ✨ Работа завершена.
