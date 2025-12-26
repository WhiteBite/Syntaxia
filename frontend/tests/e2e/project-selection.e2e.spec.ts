import ProjectSelector from '@/components/ProjectSelector.vue'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

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

// Мок для path.basename
const path = {
  basename: (path: string) => {
    const parts = path.replace(/\/$/, '').split('/')
    return parts[parts.length - 1]
  }
}

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

describe('project-selection.e2e.spec.ts', () => {
  let projectStore: ReturnType<typeof useProjectStore>

  beforeEach(() => {
    // Очистка localStorage перед каждым тестом
    localStorage.clear()

    // Установка Pinia для тестов
    setActivePinia(createPinia())
    projectStore = useProjectStore()

    // Очищаем store после инициализации
    projectStore.clearProject()
    projectStore.recentProjects = []

    // Очистка моков перед каждым тестом
    vi.clearAllMocks()

    // Reset pathExists mock to return true by default
    vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)
  })

  describe('Полный цикл работы с проектом', () => {
    it('Начальное состояние: отображение ProjectSelector с пустым списком recent projects', () => {
      const wrapper = mount(ProjectSelector)

      // Проверяем отображение заголовка
      expect(wrapper.text()).toContain('Syntaxia')

      // Проверяем, что список recent projects пуст
      expect(wrapper.findAll('.project-card')).toHaveLength(0)
    })

    it('Выбор первого проекта: клик по кнопке, мок выбора директории, отображение MainWorkspace', async () => {
      const selectedPath = '/e2e/test/project'

      // Мокируем выбор директории
      vi.spyOn(apiService, 'selectDirectory').mockResolvedValue(selectedPath)
      vi.spyOn(apiService, 'addRecentProject').mockResolvedValue()
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const wrapper = mount(ProjectSelector)

      // Клик по кнопке
      await wrapper.find('.cta-button').trigger('click')
      await new Promise(resolve => setTimeout(resolve, 50))

      // Проверяем, что вызваны необходимые методы
      expect(apiService.selectDirectory).toHaveBeenCalled()
      expect(projectStore.currentPath).toBe(selectedPath)
      expect(projectStore.currentName).toBe('project')

      // Проверяем, что эмитится событие 'opened'
      expect(wrapper.emitted('opened')).toBeTruthy()
      expect(wrapper.emitted('opened')![0]).toEqual([selectedPath])

      // Проверяем, что проект добавлен в recent
      expect(projectStore.recentProjects).toContainEqual({
        path: selectedPath,
        name: 'project',
        lastOpened: expect.any(Number)
      })
    })

    it('Закрытие и повторное открытие: возврат к ProjectSelector и открытие из recent', async () => {
      const projectPath = '/e2e/return/project'

      // Мокируем методы API
      vi.spyOn(apiService, 'selectDirectory').mockResolvedValue(projectPath)
      vi.spyOn(apiService, 'addRecentProject').mockResolvedValue()
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      // Сначала открываем проект
      const wrapper = mount(ProjectSelector)
      await wrapper.find('.cta-button').trigger('click')
      await new Promise(resolve => setTimeout(resolve, 50))

      // Проверяем, что проект открыт
      expect(projectStore.currentPath).toBe(projectPath)
      expect(projectStore.recentProjects).toContainEqual({
        path: projectPath,
        name: 'project',
        lastOpened: expect.any(Number)
      })

      // Симулируем закрытие проекта (через store)
      projectStore.clearProject()
      expect(projectStore.currentPath).toBeNull()

      // Проверяем, что проект в списке recent (может быть больше 1 из-за предыдущих тестов)
      expect(projectStore.recentProjects.length).toBeGreaterThanOrEqual(1)
      expect(projectStore.recentProjects.some(p => p.path === projectPath)).toBe(true)

      // Снова монтируем ProjectSelector для проверки списка recent
      const selectorWithRecent = mount(ProjectSelector)

      // Ждем обновления DOM
      await selectorWithRecent.vm.$nextTick()

      // Проверяем, что проект отображается в списке recent
      const recentItems = selectorWithRecent.findAll('.project-card')
      expect(recentItems.length).toBeGreaterThanOrEqual(1)
      expect(selectorWithRecent.text()).toContain('project')
      expect(selectorWithRecent.text()).toContain(projectPath)

      // Клик по первому recent проекту (который должен быть нашим)
      await recentItems[0].trigger('click')
      await new Promise(resolve => setTimeout(resolve, 50))

      // Проверяем, что проект снова открыт
      expect(projectStore.currentPath).toBe(projectPath)
      expect(projectStore.currentName).toBe('project')
    })

    it('Добавление нескольких проектов: проверка списка и порядка', async () => {
      // Очищаем store перед тестом
      projectStore.recentProjects = []

      const projects = [
        { path: '/multi/1', name: 'multi-project-1' },
        { path: '/multi/2', name: 'multi-project-2' },
        { path: '/multi/3', name: 'multi-project-3' }
      ]

      // Мокируем API вызовы
      const selectDirectoryMock = vi.spyOn(apiService, 'selectDirectory')
      vi.spyOn(apiService, 'addRecentProject').mockResolvedValue()
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const wrapper = mount(ProjectSelector)

      // Добавляем несколько проектов
      for (const project of projects) {
        selectDirectoryMock.mockResolvedValue(project.path)

        await wrapper.find('.cta-button').trigger('click')
        await new Promise(resolve => setTimeout(resolve, 50))
      }

      // Проверяем, что все проекты добавлены
      expect(projectStore.recentProjects).toHaveLength(3)

      // Проверяем порядок (последний открытый должен быть первым)
      const recentProjects = projectStore.recentProjects
      expect(recentProjects[0]).toEqual({
        path: '/multi/3',
        name: '3',
        lastOpened: expect.any(Number)
      })
      expect(recentProjects[1]).toEqual({
        path: '/multi/2',
        name: '2',
        lastOpened: expect.any(Number)
      })
      expect(recentProjects[2]).toEqual({
        path: '/multi/1',
        name: '1',
        lastOpened: expect.any(Number)
      })
    })

    it('Auto-open функционал: автоматическое открытие последнего проекта', async () => {
      const lastProject = { path: '/auto/open/last', name: 'auto-open-project', lastOpened: Date.now() }

      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      // Добавляем проект в список recent через публичный API
      await projectStore.openProjectByPath(lastProject.path)
      await new Promise(resolve => setTimeout(resolve, 50))

      // Проверяем, что проект открыт
      expect(projectStore.currentPath).toBe(lastProject.path)
      expect(projectStore.currentName).toBe('last')
    })

    it('Drag & Drop: симуляция drag & drop директории', async () => {
      // Мокируем API
      vi.spyOn(apiService, 'addRecentProject').mockResolvedValue()

      const wrapper = mount(ProjectSelector)

      // Симулируем dragover
      await wrapper.trigger('dragover')
      expect(wrapper.classes()).toContain('drag-over')

      // Create a mock FileSystemDirectoryEntry
      const mockDirEntry = {
        isDirectory: true,
        isFile: false,
        name: 'test-folder',
        fullPath: '/dropped/project/path',
        getAsEntry: () => mockDirEntry
      }

      // Create a DataTransferItem with webkitGetAsEntry
      const dataTransferItem = {
        kind: 'file',
        type: 'folder',
        webkitGetAsEntry: () => mockDirEntry
      }

      const dataTransfer = {
        files: [new File([], 'test-folder')],
        items: [dataTransferItem]
      } as any

      await wrapper.trigger('drop', { dataTransfer })

      // Проверяем, что overlay скрылся
      expect(wrapper.classes()).not.toContain('drag-over')
    })
  })
})
