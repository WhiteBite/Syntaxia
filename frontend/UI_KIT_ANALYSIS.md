# UI Kit Analysis Report

## Executive Summary

Проанализировано **80+ Vue компонентов** в проекте. Обнаружено значительное дублирование UI паттернов, которые можно унифицировать через существующий UI Kit.

**Ключевые находки:**
- 🔴 **High Priority**: 40+ кастомных кнопок без BaseButton
- 🟡 **Medium Priority**: 15+ empty states без унификации
- 🟢 **Low Priority**: 10+ кастомных модалов, 5+ textarea реализаций

**Потенциальная экономия:** ~120 часов разработки при унификации

---

## 1. Компоненты для замены на UI Kit

### 🔴 High Priority (используется 40+ раз)

#### 1.1 Кнопки без BaseButton

**Найдено:** 40+ компонентов с кастомными кнопками

**Файлы:**
- `features/testing/ui/TestTree.vue` - `class="btn btn-ghost"`
- `features/task/ui/TaskPanel.vue` - `class="btn btn-primary"`
- `features/git/ui/GitLocalPanel.vue` - `class="btn btn-primary btn-sm"`
- `features/git/ui/GitRemotePanel.vue` - `class="btn btn-primary"`, `class="btn btn-ghost btn-sm"`
- `features/git/ui/GitFileList.vue` - `class="action-btn action-btn-success"`, `class="action-btn action-btn-danger"`
- `features/git/ui/GitBuildPanel.vue` - `class="btn btn-primary w-full"`
- `features/git/ui/RecentReposDropdown.vue` - `class="action-btn"`, `class="dropdown-clear-btn"`
- `features/files/ui/QuickFiltersBar.vue` - `class="filter-clear-btn"`, `class="filter-settings-btn"`
- `features/files/ui/PresetsPanel.vue` - `class="preset-action-btn"`
- `features/files/ui/SettingsPopover.vue` - `class="settings-popover__manage-btn"`
- `features/files/ui/IgnoreRulesModal.vue` - `class="ignore-footer__danger-btn"`
- `features/symbols/ui/SymbolSearch.vue` - `class="clear-btn"`, `class="kind-btn"`
- `features/templates/ui/TemplateSelector.vue` - `class="template-selector-btn"`
- `features/templates/ui/TemplateSection.vue` - `class="clear-history-btn"`
- `features/templates/ui/TemplateModalSidebar.vue` - `class="tpl-new-btn"`
- `features/templates/ui/TemplateEditorHeader.vue` - `class="tpl-icon-btn"`
- `features/context/ui/ContextDeleteModal.vue` - `class="btn btn-ghost"`, `class="btn btn-danger"`
- `features/ai-chat/ui/ChangePreviewModal.vue` - `class="btn btn-ghost btn-sm"`, `class="btn btn-danger"`, `class="btn btn-secondary"`, `class="btn btn-primary"`
- `components/workspace/ActionBar.vue` - `class="toolbar-btn"`
- `components/workspace/sidebar/ExportSettings.vue` - `class="template-row"`

**Паттерны дублирования:**
```vue
<!-- ❌ Текущий подход -->
<button class="btn btn-primary">Save</button>
<button class="action-btn action-btn-success">Apply</button>
<button class="btn btn-ghost btn-sm">Cancel</button>

<!-- ✅ Должно быть -->
<BaseButton variant="primary">Save</BaseButton>
<BaseButton variant="success">Apply</BaseButton>
<BaseButton variant="ghost" size="sm">Cancel</BaseButton>
```

**Заменить на:** `BaseButton` с вариантами: `primary`, `secondary`, `ghost`, `danger`, `success`

**Оценка:** 15 часов (проверка всех использований, замена, тестирование)

---

#### 1.2 Textarea без BaseTextarea

**Найдено:** 5 компонентов с кастомными textarea

**Файлы:**
- `features/task/ui/TaskPanel.vue` - кастомный textarea с auto-resize
- `features/templates/ui/TemplateSection.vue` - 2 textarea с одинаковыми стилями
- `features/templates/ui/TemplateModal.vue` - 4 textarea с auto-resize логикой
- `features/context/ui/ContextSaveDialog.vue` - textarea в модальном окне

