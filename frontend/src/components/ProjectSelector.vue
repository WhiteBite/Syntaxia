<template>
  <div class="project-selector" :class="{ 'drag-over': isDragging }" @drop.prevent="handleDrop"
    @dragover.prevent="isDragging = true" @dragleave.prevent="isDragging = false">
    <!-- Animated Background -->
    <div class="bg-decoration">
      <div class="bg-glow bg-glow-1"></div>
      <div class="bg-glow bg-glow-2"></div>
      <div class="bg-glow bg-glow-3"></div>
      <div class="bg-grid"></div>
    </div>

    <!-- Top Bar -->
    <div class="absolute top-4 right-4 z-10 flex items-center gap-3">
      <!-- Auto-open toggle -->
      <label class="auto-open-toggle" :title="t('welcome.autoOpen')">
        <div class="toggle-wrapper-sm">
          <input type="checkbox" class="sr-only peer" :checked="projectStore.autoOpenLast"
            @change="onToggleAutoOpen" />
          <div class="toggle-track-sm peer-checked:bg-gradient-to-r peer-checked:from-purple-600 peer-checked:to-pink-600"></div>
          <div class="toggle-thumb-sm peer-checked:translate-x-3"></div>
        </div>
        <span>{{ t('welcome.autoOpenShort') }}</span>
      </label>
      <!-- Language Switcher -->
      <button @click="toggleLanguage" class="lang-switcher"
        :title="locale === 'ru' ? 'Switch to English' : 'Переключить на русский'">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
        </svg>
        <span class="font-medium">{{ locale === 'ru' ? 'RU' : 'EN' }}</span>
      </button>
    </div>

    <!-- Drag & Drop Overlay -->
    <Transition name="fade">
      <div v-if="isDragging" class="drop-overlay">
        <div class="drop-content">
          <div class="drop-icon">
            <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
            </svg>
          </div>
          <p class="text-2xl font-bold text-white mt-4">{{ t('welcome.dropHere') }}</p>
          <p class="text-purple-200 mt-2">{{ t('welcome.toOpenProject') }}</p>
        </div>
      </div>
    </Transition>

    <!-- Main Content -->
    <div class="content-wrapper">
      <!-- Logo & Title -->
      <div class="header-section">
        <div class="logo-container">
          <div class="logo-glow"></div>
          <div class="logo">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
            </svg>
          </div>
        </div>

        <h1 class="app-title">{{ t('welcome.title') }}</h1>
        <p class="app-subtitle">{{ t('welcome.subtitle') }}</p>
        <p class="app-hint">{{ t('welcome.dragDrop') }}</p>
      </div>

      <!-- CTA Button -->
      <div class="cta-section">
        <button @click="selectProject" class="cta-button">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          {{ t('welcome.openProject') }}
        </button>
      </div>

      <!-- Recent Projects -->
      <div class="recent-section">
        <div class="section-header">
          <h2 class="section-title">
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ t('welcome.recentProjects') }}
          </h2>
          <button v-if="recentProjects.length > 0" @click="clearAllHistory" class="clear-btn"
            :title="t('welcome.clearHistory')">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            {{ t('welcome.clearHistory') }}
          </button>
        </div>

        <div class="projects-list">
          <template v-if="recentProjects.length > 0">
            <div v-for="(project, index) in recentProjects" :key="project.path" class="project-card"
              :style="{ animationDelay: `${index * 60}ms` }" @click="openRecentProject(project.path)"
              @contextmenu.prevent="showContextMenu($event, project)">
              <div class="project-icon">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
              </div>
              <div class="project-info" :title="project.path">
                <div class="project-name">{{ project.name }}</div>
                <div class="project-path">{{ shortenPath(project.path) }}</div>
              </div>
              <button @click.stop="removeProject(project.path)" class="project-remove"
                :title="t('welcome.removeProject')">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
              <svg class="project-arrow" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </template>
          <template v-else>
            <div class="empty-projects">
              <div class="empty-icon">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
              </div>
              <p>{{ t('welcome.noRecentProjects') }}</p>
            </div>
          </template>
        </div>
      </div>

      <!-- Context Menu -->
      <Teleport to="body">
        <Transition name="fade">
          <div v-if="contextMenu.visible" class="context-menu"
            :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }" @click.stop>
            <button @click="copyProjectPath" class="context-menu-item">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
              </svg>
              {{ t('welcome.copyPath') }}
            </button>
            <button @click="removeProjectFromMenu" class="context-menu-item context-menu-item-danger">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              {{ t('welcome.removeProject') }}
            </button>
          </div>
        </Transition>
      </Teleport>

    </div>

    <!-- Version Badge (bottom left) -->
    <div class="absolute bottom-4 left-4 z-10">
      <VersionBadge />
    </div>
  </div>
