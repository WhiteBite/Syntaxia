/**
 * useCodeExplanation composable
 * Provides functionality to explain selected code via AI chat
 */

import { useI18n } from '@/composables/useI18n'
import { useChatStore } from '../model/chat.store'

export interface ExplainCodeOptions {
    code: string
    language?: string
}

export function useCodeExplanation() {
    const { t } = useI18n()
    const chatStore = useChatStore()

    /**
     * Send code to AI chat for explanation
     */
    async function explainCode(code: string, language = 'text'): Promise<void> {
        if (!code.trim()) return

        const prompt = t('chat.explain.prompt', {
            language,
            code: `\`\`\`${language}\n${code}\n\`\`\``
        })

        await chatStore.sendMessage(prompt)
    }

    /**
     * Build explanation prompt for code
     */
    function buildExplanationPrompt(code: string, language = 'text'): string {
        return t('chat.explain.prompt', {
            language,
            code: `\`\`\`${language}\n${code}\n\`\`\``
        })
    }

    return {
        explainCode,
        buildExplanationPrompt
    }
}
