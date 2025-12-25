<template>
  <div
    class="context-row group"
    :class="{ 
      'context-row--selected': isSelected,
      'context-row--active': isActive 
    }"
    :draggable="!searchQuery && !showFavoritesOnly"
    @click="$emit('select', context.id)"
    @dblclick="$emit('load', context.id)"
    @dragstart="$emit('drag-start', $event, index)"
    @dragover="$emit('drag-over', $event, index)"
    @dragleave="$emit('drag-leave')"
    @drop="$emit('drop', $event, index)"
    @dragend="$emit('drag-end')"
  >
    <!-- 1. AVATAR / CHECKBOX -->
    <div class="context-avatar-wrap">
      <!-- Avatar (always visible, dimmed on hover) -->
      <div 
        class="context-avatar"
        :class="[`context-avatar--${avatarColor}`]"
      >
        <span v-if="context.isFavorite" class="context-avatar-star">⭐</span>
        <span v-else class="context-avatar-initials">{{ initials }}</span>
      </div>
      
      <!-- Checkmark overlay (appears on hover or when selected) -->
      <div 
        class="context-check-overlay" 
        :class="{ 
          'context-check-overlay--visible': isSelected,
          'context-check-overlay--checked': isSelected 
        }"
        @click.stop="$emit('toggle-select', context.id)"
      >
        <!-- Checkmark icon when selected -->
        <svg v-if="isSelected" class="context-checkmark" viewBox="0 0 24 24" fill="none">
          <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <!-- Empty checkbox when hovering but not selected -->
        <div v-else class="context-checkbox-empty"></div>
      </div>
    </div>

    <!-- 2. INFO -->
    <div class="context-info">
      <!-- Name row -->
      <div class="context-name-row">
        <input
          v-if="isEditing"
          :value="editingName"
          @input="$emit('update:editing-name', ($event.target as HTMLInputElement).value)"
          @blur="$emit('save-rename', context.id)"
          @keyup.enter="$emit('save-rename', context.id)"
          @keyup.escape="$emit('cancel-rename')"
          @click.stop
          class="context-rename-input"
        />
        <span v-else class="context-name">{{ cleanName(context.name) || generateAutoName(context) }}</span>
        <span class="context-time">{{ formatTime(context.createdAt || '') }}</span>
      </div>
      
      <!-- Meta badges -->
      <div class="context-meta">
        <span class="context-badge">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          {{ context.fileCount }}
        </span>
        <span class="context-dot">•</span>
        <span class="context-size">{{ formatSize(context.totalSize) }}</span>
      </div>
    </div>

    <!-- 3. ACTIONS (hover only) -->
    <div class="context-actions">
      <!-- Restore Selection - главная фича -->
      <button 
        v-if="hasFiles"
        @click.stop="$emit('restore-selection', context)" 
        class="context-action context-action--restore"
        :title="t('context.restoreSelection')"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
      <button 
        @click.stop="$emit('toggle-favorite', context.id)" 
        class="context-action"
        :class="{ 'context-action--favorite': context.isFavorite }"
        :title="context.isFavorite ? t('context.unfavorite') : t('context.favorite')"
      >
        <svg class="w-4 h-4" :fill="context.isFavorite ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
        </svg>
      </button>
      <button @click.stop="$emit('copy', context.id)" class="context-action" :title="t('context.copyToClipboard')">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </button>
      <button @click.stop="$emit('start-rename', context)" class="context-action" :title="t('context.rename')">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
        </svg>
      </button>
      <button @click.stop="$emit('delete', context)" class="context-action context-action--danger" :title="t('context.delete')">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'
import { formatContextSize } from '../lib/context-utils'
import type { ContextSummary } from '../model/context.store'

const { t } = useI18n()

const props = defineProps<{
  context: ContextSummary
  index: number
  isSelected: boolean
  isActive: boolean
  isEditing: boolean
  editingName: string
  dragOverIndex: number | null
  searchQuery: string
  showFavoritesOnly: boolean
}>()

defineEmits<{
  (e: 'select', id: string): void
  (e: 'toggle-select', id: string): void
  (e: 'toggle-favorite', id: string): void
  (e: 'load', id: string): void
  (e: 'restore-selection', context: ContextSummary): void
  (e: 'start-rename', context: ContextSummary): void
  (e: 'save-rename', id: string): void
  (e: 'cancel-rename'): void
  (e: 'update:editing-name', value: string): void
  (e: 'copy', id: string): void
  (e: 'duplicate', id: string): void
  (e: 'export', id: string): void
  (e: 'delete', context: ContextSummary): void
  (e: 'drag-start', event: DragEvent, index: number): void
  (e: 'drag-over', event: DragEvent, index: number): void
  (e: 'drag-leave'): void
  (e: 'drop', event: DragEvent, index: number): void
  (e: 'drag-end'): void
}>()

// Check if context has files for restore
const hasFiles = computed(() => {
  return props.context.files && props.context.files.length > 0
})

// Generate initials from name
const initials = computed(() => {
  const name = props.context.name || ''
  
  // If no name or generic name, generate from fileCount or use emoji
  if (!name) {
    return props.context.fileCount > 0 ? String(props.context.fileCount) : '?'
  }
  
  // "Пустой контекст" -> "ПК"
  if (name.includes('Пустой')) {
    return 'ПК'
  }
  
  // "63 файлов" -> show number
  const filesMatch = name.match(/^(\d+)\s*файл/)
  if (filesMatch) {
    return filesMatch[1].length > 2 ? filesMatch[1].slice(0, 2) : filesMatch[1]
  }
  
  // Normal name: "Backend Config" -> "BC", "Vue компоненты" -> "VК"
  const words = name.split(/[\s\-_]+/).filter(w => w.length > 0)
  if (words.length >= 2) {
    return (words[0][0] + words[1][0]).toUpperCase()
  }
  
  // Single word: take first 2 chars
  return name.slice(0, 2).toUpperCase()
})

