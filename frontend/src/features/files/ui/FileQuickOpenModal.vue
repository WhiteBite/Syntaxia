<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="fileStore.isQuickOpenModalVisible" class="quick-open-overlay" @click="handleClose" @keydown.esc="handleClose">
        <div class="quick-open-modal" @click.stop>
          <input
            ref="inputRef"
            v-model="query"
            class="quick-open-input"
            :placeholder="t('files.quickOpen.placeholder')"
            @keydown="handleKeydown"
            autofocus
          />
          <div v-if="results.length > 0" class="quick-open-results">
            <div
              v-for="(file, idx) in results"
              :key="file.path"
              :class="['quick-open-item', { 'quick-open-item--active': idx === activeIndex }]"
              @click="selectFile(file)"
              @mouseenter="activeIndex = idx"
            >
              <span class="quick-open-icon">{{ getFileIcon(file.name) }}</span>
              <div class="quick-open-content">
                <span class="quick-open-name">
                  <template v-if="highlightedNames[idx] && highlightedNames[idx].length > 1">
                    <template v-for="(segment, segIdx) in highlightedNames[idx]" :key="segIdx">
                      <mark v-if="segment.isMatch" class="quick-open-highlight">{{ segment.text }}</mark>
                      <span v-else>{{ segment.text }}</span>
                    </template>
                  </template>
                  <template v-else>
                    {{ file.name }}
                  </template>
                </span>
                <span class="quick-open-path">{{ getRelativePath(file.path) }}</span>
              </div>
              <span v-if="isSelected(file.path)" class="quick-open-badge">✓</span>
            </div>
          </div>
          <div v-else-if="query" class="quick-open-empty">
            <span class="quick-open-empty-text">{{ t('files.quickOpen.noResults') }}</span>
          </div>
          <div v-else class="quick-open-hint">
            <span class="quick-open-hint-text">{{ t('files.quickOpen.hint') }}</span>
          </div>
          <div class="quick-open-footer">
            <span class="quick-open-shortcut">↑↓ {{ t('files.quickOpen.navigate') }}</span>
            <span class="quick-open-shortcut">Enter {{ t('files.quickOpen.select') }}</span>
            <span class="quick-open-shortcut">Esc {{ t('files.quickOpen.close') }}</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { getFileIcon } from '@/utils/fileIcons'
import { highlightFuzzy, type TextSegment } from '@/utils/searchHighlight'
import { useFileStore, type FileNode } from '../model/file.store'
import { computed, nextTick, ref, watch } from 'vue'

const { t } = useI18n()
const fileStore = useFileStore()

const query = ref('')
const activeIndex = ref(0)
const inputRef = ref<HTMLInputElement>()

const results = computed(() => {
  if (!query.value) return []
  
  fileStore.setSearchQuery(query.value)
  return fileStore.searchResults.slice(0, 20)
})

// Compute highlighted names for all results
const highlightedNames = computed<TextSegment[][]>(() => {
  if (!query.value) return []
  
  return results.value.map(file => highlightFuzzy(file.name, query.value))
})

watch(results, () => {
  activeIndex.value = 0
})

// Watch modal visibility to focus input and reset state
watch(() => fileStore.isQuickOpenModalVisible, (visible) => {
  if (visible) {
    query.value = ''
    activeIndex.value = 0
    nextTick(() => {
      inputRef.value?.focus()
    })
  } else {
    // Clear search when closing
    fileStore.setSearchQuery('')
  }
})

function getRelativePath(path: string): string {
  const rootPath = fileStore.rootPath
  if (!rootPath) return path
  return path.replace(rootPath + '/', '')
}

function isSelected(path: string): boolean {
  return fileStore.selectedPaths.has(path)
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    handleClose()
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = Math.min(activeIndex.value + 1, results.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = Math.max(activeIndex.value - 1, 0)
  } else if (e.key === 'Enter' && results.value[activeIndex.value]) {
    e.preventDefault()
    selectFile(results.value[activeIndex.value])
  }
}

function selectFile(file: FileNode) {
  // Toggle selection in file store
  fileStore.toggleSelect(file.path)
  
  // Expand parent folders to make file visible in tree
  const parts = file.path.split('/')
  for (let i = 1; i < parts.length; i++) {
    const folderPath = parts.slice(0, i).join('/')
    fileStore.expandPath(folderPath)
  }
  
  handleClose()
}

function handleClose() {
  fileStore.closeQuickOpenModal()
}
</script>

<style scoped>
.quick-open-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  z-index: 9999;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 20vh;
}

.quick-open-modal {
  width: min(600px, 90vw);
  max-height: 70vh;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.quick-open-input {
  width: 100%;
  padding: 1rem 1.25rem;
  font-size: 1.125rem;
  border: none;
  background: transparent;
  color: var(--text-primary);
  outline: none;
  border-bottom: 1px solid var(--border-subtle);
}

.quick-open-input::placeholder {
  color: var(--text-subtle);
}

.quick-open-results {
  flex: 1;
  overflow-y: auto;
  max-height: 400px;
}

.quick-open-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1.25rem;
  cursor: pointer;
  transition: background 100ms;
  border-left: 3px solid transparent;
}

.quick-open-item:hover,
.quick-open-item--active {
  background: var(--bg-2);
  border-left-color: var(--accent-indigo);
}

.quick-open-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.quick-open-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.quick-open-name {
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-open-highlight {
  background: rgba(168, 85, 247, 0.4);
  color: white;
  padding: 0 2px;
  border-radius: 2px;
  font-weight: 700;
}

.quick-open-path {
  font-size: 0.875rem;
  color: var(--text-subtle);
  font-family: ui-monospace, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-open-badge {
  color: var(--accent-indigo);
  font-weight: 600;
  flex-shrink: 0;
}

.quick-open-empty,
.quick-open-hint {
  padding: 3rem 1.25rem;
  text-align: center;
}

.quick-open-empty-text,
.quick-open-hint-text {
  color: var(--text-subtle);
  font-size: 0.875rem;
}

.quick-open-footer {
  display: flex;
  gap: 1rem;
  padding: 0.75rem 1.25rem;
  border-top: 1px solid var(--border-subtle);
  background: var(--bg-0);
}

.quick-open-shortcut {
  font-size: 0.75rem;
  color: var(--text-subtle);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Transitions */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 200ms ease-out;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-active .quick-open-modal,
.modal-fade-leave-active .quick-open-modal {
  transition: transform 200ms ease-out, opacity 200ms ease-out;
}

.modal-fade-enter-from .quick-open-modal,
.modal-fade-leave-to .quick-open-modal {
  transform: translateY(-20px);
  opacity: 0;
}
</style>
