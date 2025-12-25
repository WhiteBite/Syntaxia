<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="tpl-overlay" @click.self="handleClose">
        <div class="tpl-modal" @keydown="handleKeydown">
          <!-- Header -->
          <div class="tpl-header">
            <div class="tpl-header-left">
              <FileText class="w-4 h-4 text-purple-400" />
              <h2>{{ t('templates.manage') }}</h2>
            </div>
            <div class="tpl-header-center">
              <span v-if="hasChanges && !isEditingBuiltIn" class="tpl-unsaved-indicator" />
            </div>
            <button @click="handleClose" class="tpl-close"><X class="w-4 h-4" /></button>
          </div>

          <!-- Body: 3 columns -->
          <div class="tpl-body">
            <!-- Sidebar (Library) -->
            <aside class="tpl-sidebar">
              <div class="tpl-search">
                <Search class="w-3.5 h-3.5" />
                <input v-model="searchQuery" :placeholder="t('templates.search')" />
              </div>
              
              <nav class="tpl-list">
                <!-- Favorites -->
                <details v-if="filteredFavorites.length" class="tpl-group" open>
                  <summary><Star class="w-3 h-3 text-amber-400" />{{ t('templates.favorites') }}</summary>
                  <TemplateListItem v-for="tpl in filteredFavorites" :key="tpl.id" :template="tpl" 
                    :active="selectedTemplate?.id === tpl.id" :is-current="tpl.id === templateStore.activeTemplateId"
                    @select="selectTemplate(tpl)" @delete="confirmDelete(tpl.id)" 
                    @toggle-favorite="templateStore.toggleFavorite(tpl.id)" />
                </details>
                
                <!-- Built-in -->
                <details v-if="filteredBuiltIn.length" class="tpl-group" open>
                  <summary><Zap class="w-3 h-3 text-blue-400" />{{ t('templates.builtIn') }}</summary>
                  <TemplateListItem v-for="tpl in filteredBuiltIn" :key="tpl.id" :template="tpl"
                    :active="selectedTemplate?.id === tpl.id" :is-current="tpl.id === templateStore.activeTemplateId"
                    @select="selectTemplate(tpl)" @toggle-favorite="templateStore.toggleFavorite(tpl.id)" />
                </details>

                <!-- Custom -->
                <details v-if="filteredCustom.length" class="tpl-group" open>
                  <summary><User class="w-3 h-3 text-emerald-400" />{{ t('templates.custom') }}</summary>
                  <TemplateListItem v-for="tpl in filteredCustom" :key="tpl.id" :template="tpl"
                    :active="selectedTemplate?.id === tpl.id" :is-current="tpl.id === templateStore.activeTemplateId"
                    @select="selectTemplate(tpl)" @delete="confirmDelete(tpl.id)"
                    @toggle-favorite="templateStore.toggleFavorite(tpl.id)" />
                </details>
                
                <div v-if="noResults" class="tpl-no-results">{{ t('templates.noResults') }}</div>
              </nav>
              
              <button @click="createNew" class="tpl-new-btn">
                <Plus class="w-3.5 h-3.5" />{{ t('templates.create') }}
              </button>
            </aside>

            <!-- Editor (Builder) -->
            <main class="tpl-editor">
              <template v-if="editingTemplate">
                <!-- Combo Header: Icon + Name -->
                <div class="tpl-editor-header">
                  <button class="tpl-icon-btn" @click="showEmojiPicker = !showEmojiPicker" :disabled="editingTemplate.isBuiltIn">
                    <span class="tpl-icon-display">{{ editingTemplate.icon }}</span>
                  </button>
                  <div class="tpl-name-group">
                    <input v-model="editingTemplate.name" class="tpl-name-input" 
                      :placeholder="t('templates.namePlaceholder')" :disabled="editingTemplate.isBuiltIn" />
                    <input v-model="editingTemplate.description" class="tpl-desc-input" 
                      :placeholder="t('templates.descriptionPlaceholder')" :disabled="editingTemplate.isBuiltIn" />
                  </div>
                  <span v-if="editingTemplate.id === templateStore.activeTemplateId" class="tpl-active-badge">
                    <Check class="w-3 h-3" />{{ t('templates.currentTemplate') }}
                  </span>
                </div>

                <!-- Section Chips -->
                <div class="tpl-chips">
                  <button v-for="s in sectionsList" :key="s.key" @click="toggleSection(s.key)" class="tpl-chip"
                    :class="{ active: editingTemplate.sections[s.key] }" :title="t(`templates.sectionHint.${s.key}`)">
                    <Check v-if="editingTemplate.sections[s.key]" class="w-3 h-3" />
                    <span>{{ t(`templates.section.${s.key}`) }}</span>
                  </button>
                </div>

                <!-- Cards Stack -->
                <div class="tpl-cards">
                  <!-- Role Card -->
                  <TemplateCard v-if="editingTemplate.sections.role" :title="t('templates.roleContent')" :icon="UserIcon"
                    :count="editingTemplate.roleContent?.length || 0" :enabled="true">
                    <textarea ref="roleTextarea" v-model="editingTemplate.roleContent" 
                      :placeholder="t('templates.rolePlaceholder')" class="tpl-textarea"
                      @input="autoResize($event.target as HTMLTextAreaElement)" />
                  </TemplateCard>

                  <!-- Rules Card -->
                  <TemplateCard v-if="editingTemplate.sections.rules" :title="t('templates.rulesContent')" :icon="ListChecksIcon"
                    :count="editingTemplate.rulesContent?.length || 0" :enabled="true">
                    <textarea ref="rulesTextarea" v-model="editingTemplate.rulesContent"
                      :placeholder="t('templates.rulesPlaceholder')" class="tpl-textarea"
                      @input="autoResize($event.target as HTMLTextAreaElement)" />
                  </TemplateCard>

                  <!-- Context Options - Grid Tiles -->
                  <div class="tpl-options-card">
                    <div class="tpl-options-header">
                      <Settings2 class="w-3.5 h-3.5" />
                      <span>{{ t('templates.contextOptions') }}</span>
                    </div>
                    <div class="tpl-options-grid">
                      <TemplateOptionTile v-model="editingTemplate.sections.tree" :label="t('templates.section.tree')" 
                        :icon="FolderTreeIcon" hint="Structure" />
                      <TemplateOptionTile v-model="editingTemplate.sections.stats" :label="t('templates.section.stats')"
                        :icon="HashIcon" hint="Metrics" />
                      <TemplateOptionTile v-model="editingTemplate.sections.files" :label="t('templates.section.files')"
                        :icon="FileCodeIcon" hint="Content" />
                      <TemplateOptionTile v-model="editingTemplate.sections.task" :label="t('templates.section.task')"
                        :icon="ClipboardIcon" hint="Your task" />
                    </div>
                  </div>

                  <!-- Advanced (Prefix/Suffix) -->
                  <details class="tpl-advanced">
                    <summary><ChevronRight class="w-3.5 h-3.5 tpl-chevron" />{{ t('templates.additional') }}</summary>
                    <div class="tpl-advanced-content">
                      <div class="tpl-advanced-field">
                        <label>{{ t('templates.prefix') }}</label>
                        <textarea v-model="editingTemplate.customPrefix" :placeholder="t('templates.prefixPlaceholder')" rows="2" />
                      </div>
                      <div class="tpl-advanced-field">
                        <label>{{ t('templates.suffix') }}</label>
                        <textarea v-model="editingTemplate.customSuffix" :placeholder="t('templates.suffixPlaceholder')" rows="2" />
                      </div>
                    </div>
                  </details>
                </div>
              </template>
              
              <div v-else class="tpl-empty">
                <FileText class="w-10 h-10 opacity-15" />
                <p>{{ t('templates.selectToEdit') }}</p>
              </div>
            </main>

            <!-- Preview Panel -->
            <aside class="tpl-preview">
              <div class="tpl-preview-header">
                <Eye class="w-3.5 h-3.5" />
                <span>{{ t('templates.preview') }}</span>
              </div>
              <div class="tpl-preview-content" v-html="highlightedPreview"></div>
              <div class="tpl-preview-footer">
                <div class="tpl-token-bar">
                  <div class="tpl-token-fill" :style="{ width: tokenPercent + '%' }" :class="tokenBarClass" />
                </div>
                <div class="tpl-token-info">
                  <span class="tpl-token-count">{{ animatedTokens.toLocaleString() }}</span>
                  <span class="tpl-token-label">/ 32k tokens</span>
                </div>
              </div>
            </aside>
          </div>

          <!-- Footer -->
          <div class="tpl-footer">
            <div class="tpl-footer-left">
              <button @click="handleImport" class="tpl-footer-btn" :title="t('templates.import')">
                <Upload class="w-3.5 h-3.5" />
              </button>
              <button @click="handleExport" class="tpl-footer-btn" :title="t('templates.export')" :disabled="!selectedTemplate">
                <Download class="w-3.5 h-3.5" />
              </button>
              <button @click="duplicateSelected" class="tpl-footer-btn" :disabled="!selectedTemplate">
                <Copy class="w-3.5 h-3.5" /><span>{{ t('templates.duplicate') }}</span>
              </button>
            </div>
            <div class="tpl-footer-center">
              <Transition name="fade" mode="out-in">
                <span v-if="justSaved" class="tpl-saved-msg"><Check class="w-3 h-3" />{{ t('templates.saved') }}</span>
                <span v-else-if="hasChanges && !isEditingBuiltIn" class="tpl-unsaved-msg">{{ t('templates.unsavedChanges') }}</span>
              </Transition>
            </div>
            <div class="tpl-footer-right">
              <label class="tpl-autosave">
                <input type="checkbox" v-model="autoSaveOnClose" />
                <span class="tpl-autosave-track"><span class="tpl-autosave-thumb" /></span>
                <span>{{ t('templates.autosave') }}</span>
              </label>
              <button v-if="canApply" @click="applyTemplate" class="tpl-apply-btn">
                <Sparkles class="w-3.5 h-3.5" />{{ t('templates.apply') }}
              </button>
              <button @click="saveChanges" class="tpl-save-btn" :class="{ pulse: hasChanges && !isEditingBuiltIn }"
                :disabled="!hasChanges || isEditingBuiltIn">
                <Save class="w-3.5 h-3.5" />{{ t('common.save') }}
              </button>
            </div>
          </div>
        </div>

        <!-- Delete Dialog -->
        <Transition name="fade">
          <div v-if="deleteConfirmId" class="tpl-dialog-overlay" @click.self="deleteConfirmId = null">
            <div class="tpl-dialog danger">
              <Trash2 class="w-6 h-6" />
              <h3>{{ t('templates.deleteConfirm') }}</h3>
              <p>{{ t('templates.deleteConfirmText') }}</p>
              <div class="tpl-dialog-btns">
                <button @click="deleteConfirmId = null" class="tpl-cancel-btn">{{ t('common.cancel') }}</button>
                <button @click="executeDelete" class="tpl-danger-btn">{{ t('templates.delete') }}</button>
              </div>
            </div>
          </div>
        </Transition>
        
        <!-- Unsaved Dialog -->
        <Transition name="fade">
          <div v-if="showUnsavedWarning" class="tpl-dialog-overlay" @click.self="showUnsavedWarning = false">
            <div class="tpl-dialog warning">
              <AlertTriangle class="w-6 h-6" />
              <h3>{{ t('templates.unsavedChanges') }}</h3>
              <p>{{ t('templates.unsavedChangesText') }}</p>
              <div class="tpl-dialog-btns">
                <button @click="discardAndClose" class="tpl-cancel-btn">{{ t('templates.discard') }}</button>
                <button @click="saveAndClose" class="tpl-save-btn">{{ t('common.save') }}</button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>


