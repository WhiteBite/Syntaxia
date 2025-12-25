import type { FileNode } from '@/features/files/model/file.store'
import type { DomainNode } from '@/utils/fileTreeUtils'
import {
    convertDomainNodes,
    filterTreeByExtensions,
    findNode,
    walkTree,
} from '@/utils/fileTreeUtils'
import { describe, expect, it } from 'vitest'

describe('fileTreeUtils', () => {
    // Mock data
    const mockTree: FileNode[] = [
        {
            name: 'src',
            path: '/project/src',
            isDir: true,
            children: [
                { name: 'index.ts', path: '/project/src/index.ts', isDir: false },
                { name: 'utils.ts', path: '/project/src/utils.ts', isDir: false },
                { name: 'styles.css', path: '/project/src/styles.css', isDir: false },
                {
                    name: 'components',
                    path: '/project/src/components',
                    isDir: true,
                    children: [
                        { name: 'Button.vue', path: '/project/src/components/Button.vue', isDir: false },
                        { name: 'Input.vue', path: '/project/src/components/Input.vue', isDir: false },
                    ],
                },
            ],
        },
        { name: 'README.md', path: '/project/README.md', isDir: false },
    ]

    describe('findNode', () => {
        it('находит ноду по пути', () => {
            const node = findNode(mockTree, '/project/src/index.ts')

            expect(node).not.toBeNull()
            expect(node?.name).toBe('index.ts')
            expect(node?.path).toBe('/project/src/index.ts')
        })

        it('возвращает null для несуществующего пути', () => {
            const node = findNode(mockTree, '/project/nonexistent.ts')

            expect(node).toBeNull()
        })

        it('находит вложенную ноду', () => {
            const node = findNode(mockTree, '/project/src/components/Button.vue')

            expect(node).not.toBeNull()
            expect(node?.name).toBe('Button.vue')
        })

        it('использует кэш для повторных поисков', () => {
            const cache = new Map<string, FileNode>()

            // First search - populates cache
            const node1 = findNode(mockTree, '/project/src/index.ts', cache)
            expect(cache.size).toBe(1)

            // Second search - uses cache
            const node2 = findNode(mockTree, '/project/src/index.ts', cache)
            expect(node1).toBe(node2)
        })
    })

    describe('walkTree', () => {
        it('обходит все ноды', () => {
            const visited: string[] = []

            walkTree(mockTree, (node) => {
                visited.push(node.path)
            })

            expect(visited).toContain('/project/src')
            expect(visited).toContain('/project/src/index.ts')
            expect(visited).toContain('/project/src/components')
            expect(visited).toContain('/project/src/components/Button.vue')
            expect(visited).toContain('/project/README.md')
            expect(visited.length).toBe(8) // All nodes
        })

        it('обходит ноды в правильном порядке (depth-first)', () => {
            const visited: string[] = []

            walkTree(mockTree, (node) => {
                visited.push(node.name)
            })

            // src should come before README.md
            expect(visited.indexOf('src')).toBeLessThan(visited.indexOf('README.md'))
            // index.ts should come after src
            expect(visited.indexOf('index.ts')).toBeGreaterThan(visited.indexOf('src'))
        })
    })

    describe('filterTreeByExtensions', () => {
        it('фильтрует по расширению', () => {
            const filtered = filterTreeByExtensions(mockTree, ['.ts'])

            // Should have src folder with only .ts files
            expect(filtered.length).toBe(1) // Only src folder
            expect(filtered[0].name).toBe('src')

            // Check children - should only have .ts files
            const srcChildren = filtered[0].children || []
            const fileNames = srcChildren.map((c) => c.name)
            expect(fileNames).toContain('index.ts')
            expect(fileNames).toContain('utils.ts')
            expect(fileNames).not.toContain('styles.css')
        })

        it('исключает расширения из exclude списка', () => {
            const filtered = filterTreeByExtensions(mockTree, [], ['.css'])

            // Should have all files except .css
            const srcFolder = filtered.find((n) => n.name === 'src')
            expect(srcFolder).toBeDefined()

            const srcChildren = srcFolder?.children || []
            const fileNames = srcChildren.map((c) => c.name)
            expect(fileNames).toContain('index.ts')
            expect(fileNames).not.toContain('styles.css')
        })

        it('возвращает пустой массив если нет совпадений', () => {
            const filtered = filterTreeByExtensions(mockTree, ['.xyz'])

            expect(filtered.length).toBe(0)
        })

        it('сохраняет структуру папок', () => {
            const filtered = filterTreeByExtensions(mockTree, ['.vue'])

            // Should have src > components > .vue files
            expect(filtered.length).toBe(1)
            expect(filtered[0].name).toBe('src')

            const components = filtered[0].children?.find((c) => c.name === 'components')
            expect(components).toBeDefined()
            expect(components?.children?.length).toBe(2)
        })
    })

    describe('convertDomainNodes', () => {
        it('конвертирует DomainNode в FileNode', () => {
            const domainNodes: DomainNode[] = [
                {
                    name: 'test.ts',
                    path: '/test.ts',
                    isDir: false,
                    size: 1024,
                    isIgnored: false,
                },
                {
                    name: 'folder',
                    path: '/folder',
                    isDir: true,
                    children: [
                        { name: 'nested.ts', path: '/folder/nested.ts', isDir: false },
                    ],
                },
            ]

            const result = convertDomainNodes(domainNodes)

            expect(result.length).toBe(2)
            expect(result[0].name).toBe('test.ts')
            expect(result[0].isExpanded).toBe(false)
            expect(result[0].isSelected).toBe(false)
            expect(result[0].size).toBe(1024)

            expect(result[1].children?.length).toBe(1)
            expect(result[1].children?.[0].name).toBe('nested.ts')
        })

        it('обрабатывает isGitignored и isCustomIgnored', () => {
            const domainNodes: DomainNode[] = [
                {
                    name: 'ignored.ts',
                    path: '/ignored.ts',
                    isDir: false,
                    isGitignored: true,
                },
                {
                    name: 'custom.ts',
                    path: '/custom.ts',
                    isDir: false,
                    isCustomIgnored: true,
                },
            ]

            const result = convertDomainNodes(domainNodes)

            expect(result[0].isIgnored).toBe(true)
            expect(result[1].isIgnored).toBe(true)
        })
    })
})
