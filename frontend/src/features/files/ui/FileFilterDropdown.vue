<template>
  <div class="relative">
    <button @click="isOpen = !isOpen" class="input flex items-center justify-between cursor-pointer">
      <span class="flex items-center gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
        </svg>
        <span v-if="selected.length === 0">{{ t('filter.byType') }}</span>
        <span v-else class="text-indigo-400">{{ selected.length }} {{ selected.length > 1 ? t('filter.types') :
          t('filter.type') }}</span>
      </span>
      <svg class="w-4 h-4 transition-transform" :class="{ 'rotate-180': isOpen }" fill="currentColor"
        viewBox="0 0 20 20">
        <path fill-rule="evenodd"
          d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
          clip-rule="evenodd" />
      </svg>
    </button>

    <!-- Dropdown -->
    <div v-if="isOpen"
      class="absolute top-full left-0 right-0 mt-1 bg-gray-800 border border-gray-700 rounded shadow-lg z-10 max-h-80 overflow-hidden flex flex-col">
      <!-- Search -->
      <div class="p-2 border-b border-gray-700">
        <input v-model="searchQuery" type="text" :placeholder="t('filter.searchExtensions')"
          class="w-full px-2 py-1 bg-gray-900 border border-gray-700 rounded text-white text-xs focus:outline-none focus:border-blue-500" />
      </div>

      <!-- Actions -->
      <div class="flex gap-2 p-2 border-b border-gray-700">
        <button @click="selectAll" class="btn btn-ghost btn-xs flex-1">
          {{ t('filter.selectAll') }}
        </button>
        <button @click="clearAll" class="btn btn-ghost btn-xs flex-1">
          {{ t('filter.clear') }}
        </button>
      </div>

      <!-- Extension Groups -->
      <div class="overflow-y-auto flex-1">
        <!-- Code Files -->
        <div v-if="codeExtensions.length > 0" class="p-2 border-b border-gray-700">
          <div class="text-xs font-semibold text-gray-400 mb-2">{{ t('filter.code') }}</div>
          <div class="space-y-1">
            <label v-for="ext in codeExtensions" :key="ext"
              class="flex items-center gap-2 px-2 py-1 hover:bg-gray-700 rounded cursor-pointer text-xs">
              <input type="checkbox" :checked="selected.includes(ext)" @change="toggleExtension(ext)"
                class="w-3 h-3 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500" />
              <span class="text-white">{{ ext }}</span>
              <span class="text-gray-400 ml-auto">{{ getExtensionCount(ext) }}</span>
            </label>
          </div>
        </div>

        <!-- Styles -->
        <div v-if="styleExtensions.length > 0" class="p-2 border-b border-gray-700">
          <div class="text-xs font-semibold text-gray-400 mb-2">{{ t('filter.styles') }}</div>
          <div class="space-y-1">
            <label v-for="ext in styleExtensions" :key="ext"
              class="flex items-center gap-2 px-2 py-1 hover:bg-gray-700 rounded cursor-pointer text-xs">
              <input type="checkbox" :checked="selected.includes(ext)" @change="toggleExtension(ext)"
                class="w-3 h-3 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500" />
              <span class="text-white">{{ ext }}</span>
              <span class="text-gray-400 ml-auto">{{ getExtensionCount(ext) }}</span>
            </label>
          </div>
        </div>

        <!-- Config Files -->
        <div v-if="configExtensions.length > 0" class="p-2 border-b border-gray-700">
          <div class="text-xs font-semibold text-gray-400 mb-2">{{ t('filter.config') }}</div>
          <div class="space-y-1">
            <label v-for="ext in configExtensions" :key="ext"
              class="flex items-center gap-2 px-2 py-1 hover:bg-gray-700 rounded cursor-pointer text-xs">
              <input type="checkbox" :checked="selected.includes(ext)" @change="toggleExtension(ext)"
                class="w-3 h-3 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500" />
              <span class="text-white">{{ ext }}</span>
              <span class="text-gray-400 ml-auto">{{ getExtensionCount(ext) }}</span>
            </label>
          </div>
        </div>

        <!-- Documentation -->
        <div v-if="docExtensions.length > 0" class="p-2 border-b border-gray-700">
          <div class="text-xs font-semibold text-gray-400 mb-2">{{ t('filter.documentation') }}</div>
          <div class="space-y-1">
            <label v-for="ext in docExtensions" :key="ext"
              class="flex items-center gap-2 px-2 py-1 hover:bg-gray-700 rounded cursor-pointer text-xs">
              <input type="checkbox" :checked="selected.includes(ext)" @change="toggleExtension(ext)"
                class="w-3 h-3 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500" />
              <span class="text-white">{{ ext }}</span>
              <span class="text-gray-400 ml-auto">{{ getExtensionCount(ext) }}</span>
            </label>
          </div>
        </div>

        <!-- Other -->
        <div v-if="otherExtensions.length > 0" class="p-2">
          <div class="text-xs font-semibold text-gray-400 mb-2">{{ t('filter.other') }}</div>
          <div class="space-y-1">
            <label v-for="ext in otherExtensions" :key="ext"
              class="flex items-center gap-2 px-2 py-1 hover:bg-gray-700 rounded cursor-pointer text-xs">
              <input type="checkbox" :checked="selected.includes(ext)" @change="toggleExtension(ext)"
                class="w-3 h-3 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500" />
              <span class="text-white">{{ ext }}</span>
              <span class="text-gray-400 ml-auto">{{ getExtensionCount(ext) }}</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, ref, watch } from 'vue'

