<template>
    <div class="version-widget" :class="{ 'is-expanded': isExpanded }">
        <!-- Collapsed state - just version badge -->
        <button 
            v-if="!isExpanded" 
            @click="expand" 
            class="version-trigger" 
            :class="{ 'has-update': hasUpdate }"
            :title="isDev ? t('version.devModeTooltip') : t('version.clickToViewChangelog')">
            <div class="trigger-content">
                <svg class="trigger-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
                <span class="trigger-version">{{ displayVersion }}</span>
                <span v-if="hasUpdate" class="update-dot"></span>
            </div>
        </button>

        <!-- Expanded panel -->
        <Transition name="panel">
            <div v-if="isExpanded" class="version-panel">
                <!-- Panel header -->
                <div class="panel-header">
                    <div class="panel-title-row">
                        <div class="panel-icon">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <div class="panel-title-text">
                            <h3 class="panel-title">{{ t('version.changelog') }}</h3>
                            <span class="panel-version">{{ displayVersion }}</span>
                        </div>
                    </div>
                    <button @click="collapse" class="panel-close">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <!-- Panel content -->
                <div class="panel-content">
                    <!-- Loading -->
                    <div v-if="isLoading" class="panel-state">
                        <div class="spinner"></div>
                        <span>{{ t('version.checkingUpdates') }}</span>
                    </div>

                    <!-- Error -->
                    <div v-else-if="error" class="panel-state panel-state-error">
                        <span>{{ t('version.errorLoading') }}</span>
                        <button @click="loadReleases" class="retry-link">{{ t('action.retry') || 'Повторить' }}</button>
                    </div>

                    <!-- No releases -->
                    <div v-else-if="releases.length === 0" class="panel-state">
                        <span class="no-releases-text">{{ t('version.noReleases') }}</span>
                        <span class="no-releases-hint">Релизы появятся после публикации</span>
                    </div>

                    <!-- Releases list -->
                    <div v-else class="releases-list">
                        <div v-for="(release, index) in releases.slice(0, 5)" :key="release.tag_name"
                            class="release-item" :class="{ 'is-latest': index === 0 }">
                            <div class="release-row">
                                <span class="release-version">{{ release.tag_name }}</span>
                                <span v-if="index === 0" class="release-tag">Latest</span>
                                <span v-if="release.prerelease" class="release-tag release-tag-pre">Pre</span>
                            </div>
                            <span class="release-date">{{ formatDate(release.published_at) }}</span>
                        </div>
                    </div>
                </div>

                <!-- Panel footer -->
                <a href="https://github.com/WhiteBite/syntaxia/releases" target="_blank" rel="noopener noreferrer"
                    class="panel-footer">
                    <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                        <path
                            d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
                    </svg>
                    <span>{{ t('version.viewOnGithub') }}</span>
                    <svg class="w-3 h-3 ml-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                </a>
            </div>
        </Transition>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, onMounted, ref } from 'vue'
import { GetReleases, GetVersionInfo } from '../../wailsjs/go/main/App'

interface Release {
    tag_name: string
    name: string
    body: string
    published_at: string
    html_url: string
    prerelease: boolean
}

interface ReleasesResponse {
    currentVersion: string
    latestVersion: string
    hasUpdate: boolean
    releases: Release[]
    error?: string
}

const { t, locale } = useI18n()

const isExpanded = ref(false)
const isLoading = ref(false)
const error = ref<string | null>(null)
const currentVersion = ref('dev')
const latestVersion = ref('')
const hasUpdate = ref(false)
const releases = ref<Release[]>([])

const isDev = computed(() => {
    const v = currentVersion.value.toLowerCase()
    return v === 'dev' || v === 'vdev' || v.includes('dev')
})

const displayVersion = computed(() => {
    const v = currentVersion.value
    if (isDev.value) return 'vdev'
    return v.startsWith('v') ? v : `v${v}`
})

onMounted(async () => {
    try {
        const info = await GetVersionInfo()
        currentVersion.value = info.version
    } catch (e) {
        console.warn('Failed to get version info:', e)
    }
})

async function loadReleases() {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
        const response: ReleasesResponse = await GetReleases()
        currentVersion.value = response.currentVersion
        latestVersion.value = response.latestVersion
        hasUpdate.value = response.hasUpdate
        releases.value = response.releases || []

        if (response.error) {
            error.value = response.error
        }
    } catch (e) {
        console.error('Failed to load releases:', e)
        error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
        isLoading.value = false
    }
}

function expand() {
    isExpanded.value = true
    if (releases.value.length === 0 && !error.value) {
        loadReleases()
    }
}

function collapse() {
    isExpanded.value = false
}

function formatDate(dateStr: string): string {
    if (!dateStr) return ''
    const date = new Date(dateStr)
    return date.toLocaleDateString(locale.value === 'ru' ? 'ru-RU' : 'en-US', {
        month: 'short',
        day: 'numeric'
    })
}
</script>

