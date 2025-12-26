# Правило: Проверки перед коммитом

**ОБЯЗАТЕЛЬНО** выполнять проверки перед каждым коммитом и пушем.

## Backend (Go)

```bash
cd backend
go build ./...
go test ./... -count=1
golangci-lint run --config .golangci.yml
```

## Frontend (Vue/TypeScript)

```bash
cd frontend
npm run build
npm run test:run
```

## Порядок действий

1. Внести изменения в код
2. Запустить проверки для изменённых частей (backend/frontend)
3. Исправить все ошибки линтера и тестов
4. Только после успешных проверок делать `git add` и `git commit`
5. Пушить в remote

## Критично

- **НЕ ПУШИТЬ** если golangci-lint выдаёт ошибки
- **НЕ ПУШИТЬ** если `go build` или `npm run build` падают
- **НЕ ПУШИТЬ** если тесты не проходят

## Быстрая проверка (Windows)

```powershell
# Backend
cd backend; go build ./...; golangci-lint run --config .golangci.yml

# Frontend  
cd frontend; npm run build
```
