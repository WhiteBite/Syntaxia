import BaseAlert from '@/components/ui/BaseAlert.vue'
import { mount } from '@vue/test-utils'
import { AlertTriangle, CheckCircle, Info, XCircle } from 'lucide-vue-next'
import { describe, expect, it } from 'vitest'

describe('BaseAlert', () => {
    // Info variant tests
    describe('Info variant', () => {
        it('renders info alert with default icon', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: { default: 'Info message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.exists()).toBe(true)
            expect(alert.classes()).toContain('base-alert--info')
            expect(wrapper.text()).toContain('Info message')
            expect(wrapper.findComponent(Info).exists()).toBe(true)
        })

        it('renders info alert with title', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info', title: 'Information' },
                slots: { default: 'Info message' }
            })

            expect(wrapper.find('.base-alert__title').text()).toBe('Information')
            expect(wrapper.find('.base-alert__message').text()).toBe('Info message')
        })

        it('renders info alert with title slot', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: {
                    title: 'Custom Title',
                    default: 'Info message'
                }
            })

            expect(wrapper.find('.base-alert__title').text()).toBe('Custom Title')
        })
    })

    // Success variant tests
    describe('Success variant', () => {
        it('renders success alert with default icon', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'success' },
                slots: { default: 'Success message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.exists()).toBe(true)
            expect(alert.classes()).toContain('base-alert--success')
            expect(wrapper.text()).toContain('Success message')
            expect(wrapper.findComponent(CheckCircle).exists()).toBe(true)
        })

        it('renders success alert with title', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'success', title: 'Success!' },
                slots: { default: 'Operation completed' }
            })

            expect(wrapper.find('.base-alert__title').text()).toBe('Success!')
            expect(wrapper.find('.base-alert__message').text()).toBe('Operation completed')
        })

        it('applies correct CSS classes for success variant', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'success' },
                slots: { default: 'Success' }
            })

            expect(wrapper.find('.base-alert--success').exists()).toBe(true)
        })
    })

    // Warning variant tests
    describe('Warning variant', () => {
        it('renders warning alert with default icon', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'warning' },
                slots: { default: 'Warning message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.exists()).toBe(true)
            expect(alert.classes()).toContain('base-alert--warning')
            expect(wrapper.text()).toContain('Warning message')
            expect(wrapper.findComponent(AlertTriangle).exists()).toBe(true)
        })

        it('renders warning alert with title', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'warning', title: 'Warning!' },
                slots: { default: 'Be careful' }
            })

            expect(wrapper.find('.base-alert__title').text()).toBe('Warning!')
            expect(wrapper.find('.base-alert__message').text()).toBe('Be careful')
        })

        it('applies correct CSS classes for warning variant', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'warning' },
                slots: { default: 'Warning' }
            })

            expect(wrapper.find('.base-alert--warning').exists()).toBe(true)
        })
    })

    // Error variant tests
    describe('Error variant', () => {
        it('renders error alert with default icon', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'error' },
                slots: { default: 'Error message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.exists()).toBe(true)
            expect(alert.classes()).toContain('base-alert--error')
            expect(wrapper.text()).toContain('Error message')
            expect(wrapper.findComponent(XCircle).exists()).toBe(true)
        })

        it('renders error alert with title', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'error', title: 'Error!' },
                slots: { default: 'Something went wrong' }
            })

            expect(wrapper.find('.base-alert__title').text()).toBe('Error!')
            expect(wrapper.find('.base-alert__message').text()).toBe('Something went wrong')
        })

        it('applies correct CSS classes for error variant', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'error' },
                slots: { default: 'Error' }
            })

            expect(wrapper.find('.base-alert--error').exists()).toBe(true)
        })
    })

    // Dismissible functionality tests
    describe('Dismissible functionality', () => {
        it('does not show dismiss button by default', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: { default: 'Message' }
            })

            expect(wrapper.find('.base-alert__dismiss').exists()).toBe(false)
        })

        it('shows dismiss button when dismissible is true', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info', dismissible: true },
                slots: { default: 'Message' }
            })

            expect(wrapper.find('.base-alert__dismiss').exists()).toBe(true)
        })

        it('emits dismiss event when dismiss button is clicked', async () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info', dismissible: true },
                slots: { default: 'Message' }
            })

            await wrapper.find('.base-alert__dismiss').trigger('click')
            expect(wrapper.emitted('dismiss')).toBeTruthy()
            expect(wrapper.emitted('dismiss')).toHaveLength(1)
        })

        it('hides alert after dismiss button is clicked', async () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info', dismissible: true },
                slots: { default: 'Message' }
            })

            expect(wrapper.find('.base-alert').exists()).toBe(true)
            await wrapper.find('.base-alert__dismiss').trigger('click')

            // Wait for transition
            await wrapper.vm.$nextTick()
            expect(wrapper.find('.base-alert').exists()).toBe(false)
        })

        it('uses custom dismiss label', () => {
            const wrapper = mount(BaseAlert, {
                props: {
                    variant: 'info',
                    dismissible: true,
                    dismissLabel: 'Close Alert'
                },
                slots: { default: 'Message' }
            })

            const dismissBtn = wrapper.find('.base-alert__dismiss')
            expect(dismissBtn.attributes('aria-label')).toBe('Close Alert')
        })
    })

    // Custom icon tests
    describe('Custom icon', () => {
        it('renders custom icon when icon slot is provided', () => {
            const CustomIcon = {
                name: 'CustomIcon',
                template: '<svg class="custom-icon"></svg>'
            }

            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: {
                    icon: CustomIcon,
                    default: 'Message'
                }
            })

            expect(wrapper.find('.custom-icon').exists()).toBe(true)
            expect(wrapper.findComponent(Info).exists()).toBe(false)
        })
    })

    // Actions slot tests
    describe('Actions slot', () => {
        it('does not render actions container when no actions slot provided', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: { default: 'Message' }
            })

            expect(wrapper.find('.base-alert__actions').exists()).toBe(false)
        })

        it('renders actions slot when provided', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: {
                    default: 'Message',
                    actions: '<button class="test-action">Action</button>'
                }
            })

            expect(wrapper.find('.base-alert__actions').exists()).toBe(true)
            expect(wrapper.find('.test-action').exists()).toBe(true)
        })
    })

    // Accessibility tests
    describe('Accessibility', () => {
        it('has role="alert"', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: { default: 'Message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.attributes('role')).toBe('alert')
        })

        it('has aria-live="polite" for non-error variants', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info' },
                slots: { default: 'Message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.attributes('aria-live')).toBe('polite')
        })

        it('has aria-live="assertive" for error variant', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'error' },
                slots: { default: 'Error message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.attributes('aria-live')).toBe('assertive')
        })

        it('dismiss button has aria-label', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info', dismissible: true },
                slots: { default: 'Message' }
            })

            const dismissBtn = wrapper.find('.base-alert__dismiss')
            expect(dismissBtn.attributes('aria-label')).toBe('Dismiss')
        })

        it('dismiss button has type="button"', () => {
            const wrapper = mount(BaseAlert, {
                props: { variant: 'info', dismissible: true },
                slots: { default: 'Message' }
            })

            const dismissBtn = wrapper.find('.base-alert__dismiss')
            expect(dismissBtn.attributes('type')).toBe('button')
        })
    })

    // Default props tests
    describe('Default props', () => {
        it('uses info variant by default', () => {
            const wrapper = mount(BaseAlert, {
                slots: { default: 'Message' }
            })

            const alert = wrapper.find('.base-alert')
            expect(alert.classes()).toContain('base-alert--info')
        })

        it('is not dismissible by default', () => {
            const wrapper = mount(BaseAlert, {
                slots: { default: 'Message' }
            })

            expect(wrapper.find('.base-alert__dismiss').exists()).toBe(false)
        })
    })
})
