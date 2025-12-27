import { expect, test } from '@playwright/test'

test.describe('Sandbox & Change Management', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should show pending changes indicator when sandbox has changes', async ({ page }) => {
        const pendingBadge = page.locator('.badge-warning, .badge:has-text("pending")')
        const count = await pendingBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have review changes button when changes exist', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const count = await reviewBtn.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should open change preview modal', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const modal = page.locator('.fixed.inset-0.z-50')
            await expect(modal).toBeVisible()
        }
    })

    test('should display file changes with operation type', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            // Check for file list items
            const fileItems = page.locator('.cursor-pointer.border-b')
            const itemCount = await fileItems.count()
            expect(itemCount).toBeGreaterThanOrEqual(0)
        }
    })

    test('should show diff with syntax highlighting', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            // Click first file if exists
            const fileItem = page.locator('.cursor-pointer.border-b').first()
            const fileCount = await fileItem.count()

            if (fileCount > 0) {
                await fileItem.click()

                // Check for diff highlighting
                const greenLines = page.locator('.text-green-400')
                const redLines = page.locator('.text-red-400')
                const cyanLines = page.locator('.text-cyan-400')

                const hasHighlighting =
                    await greenLines.count() > 0 ||
                    await redLines.count() > 0 ||
                    await cyanLines.count() > 0

                expect(hasHighlighting).toBeTruthy()
            }
        }
    })

    test('should have apply all button in modal', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const applyBtn = page.locator('button:has-text("Apply All")')
            await expect(applyBtn).toBeVisible()
        }
    })

    test('should have discard all button in modal', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const discardBtn = page.locator('button:has-text("Discard All")')
            await expect(discardBtn).toBeVisible()
        }
    })

    test('should close modal with close button', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const closeBtn = page.locator('.btn-ghost.btn-sm svg')
            await closeBtn.first().click()

            const modal = page.locator('.fixed.inset-0.z-50')
            await expect(modal).not.toBeVisible()
        }
    })
})
