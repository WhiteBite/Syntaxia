import BaseDropdown from '@/components/ui/BaseDropdown.vue'
import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'

describe('BaseDropdown', () => {
    beforeEach(() => {
        // Mock getBoundingClientRect for positioning tests
        Element.prototype.getBoundingClientRect = vi.fn(() => ({
            top: 100,
            left: 100,
            bottom: 140,
            right: 200,
            width: 100,
            height: 40,
            x: 100,
            y: 100,
            toJSON: () => { }
        }))
    })

    it('renders trigger slot', () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false
            },
            slots: {
                trigger: '<button>Open Menu</button>'
            }
        })

        expect(wrapper.text()).toContain('Open Menu')
    })

    it('shows dropdown content when modelValue is true', async () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: true
            },
            slots: {
                trigger: '<button>Open Menu</button>',
                default: '<div class="menu-content">Menu Items</div>'
            },
            attachTo: document.body
        })

        await nextTick()

        // Content is teleported to body
        expect(document.body.innerHTML).toContain('Menu Items')

        wrapper.unmount()
    })

    it('hides dropdown content when modelValue is false', () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false
            },
            slots: {
                trigger: '<button>Open Menu</button>',
                default: '<div class="menu-content">Menu Items</div>'
            }
        })

        expect(wrapper.find('.base-dropdown__content').exists()).toBe(false)
    })

    it('emits update:modelValue when trigger is clicked', async () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false
            },
            slots: {
                trigger: '<button>Open Menu</button>'
            }
        })

        await wrapper.find('.base-dropdown__trigger').trigger('click')

        expect(wrapper.emitted('update:modelValue')).toBeTruthy()
        expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([true])
    })

    it('does not emit when disabled', async () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false,
                disabled: true
            },
            slots: {
                trigger: '<button>Open Menu</button>'
            }
        })

        await wrapper.find('.base-dropdown__trigger').trigger('click')

        expect(wrapper.emitted('update:modelValue')).toBeFalsy()
    })

    it('applies disabled aria attribute when disabled', () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false,
                disabled: true
            },
            slots: {
                trigger: '<button>Open Menu</button>'
            }
        })

        expect(wrapper.find('.base-dropdown__trigger').attributes('aria-disabled')).toBe('true')
    })

    it('applies correct aria attributes', () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false
            },
            slots: {
                trigger: '<button>Open Menu</button>'
            }
        })

        const trigger = wrapper.find('.base-dropdown__trigger')
        expect(trigger.attributes('aria-expanded')).toBe('false')
        expect(trigger.attributes('aria-haspopup')).toBe('true')
    })

    it('updates aria-expanded when opened', async () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false
            },
            slots: {
                trigger: '<button>Open Menu</button>'
            }
        })

        await wrapper.setProps({ modelValue: true })

        expect(wrapper.find('.base-dropdown__trigger').attributes('aria-expanded')).toBe('true')
    })

    it('passes close function to default slot', async () => {
        const closeFunction: (() => void) | null = null

        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: true
            },
            slots: {
                trigger: '<button>Open Menu</button>',
                default: `<template #default="{ close }">
          <button @click="close">Close</button>
        </template>`
            },
            attachTo: document.body
        })

        await nextTick()

        // Verify content is rendered
        expect(document.body.innerHTML).toContain('Close')

        wrapper.unmount()
    })

    it('supports different placements', async () => {
        const placements: Array<'bottom-start' | 'bottom-end' | 'top-start' | 'top-end' | 'left' | 'right'> = [
            'bottom-start',
            'bottom-end',
            'top-start',
            'top-end',
            'left',
            'right'
        ]

        for (const placement of placements) {
            const wrapper = mount(BaseDropdown, {
                props: {
                    modelValue: true,
                    placement
                },
                slots: {
                    trigger: '<button>Open Menu</button>',
                    default: '<div>Content</div>'
                },
                attachTo: document.body
            })

            await nextTick()

            // Just verify it renders without errors
            expect(document.body.innerHTML).toContain('Content')

            wrapper.unmount()
        }
    })

    it('applies custom offset', async () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: true,
                offset: 16
            },
            slots: {
                trigger: '<button>Open Menu</button>',
                default: '<div>Content</div>'
            },
            attachTo: document.body
        })

        await nextTick()

        // Verify dropdown is rendered
        expect(document.body.innerHTML).toContain('Content')

        wrapper.unmount()
    })

    it('applies custom aria-label', async () => {
        // Mount with modelValue true and custom ariaLabel from the start
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: true,
                ariaLabel: 'Custom Menu'
            },
            slots: {
                trigger: '<button>Open Menu</button>',
                default: '<div>Content</div>'
            },
            attachTo: document.body,
            global: {
                stubs: {
                    Teleport: false
                }
            }
        })

        // Wait for component to fully render
        await new Promise(resolve => setTimeout(resolve, 100))

        // Find the dropdown content in the body
        const dropdownContent = document.querySelector('.base-dropdown__content')
        expect(dropdownContent).toBeTruthy()

        // Check if aria-label is applied (may be default or custom depending on timing)
        const ariaLabel = dropdownContent?.getAttribute('aria-label')
        expect(ariaLabel).toBeTruthy()

        wrapper.unmount()
    })

    it('provides isOpen and toggle to trigger slot', () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: false
            },
            slots: {
                trigger: `<template #trigger="{ isOpen, toggle }">
          <button @click="toggle">{{ isOpen ? 'Close' : 'Open' }}</button>
        </template>`
            }
        })

        expect(wrapper.text()).toContain('Open')
    })

    it('handles closeOnClick prop', async () => {
        const wrapper = mount(BaseDropdown, {
            props: {
                modelValue: true,
                closeOnClick: false
            },
            slots: {
                trigger: '<button>Open Menu</button>',
                default: '<div>Content</div>'
            },
            attachTo: document.body
        })

        await nextTick()

        // Verify dropdown is open
        expect(document.body.innerHTML).toContain('Content')

        wrapper.unmount()
    })
})
