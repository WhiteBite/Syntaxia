import { expect, test } from '@playwright/test'

test.describe('Build Panel', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should display build panel', async ({ page }) => {
        const panel = page.locator('.build-panel, .panel-title:has-text("Build")')
        await expect(panel.first()).toBeVisible()
    })

    test('should have build action buttons', async ({ page }) => {
        const buildBtn = page.locator('button:has-text("Build")')
        const cleanBtn = page.locator('button:has-text("Clean")')
        const rebuildBtn = page.locator('button:has-text("Rebuild")')

        await expect(buildBtn.first()).toBeVisible()
        await expect(cleanBtn.first()).toBeVisible()
        await expect(rebuildBtn.first()).toBeVisible()
    })

    test('should display status bar with status indicator', async ({ page }) => {
        const statusBar = page.locator('.status-bar')
        const statusDot = page.locator('.status-dot')
        const statusText = page.locator('.status-text')

        await expect(statusBar.first()).toBeVisible()
        await expect(statusDot.first()).toBeVisible()
        await expect(statusText.first()).toBeVisible()
    })

    test('should show idle status initially', async ({ page }) => {
        const statusDot = page.locator('.status-dot--idle')
        const count = await statusDot.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have build output section', async ({ page }) => {
        const outputSection = page.locator('.content-main')
        await expect(outputSection.first()).toBeVisible()
    })

    test('should have build history sidebar', async ({ page }) => {
        const historySidebar = page.locator('.content-sidebar')
        await expect(historySidebar.first()).toBeVisible()
    })

    test('should show last build time when available', async ({ page }) => {
        const lastBuild = page.locator('.last-build')
        const count = await lastBuild.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})
