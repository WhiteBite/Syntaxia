<template>
  <div
    @click="$emit('select')"
    class="template-item"
    :class="{ active, 'is-current': isCurrent }"
  >
    <span class="item-icon">{{ template.icon }}</span>
    <span class="item-name">{{ template.name }}</span>
    <span v-if="isCurrent" class="current-badge">●</span>
    <div class="item-actions">
      <BaseButton
        variant="ghost"
        size="sm"
        icon-only
        @click.stop="$emit('toggle-favorite')"
        class="item-action fav"
        :class="{ active: template.isFavorite }"
        :title="t('templates.toggleFavorite')"
      >
        <template #icon>
          <Star class="w-3 h-3" :fill="template.isFavorite ? 'currentColor' : 'none'" />
        </template>
      </BaseButton>
      <BaseButton
        v-if="!template.isBuiltIn"
        variant="ghost"
        size="sm"
        icon-only
        @click.stop="$emit('delete')"
        class="item-action delete"
        :title="t('templates.delete')"
      >
        <template #icon>
          <Trash2 class="w-3 h-3" />
        </template>
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { Star, Trash2 } from 'lucide-vue-next'
import type { PromptTemplate } from '../model/template.types'
import BaseButton from '@/components/ui/BaseButton.vue'

const { t } = useI18n()

defineProps<{
  template: PromptTemplate
  active: boolean
  isCurrent?: boolean
}>()

defineEmits<{
  (e: 'select'): void
  (e: 'delete'): void
  (e: 'toggle-favorite'): void
}>()
</script>


<style scoped>
.template-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.625rem;
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.15s;
  margin-bottom: 0.125rem;
}

.template-item:hover {
  background: var(--bg-2);
}

.template-item.active {
  background: var(--accent-indigo-bg);
  border-color: var(--accent-indigo-border);
}

.item-icon {
  font-size: 1rem;
  line-height: 1;
  flex-shrink: 0;
}

.item-name {
  flex: 1;
  font-size: var(--font-size-xs);
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.current-badge {
  font-size: 0.5rem;
  color: var(--color-success);
  flex-shrink: 0;
}

.template-item.is-current .item-name {
  color: var(--color-success);
}

.item-actions {
  display: none;
  align-items: center;
  gap: 0.125rem;
}

.template-item:hover .item-actions {
  display: flex;
}

.item-action {
  width: 1.25rem;
  height: 1.25rem;
  padding: 0;
}

.item-action.fav:hover,
.item-action.fav.active {
  color: var(--color-warning);
}

.item-action.delete:hover {
  background: var(--color-danger-soft);
  color: var(--color-danger);
}
</style>
