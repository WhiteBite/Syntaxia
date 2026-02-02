import ConfirmHeavyFileModal from '@/features/files/ui/ConfirmHeavyFileModal.vue'
import { useSettingsStore } from '@/stores/settings.store'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'

describe('Size Guard - Heavy File Warnings', () => {
    let appDiv: HTMLElement

    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        // Create a div for teleport target
        appDiv = document.createElement('div')
        appDiv.id = 'app'
        document.body.appendChild(appDiv)
    })

    afterEach(() => {
        document.body.innerHTML = ''
    })

    describe('ConfirmHeavyFileModal', () => {
        it('should render modal when open', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 150000
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            expect(document.querySelector('.base-modal-backdrop')).toBeTruthy()
            expect(document.querySelector('.modal-title')).toBeTruthy()
            expect(document.body.textContent).toContain('large-file.ts')
            expect(document.body.textContent).toContain('150k')

            wrapper.unmount()
        })

        it('should not render modal when closed', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: false,
                    fileName: 'large-file.ts',
                    tokens: 150000
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            expect(document.querySelector('.base-modal-backdrop')).toBeFalsy()

            wrapper.unmount()
        })

        it('should calculate percentage correctly', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 100000 // 50% of 200k limit
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            expect(document.body.textContent).toContain('50%')

            wrapper.unmount()
        })

        it('should emit confirm event when Select Anyway is clicked', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 150000
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            const button = document.querySelector('.btn-warning') as HTMLElement
            button?.click()
            await wrapper.vm.$nextTick()

            expect(wrapper.emitted('confirm')).toBeTruthy()
            expect(wrapper.emitted('confirm')?.[0]).toEqual([false])

            wrapper.unmount()
        })

        it('should emit cancel event when Cancel is clicked', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 150000
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            const button = document.querySelector('.btn-ghost') as HTMLElement
            button?.click()
            await wrapper.vm.$nextTick()

            expect(wrapper.emitted('cancel')).toBeTruthy()

            wrapper.unmount()
        })

        it('should pass dontShowAgain flag when checkbox is checked', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 150000
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            const checkbox = document.querySelector('.checkbox-input') as HTMLInputElement
            if (checkbox) {
                checkbox.click()
                await wrapper.vm.$nextTick()
            }

            const button = document.querySelector('.btn-warning') as HTMLElement
            button?.click()
            await wrapper.vm.$nextTick()

            expect(wrapper.emitted('confirm')?.[0]).toEqual([true])

            wrapper.unmount()
        })

        it('should apply correct size class for critical files', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 150000 // Critical
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            expect(document.querySelector('.size-critical')).toBeTruthy()

            wrapper.unmount()
        })

        it('should apply correct size class for heavy files', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 75000 // Heavy
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            expect(document.querySelector('.size-heavy')).toBeTruthy()

            wrapper.unmount()
        })

        it('should format tokens correctly', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 1500
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            // 1500 tokens should be formatted as "1.5k"
            expect(document.body.textContent).toContain('1.5k')

            wrapper.unmount()
        })

        it('should format large token counts with k suffix', async () => {
            const wrapper = mount(ConfirmHeavyFileModal, {
                props: {
                    isOpen: true,
                    fileName: 'large-file.ts',
                    tokens: 150000
                },
                attachTo: appDiv
            })

            await wrapper.vm.$nextTick()

            expect(document.body.textContent).toContain('150k')

            wrapper.unmount()
        })
    })

    describe('Settings Integration', () => {
        it('should have warnHeavyFiles setting enabled by default', () => {
            const settingsStore = useSettingsStore()

            expect(settingsStore.settings.fileExplorer.warnHeavyFiles).toBe(true)
        })

        it('should have default threshold of 100k tokens', () => {
            const settingsStore = useSettingsStore()

            expect(settingsStore.settings.fileExplorer.heavyFileThreshold).toBe(100000)
        })

        it('should allow updating warnHeavyFiles setting', () => {
            const settingsStore = useSettingsStore()

            settingsStore.updateFileExplorerSettings({ warnHeavyFiles: false })

            expect(settingsStore.settings.fileExplorer.warnHeavyFiles).toBe(false)
        })

        it('should allow updating heavyFileThreshold setting', () => {
            const settingsStore = useSettingsStore()

            settingsStore.updateFileExplorerSettings({ heavyFileThreshold: 50000 })

            expect(settingsStore.settings.fileExplorer.heavyFileThreshold).toBe(50000)
        })
    })
})
