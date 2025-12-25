<template>
  <div class="flex-1 overflow-auto p-4 space-y-4">
    <!-- URL Input -->
    <div class="space-y-3">
      <div>
        <label class="block text-sm text-gray-300 mb-2">{{ t('git.repoUrl') }}</label>
        <div class="flex gap-2">
          <input 
            v-model="urlInput" 
            type="text" 
            :placeholder="t('git.urlPlaceholder')" 
            class="input flex-1"
            @keydown.enter="$emit('load-repo', urlInput)" 
          />
          <button @click="$emit('load-repo', urlInput)" :disabled="!urlInput || isLoading" class="btn btn-primary">
            <svg v-if="isLoading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            <span v-else>{{ t('git.load') }}</span>
          </button>
        </div>
      </div>
      <p class="text-xs text-gray-400">
        <span v-if="isGitHub" class="text-emerald-400">✓ GitHub API (fast, no clone)</span>
        <span v-else-if="isGitLab" class="text-orange-400">✓ GitLab API (fast, no clone)</span>
        <span v-else>{{ t('git.urlHint') }}</span>
      </p>
    </div>

    <!-- Remote Repo Loaded -->
    <div v-if="repoLoaded" class="space-y-4">
      <!-- Branch selector -->
      <div class="branch-info">
        <svg class="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <div class="flex-1">
          <select v-model="selectedBranchLocal" @change="$emit('change-branch', selectedBranchLocal)" class="input w-full text-sm">
            <option v-for="branch in branches" :key="branch.name" :value="branch.name">
              {{ branch.name }}
            </option>
          </select>
        </div>
      </div>

      <!-- Files list -->
      <GitFileList 
        v-if="files.length > 0"
        :title="`${t('git.filesAtRef')} ${selectedBranch}`"
        :files="files"
        :selected-files="selectedFiles"
        @toggle-select="$emit('toggle-file', $event)"
        @select-folder="$emit('select-folder', $event)"
        @select-all="$emit('select-all')"
        @clear-selection="$emit('clear-selection')"
        @preview-file="$emit('preview-file', $event)"
      />
    </div>

    <!-- Clone option for non-GitHub/GitLab -->
    <div v-if="!isGitHub && !isGitLab && urlInput && !repoLoaded"
      class="p-3 bg-amber-900/20 rounded-lg border border-amber-500/30">
      <p class="text-xs text-amber-300 mb-2">{{ t('git.cloneRequired') }}</p>
      <button @click="$emit('clone')" :disabled="isCloning" class="action-btn action-btn-accent w-full">
        {{ isCloning ? t('git.cloning') : t('git.clone') }}
      </button>
    </div>

    <!-- Cloned Repo Info -->
    <div v-if="clonedPath" class="p-3 bg-emerald-900/30 rounded-lg border border-emerald-500/30">
      <div class="flex items-center gap-2 mb-2">
        <svg class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span class="text-sm text-emerald-300">{{ t('git.cloned') }}</span>
      </div>
      <p class="text-xs text-gray-400 truncate">{{ clonedPath }}</p>
      <div class="flex gap-2 mt-3">
        <button @click="$emit('open-cloned')" class="btn btn-primary btn-sm flex-1">
          {{ t('git.openProject') }}
        </button>
        <button @click="$emit('cleanup-cloned')" class="btn btn-ghost btn-sm">
          {{ t('git.remove') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { ref, watch } from 'vue'
import GitFileList from './GitFileList.vue'

interface RemoteBranch {
  name: string
  commit: { sha: string }
}

const props = defineProps<{
  remoteUrl: string
  isGitHub: boolean
  isGitLab: boolean
  isLoading: boolean
  isCloning: boolean
  repoLoaded: boolean
  branches: RemoteBranch[]
  selectedBranch: string
  files: string[]
  selectedFiles: Set<string>
  clonedPath: string | null
}>()

defineEmits<{
  'load-repo': [url: string]
  'change-branch': [branch: string]
  'clone': []
  'open-cloned': []
  'cleanup-cloned': []
  'toggle-file': [path: string]
  'select-folder': [files: string[]]
  'select-all': []
  'clear-selection': []
  'preview-file': [path: string]
}>()

const { t } = useI18n()

const urlInput = ref(props.remoteUrl)
const selectedBranchLocal = ref(props.selectedBranch)

watch(() => props.remoteUrl, (val) => { urlInput.value = val })
watch(() => props.selectedBranch, (val) => { selectedBranchLocal.value = val })
</script>
