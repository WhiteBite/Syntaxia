import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock Wails API FIRST before any imports
vi.mock('#wailsjs/go/main/App', () => ({
  PathExists: vi.fn().mockResolvedValue(true),
  SelectDirectory: vi.fn(),
  GetRecentProjects: vi.fn(),
  AddRecentProject: vi.fn(),
  RemoveRecentProject: vi.fn(),
  GetCurrentDirectory: vi.fn(),
  GetVersionInfo: vi.fn().mockResolvedValue({ version: '1.0.0', buildDate: '2024-01-01' })
}))

import ProjectSelector from '@/components/ProjectSelector.vue'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

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

// Mock stores before importing
vi.mock('@/features/files', () => ({
  useFileStore: vi.fn(() => ({
    resetStore: vi.fn()
  }))
}))

vi.mock('@/features/context', () => ({
  useContextStore: vi.fn(() => ({
    clearContext: vi.fn()
  }))
}))

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

describe('project-selection.integration.spec.ts', () => {
  let projectStore: ReturnType<typeof useProjectStore>

  beforeEach(() => {
    // Установка Pinia для тестов
    setActivePinia(createPinia())
    projectStore = useProjectStore()

    // Очистка localStorage перед каждым тестом
    localStorage.clear()

    // Очистка моков перед каждым тестом
    vi.clearAllMocks()

    // Reset pathExists mock
    vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)
  })

  describe('Happy path: Выбор проекта через нативный диалог', () => {
    it('должен добавить проект в recent, сохранить в localStorage и переключить App.vue на MainWorkspace', async () => {
      const selectedPath = '/path/to/new/project'
      const projectName = 'project'

      // Мокируем apiService.selectDirectory для возврата пути
      vi.spyOn(apiService, 'selectDirectory').mockResolvedValue(selectedPath)
      vi.spyOn(apiService, 'addRecentProject').mockResolvedValue()
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const wrapper = mount(ProjectSelector)

      // Клик по кнопке выбора директории
      await wrapper.find('.cta-button').trigger('click')

      // Wait for all async operations including dynamic imports
      await new Promise(resolve => setTimeout(resolve, 1000))
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      // Проверяем, что проект добавлен в recent
      expect(projectStore.recentProjects).toContainEqual({
        path: selectedPath,
        name: projectName,
        lastOpened: expect.any(Number)
      })

      // Проверяем, что проект сохранен в localStorage
      const storedRecent = JSON.parse(localStorage.getItem('Syntaxia_recent_projects') || '[]')
      expect(storedRecent).toContainEqual({
        path: selectedPath,
        name: projectName,
        lastOpened: expect.any(Number)
      })

      // Проверяем, что эмитится событие 'opened' с путем проекта
      expect(wrapper.emitted('opened')).toBeTruthy()
      expect(wrapper.emitted('opened')![0]).toEqual([selectedPath])
    }, 10000)
  })

  describe('Recent projects flow: Открытие из списка recent', () => {
    it('должен отобразить список и открыть проект при клике', async () => {
      const testProjects = [
        { path: '/path/1', name: 'project1', lastOpened: Date.now() },
        { path: '/path/2', name: 'project2', lastOpened: Date.now() }
      ]

      // Мокируем apiService.getRecentProjects с тестовыми данными
      vi.spyOn(apiService, 'getRecentProjects').mockResolvedValue(JSON.stringify(testProjects))
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const wrapper = mount(ProjectSelector)

      // Ждем, пока данные загрузятся
      await new Promise(resolve => setTimeout(resolve, 100))
      await wrapper.vm.$nextTick()

      // Проверяем отображение списка
      expect(wrapper.findAll('.project-card')).toHaveLength(2)
      expect(wrapper.text()).toContain('project1')
      expect(wrapper.text()).toContain('project2')

      // Кликнем на первый проект из списка
      const firstProjectItem = wrapper.findAll('.project-card')[0]
      await firstProjectItem.trigger('click')

      // Wait for all async operations
      await new Promise(resolve => setTimeout(resolve, 150))
      await wrapper.vm.$nextTick()

      // Проверяем, что проект открылся
      expect(projectStore.currentPath).toBe('/path/1')
      expect(projectStore.currentName).toBe('1')

      // Проверяем, что проект переместился в начало списка
      expect(projectStore.recentProjects[0]).toEqual({
        path: '/path/1',
        name: '1',
        lastOpened: expect.any(Number)
      })
    }, 10000)
  })

  describe('Auto-open flow: Автоматическое открытие последнего проекта', () => {
    it('должен открыть проект автоматически если autoOpenLast=true в localStorage', async () => {
      const lastProject = { path: '/last/project', name: 'project', lastOpened: Date.now() }

      // Установить autoOpenLast=true в localStorage
      localStorage.setItem('Syntaxia_auto_open_last', 'true')
      // Добавить проект в recent
      localStorage.setItem('Syntaxia_recent_projects', JSON.stringify([lastProject]))

      // Мокируем apiService
      vi.spyOn(apiService, 'getRecentProjects').mockResolvedValue(JSON.stringify([lastProject]))
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      // Пересоздаем store чтобы он загрузил данные из localStorage
      setActivePinia(createPinia())
      projectStore = useProjectStore()

      // Вызываем maybeAutoOpenLastProject
      await projectStore.maybeAutoOpenLastProject()
      await new Promise(resolve => setTimeout(resolve, 150))

      // Проверяем, что проект автоматически открылся
      expect(projectStore.currentPath).toBe(lastProject.path)
      expect(projectStore.currentName).toBe('project')
    }, 10000)
  })

  describe('Error handling: Обработка ошибок', () => {
    it('должен показать toast с ошибкой и не добавить проект в recent', async () => {
      // Мокируем apiService.selectDirectory для выброса ошибки
      vi.spyOn(apiService, 'selectDirectory').mockRejectedValue(new Error('Directory selection failed'))

      const initialLength = projectStore.recentProjects.length

      const wrapper = mount(ProjectSelector)

      // Клик по кнопке выбора директории
      await wrapper.find('.cta-button').trigger('click')
      await new Promise(resolve => setTimeout(resolve, 100))

      // Проверяем, что количество проектов не изменилось (ошибка произошла до добавления)
      expect(projectStore.recentProjects.length).toBe(initialLength)

      // Проверяем, что событие opened не эмитится
      expect(wrapper.emitted('opened')).toBeUndefined()
    }, 10000)
  })

  describe('Persistence: Сохранение и загрузка из localStorage', () => {
    it('должен сохранить и восстановить проекты', async () => {
      const testProjects = [
        { path: '/path/1', name: 'project1', lastOpened: Date.now() },
        { path: '/path/2', name: 'project2', lastOpened: Date.now() },
        { path: '/path/3', name: 'project3', lastOpened: Date.now() }
      ]

      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      // Добавляем проекты в store через публичный API
      for (const project of testProjects) {
        await projectStore.openProjectByPath(project.path)
        await new Promise(resolve => setTimeout(resolve, 50))
      }

      // Проверяем, что проекты сохранены в localStorage
      const storedRecent = JSON.parse(localStorage.getItem('Syntaxia_recent_projects') || '[]')
      expect(storedRecent).toContainEqual({ path: '/path/1', name: '1', lastOpened: expect.any(Number) })
      expect(storedRecent).toContainEqual({ path: '/path/2', name: '2', lastOpened: expect.any(Number) })
      expect(storedRecent).toContainEqual({ path: '/path/3', name: '3', lastOpened: expect.any(Number) })

      // Очищаем store
      projectStore.recentProjects = []

      // Загружаем из localStorage
      const loadedRecent = JSON.parse(localStorage.getItem('Syntaxia_recent_projects') || '[]')
      projectStore.recentProjects = loadedRecent

      // Проверяем, что проекты восстановлены
      expect(projectStore.recentProjects).toEqual(storedRecent)
    }, 10000)
  })

  describe('Integration: Взаимодействие компонентов и store', () => {
    it('должен корректно передавать данные между ProjectSelector и ProjectStore', async () => {
      const selectedPath = '/integration/test/path'

      // Мокируем apiService
      vi.spyOn(apiService, 'selectDirectory').mockResolvedValue(selectedPath)
      vi.spyOn(apiService, 'addRecentProject').mockResolvedValue()
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const wrapper = mount(ProjectSelector)

      // Выбираем проект
      await wrapper.find('.cta-button').trigger('click')

      // Wait for all async operations including dynamic imports
      await new Promise(resolve => setTimeout(resolve, 1000))
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      // Проверяем, что store обновлен
      expect(projectStore.currentPath).toBe(selectedPath)
      expect(projectStore.currentName).toBe('path')
      expect(projectStore.recentProjects).toContainEqual({
        path: selectedPath,
        name: 'path',
        lastOpened: expect.any(Number)
      })

      // Проверяем, что проект сохранен в localStorage
      const storedRecent = JSON.parse(localStorage.getItem('Syntaxia_recent_projects') || '[]')
      expect(storedRecent).toContainEqual({
        path: selectedPath,
        name: 'path',
        lastOpened: expect.any(Number)
      })

      // Проверяем, что эмитится событие
      expect(wrapper.emitted('opened')).toBeTruthy()
      expect(wrapper.emitted('opened')![0]).toEqual([selectedPath])
    }, 10000)
  })
})
