# 🔍 Аудит фронтенда - Краткая сводка

## Найденные проблемы

### 🔴 Критические (High Priority)

#### 1. Модальные окна - ОГРОМНОЕ дублирование
- **30+ компонентов** с `<Teleport to="body">`
- **4-5 разных** вариантов стилей для одного и того же
- **Повторяющаяся логика:** ESC, backdrop click, focus trap
- **Оценка дублирования:** ~2000 строк кода

**Решение:** Создать `BaseModal.vue` + `useModal.ts`

#### 2. Inline SVG иконки - повторяется везде
- **Spinner SVG** дублируется 20+ раз
- **Иконки** inline вместо компонентов
- **Разные размеры:** w-3, w-4, w-5, w-6 без системы

**Решение:** Создать `BaseIcon.vue` + `BaseSpinner.vue`

#### 3. Кнопки - много кастомных классов
- Есть `BaseButton.vue`, но **используется редко**
- Вместо этого: `tpl-footer-btn`, `toolbar-btn`, `action-btn`, `run-btn`, `suggestion-btn`
- **100+ использований** кастомных классов

**Решение:** Улучшить `BaseButton.vue`, добавить loading state

### 🟡 Средние (Medium Priority)

#### 4. Dropdown/Popover - дублирование логики
- **10+ компонентов** с positioning logic
- Каждый реализует свой `Teleport` + click outside
- `FilterDropdownMenu`, `FormatDropdown`, `SettingsPopover`, `StatsPopover`

**Решение:** Создать `BaseDropdown.vue` + `BasePopover.vue`

#### 5. Инпуты - не используется BaseInput
- Есть `BaseInput.vue`, но **используется редко**
- Нет `BaseTextarea.vue`
- Нет search input компонента

**Решение:** Улучшить `BaseInput.vue`, создать `BaseTextarea.vue`

### 🟢 Низкие (Low Priority)

#### 6. Стили - частичное дублирование
- Есть хорошие базовые стили в `assets/styles/`
- Но много компонентов с `<style scoped>` дублируют их
- Особенно в feature modules

**Решение:** Использовать существующие классы, удалить дубли

## 📊 Статистика

```
Модальные окна:     30+ компонентов → 1 BaseModal
Spinner SVG:        20+ дублей → 1 BaseSpinner
Dropdown logic:     10+ компонентов → 1 BaseDropdown
Кастомные кнопки:   100+ использований → BaseButton
Inline иконки:      200+ использований → BaseIcon

Потенциальное сокращение кода: 30-40%
Уменьшение bundle size: 10-15%
```

## 🎯 Рекомендации

### Что делать СЕЙЧАС

1. **BaseModal.vue** - самое большое дублирование
   - Создать универсальный компонент
   - Мигрировать 5-10 простых модалок
   - Потом сложные (TemplateModal, IgnoreRulesModal)

2. **BaseIcon.vue + BaseSpinner.vue**
   - Обернуть lucide-vue-next
   - Заменить все inline SVG
   - Единая система размеров

3. **Улучшить BaseButton.vue**
   - Добавить loading state
   - Добавить icon support
   - Больше вариантов (outline, link)

### Что делать ПОТОМ

4. **BaseDropdown.vue + BasePopover.vue**
   - Универсальная логика позиционирования
   - Click outside, ESC handling
   - Мигрировать все dropdown компоненты

5. **BaseTextarea.vue**
   - Для форм и чатов
   - Auto-resize опция
   - Validation support

6. **Cleanup**
   - Удалить дублирующиеся стили
   - Обновить документацию
   - Создать примеры

## 🚀 Quick Start

### Создать BaseModal (1-2 часа)

```bash
# 1. Создать компонент
cat > frontend/src/components/ui/BaseModal.vue << 'EOF'
<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="modal-overlay" @click="handleBackdropClick">
        <div class="modal-container" :class="`modal-${size}`" @click.stop>
          <div v-if="$slots.header || title" class="modal-header">
            <slot name="header">
              <h2 class="modal-title">{{ title }}</h2>
            </slot>
            <button v-if="showClose" class="modal-close" @click="close">
              <XIcon />
            </button>
          </div>
          <div class="modal-body">
            <slot />
          </div>
          <div v-if="$slots.footer" class="modal-footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
EOF

# 2. Создать composable
cat > frontend/src/composables/ui/useModal.ts << 'EOF'
export function useModal() {
  const isOpen = ref(false)
  
  function open() { isOpen.value = true }
  function close() { isOpen.value = false }
  function toggle() { isOpen.value = !isOpen.value }
  
  // ESC handling
  onKeyStroke('Escape', close)
  
  return { isOpen, open, close, toggle }
}
EOF

# 3. Использовать
# В компоненте:
# <BaseModal v-model="isOpen" title="Confirm">
#   <p>Are you sure?</p>
#   <template #footer>
#     <BaseButton @click="confirm">Yes</BaseButton>
#   </template>
# </BaseModal>
```

