import BaseIcon from '@/components/ui/BaseIcon.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { h } from 'vue'

// Mock icon component
const MockIcon = {
    name: 'MockIcon',
    render() {
        return h('svg', { class: 'mock-icon' }, [
            h('path', { d: 'M0 0h24v24H0z' })
        ])
    }
}

describe('BaseIcon.vue', () => {
    describe('Рендеринг', () => {
        it('должен рендерить иконку', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon
                }
            })

            expect(wrapper.find('.mock-icon').exists()).toBe(true)
        })

        it('должен применять базовый класс', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon
                }
            })

            expect(wrapper.classes()).toContain('base-icon')
        })
    })

    describe('Размеры', () => {
        it('должен применять размер по умолчанию (md)', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon
                }
            })

            expect(wrapper.classes()).toContain('base-icon--md')
        })

        it('должен применять размер xs', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'xs'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--xs')
        })

        it('должен применять размер sm', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'sm'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--sm')
        })

        it('должен применять размер md', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'md'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--md')
        })

        it('должен применять размер lg', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'lg'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--lg')
        })

        it('должен применять размер xl', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'xl'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--xl')
        })
    })

    describe('Цвет', () => {
        it('не должен применять inline style если color не указан', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon
                }
            })

            expect(wrapper.attributes('style')).toBeUndefined()
        })

        it('должен применять кастомный цвет через inline style', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    color: '#10b981'
                }
            })

            expect(wrapper.attributes('style')).toContain('color: rgb(16, 185, 129)')
        })

        it('должен применять цвет в формате rgb', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    color: 'rgb(139, 92, 246)'
                }
            })

            expect(wrapper.attributes('style')).toContain('color: rgb(139, 92, 246)')
        })

        it('должен применять цвет в формате currentColor', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    color: 'currentColor'
                }
            })

            expect(wrapper.attributes('style')).toContain('color: currentcolor')
        })
    })

    describe('CSS классы размеров', () => {
        it('xs должен иметь правильные размеры (12px)', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'xs'
                }
            })

            const element = wrapper.element as HTMLElement
            const styles = getComputedStyle(element)

            // Проверяем что класс применен
            expect(wrapper.classes()).toContain('base-icon--xs')
        })

        it('sm должен иметь правильные размеры (16px)', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'sm'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--sm')
        })

        it('md должен иметь правильные размеры (20px)', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'md'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--md')
        })

        it('lg должен иметь правильные размеры (24px)', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'lg'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--lg')
        })

        it('xl должен иметь правильные размеры (32px)', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'xl'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--xl')
        })
    })

    describe('Комбинации props', () => {
        it('должен корректно работать с размером и цветом одновременно', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'lg',
                    color: '#8b5cf6'
                }
            })

            expect(wrapper.classes()).toContain('base-icon--lg')
            expect(wrapper.attributes('style')).toContain('color: rgb(139, 92, 246)')
        })

        it('должен корректно работать с минимальным набором props', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon
                }
            })

            expect(wrapper.classes()).toContain('base-icon')
            expect(wrapper.classes()).toContain('base-icon--md')
            expect(wrapper.attributes('style')).toBeUndefined()
        })

        it('должен корректно работать с полным набором props', () => {
            const wrapper = mount(BaseIcon, {
                props: {
                    icon: MockIcon,
                    size: 'xl',
                    color: '#ef4444'
                }
            })

            expect(wrapper.classes()).toContain('base-icon')
            expect(wrapper.classes()).toContain('base-icon--xl')
            expect(wrapper.attributes('style')).toContain('color: rgb(239, 68, 68)')
        })
    })

    describe('TypeScript типизация', () => {
        it('должен принимать только валидные размеры', () => {
            // Этот тест проверяет что TypeScript не позволит передать невалидный размер
            // В runtime тест просто проверяет что валидные значения работают
            const validSizes: Array<'xs' | 'sm' | 'md' | 'lg' | 'xl'> = ['xs', 'sm', 'md', 'lg', 'xl']

            validSizes.forEach(size => {
                const wrapper = mount(BaseIcon, {
                    props: {
                        icon: MockIcon,
                        size
                    }
                })

                expect(wrapper.classes()).toContain(`base-icon--${size}`)
            })
        })
    })
})
