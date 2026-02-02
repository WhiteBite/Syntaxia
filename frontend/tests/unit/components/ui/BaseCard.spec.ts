import BaseCard from '@/components/ui/BaseCard.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

describe('BaseCard.vue', () => {
    describe('Рендеринг', () => {
        it('должен рендерить карточку', () => {
            const wrapper = mount(BaseCard)
            expect(wrapper.find('.base-card').exists()).toBe(true)
        })

        it('должен рендерить body по умолчанию', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    default: 'Card content'
                }
            })
            expect(wrapper.find('.base-card__body').exists()).toBe(true)
            expect(wrapper.find('.base-card__body').text()).toBe('Card content')
        })

        it('должен рендерить header когда передан slot', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    header: 'Card Header'
                }
            })
            expect(wrapper.find('.base-card__header').exists()).toBe(true)
            expect(wrapper.find('.base-card__header').text()).toBe('Card Header')
        })

        it('должен рендерить footer когда передан slot', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    footer: 'Card Footer'
                }
            })
            expect(wrapper.find('.base-card__footer').exists()).toBe(true)
            expect(wrapper.find('.base-card__footer').text()).toBe('Card Footer')
        })

        it('не должен рендерить header если slot не передан', () => {
            const wrapper = mount(BaseCard)
            expect(wrapper.find('.base-card__header').exists()).toBe(false)
        })

        it('не должен рендерить footer если slot не передан', () => {
            const wrapper = mount(BaseCard)
            expect(wrapper.find('.base-card__footer').exists()).toBe(false)
        })
    })

    describe('Props', () => {
        it('должен применять класс interactive когда interactive=true', () => {
            const wrapper = mount(BaseCard, {
                props: {
                    interactive: true
                }
            })
            expect(wrapper.find('.base-card--interactive').exists()).toBe(true)
        })

        it('не должен применять класс interactive по умолчанию', () => {
            const wrapper = mount(BaseCard)
            expect(wrapper.find('.base-card--interactive').exists()).toBe(false)
        })

        it('должен применять класс glass когда glass=true', () => {
            const wrapper = mount(BaseCard, {
                props: {
                    glass: true
                }
            })
            expect(wrapper.find('.base-card--glass').exists()).toBe(true)
        })

        it('не должен применять класс glass по умолчанию', () => {
            const wrapper = mount(BaseCard)
            expect(wrapper.find('.base-card--glass').exists()).toBe(false)
        })

        it('должен применять оба класса одновременно', () => {
            const wrapper = mount(BaseCard, {
                props: {
                    interactive: true,
                    glass: true
                }
            })
            expect(wrapper.find('.base-card--interactive').exists()).toBe(true)
            expect(wrapper.find('.base-card--glass').exists()).toBe(true)
        })
    })

    describe('Slots', () => {
        it('должен рендерить все slots одновременно', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    header: '<div class="test-header">Header</div>',
                    default: '<div class="test-body">Body</div>',
                    footer: '<div class="test-footer">Footer</div>'
                }
            })
            expect(wrapper.find('.test-header').exists()).toBe(true)
            expect(wrapper.find('.test-body').exists()).toBe(true)
            expect(wrapper.find('.test-footer').exists()).toBe(true)
        })

        it('должен рендерить сложный контент в header', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    header: '<h2>Title</h2><p>Subtitle</p>'
                }
            })
            expect(wrapper.find('.base-card__header h2').text()).toBe('Title')
            expect(wrapper.find('.base-card__header p').text()).toBe('Subtitle')
        })

        it('должен рендерить сложный контент в body', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    default: '<div><p>Paragraph 1</p><p>Paragraph 2</p></div>'
                }
            })
            expect(wrapper.findAll('.base-card__body p')).toHaveLength(2)
        })

        it('должен рендерить сложный контент в footer', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    footer: '<button>Action 1</button><button>Action 2</button>'
                }
            })
            expect(wrapper.findAll('.base-card__footer button')).toHaveLength(2)
        })
    })

    describe('Структура', () => {
        it('должен иметь правильную структуру с header', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    header: 'Header',
                    default: 'Body'
                }
            })
            const card = wrapper.find('.base-card')
            const children = card.element.children
            expect(children[0].classList.contains('base-card__header')).toBe(true)
            expect(children[1].classList.contains('base-card__body')).toBe(true)
        })

        it('должен иметь правильную структуру с footer', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    default: 'Body',
                    footer: 'Footer'
                }
            })
            const card = wrapper.find('.base-card')
            const children = card.element.children
            expect(children[0].classList.contains('base-card__body')).toBe(true)
            expect(children[1].classList.contains('base-card__footer')).toBe(true)
        })

        it('должен иметь правильную структуру со всеми секциями', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    header: 'Header',
                    default: 'Body',
                    footer: 'Footer'
                }
            })
            const card = wrapper.find('.base-card')
            const children = card.element.children
            expect(children).toHaveLength(3)
            expect(children[0].classList.contains('base-card__header')).toBe(true)
            expect(children[1].classList.contains('base-card__body')).toBe(true)
            expect(children[2].classList.contains('base-card__footer')).toBe(true)
        })
    })

    describe('Edge cases', () => {
        it('должен корректно работать с пустым body', () => {
            const wrapper = mount(BaseCard)
            expect(wrapper.find('.base-card__body').text()).toBe('')
        })

        it('должен корректно работать с пустыми slots', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    header: '',
                    default: '',
                    footer: ''
                }
            })
            // Пустые slots всё равно рендерятся (Vue поведение)
            // Проверяем что компонент не падает с пустыми slots
            expect(wrapper.find('.base-card').exists()).toBe(true)
            expect(wrapper.find('.base-card__body').exists()).toBe(true)
        })

        it('должен корректно работать со всеми props и slots', () => {
            const wrapper = mount(BaseCard, {
                props: {
                    interactive: true,
                    glass: true
                },
                slots: {
                    header: 'Header',
                    default: 'Body',
                    footer: 'Footer'
                }
            })
            expect(wrapper.find('.base-card--interactive').exists()).toBe(true)
            expect(wrapper.find('.base-card--glass').exists()).toBe(true)
            expect(wrapper.find('.base-card__header').text()).toBe('Header')
            expect(wrapper.find('.base-card__body').text()).toBe('Body')
            expect(wrapper.find('.base-card__footer').text()).toBe('Footer')
        })

        it('должен корректно работать с HTML entities', () => {
            const wrapper = mount(BaseCard, {
                slots: {
                    default: '&lt;div&gt;Test&lt;/div&gt;'
                }
            })
            expect(wrapper.find('.base-card__body').text()).toContain('Test')
        })
    })

    describe('CSS классы', () => {
        it('должен всегда иметь базовый класс', () => {
            const wrapper = mount(BaseCard)
            expect(wrapper.classes()).toContain('base-card')
        })

        it('должен иметь только базовый класс без props', () => {
            const wrapper = mount(BaseCard)
            const classes = wrapper.classes()
            expect(classes).toContain('base-card')
            expect(classes).not.toContain('base-card--interactive')
            expect(classes).not.toContain('base-card--glass')
        })

        it('должен корректно комбинировать классы', () => {
            const wrapper = mount(BaseCard, {
                props: {
                    interactive: true,
                    glass: true
                }
            })
            const classes = wrapper.classes()
            expect(classes).toContain('base-card')
            expect(classes).toContain('base-card--interactive')
            expect(classes).toContain('base-card--glass')
        })
    })
})
