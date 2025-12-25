<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between">
      <p class="text-xs font-semibold text-gray-400">{{ t('history.title') }}</p>
      <span class="text-xs text-gray-400">{{ contextList.length }}</span>
    </div>

    <!-- History List -->
    <div v-if="contextList.length > 0" class="space-y-2">
      <div
        v-for="ctx in sortedContexts"
        :key="ctx.id"
        class="list-item group !p-2 cursor-pointer"
        :class="{ 'list-item-active': ctx.id === currentContextId }"
        @click="$emit('load', ctx.id)"
      >
        <div class="flex items-start gap-2">
          <div class="section-icon section-icon-purple !w-7 !h-7 flex-shrink-0 mt-0.5">
            <svg class="!w-3.5 !h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span class="text-sm text-white truncate">{{ ctx.name }}</span>
              <button
                v-if="ctx.isFavorite"
                class="text-yellow-400"
                @click.stop="toggleFavorite(ctx.id)"
              >
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </button>
            </div>
            <div class="flex items-center gap-2 text-xs text-gray-400 mt-0.5">
              <span>{{ ctx.fileCount }} {{ t('context.files') }}</span>
              <span>•</span>
              <span>{{ formatNumber(ctx.tokenCount || 0) }} tok</span>
              <span>•</span>
              <span>{{ formatTimeAgo(ctx.createdAt) }}</span>
            </div>
          </div>
          <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              class="icon-btn-sm"
              :title="ctx.isFavorite ? t('context.unfavorite') : t('context.favorite')"
              :aria-label="ctx.isFavorite ? t('context.unfavorite') : t('context.favorite')"
              @click.stop="toggleFavorite(ctx.id)"
            >
              <svg class="w-3.5 h-3.5" :class="ctx.isFavorite ? 'text-yellow-400' : 'text-gray-400'" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
              </svg>
            </button>
            <button
              class="icon-btn-sm icon-btn-danger"
              :title="t('history.delete')"
              :aria-label="t('history.delete')"
              @click.stop="deleteContext(ctx.id)"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state py-8">
      <div class="empty-state-icon !w-12 !h-12 mb-3">
        <svg class="!w-6 !h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <p class="empty-state-text">{{ t('history.empty') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context'
import { computed, onMounted } from 'vue'

const { t } = useI18n()
const contextStore = useContextStore()

defineEmits<{
  (e: 'load', contextId: string): void
}>()

const contextList = computed(() => contextStore.contextList)
const currentContextId = computed(() => contextStore.contextId)

const sortedContexts = computed(() => {
  return [...contextList.value].sort((a, b) => {
    // Favorites first
    if (a.isFavorite && !b.isFavorite) return -1
    if (!a.isFavorite && b.isFavorite) return 1
    // Then by date
    const dateA = new Date(a.createdAt || 0).getTime()
    const dateB = new Date(b.createdAt || 0).getTime()
    return dateB - dateA
  })
})

function toggleFavorite(ctxId: string) {
  contextStore.toggleFavorite(ctxId)
}

async function deleteContext(ctxId: string) {
  await contextStore.deleteContext(ctxId)
}

function formatNumber(num: number): string {
  if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`
  if (num >= 1000) return `${(num / 1000).toFixed(1)}K`
  return num.toString()
}

function formatTimeAgo(dateStr?: string): string {
  if (!dateStr) return ''
  
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return t('history.justNow')
  if (diffMins < 60) return `${diffMins} ${t('history.minutesAgo')}`
  if (diffHours < 24) return `${diffHours} ${t('history.hoursAgo')}`
  return `${diffDays} ${t('history.daysAgo')}`
}

onMounted(() => {
  contextStore.listProjectContexts()
})
</script>
