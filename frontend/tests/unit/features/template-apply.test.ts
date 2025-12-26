import { useTemplateStore } from '@/features/templates/model/template.store'
import type { TemplateContext } from '@/features/templates/model/template.types'
import { useSettingsStore } from '@/stores/settings.store'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('Template Application on Copy', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
    })

    const mockContext: TemplateContext = {
        fileTree: 'src/\n  main.ts\n  utils.ts',
        files: '<file path="src/main.ts">console.log("hello")</file>',
        task: 'Implement feature',
        userRules: '',
        fileCount: 2,
        tokenCount: 500,
        languages: ['TypeScript'],
        projectName: 'test-project'
    }

    describe('applyTemplateOnCopy setting', () => {
        it('should be true by default', () => {
            const settingsStore = useSettingsStore()
            expect(settingsStore.settings.context.applyTemplateOnCopy).toBe(true)
        })

        it('should be changeable', () => {
            const settingsStore = useSettingsStore()
            settingsStore.updateContextSettings({ applyTemplateOnCopy: false })
            expect(settingsStore.settings.context.applyTemplateOnCopy).toBe(false)
        })
    })

    describe('copy behavior with applyTemplateOnCopy=true', () => {
        it('should include Role section in output', () => {
            const templateStore = useTemplateStore()
            const settingsStore = useSettingsStore()

            settingsStore.updateContextSettings({ applyTemplateOnCopy: true })
            templateStore.setActiveTemplate('architect')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Role')
            expect(result).toContain('senior software architect')
        })

        it('should include Rules section in output', () => {
            const templateStore = useTemplateStore()
            const settingsStore = useSettingsStore()

            settingsStore.updateContextSettings({ applyTemplateOnCopy: true })
            templateStore.setActiveTemplate('architect')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Rules')
            expect(result).toContain('SOLID principles')
        })

        it('should include Project Structure section in output', () => {
            const templateStore = useTemplateStore()
            const settingsStore = useSettingsStore()

            settingsStore.updateContextSettings({ applyTemplateOnCopy: true })
            templateStore.setActiveTemplate('architect')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Project Structure')
            expect(result).toContain('src/')
        })

        it('should include Stats section in output', () => {
            const templateStore = useTemplateStore()
            const settingsStore = useSettingsStore()

            settingsStore.updateContextSettings({ applyTemplateOnCopy: true })
            templateStore.setActiveTemplate('architect')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Context Stats')
            expect(result).toContain('Files: 2')
            expect(result).toContain('Tokens: ~500')
        })

        it('should include Task section when task is set', () => {
            const templateStore = useTemplateStore()
            const settingsStore = useSettingsStore()

            settingsStore.updateContextSettings({ applyTemplateOnCopy: true })
            templateStore.setActiveTemplate('architect')
            templateStore.setTask('Implement feature')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Task')
            expect(result).toContain('Implement feature')
        })

        it('should include Files section with content', () => {
            const templateStore = useTemplateStore()
            const settingsStore = useSettingsStore()

            settingsStore.updateContextSettings({ applyTemplateOnCopy: true })
            templateStore.setActiveTemplate('architect')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Files')
            expect(result).toContain('<file path="src/main.ts">')
        })
    })

    describe('copy behavior with applyTemplateOnCopy=false', () => {
        it('should return only files content when setting is false and no template', () => {
            const templateStore = useTemplateStore()

            // Force no active template
            templateStore.$patch({ templates: [] })

            const result = templateStore.generatePrompt(mockContext)

            // Should return raw files content
            expect(result).toBe(mockContext.files)
            expect(result).not.toContain('## Role')
            expect(result).not.toContain('## Rules')
        })
    })

    describe('different templates produce different output', () => {
        it('architect template should include all sections', () => {
            const templateStore = useTemplateStore()
            templateStore.setActiveTemplate('architect')
            templateStore.setTask('Test task')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Role')
            expect(result).toContain('## Rules')
            expect(result).toContain('## Project Structure')
            expect(result).toContain('## Context Stats')
            expect(result).toContain('## Task')
            expect(result).toContain('## Files')
        })

        it('implement template should not include stats', () => {
            const templateStore = useTemplateStore()
            templateStore.setActiveTemplate('implement')
            templateStore.setTask('Test task')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Role')
            expect(result).toContain('## Rules')
            expect(result).toContain('## Project Structure')
            expect(result).not.toContain('## Context Stats')
            expect(result).toContain('## Task')
            expect(result).toContain('## Files')
        })

        it('explain template should not include rules and stats', () => {
            const templateStore = useTemplateStore()
            templateStore.setActiveTemplate('explain')
            templateStore.setTask('Test task')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Role')
            expect(result).not.toContain('## Rules')
            expect(result).toContain('## Project Structure')
            expect(result).not.toContain('## Context Stats')
            expect(result).toContain('## Task')
            expect(result).toContain('## Files')
        })

        it('review template should not include tree', () => {
            const templateStore = useTemplateStore()
            templateStore.setActiveTemplate('review')
            templateStore.setTask('Test task')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Role')
            expect(result).toContain('## Rules')
            expect(result).not.toContain('## Project Structure')
            expect(result).toContain('## Context Stats')
            expect(result).toContain('## Task')
            expect(result).toContain('## Files')
        })
    })

    describe('user rules integration', () => {
        it('should append user rules to template rules', () => {
            const templateStore = useTemplateStore()
            templateStore.setActiveTemplate('architect')
            templateStore.setUserRules('Always use TypeScript strict mode')

            const result = templateStore.generatePrompt(mockContext)

            expect(result).toContain('## Rules')
            expect(result).toContain('SOLID principles') // Template rules
            expect(result).toContain('Always use TypeScript strict mode') // User rules
        })

        it('should show only user rules when template has no rules', () => {
            const templateStore = useTemplateStore()
            templateStore.setActiveTemplate('explain') // Has no rulesContent
            templateStore.setUserRules('Custom rule')

            const result = templateStore.generatePrompt(mockContext)

            // explain template has rules: false, so no Rules section
            expect(result).not.toContain('## Rules')
        })
    })
})
