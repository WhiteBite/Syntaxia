<template>
  <div
    class="symbol-item"
    :class="{ 'symbol-item--nested': isNested }"
    @click="handleClick"
  >
    <!-- Symbol Icon -->
    <div class="symbol-icon" :class="`symbol-icon--${symbol.kind}`">
      <svg v-if="symbol.kind === 'class'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      <svg v-else-if="symbol.kind === 'function'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
      </svg>
      <svg v-else-if="symbol.kind === 'method'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <svg v-else-if="symbol.kind === 'interface'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z" />
      </svg>
      <svg v-else-if="symbol.kind === 'variable'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M4 7v10c0 2 1 3 3 3h10c2 0 3-1 3-3V7c0-2-1-3-3-3H7c-2 0-3 1-3 3z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h8" />
      </svg>
      <svg v-else-if="symbol.kind === 'constant'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
      </svg>
      <svg v-else fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
      </svg>
    </div>

    <!-- Symbol Info -->
    <div class="symbol-info">
      <span class="symbol-name">{{ symbol.name }}</span>
      <span v-if="symbol.signature" class="symbol-signature">{{ truncatedSignature }}</span>
    </div>

    <!-- Line Number -->
    <span class="symbol-line">:{{ symbol.line }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Symbol } from '../types'

const props = defineProps<{
  symbol: Symbol
  isNested?: boolean
}>()

const emit = defineEmits<{
  (e: 'click', symbol: Symbol): void
}>()

const MAX_SIGNATURE_LENGTH = 40

const truncatedSignature = computed(() => {
  if (!props.symbol.signature) return ''
  if (props.symbol.signature.length <= MAX_SIGNATURE_LENGTH) {
    return props.symbol.signature
  }
  return props.symbol.signature.slice(0, MAX_SIGNATURE_LENGTH) + '...'
})

function handleClick() {
  emit('click', props.symbol)
}
</script>

<style scoped>
.symbol-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s ease-out;
}

.symbol-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.symbol-item--nested {
  padding-left: 24px;
}

.symbol-icon {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  color: #9ca3af;
}

.symbol-icon svg {
  width: 14px;
  height: 14px;
}

.symbol-icon--class { color: #f59e0b; }
.symbol-icon--function { color: #8b5cf6; }
.symbol-icon--method { color: #3b82f6; }
.symbol-icon--interface { color: #10b981; }
.symbol-icon--variable { color: #6366f1; }
.symbol-icon--constant { color: #ec4899; }
.symbol-icon--type { color: #14b8a6; }
.symbol-icon--enum { color: #f97316; }
.symbol-icon--property { color: #64748b; }

.symbol-info {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.symbol-name {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.symbol-signature {
  font-size: 11px;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.symbol-line {
  flex-shrink: 0;
  font-size: 11px;
  color: #6b7280;
  font-family: monospace;
}
</style>
