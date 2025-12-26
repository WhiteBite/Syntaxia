/**
 * Filter configuration constants
 */
import type { SmartFilter } from '../model/types'

export const STORAGE_KEY = 'quick-filters-state'

export const languageConfig: Record<string, { icon: string }> = {
    'Go': { icon: '🐹' },
    'TypeScript': { icon: '📘' },
    'JavaScript': { icon: '📒' },
    'Vue': { icon: '💚' },
    'Python': { icon: '🐍' },
    'Java': { icon: '☕' },
    'Kotlin': { icon: '🟣' },
    'Rust': { icon: '🦀' },
    'C#': { icon: '🟦' },
    'C++': { icon: '⚡' },
    'C': { icon: '🔧' },
    'Ruby': { icon: '💎' },
    'PHP': { icon: '🐘' },
    'Swift': { icon: '🍎' },
    'Dart': { icon: '🎯' },
}

export const languageExtensions: Record<string, string[]> = {
    'Go': ['.go'],
    'TypeScript': ['.ts', '.tsx'],
    'JavaScript': ['.js', '.jsx'],
    'Vue': ['.vue'],
    'Python': ['.py'],
    'Java': ['.java'],
    'Kotlin': ['.kt', '.kts'],
    'Rust': ['.rs'],
    'C#': ['.cs'],
    'C++': ['.cpp', '.cc', '.cxx', '.hpp', '.h'],
    'C': ['.c', '.h'],
    'Ruby': ['.rb'],
    'PHP': ['.php'],
    'Swift': ['.swift'],
    'Dart': ['.dart'],
}

// i18n keys for filter labels
export const filterLabelKeys: Record<string, string> = {
    'components': 'smartFilters.components',
    'composables': 'smartFilters.composables',
    'stores': 'smartFilters.stores',
    'views': 'smartFilters.views',
    'hooks': 'smartFilters.hooks',
    'pages': 'smartFilters.pages',
    'api': 'smartFilters.api',
    'services': 'smartFilters.services',
    'modules': 'smartFilters.modules',
    'handlers': 'smartFilters.handlers',
    'domain': 'smartFilters.domain',
    'infra': 'smartFilters.infra',
    'backend': 'smartFilters.backend',
    'frontend': 'smartFilters.frontend',
    'models': 'smartFilters.models',
    'urls': 'smartFilters.urls',
    'routes': 'smartFilters.routes',
    'controllers': 'smartFilters.controllers',
    'repos': 'smartFilters.repos',
    'screens': 'smartFilters.screens',
    'widgets': 'smartFilters.widgets',
    'state': 'smartFilters.state',
}

type SmartFilterDef = Omit<SmartFilter, 'label' | 'shortLabel'> & { labelKey: string }

const createSmartFilter = (
    id: string,
    labelKey: string,
    icon: string,
    extensions: string[],
    patterns: string[],
    framework: string
): SmartFilterDef => ({
    id,
    labelKey,
    icon,
    extensions,
    patterns,
    framework,
    category: 'smart',
})