</template>


<script setup lang="ts">
import VersionBadge from '@/components/VersionBadge.vue';
import { useI18n } from '@/composables/useI18n';
import { apiService } from '@/services/api.service';
import { computed, onMounted, ref } from 'vue';
import { useProjectStore } from '../stores/project.store';
import { useUIStore } from '../stores/ui.store';

const emit = defineEmits<{
  (e: 'opened', path: string): void
}>()

const projectStore = useProjectStore()
const uiStore = useUIStore()
const { t, locale, setLocale } = useI18n()
const isDragging = ref(false)

const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  project: null as { path: string; name: string } | null
})

const recentProjects = computed(() => projectStore.recentProjects)

/**
 * Shortens a file path for display while keeping it readable.
 * Shows first part + ... + last 2 parts if path is too long.
 */
function shortenPath(path: string, maxLength = 50): string {
  if (path.length <= maxLength) return path
  
  const separator = path.includes('\\') ? '\\' : '/'
  const parts = path.split(separator)
  
  if (parts.length <= 3) return path
  
  // Keep first part (drive/root) and last 2 parts
  const first = parts[0]
  const last = parts.slice(-2).join(separator)
  const shortened = `${first}${separator}...${separator}${last}`
  
  return shortened.length < path.length ? shortened : path
}

onMounted(async () => {
  try {
    await projectStore.fetchRecentProjects()
  } catch (error) {
    console.error('Failed to load recent projects:', error)
  }
})

function toggleLanguage() {
  const newLocale = locale.value === 'ru' ? 'en' : 'ru'
  setLocale(newLocale)
  uiStore.addToast(
    newLocale === 'ru' ? 'Язык изменён на русский' : 'Language changed to English',
    'success'
  )
}

function onToggleAutoOpen(e: Event) {
  const target = e.target as HTMLInputElement
  projectStore.setAutoOpenLast(target.checked)
}

async function handleDrop(e: DragEvent) {
  isDragging.value = false
  const items = e.dataTransfer?.items
  if (!items) return

  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.kind === 'file') {
      const entry = item.webkitGetAsEntry?.()
      if (entry?.isDirectory) {
        // @ts-ignore
        const path = entry.fullPath
        const success = await projectStore.openProjectByPath(path)
        if (success) {
          emit('opened', path)
        }
        return
      }
    }
  }
  uiStore.addToast('Please drop a folder, not a file', 'warning')
}

async function selectProject() {
  try {
    const dirPath = await apiService.selectDirectory()
    if (!dirPath || dirPath === '') return

    const success = await projectStore.openProjectByPath(dirPath)
    if (success && projectStore.currentPath) {
      emit('opened', projectStore.currentPath)
    }
  } catch (error) {
    console.error('Failed to select project:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    uiStore.addToast(`Failed to select directory: ${errorMessage}`, 'error')
  }
}

async function openRecentProject(path: string) {
  try {
    await projectStore.openProjectByPath(path)
    emit('opened', path)
  } catch (error) {
    console.error('Failed to open recent project:', error)
    uiStore.addToast('Failed to open project', 'error')
  }
}

async function removeProject(path: string) {
  try {
    await projectStore.removeFromRecent(path)
    uiStore.addToast(t('welcome.projectRemoved'), 'success')
  } catch (error) {
    console.error('Failed to remove project:', error)
  }
}

async function clearAllHistory() {
  if (confirm(t('welcome.confirmClearHistory'))) {
    projectStore.clearRecent()
    uiStore.addToast(t('welcome.historyCleared'), 'success')
  }
}

function showContextMenu(event: MouseEvent, project: { path: string; name: string }) {
  contextMenu.value = { visible: true, x: event.clientX, y: event.clientY, project }
}

function hideContextMenu() {
  contextMenu.value.visible = false
  contextMenu.value.project = null
}

async function copyProjectPath() {
  if (contextMenu.value.project) {
    try {
      await navigator.clipboard.writeText(contextMenu.value.project.path)
      uiStore.addToast(t('welcome.pathCopied'), 'success')
    } catch (error) {
      console.error('Failed to copy path:', error)
    }
  }
  hideContextMenu()
}

async function removeProjectFromMenu() {
  if (contextMenu.value.project) {
    await removeProject(contextMenu.value.project.path)
  }
  hideContextMenu()
}

if (typeof window !== 'undefined') {
  window.addEventListener('click', hideContextMenu)
}
</script>


<style scoped>
.project-selector {
  @apply relative flex items-center justify-center h-screen text-white overflow-hidden;
  background: var(--bg-app);
}

/* Background Decorations */
.bg-decoration {
  @apply absolute inset-0 overflow-hidden pointer-events-none;
}

.bg-glow {
  @apply absolute rounded-full blur-3xl;
  animation: float 8s ease-in-out infinite;
}

.bg-glow-1 {
  @apply w-96 h-96 -top-48 -right-48;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.15) 0%, transparent 70%);
  animation-delay: 0s;
}