<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useContextStore } from '@/features/context'
import { 
  AlertTriangle, Check, ChevronRight, Clipboard, Copy, Download, Eye, FileCode, FileText, 
  FolderTree, Hash, ListChecks, Plus, Save, Search, Settings2, Sparkles, Star, 
  Trash2, Upload, User, X, Zap 
} from 'lucide-vue-next'
import { storeToRefs } from 'pinia'
import { computed, nextTick, ref, watch, shallowRef } from 'vue'
import { useTemplateStore } from '../model/template.store'
import { createEmptyTemplate, DEFAULT_SECTION_ORDER, SECTION_META, type PromptTemplate, type TemplateSections } from '../model/template.types'
import TemplateListItem from './TemplateListItem.vue'
import TemplateCard from './TemplateCard.vue'
import TemplateOptionTile from './TemplateOptionTile.vue'

// Icon refs for dynamic components
const UserIcon = shallowRef(User)
const ListChecksIcon = shallowRef(ListChecks)
const FolderTreeIcon = shallowRef(FolderTree)
const HashIcon = shallowRef(Hash)
const FileCodeIcon = shallowRef(FileCode)
const ClipboardIcon = shallowRef(Clipboard)

const { t } = useI18n()
const templateStore = useTemplateStore()
const contextStore = useContextStore()
const { isModalOpen, builtInTemplates, customTemplates, favoriteTemplates } = storeToRefs(templateStore)

