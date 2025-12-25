<template>
  <div class="markdown-preview">
    <div v-for="(block, index) in blocks" :key="index" class="preview-block">
      <h1 v-if="block.type === 'heading1'" v-html="block.content"></h1>
      <h2 v-if="block.type === 'heading2'" v-html="block.content"></h2>
      <h3 v-if="block.type === 'heading3'" v-html="block.content"></h3>
      <p v-if="block.type === 'paragraph'" v-html="block.content"></p>
      <pre v-if="block.type === 'code'" v-html="block.content"></pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  content: string
}>()

const blocks = computed(() => {
  if (!props.content) return []

  const lines = props.content.split('\n')
  const blocks: Array<{type: string, content: string}> = []
  let currentParagraph = ''

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i].trim()

    // Empty line - end current paragraph
    if (!line) {
      if (currentParagraph) {
        blocks.push({ type: 'paragraph', content: processInlineMarkdown(currentParagraph) })
        currentParagraph = ''
      }
      continue
    }

    // Heading
    if (line.startsWith('# ')) {
      if (currentParagraph) {
        blocks.push({ type: 'paragraph', content: processInlineMarkdown(currentParagraph) })
        currentParagraph = ''
      }
      blocks.push({ type: 'heading1', content: escapeHtml(line.substring(2)) })
    } else if (line.startsWith('## ')) {
      if (currentParagraph) {
        blocks.push({ type: 'paragraph', content: processInlineMarkdown(currentParagraph) })
        currentParagraph = ''
      }
      blocks.push({ type: 'heading2', content: escapeHtml(line.substring(3)) })
    } else if (line.startsWith('### ')) {
      if (currentParagraph) {
        blocks.push({ type: 'paragraph', content: processInlineMarkdown(currentParagraph) })
        currentParagraph = ''
      }
      blocks.push({ type: 'heading3', content: escapeHtml(line.substring(4)) })
    } else {
      // Add to current paragraph
      if (currentParagraph) {
        currentParagraph += ' ' + line
      } else {
        currentParagraph = line
      }
    }
  }

  // Add remaining paragraph
  if (currentParagraph) {
    blocks.push({ type: 'paragraph', content: processInlineMarkdown(currentParagraph) })
  }

  return blocks
})

function escapeHtml(unsafe: string): string {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;")
}

function processInlineMarkdown(text: string): string {
  // Process **bold** text
  text = text.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')

  // Process `code` text
  text = text.replace(/`(.*?)`/g, '<code>$1</code>')

  // Process line breaks within paragraphs
  text = text.replace(/\n/g, '<br>')

  return text
}
</script>

<style scoped>
.markdown-preview {
  font-family: var(--font-family);
  line-height: 1.6;
  color: var(--text-primary);
}

.preview-block {
  margin-bottom: 0.75rem;
}

.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3 {
  margin: 1rem 0 0.5rem 0;
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.markdown-preview h1 {
  font-size: 1.5rem;
  border-bottom: 1px solid var(--border-default);
  padding-bottom: 0.3rem;
}

.markdown-preview h2 {
  font-size: 1.3rem;
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: 0.2rem;
}

.markdown-preview h3 {
  font-size: 1.1rem;
}

.markdown-preview p {
  margin: 0.5rem 0;
}

.markdown-preview code {
  background: var(--bg-2);
  padding: 0.125rem 0.25rem;
  border-radius: var(--radius-sm);
  font-family: var(--font-family-mono);
  font-size: 0.9em;
}

.markdown-preview pre {
  background: var(--bg-2);
  padding: 1rem;
  border-radius: var(--radius-md);
  overflow-x: auto;
  margin: 1rem 0;
  white-space: pre-wrap;
}

.markdown-preview strong {
  font-weight: var(--font-weight-semibold);
}
</style>