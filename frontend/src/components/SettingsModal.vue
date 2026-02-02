<template>
  <BaseModal
    :model-value="modelValue"
    size="lg"
    :close-on-backdrop="true"
    :close-on-esc="true"
    :show-close="false"
    @update:model-value="(val) => !val && close()"
  >
    <div class="settings-modal-wrapper">
      <!-- Header -->
      <div class="settings-header">
        <div class="flex items-center gap-2">
          <div class="settings-header-icon">
            <Settings class="w-4 h-4" />
          </div>
          <h2 class="text-base font-semibold text-white">{{ t('settings.modal.title') }}</h2>
        </div>
        <button @click="close" class="settings-close-btn">
          <X class="w-4 h-4" />
        </button>
      </div>

          <!-- Tabs -->
          <div class="settings-tabs">
            <div class="settings-tabs-container">
              <button
                v-for="tab in tabs"
                :key="tab.id"
                @click="activeTab = tab.id"
                class="settings-tab"
                :class="{ 'settings-tab-active': activeTab === tab.id }"
              >
                <component :is="tab.icon" class="w-4 h-4" />
                <span>{{ tab.label }}</span>
              </button>
            </div>
          </div>

          <!-- Tab Content -->
          <div class="settings-content">
            <Transition name="tab-fade" mode="out-in">
              <!-- General Tab -->
              <div v-if="activeTab === 'general'" key="general" class="settings-section">
              <div class="settings-group">
                <div class="settings-group-header">
                  <Globe class="w-4 h-4 text-indigo-400" />
                  <span>{{ t('settings.modal.language') }}</span>
                </div>
                <select v-model="selectedLanguage" class="settings-select">
                  <option value="ru">Русский</option>
                  <option value="en">English</option>
                </select>
              </div>

              <div class="settings-divider"></div>

              <div class="settings-group">
                <div class="settings-group-header">
                  <Maximize class="w-4 h-4 text-emerald-400" />
                  <span>{{ t('settings.modal.uiScaling') }}</span>
                </div>
                <div class="flex items-center gap-4">
                  <input
                    type="range"
                    v-model.number="settingsStore.settings.uiScale"
                    min="0.75"
                    max="1.25"
                    step="0.05"
                    class="settings-slider"
                  />
                  <span class="settings-scale-value">{{ Math.round(settingsStore.settings.uiScale * 100) }}%</span>
                </div>
                <p class="settings-hint">{{ t('settings.modal.uiScalingHint') }}</p>
              </div>

              <div class="settings-divider"></div>


              <div class="settings-group">
                <div class="settings-group-header">
                  <HelpCircle class="w-4 h-4 text-purple-400" />
                  <span>{{ t('onboarding.startTour') }}</span>
                </div>
                <p class="settings-hint">{{ t('settings.modal.tourHint') }}</p>
                <button @click="handleStartTour" class="settings-action-btn">
                  <Sparkles class="w-4 h-4" />
                  {{ t('onboarding.startTour') }}
                </button>
              </div>
            </div>

            <!-- AI Tab -->
              <div v-else-if="activeTab === 'ai'" key="ai" class="settings-section">
                <AISettings />
              </div>

              <!-- Export Tab -->
              <div v-else-if="activeTab === 'export'" key="export" class="settings-section">
                <ExportSettings />
              </div>

              <!-- File Explorer Tab -->
              <div v-else-if="activeTab === 'fileExplorer'" key="fileExplorer" class="settings-section">
              <div class="settings-group">
                <div class="settings-group-header">
                  <Filter class="w-4 h-4 text-emerald-400" />
                  <span>{{ t('settings.modal.filterSettings') }}</span>
                </div>
                
                <div class="settings-toggle-list">
                  <label class="settings-toggle-item">
                    <div class="settings-toggle-info">
                      <span class="settings-toggle-label">{{ t('settings.useGitignore') }}</span>
                      <span class="settings-toggle-hint">{{ t('settings.useGitignoreHint') }}</span>
                    </div>
                    <div class="settings-toggle">
                      <input
                        type="checkbox"
                        v-model="settingsStore.settings.fileExplorer.useGitignore"
                        class="sr-only peer"
                      />
                      <div class="settings-toggle-track peer-checked:bg-indigo-500"></div>
                      <div class="settings-toggle-thumb peer-checked:translate-x-5"></div>
                    </div>
                  </label>

                  <label class="settings-toggle-item">
                    <div class="settings-toggle-info">
                      <span class="settings-toggle-label">{{ t('settings.useCustomIgnore') }}</span>
                      <span class="settings-toggle-hint">{{ t('settings.useCustomIgnoreHint') }}</span>
                    </div>
                    <div class="settings-toggle">
                      <input
                        type="checkbox"
                        v-model="settingsStore.settings.fileExplorer.useCustomIgnore"
                        class="sr-only peer"
                      />
                      <div class="settings-toggle-track peer-checked:bg-indigo-500"></div>
                      <div class="settings-toggle-thumb peer-checked:translate-x-5"></div>
                    </div>
                  </label>
                </div>
              </div>

              <div class="settings-divider"></div>

              <div class="settings-group">
                <div class="settings-group-header">
                  <FolderTree class="w-4 h-4 text-orange-400" />
                  <span>{{ t('settings.modal.displaySettings') }}</span>
                </div>
                
                <div class="settings-toggle-list">
                  <label class="settings-toggle-item">
                    <div class="settings-toggle-info">
                      <span class="settings-toggle-label">{{ t('settings.autoSaveSelection') }}</span>
                      <span class="settings-toggle-hint">{{ t('settings.autoSaveSelectionHint') }}</span>
                    </div>
                    <div class="settings-toggle">
                      <input
                        type="checkbox"
                        v-model="settingsStore.settings.fileExplorer.autoSaveSelection"
                        class="sr-only peer"
                      />
                      <div class="settings-toggle-track peer-checked:bg-indigo-500"></div>
                      <div class="settings-toggle-thumb peer-checked:translate-x-5"></div>
                    </div>
                  </label>

                  <label class="settings-toggle-item">
                    <div class="settings-toggle-info">
                      <span class="settings-toggle-label">{{ t('settings.compactFolders') }}</span>
                      <span class="settings-toggle-hint">{{ t('settings.compactFoldersHint') }}</span>
                    </div>
                    <div class="settings-toggle">
                      <input
                        type="checkbox"
                        v-model="settingsStore.settings.fileExplorer.compactNestedFolders"
                        class="sr-only peer"
                      />
                      <div class="settings-toggle-track peer-checked:bg-indigo-500"></div>
                      <div class="settings-toggle-thumb peer-checked:translate-x-5"></div>
                    </div>
                  </label>
                </div>
              </div>

              <div class="settings-divider"></div>

              <!-- Dependency Settings -->
              <div class="settings-group">
                <div class="settings-group-header">
                  <Share2 class="w-4 h-4 text-purple-400" />
                  <span>{{ t('settings.dependencies.title') }}</span>
                </div>
                
                <div class="settings-toggle-list">
                  <label class="settings-toggle-item">
                    <div class="settings-toggle-info">
                      <span class="settings-toggle-label">{{ t('settings.dependencies.showIndicators') }}</span>
                      <span class="settings-toggle-hint">{{ t('settings.dependencies.showIndicatorsDesc') }}</span>
                    </div>
                    <div class="settings-toggle">
                      <input
                        type="checkbox"
                        v-model="fileStore.showDependencyIndicators"
                        class="sr-only peer"
                        @change="fileStore.setShowDependencyIndicators(fileStore.showDependencyIndicators)"
                      />
                      <div class="settings-toggle-track peer-checked:bg-indigo-500"></div>
                      <div class="settings-toggle-thumb peer-checked:translate-x-5"></div>
                    </div>
                  </label>

                  <label class="settings-toggle-item" :class="{ 'opacity-50': !fileStore.showDependencyIndicators }">
                    <div class="settings-toggle-info">
                      <span class="settings-toggle-label">{{ t('settings.dependencies.autoHighlight') }}</span>
                      <span class="settings-toggle-hint">{{ t('settings.dependencies.autoHighlightDesc') }}</span>
                    </div>
                    <div class="settings-toggle">
                      <input
                        type="checkbox"
                        v-model="fileStore.autoHighlightDependencies"
                        :disabled="!fileStore.showDependencyIndicators"
                        class="sr-only peer"
                        @change="fileStore.setAutoHighlightDependencies(fileStore.autoHighlightDependencies)"
                      />
                      <div class="settings-toggle-track peer-checked:bg-indigo-500"></div>
                      <div class="settings-toggle-thumb peer-checked:translate-x-5"></div>
                    </div>
                  </label>
                </div>

                <div class="settings-subsection">
                  <label class="settings-subsection-label">
                    {{ t('settings.dependencies.includeTypes') }}
                  </label>
                  <p class="settings-subsection-hint">
                    {{ t('settings.dependencies.includeTypesDesc') }}
                  </p>
                  <div class="checkbox-group">
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        v-model="fileStore.includeDependencyTests"
                        class="checkbox"
                        @change="fileStore.setIncludeDependencyTests(fileStore.includeDependencyTests)"
                      />
                      <span>{{ t('settings.dependencies.includeTests') }}</span>
                    </label>
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        v-model="fileStore.includeDependencyStyles"
                        class="checkbox"
                        @change="fileStore.setIncludeDependencyStyles(fileStore.includeDependencyStyles)"
                      />
                      <span>{{ t('settings.dependencies.includeStyles') }}</span>
                    </label>
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        v-model="fileStore.includeDependencyTypes"
                        class="checkbox"
                        @change="fileStore.setIncludeDependencyTypes(fileStore.includeDependencyTypes)"
                      />
                      <span>{{ t('settings.dependencies.includeTypesLabel') }}</span>
                    </label>
                  </div>
                </div>
              </div>
            </div>

            <!-- System Tab -->
              <div v-else-if="activeTab === 'system'" key="system" class="settings-section">
                <ShellIntegrationSettings />
              </div>
            </Transition>
          </div>

      <!-- Footer -->
      <div class="settings-footer">
        <button @click="close" class="settings-close-action">
          {{ t('settings.modal.cancel') }}
        </button>
      </div>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import BaseModal from '@/components/ui/BaseModal.vue'
