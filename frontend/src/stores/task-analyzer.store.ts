import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export interface SuggestedFile {
  path: string
  reason: string
  confidence: number
}

export const useTaskAnalyzerStore = defineStore('taskAnalyzer', () => {
  // State
  const taskDescription = ref('')
  const suggestions = ref<SuggestedFile[]>([])
  const isAnalyzing = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const hasSuggestions = computed(() => suggestions.value.length > 0)
  const highConfidenceSuggestions = computed(() =>
    suggestions.value.filter(s => s.confidence >= 0.7)
  )

  // Actions
  async function analyzeTask(task: string) {
    if (!task.trim()) {
      suggestions.value = []
      return
    }

    taskDescription.value = task
    isAnalyzing.value = true
    error.value = null

    try {
      // TODO [Q2]: Call backend API to analyze task - see TODO-refactoring.txt
      // Backend API for task analysis not yet implemented
      suggestions.value = []
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Analysis failed'
    } finally {
      isAnalyzing.value = false
    }
  }

  function setSuggestions(newSuggestions: SuggestedFile[]) {
    suggestions.value = newSuggestions
  }

  function clearSuggestions() {
    suggestions.value = []
    taskDescription.value = ''
    error.value = null
  }

  return {
    // State
    taskDescription,
    suggestions,
    isAnalyzing,
    error,
    // Computed
    hasSuggestions,
    highConfidenceSuggestions,
    // Actions
    analyzeTask,
    setSuggestions,
    clearSuggestions
  }
})