<style scoped>
.version-widget {
    position: relative;
}

/* Trigger button */
.version-trigger {
    @apply flex items-center rounded-xl cursor-pointer;
    @apply transition-all duration-300 ease-out;
    padding: 8px 14px;
    background: rgba(15, 18, 25, 0.9);
    border: 1px solid rgba(139, 92, 246, 0.3);
    backdrop-filter: blur(12px);
    box-shadow:
        0 4px 20px rgba(0, 0, 0, 0.4),
        0 0 40px rgba(139, 92, 246, 0.1);
}

.version-trigger:hover {
    border-color: rgba(139, 92, 246, 0.5);
    transform: translateY(-2px);
    box-shadow:
        0 8px 30px rgba(0, 0, 0, 0.5),
        0 0 60px rgba(139, 92, 246, 0.2);
}

.version-trigger.has-update {
    border-color: rgba(16, 185, 129, 0.4);
    box-shadow:
        0 4px 20px rgba(0, 0, 0, 0.4),
        0 0 40px rgba(16, 185, 129, 0.15);
}

.trigger-content {
    @apply flex items-center gap-2;
}

.trigger-icon {
    @apply w-4 h-4;
    color: #a78bfa;
}

.trigger-version {
    @apply text-sm font-mono font-medium;
    color: #e2e8f0;
}

.update-dot {
    @apply w-2 h-2 rounded-full;
    background: #10b981;
    box-shadow: 0 0 8px #10b981;
    animation: pulse-glow 2s ease-in-out infinite;
}

@keyframes pulse-glow {

    0%,
    100% {
        opacity: 1;
        box-shadow: 0 0 8px #10b981;
    }

    50% {
        opacity: 0.6;
        box-shadow: 0 0 16px #10b981;
    }
}

/* Expanded panel */
.version-panel {
    @apply flex flex-col rounded-2xl overflow-hidden;
    width: 280px;
    background: rgba(15, 18, 25, 0.95);
    border: 1px solid rgba(139, 92, 246, 0.25);
    backdrop-filter: blur(16px);
    box-shadow:
        0 8px 40px rgba(0, 0, 0, 0.6),
        0 0 80px rgba(139, 92, 246, 0.15),
        inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.panel-header {
    @apply flex items-center justify-between p-4;
    background: linear-gradient(135deg, rgba(139, 92, 246, 0.1) 0%, transparent 100%);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.panel-title-row {
    @apply flex items-center gap-3;
}

.panel-icon {
    @apply w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0;
    background: linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%);
    box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
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

.panel-version {
    @apply text-xs font-mono;
    color: #a78bfa;
}

.panel-close {
    @apply p-1.5 rounded-lg;
    @apply text-gray-400 hover:text-white hover:bg-white/10;
    @apply transition-colors duration-150;
}

/* Panel content */
.panel-content {
    @apply p-3;
    max-height: 240px;
    overflow-y: auto;
}

.panel-state {
    @apply flex flex-col items-center justify-center gap-2 py-6 text-center;
    color: #94a3b8;
}

.panel-state-error {
    color: #f87171;
}

.spinner {
    @apply w-6 h-6 rounded-full;
    border: 2px solid rgba(139, 92, 246, 0.2);
    border-top-color: #8b5cf6;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.retry-link {
    @apply text-xs underline;
    color: #f87171;
}

.retry-link:hover {
    color: #fca5a5;
}

.no-releases-text {
    @apply text-sm;
}

.no-releases-hint {
    @apply text-xs;
    color: #64748b;
}

/* Releases list */
.releases-list {
    @apply space-y-1;
}

.release-item {
    @apply flex items-center justify-between p-2.5 rounded-lg;
    @apply transition-colors duration-150;
}

.release-item:hover {
    background: rgba(255, 255, 255, 0.05);
}

.release-item.is-latest {
    background: rgba(139, 92, 246, 0.1);
}

.release-row {
    @apply flex items-center gap-2;
}

.release-version {
    @apply text-sm font-mono font-medium text-white;
}

.release-tag {
    @apply px-1.5 py-0.5 rounded text-[10px] font-semibold uppercase;
    background: rgba(139, 92, 246, 0.2);
    color: #c4b5fd;
}

.release-tag-pre {
    background: rgba(251, 191, 36, 0.15);
    color: #fcd34d;
}

.release-date {
    @apply text-xs;
    color: #64748b;
}

/* Panel footer */
.panel-footer {
    @apply flex items-center gap-2 px-4 py-3 text-sm;
    @apply transition-colors duration-150;
    color: #94a3b8;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    background: rgba(0, 0, 0, 0.2);
}

.panel-footer:hover {
    color: #e2e8f0;
    background: rgba(255, 255, 255, 0.05);
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
