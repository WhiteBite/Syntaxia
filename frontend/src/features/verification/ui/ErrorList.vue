<template>
  <div class="error-list">
    <!-- Empty State -->
    <div v-if="errors.length === 0" class="error-list-empty">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none">
        <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <p class="empty-text">{{ t('verification.noIssues') }}</p>
    </div>

    <!-- Grouped by File -->
    <template v-else-if="groupBy === 'file'">
      <div 
        v-for="[file, fileErrors] in errorsByFile" 
        :key="file"
        class="error-group"
      >
        <!-- File Header -->
        <button 
          class="error-group-header"
          @click="toggleGroup(file)"
        >
          <svg 
            class="group-chevron" 
            :class="{ 'group-chevron--expanded': expandedGroups.has(file) }"
            viewBox="0 0 24 24" 
            fill="none"
          >
            <path d="M9 5l7 7-7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          
          <svg class="group-icon" viewBox="0 0 24 24" fill="none">
            <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          
          <span class="group-file">{{ getFileName(file) }}</span>
          <span class="group-path">{{ getFilePath(file) }}</span>
          
          <div class="group-badges">
            <span v-if="getErrorCount(fileErrors)" class="group-badge group-badge--error">
              {{ getErrorCount(fileErrors) }}
            </span>
            <span v-if="getWarningCount(fileErrors)" class="group-badge group-badge--warning">
              {{ getWarningCount(fileErrors) }}
            </span>
          </div>
        </button>

        <!-- File Errors -->
        <div 
          v-if="expandedGroups.has(file)"
          class="error-group-content"
        >
          <ErrorItem
            v-for="error in fileErrors"
            :key="error.id"
            :error="error"
            :is-fixing="isFixing === error.id"
            @go-to-file="$emit('go-to-file', $event)"
            @fix-with-ai="$emit('fix-with-ai', $event)"
          />
        </div>
      </div>
    </template>

    <!-- Grouped by Severity -->
    <template v-else>
      <div 
        v-for="severity in severityOrder" 
        :key="severity"
        class="error-group"
      >
        <template v-if="errorsBySeverity[severity].length > 0">
          <!-- Severity Header -->
          <button 
            class="error-group-header"
            @click="toggleGroup(severity)"
          >
            <svg 
              class="group-chevron" 
              :class="{ 'group-chevron--expanded': expandedGroups.has(severity) }"
              viewBox="0 0 24 24" 
              fill="none"
            >
              <path d="M9 5l7 7-7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            
            <span 
              class="severity-dot"
              :class="[`severity-dot--${severity}`]"
            ></span>
            
            <span class="group-title">{{ t(`verification.${severity}`) }}</span>
            
            <span class="group-count">{{ errorsBySeverity[severity].length }}</span>
          </button>

          <!-- Severity Errors -->
          <div 
            v-if="expandedGroups.has(severity)"
            class="error-group-content"
          >
            <ErrorItem
              v-for="error in errorsBySeverity[severity]"
              :key="error.id"
              :error="error"
              :is-fixing="isFixing === error.id"
              @go-to-file="$emit('go-to-file', $event)"
              @fix-with-ai="$emit('fix-with-ai', $event)"
            />
          </div>
        </template>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, ref, watch } from 'vue'
import type { Severity, VerificationError } from '../api/verification.api'
import ErrorItem from './ErrorItem.vue'

const { t } = useI18n()

const props = defineProps<{
  errors: VerificationError[]
  groupBy: 'file' | 'severity'
  isFixing: string | null
}>()

defineEmits<{
  (e: 'go-to-file', error: VerificationError): void
  (e: 'fix-with-ai', errorId: string): void
}>()

const severityOrder: Severity[] = ['error', 'warning', 'info']
const expandedGroups = ref<Set<string>>(new Set())

// Group errors by file
const errorsByFile = computed(() => {
  const grouped = new Map<string, VerificationError[]>()
  
  for (const error of props.errors) {
    const existing = grouped.get(error.file) || []
    existing.push(error)
    grouped.set(error.file, existing)
  }
  
  return grouped
})

// Group errors by severity
const errorsBySeverity = computed(() => {
  const grouped: Record<Severity, VerificationError[]> = {
    error: [],
    warning: [],
    info: []
  }
  
  for (const error of props.errors) {
    grouped[error.severity].push(error)
  }
  
  return grouped
})

// Auto-expand groups when errors change
watch(() => props.errors, () => {
  if (props.groupBy === 'file') {
    for (const file of errorsByFile.value.keys()) {
      expandedGroups.value.add(file)
    }
  } else {
    for (const severity of severityOrder) {
      if (errorsBySeverity.value[severity].length > 0) {
        expandedGroups.value.add(severity)
      }
    }
  }
}, { immediate: true })

function toggleGroup(key: string): void {
  if (expandedGroups.value.has(key)) {
    expandedGroups.value.delete(key)
  } else {
    expandedGroups.value.add(key)
  }
}

function getFileName(path: string): string {
  const parts = path.split(/[/\\]/)
  return parts[parts.length - 1] || path
}

function getFilePath(path: string): string {
  const parts = path.split(/[/\\]/)
  if (parts.length <= 1) return ''
  return parts.slice(0, -1).join('/')
}

function getErrorCount(errors: VerificationError[]): number {
  return errors.filter(e => e.severity === 'error').length
}

function getWarningCount(errors: VerificationError[]): number {
  return errors.filter(e => e.severity === 'warning').length
}

// Expose for parent
function expandAll(): void {
  if (props.groupBy === 'file') {
    for (const file of errorsByFile.value.keys()) {
      expandedGroups.value.add(file)
    }
  } else {
    for (const severity of severityOrder) {
      expandedGroups.value.add(severity)
    }
  }
}

function collapseAll(): void {
  expandedGroups.value.clear()
}

defineExpose({ expandAll, collapseAll })
</script>

<style scoped>
.error-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* Empty State */
.error-list-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
}

.empty-icon {
  width: 48px;
  height: 48px;
  color: #22c55e;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 14px;
  color: #9ca3af;
  margin: 0;
}

/* Group */
.error-group {
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  overflow: hidden;
}

.error-group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 12px 16px;
  background: none;
  border: none;
  cursor: pointer;
  transition: background 0.15s ease-out;
}

.error-group-header:hover {
  background: rgba(255, 255, 255, 0.04);
}

.group-chevron {
  width: 16px;
  height: 16px;
  color: #6b7280;
  transition: transform 0.15s ease-out;
  flex-shrink: 0;
}

.group-chevron--expanded {
  transform: rotate(90deg);
}

.group-icon {
  width: 16px;
  height: 16px;
  color: #9ca3af;
  flex-shrink: 0;
}

.group-file {
  font-size: 13px;
  font-weight: 600;
  color: #e5e7eb;
}

.group-path {
  font-size: 12px;
  color: #6b7280;
  flex: 1;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-title {
  font-size: 13px;
  font-weight: 600;
  color: #e5e7eb;
  text-transform: capitalize;
}

.group-count {
  font-size: 12px;
  color: #6b7280;
  margin-left: auto;
}

.group-badges {
  display: flex;
  gap: 4px;
  margin-left: auto;
}

.group-badge {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
}

.group-badge--error {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

.group-badge--warning {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
}

.severity-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.severity-dot--error {
  background: #ef4444;
}

.severity-dot--warning {
  background: #f59e0b;
}

.severity-dot--info {
  background: #3b82f6;
}

.error-group-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}
</style>
