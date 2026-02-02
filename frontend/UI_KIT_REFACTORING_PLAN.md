# UI Kit Refactoring Plan - Syntaxia Frontend

## 🎯 Цель
Создать централизованный UI Kit для устранения дублирования кода, улучшения консистентности и упрощения поддержки.

## 📊 Анализ текущего состояния

### Найденные проблемы

#### 1. **Модальные окна (30+ компонентов)**
- ❌ Дублирование структуры `<Teleport to="body">` + backdrop + container
- ❌ Разные стили для одинаковых элементов (modal-overlay, modal-backdrop, tpl-overlay)
- ❌ Повторяющаяся логика закрытия (ESC, backdrop click)
- ❌ Разные анимации для одного и того же

**Примеры:**
- `ConfirmHeavyFileModal.vue` - кастомные стили
- `ContextSaveDialog.vue` - другие кастомные стили
- `FilterSettingsModal.vue` - третий вариант стилей
- `TemplateModal.vue` - четвертый вариант

#### 2. **Кнопки (100+ использований)**
- ✅ Есть базовые классы в `buttons.css`
- ❌ Но много кастомных классов: `tpl-footer-btn`, `toolbar-btn`, `action-btn`, `run-btn`
- ❌ Inline SVG иконки вместо компонентов
- ❌ Дублирование loading состояний (spinner SVG повторяется 20+ раз)

#### 3. **Инпуты**
- ✅ Есть `BaseInput.vue`
- ❌ Но используется редко, вместо этого inline `<input>` с кастомными стилями
- ❌ Нет компонента для textarea
- ❌ Нет компонента для search input с иконкой

#### 4. **Иконки**
- ❌ Inline SVG повторяется везде
- ❌ Разные размеры: `w-3 h-3`, `w-4 h-4`, `w-5 h-5`, `w-6 h-6`
- ✅ Есть `lucide-vue-next`, но используется непоследовательно

#### 5. **Dropdown/Popover**
- ❌ Каждый компонент реализует свою логику позиционирования
- ❌ Дублирование Teleport + positioning logic

## 🎨 Предлагаемый UI Kit

### Структура

```
frontend/src/components/ui/
├── BaseButton.vue          ✅ Есть, нужно улучшить
├── BaseInput.vue           ✅ Есть, нужно улучшить
├── BaseTextarea.vue        ❌ Создать
├── BaseModal.vue           ❌ Создать
├── BaseDropdown.vue        ❌ Создать
├── BasePopover.vue         ❌ Создать
├── BaseIcon.vue            ❌ Создать
├── BaseSpinner.vue         ❌ Создать
├── BaseBadge.vue           ✅ Есть
├── BaseCard.vue            ✅ Есть
├── BaseTooltip.vue         ❌ Создать (есть Tooltip.vue, переименовать)
├── ToggleSwitch.vue        ✅ Есть
└── index.ts                ✅ Есть, обновить
```

### Composables для логики

```
frontend/src/composables/ui/
├── useModal.ts             ❌ Создать (ESC, backdrop, focus trap)
├── useDropdown.ts          ❌ Создать (positioning, click outside)
├── usePopover.ts           ❌ Создать (positioning, hover/click)
├── useTooltip.ts           ❌ Создать
└── index.ts                ❌ Создать
```

## 📝 Детальный план компонентов

### 1. BaseModal.vue

**Props:**
```typescript
{
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  closeOnBackdrop?: boolean
  closeOnEsc?: boolean
  showClose?: boolean
  persistent?: boolean
}
```

**Slots:**
```vue
<slot name="header" />
<slot /> <!-- body -->
<slot name="footer" />
```

**Заменит:**
- ConfirmHeavyFileModal
- ContextSaveDialog
- FilterSettingsModal
- FileQuickOpenModal
- IgnoreRulesModal
- PresetsPanel (модалки)
- И еще 20+ компонентов

### 2. BaseButton.vue (улучшить существующий)

**Добавить:**
- `loading` prop с встроенным спиннером
- Больше вариантов: `outline`, `link`, `text`
- `iconPosition: 'left' | 'right'`
- Автоматическая обработка иконок из lucide-vue-next

**Пример использования:**
```vue
<BaseButton 
  variant="primary" 
  :loading="isLoading"
  :icon="PlayIcon"
>
  Run Tests
</BaseButton>
```

### 3. BaseIcon.vue

**Props:**
```typescript
{
  name: string | Component  // lucide icon
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  color?: string
}
```

**Размеры:**
- xs: 12px (w-3 h-3)
- sm: 16px (w-4 h-4)
- md: 20px (w-5 h-5)
- lg: 24px (w-6 h-6)
- xl: 32px (w-8 h-8)

