import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

// Mock dependencies
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

vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: vi.fn(),
    }),
}))

vi.mock('@/features/files', () => ({
    useFileStore: () => ({
        selectedFilesList: [],
    }),
}))

vi.mock('@/features/context', () => ({
    useContextStore: () => ({
        buildContext: vi.fn(),
        hasContext: false,
    }),
}))

vi.mock('@/services/api.service', () => ({
    filesApi: {
        listFiles: vi.fn().mockResolvedValue([
            { path: 'src/auth/AuthService.ts', children: [] },
            { path: 'src/components/LoginForm.vue', children: [] },
            { path: 'README.md', children: [] },
        ]),
    },
    gitApi: {
        getUncommittedFiles: vi.fn().mockResolvedValue([
            { path: 'src/modified.ts' },
        ]),
        getRichCommitHistory: vi.fn().mockResolvedValue([]),
    },
    semanticApi: {
        search: vi.fn().mockResolvedValue({
            results: [
                { chunk: { filePath: 'src/auth/AuthService.ts', symbolName: 'authenticate' } },
            ],
        }),
    },
}))

import { useMentions } from '@/features/ai-chat/composables/useMentions'

describe('useMentions', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    describe('parseMentions', () => {
        it('should parse file mentions', () => {
            const inputText = ref('@AuthService.ts check this file')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions).toHaveLength(1)
            expect(mentions[0].type).toBe('file')
            expect(mentions[0].value).toBe('AuthService.ts')
        })

        it('should parse folder mentions', () => {
            const inputText = ref('@folder:src/auth/ analyze this')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions).toHaveLength(1)
            expect(mentions[0].type).toBe('folder')
            expect(mentions[0].value).toBe('src/auth/')
        })

        it('should parse symbol mentions', () => {
            const inputText = ref('@symbol:authenticate find usages')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions).toHaveLength(1)
            expect(mentions[0].type).toBe('symbol')
            expect(mentions[0].value).toBe('authenticate')
        })

        it('should parse git mentions', () => {
            const inputText = ref('@git:diff show changes')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions).toHaveLength(1)
            expect(mentions[0].type).toBe('git')
            expect(mentions[0].value).toBe('diff')
        })

        it('should parse docs mention', () => {
            const inputText = ref('@docs update readme')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions).toHaveLength(1)
            expect(mentions[0].type).toBe('docs')
        })

        it('should parse multiple mentions', () => {
            const inputText = ref('@AuthService.ts @LoginForm.vue compare these')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions).toHaveLength(2)
            expect(mentions[0].value).toBe('AuthService.ts')
            expect(mentions[1].value).toBe('LoginForm.vue')
        })

        it('should return empty array for text without mentions', () => {
            const inputText = ref('just regular text without mentions')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions).toHaveLength(0)
        })

        it('should track mention positions correctly', () => {
            const inputText = ref('check @file.ts please')
            const { parseMentions } = useMentions({ inputText })

            const mentions = parseMentions(inputText.value)

            expect(mentions[0].startIndex).toBe(6)
            expect(mentions[0].endIndex).toBe(14) // @file.ts is 8 chars: 6 + 8 = 14
        })
    })

    describe('highlightMentions', () => {
        it('should wrap mentions in span tags', () => {
            const inputText = ref('@file.ts test')
            const { highlightMentions } = useMentions({ inputText })

            const result = highlightMentions(inputText.value)

            expect(result).toContain('<span class="mention mention-file">')
            expect(result).toContain('@file.ts')
            expect(result).toContain('</span>')
        })

        it('should return original text if no mentions', () => {
            const inputText = ref('no mentions here')
            const { highlightMentions } = useMentions({ inputText })

            const result = highlightMentions(inputText.value)

            expect(result).toBe('no mentions here')
        })

        it('should highlight multiple mentions', () => {
            const inputText = ref('@a.ts and @b.ts')
            const { highlightMentions } = useMentions({ inputText })

            const result = highlightMentions(inputText.value)

            const spanCount = (result.match(/<span/g) || []).length
            expect(spanCount).toBe(2)
        })
    })

    describe('insertMention', () => {
        it('should insert mention at cursor position', () => {
            const inputText = ref('check @ please')
            const { insertMention } = useMentions({ inputText })

            const suggestion = {
                type: 'file' as const,
                value: 'test.ts',
                display: 'test.ts',
                icon: 'i-lucide-file',
            }

            const result = insertMention(suggestion, 7)

            expect(result).toContain('@test.ts')
        })

        it('should format folder mentions correctly', () => {
            const inputText = ref('@')
            const { insertMention } = useMentions({ inputText })

            const suggestion = {
                type: 'folder' as const,
                value: 'src/utils',
                display: 'utils',
                icon: 'i-lucide-folder',
            }

            const result = insertMention(suggestion, 1)

            expect(result).toContain('@folder:src/utils')
        })

        it('should format symbol mentions correctly', () => {
            const inputText = ref('@')
            const { insertMention } = useMentions({ inputText })

            const suggestion = {
                type: 'symbol' as const,
                value: 'myFunction',
                display: 'myFunction',
                icon: 'i-lucide-code',
            }

            const result = insertMention(suggestion, 1)

            expect(result).toContain('@symbol:myFunction')
        })
    })

    describe('hasMentions', () => {
        it('should return true when text has mentions', () => {
            const inputText = ref('@file.ts test')
            const { hasMentions } = useMentions({ inputText })

            expect(hasMentions.value).toBe(true)
        })

        it('should return false when text has no mentions', () => {
            const inputText = ref('no mentions')
            const { hasMentions } = useMentions({ inputText })

            expect(hasMentions.value).toBe(false)
        })
    })

    describe('getMentionAtCursor', () => {
        it('should return mention when cursor is inside it', () => {
            const inputText = ref('check @file.ts please')
            const { getMentionAtCursor } = useMentions({ inputText })

            const mention = getMentionAtCursor(10)

            expect(mention).not.toBeNull()
            expect(mention?.value).toBe('file.ts')
        })

        it('should return null when cursor is outside mentions', () => {
            const inputText = ref('check @file.ts please')
            const { getMentionAtCursor } = useMentions({ inputText })

            const mention = getMentionAtCursor(3)

            expect(mention).toBeNull()
        })
    })

    describe('processInput', () => {
        it('should parse and resolve mentions', async () => {
            const inputText = ref('@git:diff check changes')
            const { processInput } = useMentions({ inputText })

            const result = await processInput()

            expect(result.mentions).toHaveLength(1)
            expect(result.text).toBe('check changes')
        })

        it('should strip mentions from text', async () => {
            const inputText = ref('@file.ts @other.ts analyze')
            const { processInput } = useMentions({ inputText })

            const result = await processInput()

            expect(result.text).toBe('analyze')
            expect(result.mentions).toHaveLength(2)
        })
    })
})
