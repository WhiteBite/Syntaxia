<template>
  <div class="chat-welcome">
    <!-- Ambient glow background -->
    <div class="ambient-glow"></div>
    
    <!-- Hero section -->
    <div class="welcome-hero">
      <div class="welcome-icon">
        <div class="welcome-icon__glow"></div>
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
            d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      </div>
      
      <h2 class="welcome-title">Syntaxia AI</h2>
      
      <p class="welcome-subtitle">
        {{ t('chat.welcome.subtitle') }}
      </p>
    </div>

    <!-- Status badge -->
    <div class="welcome-status" :class="{ 'welcome-status--connected': isConnected }">
      <span class="welcome-status__dot"></span>
      <span v-if="isConnected">{{ providerName }} · {{ currentModel }}</span>
      <span v-else>{{ t('chat.welcome.disconnected') }}</span>
    </div>

    <!-- Onboarding tips -->
    <div class="welcome-tips">
      <div class="welcome-tip">
        <kbd>@</kbd>
        <span>{{ t('chat.welcome.tipMention') }}</span>
      </div>
      <div class="welcome-tip">
        <kbd>/</kbd>
        <span>{{ t('chat.welcome.tipCommand') }}</span>
      </div>
    </div>

    <!-- Prompt starters -->
    <div class="welcome-starters">
      <p class="welcome-starters__label">{{ t('chat.welcome.tryAsking') }}</p>
      <div class="starter-grid">
        <button
          v-for="starter in promptStarters"
          :key="starter.id"
          class="starter-card"
          @click="$emit('quick-action', { prompt: starter.prompt })"
        >
          <span class="starter-card__icon">{{ starter.icon }}</span>
          <span class="starter-card__text">{{ starter.label }}</span>
          <svg class="starter-card__arrow" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'

defineProps<{
  isConnected: boolean
  providerName: string
  currentModel: string
}>()

defineEmits<{
  'quick-action': [action: { prompt: string }]
}>()

const { t } = useI18n()

const promptStarters = computed(() => [
  { 
    id: 'analyze', 
    icon: '🔍', 
    label: t('chat.starters.analyze'),
    prompt: t('chat.prompts.analyze') 
  },
  { 
    id: 'explain', 
    icon: '💡', 
    label: t('chat.starters.explain'),
    prompt: t('chat.prompts.explain') 
  },
  { 
    id: 'refactor', 
    icon: '✨', 
    label: t('chat.starters.refactor'),
    prompt: t('chat.prompts.refactor') 
  },
  { 
    id: 'test', 
    icon: '🧪', 
    label: t('chat.starters.test'),
    prompt: t('chat.prompts.test') 
  },
])
</script>

<style scoped>
.chat-welcome {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 12px;
  height: 100%;
}

/* Ambient glow - subtle background effect */
.ambient-glow {
  position: absolute;
  top: 20%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.12) 0%, transparent 70%);
  filter: blur(40px);
  pointer-events: none;
}

/* Hero */
.welcome-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 12px;
}

.welcome-icon {
  position: relative;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(99, 102, 241, 0.15) 100%);
  border: 1px solid rgba(139, 92, 246, 0.3);
  border-radius: 14px;
  color: #a78bfa;
  margin-bottom: 10px;
}

.welcome-icon svg {
  width: 24px;
  height: 24px;
}

.welcome-icon__glow {
  position: absolute;
  inset: -4px;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.25) 0%, transparent 70%);
  border-radius: 18px;
  filter: blur(8px);
  animation: pulse-glow 3s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.05); }
}

.welcome-title {
  font-size: 16px;
  font-weight: 700;
  background: linear-gradient(135deg, #e5e7eb 0%, #a78bfa 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 4px;
}

.welcome-subtitle {
  font-size: 11px;
  color: #6b7280;
  text-align: center;
  max-width: 220px;
  line-height: 1.4;
}

/* Status */
.welcome-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  font-size: 10px;
  color: #6b7280;
  margin-bottom: 12px;
}

.welcome-status--connected {
  border-color: rgba(16, 185, 129, 0.2);
}

.welcome-status__dot {
  width: 5px;
  height: 5px;
  background: #4b5563;
  border-radius: 50%;
}

.welcome-status--connected .welcome-status__dot {
  background: #10b981;
  box-shadow: 0 0 6px rgba(16, 185, 129, 0.5);
}

/* Tips */
.welcome-tips {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.welcome-tip {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: #6b7280;
}

.welcome-tip kbd {
  padding: 2px 6px;
  background: rgba(139, 92, 246, 0.15);
  border: 1px solid rgba(139, 92, 246, 0.25);
  border-radius: 4px;
  font-family: inherit;
  font-size: 10px;
  font-weight: 600;
  color: #a78bfa;
}

/* Starters - Grid Layout */
.welcome-starters {
  width: 100%;
}

.welcome-starters__label {
  font-size: 9px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #4b5563;
  margin-bottom: 8px;
  padding-left: 2px;
}

.starter-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px;
}

.starter-card {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 8px;
  padding: 10px;
  background: #1e2330;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  position: relative;
}

.starter-card:hover {
  border-color: rgba(139, 92, 246, 0.4);
  background: rgba(255, 255, 255, 0.05);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.15);
}

.starter-card__icon {
  font-size: 16px;
  flex-shrink: 0;
}

.starter-card__text {
  font-size: 11px;
  font-weight: 500;
  color: #e5e7eb;
  line-height: 1.3;
  flex: 1;
}

.starter-card__arrow {
  display: none;
}
</style>
