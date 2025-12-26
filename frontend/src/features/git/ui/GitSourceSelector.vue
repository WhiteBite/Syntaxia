<template>
  <div class="h-full flex flex-col bg-transparent">
    <!-- Header with tabs -->
    <GitSourceHeader
      :source-type="sourceType"
      :recent-repos="recentRepos"
      @change-source="sourceType = $event"
      @select-recent="handleSelectRecentRepo"
      @clear-recent="clearRecentRepos"
    />

    <!-- Local Git Panel -->
    <GitLocalPanel 
      v-if="sourceType === 'local'"
      :is-git-repo="isGitRepo"
      :current-branch="currentBranch"
      :branches="branches"
      :commits="commits"
      :commits-loaded="commitsLoaded"
      :selected-ref="selectedRef"
      :files-at-ref="filesAtRef"
      :selected-files="selectedFiles"
      @select-ref="selectRef"
      @clear-ref="clearSelectedRef"
      @load-commits="loadCommits"
      @open-diff="diffModalOpen = true"
      @toggle-file="toggleFileSelection"
      @select-folder="handleSelectFolder"
      @select-all="selectAllFiles"
      @clear-selection="clearFileSelection"
      @preview-file="handlePreviewFile"
    />

    <!-- Remote URL Panel -->
    <GitRemotePanel 
      v-if="sourceType === 'remote'"
      :remote-url="remoteUrl"
      :is-git-hub="isGitHubRepo"
      :is-git-lab="isGitLabRepo"
      :is-loading="isLoadingRemote"
      :is-cloning="isCloning"
      :repo-loaded="remoteRepoLoaded"
      :branches="remoteBranches"
      :selected-branch="remoteSelectedBranch"
      :files="remoteFiles"
      :selected-files="remoteSelectedFiles"
      :cloned-path="clonedPath"
      @load-repo="loadRemoteRepo"
      @change-branch="handleChangeBranch"
      @clone="cloneRemote"
      @open-cloned="openClonedRepo"
      @cleanup-cloned="cleanupClonedRepo"
      @toggle-file="toggleRemoteFileSelection"
      @select-folder="handleSelectRemoteFolder"
      @select-all="selectAllRemoteFiles"
      @clear-selection="clearRemoteFileSelection"
      @preview-file="handlePreviewFile"
    />

    <!-- Bottom Panel for Remote Files -->
    <div v-if="sourceType === 'remote' && remoteSelectedFiles.size > 0" class="border-t border-gray-700 p-4">
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm text-gray-300">{{ remoteSelectedFiles.size }} {{ t('git.filesSelected') }}</span>
      </div>
      <button @click="buildContextFromRemote" :disabled="isBuilding" class="btn btn-primary w-full">
        <svg v-if="isBuilding" class="animate-spin w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        {{ t('git.buildContext') }} {{ remoteSelectedBranch }}
      </button>
    </div>

    <!-- Loading Overlay -->
    <div v-if="isLoading" class="loading-overlay absolute">
      <div class="text-center">
        <svg class="loading-spinner mx-auto mb-2" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        <span class="text-sm text-gray-400">{{ loadingMessage }}</span>
      </div>
    </div>

    <!-- Bottom Panel - Build Context (Local) -->
    <div v-if="sourceType === 'local' && selectedFiles.size > 0" class="border-t border-gray-700 p-4">
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm text-gray-300">{{ selectedFiles.size }} {{ t('git.filesSelected') }}</span>
      </div>
      <button @click="buildContextFromRef" :disabled="isBuilding" class="btn btn-primary w-full">
        <svg v-if="isBuilding" class="animate-spin w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        {{ t('git.buildContext') }} {{ selectedRef?.slice(0, 7) || 'ref' }}
      </button>
    </div>

    <!-- File Preview Modal -->
    <FilePreviewModal 
      :is-open="previewOpen" 
      :file-path="previewPath" 
      :content="previewContent"
      :is-loading="previewLoading" 
      :error="previewError" 
      @close="previewOpen = false" 
    />

    <!-- Branch Diff Modal -->
    <BranchDiffModal 
      :is-open="diffModalOpen" 
      :branches="branches" 
      :project-path="projectPath"
      :current-branch="currentBranch" 
      @close="diffModalOpen = false" 
    />
  </div>
</template>

<script setup lang="ts">
import BranchDiffModal from '@/components/BranchDiffModal.vue'
import FilePreviewModal from '@/components/FilePreviewModal.vue'
import { useGitSource, type RecentRepo } from '@/composables/useGitSource'
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useContextStore } from '@/features/context'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { onMounted, ref, watch } from 'vue'
import GitLocalPanel from './GitLocalPanel.vue'
import GitRemotePanel from './GitRemotePanel.vue'
import GitSourceHeader from './GitSourceHeader.vue'
import type { SourceType } from './GitSourceTabs.vue'

const logger = useLogger('GitSourceSelector')
const { t } = useI18n()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const contextStore = useContextStore()

const sourceType = ref<SourceType>('local')

// Use composable for git logic
const git = useGitSource()

// Destructure for template access
const {
  isGitRepo, currentBranch, branches, commits, commitsLoaded,
  selectedRef, filesAtRef, selectedFiles,
  remoteUrl, isCloning, clonedPath, isGitHubRepo, isGitLabRepo,
  isLoadingRemote, remoteRepoLoaded, remoteBranches, remoteSelectedBranch,
  remoteFiles, remoteSelectedFiles,
  isLoading, loadingMessage, isBuilding, projectPath, recentRepos,
  loadCommits, selectRef, clearSelectedRef,
  toggleFileSelection, selectAllFiles, clearFileSelection,
  loadRemoteRepo: loadRemoteRepoBase, loadRemoteFiles,
  toggleRemoteFileSelection, selectAllRemoteFiles, clearRemoteFileSelection,
  loadRecentReposFromStorage, clearRecentRepos,
} = git

