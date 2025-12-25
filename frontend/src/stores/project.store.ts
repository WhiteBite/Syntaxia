/**
 * Simple Project Store - No Clean Architecture dependencies
 * Direct localStorage integration for simplicity
 */

import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { apiService } from '@/services/api.service'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export interface RecentProject {
  path: string
  name: string
  lastOpened: number
}

const RECENT_PROJECTS_KEY = 'Syntaxia_recent_projects'
const MAX_RECENT = 10

export const useProjectStore = defineStore('project', () => {
  // State
  const currentPath = ref<string | null>(null)
  const currentName = ref<string | null>(null)
  const recentProjects = ref<RecentProject[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const autoOpenLast = ref<boolean>(false)

  // Computed
  const hasProject = computed(() => currentPath.value !== null)
  const projectName = computed(() => currentName.value || '')
  const projectPath = computed(() => currentPath.value || '')
  const hasRecentProjects = computed(() => recentProjects.value.length > 0)
  const lastProjectPath = computed(() => recentProjects.value[0]?.path || null)

  // Actions
  async function openProjectByPath(path: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      // Verify path exists before opening
      const exists = await apiService.pathExists(path)
      if (!exists) {
        error.value = `Path does not exist: ${path}`
        // Remove from recent projects if it doesn't exist
        recentProjects.value = recentProjects.value.filter(p => p.path !== path)
        saveRecentProjects()
        return false
      }

      // Simple: just set the path and name
      currentPath.value = path
      currentName.value = path.split(/[\\/]/).pop() || path

      // Reset file store and context store when project changes
      const fileStore = useFileStore()
      const contextStore = useContextStore()
      fileStore.resetStore()
      contextStore.clearContext()

      // Update recent projects
      addToRecent(path)

      // Сохраняем проект в backend через apiService
      try {
        const name = path.split(/[\\/]/).pop() || path
        await apiService.addRecentProject(path, name)
      } catch (backendError) {
        console.error('Failed to save project to backend:', backendError)
        // Продолжаем выполнение, даже если сохранение в бэкенд не удалось
      }

      saveRecentProjects()

      return true
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load project'
      return false
    } finally {
      isLoading.value = false
    }
  }

  function addToRecent(path: string) {
    const name = path.split(/[\\/]/).pop() || path
    const existing = recentProjects.value.findIndex(p => p.path === path)

    if (existing !== -1) {
      recentProjects.value.splice(existing, 1)
    }

    recentProjects.value.unshift({
      path,
      name,
      lastOpened: Date.now()
    })

    if (recentProjects.value.length > MAX_RECENT) {
      recentProjects.value = recentProjects.value.slice(0, MAX_RECENT)
    }
  }

  function loadRecentProjects() {
    try {
      const stored = localStorage.getItem(RECENT_PROJECTS_KEY)
      if (stored) {
        recentProjects.value = JSON.parse(stored)
      }
      // Load setting for auto-open
      const auto = localStorage.getItem('Syntaxia_auto_open_last')
      if (auto !== null) {
        autoOpenLast.value = auto === 'true'
      }
    } catch (err) {
      console.warn('Failed to load recent projects:', err)
    }
  }

  /**
   * Асинхронно загружает список недавних проектов из бэкенда
   * Используется для принудительного обновления списка при монтировании компонента
   */
  async function fetchRecentProjects() {
    try {
      isLoading.value = true
      error.value = null

      // Вызываем метод из нашего api.service.ts
      const projectsJson = await apiService.getRecentProjects()

      // Парсим JSON ответ от бэкенда
      if (projectsJson) {
        const projects = JSON.parse(projectsJson)
        if (Array.isArray(projects)) {
          recentProjects.value = projects
          // Также сохраняем в localStorage для кеширования
          saveRecentProjects()
        }
      }
    } catch (err) {
      console.error('Failed to fetch recent projects from backend:', err)
      error.value = err instanceof Error ? err.message : 'Failed to load recent projects'
      // В случае ошибки оставляем текущий список (из localStorage)
    } finally {
      isLoading.value = false
    }
  }

  function saveRecentProjects() {
    try {
      localStorage.setItem(RECENT_PROJECTS_KEY, JSON.stringify(recentProjects.value))
    } catch (err) {
      console.warn('Failed to save recent projects:', err)
    }
  }

  function setAutoOpenLast(value: boolean) {
    autoOpenLast.value = value
    try {
      localStorage.setItem('Syntaxia_auto_open_last', String(value))
    } catch (err) {
      console.warn('Failed to save auto-open setting:', err)
    }
  }

  async function maybeAutoOpenLastProject(): Promise<boolean> {
    if (currentPath.value) return true
    if (!autoOpenLast.value) return false
    if (!lastProjectPath.value) return false
    return await openProjectByPath(lastProjectPath.value)
  }

  async function removeFromRecent(path: string) {
    recentProjects.value = recentProjects.value.filter(p => p.path !== path)

    // Также удаляем из backend через apiService
    try {
      await apiService.removeRecentProject(path)
    } catch (backendError) {
      console.error('Failed to remove project from backend:', backendError)
      // Продолжаем выполнение, даже если удаление из бэкенда не удалось
    }

    saveRecentProjects()
  }

  function clearRecent() {
    recentProjects.value = []
    saveRecentProjects()
  }

  function clearProject() {
    currentPath.value = null
    currentName.value = null
    error.value = null

    // Also reset file store and context store when clearing project
    const fileStore = useFileStore()
    const contextStore = useContextStore()
    fileStore.resetStore()
    contextStore.clearContext()
  }

  // Initialize
  loadRecentProjects()

  return {
    // State
    currentPath,
    currentName,
    recentProjects,
    isLoading,
    error,
    autoOpenLast,
    // Computed
    hasProject,
    projectName,
    projectPath,
    hasRecentProjects,
    lastProjectPath,
    // Actions
    openProjectByPath,
    fetchRecentProjects, // Новый action для загрузки из бэкенда
    removeFromRecent,
    clearRecent,
    clearProject,
    setAutoOpenLast,
    maybeAutoOpenLastProject
  }
})