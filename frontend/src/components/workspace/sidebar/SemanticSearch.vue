<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { apiService, type SemanticIndexStats, type SemanticSearchResult } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { computed, onMounted, ref } from 'vue'

const { t } = useI18n()
const projectStore = useProjectStore()

// State
const query = ref('')
const searchType = ref<'hybrid' | 'semantic' | 'keyword'>('hybrid')
const isSearching = ref(false)
const isIndexing = ref(false)
const results = ref<SemanticSearchResult[]>([])
const stats = ref<SemanticIndexStats | null>(null)
const isAvailable = ref(false)
const isIndexed = ref(false)
const error = ref('')

// Computed
const projectRoot = computed(() => projectStore.currentPath || '')
const hasResults = computed(() => results.value.length > 0)
const canSearch = computed(() => query.value.trim().length > 0 && isIndexed.value && !isSearching.value)

// Methods
async function checkAvailability() {
  try {
    isAvailable.value = await apiService.isSemanticSearchAvailable()
    if (isAvailable.value && projectRoot.value) {
      isIndexed.value = await apiService.semanticIsIndexed(projectRoot.value)
      if (isIndexed.value) {
        await loadStats()
      }
    }
  } catch (e) {
    console.error('Failed to check semantic search availability:', e)
    isAvailable.value = false
  }
}

async function loadStats() {
  if (!projectRoot.value) return
  try {
    stats.value = await apiService.semanticGetStats(projectRoot.value)
  } catch (e) {
    console.error('Failed to load stats:', e)
  }
}

async function indexProject() {
  if (!projectRoot.value || isIndexing.value) return
  
  isIndexing.value = true
  error.value = ''
  
  try {
    await apiService.semanticIndexProject(projectRoot.value)
    isIndexed.value = true
    await loadStats()
  } catch (e: unknown) {
    error.value = (e as Error).message || 'Failed to index project'
    console.error('Indexing failed:', e)
  } finally {
    isIndexing.value = false
  }
}

async function search() {
  if (!canSearch.value || !projectRoot.value) return
  
  isSearching.value = true
  error.value = ''
  results.value = []
  
  try {
    const response = await apiService.semanticSearch({
      query: query.value,
      projectRoot: projectRoot.value,
      topK: 10,
      minScore: 0.3,
      searchType: searchType.value
    })
    results.value = response.results
  } catch (e: unknown) {
    error.value = (e as Error).message || 'Search failed'
    console.error('Search failed:', e)
  } finally {
    isSearching.value = false
  }
}

function formatScore(score: number): string {
  return (score * 100).toFixed(1) + '%'
}

function getChunkTypeIcon(type: string): string {
  switch (type) {
    case 'function': return '⚡'
    case 'class': return '📦'
    case 'method': return '🔧'
    case 'file': return '📄'
    default: return '📝'
  }
}

function truncateContent(content: string, maxLength: number = 200): string {
  if (content.length <= maxLength) return content
  return content.substring(0, maxLength) + '...'
}

// Lifecycle
onMounted(() => {
  checkAvailability()
})

// Watch for project changes
projectStore.$subscribe(() => {
  checkAvailability()
})
</script>

