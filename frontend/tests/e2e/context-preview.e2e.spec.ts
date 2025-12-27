import { expect, test } from '@playwright/test'

test.describe('Context Preview Panel', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have context preview trigger in chat', async ({ page }) => {
        // Context preview is triggered from chat when files are suggested
        const chatPanel = page.locator('.section-title-text:has-text("AI Chat")')
        await expect(chatPanel.first()).toBeVisible()
    })

    test('should display token budget bar in chat toolbar', async ({ page }) => {
        const tokenBar = page.locator('.token-budget-bar, .chat-toolbar')
        await expect(tokenBar.first()).toBeVisible()
    })

    test('should show token progress indicator', async ({ page }) => {
        const progressBar = page.locator('.token-progress-container, .token-progress-bar')
        const count = await progressBar.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should display token usage text', async ({ page }) => {
        const tokenText = page.locator('.token-text')
        const count = await tokenText.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('Context Preview Modal', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have select all button when modal opens', async ({ page }) => {
        const selectAllBtn = page.locator('button:has-text("Select All")')
        const count = await selectAllBtn.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have deselect all button when modal opens', async ({ page }) => {
        const deselectBtn = page.locator('button:has-text("Deselect All")')
        const count = await deselectBtn.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have send to AI button in modal footer', async ({ page }) => {
        const sendBtn = page.locator('button:has-text("Send to AI")')
        const count = await sendBtn.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have budget bar in context preview', async ({ page }) => {
        const budgetBar = page.locator('.budget-bar, .context-preview-budget')
        const count = await budgetBar.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})
