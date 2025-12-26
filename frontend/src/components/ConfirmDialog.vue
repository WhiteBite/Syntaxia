<template>
  <Teleport to="body">
    <Transition name="modal">
      <div 
        v-if="isOpen" 
        class="modal-container" 
        @click.self="handleCancel"
        @keydown.escape="handleCancel"
        @keydown.enter="handleConfirm"
      >
        <div class="modal-content max-w-sm" role="alertdialog" aria-modal="true" :aria-labelledby="'confirm-title'">
          <!-- Icon -->
          <div class="flex justify-center mb-4">
            <div 
              class="w-12 h-12 rounded-full flex items-center justify-center"
              :class="iconBgClass"
            >
              <AlertTriangle v-if="variant === 'warning' || variant === 'danger'" class="w-6 h-6" :class="iconClass" />
              <Info v-else class="w-6 h-6" :class="iconClass" />
            </div>
          </div>

          <!-- Title -->
          <h3 id="confirm-title" class="text-lg font-semibold text-white text-center mb-2">
            {{ options?.title }}
          </h3>

          <!-- Message -->
          <p class="text-sm text-gray-400 text-center mb-6">
            {{ options?.message }}
          </p>

          <!-- Actions -->
          <div class="flex gap-3">
            <button
              ref="cancelBtnRef"
              @click="handleCancel"
              class="flex-1 btn-unified btn-unified-secondary"
            >
              {{ options?.cancelText || t('common.cancel') }}
            </button>
            <button
              @click="handleConfirm"
              class="flex-1 btn-unified"
              :class="confirmBtnClass"
            >
              {{ options?.confirmText || t('common.confirm') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useConfirm } from '@/composables/useConfirm'
import { useI18n } from '@/composables/useI18n'
import { AlertTriangle, Info } from 'lucide-vue-next'
import { computed, nextTick, ref, watch } from 'vue'

const { t } = useI18n()
const { isOpen, options, handleConfirm, handleCancel } = useConfirm()
const cancelBtnRef = ref<HTMLButtonElement | null>(null)

// Auto-focus cancel button when dialog opens (safer default)
watch(isOpen, (open) => {
    if (open) {
        nextTick(() => cancelBtnRef.value?.focus())
    }
})

const variant = computed(() => options.value?.variant || 'info')

const iconBgClass = computed(() => {
    switch (variant.value) {
        case 'danger': return 'bg-red-500/20'
        case 'warning': return 'bg-yellow-500/20'
        default: return 'bg-blue-500/20'
    }
})

const iconClass = computed(() => {
    switch (variant.value) {
        case 'danger': return 'text-red-400'
        case 'warning': return 'text-yellow-400'
        default: return 'text-blue-400'
    }
})

const confirmBtnClass = computed(() => {
    switch (variant.value) {
        case 'danger': return 'btn-unified-danger'
        case 'warning': return 'btn-unified-warning'
        default: return 'btn-unified-primary'
    }
})
</script>
