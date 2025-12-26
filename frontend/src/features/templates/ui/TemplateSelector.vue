<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="isOpen = !isOpen"
      class="template-selector-btn"
      :title="t('templates.selectTemplate')"
    >
      <span class="template-icon">{{ activeTemplate?.icon || '📝' }}</span>
      <span class="template-name">{{ activeTemplate?.name || t('templates.select') }}</span>
      <ChevronDown class="w-3.5 h-3.5 text-gray-400 transition-transform" :class="{ 'rotate-180': isOpen }" />
    </button>

    <Transition name="dropdown">
      <div v-if="isOpen" class="template-dropdown">
        <!-- Built-in templates -->
        <div class="template-group">
          <div class="template-group-label">{{ t('templates.builtIn') }}</div>
          <button
            v-for="template in builtInTemplates"
            :key="template.id"
            @click="selectTemplate(template.id)"
            class="template-option"
            :class="{ 'template-option-active': template.id === activeTemplateId }"
          >
            <span class="template-icon">{{ template.icon }}</span>
            <div class="template-option-content">
              <span class="template-option-name">{{ template.name }}</span>
              <span class="template-option-desc">{{ template.description }}</span>
            </div>
          </button>
        </div>

        <!-- Custom templates -->
        <div v-if="customTemplates.length > 0" class="template-group">
          <div class="template-group-label">{{ t('templates.custom') }}</div>
          <button
            v-for="template in customTemplates"
            :key="template.id"
            @click="selectTemplate(template.id)"
            class="template-option"
            :class="{ 'template-option-active': template.id === activeTemplateId }"
          >
            <span class="template-icon">{{ template.icon }}</span>
            <div class="template-option-content">
              <span class="template-option-name">{{ template.name }}</span>
              <span class="template-option-desc">{{ template.description }}</span>
            </div>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { onClickOutside } from '@vueuse/core'
import { ChevronDown } from 'lucide-vue-next'
import { ref } from 'vue'
import { useTemplateStore } from '../model/template.store'
import { storeToRefs } from 'pinia'

const { t } = useI18n()
const templateStore = useTemplateStore()
const { activeTemplate, activeTemplateId, builtInTemplates, customTemplates } = storeToRefs(templateStore)

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

onClickOutside(dropdownRef, () => {
  isOpen.value = false
})

function selectTemplate(id: string) {
  templateStore.setActiveTemplate(id)
  isOpen.value = false
}
</script>

<style scoped>
.template-selector-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.625rem;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: 0.375rem;
  color: var(--text-1);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: all 0.15s;
  width: 100%;
}

.template-selector-btn:hover {
  background: var(--bg-3);
  border-color: var(--border-hover);
}

.template-icon {
  font-size: 1rem;
  line-height: 1;
}

.template-name {
  flex: 1;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.template-dropdown {
  position: absolute;
  top: calc(100% + 0.25rem);
  left: 0;
  right: 0;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: 0.5rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 50;
  max-height: 280px;
  overflow-y: auto;
}

.template-group {
  padding: 0.25rem;
}

.template-group + .template-group {
  border-top: 1px solid var(--border-default);
}

.template-group-label {
  padding: 0.375rem 0.5rem;
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--text-3);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.template-option {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  width: 100%;
  padding: 0.5rem;
  background: transparent;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  text-align: left;
  transition: background 0.15s;
}

.template-option:hover {
  background: var(--bg-3);
}

.template-option-active {
  background: var(--bg-accent-subtle);
}

.template-option-content {
  flex: 1;
  min-width: 0;
}

.template-option-name {
  display: block;
  font-size: 0.8125rem;
  color: var(--text-1);
  font-weight: 500;
}

.template-option-desc {
  display: block;
  font-size: 0.6875rem;
  color: var(--text-3);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
