<template>
  <div v-if="error" class="error-boundary p-4 rounded-lg bg-red-900/20 border border-red-500/30">
    <div class="flex items-start gap-3">
      <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-red-500/20 flex items-center justify-center">
        <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <div class="flex-1 min-w-0">
        <h3 class="text-sm font-medium text-red-400">{{ t('error.generic') }}</h3>
        <p class="mt-1 text-xs text-gray-400">{{ errorMessage }}</p>
        <button 
          v-if="canRetry"
          @click="retry"
          class="mt-2 text-xs text-red-400 hover:text-red-300 underline"
        >
          {{ t('chat.retry') }}
        </button>
      </div>
      <button @click="dismiss" class="text-gray-400 hover:text-gray-300">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { getErrorMessage } from '@/composables/useErrorHandler'
import { useI18n } from '@/composables/useI18n'
import { onErrorCaptured, ref } from 'vue'

interface Props {
  canRetry?: boolean
  onRetry?: () => void
}

const props = withDefaults(defineProps<Props>(), {
  canRetry: false
})

const { t } = useI18n()
const error = ref<Error | null>(null)
const errorMessage = ref('')

onErrorCaptured((err) => {
  error.value = err
  errorMessage.value = getErrorMessage(err)
  console.error('[ErrorBoundary] Caught error:', err)
  return false // Prevent error from propagating
})

function dismiss() {
  error.value = null
  errorMessage.value = ''
}

function retry() {
  dismiss()
  props.onRetry?.()
}
</script>
