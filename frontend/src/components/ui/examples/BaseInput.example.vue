<template>
  <div class="examples-container">
    <h1>BaseInput Examples</h1>

    <!-- Basic Input -->
    <section class="example-section">
      <h2>Базовый ввод</h2>
      <div class="input-grid">
        <BaseInput
          v-model="basicInput"
          placeholder="Введите текст..."
        />
        <p class="value-display">Значение: {{ basicInput }}</p>
      </div>
    </section>

    <!-- With Label -->
    <section class="example-section">
      <h2>С меткой (Label)</h2>
      <div class="input-grid">
        <BaseInput
          v-model="labelInput"
          label="Имя пользователя"
          placeholder="Введите имя"
        />
        <BaseInput
          v-model="emailInput"
          label="Email"
          type="email"
          placeholder="user@example.com"
        />
      </div>
    </section>

    <!-- With Icons -->
    <section class="example-section">
      <h2>С иконками</h2>
      <div class="input-grid">
        <BaseInput
          v-model="searchInput"
          label="Поиск"
          placeholder="Поиск файлов..."
          :prefix-icon="SearchIcon"
        />
        <BaseInput
          v-model="passwordInput"
          label="Пароль"
          type="password"
          placeholder="Введите пароль"
          :prefix-icon="LockIcon"
        />
        <BaseInput
          v-model="urlInput"
          label="URL"
          placeholder="https://example.com"
          :prefix-icon="LinkIcon"
        />
      </div>
    </section>

    <!-- With Suffix Icon -->
    <section class="example-section">
      <h2>С иконкой справа</h2>
      <div class="input-grid">
        <BaseInput
          v-model="clearableInput"
          label="Очищаемое поле"
          placeholder="Введите текст"
        >
          <template #suffix>
            <button
              v-if="clearableInput"
              @click="clearableInput = ''"
              class="clear-button"
            >
              <XMarkIcon class="w-4 h-4" />
            </button>
          </template>
        </BaseInput>
      </div>
    </section>

    <!-- With Error -->
    <section class="example-section">
      <h2>С ошибкой валидации</h2>
      <div class="input-grid">
        <BaseInput
          v-model="errorInput"
          label="Email"
          placeholder="user@example.com"
          :error="emailError"
          @blur="validateEmail"
        />
        <BaseInput
          v-model="requiredInput"
          label="Обязательное поле"
          placeholder="Не может быть пустым"
          :error="requiredError"
          @blur="validateRequired"
        />
      </div>
    </section>

    <!-- Disabled State -->
    <section class="example-section">
      <h2>Отключенное состояние</h2>
      <div class="input-grid">
        <BaseInput
          v-model="disabledInput"
          label="Отключенное поле"
          placeholder="Нельзя редактировать"
          :disabled="true"
        />
      </div>
    </section>

    <!-- Different Types -->
    <section class="example-section">
      <h2>Разные типы</h2>
      <div class="input-grid">
        <BaseInput
          v-model="textInput"
          label="Text"
          type="text"
          placeholder="Обычный текст"
        />
        <BaseInput
          v-model="numberInput"
          label="Number"
          type="number"
          placeholder="123"
        />
        <BaseInput
          v-model="dateInput"
          label="Date"
          type="date"
        />
        <BaseInput
          v-model="timeInput"
          label="Time"
          type="time"
        />
      </div>
    </section>

    <!-- Real Use Cases -->
    <section class="example-section">
      <h2>Реальные примеры использования</h2>

      <!-- Search -->
      <div class="use-case">
        <h3>Поиск файлов</h3>
        <BaseInput
          v-model="fileSearch"
          placeholder="Поиск по имени файла..."
          :prefix-icon="SearchIcon"
        >
          <template #suffix>
            <span v-if="fileSearch" class="text-xs text-gray-400">
              {{ filteredCount }} найдено
            </span>
          </template>
        </BaseInput>
      </div>

      <!-- Filter -->
      <div class="use-case">
        <h3>Фильтр по расширению</h3>
        <BaseInput
          v-model="extensionFilter"
          label="Расширение файла"
          placeholder=".ts, .vue, .js"
          :prefix-icon="FunnelIcon"
        />
      </div>

      <!-- Path Input -->
      <div class="use-case">
        <h3>Путь к проекту</h3>
        <BaseInput
          v-model="projectPath"
          label="Путь к проекту"
          placeholder="/path/to/project"
          :prefix-icon="FolderIcon"
        >
          <template #suffix>
            <button @click="browsePath" class="browse-button">
              Обзор
            </button>
          </template>
        </BaseInput>
      </div>

      <!-- API Key -->
      <div class="use-case">
        <h3>API ключ</h3>
        <BaseInput
          v-model="apiKey"
          label="OpenAI API Key"
          :type="showApiKey ? 'text' : 'password'"
          placeholder="sk-..."
          :prefix-icon="KeyIcon"
        >
          <template #suffix>
            <button @click="showApiKey = !showApiKey" class="toggle-button">
              <EyeIcon v-if="!showApiKey" class="w-4 h-4" />
              <EyeSlashIcon v-else class="w-4 h-4" />
            </button>
          </template>
        </BaseInput>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseInput from '../BaseInput.vue'
import {
  MagnifyingGlassIcon as SearchIcon,
  LockClosedIcon as LockIcon,
  LinkIcon,
  XMarkIcon,
  FunnelIcon,
  FolderIcon,
  KeyIcon,
  EyeIcon,
  EyeSlashIcon
} from '@heroicons/vue/24/outline'

// Basic
const basicInput = ref('')
const labelInput = ref('')
const emailInput = ref('')

// With Icons
const searchInput = ref('')
const passwordInput = ref('')
const urlInput = ref('')

// Clearable
const clearableInput = ref('')

// With Error
const errorInput = ref('')
const emailError = ref('')
const requiredInput = ref('')
const requiredError = ref('')

// Disabled
const disabledInput = ref('Это значение нельзя изменить')

// Different Types
const textInput = ref('')
const numberInput = ref('')
const dateInput = ref('')
const timeInput = ref('')

// Real Use Cases
const fileSearch = ref('')
const filteredCount = computed(() => Math.floor(Math.random() * 50) + 1)
const extensionFilter = ref('')
const projectPath = ref('')
const apiKey = ref('')
const showApiKey = ref(false)

function validateEmail() {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (errorInput.value && !emailRegex.test(errorInput.value)) {
    emailError.value = 'Неверный формат email'
  } else {
    emailError.value = ''
  }
}

function validateRequired() {
  if (!requiredInput.value.trim()) {
    requiredError.value = 'Это поле обязательно для заполнения'
  } else {
    requiredError.value = ''
  }
}

function browsePath() {
  // Simulate file browser
  projectPath.value = '/home/user/projects/syntaxia'
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

.input-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}

.value-display {
  padding: 0.75rem;
  background: var(--bg-2);
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
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

.clear-button,
.toggle-button {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.25rem;
  color: var(--text-muted);
  transition: color var(--transition-fast);
  cursor: pointer;
  background: transparent;
  border: none;
}

.clear-button:hover,
.toggle-button:hover {
  color: var(--text-primary);
}

.browse-button {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  color: var(--accent-indigo);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: color var(--transition-fast);
}

.browse-button:hover {
  color: var(--accent-purple);
}
</style>
