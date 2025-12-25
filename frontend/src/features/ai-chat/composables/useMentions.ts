/**
 * Mentions handler for chat input
 * Processes @files, @git, @problems mentions
 */

import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { gitApi } from '@/services/api/git.api'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { ref } from 'vue'

export interface MentionResult {
    type: 'files' | 'git' | 'problems'
    success: boolean
    message: string
    filesAdded?: number
}

type TranslateFunc = (key: string, params?: Record<string, string | number>) => string

export function useMentions() {
    const projectStore = useProjectStore()
    const fileStore = useFileStore()
    const contextStore = useContextStore()
    const uiStore = useUIStore()

    const isProcessing = ref(false)

    /**
     * Process @files mention - build context from selected files
     */
    async function processFilesMention(t: TranslateFunc): Promise<MentionResult> {
        const selectedFiles = fileStore.selectedFilesList

        if (selectedFiles.length === 0) {
            uiStore.addToast(t('chat.selectFilesHint'), 'warning')
            return {
                type: 'files',
                success: false,
                message: t('chat.selectFilesHint')
            }
        }

        try {
            isProcessing.value = true
            await contextStore.buildContext(selectedFiles)

            const message = t('chat.contextAttached')
            uiStore.addToast(message, 'success')

            return {
                type: 'files',
                success: true,
                message,
                filesAdded: selectedFiles.length
            }
        } catch (error) {
            const message = t('chat.contextBuildFailed')
            uiStore.addToast(message, 'error')
            return {
                type: 'files',
                success: false,
                message
            }
        } finally {
            isProcessing.value = false
        }
    }

    /**
     * Process @git mention - get uncommitted files and build context
     */
    async function processGitMention(t: TranslateFunc): Promise<MentionResult> {
        const projectPath = projectStore.currentPath

        if (!projectPath) {
            uiStore.addToast(t('chat.noApiKey'), 'warning')
            return {
                type: 'git',
                success: false,
                message: 'No project selected'
            }
        }

        try {
            isProcessing.value = true

            const uncommittedFiles = await gitApi.getUncommittedFiles(projectPath)

            if (!uncommittedFiles || uncommittedFiles.length === 0) {
                const message = t('gitContext.noChanges')
                uiStore.addToast(message, 'info')
                return {
                    type: 'git',
                    success: true,
                    message,
                    filesAdded: 0
                }
            }

            // Extract file paths from uncommitted files
            const filePaths = uncommittedFiles.map(f => f.path)

            // Build context from these files
            await contextStore.buildContext(filePaths)

            const message = `${t('chat.contextAttached')} (${filePaths.length} git files)`
            uiStore.addToast(message, 'success')

            return {
                type: 'git',
                success: true,
                message,
                filesAdded: filePaths.length
            }
        } catch (error) {
            const message = t('chat.contextBuildFailed')
            uiStore.addToast(message, 'error')
            return {
                type: 'git',
                success: false,
                message
            }
        } finally {
            isProcessing.value = false
        }
    }

    /**
     * Process @problems mention - placeholder for future implementation
     */
    async function processProblemsMention(t: TranslateFunc): Promise<MentionResult> {
        uiStore.addToast(t('chat.comingSoon'), 'info')
        return {
            type: 'problems',
            success: false,
            message: t('chat.comingSoon')
        }
    }

    /**
     * Process mention by type
     */
    async function processMention(
        mentionType: 'files' | 'git' | 'problems',
        t: TranslateFunc
    ): Promise<MentionResult> {
        switch (mentionType) {
            case 'files':
                return processFilesMention(t)
            case 'git':
                return processGitMention(t)
            case 'problems':
                return processProblemsMention(t)
            default:
                return {
                    type: mentionType,
                    success: false,
                    message: 'Unknown mention type'
                }
        }
    }

    /**
     * Parse message for mentions and process them
     * Returns cleaned message without mention tags
     */
    async function processMessageMentions(
        message: string,
        t: TranslateFunc
    ): Promise<{ cleanedMessage: string; results: MentionResult[] }> {
        const results: MentionResult[] = []
        let cleanedMessage = message

        // Process @files
        if (message.includes('@files')) {
            const result = await processMention('files', t)
            results.push(result)
            cleanedMessage = cleanedMessage.replace(/@files\s*/g, '')
        }

        // Process @git
        if (message.includes('@git')) {
            const result = await processMention('git', t)
            results.push(result)
            cleanedMessage = cleanedMessage.replace(/@git\s*/g, '')
        }

        // Process @problems
        if (message.includes('@problems')) {
            const result = await processMention('problems', t)
            results.push(result)
            cleanedMessage = cleanedMessage.replace(/@problems\s*/g, '')
        }

        return {
            cleanedMessage: cleanedMessage.trim(),
            results
        }
    }

    return {
        isProcessing,
        processMention,
        processMessageMentions,
        processFilesMention,
        processGitMention,
        processProblemsMention
    }
}
