<template>
  <Teleport to="body">
    <Transition name="dropdown">
      <div
        v-if="isOpen"
        ref="dropdownRef"
        class="mention-dropdown"
        :style="dropdownStyle"
        role="listbox"
        :aria-label="t('mentions.title')"
        @keydown="handleKeydown"
      >
        <!-- Header with categories -->
        <div class="mention-dropdown-header">
          <span class="mention-dropdown-title">{{ t('mentions.title') }}</span>
          <div class="mention-categories">
            <button
              v-for="category in categories"
              :key="category.id"
              class="category-btn"
              :class="{ active: activeCategory === category.id }"
              :title="category.label"
              @click="setCategory(category.id)"
            >
              <span :class="category.icon" />
            </button>
          </div>
        </div>

        <!-- Loading state -->
        <div v-if="isLoading" class="mention-loading">
          <span class="i-lucide-loader-2 animate-spin" />
          <span>{{ t('common.loading') }}</span>
        </div>

        <!-- Suggestions list -->
        <div v-else-if="filteredSuggestions.length > 0" class="mention-suggestions">
          <button
            v-for="(suggestion, index) in filteredSuggestions"
            :key="`${suggestion.type}-${suggestion.value}`"
            class="mention-item"
            :class="{ selected: index === selectedIndex }"
            role="option"
            :aria-selected="index === selectedIndex"
            @click="selectSuggestion(suggestion)"
            @mouseenter="selectedIndex = index"
          >
            <span class="mention-icon" :class="suggestion.icon" />
            <div class="mention-content">
              <span class="mention-display">{{ suggestion.display }}</span>
              <span v-if="suggestion.description" class="mention-description">
                {{ suggestion.description }}
              </span>
            </div>
            <span class="mention-type-badge" :class="`type-${suggestion.type}`">
              {{ suggestion.type }}
            </span>
          </button>
        </div>

        <!-- Empty state -->
        <div v-else class="mention-empty">
          <span class="i-lucide-search" />
          <span>{{ t('mentions.noResults') }}</span>
        </div>

        <!-- Footer with hints -->
        <div class="mention-dropdown-footer">
          <span class="hint">
            <kbd>↑↓</kbd> {{ t('mentions.navigate') }}
          </span>
          <span class="hint">
            <kbd>Enter</kbd> {{ t('mentions.select') }}
          </span>
          <span class="hint">
            <kbd>Esc</kbd> {{ t('mentions.close') }}
          </span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import type { MentionCategory, MentionSuggestion, MentionType } from '../types/mentions'

const props = defineProps<{
  isOpen: boolean
  isLoading: boolean
  suggestions: MentionSuggestion[]
  selectedIndex: number
  position: { top: number; left: number }
  activeCategory: MentionType | 'all'
  categories: MentionCategory[]
}>()

const emit = defineEmits<{
  select: [suggestion: MentionSuggestion]
  close: []
  navigate: [direction: 'up' | 'down']
  categoryChange: [category: MentionType | 'all']
}>()

const { t } = useI18n()
const dropdownRef = ref<HTMLElement | null>(null)
const selectedIndex = ref(props.selectedIndex)

// Sync selectedIndex with prop
watch(() => props.selectedIndex, (val) => {
  selectedIndex.value = val
})

/** Computed dropdown style */
const dropdownStyle = computed(() => ({
  position: 'fixed' as const,
  top: `${props.position.top}px`,
  left: `${props.position.left}px`,
  zIndex: '1200',
}))

/** Filter suggestions by active category */
const filteredSuggestions = computed(() => {
  if (props.activeCategory === 'all') {
    return props.suggestions
  }
  return props.suggestions.filter(s => s.type === props.activeCategory)
})

/** Select a suggestion */
function selectSuggestion(suggestion: MentionSuggestion) {
  emit('select', suggestion)
}

/** Set active category */
function setCategory(category: MentionType | 'all') {
  emit('categoryChange', category)
}

