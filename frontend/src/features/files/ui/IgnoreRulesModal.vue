<template>
  <Teleport to="body">
    <Transition name="modal-backdrop">
      <div
        v-if="isOpen"
        class="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 flex items-center justify-center p-4"
        @click.self="close"
      >
        <Transition name="modal">
          <div
            v-if="isOpen"
            class="ignore-modal"
            @click.stop
          >
            <!-- Header -->
            <div class="ignore-modal__header">
              <div class="flex items-center gap-3">
                <div class="ignore-modal__icon">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                </div>
                <div>
                  <h3 class="text-lg font-semibold text-white">{{ t('ignoreModal.title') }}</h3>
                  <p class="text-xs text-gray-400">{{ t('ignoreModal.subtitle') }}</p>
                </div>
              </div>
              
              <button @click="close" class="ignore-modal__close">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Tabs -->
            <div class="ignore-modal__tabs">
              <button
                @click="currentTab = 'gitignore'"
                :class="['ignore-tab', currentTab === 'gitignore' ? 'ignore-tab--active' : '']"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                .gitignore
              </button>
              <button
                @click="currentTab = 'custom'"
                :class="['ignore-tab', currentTab === 'custom' ? 'ignore-tab--active' : '']"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                </svg>
                {{ t('ignoreModal.customRules') }}
              </button>
            </div>

            <!-- Content -->
            <div class="ignore-modal__content">
              <!-- .gitignore Tab -->
              <div v-if="currentTab === 'gitignore'" class="ignore-custom-layout">
                <div class="ignore-editor-section">
                  <div class="ignore-editor-header">
                    <span class="ignore-editor-label">{{ t('ignoreModal.rulesCount') }}: <strong>{{ gitignoreRulesCount }}</strong></span>
                  </div>
                  <div class="ignore-editor">
                    <div ref="gitignoreLinesRef" class="ignore-editor__lines">
                      <span v-for="n in gitignoreLineCount" :key="n" :class="{ 'ignore-editor__line--comment': isGitignoreCommentLine(n) }">{{ n }}</span>
                    </div>
                    <textarea
                      ref="gitignoreTextareaRef"
                      v-model="gitignoreContent"
                      readonly
                      class="ignore-editor__textarea"
                      :placeholder="t('ignoreModal.gitignorePlaceholder')"
                      @scroll="syncGitignoreScroll"
                    ></textarea>
                  </div>
                  <p class="ignore-editor__hint">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ t('ignoreModal.gitignoreInfo') }}
                  </p>
                </div>

                <!-- Gitignore Preview -->
                <IgnorePreviewPanel :result="gitignorePreview" :loading="gitignoreLoading" />
              </div>

              <!-- Custom Rules Tab -->
              <div v-if="currentTab === 'custom'" class="ignore-custom-layout">
                <div class="ignore-editor-section">
                  <div class="ignore-editor-header">
                    <span class="ignore-editor-label">{{ t('ignoreModal.rulesCount') }}: <strong>{{ rulesCount }}</strong></span>
                  </div>
                  <div class="ignore-editor">
                    <div ref="customLinesRef" class="ignore-editor__lines">
                      <span v-for="n in customLineCount" :key="n" :class="{ 'ignore-editor__line--comment': isCommentLine(n) }">{{ n }}</span>
                    </div>
                    <textarea
                      ref="customTextareaRef"
                      v-model="customRules"
                      class="ignore-editor__textarea ignore-editor__textarea--editable"
                      :placeholder="t('ignoreModal.customPlaceholder')"
                      spellcheck="false"
                      @scroll="syncCustomScroll"
                    ></textarea>
                  </div>
                </div>

                <!-- Live Preview -->
                <IgnorePreviewPanel :result="customPreview" :loading="customLoading" />
              </div>
            </div>

            <!-- Footer -->
            <div class="ignore-modal__footer">
              <div class="ignore-footer__left">
                <button 
                  v-if="currentTab === 'custom'" 
                  @click="resetToDefaults" 
                  class="ignore-footer__danger-btn"
                >
                  {{ t('ignoreModal.reset') }}
                </button>
                <button 
                  v-if="currentTab === 'custom'" 
                  @click="clearAll" 
                  class="ignore-footer__danger-btn"
                >
                  {{ t('ignoreModal.clearAll') }}
                </button>
              </div>
              <div class="ignore-footer__right">
                <button @click="close" class="btn btn-secondary">
                  {{ t('ignoreModal.cancel') }}
                </button>
                <button
                  v-if="currentTab === 'custom'"
                  @click="save"
                  :disabled="isSaving"
                  class="ignore-footer__save-btn"
                >
                  <svg v-if="!isSaving" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                  </svg>
                  {{ t('ignoreModal.save') }}
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { apiService } from '@/services/api.service'
import type { IgnorePreviewResult } from '@/services/api/settings.api'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, ref, watch } from 'vue'
import { parseIgnoreRules } from '../lib/file-utils'
import IgnorePreviewPanel from './IgnorePreviewPanel.vue'

