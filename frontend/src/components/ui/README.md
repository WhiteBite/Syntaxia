# UI Kit Components

Коллекция переиспользуемых UI компонентов для Syntaxia проекта.

## Обзор компонентов

### Базовые компоненты

| Компонент | Назначение | Когда использовать |
|-----------|------------|-------------------|
| **BaseAlert** | Уведомления и сообщения | Информация, успех, предупреждения, ошибки |
| **BaseButton** | Кнопки с различными вариантами | Любые действия пользователя |
| **BaseInput** | Текстовые поля ввода | Формы, поиск, фильтры |
| **BaseTextarea** | Многострочный текст | Длинные тексты, комментарии |
| **BaseIcon** | Иконки с размерами | Визуальные индикаторы |
| **BaseSpinner** | Индикатор загрузки | Асинхронные операции |
| **BaseSkeleton** | Skeleton загрузка | Плейсхолдеры для загружаемого контента |
| **BaseBadge** | Метки и статусы | Счетчики, статусы, теги |
| **BaseChip** | Компактные теги с действиями | Фильтры, выбранные элементы, удаляемые теги |
| **BaseCard** | Контейнер с границами | Группировка контента |
| **BaseEmptyState** | Пустое состояние | Нет данных, нет результатов |
| **BaseTabs** | Табы с навигацией | Переключение между разделами |

### Интерактивные компоненты

| Компонент | Назначение | Когда использовать |
|-----------|------------|-------------------|
| **BaseModal** | Модальные окна | Диалоги, формы, подтверждения |
| **BaseDropdown** | Выпадающие меню | Списки действий, опции |
| **BasePopover** | Всплывающие подсказки | Дополнительная информация |

## Быстрый старт

### BaseAlert

```vue
<template>
  <!-- Информационное сообщение -->
  <BaseAlert variant="info">
    <template #title>Information</template>
    Select files from the tree to build your context.
  </BaseAlert>

  <!-- Успешное действие -->
  <BaseAlert variant="success" dismissible @dismiss="handleDismiss">
    <template #title>Success</template>
    Context built successfully with 15 files.
  </BaseAlert>

  <!-- Предупреждение с действиями -->
  <BaseAlert variant="warning">
    <template #title>Unsaved Changes</template>
    You have unsaved changes.
    <template #actions>
      <BaseButton variant="warning" size="sm" @click="save">Save</BaseButton>
      <BaseButton variant="ghost" size="sm" @click="discard">Discard</BaseButton>
    </template>
  </BaseAlert>

  <!-- Ошибка -->
  <BaseAlert variant="error" dismissible>
    <template #title>Error</template>
    Failed to build context. Please check your file selection.
  </BaseAlert>
</template>

<script setup lang="ts">
import { BaseAlert, BaseButton } from '@/components/ui'

function handleDismiss() {
  console.log('Alert dismissed')
}
</script>
```

### BaseButton

```vue
<template>
  <!-- Основная кнопка -->
  <BaseButton variant="primary" @click="handleClick">
    Сохранить
  </BaseButton>

  <!-- С иконкой -->
  <BaseButton variant="secondary" :icon="SettingsIcon">
    Настройки
  </BaseButton>

  <!-- Загрузка -->
  <BaseButton variant="primary" :loading="isLoading">
    Загрузка...
  </BaseButton>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import { SettingsIcon } from '@heroicons/vue/24/outline'
</script>
```

### BaseInput

```vue
<template>
  <BaseInput
    v-model="searchQuery"
    label="Поиск"
    placeholder="Введите текст..."
    :prefix-icon="SearchIcon"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { BaseInput } from '@/components/ui'
import { MagnifyingGlassIcon as SearchIcon } from '@heroicons/vue/24/outline'

const searchQuery = ref('')
</script>
```

### BaseModal

```vue
<template>
  <BaseButton @click="modal.open()">Открыть</BaseButton>
  
  <BaseModal v-model="modal.isOpen.value" title="Заголовок">
    <p>Содержимое модального окна</p>
    <template #footer>
      <BaseButton @click="modal.close()">Закрыть</BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { BaseModal, BaseButton } from '@/components/ui'
import { useModal } from '@/composables/useModal'

const modal = useModal()
</script>
```

