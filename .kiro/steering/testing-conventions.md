---
inclusion: fileMatch
fileMatchPattern: "**/*test*.{ts,go}|**/*spec*.ts"
---

# Тестирование - Конвенции

## Frontend (Vitest)

### Структура

```
frontend/tests/
├── unit/           # Изолированные тесты
├── integration/    # Тесты взаимодействия
└── e2e/            # End-to-end
```

### Пример

```ts
describe("useFileSelection", () => {
  it("toggleSelect добавляет файл в selection", () => {
    const { selectedPaths, toggleSelect } = useFileSelection();
    toggleSelect("/path/to/file.ts");
    expect(selectedPaths.value.has("/path/to/file.ts")).toBe(true);
  });
});
```

### Моки

```ts
vi.mock("@/wailsjs/go/main/App", () => ({
  GetFileTree: vi.fn().mockResolvedValue([]),
}));
```

### Запуск

```bash
npm run test:run
npm run test:run -- --coverage
```

## Backend (Go)

### Структура

- Unit: `*_test.go` рядом с кодом
- Integration: `tests/integration/`
- Моки: `testutils/mocks.go`

### Table-driven tests

```go
func TestSearchFiles(t *testing.T) {
    tests := []struct {
        name     string
        pattern  string
        expected int
        wantErr  bool
    }{
        {"exact match", "main.go", 1, false},
        {"invalid", "[bad", 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := SearchFiles(tt.pattern)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.Len(t, result, tt.expected)
        })
    }
}
```

### Запуск

```bash
go test ./...
go test -v -run TestName ./...
go test -cover ./...
```

## Что тестировать

**Обязательно:** бизнес-логика, utils, API handlers, error handling

**Опционально:** UI компоненты, простые getters

## Покрытие

| Слой     | Lines | Functions |
| -------- | ----- | --------- |
| Frontend | 50%   | 30%       |
| Backend  | 60%   | 40%       |
