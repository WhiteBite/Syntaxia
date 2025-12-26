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
      @preview-file="handlePreviewLocalFile"
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
      @open-cloned="handleOpenClonedRepo"
      @cleanup-cloned="cleanupClonedRepo"
      @toggle-file="toggleRemoteFileSelection"
      @select-folder="handleSelectRemoteFolder"
      @select-all="selectAllRemoteFiles"
      @clear-selection="clearRemoteFileSelection"
      @preview-file="handlePreviewRemoteFile"
    />

    <!-- Loading Overlay -->
    <GitLoadingOverlay :is-loading="isLoading" :message="loadingMessage" />

    <!-- Build Context Panel -->
    <GitBuildPanel
      :source-type="sourceType"
      :local-selected-count="selectedFiles.size"
      :remote-selected-count="remoteSelectedFiles.size"
      :local-ref="selectedRef"
      :remote-branch="remoteSelectedBranch"
      :is-building="isBuilding"
      @build-local="buildContextFromRef"
      @build-remote="buildContextFromRemote"
    />

    <!-- File Preview Modal -->
    <FilePreviewModal 
      :is-open="previewOpen" 
      :file-path="previewPath" 
      :content="previewContent"
      :is-loading="previewLoading" 
      :error="previewError" 
      @close="closePreview" 
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
import { useLogger } from '@/composables/useLogger'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { onMounted, ref, watch } from 'vue'
import { useGitOperations } from '../composables/useGitOperations'
import GitBuildPanel from './GitBuildPanel.vue'
import GitLoadingOverlay from './GitLoadingOverlay.vue'
import GitLocalPanel from './GitLocalPanel.vue'
import GitRemotePanel from './GitRemotePanel.vue'
import GitSourceHeader from './GitSourceHeader.vue'
import type { SourceType } from './GitSourceTabs.vue'

const logger = useLogger('GitSourceSelector')
const projectStore = useProjectStore()

const sourceType = ref<SourceType>('local')
const diffModalOpen = ref(false)

// Use composable for git state
const git = useGitSource()

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

// Use composable for git operations
const operations = useGitOperations({
  projectPath,
  selectedRef,
  selectedFiles,
  remoteUrl,
  remoteSelectedBranch,
  remoteSelectedFiles,
  isGitHubRepo,
  isGitLabRepo,
  isLoading,
  loadingMessage,
  isBuilding,
  isCloning,
  clonedPath,
})

const {
  previewOpen, previewPath, previewContent, previewLoading, previewError,
  buildContextFromRef, buildContextFromRemote,
  cloneRemote, openClonedRepo, cleanupClonedRepo,
  previewFile, closePreview,
} = operations

// Folder selection handlers
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

async function handleOpenClonedRepo() {
  await openClonedRepo()
  sourceType.value = 'local'
}

function handlePreviewLocalFile(filePath: string) {
  previewFile(filePath, 'local')
}

function handlePreviewRemoteFile(filePath: string) {
  previewFile(filePath, 'remote')
}

watch(() => projectStore.currentPath, checkGitRepo, { immediate: true })
onMounted(() => { checkGitRepo(); loadRecentReposFromStorage() })
</script>
