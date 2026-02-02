<template>
  <BaseDropdown
    :model-value="visible"
    placement="bottom-start"
    @update:model-value="handleVisibilityChange"
  >
    <template #trigger>
      <!-- No visible trigger - controlled externally via context menu -->
      <span></span>
    </template>

    <div v-if="node" class="context-menu-content">
      <!-- Select Actions -->
      <button
        v-if="node.isDir"
        @click="handleAction('selectAll')"
        role="menuitem"
        :aria-label="t('contextMenu.selectAll')"
        class="context-menu-item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ t('contextMenu.selectAll') }}
      </button>

      <button
        v-if="node.isDir"
        @click="handleAction('deselectAll')"
        role="menuitem"
        :aria-label="t('contextMenu.deselectAll')"
        class="context-menu-item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ t('contextMenu.deselectAll') }}
      </button>

      <div v-if="node.isDir" class="context-menu-divider"></div>

      <!-- Focus on Folder (folders only) -->
      <button
        v-if="node.isDir"
        @click="handleAction('focusOnFolder')"
        class="context-menu-item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="3" stroke-width="2"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12c0-4.5 4-8 9-8s9 3.5 9 8-4 8-9 8-9-3.5-9-8z"/>
        </svg>
        {{ t('files.focusOnFolder') }}
      </button>

      <!-- QuickLook (files only) -->
      <button
        v-if="!node.isDir"
        @click="handleAction('quickLook')"
        class="context-menu-item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        {{ t('files.quickLook') }}
        <span class="ml-auto text-xs text-gray-400">Space</span>
      </button>

      <!-- Favorites (files only) -->
      <button
        v-if="!node.isDir"
        @click="handleAction(isNodeFavorite ? 'removeFromFavorites' : 'addToFavorites')"
        class="context-menu-item"
      >
        <StarIconSolid v-if="isNodeFavorite" class="w-4 h-4 text-yellow-400" />
        <StarIconOutline v-else class="w-4 h-4" />
        {{ isNodeFavorite ? t('files.favorites.remove') : t('files.favorites.add') }}
      </button>

      <!-- Explain code (files only) -->
      <button
        v-if="!node.isDir"
        @click="handleAction('explainCode')"
        class="context-menu-item"
      >
        <LightBulbIcon class="w-4 h-4 text-yellow-400" />
        {{ t('chat.explain.menuItem') }}
      </button>

      <div v-if="!node.isDir" class="context-menu-divider"></div>

      <!-- Copy Actions -->
      <button
        @click="handleAction('copyPath')"
        class="context-menu-item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        {{ t('contextMenu.copyPath') }}
      </button>

      <button
        @click="handleAction('copyRelativePath')"
        class="context-menu-item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        {{ t('contextMenu.copyRelativePath') }}
      </button>

      <div class="context-menu-divider"></div>

      <!-- Ignore Actions -->
      <button
        v-if="!node.isIgnored"
        @click="handleAction('addToCustomIgnore')"
        class="context-menu-item"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
        </svg>
        {{ t('contextMenu.addToIgnore') }}
      </button>

      <button
        v-if="node.isIgnored"
        @click="handleAction('removeFromIgnore')"
        class="context-menu-item context-menu-item-danger"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        {{ t('contextMenu.removeFromIgnore') }}
      </button>

      <!-- Expand/Collapse Actions -->
      <template v-if="node.isDir">
        <div class="context-menu-divider"></div>

        <button
          @click="handleAction('expandAll')"
          class="context-menu-item"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
          {{ t('contextMenu.expandAll') }}
        </button>

        <button
          @click="handleAction('collapseAll')"
          class="context-menu-item"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
          {{ t('contextMenu.collapseAll') }}
        </button>
      </template>
    </div>
  </BaseDropdown>
</template>

<script setup lang="ts">
import BaseDropdown from '@/components/ui/BaseDropdown.vue'
import { useI18n } from '@/composables/useI18n'
import type { FileNode } from '@/features/files/model/file.store'
import { LightBulbIcon, StarIcon as StarIconOutline } from '@heroicons/vue/24/outline'
import { StarIcon as StarIconSolid } from '@heroicons/vue/24/solid'
import { computed } from 'vue'
import { useFavorites } from '../composables/useFavorites'

const { t } = useI18n()
const { isFavorite } = useFavorites()

interface Props {
  node: FileNode | null
  position: { x: number; y: number }
  visible: boolean
}

const props = defineProps<Props>()

const isNodeFavorite = computed(() => {
  if (!props.node || props.node.isDir) return false
  return isFavorite(props.node.path)
})

const emit = defineEmits<{
  (e: 'action', payload: { type: string; node: FileNode }): void
  (e: 'close'): void
}>()

function handleAction(type: string) {
  if (props.node) {
    emit('action', { type, node: props.node })
  }
  emit('close')
}

function handleVisibilityChange(visible: boolean) {
  if (!visible) {
    emit('close')
  }
}
</script>

<style scoped>
.context-menu-content {
  min-width: 220px;
  padding: 0.25rem 0;
}

.context-menu-item {
  @apply w-full px-4 py-2 text-left text-sm flex items-center gap-3 transition-all duration-150;
  color: white;
}

.context-menu-item:hover {
  background: var(--bg-3);
  transform: scale(1.01);
}

.context-menu-item:active {
  transform: scale(0.99);
}

.context-menu-item-danger {
  color: #fb923c;
}

.context-menu-divider {
  height: 1px;
  background: var(--border-default);
  margin: 0.25rem 0;
}
</style>
