<template>
  <div class="h-full flex flex-col bg-transparent">
    <!-- Header with tabs -->
    <div class="border-b border-gray-700/30">
      <div class="flex items-center justify-between p-3">
        <div class="section-title">
          <div class="section-icon section-icon-orange">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
            </svg>
          </div>
          <h2 class="section-title-text">{{ t('git.title') }}</h2>
        </div>
        <!-- Recent repos dropdown -->
        <RecentReposDropdown 
          v-if="recentRepos.length > 0"
          :repos="recentRepos"
          @select="handleSelectRecentRepo"
          @clear="clearRecentRepos"
        />
      </div>

      <!-- Source Type Tabs -->
      <div class="flex gap-1 px-2 pb-2">
        <button @click="sourceType = 'local'" :class="['tab-btn', sourceType === 'local' ? 'tab-btn-active tab-btn-active-indigo' : 'tab-btn-inactive']">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          {{ t('git.localGit') }}
        </button>
        <button @click="sourceType = 'remote'" :class="['tab-btn', sourceType === 'remote' ? 'tab-btn-active tab-btn-active-purple' : 'tab-btn-inactive']">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
          </svg>
          {{ t('git.remoteUrl') }}
        </button>
      </div>
    </div>

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
import { useContextStore } from '@/features/context'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { onMounted, ref, watch } from 'vue'
import GitLocalPanel from './GitLocalPanel.vue'
import GitRemotePanel from './GitRemotePanel.vue'
import RecentReposDropdown from './RecentReposDropdown.vue'

const { t } = useI18n()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const contextStore = useContextStore()

const sourceType = ref<'local' | 'remote'>('local')

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

// Handle folder selection - receives array of files from SimpleFileTree
function handleSelectFolder(files: string[]) {
  files.forEach(f => {
    if (!selectedFiles.value.has(f)) {
      selectedFiles.value.add(f)
    }
  })
  selectedFiles.value = new Set(selectedFiles.value)
}

function handleSelectRemoteFolder(files: string[]) {
  files.forEach(f => {
    if (!remoteSelectedFiles.value.has(f)) {
      remoteSelectedFiles.value.add(f)
    }
  })
  remoteSelectedFiles.value = new Set(remoteSelectedFiles.value)
}

// Additional local state
const diffModalOpen = ref(false)
const previewOpen = ref(false)
const previewPath = ref('')
const previewContent = ref('')
const previewLoading = ref(false)
const previewError = ref('')

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
  } catch {
    isGitRepo.value = false
  } finally {
    isLoading.value = false
  }
}

// Load remote repo with URL
async function loadRemoteRepo(url: string) {
  remoteUrl.value = url
  await loadRemoteRepoBase()
}

// Handle branch change
async function handleChangeBranch(branch: string) {
  remoteSelectedBranch.value = branch
  await loadRemoteFiles()
}

// Handle recent repo selection
function handleSelectRecentRepo(repo: RecentRepo) {
  remoteUrl.value = repo.url
  sourceType.value = 'remote'
  loadRemoteRepoBase()
}

// Build context from local ref
async function buildContextFromRef() {
  if (!selectedRef.value || selectedFiles.value.size === 0) return
  
  isBuilding.value = true
  try {
    const files = Array.from(selectedFiles.value)
    const content = await apiService.buildContextAtRef(projectPath.value, files, selectedRef.value)
    contextStore.setRawContext(content, files.length)
    uiStore.addToast(`Context built from ${selectedRef.value.slice(0, 7)}: ${files.length} files`, 'success')
  } catch {
    uiStore.addToast('Failed to build context from ref', 'error')
  } finally {
    isBuilding.value = false
  }
}

// Build context from remote
async function buildContextFromRemote() {
  if (!remoteUrl.value || remoteSelectedFiles.value.size === 0) return
  
  isBuilding.value = true
  try {
    const files = Array.from(remoteSelectedFiles.value)
    let content = ''
    let source = ''

    if (isGitHubRepo.value) {
      content = await apiService.gitHubBuildContext(remoteUrl.value, files, remoteSelectedBranch.value)
      source = 'GitHub'
    } else if (isGitLabRepo.value) {
      content = await apiService.gitLabBuildContext(remoteUrl.value, files, remoteSelectedBranch.value)
      source = 'GitLab'
    }

    contextStore.setRawContext(content, files.length)
    uiStore.addToast(`Context built from ${source}: ${files.length} files`, 'success')
  } catch {
    uiStore.addToast('Failed to build context', 'error')
  } finally {
    isBuilding.value = false
  }
}

// Clone remote repo
async function cloneRemote() {
  if (!remoteUrl.value) return
  
  isCloning.value = true
  isLoading.value = true
  loadingMessage.value = 'Cloning repository...'

  try {
    clonedPath.value = await apiService.cloneRepository(remoteUrl.value)
    uiStore.addToast('Repository cloned successfully', 'success')
  } catch {
    uiStore.addToast('Failed to clone repository', 'error')
  } finally {
    isCloning.value = false
    isLoading.value = false
  }
}

// Open cloned repo
async function openClonedRepo() {
  if (!clonedPath.value) return
  await projectStore.openProjectByPath(clonedPath.value)
  sourceType.value = 'local'
}

// Cleanup cloned repo
async function cleanupClonedRepo() {
  if (!clonedPath.value) return
  try {
    await apiService.cleanupTempRepository(clonedPath.value)
    clonedPath.value = null
    uiStore.addToast('Temporary repository removed', 'success')
  } catch {
    // Ignore
  }
}

// File preview
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
  } catch {
    previewError.value = t('error.loadFailed')
  } finally {
    previewLoading.value = false
  }
}

watch(() => projectStore.currentPath, checkGitRepo, { immediate: true })

onMounted(() => {
  checkGitRepo()
  loadRecentReposFromStorage()
})
</script>
