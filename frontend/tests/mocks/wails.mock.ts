/**
 * Моки для Wails API
 * Используется для тестирования компонентов, взаимодействующих с Wails backend
 */

// Типы для Wails API
interface ProjectInfo {
  path: string
  name: string
  lastOpened: string
}

interface WailsAppAPI {
  SelectDirectory(): Promise<string>
  GetRecentProjects(): Promise<ProjectInfo[]>
  AddRecentProject(path: string, name: string): Promise<void>
  RemoveRecentProject(path: string): Promise<void>
  GetCurrentDirectory(): Promise<string>
  GetProjectContext(path: string): Promise<any>
  BuildProject(path: string): Promise<any>
  RunTests(path: string): Promise<any>
  AnalyzeProject(path: string): Promise<any>
  ExportProject(path: string): Promise<any>
}

// Внутреннее хранилище для моков
let mockStorage: {
  recentProjects: ProjectInfo[]
  currentDirectory: string
  selectedDirectory: string | null
} = {
  recentProjects: [],
  currentDirectory: '/home/user',
  selectedDirectory: null
}

// Мок объекта Wails API
const wailsAppMock: WailsAppAPI = {
  SelectDirectory: async (): Promise<string> => {
    return mockStorage.selectedDirectory || '/mocked/project/path'
  },
  
  GetRecentProjects: async (): Promise<ProjectInfo[]> => {
    return [...mockStorage.recentProjects]
  },
  
  AddRecentProject: async (path: string, name: string): Promise<void> => {
    // Проверяем, есть ли уже проект в списке
    const existingIndex = mockStorage.recentProjects.findIndex(p => p.path === path)
    if (existingIndex !== -1) {
      // Если есть, обновляем lastOpened и перемещаем в начало
      const project = { ...mockStorage.recentProjects[existingIndex], lastOpened: new Date().toISOString() }
      mockStorage.recentProjects.splice(existingIndex, 1)
      mockStorage.recentProjects.unshift(project)
    } else {
      // Если нет, добавляем в начало
      mockStorage.recentProjects.unshift({
        path,
        name,
        lastOpened: new Date().toISOString()
      })
    }
    
    // Ограничиваем количество проектов до 10
    if (mockStorage.recentProjects.length > 10) {
      mockStorage.recentProjects = mockStorage.recentProjects.slice(0, 10)
    }
  },
  
  RemoveRecentProject: async (path: string): Promise<void> => {
    mockStorage.recentProjects = mockStorage.recentProjects.filter(p => p.path !== path)
  },
  
  GetCurrentDirectory: async (): Promise<string> => {
    return mockStorage.currentDirectory
  },
  
  GetProjectContext: async (path: string): Promise<any> => {
    return {
      path,
      name: path.split('/').pop() || 'unknown',
      files: [],
      settings: {}
    }
  },
  
  BuildProject: async (path: string): Promise<any> => {
    return {
      success: true,
      output: `Build successful for ${path}`,
      errors: []
    }
  },
  
  RunTests: async (path: string): Promise<any> => {
    return {
      success: true,
      passed: 10,
      failed: 0,
      total: 10
    }
  },
  
  AnalyzeProject: async (path: string): Promise<any> => {
    return {
      issues: [],
      metrics: {
        complexity: 0,
        lines: 0,
        files: 0
      }
    }
  },
  
  ExportProject: async (path: string): Promise<any> => {
    return {
      success: true,
      path: `${path}/export.zip`
    }
  }
}

/**
 * Установка моков Wails API
 */
export const setupWailsMocks = (): void => {
  // Создание мока для window.go
  if (typeof window !== 'undefined') {
    (window as any).go = {
      main: {
        App: { ...wailsAppMock }
      }
    }
  }
}

/**
 * Сброс моков Wails API к начальному состоянию
 */
export const resetWailsMocks = (): void => {
  mockStorage = {
    recentProjects: [],
    currentDirectory: '/home/user',
    selectedDirectory: null
  }
  
  // Обновляем моки с новым состоянием
  if (typeof window !== 'undefined') {
    (window as any).go.main.App = { ...wailsAppMock }
  }
}

/**
 * Установка тестовых данных для моков
 */
export const setWailsMockData = (data: Partial<typeof mockStorage>): void => {
  mockStorage = {
    ...mockStorage,
    ...data
  }
  
  // Обновляем моки с новыми данными
  if (typeof window !== 'undefined') {
    (window as any).go.main.App = { ...wailsAppMock }
  }
}

/**
 * Получение текущего состояния моков
 */
export const getWailsMockState = (): typeof mockStorage => {
  return { ...mockStorage }
}

/**
 * Установка конкретного ответа для SelectDirectory
 */
export const mockSelectDirectory = (path: string | null): void => {
  mockStorage.selectedDirectory = path
}

// Экспортируем моки по умолчанию
export default wailsAppMock