const { t } = useI18n()
const settingsStore = useSettingsStore()
const projectStore = useProjectStore()
const uiStore = useUIStore()

const isOpen = ref(false)
const currentTab = ref<'gitignore' | 'custom'>('custom')
const gitignoreContent = ref('')
const customRules = ref('')
const isSaving = ref(false)

// Preview results (detailed)
const gitignorePreview = ref<IgnorePreviewResult | null>(null)
const customPreview = ref<IgnorePreviewResult | null>(null)
const gitignoreLoading = ref(false)
const customLoading = ref(false)

// Refs for scroll sync
const gitignoreLinesRef = ref<HTMLElement | null>(null)
const gitignoreTextareaRef = ref<HTMLTextAreaElement | null>(null)
const customLinesRef = ref<HTMLElement | null>(null)
const customTextareaRef = ref<HTMLTextAreaElement | null>(null)

function syncGitignoreScroll() {
  if (gitignoreLinesRef.value && gitignoreTextareaRef.value) {
    gitignoreLinesRef.value.scrollTop = gitignoreTextareaRef.value.scrollTop
  }
}

function syncCustomScroll() {
  if (customLinesRef.value && customTextareaRef.value) {
    customLinesRef.value.scrollTop = customTextareaRef.value.scrollTop
  }
}

const rulesCount = computed(() => parseIgnoreRules(customRules.value).length)
const gitignoreRulesCount = computed(() => parseIgnoreRules(gitignoreContent.value).length)
const gitignoreLineCount = computed(() => Math.max(gitignoreContent.value.split('\n').length, 10))
const customLineCount = computed(() => Math.max(customRules.value.split('\n').length, 10))

function isCommentLine(lineNum: number): boolean {
  const lines = customRules.value.split('\n')
  return lines[lineNum - 1]?.trim().startsWith('#') || false
}

function isGitignoreCommentLine(lineNum: number): boolean {
  const lines = gitignoreContent.value.split('\n')
  return lines[lineNum - 1]?.trim().startsWith('#') || false
}

// Debounced preview for custom rules
let debounceTimer: ReturnType<typeof setTimeout> | null = null
watch(customRules, (newRules) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    if (newRules.trim() && projectStore.currentPath) {
      loadCustomPreview()
    } else {
      customPreview.value = null
    }
  }, 500)
})

async function loadCustomPreview() {
  if (!projectStore.currentPath) return
  customLoading.value = true
  try {
    customPreview.value = await apiService.testIgnoreRulesDetailed(projectStore.currentPath, customRules.value)
  } catch { customPreview.value = null }
  finally { customLoading.value = false }
}

async function loadGitignorePreview() {
  if (!projectStore.currentPath || !gitignoreContent.value.trim()) return
  gitignoreLoading.value = true
  try {
    gitignorePreview.value = await apiService.testIgnoreRulesDetailed(projectStore.currentPath, gitignoreContent.value)
  } catch { gitignorePreview.value = null }
  finally { gitignoreLoading.value = false }
}

async function open() {
  isOpen.value = true
  customRules.value = settingsStore.getCustomIgnoreRules()
  
  if (projectStore.currentPath) {
    try {
      gitignoreContent.value = await apiService.getGitignoreContent(projectStore.currentPath)
      loadGitignorePreview()
    } catch {
      gitignoreContent.value = '# No .gitignore file found'
    }
  }
  
  if (customRules.value.trim()) loadCustomPreview()
}

function close() {
  isOpen.value = false
  gitignorePreview.value = null
  customPreview.value = null
}

async function save() {
  isSaving.value = true
  try {
    await apiService.updateCustomIgnoreRules(customRules.value)
    settingsStore.setCustomIgnoreRules(customRules.value)
    uiStore.addToast(t('ignoreModal.saveSuccess'), 'success')
    close()
  } catch {
    uiStore.addToast('Failed to save ignore rules', 'error')
  } finally {
    isSaving.value = false
  }
}

