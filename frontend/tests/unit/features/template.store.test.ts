import { useTemplateStore } from '@/features/templates/model/template.store'
import type { TemplateContext } from '@/features/templates/model/template.types'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('TemplateStore', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
    })

    describe('generatePrompt', () => {
        it('should return raw files content when no template is active', () => {
            const store = useTemplateStore()
            // Force no active template
            store.$patch({ templates: [] })

            const context: TemplateContext = {
                fileTree: 'src/\n  main.ts',
                files: '<file>content</file>',
                task: 'Test task',
                userRules: '',
                fileCount: 1,
                tokenCount: 100,
                languages: ['TypeScript'],
                projectName: 'test-project'
            }

            const result = store.generatePrompt(context)
            expect(result).toBe('<file>content</file>')
        })

        it('should include Role section when enabled', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')

            const context: TemplateContext = {
                fileTree: '',
                files: 'file content',
                task: '',
                userRules: '',
                fileCount: 0,
                tokenCount: 0,
                languages: [],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            expect(result).toContain('## Role')
            expect(result).toContain('senior software architect')
        })

        it('should include Rules section when enabled', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')

            const context: TemplateContext = {
                fileTree: '',
                files: 'file content',
                task: '',
                userRules: '',
                fileCount: 0,
                tokenCount: 0,
                languages: [],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            expect(result).toContain('## Rules')
            expect(result).toContain('SOLID principles')
        })

        it('should include user rules in Rules section', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')
            store.setUserRules('Custom user rule')

            const context: TemplateContext = {
                fileTree: '',
                files: 'file content',
                task: '',
                userRules: 'Custom user rule',
                fileCount: 0,
                tokenCount: 0,
                languages: [],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            expect(result).toContain('Custom user rule')
        })

        it('should include File Tree section when enabled', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')

            const context: TemplateContext = {
                fileTree: 'src/\n  main.ts\n  utils.ts',
                files: 'file content',
                task: '',
                userRules: '',
                fileCount: 2,
                tokenCount: 100,
                languages: ['TypeScript'],
                projectName: 'test-project'
            }

            const result = store.generatePrompt(context)
            expect(result).toContain('## Project Structure')
            expect(result).toContain('src/')
            expect(result).toContain('main.ts')
        })

        it('should include Stats section when enabled', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')

            const context: TemplateContext = {
                fileTree: '',
                files: 'file content',
                task: '',
                userRules: '',
                fileCount: 5,
                tokenCount: 1500,
                languages: ['TypeScript', 'Vue'],
                projectName: 'test-project'
            }

            const result = store.generatePrompt(context)
            expect(result).toContain('## Context Stats')
            expect(result).toContain('Files: 5')
            expect(result).toContain('Tokens: ~1500')
            expect(result).toContain('TypeScript, Vue')
        })

        it('should include Task section when task is set', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')
            store.setTask('Implement new feature')

            const context: TemplateContext = {
                fileTree: '',
                files: 'file content',
                task: 'Implement new feature',
                userRules: '',
                fileCount: 0,
                tokenCount: 0,
                languages: [],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            expect(result).toContain('## Task')
            expect(result).toContain('Implement new feature')
        })

        it('should include Files section with content', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')

            const filesContent = '<file path="src/main.ts">console.log("hello")</file>'
            const context: TemplateContext = {
                fileTree: '',
                files: filesContent,
                task: '',
                userRules: '',
                fileCount: 1,
                tokenCount: 50,
                languages: ['TypeScript'],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            expect(result).toContain('## Files')
            expect(result).toContain(filesContent)
        })

        it('should include custom prefix when set', () => {
            const store = useTemplateStore()
            const templateId = store.createTemplate({
                name: 'Test Template',
                icon: '🧪',
                description: 'Test',
                tags: [],
                isBuiltIn: false,
                isFavorite: false,
                isHidden: false,
                sections: { role: true, rules: false, tree: false, stats: false, task: false, files: true },
                sectionOrder: ['role', 'files'],
                roleContent: 'Test role',
                rulesContent: '',
                customPrefix: '--- START OF CONTEXT ---',
                customSuffix: ''
            })
            store.setActiveTemplate(templateId)

            const context: TemplateContext = {
                fileTree: '',
                files: 'file content',
                task: '',
                userRules: '',
                fileCount: 0,
                tokenCount: 0,
                languages: [],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            expect(result.startsWith('--- START OF CONTEXT ---')).toBe(true)
        })

        it('should include custom suffix when set', () => {
            const store = useTemplateStore()
            const templateId = store.createTemplate({
                name: 'Test Template',
                icon: '🧪',
                description: 'Test',
                tags: [],
                isBuiltIn: false,
                isFavorite: false,
                isHidden: false,
                sections: { role: true, rules: false, tree: false, stats: false, task: false, files: true },
                sectionOrder: ['role', 'files'],
                roleContent: 'Test role',
                rulesContent: '',
                customPrefix: '',
                customSuffix: '--- END OF CONTEXT ---'
            })
            store.setActiveTemplate(templateId)

            const context: TemplateContext = {
                fileTree: '',
                files: 'file content',
                task: '',
                userRules: '',
                fileCount: 0,
                tokenCount: 0,
                languages: [],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            expect(result.endsWith('--- END OF CONTEXT ---')).toBe(true)
        })

        it('should respect section order', () => {
            const store = useTemplateStore()
            const templateId = store.createTemplate({
                name: 'Custom Order',
                icon: '🔄',
                description: 'Test order',
                tags: [],
                isBuiltIn: false,
                isFavorite: false,
                isHidden: false,
                sections: { role: true, rules: true, tree: false, stats: false, task: false, files: true },
                sectionOrder: ['files', 'rules', 'role'],
                roleContent: 'Role content',
                rulesContent: 'Rules content',
                customPrefix: '',
                customSuffix: ''
            })
            store.setActiveTemplate(templateId)

            const context: TemplateContext = {
                fileTree: '',
                files: 'Files content here',
                task: '',
                userRules: '',
                fileCount: 1,
                tokenCount: 100,
                languages: [],
                projectName: ''
            }

            const result = store.generatePrompt(context)
            const filesIndex = result.indexOf('## Files')
            const rulesIndex = result.indexOf('## Rules')
            const roleIndex = result.indexOf('## Role')

            expect(filesIndex).toBeLessThan(rulesIndex)
            expect(rulesIndex).toBeLessThan(roleIndex)
        })

        it('should not include disabled sections', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('explain')

            const context: TemplateContext = {
                fileTree: 'tree',
                files: 'files',
                task: 'task',
                userRules: '',
                fileCount: 1,
                tokenCount: 100,
                languages: ['TypeScript'],
                projectName: 'test'
            }

            const result = store.generatePrompt(context)
            expect(result).not.toContain('## Rules')
            expect(result).not.toContain('## Context Stats')
        })

        it('should handle all sections enabled (architect template)', () => {
            const store = useTemplateStore()
            store.setActiveTemplate('architect')
            store.setTask('Analyze architecture')

            const context: TemplateContext = {
                fileTree: 'src/\n  index.ts',
                files: '<file>code</file>',
                task: 'Analyze architecture',
                userRules: '',
                fileCount: 3,
                tokenCount: 500,
                languages: ['TypeScript', 'Go'],
                projectName: 'my-project'
            }

            const result = store.generatePrompt(context)

            expect(result).toContain('## Role')
            expect(result).toContain('## Rules')
            expect(result).toContain('## Project Structure')
            expect(result).toContain('## Context Stats')
            expect(result).toContain('## Task')
            expect(result).toContain('## Files')
        })
    })
})
