<template>
  <button
    @click="cycleTheme"
    class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
    :title="`Current theme: ${theme}`"
  >
    <SunIcon v-if="theme === 'light'" class="w-5 h-5 text-gray-700 dark:text-gray-300" />
    <MoonIcon v-else-if="theme === 'dark'" class="w-5 h-5 text-gray-700 dark:text-gray-300" />
    <ComputerDesktopIcon v-else class="w-5 h-5 text-gray-700 dark:text-gray-300" />
  </button>
</template>

<script setup lang="ts">
import { SunIcon, MoonIcon, ComputerDesktopIcon } from '@heroicons/vue/24/outline'
import { useTheme } from '../composables/useTheme'

const { theme, cycleTheme } = useTheme()
</script>

