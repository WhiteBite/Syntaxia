<template>
    <div class="language-badge-wrapper">
        <BaseTooltip :text="tooltipContent" position="top" multiline>
            <template #default>
                <button 
                    class="language-badge" 
                    @click="toggleExpanded"
                    :aria-label="t('languages.ariaLabel')"
                    :aria-expanded="isExpanded"
                >
                    <div class="badge-content">
                        <svg class="badge-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                            <path 
                                stroke-linecap="round" 
                                stroke-linejoin="round" 
                                stroke-width="2"
                                d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" 
                            />
                        </svg>
                        <span class="badge-text">{{ t('languages.count', { count: languageCount }) }}</span>
                    </div>
                </button>
            </template>
            <template #tooltip>
                <div class="tooltip-content">
                    <div class="tooltip-header">
                        <span class="tooltip-title">{{ t('languages.tooltipTitle') }}</span>
                    </div>
                    <div class="tooltip-languages">
                        <div 
                            v-for="lang in supportedLanguages" 
                            :key="lang.id" 
                            class="language-item"
                        >
                            <span class="language-icon" :style="{ color: lang.color }">
                                {{ lang.icon }}
                            </span>
                            <span class="language-name">{{ lang.name }}</span>
                        </div>
                    </div>
                </div>
            </template>
        </BaseTooltip>

        <!-- Expanded panel (optional detailed view) -->
        <Transition name="panel">
            <div v-if="isExpanded" class="language-panel">
                <div class="panel-header">
                    <div class="panel-title-row">
                        <div class="panel-icon">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                <path 
                                    stroke-linecap="round" 
                                    stroke-linejoin="round" 
                                    stroke-width="2"
                                    d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" 
                                />
                            </svg>
                        </div>
                        <div class="panel-title-text">
                            <h3 class="panel-title">{{ t('languages.panelTitle') }}</h3>
                            <span class="panel-subtitle">{{ t('languages.panelSubtitle') }}</span>
                        </div>
                    </div>
                    <button @click="collapse" class="panel-close" :aria-label="t('common.close')">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <div class="panel-content">
                    <div class="languages-grid">
                        <div 
                            v-for="lang in supportedLanguages" 
                            :key="lang.id" 
                            class="language-card"
                        >
                            <span class="card-icon" :style="{ color: lang.color }">
                                {{ lang.icon }}
                            </span>
                            <span class="card-name">{{ lang.name }}</span>
                            <span class="card-ext">{{ lang.extensions.join(', ') }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </Transition>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { BaseTooltip } from '@/components/ui'

interface SupportedLanguage {
    id: string
    name: string
    icon: string
    color: string
    extensions: string[]
}

const { t } = useI18n()

const isExpanded = ref(false)

const supportedLanguages: SupportedLanguage[] = [
    { id: 'typescript', name: 'TypeScript', icon: 'TS', color: '#3178c6', extensions: ['.ts', '.tsx'] },
    { id: 'javascript', name: 'JavaScript', icon: 'JS', color: '#f7df1e', extensions: ['.js', '.jsx'] },
    { id: 'python', name: 'Python', icon: 'PY', color: '#3776ab', extensions: ['.py'] },
    { id: 'go', name: 'Go', icon: 'GO', color: '#00add8', extensions: ['.go'] },
    { id: 'rust', name: 'Rust', icon: 'RS', color: '#dea584', extensions: ['.rs'] },
    { id: 'java', name: 'Java', icon: 'JV', color: '#ed8b00', extensions: ['.java'] },
    { id: 'csharp', name: 'C#', icon: 'C#', color: '#512bd4', extensions: ['.cs'] },
    { id: 'cpp', name: 'C++', icon: 'C+', color: '#00599c', extensions: ['.cpp', '.hpp', '.cc'] },
    { id: 'php', name: 'PHP', icon: 'PH', color: '#777bb4', extensions: ['.php'] },
    { id: 'ruby', name: 'Ruby', icon: 'RB', color: '#cc342d', extensions: ['.rb'] },
    { id: 'swift', name: 'Swift', icon: 'SW', color: '#f05138', extensions: ['.swift'] },
    { id: 'kotlin', name: 'Kotlin', icon: 'KT', color: '#7f52ff', extensions: ['.kt', '.kts'] },
    { id: 'vue', name: 'Vue', icon: 'VU', color: '#42b883', extensions: ['.vue'] },
    { id: 'html', name: 'HTML/CSS', icon: 'HT', color: '#e34f26', extensions: ['.html', '.css', '.scss'] },
]

const languageCount = computed(() => supportedLanguages.length)

const tooltipContent = computed(() => t('languages.clickToExpand'))

function toggleExpanded(): void {
    isExpanded.value = !isExpanded.value
}

function collapse(): void {
    isExpanded.value = false
}
</script>

<style scoped>
.language-badge-wrapper {
    position: relative;
}

/* Badge button */
.language-badge {
    @apply flex items-center rounded-xl cursor-pointer;
    @apply transition-all duration-300 ease-out;
    padding: 8px 14px;
    background: rgba(15, 18, 25, 0.9);
    border: 1px solid rgba(16, 185, 129, 0.3);
    backdrop-filter: blur(12px);
    box-shadow:
        0 4px 20px rgba(0, 0, 0, 0.4),
        0 0 40px rgba(16, 185, 129, 0.1);
}

.language-badge:hover {
    border-color: rgba(16, 185, 129, 0.5);
    transform: translateY(-2px);
    box-shadow:
        0 8px 30px rgba(0, 0, 0, 0.5),
        0 0 60px rgba(16, 185, 129, 0.2);
}

.badge-content {
    @apply flex items-center gap-2;
}

.badge-icon {
    @apply w-4 h-4;
    color: #10b981;
}

.badge-text {
    @apply text-sm font-medium;
    color: #e2e8f0;
}

/* Tooltip content */
.tooltip-content {
    @apply flex flex-col gap-2;
}

.tooltip-header {
    @apply pb-2 border-b border-gray-700;
}

.tooltip-title {
    @apply text-xs font-semibold text-white;
}

.tooltip-languages {
    @apply grid grid-cols-2 gap-1.5;
}

.language-item {
    @apply flex items-center gap-1.5 py-0.5;
}

.language-icon {
    @apply text-[10px] font-bold w-5 text-center;
}

.language-name {
    @apply text-xs text-gray-300;
}

/* Expanded panel */
.language-panel {
    @apply flex flex-col rounded-2xl overflow-hidden;
    position: absolute;
    bottom: calc(100% + 8px);
    left: 0;
    width: min(320px, 90vw);
    background: rgba(15, 18, 25, 0.95);
    border: 1px solid rgba(16, 185, 129, 0.25);
    backdrop-filter: blur(16px);
    box-shadow:
        0 8px 40px rgba(0, 0, 0, 0.6),
        0 0 80px rgba(16, 185, 129, 0.15),
        inset 0 1px 0 rgba(255, 255, 255, 0.05);
    z-index: var(--z-dropdown);
}

.panel-header {
    @apply flex items-center justify-between p-4;
    background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, transparent 100%);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.panel-title-row {
    @apply flex items-center gap-3;
}

.panel-icon {
    @apply w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0;
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.panel-icon svg {
    @apply w-4 h-4 text-white;
}

.panel-title-text {
    @apply flex flex-col;
}

.panel-title {
    @apply text-sm font-semibold text-white leading-tight;
}

.panel-subtitle {
    @apply text-xs;
    color: #10b981;
}

.panel-close {
    @apply p-1.5 rounded-lg;
    @apply text-gray-400 hover:text-white hover:bg-white/10;
    @apply transition-colors duration-150;
}

/* Panel content */
.panel-content {
    @apply p-3;
    max-height: min(280px, 40vh);
    overflow-y: auto;
}

.languages-grid {
    @apply grid grid-cols-2 gap-2;
}

.language-card {
    @apply flex flex-col items-center p-2.5 rounded-lg;
    @apply transition-colors duration-150;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.language-card:hover {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba(16, 185, 129, 0.2);
}

.card-icon {
    @apply text-sm font-bold mb-1;
}

.card-name {
    @apply text-xs font-medium text-white;
}

.card-ext {
    @apply text-[10px] text-gray-500 mt-0.5;
}

/* Panel animation */
.panel-enter-active {
    animation: panel-in 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.panel-leave-active {
    animation: panel-out 0.2s ease-in;
}

@keyframes panel-in {
    from {
        opacity: 0;
        transform: translateY(10px) scale(0.95);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

@keyframes panel-out {
    from {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
    to {
        opacity: 0;
        transform: translateY(10px) scale(0.95);
    }
}
</style>
