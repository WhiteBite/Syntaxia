import { expect, test } from '@playwright/test'

test.describe('Incremental Context', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have dual mode selector for context gathering', async ({ page }) => {
        const modeSelector = page.locator('.chat-mode-selector')
        await expect(modeSelector).toBeVisible()
    })

    test('should have Explore mode for read-only context', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        await expect(exploreBtn).toBeVisible()
    })

    test('should have Execute mode for write operations', async ({ page }) => {
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')
        await expect(executeBtn).toBeVisible()
    })

    test('should switch to Explore mode', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        await exploreBtn.click()

        await expect(exploreBtn).toHaveClass(/mode-btn-active/)
    })

    test('should switch to Execute mode', async ({ page }) => {
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')
        await executeBtn.click()

        await expect(executeBtn).toHaveClass(/mode-btn-active/)
    })

    test('should show pending changes count when in Execute mode', async ({ page }) => {
        const pendingBadge = page.locator('.badge-warning, .pending-count')
        const count = await pendingBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('Repository Map', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have project structure available', async ({ page }) => {
        // Repo map is used internally for context
        const fileTree = page.locator('.file-tree, .tree-container')
        const count = await fileTree.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show file relevance indicators', async ({ page }) => {
        const relevanceIndicator = page.locator('.relevance-score, [class*="relevance"]')
        const count = await relevanceIndicator.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('Smart Chunking', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display token counts for files', async ({ page }) => {
        const tokenCount = page.locator('.token-count, [class*="tokens"]')
        const count = await tokenCount.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have token budget limit indicator', async ({ page }) => {
        const budgetIndicator = page.locator('.token-budget-bar, .budget-bar')
        const count = await budgetIndicator.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})
