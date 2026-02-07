<template>
  <BasePopover v-model="isOpen" trigger="click" placement="bottom-end" :arrow="false">
    <template #trigger>
      <BaseButton
        variant="ghost"
        size="sm"
        icon-only
        :title="t('files.settings')"
      >
        <template #icon>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </template>
      </BaseButton>
    </template>

    <div class="settings-popover-content">
      <div class="settings-popover__header">
        <span>{{ t('settings.title') }}</span>
      </div>

      <div class="settings-popover__body">
        <!-- Quick Toggles -->
        <div class="settings-popover__toggles">
          <div class="settings-toggle" @click.stop="toggleSetting('useGitignore')">
            <span class="settings-toggle__label">{{ t('settings.useGitignore') }}</span>
            <div class="pointer-events-none">
              <ToggleSwitch 
                :model-value="settings.fileExplorer.useGitignore"
              />
            </div>
          </div>

          <div class="settings-toggle" @click.stop="toggleSetting('useCustomIgnore')">
            <span class="settings-toggle__label">{{ t('settings.useCustomIgnore') }}</span>
            <div class="pointer-events-none">
              <ToggleSwitch 
                :model-value="settings.fileExplorer.useCustomIgnore"
              />
            </div>
          </div>

          <div class="settings-toggle" @click.stop="toggleSetting('compactNestedFolders')">
            <span class="settings-toggle__label">{{ t('settings.compactFolders') }}</span>
            <div class="pointer-events-none">
              <ToggleSwitch 
                :model-value="settings.fileExplorer.compactNestedFolders"
              />
            </div>
          </div>

          <div class="settings-toggle" @click.stop="toggleSetting('foldersFirst')">
            <span class="settings-toggle__label">{{ t('settings.foldersFirst') }}</span>
            <div class="pointer-events-none">
              <ToggleSwitch 
                :model-value="settings.fileExplorer.foldersFirst"
              />
            </div>
          </div>

          <div class="settings-toggle" @click.stop="toggleSetting('allowSelectBinary')">
            <span class="settings-toggle__label">{{ t('settings.allowSelectBinary') }}</span>
            <div class="pointer-events-none">
              <ToggleSwitch 
                :model-value="settings.fileExplorer.allowSelectBinary"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Manage Rules Button -->
      <div class="settings-popover__footer">
        <BaseButton
          variant="ghost"
          size="sm"
          class="w-full"
          @click="openIgnoreRules"
        >
          <template #icon>
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
            </svg>
          </template>
          {{ t('settings.manageRules') }}
        </BaseButton>
      </div>
    </div>
  </BasePopover>
</template>

<script setup lang="ts">
import BasePopover from '@/components/ui/BasePopover.vue'
import ToggleSwitch from '@/components/ui/ToggleSwitch.vue'
import { BaseButton } from '@/components/ui'
import { useI18n } from '@/composables/useI18n'
import { useSettingsStore } from '@/stores/settings.store'
import { computed, ref } from 'vue'

const { t } = useI18n()
const settingsStore = useSettingsStore()

const emit = defineEmits<{
  (e: 'open-ignore-rules'): void
  (e: 'settings-changed'): void
}>()

const isOpen = ref(false)
const settings = computed(() => settingsStore.settings)

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


.settings-popover__header {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 0.75rem;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.settings-popover__body {
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

.settings-popover-content {
  width: 16rem;
  background: #1c1f2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.75rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}
</style>
