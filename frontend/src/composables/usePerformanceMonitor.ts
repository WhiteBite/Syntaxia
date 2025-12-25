/**
 * Composable for monitoring component performance
 * Logs warnings for slow-rendering components
 */

import { ref, onBeforeMount, onMounted } from 'vue'

const SLOW_RENDER_THRESHOLD = 100 // ms

export function usePerformanceMonitor(componentName: string) {
  const renderTime = ref(0)
  const mountTime = ref(0)

  onBeforeMount(() => {
    mountTime.value = performance.now()
  })

  onMounted(() => {
    const duration = performance.now() - mountTime.value
    renderTime.value = duration

    // Log warning for slow renders
    if (duration > SLOW_RENDER_THRESHOLD) {
      console.warn(
        `[Performance] Slow render detected: ${componentName} took ${duration.toFixed(2)}ms to mount`
      )
    } else if (import.meta.env.DEV) {
      console.debug(
        `[Performance] ${componentName} mounted in ${duration.toFixed(2)}ms`
      )
    }
  })

  return {
    renderTime
  }
}
