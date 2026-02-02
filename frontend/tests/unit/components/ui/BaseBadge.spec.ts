import BaseBadge from '@/components/ui/BaseBadge.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

describe('BaseBadge.vue', () => {
    describe('Рендеринг', () => {
        it('должен рендерить badge', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge').exists()).toBe(true)
        })

        it('должен рендерить текст из default slot', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Test Badge'
                }
            })
            expect(wrapper.find('.base-badge__text').text()).toBe('Test Badge')
        })

        it('должен рендерить иконку из icon prop', () => {
            const MockIcon = { name: 'MockIcon', template: '<svg><path /></svg>' }
            const wrapper = mount(BaseBadge, {
                props: {
                    icon: MockIcon
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge__icon').exists()).toBe(true)
        })

        it('должен рендерить иконку из icon slot', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Badge',
                    icon: '<svg class="test-icon"><path /></svg>'
                }
            })
            expect(wrapper.find('.base-badge__icon').exists()).toBe(true)
            expect(wrapper.find('.test-icon').exists()).toBe(true)
        })

        it('не должен рендерить иконку если не передана', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge__icon').exists()).toBe(false)
        })
    })

    describe('Variants', () => {
        it('должен применять variant по умолчанию (secondary)', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--secondary').exists()).toBe(true)
        })

        it('должен применять variant primary', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'primary'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--primary').exists()).toBe(true)
        })

        it('должен применять variant success', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'success'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--success').exists()).toBe(true)
        })

        it('должен применять variant warning', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'warning'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--warning').exists()).toBe(true)
        })

        it('должен применять variant danger', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'danger'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--danger').exists()).toBe(true)
        })

        it('должен применять variant info', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'info'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--info').exists()).toBe(true)
        })

        it('должен применять variant accent', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'accent'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--accent').exists()).toBe(true)
        })
    })

    describe('Sizes', () => {
        it('должен применять size по умолчанию (sm)', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--sm').exists()).toBe(true)
        })

        it('должен применять size xs', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    size: 'xs'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--xs').exists()).toBe(true)
        })

        it('должен применять size sm', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    size: 'sm'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--sm').exists()).toBe(true)
        })

        it('должен применять size md', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    size: 'md'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--md').exists()).toBe(true)
        })
    })

    describe('Icon slot приоритет', () => {
        it('icon slot должен иметь приоритет над icon prop', () => {
            const MockIcon = { name: 'MockIcon', template: '<svg><path /></svg>' }
            const wrapper = mount(BaseBadge, {
                props: {
                    icon: MockIcon
                },
                slots: {
                    default: 'Badge',
                    icon: '<svg class="slot-icon"><path /></svg>'
                }
            })
            expect(wrapper.find('.slot-icon').exists()).toBe(true)
        })
    })

    describe('Комбинации props', () => {
        it('должен корректно работать с variant и size', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'success',
                    size: 'md'
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--success').exists()).toBe(true)
            expect(wrapper.find('.base-badge--md').exists()).toBe(true)
        })

        it('должен корректно работать с variant, size и icon', () => {
            const MockIcon = { name: 'MockIcon', template: '<svg><path /></svg>' }
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'primary',
                    size: 'xs',
                    icon: MockIcon
                },
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge--primary').exists()).toBe(true)
            expect(wrapper.find('.base-badge--xs').exists()).toBe(true)
            expect(wrapper.find('.base-badge__icon').exists()).toBe(true)
        })

        it('должен корректно работать со всеми props', () => {
            const MockIcon = { name: 'MockIcon', template: '<svg><path /></svg>' }
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'warning',
                    size: 'md',
                    icon: MockIcon
                },
                slots: {
                    default: 'Important'
                }
            })
            expect(wrapper.find('.base-badge--warning').exists()).toBe(true)
            expect(wrapper.find('.base-badge--md').exists()).toBe(true)
            expect(wrapper.find('.base-badge__icon').exists()).toBe(true)
            expect(wrapper.find('.base-badge__text').text()).toBe('Important')
        })
    })

    describe('Edge cases', () => {
        it('должен корректно работать с пустым текстом', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: ''
                }
            })
            expect(wrapper.find('.base-badge__text').text()).toBe('')
        })

        it('должен корректно работать с длинным текстом', () => {
            const longText = 'Very Long Badge Text That Should Not Break Layout'
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: longText
                }
            })
            expect(wrapper.find('.base-badge__text').text()).toBe(longText)
        })

        it('должен корректно работать с числами', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: '42'
                }
            })
            expect(wrapper.find('.base-badge__text').text()).toBe('42')
        })

        it('должен корректно работать с Unicode символами', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: '🚀 Rocket'
                }
            })
            expect(wrapper.find('.base-badge__text').text()).toBe('🚀 Rocket')
        })

        it('должен корректно работать с HTML entities', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: '&lt;Tag&gt;'
                }
            })
            expect(wrapper.find('.base-badge__text').text()).toContain('Tag')
        })
    })

    describe('CSS классы', () => {
        it('должен всегда иметь базовый класс', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.classes()).toContain('base-badge')
        })

        it('должен иметь правильную комбинацию классов', () => {
            const wrapper = mount(BaseBadge, {
                props: {
                    variant: 'danger',
                    size: 'md'
                },
                slots: {
                    default: 'Badge'
                }
            })
            const classes = wrapper.classes()
            expect(classes).toContain('base-badge')
            expect(classes).toContain('base-badge--danger')
            expect(classes).toContain('base-badge--md')
        })
    })

    describe('Структура', () => {
        it('должен иметь правильную структуру без иконки', () => {
            const wrapper = mount(BaseBadge, {
                slots: {
                    default: 'Badge'
                }
            })
            expect(wrapper.find('.base-badge__text').exists()).toBe(true)
            expect(wrapper.find('.base-badge__icon').exists()).toBe(false)
        })

        it('должен иметь правильную структуру с иконкой', () => {
            const MockIcon = { name: 'MockIcon', template: '<svg><path /></svg>' }
            const wrapper = mount(BaseBadge, {
                props: {
                    icon: MockIcon
                },
                slots: {
                    default: 'Badge'
                }
            })
            const badge = wrapper.find('.base-badge')
            expect(badge.find('.base-badge__icon').exists()).toBe(true)
            expect(badge.find('.base-badge__text').exists()).toBe(true)
        })

        it('иконка должна быть перед текстом', () => {
            const MockIcon = { name: 'MockIcon', template: '<svg><path /></svg>' }
            const wrapper = mount(BaseBadge, {
                props: {
                    icon: MockIcon
                },
                slots: {
                    default: 'Badge'
                }
            })
            const badge = wrapper.find('.base-badge')
            const children = Array.from(badge.element.children)
            const iconIndex = children.findIndex(el => el.classList.contains('base-badge__icon'))
            const textIndex = children.findIndex(el => el.classList.contains('base-badge__text'))
            expect(iconIndex).toBeLessThan(textIndex)
        })
    })

    describe('TypeScript типизация', () => {
        it('должен принимать все валидные variants', () => {
            const variants: Array<'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info' | 'accent'> = [
                'primary', 'secondary', 'success', 'warning', 'danger', 'info', 'accent'
            ]

            variants.forEach(variant => {
                const wrapper = mount(BaseBadge, {
                    props: { variant },
                    slots: { default: 'Badge' }
                })
                expect(wrapper.find(`.base-badge--${variant}`).exists()).toBe(true)
            })
        })

        it('должен принимать все валидные sizes', () => {
            const sizes: Array<'xs' | 'sm' | 'md'> = ['xs', 'sm', 'md']

            sizes.forEach(size => {
                const wrapper = mount(BaseBadge, {
                    props: { size },
                    slots: { default: 'Badge' }
                })
                expect(wrapper.find(`.base-badge--${size}`).exists()).toBe(true)
            })
        })
    })
})
