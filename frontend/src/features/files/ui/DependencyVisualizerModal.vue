<template>
  <BaseModal
    v-model="isOpen"
    :title="modalTitle"
    size="lg"
    @close="handleClose"
  >
    <div class="dep-visualizer">
      <!-- Left panel: Dependencies list -->
      <div class="dep-list">
        <!-- Outgoing section -->
        <div v-if="outgoingDeps.length > 0" class="dep-section">
          <div class="dep-section-header">
            <BaseIcon :icon="ArrowRightCircle" size="sm" class="text-success" />
            <h3>{{ t('files.dependencies.outgoing') }} ({{ outgoingDeps.length }})</h3>
          </div>
          
          <div class="dep-items">
            <div
              v-for="dep in outgoingDeps"
              :key="dep.targetPath"
              class="dep-item"
              @click="handleFileClick(dep.targetPath)"
            >
              <BaseIcon :icon="getFileIcon(dep.targetPath)" size="sm" />
              <span class="dep-item-name">{{ getFileName(dep.targetPath) }}</span>
              <BaseBadge :variant="getDepTypeBadge(dep.type)">
                {{ t(`files.dependencies.types.${dep.type}`) }}
              </BaseBadge>
              <button
                v-if="!isSelected(dep.targetPath)"
                class="btn btn-ghost btn-sm"
                @click.stop="handleAddFile(dep.targetPath)"
              >
                <BaseIcon :icon="Plus" size="sm" />
              </button>
              <BaseIcon
                v-else
                :icon="Check"
                size="sm"
                class="text-success"
              />
            </div>
          </div>
        </div>

        <!-- Incoming section -->
        <div v-if="incomingDeps.length > 0" class="dep-section">
          <div class="dep-section-header">
            <BaseIcon :icon="ArrowLeftCircle" size="sm" class="text-info" />
            <h3>{{ t('files.dependencies.incoming') }} ({{ incomingDeps.length }})</h3>
          </div>
          
          <div class="dep-items">
            <div
              v-for="dep in incomingDeps"
              :key="dep.path"
              class="dep-item"
              @click="handleFileClick(dep.path)"
            >
              <BaseIcon :icon="getFileIcon(dep.path)" size="sm" />
              <span class="dep-item-name">{{ getFileName(dep.path) }}</span>
              <button
                v-if="!isSelected(dep.path)"
                class="btn btn-ghost btn-sm"
                @click.stop="handleAddFile(dep.path)"
              >
                <BaseIcon :icon="Plus" size="sm" />
              </button>
              <BaseIcon
                v-else
                :icon="Check"
                size="sm"
                class="text-success"
              />
            </div>
          </div>
        </div>

        <!-- No dependencies message -->
        <div v-if="outgoingDeps.length === 0 && incomingDeps.length === 0" class="dep-empty">
          <BaseIcon :icon="FileQuestion" size="lg" class="text-secondary" />
          <p>{{ t('files.dependencies.noDependencies') }}</p>
        </div>
      </div>

      <!-- Right panel: Actions & Stats -->
      <div class="dep-sidebar">
        <div class="dep-stats">
          <h4>{{ t('files.dependencies.statistics') }}</h4>
          <div class="stat-item">
            <span class="stat-label">{{ t('files.dependencies.totalDeps') }}</span>
            <span class="stat-value">{{ totalDeps }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">{{ t('files.dependencies.selected') }}</span>
            <span class="stat-value">{{ selectedDepsCount }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">{{ t('files.dependencies.missing') }}</span>
            <span class="stat-value text-warning">{{ missingDepsCount }}</span>
          </div>
        </div>

        <div class="dep-actions">
          <button
            class="btn btn-primary"
            :disabled="missingDepsCount === 0"
            @click="handleAddAllMissing"
          >
            <BaseIcon :icon="PlusCircle" size="sm" />
            {{ t('files.dependencies.addAllMissing') }}
          </button>
          
          <button
            class="btn btn-ghost"
            @click="handleNavigateToFile"
          >
            <BaseIcon :icon="Eye" size="sm" />
            {{ t('files.dependencies.viewInTree') }}
          </button>
        </div>
      </div>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useFileStore } from '../model/file.store'
import { useDependencyGraph } from '@/composables/useDependencyGraph'
import { useI18n } from '@/composables/useI18n'
import BaseModal from '@/components/ui/BaseModal.vue'
import BaseIcon from '@/components/ui/BaseIcon.vue'
import BaseBadge from '@/components/ui/BaseBadge.vue'
import { 
  ArrowRightCircle, 
  ArrowLeftCircle, 
  Plus, 
  Check, 
  PlusCircle, 
  Eye,
  FileQuestion,
  FileCode,
  FileText
} from 'lucide-vue-next'

interface Props {
  filePath: string
  direction?: 'incoming' | 'outgoing' | 'all'
}

const props = withDefaults(defineProps<Props>(), {
  direction: 'all'
})

const isOpen = defineModel<boolean>({ required: true })

const fileStore = useFileStore()
const { findAllDependencies } = useDependencyGraph()
const { t } = useI18n()

// Get all file nodes for dependency analysis
const allNodes = computed(() => {
  const nodes: any[] = []
  function traverse(node: any) {
    if (!node.isDir) nodes.push(node)
    if (node.children) {
      node.children.forEach(traverse)
    }
  }
  fileStore.nodes.forEach(traverse)
  return nodes
})

const dependencies = computed(() => 
  findAllDependencies(props.filePath, allNodes.value)
)

const outgoingDeps = computed(() => 
  dependencies.value.outgoing.map(dep => ({
    targetPath: dep.path,
    type: dep.type
  }))
)

const incomingDeps = computed(() => dependencies.value.incoming)

const totalDeps = computed(() => outgoingDeps.value.length + incomingDeps.value.length)

const selectedDepsCount = computed(() => {
  let count = 0
  outgoingDeps.value.forEach(dep => {
    if (fileStore.selectedPaths.has(dep.targetPath)) count++
  })
  incomingDeps.value.forEach(dep => {
    if (fileStore.selectedPaths.has(dep.path)) count++
  })
  return count
})

const missingDepsCount = computed(() => totalDeps.value - selectedDepsCount.value)

const modalTitle = computed(() => {
  const fileName = getFileName(props.filePath)
  return t('files.dependencies.modalTitle', { file: fileName })
})

function isSelected(path: string): boolean {
  return fileStore.selectedPaths.has(path)
}

function getFileName(path: string): string {
  return path.substring(path.lastIndexOf('/') + 1)
}

function getFileIcon(path: string) {
  if (path.match(/\.(ts|tsx|js|jsx|vue)$/)) return FileCode
  return FileText
}

function getDepTypeBadge(type: string): 'primary' | 'success' | 'warning' | 'info' {
  switch (type) {
    case 'import': return 'primary'
    case 'style': return 'success'
    case 'test': return 'warning'
    case 'type': return 'info'
    default: return 'primary'
  }
}

function handleFileClick(path: string) {
  // Navigate to file in tree
  fileStore.setFocusedPath(path)
}

function handleAddFile(path: string) {
  fileStore.selectPath(path)
}

function handleAddAllMissing() {
  const missing: string[] = []
  
  outgoingDeps.value.forEach(dep => {
    if (!fileStore.selectedPaths.has(dep.targetPath)) {
      missing.push(dep.targetPath)
    }
  })
  
  incomingDeps.value.forEach(dep => {
    if (!fileStore.selectedPaths.has(dep.path)) {
      missing.push(dep.path)
    }
  })
  
  missing.forEach(path => fileStore.selectPath(path))
}

function handleNavigateToFile() {
  fileStore.setFocusedPath(props.filePath)
  isOpen.value = false
}

function handleClose() {
  isOpen.value = false
}
</script>

<style scoped>
.dep-visualizer {
  display: flex;
  gap: 16px;
  height: min(600px, 70vh);
}

.dep-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.dep-section {
  margin-bottom: 24px;
}

.dep-section:last-child {
  margin-bottom: 0;
}

.dep-section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.dep-section-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.dep-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.dep-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background 150ms ease-out;
}

.dep-item:hover {
  background: var(--bg-hover);
}

.dep-item-name {
  flex: 1;
  font-size: 13px;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dep-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
  color: var(--text-secondary);
}

.dep-empty p {
  margin-top: 12px;
  font-size: 14px;
}

.dep-sidebar {
  width: min(280px, 30%);
  border-left: 1px solid var(--border-default);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.dep-stats h4 {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  margin-bottom: 12px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-default);
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  font-variant-numeric: tabular-nums;
}

.dep-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.text-success {
  color: var(--color-success);
}

.text-info {
  color: var(--color-info);
}

.text-secondary {
  color: var(--text-secondary);
}

.text-warning {
  color: var(--color-warning);
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 150ms ease-out;
  border: none;
}

.btn-primary {
  background: var(--accent-purple-bg);
  color: white;
  border: 1px solid var(--accent-purple-border);
}

.btn-primary:hover:not(:disabled) {
  background: var(--accent-purple-hover);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-ghost {
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border-default);
}

.btn-ghost:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}
</style>
