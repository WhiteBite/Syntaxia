<template>
  <div id="app" class="app-container layout-root">
    <!-- Skip link for keyboard navigation -->
    <a href="#main-content" class="skip-link">{{ t('accessibility.skipToContent') }}</a>

    <!-- Global Error Handler -->
    <div v-if="globalError" class="modal-container">
      <div class="modal-overlay" @click="clearGlobalError"></div>
      <div class="modal-content modal-error">
        <h3 class="modal-title">{{ t('common.criticalError') }}</h3>
        <p class="modal-text">{{ globalError }}</p>
        <button @click="clearGlobalError" class="btn btn-danger">
          {{ t('common.close') }}
        </button>
      </div>
    </div>

    <!-- Global Loading Overlay -->
    <div v-if="projectStore.isLoading" class="loading-overlay">
      <div class="loading-card">
        <div class="flex items-center justify-center mb-4">
          <svg class="loading-spinner" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
        </div>
        <h3 class="loading-title">{{ t('common.loadingProject') }}</h3>
        <p class="loading-text">{{ t('common.pleaseWait') }}</p>
      </div>
    </div>

    <!-- Main Content: Show ProjectSelector if no project, otherwise show workspace -->
    <main id="main-content" class="layout-fill layout-column layout-clip">
      <ProjectSelector v-if="!projectStore.hasProject" @opened="onProjectOpened" />
      <MainWorkspace v-else />
    </main>

    <!-- Toast Notifications (Bottom Center) - compact glassmorphism style -->
    <div class="fixed bottom-3 left-1/2 -translate-x-1/2 z-[100] flex flex-col items-center gap-2">
      <TransitionGroup name="toast">
        <div
          v-for="toast in uiStore.toasts"
          :key="toast.id"
          :class="[
            'toast-glass',
            toast.type === 'error' ? 'toast-glass-error' :
            toast.type === 'success' ? 'toast-glass-success' :
            toast.type === 'warning' ? 'toast-glass-warning' : 'toast-glass-info'
          ]"
        >
          <!-- Icon with glow wrapper -->
          <div class="toast-glass-icon-wrap">
            <svg v-if="toast.type === 'success'" class="toast-glass-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
            </svg>
            <svg v-else-if="toast.type === 'error'" class="toast-glass-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
            </svg>
            <svg v-else-if="toast.type === 'warning'" class="toast-glass-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <svg v-else class="toast-glass-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span class="toast-glass-message">{{ toast.message }}</span>
          <!-- Divider + Action Button -->
          <template v-if="toast.action">
            <div class="toast-glass-divider" />
            <button 
              @click="toast.action.onClick(); uiStore.removeToast(toast.id)"
              class="toast-glass-action"
            >
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              {{ toast.action.label }}
            </button>
          </template>
          <button @click="uiStore.removeToast(toast.id)" class="toast-glass-close">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>

    <!-- Global components (lazy loaded) -->
    <CommandPalette 
      v-if="projectStore.hasProject" 
      :is-open="isCommandPaletteOpen" 
      @close="isCommandPaletteOpen = false" 
    />
    <KeyboardShortcutsModal 
      v-if="projectStore.hasProject" 
      :is-open="isShortcutsModalOpen" 
      @close="isShortcutsModalOpen = false" 
    />
    
    <!-- Memory Dashboard -->
    <MemoryDashboard
      v-if="showMemoryDashboard && projectStore.hasProject"
      class="fixed top-20 right-4 z-40"
      @close="showMemoryDashboard = false"
    />
    
    <!-- Settings Modal -->
    <SettingsModal v-model="uiStore.showSettingsModal" />

    <!-- Confirm Dialog (global) -->
    <ConfirmDialog />

  </div>
</template>

<script setup lang="ts">
import ProjectSelector from '@/components/ProjectSelector.vue'
import MainWorkspace from '@/components/workspace/MainWorkspace.vue'
import { useI18n } from '@/composables/useI18n'
import { useMemoryMonitor } from '@/composables/useMemoryMonitor'
import { useOnboarding } from '@/composables/useOnboarding'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { shellApi } from '@/services/api/shell.api'
import { useMagicKeys } from '@vueuse/core'
import { defineAsyncComponent, onMounted, onUnmounted, ref, watch } from 'vue'

// Lazy load heavy components
const CommandPalette = defineAsyncComponent(() => import('@/components/CommandPalette.vue'))
const KeyboardShortcutsModal = defineAsyncComponent(() => import('@/components/KeyboardShortcutsModal.vue'))
// ThemeToggle removed - app uses dark theme only
const MemoryDashboard = defineAsyncComponent(() => import('@/components/MemoryDashboard.vue'))
const SettingsModal = defineAsyncComponent(() => import('@/components/SettingsModal.vue'))
const ConfirmDialog = defineAsyncComponent(() => import('@/components/ConfirmDialog.vue'))