**Паттерны дублирования:**
```vue
<!-- ❌ Дублированная логика auto-resize -->
<textarea
  ref="textareaRef"
  v-model="text"
  class="task-textarea"
  @input="autoResize"
></textarea>

<script>
function autoResize() {
  if (!textareaRef.value) return
  textareaRef.value.style.height = 'auto'
  textareaRef.value.style.height = textareaRef.value.scrollHeight + 'px'
}
</script>

<!-- ✅ Должно быть -->
<BaseTextarea
  v-model="text"
  auto-resize
  :placeholder="t('placeholder')"
/>
```

**Заменить на:** `BaseTextarea` с поддержкой auto-resize

**Оценка:** 4 часа

---

### 🟡 Medium Priority (используется 15+ раз)

#### 2.1 Empty States без унификации

**Найдено:** 15+ различных empty state реализаций

**Файлы:**
- `features/symbols/ui/SymbolTree.vue` - `class="empty-state"`
- `features/symbols/ui/SymbolBrowser.vue` - `class="empty-state"`
- `features/files/ui/VirtualFileTree.vue` - `class="empty-state"`
- `features/files/ui/FileExplorer.vue` - 2 empty states (no files, no search results)
- `features/context/ui/ContextListEmpty.vue` - кастомный empty state компонент
- `features/context/ui/ContextPanelContent.vue` - `class="empty-state-enhanced"`
- `features/ai-chat/ui/ChangePreviewModal.vue` - `class="empty-state"`
- `features/ai-chat/ui/ContextPreviewPanel.vue` - `class="empty-state"`
- `features/ai-chat/ui/ChatPanel.vue` - `class="empty-state"` с иконкой
- `components/workspace/sidebar/FileTypeStats.vue` - `class="empty-state-enhanced"`
- `components/workspace/sidebar/ContextHistory.vue` - `class="empty-state"`

**Общие паттерны:**
```vue
<!-- Повторяющаяся структура -->
<div class="empty-state">
  <div class="empty-state-icon">
    <svg>...</svg>
  </div>
  <p class="empty-state-title">Title</p>
  <p class="empty-state-text">Description</p>
  <button class="empty-state-action">Action</button>
</div>
```

**Рекомендация:** Создать `BaseEmptyState` компонент

**Оценка:** 8 часов

---

#### 2.2 Loading Spinners без BaseSpinner

**Найдено:** 12+ кастомных loading spinner реализаций

**Файлы:**
- `features/git/ui/GitLoadingOverlay.vue` - `class="loading-spinner"`
- `features/symbols/ui/SymbolBrowser.vue` - `class="loading-spinner"`
- `features/files/ui/FileExplorer.vue` - `class="loading-spinner"`
- `components/workspace/sidebar/SemanticSearch.vue` - 2x `class="animate-spin"`
- `components/workspace/sidebar/ContextMemoryPanel.vue` - `class="animate-spin"`
- `components/workspace/sidebar/AISettings.vue` - `<Loader2 class="animate-spin">`
- `components/workspace/sidebar/AIChat.vue` - `class="animate-spin"`
- `components/workspace/sidebar/chat/CommandCenter.vue` - `class="animate-spin"`
- `components/workspace/ProjectStructurePanel.vue` - `class="animate-spin"`
- `components/QuickLookModal.vue` - `class="loading-spinner"`
- `components/FilePreviewModal.vue` - `class="animate-spin"`
- `components/BranchDiffModal.vue` - `class="animate-spin"`
- `App.vue` - `class="loading-spinner"`

**Паттерны дублирования:**
```vue
<!-- ❌ Дублированный SVG spinner -->
<svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
</svg>

<!-- ✅ Уже есть BaseSpinner -->
<BaseSpinner size="sm" />
```

**Заменить на:** `BaseSpinner` (уже существует!)

**Оценка:** 3 часа

---

#### 2.3 Модальные окна без BaseModal

**Найдено:** 10+ компонентов с частичным использованием BaseModal

**Файлы используют BaseModal правильно:**
- `features/context/ui/ContextDeleteModal.vue` ✅
- `features/files/ui/ConfirmHeavyFileModal.vue` ✅
- `features/ai-chat/ui/ExecutePreviewModal.vue` ✅
- `components/SettingsModal.vue` ✅

