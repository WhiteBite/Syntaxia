import BaseEmptyState from '@/components/ui/BaseEmptyState.vue'
import BaseIcon from '@/components/ui/BaseIcon.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { h } from 'vue'

// Mock icon component
const MockIcon = {
    name: 'MockIcon',
    render: () => h('svg', { class: 'mock-icon' }, [h('path')])
}

describe('BaseEmptyState', () => {
    it('renders with title and description', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'No data',
                description: 'No data available'
            }
        })

        expect(wrapper.find('.base-empty-state__title').text()).toBe('No data')
        expect(wrapper.find('.base-empty-state__description').text()).toBe('No data available')
    })

    it('renders with icon component', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                icon: MockIcon,
                title: 'Empty'
            },
            global: {
                components: {
                    BaseIcon
                }
            }
        })

        expect(wrapper.findComponent(BaseIcon).exists()).toBe(true)
    })

    it('renders custom icon slot', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty'
            },
            slots: {
                icon: '<div class="custom-icon">Custom Icon</div>'
            }
        })

        expect(wrapper.find('.custom-icon').exists()).toBe(true)
        expect(wrapper.find('.custom-icon').text()).toBe('Custom Icon')
    })

    it('renders action slot', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty'
            },
            slots: {
                action: '<button class="test-button">Click me</button>'
            }
        })

        expect(wrapper.find('.base-empty-state__action').exists()).toBe(true)
        expect(wrapper.find('.test-button').exists()).toBe(true)
    })

    it('applies small size class', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty',
                size: 'sm'
            }
        })

        expect(wrapper.find('.base-empty-state--sm').exists()).toBe(true)
    })

    it('applies medium size class by default', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty'
            }
        })

        expect(wrapper.find('.base-empty-state--md').exists()).toBe(true)
    })

    it('applies large size class', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty',
                size: 'lg'
            }
        })

        expect(wrapper.find('.base-empty-state--lg').exists()).toBe(true)
    })

    it('does not render title when not provided', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                description: 'Only description'
            }
        })

        expect(wrapper.find('.base-empty-state__title').exists()).toBe(false)
    })

    it('does not render description when not provided', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Only title'
            }
        })

        expect(wrapper.find('.base-empty-state__description').exists()).toBe(false)
    })

    it('does not render action container when no action slot provided', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty'
            }
        })

        expect(wrapper.find('.base-empty-state__action').exists()).toBe(false)
    })

    it('applies with-action class when action slot is provided', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty'
            },
            slots: {
                action: '<button>Action</button>'
            }
        })

        expect(wrapper.find('.base-empty-state--with-action').exists()).toBe(true)
    })

    it('renders multiple action buttons', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                title: 'Empty'
            },
            slots: {
                action: `
          <button class="btn-1">Action 1</button>
          <button class="btn-2">Action 2</button>
        `
            }
        })

        expect(wrapper.find('.btn-1').exists()).toBe(true)
        expect(wrapper.find('.btn-2').exists()).toBe(true)
    })

    it('computes correct icon size for small variant', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                icon: MockIcon,
                title: 'Empty',
                size: 'sm'
            },
            global: {
                components: {
                    BaseIcon
                }
            }
        })

        const iconComponent = wrapper.findComponent(BaseIcon)
        expect(iconComponent.props('size')).toBe('md')
    })

    it('computes correct icon size for medium variant', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                icon: MockIcon,
                title: 'Empty',
                size: 'md'
            },
            global: {
                components: {
                    BaseIcon
                }
            }
        })

        const iconComponent = wrapper.findComponent(BaseIcon)
        expect(iconComponent.props('size')).toBe('lg')
    })

    it('computes correct icon size for large variant', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                icon: MockIcon,
                title: 'Empty',
                size: 'lg'
            },
            global: {
                components: {
                    BaseIcon
                }
            }
        })

        const iconComponent = wrapper.findComponent(BaseIcon)
        expect(iconComponent.props('size')).toBe('xl')
    })

    it('renders with all props and slots', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                icon: MockIcon,
                title: 'No results',
                description: 'Try adjusting your filters',
                size: 'lg'
            },
            slots: {
                action: '<button class="action-btn">Clear Filters</button>'
            },
            global: {
                components: {
                    BaseIcon
                }
            }
        })

        expect(wrapper.find('.base-empty-state--lg').exists()).toBe(true)
        expect(wrapper.findComponent(BaseIcon).exists()).toBe(true)
        expect(wrapper.find('.base-empty-state__title').text()).toBe('No results')
        expect(wrapper.find('.base-empty-state__description').text()).toBe('Try adjusting your filters')
        expect(wrapper.find('.action-btn').exists()).toBe(true)
    })

    it('has correct structure and classes', () => {
        const wrapper = mount(BaseEmptyState, {
            props: {
                icon: MockIcon,
                title: 'Empty',
                description: 'Description'
            },
            global: {
                components: {
                    BaseIcon
                }
            }
        })

        expect(wrapper.find('.base-empty-state').exists()).toBe(true)
        expect(wrapper.find('.base-empty-state__icon').exists()).toBe(true)
        expect(wrapper.find('.base-empty-state__title').exists()).toBe(true)
        expect(wrapper.find('.base-empty-state__description').exists()).toBe(true)
    })
})
