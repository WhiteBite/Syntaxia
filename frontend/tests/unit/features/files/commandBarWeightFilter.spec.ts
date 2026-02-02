import { useFileStore } from '@/features/files/model/file.store'
import CommandBar from '@/features/files/ui/CommandBar.vue'
import type { FileNode } from '@/types/domain'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('CommandBar - Weight Filter UI', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    const createFileNode = (path: string, name: string, size: number): FileNode => ({
        path,
        name,
        isDir: false,
        size,
        depth: 0,
    })

    it('should display token stats bar when files exist', () => {
        const store = useFileStore()

        // Create files with different weights
        const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000) // 12.5K tokens
        const heavyFile = createFileNode('/heavy.txt', 'heavy.txt', 250000) // 62.5K tokens
        const criticalFile = createFileNode('/critical.txt', 'critical.txt', 500000) // 125K tokens

        store.setFileTree([mediumFile, heavyFile, criticalFile])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        // Token stats bar should be visible
        expect(wrapper.find('.token-stats-bar').exists()).toBe(true)

        // Should show all three weight segments
        expect(wrapper.find('.token-stat-segment--medium').exists()).toBe(true)
        expect(wrapper.find('.token-stat-segment--heavy').exists()).toBe(true)
        expect(wrapper.find('.token-stat-segment--critical').exists()).toBe(true)
    })

    it('should not display token stats bar when no files exist', () => {
        const store = useFileStore()
        store.setFileTree([])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        expect(wrapper.find('.token-stats-bar').exists()).toBe(false)
    })

    it('should show correct counts for each weight category', () => {
        const store = useFileStore()

        // 2 medium, 1 heavy, 1 critical
        const medium1 = createFileNode('/m1.txt', 'm1.txt', 50000)
        const medium2 = createFileNode('/m2.txt', 'm2.txt', 60000)
        const heavy1 = createFileNode('/h1.txt', 'h1.txt', 250000)
        const critical1 = createFileNode('/c1.txt', 'c1.txt', 500000)

        store.setFileTree([medium1, medium2, heavy1, critical1])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        // Check counts
        const mediumSegment = wrapper.find('.token-stat-segment--medium .token-stat-count')
        const heavySegment = wrapper.find('.token-stat-segment--heavy .token-stat-count')
        const criticalSegment = wrapper.find('.token-stat-segment--critical .token-stat-count')

        expect(mediumSegment.text()).toBe('2')
        expect(heavySegment.text()).toBe('1')
        expect(criticalSegment.text()).toBe('1')
    })

    it('should apply active class when filter is set', async () => {
        const store = useFileStore()

        const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000)
        store.setFileTree([mediumFile])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        // Initially no active filter
        expect(wrapper.find('.token-stat-segment--medium.active').exists()).toBe(false)

        // Set filter
        store.setWeightFilter('medium')
        await wrapper.vm.$nextTick()

        // Should have active class
        expect(wrapper.find('.token-stat-segment--medium.active').exists()).toBe(true)
    })

    it('should show clear button when filter is active', async () => {
        const store = useFileStore()

        const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000)
        store.setFileTree([mediumFile])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        // Initially no clear button
        expect(wrapper.find('.clear-weight-filter').exists()).toBe(false)

        // Set filter
        store.setWeightFilter('medium')
        await wrapper.vm.$nextTick()

        // Clear button should appear
        expect(wrapper.find('.clear-weight-filter').exists()).toBe(true)
    })

    it('should show visible files count when filter is active', async () => {
        const store = useFileStore()

        const small = createFileNode('/small.txt', 'small.txt', 20000) // 5K tokens
        const medium = createFileNode('/medium.txt', 'medium.txt', 50000) // 12.5K tokens
        const heavy = createFileNode('/heavy.txt', 'heavy.txt', 250000) // 62.5K tokens

        store.setFileTree([small, medium, heavy])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        // Set filter to medium (should show 2 files: medium + heavy)
        store.setWeightFilter('medium')
        await wrapper.vm.$nextTick()

        // Should show filter info
        expect(wrapper.find('.weight-filter-info').exists()).toBe(true)
        expect(wrapper.find('.filter-info-text').text()).toContain('2')
    })

    it('should toggle filter when clicking same segment twice', async () => {
        const store = useFileStore()

        const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000)
        store.setFileTree([mediumFile])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        const mediumButton = wrapper.find('.token-stat-segment--medium')

        // First click: activate
        await mediumButton.trigger('click')
        expect(store.weightFilter).toBe('medium')

        // Second click: deactivate (handled by handleWeightFilterClick)
        await mediumButton.trigger('click')
        expect(store.weightFilter).toBe('none')
    })

    it('should clear filter when clicking clear button', async () => {
        const store = useFileStore()

        const mediumFile = createFileNode('/medium.txt', 'medium.txt', 50000)
        store.setFileTree([mediumFile])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        // Set filter
        store.setWeightFilter('medium')
        await wrapper.vm.$nextTick()

        // Click clear button
        const clearButton = wrapper.find('.clear-weight-filter')
        await clearButton.trigger('click')

        expect(store.weightFilter).toBe('none')
    })

    it('should only show segments for weight categories that have files', () => {
        const store = useFileStore()

        // Only medium files, no heavy or critical
        const medium1 = createFileNode('/m1.txt', 'm1.txt', 50000)
        const medium2 = createFileNode('/m2.txt', 'm2.txt', 60000)

        store.setFileTree([medium1, medium2])

        const wrapper = mount(CommandBar, {
            props: {
                selectedCount: 0,
                isBuilding: false,
            },
        })

        expect(wrapper.find('.token-stat-segment--medium').exists()).toBe(true)
        expect(wrapper.find('.token-stat-segment--heavy').exists()).toBe(false)
        expect(wrapper.find('.token-stat-segment--critical').exists()).toBe(false)
    })
})