/** Handle keyboard navigation */
function handleKeydown(e: KeyboardEvent) {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      emit('navigate', 'down')
      break
    case 'ArrowUp':
      e.preventDefault()
      emit('navigate', 'up')
      break
    case 'Enter':
      e.preventDefault()
      if (filteredSuggestions.value[selectedIndex.value]) {
        selectSuggestion(filteredSuggestions.value[selectedIndex.value])
      }
      break
    case 'Escape':
      e.preventDefault()
      emit('close')
      break
    case 'Tab':
      emit('close')
      break
  }
}

/** Handle click outside */
function handleClickOutside(e: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.mention-dropdown {
  @apply rounded-xl overflow-hidden;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  box-shadow: var(--shadow-xl);
  min-width: min(320px, 80vw);
  max-width: min(400px, 90vw);
  max-height: min(400px, 50vh);
  display: flex;
  flex-direction: column;
}

.mention-dropdown-header {
  @apply flex items-center justify-between px-3 py-2;
  background: var(--bg-2);
  border-bottom: 1px solid var(--border-default);
}

.mention-dropdown-title {
  @apply text-xs font-semibold;
  color: var(--text-primary);
}

.mention-categories {
  @apply flex gap-1;
}

.category-btn {
  @apply p-1.5 rounded-md transition-all;
  color: var(--text-muted);
}

.category-btn:hover {
  background: var(--bg-3);
  color: var(--text-primary);
}

.category-btn.active {
  background: var(--accent-indigo-bg);
  color: var(--accent-indigo);
}

.mention-loading {
  @apply flex items-center justify-center gap-2 py-8;
  color: var(--text-muted);
}

.mention-suggestions {
  @apply py-1 overflow-y-auto;
  flex: 1;
  min-height: 0;
}

.mention-item {
  @apply flex items-center gap-3 px-3 py-2 w-full text-left transition-all;
  color: var(--text-primary);
}

.mention-item:hover,
.mention-item.selected {
  background: var(--bg-2);
}

.mention-item.selected {
  background: var(--accent-indigo-bg);
}

.mention-icon {
  @apply text-base flex-shrink-0;
  color: var(--text-muted);
}

.mention-item.selected .mention-icon {
  color: var(--accent-indigo);
}

.mention-content {
  @apply flex flex-col flex-1 min-w-0;
}

.mention-display {
  @apply text-sm font-medium truncate;
}

.mention-description {
  @apply text-xs truncate;
  color: var(--text-muted);
}

.mention-type-badge {
  @apply text-[10px] px-1.5 py-0.5 rounded font-medium flex-shrink-0;
  background: var(--bg-3);
  color: var(--text-muted);
}

.mention-type-badge.type-file {
  background: var(--accent-blue-bg);
  color: var(--accent-blue);
}

.mention-type-badge.type-folder {
  background: var(--accent-amber-bg);
  color: var(--accent-amber);
}

.mention-type-badge.type-symbol {
  background: var(--accent-purple-bg);
  color: var(--accent-purple);
}

.mention-type-badge.type-git {
  background: var(--accent-orange-bg);
  color: var(--accent-orange);
}

.mention-type-badge.type-docs {
  background: var(--accent-emerald-bg);
  color: var(--accent-emerald);
}

.mention-empty {
  @apply flex flex-col items-center justify-center gap-2 py-8;
  color: var(--text-muted);
}

.mention-dropdown-footer {
  @apply flex items-center justify-center gap-4 px-3 py-2 text-[10px];
  background: var(--bg-2);
  border-top: 1px solid var(--border-default);
  color: var(--text-muted);
}

.hint {
  @apply flex items-center gap-1;
}

.hint kbd {
  @apply px-1 py-0.5 rounded text-[9px] font-mono;
  background: var(--bg-3);
  border: 1px solid var(--border-default);
}

/* Transitions */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 150ms ease-out;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
