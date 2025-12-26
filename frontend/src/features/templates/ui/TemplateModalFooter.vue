<template>
  <div class="tpl-footer">
    <div class="tpl-footer-left">
      <button @click="$emit('import')" class="tpl-footer-btn" :title="t('templates.import')">
        <Upload class="w-3.5 h-3.5" />
      </button>
      <button 
        @click="$emit('export')" 
        class="tpl-footer-btn" 
        :title="t('templates.export')" 
        :disabled="!hasSelectedTemplate"
      >
        <Download class="w-3.5 h-3.5" />
      </button>
      <button 
        @click="$emit('duplicate')" 
        class="tpl-footer-btn" 
        :disabled="!hasSelectedTemplate"
      >
        <Copy class="w-3.5 h-3.5" />
        <span>{{ t('templates.duplicate') }}</span>
      </button>
    </div>
    <div class="tpl-footer-center">
      <Transition name="fade" mode="out-in">
        <span v-if="justSaved" class="tpl-saved-msg">
          <Check class="w-3 h-3" />{{ t('templates.saved') }}
        </span>
        <span v-else-if="hasChanges && !isEditingBuiltIn" class="tpl-unsaved-msg">
          {{ t('templates.unsavedChanges') }}
        </span>
      </Transition>
    </div>
    <div class="tpl-footer-right">
      <label class="tpl-autosave">
        <input 
          type="checkbox" 
          :checked="autoSaveOnClose"
          @change="$emit('update:autoSaveOnClose', ($event.target as HTMLInputElement).checked)" 
        />
        <span class="tpl-autosave-track">
          <span class="tpl-autosave-thumb" />
        </span>
        <span>{{ t('templates.autosave') }}</span>
      </label>
      <button 
        v-if="canApply" 
        @click="$emit('apply')" 
        class="tpl-apply-btn"
      >
        <Sparkles class="w-3.5 h-3.5" />{{ t('templates.apply') }}
      </button>
      <button 
        @click="$emit('save')" 
        class="tpl-save-btn" 
        :class="{ pulse: hasChanges && !isEditingBuiltIn }"
        :disabled="!hasChanges || isEditingBuiltIn"
      >
        <Save class="w-3.5 h-3.5" />{{ t('common.save') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { Check, Copy, Download, Save, Sparkles, Upload } from 'lucide-vue-next'

const { t } = useI18n()

defineProps<{
  hasSelectedTemplate: boolean
  hasChanges: boolean
  isEditingBuiltIn: boolean
  canApply: boolean
  autoSaveOnClose: boolean
  justSaved: boolean
}>()

defineEmits<{
  (e: 'import'): void
  (e: 'export'): void
  (e: 'duplicate'): void
  (e: 'apply'): void
  (e: 'save'): void
  (e: 'update:autoSaveOnClose', value: boolean): void
}>()
</script>

<style scoped>
.tpl-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.625rem 1rem;
  background: rgba(255, 255, 255, 0.02);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.tpl-footer-left,
.tpl-footer-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tpl-footer-center {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tpl-footer-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.625rem;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.15s;
}

.tpl-footer-btn:hover:not(:disabled) { 
  background: rgba(255, 255, 255, 0.05); 
  border-color: rgba(255, 255, 255, 0.2);
  color: var(--text-primary);
}

.tpl-footer-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.tpl-saved-msg {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 11px;
  color: var(--color-success);
}

.tpl-unsaved-msg {
  font-size: 11px;
  color: var(--color-warning);
}

/* Autosave Toggle */
.tpl-autosave {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 11px;
  color: var(--text-muted);
}

.tpl-autosave input {
  display: none;
}

.tpl-autosave-track {
  position: relative;
  width: 28px;
  height: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  transition: background 0.2s;
}

.tpl-autosave input:checked + .tpl-autosave-track {
  background: var(--accent-indigo);
}

.tpl-autosave-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 12px;
  height: 12px;
  background: white;
  border-radius: 50%;
  transition: transform 0.2s;
}

.tpl-autosave input:checked + .tpl-autosave-track .tpl-autosave-thumb {
  transform: translateX(12px);
}

/* Action Buttons */
.tpl-apply-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  background: var(--accent-purple-bg);
  border: 1px solid var(--accent-purple-border);
  border-radius: var(--radius-md);
  color: var(--accent-purple);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tpl-apply-btn:hover { 
  background: rgba(168, 85, 247, 0.25);
  transform: translateY(-1px);
}

.tpl-save-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  background: var(--accent-indigo);
  border: none;
  border-radius: var(--radius-md);
  color: white;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tpl-save-btn:hover:not(:disabled) { 
  background: #818cf8;
  transform: translateY(-1px);
}

.tpl-save-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  transform: none;
}

.tpl-save-btn.pulse {
  animation: btnPulse 0.5s;
}

/* Fade Transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes btnPulse {
  0% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0.5); }
  70% { box-shadow: 0 0 0 8px rgba(99, 102, 241, 0); }
  100% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0); }
}
</style>
