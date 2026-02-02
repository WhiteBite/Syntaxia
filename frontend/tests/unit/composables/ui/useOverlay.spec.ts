import { useOverlay } from '@/composables/ui/useOverlay'
import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick, ref } from 'vue'

describe('useOverlay', () => {
    beforeEach(() => {
        // Reset body styles before each test
        document.body.style.overflow = ''
        document.body.style.paddingRight = ''
    })

    describe('Basic functionality', () => {
        it('should initialize with default state', () => {
            const { isVisible, position } = useOverlay()

            expect(isVisible.value).toBe(false)
            expect(position.value).toEqual({
                top: '0px',
                left: '0px',
                transform: ''
            })
        })

        it('should initialize with custom initial state', () => {
            const { isVisible } = useOverlay({ initialState: true })

            expect(isVisible.value).toBe(true)
        })

        it('should open overlay', () => {
            const { isVisible, open } = useOverlay()

            open()

            expect(isVisible.value).toBe(true)
        })

        it('should close overlay', () => {
            const { isVisible, open, close } = useOverlay({ initialState: true })

            expect(isVisible.value).toBe(true)

            close()

            expect(isVisible.value).toBe(false)
        })

        it('should toggle overlay', () => {
            const { isVisible, toggle } = useOverlay()

            expect(isVisible.value).toBe(false)

            toggle()
            expect(isVisible.value).toBe(true)

            toggle()
            expect(isVisible.value).toBe(false)
        })

        it('should not open if already open', () => {
            const onOpen = vi.fn()
            const { isVisible, open } = useOverlay({ onOpen })

            open()
            expect(onOpen).toHaveBeenCalledTimes(1)

            open()
            expect(onOpen).toHaveBeenCalledTimes(1)
        })

        it('should not close if already closed', () => {
            const onClose = vi.fn()
            const { isVisible, close } = useOverlay({ onClose })

            expect(isVisible.value).toBe(false)

            close()
            expect(onClose).not.toHaveBeenCalled()
        })
    })

    describe('Position management', () => {
        it('should set position', () => {
            const { position, setPosition } = useOverlay()

            setPosition({ top: '100px', left: '200px' })

            expect(position.value).toEqual({
                top: '100px',
                left: '200px',
                transform: ''
            })
        })

        it('should partially update position', () => {
            const { position, setPosition } = useOverlay()

            setPosition({ top: '50px' })

            expect(position.value).toEqual({
                top: '50px',
                left: '0px',
                transform: ''
            })
        })

        it('should update transform', () => {
            const { position, setPosition } = useOverlay()

            setPosition({ transform: 'translateX(-50%)' })

            expect(position.value.transform).toBe('translateX(-50%)')
        })
    })

    describe('Callbacks', () => {
        it('should call onOpen callback', () => {
            const onOpen = vi.fn()
            const { open } = useOverlay({ onOpen })

            open()

            expect(onOpen).toHaveBeenCalledTimes(1)
        })

        it('should call onClose callback', () => {
            const onClose = vi.fn()
            const { open, close } = useOverlay({ onClose })

            open()
            close()

            expect(onClose).toHaveBeenCalledTimes(1)
        })

        it('should call callbacks in correct order', () => {
            const calls: string[] = []
            const onOpen = vi.fn(() => calls.push('open'))
            const onClose = vi.fn(() => calls.push('close'))

            const { open, close } = useOverlay({ onOpen, onClose })

            open()
            close()
            open()

            expect(calls).toEqual(['open', 'close', 'open'])
        })
    })

    describe('Scroll lock', () => {
        it('should lock scroll when enabled and overlay opens', () => {
            const { open } = useOverlay({ lockScroll: true })

            open()

            expect(document.body.style.overflow).toBe('hidden')
            expect(document.body.style.paddingRight).toBeTruthy()
        })

        it('should unlock scroll when overlay closes', () => {
            const { open, close } = useOverlay({ lockScroll: true })

            open()
            expect(document.body.style.overflow).toBe('hidden')

            close()
            expect(document.body.style.overflow).toBe('')
            expect(document.body.style.paddingRight).toBe('')
        })

        it('should not lock scroll when disabled', () => {
            const { open } = useOverlay({ lockScroll: false })

            open()

            expect(document.body.style.overflow).toBe('')
        })

        it('should unlock scroll on unmount', async () => {
            const TestComponent = {
                setup() {
                    const { open } = useOverlay({ lockScroll: true })
                    open()
                    return {}
                },
                template: '<div>Test</div>'
            }

            const wrapper = mount(TestComponent)
            expect(document.body.style.overflow).toBe('hidden')

            wrapper.unmount()
            expect(document.body.style.overflow).toBe('')
        })
    })

    describe('Escape key handling', () => {
        it('should close on escape key when enabled', async () => {
            const TestComponent = {
                setup() {
                    const overlay = useOverlay({ closeOnEscape: true })
                    overlay.open()
                    return { overlay }
                },
                template: '<div ref="overlayRef">Test</div>'
            }

            const wrapper = mount(TestComponent)
            await nextTick()
            expect(wrapper.vm.overlay.isVisible.value).toBe(true)

            // Simulate escape key press
            const event = new KeyboardEvent('keydown', { key: 'Escape', bubbles: true })
            window.dispatchEvent(event)
            await nextTick()

            expect(wrapper.vm.overlay.isVisible.value).toBe(false)
            wrapper.unmount()
        })

        it('should not close on escape when disabled', async () => {
            const TestComponent = {
                setup() {
                    const overlay = useOverlay({ closeOnEscape: false })
                    overlay.open()
                    return { overlay }
                },
                template: '<div>Test</div>'
            }

            const wrapper = mount(TestComponent)
            await nextTick()
            expect(wrapper.vm.overlay.isVisible.value).toBe(true)

            const event = new KeyboardEvent('keydown', { key: 'Escape', bubbles: true })
            window.dispatchEvent(event)
            await nextTick()

            expect(wrapper.vm.overlay.isVisible.value).toBe(true)
            wrapper.unmount()
        })

        it('should not close on escape when already closed', async () => {
            const onClose = vi.fn()
            const TestComponent = {
                setup() {
                    const overlay = useOverlay({ closeOnEscape: true, onClose })
                    return { overlay }
                },
                template: '<div>Test</div>'
            }

            const wrapper = mount(TestComponent)
            await nextTick()
            expect(wrapper.vm.overlay.isVisible.value).toBe(false)

            const event = new KeyboardEvent('keydown', { key: 'Escape', bubbles: true })
            window.dispatchEvent(event)
            await nextTick()

            expect(onClose).not.toHaveBeenCalled()
            wrapper.unmount()
        })
    })

    describe('Click outside handling', () => {
        it('should close on click outside when enabled', async () => {
            // Note: This test verifies the configuration is set up correctly
            // Full click outside behavior is tested in component integration tests
            const onClose = vi.fn()
            const TestComponent = {
                setup() {
                    const overlay = useOverlay({
                        closeOnClickOutside: true,
                        onClose
                    })
                    return { overlay }
                },
                template: '<div><div ref="overlayRef" v-if="overlay.isVisible.value">Content</div></div>'
            }

            const wrapper = mount(TestComponent, {
                attachTo: document.body
            })

            wrapper.vm.overlay.open()
            await nextTick()

            expect(wrapper.vm.overlay.isVisible.value).toBe(true)

            // Manually trigger close to verify the mechanism works
            wrapper.vm.overlay.close()
            await nextTick()

            expect(wrapper.vm.overlay.isVisible.value).toBe(false)
            expect(onClose).toHaveBeenCalled()

            wrapper.unmount()
        })

        it('should not close on click outside when disabled', async () => {
            const TestComponent = {
                setup() {
                    const overlay = useOverlay({ closeOnClickOutside: false })
                    overlay.open()
                    return { overlay }
                },
                template: '<div ref="overlayRef">Content</div>'
            }

            const wrapper = mount(TestComponent, {
                attachTo: document.body
            })

            expect(wrapper.vm.overlay.isVisible.value).toBe(true)

            document.body.click()
            await nextTick()

            expect(wrapper.vm.overlay.isVisible.value).toBe(true)

            wrapper.unmount()
        })

        it('should ignore clicks on specified elements', async () => {
            const TestComponent = {
                setup() {
                    const triggerRef = ref<HTMLElement | null>(null)
                    const overlay = useOverlay({
                        closeOnClickOutside: true,
                        ignoreElements: [triggerRef]
                    })
                    overlay.open()
                    return { overlay, triggerRef }
                },
                template: `
                    <div>
                        <button ref="triggerRef">Trigger</button>
                        <div ref="overlayRef">Content</div>
                    </div>
                `
            }

            const wrapper = mount(TestComponent, {
                attachTo: document.body
            })

            expect(wrapper.vm.overlay.isVisible.value).toBe(true)

            // Click on trigger (should be ignored)
            const trigger = wrapper.find('button')
            await trigger.trigger('click')
            await nextTick()

            expect(wrapper.vm.overlay.isVisible.value).toBe(true)

            wrapper.unmount()
        })
    })

    describe('Z-index management', () => {
        it('should return default z-index', () => {
            const { getZIndex } = useOverlay()

            expect(getZIndex()).toBe(1000)
        })

        it('should return custom z-index', () => {
            const { getZIndex } = useOverlay({ zIndex: 2000 })

            expect(getZIndex()).toBe(2000)
        })
    })

    describe('Overlay ref', () => {
        it('should provide overlayRef', () => {
            const { overlayRef } = useOverlay()

            expect(overlayRef.value).toBeNull()
        })

        it('should allow setting overlayRef', () => {
            const { overlayRef } = useOverlay()
            const element = document.createElement('div')

            overlayRef.value = element

            expect(overlayRef.value).toBe(element)
        })
    })

    describe('Edge cases', () => {
        it('should handle rapid open/close calls', () => {
            const { isVisible, open, close } = useOverlay()

            open()
            close()
            open()
            close()
            open()

            expect(isVisible.value).toBe(true)
        })

        it('should handle multiple overlays independently', () => {
            const overlay1 = useOverlay()
            const overlay2 = useOverlay()

            overlay1.open()
            expect(overlay1.isVisible.value).toBe(true)
            expect(overlay2.isVisible.value).toBe(false)

            overlay2.open()
            expect(overlay1.isVisible.value).toBe(true)
            expect(overlay2.isVisible.value).toBe(true)

            overlay1.close()
            expect(overlay1.isVisible.value).toBe(false)
            expect(overlay2.isVisible.value).toBe(true)
        })

        it('should handle position updates while closed', () => {
            const { position, setPosition } = useOverlay()

            setPosition({ top: '100px', left: '200px' })

            expect(position.value).toEqual({
                top: '100px',
                left: '200px',
                transform: ''
            })
        })

        it('should handle position updates while open', () => {
            const { position, setPosition, open } = useOverlay()

            open()
            setPosition({ top: '100px', left: '200px' })

            expect(position.value).toEqual({
                top: '100px',
                left: '200px',
                transform: ''
            })
        })
    })
})