// Generate color based on ID hash (always unique)
const avatarColor = computed(() => {
  const colors = ['indigo', 'purple', 'pink', 'blue', 'cyan', 'teal', 'green', 'amber']
  // Use ID for hash to ensure different colors even with same names
  const source = props.context.id || props.context.name || ''
  let hash = 0
  for (let i = 0; i < source.length; i++) {
    hash = source.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
})

function formatSize(bytes: number): string {
  return formatContextSize(bytes)
}

function formatTime(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const mins = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (mins < 1) return 'сейчас'
  if (mins < 60) return `${mins}м`
  if (hours < 24) return `${hours}ч`
  if (days < 7) return `${days}д`
  return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
}

function generateAutoName(ctx: ContextSummary): string {
  if (ctx.fileCount > 0) {
    return `${ctx.fileCount} файлов`
  }
  return 'Без названия'
}

// Remove [N] suffix from context names
function cleanName(name: string | undefined): string {
  if (!name) return ''
  return name.replace(/\s*\[\d+\]\s*$/, '').trim()
}
</script>

<style scoped>
/* Row */
.context-row {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  margin: 0 8px 4px 8px;
  border-radius: 12px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.context-row:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.06);
}

.context-row--selected {
  background: rgba(139, 92, 246, 0.08);
  border-color: rgba(139, 92, 246, 0.2);
}

.context-row--active {
  background: rgba(139, 92, 246, 0.15);
  border-color: rgba(139, 92, 246, 0.1);
  border-left: 3px solid #8b5cf6;
  padding-left: 9px;
  border-radius: 0 12px 12px 0;
  margin-left: 8px;
  box-shadow: inset 0 0 0 1px rgba(139, 92, 246, 0.08);
}

.context-row--active .context-name {
  color: white;
}

/* Avatar */
.context-avatar-wrap {
  position: relative;
  width: 40px;
  height: 40px;
  flex-shrink: 0;
}

.context-avatar {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 800;
  color: white;
  transition: all 0.15s ease-out;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.context-row:hover .context-avatar {
  filter: brightness(0.7);
}

.context-avatar-initials {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  letter-spacing: 0.5px;
}

.context-avatar-star {
  font-size: 16px;
}

/* Avatar colors */
.context-avatar--indigo { background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%); }
.context-avatar--purple { background: linear-gradient(135deg, #a855f7 0%, #9333ea 100%); }
.context-avatar--pink { background: linear-gradient(135deg, #ec4899 0%, #db2777 100%); }
.context-avatar--blue { background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); }
.context-avatar--cyan { background: linear-gradient(135deg, #06b6d4 0%, #0891b2 100%); }
.context-avatar--teal { background: linear-gradient(135deg, #14b8a6 0%, #0d9488 100%); }
.context-avatar--green { background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%); }
.context-avatar--orange { background: linear-gradient(135deg, #f97316 0%, #ea580c 100%); }
.context-avatar--amber { background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); }

/* Check Overlay */
.context-check-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  opacity: 0;
  cursor: pointer;
  transition: all 0.15s ease-out;
  background: rgba(0, 0, 0, 0.5);
}

.context-row:hover .context-check-overlay {
  opacity: 1;
}

.context-check-overlay--visible {
  opacity: 1;
}

.context-check-overlay--checked {
  background: rgba(139, 92, 246, 0.85);
  box-shadow: 0 2px 8px rgba(139, 92, 246, 0.4);
}

.context-checkmark {
  width: 22px;
  height: 22px;
  color: white;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

.context-checkbox-empty {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.5);
  border-radius: 5px;
  transition: all 0.15s ease-out;
}

.context-check-overlay:hover .context-checkbox-empty {
  border-color: white;
  background: rgba(255, 255, 255, 0.1);
}

/* Info */
.context-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.context-name-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.context-name {
  font-size: 13px;
  font-weight: 600;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.15s ease-out;
}

.context-row:hover .context-name {
  color: white;
}

.context-time {
  font-size: 10px;
  font-family: ui-monospace, monospace;
  color: #4b5563;
  flex-shrink: 0;
}

.context-rename-input {
  flex: 1;
  font-size: 13px;
  padding: 2px 6px;
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid #8b5cf6;
  border-radius: 4px;
  color: white;
  outline: none;
}

/* Meta */
.context-meta {
  display: flex;
  align-items: center;
  gap: 6px;
}

.context-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px;
  font-size: 10px;
  font-weight: 500;
  color: #9ca3af;
}

.context-dot {
  font-size: 10px;
  color: #4b5563;
}

.context-size {
  font-size: 10px;
  color: #6b7280;
}

/* Actions */
.context-actions {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 4px 6px;
  background: #1a1d2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  opacity: 0;
  transform: translateX(8px);
  transition: all 0.15s ease-out;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.context-row:hover .context-actions {
  opacity: 1;
  transform: translateX(0);
}

.context-action {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: none;
  border: none;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.context-action:hover {
  background: rgba(255, 255, 255, 0.08);
  color: white;
}

.context-action--restore {
  color: #22c55e;
}

.context-action--restore:hover {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

.context-action--favorite {
  color: #facc15;
}

.context-action--danger:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}
</style>
