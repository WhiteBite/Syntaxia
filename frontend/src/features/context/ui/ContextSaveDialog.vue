<template>
  <BaseModal
    v-model="localShow"
    size="sm"
    :close-on-backdrop="true"
    @close="emit('close')"
  >
    <!-- Gradient accent line -->
    <div class="save-modal-accent"></div>
    
    <h3 class="save-modal-title">{{ t('context.saveContext') }}</h3>
    
    <form @submit.prevent="handleSubmit" class="save-modal-content">
      <div class="save-modal-field">
        <label class="save-modal-label">{{ t('context.topic') }}</label>
        <input 
          ref="topicInputRef"
          v-model="localTopic" 
          type="text" 
          class="save-modal-input" 
          :placeholder="t('context.topicPlaceholder')"
          @keyup.enter="localTopic.trim() && handleSubmit()"
        />
      </div>
      <div class="save-modal-field">
        <label class="save-modal-label">{{ t('context.summary') }}</label>
        <textarea 
          v-model="localSummary" 
          class="save-modal-textarea" 
          :placeholder="t('context.summaryPlaceholder')" 
        />
      </div>
    </form>
    
    <template #footer>
      <!-- File count badge -->
      <div class="save-modal-badge">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <span>{{ fileCount }} {{ t('context.filesShort') }}</span>
      </div>
      
      <div class="save-modal-actions">
        <button type="button" @click="emit('close')" class="save-modal-cancel">
          {{ t('context.cancel') }}
        </button>
        <button 
          type="submit" 
          @click="handleSubmit" 
          class="save-modal-submit" 
          :disabled="!localTopic.trim()"
        >
          {{ t('common.save') }}
        </button>
      </div>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import BaseModal from '@/components/ui/BaseModal.vue'
import { useI18n } from '@/composables/useI18n'
import { computed, nextTick, ref, watch } from 'vue'

const { t } = useI18n()

const props = defineProps<{
  show: boolean
  fileCount: number
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', topic: string, summary: string): void
}>()

const localShow = computed({
  get: () => props.show,
  set: (value) => {
    if (!value) emit('close')
  }
})

const localTopic = ref('')
const localSummary = ref('')
const topicInputRef = ref<HTMLInputElement | null>(null)

// Auto-focus topic input when dialog opens
watch(() => props.show, (isOpen) => {
  if (isOpen) {
    localTopic.value = ''
    localSummary.value = ''
    nextTick(() => {
      topicInputRef.value?.focus()
    })
  }
})

function handleSubmit() {
  if (!localTopic.value.trim()) return
  emit('save', localTopic.value.trim(), localSummary.value.trim())
}
</script>

<style scoped>
/* Save Modal - Premium Design */
.save-modal-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #8b5cf6 0%, #6366f1 50%, #3b82f6 100%);
  margin: calc(var(--space-6) * -1);
  margin-bottom: 0;
  width: calc(100% + var(--space-6) * 2);
}

.save-modal-title {
  font-size: 18px;
  font-weight: 700;
  color: white;
  margin-bottom: 24px;
}

.save-modal-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-top: 24px;
}

.save-modal-field {
  display: flex;
  flex-direction: column;
}

.save-modal-label {
  display: block;
  font-size: 10px;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 8px;
}

.save-modal-input {
  width: 100%;
  padding: 14px 16px;
  background: #0f111a;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  font-size: 14px;
  color: white;
  outline: none;
  transition: all 0.2s ease-out;
}

.save-modal-input::placeholder {
  color: #4b5563;
}

.save-modal-input:focus {
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
}

.save-modal-textarea {
  width: 100%;
  height: 90px;
  padding: 14px 16px;
  background: #0f111a;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  font-size: 14px;
  color: white;
  outline: none;
  resize: none;
  transition: all 0.2s ease-out;
}

.save-modal-textarea::placeholder {
  color: #4b5563;
}

.save-modal-textarea:focus {
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
}

.save-modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.save-modal-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid rgba(139, 92, 246, 0.2);
  border-radius: 8px;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #a78bfa;
  font-weight: 500;
}

.save-modal-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.save-modal-cancel {
  padding: 10px 18px;
  background: none;
  border: none;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 500;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.save-modal-cancel:hover {
  color: white;
  background: rgba(255, 255, 255, 0.05);
}

.save-modal-submit {
  padding: 12px 28px;
  background: linear-gradient(135deg, #9333ea 0%, #7c3aed 50%, #6366f1 100%);
  border: none;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 700;
  color: white;
  cursor: pointer;
  transition: all 0.2s ease-out;
  box-shadow: 
    0 4px 16px rgba(139, 92, 246, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1) inset,
    0 1px 0 rgba(255, 255, 255, 0.15) inset;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.save-modal-submit:hover:not(:disabled) {
  background: linear-gradient(135deg, #a855f7 0%, #8b5cf6 50%, #818cf8 100%);
  transform: translateY(-2px);
  box-shadow: 
    0 8px 24px rgba(139, 92, 246, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.15) inset,
    0 1px 0 rgba(255, 255, 255, 0.2) inset;
}

.save-modal-submit:active:not(:disabled) {
  transform: translateY(0);
}

.save-modal-submit:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: linear-gradient(135deg, #6b7280 0%, #4b5563 100%);
  box-shadow: none;
}
</style>
