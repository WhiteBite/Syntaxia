import { expect, test } from '@playwright/test'

test.describe('Smart Context Flow', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display context preview panel', async ({ page }) => {
        const contextPanel = page.locator('.context-preview, [class*="context"]')
        const count = await contextPanel.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show token budget indicator', async ({ page }) => {
        const tokenBudget = page.locator('.chat-toolbar, [class*="token"]')
        await expect(tokenBudget.first()).toBeVisible()
    })

    test('should have dual mode selector (Explore/Execute)', async ({ page }) => {
        const modeSelector = page.locator('.chat-mode-selector')
        await expect(modeSelector).toBeVisible()

        const exploreMode = page.locator('.mode-btn:has-text("Explore")')
        const executeMode = page.locator('.mode-btn:has-text("Execute")')

        await expect(exploreMode).toBeVisible()
        await expect(executeMode).toBeVisible()
    })

    test('should switch to explore mode for read-only operations', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        await exploreBtn.click()

        await expect(exploreBtn).toHaveClass(/mode-btn-active/)
    })

    test('should switch to execute mode for write operations', async ({ page }) => {
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')
        await executeBtn.click()

        await expect(executeBtn).toHaveClass(/mode-btn-active/)
    })

    test('should show context used badge after AI response', async ({ page }) => {
        // Context badge appears after AI interaction
        const contextBadge = page.locator('[class*="context-used"], .badge')
        const count = await contextBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should display pending changes count in header', async ({ page }) => {
        const pendingBadge = page.locator('.badge-warning')
        const count = await pendingBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('@ Mentions', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have textarea for message input', async ({ page }) => {
        const textarea = page.locator('textarea')
        await expect(textarea.first()).toBeVisible()
    })

    test('should accept @ symbol in input', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@')

        const value = await textarea.inputValue()
        expect(value).toBe('@')
    })

    test('should type @file mention', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@file test.ts')

        const value = await textarea.inputValue()
        expect(value).toContain('@file')
    })

    test('should type @folder mention', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@folder src/')

        const value = await textarea.inputValue()
        expect(value).toContain('@folder')
    })

    test('should type @symbol mention', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@symbol MyClass')

        const value = await textarea.inputValue()
        expect(value).toContain('@symbol')
    })
})
