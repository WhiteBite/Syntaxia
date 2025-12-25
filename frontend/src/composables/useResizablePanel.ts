/**
 * Composable for creating resizable panels
 * Allows users to drag and resize panels with mouse
 * Saves panel width to localStorage for persistence
 * 
 * Uses viewport-relative constraints to ensure panels adapt to screen size
 */

import { useStorage } from '@vueuse/core'
import { onMounted, onUnmounted, ref } from 'vue'

/** Maximum percentage of viewport width a single panel can occupy */
const DEFAULT_MAX_PANEL_WIDTH_PERCENT = 0.4

export interface ResizablePanelOptions {
  minWidth?: number
  maxWidth?: number
  defaultWidth?: number
  storageKey: string
  invertDirection?: boolean // Для правой панели - тянем влево = увеличиваем
  maxWidthPercent?: number // Custom max width as percentage of viewport (0-1)
}

/**
 * Calculate effective max width based on viewport
 * Ensures panel never exceeds maxWidthPercent of viewport
 */
function getEffectiveMaxWidth(maxWidth: number, maxWidthPercent: number): number {
  const viewportMax = Math.floor(window.innerWidth * maxWidthPercent)
  return Math.min(maxWidth, viewportMax)
}

export function useResizablePanel(options: ResizablePanelOptions) {
  const {
    minWidth = 200,
    maxWidth = 800,
    defaultWidth = 300,
    storageKey,
    invertDirection = false,
    maxWidthPercent = DEFAULT_MAX_PANEL_WIDTH_PERCENT
  } = options

  const panelRef = ref<HTMLElement>()

  // Используем useStorage для автоматического сохранения и загрузки
  const width = useStorage(storageKey, defaultWidth, localStorage, {
    mergeDefaults: true,
    serializer: {
      read: (v: string) => {
        const parsed = Number(v)
        if (isNaN(parsed)) return defaultWidth
        // Validate against viewport-relative max width
        const effectiveMax = getEffectiveMaxWidth(maxWidth, maxWidthPercent)
        return Math.max(minWidth, Math.min(effectiveMax, parsed))
      },
      write: (v: number) => String(v)
    }
  })
  const isResizing = ref(false)

  // Re-validate width on mount and window resize
  const validateWidth = () => {
    const effectiveMax = getEffectiveMaxWidth(maxWidth, maxWidthPercent)
    if (width.value > effectiveMax) {
      width.value = effectiveMax
    }
  }

  onMounted(() => {
    validateWidth()
    window.addEventListener('resize', validateWidth)
  })

  let startX = 0
  let startWidth = 0

  const onMouseDown = (e: MouseEvent) => {
    if (!panelRef.value) return

    isResizing.value = true
    startX = e.clientX
    startWidth = width.value

    // Prevent text selection during resize
    document.body.style.userSelect = 'none'
    document.body.style.cursor = 'col-resize'

    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  const onMouseMove = (e: MouseEvent) => {
    if (!isResizing.value) return

    const delta = e.clientX - startX
    // Для правой панели инвертируем: тянем влево = увеличиваем ширину
    const adjustedDelta = invertDirection ? -delta : delta
    // Use viewport-relative max width
    const effectiveMax = getEffectiveMaxWidth(maxWidth, maxWidthPercent)
    const newWidth = Math.max(minWidth, Math.min(effectiveMax, startWidth + adjustedDelta))
    width.value = newWidth
  }

  const onMouseUp = () => {
    if (!isResizing.value) return

    isResizing.value = false

    // Restore default cursor and text selection
    document.body.style.userSelect = ''
    document.body.style.cursor = ''

    // useStorage автоматически сохраняет изменения - ручное сохранение удалено

    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }

  onUnmounted(() => {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    window.removeEventListener('resize', validateWidth)
  })

  // Reset to default width
  const resetToDefault = () => {
    width.value = defaultWidth
  }

  return {
    panelRef,
    width,
    isResizing,
    onMouseDown,
    resetToDefault,
    defaultWidth
  }
}
