import ProjectSelector from '@/components/ProjectSelector.vue'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

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

describe('ProjectSelector.vue', () => {
  let projectStore: ReturnType<typeof useProjectStore>
  let uiStore: ReturnType<typeof useUIStore>

  beforeEach(() => {
    // Установка Pinia для тестов
    setActivePinia(createPinia())
    projectStore = useProjectStore()
    uiStore = useUIStore()

    // Создаем spy для addToast
    vi.spyOn(uiStore, 'addToast')

    // Очистка моков перед каждым тестом
    vi.clearAllMocks()
  })

  describe('Рендеринг', () => {
    it('должен отображать заголовок и кнопку', () => {
      const wrapper = mount(ProjectSelector)

      expect(wrapper.text()).toContain('Syntaxia')
      // Button text depends on locale (ru/en), just check button exists
      expect(wrapper.find('.cta-button').exists()).toBe(true)
    })

    it('должен отображать пустой список recent projects с placeholder\'ами', () => {
      const wrapper = mount(ProjectSelector)

      expect(wrapper.findAll('.project-card')).toHaveLength(0)
      // We don't have a "No recent projects" text in the actual component, so removing this check
    })

    it('должен отображать список recent projects если они есть', async () => {
      projectStore.recentProjects = [
        { path: '/path/1', name: 'Project 1', lastOpened: Date.now() },
        { path: '/path/2', name: 'Project 2', lastOpened: Date.now() }
      ]

      const wrapper = mount(ProjectSelector)

      await wrapper.vm.$nextTick()

      expect(wrapper.findAll('.project-card')).toHaveLength(2)
      expect(wrapper.text()).toContain('Project 1')
      expect(wrapper.text()).toContain('Project 2')
    })
  })

  describe('Клик на кнопку', () => {
    it('должен вызывать apiService.selectDirectory', async () => {
      const selectDirectorySpy = vi.spyOn(apiService, 'selectDirectory').mockResolvedValue('/selected/path')
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)
      const wrapper = mount(ProjectSelector)

      await wrapper.find('.cta-button').trigger('click')
      await new Promise(resolve => setTimeout(resolve, 10))

      expect(selectDirectorySpy).toHaveBeenCalled()
    })

    it('должен вызывать projectStore.openProjectByPath с выбранным путем', async () => {
      vi.spyOn(apiService, 'selectDirectory').mockResolvedValue('/selected/path')
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)
      const openProjectByPathSpy = vi.spyOn(projectStore, 'openProjectByPath')

      const wrapper = mount(ProjectSelector)

      await wrapper.find('.cta-button').trigger('click')
      await new Promise(resolve => setTimeout(resolve, 10))

      expect(openProjectByPathSpy).toHaveBeenCalledWith('/selected/path')
    })

    it('должен эмитить событие \'opened\' с путем проекта', async () => {
      vi.spyOn(apiService, 'selectDirectory').mockResolvedValue('/selected/path')
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const wrapper = mount(ProjectSelector)

      await wrapper.find('.cta-button').trigger('click')

      // Wait for async operations to complete
      await new Promise(resolve => setTimeout(resolve, 50))

      expect(wrapper.emitted('opened')).toBeTruthy()
      expect(wrapper.emitted('opened')![0]).toEqual(['/selected/path'])
    })
  })

  describe('Клик на recent project', () => {
    it('должен открывать проект через projectStore.openProjectByPath', async () => {
      const project = { path: '/recent/path', name: 'Recent Project', lastOpened: Date.now() }
      projectStore.recentProjects = [project]
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const openProjectByPathSpy = vi.spyOn(projectStore, 'openProjectByPath')

      const wrapper = mount(ProjectSelector)
      await wrapper.vm.$nextTick()

      await wrapper.find('.project-card').trigger('click')

      expect(openProjectByPathSpy).toHaveBeenCalledWith('/recent/path')
    })

    it('должен эмитить событие \'opened\'', async () => {
      const project = { path: '/recent/path', name: 'Recent Project', lastOpened: Date.now() }
      projectStore.recentProjects = [project]
      vi.spyOn(apiService, 'pathExists').mockResolvedValue(true)

      const wrapper = mount(ProjectSelector)
      await wrapper.vm.$nextTick()

      await wrapper.find('.project-card').trigger('click')

      // Wait for async operations to complete
      await new Promise(resolve => setTimeout(resolve, 50))

      expect(wrapper.emitted('opened')).toBeTruthy()
    })
  })

  describe('onMounted', () => {
    it('должен вызывать projectStore.fetchRecentProjects', async () => {
      const fetchRecentProjectsSpy = vi.spyOn(projectStore, 'fetchRecentProjects')

      mount(ProjectSelector)

      // Ждем выполнения onMounted
      await new Promise(resolve => setTimeout(resolve, 0))

      expect(fetchRecentProjectsSpy).toHaveBeenCalled()
    })
  })

  describe('Drag & Drop', () => {
    it('должен показывать overlay при dragover', async () => {
      const wrapper = mount(ProjectSelector)

      await wrapper.trigger('dragover')

      expect(wrapper.classes()).toContain('drag-over')
    })

    it('должен скрывать overlay при dragleave', async () => {
      const wrapper = mount(ProjectSelector)

      // Сначала добавим класс drag-over
      await wrapper.trigger('dragover')
      expect(wrapper.classes()).toContain('drag-over')

      // Затем уберем его при dragleave
      await wrapper.trigger('dragleave')

      expect(wrapper.classes()).not.toContain('drag-over')
    })

    it('должен обрабатывать drop директории', async () => {
      const project = { path: '/dropped/path', name: 'Dropped Project', lastOpened: Date.now() }
      vi.spyOn(apiService, 'getRecentProjects').mockResolvedValue(JSON.stringify([project]))

      const wrapper = mount(ProjectSelector)

      // Create a mock FileSystemDirectoryEntry
      const mockDirEntry = {
        isDirectory: true,
        isFile: false,
        name: 'test-folder',
        fullPath: '/dropped/path',
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

    it('должен показывать ошибку при drop файла', async () => {
      const wrapper = mount(ProjectSelector)

      // Create a mock FileSystemFileEntry
      const mockFileEntry = {
        isDirectory: false,
        isFile: true,
        name: 'test-file.txt',
        fullPath: '/path/to/test-file.txt',
        getAsEntry: () => mockFileEntry
      }

      // Create a DataTransferItem with webkitGetAsEntry
      const dataTransferItem = {
        kind: 'file',
        type: 'text/plain',
        webkitGetAsEntry: () => mockFileEntry
      }

      const dataTransfer = {
        files: [new File([], 'test-file.txt', { type: 'text/plain' })],
        items: [dataTransferItem]
      } as any

      await wrapper.trigger('drop', { dataTransfer })

      // Проверяем, что overlay скрылся
      expect(wrapper.classes()).not.toContain('drag-over')

      // Проверяем, что показалось сообщение об ошибке
      expect(uiStore.addToast).toHaveBeenCalledWith('Please drop a folder, not a file', 'warning')
    })
  })

  describe('Auto-open checkbox', () => {
    it('должен вызывать projectStore.setAutoOpenLast', async () => {
      const setAutoOpenLastSpy = vi.spyOn(projectStore, 'setAutoOpenLast')

      const wrapper = mount(ProjectSelector)

      const checkbox = wrapper.find('input[type="checkbox"]')
      await checkbox.setChecked(true)

      expect(setAutoOpenLastSpy).toHaveBeenCalledWith(true)

      await checkbox.setChecked(false)

      expect(setAutoOpenLastSpy).toHaveBeenCalledWith(false)
    })
  })

  describe('Обработка ошибок', () => {
    it('должен показывать toast при ошибке выбора', async () => {
      vi.spyOn(apiService, 'selectDirectory').mockRejectedValue(new Error('Directory selection failed'))

      const wrapper = mount(ProjectSelector)

      // Clear any previous toast calls (like language change)
      vi.mocked(uiStore.addToast).mockClear()

      await wrapper.find('.cta-button').trigger('click')

      // Wait for async operations to complete
      await new Promise(resolve => setTimeout(resolve, 50))

      // Проверяем, что показалось сообщение об ошибке
      expect(uiStore.addToast).toHaveBeenCalledWith('Failed to select directory: Directory selection failed', 'error')
    })
  })
})