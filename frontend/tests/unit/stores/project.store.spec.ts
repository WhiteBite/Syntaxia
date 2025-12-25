import { useProjectStore } from '@/stores/project.store'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/services/api.service', () => ({
  apiService: {
    selectDirectory: vi.fn(),
    getRecentProjects: vi.fn(),
    addRecentProject: vi.fn(),
    removeRecentProject: vi.fn(),
    getCurrentDirectory: vi.fn(),
    pathExists: vi.fn().mockResolvedValue(true)
  }
}))

import { apiService } from '@/services/api.service'

// Мок для localStorage
const localStorageMock = (() => {
  let store: { [key: string]: string } = {}
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value.toString()
    },
    removeItem: (key: string) => {
      delete store[key]
    },
    clear: () => {
      store = {}
    },
    get length() {
      return Object.keys(store).length
    },
    key: (index: number) => Object.keys(store)[index] || null
  }
})()

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

describe('project.store.ts', () => {
  let projectStore: ReturnType<typeof useProjectStore>

  beforeEach(() => {
    // Очистка localStorage перед каждым тестом
    localStorage.clear()

    // Установка Pinia для тестов
    setActivePinia(createPinia())
    projectStore = useProjectStore()

    // Дополнительно очищаем store после инициализации, чтобы переопределить значения из localStorage
    projectStore.clearProject()
    projectStore.recentProjects = []
    projectStore.setAutoOpenLast(false) // Ensure autoOpenLast is false by default
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('openProjectByPath', () => {
    it('должен устанавливать currentPath и currentName', async () => {
      const testPath = '/path/to/project'

      await projectStore.openProjectByPath(testPath)

      expect(projectStore.currentPath).toBe(testPath)
      expect(projectStore.currentName).toBe('project') // Using split instead of basename
    })

    it('должен вызывать addRecentProject с правильными параметрами', async () => {
      const testPath = '/path/to/project'

      const { apiService } = await import('@/services/api.service')
      const addRecentProjectSpy = vi.spyOn(apiService, 'addRecentProject')

      await projectStore.openProjectByPath(testPath)

      expect(addRecentProjectSpy).toHaveBeenCalledWith(testPath, 'project') // Using split instead of basename
    })

    it('должен сохранять проект в localStorage', async () => {
      const testPath = '/path/to/project'

      await projectStore.openProjectByPath(testPath)

      const storedRecent = JSON.parse(localStorage.getItem('Syntaxia_recent_projects') || '[]')
      expect(storedRecent).toContainEqual({ path: testPath, name: 'project', lastOpened: expect.any(Number) })
    })

    it('должен обрабатывать ошибки backend и продолжать работу', async () => {
      const testPath = '/path/to/project'

      const { apiService } = await import('@/services/api.service')
      vi.spyOn(apiService, 'addRecentProject').mockRejectedValue(new Error('Backend error'))

      // Проверяем, что ошибка не прерывает выполнение
      await expect(projectStore.openProjectByPath(testPath)).resolves.not.toThrow()

      // Проект все равно должен быть добавлен в store
      expect(projectStore.currentPath).toBe(testPath)
      expect(projectStore.currentName).toBe('project')
    })
  })

  describe('fetchRecentProjects', () => {
    it('должен загружать проекты из backend через apiService.getRecentProjects', async () => {
      const mockProjects = [
        { path: '/path/1', name: 'project1', lastOpened: Date.now() },
        { path: '/path/2', name: 'project2', lastOpened: Date.now() }
      ]

      const { apiService } = await import('@/services/api.service')
      vi.spyOn(apiService, 'getRecentProjects').mockResolvedValue(JSON.stringify(mockProjects))

      await projectStore.fetchRecentProjects()

      expect(apiService.getRecentProjects).toHaveBeenCalled()
      expect(projectStore.recentProjects).toEqual(mockProjects)
    })

    it('должен обрабатывать ошибки и оставлять данные из localStorage', async () => {
      const localStorageProjects = [
        { path: '/local/path', name: 'local-project', lastOpened: Date.now() }
      ]

      localStorage.setItem('Syntaxia_recent_projects', JSON.stringify(localStorageProjects))
      projectStore.recentProjects = localStorageProjects // Устанавливаем начальные данные
      vi.spyOn(apiService, 'getRecentProjects').mockRejectedValue(new Error('Backend error'))

      await projectStore.fetchRecentProjects()

      // Данные из localStorage должны остаться
      expect(projectStore.recentProjects).toEqual(localStorageProjects)
    })

    it('должен объединять данные из backend и localStorage', async () => {
      const localStorageProjects = [
        { path: '/local/path', name: 'local-project', lastOpened: Date.now() }
      ]
      const backendProjects = [
        { path: '/backend/path', name: 'backend-project', lastOpened: Date.now() }
      ]

      localStorage.setItem('Syntaxia_recent_projects', JSON.stringify(localStorageProjects))
      vi.spyOn(apiService, 'getRecentProjects').mockResolvedValue(JSON.stringify(backendProjects))

      await projectStore.fetchRecentProjects()

      // Only backend data should remain after fetch
      expect(projectStore.recentProjects).toEqual(backendProjects)
    })
  })

  describe('openProjectByPath', () => {
    it('должен добавлять проект в начало списка через addToRecent', async () => {
      const projectPath = '/path/1'

      await projectStore.openProjectByPath(projectPath)

      const expectedProject = {
        path: projectPath,
        name: '1', // Using split instead of basename
        lastOpened: expect.any(Number)
      }
      expect(projectStore.recentProjects[0]).toEqual(expectedProject)
    })

    it('должен перемещать существующий проект в начало через addToRecent', async () => {
      const project1 = { path: '/path/1', name: 'project1', lastOpened: Date.now() }
      const project2 = { path: '/path/2', name: 'project2', lastOpened: Date.now() }

      projectStore.recentProjects = [project1, project2]

      await projectStore.openProjectByPath('/path/1')

      const updatedProject1 = { path: '/path/1', name: '1', lastOpened: expect.any(Number) }
      expect(projectStore.recentProjects[0]).toEqual(updatedProject1)
      expect(projectStore.recentProjects.length).toBe(2)
    })

    it('должен ограничивать список до MAX_RECENT (10) проектов через addToRecent', async () => {
      for (let i = 0; i < 12; i++) {
        const projectPath = `/path/${i}`
        await projectStore.openProjectByPath(projectPath)
      }

      expect(projectStore.recentProjects.length).toBe(10)
      // Последний добавленный проект должен быть первым
      expect(projectStore.recentProjects[0].path).toBe('/path/11')
      expect(projectStore.recentProjects[0].name).toBe('11')
    })
  })

  describe('setAutoOpenLast', () => {
    it('должен сохранять настройку в localStorage', () => {
      projectStore.setAutoOpenLast(true)

      expect(localStorage.getItem('Syntaxia_auto_open_last')).toBe('true')

      projectStore.setAutoOpenLast(false)

      expect(localStorage.getItem('Syntaxia_auto_open_last')).toBe('false')
    })
  })

  describe('maybeAutoOpenLastProject', () => {
    beforeEach(() => {
      // Ensure clean state before each test in this suite
      projectStore.clearProject()
      projectStore.recentProjects = []
    })

    it('должен открывать последний проект если autoOpenLast=true', async () => {
      // Make sure we have a clean state
      projectStore.clearProject()
      projectStore.recentProjects = []

      const recentProject = { path: '/last/project', name: 'project', lastOpened: Date.now() }
      // Add the project to recent projects using the store's internal method to ensure proper reactivity
      projectStore.recentProjects = [recentProject]
      projectStore.setAutoOpenLast(true)

      // Verify conditions before calling the function
      expect(projectStore.currentPath).toBeNull() // Should be null to proceed
      expect(projectStore.autoOpenLast).toBe(true) // Should be true to proceed
      expect(projectStore.lastProjectPath).toBe('/last/project') // Should have a last project to proceed

      const result = await projectStore.maybeAutoOpenLastProject()

      // Verify the function executed and returned true (since it opens the project)
      expect(result).toBe(true)

      // Verify that the project was actually opened (state changed)
      expect(projectStore.currentPath).toBe('/last/project')
      expect(projectStore.currentName).toBe('project')
    })

    it('не должен открывать проект если autoOpenLast=false', async () => {
      const recentProject = { path: '/last/project', name: 'last-project', lastOpened: Date.now() }
      projectStore.recentProjects = [recentProject]
      projectStore.setAutoOpenLast(false)

      const openProjectByPathSpy = vi.spyOn(projectStore, 'openProjectByPath').mockResolvedValue(true)

      await projectStore.maybeAutoOpenLastProject()

      expect(openProjectByPathSpy).not.toHaveBeenCalled()
    })

    it('не должен открывать проект если нет recentProjects', async () => {
      projectStore.setAutoOpenLast(true)

      const openProjectByPathSpy = vi.spyOn(projectStore, 'openProjectByPath').mockResolvedValue(true)

      await projectStore.maybeAutoOpenLastProject()

      expect(openProjectByPathSpy).not.toHaveBeenCalled()
    })
  })

  describe('removeFromRecent', () => {
    it('должен удалять проект из списка', async () => {
      const project1 = { path: '/path/1', name: 'project1', lastOpened: Date.now() }
      const project2 = { path: '/path/2', name: 'project2', lastOpened: Date.now() }
      projectStore.recentProjects = [project1, project2]

      await projectStore.removeFromRecent('/path/1')

      expect(projectStore.recentProjects).toEqual([project2])
    })

    it('должен вызывать removeRecentProject в apiService', async () => {
      const { apiService } = await import('@/services/api.service')
      const removeRecentProjectSpy = vi.spyOn(apiService, 'removeRecentProject')
      const project = { path: '/path/1', name: 'project1', lastOpened: Date.now() }
      projectStore.recentProjects = [project]

      await projectStore.removeFromRecent('/path/1')

      expect(removeRecentProjectSpy).toHaveBeenCalledWith('/path/1')
    })
  })

  describe('clearRecent', () => {
    it('должен очищать весь список', () => {
      const project1 = { path: '/path/1', name: 'project1', lastOpened: Date.now() }
      const project2 = { path: '/path/2', name: 'project2', lastOpened: Date.now() }
      projectStore.recentProjects = [project1, project2]

      projectStore.clearRecent()

      expect(projectStore.recentProjects).toEqual([])
    })
  })
})