**Заменит:** Все inline SVG

### 4. BaseSpinner.vue

**Props:**
```typescript
{
  size?: 'xs' | 'sm' | 'md' | 'lg'
  color?: string
}
```

**Заменит:** 20+ дублирований spinner SVG

### 5. BaseTextarea.vue

**Props:**
```typescript
{
  modelValue: string
  label?: string
  placeholder?: string
  rows?: number
  maxLength?: number
  error?: string
  resize?: boolean
  autoResize?: boolean
}
```

### 6. BaseDropdown.vue

**Props:**
```typescript
{
  modelValue: boolean
  placement?: 'bottom-start' | 'bottom-end' | 'top-start' | 'top-end'
  offset?: number
  closeOnClick?: boolean
}
```

**Slots:**
```vue
<slot name="trigger" />
<slot /> <!-- content -->
```

**Заменит:**
- FilterDropdownMenu
- FormatDropdown
- RecentReposDropdown
- MentionDropdown

### 7. BasePopover.vue

**Props:**
```typescript
{
  modelValue: boolean
  trigger?: 'click' | 'hover'
  placement?: 'top' | 'bottom' | 'left' | 'right'
  offset?: number
}
```

**Заменит:**
- SettingsPopover
- StatsPopover
- ImpactAnalysisPopup

## 🔄 План миграции

### Фаза 1: Создание базовых компонентов (1-2 дня)
1. ✅ BaseButton - улучшить
2. ✅ BaseInput - улучшить
3. ❌ BaseTextarea - создать
4. ❌ BaseIcon - создать
5. ❌ BaseSpinner - создать

### Фаза 2: Модальные окна (2-3 дня)
1. ❌ BaseModal - создать
2. ❌ useModal composable - создать
3. ❌ Мигрировать 5-10 простых модалок
4. ❌ Мигрировать сложные модалки (TemplateModal, IgnoreRulesModal)

### Фаза 3: Dropdown/Popover (1-2 дня)
1. ❌ BaseDropdown - создать
2. ❌ BasePopover - создать
3. ❌ useDropdown, usePopover composables
4. ❌ Мигрировать все dropdown компоненты

### Фаза 4: Cleanup (1 день)
1. ❌ Удалить дублирующиеся стили
2. ❌ Обновить документацию
3. ❌ Создать Storybook/примеры

## 📈 Ожидаемые результаты

### Метрики
- **Уменьшение кода:** ~30-40% в компонентах
- **Удаление дублирования:** 
  - Модалки: 30+ → 1 компонент
  - Spinner SVG: 20+ → 1 компонент
  - Dropdown logic: 10+ → 1 composable
- **Консистентность:** 100% единообразие UI
- **Bundle size:** -10-15% за счет tree-shaking

### Преимущества
- ✅ Единый источник правды для UI
- ✅ Легче поддерживать и обновлять
- ✅ Автоматическая accessibility
- ✅ Меньше багов из-за дублирования
- ✅ Быстрее разработка новых фич

## 🎯 Приоритеты

### High Priority (делать первым)
1. **BaseModal** - самое большое дублирование
2. **BaseIcon** - используется везде
3. **BaseSpinner** - дублируется 20+ раз

### Medium Priority
4. **BaseButton** - улучшить существующий
5. **BaseTextarea** - нужен для форм
6. **BaseDropdown** - много дублирования

### Low Priority
7. **BasePopover** - меньше использований
8. **BaseTooltip** - уже есть, нужно улучшить

## 🚀 Начало работы

### Шаг 1: Создать BaseModal
```bash
# Создать компонент
touch frontend/src/components/ui/BaseModal.vue

# Создать composable
touch frontend/src/composables/ui/useModal.ts

# Создать тесты
touch frontend/tests/unit/components/ui/BaseModal.spec.ts
```

### Шаг 2: Мигрировать простую модалку
Начать с `ConfirmHeavyFileModal.vue` - простая структура, хороший пример.

### Шаг 3: Документация
Создать примеры использования в `frontend/src/components/ui/README.md`

## 📚 Дополнительно

### Стили
- Использовать существующие CSS переменные из `design-tokens.css`
- Не создавать новые классы, использовать Tailwind + CSS vars
- Все анимации через `animations.css`

### TypeScript
- Строгая типизация всех props
- Экспорт типов для переиспользования
- Generic типы где нужно

### Accessibility
- ARIA атрибуты по умолчанию
- Keyboard navigation
- Focus management
- Screen reader support

### Testing
- Unit тесты для всех компонентов
- Integration тесты для composables
- Visual regression тесты (опционально)
