<template>
  <div class="template-section">
    <!-- Header: Template selector + quick access + settings -->
    <div class="template-header">
      <TemplateSelector class="flex-1" />
      <div class="template-quick-actions">
        <button
          v-for="tpl in quickTemplates"
          :key="tpl.id"
          @click="templateStore.setActiveTemplate(tpl.id)"
          class="quick-template-btn"
          :class="{ active: tpl.id === templateStore.activeTemplateId }"
          :title="tpl.name"
        >
          {{ tpl.icon }}
        </button>
        <button
          @click="templateStore.openModal()"
          class="template-settings-btn"
          :title="t('templates.settings')"
        >
          <Settings class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Smart suggestion -->
    <Transition name="fade">
      <div v-if="suggestedTemplate && suggestedTemplate !== templateStore.activeTemplateId" class="template-suggestion">
        <Lightbulb class="w-3.5 h-3.5" />
        <span>{{ t('templates.suggestion') }}</span>
        <button @click="applySuggestion" class="suggestion-btn">
          {{ suggestedTemplateName }}
        </button>
      </div>
    </Transition>

    <!-- Task input -->
    <div class="template-task">
      <div class="task-header">
        <label class="template-label">
          <FileText class="w-3.5 h-3.5" />
          {{ t('templates.task') }}
        </label>
        <div class="task-actions">
          <button 
            v-if="task" 
            @click="clearTask" 
            class="task-action-btn" 
            :title="t('templates.clearTask')"
          >
            <X class="w-3.5 h-3.5" />
          </button>
          <button 
            v-if="templateStore.taskHistory.length > 0" 
            @click="showHistory = !showHistory" 
            class="task-action-btn" 
            :class="{ active: showHistory }"
            :title="t('templates.history')"
          >
            <History class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
      
      <div class="task-input-wrapper">
        <textarea
          v-model="task"
          :placeholder="t('templates.taskPlaceholder')"
          class="template-textarea"
          rows="3"
        />
        <span v-if="task" class="task-char-count">{{ task.length }}</span>
      </div>

      <!-- Task history dropdown -->
      <Transition name="slide">
        <div v-if="showHistory" class="task-history">
          <div class="history-header">
            <span>{{ t('templates.recentTasks') }}</span>
            <button @click="templateStore.clearTaskHistory()" class="clear-history-btn" :title="t('templates.clearHistory')">
              <Trash2 class="w-3 h-3" />
            </button>
          </div>
          <button
            v-for="item in templateStore.taskHistory"
            :key="item.id"
            @click="selectHistoryItem(item)"
            class="history-item"
          >
            <span class="history-text">{{ item.text }}</span>
            <span class="history-time">{{ formatRelativeTime(item.timestamp) }}</span>
          </button>
        </div>
      </Transition>
    </div>

    <!-- User Rules (collapsible) -->
    <div class="expand-section">
      <button @click="showRules = !showRules" class="expand-header">
        <ChevronRight class="w-4 h-4 expand-icon" :class="{ expanded: showRules }" />
        <span>{{ t('templates.userRules') }}</span>
        <span v-if="templateStore.userRules" class="expand-badge">✓</span>
      </button>
      <Transition name="expand">
        <div v-if="showRules" class="expand-content">
          <textarea
            v-model="userRulesLocal"
            :placeholder="t('templates.userRulesPlaceholder')"
            class="template-textarea"
            rows="3"
          />
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { ChevronRight, FileText, History, Lightbulb, Settings, Trash2, X } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useTemplateStore } from '../model/template.store'
import type { TaskHistoryItem } from '../model/template.types'
import TemplateSelector from './TemplateSelector.vue'

const { t } = useI18n()
const templateStore = useTemplateStore()

const showHistory = ref(false)
const showRules = ref(false)

// Quick access templates (favorites or first 3)
const quickTemplates = computed(() => {
  const favs = templateStore.favoriteTemplates
  if (favs.length > 0) return favs.slice(0, 3)
  return templateStore.visibleTemplates.slice(0, 3)
})

// Task with two-way binding
const task = computed({
  get: () => templateStore.currentTask,
  set: (v: string) => templateStore.setTask(v)
})

// User rules with local state
const userRulesLocal = computed({
  get: () => templateStore.userRules,
  set: (v: string) => templateStore.setUserRules(v)
})

// Smart suggestion
const suggestedTemplate = computed(() => templateStore.suggestTemplate(task.value))
const suggestedTemplateName = computed(() => {
  if (!suggestedTemplate.value) return ''
  const tpl = templateStore.templates.find(t => t.id === suggestedTemplate.value)
  return tpl ? `${tpl.icon} ${tpl.name}` : ''
})

function applySuggestion() {
  if (suggestedTemplate.value) {
    templateStore.setActiveTemplate(suggestedTemplate.value)
  }
}

function clearTask() {
  templateStore.setTask('')
}

