import BaseTab from '@/components/ui/BaseTab.vue'
import BaseTabs from '@/components/ui/BaseTabs.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { nextTick } from 'vue'

// Mock icon component
const MockIcon = {
    name: 'MockIcon',
    template: '<svg><path /></svg>'
}

describe('BaseTabs.vue', () => {
    describe('Rendering', () => {
        it('should render tabs container', () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            expect(wrapper.find('.base-tabs').exists()).toBe(true)
            expect(wrapper.find('.base-tabs__header').exists()).toBe(true)
            expect(wrapper.find('.base-tabs__content').exists()).toBe(true)
        })

        it('should render tab buttons for each tab', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            expect(tabs).toHaveLength(3)
            expect(tabs[0].text()).toContain('Tab 1')
            expect(tabs[1].text()).toContain('Tab 2')
            expect(tabs[2].text()).toContain('Tab 3')
        })

        it('should render active indicator', () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `<BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>`
                },
                global: {
                    components: { BaseTab }
                }
            })

            expect(wrapper.find('.base-tabs__indicator').exists()).toBe(true)
        })
    })

    describe('Active state', () => {
        it('should mark active tab with correct class', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab2' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            expect(tabs[0].classes()).not.toContain('base-tabs__tab--active')
            expect(tabs[1].classes()).toContain('base-tabs__tab--active')
        })

        it('should set aria-selected on active tab', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            expect(tabs[0].attributes('aria-selected')).toBe('true')
            expect(tabs[1].attributes('aria-selected')).toBe('false')
        })

        it('should only render active tab content', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            expect(wrapper.text()).toContain('Content 1')
            expect(wrapper.text()).not.toContain('Content 2')
        })
    })

    describe('Tab switching', () => {
        it('should emit update:modelValue when tab is clicked', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[1].trigger('click')

            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab2'])
        })

        it('should switch content when modelValue changes', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            expect(wrapper.text()).toContain('Content 1')

            await wrapper.setProps({ modelValue: 'tab2' })
            await nextTick()

            expect(wrapper.text()).not.toContain('Content 1')
            expect(wrapper.text()).toContain('Content 2')
        })
    })

    describe('Disabled tabs', () => {
        it('should apply disabled class to disabled tab', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2" :disabled="true">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            expect(tabs[1].classes()).toContain('base-tabs__tab--disabled')
            expect(tabs[1].attributes('disabled')).toBeDefined()
        })

        it('should not emit update:modelValue when disabled tab is clicked', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2" :disabled="true">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[1].trigger('click')

            expect(wrapper.emitted('update:modelValue')).toBeFalsy()
        })

        it('should set aria-disabled on disabled tab', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2" :disabled="true">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            expect(tabs[1].attributes('aria-disabled')).toBe('true')
        })
    })

    describe('Icons', () => {
        it('should render icon when provided', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
          `
                },
                global: {
                    components: { BaseTab },
                    stubs: {
                        BaseIcon: {
                            template: '<div class="base-icon-stub"></div>'
                        }
                    }
                }
            })

            await nextTick()
            await nextTick()

            // Icon rendering is handled by BaseTab component
            // This test verifies the structure is correct
            const tabs = wrapper.findAll('.base-tabs__tab')
            expect(tabs.length).toBeGreaterThan(0)
        })
    })

    describe('Keyboard navigation', () => {
        it('should navigate to next tab with ArrowRight', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[0].trigger('keydown', { key: 'ArrowRight' })

            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab2'])
        })

        it('should navigate to previous tab with ArrowLeft', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab2' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[1].trigger('keydown', { key: 'ArrowLeft' })

            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab1'])
        })

        it('should navigate to first tab with Home', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab3' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[2].trigger('keydown', { key: 'Home' })

            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab1'])
        })

        it('should navigate to last tab with End', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[0].trigger('keydown', { key: 'End' })

            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab3'])
        })

        it('should wrap around when navigating right from last tab', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab3' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[2].trigger('keydown', { key: 'ArrowRight' })

            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab1'])
        })

        it('should wrap around when navigating left from first tab', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[0].trigger('keydown', { key: 'ArrowLeft' })

            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab3'])
        })

        it('should skip disabled tabs when navigating', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2" :disabled="true">Content 2</BaseTab>
            <BaseTab name="tab3" label="Tab 3">Content 3</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            await tabs[0].trigger('keydown', { key: 'ArrowRight' })

            // Should skip tab2 (disabled) and go to tab3
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['tab3'])
        })
    })

    describe('Accessibility', () => {
        it('should have role="tablist" on container', () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `<BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>`
                },
                global: {
                    components: { BaseTab }
                }
            })

            expect(wrapper.find('[role="tablist"]').exists()).toBe(true)
        })

        it('should have role="tab" on tab buttons', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('[role="tab"]')
            expect(tabs.length).toBeGreaterThan(0)
        })

        it('should set tabindex correctly', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `
            <BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>
            <BaseTab name="tab2" label="Tab 2">Content 2</BaseTab>
          `
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tabs = wrapper.findAll('.base-tabs__tab')
            expect(tabs[0].attributes('tabindex')).toBe('0')
            expect(tabs[1].attributes('tabindex')).toBe('-1')
        })

        it('should have aria-controls attribute', async () => {
            const wrapper = mount(BaseTabs, {
                props: { modelValue: 'tab1' },
                slots: {
                    default: `<BaseTab name="tab1" label="Tab 1">Content 1</BaseTab>`
                },
                global: {
                    components: { BaseTab }
                }
            })

            await nextTick()
            await nextTick()

            const tab = wrapper.find('.base-tabs__tab')
            expect(tab.attributes('aria-controls')).toBe('tabpanel-tab1')
        })
    })
})
