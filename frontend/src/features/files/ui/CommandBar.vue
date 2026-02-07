<template>
  <div class="magic-bar-wrapper">
    <!-- Compact single-row Command Bar -->
    <div class="magic-bar-container px-2 pb-2">
      <div class="magic-bar flex items-center bg-[#1c1f2e] border border-white/10 rounded-xl p-1 shadow-2xl h-11 relative overflow-visible">
        
        <!-- LEFT: Token Limit & Presets -->
        <div class="flex items-center gap-1 px-2 border-r border-white/10" ref="limitRef">
          <BaseButton 
            variant="ghost" size="sm" class="h-8 px-2 flex items-center gap-1 hover:bg-white/5 no-drag"
            @click="toggleDropdown"
          >
            <span class="text-[11px] font-bold text-gray-300">{{ formatTokens(settings.maxTokens) }}</span>
            <ChevronDown class="w-3 h-3 text-gray-500 transition-transform" :class="{ 'rotate-180': showDropdown }" />
          </BaseButton>

          <PresetsDropdown>
            <template #trigger>
              <BaseButton 
                variant="ghost" size="sm" icon-only 
                class="w-8 h-8 text-gray-400 hover:text-indigo-400 transition-colors"
                :title="t('presets.title')"
              >
                <Bookmark class="w-4 h-4" />
              </BaseButton>
            </template>
          </PresetsDropdown>
        </div>

        <!-- CENTER: Stats & Recommendations -->
        <div class="flex-1 flex items-center justify-around px-2 min-w-0">
          <!-- Selected Counter -->
          <div class="flex flex-col items-center justify-center min-w-[60px]" :title="t('commandBar.selected')">
            <span class="text-[10px] font-black text-indigo-400 leading-none">{{ selectedCount }}</span>
            <span class="text-[8px] uppercase text-gray-500 font-bold tracking-tighter">{{ t('files.selected') }}</span>
          </div>

          <!-- Token Stats (Heatmap Buttons) - Simplified to tiny dots -->
          <div class="flex items-center gap-1.5">
            <button
              v-for="level in (['medium', 'heavy', 'critical'] as const)"
              :key="level"
              @click="handleWeightFilterClick(level)"
              class="flex flex-col items-center group transition-all"
              :class="[weightFilter === level ? 'opacity-100 scale-110' : 'opacity-40 hover:opacity-100']"
              :title="levelLabels[level]"
            >
              <span class="text-[9px] font-bold mb-0.5" :class="levelClasses[level].text">{{ tokenStats[level] }}</span>
              <div class="w-2 h-1 rounded-full" :class="levelClasses[level].dot"></div>
            </button>
          </div>

          <!-- Recommendations Indicator (Lamp) -->
          <button 
            v-if="recommendationsCount > 0"
            class="flex items-center justify-center w-8 h-8 rounded-full bg-amber-500/10 text-amber-500 hover:bg-amber-500/20 transition-all relative"
            @click="toggleAnalysisPopup"
            :title="t('context.recommendations')"
          >
            <Sparkles class="w-4 h-4" />
            <span class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-amber-600 text-[8px] font-black text-white border border-[#1c1f2e] animate-pulse">
              {{ recommendationsCount }}
            </span>
          </button>
        </div>

        <!-- RIGHT: Build Button -->
        <div class="pl-1 border-l border-white/10 flex items-center">
          <button 
            class="build-btn relative group h-9 px-5 rounded-lg overflow-hidden transition-all active:scale-95 flex items-center justify-center"
            :class="{ 'opacity-50 grayscale cursor-not-allowed': isButtonDisabled }"
            :disabled="isButtonDisabled"
            @click="handleBuild"
          >
            <!-- Animated Background -->
            <div class="absolute inset-0 bg-gradient-to-r from-indigo-600 via-purple-600 to-indigo-600 bg-[length:200%_auto] animate-gradient group-hover:from-indigo-500 group-hover:to-purple-500 transition-all"></div>
            
            <div class="relative flex items-center gap-2 text-white">
              <Zap v-if="!isBuilding" class="w-4 h-4 fill-current drop-shadow-md" />
              <RefreshCw v-else class="w-4 h-4 animate-spin" />
              <span class="text-[11px] font-black uppercase tracking-widest drop-shadow-md">{{ isBuilding ? t('context.building') : t('commandBar.build') }}</span>
            </div>
            
            <!-- Shimmer effect -->
            <div class="absolute inset-0 opacity-0 group-hover:opacity-20 transition-opacity bg-gradient-to-r from-transparent via-white to-transparent -translate-x-full group-hover:translate-x-full duration-1000"></div>
          </button>
        </div>

        <!-- Limit Dropdown (absolute) -->
        <Transition name="dropdown">
          <div v-if="showDropdown" class="limit-dropdown absolute bottom-full left-1 mb-2 w-48 bg-[#1c1f2e] border border-white/10 rounded-xl shadow-2xl p-1 z-50">
            <div class="px-2 py-1.5 text-[9px] font-bold text-gray-500 uppercase tracking-widest">{{ t('commandBar.customLimit') }}</div>
            <button 
              v-for="preset in tokenPresets" 
              :key="preset.value"
              @click.stop="selectPreset(preset.value)"
              class="flex items-center justify-between w-full px-3 py-2 rounded-lg text-[11px] font-medium transition-colors"
              :class="isPresetActive(preset.value) ? 'bg-indigo-500/20 text-indigo-300' : 'text-gray-400 hover:bg-white/5'"
            >
              <span>{{ preset.label }}</span>
              <span class="text-[9px] opacity-40">{{ preset.model }}</span>
            </button>
            
            <div class="border-t border-white/5 mt-1 pt-1 p-1">
              <div class="relative">
                <input
                  ref="customInputRef"
                  v-model="customTokenValue"
                  type="text"
                  inputmode="numeric"
                  class="w-full bg-black/20 border border-white/5 rounded-lg px-2 py-1.5 text-[11px] text-white focus:border-indigo-500/50 outline-none pr-8"
                  placeholder="Custom (K)..."
                  @click.stop
                  @keydown.enter="applyCustomLimit"
                />
                <Check v-if="customTokenValue" class="absolute right-2 top-1/2 -translate-y-1/2 w-3 h-3 text-emerald-500 cursor-pointer" @click.stop="applyCustomLimit" />
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </div>

    <!-- Selection Bar (Bottom-most strip) -->
    <div v-if="selectedCount > 0" class="flex items-center justify-between px-4 py-1.5 bg-black/40 border-t border-white/5">
       <div class="flex items-center gap-3">
         <div class="flex items-center gap-1.5">
           <div class="w-1.5 h-1.5 rounded-full bg-indigo-500 shadow-[0_0_8px_rgba(99,102,241,0.6)]"></div>
           <span class="text-[10px] text-gray-400 font-medium">~{{ estimatedTokens }}k tokens</span>
         </div>
       </div>

       <div class="flex items-center gap-2">
         <button 
           class="flex items-center gap-1.5 px-2 py-1 rounded text-[9px] font-black tracking-tighter uppercase transition-all border border-transparent hover:bg-white/5"
           :class="noiseReductionEnabled ? 'text-emerald-400' : 'text-gray-500'"
           @click="toggleNoiseReduction"
           :title="noiseReductionTooltip"
         >
           <Sparkles class="w-3 h-3" />
           {{ t('context.noiseReduction') }}
         </button>
         
         <div class="w-px h-3 bg-white/10"></div>

         <button 
           class="p-1 text-gray-600 hover:text-red-400 transition-colors"
           @click="handleClear"
         >
           <X class="w-3.5 h-3.5" />
         </button>
       </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useFileStore } from '@/features/files/model/file.store'