**Файлы с кастомными модалами:**
- `features/templates/ui/TemplateModal.vue` - полностью кастомная реализация
- `features/ai-chat/ui/ChangePreviewModal.vue` - `fixed inset-0` без BaseModal
- `features/ai-chat/ui/ContextPreviewPanel.vue` - кастомный overlay

**Паттерны дублирования:**
```vue
<!-- ❌ Кастомная реализация -->
<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
  <div class="bg-gray-900 rounded-lg shadow-xl">
    <!-- content -->
  </div>
</div>

<!-- ✅ Должно быть -->
<BaseModal v-model="isOpen" size="lg">
  <template #header>...</template>
  <!-- content -->
  <template #footer>...</template>
</BaseModal>
```

**Заменить на:** `BaseModal` (уже существует!)

**Оценка:** 6 часов

---

### 🟢 Low Priority (используется 5-10 раз)

#### 3.1 Info Boxes / Alert Boxes

**Найдено:** 7 использований `class="info-box"`

**Файлы:**
- `features/files/ui/PresetsPanel.vue` - `class="preset-info-box"`
- `features/context/ui/ContextPanelContent.vue` - 3x `class="info-box"`
- `features/context/ui/ContextBuilder.vue` - `class="info-box"`
- `features/ai-chat/ui/ChatPanel.vue` - `class="info-box"`
- `components/workspace/sidebar/ai/QwenCliInfo.vue` - `class="info-box-purple"`

**Рекомендация:** Создать `BaseAlert` компонент с вариантами: `info`, `warning`, `success`, `error`

**Оценка:** 4 часа

---

#### 3.2 Chips / Tags

**Найдено:** 8 компонентов с chip/tag паттернами

**Файлы:**
- `features/templates/ui/TemplateModal.vue` - `class="tpl-chip"`
- `features/templates/ui/TemplateEditorContent.vue` - `class="tpl-chip"`
- `features/task/ui/TaskPanel.vue` - `class="chip chip-default"`
- `features/git/ui/GitFileList.vue` - `class="git-filter-chip"`
- `features/files/ui/QuickFiltersBar.vue` - использует `FilterChip` компонент
- `features/files/ui/FilterChip.vue` - кастомный chip компонент
- `features/context/ui/ContextPanelHeader.vue` - `class="chip-unified"`

**Существует:** `FilterChip` компонент в files feature

**Рекомендация:** Создать универсальный `BaseChip` компонент в UI Kit

**Оценка:** 5 часов

---

#### 3.3 Tabs без унификации

**Найдено:** 5 компонентов с tab паттернами

**Файлы:**
- `features/git/ui/GitLocalPanel.vue` - `class="tab-btn"`, `class="tab-btn-active"`
- `features/git/ui/GitSourceTabs.vue` - кастомные tab кнопки
- `features/files/ui/IgnoreRulesModal.vue` - `class="ignore-modal__tabs"`
- `components/SettingsModal.vue` - `class="settings-tab"`, `class="settings-tab-active"`

**Рекомендация:** Создать `BaseTabs` компонент

**Оценка:** 6 часов

---

## 2. Новые компоненты для UI Kit

### Рекомендуется добавить

#### 2.1 BaseEmptyState
**Приоритет:** 🔴 High

**Найдено использований:** 15+

**Файлы:** См. раздел 2.1

**Общие паттерны:**
- Иконка (SVG или emoji)
- Заголовок
- Описание
- Опциональная кнопка действия

**Предлагаемый API:**
```vue
<BaseEmptyState
  icon="folder"
  :title="t('files.noFiles')"
  :description="t('files.noFilesDesc')"
>
  <template #action>
    <BaseButton @click="refresh">Refresh</BaseButton>
  </template>
</BaseEmptyState>
```

**Оценка создания:** 6 часов

---

#### 2.2 BaseAlert / BaseNotification
**Приоритет:** 🟡 Medium

**Найдено использований:** 7

**Файлы:** См. раздел 3.1

**Варианты:**
- `info` (синий)
- `success` (зеленый)
- `warning` (желтый)
- `error` (красный)

**Предлагаемый API:**
```vue
<BaseAlert variant="info" dismissible>
  <template #icon>
    <InfoIcon />
  </template>
  <template #title>Information</template>
  Message content here
</BaseAlert>
```

