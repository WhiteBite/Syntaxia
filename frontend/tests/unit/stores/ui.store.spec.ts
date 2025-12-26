import { useUIStore, type Toast, type ToastAction } from '@/stores/ui.store'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

describe('UIStore', () => {
    let store: ReturnType<typeof useUIStore>
    let consoleSpy: {
        log: ReturnType<typeof vi.spyOn>
        info: ReturnType<typeof vi.spyOn>
        warn: ReturnType<typeof vi.spyOn>
        error: ReturnType<typeof vi.spyOn>
    }

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useUIStore()

        consoleSpy = {
            log: vi.spyOn(console, 'log').mockImplementation(() => { }),
            info: vi.spyOn(console, 'info').mockImplementation(() => { }),
            warn: vi.spyOn(console, 'warn').mockImplementation(() => { }),
            error: vi.spyOn(console, 'error').mockImplementation(() => { })
        }

        vi.useFakeTimers()
    })

    afterEach(() => {
        vi.restoreAllMocks()
        vi.useRealTimers()
    })

    describe('initial state', () => {
        it('should have empty toasts array', () => {
            expect(store.toasts).toEqual([])
        })

        it('should have showSettingsModal as false', () => {
            expect(store.showSettingsModal).toBe(false)
        })

        it('should have showKeyboardShortcutsModal as false', () => {
            expect(store.showKeyboardShortcutsModal).toBe(false)
        })
    })

    describe('addToast', () => {
        it('should add toast with default type info', () => {
            store.addToast('Test message')

            expect(store.toasts).toHaveLength(1)
            expect(store.toasts[0].message).toBe('Test message')
            expect(store.toasts[0].type).toBe('info')
        })

        it('should add toast with specified type', () => {
            store.addToast('Success message', 'success')
            store.addToast('Error message', 'error')
            store.addToast('Warning message', 'warning')

            expect(store.toasts).toHaveLength(3)
            expect(store.toasts[0].type).toBe('success')
            expect(store.toasts[1].type).toBe('error')
            expect(store.toasts[2].type).toBe('warning')
        })

        it('should generate unique toast ids', () => {
            store.addToast('Message 1')
            store.addToast('Message 2')
            store.addToast('Message 3')

            const ids = store.toasts.map(t => t.id)
            const uniqueIds = new Set(ids)
            expect(uniqueIds.size).toBe(3)
        })

        it('should return toast id', () => {
            const id = store.addToast('Test message')

            expect(id).toBe('toast-0')
        })

        it('should add toast with action', () => {
            const action: ToastAction = {
                label: 'Undo',
                icon: 'undo',
                onClick: vi.fn()
            }

            store.addToast('Action message', 'info', 3000, action)

            expect(store.toasts[0].action).toEqual(action)
        })

        it('should auto-remove toast after duration', () => {
            store.addToast('Auto-remove message', 'info', 1000)

            expect(store.toasts).toHaveLength(1)

            vi.advanceTimersByTime(1000)

            expect(store.toasts).toHaveLength(0)
        })

        it('should not auto-remove toast when duration is 0', () => {
            store.addToast('Persistent message', 'info', 0)

            expect(store.toasts).toHaveLength(1)

            vi.advanceTimersByTime(10000)

            expect(store.toasts).toHaveLength(1)
        })

        it('should log error toast to console.error', () => {
            store.addToast('Error occurred', 'error')

            expect(consoleSpy.error).toHaveBeenCalledWith('[UIStore]', '[Toast ERROR] Error occurred')
        })

        it('should log warning toast to console.warn', () => {
            store.addToast('Warning message', 'warning')

            expect(consoleSpy.warn).toHaveBeenCalledWith('[UIStore]', '[Toast WARNING] Warning message')
        })

        it('should log info toast to console.info', () => {
            store.addToast('Info message', 'info')

            expect(consoleSpy.info).toHaveBeenCalledWith('[UIStore]', '[Toast INFO] Info message')
        })
    })

    describe('removeToast', () => {
        it('should remove toast by id', () => {
            const id = store.addToast('Test message')

            expect(store.toasts).toHaveLength(1)

            store.removeToast(id)

            expect(store.toasts).toHaveLength(0)
        })

        it('should not throw when removing non-existent toast', () => {
            expect(() => store.removeToast('non-existent-id')).not.toThrow()
        })

        it('should only remove specified toast', () => {
            store.addToast('Message 1')
            const id2 = store.addToast('Message 2')
            store.addToast('Message 3')

            store.removeToast(id2)

            expect(store.toasts).toHaveLength(2)
            expect(store.toasts.find(t => t.message === 'Message 2')).toBeUndefined()
        })
    })

    describe('clearToasts', () => {
        it('should remove all toasts', () => {
            store.addToast('Message 1')
            store.addToast('Message 2')
            store.addToast('Message 3')

            expect(store.toasts).toHaveLength(3)

            store.clearToasts()

            expect(store.toasts).toHaveLength(0)
        })

        it('should work on empty toasts array', () => {
            expect(() => store.clearToasts()).not.toThrow()
            expect(store.toasts).toHaveLength(0)
        })
    })

    describe('openSettingsModal', () => {
        it('should set showSettingsModal to true', () => {
            expect(store.showSettingsModal).toBe(false)

            store.openSettingsModal()

            expect(store.showSettingsModal).toBe(true)
        })
    })

    describe('openKeyboardShortcutsModal', () => {
        it('should set showKeyboardShortcutsModal to true', () => {
            expect(store.showKeyboardShortcutsModal).toBe(false)

            store.openKeyboardShortcutsModal()

            expect(store.showKeyboardShortcutsModal).toBe(true)
        })
    })

    describe('toast types', () => {
        it('should correctly type Toast interface', () => {
            const toast: Toast = {
                id: 'test-id',
                message: 'Test message',
                type: 'success',
                duration: 3000,
                action: {
                    label: 'Action',
                    onClick: () => { }
                }
            }

            expect(toast.id).toBe('test-id')
            expect(toast.type).toBe('success')
        })
    })
})
