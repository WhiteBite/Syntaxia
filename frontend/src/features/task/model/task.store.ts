import { useFileStore } from '@/features/files'
import { apiService } from '@/services/api.service'
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export interface TaskDraft {
    description: string
    type?: string
    timestamp: string
}

export interface AnalysisResult {
    suggestedFiles: string[]
    complexity: 'low' | 'medium' | 'high'
    estimatedTime?: number
    recommendations: string[]
}

export const useTaskStore = defineStore('task', () => {
    // State
    const taskDescription = ref('')
    const taskType = ref<string>('feature')
    const isAnalyzing = ref(false)
    const analysisResult = ref<AnalysisResult | null>(null)
    const suggestions = ref<string[]>([
        'Add user authentication',
        'Implement data caching',
        'Refactor API calls',
        'Add unit tests'
    ])
    const error = ref<string | null>(null)

    // Auto-save to localStorage
    let saveTimeout: number | null = null
    watch(taskDescription, () => {
        if (saveTimeout) clearTimeout(saveTimeout)
        saveTimeout = setTimeout(() => {
            saveTaskDraft()
        }, 500) as unknown as number
    })

    // Actions
    async function analyzeTask() {
        if (taskDescription.value.trim().length < 10) {
            error.value = 'Task description is too short'
            return
        }

        isAnalyzing.value = true
        error.value = null

        try {
            // Get current project path and file tree
            const projectPath = await apiService.getCurrentDirectory()
            const fileStore = useFileStore()

            // Load file tree if not already loaded
            if (fileStore.nodes.length === 0) {
                await fileStore.loadFileTree(projectPath)
            }

            // Convert file tree to domain format for API
            const domainFiles = convertToDomainFiles(fileStore.nodes)

            // Call backend API to get suggested files
            const suggestedFilePaths = await apiService.suggestContextFiles(
                taskDescription.value,
                domainFiles as unknown as import('#wailsjs/go/models').domain.FileNode[]
            )

            // Calculate complexity based on suggested files
            let complexity: 'low' | 'medium' | 'high' = 'low'
            if (suggestedFilePaths.length > 10) {
                complexity = 'high'
            } else if (suggestedFilePaths.length > 5) {
                complexity = 'medium'
            }

            // Estimate time based on complexity and file count
            const estimatedTime = suggestedFilePaths.length * 5 + (complexity === 'high' ? 30 : complexity === 'medium' ? 15 : 5)

            analysisResult.value = {
                suggestedFiles: suggestedFilePaths,
                complexity,
                estimatedTime,
                recommendations: generateRecommendations(complexity, suggestedFilePaths.length)
            }

            // Auto-select suggested files in file store
            if (suggestedFilePaths.length > 0) {
                fileStore.clearSelection()
                fileStore.selectMultiple(suggestedFilePaths)
            }
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to analyze task'
            throw err
        } finally {
            isAnalyzing.value = false
        }
    }

    interface TaskFileNode {
        name: string
        path: string
        relPath: string
        isDir: boolean
        children?: TaskFileNode[]
        isGitignored: boolean
        isCustomIgnored: boolean
        isIgnored: boolean
        size: number
    }

    interface InputNode {
        name: string
        path: string
        isDir: boolean
        children?: InputNode[]
        isIgnored?: boolean
        size?: number
    }

    function convertToDomainFiles(nodes: InputNode[]): TaskFileNode[] {
        return nodes.map(node => ({
            name: node.name,
            path: node.path,
            relPath: node.path,
            isDir: node.isDir,
            children: node.children ? convertToDomainFiles(node.children) : undefined,
            isGitignored: node.isIgnored || false,
            isCustomIgnored: false,
            isIgnored: node.isIgnored || false,
            size: node.size || 0
        }))
    }

    function generateRecommendations(complexity: string, fileCount: number): string[] {
        const recommendations: string[] = []

        if (complexity === 'high') {
            recommendations.push('Consider breaking down into smaller tasks')
            recommendations.push('This is a complex task that may require multiple iterations')
        }

        if (fileCount > 15) {
            recommendations.push('Large number of files involved - ensure proper testing')
        }

        if (fileCount === 0) {
            recommendations.push('No relevant files found - try refining your task description')
        } else {
            recommendations.push(`${fileCount} relevant files identified`)
        }

        recommendations.push('Review existing similar implementations')
        recommendations.push('Add tests for new functionality')

        return recommendations
    }

    function saveTaskDraft() {
        try {
            const draft: TaskDraft = {
                description: taskDescription.value,
                type: taskType.value,
                timestamp: new Date().toISOString()
            }
            localStorage.setItem('task-draft', JSON.stringify(draft))
        } catch (err) {
            console.warn('Failed to save task draft:', err)
        }
    }

    function loadTaskDraft() {
        try {
            const draftJson = localStorage.getItem('task-draft')
            if (draftJson) {
                const draft: TaskDraft = JSON.parse(draftJson)
                taskDescription.value = draft.description
                taskType.value = draft.type || 'feature'
            }
        } catch (err) {
            console.warn('Failed to load task draft:', err)
        }
    }

    function clearTask() {
        taskDescription.value = ''
        taskType.value = 'feature'
        analysisResult.value = null
        error.value = null
        localStorage.removeItem('task-draft')
    }

    function applySuggestion(suggestion: string) {
        taskDescription.value = suggestion
    }

    return {
        // State
        taskDescription,
        taskType,
        isAnalyzing,
        analysisResult,
        suggestions,
        error,
        // Actions
        analyzeTask,
        saveTaskDraft,
        loadTaskDraft,
        clearTask,
        applySuggestion
    }
})