import AISettings from '@/components/workspace/sidebar/AISettings.vue'
import ExportSettings from '@/components/workspace/sidebar/ExportSettings.vue'
import ShellIntegrationSettings from '@/components/ShellIntegrationSettings.vue'
import { useI18n } from '@/composables/useI18n'
import { useOnboarding } from '@/composables/useOnboarding'
import { useSettingsStore } from '@/stores/settings.store'
import { useFileStore } from '@/features/files/model/file.store'
import { FileText, Filter, FolderTree, Globe, HelpCircle, Lightbulb, Maximize, Monitor, Settings, Share2, Sparkles, X } from 'lucide-vue-next'

import { computed, ref, watch } from 'vue'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const { t, setLocale, locale } = useI18n()
const settingsStore = useSettingsStore()
const fileStore = useFileStore()
const { startTour, resetTour } = useOnboarding()

const tabs = computed(() => [
  { id: 'general', label: t('settings.modal.general'), icon: Settings },
  { id: 'ai', label: t('settings.modal.ai'), icon: Lightbulb },
  { id: 'export', label: t('settings.modal.export'), icon: FileText },
  { id: 'fileExplorer', label: t('settings.modal.fileExplorer'), icon: FolderTree },
  { id: 'system', label: t('settings.modal.system'), icon: Monitor },
])

