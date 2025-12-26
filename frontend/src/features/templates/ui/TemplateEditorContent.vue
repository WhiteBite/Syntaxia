<template>
  <div class="tpl-editor-content">
    <!-- Section Chips -->
    <div class="tpl-chips">
      <button 
        v-for="s in sectionsList" 
        :key="s.key" 
        @click="$emit('toggle-section', s.key)" 
        class="tpl-chip"
        :class="{ active: template.sections[s.key] }" 
        :title="t(`templates.sectionHint.${s.key}`)"
      >
        <Check v-if="template.sections[s.key]" class="w-3 h-3" />
        <span>{{ t(`templates.section.${s.key}`) }}</span>
      </button>
    </div>

    <!-- Cards Stack -->
    <div class="tpl-cards">
      <!-- Role Card -->
      <TemplateCard 
        v-if="template.sections.role" 
        :title="t('templates.roleContent')" 
        :icon="UserIcon"
        :count="template.roleContent?.length || 0" 
        :enabled="true"
      >
        <textarea 
          ref="roleTextarea" 
          :value="template.roleContent"
          @input="handleTextareaInput('roleContent', $event)" 
          :placeholder="t('templates.rolePlaceholder')" 
          class="tpl-textarea"
        />
      </TemplateCard>

      <!-- Rules Card -->
      <TemplateCard 
        v-if="template.sections.rules" 
        :title="t('templates.rulesContent')" 
        :icon="ListChecksIcon"
        :count="template.rulesContent?.length || 0" 
        :enabled="true"
      >
        <textarea 
          ref="rulesTextarea" 
          :value="template.rulesContent"
          @input="handleTextareaInput('rulesContent', $event)"
          :placeholder="t('templates.rulesPlaceholder')" 
          class="tpl-textarea"
        />
      </TemplateCard>

      <!-- Context Options - Grid Tiles -->
      <div class="tpl-options-card">
        <div class="tpl-options-header">
          <Settings2 class="w-3.5 h-3.5" />
          <span>{{ t('templates.contextOptions') }}</span>
        </div>
        <div class="tpl-options-grid">
          <TemplateOptionTile 
            :model-value="template.sections.tree"
            @update:model-value="$emit('toggle-section', 'tree')"
            :label="t('templates.section.tree')" 
            :icon="FolderTreeIcon" 
            hint="Structure" 
          />
          <TemplateOptionTile 
            :model-value="template.sections.stats"
            @update:model-value="$emit('toggle-section', 'stats')"
            :label="t('templates.section.stats')"
            :icon="HashIcon" 
            hint="Metrics" 
          />
          <TemplateOptionTile 
            :model-value="template.sections.files"
            @update:model-value="$emit('toggle-section', 'files')"
            :label="t('templates.section.files')"
            :icon="FileCodeIcon" 
            hint="Content" 
          />
          <TemplateOptionTile 
            :model-value="template.sections.task"
            @update:model-value="$emit('toggle-section', 'task')"
            :label="t('templates.section.task')"
            :icon="ClipboardIcon" 
            hint="Your task" 
          />
        </div>
      </div>

      <!-- Advanced (Prefix/Suffix) -->
      <details class="tpl-advanced">
        <summary>
          <ChevronRight class="w-3.5 h-3.5 tpl-chevron" />{{ t('templates.additional') }}
        </summary>
        <div class="tpl-advanced-content">
          <div class="tpl-advanced-field">
            <label>{{ t('templates.prefix') }}</label>
            <textarea 
              :value="template.customPrefix"
              @input="$emit('update-field', 'customPrefix', ($event.target as HTMLTextAreaElement).value)"
              :placeholder="t('templates.prefixPlaceholder')" 
              rows="2" 
            />
          </div>
          <div class="tpl-advanced-field">
            <label>{{ t('templates.suffix') }}</label>
            <textarea 
              :value="template.customSuffix"
              @input="$emit('update-field', 'customSuffix', ($event.target as HTMLTextAreaElement).value)"
              :placeholder="t('templates.suffixPlaceholder')" 
              rows="2" 
            />
          </div>
        </div>
      </details>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { 
  Check, ChevronRight, Clipboard, FileCode, 
  FolderTree, Hash, ListChecks, Settings2, User 
} from 'lucide-vue-next'
import { nextTick, ref, shallowRef, watch } from 'vue'
import { SECTION_META, type PromptTemplate, type TemplateSections } from '../model/template.types'
import TemplateCard from './TemplateCard.vue'
import TemplateOptionTile from './TemplateOptionTile.vue'

const { t } = useI18n()

// Icon refs for dynamic components
const UserIcon = shallowRef(User)
const ListChecksIcon = shallowRef(ListChecks)
const FolderTreeIcon = shallowRef(FolderTree)
const HashIcon = shallowRef(Hash)
const FileCodeIcon = shallowRef(FileCode)
const ClipboardIcon = shallowRef(Clipboard)

const sectionsList = SECTION_META

const props = defineProps<{
  template: PromptTemplate
}>()

