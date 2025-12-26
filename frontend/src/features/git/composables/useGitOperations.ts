import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useContextStore } from '@/features/context'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { ref, type Ref } from 'vue'

export interface UseGitOperationsOptions {
    projectPath: Ref<string>
    selectedRef: Ref<string | null>
    selectedFiles: Ref<Set<string>>
    remoteUrl: Ref<string>
    remoteSelectedBranch: Ref<string>
    remoteSelectedFiles: Ref<Set<string>>
    isGitHubRepo: Ref<boolean>
    isGitLabRepo: Ref<boolean>
    isLoading: Ref<boolean>
    loadingMessage: Ref<string>
    isBuilding: Ref<boolean>
    isCloning: Ref<boolean>
    clonedPath: Ref<string | null>
}

export function useGitOperations(options: UseGitOperationsOptions) {
    const logger = useLogger('useGitOperations')
    const { t } = useI18n()
    const uiStore = useUIStore()
    const contextStore = useContextStore()
    const projectStore = useProjectStore()

    const {
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
    } = options

    // Preview state
    const previewOpen = ref(false)
    const previewPath = ref('')
    const previewContent = ref('')
    const previewLoading = ref(false)
    const previewError = ref('')

    async function buildContextFromRef(): Promise<void> {
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

    async function buildContextFromRemote(): Promise<void> {
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
        } catch (err) {
            logger.error('Failed to build context from remote', err)
            uiStore.addToast('Failed to build context', 'error')
        } finally {
            isBuilding.value = false
        }
    }

    async function cloneRemote(): Promise<void> {
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

    async function openClonedRepo(): Promise<void> {
        if (!clonedPath.value) return
        await projectStore.openProjectByPath(clonedPath.value)
    }

    async function cleanupClonedRepo(): Promise<void> {
        if (!clonedPath.value) return

        try {
            await apiService.cleanupTempRepository(clonedPath.value)
            clonedPath.value = null
            uiStore.addToast('Temporary repository removed', 'success')
        } catch (err) {
            logger.warn('Failed to cleanup cloned repo', err)
        }
    }

    async function previewFile(filePath: string, sourceType: 'local' | 'remote'): Promise<void> {
        previewPath.value = filePath
        previewContent.value = ''
        previewError.value = ''
        previewLoading.value = true
        previewOpen.value = true

        try {
            let content = ''

            if (sourceType === 'local' && selectedRef.value) {
                content = await apiService.getFileAtRef(projectPath.value, filePath, selectedRef.value)
            } else if (sourceType === 'remote') {
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

    function closePreview(): void {
        previewOpen.value = false
    }

    return {
        // Preview state
        previewOpen,
        previewPath,
        previewContent,
        previewLoading,
        previewError,

        // Methods
        buildContextFromRef,
        buildContextFromRemote,
        cloneRemote,
        openClonedRepo,
        cleanupClonedRepo,
        previewFile,
        closePreview,
    }
}
