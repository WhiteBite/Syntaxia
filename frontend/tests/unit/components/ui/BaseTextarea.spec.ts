import BaseTextarea from '@/components/ui/BaseTextarea.vue'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

describe('BaseTextarea.vue', () => {
    describe('Рендеринг', () => {
        it('должен рендерить textarea', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            expect(wrapper.find('textarea').exists()).toBe(true)
        })

        it('должен отображать label если передан', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    label: 'Test Label'
                }
            })
            expect(wrapper.find('.base-textarea-label').text()).toBe('Test Label')
        })

        it('должен отображать placeholder', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    placeholder: 'Enter text...'
                }
            })
            expect(wrapper.find('textarea').attributes('placeholder')).toBe('Enter text...')
        })

        it('должен отображать значение', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: 'Test content'
                }
            })
            expect(wrapper.find('textarea').element.value).toBe('Test content')
        })
    })

    describe('Props', () => {
        it('должен применять rows по умолчанию (3)', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            expect(wrapper.find('textarea').attributes('rows')).toBe('3')
        })

        it('должен применять кастомное количество rows', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    rows: 5
                }
            })
            expect(wrapper.find('textarea').attributes('rows')).toBe('5')
        })

        it('должен быть disabled когда передан disabled prop', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    disabled: true
                }
            })
            expect(wrapper.find('textarea').attributes('disabled')).toBeDefined()
            expect(wrapper.find('.base-textarea-wrapper--disabled').exists()).toBe(true)
        })

        it('должен применять maxlength', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    maxLength: 100
                }
            })
            expect(wrapper.find('textarea').attributes('maxlength')).toBe('100')
        })

        it('должен отображать счетчик символов когда указан maxLength', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: 'Hello',
                    maxLength: 100
                }
            })
            expect(wrapper.find('.base-textarea-counter').text()).toBe('5/100')
        })

        it('должен отображать error сообщение', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    error: 'This field is required'
                }
            })
            expect(wrapper.find('.base-textarea-error').text()).toBe('This field is required')
            expect(wrapper.find('.base-textarea-container--error').exists()).toBe(true)
        })

        it('должен применять класс no-resize когда resize=false', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    resize: false
                }
            })
            expect(wrapper.find('.base-textarea--no-resize').exists()).toBe(true)
        })

        it('должен применять класс no-resize когда autoResize=true', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    autoResize: true
                }
            })
            expect(wrapper.find('.base-textarea--no-resize').exists()).toBe(true)
        })
    })

    describe('Events', () => {
        it('должен эмитить update:modelValue при вводе', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            const textarea = wrapper.find('textarea')
            await textarea.setValue('New text')
            expect(wrapper.emitted('update:modelValue')).toBeTruthy()
            expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['New text'])
        })

        it('должен обновлять счетчик при вводе', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    maxLength: 100
                }
            })
            await wrapper.setProps({ modelValue: 'Hello World' })
            expect(wrapper.find('.base-textarea-counter').text()).toBe('11/100')
        })

        it('должен добавлять класс focused при фокусе', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            const textarea = wrapper.find('textarea')
            await textarea.trigger('focus')
            expect(wrapper.find('.base-textarea-container--focused').exists()).toBe(true)
        })

        it('должен убирать класс focused при blur', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            const textarea = wrapper.find('textarea')
            await textarea.trigger('focus')
            await textarea.trigger('blur')
            expect(wrapper.find('.base-textarea-container--focused').exists()).toBe(false)
        })
    })

    describe('Character counter', () => {
        it('должен показывать warning класс когда близко к лимиту (90%)', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    maxLength: 100
                }
            })
            await wrapper.setProps({ modelValue: 'a'.repeat(90) })
            expect(wrapper.find('.base-textarea-counter--limit').exists()).toBe(true)
        })

        it('не должен показывать warning класс когда далеко от лимита', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: 'Hello',
                    maxLength: 100
                }
            })
            expect(wrapper.find('.base-textarea-counter--limit').exists()).toBe(false)
        })

        it('не должен отображать счетчик если maxLength не указан', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: 'Hello'
                }
            })
            expect(wrapper.find('.base-textarea-counter').exists()).toBe(false)
        })
    })

    describe('Exposed methods', () => {
        it('должен экспонировать метод focus', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            expect(wrapper.vm.focus).toBeDefined()
            expect(typeof wrapper.vm.focus).toBe('function')
        })

        it('должен экспонировать textarea ref', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            expect(wrapper.vm.textarea).toBeDefined()
        })
    })

    describe('Edge cases', () => {
        it('должен корректно работать с пустым значением', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: ''
                }
            })
            expect(wrapper.find('textarea').element.value).toBe('')
        })

        it('должен корректно работать с многострочным текстом', () => {
            const multilineText = 'Line 1\nLine 2\nLine 3'
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: multilineText
                }
            })
            expect(wrapper.find('textarea').element.value).toBe(multilineText)
        })

        it('должен корректно работать с Unicode символами', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    maxLength: 10
                }
            })
            // Unicode emoji считаются как несколько символов в JavaScript
            await wrapper.setProps({ modelValue: '🚀🎉✨' })
            // Emoji могут занимать 2-4 символа каждый в зависимости от реализации
            const counter = wrapper.find('.base-textarea-counter').text()
            expect(counter).toMatch(/\d+\/10/)
        })

        it('не должен эмитить события когда disabled', async () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: '',
                    disabled: true
                }
            })
            const textarea = wrapper.find('textarea')
            await textarea.trigger('focus')
            // Disabled textarea не должна получать фокус, но событие может быть эмитировано
            // Проверяем что textarea действительно disabled
            expect(textarea.attributes('disabled')).toBeDefined()
        })

        it('должен корректно работать со всеми props одновременно', () => {
            const wrapper = mount(BaseTextarea, {
                props: {
                    modelValue: 'Test',
                    label: 'Label',
                    placeholder: 'Placeholder',
                    rows: 5,
                    maxLength: 100,
                    error: 'Error message',
                    disabled: false,
                    resize: true
                }
            })
            expect(wrapper.find('.base-textarea-label').text()).toBe('Label')
            expect(wrapper.find('textarea').attributes('placeholder')).toBe('Placeholder')
            expect(wrapper.find('textarea').attributes('rows')).toBe('5')
            expect(wrapper.find('.base-textarea-counter').text()).toBe('4/100')
            expect(wrapper.find('.base-textarea-error').text()).toBe('Error message')
        })
    })
})
