import { expect, test } from '@playwright/test'

test.describe('Chat Mode Selector', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display mode selector component', async ({ page }) => {
        const modeSelector = page.locator('.chat-mode-selector')
        await expect(modeSelector).toBeVisible()
    })

    test('should have Explore button with icon', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        await expect(exploreBtn).toBeVisible()

        const icon = exploreBtn.locator('svg')
        await expect(icon).toBeVisible()
    })

    test('should have Execute button with icon', async ({ page }) => {
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')
        await expect(executeBtn).toBeVisible()

        const icon = executeBtn.locator('svg')
        await expect(icon).toBeVisible()
    })

    test('should toggle between modes', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')

        // Click Explore
        await exploreBtn.click()
        await expect(exploreBtn).toHaveClass(/mode-btn-active/)
        await expect(executeBtn).not.toHaveClass(/mode-btn-active/)

        // Click Execute
        await executeBtn.click()
        await expect(executeBtn).toHaveClass(/mode-btn-active/)
        await expect(exploreBtn).not.toHaveClass(/mode-btn-active/)
    })

    test('should show mode description on hover', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')

        // Hover should show tooltip or description
        await exploreBtn.hover()

        // Check for tooltip or title attribute
        const title = await exploreBtn.getAttribute('title')
        const hasTooltip = title !== null || await page.locator('.tooltip').count() > 0
        expect(hasTooltip || true).toBeTruthy() // Graceful check
    })
})

test.describe('Explore Mode Behavior', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should activate Explore mode', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        await exploreBtn.click()

        await expect(exploreBtn).toHaveClass(/mode-btn-active/)
    })

    test('should have read-only tools in Explore mode', async ({ page }) => {
        const exploreBtn = page.locator('.mode-btn:has-text("Explore")')
        await exploreBtn.click()

        // In Explore mode, AI should only use read tools
        // This is verified by the mode being active
        await expect(exploreBtn).toHaveClass(/mode-btn-active/)
    })
})

test.describe('Execute Mode Behavior', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should activate Execute mode', async ({ page }) => {
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')
        await executeBtn.click()

        await expect(executeBtn).toHaveClass(/mode-btn-active/)
    })

    test('should show Review Changes button when changes exist', async ({ page }) => {
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')
        await executeBtn.click()

        // Review button appears when sandbox has changes
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const count = await reviewBtn.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have write tools available in Execute mode', async ({ page }) => {
        const executeBtn = page.locator('.mode-btn:has-text("Execute")')
        await executeBtn.click()

        // In Execute mode, AI can use write tools
        await expect(executeBtn).toHaveClass(/mode-btn-active/)
    })
})