const activeTab = ref<string>('general')
const selectedLanguage = ref(locale.value)

// Sync language selection with locale
watch(
  () => props.modelValue,
  (isOpen) => {
    if (isOpen) {
      selectedLanguage.value = locale.value
      activeTab.value = 'general'
    }
  }
)

// Apply language change immediately
watch(selectedLanguage, (newLang) => {
  if (newLang !== locale.value) {
    setLocale(newLang as 'ru' | 'en')
  }
})

function close() {
  emit('update:modelValue', false)
}

function handleStartTour() {
  close()
  resetTour()
  setTimeout(() => {
    startTour()
  }, 300)
}
</script>

<style scoped>
.settings-modal-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  max-height: 70vh;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: var(--bg-1);
  border-bottom: 1px solid var(--border-default);
}

.settings-header-icon {
  width: 1.75rem;
  height: 1.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent-indigo-bg);
  border: 1px solid var(--accent-indigo-border);
  border-radius: var(--radius-md);
  color: #a5b4fc;
}

.settings-header-icon svg {
  width: 1rem;
  height: 1rem;
}

.settings-close-btn {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 150ms ease-out;
}

.settings-close-btn:hover {
  background: var(--bg-3);
  color: var(--text-primary);
}

.settings-tabs {
  padding: 0.5rem 1rem;
  background: var(--bg-0);
  border-bottom: 1px solid var(--border-subtle);
}

.settings-tabs-container {
  display: flex;
  gap: 0.25rem;
  padding: 0.1875rem;
  background: var(--bg-1);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-subtle);
}

.settings-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
  padding: 0.5rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 200ms ease-out;
  white-space: nowrap;
}

.settings-tab:hover:not(.settings-tab-active) {
  color: var(--text-secondary);
  background: var(--bg-2);
}

