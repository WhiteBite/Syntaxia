<template>
  <aside class="tpl-preview">
    <div class="tpl-preview-header">
      <Eye class="w-3.5 h-3.5" />
      <span>{{ t('templates.preview') }}</span>
    </div>
    <div class="tpl-preview-content" v-html="highlightedPreview"></div>
    <div class="tpl-preview-footer">
      <div class="tpl-token-bar">
        <div 
          class="tpl-token-fill" 
          :style="{ width: tokenPercent + '%' }" 
          :class="tokenBarClass" 
        />
      </div>
      <div class="tpl-token-info">
        <span class="tpl-token-count">{{ animatedTokens.toLocaleString() }}</span>
        <span class="tpl-token-label">/ 32k tokens</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context'
import { Eye } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import type { PromptTemplate } from '../model/template.types'

const { t } = useI18n()
const contextStore = useContextStore()

const props = defineProps<{
  editingTemplate: PromptTemplate | null
}>()

const animatedTokens = ref(0)

const previewTokens = computed(() => {
  if (!props.editingTemplate) return 0
  let count = 0
  const tpl = props.editingTemplate
  if (tpl.roleContent) count += Math.round(tpl.roleContent.length / 4)
  if (tpl.rulesContent) count += Math.round(tpl.rulesContent.length / 4)
  if (tpl.customPrefix) count += Math.round(tpl.customPrefix.length / 4)
  if (tpl.customSuffix) count += Math.round(tpl.customSuffix.length / 4)
  count += contextStore.tokenCount || 0
  return count
})

const tokenPercent = computed(() => Math.min(100, (previewTokens.value / 32000) * 100))

const tokenBarClass = computed(() => {
  if (tokenPercent.value > 90) return 'danger'
  if (tokenPercent.value > 70) return 'warning'
  return ''
})

// Animate token count
watch(previewTokens, (newVal) => {
  const start = animatedTokens.value
  const diff = newVal - start
  const duration = 300
  const startTime = performance.now()
  
  function animate(currentTime: number) {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)
    animatedTokens.value = Math.round(start + diff * progress)
    if (progress < 1) requestAnimationFrame(animate)
  }
  requestAnimationFrame(animate)
}, { immediate: true })

const previewContent = computed(() => {
  if (!props.editingTemplate) return ''
  const tpl = props.editingTemplate
  const parts: string[] = []
  if (tpl.customPrefix) parts.push(tpl.customPrefix)
  
  const files = contextStore.summary?.files || []
  const hasContext = contextStore.hasContext && files.length > 0
  
  for (const sec of tpl.sectionOrder) {
    if (!tpl.sections[sec]) continue
    switch (sec) {
      case 'role': 
        if (tpl.roleContent) parts.push(`## Role\n${tpl.roleContent}`)
        break
      case 'rules': 
        if (tpl.rulesContent) parts.push(`## Rules\n${tpl.rulesContent}`)
        break
      case 'tree': 
        if (hasContext) {
          const tree = files.slice(0, 15).map(f => `  ${f}`).join('\n')
          parts.push(`## Project Structure\n${tree}${files.length > 15 ? `\n  ... +${files.length - 15} files` : ''}`)
        } else {
          parts.push(`## Project Structure\n[Build context to see file tree]`)
        }
        break
      case 'stats': 
        parts.push(`## Stats\n- Files: ${contextStore.fileCount || 0}\n- Lines: ${contextStore.lineCount || 0}\n- Tokens: ~${contextStore.tokenCount || 0}`)
        break
      case 'task': 
        parts.push(`## Task\n[Your task description]`)
        break
      case 'files': 
        if (hasContext) {
          const fileList = files.slice(0, 5).map(f => `- ${f}`).join('\n')
          parts.push(`## Files\n${fileList}${files.length > 5 ? `\n... +${files.length - 5} more files` : ''}\n\n[File contents]`)
        } else {
          parts.push(`## Files\n[Build context to see files]`)
        }
        break
    }
  }
  if (tpl.customSuffix) parts.push(tpl.customSuffix)
  return parts.join('\n\n')
})

const highlightedPreview = computed(() => {
  if (!previewContent.value) return `<span class="tpl-preview-empty">${t('templates.previewEmpty')}</span>`
  return previewContent.value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/^(## .+)$/gm, '<span class="tpl-h2">$1</span>')
    .replace(/^(- .+)$/gm, '<span class="tpl-list">$1</span>')
    .replace(/\[([^\]]+)\]/g, '<span class="tpl-placeholder">[$1]</span>')
})
</script>

<style scoped>
.tpl-preview {
  display: flex;
  flex-direction: column;
  background: #08090C;
  border-left: 1px solid rgba(255, 255, 255, 0.06);
}

.tpl-preview-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}

.tpl-preview-header svg {
  color: var(--text-subtle);
}

.tpl-preview-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  font-size: 12px;
  font-family: var(--font-mono);
  line-height: 1.7;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  min-height: 0;
}

.tpl-preview-content::-webkit-scrollbar { width: 4px; }
.tpl-preview-content::-webkit-scrollbar-track { background: transparent; }
.tpl-preview-content::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 2px; }

/* Syntax highlighting in preview */
:deep(.tpl-h2) { color: #60a5fa; font-weight: 600; }
:deep(.tpl-list) { color: #a78bfa; }
:deep(.tpl-placeholder) { color: var(--text-subtle); font-style: italic; }
:deep(.tpl-preview-empty) { color: var(--text-subtle); font-style: italic; }

.tpl-preview-footer {
  padding: 0.75rem 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.tpl-token-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.tpl-token-fill {
  height: 100%;
  background: var(--accent-indigo);
  border-radius: 2px;
  transition: width 0.3s ease-out, background 0.3s;
}

.tpl-token-fill.warning { background: var(--color-warning); }
.tpl-token-fill.danger { background: var(--color-danger); animation: shake 0.5s; }

.tpl-token-info {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
}

.tpl-token-count {
  font-size: 14px;
  font-weight: 600;
  font-family: var(--font-mono);
  color: var(--text-primary);
}

.tpl-token-label {
  font-size: 10px;
  color: var(--text-subtle);
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-2px); }
  75% { transform: translateX(2px); }
}
</style>
