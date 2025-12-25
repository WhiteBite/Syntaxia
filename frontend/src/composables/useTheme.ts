/**
 * Composable for theme management (dark/light/system)
 * Persists theme preference to localStorage
 */

import { ref, watch, onMounted } from 'vue'

const THEME_KEY = 'Syntaxia_theme'

export type Theme = 'dark' | 'light' | 'system'

// Singleton state shared across all instances
const currentTheme = ref<Theme>('dark')
let initialized = false

export function useTheme() {
  const applyTheme = () => {
    const root = document.documentElement
    
    if (currentTheme.value === 'system') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      root.classList.toggle('dark', prefersDark)
    } else {
      root.classList.toggle('dark', currentTheme.value === 'dark')
    }
  }

  const setTheme = (newTheme: Theme) => {
    currentTheme.value = newTheme
  }

  const cycleTheme = () => {
    const themes: Theme[] = ['dark', 'light', 'system']
    const currentIndex = themes.indexOf(currentTheme.value)
    const nextIndex = (currentIndex + 1) % themes.length
    currentTheme.value = themes[nextIndex]
  }

  // Initialize only once
  if (!initialized) {
    initialized = true
    
    // Load saved theme from localStorage
    onMounted(() => {
      try {
        const saved = localStorage.getItem(THEME_KEY)
        if (saved && ['dark', 'light', 'system'].includes(saved)) {
          currentTheme.value = saved as Theme
        }
      } catch (err) {
        console.warn('Failed to load theme from localStorage:', err)
      }
    })

    // Watch for theme changes and apply them
    watch(currentTheme, (newTheme) => {
      try {
        localStorage.setItem(THEME_KEY, newTheme)
      } catch (err) {
        console.warn('Failed to save theme to localStorage:', err)
      }
      applyTheme()
    }, { immediate: true })

    // Listen for system theme changes
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', () => {
      if (currentTheme.value === 'system') {
        applyTheme()
      }
    })
  }

  return {
    theme: currentTheme,
    setTheme,
    cycleTheme,
    applyTheme
  }
}
