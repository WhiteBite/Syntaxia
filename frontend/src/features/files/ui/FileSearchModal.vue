<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="modal-overlay file-search-overlay"
        @click.self="quickOpen.close()"
        @keydown="quickOpen.handleKeyDown"
      >
        <div class="file-search-modal">
          <!-- Search Input -->
          <div class="search-header">
            <MagnifyingGlassIcon class="search-icon" />
            <input
              ref="inputRef"
              v-model="quickOpen.query.value"
              type="text"
              class="search-input"
              :placeholder="t('fileSearch.placeholder')"
              autocomplete="off"
              spellcheck="false"
            />
            <kbd class="search-hint">ESC</kbd>
          </div>

          <!-- Results List -->
          <div class="search-results" v-if="quickOpen.results.value.length > 0">
            <div
              v-for="(result, index) in quickOpen.results.value"
              :key="result.item.path"
              :class="['search-result-item', { 'is-selected': index === quickOpen.selectedIndex.value }]"
              @click="quickOpen.selectItem(index)"
              @mouseenter="quickOpen.selectedIndex.value = index"
            >
              <span class="file-icon">{{ getFileIcon(result.item.name) }}</span>
              <div class="file-info">
                <span class="file-name" v-html="highlightMatches(result.item.name, result.matches, 'name')"></span>
                <span class="file-path" v-html="highlightMatches(result.item.path, result.matches, 'path')"></span>
              </div>
              <span v-if="isFileSelected(result.item.path)" class="selected-badge">
                {{ t('fileSearch.inContext') }}
              </span>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else-if="quickOpen.query.value && !quickOpen.isSearching.value" class="search-empty">
            <span>{{ t('fileSearch.noResults') }}</span>
          </div>

          <!-- Initial State -->
          <div v-else-if="!quickOpen.query.value" class="search-empty search-hint-text">
            <span>{{ t('fileSearch.hint') }}</span>
          </div>

          <!-- Loading -->
          <div v-else-if="quickOpen.isSearching.value" class="search-empty">
            <span>{{ t('common.loading') }}</span>
          </div>

          <!-- Footer -->
          <div class="search-footer">
            <div class="footer-hint">
              <kbd>↑↓</kbd>
              <span>{{ t('fileSearch.navigate') }}</span>
            </div>
            <div class="footer-hint">
              <kbd>Enter</kbd>
              <span>{{ t('fileSearch.select') }}</span>
            </div>
            <div class="footer-hint">
              <kbd>Esc</kbd>
              <span>{{ t('fileSearch.close') }}</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useFileStore } from '@/features/files/model/file.store'
import { getFileIcon } from '@/utils/fileIcons'
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline'
import type { FuseResultMatch } from 'fuse.js'
import { nextTick, watch } from 'vue'
import { useFileQuickOpen } from '../composables/useFileQuickOpen'

const props = defineProps<{
  isOpen: boolean
}>()

const { t } = useI18n()
const fileStore = useFileStore()
const quickOpen = useFileQuickOpen()

// Focus input when modal opens
watch(() => props.isOpen, async (open) => {
  if (open) {
    quickOpen.reset()
    await nextTick()
    quickOpen.focusInput()
  }
})

// Check if file is already selected
function isFileSelected(path: string): boolean {
  return fileStore.selectedPaths.has(path)
}

// Highlight matched characters in text
function highlightMatches(
  text: string,
  matches: ReadonlyArray<FuseResultMatch> | undefined,
  key: string
): string {
  if (!matches) return escapeHtml(text)

  const match = matches.find(m => m.key === key)
  if (!match || !match.indices || match.indices.length === 0) {
    return escapeHtml(text)
  }

  const indices = [...match.indices].sort((a, b) => a[0] - b[0])
  let result = ''
  let lastIndex = 0

  for (const [start, end] of indices) {
    if (start > lastIndex) {
      result += escapeHtml(text.slice(lastIndex, start))
    }
    result += `<mark class="search-highlight">${escapeHtml(text.slice(start, end + 1))}</mark>`
    lastIndex = end + 1
  }

  if (lastIndex < text.length) {
    result += escapeHtml(text.slice(lastIndex))
  }

  return result
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}
</script>

<style scoped>
.file-search-overlay { @apply flex items-start justify-center pt-[15vh]; }

.file-search-modal {
  @apply w-full flex flex-col rounded-xl overflow-hidden;
  max-width: min(600px, 90vw);
  max-height: 60vh;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.search-header {
  @apply flex items-center gap-3 px-4 py-3;
  border-bottom: 1px solid var(--border-default);
}

.search-icon { @apply w-5 h-5 flex-shrink-0; color: var(--text-muted); }

.search-input {
  @apply flex-1 bg-transparent border-none outline-none text-base;
  color: var(--text-primary);
}

.search-input::placeholder { color: var(--text-muted); }

.search-hint {
  @apply px-1.5 py-0.5 text-xs rounded;
  background: var(--bg-3);
  color: var(--text-muted);
  font-family: inherit;
}

.search-results { @apply flex-1 overflow-y-auto py-2; min-height: 0; }

.search-result-item {
  @apply flex items-center gap-3 px-4 py-2 cursor-pointer;
  transition: background-color 150ms ease;
}

.search-result-item:hover,
.search-result-item.is-selected { background: var(--bg-2); }

.file-icon { @apply text-base flex-shrink-0; }

.file-info { @apply flex-1 min-w-0 flex flex-col; }

.file-name { @apply text-sm font-medium truncate; color: var(--text-primary); }

.file-path { @apply text-xs truncate; color: var(--text-muted); }

.selected-badge {
  @apply px-2 py-0.5 text-xs rounded-full flex-shrink-0;
  background: var(--accent-indigo-muted);
  color: var(--accent-indigo);
}

.search-empty { @apply flex items-center justify-center py-8 text-sm; color: var(--text-muted); }

.search-hint-text { @apply text-center px-4; }

.search-footer {
  @apply flex items-center justify-center gap-6 px-4 py-2;
  border-top: 1px solid var(--border-default);
  background: var(--bg-2);
}

.footer-hint { @apply flex items-center gap-1.5 text-xs; color: var(--text-muted); }

.footer-hint kbd { @apply px-1.5 py-0.5 rounded; background: var(--bg-3); font-family: inherit; }

:deep(.search-highlight) {
  @apply px-0.5 rounded;
  background: var(--accent-amber-muted);
  color: var(--accent-amber);
}

.modal-enter-active, .modal-leave-active { transition: opacity 150ms ease; }
.modal-enter-active .file-search-modal, .modal-leave-active .file-search-modal {
  transition: transform 150ms ease, opacity 150ms ease;
}
.modal-enter-from, .modal-leave-to { opacity: 0; }
.modal-enter-from .file-search-modal, .modal-leave-to .file-search-modal {
  transform: scale(0.95) translateY(-10px);
  opacity: 0;
}
</style>
