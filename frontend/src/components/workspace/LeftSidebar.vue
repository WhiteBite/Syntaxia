<template>
  <div class="sidebar-container">
    <!-- Tab Switcher with Sliding Indicator -->
    <div class="tabs-container">
      <!-- Sliding Indicator -->
      <div 
        class="tabs-indicator"
        :class="tabIndicatorClass"
        :style="{ transform: `translateX(${tabIndex * 100}%)` }"
      ></div>
      
      <button
        @click="currentTab = 'files'"
        :class="['sidebar-tab', currentTab === 'files' ? 'sidebar-tab-active text-indigo-300' : '']"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        <span>{{ t('tabs.files') }}</span>
        <span v-if="fileStore.selectedPaths.size > 0" class="tab-badge">
          {{ fileStore.selectedPaths.size }}
        </span>
      </button>
      
      <button
        @click="currentTab = 'git'"
        :class="['sidebar-tab', currentTab === 'git' ? 'sidebar-tab-active text-orange-300' : '']"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
        </svg>
        <span>Git</span>
      </button>
      
      <button
        @click="currentTab = 'contexts'"
        :class="['sidebar-tab', currentTab === 'contexts' ? 'sidebar-tab-active text-purple-300' : '']"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <span>{{ t('tabs.contexts') }}</span>
      </button>
    </div>

    <!-- Content - using v-show to preserve component state when switching tabs -->
    <div class="sidebar-content">
      <FileExplorer 
        v-show="currentTab === 'files'" 
        class="sidebar-tab-content"
        :class="{ 'sidebar-tab-content--active': currentTab === 'files' }"
        @preview-file="handlePreviewFile" 
        @build-context="handleBuildContext" 
      />
      <GitSourceSelector 
        v-show="currentTab === 'git'" 
        class="sidebar-tab-content"
        :class="{ 'sidebar-tab-content--active': currentTab === 'git' }"
      />
      <ContextList 
        v-show="currentTab === 'contexts'" 
        class="sidebar-tab-content"
        :class="{ 'sidebar-tab-content--active': currentTab === 'contexts' }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
// Direct imports to avoid circular dependency with barrel exports
import { useFileStore } from '@/features/files/model/file.store'
import FileExplorer from '@/features/files/ui/FileExplorer.vue'
import ContextList from '@/features/context/ui/ContextList.vue'
import GitSourceSelector from '@/features/git/ui/GitSourceSelector.vue'
import { computed, ref, watch } from 'vue'

const fileStore = useFileStore()
const { t } = useI18n()
const currentTab = ref<'files' | 'git' | 'contexts'>('files')

const tabIndex = computed(() => {
  const tabs = ['files', 'git', 'contexts']
  return tabs.indexOf(currentTab.value)
})

const tabIndicatorClass = computed(() => {
  const classes: Record<string, string> = {
    files: 'tabs-indicator-indigo',
    git: 'tabs-indicator-orange',
    contexts: 'tabs-indicator-purple'
  }
  return classes[currentTab.value]
})

const emit = defineEmits<{
  (e: 'preview-file', filePath: string): void
  (e: 'build-context'): void
}>()

function handlePreviewFile(filePath: string) {
  emit('preview-file', filePath)
}

function handleBuildContext() {
  emit('build-context')
}

// Persist tab selection
watch(currentTab, (tab) => {
  try {
    localStorage.setItem('left-sidebar-tab', tab)
  } catch (err) {
    console.warn('Failed to save sidebar tab:', err)
  }
})

// Restore tab selection
try {
  const savedTab = localStorage.getItem('left-sidebar-tab')
  if (savedTab === 'files' || savedTab === 'contexts' || savedTab === 'git') {
    currentTab.value = savedTab
  }
} catch (err) {
  console.warn('Failed to load sidebar tab:', err)
}
</script>

<style scoped>
.sidebar-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background: transparent;
}

.tabs-container {
  position: relative;
  display: flex;
  gap: 4px;
  padding: 8px 8px 0;
  background: var(--bg-1);
}

/* Sliding Indicator - Bottom Glow Style */
.tabs-indicator {
  position: absolute;
  bottom: 0;
  height: 2px;
  width: calc(33.333% - 8px);
  left: 8px;
  border-radius: 2px 2px 0 0;
  transition: transform 250ms cubic-bezier(0.4, 0, 0.2, 1),
              background 200ms ease-out,
              box-shadow 200ms ease-out;
  z-index: 1;
  pointer-events: none;
}

.tabs-indicator-indigo {
  background: #6366f1;
  box-shadow: 0 0 12px rgba(99, 102, 241, 0.6);
}

.tabs-indicator-orange {
  background: #f97316;
  box-shadow: 0 0 12px rgba(249, 115, 22, 0.6);
}

.tabs-indicator-purple {
  background: #a855f7;
  box-shadow: 0 0 12px rgba(168, 85, 247, 0.6);
}

.tabs-indicator-cyan {
  background: #06b6d4;
  box-shadow: 0 0 12px rgba(6, 182, 212, 0.6);
}

.tabs-indicator-emerald {
  background: #10b981;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.6);
}

.sidebar-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 12px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 8px 8px 0 0;
  position: relative;
  z-index: 0;
  color: var(--text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: color 150ms ease-out, background 150ms ease-out;
}

.sidebar-tab:hover:not(.sidebar-tab-active) {
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.03);
}

.sidebar-tab:active {
  transform: scale(0.98);
}

.sidebar-tab-active {
  font-weight: 600;
  color: white;
  background: linear-gradient(180deg, transparent 0%, rgba(255, 255, 255, 0.03) 100%);
}

.tab-badge {
  @apply px-1.5 py-0.5 text-xs font-bold rounded-full;
  background: var(--accent-indigo-bg);
  color: white;
  border: 1px solid var(--accent-indigo-border);
}

/* Content area - fills remaining space */
.sidebar-content {
  flex: 1 1 0;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* Tab content - hidden by v-show, active one fills container */
.sidebar-tab-content {
  flex: 0 0 0;
  min-height: 0;
  overflow: hidden;
}

/* Active tab content fills the container */
.sidebar-tab-content--active {
  flex: 1 1 0;
}
</style>
