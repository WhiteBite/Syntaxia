<template>
  <div class="relative">
    <button
      @click="isOpen = !isOpen"
      class="action-btn text-gray-400 hover:text-white"
      :title="t('git.recentRepos')"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </button>

    <!-- Dropdown -->
    <div
      v-if="isOpen"
      class="absolute right-0 top-full mt-1 w-64 bg-gray-800 border border-gray-700 rounded-lg shadow-xl z-50"
    >
      <div class="p-2 border-b border-gray-700 flex items-center justify-between">
        <span class="text-xs text-gray-400">{{ t('git.recentRepos') }}</span>
        <button
          @click="handleClear"
          class="text-xs text-gray-400 hover:text-red-400"
        >
          {{ t('common.clear') }}
        </button>
      </div>
      <div class="max-h-48 overflow-y-auto">
        <button
          v-for="repo in repos"
          :key="repo.url"
          @click="handleSelect(repo)"
          class="w-full px-3 py-2 text-left hover:bg-gray-700/50 flex items-center gap-2"
        >
          <svg v-if="repo.isGitHub" class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
          <svg v-else class="w-4 h-4 text-orange-400" fill="currentColor" viewBox="0 0 24 24">
            <path d="M22.65 14.39L12 22.13 1.35 14.39a.84.84 0 01-.3-.94l1.22-3.78 2.44-7.51A.42.42 0 014.82 2a.43.43 0 01.58 0 .42.42 0 01.11.18l2.44 7.49h8.1l2.44-7.51A.42.42 0 0118.6 2a.43.43 0 01.58 0 .42.42 0 01.11.18l2.44 7.51L23 13.45a.84.84 0 01-.35.94z"/>
          </svg>
          <span class="text-sm text-gray-300 truncate">{{ repo.name }}</span>
        </button>
      </div>
    </div>

    <!-- Click outside to close -->
    <div v-if="isOpen" class="fixed inset-0 z-40" @click="isOpen = false" />
  </div>
</template>

<script setup lang="ts">
import type { RecentRepo } from '@/composables/useGitSource';
import { useI18n } from '@/composables/useI18n';
import { ref } from 'vue';

const { t } = useI18n()

defineProps<{
  repos: RecentRepo[]
}>()

const emit = defineEmits<{
  select: [repo: RecentRepo]
  clear: []
}>()

const isOpen = ref(false)

function handleSelect(repo: RecentRepo) {
  emit('select', repo)
  isOpen.value = false
}

function handleClear() {
  emit('clear')
  isOpen.value = false
}
</script>