import { useSettingsStore } from '@/stores/settings.store'
import { useAnalysisStatus } from '@/features/files/composables/useAnalysisStatus'
import { BaseButton } from '@/components/ui'
import { Check, ChevronDown, Zap, X, Sparkles, Bookmark, RefreshCw } from 'lucide-vue-next'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { TOKEN_THRESHOLDS } from '@/config/constants'
import type { WeightFilterLevel } from '@/composables/useFileFilter'
import type { FileNode } from '@/types/domain'
import PresetsDropdown from './PresetsDropdown.vue'

const props = defineProps<{
  selectedCount: number
  isBuilding: boolean
}>()

const emit = defineEmits<{
  (e: 'build'): void
}>()

const { t } = useI18n()
const settingsStore = useSettingsStore()
const fileStore = useFileStore()
const settings = computed(() => settingsStore.settings.context)

// Real analysis status integration
const selectedFilesRef = computed(() => Array.from(fileStore.selectedPaths))
const analysis = useAnalysisStatus({
  selectedFiles: selectedFilesRef,
  onAddFiles: (files) => fileStore.selectMultiple(files)
})

const recommendationsCount = computed(() => (analysis.relatedCount.value || 0) + (analysis.dependentCount.value || 0))
const toggleAnalysisPopup = () => {
    if (analysis.relatedCount.value > 0) analysis.showRelatedPopup.value = true
    else if (analysis.dependentCount.value > 0) analysis.showImpactPopup.value = true
}