const isOpen = computed(() => isModalOpen.value)
const selectedTemplate = ref<PromptTemplate | null>(null)
const editingTemplate = ref<PromptTemplate | null>(null)
const originalJson = ref('')
const deleteConfirmId = ref<string | null>(null)
const showUnsavedWarning = ref(false)
const searchQuery = ref('')
const autoSaveOnClose = ref(localStorage.getItem('template-autosave') === 'true')
const showEmojiPicker = ref(false)
const justSaved = ref(false)

const roleTextarea = ref<HTMLTextAreaElement | null>(null)
const rulesTextarea = ref<HTMLTextAreaElement | null>(null)

const sectionsList = SECTION_META
const hasChanges = computed(() => !editingTemplate.value ? false : !selectedTemplate.value ? true : JSON.stringify(editingTemplate.value) !== originalJson.value)
const isEditingBuiltIn = computed(() => editingTemplate.value?.isBuiltIn ?? false)
const canApply = computed(() => selectedTemplate.value && selectedTemplate.value.id !== templateStore.activeTemplateId)

// Animated token counter
const animatedTokens = ref(0)
const previewTokens = computed(() => {
  if (!editingTemplate.value) return 0
  let count = 0
  const tpl = editingTemplate.value
  if (tpl.roleContent) count += Math.round(tpl.roleContent.length / 4)
  if (tpl.rulesContent) count += Math.round(tpl.rulesContent.length / 4)
  if (tpl.customPrefix) count += Math.round(tpl.customPrefix.length / 4)
  if (tpl.customSuffix) count += Math.round(tpl.customSuffix.length / 4)
  count += contextStore.tokenCount || 0
  return count
})

