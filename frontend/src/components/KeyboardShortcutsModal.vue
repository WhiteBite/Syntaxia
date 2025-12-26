<template>
  <TransitionRoot :show="isOpen" as="template">
    <Dialog @close="$emit('close')" class="relative z-50">
      <TransitionChild
        as="template"
        enter="ease-out duration-200"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-150"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4">
          <TransitionChild
            as="template"
            enter="ease-out duration-200"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="ease-in duration-150"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel class="w-full max-w-4xl transform overflow-hidden rounded-xl bg-white dark:bg-gray-800 shadow-2xl transition-all">
              <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                <DialogTitle class="text-xl font-semibold text-gray-900 dark:text-gray-100">
                  Keyboard Shortcuts
                </DialogTitle>
              </div>

              <div class="px-6 py-4 max-h-[70vh] overflow-y-auto">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div v-for="category in categories" :key="category.name">
                    <h3 class="text-sm font-semibold text-gray-400 dark:text-gray-400 uppercase tracking-wider mb-3">
                      {{ category.name }}
                    </h3>
                    <div class="space-y-2">
                      <div
                        v-for="shortcut in category.shortcuts"
                        :key="shortcut.key"
                        class="flex items-center justify-between py-2 px-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
                      >
                        <span class="text-sm text-gray-700 dark:text-gray-300">
                          {{ shortcut.description }}
                        </span>
                        <div class="flex gap-1">
                          <kbd
                            v-for="(key, index) in shortcut.key.split('+')"
                            :key="index"
                            class="px-2 py-1 text-xs font-semibold text-gray-600 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded shadow-sm"
                          >
                            {{ key }}
                          </kbd>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div class="px-6 py-4 bg-gray-50 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 flex justify-end">
                <button
                  @click="$emit('close')"
                  class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
                >
                  Close
                </button>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { TransitionRoot, TransitionChild, Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'

defineProps<{
  isOpen: boolean
}>()

defineEmits<{
  close: []
}>()

const categories = [
  {
    name: 'General',
    shortcuts: [
      { key: 'Ctrl+K', description: 'Open command palette' },
      { key: 'Ctrl+/', description: 'Show keyboard shortcuts' },
      { key: 'Escape', description: 'Close modal' },
    ]
  },
  {
    name: 'File',
    shortcuts: [
      { key: 'Ctrl+N', description: 'New file' },
      { key: 'Ctrl+O', description: 'Open file' },
      { key: 'Ctrl+S', description: 'Save file' },
      { key: 'Ctrl+P', description: 'Quick open file' },
    ]
  },
  {
    name: 'Tasks',
    shortcuts: [
      { key: 'Ctrl+Enter', description: 'Analyze task' },
      { key: 'Ctrl+B', description: 'Build context' },
      { key: 'Ctrl+Shift+A', description: 'New task' },
      { key: 'Ctrl+Shift+D', description: 'Duplicate task' },
    ]
  },
  {
    name: 'Navigation',
    shortcuts: [
      { key: 'Ctrl+1', description: 'Focus file explorer' },
      { key: 'Ctrl+2', description: 'Focus task panel' },
      { key: 'Ctrl+3', description: 'Focus output' },
      { key: 'Alt+←', description: 'Go back' },
      { key: 'Alt+→', description: 'Go forward' },
    ]
  },
  {
    name: 'View',
    shortcuts: [
      { key: 'Ctrl++', description: 'Zoom in' },
      { key: 'Ctrl+-', description: 'Zoom out' },
      { key: 'Ctrl+0', description: 'Reset zoom' },
      { key: 'F11', description: 'Toggle fullscreen' },
    ]
  },
  {
    name: 'Window',
    shortcuts: [
      { key: 'Ctrl+R', description: 'Reload window' },
      { key: 'Ctrl+,', description: 'Open settings' },
    ]
  }
]
</script>