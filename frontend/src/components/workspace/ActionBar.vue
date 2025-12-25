<template>
  <div class="action-bar">
    <!-- Left: Project & Stats -->
    <div class="action-bar-left">
      <!-- Project Name -->
      <button @click="changeProject" class="project-btn" :title="t('hotkey.changeProject')">
        <div class="project-icon">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
        </div>
        <span class="project-name">{{ projectStore.projectName }}</span>
        <svg class="w-3 h-3 text-gray-400 group-hover:text-gray-300" fill="none" stroke="currentColor"
          viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      <!-- Quick Template Switch -->
      <div class="template-quick-switch">
        <button
          v-for="tpl in quickTemplates"
          :key="tpl.id"
          @click="templateStore.setActiveTemplate(tpl.id)"
          class="quick-tpl-btn"
          :class="{ active: tpl.id === templateStore.activeTemplateId }"
          :title="tpl.name"
        >
          {{ tpl.icon }}
        </button>
      </div>
    </div>

    <!-- Separator -->
    <div class="action-bar-separator"></div>

    <!-- Right: Actions -->
    <div class="action-bar-right">
      <!-- Settings -->
      <button @click="openSettings" class="toolbar-btn" :title="t('settings.modal.title') + ' (Ctrl+,)'">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      </button>

      <!-- Reset Layout -->
      <button @click="handleResetLayout" class="toolbar-btn" :title="t('workspace.resetLayout')">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>

      <!-- Language -->
      <button @click="toggleLanguage" class="toolbar-btn"
        :title="locale === 'ru' ? 'Switch to English' : 'Переключить на русский'">
        {{ locale.toUpperCase() }}
      </button>
    </div>
  </div>
</template>


<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useTemplateStore } from '@/features/templates'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed } from 'vue'

const projectStore = useProjectStore()
const templateStore = useTemplateStore()
const uiStore = useUIStore()
const { t, locale, setLocale } = useI18n()

// Quick access templates (favorites or first 4)
const quickTemplates = computed(() => {
  const favs = templateStore.favoriteTemplates
  if (favs.length > 0) return favs.slice(0, 4)
  return templateStore.visibleTemplates.slice(0, 4)
})

function toggleLanguage() {
  const newLocale = locale.value === 'ru' ? 'en' : 'ru'
  setLocale(newLocale)
  uiStore.addToast(newLocale === 'ru' ? 'Язык изменён на русский' : 'Language changed to English', 'success')
}

function changeProject() {
  projectStore.clearProject()
  uiStore.addToast(t('action.changeProject'), 'info')
}

function openSettings() {
  uiStore.openSettingsModal()
}

const emit = defineEmits<{
  (e: 'open-export'): void
  (e: 'reset-layout'): void
}>()

async function handleResetLayout() {
  // Reset panel sizes
  emit('reset-layout')
  
  // Reset window size via backend
  try {
    // @ts-ignore - Wails bindings
    if (window.go?.main?.App?.ResetWindowState) {
      await window.go.main.App.ResetWindowState()
    }
  } catch (error) {
    console.warn('Failed to reset window state:', error)
  }
}
</script>

<style scoped>
.action-bar {
  @apply h-10 flex items-center justify-between px-3 gap-2;
}

.action-bar-left {
  @apply flex items-center gap-4 min-w-0 flex-shrink;
}

.action-bar-right {
  @apply flex items-center gap-2 flex-shrink-0;
}

/* Separator */
.action-bar-separator {
  @apply flex-1;
}

/* Project Button */
.project-btn {
  @apply flex items-center gap-2 px-3 py-1.5 rounded-xl;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  transition: all 200ms ease-out;
}

.project-btn:hover {
  background: var(--bg-2);
  border-color: var(--border-strong);
}

.project-icon {
  @apply w-6 h-6 rounded-lg flex items-center justify-center;
  background: var(--accent-purple-bg);
  color: #d8b4fe;
}

.project-name {
  @apply text-xs font-medium truncate max-w-[140px];
  color: var(--text-primary);
}

/* Quick Template Switch */
.template-quick-switch {
  @apply flex items-center gap-1 px-2 py-1 rounded-lg;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
}

.quick-tpl-btn {
  @apply flex items-center justify-center w-7 h-7 rounded-md text-sm;
  background: transparent;
  border: 1px solid transparent;
  cursor: pointer;
  opacity: 0.6;
  transition: all 0.15s;
}

.quick-tpl-btn:hover {
  opacity: 1;
  background: var(--bg-2);
}

.quick-tpl-btn.active {
  opacity: 1;
  background: var(--bg-accent-subtle);
  border-color: var(--color-primary);
}

/* Toolbar Button */
.toolbar-btn {
  @apply px-2.5 py-1.5 rounded-lg text-xs font-medium;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  color: var(--text-muted);
  transition: all 150ms ease-out;
}

.toolbar-btn:hover {
  color: var(--text-primary);
  background: var(--bg-2);
  border-color: var(--border-strong);
}
</style>
