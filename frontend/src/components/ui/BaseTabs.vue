<template>
  <div class="base-tabs" role="tablist">
    <!-- Tab Headers -->
    <div class="base-tabs__header">
      <button
        v-for="(tab, index) in tabs"
        :key="tab.name"
        :ref="el => tabRefs[index] = el as HTMLButtonElement"
        role="tab"
        :aria-selected="modelValue === tab.name"
        :aria-controls="`tabpanel-${tab.name}`"
        :aria-disabled="tab.disabled"
        :tabindex="modelValue === tab.name ? 0 : -1"
        :class="[
          'base-tabs__tab',
          { 'base-tabs__tab--active': modelValue === tab.name },
          { 'base-tabs__tab--disabled': tab.disabled }
        ]"
        :disabled="tab.disabled"
        @click="selectTab(tab.name)"
        @keydown="handleKeyDown($event, index)"
      >
        <BaseIcon
          v-if="tab.icon"
          :icon="tab.icon"
          size="sm"
          class="base-tabs__tab-icon"
        />
        <span class="base-tabs__tab-label">
          <slot :name="`label-${tab.name}`">
            {{ tab.label }}
          </slot>
        </span>
      </button>
      
      <!-- Active indicator -->
      <div
        class="base-tabs__indicator"
        :style="indicatorStyle"
      />
    </div>

    <!-- Tab Panels -->
    <div class="base-tabs__content">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * BaseTabs - Tab navigation component with keyboard support
 * 
 * @example
 * // Basic usage
 * <BaseTabs v-model="activeTab">
 *   <BaseTab name="general" label="General">
 *     General content
 *   </BaseTab>
 *   <BaseTab name="advanced" label="Advanced">
 *     Advanced content
 *   </BaseTab>
 * </BaseTabs>
 * 
 * @example
 * // With icons
 * <BaseTabs v-model="activeTab">
 *   <BaseTab name="general" label="General" :icon="SettingsIcon">
 *     General content
 *   </BaseTab>
 *   <BaseTab name="advanced" label="Advanced" :icon="CodeIcon">
 *     Advanced content
 *   </BaseTab>
 * </BaseTabs>
 * 
 * @example
 * // With custom label slot
 * <BaseTabs v-model="activeTab">
 *   <BaseTab name="general" label="General">
 *     <template #label>
 *       <span>Custom Label</span>
 *     </template>
 *     Content
 *   </BaseTab>
 * </BaseTabs>
 */
import { computed, nextTick, onMounted, provide, ref, watch, type Component } from 'vue'
import BaseIcon from './BaseIcon.vue'

interface TabInfo {
  name: string
  label: string
  icon?: Component
  disabled?: boolean
}

