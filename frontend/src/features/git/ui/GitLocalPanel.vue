<template>
  <div class="flex-1 overflow-auto">
    <!-- Not a Git Repo -->
    <div v-if="!isGitRepo" class="p-4 text-center text-gray-400">
      <svg class="w-12 h-12 mx-auto mb-3 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <p>{{ t('git.notGitRepo') }}</p>
      <p class="text-sm mt-2">{{ t('git.openGitProject') }}</p>
    </div>

    <div v-else class="p-4 space-y-4">
      <!-- Current Branch Info -->
      <div class="branch-info">
        <svg class="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <div class="flex-1">
          <span class="text-xs text-gray-400">{{ t('git.currentBranch') }}:</span>
          <span class="ml-2 text-sm text-white font-medium">{{ currentBranch }}</span>
        </div>
        <button @click="$emit('select-ref', currentBranch)" class="btn btn-primary btn-sm" :disabled="selectedRef === currentBranch">
          {{ t('git.load') }}
        </button>
      </div>

      <!-- Branch/Commit Selection -->
      <div class="space-y-3">
        <div class="flex border-b border-gray-700">
          <button @click="refType = 'branches'" :class="['tab-btn text-sm', refType === 'branches' ? 'tab-btn-active' : 'tab-btn-inactive']">
            {{ t('git.branches') }} ({{ branches.length }})
          </button>
          <button @click="switchToCommits" :class="['tab-btn text-sm', refType === 'commits' ? 'tab-btn-active' : 'tab-btn-inactive']">
            {{ t('git.commits') }} {{ commitsLoaded ? `(${commits.length})` : '' }}
          </button>
          <button v-if="branches.length >= 2" @click="$emit('open-diff')"
            class="ml-auto px-2 py-1 text-xs text-purple-400 hover:text-purple-300 hover:bg-purple-500/10 rounded transition-colors flex items-center gap-1">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
            {{ t('git.diff') }}
          </button>
        </div>

        <!-- Selected Ref Info -->
        <div v-if="selectedRef" class="p-3 bg-indigo-900/30 rounded-lg border border-indigo-500/30">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-xs text-indigo-300">{{ t('git.buildingFrom') }}:</span>
              <span class="ml-2 text-sm text-white font-medium">{{ selectedRef }}</span>
              <span v-if="selectedRef === currentBranch" class="text-xs text-emerald-400 ml-2">({{ t('git.current') }})</span>
            </div>
            <button @click="$emit('clear-ref')" class="text-xs text-gray-400 hover:text-white">
              {{ t('git.clear') }}
            </button>
          </div>
          <p class="text-xs text-gray-400 mt-1">{{ t('git.noCheckout') }}</p>
        </div>

        <!-- Branches List -->
        <div v-if="refType === 'branches'" class="branches-list space-y-1">
          <button v-for="branch in branches" :key="branch" @click="$emit('select-ref', branch)" :class="[
            'w-full px-3 py-2 text-left text-sm rounded transition-colors flex items-center gap-2',
            selectedRef === branch ? 'bg-indigo-600/20 text-indigo-300 border border-indigo-500/30'
              : branch === currentBranch ? 'bg-emerald-900/20 text-emerald-300 border border-emerald-500/20'
              : 'text-gray-300 hover:bg-gray-700'
          ]">
            <svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            <span class="truncate">{{ branch }}</span>
            <span v-if="branch === currentBranch" class="text-xs text-emerald-400 ml-auto">({{ t('git.current') }})</span>
            <span v-else-if="selectedRef === branch" class="text-xs text-indigo-400 ml-auto">({{ t('git.selected') }})</span>
          </button>
        </div>

        <!-- Commits List -->
        <div v-if="refType === 'commits'" class="branches-list space-y-1">
          <button v-for="commit in commits" :key="commit.hash" @click="$emit('select-ref', commit.hash)" :class="[
            'w-full px-3 py-2 text-left text-sm rounded transition-colors',
            selectedRef === commit.hash ? 'bg-indigo-600/20 text-indigo-300 border border-indigo-500/30' : 'text-gray-300 hover:bg-gray-700'
          ]">
            <div class="flex items-start gap-2">
              <code class="text-xs text-amber-400 font-mono flex-shrink-0">{{ commit.hash.slice(0, 7) }}</code>
              <div class="flex-1 min-w-0">
                <p class="truncate text-white">{{ commit.subject }}</p>
                <p class="text-xs text-gray-400">{{ commit.author }} • {{ formatDate(commit.date) }}</p>
              </div>
            </div>
          </button>
        </div>
      </div>

      <!-- File Tree at Selected Ref -->
      <GitFileList 
        v-if="selectedRef && filesAtRef.length > 0"
        :title="`${t('git.filesAtRef')} ${selectedRef?.slice(0, 7)}`"
        :files="filesAtRef"
        :selected-files="selectedFiles"
        @toggle-select="$emit('toggle-file', $event)"
        @select-folder="$emit('select-folder', $event)"
        @select-all="$emit('select-all')"
        @clear-selection="$emit('clear-selection')"
        @preview-file="$emit('preview-file', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import type { CommitInfo } from '@/services/api.service';
import { ref } from 'vue';
import GitFileList from './GitFileList.vue';

defineProps<{
  isGitRepo: boolean
  currentBranch: string
  branches: string[]
  commits: CommitInfo[]
  commitsLoaded: boolean
  selectedRef: string | null
  filesAtRef: string[]
  selectedFiles: Set<string>
}>()

const emit = defineEmits<{
  'select-ref': [ref: string]
  'clear-ref': []
  'load-commits': []
  'open-diff': []
  'toggle-file': [path: string]
  'select-folder': [files: string[]]
  'select-all': []
  'clear-selection': []
  'preview-file': [path: string]
}>()

const { t } = useI18n()
const refType = ref<'branches' | 'commits'>('branches')

function switchToCommits() {
  refType.value = 'commits'
  emit('load-commits')
}

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleDateString()
  } catch {
    return dateStr
  }
}
</script>

<style scoped>
.branches-list {
  max-height: 40vh;
  min-height: 0;
  overflow-y: auto;
}
</style>
