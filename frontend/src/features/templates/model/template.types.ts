// Template system types

export interface TemplateSections {
    role: boolean
    rules: boolean
    tree: boolean
    stats: boolean
    task: boolean
    files: boolean
}

export interface TemplateVariable {
    key: string
    label: string
    description: string
    example: string
}

export interface PromptTemplate {
    id: string
    name: string
    icon: string
    description: string
    tags: string[]
    sections: TemplateSections
    sectionOrder: (keyof TemplateSections)[]
    roleContent: string
    rulesContent: string
    customPrefix: string
    customSuffix: string
    isBuiltIn: boolean
    isFavorite: boolean
    isHidden: boolean
    createdAt: string
    updatedAt: string
}

export interface TemplateContext {
    fileTree: string
    files: string
    task: string
    userRules: string
    fileCount: number
    tokenCount: number
    languages: string[]
    projectName: string
}

export interface TaskHistoryItem {
    id: string
    text: string
    templateId: string
    timestamp: string
}

// Default section order
export const DEFAULT_SECTION_ORDER: (keyof TemplateSections)[] = [
    'role',
    'rules',
    'tree',
    'stats',
    'task',
    'files'
]

// Section metadata for UI
export interface SectionMeta {
    key: keyof TemplateSections
    icon: string
    color: string
}

export const SECTION_META: SectionMeta[] = [
    { key: 'role', icon: '🎭', color: 'purple' },
    { key: 'rules', icon: '📋', color: 'blue' },
    { key: 'tree', icon: '🌳', color: 'green' },
    { key: 'stats', icon: '📊', color: 'orange' },
    { key: 'task', icon: '📝', color: 'yellow' },
    { key: 'files', icon: '📁', color: 'cyan' }
]

// Template tags for smart suggestions
export const TEMPLATE_TAGS = [
    'architecture',
    'implementation',
    'review',
    'bugfix',
    'refactor',
    'documentation',
    'testing',
    'optimization'
] as const

export type TemplateTag = typeof TEMPLATE_TAGS[number]

// Available template variables
export const TEMPLATE_VARIABLES: TemplateVariable[] = [
    { key: '{{file_tree}}', label: 'File Tree', description: 'Project file structure', example: '├── src\n│   └── main.ts' },
    { key: '{{files}}', label: 'Files Content', description: 'Selected files content', example: '<file path="src/main.ts">...</file>' },
    { key: '{{task}}', label: 'Task', description: 'User task description', example: 'Implement feature X' },
    { key: '{{user_rules}}', label: 'User Rules', description: 'Custom rules from settings', example: 'Follow clean code principles' },
    { key: '{{file_count}}', label: 'File Count', description: 'Number of selected files', example: '15' },
    { key: '{{token_count}}', label: 'Token Count', description: 'Estimated token count', example: '12500' },
    { key: '{{languages}}', label: 'Languages', description: 'Detected programming languages', example: 'TypeScript, Go, Vue' },
    { key: '{{project_name}}', label: 'Project Name', description: 'Current project name', example: 'syntaxia' },
]

// Default templates
export const DEFAULT_TEMPLATES: Omit<PromptTemplate, 'createdAt' | 'updatedAt'>[] = [
    {
        id: 'architect',
        name: 'Architect',
        icon: '🏗️',
        description: 'System design and architecture analysis',
        tags: ['architecture', 'refactor'],
        isBuiltIn: true,
        isFavorite: false,
        isHidden: false,
        sections: { role: true, rules: true, tree: true, stats: true, task: true, files: true },
        sectionOrder: ['role', 'rules', 'tree', 'stats', 'task', 'files'],
        roleContent: 'You are a senior software architect. Analyze the codebase structure, identify patterns, and provide architectural recommendations.',
        rulesContent: '- Focus on scalability and maintainability\n- Consider SOLID principles\n- Identify potential technical debt\n- Suggest improvements with rationale',
        customPrefix: '',
        customSuffix: ''
    },
    {
        id: 'implement',
        name: 'Implement',
        icon: '⚡',
        description: 'Code implementation and feature development',
        tags: ['implementation'],
        isBuiltIn: true,
        isFavorite: false,
        isHidden: false,
        sections: { role: true, rules: true, tree: true, stats: false, task: true, files: true },
        sectionOrder: ['role', 'rules', 'tree', 'task', 'files'],
        roleContent: 'You are an expert developer. Implement the requested feature following the existing code patterns and conventions.',
        rulesContent: '- Follow existing code style\n- Write clean, readable code\n- Add appropriate comments\n- Consider edge cases',
        customPrefix: '',
        customSuffix: ''
    },
    {
        id: 'review',
        name: 'Review',
        icon: '🔍',
        description: 'Code review and quality analysis',
        tags: ['review', 'bugfix'],
        isBuiltIn: true,
        isFavorite: false,
        isHidden: false,
        sections: { role: true, rules: true, tree: false, stats: true, task: true, files: true },
        sectionOrder: ['role', 'rules', 'stats', 'task', 'files'],
        roleContent: 'You are a code reviewer. Analyze the code for bugs, security issues, performance problems, and best practices violations.',
        rulesContent: '- Check for security vulnerabilities\n- Identify performance bottlenecks\n- Verify error handling\n- Suggest improvements',
        customPrefix: '',
        customSuffix: ''
    },
    {
        id: 'explain',
        name: 'Explain',
        icon: '📚',
        description: 'Code explanation and documentation',
        tags: ['documentation'],
        isBuiltIn: true,
        isFavorite: false,
        isHidden: false,
        sections: { role: true, rules: false, tree: true, stats: false, task: true, files: true },
        sectionOrder: ['role', 'tree', 'task', 'files'],
        roleContent: 'You are a technical writer. Explain the code clearly, document its purpose, and describe how it works.',
        rulesContent: '',
        customPrefix: '',
        customSuffix: ''
    }
]

export function createEmptyTemplate(): Omit<PromptTemplate, 'id' | 'createdAt' | 'updatedAt'> {
    return {
        name: '',
        icon: '📝',
        description: '',
        tags: [],
        isBuiltIn: false,
        isFavorite: false,
        isHidden: false,
        sections: { role: true, rules: true, tree: true, stats: true, task: true, files: true },
        sectionOrder: [...DEFAULT_SECTION_ORDER],
        roleContent: '',
        rulesContent: '',
        customPrefix: '',
        customSuffix: ''
    }
}

// Smart suggestion keywords
export const SUGGESTION_KEYWORDS: Record<string, string[]> = {
    architect: ['architecture', 'design', 'structure', 'refactor', 'модуль', 'архитектур', 'структур'],
    implement: ['implement', 'add', 'create', 'build', 'feature', 'добавить', 'создать', 'реализовать', 'функци'],
    review: ['review', 'check', 'bug', 'fix', 'error', 'проверить', 'баг', 'ошибк', 'исправить'],
    explain: ['explain', 'what', 'how', 'why', 'document', 'объясни', 'как', 'почему', 'документ']
}