const { t } = useI18n()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const globalError = ref<string | null>(null)
const isCommandPaletteOpen = ref(false)
const isShortcutsModalOpen = ref(false)
const showMemoryDashboard = ref(false)

// Keyboard shortcuts
const keys = useMagicKeys()
const ctrlK = keys['Ctrl+K']
const ctrlP = keys['Ctrl+P']
const ctrlSlash = keys['Ctrl+/']
const ctrlB = keys['Ctrl+B']
const ctrlEnter = keys['Ctrl+Enter']
const ctrlE = keys['Ctrl+E']
const ctrlShiftC = keys['Ctrl+Shift+C']
const ctrlShiftM = keys['Ctrl+Shift+M']
const ctrlComma = keys['Ctrl+,']

watch(ctrlK, (v) => {
  if (v) isCommandPaletteOpen.value = !isCommandPaletteOpen.value
})
watch(ctrlP, (v) => {
  if (v) {
    isCommandPaletteOpen.value = true
    // Here you might want to pre-fill the command palette with a file search prefix
  }
})
watch(ctrlSlash, (v) => {
  if (v) isShortcutsModalOpen.value = !isShortcutsModalOpen.value
})

// Build context shortcuts
watch([ctrlB, ctrlEnter], ([b, enter]) => {
  if ((b || enter) && projectStore.hasProject) {
    // Trigger build context
    const event = new CustomEvent('global-build-context')
    window.dispatchEvent(event)
  }
})

// Export shortcut
watch(ctrlE, (v) => {
  if (v && projectStore.hasProject) {
    const event = new CustomEvent('global-open-export')
    window.dispatchEvent(event)
  }
})

// Copy context shortcut
watch(ctrlShiftC, (v) => {
  if (v && projectStore.hasProject) {
    const event = new CustomEvent('global-copy-context')
    window.dispatchEvent(event)
  }
})

// Undo/Redo selection shortcuts
const ctrlZ = keys['Ctrl+Z']
const ctrlY = keys['Ctrl+Y']

watch(ctrlZ, (v) => {
  if (v && projectStore.hasProject) {
    // Only trigger if not in an input field
    const activeEl = document.activeElement
    if (activeEl?.tagName !== 'INPUT' && activeEl?.tagName !== 'TEXTAREA') {
      const event = new CustomEvent('global-undo-selection')
      window.dispatchEvent(event)
    }
  }
})

watch(ctrlY, (v) => {
  if (v && projectStore.hasProject) {
    const activeEl = document.activeElement
    if (activeEl?.tagName !== 'INPUT' && activeEl?.tagName !== 'TEXTAREA') {
      const event = new CustomEvent('global-redo-selection')
      window.dispatchEvent(event)
    }
  }
})

// Memory dashboard shortcut
watch(ctrlShiftM, (v) => {
  if (v && projectStore.hasProject) {
    showMemoryDashboard.value = !showMemoryDashboard.value
  }
})

// Settings shortcut
watch(ctrlComma, (v) => {
  if (v) {
    uiStore.openSettingsModal()
  }
})

// Start memory monitoring
const memoryMonitor = useMemoryMonitor()

// Onboarding tour
const onboarding = useOnboarding()

function clearGlobalError() {
  globalError.value = null
}

function onProjectOpened(_path: string) {
  uiStore.addToast('Project loaded successfully', 'success')
  
  // Start onboarding tour for new users after project is loaded
  if (onboarding.shouldShowTour()) {
    // Small delay to ensure UI is ready
    setTimeout(() => {
      onboarding.startTour()
    }, 500)
  }
}

// Store cleanup interval ref
const cleanupIntervalRef = ref<number | null>(null)