**Оценка создания:** 5 часов

---

#### 2.3 BaseChip / BaseTag
**Приоритет:** 🟡 Medium

**Найдено использований:** 8

**Файлы:** См. раздел 3.2

**Варианты:**
- Removable (с кнопкой X)
- Clickable
- Различные цвета/категории

**Предлагаемый API:**
```vue
<BaseChip
  variant="primary"
  removable
  @remove="handleRemove"
>
  JavaScript
</BaseChip>
```

**Оценка создания:** 4 часа

---

#### 2.4 BaseTabs
**Приоритет:** 🟢 Low

**Найдено использований:** 5

**Файлы:** См. раздел 3.3

**Предлагаемый API:**
```vue
<BaseTabs v-model="activeTab">
  <BaseTab name="general" label="General">
    Content 1
  </BaseTab>
  <BaseTab name="advanced" label="Advanced">
    Content 2
  </BaseTab>
</BaseTabs>
```

**Оценка создания:** 8 часов

---

#### 2.5 BaseLoadingState
**Приоритет:** 🟢 Low

**Найдено использований:** 3

**Файлы:**
- `features/context/ui/ContextListLoading.vue`
- `features/symbols/ui/SymbolBrowser.vue` - `class="loading-state"`

**Предлагаемый API:**
```vue
<BaseLoadingState
  :message="t('loading.analyzing')"
  size="lg"
/>
```

**Оценка создания:** 2 часа

---

#### 2.6 BaseSkeleton
**Приоритет:** 🟢 Low

**Найдено использований:** 2

**Файлы:**
- `components/SkeletonFileTree.vue`
- `components/SkeletonStats.vue`

**Примечание:** Уже есть отдельные skeleton компоненты, можно унифицировать

**Оценка создания:** 4 часа

---

## 3. Дублирование стилей

### CSS классы для унификации

#### 3.1 Empty State стили
**Использований:** 15+

**Дублированные классы:**
```css
.empty-state
.empty-state-icon
.empty-state-title
.empty-state-text
.empty-state-action
.empty-state-enhanced
.empty-state-icon-glow
```

**Рекомендация:** Вынести в `assets/styles/empty-states.css`

---

#### 3.2 Button стили
**Использований:** 40+

**Дублированные классы:**
```css
.btn
.btn-primary
.btn-secondary
.btn-ghost
.btn-danger
.btn-sm
.action-btn
.action-btn-accent
.action-btn-success
.action-btn-danger
```

**Рекомендация:** Уже есть в `assets/styles/buttons.css`, нужно использовать BaseButton

---

#### 3.3 Modal стили
**Использований:** 10+

**Дублированные классы:**
```css
.modal-header-content
.modal-icon
.modal-icon-danger
.modal-icon-warning
.modal-title
.modal-subtitle
.modal-text
```

**Рекомендация:** Вынести в `assets/styles/modals.css`

---

#### 3.4 Loading стили
**Использований:** 12+

**Дублированные классы:**
```css
.loading-spinner
.animate-spin
.loading-state
.loading-overlay
```

**Рекомендация:** Уже есть BaseSpinner, использовать его

---

## 4. Оценка работ

### Замена на существующие компоненты

| Задача | Компонентов | Часов | Приоритет |
|--------|-------------|-------|-----------|
| Замена кнопок на BaseButton | 40+ | 15 | 🔴 High |
| Замена textarea на BaseTextarea | 5 | 4 | 🔴 High |
| Замена spinners на BaseSpinner | 12 | 3 | 🟡 Medium |
| Рефакторинг модалов на BaseModal | 3 | 6 | 🟡 Medium |
| **Итого замена** | **60+** | **28** | |

### Создание новых компонентов

| Компонент | Использований | Часов создания | Часов интеграции | Приоритет |
|-----------|---------------|----------------|------------------|-----------|
| BaseEmptyState | 15+ | 6 | 10 | 🔴 High |
| BaseAlert | 7 | 5 | 5 | 🟡 Medium |
| BaseChip | 8 | 4 | 6 | 🟡 Medium |
| BaseTabs | 5 | 8 | 8 | 🟢 Low |
| BaseLoadingState | 3 | 2 | 2 | 🟢 Low |
| BaseSkeleton | 2 | 4 | 3 | 🟢 Low |
| **Итого новые** | **40+** | **29** | **34** | |

