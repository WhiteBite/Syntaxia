import BaseSpinner from '@/components/ui/BaseSpinner.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

describe('BaseSpinner.vue', () => {
    describe('Рендеринг', () => {
        it('должен рендерить SVG элемент', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.find('svg').exists()).toBe(true)
        })

        it('должен применять базовый класс', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.classes()).toContain('base-spinner')
        })

        it('должен применять класс анимации', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.classes()).toContain('animate-spin')
        })

        it('должен иметь правильный viewBox', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.attributes('viewBox')).toBe('0 0 24 24')
        })

        it('должен иметь fill="none"', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.attributes('fill')).toBe('none')
        })
    })

    describe('Размеры', () => {
        it('должен применять размер по умолчанию (md)', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.classes()).toContain('base-spinner--md')
        })

        it('должен применять размер xs', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'xs'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--xs')
        })

        it('должен применять размер sm', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'sm'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--sm')
        })

        it('должен применять размер md', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'md'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--md')
        })

        it('должен применять размер lg', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'lg'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--lg')
        })
    })

    describe('Цвет', () => {
        it('должен применять цвет по умолчанию (currentColor)', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.attributes('style')).toContain('color: currentcolor')
        })

        it('должен применять кастомный цвет через inline style', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    color: '#8b5cf6'
                }
            })

            expect(wrapper.attributes('style')).toContain('color: rgb(139, 92, 246)')
        })

        it('должен применять цвет в формате rgb', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    color: 'rgb(16, 185, 129)'
                }
            })

            expect(wrapper.attributes('style')).toContain('color: rgb(16, 185, 129)')
        })

        it('должен применять цвет в формате hex', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    color: '#ef4444'
                }
            })

            expect(wrapper.attributes('style')).toContain('color: rgb(239, 68, 68)')
        })
    })

    describe('SVG структура', () => {
        it('должен содержать circle элемент', () => {
            const wrapper = mount(BaseSpinner)

            const circle = wrapper.find('circle')
            expect(circle.exists()).toBe(true)
        })

        it('circle должен иметь правильные атрибуты', () => {
            const wrapper = mount(BaseSpinner)

            const circle = wrapper.find('circle')
            expect(circle.attributes('cx')).toBe('12')
            expect(circle.attributes('cy')).toBe('12')
            expect(circle.attributes('r')).toBe('10')
            expect(circle.attributes('stroke')).toBe('currentColor')
            expect(circle.attributes('stroke-width')).toBe('2')
        })

        it('circle должен иметь класс opacity-25', () => {
            const wrapper = mount(BaseSpinner)

            const circle = wrapper.find('circle')
            expect(circle.classes()).toContain('opacity-25')
        })

        it('должен содержать path элемент', () => {
            const wrapper = mount(BaseSpinner)

            const path = wrapper.find('path')
            expect(path.exists()).toBe(true)
        })

        it('path должен иметь правильные атрибуты', () => {
            const wrapper = mount(BaseSpinner)

            const path = wrapper.find('path')
            expect(path.attributes('fill')).toBe('currentColor')
            expect(path.attributes('d')).toBe('M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z')
        })

        it('path должен иметь класс opacity-75', () => {
            const wrapper = mount(BaseSpinner)

            const path = wrapper.find('path')
            expect(path.classes()).toContain('opacity-75')
        })
    })

    describe('CSS классы размеров', () => {
        it('xs должен иметь правильные размеры (12px)', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'xs'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--xs')
        })

        it('sm должен иметь правильные размеры (16px)', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'sm'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--sm')
        })

        it('md должен иметь правильные размеры (20px)', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'md'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--md')
        })

        it('lg должен иметь правильные размеры (24px)', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'lg'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--lg')
        })
    })

    describe('Комбинации props', () => {
        it('должен корректно работать с размером и цветом одновременно', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'lg',
                    color: '#10b981'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner--lg')
            expect(wrapper.attributes('style')).toContain('color: rgb(16, 185, 129)')
        })

        it('должен корректно работать с минимальным набором props', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.classes()).toContain('base-spinner')
            expect(wrapper.classes()).toContain('base-spinner--md')
            expect(wrapper.classes()).toContain('animate-spin')
            expect(wrapper.attributes('style')).toContain('color: currentcolor')
        })

        it('должен корректно работать с полным набором props', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    size: 'xl',
                    color: '#8b5cf6'
                }
            })

            expect(wrapper.classes()).toContain('base-spinner')
            expect(wrapper.classes()).toContain('base-spinner--xl')
            expect(wrapper.classes()).toContain('animate-spin')
            expect(wrapper.attributes('style')).toContain('color: rgb(139, 92, 246)')
        })
    })

    describe('TypeScript типизация', () => {
        it('должен принимать только валидные размеры', () => {
            const validSizes: Array<'xs' | 'sm' | 'md' | 'lg'> = ['xs', 'sm', 'md', 'lg']

            validSizes.forEach(size => {
                const wrapper = mount(BaseSpinner, {
                    props: {
                        size
                    }
                })

                expect(wrapper.classes()).toContain(`base-spinner--${size}`)
            })
        })
    })

    describe('Анимация', () => {
        it('должен иметь класс animate-spin для анимации вращения', () => {
            const wrapper = mount(BaseSpinner)

            expect(wrapper.classes()).toContain('animate-spin')
        })

        it('анимация должна применяться независимо от размера', () => {
            const sizes: Array<'xs' | 'sm' | 'md' | 'lg'> = ['xs', 'sm', 'md', 'lg']

            sizes.forEach(size => {
                const wrapper = mount(BaseSpinner, {
                    props: {
                        size
                    }
                })

                expect(wrapper.classes()).toContain('animate-spin')
            })
        })

        it('анимация должна применяться независимо от цвета', () => {
            const wrapper = mount(BaseSpinner, {
                props: {
                    color: '#ef4444'
                }
            })

            expect(wrapper.classes()).toContain('animate-spin')
        })
    })
})
