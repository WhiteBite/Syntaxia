<template>
  <div v-if="hasAnyDependencies" class="tree-row-deps">
    <!-- Outgoing dependencies (files this file imports) -->
    <BaseTooltip
      v-if="outgoingCount > 0"
      :text="getTooltipText('outgoing')"
      position="top"
      :multiline="true"
    >
      <div class="dep-group" @click.stop="handleClick('outgoing')">
        <BaseIcon 
          :icon="ArrowRightCircle"
          size="sm"
          class="dep-icon dep-outgoing"
        />
        <span class="dep-count">{{ outgoingCount }}</span>
      </div>
    </BaseTooltip>

    <!-- Incoming dependencies (files that import this file) -->
    <BaseTooltip
      v-if="incomingCount > 0"
      :text="getTooltipText('incoming')"
      position="top"
      :multiline="true"
    >
      <div class="dep-group" @click.stop="handleClick('incoming')">
        <BaseIcon 
          :icon="ArrowLeftCircle"
          size="sm"
          class="dep-icon dep-incoming"
        />
        <span class="dep-count">{{ incomingCount }}</span>
      </div>
    </BaseTooltip>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useFileStore } from '../model/file.store'
import BaseIcon from '@/components/ui/BaseIcon.vue'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import { ArrowRightCircle, ArrowLeftCircle } from 'lucide-vue-next'
import { useI18n } from '@/composables/useI18n'

const props = defineProps<{
  filePath: string
}>()

const emit = defineEmits<{
  showDependencies: [direction: 'incoming' | 'outgoing']
}>()

const fileStore = useFileStore()
const { t } = useI18n()

const dependencies = computed(() => 
  fileStore.allFileDependencies.get(props.filePath) || { incoming: [], outgoing: [] }
)

const outgoingCount = computed(() => dependencies.value.outgoing.length)
const incomingCount = computed(() => dependencies.value.incoming.length)
const hasAnyDependencies = computed(() => outgoingCount.value > 0 || incomingCount.value > 0)

function handleClick(direction: 'incoming' | 'outgoing') {
  emit('showDependencies', direction)
}

function getTooltipText(direction: 'incoming' | 'outgoing'): string {
  const deps = direction === 'incoming' ? dependencies.value.incoming : dependencies.value.outgoing
  const count = deps.length
  
  if (count === 0) return ''
  
  const label = direction === 'incoming' 
    ? t('files.incomingDependencies', { count })
    : t('files.outgoingDependencies', { count })
  
  // Show first 3 file names
  const fileNames = deps.slice(0, 3).map(path => {
    const name = path.substring(path.lastIndexOf('/') + 1)
    return `• ${name}`
  })
  
  if (count > 3) {
    fileNames.push(`• +${count - 3} more...`)
  }
  
  return `${label}\n\n${fileNames.join('\n')}`
}
</script>

<style scoped>
.tree-row-deps {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-left: auto;
  opacity: 0;
  transition: opacity 200ms ease-out;
}

.tree-row:hover .tree-row-deps,
.tree-row:focus-within .tree-row-deps {
  opacity: 1;
}

.dep-group {
  display: flex;
  align-items: center;
  gap: 2px;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: var(--radius-sm);
  transition: background 150ms ease-out;
}

.dep-group:hover {
  background: var(--bg-hover);
}

.dep-icon {
  transition: transform 150ms ease-out;
  flex-shrink: 0;
}

.dep-icon.dep-outgoing {
  color: var(--color-success);
}

.dep-icon.dep-incoming {
  color: var(--color-info);
}

.dep-group:hover .dep-icon {
  transform: scale(1.15);
}

.dep-count {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-secondary);
  min-width: 12px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}
</style>
