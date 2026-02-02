<template>
  <div class="examples-container">
    <h1>BaseChip Examples</h1>
    
    <!-- Variants -->
    <section class="example-section">
      <h2>Variants</h2>
      <div class="example-grid">
        <div class="example-item">
          <label>Default</label>
          <BaseChip variant="default">Default Chip</BaseChip>
        </div>
        <div class="example-item">
          <label>Primary</label>
          <BaseChip variant="primary">Primary Chip</BaseChip>
        </div>
        <div class="example-item">
          <label>Success</label>
          <BaseChip variant="success">Success Chip</BaseChip>
        </div>
        <div class="example-item">
          <label>Warning</label>
          <BaseChip variant="warning">Warning Chip</BaseChip>
        </div>
        <div class="example-item">
          <label>Danger</label>
          <BaseChip variant="danger">Danger Chip</BaseChip>
        </div>
      </div>
    </section>

    <!-- Sizes -->
    <section class="example-section">
      <h2>Sizes</h2>
      <div class="example-grid">
        <div class="example-item">
          <label>Extra Small (xs)</label>
          <BaseChip variant="primary" size="xs">Extra Small</BaseChip>
        </div>
        <div class="example-item">
          <label>Small (sm)</label>
          <BaseChip variant="primary" size="sm">Small</BaseChip>
        </div>
        <div class="example-item">
          <label>Medium (md)</label>
          <BaseChip variant="primary" size="md">Medium</BaseChip>
        </div>
      </div>
    </section>

    <!-- With Icons -->
    <section class="example-section">
      <h2>With Icons</h2>
      <div class="example-grid">
        <div class="example-item">
          <label>Icon via prop</label>
          <BaseChip variant="primary" :icon="FileCode">JavaScript</BaseChip>
        </div>
        <div class="example-item">
          <label>Icon via slot (emoji)</label>
          <BaseChip variant="success">
            <template #icon>
              <span>✓</span>
            </template>
            Completed
          </BaseChip>
        </div>
        <div class="example-item">
          <label>Icon via slot (SVG)</label>
          <BaseChip variant="warning">
            <template #icon>
              <AlertTriangle class="w-3 h-3" />
            </template>
            Warning
          </BaseChip>
        </div>
      </div>
    </section>

    <!-- Removable -->
    <section class="example-section">
      <h2>Removable Chips</h2>
      <div class="example-grid">
        <div class="example-item">
          <label>Default removable</label>
          <BaseChip variant="default" removable @remove="handleRemove('Default')">
            Removable
          </BaseChip>
        </div>
        <div class="example-item">
          <label>Primary removable</label>
          <BaseChip variant="primary" removable @remove="handleRemove('Primary')">
            TypeScript
          </BaseChip>
        </div>
        <div class="example-item">
          <label>Success removable</label>
          <BaseChip variant="success" removable @remove="handleRemove('Success')">
            Selected
          </BaseChip>
        </div>
        <div class="example-item">
          <label>With icon + removable</label>
          <BaseChip variant="primary" :icon="Filter" removable @remove="handleRemove('Filter')">
            Active Filter
          </BaseChip>
        </div>
      </div>
    </section>

    <!-- Clickable -->
    <section class="example-section">
      <h2>Clickable Chips</h2>
      <div class="example-grid">
        <div class="example-item">
          <label>Clickable default</label>
          <BaseChip variant="default" clickable @click="handleClick('Default')">
            Click me
          </BaseChip>
        </div>
        <div class="example-item">
          <label>Clickable primary</label>
          <BaseChip variant="primary" clickable @click="handleClick('Primary')">
            Click me
          </BaseChip>
        </div>
        <div class="example-item">
          <label>Clickable with icon</label>
          <BaseChip variant="success" :icon="Check" clickable @click="handleClick('Success')">
            Toggle
          </BaseChip>
        </div>
      </div>
    </section>

    <!-- Combined Features -->
    <section class="example-section">
      <h2>Combined Features</h2>
      <div class="example-grid">
        <div class="example-item">
          <label>Clickable + Removable</label>
          <BaseChip 
            variant="primary" 
            clickable 
            removable 
            @click="handleClick('Combined')"
            @remove="handleRemove('Combined')"
          >
            Multi-action
          </BaseChip>
        </div>
        <div class="example-item">
          <label>Icon + Clickable + Removable</label>
          <BaseChip 
            variant="warning" 
            :icon="Tag"
            clickable 
            removable 
            @click="handleClick('Full')"
            @remove="handleRemove('Full')"
          >
            Full Featured
          </BaseChip>
        </div>
      </div>
    </section>

    <!-- Use Cases -->
    <section class="example-section">
      <h2>Real-world Use Cases</h2>
      
      <div class="use-case">
        <h3>File Type Filters</h3>
        <div class="flex flex-wrap gap-2">
          <BaseChip 
            v-for="filter in fileFilters" 
            :key="filter.id"
            :variant="filter.active ? 'primary' : 'default'"
            clickable
            @click="toggleFilter(filter.id)"
          >
            <template #icon>
              <span>{{ filter.icon }}</span>
            </template>
            {{ filter.label }}
          </BaseChip>
        </div>
      </div>

      <div class="use-case">
        <h3>Selected Tags</h3>
        <div class="flex flex-wrap gap-2">
          <BaseChip 
            v-for="tag in selectedTags" 
            :key="tag"
            variant="success"
            removable
            @remove="removeTag(tag)"
          >
            {{ tag }}
          </BaseChip>
        </div>
      </div>

      <div class="use-case">
        <h3>Status Indicators</h3>
        <div class="flex flex-wrap gap-2">
          <BaseChip variant="success" size="xs">
            <template #icon>
              <Check class="w-3 h-3" />
            </template>
            Completed
          </BaseChip>
          <BaseChip variant="warning" size="xs">
            <template #icon>
              <Clock class="w-3 h-3" />
            </template>
            In Progress
          </BaseChip>
          <BaseChip variant="danger" size="xs">
            <template #icon>
              <X class="w-3 h-3" />
            </template>
            Failed
          </BaseChip>
        </div>
      </div>

      <div class="use-case">
        <h3>Quick Suggestions</h3>
        <div class="flex flex-wrap gap-2">
          <BaseChip 
            v-for="suggestion in suggestions" 
            :key="suggestion"
            variant="default"
            clickable
            size="sm"
            @click="applySuggestion(suggestion)"
          >
            {{ suggestion }}
          </BaseChip>
        </div>
      </div>
    </section>

    <!-- Event Log -->
    <section class="example-section">
      <h2>Event Log</h2>
      <div class="event-log">
        <div v-for="(event, index) in events" :key="index" class="event-item">
          {{ event }}
        </div>
        <div v-if="events.length === 0" class="text-gray-500">
          No events yet. Try clicking or removing chips above.
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import BaseChip from '../BaseChip.vue'
import { FileCode, AlertTriangle, Filter, Check, Tag, Clock, X } from 'lucide-vue-next'

