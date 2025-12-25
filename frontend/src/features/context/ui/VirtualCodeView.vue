<template>
  <div 
    ref="containerRef" 
    class="virtual-code-view"
  >
    <div 
      class="virtual-code-view__spacer"
      :style="{ height: `${totalHeight}px` }"
    >
      <div 
        class="virtual-code-view__content"
        :style="{ transform: `translateY(${offsetY}px)` }"
      >
        <template v-for="item in virtualItems" :key="item.index">
          <!-- Chunk boundary marker -->
          <div 
            v-if="chunkBoundaries?.has(item.index)"
            class="chunk-boundary"
          >
            <span class="chunk-boundary__line"></span>
            <span class="chunk-boundary__label">✂️ Chunk {{ getChunkNumber(item.index) }}</span>
            <span class="chunk-boundary__line"></span>
          </div>
          <div
            class="code-line"
            :class="{ 'code-line-highlight': isLineHighlighted(item.index) }"
            :style="{ height: `${LINE_HEIGHT}px` }"
          >
            <span class="line-number">{{ item.index + 1 }}</span>
            <span v-html="highlightLine(lines[item.index] || '')"></span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useVirtualizer } from '@tanstack/vue-virtual'
import { computed, ref, watch } from 'vue'

const LINE_HEIGHT = 20

const props = defineProps<{
  lines: string[]
  highlightedLines?: Set<number>
  searchQuery?: string
  chunkBoundaries?: Set<number>
  outputFormat?: string // 'xml' | 'markdown' | 'plain'
}>()

const containerRef = ref<HTMLElement | null>(null)
const lineCount = computed(() => props.lines.length)

// Virtual scroller setup
const virtualizer = useVirtualizer({
  get count() { return lineCount.value },
  getScrollElement: () => containerRef.value,
  estimateSize: () => LINE_HEIGHT,
  overscan: 30, // Render extra items above/below viewport for smooth scroll
})

const virtualItems = computed(() => virtualizer.value.getVirtualItems())
const totalHeight = computed(() => virtualizer.value.getTotalSize())
const offsetY = computed(() => virtualItems.value[0]?.start ?? 0)

function isLineHighlighted(lineIndex: number): boolean {
  return props.highlightedLines?.has(lineIndex) ?? false
}

function getChunkNumber(lineIndex: number): number {
  if (!props.chunkBoundaries) return 1
  let chunkNum = 1
  for (const boundary of Array.from(props.chunkBoundaries).sort((a, b) => a - b)) {
    if (boundary <= lineIndex) chunkNum++
    else break
  }
  return chunkNum
}

function highlightLine(line: string): string {
  let result = escapeHtml(line)
  
  // Apply syntax highlighting based on format
  result = applySyntaxHighlight(result, props.outputFormat || 'plain')
  
  // Apply search highlighting on top
  if (props.searchQuery) {
    const query = escapeHtml(props.searchQuery)
    const regex = new RegExp(`(${escapeRegex(query)})`, 'gi')
    result = result.replace(regex, '<mark class="search-highlight">$1</mark>')
  }
  
  return result
}

function applySyntaxHighlight(line: string, format: string): string {
  switch (format) {
    case 'xml': return highlightXml(line)
    case 'markdown': return highlightMarkdown(line)
    default: return highlightPlain(line)
  }
}

function highlightXml(line: string): string {
  return line
    .replace(/(&lt;\/?)(file|content)(&gt;)/g, '<span class="syntax-tag">$1$2$3</span>')
    .replace(/(path)(=)(&quot;[^&]*&quot;)/g,
      '<span class="syntax-attr">$1</span><span class="syntax-punct">$2</span><span class="syntax-string">$3</span>')
}

function highlightMarkdown(line: string): string {
  if (line.startsWith('## File:') || line.startsWith('# ')) {
    return `<span class="syntax-heading">${line}</span>`
  }
  if (line.startsWith('```')) {
    return `<span class="syntax-fence">${line}</span>`
  }
  return line
}

function highlightPlain(line: string): string {
  if (line.startsWith('--- File:') && line.endsWith('---')) {
    return `<span class="syntax-separator">${line}</span>`
  }
  return line
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

function escapeRegex(text: string): string {
  return text.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

// Scroll to specific line
function scrollToLine(lineIndex: number) {
  virtualizer.value.scrollToIndex(lineIndex, { align: 'center' })
}

// Expose for parent component
defineExpose({ scrollToLine, containerRef })

// Re-measure when lines change
watch(() => props.lines.length, () => virtualizer.value.measure())
</script>

<style scoped>
.virtual-code-view {
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  font-family: ui-monospace, 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  font-size: 0.8125rem;
  line-height: 1.25rem;
  background: var(--bg-1);
}

.virtual-code-view__spacer {
  position: relative;
  width: 100%;
}

.virtual-code-view__content {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
}

.code-line {
  display: flex;
  align-items: center;
  padding: 0 0.75rem;
  white-space: pre;
  overflow: hidden;
  text-overflow: ellipsis;
}

.code-line-highlight {
  background: rgba(234, 179, 8, 0.15);
}

.line-number {
  flex-shrink: 0;
  width: 3.5rem;
  padding-right: 1rem;
  text-align: right;
  color: #4b5563;
  user-select: none;
}

:deep(.search-highlight) {
  background: rgba(234, 179, 8, 0.4);
  color: inherit;
  border-radius: 2px;
  padding: 0 2px;
}

.chunk-boundary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0.75rem;
  color: #9333ea;
  font-size: 0.6875rem;
}

.chunk-boundary__line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, transparent, #9333ea40, transparent);
}

.chunk-boundary__label {
  flex-shrink: 0;
  padding: 0.125rem 0.5rem;
  background: rgba(147, 51, 234, 0.1);
  border-radius: 0.25rem;
}

/* Syntax highlighting */
:deep(.syntax-tag) { color: #f472b6; }
:deep(.syntax-attr) { color: #a78bfa; }
:deep(.syntax-string) { color: #34d399; }
:deep(.syntax-punct) { color: #94a3b8; }
:deep(.syntax-heading) { color: #60a5fa; font-weight: 600; }
:deep(.syntax-fence) { color: #fbbf24; }
:deep(.syntax-separator) { color: #6366f1; font-weight: 500; }
</style>