interface Props {
  modelValue: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

// Tab registration
const tabs = ref<TabInfo[]>([])
const tabRefs = ref<(HTMLButtonElement | null)[]>([])

function registerTab(tab: TabInfo) {
  if (!tabs.value.find(t => t.name === tab.name)) {
    tabs.value.push(tab)
  }
}

function unregisterTab(name: string) {
  const index = tabs.value.findIndex(t => t.name === name)
  if (index !== -1) {
    tabs.value.splice(index, 1)
  }
}

function updateTab(name: string, updates: Partial<TabInfo>) {
  const tab = tabs.value.find(t => t.name === name)
  if (tab) {
    Object.assign(tab, updates)
  }
}

// Provide registration functions to child tabs
provide('registerTab', registerTab)
provide('unregisterTab', unregisterTab)
provide('updateTab', updateTab)
provide('activeTab', computed(() => props.modelValue))

// Tab selection
function selectTab(name: string) {
  const tab = tabs.value.find(t => t.name === name)
  if (tab && !tab.disabled) {
    emit('update:modelValue', name)
  }
}

// Keyboard navigation
function handleKeyDown(event: KeyboardEvent, currentIndex: number) {
  let targetIndex = currentIndex

  switch (event.key) {
    case 'ArrowLeft':
      event.preventDefault()
      targetIndex = findPreviousEnabledTab(currentIndex)
      break
    case 'ArrowRight':
      event.preventDefault()
      targetIndex = findNextEnabledTab(currentIndex)
      break
    case 'Home':
      event.preventDefault()
      targetIndex = findFirstEnabledTab()
      break
    case 'End':
      event.preventDefault()
      targetIndex = findLastEnabledTab()
      break
    default:
      return
  }

  if (targetIndex !== -1 && targetIndex !== currentIndex) {
    const targetTab = tabs.value[targetIndex]
    if (targetTab) {
      selectTab(targetTab.name)
      nextTick(() => {
        tabRefs.value[targetIndex]?.focus()
      })
    }
  }
}

function findNextEnabledTab(currentIndex: number): number {
  for (let i = currentIndex + 1; i < tabs.value.length; i++) {
    if (!tabs.value[i].disabled) return i
  }
  // Wrap around
  for (let i = 0; i < currentIndex; i++) {
    if (!tabs.value[i].disabled) return i
  }
  return currentIndex
}

function findPreviousEnabledTab(currentIndex: number): number {
  for (let i = currentIndex - 1; i >= 0; i--) {
    if (!tabs.value[i].disabled) return i
  }
  // Wrap around
  for (let i = tabs.value.length - 1; i > currentIndex; i--) {
    if (!tabs.value[i].disabled) return i
  }
  return currentIndex
}

function findFirstEnabledTab(): number {
  return tabs.value.findIndex(t => !t.disabled)
}

function findLastEnabledTab(): number {
  for (let i = tabs.value.length - 1; i >= 0; i--) {
    if (!tabs.value[i].disabled) return i
  }
  return -1
}

// Active indicator positioning
const indicatorStyle = ref({})

function updateIndicator() {
  const activeIndex = tabs.value.findIndex(t => t.name === props.modelValue)
  if (activeIndex !== -1 && tabRefs.value[activeIndex]) {
    const tabElement = tabRefs.value[activeIndex]
    if (tabElement) {
      indicatorStyle.value = {
        left: `${tabElement.offsetLeft}px`,
        width: `${tabElement.offsetWidth}px`
      }
    }
  }
}

watch(() => props.modelValue, () => {
  nextTick(updateIndicator)
})

watch(() => tabs.value.length, () => {
  nextTick(updateIndicator)
})

onMounted(() => {
  nextTick(updateIndicator)
})
</script>

<style scoped>
.base-tabs {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.base-tabs__header {
  position: relative;
  display: flex;
  gap: var(--space-1);
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: 2px;
}

.base-tabs__tab {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  cursor: pointer;
  transition: all var(--transition-fast);
  position: relative;
  white-space: nowrap;
  outline: none;
}

.base-tabs__tab:hover:not(.base-tabs__tab--disabled):not(.base-tabs__tab--active) {
  color: var(--text-secondary);
  background: var(--bg-2);
}

.base-tabs__tab--active {
  color: var(--text-primary);
  font-weight: var(--font-weight-semibold);
}

.base-tabs__tab--disabled {
  color: var(--text-disabled);
  cursor: not-allowed;
  opacity: 0.5;
}

.base-tabs__tab:focus-visible {
  box-shadow: 0 0 0 2px var(--bg-0), 0 0 0 4px var(--accent-indigo);
  z-index: 1;
}

.base-tabs__tab-icon {
  flex-shrink: 0;
}

.base-tabs__tab-label {
  flex: 1;
  min-width: 0;
}

.base-tabs__indicator {
  position: absolute;
  bottom: 0;
  height: 2px;
  background: var(--gradient-primary);
  border-radius: 1px;
  transition: all var(--transition-normal);
  pointer-events: none;
}

.base-tabs__content {
  flex: 1;
  min-height: 0;
}
</style>