export const frameworkFiltersConfig: Record<string, SmartFilterDef[]> = {
    'Vue.js': [
        createSmartFilter('vue-components', 'components', '🧩', ['.vue'], ['**/components/**'], 'Vue.js'),
        createSmartFilter('vue-composables', 'composables', '🪝', ['.ts'], ['**/composables/**', '**/use*.ts'], 'Vue.js'),
        createSmartFilter('vue-stores', 'stores', '🗄️', ['.ts'], ['**/stores/**', '**/*.store.ts'], 'Vue.js'),
        createSmartFilter('vue-views', 'views', '📄', ['.vue'], ['**/views/**', '**/pages/**'], 'Vue.js'),
    ],
    'React': [
        createSmartFilter('react-components', 'components', '🧩', ['.tsx', '.jsx'], ['**/components/**'], 'React'),
        createSmartFilter('react-hooks', 'hooks', '🪝', ['.ts', '.tsx'], ['**/hooks/**', '**/use*.ts', '**/use*.tsx'], 'React'),
        createSmartFilter('react-pages', 'pages', '📄', ['.tsx', '.jsx'], ['**/pages/**', '**/app/**'], 'React'),
    ],
    'Next.js': [
        createSmartFilter('next-pages', 'pages', '📄', ['.tsx', '.jsx'], ['**/app/**', '**/pages/**'], 'Next.js'),
        createSmartFilter('next-components', 'components', '🧩', ['.tsx', '.jsx'], ['**/components/**'], 'Next.js'),
        createSmartFilter('next-api', 'api', '🔌', ['.ts', '.tsx'], ['**/api/**', '**/route.ts'], 'Next.js'),
    ],
    'Angular': [
        createSmartFilter('angular-components', 'components', '🧩', ['.ts'], ['**/*.component.ts'], 'Angular'),
        createSmartFilter('angular-services', 'services', '⚙️', ['.ts'], ['**/*.service.ts'], 'Angular'),
        createSmartFilter('angular-modules', 'modules', '📦', ['.ts'], ['**/*.module.ts'], 'Angular'),
    ],
    'Gin': [
        createSmartFilter('go-handlers', 'handlers', '🎯', ['.go'], ['**/handlers/**', '**/*_handler.go'], 'Gin'),
        createSmartFilter('go-services', 'services', '⚙️', ['.go'], ['**/services/**', '**/*_service.go', '**/application/**'], 'Gin'),
        createSmartFilter('go-domain', 'domain', '🏛️', ['.go'], ['**/domain/**', '**/entities/**', '**/models/**'], 'Gin'),
        createSmartFilter('go-infra', 'infra', '🔧', ['.go'], ['**/infrastructure/**', '**/repository/**', '**/adapters/**'], 'Gin'),
    ],
    'Echo': [
        createSmartFilter('go-handlers', 'handlers', '🎯', ['.go'], ['**/handlers/**', '**/*_handler.go'], 'Echo'),
        createSmartFilter('go-services', 'services', '⚙️', ['.go'], ['**/services/**', '**/*_service.go'], 'Echo'),
    ],
    'Wails': [
        createSmartFilter('wails-backend', 'backend', '🐹', ['.go'], ['**/backend/**', '**/*.go'], 'Wails'),
        createSmartFilter('wails-frontend', 'frontend', '🎨', ['.vue', '.tsx', '.ts'], ['**/frontend/**'], 'Wails'),
    ],
    'Django': [
        createSmartFilter('django-views', 'views', '👁️', ['.py'], ['**/views.py', '**/views/**'], 'Django'),
        createSmartFilter('django-models', 'models', '🗃️', ['.py'], ['**/models.py', '**/models/**'], 'Django'),
        createSmartFilter('django-urls', 'urls', '🔗', ['.py'], ['**/urls.py'], 'Django'),
    ],
    'FastAPI': [
        createSmartFilter('fastapi-routes', 'routes', '🔌', ['.py'], ['**/routes/**', '**/routers/**', '**/api/**'], 'FastAPI'),
        createSmartFilter('fastapi-models', 'models', '🗃️', ['.py'], ['**/models/**', '**/schemas/**'], 'FastAPI'),
    ],
    'Spring Boot': [
        createSmartFilter('spring-controllers', 'controllers', '🎯', ['.java', '.kt'], ['**/*Controller.java', '**/*Controller.kt'], 'Spring'),
        createSmartFilter('spring-services', 'services', '⚙️', ['.java', '.kt'], ['**/*Service.java', '**/*Service.kt'], 'Spring'),
        createSmartFilter('spring-repos', 'repos', '🗄️', ['.java', '.kt'], ['**/*Repository.java', '**/*Repository.kt'], 'Spring'),
    ],
    'Flutter': [
        createSmartFilter('flutter-screens', 'screens', '📱', ['.dart'], ['**/screens/**', '**/pages/**'], 'Flutter'),
        createSmartFilter('flutter-widgets', 'widgets', '🧩', ['.dart'], ['**/widgets/**', '**/components/**'], 'Flutter'),
        createSmartFilter('flutter-bloc', 'state', '🔄', ['.dart'], ['**/bloc/**', '**/cubit/**', '**/providers/**'], 'Flutter'),
    ],
}
