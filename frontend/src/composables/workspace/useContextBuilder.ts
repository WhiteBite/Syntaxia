/**
 * Context Builder Logic
 * 
 * Handles context building, copying, and automatic rebuilding on settings changes.
 */

import { useI18n } from '@/composables/useI18n';
import { useLogger } from '@/composables/useLogger';
import { useContextStore } from '@/features/context';
import { useFileStore } from '@/features/files';
import { detectLanguages, generateFileTree, useTemplateStore } from '@/features/templates';
import { useProjectStore } from '@/stores/project.store';
import { useSettingsStore } from '@/stores/settings.store';
import { useUIStore } from '@/stores/ui.store';
import { watch } from 'vue';

const logger = useLogger('ContextBuilder');

export function useContextBuilder() {
    const { t } = useI18n();
    const contextStore = useContextStore();
    const fileStore = useFileStore();
    const uiStore = useUIStore();
    const settingsStore = useSettingsStore();
    const templateStore = useTemplateStore();
    const projectStore = useProjectStore();

    /**
     * Build context from selected files
     */
    async function buildContext() {
        if (fileStore.selectedPaths.size === 0) {
            uiStore.addToast('Выберите файлы для построения контекста', 'warning');
            return;
        }

        if (contextStore.isBuilding) return;

        try {
            const filePaths = Array.from(fileStore.selectedPaths);
            const options = {
                maxTokens: settingsStore.settings.context.maxTokens,
                stripComments: settingsStore.settings.context.stripComments,
                includeTests: settingsStore.settings.context.includeTests,
                splitStrategy: settingsStore.settings.context.splitStrategy,
                outputFormat: settingsStore.settings.context.outputFormat,
                includeLineNumbers: settingsStore.settings.context.includeLineNumbers,
                excludeTests: settingsStore.settings.context.excludeTests,
                collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
                stripLicense: settingsStore.settings.context.stripLicense,
                compactDataFiles: settingsStore.settings.context.compactDataFiles,
                trimWhitespace: settingsStore.settings.context.trimWhitespace
            };

            await contextStore.buildContext(filePaths, options);

            // Show success toast with copy action
            uiStore.addToast(t('toast.contextBuilt'), 'success', 5000, {
                label: t('context.copy'),
                icon: '📋',
                onClick: async () => {
                    try {
                        const content = await contextStore.getFullContextContent();
                        await navigator.clipboard.writeText(content);
                        uiStore.addToast(t('toast.contextCopied'), 'success');
                    } catch {
                        uiStore.addToast(t('toast.copyError'), 'error');
                    }
                }
            });
        } catch (error) {
            logger.error('Failed to build context:', error);

            // Handle token limit exceeded error with detailed message
            if (error instanceof Error && error.message === 'TOKEN_LIMIT_EXCEEDED') {
                const storeError = contextStore.error;
                if (storeError?.startsWith('TOKEN_LIMIT_EXCEEDED:')) {
                    const parts = storeError.split(':');
                    const actual = Number(parts[1]);
                    const limit = Number(parts[2]);
                    const actualK = Math.round(actual / 1000);
                    const limitK = Math.round(limit / 1000);
                    uiStore.addToast(
                        t('error.tokenLimitExceeded', { actual: actualK, limit: limitK }),
                        'error'
                    );
                } else {
                    uiStore.addToast(t('error.tokenLimitGeneric'), 'error');
                }
                return;
            }

            const errorMsg = error instanceof Error ? error.message : 'Unknown error';
            uiStore.addToast(`${t('toast.contextError')}: ${errorMsg}`, 'error');
        }
    }

    /**
     * Copy context to clipboard with optional template
     */
    async function copyContext() {
        if (contextStore.hasContext && contextStore.contextId) {
            try {
                const filesContent = await contextStore.getFullContextContent();

                let content: string;
                if (settingsStore.settings.context.applyTemplateOnCopy && templateStore.activeTemplate) {
                    const files = contextStore.summary?.files || [];
                    const templateContext = {
                        fileTree: generateFileTree(files, projectStore.projectName),
                        files: filesContent,
                        task: templateStore.currentTask,
                        userRules: templateStore.userRules,
                        fileCount: contextStore.fileCount,
                        tokenCount: contextStore.tokenCount,
                        languages: detectLanguages(files),
                        projectName: projectStore.projectName
                    };
                    content = templateStore.generatePrompt(templateContext);
                } else {
                    content = filesContent;
                }

                await navigator.clipboard.writeText(content);
                uiStore.addToast(t('toast.contextCopied'), 'success');
            } catch (error) {
                logger.error('Failed to copy context:', error);
                uiStore.addToast(t('toast.copyError'), 'error');
            }
        }
    }

    /**
     * Watch for format/settings changes and rebuild context automatically
     */
    function setupAutoRebuild() {
        watch(
            () => [
                settingsStore.settings.context.outputFormat,
                settingsStore.settings.context.stripComments
            ],
            async ([newFormat, newStripComments], [oldFormat, oldStripComments]) => {
                // Only rebuild if context exists and settings actually changed
                if (!contextStore.hasContext || contextStore.isBuilding) return;
                if (newFormat === oldFormat && newStripComments === oldStripComments) return;

                // Get the files from current context summary
                const currentFiles = contextStore.summary?.files;
                if (!currentFiles || currentFiles.length === 0) {
                    // Fallback to selected files if no files in summary
                    if (fileStore.selectedPaths.size === 0) return;
                }

                const filePaths = currentFiles && currentFiles.length > 0
                    ? currentFiles
                    : Array.from(fileStore.selectedPaths);

                try {
                    const options = {
                        maxTokens: settingsStore.settings.context.maxTokens,
                        stripComments: settingsStore.settings.context.stripComments,
                        includeTests: settingsStore.settings.context.includeTests,
                        splitStrategy: settingsStore.settings.context.splitStrategy,
                        outputFormat: settingsStore.settings.context.outputFormat,
                        excludeTests: settingsStore.settings.context.excludeTests,
                        collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
                        stripLicense: settingsStore.settings.context.stripLicense,
                        compactDataFiles: settingsStore.settings.context.compactDataFiles,
                        trimWhitespace: settingsStore.settings.context.trimWhitespace
                    };

                    await contextStore.buildContext(filePaths, options);
                    uiStore.addToast(t('context.rebuilt'), 'success');
                } catch (error) {
                    logger.error('Failed to rebuild context:', error);
                }
            }
        );
    }

    return {
        buildContext,
        copyContext,
        setupAutoRebuild
    };
}
