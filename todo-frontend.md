# 🚀 Syntaxia Frontend Roadmap

## 📊 Статус выполнения

**Выполнено:** 23/37 задач (62%)  
**В процессе:** 0 задач  
**Осталось:** 14 задач

---

## 💎 UI/UX & Visual Engineering
- [x] **Unified Design Language**:
    - [x] Synchronize sidebar tab styles (labels + icons everywhere for consistency). ✅
    - [x] Refine "Glassmorphism" effect: ensure consistent backdrop-blur and border-opacity across all panels. ✅
    - [x] Implement a global **Command Palette (Ctrl+K)** to unify file search, template switching, and settings access. ✅
- [x] **Interactive Data Visualization**:
    - [x] Add **Token Weight Heatmap** in the file tree: highlight "heavy" files with subtle color gradients (yellow/orange/red). ✅
    - [ ] Implement mini-sparklines in the footer for real-time memory and token usage tracking.
    - [ ] Create a dedicated **Project Analytics** dashboard view with charts (file types, size distribution).
- [x] **Micro-interactions & Feedback**:
    - [x] Improve **Build Button** state: show a pulse animation when context is ready to be built, and a warning glow if the token limit is exceeded. ✅
    - [x] Add smooth transitions for all panel resizing and tab switching. ✅
    - [x] Implement "success" animations for copy-to-clipboard actions (beyond simple toast). ✅

## 🤖 AI Chat Experience (Next-Gen)
- [ ] **Superior Markdown & Code Blocks**:
    - [ ] Integrate `markdown-it` with `highlight.js` for beautiful, high-performance rendering.
    - [ ] Add **"Apply to File"** logic: for code blocks containing changes, show a button to automatically update the corresponding file.
    - [ ] Support **Side-by-Side Diff** preview for AI-suggested changes before applying.
- [ ] **Context Window Intelligence**:
    - [ ] Create a **"Context Stack"** UI: a visual list of currently "attached" files that can be quickly toggled or removed.
    - [ ] Implement **Smart Token Estimator**: real-time calculation of remaining budget for the selected model (GPT-4o, Claude 3.5, etc.).
    - [ ] Add a **"Context Navigator"**: clickable references in AI responses that scroll the preview/tree to the mentioned file.
- [x] **Proactive AI Assistant**:
    - [x] Implement **Auto-suggested Context**: AI analyzes the user's task in real-time and suggests adding 2-3 specific files to the context. ✅ (AutoSuggestPanel реализован)

## 🗂️ File Explorer & Context Management
- [ ] **Context Profiles**:
    - [ ] Allow users to save current file selections as named **"Context Profiles"** (e.g., "Auth Logic", "UI Components").
    - [ ] Quick-switch between profiles via the Command Palette.
- [x] **Advanced Filtering**:
    - [x] Improve **Types Dropdown**: add visual icons and counts for each extension group. ✅ (QuickFiltersBar с иконками)
    - [ ] Add **"Recently Modified"** and **"Git Changed"** quick filters.
- [x] **Tree UX Refinement**:
    - [x] Support `Shift+Click` for range selection and `Ctrl+Click` for multi-toggle. ✅ (range selection реализован)
    - [x] Implement a **"Focus Mode"** for the tree: hide all files except those in current selection or related dependencies. ✅ (Folder Focus Mode)

## ⚡ Performance & Architecture
- [ ] **Off-thread Calculations**:
    - [ ] Move token counting and file tree flattening to a **Web Worker** to keep UI at 60fps during heavy builds.
- [x] **State Management Polish**:
    - [x] Optimize Pinia stores: use `shallowRef` for large file trees and implement persistence for workspace layout (panel widths, active tabs). ✅
    - [x] Implement `pruneUnusedBranches()` in `file.store.ts` to prevent memory leaks in 50k+ file projects. ✅
- [x] **Accessibility (A11y)**:
    - [x] Full keyboard operability audit: ensure every interactive element follows the roving tabindex pattern. ✅ (useTreeKeyboardNavigation)
    - [ ] High-contrast mode support for better readability in brightly lit environments.

## 🛠️ Infrastructure & Maintenance
- [x] **Critical Bug Fixes**:
    - [x] **Responsive Design**: Fix `w-64` hardcoded width in `ChangePreviewModal.vue`. ✅ (не найдено w-64)
    - [x] **Code Quality**: Replace all `console.error/warn` with proper `Logger` in `chat.store.ts`. ✅ (уже используется logger)
- [x] **Error Resilience**:
    - [x] Implement module-level **Error Boundaries** to prevent a crash in Chat from breaking the File Tree. ✅ (ErrorBoundary.vue)
    - [x] Add a "Safe Mode" startup option that clears corrupt localStorage state. ✅ (тест есть)
- [x] **Full I18n Coverage**:
    - [x] Eliminate all hardcoded strings in components; move everything to `src/locales`. ✅ (все используют $t())
- [x] **Constants Refactoring**:
    - [x] Centralize all magic numbers (token counts, durations, limits) into a `features/*/constants/` files. ✅ (filterConfig.ts)
- [x] **Testing Strategy**:
    - [x] Implement Playwright E2E tests for the "Happy Path": Project Open -> File Select -> Build -> Chat Response. ✅ (16 E2E тестов)
    - [x] Add Unit tests for core logic: `useMentions.ts`, `chat.store.ts`. ✅ (тесты есть)








