import BaseChip from '@/components/ui/BaseChip.vue'
import { mount } from '@vue/test-utils'
import { FileCode, X } from 'lucide-vue-next'
import { describe, expect, it } from 'vitest'

describe('BaseChip', () => {
    describe('Rendering', () => {
        it('renders with default props', () => {
            const wrapper = mount(BaseChip, {
                slots: {
                    default: 'Test Chip'
                }
            })

            expect(wrapper.text()).toContain('Test Chip')
            expect(wrapper.classes()).toContain('base-chip')
            expect(wrapper.classes()).toContain('base-chip--default')
            expect(wrapper.classes()).toContain('base-chip--sm')
        })

        it('renders as span by default', () => {
            const wrapper = mount(BaseChip, {
                slots: { default: 'Test' }
            })

            expect(wrapper.element.tagName).toBe('SPAN')
        })

        it('renders as button when clickable', () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            expect(wrapper.element.tagName).toBe('BUTTON')
            expect(wrapper.attributes('type')).toBe('button')
        })

        it('renders slot content correctly', () => {
            const wrapper = mount(BaseChip, {
                slots: {
                    default: '<strong>Bold Text</strong>'
                }
            })

            expect(wrapper.html()).toContain('<strong>Bold Text</strong>')
        })
    })

    describe('Variants', () => {
        it.each([
            ['default', 'base-chip--default'],
            ['primary', 'base-chip--primary'],
            ['success', 'base-chip--success'],
            ['warning', 'base-chip--warning'],
            ['danger', 'base-chip--danger']
        ])('applies %s variant class', (variant, expectedClass) => {
            const wrapper = mount(BaseChip, {
                props: { variant: variant as any },
                slots: { default: 'Test' }
            })

            expect(wrapper.classes()).toContain(expectedClass)
        })
    })

    describe('Sizes', () => {
        it.each([
            ['xs', 'base-chip--xs'],
            ['sm', 'base-chip--sm'],
            ['md', 'base-chip--md']
        ])('applies %s size class', (size, expectedClass) => {
            const wrapper = mount(BaseChip, {
                props: { size: size as any },
                slots: { default: 'Test' }
            })

            expect(wrapper.classes()).toContain(expectedClass)
        })
    })

    describe('Icon Support', () => {
        it('renders icon from prop', () => {
            const wrapper = mount(BaseChip, {
                props: { icon: FileCode },
                slots: { default: 'Test' }
            })

            expect(wrapper.find('.base-chip__icon').exists()).toBe(true)
            expect(wrapper.findComponent(FileCode).exists()).toBe(true)
        })

        it('renders icon from slot', () => {
            const wrapper = mount(BaseChip, {
                slots: {
                    icon: '<span class="custom-icon">🔥</span>',
                    default: 'Test'
                }
            })

            expect(wrapper.find('.base-chip__icon').exists()).toBe(true)
            expect(wrapper.find('.custom-icon').text()).toBe('🔥')
        })

        it('prioritizes icon slot over icon prop', () => {
            const wrapper = mount(BaseChip, {
                props: { icon: FileCode },
                slots: {
                    icon: '<span class="slot-icon">✓</span>',
                    default: 'Test'
                }
            })

            expect(wrapper.find('.slot-icon').exists()).toBe(true)
            expect(wrapper.findComponent(FileCode).exists()).toBe(false)
        })

        it('does not render icon container when no icon provided', () => {
            const wrapper = mount(BaseChip, {
                slots: { default: 'Test' }
            })

            expect(wrapper.find('.base-chip__icon').exists()).toBe(false)
        })
    })

    describe('Clickable Behavior', () => {
        it('applies clickable class when clickable prop is true', () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            expect(wrapper.classes()).toContain('base-chip--clickable')
        })

        it('emits click event when clicked', async () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            await wrapper.trigger('click')

            expect(wrapper.emitted('click')).toBeTruthy()
            expect(wrapper.emitted('click')).toHaveLength(1)
        })

        it('does not emit click event when not clickable', async () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: false },
                slots: { default: 'Test' }
            })

            await wrapper.trigger('click')

            expect(wrapper.emitted('click')).toBeFalsy()
        })

        it('emits click event on Enter key', async () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            await wrapper.trigger('keydown', { key: 'Enter' })

            expect(wrapper.emitted('click')).toBeTruthy()
            expect(wrapper.emitted('click')).toHaveLength(1)
        })

        it('emits click event on Space key', async () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            await wrapper.trigger('keydown', { key: ' ' })

            expect(wrapper.emitted('click')).toBeTruthy()
            expect(wrapper.emitted('click')).toHaveLength(1)
        })

        it('has correct accessibility attributes when clickable', () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            expect(wrapper.attributes('role')).toBe('button')
            expect(wrapper.attributes('tabindex')).toBe('0')
        })
    })

    describe('Removable Behavior', () => {
        it('applies removable class when removable prop is true', () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            expect(wrapper.classes()).toContain('base-chip--removable')
        })

        it('renders remove button when removable', () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            expect(wrapper.find('.base-chip__remove').exists()).toBe(true)
            expect(wrapper.findComponent(X).exists()).toBe(true)
        })

        it('does not render remove button when not removable', () => {
            const wrapper = mount(BaseChip, {
                props: { removable: false },
                slots: { default: 'Test' }
            })

            expect(wrapper.find('.base-chip__remove').exists()).toBe(false)
        })

        it('emits remove event when remove button clicked', async () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            await wrapper.find('.base-chip__remove').trigger('click')

            expect(wrapper.emitted('remove')).toBeTruthy()
            expect(wrapper.emitted('remove')).toHaveLength(1)
        })

        it('stops propagation when remove button clicked', async () => {
            const wrapper = mount(BaseChip, {
                props: {
                    removable: true,
                    clickable: true
                },
                slots: { default: 'Test' }
            })

            await wrapper.find('.base-chip__remove').trigger('click')

            expect(wrapper.emitted('remove')).toBeTruthy()
            expect(wrapper.emitted('click')).toBeFalsy()
        })

        it('emits remove event on Delete key', async () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            await wrapper.trigger('keydown', { key: 'Delete' })

            expect(wrapper.emitted('remove')).toBeTruthy()
            expect(wrapper.emitted('remove')).toHaveLength(1)
        })

        it('emits remove event on Enter key when remove button focused', async () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            await wrapper.find('.base-chip__remove').trigger('keydown', { key: 'Enter' })

            expect(wrapper.emitted('remove')).toBeTruthy()
        })

        it('emits remove event on Space key when remove button focused', async () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            await wrapper.find('.base-chip__remove').trigger('keydown', { key: ' ' })

            expect(wrapper.emitted('remove')).toBeTruthy()
        })

        it('has aria-label on remove button', () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            const removeButton = wrapper.find('.base-chip__remove')
            expect(removeButton.attributes('aria-label')).toBeTruthy()
        })
    })

    describe('Combined Features', () => {
        it('supports both clickable and removable', () => {
            const wrapper = mount(BaseChip, {
                props: {
                    clickable: true,
                    removable: true
                },
                slots: { default: 'Test' }
            })

            expect(wrapper.classes()).toContain('base-chip--clickable')
            expect(wrapper.classes()).toContain('base-chip--removable')
            expect(wrapper.find('.base-chip__remove').exists()).toBe(true)
        })

        it('supports icon + clickable + removable', async () => {
            const wrapper = mount(BaseChip, {
                props: {
                    icon: FileCode,
                    clickable: true,
                    removable: true
                },
                slots: { default: 'Test' }
            })

            expect(wrapper.find('.base-chip__icon').exists()).toBe(true)
            expect(wrapper.find('.base-chip__remove').exists()).toBe(true)

            await wrapper.trigger('click')
            expect(wrapper.emitted('click')).toBeTruthy()

            await wrapper.find('.base-chip__remove').trigger('click')
            expect(wrapper.emitted('remove')).toBeTruthy()
        })
    })

    describe('Accessibility', () => {
        it('has proper role when clickable', () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            expect(wrapper.attributes('role')).toBe('button')
        })

        it('has tabindex when clickable', () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            expect(wrapper.attributes('tabindex')).toBe('0')
        })

        it('does not have role when not clickable', () => {
            const wrapper = mount(BaseChip, {
                slots: { default: 'Test' }
            })

            expect(wrapper.attributes('role')).toBeUndefined()
        })

        it('does not have tabindex when not clickable', () => {
            const wrapper = mount(BaseChip, {
                slots: { default: 'Test' }
            })

            expect(wrapper.attributes('tabindex')).toBeUndefined()
        })
    })

    describe('Edge Cases', () => {
        it('handles empty content', () => {
            const wrapper = mount(BaseChip)

            expect(wrapper.find('.base-chip__text').exists()).toBe(true)
        })

        it('handles long text content', () => {
            const longText = 'This is a very long chip text that should be handled properly'
            const wrapper = mount(BaseChip, {
                slots: { default: longText }
            })

            expect(wrapper.text()).toContain(longText)
        })

        it('handles multiple rapid clicks', async () => {
            const wrapper = mount(BaseChip, {
                props: { clickable: true },
                slots: { default: 'Test' }
            })

            await wrapper.trigger('click')
            await wrapper.trigger('click')
            await wrapper.trigger('click')

            expect(wrapper.emitted('click')).toHaveLength(3)
        })

        it('handles multiple rapid removes', async () => {
            const wrapper = mount(BaseChip, {
                props: { removable: true },
                slots: { default: 'Test' }
            })

            const removeButton = wrapper.find('.base-chip__remove')
            await removeButton.trigger('click')
            await removeButton.trigger('click')

            expect(wrapper.emitted('remove')).toHaveLength(2)
        })
    })
})
