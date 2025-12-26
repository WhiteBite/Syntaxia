<template>
  <div class="settings-popover-wrapper">
    <!-- Trigger Button -->
    <button
      ref="triggerRef"
      @click="toggle"
      class="settings-trigger"
      :class="{ 'settings-trigger--active': isOpen }"
      :title="t('files.settings')"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
    </button>

    <!-- Popover via Teleport -->
    <Teleport to="body">
      <Transition name="popover">
        <div v-if="isOpen" ref="popoverRef" class="settings-popover" :style="popoverStyle">
        <div class="settings-popover__header">
          <span>{{ t('settings.title') }}</span>
        </div>

        <div class="settings-popover__content">
          <!-- Quick Toggles -->
          <div class="settings-popover__toggles">
            <div class="settings-toggle" @click="toggleSetting('useGitignore')">
              <span class="settings-toggle__label">{{ t('settings.useGitignore') }}</span>
              <ToggleSwitch 
                :model-value="settings.fileExplorer.useGitignore"
                @update:model-value="v => updateSetting('useGitignore', v)"
              />
            </div>

            <div class="settings-toggle" @click="toggleSetting('useCustomIgnore')">
              <span class="settings-toggle__label">{{ t('settings.useCustomIgnore') }}</span>
              <ToggleSwitch 
                :model-value="settings.fileExplorer.useCustomIgnore"
                @update:model-value="v => updateSetting('useCustomIgnore', v)"
              />
            </div>

            <div class="settings-toggle" @click="toggleSetting('compactNestedFolders')">
              <span class="settings-toggle__label">{{ t('settings.compactFolders') }}</span>
              <ToggleSwitch 
                :model-value="settings.fileExplorer.compactNestedFolders"
                @update:model-value="v => updateSetting('compactNestedFolders', v)"
              />
            </div>

            <div class="settings-toggle" @click="toggleSetting('foldersFirst')">
              <span class="settings-toggle__label">{{ t('settings.foldersFirst') }}</span>
              <ToggleSwitch 
                :model-value="settings.fileExplorer.foldersFirst"
                @update:model-value="v => updateSetting('foldersFirst', v)"
              />
            </div>

            <div class="settings-toggle" @click="toggleSetting('allowSelectBinary')">
              <span class="settings-toggle__label">{{ t('settings.allowSelectBinary') }}</span>
              <ToggleSwitch 
                :model-value="settings.fileExplorer.allowSelectBinary"
                @update:model-value="v => updateSetting('allowSelectBinary', v)"
              />
            </div>
          </div>
        </div>

        <!-- Manage Rules Button -->
        <div class="settings-popover__footer">
          <button @click="openIgnoreRules" class="settings-popover__manage-btn">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
            </svg>
            {{ t('settings.manageRules') }}
          </button>
        </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import ToggleSwitch from '@/components/ui/ToggleSwitch.vue'
import { useI18n } from '@/composables/useI18n'
import { useSettingsStore } from '@/stores/settings.store'
import { onClickOutside } from '@vueuse/core'
import { ref, computed, watch, nextTick } from 'vue'

const { t } = useI18n()
const settingsStore = useSettingsStore()

const emit = defineEmits<{
  (e: 'open-ignore-rules'): void
  (e: 'settings-changed'): void
}>()

const isOpen = ref(false)
const triggerRef = ref<HTMLElement | null>(null)
const popoverRef = ref<HTMLElement | null>(null)
const popoverPosition = ref({ top: 0, left: 0 })

const settings = computed(() => settingsStore.settings)

const popoverStyle = computed(() => ({
  top: `${popoverPosition.value.top}px`,
  left: `${popoverPosition.value.left}px`,
}))

onClickOutside(popoverRef, (event) => {
  // Don't close if clicking the trigger
  if (triggerRef.value?.contains(event.target as Node)) return
  isOpen.value = false
})

function updatePosition() {
  if (!triggerRef.value) return
  const rect = triggerRef.value.getBoundingClientRect()
  const popoverWidth = 256 // 16rem
  
  // Position below trigger, aligned to right edge
  let left = rect.right - popoverWidth
  const top = rect.bottom + 8
  
  // Keep within viewport
  if (left < 8) left = 8
  
  popoverPosition.value = { top, left }
}

watch(isOpen, async (open) => {
  if (open) {
    await nextTick()
    updatePosition()
  }
})

function toggle() {
  isOpen.value = !isOpen.value
}

function toggleSetting(key: keyof typeof settings.value.fileExplorer) {
  const current = settings.value.fileExplorer[key]
  if (typeof current === 'boolean') {
    updateSetting(key, !current)
  }
}

function updateSetting(key: keyof typeof settings.value.fileExplorer, value: boolean) {
  settingsStore.updateFileExplorerSettings({ [key]: value })
  emit('settings-changed')
}

function openIgnoreRules() {
  isOpen.value = false
  emit('open-ignore-rules')
}
</script>

<style scoped>
.settings-popover-wrapper {
  position: relative;
}

.settings-trigger {
  padding: 0.5rem;
  border-radius: 0.5rem;
  color: #64748b;
  transition: all 150ms;
}

.settings-trigger:hover,
.settings-trigger--active {
  background: rgba(99, 102, 241, 0.15);
  color: #a5b4fc;
}

.settings-popover {
  position: fixed;
  width: 16rem;
  background: rgba(15, 23, 42, 0.98);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.75rem;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.5);
  z-index: 99999;
  overflow: hidden;
}

.settings-popover__header {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 0.75rem;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.settings-popover__content {
  padding: 0.5rem;
}

.settings-popover__toggles {
  display: flex;
  flex-direction: column;
}

.settings-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.625rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: background 150ms;
}

.settings-toggle:hover {
  background: rgba(255, 255, 255, 0.03);
}

.settings-toggle__label {
  font-size: 0.8125rem;
  color: #e2e8f0;
}

.settings-popover__footer {
  padding: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.settings-popover__manage-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.625rem;
  background: rgba(168, 85, 247, 0.1);
  border: 1px solid rgba(168, 85, 247, 0.2);
  border-radius: 0.5rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: #c084fc;
  transition: all 150ms;
}

.settings-popover__manage-btn:hover {
  background: rgba(168, 85, 247, 0.15);
  border-color: rgba(168, 85, 247, 0.3);
}

/* Transitions */
.popover-enter-active, .popover-leave-active {
  transition: all 200ms ease;
}
.popover-enter-from, .popover-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.95);
}
</style>
