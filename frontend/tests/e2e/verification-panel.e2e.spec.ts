import { expect, test } from '@playwright/test'

test.describe('Verification Panel', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display verification panel', async ({ page }) => {
        const panel = page.locator('.verification-panel, .panel-title:has-text("Verification")')
        await expect(panel.first()).toBeVisible()
    })

    test('should have build/lint/test tabs', async ({ page }) => {
        const buildTab = page.locator('.tab-btn:has-text("Build")')
        const lintTab = page.locator('.tab-btn:has-text("Lint")')
        const testTab = page.locator('.tab-btn:has-text("Test")')

        await expect(buildTab.first()).toBeVisible()
        await expect(lintTab.first()).toBeVisible()
        await expect(testTab.first()).toBeVisible()
    })

    test('should switch between tabs', async ({ page }) => {
        const lintTab = page.locator('.tab-btn:has-text("Lint")')

        await lintTab.first().click()
        await expect(lintTab.first()).toHaveClass(/tab-btn--active/)
    })

    test('should have run all button', async ({ page }) => {
        const runAllBtn = page.locator('.run-all-btn, button:has-text("Run All")')
        await expect(runAllBtn.first()).toBeVisible()
    })

    test('should have group by file/severity buttons', async ({ page }) => {
        const groupByFile = page.locator('.toolbar-btn[title*="file"]')
        const groupBySeverity = page.locator('.toolbar-btn[title*="severity"]')

        await expect(groupByFile.first()).toBeVisible()
        await expect(groupBySeverity.first()).toBeVisible()
    })

    test('should toggle group by mode', async ({ page }) => {
        const groupBySeverity = page.locator('.toolbar-btn[title*="severity"]')

        await groupBySeverity.first().click()
        await expect(groupBySeverity.first()).toHaveClass(/toolbar-btn--active/)
    })

    test('should have expand/collapse/clear buttons', async ({ page }) => {
        const expandBtn = page.locator('.toolbar-btn[title*="Expand"]')
        const collapseBtn = page.locator('.toolbar-btn[title*="Collapse"]')
        const clearBtn = page.locator('.toolbar-btn[title*="Clear"]')

        await expect(expandBtn.first()).toBeVisible()
        await expect(collapseBtn.first()).toBeVisible()
        await expect(clearBtn.first()).toBeVisible()
    })

    test('should show status indicators on tabs', async ({ page }) => {
        const tabStatus = page.locator('.tab-status')
        await expect(tabStatus.first()).toBeVisible()
    })
})