// Additional local state
const diffModalOpen = ref(false)
const previewOpen = ref(false)
const previewPath = ref('')
const previewContent = ref('')
const previewLoading = ref(false)
const previewError = ref('')

// Handle folder selection
function handleSelectFolder(files: string[]) {
  files.forEach(f => { if (!selectedFiles.value.has(f)) selectedFiles.value.add(f) })
  selectedFiles.value = new Set(selectedFiles.value)
}

function handleSelectRemoteFolder(files: string[]) {
  files.forEach(f => { if (!remoteSelectedFiles.value.has(f)) remoteSelectedFiles.value.add(f) })
  remoteSelectedFiles.value = new Set(remoteSelectedFiles.value)
}

// Check git repo on project change
async function checkGitRepo() {
  if (!projectPath.value) return
  isLoading.value = true
  loadingMessage.value = 'Checking repository...'
  commitsLoaded.value = false
  branches.value = []
  commits.value = []

  try {
    isGitRepo.value = await apiService.isGitRepository(projectPath.value)
    if (isGitRepo.value) {
      currentBranch.value = await apiService.getCurrentBranch(projectPath.value)
      const result = await apiService.getBranches(projectPath.value)
      branches.value = JSON.parse(result)
    }
  } catch (err) {
    logger.error('Failed to check git repo', err)
    isGitRepo.value = false
  } finally {
    isLoading.value = false
  }
}

async function loadRemoteRepo(url: string) {
  remoteUrl.value = url
  await loadRemoteRepoBase()
}

async function handleChangeBranch(branch: string) {
  remoteSelectedBranch.value = branch
  await loadRemoteFiles()
}

function handleSelectRecentRepo(repo: RecentRepo) {
  remoteUrl.value = repo.url
  sourceType.value = 'remote'
  loadRemoteRepoBase()
}

async function buildContextFromRef() {
  if (!selectedRef.value || selectedFiles.value.size === 0) return
  isBuilding.value = true
  try {
    const files = Array.from(selectedFiles.value)
    const content = await apiService.buildContextAtRef(projectPath.value, files, selectedRef.value)
    contextStore.setRawContext(content, files.length)
    uiStore.addToast(`Context built from ${selectedRef.value.slice(0, 7)}: ${files.length} files`, 'success')
  } catch (err) {
    logger.error('Failed to build context from ref', err)
    uiStore.addToast('Failed to build context from ref', 'error')
  } finally {
    isBuilding.value = false
  }
}

async function buildContextFromRemote() {
  if (!remoteUrl.value || remoteSelectedFiles.value.size === 0) return
  isBuilding.value = true
  try {
    const files = Array.from(remoteSelectedFiles.value)
    let content = '', source = ''
    if (isGitHubRepo.value) {
      content = await apiService.gitHubBuildContext(remoteUrl.value, files, remoteSelectedBranch.value)
      source = 'GitHub'
    } else if (isGitLabRepo.value) {
      content = await apiService.gitLabBuildContext(remoteUrl.value, files, remoteSelectedBranch.value)
      source = 'GitLab'
    }
    contextStore.setRawContext(content, files.length)
    uiStore.addToast(`Context built from ${source}: ${files.length} files`, 'success')
  } catch (err) {
    logger.error('Failed to build context from remote', err)
    uiStore.addToast('Failed to build context', 'error')
  } finally {
    isBuilding.value = false
  }
}

async function cloneRemote() {
  if (!remoteUrl.value) return
  isCloning.value = true
  isLoading.value = true
  loadingMessage.value = 'Cloning repository...'
  try {
    clonedPath.value = await apiService.cloneRepository(remoteUrl.value)
    uiStore.addToast('Repository cloned successfully', 'success')
  } catch (err) {
    logger.error('Failed to clone repository', err)
    uiStore.addToast('Failed to clone repository', 'error')
  } finally {
    isCloning.value = false
    isLoading.value = false
  }
}

async function openClonedRepo() {
  if (!clonedPath.value) return
  await projectStore.openProjectByPath(clonedPath.value)
  sourceType.value = 'local'
}

async function cleanupClonedRepo() {
  if (!clonedPath.value) return
  try {
    await apiService.cleanupTempRepository(clonedPath.value)
    clonedPath.value = null
    uiStore.addToast('Temporary repository removed', 'success')
  } catch (err) {
    logger.warn('Failed to cleanup cloned repo', err)
  }
}

async function handlePreviewFile(filePath: string) {
  previewPath.value = filePath
  previewContent.value = ''
  previewError.value = ''
  previewLoading.value = true
  previewOpen.value = true
  try {
    let content = ''
    if (sourceType.value === 'local' && selectedRef.value) {
      content = await apiService.getFileAtRef(projectPath.value, filePath, selectedRef.value)
    } else if (sourceType.value === 'remote') {
      if (isGitHubRepo.value) {
        content = await apiService.gitHubGetFileContent(remoteUrl.value, filePath, remoteSelectedBranch.value)
      } else if (isGitLabRepo.value) {
        content = await apiService.gitLabGetFileContent(remoteUrl.value, filePath, remoteSelectedBranch.value)
      }
    }
    previewContent.value = content
  } catch (err) {
    logger.error('Failed to preview file', err)
    previewError.value = t('error.loadFailed')
  } finally {
    previewLoading.value = false
  }
}

watch(() => projectStore.currentPath, checkGitRepo, { immediate: true })
onMounted(() => { checkGitRepo(); loadRecentReposFromStorage() })
</script>
