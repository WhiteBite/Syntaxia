<template>
  <div class="collapsible-section">
    <button 
      class="collapsible-header"
      :class="{ 'collapsible-header-open': isOpen }"
      @click="toggle"
    >
      <div class="collapsible-title">
        <div class="collapsible-icon" :class="iconClass">
          <component :is="icon" class="w-3.5 h-3.5" />
        </div>
        <span class="collapsible-label">{{ label }}</span>
      </div>
      <ChevronDown 
        class="collapsible-arrow" 
        :class="{ 'collapsible-arrow-open': isOpen }"
      />
    </button>
    <Transition name="collapsible">
      <div v-show="isOpen" class="collapsible-content">
        <slot />
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ChevronDown } from 'lucide-vue-next';
import { ref, type Component } from 'vue';

const props = withDefaults(defineProps<{
  label: string
  icon: Component
  iconColor?: 'indigo' | 'purple' | 'emerald' | 'orange' | 'pink'
  defaultOpen?: boolean
}>(), {
  iconColor: 'indigo',
  defaultOpen: true
})

const isOpen = ref(props.defaultOpen)

const iconClass = `collapsible-icon-${props.iconColor}`

function toggle() {
  isOpen.value = !isOpen.value
}
</script>
