<template>
  <div class="ai-settings-panel">
    <!-- Section Header -->
    <div class="ai-settings-header">
      <div class="ai-settings-icon">
        <Lightbulb class="w-4 h-4" />
      </div>
      <span class="ai-settings-title">{{ t('settings.aiProvider') }}</span>
    </div>

    <!-- Settings Content -->
    <div class="ai-settings-content">
      <!-- Provider Selection -->
      <ProviderSelector
        :model-value="settings.selectedProvider"
        :providers="providers"
        :description="currentProvider?.description"
        @update:model-value="selectProvider"
      />

      <!-- Divider -->
      <div class="ai-settings-divider"></div>

      <!-- API Key Input (not for CLI) -->
      <ApiKeyInput
        v-if="needsApiKey"
        :model-value="currentApiKey"
        :provider-name="currentProvider?.name || ''"
        :show-key="showApiKey"
        :hint="currentProviderHint"
        @update:model-value="updateApiKey"
        @toggle-visibility="toggleShowApiKey"
        @clear="clearApiKey"
      />

      <!-- Qwen CLI Info -->
      <QwenCliInfo v-if="settings.selectedProvider === 'qwen-cli'" />

      <!-- Model Selection -->
      <ModelSelector
        v-if="settings.selectedProvider && availableModels.length > 0"
        :model-value="selectedModel"
        :models="availableModels"
        @update:model-value="updateModel"
      />

      <!-- Host URL (for LocalAI, Qwen API) -->
      <HostUrlInput
        v-if="needsHostUrl"
        :model-value="currentHostUrl"
        :placeholder="hostPlaceholder"
        @update:model-value="updateHost"
      />
    </div>

    <!-- Sticky Footer -->
    <div class="ai-settings-footer">
      <button
        @click="saveSettings"
        :disabled="isSaving"
        class="ai-settings-save-btn"
      >
        <Loader2 v-if="isSaving" class="animate-spin w-4 h-4" />
        {{ isSaving ? t('settings.saving') : t('settings.save') }}
      </button>

      <!-- Status Message -->
      <StatusMessage :message="statusMessage" :type="statusType" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAISettings } from '@/composables/useAISettings'
import { useI18n } from '@/composables/useI18n'
import { Lightbulb, Loader2 } from 'lucide-vue-next'
import { onMounted } from 'vue'
import {
    ApiKeyInput,
    HostUrlInput,
    ModelSelector,
    ProviderSelector,
    QwenCliInfo,
    StatusMessage
} from './ai'

const { t } = useI18n()

const {
  settings,
  showApiKey,
  isSaving,
  statusMessage,
  statusType,
  providers,
  currentProvider,
  currentProviderHint,
  currentApiKey,
  availableModels,
  selectedModel,
  needsApiKey,
  needsHostUrl,
  currentHostUrl,
  hostPlaceholder,
  selectProvider,
  updateApiKey,
  clearApiKey,
  updateModel,
  updateHost,
  toggleShowApiKey,
  loadSettings,
  saveSettings
} = useAISettings()

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.ai-settings-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: linear-gradient(180deg, #131620 0%, #0f111a 100%);
}

.ai-settings-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.ai-settings-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(99, 102, 241, 0.15) 100%);
  border: 1px solid rgba(139, 92, 246, 0.25);
  border-radius: 8px;
  color: #a78bfa;
}

.ai-settings-title {
  font-size: 14px;
  font-weight: 600;
  color: white;
}

.ai-settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.1) transparent;
}

.ai-settings-content::-webkit-scrollbar {
  width: 6px;
}

.ai-settings-content::-webkit-scrollbar-track {
  background: transparent;
}

.ai-settings-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
}

.ai-settings-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.05);
  margin: 4px 0;
}

.ai-settings-footer {
  padding: 12px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(0, 0, 0, 0.2);
}

.ai-settings-save-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 10px 16px;
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.15s ease-out;
  box-shadow: 0 2px 8px rgba(139, 92, 246, 0.25);
}

.ai-settings-save-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #a78bfa 0%, #818cf8 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.35);
}

.ai-settings-save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
