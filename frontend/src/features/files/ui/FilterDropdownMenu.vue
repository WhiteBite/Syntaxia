<template>
  <Teleport to="body">
    <Transition name="dropdown">
      <div v-if="isOpen" class="filter-dropdown-menu" :class="menuClass" :style="style">
        <div class="filter-dropdown-header" :class="headerClass">
          <span>{{ title }}</span>
          <button v-if="hasActiveFilters" @click="$emit('clear')" class="filter-clear-btn">
            {{ t('quickFilters.clearAll') }}
          </button>
        </div>
        <div class="filter-dropdown-items">
          <slot />
        </div>
        <div v-if="hint" class="filter-dropdown-hint">
          <kbd>Ctrl</kbd> {{ t('quickFilters.multiSelect') }} · <kbd>Shift</kbd> {{ t('quickFilters.exclude') }}
        </div>
        <div v-if="footer" class="filter-dropdown-footer" :class="footerClass">
          <span class="filter-dropdown-footer-text">{{ footer }}</span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';

const { t } = useI18n()

defineProps<{
  isOpen: boolean
  title: string
  style: Record<string, string>
  hasActiveFilters?: boolean
  hint?: boolean
  footer?: string
  menuClass?: string
  headerClass?: string
  footerClass?: string
}>()

defineEmits<{
  clear: []
}>()
</script>

<style scoped>
.filter-dropdown-menu {
  @apply rounded-xl overflow-hidden;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  box-shadow: var(--shadow-xl);
  min-width: min(260px, 80vw);
  max-width: min(320px, 40vw);
}
.filter-dropdown-header {
  @apply flex items-center justify-between px-3 py-2 text-xs font-semibold;
  background: var(--bg-2);
  border-bottom: 1px solid var(--border-default);
  color: var(--text-primary);
}
.filter-clear-btn {
  @apply text-[10px] px-2 py-0.5 rounded;
  color: var(--accent-indigo);
  transition: all 150ms ease-out;
}
.filter-clear-btn:hover { background: var(--accent-indigo-bg); }
.filter-dropdown-items { @apply py-1 max-h-[280px] overflow-y-auto; }
.filter-dropdown-hint {
  @apply px-3 py-2 text-[10px] text-center;
  background: var(--bg-2);
  border-top: 1px solid var(--border-default);
  color: var(--text-muted);
}
.filter-dropdown-hint kbd {
  @apply px-1.5 py-0.5 rounded text-[9px] font-mono;
  background: var(--bg-3);
  border: 1px solid var(--border-default);
}
.filter-dropdown-footer {
  @apply px-3 py-2 text-center;
  background: var(--bg-2);
  border-top: 1px solid var(--border-default);
}
.filter-dropdown-footer-text { @apply text-[10px]; color: var(--accent-emerald); }

.dropdown-enter-active, .dropdown-leave-active { transition: all 150ms ease-out; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
