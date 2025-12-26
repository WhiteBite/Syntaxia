<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
        @click.self="close" @keydown.esc="close">
        <div
          class="bg-gray-800/95 backdrop-blur-md rounded-2xl shadow-2xl w-full max-w-3xl max-h-[90vh] flex flex-col border border-gray-700/50"
          @click.stop>
          <!-- Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-gray-700/50">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-indigo-500/20 flex items-center justify-center">
                <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                </svg>
              </div>
              <div>
                <h3 class="text-lg font-semibold text-white">Экспорт контекста</h3>
                <p class="text-xs text-gray-400">Выберите формат и параметры экспорта</p>
              </div>
            </div>
            <button @click="close" class="p-2 hover:bg-gray-700/50 rounded-xl transition-colors" aria-label="Закрыть">
              <svg class="w-5 h-5 text-gray-400 hover:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Content -->
          <div class="flex-1 overflow-y-auto p-6">
            <!-- Export Mode Selection -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-300 mb-3">Режим экспорта</label>
              <div class="grid grid-cols-3 gap-3">
                <button v-for="mode in exportModes" :key="mode.value" @click="selectedMode = mode.value" :class="[
                  'p-4 rounded-xl border transition-all text-left',
                  selectedMode === mode.value
                    ? 'border-indigo-500/50 bg-indigo-500/10 shadow-lg shadow-indigo-500/10'
                    : 'border-gray-700/50 hover:border-gray-600/50 bg-gray-800/50 hover:bg-gray-700/50'
                ]">
                  <div class="flex items-start gap-3">
                    <component :is="mode.icon" class="w-5 h-5 flex-shrink-0"
                      :class="selectedMode === mode.value ? 'text-indigo-400' : 'text-gray-400'" />
                    <div class="flex-1 min-w-0">
                      <div class="text-sm font-medium text-white mb-1">{{ mode.label }}</div>
                      <div class="text-xs text-gray-400">{{ mode.description }}</div>
                    </div>
                  </div>
                </button>
              </div>
            </div>

            <!-- Clipboard Mode Settings -->
            <div v-if="selectedMode === 'clipboard'" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Формат</label>
                <select v-model="settings.exportFormat" class="input">
                  <option value="plain">Plain Text</option>
                  <option value="manifest">With Manifest</option>
                  <option value="json">JSON</option>
                </select>
              </div>

              <label class="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
                <input v-model="settings.stripComments" type="checkbox"
                  class="w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0" />
                Удалить комментарии
              </label>

              <label class="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
                <input v-model="settings.includeManifest" type="checkbox"
                  class="w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0" />
                Включить манифест файлов
              </label>
            </div>

            <!-- AI Mode Settings -->
            <div v-if="selectedMode === 'ai'" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">AI Профиль</label>
                <select v-model="settings.aiProfile" class="input">
                  <option value="gpt-4">GPT-4</option>
                  <option value="gpt-3.5">GPT-3.5 Turbo</option>
                  <option value="claude">Claude</option>
                  <option value="gemini">Gemini Pro</option>
                </select>
              </div>

              <div>
                <label class="flex items-center justify-between text-sm font-medium text-gray-300 mb-2">
                  <span>Лимит токенов: {{ formatNumber(settings.tokenLimit) }}</span>
                </label>
                <input v-model.number="settings.tokenLimit" type="range" min="1000" max="128000" step="1000"
                  class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-blue-500" />
              </div>

              <label class="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
                <input v-model="settings.enableAutoSplit" type="checkbox"
                  class="w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0" />
                Автоматически разбивать на чанки
              </label>

              <div v-if="settings.enableAutoSplit"
                class="ml-6 space-y-3 p-3 bg-gray-800/50 rounded-xl border border-gray-700/30">
                <div>
                  <label class="block text-xs text-gray-400 mb-1">Токенов на чанк</label>
                  <input v-model.number="settings.maxTokensPerChunk" type="number" min="500" max="32000"
                    class="input text-sm py-1.5" />
                </div>
                <div>
                  <label class="block text-xs text-gray-400 mb-1">Перекрытие токенов</label>
                  <input v-model.number="settings.overlapTokens" type="number" min="0" max="1000"
                    class="input text-sm py-1.5" />
                </div>
                <div>
                  <label class="block text-xs text-gray-400 mb-1">Стратегия разбиения</label>
                  <select v-model="settings.splitStrategy" class="input text-sm py-1.5">
                    <option value="smart">Smart (рекомендуется)</option>
                    <option value="file">По файлам</option>
                    <option value="token">По токенам</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Human Mode Settings -->
            <div v-if="selectedMode === 'human'" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">Тема оформления</label>
                <select v-model="settings.theme" class="input">
                  <option value="default">По умолчанию</option>
                  <option value="dark">Тёмная</option>
                  <option value="light">Светлая</option>
                  <option value="minimal">Минимальная</option>
                </select>
              </div>

              <label class="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
                <input v-model="settings.includeLineNumbers" type="checkbox"
                  class="w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0" />
                Показывать номера строк
              </label>

              <label class="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
                <input v-model="settings.includePageNumbers" type="checkbox"
                  class="w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0" />
                Показывать номера страниц
              </label>
            </div>

            <!-- Context Preview -->
            <div v-if="contextStore.hasContext" class="mt-6 p-4 bg-gray-800/50 rounded-xl border border-gray-700/30">
              <div class="text-sm font-medium text-gray-300 mb-3">Предпросмотр контекста</div>
              <div class="grid grid-cols-3 gap-3">
                <div class="stat-card">
                  <div class="stat-card-value">{{ contextStore.fileCount }}</div>
                  <div class="stat-label">Файлов</div>
                </div>
                <div class="stat-card">
                  <div class="stat-card-value">{{ formatNumber(contextStore.lineCount) }}</div>
                  <div class="stat-label">Строк</div>
                </div>
                <div class="stat-card">
                  <div class="stat-card-value stat-card-value-indigo">{{ formatNumber(contextStore.tokenCount) }}</div>
                  <div class="stat-label">Токенов</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="flex items-center justify-between px-6 py-4 border-t border-gray-700/50">
            <div class="text-xs text-gray-400">
              <span v-if="exportResult">
                Экспорт завершен успешно
              </span>
              <span v-else-if="isExporting">
                Экспорт в процессе...
              </span>
            </div>
            <div class="flex items-center gap-3">
              <button @click="close" class="btn btn-secondary">
                Отмена
              </button>
              <button @click="handleExport" :disabled="isExporting || !contextStore.hasContext"
                class="btn btn-primary flex items-center gap-2">
                <svg v-if="isExporting" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                  </path>
                </svg>
                {{ isExporting ? 'Экспорт...' : 'Экспортировать' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useExport } from '@/composables/useExport'
import { h } from 'vue'

import { useContextStore } from '@/features/context/model/context.store'

const exportComposable = useExport()
const contextStore = useContextStore() // Get contextStore directly from the store

// Create local refs that sync with the composable
const isOpen = exportComposable.isOpen
const isExporting = exportComposable.isExporting
const exportResult = exportComposable.exportResult
type ExportMode = 'clipboard' | 'ai' | 'human'
const selectedMode = exportComposable.selectedMode

// Sync settings with the composable
const settings = exportComposable.settings

const exportModes: { value: ExportMode; label: string; description: string; icon: ReturnType<typeof h> }[] = [
  {
    value: 'clipboard',
    label: 'Буфер обмена',
    description: 'Текстовый экспорт',
    icon: h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2' })
    ])
  },
  {
    value: 'ai',
    label: 'AI оптимизация',
    description: 'С разбиением на чанки',
    icon: h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z' })
    ])
  },
  {
    value: 'human',
    label: 'PDF/ZIP',
    description: 'Для печати',
    icon: h('svg', { class: 'w-5 h-5', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z' })
    ])
  }
]

function close() {
  exportComposable.close()
}

async function handleExport() {
  if (!contextStore.hasContext) return

  // Execute export using the composable
  await exportComposable.executeExport()
}

function formatNumber(num: number): string {
  if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`
  if (num >= 1000) return `${(num / 1000).toFixed(1)}K`
  return num.toString()
}

defineExpose({
  open: exportComposable.open,
  close: exportComposable.close
})
</script>

<style scoped>
.modal-enter-active>div,
.modal-leave-active>div {
  transition: transform 0.2s ease;
}

.modal-enter-from>div,
.modal-leave-to>div {
  transform: scale(0.95);
}
</style>
