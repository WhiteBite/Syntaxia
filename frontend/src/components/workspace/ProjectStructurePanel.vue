<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useProjectStore } from '@/stores/project.store'
import type { ProjectStructure } from '@/types/dto'
import { onMounted, ref, watch } from 'vue'
import { GetProjectStructure } from '../../../wailsjs/go/main/App'

const { t } = useI18n()
const projectStore = useProjectStore()

const structure = ref<ProjectStructure | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const activeTab = ref<'architecture' | 'frameworks' | 'conventions' | 'languages'>('architecture')

async function loadStructure() {
  if (!projectStore.currentPath) return

  loading.value = true
  error.value = null

  try {
    structure.value = await GetProjectStructure(projectStore.currentPath) as unknown as ProjectStructure
  } catch (e) {
    error.value = String(e)
    console.error('Failed to load project structure:', e)
  } finally {
    loading.value = false
  }
}

function getArchitectureIcon(type: string): string {
  const icons: Record<string, string> = {
    'clean': '🏛️',
    'hexagonal': '⬡',
    'mvc': '📐',
    'mvvm': '🔗',
    'layered': '📚',
    'ddd': '🎯',
    'microservices': '🔌',
    'monolith': '🏢',
    'unknown': '❓'
  }
  return icons[type] || '📁'
}

function getConfidenceColor(confidence: number): string {
  if (confidence >= 0.7) return 'text-emerald-400'
  if (confidence >= 0.4) return 'text-yellow-400'
  return 'text-gray-400'
}

function getCategoryIcon(category: string): string {
  const icons: Record<string, string> = {
    'web': '🌐',
    'mobile': '📱',
    'desktop': '🖥️',
    'testing': '🧪',
    'cli': '⌨️'
  }
  return icons[category] || '📦'
}

onMounted(() => {
  if (projectStore.currentPath) {
    loadStructure()
  }
})

watch(() => projectStore.currentPath, (newDir) => {
  if (newDir) {
    loadStructure()
  }
})
</script>

