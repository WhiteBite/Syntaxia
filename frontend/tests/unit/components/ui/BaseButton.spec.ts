import BaseButton from '@/components/ui/BaseButton.vue'
import BaseIcon from '@/components/ui/BaseIcon.vue'
import BaseSpinner from '@/components/ui/BaseSpinner.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { h } from 'vue'

// Mock icon component
const MockIcon = {
    name: 'MockIcon',
    template: '<svg><path /></svg>'
}

describe('BaseButton.vue', () => {
    describe('Rendering', () => {
        it('should render button with default props', () => {
            const wrapper = mount(BaseButton, {
                slots: {
                    default: 'Click me'
                }
            })

            expect(wrapper.find('button').exists()).toBe(true)
            expect(wrapper.text()).toBe('Click me')
            expect(wrapper.classes()).toContain('base-button')
            expect(wrapper.classes()).toContain('base-button--secondary')
            expect(wrapper.classes()).toContain('base-button--md')
        })

        it('should render with custom text', () => {
            const wrapper = mount(BaseButton, {
                slots: {
                    default: 'Custom Text'
                }
            })

            expect(wrapper.text()).toBe('Custom Text')
        })

        it('should apply correct type attribute', () => {
            const wrapper = mount(BaseButton, {
                props: { type: 'submit' }
            })

            expect(wrapper.find('button').attributes('type')).toBe('submit')
        })
    })

    describe('Variants', () => {
        const variants = ['primary', 'secondary', 'ghost', 'danger', 'success', 'warning', 'link', 'text'] as const

        variants.forEach(variant => {
            it(`should apply ${variant} variant class`, () => {
                const wrapper = mount(BaseButton, {
                    props: { variant },
                    slots: { default: 'Button' }
                })

                expect(wrapper.classes()).toContain(`base-button--${variant}`)
            })
        })
    })

    describe('Sizes', () => {
        const sizes = ['xs', 'sm', 'md', 'lg'] as const

        sizes.forEach(size => {
            it(`should apply ${size} size class`, () => {
                const wrapper = mount(BaseButton, {
                    props: { size },
                    slots: { default: 'Button' }
                })

                expect(wrapper.classes()).toContain(`base-button--${size}`)
            })
        })
    })

    describe('Disabled state', () => {
        it('should be disabled when disabled prop is true', () => {
            const wrapper = mount(BaseButton, {
                props: { disabled: true },
                slots: { default: 'Button' }
            })

            expect(wrapper.find('button').attributes('disabled')).toBeDefined()
        })

        it('should not emit click event when disabled', async () => {
            const wrapper = mount(BaseButton, {
                props: { disabled: true },
                slots: { default: 'Button' }
            })

            await wrapper.find('button').trigger('click')

            expect(wrapper.emitted('click')).toBeFalsy()
        })
    })

    describe('Loading state', () => {
        it('should show spinner when loading is true', () => {
            const wrapper = mount(BaseButton, {
                props: { loading: true },
                slots: { default: 'Button' }
            })

            expect(wrapper.findComponent(BaseSpinner).exists()).toBe(true)
            expect(wrapper.classes()).toContain('base-button--loading')
        })

        it('should be disabled when loading', () => {
            const wrapper = mount(BaseButton, {
                props: { loading: true },
                slots: { default: 'Button' }
            })

            expect(wrapper.find('button').attributes('disabled')).toBeDefined()
        })

        it('should hide icon when loading', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    loading: true,
                    icon: MockIcon
                },
                slots: { default: 'Button' }
            })

            expect(wrapper.findComponent(BaseSpinner).exists()).toBe(true)
            expect(wrapper.findComponent(BaseIcon).exists()).toBe(false)
        })

        it('should preserve button text when loading', () => {
            const wrapper = mount(BaseButton, {
                props: { loading: true },
                slots: { default: 'Loading...' }
            })

            expect(wrapper.text()).toBe('Loading...')
        })
    })

    describe('Icon support', () => {
        it('should render icon when icon prop is provided', () => {
            const wrapper = mount(BaseButton, {
                props: { icon: MockIcon },
                slots: { default: 'Button' }
            })

            expect(wrapper.findComponent(BaseIcon).exists()).toBe(true)
        })

        it('should render icon on the left by default', () => {
            const wrapper = mount(BaseButton, {
                props: { icon: MockIcon },
                slots: { default: 'Button' }
            })

            const iconWrapper = wrapper.find('.base-button__icon--left')
            expect(iconWrapper.exists()).toBe(true)
        })

        it('should render icon on the right when iconPosition is right', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    icon: MockIcon,
                    iconPosition: 'right'
                },
                slots: { default: 'Button' }
            })

            const iconWrapper = wrapper.find('.base-button__icon--right')
            expect(iconWrapper.exists()).toBe(true)
        })

        it('should render icon from slot', () => {
            const wrapper = mount(BaseButton, {
                slots: {
                    default: 'Button',
                    icon: h(MockIcon)
                }
            })

            expect(wrapper.find('.base-button__icon').exists()).toBe(true)
        })

        it('should hide text when iconOnly is true', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    icon: MockIcon,
                    iconOnly: true
                },
                slots: { default: 'Button' }
            })

            expect(wrapper.find('.base-button__text').exists()).toBe(false)
            expect(wrapper.classes()).toContain('base-button--icon-only')
        })
    })

    describe('Block prop', () => {
        it('should apply block class when block is true', () => {
            const wrapper = mount(BaseButton, {
                props: { block: true },
                slots: { default: 'Button' }
            })

            expect(wrapper.classes()).toContain('base-button--block')
        })
    })

    describe('Outline prop', () => {
        it('should apply outline class when outline is true', () => {
            const wrapper = mount(BaseButton, {
                props: { outline: true },
                slots: { default: 'Button' }
            })

            expect(wrapper.classes()).toContain('base-button--outline')
        })

        it('should work with different variants', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    variant: 'primary',
                    outline: true
                },
                slots: { default: 'Button' }
            })

            expect(wrapper.classes()).toContain('base-button--outline')
            expect(wrapper.classes()).toContain('base-button--primary')
        })
    })

    describe('Click events', () => {
        it('should emit click event when clicked', async () => {
            const wrapper = mount(BaseButton, {
                slots: { default: 'Button' }
            })

            await wrapper.find('button').trigger('click')

            expect(wrapper.emitted('click')).toBeTruthy()
            expect(wrapper.emitted('click')).toHaveLength(1)
        })

        it('should pass event object to click handler', async () => {
            const wrapper = mount(BaseButton, {
                slots: { default: 'Button' }
            })

            await wrapper.find('button').trigger('click')

            const clickEvent = wrapper.emitted('click')?.[0]?.[0]
            expect(clickEvent).toBeInstanceOf(Event)
        })
    })

    describe('Icon size mapping', () => {
        it('should use xs icon for xs button', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    size: 'xs',
                    icon: MockIcon
                },
                slots: { default: 'Button' }
            })

            const icon = wrapper.findComponent(BaseIcon)
            expect(icon.props('size')).toBe('xs')
        })

        it('should use sm icon for sm button', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    size: 'sm',
                    icon: MockIcon
                },
                slots: { default: 'Button' }
            })

            const icon = wrapper.findComponent(BaseIcon)
            expect(icon.props('size')).toBe('sm')
        })

        it('should use sm icon for md button', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    size: 'md',
                    icon: MockIcon
                },
                slots: { default: 'Button' }
            })

            const icon = wrapper.findComponent(BaseIcon)
            expect(icon.props('size')).toBe('sm')
        })

        it('should use md icon for lg button', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    size: 'lg',
                    icon: MockIcon
                },
                slots: { default: 'Button' }
            })

            const icon = wrapper.findComponent(BaseIcon)
            expect(icon.props('size')).toBe('md')
        })
    })

    describe('Spinner size mapping', () => {
        it('should match spinner size to icon size', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    size: 'lg',
                    loading: true
                },
                slots: { default: 'Button' }
            })

            const spinner = wrapper.findComponent(BaseSpinner)
            expect(spinner.props('size')).toBe('md')
        })
    })

    describe('Complex scenarios', () => {
        it('should handle loading with icon position right', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    loading: true,
                    icon: MockIcon,
                    iconPosition: 'right'
                },
                slots: { default: 'Button' }
            })

            // When loading, spinner appears but icon on right is hidden
            expect(wrapper.find('.base-button__icon--left').exists()).toBe(true)
            expect(wrapper.findComponent(BaseSpinner).exists()).toBe(true)
            // Right icon should not exist when loading
            expect(wrapper.findComponent(BaseIcon).exists()).toBe(false)
        })

        it('should handle all props together', () => {
            const wrapper = mount(BaseButton, {
                props: {
                    variant: 'primary',
                    size: 'lg',
                    icon: MockIcon,
                    iconPosition: 'right',
                    block: true,
                    outline: true,
                    disabled: false,
                    loading: false
                },
                slots: { default: 'Complex Button' }
            })

            expect(wrapper.classes()).toContain('base-button--primary')
            expect(wrapper.classes()).toContain('base-button--lg')
            expect(wrapper.classes()).toContain('base-button--block')
            expect(wrapper.classes()).toContain('base-button--outline')
            expect(wrapper.find('.base-button__icon--right').exists()).toBe(true)
            expect(wrapper.text()).toBe('Complex Button')
        })

        it('should handle link variant without icon', () => {
            const wrapper = mount(BaseButton, {
                props: { variant: 'link' },
                slots: { default: 'Link Button' }
            })

            expect(wrapper.classes()).toContain('base-button--link')
            expect(wrapper.text()).toBe('Link Button')
        })

        it('should handle text variant', () => {
            const wrapper = mount(BaseButton, {
                props: { variant: 'text' },
                slots: { default: 'Text Button' }
            })

            expect(wrapper.classes()).toContain('base-button--text')
            expect(wrapper.text()).toBe('Text Button')
        })
    })
})
