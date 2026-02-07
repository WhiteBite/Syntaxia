<template>
  <div class="examples-container p-6">
    <h1 class="text-2xl font-bold mb-6">BaseIcon Examples</h1>

    <!-- Sizes -->
    <section class="example-section mb-10">
      <h2 class="text-xl font-semibold mb-4">Размеры (Sizes)</h2>
      <div class="flex flex-wrap gap-6 items-end">
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="StarIcon" size="xs" />
          <span class="text-xs text-gray-500">xs (12px)</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="StarIcon" size="sm" />
          <span class="text-xs text-gray-500">sm (16px)</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="StarIcon" size="md" />
          <span class="text-xs text-gray-500">md (20px)</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="StarIcon" size="lg" />
          <span class="text-xs text-gray-500">lg (24px)</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="StarIcon" size="xl" />
          <span class="text-xs text-gray-500">xl (32px)</span>
        </div>
      </div>
    </section>

    <!-- Colors -->
    <section class="example-section mb-10">
      <h2 class="text-xl font-semibold mb-4">Цвета</h2>
      <div class="flex flex-wrap gap-6">
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="HeartIcon" color="var(--color-danger)" size="lg" />
          <span class="text-xs text-gray-500">Danger</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="CheckCircleIcon" color="var(--color-success)" size="lg" />
          <span class="text-xs text-gray-500">Success</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="ExclamationTriangleIcon" color="var(--color-warning)" size="lg" />
          <span class="text-xs text-gray-500">Warning</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="InformationCircleIcon" color="var(--color-info)" size="lg" />
          <span class="text-xs text-gray-500">Info</span>
        </div>
        <div class="flex flex-col items-center gap-2">
          <BaseIcon :icon="SparklesIcon" color="var(--accent-purple)" size="lg" />
          <span class="text-xs text-gray-500">Accent</span>
        </div>
      </div>
    </section>

    <!-- Common Icons -->
    <section class="example-section">
      <h2 class="text-xl font-semibold mb-6">Часто используемые иконки</h2>
      
      <div class="mb-8">
        <h3 class="text-lg font-medium mb-4 text-gray-400">Действия</h3>
        <div class="flex flex-wrap gap-6">
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="PlusIcon" />
            <span class="text-xs text-gray-500">Plus</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="PencilIcon" />
            <span class="text-xs text-gray-500">Edit</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="TrashIcon" />
            <span class="text-xs text-gray-500">Delete</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="XMarkIcon" />
            <span class="text-xs text-gray-500">Close</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="CheckIcon" />
            <span class="text-xs text-gray-500">Check</span>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-medium mb-4 text-gray-400">Навигация</h3>
        <div class="flex flex-wrap gap-6">
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="ChevronLeftIcon" />
            <span class="text-xs text-gray-500">Left</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="ChevronRightIcon" />
            <span class="text-xs text-gray-500">Right</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="ChevronUpIcon" />
            <span class="text-xs text-gray-500">Up</span>
          </div>
          <div class="flex flex-col items-center gap-2">
            <BaseIcon :icon="ChevronDownIcon" />
            <span class="text-xs text-gray-500">Down</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { BaseIcon } from '@/components/ui'
import { 
  StarIcon, 
  HeartIcon, 
  CheckCircleIcon, 
  ExclamationTriangleIcon, 
  InformationCircleIcon, 
  SparklesIcon,
  PlusIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
  CheckIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronUpIcon,
  ChevronDownIcon
} from '@heroicons/vue/24/outline'
</script>