const { t } = useI18n()

interface Props {
  extensions: string[]
  selected: string[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:selected', value: string[]): void
}>()

const isOpen = ref(false)
const searchQuery = ref('')

const filteredExtensions = computed(() => {
  if (!searchQuery.value) return props.extensions
  return props.extensions.filter(ext =>
    ext.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const codeExtensions = computed(() => {
  const codeExts = ['.ts', '.js', '.tsx', '.jsx', '.vue', '.go', '.py', '.java', '.cpp', '.c', '.rs', '.php', '.rb', '.swift', '.kt']
  return filteredExtensions.value.filter(ext => codeExts.includes(ext))
})

const styleExtensions = computed(() => {
  const styleExts = ['.css', '.scss', '.sass', '.less', '.styl']
  return filteredExtensions.value.filter(ext => styleExts.includes(ext))
})

const configExtensions = computed(() => {
  const configExts = ['.json', '.yaml', '.yml', '.toml', '.xml', '.ini', '.env']
  return filteredExtensions.value.filter(ext => configExts.includes(ext))
})

const docExtensions = computed(() => {
  const docExts = ['.md', '.txt', '.rst', '.adoc']
  return filteredExtensions.value.filter(ext => docExts.includes(ext))
})

const otherExtensions = computed(() => {
  const categorized = [
    ...codeExtensions.value,
    ...styleExtensions.value,
    ...configExtensions.value,
    ...docExtensions.value
  ]
  return filteredExtensions.value.filter(ext => !categorized.includes(ext))
})

function toggleExtension(ext: string) {
  const newSelected = props.selected.includes(ext)
    ? props.selected.filter(e => e !== ext)
    : [...props.selected, ext]
  emit('update:selected', newSelected)
}

function selectAll() {
  emit('update:selected', [...props.extensions])
}

function clearAll() {
  emit('update:selected', [])
  isOpen.value = false
}

function getExtensionCount(_ext: string): number {
  // This would ideally come from the store, but for now return placeholder
  return 0
}

// Close dropdown when clicking outside
watch(isOpen, (value) => {
  if (value) {
    const closeHandler = (e: MouseEvent) => {
      const target = e.target as HTMLElement
      if (!target.closest('.relative')) {
        isOpen.value = false
        document.removeEventListener('click', closeHandler)
      }
    }
    setTimeout(() => {
      document.addEventListener('click', closeHandler)
    }, 0)
  }
})
</script>
