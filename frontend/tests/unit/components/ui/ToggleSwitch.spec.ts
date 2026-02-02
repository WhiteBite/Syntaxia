import ToggleSwitch from '@/components/ui/ToggleSwitch.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

describe('ToggleSwitch.vue', () => {
    describe('Рендеринг', () => {
        it('должен рендерить toggle switch', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.find('.toggle-switch').exists()).toBe(true)
        })

        it('должен рендерить track', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.find('.toggle-switch__track').exists()).toBe(true)
        })

        it('должен рендерить thumb', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.find('.toggle-switch__thumb').exists()).toBe(true)
        })

        it('должен быть button элементом', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.element.tagName).toBe('BUTTON')
        })

        it('должен иметь type="button"', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.attributes('type')).toBe('button')
        })

        it('должен иметь role="switch"', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.attributes('role')).toBe('switch')
        })
    })

    describe('Props', () => {
        it('должен применять класс active когда modelValue=true', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: true
                }
            })
            expect(wrapper.find('.toggle-switch--active').exists()).toBe(true)
        })

        it('не должен применять класс active когда modelValue=false', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.find('.toggle-switch--active').exists()).toBe(false)
        })

        it('должен применять класс disabled когда disabled=true', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: true
                }
            })
            expect(wrapper.find('.toggle-switch--disabled').exists()).toBe(true)
        })

        it('не должен применять класс disabled по умолчанию', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.find('.toggle-switch--disabled').exists()).toBe(false)
        })

        it('должен иметь атрибут disabled когда disabled=true', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: true
                }
            })
            expect(wrapper.attributes('disabled')).toBeDefined()
        })

        it('не должен иметь атрибут disabled когда disabled=false', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: false
                }
            })
            expect(wrapper.attributes('disabled')).toBeUndefined()
        })
    })

    describe('ARIA атрибуты', () => {
        it('должен иметь aria-checked="true" когда modelValue=true', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: true
                }
            })
            expect(wrapper.attributes('aria-checked')).toBe('true')
        })

        it('должен иметь aria-checked="false" когда modelValue=false', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.attributes('aria-checked')).toBe('false')
        })

        it('должен обновлять aria-checked при изменении modelValue', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.attributes('aria-checked')).toBe('false')

            await wrapper.setProps({ modelValue: true })
            expect(wrapper.attributes('aria-checked')).toBe('true')
        })
    })

    describe('Events', () => {
        it('должен эмитить update:modelValue при клике', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([true])
        })

        it('должен эмитить change при клике', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            await wrapper.trigger('click')
            expect(wrapper.emitted('change')).toBeTruthy()
            expect(wrapper.emitted('change')?.[0]).toEqual([true])
        })

        it('должен переключать значение с false на true', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([true])
        })

        it('должен переключать значение с true на false', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: true
                }
            })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([false])
        })

        it('не должен эмитить события когда disabled', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: true
                }
            })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')).toBeFalsy()
            expect(wrapper.emitted('change')).toBeFalsy()
        })

        it('должен эмитить оба события одновременно', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
            expect(wrapper.emitted('change')).toBeTruthy()
        })

        it('должен эмитить события с правильным значением при множественных кликах', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })

            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([true])

            await wrapper.setProps({ modelValue: true })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')?.[1]).toEqual([false])
        })
    })

    describe('Визуальные состояния', () => {
        it('должен иметь класс active и disabled одновременно', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: true,
                    disabled: true
                }
            })
            expect(wrapper.find('.toggle-switch--active').exists()).toBe(true)
            expect(wrapper.find('.toggle-switch--disabled').exists()).toBe(true)
        })

        it('должен корректно обновлять классы при изменении props', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: false
                }
            })

            expect(wrapper.find('.toggle-switch--active').exists()).toBe(false)
            expect(wrapper.find('.toggle-switch--disabled').exists()).toBe(false)

            await wrapper.setProps({ modelValue: true, disabled: true })

            expect(wrapper.find('.toggle-switch--active').exists()).toBe(true)
            expect(wrapper.find('.toggle-switch--disabled').exists()).toBe(true)
        })
    })

    describe('Edge cases', () => {
        it('должен корректно работать с быстрыми множественными кликами', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })

            // Simulate rapid clicks with prop updates
            await wrapper.trigger('click')
            await wrapper.setProps({ modelValue: true })

            await wrapper.trigger('click')
            await wrapper.setProps({ modelValue: false })

            await wrapper.trigger('click')
            await wrapper.setProps({ modelValue: true })

            expect(wrapper.emitted('update:modelValue')).toHaveLength(3)
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([true])
            expect(wrapper.emitted('update:modelValue')?.[1]).toEqual([false])
            expect(wrapper.emitted('update:modelValue')?.[2]).toEqual([true])
        })

        it('должен игнорировать клики когда disabled', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: true
                }
            })

            await wrapper.trigger('click')
            await wrapper.trigger('click')
            await wrapper.trigger('click')

            expect(wrapper.emitted('update:modelValue')).toBeFalsy()
        })

        it('должен корректно работать при изменении disabled во время работы', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: false
                }
            })

            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')).toHaveLength(1)

            await wrapper.setProps({ disabled: true })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')).toHaveLength(1)

            await wrapper.setProps({ disabled: false })
            await wrapper.trigger('click')
            expect(wrapper.emitted('update:modelValue')).toHaveLength(2)
        })
    })

    describe('CSS классы', () => {
        it('должен всегда иметь базовый класс', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.classes()).toContain('toggle-switch')
        })

        it('должен иметь только базовый класс без props', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            const classes = wrapper.classes()
            expect(classes).toContain('toggle-switch')
            expect(classes).not.toContain('toggle-switch--active')
            expect(classes).not.toContain('toggle-switch--disabled')
        })

        it('должен корректно комбинировать все классы', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: true,
                    disabled: true
                }
            })
            const classes = wrapper.classes()
            expect(classes).toContain('toggle-switch')
            expect(classes).toContain('toggle-switch--active')
            expect(classes).toContain('toggle-switch--disabled')
        })
    })

    describe('Структура DOM', () => {
        it('должен иметь правильную вложенность элементов', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })

            const button = wrapper.find('.toggle-switch')
            expect(button.exists()).toBe(true)

            const track = button.find('.toggle-switch__track')
            expect(track.exists()).toBe(true)

            const thumb = track.find('.toggle-switch__thumb')
            expect(thumb.exists()).toBe(true)
        })

        it('thumb должен быть внутри track', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })

            const track = wrapper.find('.toggle-switch__track')
            const thumb = track.find('.toggle-switch__thumb')

            expect(thumb.exists()).toBe(true)
        })
    })

    describe('Accessibility', () => {
        it('должен быть доступен для клавиатуры (button)', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.element.tagName).toBe('BUTTON')
        })

        it('должен иметь правильную ARIA роль', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })
            expect(wrapper.attributes('role')).toBe('switch')
        })

        it('должен корректно обновлять ARIA состояние', async () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false
                }
            })

            expect(wrapper.attributes('aria-checked')).toBe('false')

            await wrapper.setProps({ modelValue: true })
            expect(wrapper.attributes('aria-checked')).toBe('true')

            await wrapper.setProps({ modelValue: false })
            expect(wrapper.attributes('aria-checked')).toBe('false')
        })

        it('должен быть недоступен когда disabled', () => {
            const wrapper = mount(ToggleSwitch, {
                props: {
                    modelValue: false,
                    disabled: true
                }
            })
            expect(wrapper.attributes('disabled')).toBeDefined()
        })
    })
})
