import { useAnalysisStatus } from '@/features/files/composables/useAnalysisStatus'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

// Mock stores and services
vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => key,
    }),
}))

vi.mock('@/stores/project.store', () => ({
    useProjectStore: () => ({
        currentPath: '/test/project',
    }),
}))

const mockGetSmartSuggestions = vi.fn()
const mockGetImpactPreview = vi.fn()

vi.mock('@/services/api.service', () => ({
    apiService: {
        getSmartSuggestions: (...args: unknown[]) => mockGetSmartSuggestions(...args),
        getImpactPreview: (...args: unknown[]) => mockGetImpactPreview(...args),
    },
}))

describe('useAnalysisStatus', () => {
    const selectedFiles = ref<string[]>([])
    const mockOnAddFiles = vi.fn()

    beforeEach(() => {
        vi.clearAllMocks()
        vi.useFakeTimers()
        selectedFiles.value = []
        mockGetSmartSuggestions.mockResolvedValue({
            suggestions: [
                { path: 'src/utils.ts', source: 'git', score: 0.9 },
                { path: 'src/helpers.ts', source: 'arch', score: 0.8 },
            ]
        })
        mockGetImpactPreview.mockResolvedValue({
            totalDependents: 5,
            aggregateRisk: 0.3,
            riskLevel: 'low',
            affectedFiles: [
                { path: 'src/app.ts', type: 'direct', dependents: 2 },
            ],
            relatedTests: ['test/app.spec.ts'],
        })
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it('should initialize with default state', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.suggestions.value).toEqual([])
        expect(status.impactResult.value).toBeNull()
        expect(status.isLoadingRelated.value).toBe(false)
        expect(status.showRelatedPopup.value).toBe(false)
        expect(status.showImpactPopup.value).toBe(false)
    })

    it('should not show bar when no files selected', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.shouldShowBar.value).toBe(false)
    })

    it('should fetch suggestions when files are selected', async () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        expect(mockGetSmartSuggestions).toHaveBeenCalledWith('/test/project', ['src/app.ts'])
        expect(status.suggestions.value).toHaveLength(2)
    })

    it('should fetch impact when files are selected', async () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        expect(mockGetImpactPreview).toHaveBeenCalledWith('/test/project', ['src/app.ts'])
        expect(status.impactResult.value).not.toBeNull()
        expect(status.dependentCount.value).toBe(5)
    })

    it('should clear data when files are deselected', async () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        expect(status.suggestions.value).toHaveLength(2)

        selectedFiles.value = []
        await vi.advanceTimersByTimeAsync(100)

        expect(status.suggestions.value).toEqual([])
        expect(status.impactResult.value).toBeNull()
    })

    it('should toggle related file selection', async () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        // Initially all suggestions are selected
        expect(status.selectedRelated.value.has('src/utils.ts')).toBe(true)

        status.toggleRelated('src/utils.ts')
        expect(status.selectedRelated.value.has('src/utils.ts')).toBe(false)

        status.toggleRelated('src/utils.ts')
        expect(status.selectedRelated.value.has('src/utils.ts')).toBe(true)
    })

    it('should add selected related files', async () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        status.addSelectedRelated()

        expect(mockOnAddFiles).toHaveBeenCalledWith(['src/utils.ts', 'src/helpers.ts'])
        expect(status.selectedRelated.value.size).toBe(0)
        expect(status.showRelatedPopup.value).toBe(false)
    })

    it('should not add files when none selected', async () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        // Deselect all
        status.toggleRelated('src/utils.ts')
        status.toggleRelated('src/helpers.ts')

        status.addSelectedRelated()

        expect(mockOnAddFiles).not.toHaveBeenCalled()
    })

    it('should return correct source labels', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getSourceLabel('git')).toBe('context.sourceGitShort')
        expect(status.getSourceLabel('arch')).toBe('context.sourceArchShort')
        expect(status.getSourceLabel('semantic')).toBe('context.sourceSemanticShort')
        expect(status.getSourceLabel('unknown')).toBe('')
    })

    it('should return correct source badge classes', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getSourceBadgeClass('git')).toBe('badge-git')
        expect(status.getSourceBadgeClass('arch')).toBe('badge-arch')
        expect(status.getSourceBadgeClass('semantic')).toBe('badge-semantic')
        expect(status.getSourceBadgeClass('unknown')).toBe('badge-default')
    })

    it('should return correct file icon classes', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getFileIconClass('app.vue')).toBe('icon-vue')
        expect(status.getFileIconClass('utils.ts')).toBe('icon-ts')
        expect(status.getFileIconClass('main.go')).toBe('icon-go')
        expect(status.getFileIconClass('script.py')).toBe('icon-py')
        expect(status.getFileIconClass('unknown.xyz')).toBe('icon-default')
    })

    it('should extract file name from path', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getFileName('src/components/Button.vue')).toBe('Button.vue')
        expect(status.getFileName('app.ts')).toBe('app.ts')
    })

    it('should extract file path without name', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getFilePath('src/components/Button.vue')).toBe('src/components')
        expect(status.getFilePath('app.ts')).toBe('')
    })

    it('should return correct risk classes', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getRiskClass('low')).toBe('risk-low')
        expect(status.getRiskClass('medium')).toBe('risk-medium')
        expect(status.getRiskClass('high')).toBe('risk-high')
        expect(status.getRiskClass('unknown')).toBe('')
    })

    it('should return correct risk bar classes', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getRiskBarClass('low')).toBe('bg-green-500')
        expect(status.getRiskBarClass('medium')).toBe('bg-amber-500')
        expect(status.getRiskBarClass('high')).toBe('bg-red-500')
        expect(status.getRiskBarClass('unknown')).toBe('bg-gray-500')
    })

    it('should return correct risk labels', () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.getRiskLabel('low')).toBe('context.riskLow')
        expect(status.getRiskLabel('medium')).toBe('context.riskMedium')
        expect(status.getRiskLabel('high')).toBe('context.riskHigh')
        expect(status.getRiskLabel('unknown')).toBe('')
    })

    it('should handle API errors gracefully', async () => {
        mockGetSmartSuggestions.mockRejectedValue(new Error('API Error'))
        mockGetImpactPreview.mockRejectedValue(new Error('API Error'))

        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        expect(status.suggestions.value).toEqual([])
        expect(status.impactResult.value).toBeNull()
    })

    it('should handle timeout gracefully', async () => {
        mockGetSmartSuggestions.mockImplementation(
            () => new Promise(resolve => setTimeout(() => resolve({ suggestions: [] }), 10000))
        )

        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(6000)

        expect(status.suggestions.value).toEqual([])
    })

    it('should compute relatedCount correctly', async () => {
        const status = useAnalysisStatus({
            selectedFiles,
            onAddFiles: mockOnAddFiles,
        })

        expect(status.relatedCount.value).toBe(0)

        selectedFiles.value = ['src/app.ts']
        await vi.advanceTimersByTimeAsync(600)

        expect(status.relatedCount.value).toBe(2)
    })
})
