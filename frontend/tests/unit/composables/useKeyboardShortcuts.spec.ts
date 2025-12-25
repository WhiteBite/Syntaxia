/**
 * Keyboard Shortcuts Tests
 * Tests for useKeyboardShortcuts composable
 */

import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts'
import { mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { defineComponent, ref } from 'vue'

// Helper to create keyboard events
function createKeyboardEvent(
    key: string,
    options: { ctrl?: boolean; shift?: boolean; alt?: boolean } = {}
): KeyboardEvent {
    return new KeyboardEvent('keydown', {
        key,
        ctrlKey: options.ctrl || false,
        shiftKey: options.shift || false,
        altKey: options.alt || false,
        bubbles: true,
    })
}

describe('useKeyboardShortcuts', () => {
    let cleanup: (() => void) | null = null

    afterEach(() => {
        cleanup?.()
        cleanup = null
    })

    it('should trigger action on matching shortcut', async () => {
        const action = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'b', ctrl: true, action, description: 'Test action' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))

        expect(action).toHaveBeenCalledTimes(1)

        wrapper.unmount()
    })

    it('should trigger buildContext on Ctrl+B', async () => {
        const buildContext = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'b', ctrl: true, action: buildContext, description: 'Build context' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))

        expect(buildContext).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should trigger copyContext on Ctrl+Shift+C', async () => {
        const copyContext = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'c', ctrl: true, shift: true, action: copyContext, description: 'Copy context' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        document.dispatchEvent(createKeyboardEvent('c', { ctrl: true, shift: true }))

        expect(copyContext).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should close modal on Escape', async () => {
        const closeModal = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'Escape', action: closeModal, description: 'Close modal' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        document.dispatchEvent(createKeyboardEvent('Escape'))

        expect(closeModal).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should switch tabs on Ctrl+1/2/3', async () => {
        const switchTab1 = vi.fn()
        const switchTab2 = vi.fn()
        const switchTab3 = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: '1', ctrl: true, action: switchTab1, description: 'Tab 1' },
                    { key: '2', ctrl: true, action: switchTab2, description: 'Tab 2' },
                    { key: '3', ctrl: true, action: switchTab3, description: 'Tab 3' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        document.dispatchEvent(createKeyboardEvent('1', { ctrl: true }))
        expect(switchTab1).toHaveBeenCalled()

        document.dispatchEvent(createKeyboardEvent('2', { ctrl: true }))
        expect(switchTab2).toHaveBeenCalled()

        document.dispatchEvent(createKeyboardEvent('3', { ctrl: true }))
        expect(switchTab3).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should NOT trigger when input is focused (except Escape)', async () => {
        const action = vi.fn()
        const escapeAction = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'b', ctrl: true, action, description: 'Build' },
                    { key: 'Escape', action: escapeAction, description: 'Close' },
                ])
                return () => null
            },
            template: '<input id="test-input" />',
        })

        const wrapper = mount(TestComponent)

        // Create input element and focus it
        const input = document.createElement('input')
        document.body.appendChild(input)
        input.focus()

        // Create event with input as target
        const event = new KeyboardEvent('keydown', {
            key: 'b',
            ctrlKey: true,
            bubbles: true,
        })
        Object.defineProperty(event, 'target', { value: input })

        document.dispatchEvent(event)

        // Should NOT trigger because input is focused
        expect(action).not.toHaveBeenCalled()

        // Escape should still work
        const escapeEvent = new KeyboardEvent('keydown', {
            key: 'Escape',
            bubbles: true,
        })
        Object.defineProperty(escapeEvent, 'target', { value: input })

        document.dispatchEvent(escapeEvent)
        expect(escapeAction).toHaveBeenCalled()

        document.body.removeChild(input)
        wrapper.unmount()
    })

    it('should cleanup listeners on unmount', async () => {
        const action = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'b', ctrl: true, action, description: 'Test' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        // Trigger once
        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))
        expect(action).toHaveBeenCalledTimes(1)

        // Unmount
        wrapper.unmount()

        // Trigger again - should not call action
        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))
        expect(action).toHaveBeenCalledTimes(1) // Still 1, not 2
    })

    it('should respect enabled condition', async () => {
        const action = vi.fn()
        const isEnabled = ref(false)

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'b', ctrl: true, action, description: 'Test', enabled: isEnabled },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        // Should not trigger when disabled
        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))
        expect(action).not.toHaveBeenCalled()

        // Enable and try again
        isEnabled.value = true
        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))
        expect(action).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should handle case-insensitive key matching', async () => {
        const action = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'B', ctrl: true, action, description: 'Test' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        // Lowercase 'b' should match uppercase 'B' definition
        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))
        expect(action).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should prevent default on matched shortcuts', async () => {
        const action = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'b', ctrl: true, action, description: 'Test' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        const event = createKeyboardEvent('b', { ctrl: true })
        const preventDefaultSpy = vi.spyOn(event, 'preventDefault')

        document.dispatchEvent(event)

        expect(preventDefaultSpy).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should not trigger without required modifier', async () => {
        const action = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'b', ctrl: true, action, description: 'Test' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        // Press 'b' without Ctrl
        document.dispatchEvent(createKeyboardEvent('b'))
        expect(action).not.toHaveBeenCalled()

        // Press 'b' with Ctrl
        document.dispatchEvent(createKeyboardEvent('b', { ctrl: true }))
        expect(action).toHaveBeenCalled()

        wrapper.unmount()
    })

    it('should handle Alt modifier', async () => {
        const action = vi.fn()

        const TestComponent = defineComponent({
            setup() {
                useKeyboardShortcuts([
                    { key: 'p', alt: true, action, description: 'Test' },
                ])
                return () => null
            },
        })

        const wrapper = mount(TestComponent)

        document.dispatchEvent(createKeyboardEvent('p', { alt: true }))
        expect(action).toHaveBeenCalled()

        wrapper.unmount()
    })
})
