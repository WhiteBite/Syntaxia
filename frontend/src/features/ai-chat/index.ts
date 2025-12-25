// AI Chat module
export { useChatStore } from './model/chat.store'
export type { ChatHistory, Message } from './model/chat.store'
export { default as ChatPanel } from './ui/ChatPanel.vue'
export { default as MessageItem } from './ui/MessageItem.vue'

// Composables
export { useChatMessages } from './composables/useChatMessages'
export type {
    Message as ChatMessage,
    SmartContextPreview,
    ToolCallLog as ToolCallInfo
} from './composables/useChatMessages'
export { useMentions } from './composables/useMentions'
export type { MentionResult } from './composables/useMentions'