const tokenPercent = computed(() => Math.min(100, (previewTokens.value / 32000) * 100))
const tokenBarClass = computed(() => {
  if (tokenPercent.value > 90) return 'danger'
  if (tokenPercent.value > 70) return 'warning'
  return ''
})

// Animate token count
watch(previewTokens, (newVal) => {
  const start = animatedTokens.value
  const diff = newVal - start
  const duration = 300
  const startTime = performance.now()
  
  function animate(currentTime: number) {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)
    animatedTokens.value = Math.round(start + diff * progress)
    if (progress < 1) requestAnimationFrame(animate)
  }
  requestAnimationFrame(animate)
}, { immediate: true })

const previewContent = computed(() => {
  if (!editingTemplate.value) return ''
  const tpl = editingTemplate.value
  const parts: string[] = []
  if (tpl.customPrefix) parts.push(tpl.customPrefix)
  
  const files = contextStore.summary?.files || []
  const hasContext = contextStore.hasContext && files.length > 0
  
  for (const sec of tpl.sectionOrder) {
    if (!tpl.sections[sec]) continue
    switch (sec) {
      case 'role': if (tpl.roleContent) parts.push(`## Role\n${tpl.roleContent}`); break
      case 'rules': if (tpl.rulesContent) parts.push(`## Rules\n${tpl.rulesContent}`); break
      case 'tree': 
        if (hasContext) {
          const tree = files.slice(0, 15).map(f => `  ${f}`).join('\n')
          parts.push(`## Project Structure\n${tree}${files.length > 15 ? `\n  ... +${files.length - 15} files` : ''}`)
        } else parts.push(`## Project Structure\n[Build context to see file tree]`)
        break
      case 'stats': parts.push(`## Stats\n- Files: ${contextStore.fileCount || 0}\n- Lines: ${contextStore.lineCount || 0}\n- Tokens: ~${contextStore.tokenCount || 0}`); break
      case 'task': parts.push(`## Task\n[Your task description]`); break
      case 'files': 
        if (hasContext) {
          const fileList = files.slice(0, 5).map(f => `- ${f}`).join('\n')
          parts.push(`## Files\n${fileList}${files.length > 5 ? `\n... +${files.length - 5} more files` : ''}\n\n[File contents]`)
        } else parts.push(`## Files\n[Build context to see files]`)
        break
    }
  }
  if (tpl.customSuffix) parts.push(tpl.customSuffix)
  return parts.join('\n\n')
})

