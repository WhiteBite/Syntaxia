<template>
  <label 
    class="smart-hud-item"
    :class="{ 'smart-hud-item-selected': isSelected }"
  >
    <!-- Checkbox -->
    <div class="smart-hud-checkbox" :class="{ checked: isSelected }">
      <input 
        type="checkbox" 
        :checked="isSelected"
        @change="$emit('toggle')"
      />
      <svg v-if="isSelected" viewBox="0 0 12 12" fill="none">
        <path d="M2 6L5 9L10 3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </div>

    <!-- File Icon (SVG) -->
    <div class="smart-hud-file-icon" :class="iconClass">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z"/>
        <path d="M14 2v6h6M10 12l-2 2 2 2M14 12l2 2-2 2"/>
      </svg>
    </div>

    <!-- File Info -->
    <div class="smart-hud-file-info">
      <span class="smart-hud-file-name">{{ fileName }}</span>
      <span class="smart-hud-file-path">{{ filePath }}</span>
    </div>

    <!-- Source Badge -->
    <span class="smart-hud-badge" :class="badgeClass">
      <span class="smart-hud-badge-icon">🔗</span>
      {{ sourceLabel }}
    </span>
  </label>
</template>

<script setup lang="ts">
defineProps<{
  isSelected: boolean
  iconClass: string
  fileName: string
  filePath: string
  badgeClass: string
  sourceLabel: string
}>()

defineEmits<{
  (e: 'toggle'): void
}>()
</script>

<style scoped>
.smart-hud-item {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.625rem 0.75rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: background 100ms;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.smart-hud-item:last-child {
  border-bottom: none;
}

.smart-hud-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.smart-hud-item-selected {
  background: rgba(139, 92, 246, 0.1);
}

.smart-hud-item-selected:hover {
  background: rgba(139, 92, 246, 0.15);
}

/* Checkbox */
.smart-hud-checkbox {
  position: relative;
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
  border: 1.5px solid rgba(255, 255, 255, 0.15);
  border-radius: 0.25rem;
  background: rgba(0, 0, 0, 0.3);
  transition: all 100ms;
  display: flex;
  align-items: center;
  justify-content: center;
}

.smart-hud-checkbox.checked {
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  border-color: transparent;
}

.smart-hud-checkbox input {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.smart-hud-checkbox svg {
  width: 0.625rem;
  height: 0.625rem;
  color: white;
}

/* File Icon */
.smart-hud-file-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.smart-hud-file-icon svg {
  width: 1rem;
  height: 1rem;
}

/* File icon colors */
.smart-hud-file-icon.icon-vue { color: #6ee7b7; }
.smart-hud-file-icon.icon-ts { color: #7dd3fc; }
.smart-hud-file-icon.icon-js { color: #fde047; }
.smart-hud-file-icon.icon-go { color: #67e8f9; }
.smart-hud-file-icon.icon-py { color: #7dd3fc; }
.smart-hud-file-icon.icon-rs { color: #fdba74; }
.smart-hud-file-icon.icon-json { color: #fde047; }
.smart-hud-file-icon.icon-css { color: #c4b5fd; }
.smart-hud-file-icon.icon-html { color: #fca5a5; }
.smart-hud-file-icon.icon-default { color: #d1d5db; }

/* File Info */
.smart-hud-file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.smart-hud-file-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #e5e7eb;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.3;
}

.smart-hud-file-path {
  font-size: 0.6875rem;
  color: #a1a1aa;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.3;
}

/* Source Badges */
.smart-hud-badge {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.5625rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  border: none;
}

.smart-hud-badge-icon {
  font-size: 0.5rem;
  opacity: 0.8;
}

.badge-git {
  color: #93c5fd;
  background: rgba(59, 130, 246, 0.12);
}

.badge-arch {
  color: #86efac;
  background: rgba(34, 197, 94, 0.12);
}

.badge-semantic {
  color: #d8b4fe;
  background: rgba(168, 85, 247, 0.12);
}

.badge-default {
  color: #a1a1aa;
  background: rgba(107, 114, 128, 0.12);
}
</style>
