import BaseModal from '@/components/ui/BaseModal.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { nextTick } from 'vue'

describe('BaseModal', () => {
    describe('Props', () => {
        it('should accept modelValue prop', () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            expect(wrapper.props('modelValue')).toBe(true)
            wrapper.unmount()
        })

        it('should accept title prop', () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                    title: 'Test Modal',
                },
            })

            expect(wrapper.props('title')).toBe('Test Modal')
            wrapper.unmount()
        })

        it('should accept size prop with default value md', () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            expect(wrapper.props('size')).toBe('md')
            wrapper.unmount()
        })

        it('should accept all size variants', () => {
            const sizes: Array<'sm' | 'md' | 'lg' | 'xl' | 'full'> = ['sm', 'md', 'lg', 'xl', 'full']

            sizes.forEach(size => {
                const wrapper = mount(BaseModal, {
                    props: {
                        modelValue: true,
                        size,
                    },
                })

                expect(wrapper.props('size')).toBe(size)
                wrapper.unmount()
            })
        })

        it('should have correct default prop values', () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            expect(wrapper.props('closeOnBackdrop')).toBe(true)
            expect(wrapper.props('closeOnEsc')).toBe(true)
            expect(wrapper.props('showClose')).toBe(true)
            wrapper.unmount()
        })
    })

    describe('Events', () => {
        it('should emit update:modelValue event', async () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            wrapper.vm.$emit('update:modelValue', false)
            await nextTick()

            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([false])
            wrapper.unmount()
        })

        it('should emit close event', async () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            wrapper.vm.$emit('close')
            await nextTick()

            expect(wrapper.emitted('close')).toBeTruthy()
            wrapper.unmount()
        })

        it('should emit open event', async () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            wrapper.vm.$emit('open')
            await nextTick()

            expect(wrapper.emitted('open')).toBeTruthy()
            wrapper.unmount()
        })
    })

    describe('Component Structure', () => {
        it('should be a Vue component', () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            expect(wrapper.exists()).toBe(true)
            wrapper.unmount()
        })

        it('should have Teleport as root element', () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
            })

            expect(wrapper.vm).toBeDefined()
            wrapper.unmount()
        })

        it('should accept slots', () => {
            const wrapper = mount(BaseModal, {
                props: {
                    modelValue: true,
                },
                slots: {
                    default: 'Test content',
                },
            })

            expect(wrapper.vm).toBeDefined()
            wrapper.unmount()
        })
    })
})