### Создать BaseIcon (30 минут)

```vue
<template>
  <component 
    :is="icon" 
    :class="['base-icon', `base-icon--${size}`]"
    :style="{ color }"
  />
</template>

<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  icon: Component
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  color?: string
}

withDefaults(defineProps<Props>(), {
  size: 'md'
})
</script>

<style scoped>
.base-icon--xs { @apply w-3 h-3; }
.base-icon--sm { @apply w-4 h-4; }
.base-icon--md { @apply w-5 h-5; }
.base-icon--lg { @apply w-6 h-6; }
.base-icon--xl { @apply w-8 h-8; }
</style>
```

## 📈 Ожидаемый результат

### До
```vue
<!-- 30 разных модалок с дублированием -->
<Teleport to="body">
  <div v-if="isOpen" class="modal-backdrop" @click="close">
    <div class="modal-container" @click.stop>
      <!-- ... 50 строк кода ... -->
    </div>
  </div>
</Teleport>

<!-- 20 раз повторяется spinner -->
<svg class="animate-spin w-4 h-4" viewBox="0 0 24 24">
  <circle cx="12" cy="12" r="10" stroke="currentColor" />
  <!-- ... -->
</svg>
```

### После
```vue
<!-- Одна строка -->
<BaseModal v-model="isOpen" title="Confirm">
  <p>Content</p>
</BaseModal>

<!-- Одна строка -->
<BaseSpinner size="sm" />
```

## 🎨 UI Kit структура

```
components/ui/
├── BaseButton.vue       ✅ Улучшить
├── BaseInput.vue        ✅ Улучшить
├── BaseTextarea.vue     ❌ Создать
├── BaseModal.vue        ❌ Создать (HIGH PRIORITY)
├── BaseIcon.vue         ❌ Создать (HIGH PRIORITY)
├── BaseSpinner.vue      ❌ Создать (HIGH PRIORITY)
├── BaseDropdown.vue     ❌ Создать
├── BasePopover.vue      ❌ Создать
├── BaseBadge.vue        ✅ Есть
├── BaseCard.vue         ✅ Есть
└── index.ts             ✅ Обновить

composables/ui/
├── useModal.ts          ❌ Создать
├── useDropdown.ts       ❌ Создать
├── usePopover.ts        ❌ Создать
└── index.ts             ❌ Создать
```

## ⏱️ Оценка времени

| Задача | Время | Приоритет |
|--------|-------|-----------|
| BaseModal + useModal | 2-3 часа | 🔴 HIGH |
| BaseIcon + BaseSpinner | 1-2 часа | 🔴 HIGH |
| Улучшить BaseButton | 1-2 часа | 🔴 HIGH |
| Мигрировать 10 модалок | 3-4 часа | 🔴 HIGH |
| BaseDropdown + useDropdown | 2-3 часа | 🟡 MEDIUM |
| BaseTextarea | 1 час | 🟡 MEDIUM |
| BasePopover + usePopover | 2 часа | 🟢 LOW |
| Cleanup + docs | 2-3 часа | 🟢 LOW |

**Итого:** 14-20 часов работы для полного рефакторинга

## 💡 Дополнительные улучшения

### Composables для переиспользования
- `useClickOutside` - для dropdown/popover
- `useKeyboardNav` - для списков
- `useFocusTrap` - для модалок
- `useScrollLock` - для модалок

### Accessibility
- Все компоненты с ARIA атрибутами
- Keyboard navigation из коробки
- Focus management автоматически
- Screen reader support

### Performance
- Lazy loading для модалок
- Virtual scrolling для больших списков
- Мemoization где нужно
- Tree-shaking friendly

## 📝 Следующие шаги

1. **Прочитать** `UI_KIT_REFACTORING_PLAN.md` для деталей
2. **Создать** BaseModal, BaseIcon, BaseSpinner
3. **Мигрировать** 5-10 простых компонентов
4. **Тестировать** и итерировать
5. **Документировать** примеры использования
6. **Cleanup** старый код

---

**Вопросы?** Смотри детальный план в `UI_KIT_REFACTORING_PLAN.md`
