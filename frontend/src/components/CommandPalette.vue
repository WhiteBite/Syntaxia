<template>
<TransitionRoot :show="isOpen" as="template">
<Dialog @close="close" class="relative z-50">
  <TransitionChild
      as="template"
      enter="ease-out duration-200"
      enter-from="opacity-0"
      enter-to="opacity-100"
      leave="ease-in duration-150"
      leave-from="opacity-100"
      leave-to="opacity-0"
  >
    <div class="fixed inset-0 bg-black/50 backdrop-blur-sm"/>
  </TransitionChild>

  <div class="fixed inset-0 overflow-y-auto">
    <div class="flex min-h-full items-start justify-center p-4 pt-[15vh]">
      <TransitionChild
          as="template"
          enter="ease-out duration-200"
          enter-from="opacity-0 scale-95"
          enter-to="opacity-100 scale-100"
          leave="ease-in duration-150"
          leave-from="opacity-100 scale-100"
          leave-to="opacity-0 scale-95"
      >
        <DialogPanel
            class="w-full max-w-2xl transform overflow-hidden rounded-xl bg-white dark:bg-gray-800 shadow-2xl transition-all">
          <div class="relative">
            <input
                ref="searchInput"
                v-model="query"
                type="text"
                placeholder="Type a command or search..."
                class="w-full border-0 bg-transparent px-4 py-4 text-gray-900 dark:text-gray-100 placeholder:text-gray-400 focus:ring-0 text-lg"
                @keydown.down.prevent="selectNext"
                @keydown.up.prevent="selectPrevious"
                @keydown.enter.prevent="executeSelected"
            />
          </div>

          <div v-if="filteredCommands.length > 0"
               class="border-t border-gray-200 dark:border-gray-700">
            <ul class="max-h-96 overflow-y-auto py-2">
              <li
                  v-for="(command, index) in filteredCommands"
                  :key="command.id"
                  :class="[
                      'px-4 py-3 cursor-pointer flex items-center justify-between transition-colors',
                      selectedIndex === index
                        ? 'bg-blue-50 dark:bg-blue-900/20'
                        : 'hover:bg-gray-50 dark:hover:bg-gray-700/50'
                    ]"
                  @click="execute(command)"
                  @mouseenter="selectedIndex = index"
              >
                <div class="flex items-center gap-3">
                  <component :is="command.icon" class="w-5 h-5 text-gray-400"/>
                  <div>
                    <div class="text-sm font-medium text-gray-900 dark:text-gray-100">
                      {{ command.name }}
                    </div>
                    <div class="text-xs text-gray-400 dark:text-gray-400">
                      {{ command.description }}
                    </div>
                  </div>
                </div>
                <kbd
                    class="px-2 py-1 text-xs font-semibold text-gray-400 bg-gray-100 dark:bg-gray-700 dark:text-gray-400 rounded">
                  {{ command.shortcut }}
                </kbd>
              </li>
            </ul>
          </div>

          <div v-else class="px-4 py-8 text-center text-sm text-gray-400 dark:text-gray-400">
            No commands found
          </div>
        </DialogPanel>
      </TransitionChild>
    </div>
  </div>
</Dialog>
</TransitionRoot>
</template>

<script setup lang="ts">
import { useUIStore } from '@/stores/ui.store'
import { Dialog, DialogPanel, TransitionChild, TransitionRoot } from '@headlessui/vue'
import {
    ArrowPathIcon,
    CogIcon,
    DocumentTextIcon,
    FolderOpenIcon,
    MagnifyingGlassIcon,
    QuestionMarkCircleIcon
} from '@heroicons/vue/24/outline'
import Fuse from 'fuse.js'
import type { Component } from 'vue'
import { computed, nextTick, ref, watch } from 'vue'

interface Command {
  id: string
  name: string
  description: string
  shortcut: string
  icon: Component
  action: () => void
}

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const uiStore = useUIStore()

const query = ref('')
const searchInput = ref<HTMLInputElement>()
const selectedIndex = ref(0)

const commands = ref<Command[]>([
  {
    id: '1',
    name: 'Search Files',
    description: 'Search for files in current project',
    shortcut: 'Ctrl+P',
    icon: MagnifyingGlassIcon,
    action: () => uiStore.addToast('File search: use the file tree filter', 'info')
  },
  {
    id: '2',
    name: 'New Document',
    description: 'Create a new document',
    shortcut: 'Ctrl+N',
    icon: DocumentTextIcon,
    action: () => uiStore.addToast('New document: not implemented yet', 'info')
  },
  {
    id: '3',
    name: 'Settings',
    description: 'Open application settings',
    shortcut: 'Ctrl+,',
    icon: CogIcon,
    action: () => uiStore.openSettingsModal()
  },
  {
    id: '4',
    name: 'Open Project',
    description: 'Open a project folder',
    shortcut: 'Ctrl+O',
    icon: FolderOpenIcon,
    action: () => uiStore.addToast('Open project: use File menu', 'info')
  },
  {
    id: '5',
    name: 'Reload Window',
    description: 'Reload the application',
    shortcut: 'Ctrl+R',
    icon: ArrowPathIcon,
    action: () => window.location.reload()
  },
  {
    id: '6',
    name: 'Help',
    description: 'Show keyboard shortcuts',
    shortcut: 'Ctrl+/',
    icon: QuestionMarkCircleIcon,
    action: () => uiStore.openKeyboardShortcutsModal()
  },
])

const fuse = computed(() => new Fuse(commands.value, {
  keys: ['name', 'description'],
  threshold: 0.3
}))

const filteredCommands = computed(() => {
  if (!query.value) return commands.value
  return fuse.value.search(query.value).map(result => result.item)
})

const selectNext = () => {
  selectedIndex.value = Math.min(selectedIndex.value + 1, filteredCommands.value.length - 1)
}

const selectPrevious = () => {
  selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
}

const execute = (command: Command) => {
  command.action()
  close()
}

const executeSelected = () => {
  if (filteredCommands.value[selectedIndex.value]) {
    execute(filteredCommands.value[selectedIndex.value])
  }
}

const close = () => {
  query.value = ''
  selectedIndex.value = 0
  emit('close')
}

watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    nextTick(() => {
      searchInput.value?.focus()
    })
  }
})
</script>