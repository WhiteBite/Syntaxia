/**
 * Confirm Dialog Composable
 * Provides a promise-based confirmation dialog
 */

import { ref } from 'vue'

export interface ConfirmOptions {
    title: string
    message: string
    confirmText?: string
    cancelText?: string
    variant?: 'info' | 'warning' | 'danger'
}

const isOpen = ref(false)
const options = ref<ConfirmOptions | null>(null)
let resolvePromise: ((value: boolean) => void) | null = null

export function useConfirm() {
    function confirm(opts: ConfirmOptions): Promise<boolean> {
        options.value = opts
        isOpen.value = true

        return new Promise((resolve) => {
            resolvePromise = resolve
        })
    }

    function handleConfirm() {
        isOpen.value = false
        resolvePromise?.(true)
        resolvePromise = null
    }

    function handleCancel() {
        isOpen.value = false
        resolvePromise?.(false)
        resolvePromise = null
    }

    return {
        isOpen,
        options,
        confirm,
        handleConfirm,
        handleCancel,
    }
}
