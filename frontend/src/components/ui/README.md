# UI Kit Components

Коллекция переиспользуемых UI компонентов для Syntaxia проекта.

## Обзор компонентов

### Базовые компоненты

| Компонент | Назначение | Когда использовать |
|-----------|------------|-------------------|
| **BaseButton** | Кнопки с различными вариантами | Любые действия пользователя |
| **BaseInput** | Текстовые поля ввода | Формы, поиск, фильтры |
| **BaseTextarea** | Многострочный текст | Длинные тексты, комментарии |
| **BaseIcon** | Иконки с размерами | Визуальные индикаторы |
| **BaseSpinner** | Индикатор загрузки | Асинхронные операции |
| **BaseBadge** | Метки и статусы | Счетчики, статусы, теги |
| **BaseCard** | Контейнер с границами | Группировка контента |

### Интерактивные компоненты

| Компонент | Назначение | Когда использовать |
|-----------|------------|-------------------|
| **BaseModal** | Модальные окна | Диалоги, формы, подтверждения |
| **BaseDropdown** | Выпадающие меню | Списки действий, опции |
| **BasePopover** | Всплывающие подсказки | Дополнительная информация |

## Быстрый старт

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

## Варианты (Variants)

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

- [BaseButton.example.vue](./examples/BaseButton.example.vue)
- [BaseInput.example.vue](./examples/BaseInput.example.vue)
- [BaseTextarea.example.vue](./examples/BaseTextarea.example.vue)
- [BaseIcon.example.vue](./examples/BaseIcon.example.vue)
- [BaseSpinner.example.vue](./examples/BaseSpinner.example.vue)
- [BaseModal.example.vue](./examples/BaseModal.example.vue)
- [BaseDropdown.example.vue](./examples/BaseDropdown.example.vue)
- [BasePopover.example.vue](./examples/BasePopover.example.vue)
- [BaseBadge.example.vue](./examples/BaseBadge.example.vue)
- [BaseCard.example.vue](./examples/BaseCard.example.vue)

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