const highlightedPreview = computed(() => {
  if (!previewContent.value) return `<span class="tpl-preview-empty">${t('templates.previewEmpty')}</span>`
  return previewContent.value
    .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
    .replace(/^(## .+)$/gm, '<span class="tpl-h2">$1</span>')
    .replace(/^(- .+)$/gm, '<span class="tpl-list">$1</span>')
    .replace(/\[([^\]]+)\]/g, '<span class="tpl-placeholder">[$1]</span>')
})

const filteredFavorites = computed(() => {
  let list = favoriteTemplates.value
  if (searchQuery.value) { const q = searchQuery.value.toLowerCase(); list = list.filter(t => t.name.toLowerCase().includes(q)) }
  return list
})

const filteredBuiltIn = computed(() => {
  let list = builtInTemplates.value.filter(t => !t.isFavorite && !t.isHidden)
  if (searchQuery.value) { const q = searchQuery.value.toLowerCase(); list = list.filter(t => t.name.toLowerCase().includes(q)) }
  return list
})

const filteredCustom = computed(() => {
  let list = customTemplates.value.filter(t => !t.isHidden)
  if (searchQuery.value) { const q = searchQuery.value.toLowerCase(); list = list.filter(t => t.name.toLowerCase().includes(q)) }
  return list
})

const noResults = computed(() => searchQuery.value && !filteredFavorites.value.length && !filteredBuiltIn.value.length && !filteredCustom.value.length)

watch(autoSaveOnClose, v => localStorage.setItem('template-autosave', v ? 'true' : 'false'))
watch(isOpen, open => { if (open && templateStore.activeTemplate) selectTemplate(templateStore.activeTemplate) })

function autoResize(el: HTMLTextAreaElement) {
  el.style.height = 'auto'
  el.style.height = Math.max(80, el.scrollHeight) + 'px'
}

function selectTemplate(tpl: PromptTemplate) {
  if (hasChanges.value && !isEditingBuiltIn.value) { showUnsavedWarning.value = true; return }
  selectedTemplate.value = tpl
  editingTemplate.value = JSON.parse(JSON.stringify(tpl))
  originalJson.value = JSON.stringify(tpl)
  nextTick(() => {
    if (roleTextarea.value) autoResize(roleTextarea.value)
    if (rulesTextarea.value) autoResize(rulesTextarea.value)
  })
}

function toggleSection(key: keyof TemplateSections) { 
  if (editingTemplate.value) editingTemplate.value.sections[key] = !editingTemplate.value.sections[key] 
}

function createNew() {
  editingTemplate.value = { 
    ...createEmptyTemplate(), id: `new-${Date.now()}`, name: 'New Template', icon: '✨', 
    createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() 
  } as PromptTemplate
  selectedTemplate.value = null
  originalJson.value = ''
}

function saveChanges() {
  if (!editingTemplate.value) return
  if (selectedTemplate.value) {
    templateStore.updateTemplate(selectedTemplate.value.id, editingTemplate.value)
    selectTemplate(editingTemplate.value)
  } else {
    const id = templateStore.createTemplate(editingTemplate.value)
    const created = templateStore.templates.find(t => t.id === id)
    if (created) selectTemplate(created)
  }
  justSaved.value = true
  setTimeout(() => justSaved.value = false, 2000)
}

function handleClose() {
  if (hasChanges.value && !isEditingBuiltIn.value) { 
    autoSaveOnClose.value ? (saveChanges(), setTimeout(close, 100)) : showUnsavedWarning.value = true 
  } else close()
}

function close() { 
  templateStore.closeModal()
  selectedTemplate.value = null
  editingTemplate.value = null
  showUnsavedWarning.value = false 
}

function discardAndClose() { showUnsavedWarning.value = false; close() }
function saveAndClose() { saveChanges(); showUnsavedWarning.value = false; setTimeout(close, 100) }
function confirmDelete(id: string) { deleteConfirmId.value = id }

function executeDelete() { 
  if (deleteConfirmId.value) { 
    templateStore.deleteTemplate(deleteConfirmId.value)
    if (selectedTemplate.value?.id === deleteConfirmId.value) { selectedTemplate.value = null; editingTemplate.value = null } 
    deleteConfirmId.value = null 
  } 
}

function duplicateSelected() { 
  if (!selectedTemplate.value) return
  const id = templateStore.duplicateTemplate(selectedTemplate.value.id)
  if (id) { const c = templateStore.templates.find(t => t.id === id); if (c) { originalJson.value = ''; selectTemplate(c) } } 
}

function applyTemplate() { if (selectedTemplate.value) templateStore.setActiveTemplate(selectedTemplate.value.id) }

function handleKeydown(e: KeyboardEvent) { 
  if (e.key === 'Escape') handleClose()
  if ((e.ctrlKey || e.metaKey) && e.key === 's') { e.preventDefault(); if (hasChanges.value && !isEditingBuiltIn.value) saveChanges() } 
}

function handleExport() { 
  if (!selectedTemplate.value) return
  const blob = new Blob([JSON.stringify(selectedTemplate.value, null, 2)], { type: 'application/json' })
  const a = document.createElement('a'); a.href = URL.createObjectURL(blob)
  a.download = `template-${selectedTemplate.value.id}.json`; a.click() 
}

function handleImport() { 
  const input = document.createElement('input'); input.type = 'file'; input.accept = '.json'
  input.onchange = async (e: Event) => { 
    const file = (e.target as HTMLInputElement).files?.[0]; if (!file) return
    try { 
      const data = JSON.parse(await file.text()) as PromptTemplate
      data.id = `imported-${Date.now()}`; data.isBuiltIn = false
      data.sectionOrder = data.sectionOrder || DEFAULT_SECTION_ORDER
      const id = templateStore.createTemplate(data)
      const c = templateStore.templates.find(t => t.id === id); if (c) selectTemplate(c) 
    } catch { /* ignore */ } 
  }
  input.click() 
}
</script>


<style scoped>
/* ============================================
   TEMPLATE MODAL - Engineering Aesthetic
   ============================================ */

/* Overlay with blur */
.tpl-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
}

