<template>
  <BaseTooltip :text="tooltipContent" :position="position" :multiline="true">
    <slot />
  </BaseTooltip>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useFileStore } from '../model/file.store'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import { useI18n } from '@/composables/useI18n'

const props = withDefaults(defineProps<{
  filePath: string
  position?: 'top' | 'bottom' | 'left' | 'right'
}>(), {
  position: 'top'
})

const fileStore = useFileStore()
const { t } = useI18n()

const dependencies = computed(() => 
  fileStore.allFileDependencies.get(props.filePath) || { incoming: [], outgoing: [] }
)

const tooltipContent = computed(() => {
  const { incoming, outgoing } = dependencies.value
  const parts: string[] = []
  
  if (incoming.length > 0) {
    parts.push(`${t('files.incomingDependencies', { count: incoming.length })}:`)
    incoming.slice(0, 5).forEach(path => {
      const name = path.substring(path.lastIndexOf('/') + 1)
      parts.push(`  • ${name}`)
    })
    if (incoming.length > 5) {
      parts.push(`  • +${incoming.length - 5} more...`)
    }
  }
  
  if (outgoing.length > 0) {
    if (parts.length > 0) parts.push('') // Empty line separator
    parts.push(`${t('files.outgoingDependencies', { count: outgoing.length })}:`)
    outgoing.slice(0, 5).forEach(path => {
      const name = path.substring(path.lastIndexOf('/') + 1)
      parts.push(`  • ${name}`)
    })
    if (outgoing.length > 5) {
      parts.push(`  • +${outgoing.length - 5} more...`)
    }
  }
  
  if (parts.length === 0) {
    return t('files.noDependencies')
  }
  
  return parts.join('\n')
})
</script>