### BaseEmptyState

```vue
<template>
  <!-- Базовое использование -->
  <BaseEmptyState
    :icon="FolderIcon"
    title="Нет файлов"
    description="Выберите файлы из дерева для построения контекста"
  />

  <!-- С кнопкой действия -->
  <BaseEmptyState
    :icon="RefreshIcon"
    title="Нет данных"
    description="Нажмите обновить для загрузки данных"
  >
    <template #action>
      <BaseButton variant="primary" @click="refresh">
        Обновить
      </BaseButton>
    </template>
  </BaseEmptyState>

  <!-- Компактный размер -->
  <BaseEmptyState
    :icon="InboxIcon"
    title="Пусто"
    size="sm"
  />
</template>

<script setup lang="ts">
import { BaseEmptyState, BaseButton } from '@/components/ui'
import { FolderIcon, RefreshIcon, InboxIcon } from 'lucide-vue-next'

function refresh() {
  // Логика обновления
}
</script>
```

### BaseChip

```vue
<template>
  <!-- Простой chip -->
  <BaseChip variant="primary">JavaScript</BaseChip>

  <!-- Удаляемый chip -->
  <BaseChip 
    variant="success" 
    removable 
    @remove="handleRemove"
  >
    Selected
  </BaseChip>

  <!-- Кликабельный chip с иконкой -->
  <BaseChip 
    clickable 
    :icon="FilterIcon"
    @click="handleClick"
  >
    Filter
  </BaseChip>

  <!-- Chip с кастомной иконкой через slot -->
  <BaseChip variant="warning">
    <template #icon>
      <span>⚠️</span>
    </template>
    Warning
  </BaseChip>

  <!-- Все возможности вместе -->
  <BaseChip 
    variant="primary" 
    size="sm"
    clickable 
    removable 
    :icon="TagIcon"
    @click="handleClick"
    @remove="handleRemove"
  >
    Multi-action
  </BaseChip>

  <!-- Использование в фильтрах -->
  <div class="flex gap-2">
    <BaseChip
      v-for="filter in activeFilters"
      :key="filter.id"
      :variant="filter.active ? 'primary' : 'default'"
      clickable
      removable
      @click="toggleFilter(filter)"
      @remove="removeFilter(filter)"
    >
      <template #icon>
        <span>{{ filter.icon }}</span>
      </template>
      {{ filter.label }}
    </BaseChip>
  </div>
</template>

<script setup lang="ts">
import { BaseChip } from '@/components/ui'
import { Filter as FilterIcon, Tag as TagIcon } from 'lucide-vue-next'

function handleClick() {
  console.log('Chip clicked')
}

function handleRemove() {
  console.log('Chip removed')
}
</script>
```

**Особенности:**
- ✅ 5 вариантов цветов (default, primary, success, warning, danger)
- ✅ 3 размера (xs, sm, md)
- ✅ Removable с кнопкой удаления
- ✅ Clickable для интерактивности
- ✅ Поддержка иконок (prop или slot)
- ✅ Keyboard support (Enter, Space, Delete)
- ✅ Accessibility (ARIA labels, roles)

### BaseSkeleton

```vue
<template>
  <!-- Text skeleton -->
  <BaseSkeleton variant="text" width="200px" />

  <!-- Multiple lines -->
  <BaseSkeleton variant="text" :count="3" />

  <!-- Circle (avatar) -->
  <BaseSkeleton variant="circle" width="40px" height="40px" />

  <!-- Card skeleton -->
  <BaseSkeleton variant="card" height="200px" />

  <!-- Rectangle without animation -->
  <BaseSkeleton variant="rect" width="100%" height="100px" :animated="false" />

  <!-- Complex example: User card -->
  <div class="user-card-skeleton">
    <BaseSkeleton variant="circle" width="64px" height="64px" />
    <div class="user-card-content">
      <BaseSkeleton variant="text" width="120px" height="20px" />
      <BaseSkeleton variant="text" width="180px" height="16px" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseSkeleton } from '@/components/ui'
</script>
```