<template>
  <div class="semantic-search">
    <!-- Header -->
    <div class="panel-header">
      <div class="flex items-center gap-2">
        <div class="panel-icon panel-icon-purple">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <span class="panel-title">{{ t('semanticSearch.title') }}</span>
      </div>
    </div>

    <!-- Not Available Warning -->
    <div v-if="!isAvailable" class="p-4">
      <div class="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-3 text-sm text-yellow-400">
        <p class="font-medium mb-1">{{ t('semanticSearch.notAvailable') }}</p>
        <p class="text-xs text-gray-400">{{ t('semanticSearch.configureApiKey') }}</p>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else class="p-4 space-y-4">
      <!-- Index Status -->
      <div v-if="!isIndexed" class="space-y-3">
        <div class="bg-blue-500/10 border border-blue-500/30 rounded-lg p-3 text-sm">
          <p class="text-blue-400 mb-2">{{ t('semanticSearch.notIndexed') }}</p>
          <button 
            @click="indexProject"
            :disabled="isIndexing"
            class="btn btn-primary btn-sm w-full"
          >
            <span v-if="isIndexing" class="flex items-center gap-2">
              <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ t('semanticSearch.indexing') }}
            </span>
            <span v-else>{{ t('semanticSearch.indexProject') }}</span>
          </button>
        </div>
      </div>

      <!-- Search Form -->
      <div v-else class="space-y-3">
        <!-- Stats -->
        <div v-if="stats" class="stats-grid grid-cols-3 text-xs">
          <div class="stat-card">
            <div class="stat-value">{{ stats.totalChunks }}</div>
            <div class="stat-label">{{ t('semanticSearch.chunks') }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.totalFiles }}</div>
            <div class="stat-label">{{ t('semanticSearch.files') }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ (stats.indexSize / 1024 / 1024).toFixed(1) }}MB</div>
            <div class="stat-label">{{ t('semanticSearch.size') }}</div>
          </div>
        </div>

        <!-- Search Input -->
        <div class="relative">
          <input
            v-model="query"
            type="text"
            :placeholder="t('semanticSearch.placeholder')"
            class="search-input w-full pr-10"
            @keyup.enter="search"
          />
          <button 
            @click="search"
            :disabled="!canSearch"
            class="absolute right-2 top-1/2 -translate-y-1/2 p-1 rounded hover:bg-gray-700/50 disabled:opacity-50"
          >
            <svg v-if="isSearching" class="animate-spin h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </button>
        </div>

        <!-- Search Type -->
        <div class="flex gap-1">
          <button 
            v-for="type in ['hybrid', 'semantic', 'keyword'] as const"
            :key="type"
            @click="searchType = type"
            :class="[
              'chip text-xs',
              searchType === type ? 'chip-active' : 'chip-default'
            ]"
          >
            {{ t(`semanticSearch.type.${type}`) }}
          </button>
        </div>

        <!-- Re-index button -->
        <button 
          @click="indexProject"
          :disabled="isIndexing"
          class="text-xs text-gray-400 hover:text-gray-400 flex items-center gap-1"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ t('semanticSearch.reindex') }}
        </button>
      </div>

      <!-- Error -->
      <div v-if="error" class="bg-red-500/10 border border-red-500/30 rounded-lg p-2 text-xs text-red-400">
        {{ error }}
      </div>

      <!-- Results -->
      <div v-if="hasResults" class="space-y-2">
        <div class="text-xs text-gray-400 mb-2">
          {{ t('semanticSearch.resultsCount').replace('{count}', String(results.length)) }}
        </div>
        
        <div 
          v-for="result in results" 
          :key="result.chunk.id"
          class="card p-3 space-y-2 cursor-pointer hover:border-purple-500/50 transition-colors"
        >
          <!-- Header -->
          <div class="flex items-start justify-between gap-2">
            <div class="flex items-center gap-2 min-w-0">
              <span class="text-sm">{{ getChunkTypeIcon(result.chunk.chunkType) }}</span>
              <span class="text-xs text-gray-300 truncate">
                {{ result.chunk.symbolName || result.chunk.filePath }}
              </span>
            </div>
            <span class="badge badge-primary text-xs shrink-0">
              {{ formatScore(result.score) }}
            </span>
          </div>
          
          <!-- File path -->
          <div class="text-xs text-gray-400 truncate">
            {{ result.chunk.filePath }}:{{ result.chunk.startLine }}-{{ result.chunk.endLine }}
          </div>
          
          <!-- Content preview -->
          <pre class="text-xs text-gray-400 bg-gray-800/50 rounded p-2 overflow-hidden whitespace-pre-wrap font-mono">{{ truncateContent(result.chunk.content) }}</pre>
        </div>
      </div>

      <!-- No Results -->
      <div v-else-if="query && !isSearching && isIndexed" class="text-center text-gray-400 text-sm py-4">
        {{ t('semanticSearch.noResults') }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.semantic-search {
  @apply bg-transparent;
}
</style>