onMounted(async () => {
  // Start memory monitoring FIRST (before any heavy operations)
  // In dev mode, monitoring is automatically throttled
  memoryMonitor.startMonitoring()
  
  // Register critical memory callback to show dashboard
  memoryMonitor.onCritical(() => {
    showMemoryDashboard.value = true
  })
  
  // Check for startup path from context menu
  try {
    const startupPath = await shellApi.getStartupPath()
    if (startupPath) {
      await shellApi.clearStartupPath()
      await projectStore.openProjectByPath(startupPath)
      uiStore.addToast('Project opened from context menu', 'success')
      return // Skip auto-open if we have startup path
    }
  } catch {
    // Ignore startup path errors
  }
  
  // Note: fetchRecentProjects is called in ProjectSelector.vue onMounted
  // We don't call it here to avoid duplicate calls
  
  // Try to auto-open last project if setting is enabled
  // This will only work if recent projects are already loaded from localStorage
  projectStore.maybeAutoOpenLastProject()
  
  // Periodic cleanup - DISABLED in dev mode to prevent memory leaks from dynamic imports
  // In dev mode, Vite's HMR + dynamic imports can cause memory accumulation
  if (!import.meta.env.DEV) {
    cleanupIntervalRef.value = window.setInterval(async () => {
      // Clear old caches
      try {
        const { clearAllCaches, getCacheStats } = await import('@/composables/useApiCache')
        const stats = getCacheStats()
        if (stats.size > stats.maxSize * 0.5) {
          clearAllCaches()
        }
      } catch {
        // Ignore cache clear errors
      }
      
      // Memory stats checked silently
      await memoryMonitor.getMemoryStats()
    }, 10 * 60 * 1000) // Every 10 minutes
  }
})

onUnmounted(() => {
  // Clear cleanup interval
  if (cleanupIntervalRef.value) {
    clearInterval(cleanupIntervalRef.value)
    cleanupIntervalRef.value = null
  }
  
  // Stop memory monitoring
  memoryMonitor.stopMonitoring()
  
  // Clear all caches
  import('@/composables/useApiCache').then(({ clearAllCaches }) => {
    clearAllCaches()
  }).catch(() => {
    // Ignore
  })
})

// Global error handler for memory errors (moved outside onMounted)
onMounted(() => {
  window.addEventListener('error', (event) => {
    if (event.message && event.message.includes('out of memory')) {
      console.error('[App] Out of memory error detected!')
      uiStore.addToast('Critical memory error. Clearing caches...', 'error')
      
      // Emergency cleanup
      try {
        const { clearAllCaches } = require('@/composables/useApiCache')
        clearAllCaches()
        
        const { useContextStore } = require('@/features/context/model/context.store')
        const { useFileStore } = require('@/features/files/model/file.store')
        const contextStore = useContextStore()
        const fileStore = useFileStore()
        
        contextStore.clearContext()
        fileStore.resetStore()
      } catch (e) {
        console.error('Emergency cleanup failed:', e)
      }
    }
  })
})
</script>

<style>
/* Корневой контейнер - использует layout-root из layout.css */
.app-container {
  /* Layout: см. layout.css */
  color: white;
  background: var(--bg-app);
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-feature-settings: 'cv02', 'cv03', 'cv04', 'cv11';
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Global scrollbar styles */
* {
  scrollbar-width: thin;
  scrollbar-color: #6b7280 rgba(31, 41, 55, 0.4);
}

*::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

*::-webkit-scrollbar-track {
  background: rgba(31, 41, 55, 0.4);
  border-radius: 4px;
}

*::-webkit-scrollbar-thumb {
  background: #6b7280;
  border-radius: 4px;
  border: 1px solid rgba(31, 41, 55, 0.6);
}

*::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

*::-webkit-scrollbar-corner {
  background: rgba(31, 41, 55, 0.4);
}

/* Focus styles */
*:focus {
  outline: none;
}

/* Scope focus outline to interactive controls to avoid artifacts during resize */
:where(
  a,
  button,
  input,
  select,
  textarea,
  [role="button"],
  [role="checkbox"],
  [role="menuitem"],
  [role="tab"],
  [role="switch"],
  [contenteditable],
  [tabindex]:not([tabindex="-1"]) 
):focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}


/* Transition defaults */
* {
  transition-property: color, background-color, border-color, text-decoration-color, fill, stroke, opacity, box-shadow, transform, filter, backdrop-filter;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

/* Selection styling */
::selection {
  background-color: rgba(59, 130, 246, 0.3);
  color: #ffffff;
}

::-moz-selection {
  background-color: rgba(59, 130, 246, 0.3);
  color: #ffffff;
}

/* Animations */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* Improved button focus and hover states */
button:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}

/* Better form input styles */
input, textarea, select {
  transition: all 0.2s ease-in-out;
}

input:focus, textarea:focus, select:focus {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* Loading animation improvements */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Toast animations */
.toast-enter-active {
  animation: toast-in 0.3s ease-out;
}

.toast-leave-active {
  animation: toast-out 0.2s ease-in;
}

@keyframes toast-in {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes toast-out {
  from {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
  to {
    opacity: 0;
    transform: translateY(-10px) scale(0.95);
  }
}
</style>