<template>
  <div v-if="favorites.length > 0 || !collapsed" class="favorites-panel">
    <!-- Header -->
    <button 
      class="favorites-header"
      @click="collapsed = !collapsed"
      :aria-expanded="!collapsed"
    >
      <StarIcon class="w-4 h-4 text-yellow-400" />
      <span class="favorites-title">{{ t('files.favorites.title') }}</span>
      <span class="favorites-count">{{ favorites.length }}</span>
      <ChevronDownIcon 
        class="w-4 h-4 text-gray-400 transition-transform" 
        :class="{ 'rotate-180': collapsed }"
      />
    </button>

    <!-- Content -->
    <Transition name="collapse">
      <div v-if="!collapsed" class="favorites-content">
        <div v-if="favorites.length === 0" class="favorites-empty">
          {{ t('files.favorites.empty') }}
        </div>
        
        <div v-else class="favorites-list">
          <div
            v-for="path in favorites"
            :key="path"
            class="favorite-item group"
            @click="handleClick(path)"
            :title="path"
          >
            <span class="favorite-icon">{{ getFileIcon(getFileName(path)) }}</span>
            <span class="favorite-name">{{ getFileName(path) }}</span>
            <button
              class="favorite-remove"
              @click.stop="handleRemove(path)"
              :title="t('files.favorites.remove')"
            >
              <XMarkIcon class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { getFileIcon } from '@/utils/fileIcons'
import { ChevronDownIcon, StarIcon, XMarkIcon } from '@heroicons/vue/24/solid'
import { ref } from 'vue'
import { useFavorites } from '../composables/useFavorites'
import { useFileStore } from '../model/file.store'

const { t } = useI18n()
const fileStore = useFileStore()
const { favorites, removeFavorite } = useFavorites()

const collapsed = ref(false)

const emit = defineEmits<{
  (e: 'select', path: string): void
}>()

function getFileName(path: string): string {
  return path.split('/').pop() || path
}

function handleClick(path: string) {
  // Add to context selection
  if (!fileStore.selectedPaths.has(path)) {
    fileStore.toggleSelect(path)
  }
  emit('select', path)
}

function handleRemove(path: string) {
  removeFavorite(path)
}
</script>

<style scoped>
.favorites-panel {
  border-bottom: 1px solid var(--border-default);
  background: rgba(255, 255, 255, 0.02);
}

.favorites-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.5rem 0.75rem;
  text-align: left;
  background: transparent;
  border: none;
  cursor: pointer;
  transition: background 0.15s;
}

.favorites-header:hover {
  background: rgba(255, 255, 255, 0.05);
}

.favorites-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.favorites-count {
  font-size: 0.65rem;
  padding: 0.125rem 0.375rem;
  background: rgba(234, 179, 8, 0.2);
  color: rgb(234, 179, 8);
  border-radius: 9999px;
  margin-left: auto;
}

.favorites-content {
  padding: 0.25rem 0.5rem 0.5rem;
}

.favorites-empty {
  padding: 0.75rem;
  text-align: center;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.favorites-list {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.favorite-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.5rem;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: background 0.15s;
}

.favorite-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.favorite-icon {
  font-size: 0.875rem;
  flex-shrink: 0;
}

.favorite-name {
  flex: 1;
  font-size: 0.8125rem;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.favorite-remove {
  opacity: 0;
  padding: 0.25rem;
  border-radius: 0.25rem;
  background: transparent;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
}

.favorite-item:hover .favorite-remove {
  opacity: 1;
}

.favorite-remove:hover {
  background: rgba(239, 68, 68, 0.2);
  color: rgb(239, 68, 68);
}

/* Collapse animation */
.collapse-enter-active,
.collapse-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}

.collapse-enter-from,
.collapse-leave-to {
  opacity: 0;
  max-height: 0;
}

.collapse-enter-to,
.collapse-leave-from {
  opacity: 1;
  max-height: 20rem;
}
</style>
