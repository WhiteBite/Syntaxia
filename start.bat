@echo off
chcp 65001 > nul

echo 🚀 Запуск Syntaxia...

set "APP_PATH=build\bin\Syntaxia.exe"

if exist "%APP_PATH%" (
    echo ✅ Запуск приложения...
    start "" "%APP_PATH%"
) else (
    echo ❌ Приложение не собрано. Запустите dev.ps1 для разработки или соберите проект.
    exit /b 1
)
