import BaseSkeleton from '@/components/ui/BaseSkeleton.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

describe('BaseSkeleton.vue', () => {
    describe('Rendering', () => {
        it('should render skeleton with default props', () => {
            const wrapper = mount(BaseSkeleton)

            expect(wrapper.find('.base-skeleton').exists()).toBe(true)
            expect(wrapper.classes()).toContain('base-skeleton')
            expect(wrapper.classes()).toContain('base-skeleton--text')
            expect(wrapper.classes()).toContain('base-skeleton--animated')
        })

        it('should render with role and aria attributes', () => {
            const wrapper = mount(BaseSkeleton)

            const skeleton = wrapper.find('.base-skeleton')
            expect(skeleton.attributes('role')).toBe('status')
            expect(skeleton.attributes('aria-busy')).toBe('true')
            expect(skeleton.attributes('aria-label')).toBe('Loading content')
        })

        it('should render multiple skeletons when count > 1', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 3 }
            })

            const skeletons = wrapper.findAll('.base-skeleton')
            expect(skeletons).toHaveLength(3)
            expect(wrapper.find('.base-skeleton-group').exists()).toBe(true)
        })

        it('should render single skeleton when count is 1', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 1 }
            })

            expect(wrapper.find('.base-skeleton-group').exists()).toBe(false)
            expect(wrapper.find('.base-skeleton').exists()).toBe(true)
        })
    })

    describe('Variants', () => {
        const variants = ['text', 'circle', 'rect', 'card'] as const

        variants.forEach(variant => {
            it(`should apply ${variant} variant class`, () => {
                const wrapper = mount(BaseSkeleton, {
                    props: { variant }
                })

                expect(wrapper.classes()).toContain(`base-skeleton--${variant}`)
            })
        })

        it('should apply default height for text variant', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { variant: 'text' }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('height: 1rem')
        })

        it('should apply default height for card variant', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { variant: 'card' }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('height: 200px')
        })

        it('should apply width as height for circle variant', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { variant: 'circle', width: '50px' }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('width: 50px')
            expect(style).toContain('height: 50px')
        })
    })

    describe('Width and Height', () => {
        it('should apply width as string', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { width: '200px' }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('width: 200px')
        })

        it('should apply width as number', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { width: 150 }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('width: 150px')
        })

        it('should apply height as string', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { height: '100px' }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('height: 100px')
        })

        it('should apply height as number', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { height: 80 }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('height: 80px')
        })

        it('should apply percentage width', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { width: '100%' }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('width: 100%')
        })

        it('should apply custom width and height together', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { width: '300px', height: '150px' }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('width: 300px')
            expect(style).toContain('height: 150px')
        })
    })

    describe('Animation', () => {
        it('should have animated class by default', () => {
            const wrapper = mount(BaseSkeleton)

            expect(wrapper.classes()).toContain('base-skeleton--animated')
        })

        it('should not have animated class when animated is false', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { animated: false }
            })

            expect(wrapper.classes()).not.toContain('base-skeleton--animated')
        })

        it('should apply animated class when animated is true', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { animated: true }
            })

            expect(wrapper.classes()).toContain('base-skeleton--animated')
        })
    })

    describe('Count and Gap', () => {
        it('should render correct number of skeletons', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 5 }
            })

            const skeletons = wrapper.findAll('.base-skeleton')
            expect(skeletons).toHaveLength(5)
        })

        it('should apply default gap', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 3 }
            })

            const group = wrapper.find('.base-skeleton-group')
            const style = group.attributes('style')
            expect(style).toContain('gap: 8px')
        })

        it('should apply custom gap', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 3, gap: 16 }
            })

            const group = wrapper.find('.base-skeleton-group')
            const style = group.attributes('style')
            expect(style).toContain('gap: 16px')
        })
    })

    describe('Aria Labels', () => {
        it('should have correct aria-label for single skeleton', () => {
            const wrapper = mount(BaseSkeleton)

            const skeleton = wrapper.find('.base-skeleton')
            expect(skeleton.attributes('aria-label')).toBe('Loading content')
        })

        it('should have correct aria-label for multiple skeletons', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 3 }
            })

            const skeleton = wrapper.find('.base-skeleton')
            expect(skeleton.attributes('aria-label')).toBe('Loading 3 items')
        })

        it('should update aria-label when count changes', async () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 2 }
            })

            expect(wrapper.find('.base-skeleton').attributes('aria-label')).toBe('Loading 2 items')

            await wrapper.setProps({ count: 5 })

            expect(wrapper.find('.base-skeleton').attributes('aria-label')).toBe('Loading 5 items')
        })
    })

    describe('Complex scenarios', () => {
        it('should handle all props together', () => {
            const wrapper = mount(BaseSkeleton, {
                props: {
                    variant: 'rect',
                    width: '250px',
                    height: '120px',
                    count: 4,
                    animated: false,
                    gap: 12
                }
            })

            const skeletons = wrapper.findAll('.base-skeleton')
            expect(skeletons).toHaveLength(4)
            expect(skeletons[0].classes()).toContain('base-skeleton--rect')
            expect(skeletons[0].classes()).not.toContain('base-skeleton--animated')

            const group = wrapper.find('.base-skeleton-group')
            expect(group.attributes('style')).toContain('gap: 12px')

            const skeleton = wrapper.find('.base-skeleton')
            const style = skeleton.attributes('style')
            expect(style).toContain('width: 250px')
            expect(style).toContain('height: 120px')
        })

        it('should handle circle variant with only width', () => {
            const wrapper = mount(BaseSkeleton, {
                props: {
                    variant: 'circle',
                    width: '64px'
                }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('width: 64px')
            expect(style).toContain('height: 64px')
        })

        it('should override default height when provided', () => {
            const wrapper = mount(BaseSkeleton, {
                props: {
                    variant: 'text',
                    height: '24px'
                }
            })

            const style = wrapper.find('.base-skeleton').attributes('style')
            expect(style).toContain('height: 24px')
        })

        it('should handle zero count gracefully', () => {
            const wrapper = mount(BaseSkeleton, {
                props: { count: 0 }
            })

            // Should render single skeleton when count is 0 or 1
            expect(wrapper.find('.base-skeleton').exists()).toBe(true)
            expect(wrapper.find('.base-skeleton-group').exists()).toBe(false)
        })
    })
})
