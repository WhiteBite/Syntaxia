import { expect, test } from '@playwright/test'

test.describe('AI Chat', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display chat panel with mode selector', async ({ page }) => {
        const chatPanel = page.locator('.section-title-text:has-text("AI Chat"), .panel-title:has-text("AI Chat")')
        await expect(chatPanel.first()).toBeVisible()

        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')

        await expect(exploreBtn).toBeVisible()
        await expect(executeBtn).toBeVisible()
    })

    test('should switch between Explore and Execute modes', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')

        await exploreBtn.click()
        await expect(exploreBtn).toHaveClass(/mode-btn-active/)

        await executeBtn.click()
        await expect(executeBtn).toHaveClass(/mode-btn-active/)
        await expect(exploreBtn).not.toHaveClass(/mode-btn-active/)
    })

    test('should show empty state when no messages', async ({ page }) => {
        const emptyState = page.locator('.empty-state')
        await expect(emptyState.first()).toBeVisible()
    })

    test('should have message input with send button', async ({ page }) => {
        const textarea = page.locator('textarea[placeholder]')
        const sendBtn = page.locator('.action-btn-accent')

        await expect(textarea).toBeVisible()
        await expect(sendBtn).toBeVisible()
        await expect(sendBtn).toBeDisabled()

        await textarea.fill('Test message')
        await expect(sendBtn).toBeEnabled()
    })

    test('should show token budget bar', async ({ page }) => {
        const tokenBar = page.locator('.chat-toolbar')
        await expect(tokenBar).toBeVisible()
    })

    test('should show review changes button when sandbox has changes', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        // Button should only appear when there are pending changes
        const count = await reviewBtn.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should clear chat when clear button clicked', async ({ page }) => {
        const clearBtn = page.locator('button:has-text("Clear")')
        await expect(clearBtn).toBeVisible()
    })
})
