/**
 * Global error handler plugin for Vue
 * Catches unhandled errors and displays user-friendly messages
 */

import type { App } from 'vue'
import { useUIStore } from '@/stores/ui.store'

export interface ErrorHandlerOptions {
  showNotification?: boolean
  logToConsole?: boolean
}

export function setupErrorHandler(app: App, options: ErrorHandlerOptions = {}) {
  const {
    showNotification = true,
    logToConsole = true
  } = options

  // Vue error handler
  app.config.errorHandler = (err, instance, info) => {
    if (logToConsole) {
      console.error('[Error Handler]', err)
      console.error('Component:', instance)
      console.error('Error info:', info)
    }

    if (showNotification) {
      try {
        const uiStore = useUIStore()
        const message = err instanceof Error ? err.message : String(err)
        uiStore.addToast(`Application error: ${message}`, 'error')
      } catch (toastError) {
        // Fallback to console if toast system fails
        console.error('Failed to show error toast:', toastError)
        if (import.meta.env.PROD) {
          alert(`An error occurred: ${err instanceof Error ? err.message : String(err)}`)
        }
      }
    }

    // Additional info in development
    if (import.meta.env.DEV) {
      console.error('Error details:', { err, instance, info })
    }
  }

  // Handle unhandled promise rejections
  window.addEventListener('unhandledrejection', (event) => {
    if (logToConsole) {
      console.error('[Unhandled Promise Rejection]', event.reason)
    }

    if (showNotification) {
      try {
        const uiStore = useUIStore()
        const message = event.reason instanceof Error ? event.reason.message : String(event.reason)
        uiStore.addToast(`Unhandled promise rejection: ${message}`, 'error')
      } catch (toastError) {
        console.error('Failed to show rejection toast:', toastError)
      }
    }

    // Prevent default error handling
    event.preventDefault()
  })
}
