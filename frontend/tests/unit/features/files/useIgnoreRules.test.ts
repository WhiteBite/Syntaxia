/**
 * Unit tests for useIgnoreRules.ts composable
 * Tests ignore rules management (gitignore and custom)
 */
import { useIgnoreRules } from '@/features/files/composables/useIgnoreRules'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock logger
vi.mock('@/composables/useLogger', () => ({
    useLogger: () => ({
        debug: vi.fn(),
        info: vi.fn(),
        warn: vi.fn(),
        error: vi.fn(),
    }),
}))

// Shared mock functions
const mockAddToast = vi.fn()
const mockUpdateCustomIgnoreRules = vi.fn().mockResolvedValue(undefined)
const mockClearFileTreeCache = vi.fn().mockResolvedValue(undefined)
const mockClearCache = vi.fn()
const mockRemoveNode = vi.fn().mockReturnValue(true)
const mockSetCustomIgnoreRules = vi.fn()

// State for settings store mock
let customIgnoreRules = ''

// Mock API service
vi.mock('@/services/api.service', async () => {
    return {
        apiService: {
            get updateCustomIgnoreRules() { return mockUpdateCustomIgnoreRules },
            get clearFileTreeCache() { return mockClearFileTreeCache },
        },
    }
})

// Mock files API
vi.mock('@/features/files/api/files.api', async () => {
    return {
        filesApi: {
            get clearCache() { return mockClearCache },
        },
    }
})

// Mock file store
vi.mock('@/features/files/model/file.store', async () => {
    return {
        useFileStore: () => ({
            removeNode: mockRemoveNode,
        }),
    }
})

// Mock settings store
vi.mock('@/stores/settings.store', async () => {
    return {
        useSettingsStore: () => ({
            getCustomIgnoreRules: () => customIgnoreRules,
            setCustomIgnoreRules: (rules: string) => {
                customIgnoreRules = rules
                mockSetCustomIgnoreRules(rules)
            },
        }),
    }
})

// Mock UI store
vi.mock('@/stores/ui.store', async () => {
    return {
        useUIStore: () => ({
            addToast: mockAddToast,
        }),
    }
})

describe('useIgnoreRules', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        customIgnoreRules = ''
        mockUpdateCustomIgnoreRules.mockResolvedValue(undefined)
    })

    describe('addToIgnore', () => {
        it('should add file to ignore rules', async () => {
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            const result = await ignoreRules.addToIgnore(node)

            expect(result).toBe(true)
            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith('test.log')
            expect(customIgnoreRules).toBe('test.log')
        })

        it('should add directory to ignore rules with trailing slash', async () => {
            const ignoreRules = useIgnoreRules()
            const node = { name: 'node_modules', path: '/project/node_modules', isDir: true }

            const result = await ignoreRules.addToIgnore(node)

            expect(result).toBe(true)
            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith('node_modules/')
        })

        it('should append to existing rules', async () => {
            customIgnoreRules = '*.log'
            const ignoreRules = useIgnoreRules()
            const node = { name: 'dist', path: '/project/dist', isDir: true }

            await ignoreRules.addToIgnore(node)

            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith('*.log\ndist/')
        })

        it('should clear caches after adding to ignore', async () => {
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            await ignoreRules.addToIgnore(node)

            expect(mockClearFileTreeCache).toHaveBeenCalled()
            expect(mockClearCache).toHaveBeenCalled()
        })

        it('should show success toast', async () => {
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            await ignoreRules.addToIgnore(node)

            expect(mockAddToast).toHaveBeenCalledWith('Добавлено в исключения', 'success')
        })

        it('should return false and show error toast on failure', async () => {
            mockUpdateCustomIgnoreRules.mockRejectedValueOnce(new Error('API Error'))
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            const result = await ignoreRules.addToIgnore(node)

            expect(result).toBe(false)
            expect(mockAddToast).toHaveBeenCalledWith('Ошибка добавления в исключения', 'error')
        })

        it('should remove node from file tree', async () => {
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            await ignoreRules.addToIgnore(node)

            expect(mockRemoveNode).toHaveBeenCalledWith('/project/test.log')
        })
    })

    describe('removeFromIgnore', () => {
        it('should remove file from ignore rules', async () => {
            customIgnoreRules = 'test.log\n*.tmp'
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            const result = await ignoreRules.removeFromIgnore(node)

            expect(result).toBe(true)
            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith('*.tmp')
        })

        it('should remove directory from ignore rules', async () => {
            customIgnoreRules = 'node_modules/\ndist/'
            const ignoreRules = useIgnoreRules()
            const node = { name: 'node_modules', path: '/project/node_modules', isDir: true }

            const result = await ignoreRules.removeFromIgnore(node)

            expect(result).toBe(true)
            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith('dist/')
        })

        it('should filter out comments when removing', async () => {
            customIgnoreRules = '# Comment\ntest.log\n*.tmp'
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            await ignoreRules.removeFromIgnore(node)

            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith('*.tmp')
        })

        it('should show success toast', async () => {
            customIgnoreRules = 'test.log'
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            await ignoreRules.removeFromIgnore(node)

            expect(mockAddToast).toHaveBeenCalledWith('Удалено из исключений', 'success')
        })

        it('should return false and show error toast on failure', async () => {
            mockUpdateCustomIgnoreRules.mockRejectedValueOnce(new Error('API Error'))
            customIgnoreRules = 'test.log'
            const ignoreRules = useIgnoreRules()
            const node = { name: 'test.log', path: '/project/test.log', isDir: false }

            const result = await ignoreRules.removeFromIgnore(node)

            expect(result).toBe(false)
            expect(mockAddToast).toHaveBeenCalledWith('Ошибка удаления из исключений', 'error')
        })
    })

    describe('getCustomRules', () => {
        it('should return current custom rules', () => {
            customIgnoreRules = '*.log\nnode_modules/'
            const ignoreRules = useIgnoreRules()

            const rules = ignoreRules.getCustomRules()

            expect(rules).toBe('*.log\nnode_modules/')
        })

        it('should return empty string when no rules', () => {
            customIgnoreRules = ''
            const ignoreRules = useIgnoreRules()

            const rules = ignoreRules.getCustomRules()

            expect(rules).toBe('')
        })
    })

    describe('updateCustomRules', () => {
        it('should update custom rules', async () => {
            const ignoreRules = useIgnoreRules()
            const newRules = '*.log\n*.tmp\nnode_modules/'

            const result = await ignoreRules.updateCustomRules(newRules)

            expect(result).toBe(true)
            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith(newRules)
            expect(mockSetCustomIgnoreRules).toHaveBeenCalledWith(newRules)
        })

        it('should return false on failure', async () => {
            mockUpdateCustomIgnoreRules.mockRejectedValueOnce(new Error('API Error'))
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.updateCustomRules('*.log')

            expect(result).toBe(false)
        })

        it('should handle empty rules', async () => {
            const ignoreRules = useIgnoreRules()

            const result = await ignoreRules.updateCustomRules('')

            expect(result).toBe(true)
            expect(mockUpdateCustomIgnoreRules).toHaveBeenCalledWith('')
        })
    })
})
