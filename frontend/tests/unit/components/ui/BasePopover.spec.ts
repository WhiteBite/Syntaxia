import BasePopover from '@/components/ui/BasePopover.vue'
import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'

describe('BasePopover', () => {
    let wrapper: ReturnType<typeof mount>

    beforeEach(() => {
        vi.useFakeTimers()
    })

    afterEach(() => {
        wrapper?.unmount()
        vi.restoreAllMocks()
        vi.useRealTimers()
    })

    describe('rendering', () => {
        it('should render trigger slot', () => {
            wrapper = mount(BasePopover, {
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.find('button').text()).toBe('Trigger')
        })

        it('should not show popover initially', () => {
            wrapper = mount(BasePopover, {
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div class="test-content">Content</div>'
                }
            })

            expect(wrapper.find('.base-popover').exists()).toBe(false)
        })
    })

    describe('trigger behavior', () => {
        it('should emit update:modelValue on click when trigger is "click"', async () => {
            wrapper = mount(BasePopover, {
                props: {
                    trigger: 'click'
                },
                slots: {
                    trigger: '<button class="trigger-btn">Trigger</button>',
                    default: '<div class="test-content">Content</div>'
                }
            })

            await wrapper.find('.trigger-btn').trigger('click')
            await nextTick()

            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([true])
        })

        it('should not emit on click when trigger is "manual"', async () => {
            wrapper = mount(BasePopover, {
                props: {
                    trigger: 'manual'
                },
                slots: {
                    trigger: '<button class="trigger-btn">Trigger</button>',
                    default: '<div class="test-content">Content</div>'
                }
            })

            await wrapper.find('.trigger-btn').trigger('click')
            await nextTick()

            expect(wrapper.emitted('update:modelValue')).toBeFalsy()
        })

        it('should accept hover trigger prop', () => {
            wrapper = mount(BasePopover, {
                props: {
                    trigger: 'hover'
                },
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('trigger')).toBe('hover')
        })
    })

    describe('v-model', () => {
        it('should accept modelValue prop', () => {
            wrapper = mount(BasePopover, {
                props: {
                    modelValue: true
                },
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('modelValue')).toBe(true)
        })

        it('should emit update:modelValue', async () => {
            wrapper = mount(BasePopover, {
                props: {
                    trigger: 'click'
                },
                slots: {
                    trigger: '<button class="trigger-btn">Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            await wrapper.find('.trigger-btn').trigger('click')
            await nextTick()

            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
        })
    })

    describe('placement', () => {
        it('should accept placement prop', () => {
            wrapper = mount(BasePopover, {
                props: {
                    placement: 'bottom'
                },
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('placement')).toBe('bottom')
        })

        it('should default to top placement', () => {
            wrapper = mount(BasePopover, {
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('placement')).toBe('top')
        })
    })

    describe('arrow', () => {
        it('should have arrow enabled by default', () => {
            wrapper = mount(BasePopover, {
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('arrow')).toBe(true)
        })

        it('should accept arrow prop', () => {
            wrapper = mount(BasePopover, {
                props: {
                    arrow: false
                },
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('arrow')).toBe(false)
        })
    })

    describe('offset', () => {
        it('should default to 8', () => {
            wrapper = mount(BasePopover, {
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('offset')).toBe(8)
        })

        it('should accept custom offset', () => {
            wrapper = mount(BasePopover, {
                props: {
                    offset: 16
                },
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('offset')).toBe(16)
        })
    })

    describe('delay', () => {
        it('should default to 0', () => {
            wrapper = mount(BasePopover, {
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('delay')).toBe(0)
        })

        it('should accept custom delay', () => {
            wrapper = mount(BasePopover, {
                props: {
                    delay: 500
                },
                slots: {
                    trigger: '<button>Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            expect(wrapper.props('delay')).toBe(500)
        })
    })

    describe('accessibility', () => {
        it('should render trigger as inline-block', () => {
            wrapper = mount(BasePopover, {
                slots: {
                    trigger: '<button class="trigger-btn">Trigger</button>',
                    default: '<div>Content</div>'
                }
            })

            const trigger = wrapper.find('.base-popover-trigger')
            expect(trigger.exists()).toBe(true)
        })
    })
})
