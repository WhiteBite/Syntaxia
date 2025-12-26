---
inclusion: fileMatch
fileMatchPattern: "frontend/**/*.{vue,css,ts}"
---

# Responsive Design - Правила

## Основной принцип

**НЕ использовать фиксированные пиксели для layout-критичных размеров.**

## Единицы измерения

| Используй    | Для чего            |
| ------------ | ------------------- |
| `%`          | Процент от родителя |
| `vw` / `vh`  | Процент от viewport |
| `rem` / `em` | Относительно шрифта |
| `fr`         | CSS Grid            |

## Пиксели разрешены

- border, shadow, border-radius
- Иконки фиксированного размера
- padding/margin до 16px
- gap до 16px

## Пиксели запрещены

- width/height панелей
- min-width/max-width контейнеров
- min-height/max-height скроллируемых областей

## Примеры

```css
/* ❌ Неправильно */
max-width: 800px;
min-width: 300px;
max-height: 200px;

/* ✅ Правильно */
max-width: min(800px, 40vw);
min-width: 20%;
max-height: 25vh;
```

## localStorage размеры

При загрузке валидировать против viewport:

```typescript
const MAX_PANEL_WIDTH_PERCENT = 0.4;

function getEffectiveMaxWidth(maxWidth: number): number {
  const viewportMax = Math.floor(window.innerWidth * MAX_PANEL_WIDTH_PERCENT);
  return Math.min(maxWidth, viewportMax);
}
```

## Flexbox overflow

```css
.container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.scrollable-child {
  flex: 1;
  min-height: 0; /* КРИТИЧЕСКИ ВАЖНО */
  overflow-y: auto;
}
```

## Тестирование

Проверять на: 1920x1080, 1366x768, 1280x720, 2560x1440
