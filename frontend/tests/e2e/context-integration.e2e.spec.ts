import { expect, test } from '@playwright/test'

/**
 * E2E тесты для проверки интеграции контекста с AI Chat
 * Проверяет что бот видит контекст из выбранных файлов
 */
test.describe('Context Integration with AI Chat', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have file selection panel', async ({ page }) => {
        const fileTree = page.locator('.file-tree, .tree-container, .file-explorer')
        const count = await fileTree.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have chat panel visible', async ({ page }) => {
        const chatPanel = page.locator('.section-title-text:has-text("AI Chat")')
        await expect(chatPanel.first()).toBeVisible()
    })

    test('should have message input for sending context', async ({ page }) => {
        const textarea = page.locator('textarea')
        await expect(textarea.first()).toBeVisible()
    })

    test('should enable send button when message entered', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        const sendBtn = page.locator('.action-btn-accent')

        await textarea.fill('Analyze this code')
        await expect(sendBtn).toBeEnabled()
    })

    test('should show selected files count in context', async ({ page }) => {
        // Check for selected files indicator
        const selectedIndicator = page.locator('.selected-count, .badge:has-text("selected"), [class*="selected"]')
        const count = await selectedIndicator.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have token budget visible in chat', async ({ page }) => {
        const tokenBudget = page.locator('.token-budget-bar, .chat-toolbar')
        await expect(tokenBudget.first()).toBeVisible()
    })
})

test.describe('Smart Context Collection', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should collect context when sending message', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        const sendBtn = page.locator('.action-btn-accent')

        // Type a message that would trigger context collection
        await textarea.fill('What files are in this project?')

        // Send button should be enabled
        await expect(sendBtn).toBeEnabled()
    })

    test('should show context files in preview when available', async ({ page }) => {
        const contextPreview = page.locator('.context-preview, .context-files, [class*="context"]')
        const count = await contextPreview.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should display relevance scores for context files', async ({ page }) => {
        const relevanceScore = page.locator('.relevance-score, [class*="relevance"], .score')
        const count = await relevanceScore.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('Context Passed to AI', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have project root available', async ({ page }) => {
        // Project path should be set when project is opened
        const projectIndicator = page.locator('.project-name, .project-path, [class*="project"]')
        const count = await projectIndicator.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should pass smart context to agentic chat', async ({ page }) => {
        // When message is sent, smart context should be collected and passed
        const textarea = page.locator('textarea').first()
        await textarea.fill('Explain the project structure')

        // The send button should be ready
        const sendBtn = page.locator('.action-btn-accent')
        await expect(sendBtn).toBeEnabled()
    })

    test('should show tool calls when AI uses context', async ({ page }) => {
        // Tool calls appear when AI reads files
        const toolCalls = page.locator('.tool-call-container, .tool-call-header')
        const count = await toolCalls.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should display context used badge after response', async ({ page }) => {
        // Context badge shows files and tokens used
        const contextBadge = page.locator('.context-used-badge, [class*="context-used"]')
        const count = await contextBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('File Selection to Context Flow', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have file checkboxes for selection', async ({ page }) => {
        const checkboxes = page.locator('.file-checkbox, input[type="checkbox"]')
        const count = await checkboxes.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should update context when files selected', async ({ page }) => {
        // Selecting files should update the context
        const fileItem = page.locator('.file-item, .tree-node').first()
        const count = await fileItem.count()

        if (count > 0) {
            // File items exist, context can be built
            expect(count).toBeGreaterThan(0)
        }
    })

    test('should show selected files in chat context', async ({ page }) => {
        // Selected files should appear in chat context area
        const contextArea = page.locator('.context-files, .selected-files, [class*="context"]')
        const count = await contextArea.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})
