# Тестирование

Документация по тестовой инфраструктуре проекта Syntaxia.

## Обзор тестовой инфраструктуры

Проект использует следующие инструменты для тестирования:

- **Vitest** - фреймворк для тестирования
- **@vue/test-utils** - утилиты для тестирования Vue компонентов
- **jsdom** - среда выполнения тестов в Node.js
- **Pinia** - для тестирования stores
- **TypeScript** - с типобезопасностью в тестах

Структура директорий тестов:

```
frontend/tests/
├── setup.ts              # Файл настройки для Vitest
├── README.md             # Документация по тестированию (этот файл)
├── mocks/                # Моки для тестирования
│   └── wails.mock.ts     # Моки для Wails API
├── unit/                 # Unit тесты
│   ├── stores/           # Тесты для stores
│   │   └── project.store.spec.ts
│   └── components/       # Тесты для компонентов
│       └── ProjectSelector.spec.ts
├── integration/          # Интеграционные тесты
│   └── project-selection.spec.ts
└── e2e/                  # E2E тесты
    └── project-selection.e2e.spec.ts
```

### Соглашения по именованию

- Файлы тестов должны заканчиваться на `.spec.ts` или `.test.ts`
- Unit тесты для компонентов: `ComponentName.spec.ts`
- Unit тесты для stores: `store-name.spec.ts`
- Интеграционные тесты: `integration-test-name.spec.ts`
- E2E тесты: `feature-name.e2e.spec.ts`

## Запуск тестов

### Основные команды

```bash
# Запуск в watch mode
npm run test

# Однократный запуск
npm run test:run

# UI для просмотра результатов
npm run test:ui

# Запуск с coverage
npm run test:coverage
```

### Запуск конкретного теста

```bash
# Запуск конкретного файла
npm run test tests/unit/stores/project.store.spec.ts

# Запуск тестов с определенным названием
npm run test -t "openProjectByPath"
```

## Типы тестов

### Unit тесты

Unit тесты тестируют отдельные функции, компоненты и stores изолированно:

- **Stores**: Тестирование действий, геттеров и состояния
- **Composables**: Тестирование пользовательских хуков
- **Utils**: Тестирование вспомогательных функций
- **Components (изолированно)**: Тестирование компонентов без зависимостей

### Component тесты

Тестирование Vue компонентов с фокусом на UI поведение:

- Рендеринг
- Обработка событий
- Взаимодействие с пользователем
- Мокирование зависимостей

### Integration тесты

Тестирование взаимодействия между несколькими модулями:

- Взаимодействие компонентов и stores
- Взаимодействие с API сервисами
- Проверка полных пользовательских сценариев

### E2E тесты

Тестирование полных пользовательских сценариев:

- Тестирование от начала до конца
- Имитация действий пользователя
- Проверка состояния приложения

## Моки и утилиты

### wails.mock.ts

Файл содержит моки для Wails API, включая:

- `setupWailsMocks()` - установка моков
- `resetWailsMocks()` - сброс моков
- `setWailsMockData()` - установка тестовых данных
- `getWailsMockState()` - получение текущего состояния
- `mockSelectDirectory()` - мок для выбора директории

### Использование моков

```typescript
import { setupWailsMocks, resetWailsMocks } from '@/tests/mocks/wails.mock'

beforeEach(() => {
  setupWailsMocks()
})

afterEach(() => {
  resetWailsMocks()
})
```

## Примеры

### Пример unit теста для store

```typescript
import { describe, it, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useProjectStore } from '@/stores/project.store'

describe('project.store.ts', () => {
 let projectStore

  beforeEach(() => {
    setActivePinia(createPinia())
    projectStore = useProjectStore()
  })

  it('должен открыть проект по пути', async () => {
    await projectStore.openProjectByPath('/test/path')
    
    expect(projectStore.currentPath).toBe('/test/path')
  })
})
```

### Пример component теста

```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ProjectSelector from '@/components/ProjectSelector.vue'

describe('ProjectSelector.vue', () => {
  it('должен отображать заголовок', () => {
    const wrapper = mount(ProjectSelector)
    
    expect(wrapper.text()).toContain('Select a Project')
  })
})
```

### Пример integration теста

```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ProjectSelector from '@/components/ProjectSelector.vue'
import { useProjectStore } from '@/stores/project.store'

describe('project-selection.integration.spec.ts', () => {
  it('должен выбрать проект и обновить store', async () => {
    setActivePinia(createPinia())
    const projectStore = useProjectStore()
    const wrapper = mount(ProjectSelector)
    
    // Симуляция выбора проекта
    await wrapper.find('button').trigger('click')
    
    expect(projectStore.hasProject).toBe(true)
  })
})
```

## Best Practices

### Изоляция тестов

Каждый тест должен быть независимым:

```typescript
beforeEach(() => {
  // Подготовка тестового окружения
  setActivePinia(createPinia())
  localStorage.clear()
})

afterEach(() => {
  // Очистка после теста
  vi.clearAllMocks()
})
```

### Именование тестов

Тесты должны быть понятными и описательными:

```typescript
// Хорошо
it('должен устанавливать currentPath и currentName при открытии проекта', async () => {
  // ...
})

// Плохо
it('проверяет функцию', async () => {
  // ...
})
```

### Организация test suites

Группируйте тесты по функциональности:

```typescript
describe('openProjectByPath', () => {
  it('должен устанавливать currentPath', () => {
    // ...
  })
  
  it('должен вызывать addRecentProject', () => {
    // ...
  })
})
```

## Troubleshooting

### Частые проблемы и их решения

1. **Тесты падают из-за асинхронности**

   Используйте `await wrapper.vm.$nextTick()` или `await new Promise(resolve => setTimeout(resolve, 0))` для ожидания обновления DOM.

2. **Проблемы с моками Wails API**

   Убедитесь, что `setupWailsMocks()` вызывается в `beforeEach()`.

3. **Проблемы с Pinia store**

   Используйте `setActivePinia(createPinia())` в `beforeEach()` для изоляции store между тестами.

4. **localStorage не очищается между тестами**

   Добавьте `localStorage.clear()` в `beforeEach()` или `afterEach()`.

### Debugging тестов

Для отладки тестов можно использовать:

```typescript
// Включить логирование
console.log('Debug info:', variable)

// Использовать debugger
debugger

// Запустить тест в debug режиме
npm run test:debug
```

## CI/CD интеграция

Тесты запускаются в GitHub Actions при каждом push/PR:

- Все unit тесты
- Интеграционные тесты
- Проверка покрытия кода

Минимальное покрытие кода: 70%

Покрытие отображается в виде отчета в GitHub после каждого запуска.