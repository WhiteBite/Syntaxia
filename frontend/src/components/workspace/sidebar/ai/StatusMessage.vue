<template>
  <div v-if="message" :class="['status-message', type === 'success' ? 'status-message-success' : 'status-message-error']">
    <Check v-if="type === 'success'" class="w-4 h-4" />
    <X v-else class="w-4 h-4" />
    {{ message }}
  </div>
</template>

<script setup lang="ts">
import { Check, X } from 'lucide-vue-next';

defineProps<{
  message: string
  type: 'success' | 'error'
}>()
</script>
