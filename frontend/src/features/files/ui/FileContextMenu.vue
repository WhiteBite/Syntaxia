<template>
  <Teleport to="body">
    <Transition name="context-menu">
      <div
        v-if="visible && node"
        role="menu"
        class="context-menu fixed z-[1060] bg-gray-800 border-2 border-gray-600 rounded-lg shadow-2xl py-1 min-w-[220px]"
        :style="{ left: `${position.x}px`, top: `${position.y}px` }"
        @click.stop
      >
        <!-- Select Actions -->
        <button
          v-if="node.isDir"
          @click="handleAction('selectAll')"
          role="menuitem"
          :aria-label="t('contextMenu.selectAll')"
          class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 hover:scale-[1.01] active:scale-[0.99] flex items-center gap-3 transition-all duration-150"
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
          class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 hover:scale-[1.01] active:scale-[0.99] flex items-center gap-3 transition-all duration-150"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ t('contextMenu.deselectAll') }}
        </button>

        <div v-if="node.isDir" class="h-px bg-gray-700 my-1"></div>

        <!-- Copy Actions -->
        <!-- QuickLook (files only) -->
        <button
          v-if="!node.isDir"
          @click="handleAction('quickLook')"
          class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
          {{ t('files.quickLook') }}
          <span class="ml-auto text-xs text-gray-400">Space</span>
        </button>

        <div v-if="!node.isDir" class="h-px bg-gray-700 my-1"></div>

        <!-- Copy Actions -->
        <button
          @click="handleAction('copyPath')"
          class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          {{ t('contextMenu.copyPath') }}
        </button>

        <button
          @click="handleAction('copyRelativePath')"
          class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          {{ t('contextMenu.copyRelativePath') }}
        </button>

        <div class="h-px bg-gray-700 my-1"></div>

        <!-- Ignore Actions -->
        <button
          v-if="!node.isIgnored"
          @click="handleAction('addToCustomIgnore')"
          class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
          </svg>
          {{ t('contextMenu.addToIgnore') }}
        </button>

        <button
          v-if="node.isIgnored"
          @click="handleAction('removeFromIgnore')"
          class="w-full px-4 py-2 text-left text-sm text-orange-400 hover:bg-gray-700 flex items-center gap-3 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
          {{ t('contextMenu.removeFromIgnore') }}
        </button>

        <!-- Expand/Collapse Actions -->
        <template v-if="node.isDir">
          <div class="h-px bg-gray-700 my-1"></div>

          <button
            @click="handleAction('expandAll')"
            class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
            {{ t('contextMenu.expandAll') }}
          </button>

          <button
            @click="handleAction('collapseAll')"
            class="w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
            </svg>
            {{ t('contextMenu.collapseAll') }}
          </button>
        </template>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { FileNode } from '@/features/files/model/file.store'

const { t } = useI18n()

interface Props {
  node: FileNode | null
  position: { x: number; y: number }
  visible: boolean
}

const props = defineProps<Props>()

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
</script>

<style scoped>
.context-menu-enter-active {
  transition: all 0.2s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

.context-menu-enter-from {
  opacity: 0;
  transform: scale(0.9) translateY(-8px);
}

.context-menu-leave-active {
  transition: all 0.15s ease-out;
}

.context-menu-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
