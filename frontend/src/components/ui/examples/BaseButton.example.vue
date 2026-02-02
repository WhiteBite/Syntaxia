<template>
  <div class="examples-container">
    <h1>BaseButton Examples</h1>

    <!-- Variants -->
    <section class="example-section">
      <h2>Варианты (Variants)</h2>
      <div class="button-grid">
        <BaseButton variant="primary">Primary</BaseButton>
        <BaseButton variant="secondary">Secondary</BaseButton>
        <BaseButton variant="ghost">Ghost</BaseButton>
        <BaseButton variant="danger">Danger</BaseButton>
        <BaseButton variant="success">Success</BaseButton>
        <BaseButton variant="warning">Warning</BaseButton>
        <BaseButton variant="link">Link</BaseButton>
        <BaseButton variant="text">Text</BaseButton>
      </div>
    </section>

    <!-- Sizes -->
    <section class="example-section">
      <h2>Размеры (Sizes)</h2>
      <div class="button-grid">
        <BaseButton size="xs" variant="primary">Extra Small</BaseButton>
        <BaseButton size="sm" variant="primary">Small</BaseButton>
        <BaseButton size="md" variant="primary">Medium</BaseButton>
        <BaseButton size="lg" variant="primary">Large</BaseButton>
      </div>
    </section>

    <!-- With Icons -->
    <section class="example-section">
      <h2>С иконками</h2>
      <div class="button-grid">
        <BaseButton variant="primary" :icon="PlayIcon">
          Запустить
        </BaseButton>
        <BaseButton variant="secondary" :icon="SettingsIcon" icon-position="right">
          Настройки
        </BaseButton>
        <BaseButton variant="danger" :icon="TrashIcon">
          Удалить
        </BaseButton>
        <BaseButton variant="success" :icon="CheckIcon">
          Готово
        </BaseButton>
      </div>
    </section>

    <!-- Icon Only -->
    <section class="example-section">
      <h2>Только иконка</h2>
      <div class="button-grid">
        <BaseButton variant="ghost" icon-only :icon="SettingsIcon" />
        <BaseButton variant="ghost" icon-only :icon="TrashIcon" />
        <BaseButton variant="ghost" icon-only :icon="PencilIcon" />
        <BaseButton variant="ghost" icon-only :icon="EllipsisIcon" />
      </div>
    </section>

    <!-- Loading State -->
    <section class="example-section">
      <h2>Состояние загрузки</h2>
      <div class="button-grid">
        <BaseButton variant="primary" :loading="true">
          Загрузка...
        </BaseButton>
        <BaseButton variant="secondary" :loading="true" :icon="SaveIcon">
          Сохранение...
        </BaseButton>
        <BaseButton variant="success" :loading="isLoading" @click="simulateLoading">
          Кликни меня
        </BaseButton>
      </div>
    </section>

    <!-- Disabled State -->
    <section class="example-section">
      <h2>Отключенное состояние</h2>
      <div class="button-grid">
        <BaseButton variant="primary" :disabled="true">
          Disabled Primary
        </BaseButton>
        <BaseButton variant="secondary" :disabled="true">
          Disabled Secondary
        </BaseButton>
        <BaseButton variant="danger" :disabled="true" :icon="TrashIcon">
          Disabled Danger
        </BaseButton>
      </div>
    </section>

    <!-- Outline Variants -->
    <section class="example-section">
      <h2>Outline варианты</h2>
      <div class="button-grid">
        <BaseButton variant="primary" :outline="true">
          Primary Outline
        </BaseButton>
        <BaseButton variant="secondary" :outline="true">
          Secondary Outline
        </BaseButton>
        <BaseButton variant="danger" :outline="true">
          Danger Outline
        </BaseButton>
        <BaseButton variant="success" :outline="true">
          Success Outline
        </BaseButton>
      </div>
    </section>

    <!-- Block (Full Width) -->
    <section class="example-section">
      <h2>Полная ширина (Block)</h2>
      <div class="flex flex-col gap-2">
        <BaseButton variant="primary" :block="true">
          Full Width Primary
        </BaseButton>
        <BaseButton variant="secondary" :block="true" :icon="SettingsIcon">
          Full Width with Icon
        </BaseButton>
      </div>
    </section>

    <!-- Real Use Cases -->
    <section class="example-section">
      <h2>Реальные примеры использования</h2>
      
      <!-- Form Actions -->
      <div class="use-case">
        <h3>Действия формы</h3>
        <div class="flex gap-2">
          <BaseButton variant="primary" :icon="CheckIcon">
            Сохранить
          </BaseButton>
          <BaseButton variant="ghost">
            Отмена
          </BaseButton>
        </div>
      </div>

      <!-- Destructive Action -->
      <div class="use-case">
        <h3>Опасное действие</h3>
        <div class="flex gap-2">
          <BaseButton variant="danger" :icon="TrashIcon">
            Удалить файл
          </BaseButton>
          <BaseButton variant="ghost">
            Отмена
          </BaseButton>
        </div>
      </div>

      <!-- Toolbar Actions -->
      <div class="use-case">
        <h3>Панель инструментов</h3>
        <div class="flex gap-1">
          <BaseButton variant="ghost" icon-only :icon="PencilIcon" size="sm" />
          <BaseButton variant="ghost" icon-only :icon="TrashIcon" size="sm" />
          <BaseButton variant="ghost" icon-only :icon="SettingsIcon" size="sm" />
          <BaseButton variant="ghost" icon-only :icon="EllipsisIcon" size="sm" />
        </div>
      </div>

      <!-- Loading Action -->
      <div class="use-case">
        <h3>Асинхронное действие</h3>
        <BaseButton 
          variant="primary" 
          :loading="isSaving" 
          :icon="SaveIcon"
          @click="handleSave"
        >
          {{ isSaving ? 'Сохранение...' : 'Сохранить изменения' }}
        </BaseButton>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import BaseButton from '../BaseButton.vue'
import {
  PlayIcon,
  Cog6ToothIcon as SettingsIcon,
  TrashIcon,
  CheckIcon,
  PencilIcon,
  EllipsisHorizontalIcon as EllipsisIcon,
  ArrowDownTrayIcon as SaveIcon
} from '@heroicons/vue/24/outline'

const isLoading = ref(false)
const isSaving = ref(false)

function simulateLoading() {
  isLoading.value = true
  setTimeout(() => {
    isLoading.value = false
  }, 2000)
}

function handleSave() {
  isSaving.value = true
  setTimeout(() => {
    isSaving.value = false
  }, 1500)
}
</script>

<style scoped>
.examples-container {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

h1 {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 2rem;
  color: var(--text-primary);
}

.example-section {
  margin-bottom: 3rem;
  padding: 1.5rem;
  background: var(--bg-1);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
}

.example-section h2 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: var(--text-primary);
}

.button-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  align-items: center;
}

.use-case {
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: var(--bg-2);
  border-radius: var(--radius-md);
}

.use-case h3 {
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.75rem;
  color: var(--text-muted);
}
</style>