/* Modal container */
.tpl-modal {
  width: min(1400px, 96vw);
  height: min(850px, 92vh);
  display: flex;
  flex-direction: column;
  background: #0D0E12;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

/* Header */
.tpl-header {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.tpl-header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tpl-header h2 {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.tpl-header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.tpl-unsaved-indicator {
  width: 6px;
  height: 6px;
  background: var(--color-warning);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.tpl-close {
  padding: 0.375rem;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
}
.tpl-close:hover { background: rgba(255, 255, 255, 0.1); color: var(--text-primary); }

/* Body - 3 columns with CSS Grid */
.tpl-body {
  flex: 1;
  display: grid;
  grid-template-columns: 200px 1fr 320px;
  min-height: 0;
}

/* Sidebar */
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
.tpl-search:focus-within { border-color: var(--accent-indigo-border); }
.tpl-search svg { color: var(--text-subtle); flex-shrink: 0; }
.tpl-search input { 
  flex: 1; background: transparent; border: none; 
  color: var(--text-primary); font-size: 12px; outline: none; min-width: 0;
}
.tpl-search input::placeholder { color: var(--text-subtle); }

.tpl-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 0.5rem;
}

/* Custom scrollbar */
.tpl-list::-webkit-scrollbar { width: 4px; }
.tpl-list::-webkit-scrollbar-track { background: transparent; }
.tpl-list::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 2px; }
.tpl-list::-webkit-scrollbar-thumb:hover { background: rgba(255, 255, 255, 0.2); }

.tpl-group { margin-bottom: 0.5rem; }
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
.tpl-group summary::-webkit-details-marker { display: none; }
.tpl-group summary:hover { color: var(--text-muted); }

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


/* Editor (Builder) */
.tpl-editor {
  display: flex;
  flex-direction: column;
  background: #0D0E12;
  min-width: 0;
  overflow: hidden;
}

.tpl-editor-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.tpl-icon-btn {
  width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.15s;
}
.tpl-icon-btn:hover:not(:disabled) { 
  background: rgba(255, 255, 255, 0.08); 
  border-color: rgba(255, 255, 255, 0.2); 
}
.tpl-icon-btn:disabled { cursor: default; opacity: 0.6; }

.tpl-icon-display { font-size: 1.25rem; line-height: 1; }

.tpl-name-group { flex: 1; display: flex; flex-direction: column; gap: 0.125rem; min-width: 0; }

.tpl-name-input {
  background: transparent;
  border: none;
  border-bottom: 1px solid transparent;
  color: var(--text-primary);
  font-size: 15px;
  font-weight: 600;
  padding: 0.125rem 0;
  transition: border-color 0.15s;
}
.tpl-name-input:hover:not(:disabled) { border-bottom-color: rgba(255, 255, 255, 0.1); }
.tpl-name-input:focus { outline: none; border-bottom-color: var(--accent-indigo); }

.tpl-desc-input {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 11px;
  padding: 0.125rem 0;
}
.tpl-desc-input:focus { outline: none; color: var(--text-secondary); }

.tpl-active-badge {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.625rem;
  background: var(--color-success-soft);
  border: 1px solid var(--color-success-border);
  border-radius: var(--radius-full);
  font-size: 10px;
  font-weight: 500;
  color: var(--color-success);
  white-space: nowrap;
}

/* Section Chips */
.tpl-chips {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.625rem 1rem;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-wrap: wrap;
}

.tpl-chip {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.75rem;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: var(--radius-full);
  color: var(--text-muted);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.tpl-chip:hover { 
  background: rgba(255, 255, 255, 0.04); 
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
}
.tpl-chip.active { 
  background: var(--accent-indigo-bg); 
  border-color: var(--accent-indigo-border); 
  color: var(--accent-indigo); 
}

/* Cards Stack */
.tpl-cards {
  flex: 1;
  overflow-y: auto;
  padding: 0.75rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
}

.tpl-cards::-webkit-scrollbar { width: 4px; }
.tpl-cards::-webkit-scrollbar-track { background: transparent; }
.tpl-cards::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 2px; }

.tpl-textarea {
  width: 100%;
  min-height: 80px;
  max-height: 200px;
  padding: 0;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 13px;
  font-family: var(--font-mono);
  line-height: 1.6;
  resize: none;
  outline: none;
}
.tpl-textarea::placeholder { color: var(--text-subtle); font-style: italic; }

/* Context Options Grid */
.tpl-options-card {
  background: #0d1117;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.tpl-options-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 0.75rem;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}
.tpl-options-header svg { color: var(--text-subtle); }

.tpl-options-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.5rem;
  padding: 0.75rem;
}

/* Advanced Section */
.tpl-advanced {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: var(--radius-md);
}

.tpl-advanced summary {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.625rem 0.75rem;
  font-size: 11px;
  font-weight: 500;
  color: var(--text-muted);
  cursor: pointer;
  user-select: none;
  list-style: none;
}
.tpl-advanced summary::-webkit-details-marker { display: none; }
.tpl-advanced summary:hover { color: var(--text-primary); }

.tpl-advanced .tpl-chevron { transition: transform 0.2s; }
.tpl-advanced[open] .tpl-chevron { transform: rotate(90deg); }

.tpl-advanced-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
  padding: 0 0.75rem 0.75rem;
}

.tpl-advanced-field label {
  display: block;
  margin-bottom: 0.375rem;
  font-size: 10px;
  font-weight: 500;
  color: var(--text-subtle);
  text-transform: uppercase;
}

.tpl-advanced-field textarea {
  width: 100%;
  padding: 0.5rem;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  font-size: 11px;
  font-family: var(--font-mono);
  resize: none;
  outline: none;
  transition: border-color 0.15s;
}
.tpl-advanced-field textarea:focus { border-color: var(--accent-indigo-border); }

.tpl-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  color: var(--text-subtle);
  font-size: 12px;
}


