# 🚀 Syntaxia Frontend Roadmap

## 💎 UI/UX & Visual Engineering
- [ ] **Unified Design Language**:
    - [ ] Synchronize sidebar tab styles (labels + icons everywhere for consistency).
    - [ ] Refine "Glassmorphism" effect: ensure consistent backdrop-blur and border-opacity across all panels.
    - [ ] Implement a global **Command Palette (Ctrl+K)** to unify file search, template switching, and settings access.
- [ ] **Interactive Data Visualization**:
    - [ ] Add **Token Weight Heatmap** in the file tree: highlight "heavy" files with subtle color gradients (yellow/orange/red).
    - [ ] Implement mini-sparklines in the footer for real-time memory and token usage tracking.
    - [ ] Create a dedicated **Project Analytics** dashboard view with charts (file types, size distribution).
- [ ] **Micro-interactions & Feedback**:
    - [ ] Improve **Build Button** state: show a pulse animation when context is ready to be built, and a warning glow if the token limit is exceeded.
    - [ ] Add smooth transitions for all panel resizing and tab switching.
    - [ ] Implement "success" animations for copy-to-clipboard actions (beyond simple toast).

## 🤖 AI Chat Experience (Next-Gen)
- [ ] **Superior Markdown & Code Blocks**:
    - [ ] Integrate `markdown-it` with `highlight.js` for beautiful, high-performance rendering.
    - [ ] Add **"Apply to File"** logic: for code blocks containing changes, show a button to automatically update the corresponding file.
    - [ ] Support **Side-by-Side Diff** preview for AI-suggested changes before applying.
- [ ] **Context Window Intelligence**:
    - [ ] Create a **"Context Stack"** UI: a visual list of currently "attached" files that can be quickly toggled or removed.
    - [ ] Implement **Smart Token Estimator**: real-time calculation of remaining budget for the selected model (GPT-4o, Claude 3.5, etc.).
    - [ ] Add a **"Context Navigator"**: clickable references in AI responses that scroll the preview/tree to the mentioned file.
- [ ] **Proactive AI Assistant**:
    - [ ] Implement **Auto-suggested Context**: AI analyzes the user's task in real-time and suggests adding 2-3 specific files to the context.

## 🗂️ File Explorer & Context Management
- [ ] **Context Profiles**:
    - [ ] Allow users to save current file selections as named **"Context Profiles"** (e.g., "Auth Logic", "UI Components").
    - [ ] Quick-switch between profiles via the Command Palette.
- [ ] **Advanced Filtering**:
    - [ ] Improve **Types Dropdown**: add visual icons and counts for each extension group.
    - [ ] Add **"Recently Modified"** and **"Git Changed"** quick filters.
- [ ] **Tree UX Refinement**:
    - [ ] Support `Shift+Click` for range selection and `Ctrl+Click` for multi-toggle.
    - [ ] Implement a **"Focus Mode"** for the tree: hide all files except those in current selection or related dependencies.

## ⚡ Performance & Architecture
- [ ] **Off-thread Calculations**:
    - [ ] Move token counting and file tree flattening to a **Web Worker** to keep UI at 60fps during heavy builds.
- [ ] **State Management Polish**:
    - [ ] Optimize Pinia stores: use `shallowRef` for large file trees and implement persistence for workspace layout (panel widths, active tabs).
    - [ ] Implement `pruneUnusedBranches()` in `file.store.ts` to prevent memory leaks in 50k+ file projects.
- [ ] **Accessibility (A11y)**:
    - [ ] Full keyboard operability audit: ensure every interactive element follows the roving tabindex pattern.
    - [ ] High-contrast mode support for better readability in brightly lit environments.

## 🛠️ Infrastructure & Maintenance
- [ ] **Critical Bug Fixes**:
    - [ ] **Responsive Design**: Fix `w-64` hardcoded width in `ChangePreviewModal.vue`.
    - [ ] **Code Quality**: Replace all `console.error/warn` with proper `Logger` in `chat.store.ts`.
- [ ] **Error Resilience**:
    - [ ] Implement module-level **Error Boundaries** to prevent a crash in Chat from breaking the File Tree.
    - [ ] Add a "Safe Mode" startup option that clears corrupt localStorage state.
- [ ] **Full I18n Coverage**:
    - [ ] Eliminate all hardcoded strings in components; move everything to `src/locales`.
- [ ] **Constants Refactoring**:
    - [ ] Centralize all magic numbers (token counts, durations, limits) into a `features/ai-chat/constants/index.ts` file.
- [ ] **Testing Strategy**:
    - [ ] Implement Playwright E2E tests for the "Happy Path": Project Open -> File Select -> Build -> Chat Response.
    - [ ] Add Unit tests for core logic: `useChatMessages.ts`, `useMentions.ts`, `chat.store.ts`.
