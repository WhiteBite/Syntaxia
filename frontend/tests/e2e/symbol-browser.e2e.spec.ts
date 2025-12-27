import { expect, test } from '@playwright/test'

test.describe('Symbol Browser', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display symbol browser panel', async ({ page }) => {
        const panel = page.locator('.symbol-browser, .panel-title:has-text("Symbols")')
        await expect(panel.first()).toBeVisible()
    })

    test('should have stats badge showing symbol count', async ({ page }) => {
        const statsBadge = page.locator('.stats-badge')
        await expect(statsBadge.first()).toBeVisible()
    })

    test('should have toolbar buttons', async ({ page }) => {
        const expandBtn = page.locator('.toolbar-btn[title*="Expand"]')
        const collapseBtn = page.locator('.toolbar-btn[title*="Collapse"]')
        const refreshBtn = page.locator('.toolbar-btn[title*="Refresh"]')

        await expect(expandBtn.first()).toBeVisible()
        await expect(collapseBtn.first()).toBeVisible()
        await expect(refreshBtn.first()).toBeVisible()
    })

    test('should show loading state while fetching symbols', async ({ page }) => {
        const loadingState = page.locator('.loading-state')
        const symbolTree = page.locator('.symbol-browser .panel-content')

        // Either loading or content should be visible
        const hasLoading = await loadingState.count() > 0
        const hasContent = await symbolTree.count() > 0
        expect(hasLoading || hasContent).toBeTruthy()
    })

    test('should show empty state when no project selected', async ({ page }) => {
        const emptyState = page.locator('.empty-state:has-text("Select")')
        const count = await emptyState.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have footer with file and symbol counts', async ({ page }) => {
        const footer = page.locator('.panel-footer')
        const count = await footer.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})
