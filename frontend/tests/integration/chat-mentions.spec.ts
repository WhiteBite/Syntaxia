import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { gitApi } from '@/services/api/git.api'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock localStorage
const localStorageMock = (() => {
    let store: Record<string, string> = {}
    return {
        getItem: (key: string) => store[key] || null,
        setItem: (key: string, value: string) => { store[key] = value },
        removeItem: (key: string) => { delete store[key] },
        clear: () => { store = {} },
        get length() { return Object.keys(store).length },
        key: (index: number) => Object.keys(store)[index] || null
    }
})()

Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock git API
vi.mock('@/services/api/git.api', () => ({
    gitApi: {
        getUncommittedFiles: vi.fn()
    }
}))

// Translation mock
const mockT = (key: string) => {
    const translations: Record<string, string> = {
        'chat.selectFilesHint': 'Select files first',
        'chat.contextAttached': 'Context attached',
        'chat.contextBuildFailed': 'Failed to build context',
        'chat.comingSoon': 'Coming soon',
        'gitContext.noChanges': 'No uncommitted changes'
    }
    return translations[key] || key
}


describe('chat-mentions.integration.spec.ts', () => {
    let projectStore: ReturnType<typeof useProjectStore>
    let fileStore: ReturnType<typeof useFileStore>
    let contextStore: ReturnType<typeof useContextStore>

    beforeEach(() => {
        setActivePinia(createPinia())
        projectStore = useProjectStore()
        useUIStore()
        fileStore = useFileStore()
        contextStore = useContextStore()
        localStorage.clear()
        vi.clearAllMocks()
        projectStore.currentPath = '/test/project'
    })

    describe('@files mention', () => {
        it('should show warning when no files selected', async () => {
            const { useMentions } = await import('@/features/ai-chat')
            const { processFilesMention } = useMentions()
            const result = await processFilesMention(mockT)
            expect(result.success).toBe(false)
            expect(result.type).toBe('files')
        })

        it('should build context from selected files', async () => {
            const buildContextSpy = vi.spyOn(contextStore, 'buildContext').mockResolvedValue()
            fileStore.selectMultiple(['src/app.ts', 'src/utils.ts'])
            const { useMentions } = await import('@/features/ai-chat')
            const { processFilesMention } = useMentions()
            const result = await processFilesMention(mockT)
            expect(buildContextSpy).toHaveBeenCalled()
            expect(result.success).toBe(true)
            expect(result.type).toBe('files')
        })
    })

    describe('@git mention', () => {
        it('should show info when no uncommitted files', async () => {
            vi.mocked(gitApi.getUncommittedFiles).mockResolvedValue([])
            const { useMentions } = await import('@/features/ai-chat')
            const { processGitMention } = useMentions()
            const result = await processGitMention(mockT)
            expect(result.success).toBe(true)
            expect(result.type).toBe('git')
            expect(result.filesAdded).toBe(0)
        })


        it('should build context from uncommitted files', async () => {
            const uncommittedFiles = [
                { path: 'src/modified.ts', status: 'modified' },
                { path: 'src/new.ts', status: 'added' }
            ]
            vi.mocked(gitApi.getUncommittedFiles).mockResolvedValue(uncommittedFiles as never)
            const buildContextSpy = vi.spyOn(contextStore, 'buildContext').mockResolvedValue()
            const { useMentions } = await import('@/features/ai-chat')
            const { processGitMention } = useMentions()
            const result = await processGitMention(mockT)
            expect(buildContextSpy).toHaveBeenCalledWith(['src/modified.ts', 'src/new.ts'])
            expect(result.success).toBe(true)
            expect(result.filesAdded).toBe(2)
        })

        it('should handle git API error', async () => {
            vi.mocked(gitApi.getUncommittedFiles).mockRejectedValue(new Error('Git error'))
            const { useMentions } = await import('@/features/ai-chat')
            const { processGitMention } = useMentions()
            const result = await processGitMention(mockT)
            expect(result.success).toBe(false)
        })
    })

    describe('@problems mention', () => {
        it('should show coming soon message', async () => {
            const { useMentions } = await import('@/features/ai-chat')
            const { processProblemsMention } = useMentions()
            const result = await processProblemsMention(mockT)
            expect(result.success).toBe(false)
            expect(result.type).toBe('problems')
        })
    })


    describe('processMessageMentions', () => {
        it('should process @files mention in message', async () => {
            fileStore.selectMultiple(['src/app.ts'])
            vi.spyOn(contextStore, 'buildContext').mockResolvedValue()
            const { useMentions } = await import('@/features/ai-chat')
            const { processMessageMentions } = useMentions()
            const { cleanedMessage, results } = await processMessageMentions('@files analyze this', mockT)
            expect(cleanedMessage).toBe('analyze this')
            expect(results).toHaveLength(1)
            expect(results[0].type).toBe('files')
        })

        it('should process @git mention in message', async () => {
            vi.mocked(gitApi.getUncommittedFiles).mockResolvedValue([])
            const { useMentions } = await import('@/features/ai-chat')
            const { processMessageMentions } = useMentions()
            const { cleanedMessage, results } = await processMessageMentions('@git show changes', mockT)
            expect(cleanedMessage).toBe('show changes')
            expect(results).toHaveLength(1)
            expect(results[0].type).toBe('git')
        })

        it('should handle message with only mention', async () => {
            fileStore.selectMultiple(['src/app.ts'])
            vi.spyOn(contextStore, 'buildContext').mockResolvedValue()
            const { useMentions } = await import('@/features/ai-chat')
            const { processMessageMentions } = useMentions()
            const { cleanedMessage } = await processMessageMentions('@files', mockT)
            expect(cleanedMessage).toBe('')
        })
    })
})
