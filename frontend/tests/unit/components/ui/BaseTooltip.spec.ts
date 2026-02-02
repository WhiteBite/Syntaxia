import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

describe('BaseTooltip', () => {
    beforeEach(() => {
        // Mock getBoundingClientRect
        Element.prototype.getBoundingClientRect = vi.fn(() => ({
            top: 100,
            left: 100,
            bottom: 120,
            right: 200,
            width: 100,
            height: 20,
            x: 100,
            y: 100,
            toJSON: () => { }
        }))
    })

    it('renders slot content', () => {
        const wrapper = mount(BaseTooltip, {
            props: { text: 'Tooltip text' },
            slots: {
                default: '<button>Hover me</button>'
            }
        })

        expect(wrapper.find('button').text()).toBe('Hover me')
    })

    it('shows tooltip on mouseenter', async () => {
        const wrapper = mount(BaseTooltip, {
            props: { text: 'Tooltip text' },
            slots: {
                default: '<button>Hover me</button>'
            },
            attachTo: document.body
        })

        await wrapper.trigger('mouseenter')
        await flushPromises()

        const tooltip = document.querySelector('.tooltip')
        expect(tooltip).toBeTruthy()
        expect(tooltip?.textContent).toBe('Tooltip text')

        wrapper.unmount()
    })

    it('hides tooltip on mouseleave', async () => {
        const wrapper = mount(BaseTooltip, {
            props: { text: 'Tooltip text' },
            slots: {
                default: '<button>Hover me</button>'
            },
            attachTo: document.body
        })

        await wrapper.trigger('mouseenter')
        await flushPromises()
        expect(document.querySelector('.tooltip')).toBeTruthy()

        await wrapper.trigger('mouseleave')
        await flushPromises()
        expect(document.querySelector('.tooltip')).toBeFalsy()

        wrapper.unmount()
    })

    it('applies correct position class', async () => {
        const positions = ['top', 'bottom', 'left', 'right'] as const

        for (const position of positions) {
            const wrapper = mount(BaseTooltip, {
                props: { text: 'Tooltip text', position },
                slots: {
                    default: '<button>Hover me</button>'
                },
                attachTo: document.body
            })

            await wrapper.trigger('mouseenter')
            await flushPromises()

            const tooltip = document.querySelector(`.tooltip--${position}`)
            expect(tooltip).toBeTruthy()

            wrapper.unmount()
        }
    })

    it('applies multiline class when multiline prop is true', async () => {
        const wrapper = mount(BaseTooltip, {
            props: { text: 'Long tooltip text', multiline: true },
            slots: {
                default: '<button>Hover me</button>'
            },
            attachTo: document.body
        })

        await wrapper.trigger('mouseenter')
        await flushPromises()

        const tooltip = document.querySelector('.tooltip--multiline')
        expect(tooltip).toBeTruthy()

        wrapper.unmount()
    })

    it('uses top position by default', async () => {
        const wrapper = mount(BaseTooltip, {
            props: { text: 'Tooltip text' },
            slots: {
                default: '<button>Hover me</button>'
            },
            attachTo: document.body
        })

        await wrapper.trigger('mouseenter')
        await flushPromises()

        const tooltip = document.querySelector('.tooltip--top')
        expect(tooltip).toBeTruthy()

        wrapper.unmount()
    })

    it('calculates tooltip position based on trigger element', async () => {
        const wrapper = mount(BaseTooltip, {
            props: { text: 'Tooltip text', position: 'bottom' },
            slots: {
                default: '<button>Hover me</button>'
            },
            attachTo: document.body
        })

        await wrapper.trigger('mouseenter')
        await flushPromises()

        const tooltip = document.querySelector('.tooltip') as HTMLElement
        expect(tooltip).toBeTruthy()

        const style = tooltip?.getAttribute('style') || ''
        expect(style).toContain('position: fixed')
        expect(style).toContain('z-index: 10000')

        wrapper.unmount()
    })

    it('uses Teleport to render tooltip in body', async () => {
        const wrapper = mount(BaseTooltip, {
            props: { text: 'Tooltip text' },
            slots: {
                default: '<button>Hover me</button>'
            },
            attachTo: document.body
        })

        await wrapper.trigger('mouseenter')
        await flushPromises()

        // Tooltip should be rendered directly in body, not inside wrapper
        const tooltipInBody = document.querySelector('.tooltip')
        expect(tooltipInBody).toBeTruthy()
        expect(tooltipInBody?.textContent).toBe('Tooltip text')

        wrapper.unmount()
    })
})