**Специализированные skeleton компоненты:**

```vue
<template>
  <!-- File tree skeleton -->
  <SkeletonFileTree :rows="10" />

  <!-- Stats cards skeleton -->
  <SkeletonStats :cards="3" />
</template>

<script setup lang="ts">
import SkeletonFileTree from '@/components/SkeletonFileTree.vue'
import SkeletonStats from '@/components/SkeletonStats.vue'
</script>
```

### BaseTabs

```vue
<template>
  <!-- Базовое использование -->
  <BaseTabs v-model="activeTab">
    <BaseTab name="general" label="General">
      <p>General settings content</p>
    </BaseTab>
    <BaseTab name="advanced" label="Advanced">
      <p>Advanced settings content</p>
    </BaseTab>
  </BaseTabs>

  <!-- С иконками -->
  <BaseTabs v-model="activeTab">
    <BaseTab name="general" label="General" :icon="SettingsIcon">
      <p>General settings with icon</p>
    </BaseTab>
    <BaseTab name="advanced" label="Advanced" :icon="CodeIcon">
      <p>Advanced settings with icon</p>
    </BaseTab>
  </BaseTabs>

  <!-- С отключенным табом -->
  <BaseTabs v-model="activeTab">
    <BaseTab name="available" label="Available">
      <p>Available content</p>
    </BaseTab>
    <BaseTab name="locked" label="Locked" :disabled="true">
      <p>This content is not accessible</p>
    </BaseTab>
  </BaseTabs>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { BaseTabs, BaseTab } from '@/components/ui'
import { Settings as SettingsIcon, Code as CodeIcon } from 'lucide-vue-next'

const activeTab = ref('general')
</script>
```

**Особенности:**
- ✅ Lazy loading контента (рендерится только активный таб)
- ✅ Keyboard navigation (Arrow keys, Home, End)
- ✅ Accessibility (ARIA roles, tabindex)
- ✅ Animated indicator для активного таба
- ✅ Поддержка иконок
- ✅ Disabled состояние

## Варианты (Variants)

### Alerts

- `info` - Информационное сообщение (синий)
- `success` - Успешное действие (зеленый)
- `warning` - Предупреждение (желтый)
- `error` - Ошибка (красный)

### Кнопки

- `primary` - Основное действие (градиент, акцент)
- `secondary` - Второстепенное действие (серый фон)
- `ghost` - Прозрачная кнопка
- `danger` - Опасное действие (красный)
- `success` - Успешное действие (зеленый)
- `warning` - Предупреждение (желтый)
- `link` - Текстовая ссылка
- `text` - Простой текст

### Размеры

- `xs` - Очень маленький (24px)
- `sm` - Маленький (32px)
- `md` - Средний (40px) - по умолчанию
- `lg` - Большой (48px)

### Бейджи

- `primary` - Фиолетовый акцент
- `secondary` - Серый
- `success` - Зеленый (успех)
- `warning` - Желтый (предупреждение)
- `danger` - Красный (ошибка)
- `info` - Синий (информация)
- `accent` - Белый на прозрачном

## Best Practices

### ✅ Правильно

```vue
<!-- Используй семантические варианты -->
<BaseButton variant="danger" @click="deleteFile">Удалить</BaseButton>
<BaseButton variant="primary" @click="save">Сохранить</BaseButton>

<!-- Показывай состояние загрузки -->
<BaseButton :loading="isSaving" @click="save">Сохранить</BaseButton>

<!-- Используй иконки для улучшения UX -->
<BaseButton :icon="TrashIcon" variant="danger">Удалить</BaseButton>
```

### ❌ Неправильно

```vue
<!-- Не используй primary для всех кнопок -->
<BaseButton variant="primary">Отмена</BaseButton>

<!-- Не забывай про disabled состояние -->
<BaseButton @click="save">Сохранить</BaseButton> <!-- Должно быть :disabled="!isValid" -->

<!-- Не дублируй стили inline -->
<BaseButton style="background: red">Удалить</BaseButton> <!-- Используй variant="danger" -->
```

