<template>
  <div class="examples-container">
    <h1>BaseTextarea Examples</h1>

    <!-- Basic Textarea -->
    <section class="example-section">
      <h2>Базовое текстовое поле</h2>
      <BaseTextarea
        v-model="basicText"
        placeholder="Введите многострочный текст..."
      />
      <p class="char-count">Символов: {{ basicText.length }}</p>
    </section>

    <!-- With Label -->
    <section class="example-section">
      <h2>С меткой (Label)</h2>
      <BaseTextarea
        v-model="labelText"
        label="Описание"
        placeholder="Введите описание проекта..."
      />
    </section>

    <!-- With Error -->
    <section class="example-section">
      <h2>С ошибкой валидации</h2>
      <BaseTextarea
        v-model="errorText"
        label="Комментарий (минимум 10 символов)"
        placeholder="Введите комментарий..."
        :error="commentError"
        @blur="validateComment"
      />
    </section>

    <!-- Disabled State -->
    <section class="example-section">
      <h2>Отключенное состояние</h2>
      <BaseTextarea
        v-model="disabledText"
        label="Только для чтения"
        :disabled="true"
      />
    </section>

    <!-- Different Rows -->
    <section class="example-section">
      <h2>Разное количество строк</h2>
      <div class="textarea-grid">
        <div>
          <label class="label">3 строки (по умолчанию)</label>
          <BaseTextarea
            v-model="rows3"
            placeholder="3 строки..."
          />
        </div>
        <div>
          <label class="label">5 строк</label>
          <BaseTextarea
            v-model="rows5"
            placeholder="5 строк..."
            :rows="5"
          />
        </div>
        <div>
          <label class="label">10 строк</label>
          <BaseTextarea
            v-model="rows10"
            placeholder="10 строк..."
            :rows="10"
          />
        </div>
      </div>
    </section>

    <!-- Real Use Cases -->
    <section class="example-section">
      <h2>Реальные примеры использования</h2>

      <!-- AI Prompt -->
      <div class="use-case">
        <h3>AI промпт</h3>
        <BaseTextarea
          v-model="aiPrompt"
          label="Промпт для AI"
          placeholder="Опишите задачу для AI ассистента..."
          :rows="5"
        />
        <div class="flex justify-between items-center mt-2">
          <span class="text-xs text-gray-400">
            {{ aiPrompt.length }} / 2000 символов
          </span>
          <BaseButton
            variant="primary"
            size="sm"
            :disabled="!aiPrompt.trim()"
            @click="sendPrompt"
          >
            Отправить
          </BaseButton>
        </div>
      </div>

      <!-- Code Snippet -->
      <div class="use-case">
        <h3>Фрагмент кода</h3>
        <BaseTextarea
          v-model="codeSnippet"
          label="Код"
          placeholder="Вставьте код..."
          :rows="8"
          class="font-mono"
        />
        <div class="flex gap-2 mt-2">
          <BaseButton variant="secondary" size="sm" @click="formatCode">
            Форматировать
          </BaseButton>
          <BaseButton variant="ghost" size="sm" @click="copyCode">
            Копировать
          </BaseButton>
        </div>
      </div>

      <!-- Commit Message -->
      <div class="use-case">
        <h3>Сообщение коммита</h3>
        <BaseTextarea
          v-model="commitMessage"
          label="Commit Message"
          placeholder="feat: добавить новую функцию&#10;&#10;Подробное описание изменений..."
          :rows="4"
          :error="commitError"
          @blur="validateCommit"
        />
        <p class="text-xs text-gray-400 mt-1">
          Первая строка - краткое описание (до 72 символов)
        </p>
      </div>

      <!-- Notes -->
      <div class="use-case">
        <h3>Заметки</h3>
        <BaseTextarea
          v-model="notes"
          label="Заметки к проекту"
          placeholder="Добавьте заметки..."
          :rows="6"
        />
        <div class="flex justify-between items-center mt-2">
          <span class="text-xs text-gray-400">
            Последнее изменение: {{ lastModified }}
          </span>
          <BaseButton
            variant="success"
            size="sm"
            :loading="isSaving"
            @click="saveNotes"
          >
            Сохранить
          </BaseButton>
        </div>
      </div>

      <!-- JSON Editor -->
      <div class="use-case">
        <h3>JSON конфигурация</h3>
        <BaseTextarea
          v-model="jsonConfig"
          label="Конфигурация"
          placeholder='{ "key": "value" }'
          :rows="8"
          :error="jsonError"
          class="font-mono"
          @input="validateJson"
        />
        <div v-if="!jsonError && jsonConfig" class="mt-2 p-2 bg-green-900/20 rounded text-xs text-green-400">
          ✓ Валидный JSON
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseTextarea from '../BaseTextarea.vue'
import BaseButton from '../BaseButton.vue'
import { useLogger } from '@/composables/useLogger'

const logger = useLogger('BaseTextareaExample')

// Basic
const basicText = ref('')
const labelText = ref('')

// With Error
const errorText = ref('')
const commentError = ref('')

// Disabled
const disabledText = ref('Это текст нельзя редактировать.\nОн доступен только для чтения.')

// Different Rows
const rows3 = ref('')
const rows5 = ref('')
const rows10 = ref('')

// Real Use Cases
const aiPrompt = ref('')
const codeSnippet = ref(`function example() {\n  console.log('Hello World');\n  return true;\n}`)
const commitMessage = ref('')
const commitError = ref('')
const notes = ref('')
const isSaving = ref(false)
const lastModified = computed(() => new Date().toLocaleString('ru-RU'))
const jsonConfig = ref('{\n  "theme": "dark",\n  "fontSize": 14\n}')
const jsonError = ref('')

function validateComment() {
  if (errorText.value.length < 10 && errorText.value.length > 0) {
    commentError.value = 'Комментарий должен содержать минимум 10 символов'
  } else {
    commentError.value = ''
  }
}

function validateCommit() {
  const firstLine = commitMessage.value.split('\n')[0]
  if (firstLine.length > 72) {
    commitError.value = 'Первая строка не должна превышать 72 символа'
  } else {
    commitError.value = ''
  }
}

function validateJson() {
  if (!jsonConfig.value.trim()) {
    jsonError.value = ''
    return
  }
  
  try {
    JSON.parse(jsonConfig.value)
    jsonError.value = ''
  } catch (e) {
    jsonError.value = 'Невалидный JSON формат'
  }
}

function sendPrompt() {
  logger.debug('Sending prompt:', aiPrompt.value)
}

function formatCode() {
  // Simulate code formatting
  logger.debug('Formatting code...')
}

function copyCode() {
  navigator.clipboard.writeText(codeSnippet.value)
}

function saveNotes() {
  isSaving.value = true
  setTimeout(() => {
    isSaving.value = false
  }, 1000)
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

.textarea-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}

.label {
  display: block;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}

.char-count {
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: var(--text-muted);
  text-align: right;
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

.font-mono {
  font-family: var(--font-mono);
}
</style>
