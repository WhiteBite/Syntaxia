<template>
  <div class="border-b border-gray-700/30">
    <div class="flex items-center justify-between p-3">
      <div class="section-title">
        <div class="section-icon section-icon-orange">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
          </svg>
        </div>
        <h2 class="section-title-text">{{ t('git.title') }}</h2>
      </div>
      <!-- Recent repos dropdown -->
      <RecentReposDropdown 
        v-if="recentRepos.length > 0"
        :repos="recentRepos"
        @select="$emit('select-recent', $event)"
        @clear="$emit('clear-recent')"
      />
    </div>

    <!-- Source Type Tabs -->
    <GitSourceTabs 
      :model-value="sourceType" 
      @change="$emit('change-source', $event)" 
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { RecentRepo } from '@/composables/useGitSource'
import RecentReposDropdown from './RecentReposDropdown.vue'
import GitSourceTabs, { type SourceType } from './GitSourceTabs.vue'

defineProps<{
  sourceType: SourceType
  recentRepos: RecentRepo[]
}>()

defineEmits<{
  'change-source': [value: SourceType]
  'select-recent': [repo: RecentRepo]
  'clear-recent': []
}>()

const { t } = useI18n()
</script>
