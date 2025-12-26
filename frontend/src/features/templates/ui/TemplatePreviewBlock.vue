<template>
  <!-- In-Editor Template Block - Sticky HUD -->
  <div 
    v-if="templateStore.activeTemplate && contextStore.hasContext && settingsStore.settings.context.applyTemplateOnCopy"
    class="tpl-block"
    :class="{ expanded: isExpanded }"
    :style="{ '--role-color': roleColor }"
  >
    <!-- Collapsed Header (42px) -->
    <div class="tpl-block-header" @click="toggleExpand">
      <span class="tpl-block-icon">{{ templateStore.activeTemplate.icon }}</span>
      <span class="tpl-block-name">{{ templateStore.activeTemplate.name }}</span>
      <span class="tpl-block-tokens">{{ tokenCount.toLocaleString() }} tok</span>
      <span class="tpl-block-separator">|</span>
      <span class="tpl-block-preview">{{ truncatedRole }}</span>
      <div class="tpl-block-actions">
        <button @click.stop="openEditor" class="tpl-block-btn" :title="t('templates.settings')">
          <Pencil class="w-3 h-3" />
        </button>
      </div>
      <ChevronDown class="tpl-block-chevron" :class="{ rotated: isExpanded }" />
    </div>

    <!-- Expanded Content -->
    <Transition name="slide">
      <div v-if="isExpanded" class="tpl-block-content">
        <div v-if="templateStore.activeTemplate.sections.role && templateStore.activeTemplate.roleContent" class="tpl-block-section">
          <span class="tpl-block-label">ROLE</span>
          <pre class="tpl-block-text">{{ templateStore.activeTemplate.roleContent }}</pre>
        </div>
        <div v-if="templateStore.activeTemplate.sections.rules && templateStore.activeTemplate.rulesContent" class="tpl-block-section">
          <span class="tpl-block-label">RULES</span>
          <pre class="tpl-block-text">{{ templateStore.activeTemplate.rulesContent }}</pre>
        </div>
        <div v-if="templateStore.activeTemplate.sections.tree" class="tpl-block-section">
          <span class="tpl-block-label">TREE</span>
          <span class="tpl-block-meta">{{ contextStore.fileCount }} files</span>
        </div>
        <div v-if="templateStore.activeTemplate.sections.stats" class="tpl-block-section">
          <span class="tpl-block-label">STATS</span>
          <span class="tpl-block-meta">{{ contextStore.tokenCount }} tokens</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context'
import { useTemplateStore } from '../model/template.store'
import { useSettingsStore } from '@/stores/settings.store'
import { ChevronDown, Pencil } from 'lucide-vue-next'
import { computed, ref } from 'vue'

const { t } = useI18n()
const templateStore = useTemplateStore()
const contextStore = useContextStore()
const settingsStore = useSettingsStore()

const isExpanded = ref(false)

// Role color based on template
const roleColor = computed(() => {
  const id = templateStore.activeTemplate?.id || ''
  const colors: Record<string, string> = {
    architect: '#f59e0b', // amber
    review: '#a855f7',    // purple
    implement: '#3b82f6', // blue
    explain: '#10b981'    // emerald
  }
  return colors[id] || '#6366f1' // default indigo
})

const tokenCount = computed(() => {
  const tpl = templateStore.activeTemplate
  if (!tpl) return 0
  let count = 0
  if (tpl.roleContent) count += Math.round(tpl.roleContent.length / 4)
  if (tpl.rulesContent) count += Math.round(tpl.rulesContent.length / 4)
  count += contextStore.tokenCount
  return count
})

const truncatedRole = computed(() => {
  const role = templateStore.activeTemplate?.roleContent || ''
  return role.length > 50 ? role.slice(0, 50) + '...' : role
})

function toggleExpand() {
  isExpanded.value = !isExpanded.value
}

function openEditor() {
  templateStore.openModal()
}
</script>

<style scoped>
/* Block container with left accent border */
.tpl-block {
  position: sticky;
  top: 0;
  z-index: 10;
  background: rgba(13, 17, 23, 0.97);
  backdrop-filter: blur(12px);
  border-left: 3px solid var(--role-color);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  /* Gradient shadow fade below */
  box-shadow: 0 4px 12px -2px rgba(0, 0, 0, 0.4);
  transition: box-shadow 0.2s ease-out;
}

.tpl-block:hover {
  box-shadow: 0 6px 16px -2px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.08);
}

/* Collapsed header - exactly 42px */
.tpl-block-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  height: 42px;
  padding: 0 0.75rem;
  cursor: pointer;
  user-select: none;
}

.tpl-block-icon {
  font-size: 1rem;
  line-height: 1;
  flex-shrink: 0;
}

.tpl-block-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-primary);
  flex-shrink: 0;
}

/* Token pill */
.tpl-block-tokens {
  padding: 0.125rem 0.5rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 9999px;
  font-size: 10px;
  font-family: var(--font-mono);
  color: var(--text-muted);
  flex-shrink: 0;
}

/* Separator */
.tpl-block-separator {
  color: rgba(255, 255, 255, 0.15);
  font-size: 12px;
  flex-shrink: 0;
}

/* Preview text - truncated */
.tpl-block-preview {
  flex: 1;
  min-width: 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Actions (edit button) */
.tpl-block-actions {
  display: flex;
  align-items: center;
  opacity: 0;
  transition: opacity 0.15s;
  flex-shrink: 0;
}

.tpl-block:hover .tpl-block-actions {
  opacity: 1;
}

.tpl-block-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: var(--text-subtle);
  cursor: pointer;
  transition: all 0.15s;
}

.tpl-block-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

/* Chevron - always visible */
.tpl-block-chevron {
  width: 14px;
  height: 14px;
  color: var(--text-subtle);
  flex-shrink: 0;
  transition: transform 0.2s ease-out, color 0.15s;
}

.tpl-block:hover .tpl-block-chevron {
  color: var(--text-primary);
}

.tpl-block-chevron.rotated {
  transform: rotate(180deg);
}

/* Expanded Content */
.tpl-block-content {
  padding: 0 0.75rem 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  border-top: 1px dashed rgba(255, 255, 255, 0.1);
}

.tpl-block-section {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.tpl-block-label {
  flex-shrink: 0;
  width: 3rem;
  padding: 0.125rem 0;
  font-size: 9px;
  font-weight: 600;
  color: var(--text-subtle);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.tpl-block-text {
  flex: 1;
  margin: 0;
  padding: 0;
  font-size: 11px;
  font-family: var(--font-mono);
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.5;
  max-height: 100px;
  overflow-y: auto;
}

.tpl-block-meta {
  font-size: 11px;
  font-family: var(--font-mono);
  color: var(--text-muted);
}

/* Slide animation */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
}

.slide-enter-to,
.slide-leave-from {
  opacity: 1;
  max-height: 300px;
}
</style>
