<template>
  <div class="ai-chat">
    <!-- Compact header with model info -->
    <div class="ai-chat__header">
      <ModelLimitsIndicator
        v-if="isConnected"
        :used-tokens="totalUsedTokens"
        :model-name="currentModel"
        :provider-name="providerName"
      />
    </div>

    <!-- Messages area -->
    <div ref="messagesContainer" class="ai-chat__messages">
      <ChatWelcome 
        v-if="messages.length === 0"
        :is-connected="isConnected"
        :provider-name="providerName"
        :current-model="currentModel"
        @quick-action="handleQuickAction"
      />

      <ChatMessageItem
        v-for="(msg, index) in messages"
        :key="index"
        :message="msg"
        :index="index"
        :expanded-tool-calls="expandedToolCalls"
        @copy="copyMessage(msg.content, t)"
        @toggle-tools="toggleToolCalls"
      />

      <ThinkingIndicator v-if="isThinking" @stop="stopGeneration" />
    </div>

    <!-- Smart Context Preview (floating) -->
    <Transition name="slide-up">
      <SmartContextPreviewPanel
        v-if="smartContextPreview && smartContextPreview.files.length > 0"
        :preview="smartContextPreview"
        @confirm="(files) => confirmSmartContext(files, t, scrollToBottom)"
        @cancel="cancelSmartContext"
      />
    </Transition>

    <!-- Analyzing Indicator -->
    <Transition name="fade">
      <div v-if="isAnalyzing" class="ai-chat__analyzing">
        <div class="analyzing-pulse"></div>
        <svg class="w-4 h-4 text-purple-400 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
        <span>{{ t('chat.analyzing') }}</span>
      </div>
    </Transition>

    <!-- Command Center (Premium Input) -->
    <CommandCenter
      v-model="inputMessage"
      :is-thinking="isThinking"
      :is-analyzing="isAnalyzing"
      :has-messages="messages.length > 0"
      @send="handleSend"
      @clear="clearChat"
      @attach="handleAttach"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useChatMessages } from '@/features/ai-chat/composables/useChatMessages'
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import ChatMessageItem from './chat/ChatMessageItem.vue'
import ChatWelcome from './chat/ChatWelcome.vue'
import CommandCenter from './chat/CommandCenter.vue'
import SmartContextPreviewPanel from './chat/SmartContextPreviewPanel.vue'
import ThinkingIndicator from './chat/ThinkingIndicator.vue'
import ModelLimitsIndicator from './ModelLimitsIndicator.vue'

const { t } = useI18n()

const {
  messages,
  inputMessage,
  isThinking,
  isAnalyzing,
  smartContextPreview,
  expandedToolCalls,
  currentModel,
  providerName,
  isConnected,
  totalUsedTokens,
  sendSmartMessage,
  confirmSmartContext,
  cancelSmartContext,
  clearChat,
  toggleToolCalls,
  stopGeneration,
  copyMessage,
  initialize,
  cleanup,
} = useChatMessages()

const messagesContainer = ref<HTMLElement | null>(null)

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

async function handleSend() {
  await sendSmartMessage(scrollToBottom, t)
}

function handleQuickAction(action: { prompt: string }) {
  inputMessage.value = action.prompt
  handleSend()
}

function handleAttach() {
  // Will open file picker in future iteration
}

onMounted(() => {
  initialize()
})

onUnmounted(() => {
  cleanup()
})
</script>

<style scoped>
/* Root: Flex column, no overflow on container */
.ai-chat {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  background: linear-gradient(180deg, #131620 0%, #0f111a 50%, #0b0c10 100%);
  border-left: 1px solid rgba(255, 255, 255, 0.05);
}

/* Header: Fixed height */
.ai-chat__header {
  flex: none;
  padding: 12px 16px 8px;
}

/* Messages: Scrollable area with custom scrollbar */
.ai-chat__messages {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Custom scrollbar */
.ai-chat__messages::-webkit-scrollbar {
  width: 4px;
}

.ai-chat__messages::-webkit-scrollbar-track {
  background: transparent;
}

.ai-chat__messages::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
}

.ai-chat__messages::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* Analyzing indicator */
.ai-chat__analyzing {
  flex: none;
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 16px 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.1) 0%, rgba(99, 102, 241, 0.05) 100%);
  border: 1px solid rgba(139, 92, 246, 0.2);
  border-radius: 12px;
  font-size: 13px;
  color: #c4b5fd;
  position: relative;
  overflow: hidden;
}

.analyzing-pulse {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(139, 92, 246, 0.1) 50%, transparent 100%);
  animation: pulse-slide 1.5s ease-in-out infinite;
}

@keyframes pulse-slide {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

/* Transitions */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.25s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(16px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
