<template>
  <div class="info-box space-y-4">
    <div class="flex items-center justify-between">
      <h4 class="text-sm font-semibold text-white">{{ t('context.buildOptions') }}</h4>
      <button
        @click="resetToDefaults"
        class="text-xs text-gray-400 hover:text-white transition-colors px-2 py-1 rounded-lg hover:bg-gray-700/50"
      >
        {{ t('context.reset') }}
      </button>
    </div>
    
    <div class="space-y-3">
      <div>
        <label class="block text-xs text-gray-400 mb-1.5">{{ t('context.maxTokens') }}</label>
        <input
          v-model.number="localOptions.maxTokens"
          type="number"
          min="1000"
          max="10000000"
          step="10000"
          class="input text-sm"
          @change="updateStoreOptions"
        />
        <p class="text-xs text-gray-400 mt-1.5">{{ t('context.default') }}: {{ formatNumber(settingsStore.settings.context.maxTokens) }}</p>
      </div>

      <div class="flex items-center justify-between" :title="t('context.stripCommentsTooltip')">
        <div class="flex items-center gap-2 flex-1">
          <input
            v-model="localOptions.stripComments"
            type="checkbox"
            id="stripComments"
            class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500"
            @change="updateStoreOptions"
          />
          <label for="stripComments" class="text-xs text-gray-300">{{ t('context.stripComments') }}</label>
        </div>
        <div class="text-green-500 text-xs flex items-center" :title="t('context.working')">
          <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 00 16zm3.707-9.293a1 1 0 0-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
        </div>
      </div>

      <div class="flex items-center justify-between" :title="t('context.includeTestsTooltip')">
        <div class="flex items-center gap-2 flex-1">
          <input
            v-model="localOptions.includeTests"
            type="checkbox"
            id="includeTests"
            class="w-4 h-4 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500"
            @change="updateStoreOptions"
          />
          <label for="includeTests" class="text-xs text-gray-300">{{ t('context.includeTests') }}</label>
        </div>
        <div class="text-yellow-500 text-xs flex items-center" :title="t('context.pendingBackend')">
          <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 0-1-1H9z" clip-rule="evenodd" />
          </svg>
        </div>
      </div>

      <div class="space-y-1.5" :title="t('context.splitStrategyTooltip')">
        <div class="flex items-center justify-between">
          <label class="text-xs text-gray-400">{{ t('context.splitStrategy') }}</label>
          <span class="badge badge-warning text-[10px]">Soon</span>
        </div>
        <select
          v-model="localOptions.splitStrategy"
          class="input text-sm"
          @change="updateStoreOptions"
        >
          <option value="semantic">{{ t('context.semantic') }} ({{ t('context.semanticHint') }})</option>
          <option value="fixed">{{ t('context.fixed') }} ({{ t('context.fixedHint') }})</option>
          <option value="adaptive">{{ t('context.adaptive') }} ({{ t('context.adaptiveHint') }})</option>
        </select>
      </div>

      <div class="space-y-1.5" :title="t('context.outputFormatTooltip')">
        <div class="flex items-center justify-between">
          <label class="text-xs text-gray-400">{{ t('context.outputFormat') }}</label>
          <span class="badge badge-success text-[10px]">✓</span>
        </div>
        <select
          v-model="localOptions.outputFormat"
          class="input text-sm"
          @change="updateStoreOptions"
        >
          <option value="markdown">Markdown (```code blocks```)</option>
          <option value="xml">XML (&lt;file&gt;&lt;content&gt;)</option>
          <option value="json">JSON (structured)</option>
          <option value="plain">Plain Text (--- separators)</option>
        </select>
      </div>
    </div>

    <div v-if="fileStore.selectedCount > 0" class="text-xs pt-3 border-t border-gray-700/30">
      <p class="text-gray-400">{{ fileStore.selectedCount }} {{ t('context.filesSelected') }}</p>
      <p class="text-indigo-400 mt-1">{{ t('context.estimated') }}: ~{{ estimatedTokens }} {{ t('context.tokens') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useFileStore } from '@/features/files'
import { useSettingsStore } from '@/stores/settings.store'
import { computed, onMounted, ref } from 'vue'

const fileStore = useFileStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()

// Expose options for parent component to use
const localOptions = ref({
  maxTokens: settingsStore.settings.context.maxTokens,
  stripComments: settingsStore.settings.context.stripComments,
  includeTests: settingsStore.settings.context.includeTests,
  splitStrategy: settingsStore.settings.context.splitStrategy,
  outputFormat: settingsStore.settings.context.outputFormat
})

const estimatedTokens = computed(() => {
  return fileStore.selectedCount * 500 // Rough estimate
})

function resetToDefaults() {
  localOptions.value = {
    maxTokens: settingsStore.settings.context.maxTokens,
    stripComments: settingsStore.settings.context.stripComments,
    includeTests: settingsStore.settings.context.includeTests,
    splitStrategy: settingsStore.settings.context.splitStrategy,
    outputFormat: settingsStore.settings.context.outputFormat
  }
  updateStoreOptions()
}

function updateStoreOptions() {
  // Save to both localStorage AND settings store for persistence
  try {
    localStorage.setItem('context-build-options', JSON.stringify(localOptions.value))
    
    // Also update settings store
    settingsStore.updateContextSettings({
      maxTokens: localOptions.value.maxTokens || settingsStore.settings.context.maxTokens,
      stripComments: localOptions.value.stripComments ?? settingsStore.settings.context.stripComments,
      includeTests: localOptions.value.includeTests ?? settingsStore.settings.context.includeTests,
      splitStrategy: localOptions.value.splitStrategy || settingsStore.settings.context.splitStrategy,
      outputFormat: localOptions.value.outputFormat || settingsStore.settings.context.outputFormat
    })
  } catch (err) {
    console.warn('Failed to save build options:', err)
  }
}

// Load saved options on mount
onMounted(() => {
  try {
    const saved = localStorage.getItem('context-build-options')
    if (saved) {
      const parsed = JSON.parse(saved)
      localOptions.value = {
        maxTokens: parsed.maxTokens || settingsStore.settings.context.maxTokens,
        stripComments: parsed.stripComments ?? settingsStore.settings.context.stripComments,
        includeTests: parsed.includeTests ?? settingsStore.settings.context.includeTests,
        splitStrategy: parsed.splitStrategy || settingsStore.settings.context.splitStrategy,
        outputFormat: parsed.outputFormat || settingsStore.settings.context.outputFormat
      }
    }
  } catch (err) {
    console.warn('Failed to load build options:', err)
    // Fallback to settings store
    localOptions.value = {
      maxTokens: settingsStore.settings.context.maxTokens,
      stripComments: settingsStore.settings.context.stripComments,
      includeTests: settingsStore.settings.context.includeTests,
      splitStrategy: settingsStore.settings.context.splitStrategy,
      outputFormat: settingsStore.settings.context.outputFormat
    }
  }
})

// Helper function
function formatNumber(num: number): string {
  if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`
  if (num >= 1000) return `${(num / 1000).toFixed(0)}K`
  return num.toString()
}

// Expose options to parent
defineExpose({
  options: localOptions
})
</script>