.bg-glow-2 {
  @apply w-80 h-80 -bottom-40 -left-40;
  background: radial-gradient(circle, rgba(236, 72, 153, 0.12) 0%, transparent 70%);
  animation-delay: -3s;
}

.bg-glow-3 {
  @apply w-[500px] h-[500px] top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.05) 0%, transparent 70%);
  animation-delay: -5s;
}

.bg-grid {
  @apply absolute inset-0 opacity-[0.02];
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.1) 1px, transparent 1px);
  background-size: 50px 50px;
}

@keyframes float {

  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }

  50% {
    transform: translate(20px, -20px) scale(1.05);
  }
}

/* Top Bar Controls */
.auto-open-toggle {
  @apply flex items-center gap-2 px-3 py-2 rounded-xl cursor-pointer;
  @apply bg-gray-800/60 border border-gray-700/50;
  @apply text-xs text-gray-400;
  @apply backdrop-blur-sm transition-all duration-200;
}

.auto-open-toggle:hover {
  @apply bg-gray-700/60 border-gray-600/50 text-gray-200;
}

/* Language Switcher */
.lang-switcher {
  @apply px-3 py-2 text-sm rounded-xl;
  @apply bg-gray-800/60 border border-gray-700/50;
  @apply text-gray-300 hover:text-white;
  @apply flex items-center gap-2;
  @apply backdrop-blur-sm transition-all duration-200;
}

.lang-switcher:hover {
  @apply bg-gray-700/60 border-gray-600/50;
  transform: scale(1.02);
}

/* Drop Overlay */
.drop-overlay {
  @apply fixed inset-8 z-50 rounded-2xl;
  @apply flex items-center justify-center;
  @apply border-2 border-dashed border-purple-500;
  background: rgba(139, 92, 246, 0.1);
  backdrop-filter: blur(8px);
}

.drop-content {
  @apply text-center;
}

.drop-icon {
  @apply w-24 h-24 mx-auto rounded-2xl;
  @apply flex items-center justify-center;
  @apply text-purple-400;
  background: var(--gradient-primary-soft);
  animation: bounce 1s ease-in-out infinite;
}

@keyframes bounce {

  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-10px);
  }
}

/* Content Wrapper */
.content-wrapper {
  @apply relative w-full max-w-2xl px-6 py-6 h-full flex flex-col;
}

/* Header Section */
.header-section {
  @apply text-center mb-6 flex-shrink-0;
}

.logo-container {
  @apply relative inline-flex items-center justify-center mb-4;
}

