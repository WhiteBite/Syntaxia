import { useLogger } from '@/composables/useLogger'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
import {
    DEFAULT_SECTION_ORDER,
    DEFAULT_TEMPLATES,
    SECTION_META,
    SUGGESTION_KEYWORDS,
    type PromptTemplate,
    type TaskHistoryItem,
    type TemplateContext,
    type TemplateSections
} from './template.types'

export { DEFAULT_SECTION_ORDER as SECTION_ORDER }

const logger = useLogger('TemplateStore')
const STORAGE_KEY = 'prompt-templates-v3'
const MAX_TASK_HISTORY = 10

export const useTemplateStore = defineStore('templates', () => {
    const templates = ref<PromptTemplate[]>([])
    const activeTemplateId = ref<string>('architect')
    const currentTask = ref<string>('')
    const userRules = ref<string>('')
    const taskHistory = ref<TaskHistoryItem[]>([])
    const isModalOpen = ref(false)

    function initTemplates() {
        const saved = loadFromStorage()
        const now = new Date().toISOString()
        const builtIn = DEFAULT_TEMPLATES.map(t => ({
            ...t, createdAt: now, updatedAt: now,
            isFavorite: saved.find((s: PromptTemplate) => s.id === t.id)?.isFavorite ?? false,
            isHidden: saved.find((s: PromptTemplate) => s.id === t.id)?.isHidden ?? false
        }))
        templates.value = [...builtIn, ...saved.filter((t: PromptTemplate) => !t.isBuiltIn)]
        const savedActiveId = localStorage.getItem('active-template-id')
        if (savedActiveId && templates.value.some(t => t.id === savedActiveId)) activeTemplateId.value = savedActiveId
        const savedTask = localStorage.getItem('template-task')
        if (savedTask) currentTask.value = savedTask
        const savedRules = localStorage.getItem('template-user-rules')
        if (savedRules) userRules.value = savedRules
        try {
            const h = localStorage.getItem('template-task-history')
            if (h) taskHistory.value = JSON.parse(h)
        } catch { /* ignore */ }
    }

    const activeTemplate = computed(() => templates.value.find(t => t.id === activeTemplateId.value) || templates.value[0])
    const builtInTemplates = computed(() => templates.value.filter(t => t.isBuiltIn && !t.isHidden))
    const customTemplates = computed(() => templates.value.filter(t => !t.isBuiltIn))
    const favoriteTemplates = computed(() => templates.value.filter(t => t.isFavorite && !t.isHidden))
    const visibleTemplates = computed(() => templates.value.filter(t => !t.isHidden))

    function loadFromStorage(): PromptTemplate[] {
        try { const s = localStorage.getItem(STORAGE_KEY); if (s) return JSON.parse(s) } catch (e) { logger.warn('Load failed:', e) }
        return []
    }

    function saveToStorage() {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(templates.value))
            localStorage.setItem('active-template-id', activeTemplateId.value)
            localStorage.setItem('template-task', currentTask.value)
            localStorage.setItem('template-user-rules', userRules.value)
            localStorage.setItem('template-task-history', JSON.stringify(taskHistory.value))
        } catch (e) { logger.warn('Save failed:', e) }
    }

    watch([templates, activeTemplateId, currentTask, userRules, taskHistory], saveToStorage, { deep: true })

    function setActiveTemplate(id: string) { if (templates.value.some(t => t.id === id)) activeTemplateId.value = id }
    function setTask(task: string) { currentTask.value = task }
    function setUserRules(rules: string) { userRules.value = rules }

    function addToTaskHistory(task: string) {
        if (!task.trim()) return
        taskHistory.value = taskHistory.value.filter(h => h.text !== task)
        taskHistory.value.unshift({ id: `task-${Date.now()}`, text: task, templateId: activeTemplateId.value, timestamp: new Date().toISOString() })
        if (taskHistory.value.length > MAX_TASK_HISTORY) taskHistory.value = taskHistory.value.slice(0, MAX_TASK_HISTORY)
    }

    function clearTaskHistory() { taskHistory.value = [] }
    function toggleFavorite(id: string) { const t = templates.value.find(x => x.id === id); if (t) { t.isFavorite = !t.isFavorite; t.updatedAt = new Date().toISOString() } }
    function toggleHidden(id: string) { const t = templates.value.find(x => x.id === id); if (t) { t.isHidden = !t.isHidden; t.updatedAt = new Date().toISOString() } }

    function createTemplate(tpl: Omit<PromptTemplate, 'id' | 'createdAt' | 'updatedAt'>): string {
        const now = new Date().toISOString(), id = `custom-${Date.now()}`
        templates.value.push({ ...tpl, id, isBuiltIn: false, createdAt: now, updatedAt: now })
        return id
    }

    function updateTemplate(id: string, updates: Partial<PromptTemplate>) {
        const idx = templates.value.findIndex(t => t.id === id)
        if (idx !== -1) templates.value[idx] = { ...templates.value[idx], ...updates, updatedAt: new Date().toISOString() }
    }

    function deleteTemplate(id: string): boolean {
        const t = templates.value.find(x => x.id === id)
        if (t?.isBuiltIn) return false
        templates.value = templates.value.filter(x => x.id !== id)
        if (activeTemplateId.value === id) activeTemplateId.value = templates.value[0]?.id || 'architect'
        return true
    }

    function duplicateTemplate(id: string): string | null {
        const t = templates.value.find(x => x.id === id)
        if (!t) return null
        return createTemplate({ ...t, name: `${t.name} (copy)`, isBuiltIn: false, isFavorite: false, isHidden: false })
    }

    function resetToDefault(id: string): boolean {
        const t = templates.value.find(x => x.id === id)
        if (!t?.isBuiltIn) return false
        const defaultTpl = DEFAULT_TEMPLATES.find(d => d.id === id)
        if (!defaultTpl) return false
        // Reset customizations
        t.isFavorite = false
        t.isHidden = false
        t.updatedAt = new Date().toISOString()
        return true
    }

    function openModal() { isModalOpen.value = true }
    function closeModal() { isModalOpen.value = false }

    function suggestTemplate(taskText: string): string | null {
        if (!taskText.trim()) return null
        const lower = taskText.toLowerCase()
        for (const [tid, kws] of Object.entries(SUGGESTION_KEYWORDS)) { if (kws.some(kw => lower.includes(kw))) return tid }
        return null
    }

    function generatePrompt(context: TemplateContext): string {
        const tpl = activeTemplate.value
        if (!tpl) return context.files
        const parts: string[] = []
        if (tpl.customPrefix) parts.push(tpl.customPrefix)
        for (const sec of (tpl.sectionOrder || DEFAULT_SECTION_ORDER)) {
            if (!tpl.sections[sec]) continue
            switch (sec) {
                case 'role': if (tpl.roleContent) parts.push(`## Role\n${tpl.roleContent}`); break
                case 'rules': if (tpl.rulesContent || userRules.value) parts.push(`## Rules\n${[tpl.rulesContent, userRules.value].filter(Boolean).join('\n\n')}`); break
                case 'tree': if (context.fileTree) parts.push(`## Project Structure\n\`\`\`\n${context.fileTree}\n\`\`\``); break
                case 'stats': parts.push(`## Context Stats\n- Files: ${context.fileCount}\n- Tokens: ~${context.tokenCount}\n- Languages: ${context.languages.join(', ') || 'N/A'}`); break
                case 'task': if (currentTask.value) parts.push(`## Task\n${currentTask.value}`); break
                case 'files': if (context.files) parts.push(`## Files\n${context.files}`); break
            }
        }
        if (tpl.customSuffix) parts.push(tpl.customSuffix)
        return parts.join('\n\n')
    }

    function generatePreview(context: Partial<TemplateContext>): string {
        const tpl = activeTemplate.value
        if (!tpl) return ''
        const parts: string[] = []
        if (tpl.customPrefix) parts.push(tpl.customPrefix)
        for (const sec of (tpl.sectionOrder || DEFAULT_SECTION_ORDER)) {
            if (!tpl.sections[sec]) continue
            switch (sec) {
                case 'role': if (tpl.roleContent) parts.push(`## Role\n${tpl.roleContent}`); break
                case 'rules': if (tpl.rulesContent || userRules.value) parts.push(`## Rules\n${[tpl.rulesContent, userRules.value].filter(Boolean).join('\n\n')}`); break
                case 'tree': parts.push(`## Project Structure\n\`\`\`\n${context.fileTree || '[File tree]'}\n\`\`\``); break
                case 'stats': parts.push(`## Context Stats\n- Files: ${context.fileCount || 0}\n- Tokens: ~${context.tokenCount || 0}\n- Languages: ${context.languages?.join(', ') || 'N/A'}`); break
                case 'task': parts.push(`## Task\n${currentTask.value || '[Your task]'}`); break
                case 'files': parts.push(`## Files\n[${context.fileCount || 0} files]`); break
            }
        }
        if (tpl.customSuffix) parts.push(tpl.customSuffix)
        return parts.join('\n\n')
    }

    function getSectionLabel(sec: keyof TemplateSections): string {
        return { role: 'Role', rules: 'Rules', tree: 'File Tree', stats: 'Stats', task: 'Task', files: 'Files' }[sec]
    }

    function getSectionMeta(sec: keyof TemplateSections) { return SECTION_META.find(m => m.key === sec) }

    initTemplates()

    return {
        templates, activeTemplateId, currentTask, userRules, taskHistory, isModalOpen,
        activeTemplate, builtInTemplates, customTemplates, favoriteTemplates, visibleTemplates,
        setActiveTemplate, setTask, setUserRules, addToTaskHistory, clearTaskHistory,
        toggleFavorite, toggleHidden, createTemplate, updateTemplate, deleteTemplate, duplicateTemplate, resetToDefault,
        generatePrompt, generatePreview, suggestTemplate, getSectionLabel, getSectionMeta, openModal, closeModal
    }
})