function selectHistoryItem(item: TaskHistoryItem) {
  templateStore.setTask(item.text)
  showHistory.value = false
}

function formatRelativeTime(timestamp: string): string {
  const now = Date.now()
  const then = new Date(timestamp).getTime()
  const diff = now - then
  const mins = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (mins < 1) return t('common.justNow') || 'just now'
  if (mins < 60) return `${mins}m`
  if (hours < 24) return `${hours}h`
  return `${days}d`
}

// Close history on click outside
watch(showHistory, (v) => {
  if (v) {
    const handler = (e: MouseEvent) => {
      const target = e.target as HTMLElement
      if (!target.closest('.task-history') && !target.closest('.task-action-btn')) {
        showHistory.value = false
        document.removeEventListener('click', handler)
      }
    }
    setTimeout(() => document.addEventListener('click', handler), 0)
  }
})
</script>

<style scoped>
.template-section {
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
}

.template-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.template-quick-actions {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.quick-template-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  background: var(--bg-2);
  border: 1px solid transparent;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.15s;
  opacity: 0.6;
}

.quick-template-btn:hover {
  opacity: 1;
  background: var(--bg-3);
}

.quick-template-btn.active {
  opacity: 1;
  border-color: var(--color-primary);
  background: var(--bg-accent-subtle);
}

.template-settings-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  color: var(--text-2);
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
}

.template-settings-btn:hover {
  background: var(--bg-3);
  color: var(--text-1);
}

/* Suggestion */
.template-suggestion {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.5rem;
  background: var(--bg-accent-subtle);
  border-radius: 0.375rem;
  font-size: 0.6875rem;
  color: var(--text-2);
}

.suggestion-btn {
  padding: 0.125rem 0.375rem;
  background: var(--color-primary);
  border: none;
  border-radius: 0.25rem;
  color: white;
  font-size: 0.6875rem;
  font-weight: 500;
  cursor: pointer;
}

.suggestion-btn:hover {
  opacity: 0.9;
}

/* Task */
.template-task {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  position: relative;
}

.task-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.template-label {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-2);
}

.task-actions {
  display: flex;
  align-items: center;
  gap: 0.125rem;
}

.task-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  background: transparent;
  border: none;
  border-radius: 0.25rem;
  color: var(--text-3);
  cursor: pointer;
  transition: all 0.15s;
}

.task-action-btn:hover {
  background: var(--bg-3);
  color: var(--text-1);
}

.task-action-btn.active {
  background: var(--bg-accent-subtle);
  color: var(--color-primary);
}

.task-input-wrapper {
  position: relative;
}

.template-textarea {
  width: 100%;
  padding: 0.5rem;
  padding-right: 2.5rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  color: var(--text-1);
  font-size: 0.8125rem;
  font-family: inherit;
  resize: vertical;
  min-height: 60px;
  transition: border-color 0.15s;
}

.template-textarea:focus {
  outline: none;
  border-color: var(--border-focus);
}

.template-textarea::placeholder {
  color: var(--text-3);
}

.task-char-count {
  position: absolute;
  bottom: 0.375rem;
  right: 0.5rem;
  font-size: 0.625rem;
  color: var(--text-3);
  pointer-events: none;
}

/* Task history */
.task-history {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 0.25rem;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 20;
  max-height: 200px;
  overflow-y: auto;
}

.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem;
  border-bottom: 1px solid var(--border-default);
  font-size: 0.6875rem;
  color: var(--text-3);
}

.clear-history-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.25rem;
  background: transparent;
  border: none;
  border-radius: 0.25rem;
  color: var(--text-3);
  cursor: pointer;
}

.clear-history-btn:hover {
  background: var(--color-danger-soft);
  color: var(--color-danger);
}

.history-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.5rem;
  background: transparent;
  border: none;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s;
}

.history-item:hover {
  background: var(--bg-3);
}

.history-text {
  flex: 1;
  font-size: 0.75rem;
  color: var(--text-1);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.history-time {
  font-size: 0.625rem;
  color: var(--text-3);
  flex-shrink: 0;
}

/* Expandable section */
.expand-section {
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  overflow: hidden;
}

.expand-header {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  width: 100%;
  padding: 0.5rem;
  background: var(--bg-2);
  border: none;
  color: var(--text-2);
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.15s;
}

.expand-header:hover {
  background: var(--bg-3);
  color: var(--text-1);
}

.expand-icon {
  transition: transform 0.2s;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.expand-badge {
  margin-left: auto;
  color: var(--color-success);
  font-size: 0.625rem;
}

.expand-content {
  padding: 0.5rem;
  background: var(--bg-1);
  border-top: 1px solid var(--border-default);
}

/* Animations */
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.slide-enter-active, .slide-leave-active { transition: all 0.2s; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-4px); }

.expand-enter-active, .expand-leave-active { transition: all 0.2s; }
.expand-enter-from, .expand-leave-to { opacity: 0; max-height: 0; padding: 0 0.5rem; }
</style>
