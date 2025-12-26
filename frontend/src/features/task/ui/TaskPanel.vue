<template>
  <div class="h-full flex flex-col bg-gray-800/50 backdrop-blur-sm p-4">
    <!-- Header with character counter -->
    <div class="flex items-center justify-between mb-3">
      <div class="flex items-center gap-2">
        <div class="w-1 h-5 bg-blue-500 rounded-full"></div>
        <h2 class="text-lg font-semibold text-white">{{ t('task.title') }}</h2>
      </div>
      <span class="text-xs text-gray-400" :class="{ 'text-yellow-400': taskStore.taskDescription.length > 4500, 'text-red-400': taskStore.taskDescription.length >= 5000 }">
        {{ taskStore.taskDescription.length }} / 5000
      </span>
    </div>
    
    <!-- Textarea with auto-resize -->
    <textarea
      ref="textareaRef"
      v-model="taskStore.taskDescription"
      :placeholder="t('task.placeholder')"
      class="task-textarea flex-1 w-full px-4 py-3 bg-gray-800/90 border border-gray-700/50 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500/50 focus:bg-gray-800 resize-none transition-all duration-200"
      maxlength="5000"
      @input="autoResize"
    ></textarea>

    <!-- AI Suggestions -->
    <div v-if="taskStore.suggestions.length > 0 && !taskStore.taskDescription" class="mt-3 space-y-2">
      <p class="text-xs text-gray-400 font-medium">{{ t('task.quickSuggestions') }}</p>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="suggestion in taskStore.suggestions"
          :key="suggestion"
          @click="taskStore.applySuggestion(suggestion)"
          class="chip chip-default"
        >
          {{ suggestion }}
        </button>
      </div>
    </div>

    <!-- Analysis Result -->
    <div v-if="taskStore.analysisResult" class="mt-3 p-3 bg-gray-700/50 rounded-lg text-sm">
      <p class="text-gray-300 font-medium mb-2">{{ t('task.analysisResult') }}</p>
      <div class="space-y-1 text-xs text-gray-400">
        <p>{{ t('task.complexity') }}: <span class="text-white">{{ taskStore.analysisResult.complexity }}</span></p>
        <p v-if="taskStore.analysisResult.estimatedTime">
          Estimated time: <span class="text-white">{{ taskStore.analysisResult.estimatedTime }}min</span>
        </p>
        <div v-if="taskStore.analysisResult.recommendations.length > 0" class="mt-2">
          <p class="font-medium text-gray-300 mb-1">Recommendations:</p>
          <ul class="list-disc list-inside space-y-0.5">
            <li v-for="rec in taskStore.analysisResult.recommendations" :key="rec">{{ rec }}</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Action buttons -->
    <div class="mt-4 flex gap-2">
      <button
        @click="handleAnalyze"
        :disabled="!canAnalyze"
        class="btn btn-primary flex-1 py-2"
      >
        <svg v-if="!taskStore.isAnalyzing" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
        <svg v-else class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 5.373 0 12h4z"></path>
        </svg>
        <span>{{ taskStore.isAnalyzing ? 'Analyzing...' : 'Analyze Task' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { ref, computed, onMounted } from 'vue'
import { useTaskStore } from '../model/task.store'
import { useUIStore } from '@/stores/ui.store'

const { t } = useI18n()
const taskStore = useTaskStore()
const uiStore = useUIStore()

const textareaRef = ref<HTMLTextAreaElement>()

const canAnalyze = computed(() => {
  return taskStore.taskDescription.trim().length > 10 && !taskStore.isAnalyzing
})

// Auto-resize textarea
function autoResize() {
  if (!textareaRef.value) return
  textareaRef.value.style.height = 'auto'
  textareaRef.value.style.height = textareaRef.value.scrollHeight + 'px'
}

async function handleAnalyze() {
  if (!canAnalyze.value) return
  
  try {
    await taskStore.analyzeTask()
    uiStore.addToast('Task analyzed successfully', 'success')
  } catch (error) {
    uiStore.addToast('Failed to analyze task', 'error')
  }
}

onMounted(() => {
  taskStore.loadTaskDraft()
  autoResize()
})
</script>

<style scoped>
.task-textarea {
  font-family: inherit;
  line-height: 1.6;
}

.task-textarea::placeholder {
  font-style: italic;
}

.task-textarea:focus {
  box-shadow: 0 0 0 3px rgba(107, 114, 128, 0.1);
}

.task-textarea::-webkit-scrollbar {
  width: 6px;
}

.task-textarea::-webkit-scrollbar-track {
  background: transparent;
}

.task-textarea::-webkit-scrollbar-thumb {
  background: rgba(107, 114, 128, 0.3);
  border-radius: 3px;
}

.task-textarea::-webkit-scrollbar-thumb:hover {
  background: rgba(107, 114, 128, 0.5);
}
</style>