const showDropdown = ref(false)
const limitRef = ref<HTMLElement | null>(null)
const customTokenValue = ref('')

const isDisabled = computed(() => props.selectedCount === 0)
const isButtonDisabled = computed(() => props.selectedCount === 0 || props.isBuilding)
const estimatedTokens = computed(() => Math.round(fileStore.estimatedTokenCount / 1000))
const weightFilter = computed(() => fileStore.weightFilter)
const noiseReductionEnabled = computed(() => settings.value.enableNoiseReduction)

const levelClasses = {
  medium: { text: 'text-amber-500', dot: 'bg-amber-500' },
  heavy: { text: 'text-orange-500', dot: 'bg-orange-500' },
  critical: { text: 'text-red-500', dot: 'bg-red-500' }
}

const levelLabels = {
  medium: 'Medium (10K+)',
  heavy: 'Heavy (50K+)',
  critical: 'Critical (100K+)'
}

const estimatedSavings = computed(() => {
  if (!noiseReductionEnabled.value) return 0
  let s = 0
  if (settings.value.noiseReductionCollapseImports) s += 10
  if (settings.value.noiseReductionRemoveComments) s += 8
  if (settings.value.noiseReductionRemoveTypeDefinitions) s += 5
  return Math.min(s, 30)
})

const noiseReductionTooltip = computed(() => {
  if (!noiseReductionEnabled.value) return t('context.noiseReductionTooltip')
  return `Saving ~${estimatedSavings.value}% tokens`
})

function toggleNoiseReduction() {
  settingsStore.updateContextSettings({ 
    enableNoiseReduction: !settings.value.enableNoiseReduction 
  })
}

const tokenStats = computed(() => {
  const stats = { medium: 0, heavy: 0, critical: 0, total: 0 }
  const count = (nodes: FileNode[]) => {
    for (const n of nodes) {
      if (!n.isDir && n.size) {
        const tokens = Math.round(n.size / TOKEN_THRESHOLDS.BYTES_PER_TOKEN)
        stats.total++
        if (tokens >= TOKEN_THRESHOLDS.CRITICAL) stats.critical++
        else if (tokens >= TOKEN_THRESHOLDS.HEAVY) stats.heavy++
        else if (tokens >= TOKEN_THRESHOLDS.MEDIUM) stats.medium++
      }
      if (n.children) count(n.children)
    }
  }
  count(fileStore.nodes)
  return stats
})

const tokenPresets = [
  { value: 32000, label: '32K', model: 'GPT-4' },
  { value: 128000, label: '128K', model: 'Turbo' },
  { value: 200000, label: '200K', model: 'Claude' },
  { value: 1000000, label: '1M', model: 'Gemini' },
]

function formatTokens(n: number): string {
  if (n >= 1000000) return `${(n / 1000000).toFixed(1)}M`
  if (n >= 1000) return `${Math.round(n / 1000)}K`
  return n.toString()
}

function isPresetActive(v: number): boolean {
  return Math.abs(settings.value.maxTokens - v) <= v * 0.05
}

function toggleDropdown() { showDropdown.value = !showDropdown.value }

function selectPreset(v: number) {
  settingsStore.updateContextSettings({ maxTokens: v })
  showDropdown.value = false
}

function applyCustomLimit() {
  const val = parseFloat(customTokenValue.value.replace(/[^\d.]/g, ''))
  if (!isNaN(val) && val > 0) {
    settingsStore.updateContextSettings({ maxTokens: Math.round(val * 1000) })
    customTokenValue.value = ''
    showDropdown.value = false
  }
}

function handleBuild() { if (!isDisabled.value) emit('build') }
function handleClear() { fileStore.clearSelection() }

function handleWeightFilterClick(level: WeightFilterLevel) {
  fileStore.weightFilter === level ? fileStore.clearWeightFilter() : fileStore.setWeightFilter(level)
}

function handleClickOutside(event: MouseEvent) {
  if (showDropdown.value && limitRef.value && !limitRef.value.contains(event.target as Node)) {
    showDropdown.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>

<style scoped>
.animate-gradient {
  animation: gradient 3s ease infinite;
}

@keyframes gradient {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.dropdown-enter-active, .dropdown-leave-active { transition: all 0.2s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(10px); }

.no-drag { -webkit-app-region: no-drag; }
</style>
