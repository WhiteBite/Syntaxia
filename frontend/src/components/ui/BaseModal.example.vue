<template>
  <div>
    <!-- Example 1: Basic Modal -->
    <BaseButton @click="basicModal.open()">Open Basic Modal</BaseButton>
    <BaseModal v-model="basicModal.isOpen.value" title="Basic Modal">
      <p>This is a basic modal with default settings.</p>
    </BaseModal>

    <!-- Example 2: Modal with Footer -->
    <BaseButton @click="footerModal.open()">Open Modal with Footer</BaseButton>
    <BaseModal v-model="footerModal.isOpen.value" title="Confirm Action">
      <p>Are you sure you want to proceed with this action?</p>
      <template #footer>
        <BaseButton variant="ghost" @click="footerModal.close()">Cancel</BaseButton>
        <BaseButton variant="primary" @click="handleConfirm">Confirm</BaseButton>
      </template>
    </BaseModal>

    <!-- Example 3: Large Modal -->
    <BaseButton @click="largeModal.open()">Open Large Modal</BaseButton>
    <BaseModal v-model="largeModal.isOpen.value" title="Large Modal" size="lg">
      <p>This is a large modal with more content space.</p>
      <div class="space-y-4">
        <div v-for="i in 10" :key="i" class="p-4 bg-gray-800 rounded">
          Content block {{ i }}
        </div>
      </div>
    </BaseModal>

    <!-- Example 4: Persistent Modal -->
    <BaseButton @click="persistentModal.open()">Open Persistent Modal</BaseButton>
    <BaseModal
      v-model="persistentModal.isOpen.value"
      title="Persistent Modal"
      :persistent="true"
      :close-on-backdrop="false"
      :close-on-esc="false"
    >
      <p>This modal cannot be closed by clicking outside or pressing ESC.</p>
      <p>You must click the close button or the cancel button below.</p>
      <template #footer>
        <BaseButton @click="persistentModal.close()">Close</BaseButton>
      </template>
    </BaseModal>

    <!-- Example 5: Custom Header -->
    <BaseButton @click="customModal.open()">Open Custom Header Modal</BaseButton>
    <BaseModal v-model="customModal.isOpen.value" :show-close="false">
      <template #header>
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-purple-500/20 flex items-center justify-center">
            <span class="text-purple-400">🎉</span>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-white">Custom Header</h3>
            <p class="text-sm text-gray-400">With custom styling</p>
          </div>
        </div>
      </template>
      <p>This modal has a custom header with icon and description.</p>
      <template #footer>
        <BaseButton @click="customModal.close()">Got it</BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { BaseModal, BaseButton } from '@/components/ui'
import { useModal } from '@/composables/useModal'

// Create modal instances
const basicModal = useModal()
const footerModal = useModal()
const largeModal = useModal()
const persistentModal = useModal()
const customModal = useModal()

function handleConfirm() {
  // Handle confirmation logic
  footerModal.close()
}
</script>