.settings-tab-active {
  color: white;
  background: var(--accent-indigo-bg);
  border-color: var(--accent-indigo-border);
  font-weight: 600;
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 1.25rem;
  min-height: 0;
}

.settings-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.settings-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.settings-group-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.settings-hint {
  font-size: 0.75rem;
  color: var(--text-muted);
  line-height: 1.5;
}

.settings-select {
  width: 100%;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 150ms ease-out;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%239ca3af' stroke-width='2'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
}

.settings-select:hover {
  background-color: var(--bg-3);
  border-color: var(--border-strong);
}

.settings-select:focus {
  outline: none;
  border-color: var(--accent-indigo);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.settings-slider {
  flex: 1;
  height: 6px;
  background: var(--bg-3);
  border-radius: 3px;
  appearance: none;
  outline: none;
}

.settings-slider::-webkit-slider-thumb {
  appearance: none;
  width: 18px;
  height: 18px;
  background: var(--accent-indigo);
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.15s ease-out;
  border: 2px solid white;
}

.settings-slider::-webkit-slider-thumb:hover {
  transform: scale(1.1);
  box-shadow: 0 0 10px rgba(99, 102, 241, 0.4);
}

.settings-scale-value {
  min-width: 3rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--accent-indigo);
  font-family: var(--font-mono);
}

.settings-action-btn {

  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: white;
  background: var(--gradient-primary);
  border: none;
  border-radius: var(--radius-xl);
  cursor: pointer;
  transition: all 200ms ease-out;
  box-shadow: var(--shadow-glow-accent);
}

.settings-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(139, 92, 246, 0.4);
}

.settings-divider {
  height: 1px;
  background: var(--border-subtle);
  margin: 0.5rem 0;
}

.settings-toggle-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.settings-toggle-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-xl);
  cursor: pointer;
  transition: all 150ms ease-out;
}

.settings-toggle-item:hover {
  background: var(--bg-3);
  border-color: var(--border-default);
}

.settings-toggle-info {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.settings-toggle-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.settings-toggle-hint {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.settings-toggle {
  position: relative;
  width: 2.75rem;
  height: 1.5rem;
  flex-shrink: 0;
}

.settings-toggle-track {
  position: absolute;
  inset: 0;
  background: var(--bg-0);
  border: 1px solid var(--border-strong);
  border-radius: 9999px;
  transition: all 200ms ease-out;
}

.settings-toggle-thumb {
  position: absolute;
  top: 0.125rem;
  left: 0.125rem;
  width: 1.25rem;
  height: 1.25rem;
  background: var(--text-muted);
  border-radius: 9999px;
  transition: all 200ms ease-out;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.settings-toggle input:checked ~ .settings-toggle-track {
  background: var(--accent-indigo);
  border-color: var(--accent-indigo);
}

.settings-toggle input:checked ~ .settings-toggle-thumb {
  transform: translateX(1.25rem);
  background: white;
}

.settings-footer {
  display: flex;
  justify-content: flex-end;
  padding: 0.5rem 1rem;
  background: var(--bg-1);
  border-top: 1px solid var(--border-subtle);
}

.settings-close-action {
  padding: 0.5rem 1rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all 150ms ease-out;
}

.settings-close-action:hover {
  color: var(--text-primary);
  background: var(--bg-2);
}

/* Tab Content Animation */
.tab-fade-enter-active,
.tab-fade-leave-active {
  transition: opacity 0.15s ease;
}

.tab-fade-enter-from,
.tab-fade-leave-to {
  opacity: 0;
}

/* Dependency Settings Styles */
.settings-subsection {
  margin-top: 1rem;
  padding: 1rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
}

.settings-subsection-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.25rem;
}

.settings-subsection-hint {
  font-size: 0.75rem;
  color: var(--text-muted);
  line-height: 1.5;
  margin-bottom: 0.75rem;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 8px;
  border-radius: var(--radius-sm);
  transition: background 150ms ease-out;
}

.checkbox-label:hover {
  background: var(--bg-hover);
}

.checkbox {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  border: 2px solid var(--border-default);
  cursor: pointer;
  transition: all 150ms ease-out;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  background: transparent;
}

.checkbox:checked {
  background: var(--accent-indigo);
  border-color: var(--accent-indigo);
  position: relative;
}

.checkbox:checked::after {
  content: '';
  position: absolute;
  left: 4px;
  top: 1px;
  width: 4px;
  height: 8px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.checkbox:hover {
  border-color: var(--accent-indigo);
}

.checkbox-label span {
  font-size: 14px;
  color: var(--text-primary);
}
</style>
