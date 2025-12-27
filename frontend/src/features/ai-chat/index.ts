// AI Chat module
export { useChatStore } from './model/chat.store'
export type { ChatHistory, ChatMode, Message, TokenBudget } from './model/chat.store'
export { default as ChatHistoryPanel } from './ui/ChatHistoryPanel.vue'
export { default as ChatModeSelector } from './ui/ChatModeSelector.vue'
export { default as ChatPanel } from './ui/ChatPanel.vue'
export { default as ContextFileItem } from './ui/ContextFileItem.vue'
export { default as ContextPreviewPanel } from './ui/ContextPreviewPanel.vue'
export { default as ExecutePreviewModal } from './ui/ExecutePreviewModal.vue'
export { default as MentionDropdown } from './ui/MentionDropdown.vue'
export { default as MessageItem } from './ui/MessageItem.vue'
export { default as TokenBudgetBar } from './ui/TokenBudgetBar.vue'

// API
export { chatApi } from './api/chat.api'
export type {
    StreamChatCallbacks,
    StreamChatRequest, StreamChunk, StreamSubscription,
    StreamToolCall,
    TokenCountInfo
} from './api/chat.api'

// Composables
export { useChat } from './composables/useChat'
export type {
    TokenBudget as ChatTokenBudget, ChatToolCall,
    DiffHunk,
    DiffLine,
    DiffPreview, ChatMessage as StreamingChatMessage
} from './composables/useChat'
export { useChatHistory } from './composables/useChatHistory'
export { useChatMessages } from './composables/useChatMessages'
export type {
    Message as ChatMessage,
    SmartContextPreview,
    ToolCallLog as ToolCallInfo
} from './composables/useChatMessages'
export { useCodeExplanation } from './composables/useCodeExplanation'
export type { ExplainCodeOptions } from './composables/useCodeExplanation'
export { useContextPreview } from './composables/useContextPreview'
export type { ContextFilePreview } from './composables/useContextPreview'
export { useMentionAutocomplete } from './composables/useMentionAutocomplete'
export type { UseMentionAutocompleteOptions } from './composables/useMentionAutocomplete'
export { useMentions } from './composables/useMentions'
export type { MentionResult, UseMentionsOptions } from './composables/useMentions'

// UI Components
export { default as ExplainCodeButton } from './ui/ExplainCodeButton.vue'

// Types
export { GIT_MENTION_OPTIONS, MENTION_PATTERNS } from './types/mentions'
export type {
    AutocompleteState,
    Mention,
    MentionCategory,
    MentionPattern,
    MentionSuggestion,
    MentionType,
    ParsedInput
} from './types/mentions'