/* Preview Panel */
.tpl-preview {
  display: flex;
  flex-direction: column;
  background: #08090C;
  border-left: 1px solid rgba(255, 255, 255, 0.06);
}

.tpl-preview-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}
.tpl-preview-header svg { color: var(--text-subtle); }

.tpl-preview-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  font-size: 12px;
  font-family: var(--font-mono);
  line-height: 1.7;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
}

.tpl-preview-content::-webkit-scrollbar { width: 4px; }
.tpl-preview-content::-webkit-scrollbar-track { background: transparent; }
.tpl-preview-content::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 2px; }

/* Syntax highlighting in preview */
:deep(.tpl-h2) { color: #60a5fa; font-weight: 600; }
:deep(.tpl-list) { color: #a78bfa; }
:deep(.tpl-placeholder) { color: var(--text-subtle); font-style: italic; }
:deep(.tpl-preview-empty) { color: var(--text-subtle); font-style: italic; }

.tpl-preview-footer {
  padding: 0.75rem 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.tpl-token-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.tpl-token-fill {
  height: 100%;
  background: var(--accent-indigo);
  border-radius: 2px;
  transition: width 0.3s ease-out, background 0.3s;
}
.tpl-token-fill.warning { background: var(--color-warning); }
.tpl-token-fill.danger { background: var(--color-danger); animation: shake 0.5s; }

.tpl-token-info {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
}

.tpl-token-count {
  font-size: 14px;
  font-weight: 600;
  font-family: var(--font-mono);
  color: var(--text-primary);
}

.tpl-token-label {
  font-size: 10px;
  color: var(--text-subtle);
}

/* Footer */
.tpl-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.625rem 1rem;
  background: rgba(255, 255, 255, 0.02);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.tpl-footer-left,
.tpl-footer-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tpl-footer-center {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tpl-footer-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.625rem;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.15s;
}
.tpl-footer-btn:hover:not(:disabled) { 
  background: rgba(255, 255, 255, 0.05); 
  border-color: rgba(255, 255, 255, 0.2);
  color: var(--text-primary);
}
.tpl-footer-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.tpl-saved-msg {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 11px;
  color: var(--color-success);
}

.tpl-unsaved-msg {
  font-size: 11px;
  color: var(--color-warning);
}

/* Autosave Toggle */
.tpl-autosave {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 11px;
  color: var(--text-muted);
}
.tpl-autosave input { display: none; }

.tpl-autosave-track {
  position: relative;
  width: 28px;
  height: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  transition: background 0.2s;
}
.tpl-autosave input:checked + .tpl-autosave-track { background: var(--accent-indigo); }

.tpl-autosave-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 12px;
  height: 12px;
  background: white;
  border-radius: 50%;
  transition: transform 0.2s;
}
.tpl-autosave input:checked + .tpl-autosave-track .tpl-autosave-thumb { transform: translateX(12px); }

/* Action Buttons */
.tpl-apply-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  background: var(--accent-purple-bg);
  border: 1px solid var(--accent-purple-border);
  border-radius: var(--radius-md);
  color: var(--accent-purple);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}
.tpl-apply-btn:hover { 
  background: rgba(168, 85, 247, 0.25);
  transform: translateY(-1px);
}

.tpl-save-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  background: var(--accent-indigo);
  border: none;
  border-radius: var(--radius-md);
  color: white;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}
