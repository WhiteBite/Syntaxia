import { beforeEach, describe, expect, it, vi } from 'vitest'

// Simple mock state
let mockCustomRules = ''
let mockUpdateCalled = false
let mockUpdateArg = ''
let mockToastCalled = false
let mockToastArgs: [string, string] = ['', '']
let mockRemoveNodeCalled = false
let mockShouldFail = false

// Mock dependencies
vi.mock('@/features/files/model/file.store', () => ({
    useFileStore: () => ({
        removeNode: () => { mockRemoveNodeCalled = true },
    }),
}))

vi.mock('@/stores/settings.store', () => ({
    useSettingsStore: () => ({
        getCustomIgnoreRules: () => mockCustomRules,
        setCustomIgnoreRules: vi.fn(),
    }),
}))

vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: (msg: string, type: string) => {
            mockToastCalled = true
            mockToastArgs = [msg, type]
        },
    }),
}))

vi.mock('@/services/api.service', () => ({
    apiService: {
        updateCustomIgnoreRules: async (rules: string) => {
            if (mockShouldFail) throw new Error('API error')
            mockUpdateCalled = true
            mockUpdateArg = rules
        },
        clearFileTreeCache: vi.fn().mockResolvedValue(undefined),
    },
}))

vi.mock('@/features/files/api/files.api', () => ({
    filesApi: {
        clearCache: vi.fn(),
    },
}))

vi.mock('@/composables/useLogger', () => ({
    useLogger: () => ({
        error: vi.fn(),
    }),
}))

// Import after mocks
import { useIgnoreRules } from '@/features/files/composables/useIgnoreRules'

describe('useIgnoreRules', () => {
    beforeEach(() => {
        mockCustomRules = ''
        mockUpdateCalled = false
        mockUpdateArg = ''
        mockToastCalled = false
        mockToastArgs = ['', '']
        mockRemoveNodeCalled = false
        mockShouldFail = false
    })

    describe('addToIgnore', () => {
        it('should add file to ignore rules', async () => {
            mockCustomRules = 'existing.txt'
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.addToIgnore({
                name: 'test.log',
                path: '/project/test.log',
                isDir: false,
            })

            expect(result).toBe(true)
            expect(mockUpdateCalled).toBe(true)
            expect(mockUpdateArg).toBe('existing.txt\ntest.log')
            expect(mockRemoveNodeCalled).toBe(true)
            expect(mockToastArgs[1]).toBe('success')
        })

        it('should add folder with trailing slash', async () => {
            mockCustomRules = ''
            const ignoreRules = useIgnoreRules()

            await ignoreRules.addToIgnore({
                name: 'node_modules',
                path: '/project/node_modules',
                isDir: true,
            })

            expect(mockUpdateArg).toBe('node_modules/')
        })

        it('should handle empty existing rules', async () => {
            mockCustomRules = ''
            const ignoreRules = useIgnoreRules()

            await ignoreRules.addToIgnore({
                name: 'test.log',
                path: '/project/test.log',
                isDir: false,
            })

            expect(mockUpdateArg).toBe('test.log')
        })

        it('should return false on error', async () => {
            mockShouldFail = true
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.addToIgnore({
                name: 'test.log',
                path: '/project/test.log',
                isDir: false,
            })

            expect(result).toBe(false)
            expect(mockToastArgs[1]).toBe('error')
        })
    })

    describe('removeFromIgnore', () => {
        it('should remove file from ignore rules', async () => {
            mockCustomRules = 'test.log\nother.txt'
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.removeFromIgnore({
                name: 'test.log',
                path: '/project/test.log',
                isDir: false,
            })

            expect(result).toBe(true)
            expect(mockUpdateCalled).toBe(true)
            expect(mockUpdateArg).toBe('other.txt')
            expect(mockToastArgs[1]).toBe('success')
        })

        it('should remove folder with trailing slash', async () => {
            mockCustomRules = 'node_modules/\nsrc/'
            const ignoreRules = useIgnoreRules()

            await ignoreRules.removeFromIgnore({
                name: 'node_modules',
                path: '/project/node_modules',
                isDir: true,
            })

            expect(mockUpdateArg).toBe('src/')
        })

        it('should return false on error', async () => {
            mockCustomRules = 'test.log'
            mockShouldFail = true
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.removeFromIgnore({
                name: 'test.log',
                path: '/project/test.log',
                isDir: false,
            })

            expect(result).toBe(false)
            expect(mockToastArgs[1]).toBe('error')
        })
    })

    describe('getCustomRules', () => {
        it('should return current custom rules', () => {
            mockCustomRules = '*.log\nnode_modules/'
            const ignoreRules = useIgnoreRules()

            const rules = ignoreRules.getCustomRules()

            expect(rules).toBe('*.log\nnode_modules/')
        })
    })

    describe('updateCustomRules', () => {
        it('should update custom rules', async () => {
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.updateCustomRules('new rules')

            expect(result).toBe(true)
            expect(mockUpdateCalled).toBe(true)
            expect(mockUpdateArg).toBe('new rules')
        })

        it('should return false on error', async () => {
            mockShouldFail = true
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.updateCustomRules('new rules')

            expect(result).toBe(false)
        })
    })
})