const emit = defineEmits<{
  (e: 'toggle-section', key: keyof TemplateSections): void
  (e: 'update-field', field: keyof PromptTemplate, value: string): void
}>()

const roleTextarea = ref<HTMLTextAreaElement | null>(null)
const rulesTextarea = ref<HTMLTextAreaElement | null>(null)

function autoResize(el: HTMLTextAreaElement) {
  el.style.height = 'auto'
  el.style.height = Math.max(80, el.scrollHeight) + 'px'
}

function handleTextareaInput(field: 'roleContent' | 'rulesContent', event: Event) {
  const target = event.target as HTMLTextAreaElement
  autoResize(target)
  emit('update-field', field, target.value)
}

// Auto-resize textareas when template changes
watch(() => props.template, () => {
  nextTick(() => {
    if (roleTextarea.value) autoResize(roleTextarea.value)
    if (rulesTextarea.value) autoResize(rulesTextarea.value)
  })
})

// Expose refs for parent to trigger resize
defineExpose({
  resizeTextareas: () => {
    nextTick(() => {
      if (roleTextarea.value) autoResize(roleTextarea.value)
      if (rulesTextarea.value) autoResize(rulesTextarea.value)
    })
  }
})
</script>

<style scoped>
.tpl-editor-content { display: flex; flex-direction: column; flex: 1; min-height: 0; overflow: hidden; }

.tpl-chips { display: flex; align-items: center; gap: 0.375rem; padding: 0.625rem 1rem; background: rgba(255, 255, 255, 0.02); border-bottom: 1px solid rgba(255, 255, 255, 0.06); flex-wrap: wrap; }
.tpl-chip { display: flex; align-items: center; gap: 0.25rem; padding: 0.375rem 0.75rem; background: transparent; border: 1px solid rgba(255, 255, 255, 0.12); border-radius: var(--radius-full); color: var(--text-muted); font-size: 11px; font-weight: 500; cursor: pointer; transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1); }
.tpl-chip:hover { background: rgba(255, 255, 255, 0.04); border-color: rgba(255, 255, 255, 0.2); transform: translateY(-1px); }
.tpl-chip.active { background: var(--accent-indigo-bg); border-color: var(--accent-indigo-border); color: var(--accent-indigo); }

.tpl-cards { flex: 1; overflow-y: auto; padding: 0.75rem 1rem; display: flex; flex-direction: column; gap: 0.625rem; min-height: 0; }
.tpl-cards::-webkit-scrollbar { width: 4px; }
.tpl-cards::-webkit-scrollbar-track { background: transparent; }
.tpl-cards::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 2px; }

.tpl-textarea { width: 100%; min-height: 80px; max-height: 200px; padding: 0; background: transparent; border: none; color: var(--text-secondary); font-size: 13px; font-family: var(--font-mono); line-height: 1.6; resize: none; outline: none; }
.tpl-textarea::placeholder { color: var(--text-subtle); font-style: italic; }

.tpl-options-card { background: #0d1117; border: 1px solid rgba(255, 255, 255, 0.08); border-radius: var(--radius-md); overflow: hidden; }
.tpl-options-header { display: flex; align-items: center; gap: 0.5rem; padding: 0.625rem 0.75rem; background: rgba(255, 255, 255, 0.02); border-bottom: 1px solid rgba(255, 255, 255, 0.06); font-size: 11px; font-weight: 600; color: var(--text-muted); text-transform: uppercase; }
.tpl-options-header svg { color: var(--text-subtle); }
.tpl-options-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.5rem; padding: 0.75rem; }

.tpl-advanced { background: rgba(255, 255, 255, 0.02); border: 1px solid rgba(255, 255, 255, 0.06); border-radius: var(--radius-md); }
.tpl-advanced summary { display: flex; align-items: center; gap: 0.375rem; padding: 0.625rem 0.75rem; font-size: 11px; font-weight: 500; color: var(--text-muted); cursor: pointer; user-select: none; list-style: none; }
.tpl-advanced summary::-webkit-details-marker { display: none; }
.tpl-advanced summary:hover { color: var(--text-primary); }
.tpl-advanced .tpl-chevron { transition: transform 0.2s; }
.tpl-advanced[open] .tpl-chevron { transform: rotate(90deg); }
.tpl-advanced-content { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; padding: 0 0.75rem 0.75rem; }
.tpl-advanced-field label { display: block; margin-bottom: 0.375rem; font-size: 10px; font-weight: 500; color: var(--text-subtle); text-transform: uppercase; }
.tpl-advanced-field textarea { width: 100%; padding: 0.5rem; background: rgba(0, 0, 0, 0.3); border: 1px solid rgba(255, 255, 255, 0.08); border-radius: var(--radius-sm); color: var(--text-secondary); font-size: 11px; font-family: var(--font-mono); resize: none; outline: none; transition: border-color 0.15s; }
.tpl-advanced-field textarea:focus { border-color: var(--accent-indigo-border); }
</style>