const events = ref<string[]>([])

const fileFilters = ref([
  { id: 'code', label: 'Code', icon: '💻', active: false },
  { id: 'test', label: 'Tests', icon: '🧪', active: true },
  { id: 'docs', label: 'Docs', icon: '📝', active: false },
  { id: 'config', label: 'Config', icon: '⚙️', active: false }
])

const selectedTags = ref(['TypeScript', 'Vue', 'Pinia', 'Tailwind'])

const suggestions = ref([
  'Add unit tests',
  'Refactor component',
  'Update documentation',
  'Fix linting errors'
])

function handleClick(label: string) {
  const event = `Clicked: ${label} at ${new Date().toLocaleTimeString()}`
  events.value.unshift(event)
  if (events.value.length > 10) events.value.pop()
}

function handleRemove(label: string) {
  const event = `Removed: ${label} at ${new Date().toLocaleTimeString()}`
  events.value.unshift(event)
  if (events.value.length > 10) events.value.pop()
}

function toggleFilter(id: string) {
  const filter = fileFilters.value.find(f => f.id === id)
  if (filter) {
    filter.active = !filter.active
    handleClick(`Filter ${filter.label}`)
  }
}

function removeTag(tag: string) {
  selectedTags.value = selectedTags.value.filter(t => t !== tag)
  handleRemove(`Tag ${tag}`)
}

function applySuggestion(suggestion: string) {
  handleClick(`Suggestion: ${suggestion}`)
}
</script>

<style scoped>
.examples-container {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  background: var(--bg-0);
  color: var(--text-primary);
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
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-default);
}

h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: var(--text-primary);
}

h3 {
  font-size: 1.125rem;
  font-weight: 500;
  margin-bottom: 0.75rem;
  color: var(--text-secondary);
}

.example-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.example-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.example-item label {
  font-size: 0.875rem;
  color: var(--text-muted);
  font-weight: 500;
}

.use-case {
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: var(--bg-2);
  border-radius: var(--radius-md);
}

.use-case:last-child {
  margin-bottom: 0;
}

.event-log {
  max-height: 300px;
  overflow-y: auto;
  padding: 1rem;
  background: var(--bg-2);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 0.875rem;
}

.event-item {
  padding: 0.5rem;
  margin-bottom: 0.25rem;
  background: var(--bg-3);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
}

.event-item:last-child {
  margin-bottom: 0;
}
</style>
