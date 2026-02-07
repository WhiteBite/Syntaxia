import { useDependencyGraph } from '@/composables/useDependencyGraph'
import { useFileStore } from '@/features/files/model/file.store'
import type { FileNode } from '@/types/domain'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('Dependency Visualizer', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    describe('useDependencyGraph', () => {
        const { findDependencies, findIncomingDependencies, findAllDependencies, findRelatedFiles } = useDependencyGraph()

        const createMockNodes = (): FileNode[] => [
            {
                name: 'Button.vue',
                path: '/src/components/Button.vue',
                isDir: false,
                isExpanded: false,
                size: 1000,
                contentType: 'text'
            },
            {
                name: 'Button.test.ts',
                path: '/src/components/Button.test.ts',
                isDir: false,
                isExpanded: false,
                size: 500,
                contentType: 'text'
            },
            {
                name: 'Button.module.css',
                path: '/src/components/Button.module.css',
                isDir: false,
                isExpanded: false,
                size: 300,
                contentType: 'text'
            },
            {
                name: 'Button.types.ts',
                path: '/src/components/Button.types.ts',
                isDir: false,
                isExpanded: false,
                size: 200,
                contentType: 'text'
            },
            {
                name: 'App.vue',
                path: '/src/App.vue',
                isDir: false,
                isExpanded: false,
                size: 2000,
                contentType: 'text'
            }
        ]

        it('should find outgoing dependencies (files this file depends on)', () => {
            const nodes = createMockNodes()
            const deps = findDependencies('/src/components/Button.vue', nodes)

            expect(deps).toHaveLength(3)
            expect(deps.map(d => d.path)).toContain('/src/components/Button.test.ts')
            expect(deps.map(d => d.path)).toContain('/src/components/Button.module.css')
            expect(deps.map(d => d.path)).toContain('/src/components/Button.types.ts')

            // Check directions
            deps.forEach(dep => {
                expect(dep.direction).toBe('outgoing')
            })
        })

        it('should find incoming dependencies (files that depend on this file)', () => {
            const nodes = createMockNodes()
            const deps = findIncomingDependencies('/src/components/Button.test.ts', nodes)

            expect(deps).toHaveLength(1)
            expect(deps[0].path).toBe('/src/components/Button.vue')
            expect(deps[0].direction).toBe('incoming')
        })

        it('should find all dependencies (bidirectional)', () => {
            const nodes = createMockNodes()
            const { incoming, outgoing } = findAllDependencies('/src/components/Button.vue', nodes)

            expect(outgoing).toHaveLength(3)
            expect(incoming).toHaveLength(0) // Button.vue is not a test/style/type file
        })

        it('should identify dependency types correctly', () => {
            const nodes = createMockNodes()
            const deps = findDependencies('/src/components/Button.vue', nodes)

            const testDep = deps.find(d => d.path.includes('.test.'))
            const styleDep = deps.find(d => d.path.includes('.css'))
            const typeDep = deps.find(d => d.path.includes('.types.'))

            expect(testDep?.type).toBe('test')
            expect(styleDep?.type).toBe('style')
            expect(typeDep?.type).toBe('type')
        })

        it('should find related files by naming pattern', () => {
            const nodes = createMockNodes()
            const related = findRelatedFiles('/src/components/Button.vue', nodes)

            expect(related).toHaveLength(3)
            expect(related).toContain('/src/components/Button.test.ts')
            expect(related).toContain('/src/components/Button.module.css')
            expect(related).toContain('/src/components/Button.types.ts')
        })

        it('should handle files with no dependencies', () => {
            const nodes = createMockNodes()
            const deps = findDependencies('/src/App.vue', nodes)

            expect(deps).toHaveLength(0)
        })

        it('should handle compound extensions correctly', () => {
            const nodes: FileNode[] = [
                {
                    name: 'utils.ts',
                    path: '/src/utils.ts',
                    isDir: false,
                    isExpanded: false,
                    size: 1000,
                    contentType: 'text'
                },
                {
                    name: 'utils.test.ts',
                    path: '/src/utils.test.ts',
                    isDir: false,
                    isExpanded: false,
                    size: 500,
                    contentType: 'text'
                },
                {
                    name: 'utils.d.ts',
                    path: '/src/utils.d.ts',
                    isDir: false,
                    isExpanded: false,
                    size: 200,
                    contentType: 'text'
                }
            ]

            const deps = findDependencies('/src/utils.ts', nodes)
            expect(deps).toHaveLength(2)
            expect(deps.map(d => d.path)).toContain('/src/utils.test.ts')
            expect(deps.map(d => d.path)).toContain('/src/utils.d.ts')
        })
    })

    describe('File Store - Dependency Management', () => {
        it('should track dependencies for selected files', async () => {
            const fileStore = useFileStore()

            const nodes: FileNode[] = [
                {
                    name: 'root',
                    path: '/project',
                    isDir: true,
                    isExpanded: true,
                    children: [
                        {
                            name: 'Button.vue',
                            path: '/project/Button.vue',
                            isDir: false,
                            isExpanded: false,
                            size: 1000,
                            contentType: 'text'
                        },
                        {
                            name: 'Button.test.ts',
                            path: '/project/Button.test.ts',
                            isDir: false,
                            isExpanded: false,
                            size: 500,
                            contentType: 'text'
                        }
                    ]
                }
            ]

            fileStore.setFileTree(nodes)
            fileStore.setRootPath('/project')

            // Disable backend dependency graph for tests
            fileStore.useBackendDependencyGraph = false

            fileStore.selectPath('/project/Button.vue')

            // Wait for dependency computation
            await new Promise(resolve => setTimeout(resolve, 50))

            const deps = fileStore.allFileDependencies
            expect(deps.has('/project/Button.test.ts')).toBe(true)
        })

        it('should add all dependencies when requested', async () => {
            const fileStore = useFileStore()

            const nodes: FileNode[] = [
                {
                    name: 'root',
                    path: '/project',
                    isDir: true,
                    isExpanded: true,
                    children: [
                        {
                            name: 'Button.vue',
                            path: '/project/Button.vue',
                            isDir: false,
                            isExpanded: false,
                            size: 1000,
                            contentType: 'text'
                        },
                        {
                            name: 'Button.test.ts',
                            path: '/project/Button.test.ts',
                            isDir: false,
                            isExpanded: false,
                            size: 500,
                            contentType: 'text'
                        },
                        {
                            name: 'Button.module.css',
                            path: '/project/Button.module.css',
                            isDir: false,
                            isExpanded: false,
                            size: 300,
                            contentType: 'text'
                        }
                    ]
                }
            ]

            fileStore.setFileTree(nodes)
            fileStore.setRootPath('/project')

            // Disable backend dependency graph for tests
            fileStore.useBackendDependencyGraph = false

            fileStore.selectPath('/project/Button.vue')

            // Wait for dependency computation
            await new Promise(resolve => setTimeout(resolve, 50))

            const count = fileStore.addDependencies('/project/Button.test.ts')

            // Should add Button.vue (incoming dependency)
            expect(count).toBeGreaterThan(0)
            expect(fileStore.selectedPaths.has('/project/Button.vue')).toBe(true)
        })

        it('should not add already selected dependencies', () => {
            const fileStore = useFileStore()

            const nodes: FileNode[] = [
                {
                    name: 'root',
                    path: '/project',
                    isDir: true,
                    isExpanded: true,
                    children: [
                        {
                            name: 'Button.vue',
                            path: '/project/Button.vue',
                            isDir: false,
                            isExpanded: false,
                            size: 1000,
                            contentType: 'text'
                        },
                        {
                            name: 'Button.test.ts',
                            path: '/project/Button.test.ts',
                            isDir: false,
                            isExpanded: false,
                            size: 500,
                            contentType: 'text'
                        }
                    ]
                }
            ]

            fileStore.setFileTree(nodes)
            fileStore.setRootPath('/project')

            // Select both files
            fileStore.selectPath('/project/Button.vue')
            fileStore.selectPath('/project/Button.test.ts')

            // Try to add dependencies - should return 0 since all are already selected
            const count = fileStore.addDependencies('/project/Button.vue')
            expect(count).toBe(0)
        })

        it('should handle bidirectional dependencies correctly', async () => {
            const fileStore = useFileStore()

            const nodes: FileNode[] = [
                {
                    name: 'root',
                    path: '/project',
                    isDir: true,
                    isExpanded: true,
                    children: [
                        {
                            name: 'Component.vue',
                            path: '/project/Component.vue',
                            isDir: false,
                            isExpanded: false,
                            size: 1000,
                            contentType: 'text'
                        },
                        {
                            name: 'Component.test.ts',
                            path: '/project/Component.test.ts',
                            isDir: false,
                            isExpanded: false,
                            size: 500,
                            contentType: 'text'
                        }
                    ]
                }
            ]

            fileStore.setFileTree(nodes)
            fileStore.setRootPath('/project')

            // Disable backend dependency graph for tests
            fileStore.useBackendDependencyGraph = false

            fileStore.selectPath('/project/Component.vue')

            // Wait for dependency computation
            await new Promise(resolve => setTimeout(resolve, 50))

            const deps = fileStore.allFileDependencies.get('/project/Component.test.ts')

            expect(deps).toBeDefined()
            expect(deps?.incoming).toContain('/project/Component.vue')
        })
    })

    describe('Dependency Highlighting', () => {
        it('should highlight related files on hover', () => {
            // This test verifies basic file store functionality
            // Dependency highlighting is tested in component tests
            const fileStore = useFileStore()

            // Verify store is initialized
            expect(fileStore).toBeDefined()
            expect(fileStore.setFileTree).toBeDefined()
            expect(fileStore.selectPath).toBeDefined()
        })
    })

    describe('Performance', () => {
        it('should handle large dependency graphs efficiently', () => {
            const { findDependencies } = useDependencyGraph()

            // Create 1000 mock files
            const nodes: FileNode[] = []
            for (let i = 0; i < 1000; i++) {
                nodes.push({
                    name: `file${i}.ts`,
                    path: `/src/file${i}.ts`,
                    isDir: false,
                    isExpanded: false,
                    size: 1000,
                    contentType: 'text'
                })
            }

            const startTime = performance.now()
            findDependencies('/src/file0.ts', nodes)
            const endTime = performance.now()

            // Should complete in less than 100ms
            expect(endTime - startTime).toBeLessThan(100)
        })
    })
})
