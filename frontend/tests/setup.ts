import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, vi } from 'vitest'

// Настройка глобальных моков для Wails API
beforeEach(() => {
  // Создание мока для window.go
  Object.defineProperty(window, 'go', {
    value: {
      main: {
        App: {
          SelectDirectory: vi.fn(),
          GetRecentProjects: vi.fn(),
          AddRecentProject: vi.fn(),
          RemoveRecentProject: vi.fn(),
          GetCurrentDirectory: vi.fn(),
          GetProjectContext: vi.fn(),
          BuildProject: vi.fn(),
          RunTests: vi.fn(),
          AnalyzeProject: vi.fn(),
          ExportProject: vi.fn()
        }
      }
    },
    writable: true
  })

  // Мок для runtime.OpenDirectoryDialog если используется напрямую
  if (!window.go.main.App.SelectDirectory) {
    window.go.main.App.SelectDirectory = vi.fn()
  }
})

afterEach(() => {
  // Очистка моков после каждого теста
  if (window.go?.main?.App) {
    Object.keys(window.go.main.App).forEach(key => {
      if (typeof window.go.main.App[key] === 'function' && window.go.main.App[key]._isMockFunction) {
        window.go.main.App[key].mockClear()
      }
    })
  }
})

// Глобальные утилиты для тестирования
export const createWrapper = (component: any, options?: any) => {
  return mount(component, options)
}

export const mockApiService = {
  selectDirectory: vi.fn(),
  getRecentProjects: vi.fn(),
  addRecentProject: vi.fn(),
  removeRecentProject: vi.fn(),
  getCurrentDirectory: vi.fn(),
  getProjectContext: vi.fn(),
  buildProject: vi.fn(),
  runTests: vi.fn(),
  analyzeProject: vi.fn(),
  exportProject: vi.fn()
}

// Настройка jsdom environment
global.ResizeObserver = require('resize-observer-polyfill')
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})