import { useFileStore } from '@/features/files/model/file.store'
import DependencyVisualizerModal from '@/features/files/ui/DependencyVisualizerModal.vue'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock lucide-vue-next icons
vi.mock('lucide-vue-next', () => ({
    ArrowRightCircle: { name: 'ArrowRightCircle' },
    ArrowLeftCircle: { name: 'ArrowLeftCircle' },
    Plus: { name: 'Plus' },
    Check: { name: 'Check' },
    PlusCircle: { name: 'PlusCircle' },
    Eye: { name: 'Eye' },
    FileQuestion: { name: 'FileQuestion' },
    FileCode: { name: 'FileCode' },
    FileText: { name: 'FileText' }
}))

// Mock composables
vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string, params?: Record<string, any>) => {
            if (params) {
                return `${key}:${JSON.stringify(params)}`
            }
            return key
        }
    })
}))

vi.mock('@/composables/useDependencyGraph', () => ({
    useDependencyGraph: () => ({
        findAllDependencies: vi.fn((filePath: string) => {
            if (filePath === 'src/components/Button.vue') {
                return {
                    incoming: [
                        { path: 'src/pages/Home.vue', type: 'import', confidence: 'high', direction: 'incoming' }
                    ],
                    outgoing: [
                        { path: 'src/components/Button.css', type: 'style', confidence: 'high', direction: 'outgoing' },
                        { path: 'src/types/button.ts', type: 'type', confidence: 'high', direction: 'outgoing' }
                    ]
                }
            }
            return { incoming: [], outgoing: [] }
        })
    })
}))

describe('DependencyVisualizerModal', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('renders modal with dependencies', () => {
        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        expect(wrapper.find('.dep-visualizer').exists()).toBe(true)
    })

    it('displays outgoing dependencies', () => {
        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        const depItems = wrapper.findAll('.dep-item')
        expect(depItems.length).toBeGreaterThan(0)
    })

    it('displays incoming dependencies', () => {
        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        const sections = wrapper.findAll('.dep-section')
        expect(sections.length).toBeGreaterThan(0)
    })

    it('shows statistics correctly', () => {
        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        expect(wrapper.find('.dep-stats').exists()).toBe(true)
        expect(wrapper.findAll('.stat-item').length).toBeGreaterThan(0)
    })

    it('shows empty state when no dependencies', () => {
        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Empty.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        expect(wrapper.find('.dep-empty').exists()).toBe(true)
    })

    it('can add missing dependencies', async () => {
        const fileStore = useFileStore()

        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        const addAllButton = wrapper.find('.btn-primary')
        expect(addAllButton.exists()).toBe(true)

        await addAllButton.trigger('click')

        // Verify files were added to selection
        expect(fileStore.selectedPaths.has('src/components/Button.css')).toBe(true)
        expect(fileStore.selectedPaths.has('src/types/button.ts')).toBe(true)
        expect(fileStore.selectedPaths.has('src/pages/Home.vue')).toBe(true)
    })

    it('disables add all button when no missing dependencies', async () => {
        const fileStore = useFileStore()

        // Pre-select all dependencies
        fileStore.selectPath('src/components/Button.css')
        fileStore.selectPath('src/types/button.ts')
        fileStore.selectPath('src/pages/Home.vue')

        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        const addAllButton = wrapper.find('.btn-primary')
        expect(addAllButton.attributes('disabled')).toBeDefined()
    })

    it('can navigate to file in tree', async () => {
        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>',
                        props: ['modelValue'],
                        emits: ['update:modelValue', 'close']
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        // Call the method directly
        await wrapper.vm.handleNavigateToFile()

        // Verify modal was closed
        expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([false])
    })

    it('closes modal on close button click', async () => {
        const wrapper = mount(DependencyVisualizerModal, {
            props: {
                filePath: 'src/components/Button.vue',
                modelValue: true
            },
            global: {
                stubs: {
                    BaseModal: {
                        template: '<div class="base-modal"><slot /></div>'
                    },
                    BaseIcon: {
                        template: '<span class="base-icon" />'
                    },
                    BaseBadge: {
                        template: '<span class="base-badge"><slot /></span>'
                    }
                }
            }
        })

        await wrapper.vm.handleClose()

        expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([false])
    })
})
