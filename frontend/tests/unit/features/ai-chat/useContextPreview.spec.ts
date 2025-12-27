import { useContextPreview, type ContextFilePreview } from '@/features/ai-chat/composables/useContextPreview'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock template store
vi.mock('@/features/templates/model/template.store', () => ({
    useTemplateStore: () => ({
        activeTemplate: {
            name: 'Test Template',
            icon: '🧪',
            roleContent: 'You are a test assistant',
            rulesContent: 'Follow test rules'
        },
        userRules: ''
    })
}))

// Mock context utils
vi.mock('@/features/context/lib/context-utils', () => ({
    estimateTokens: (text: string) => Math.ceil(text.length / 4)
}))

// Mock localStorage
const localStorageMock = {
    store: {} as Record<string, string>,
    getItem: vi.fn((key: string) => localStorageMock.store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
        localStorageMock.store[key] = value
    }),
    clear: vi.fn(() => {
        localStorageMock.store = {}
    }),
}
Object.defineProperty(global, 'localStorage', { value: localStorageMock })

describe('useContextPreview', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorageMock.clear()
        vi.clearAllMocks()
    })

    it('should initialize with empty files', () => {
        const { files } = useContextPreview()
        expect(files.value).toEqual([])
    })

    it('should initialize with default token limit', () => {
        const { tokenLimit } = useContextPreview()
        expect(tokenLimit.value).toBe(128000)
    })

    it('should initialize with showPrompt false', () => {
        const { showPrompt } = useContextPreview()
        expect(showPrompt.value).toBe(false)
    })

    it('should calculate selected tokens correctly', () => {
        const { setFiles, selectedTokens } = useContextPreview()
        setFiles([
            { path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true },
            { path: 'b.ts', tokens: 200, relevance: 1, reason: '', selected: false },
        ])
        expect(selectedTokens.value).toBe(100)
    })

    it('should calculate total tokens correctly', () => {
        const { setFiles, totalTokens } = useContextPreview()
        setFiles([
            { path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true },
            { path: 'b.ts', tokens: 200, relevance: 1, reason: '', selected: false },
        ])
        expect(totalTokens.value).toBe(300)
    })

    it('should toggle file selection', () => {
        const { setFiles, toggleFile, files } = useContextPreview()
        setFiles([{ path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: false }])

        toggleFile('a.ts')
        expect(files.value[0].selected).toBe(true)

        toggleFile('a.ts')
        expect(files.value[0].selected).toBe(false)
    })

    it('should detect over limit', () => {
        const { setFiles, setTokenLimit, isOverLimit } = useContextPreview()
        setTokenLimit(100)
        setFiles([{ path: 'a.ts', tokens: 200, relevance: 1, reason: '', selected: true }])
        expect(isOverLimit.value).toBe(true)
    })

    it('should not be over limit when under', () => {
        const { setFiles, setTokenLimit, isOverLimit } = useContextPreview()
        setTokenLimit(500)
        setFiles([{ path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true }])
        expect(isOverLimit.value).toBe(false)
    })

    it('should select all files', () => {
        const { setFiles, selectAll, files } = useContextPreview()
        setFiles([
            { path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: false },
            { path: 'b.ts', tokens: 200, relevance: 1, reason: '', selected: false },
        ])

        selectAll()

        expect(files.value[0].selected).toBe(true)
        expect(files.value[1].selected).toBe(true)
    })

    it('should deselect all files', () => {
        const { setFiles, deselectAll, files } = useContextPreview()
        setFiles([
            { path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true },
            { path: 'b.ts', tokens: 200, relevance: 1, reason: '', selected: true },
        ])

        deselectAll()

        expect(files.value[0].selected).toBe(false)
        expect(files.value[1].selected).toBe(false)
    })

    it('should get selected paths', () => {
        const { setFiles, getSelectedPaths } = useContextPreview()
        setFiles([
            { path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true },
            { path: 'b.ts', tokens: 200, relevance: 1, reason: '', selected: false },
            { path: 'c.ts', tokens: 150, relevance: 1, reason: '', selected: true },
        ])

        const paths = getSelectedPaths()
        expect(paths).toEqual(['a.ts', 'c.ts'])
    })

    it('should toggle prompt expanded state', () => {
        const { showPrompt, togglePromptExpanded } = useContextPreview()
        expect(showPrompt.value).toBe(false)

        togglePromptExpanded()
        expect(showPrompt.value).toBe(true)

        togglePromptExpanded()
        expect(showPrompt.value).toBe(false)
    })

    it('should reset files but preserve showPrompt (persisted state)', () => {
        const preview = useContextPreview()
        preview.setFiles([{ path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true }])
        const showPromptBefore = preview.showPrompt.value

        preview.reset()

        expect(preview.files.value).toEqual([])
        // showPrompt is persisted in localStorage, so reset() doesn't change it
        expect(preview.showPrompt.value).toBe(showPromptBefore)
    })

    it('should calculate selected count', () => {
        const { setFiles, selectedCount } = useContextPreview()
        setFiles([
            { path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true },
            { path: 'b.ts', tokens: 200, relevance: 1, reason: '', selected: false },
            { path: 'c.ts', tokens: 150, relevance: 1, reason: '', selected: true },
        ])

        expect(selectedCount.value).toBe(2)
    })

    it('should calculate token usage percent', () => {
        const { setFiles, setTokenLimit, tokenUsagePercent } = useContextPreview()
        setTokenLimit(1000)
        setFiles([{ path: 'a.ts', tokens: 250, relevance: 1, reason: '', selected: true }])

        // 250 tokens + prompt tokens (from mock) / 1000 * 100
        expect(tokenUsagePercent.value).toBeGreaterThan(0)
        expect(tokenUsagePercent.value).toBeLessThanOrEqual(100)
    })

    it('should cap token usage percent at 100', () => {
        const { setFiles, setTokenLimit, tokenUsagePercent } = useContextPreview()
        setTokenLimit(100)
        setFiles([{ path: 'a.ts', tokens: 500, relevance: 1, reason: '', selected: true }])

        expect(tokenUsagePercent.value).toBe(100)
    })

    it('should determine canSend based on selection and limit', () => {
        const { setFiles, setTokenLimit, canSend } = useContextPreview()

        // No files selected - cannot send
        setFiles([{ path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: false }])
        expect(canSend.value).toBe(false)

        // Files selected, under limit - can send
        setTokenLimit(10000)
        setFiles([{ path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true }])
        expect(canSend.value).toBe(true)

        // Files selected, over limit - cannot send
        setTokenLimit(50)
        expect(canSend.value).toBe(false)
    })

    it('should compute system prompt from template', () => {
        const { systemPrompt } = useContextPreview()

        expect(systemPrompt.value.name).toBe('Test Template')
        expect(systemPrompt.value.icon).toBe('🧪')
        expect(systemPrompt.value.content).toContain('Role')
        expect(systemPrompt.value.tokens).toBeGreaterThan(0)
    })

    it('should include prompt tokens in totalWithPrompt', () => {
        const { setFiles, totalWithPrompt, selectedTokens, promptTokens } = useContextPreview()
        setFiles([{ path: 'a.ts', tokens: 100, relevance: 1, reason: '', selected: true }])

        expect(totalWithPrompt.value).toBe(selectedTokens.value + promptTokens.value)
    })

    it('should return selected files', () => {
        const { setFiles, selectedFiles } = useContextPreview()
        const testFiles: ContextFilePreview[] = [
            { path: 'a.ts', tokens: 100, relevance: 1, reason: 'test', selected: true },
            { path: 'b.ts', tokens: 200, relevance: 0.5, reason: 'test2', selected: false },
        ]
        setFiles(testFiles)

        expect(selectedFiles.value).toHaveLength(1)
        expect(selectedFiles.value[0].path).toBe('a.ts')
    })
})