<template>
  <div class="project-structure-panel">
    <!-- Header -->
    <div class="panel-header">
      <div class="flex items-center gap-2">
        <div class="panel-icon panel-icon-indigo">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <span class="panel-title">{{ t('projectStructure.title') }}</span>
      </div>
      <button @click="loadStructure" class="action-btn" :disabled="loading" :title="t('projectStructure.refresh')">
        <svg class="w-4 h-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="p-4 text-center text-gray-400">
      <div class="animate-pulse">{{ t('projectStructure.detecting') }}</div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="p-4 text-center text-red-400">
      {{ error }}
    </div>

    <!-- Content -->
    <div v-else-if="structure" class="flex-1 overflow-hidden flex flex-col">
      <!-- Tabs -->
      <div class="flex border-b border-gray-700/30 px-2">
        <button v-for="tab in ['architecture', 'frameworks', 'conventions', 'languages']" :key="tab"
          @click="activeTab = tab as any" class="tab-btn"
          :class="activeTab === tab ? 'tab-btn-active-indigo' : 'tab-btn-inactive'">
          {{ t(`projectStructure.${tab}`) }}
        </button>
      </div>

      <!-- Tab Content -->
      <div class="flex-1 overflow-y-auto p-3 space-y-3">
        <!-- Architecture Tab -->
        <template v-if="activeTab === 'architecture'">
          <div v-if="structure.architecture" class="space-y-3">
            <!-- Architecture Type -->
            <div class="card p-3">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-2xl">{{ getArchitectureIcon(structure.architecture.type) }}</span>
                <div>
                  <div class="font-medium text-white capitalize">{{ structure.architecture.type }}</div>
                  <div class="text-xs" :class="getConfidenceColor(structure.architecture.confidence)">
                    {{ t('projectStructure.confidence') }}: {{ Math.round(structure.architecture.confidence * 100) }}%
                  </div>
                </div>
              </div>
              <p class="text-sm text-gray-400">{{ structure.architecture.description }}</p>
            </div>

            <!-- Indicators -->
            <div v-if="structure.architecture.indicators?.length" class="card p-3">
              <div class="text-xs text-gray-400 mb-2">{{ t('projectStructure.indicators') }}</div>
              <ul class="space-y-1">
                <li v-for="ind in structure.architecture.indicators" :key="ind"
                  class="text-sm text-gray-300 flex items-start gap-2">
                  <span class="text-emerald-400 mt-0.5">✓</span>
                  <span>{{ ind }}</span>
                </li>
              </ul>
            </div>

            <!-- Layers -->
            <div v-if="structure.architecture.layers?.length" class="card p-3">
              <div class="text-xs text-gray-400 mb-2">{{ t('projectStructure.layers') }}</div>
              <div class="space-y-2">
                <div v-for="layer in structure.architecture.layers" :key="layer.name"
                  class="bg-gray-800/50 rounded p-2">
                  <div class="flex items-center justify-between">
                    <span class="font-medium text-indigo-400">{{ layer.name }}</span>
                    <span class="text-xs text-gray-400">{{ layer.path }}</span>
                  </div>
                  <p v-if="layer.description" class="text-xs text-gray-400 mt-1">{{ layer.description }}</p>
                  <div v-if="layer.dependencies?.length" class="text-xs text-gray-400 mt-1">
                    {{ t('projectStructure.dependencies') }}: {{ layer.dependencies.join(', ') }}
                  </div>
                  <div v-if="layer.patterns?.length" class="flex flex-wrap gap-1 mt-1">
                    <span v-for="p in layer.patterns" :key="p" class="badge badge-primary text-xs">{{ p }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Project Type -->
            <div class="card p-3">
              <div class="text-xs text-gray-400 mb-1">{{ t('projectStructure.projectType') }}</div>
              <div class="font-medium text-white capitalize">{{ structure.projectType }}</div>
            </div>
          </div>
          <div v-else class="text-center text-gray-400 py-8">
            {{ t('projectStructure.noArchitecture') }}
          </div>
        </template>

        <!-- Frameworks Tab -->
        <template v-else-if="activeTab === 'frameworks'">
          <div v-if="structure.frameworks?.length" class="space-y-2">
            <div v-for="fw in structure.frameworks" :key="fw.name" class="card p-3">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-lg">{{ getCategoryIcon(fw.category) }}</span>
                <div class="flex-1">
                  <div class="font-medium text-white">
                    {{ fw.name }}
                    <span v-if="fw.version" class="text-xs text-gray-400 ml-1">v{{ fw.version }}</span>
                  </div>
                  <div class="text-xs text-gray-400">{{ fw.language }} • {{ fw.category }}</div>
                </div>
              </div>

              <div v-if="fw.configFiles?.length" class="text-xs text-gray-400 mb-2">
                Config: {{ fw.configFiles.join(', ') }}
              </div>

              <div v-if="fw.bestPractices?.length" class="mt-2">
                <div class="text-xs text-gray-400 mb-1">{{ t('projectStructure.bestPractices') }}</div>
                <ul class="space-y-1">
                  <li v-for="bp in fw.bestPractices.slice(0, 3)" :key="bp"
                    class="text-xs text-gray-400 flex items-start gap-1">
                    <span class="text-yellow-400">•</span>
                    <span>{{ bp }}</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div v-else class="text-center text-gray-400 py-8">
            {{ t('projectStructure.noFrameworks') }}
          </div>
        </template>

        <!-- Conventions Tab -->
        <template v-else-if="activeTab === 'conventions'">
          <div v-if="structure.conventions" class="space-y-3">
            <div class="card p-3">
              <div class="text-xs text-gray-400 mb-1">{{ t('projectStructure.namingStyle') }}</div>
              <div class="font-medium text-white">{{ structure.conventions.namingStyle }}</div>
            </div>

            <div class="card p-3">
              <div class="text-xs text-gray-400 mb-1">{{ t('projectStructure.folderStructure') }}</div>
              <div class="font-medium text-white capitalize">{{ structure.conventions.folderStructure }}</div>
            </div>

            <div v-if="structure.conventions.testConventions" class="card p-3">
              <div class="text-xs text-gray-400 mb-1">{{ t('projectStructure.testConventions') }}</div>
              <div class="text-sm text-gray-300">
                <div v-if="structure.conventions.testConventions.framework">
                  Framework: <span class="text-white">{{ structure.conventions.testConventions.framework }}</span>
                </div>
                <div>Location: <span class="text-white">{{ structure.conventions.testConventions.location }}</span>
                </div>
                <div>Suffix: <span class="text-white">{{ structure.conventions.testConventions.fileSuffix }}</span>
                </div>
              </div>
            </div>

            <div v-if="structure.conventions.codeStyle" class="card p-3">
              <div class="text-xs text-gray-400 mb-1">{{ t('projectStructure.codeStyle') }}</div>
              <div class="text-sm text-gray-300">
                <div>Indent: <span class="text-white">{{ structure.conventions.codeStyle.indentStyle }} ({{
                  structure.conventions.codeStyle.indentSize }})</span></div>
                <div v-if="structure.conventions.codeStyle.configFile">
                  Config: <span class="text-white">{{ structure.conventions.codeStyle.configFile }}</span>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- Languages Tab -->
        <template v-else-if="activeTab === 'languages'">
          <div v-if="structure.languages?.length" class="space-y-2">
            <div v-for="lang in structure.languages" :key="lang.name" class="card p-3">
              <div class="flex items-center justify-between mb-1">
                <div class="font-medium text-white">
                  {{ lang.name }}
                  <span v-if="lang.primary" class="badge badge-success text-xs ml-1">{{ t('projectStructure.primary')
                    }}</span>
                </div>
                <span class="text-sm text-gray-400">{{ lang.fileCount }} files</span>
              </div>
              <div class="w-full bg-gray-700 rounded-full h-1.5">
                <div class="h-1.5 rounded-full transition-all" :class="lang.primary ? 'bg-indigo-500' : 'bg-gray-500'"
                  :style="{ width: `${lang.percentage}%` }"></div>
              </div>
              <div class="text-xs text-gray-400 mt-1">{{ lang.percentage.toFixed(1) }}%</div>
            </div>
          </div>

          <!-- Build Systems -->
          <div v-if="structure.buildSystems?.length" class="mt-4">
            <div class="text-xs text-gray-400 mb-2">{{ t('projectStructure.buildSystems') }}</div>
            <div class="space-y-2">
              <div v-for="bs in structure.buildSystems" :key="bs.name" class="card p-2">
                <div class="font-medium text-white text-sm">{{ bs.name }}</div>
                <div class="text-xs text-gray-400">{{ bs.configFile }}</div>
                <div v-if="bs.scripts?.length" class="flex flex-wrap gap-1 mt-1">
                  <span v-for="s in bs.scripts.slice(0, 5)" :key="s" class="badge badge-primary text-xs">{{ s }}</span>
                  <span v-if="bs.scripts.length > 5" class="text-xs text-gray-400">+{{ bs.scripts.length - 5 }}</span>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- No project -->
    <div v-else class="p-4 text-center text-gray-400">
      {{ t('files.noFiles') }}
    </div>
  </div>
</template>

<style scoped>
.project-structure-panel {
  @apply flex flex-col h-full bg-transparent;
}
</style>
