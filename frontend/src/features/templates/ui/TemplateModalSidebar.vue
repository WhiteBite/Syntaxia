<template>
  <aside class="tpl-sidebar">
    <div class="tpl-search">
      <Search class="w-3.5 h-3.5" />
      <input 
        :value="searchQuery" 
        @input="$emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
        :placeholder="t('templates.search')" 
      />
    </div>
    
    <nav class="tpl-list">
      <!-- Favorites -->
      <details v-if="filteredFavorites.length" class="tpl-group" open>
        <summary>
          <Star class="w-3 h-3 text-amber-400" />{{ t('templates.favorites') }}
        </summary>
        <TemplateListItem 
          v-for="tpl in filteredFavorites" 
          :key="tpl.id" 
          :template="tpl" 
          :active="selectedTemplateId === tpl.id" 
          :is-current="tpl.id === activeTemplateId"
          @select="$emit('select', tpl)" 
          @delete="$emit('delete', tpl.id)" 
          @toggle-favorite="$emit('toggle-favorite', tpl.id)" 
        />
      </details>
      
      <!-- Built-in -->
      <details v-if="filteredBuiltIn.length" class="tpl-group" open>
        <summary>
          <Zap class="w-3 h-3 text-blue-400" />{{ t('templates.builtIn') }}
        </summary>
        <TemplateListItem 
          v-for="tpl in filteredBuiltIn" 
          :key="tpl.id" 
          :template="tpl"
          :active="selectedTemplateId === tpl.id" 
          :is-current="tpl.id === activeTemplateId"
          @select="$emit('select', tpl)" 
          @toggle-favorite="$emit('toggle-favorite', tpl.id)" 
        />
      </details>

      <!-- Custom -->
      <details v-if="filteredCustom.length" class="tpl-group" open>
        <summary>
          <User class="w-3 h-3 text-emerald-400" />{{ t('templates.custom') }}
        </summary>
        <TemplateListItem 
          v-for="tpl in filteredCustom" 
          :key="tpl.id" 
          :template="tpl"
          :active="selectedTemplateId === tpl.id" 
          :is-current="tpl.id === activeTemplateId"
          @select="$emit('select', tpl)" 
          @delete="$emit('delete', tpl.id)"
          @toggle-favorite="$emit('toggle-favorite', tpl.id)" 
        />
      </details>
      
      <div v-if="noResults" class="tpl-no-results">{{ t('templates.noResults') }}</div>
    </nav>
    
    <button @click="$emit('create-new')" class="tpl-new-btn">
      <Plus class="w-3.5 h-3.5" />{{ t('templates.create') }}
    </button>
  </aside>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { Plus, Search, Star, User, Zap } from 'lucide-vue-next'
import { computed } from 'vue'
import type { PromptTemplate } from '../model/template.types'
import TemplateListItem from './TemplateListItem.vue'

const { t } = useI18n()

const props = defineProps<{
  searchQuery: string
  selectedTemplateId: string | null
  activeTemplateId: string | null
  favoriteTemplates: PromptTemplate[]
  builtInTemplates: PromptTemplate[]
  customTemplates: PromptTemplate[]
}>()

defineEmits<{
  (e: 'update:searchQuery', value: string): void
  (e: 'select', template: PromptTemplate): void
  (e: 'delete', id: string): void
  (e: 'toggle-favorite', id: string): void
  (e: 'create-new'): void
}>()

const filteredFavorites = computed(() => {
  let list = props.favoriteTemplates
  if (props.searchQuery) {
    const q = props.searchQuery.toLowerCase()
    list = list.filter(t => t.name.toLowerCase().includes(q))
  }
  return list
})

const filteredBuiltIn = computed(() => {
  let list = props.builtInTemplates.filter(t => !t.isFavorite && !t.isHidden)
  if (props.searchQuery) {
    const q = props.searchQuery.toLowerCase()
    list = list.filter(t => t.name.toLowerCase().includes(q))
  }
  return list
})

const filteredCustom = computed(() => {
  let list = props.customTemplates.filter(t => !t.isHidden)
  if (props.searchQuery) {
    const q = props.searchQuery.toLowerCase()
    list = list.filter(t => t.name.toLowerCase().includes(q))
  }
  return list
})

const noResults = computed(() => 
  props.searchQuery && 
  !filteredFavorites.value.length && 
  !filteredBuiltIn.value.length && 
  !filteredCustom.value.length
)
</script>

<style scoped>
.tpl-sidebar {
  display: flex;
  flex-direction: column;
  background: #08090C;
  border-right: 1px solid rgba(255, 255, 255, 0.06);
}

.tpl-search {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0.75rem;
  padding: 0.5rem 0.625rem;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-md);
  transition: border-color 0.15s;
}

.tpl-search:focus-within {
  border-color: var(--accent-indigo-border);
}

.tpl-search svg {
  color: var(--text-subtle);
  flex-shrink: 0;
}

.tpl-search input { 
  flex: 1;
  background: transparent;
  border: none; 
  color: var(--text-primary);
  font-size: 12px;
  outline: none;
  min-width: 0;
}

.tpl-search input::placeholder {
  color: var(--text-subtle);
}

.tpl-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 0.5rem;
  min-height: 0;
}

.tpl-list::-webkit-scrollbar { width: 4px; }
.tpl-list::-webkit-scrollbar-track { background: transparent; }
.tpl-list::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 2px; }
.tpl-list::-webkit-scrollbar-thumb:hover { background: rgba(255, 255, 255, 0.2); }

.tpl-group {
  margin-bottom: 0.5rem;
}

.tpl-group summary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  font-size: 10px;
  font-weight: 600;
  color: var(--text-subtle);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  cursor: pointer;
  user-select: none;
  list-style: none;
  border-radius: var(--radius-sm);
  transition: color 0.15s;
}

.tpl-group summary::-webkit-details-marker {
  display: none;
}

.tpl-group summary:hover {
  color: var(--text-muted);
}

.tpl-no-results {
  padding: 1.5rem;
  text-align: center;
  font-size: 11px;
  color: var(--text-subtle);
}

.tpl-new-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
  margin: 0.75rem;
  padding: 0.625rem;
  background: transparent;
  border: 1px dashed rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-md);
  color: var(--text-muted);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}

.tpl-new-btn:hover { 
  background: rgba(255, 255, 255, 0.04); 
  border-color: var(--accent-indigo-border); 
  color: var(--accent-indigo); 
}
</style>