function resetToDefaults() {
  customRules.value = ''
  customPreview.value = null
}

function clearAll() {
  customRules.value = ''
  customPreview.value = null
}

defineExpose({ open, close })
</script>

<style scoped>
.ignore-modal {
  background: rgba(15, 23, 42, 0.95);
  backdrop-filter: blur(16px);
  border-radius: 1rem;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  width: 100%;
  max-width: 56rem;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
}

.ignore-modal__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.ignore-modal__icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.75rem;
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.2) 0%, rgba(139, 92, 246, 0.2) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c084fc;
}

.ignore-modal__close {
  padding: 0.5rem;
  border-radius: 0.5rem;
  color: #64748b;
  transition: all 150ms;
}

.ignore-modal__close:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #f1f5f9;
}

/* Tabs */
.ignore-modal__tabs {
  display: flex;
  gap: 0.25rem;
  padding: 0.75rem 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(0, 0, 0, 0.1);
}

.ignore-tab {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: #94a3b8;
  transition: all 150ms;
  position: relative;
}

.ignore-tab:hover {
  color: #e2e8f0;
  background: rgba(255, 255, 255, 0.03);
}

.ignore-tab--active {
  color: #c084fc;
  background: rgba(168, 85, 247, 0.1);
}

.ignore-tab--active::after {
  content: '';
  position: absolute;
  bottom: -0.75rem;
  left: 0.5rem;
  right: 0.5rem;
  height: 2px;
  background: linear-gradient(90deg, #a855f7, #ec4899);
  border-radius: 1px;
}

/* Content */
.ignore-modal__content {
  padding: 1.5rem;
  height: 400px;
  overflow: hidden;
}

.ignore-editor-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  height: 100%;
}

.ignore-custom-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  height: 100%;
}

.ignore-editor-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  height: 100%;
  min-height: 0;
}

.ignore-editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ignore-editor-label {
  font-size: 0.75rem;
  color: #64748b;
}

.ignore-editor-label strong {
  color: #e2e8f0;
}

/* Editor with line numbers */
.ignore-editor {
  flex: 1;
  min-height: 200px;
  display: flex;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 0.75rem;
  overflow: hidden;
  font-family: ui-monospace, 'SF Mono', 'Cascadia Code', monospace;
  font-size: 0.8125rem;
  line-height: 1.5;
}

.ignore-editor__lines {
  display: flex;
  flex-direction: column;
  padding: 0.75rem 0;
  background: rgba(0, 0, 0, 0.2);
  border-right: 1px solid rgba(255, 255, 255, 0.06);
  user-select: none;
  min-width: 2.5rem;
  overflow-y: auto;
  overflow-x: hidden;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.ignore-editor__lines::-webkit-scrollbar {
  display: none;
}

.ignore-editor__lines span {
  padding: 0 0.625rem;
  text-align: right;
  color: #475569;
  font-size: 0.75rem;
  height: 1.5em;
  flex-shrink: 0;
}

.ignore-editor__line--comment {
  color: #22c55e !important;
}

.ignore-editor__textarea {
  flex: 1;
  padding: 0.75rem;
  background: transparent;
  border: none;
  color: #e2e8f0;
  resize: none;
  outline: none;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  overflow-y: auto;
  overflow-x: hidden;
  white-space: pre;
}

.ignore-editor__textarea::placeholder {
  color: #475569;
}

.ignore-editor__textarea--editable:focus {
  background: rgba(168, 85, 247, 0.03);
}

.ignore-editor__hint {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.75rem;
  color: #64748b;
}



/* Footer */
.ignore-modal__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(0, 0, 0, 0.1);
}

.ignore-footer__left {
  display: flex;
  gap: 0.5rem;
}

.ignore-footer__right {
  display: flex;
  gap: 0.75rem;
}

.ignore-footer__danger-btn {
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.8125rem;
  color: #64748b;
  transition: all 150ms;
}

.ignore-footer__danger-btn:hover {
  color: #f87171;
  background: rgba(239, 68, 68, 0.1);
}

.ignore-footer__save-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1.25rem;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 0.625rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: white;
  transition: all 150ms;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.25);
}

.ignore-footer__save-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.35);
}

.ignore-footer__save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Transitions */
.modal-backdrop-enter-active, .modal-backdrop-leave-active { transition: opacity 0.2s ease; }
.modal-backdrop-enter-from, .modal-backdrop-leave-to { opacity: 0; }
.modal-enter-active, .modal-leave-active { transition: all 0.3s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; transform: scale(0.95); }
</style>
