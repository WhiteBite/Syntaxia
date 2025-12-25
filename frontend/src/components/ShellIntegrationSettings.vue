<template>
  <div class="space-y-4">
    <div class="flex items-start gap-3 p-4 rounded-lg bg-gray-800/50 border border-gray-700/30">
      <Monitor class="w-5 h-5 text-indigo-400 mt-0.5 flex-shrink-0" />
      <div class="flex-1 min-w-0">
        <h3 class="text-sm font-medium text-white mb-1">
          {{ t('settings.shellIntegration.title') }}
        </h3>
        <p class="text-xs text-gray-400 mb-3">
          {{ t('settings.shellIntegration.description') }}
        </p>

        <div class="flex items-center gap-3">
          <button
            v-if="!isRegistered"
            @click="handleRegister"
            :disabled="isLoading"
            class="btn-unified btn-unified-primary text-sm"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin" />
            <Plus v-else class="w-4 h-4" />
            {{ t('settings.shellIntegration.enable') }}
          </button>

          <button
            v-else
            @click="handleUnregister"
            :disabled="isLoading"
            class="btn-unified btn-unified-secondary text-sm"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin" />
            <Trash2 v-else class="w-4 h-4" />
            {{ t('settings.shellIntegration.disable') }}
          </button>

          <span
            class="text-xs px-2 py-1 rounded"
            :class="isRegistered ? 'bg-green-500/20 text-green-400' : 'bg-gray-600/30 text-gray-400'"
          >
            {{ isRegistered ? t('settings.shellIntegration.enabled') : t('settings.shellIntegration.disabled') }}
          </span>
        </div>

        <p v-if="currentOS === 'windows'" class="text-xs text-gray-500 mt-2">
          {{ t('settings.shellIntegration.requiresAdmin') }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { shellApi } from '@/services/api/shell.api'
import { useUIStore } from '@/stores/ui.store'
import { Loader2, Monitor, Plus, Trash2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'

const { t } = useI18n()
const uiStore = useUIStore()

const isRegistered = ref(false)
const isLoading = ref(false)
const currentOS = ref('')

async function loadStatus() {
  try {
    const status = await shellApi.getStatus()
    isRegistered.value = status.isRegistered
    currentOS.value = status.currentOS
  } catch {
    uiStore.addToast(t('settings.shellIntegration.error'), 'error')
  }
}

async function handleRegister() {
  isLoading.value = true
  try {
    await shellApi.register()
    isRegistered.value = true
    uiStore.addToast(t('settings.shellIntegration.enableSuccess'), 'success')
  } catch {
    uiStore.addToast(t('settings.shellIntegration.error'), 'error')
  } finally {
    isLoading.value = false
  }
}

async function handleUnregister() {
  isLoading.value = true
  try {
    await shellApi.unregister()
    isRegistered.value = false
    uiStore.addToast(t('settings.shellIntegration.disableSuccess'), 'success')
  } catch {
    uiStore.addToast(t('settings.shellIntegration.error'), 'error')
  } finally {
    isLoading.value = false
  }
}

onMounted(loadStatus)
</script>
