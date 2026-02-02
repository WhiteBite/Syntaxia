<template>
  <div
    v-if="isActive"
    :id="`tabpanel-${name}`"
    role="tabpanel"
    :aria-labelledby="`tab-${name}`"
    class="base-tab"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
/**
 * BaseTab - Individual tab panel component
 * 
 * Must be used as a child of BaseTabs component.
 * Content is lazy-loaded - only rendered when the tab is active.
 * 
 * @example
 * // Basic usage
 * <BaseTab name="general" label="General Settings">
 *   <p>General settings content</p>
 * </BaseTab>
 * 
 * @example
 * // With icon
 * <BaseTab name="advanced" label="Advanced" :icon="CodeIcon">
 *   <p>Advanced settings content</p>
 * </BaseTab>
 * 
 * @example
 * // With custom label slot
 * <BaseTab name="custom" label="Custom">
 *   <template #label>
 *     <span class="custom-label">Custom Label</span>
 *   </template>
 *   <p>Tab content</p>
 * </BaseTab>
 * 
 * @example
 * // Disabled tab
 * <BaseTab name="disabled" label="Disabled" :disabled="true">
 *   <p>This content won't be accessible</p>
 * </BaseTab>
 */
import { computed, inject, onBeforeUnmount, onMounted, watch, type Component, type ComputedRef } from 'vue'

interface Props {
  name: string
  label: string
  icon?: Component
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false
})

// Inject functions from parent BaseTabs
const registerTab = inject<(tab: { name: string; label: string; icon?: Component; disabled?: boolean }) => void>('registerTab')
const unregisterTab = inject<(name: string) => void>('unregisterTab')
const updateTab = inject<(name: string, updates: Partial<{ name: string; label: string; icon?: Component; disabled?: boolean }>) => void>('updateTab')
const activeTab = inject<ComputedRef<string>>('activeTab')

if (!registerTab || !unregisterTab || !updateTab || !activeTab) {
  throw new Error('BaseTab must be used within BaseTabs component')
}

// Check if this tab is active
const isActive = computed(() => activeTab.value === props.name)

// Register tab on mount
onMounted(() => {
  registerTab({
    name: props.name,
    label: props.label,
    icon: props.icon,
    disabled: props.disabled
  })
})

// Unregister tab on unmount
onBeforeUnmount(() => {
  unregisterTab(props.name)
})

// Update tab info when props change
watch(() => [props.label, props.icon, props.disabled], () => {
  updateTab(props.name, {
    label: props.label,
    icon: props.icon,
    disabled: props.disabled
  })
}, { deep: true })
</script>

<style scoped>
.base-tab {
  animation: fadeIn var(--transition-fast);
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