.logo {
  @apply w-16 h-16 rounded-2xl;
  @apply flex items-center justify-center;
  @apply relative z-10;
  background: var(--gradient-primary);
  box-shadow: var(--shadow-glow-accent);
  animation: logo-appear 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

.logo-glow {
  @apply absolute inset-0 rounded-2xl;
  background: var(--gradient-primary);
  filter: blur(20px);
  opacity: 0.5;
  animation: pulse-glow 3s ease-in-out infinite;
}

@keyframes logo-appear {
  from {
    transform: scale(0) rotate(-180deg);
    opacity: 0;
  }

  to {
    transform: scale(1) rotate(0);
    opacity: 1;
  }
}

@keyframes pulse-glow {

  0%,
  100% {
    opacity: 0.4;
    transform: scale(1);
  }

  50% {
    opacity: 0.6;
    transform: scale(1.1);
  }
}

.app-title {
  @apply text-4xl font-bold mb-2;
  background: linear-gradient(135deg, #ffffff 0%, #f1f5f9 40%, #e2e8f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 0 40px rgba(255, 255, 255, 0.1);
  animation: fade-up 0.5s ease-out 0.2s both;
}

.app-subtitle {
  @apply text-base text-gray-300 mb-1;
  animation: fade-up 0.5s ease-out 0.3s both;
}

.app-hint {
  @apply text-xs text-gray-400;
  animation: fade-up 0.5s ease-out 0.4s both;
}

@keyframes fade-up {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}


/* CTA Section */
.cta-section {
  @apply flex justify-center mb-6 flex-shrink-0;
  animation: fade-up 0.5s ease-out 0.5s both;
}

.cta-button {
  @apply px-6 py-3 text-sm font-semibold rounded-xl;
  @apply text-white;
  @apply flex items-center gap-2;
  background: var(--gradient-primary);
  box-shadow: var(--shadow-glow-accent);
  transition: all 0.2s ease-out;
}

.cta-button:hover {
  background: var(--gradient-primary-hover);
  box-shadow: 0 8px 32px rgba(139, 92, 246, 0.5), 0 0 48px rgba(236, 72, 153, 0.2);
  transform: translateY(-3px) scale(1.02);
}

.cta-button:active {
  transform: translateY(-1px) scale(1.01);
}

/* Recent Section */
.recent-section {
  @apply mb-4 flex-1 flex flex-col min-h-0;
  animation: fade-up 0.5s ease-out 0.6s both;
}

.section-header {
  @apply flex items-center justify-between mb-4;
}

.section-title {
  @apply text-base font-semibold text-gray-200;
  @apply flex items-center gap-2;
}

.clear-btn {
  @apply text-xs text-gray-400 hover:text-red-400;
  @apply flex items-center gap-1;
  @apply transition-colors duration-200;
}

/* Projects List */
.projects-list {
  @apply space-y-2 overflow-y-auto pr-1 flex-1 min-h-0;
}

.project-card {
  @apply relative p-3 rounded-xl cursor-pointer;
  @apply flex items-center gap-3;
  @apply bg-gray-800/40 border border-gray-700/40;
  @apply transition-all duration-200;
  animation: slide-in 0.4s ease-out both;
}

.project-card:hover {
  @apply bg-gray-700/50 border-gray-600/50;
  transform: translateX(4px);
  box-shadow: var(--shadow-md);
}

.project-card:hover .project-icon {
  @apply bg-purple-500/30 border-purple-400/40;
}

.project-card:hover .project-icon svg {
  @apply text-purple-300;
}

.project-card:hover .project-arrow {
  @apply opacity-100 text-purple-400;
}

.project-card:hover .project-remove {
  @apply opacity-100;
}

@keyframes slide-in {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }

  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.project-icon {
  @apply w-10 h-10 rounded-lg flex-shrink-0;
  @apply flex items-center justify-center;
  @apply bg-gray-700/50 border border-gray-600/30;
  @apply transition-all duration-200;
}

.project-icon svg {
  @apply text-gray-400 transition-colors duration-200;
}

.project-info {
  @apply flex-1 min-w-0;
}

.project-name {
  @apply font-medium text-white mb-0.5;
}

.project-path {
  @apply text-xs text-gray-400 truncate;
  max-width: 100%;
}

.project-remove {
  @apply p-2 rounded-lg;
  @apply text-gray-600 hover:text-red-400 hover:bg-red-500/10;
  @apply opacity-0 transition-all duration-200;
}

.project-arrow {
  @apply w-5 h-5 text-gray-600;
  @apply opacity-0 transition-all duration-200;
}

/* Empty Projects */
.empty-projects {
  @apply py-12 text-center rounded-xl;
  @apply bg-gray-800/30 border border-gray-700/30;
}

.empty-icon {
  @apply w-14 h-14 mx-auto mb-3 rounded-xl;
  @apply flex items-center justify-center;
  @apply bg-gray-700/50 text-gray-400;
}

.empty-projects p {
  @apply text-sm text-gray-400;
}

/* Context Menu */
.context-menu {
  @apply fixed z-50 py-1 min-w-[180px] rounded-xl;
  @apply bg-gray-800 border border-gray-700/50;
  @apply shadow-2xl backdrop-blur-md;
}

.context-menu-item {
  @apply w-full px-3 py-2 text-left text-sm;
  @apply flex items-center gap-2;
  @apply text-gray-300 hover:bg-gray-700/60 hover:text-white;
  @apply transition-colors duration-150;
}

.context-menu-item-danger {
  @apply text-red-400 hover:text-red-300;
}

/* Small toggle for header */
.toggle-wrapper-sm {
  @apply relative inline-flex flex-shrink-0;
}

.toggle-track-sm {
  @apply w-7 h-4 rounded-full;
  @apply bg-gray-600 transition-all duration-200;
}

.toggle-thumb-sm {
  @apply absolute left-0.5 top-0.5 w-3 h-3 rounded-full;
  @apply bg-white shadow-sm;
  @apply transition-transform duration-200;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
