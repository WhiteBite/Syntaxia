<template>
  <div class="skeleton-loader">
    <!-- Skeleton Lines (simulating real code structure) -->
    <div class="skeleton-content">
      <div 
        v-for="(line, i) in codeLines" 
        :key="i" 
        class="skeleton-line"
        :style="{ 
          width: line.width,
          marginLeft: line.indent,
          animationDelay: `${(i % 6) * 0.08}s`
        }"
      />
      
      <!-- Centered Status Text overlay -->
      <div class="skeleton-overlay">
        <span class="skeleton-status-text">{{ statusText }}</span>
        <span class="skeleton-status-dots">{{ dots }}</span>
      </div>
    </div>
    
    <!-- Progress Bar at bottom (browser-style) -->
    <div class="skeleton-progress">
      <div 
        class="skeleton-progress-fill" 
        :style="{ width: `${progress}%` }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

defineProps<{
  progress: number
  statusText: string
}>()

const dots = ref('')

let dotsInterval: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  dotsInterval = setInterval(() => {
    dots.value = dots.value.length >= 3 ? '' : dots.value + '.'
  }, 400)
})

onUnmounted(() => {
  if (dotsInterval) clearInterval(dotsInterval)
})

// Realistic code structure with indentation
const codeLines = [
  { width: '25%', indent: '0' },      // import
  { width: '35%', indent: '0' },      // import
  { width: '20%', indent: '0' },      // empty/short
  { width: '40%', indent: '0' },      // function declaration
  { width: '60%', indent: '1.5rem' }, // indented code
  { width: '75%', indent: '1.5rem' }, // longer line
  { width: '45%', indent: '1.5rem' }, // medium line
  { width: '55%', indent: '3rem' },   // nested
  { width: '80%', indent: '3rem' },   // long nested
  { width: '30%', indent: '3rem' },   // short nested
  { width: '50%', indent: '1.5rem' }, // back to first indent
  { width: '15%', indent: '0' },      // closing brace
  { width: '10%', indent: '0' },      // empty
  { width: '45%', indent: '0' },      // new function
  { width: '70%', indent: '1.5rem' }, // code
  { width: '35%', indent: '1.5rem' }, // return
]
</script>

<style scoped>
.skeleton-loader {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-1);
  overflow: hidden;
  position: relative;
}

.skeleton-content {
  flex: 1;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
  overflow: hidden;
  position: relative;
}

.skeleton-line {
  height: 0.75rem;
  background: linear-gradient(
    90deg,
    rgba(139, 92, 246, 0.12) 0%,
    rgba(139, 92, 246, 0.22) 50%,
    rgba(139, 92, 246, 0.12) 100%
  );
  background-size: 200% 100%;
  border-radius: 0.25rem;
  animation: shimmer 1.8s ease-in-out infinite;
  flex-shrink: 0;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* Centered overlay with status - no blur to keep lines visible */
.skeleton-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.125rem;
  background: rgba(15, 17, 26, 0.5);
}

.skeleton-status-text {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #e5e7eb;
}

.skeleton-status-dots {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #8b5cf6;
  width: 1.25rem;
}

/* Progress Bar - browser style, full width */
.skeleton-progress {
  height: 3px;
  background: rgba(139, 92, 246, 0.1);
  flex-shrink: 0;
}

.skeleton-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #8b5cf6, #ec4899);
  transition: width 0.15s ease-out;
}
</style>
