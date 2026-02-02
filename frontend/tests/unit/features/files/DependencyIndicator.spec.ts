import { useFileStore } from '@/features/files/model/file.store'
import DependencyIndicator from '@/features/files/ui/DependencyIndicator.vue'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

describe('DependencyIndicator', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('should not render when file has no dependencies', () => {
        const wrapper = mount(DependencyIndicator, {
            props: {
                filePath: '/test/file.ts'
            }
        })

        expect(wrapper.find('.tree-row-deps').exists()).toBe(false)
    })

    it('should render outgoing dependencies indicator', () => {
        const fileStore = useFileStore()

        // Mock dependencies
        fileStore.allFileDependencies.set('/test/file.ts', {
            incoming: [],
            outgoing: ['/test/dep1.ts', '/test/dep2.ts']
        })

        const wrapper = mount(DependencyIndicator, {
            props: {
                filePath: '/test/file.ts'
            }
        })

        expect(wrapper.find('.tree-row-deps').exists()).toBe(true)
        expect(wrapper.find('.dep-outgoing').exists()).toBe(true)
        expect(wrapper.find('.dep-count').text()).toBe('2')
    })

    it('should render incoming dependencies indicator', () => {
        const fileStore = useFileStore()

        // Mock dependencies
        fileStore.allFileDependencies.set('/test/file.ts', {
            incoming: ['/test/parent1.ts', '/test/parent2.ts', '/test/parent3.ts'],
            outgoing: []
        })

        const wrapper = mount(DependencyIndicator, {
            props: {
                filePath: '/test/file.ts'
            }
        })

        expect(wrapper.find('.tree-row-deps').exists()).toBe(true)
        expect(wrapper.find('.dep-incoming').exists()).toBe(true)
        expect(wrapper.find('.dep-count').text()).toBe('3')
    })

    it('should render both incoming and outgoing indicators', () => {
        const fileStore = useFileStore()

        // Mock dependencies
        fileStore.allFileDependencies.set('/test/file.ts', {
            incoming: ['/test/parent.ts'],
            outgoing: ['/test/dep.ts']
        })

        const wrapper = mount(DependencyIndicator, {
            props: {
                filePath: '/test/file.ts'
            }
        })

        expect(wrapper.find('.tree-row-deps').exists()).toBe(true)
        expect(wrapper.find('.dep-incoming').exists()).toBe(true)
        expect(wrapper.find('.dep-outgoing').exists()).toBe(true)
    })

    it('should emit showDependencies event on click', async () => {
        const fileStore = useFileStore()

        // Mock dependencies
        fileStore.allFileDependencies.set('/test/file.ts', {
            incoming: ['/test/parent.ts'],
            outgoing: ['/test/dep.ts']
        })

        const wrapper = mount(DependencyIndicator, {
            props: {
                filePath: '/test/file.ts'
            }
        })

        // Click on outgoing indicator
        await wrapper.find('.dep-outgoing').trigger('click')
        expect(wrapper.emitted('showDependencies')).toBeTruthy()
        expect(wrapper.emitted('showDependencies')?.[0]).toEqual(['outgoing'])

        // Click on incoming indicator
        await wrapper.find('.dep-incoming').trigger('click')
        expect(wrapper.emitted('showDependencies')?.[1]).toEqual(['incoming'])
    })

    it('should display correct count for multiple dependencies', () => {
        const fileStore = useFileStore()

        // Mock many dependencies
        const outgoing = Array.from({ length: 10 }, (_, i) => `/test/dep${i}.ts`)
        fileStore.allFileDependencies.set('/test/file.ts', {
            incoming: [],
            outgoing
        })

        const wrapper = mount(DependencyIndicator, {
            props: {
                filePath: '/test/file.ts'
            }
        })

        expect(wrapper.find('.dep-count').text()).toBe('10')
    })

    it('should stop event propagation on click', async () => {
        const fileStore = useFileStore()

        fileStore.allFileDependencies.set('/test/file.ts', {
            incoming: [],
            outgoing: ['/test/dep.ts']
        })

        const wrapper = mount(DependencyIndicator, {
            props: {
                filePath: '/test/file.ts'
            }
        })

        const clickHandler = vi.fn()
        wrapper.element.parentElement?.addEventListener('click', clickHandler)

        await wrapper.find('.dep-group').trigger('click')

        // Parent click handler should not be called due to .stop modifier
        expect(clickHandler).not.toHaveBeenCalled()
    })
})