## Композиция компонентов

### Форма с валидацией

```vue
<template>
  <form @submit.prevent="handleSubmit">
    <BaseInput
      v-model="form.name"
      label="Имя"
      :error="errors.name"
      placeholder="Введите имя"
    />
    
    <BaseTextarea
      v-model="form.description"
      label="Описание"
      :error="errors.description"
      placeholder="Введите описание"
    />
    
    <div class="flex gap-2">
      <BaseButton type="submit" variant="primary" :loading="isSubmitting">
        Сохранить
      </BaseButton>
      <BaseButton variant="ghost" @click="cancel">
        Отмена
      </BaseButton>
    </div>
  </form>
</template>
```

### Карточка с действиями

```vue
<template>
  <BaseCard>
    <template #header>
      <div class="flex items-center justify-between">
        <h3>Заголовок</h3>
        <BaseDropdown>
          <template #trigger>
            <BaseButton variant="ghost" icon-only :icon="EllipsisIcon" />
          </template>
          <template #content>
            <button @click="edit">Редактировать</button>
            <button @click="delete">Удалить</button>
          </template>
        </BaseDropdown>
      </div>
    </template>
    
    <p>Содержимое карточки</p>
    
    <template #footer>
      <BaseBadge variant="success">Активно</BaseBadge>
    </template>
  </BaseCard>
</template>
```

## Доступность (A11y)

Все компоненты поддерживают:

- ✅ Навигация с клавиатуры (Tab, Enter, Escape)
- ✅ Focus states (outline при фокусе)
- ✅ ARIA атрибуты
- ✅ Screen reader поддержка
- ✅ Disabled состояния

## Темизация

Компоненты используют CSS переменные из `:root`:

```css
/* Цвета */
--accent-indigo
--accent-purple
--color-success
--color-danger
--color-warning

/* Фоны */
--bg-0, --bg-1, --bg-2, --bg-3

/* Текст */
--text-primary, --text-secondary, --text-muted

/* Границы */
--border-default, --border-strong, --border-subtle

/* Размеры */
--space-1, --space-2, --space-3, --space-4
--radius-md, --radius-lg, --radius-xl
```

## Детальная документация

Для подробных примеров смотри файлы в `examples/`:

- [BaseAlert.example.vue](./examples/BaseAlert.example.vue)
- [BaseButton.example.vue](./examples/BaseButton.example.vue)
- [BaseInput.example.vue](./examples/BaseInput.example.vue)
- [BaseTextarea.example.vue](./examples/BaseTextarea.example.vue)
- [BaseIcon.example.vue](./examples/BaseIcon.example.vue)
- [BaseSpinner.example.vue](./examples/BaseSpinner.example.vue)
- [BaseModal.example.vue](./examples/BaseModal.example.vue)
- [BaseDropdown.example.vue](./examples/BaseDropdown.example.vue)
- [BasePopover.example.vue](./examples/BasePopover.example.vue)
- [BaseBadge.example.vue](./examples/BaseBadge.example.vue)
- [BaseChip.example.vue](./examples/BaseChip.example.vue)
- [BaseCard.example.vue](./examples/BaseCard.example.vue)
- [BaseEmptyState.example.vue](./examples/BaseEmptyState.example.vue)
- [BaseTabs.example.vue](./examples/BaseTabs.example.vue)

## Composables

Документация по composables: [composables/ui/README.md](../../composables/ui/README.md)

## Миграция

Гайд по миграции старых компонентов: [UI_KIT_MIGRATION_GUIDE.md](../../../docs/UI_KIT_MIGRATION_GUIDE.md)

## Тестирование

Все компоненты покрыты unit тестами в `tests/unit/components/ui/`.

Запуск тестов:

```bash
cd frontend
npm run test:run
```

## Вопросы?

Если компонент не подходит под твой use case:

1. Проверь примеры в `examples/`
2. Посмотри существующее использование в проекте
3. Расширь компонент через props/slots
4. Создай feature-specific компонент в `features/[name]/ui/`