.tpl-save-btn:hover:not(:disabled) { 
  background: #818cf8;
  transform: translateY(-1px);
}
.tpl-save-btn:disabled { opacity: 0.4; cursor: not-allowed; transform: none; }
.tpl-save-btn.pulse { animation: btnPulse 0.5s; }

.tpl-cancel-btn {
  padding: 0.5rem 0.75rem;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-md);
  color: var(--text-muted);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.15s;
}
.tpl-cancel-btn:hover { background: rgba(255, 255, 255, 0.05); color: var(--text-primary); }

.tpl-danger-btn {
  padding: 0.5rem 0.875rem;
  background: var(--color-danger-soft);
  border: 1px solid var(--color-danger-border);
  border-radius: var(--radius-md);
  color: var(--color-danger);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}
.tpl-danger-btn:hover { background: rgba(248, 113, 113, 0.25); }


/* Dialogs */
.tpl-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: calc(var(--z-modal) + 10);
}

.tpl-dialog {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1.5rem 2rem;
  background: #0D0E12;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-lg);
  text-align: center;
  max-width: 320px;
}
.tpl-dialog.danger svg { color: var(--color-danger); }
.tpl-dialog.warning svg { color: var(--color-warning); }
.tpl-dialog h3 { font-size: 14px; font-weight: 600; color: var(--text-primary); margin: 0; }
.tpl-dialog p { font-size: 12px; color: var(--text-muted); margin: 0; }
.tpl-dialog-btns { display: flex; gap: 0.5rem; margin-top: 0.5rem; }

/* Animations */
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

@keyframes btnPulse {
  0% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0.5); }
  70% { box-shadow: 0 0 0 8px rgba(99, 102, 241, 0); }
  100% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0); }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-2px); }
  75% { transform: translateX(2px); }
}

/* Modal Transitions */
.modal-enter-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.modal-leave-active {
  transition: all 0.2s ease-out;
}
.modal-enter-from {
  opacity: 0;
}
.modal-enter-from .tpl-modal {
  transform: scale(0.95);
  opacity: 0;
}
.modal-leave-to {
  opacity: 0;
}
.modal-leave-to .tpl-modal {
  transform: scale(0.98);
  opacity: 0;
}

/* Fade Transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease-out;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
