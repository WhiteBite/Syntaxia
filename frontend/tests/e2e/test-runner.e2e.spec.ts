import { expect, test } from '@playwright/test'

test.describe('Test Runner Panel', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display test runner panel', async ({ page }) => {
        const panel = page.locator('.test-runner-panel, .panel-title:has-text("Tests")')
        await expect(panel.first()).toBeVisible()
    })

    test('should have run all tests button', async ({ page }) => {
        const runAllBtn = page.locator('button:has-text("Run All")')
        await expect(runAllBtn.first()).toBeVisible()
    })

    test('should have run failed tests button', async ({ page }) => {
        const runFailedBtn = page.locator('button:has-text("Run Failed")')
        await expect(runFailedBtn.first()).toBeVisible()
    })

    test('should display test statistics', async ({ page }) => {
        const passedStat = page.locator('.stat-item--passed')
        const failedStat = page.locator('.stat-item--failed')
        const skippedStat = page.locator('.stat-item--skipped')

        await expect(passedStat.first()).toBeVisible()
        await expect(failedStat.first()).toBeVisible()
        await expect(skippedStat.first()).toBeVisible()
    })

    test('should have filter tabs for test results', async ({ page }) => {
        const allTab = page.locator('.filter-tab:has-text("All")')
        const passedTab = page.locator('.filter-tab:has-text("Passed")')
        const failedTab = page.locator('.filter-tab:has-text("Failed")')

        await expect(allTab.first()).toBeVisible()
        await expect(passedTab.first()).toBeVisible()
        await expect(failedTab.first()).toBeVisible()
    })

    test('should switch filter tabs', async ({ page }) => {
        const failedTab = page.locator('.filter-tab:has-text("Failed")')

        await failedTab.first().click()
        await expect(failedTab.first()).toHaveClass(/filter-tab--active/)
    })

    test('should have expand/collapse buttons', async ({ page }) => {
        const expandBtn = page.locator('.toolbar-btn[title*="Expand"]')
        const collapseBtn = page.locator('.toolbar-btn[title*="Collapse"]')

        await expect(expandBtn.first()).toBeVisible()
        await expect(collapseBtn.first()).toBeVisible()
    })

    test('should have clear results button', async ({ page }) => {
        const clearBtn = page.locator('.toolbar-btn[title*="Clear"]')
        await expect(clearBtn.first()).toBeVisible()
    })
})