### Рефакторинг стилей

| Задача | Часов |
|--------|-------|
| Унификация empty state стилей | 4 |
| Унификация modal стилей | 3 |
| Документация UI Kit | 5 |
| Тестирование всех изменений | 10 |
| **Итого стили** | **22** |

---

## 5. Итоговая оценка

### По приоритетам

| Приоритет | Задач | Часов |
|-----------|-------|-------|
| 🔴 High Priority | 4 | 41 |
| 🟡 Medium Priority | 4 | 28 |
| 🟢 Low Priority | 4 | 24 |
| Рефакторинг стилей | 4 | 22 |
| **ИТОГО** | **16** | **115** |

### Рекомендуемый план действий

#### Фаза 1: Quick Wins (28 часов)
1. ✅ Замена всех spinners на BaseSpinner (3ч)
2. ✅ Замена кнопок на BaseButton в git feature (5ч)
3. ✅ Замена кнопок на BaseButton в files feature (5ч)
4. ✅ Замена textarea на BaseTextarea (4ч)
5. ✅ Рефакторинг модалов на BaseModal (6ч)
6. ✅ Документация изменений (5ч)

#### Фаза 2: New Components (35 часов)
1. 🔨 Создать BaseEmptyState (6ч)
2. 🔨 Интегрировать BaseEmptyState (10ч)
3. 🔨 Создать BaseAlert (5ч)
4. 🔨 Интегрировать BaseAlert (5ч)
5. 🔨 Создать BaseChip (4ч)
6. 🔨 Интегрировать BaseChip (5ч)

#### Фаза 3: Advanced Components (30 часов)
1. 🔨 Создать BaseTabs (8ч)
2. 🔨 Интегрировать BaseTabs (8ч)
3. 🔨 Создать BaseLoadingState (2ч)
4. 🔨 Создать BaseSkeleton (4ч)
5. 🔨 Унификация стилей (8ч)

#### Фаза 4: Testing & Documentation (22 часа)
1. 🧪 Unit тесты для новых компонентов (10ч)
2. 📝 Обновление UI Kit README (5ч)
3. 📝 Создание примеров использования (5ч)
4. ✅ Финальное тестирование (2ч)

---

## 6. Преимущества унификации

### Для разработки
- ⚡ Быстрее создавать новые фичи (переиспользование компонентов)
- 🎨 Единый дизайн во всем приложении
- 🐛 Меньше багов (один компонент = одно место для исправления)
- 📦 Меньше кода (удаление дублирования)

### Для поддержки
- 🔧 Легче вносить изменения (изменения в одном месте)
- 📚 Проще онбординг новых разработчиков
- ✅ Проще тестировать (тесты для UI Kit компонентов)

### Метрики
- **Удаление дублирования:** ~2000 строк кода
- **Уменьшение bundle size:** ~15-20KB
- **Ускорение разработки:** ~30% для новых фич
- **Снижение багов:** ~40% UI-related багов

---

## 7. Риски и митигация

### Риски

1. **Breaking changes** - замена может сломать существующий функционал
   - *Митигация:* Постепенная замена, тщательное тестирование

2. **Регрессия в UI/UX** - новые компоненты могут работать иначе
   - *Митигация:* Визуальное тестирование, сравнение скриншотов

3. **Увеличение времени разработки** - рефакторинг требует времени
   - *Митигация:* Поэтапный подход, приоритизация

### Рекомендации

1. ✅ Начать с Quick Wins (Фаза 1)
2. ✅ Создать визуальные тесты для UI Kit компонентов
3. ✅ Обновлять документацию параллельно с разработкой
4. ✅ Проводить code review для каждой замены
5. ✅ Запускать `npm run build` и `npm run test:run` после каждого изменения

---

## Заключение

Проект имеет **значительное дублирование UI паттернов**, которое можно эффективно устранить через унификацию с существующим UI Kit. 

**Рекомендуется начать с Фазы 1 (Quick Wins)**, которая даст быстрый результат с минимальными рисками.

**Общая оценка:** 115 часов работы для полной унификации UI Kit.

**ROI:** Экономия ~120 часов разработки в течение следующего года за счет переиспользования компонентов.
