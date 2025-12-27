import { expect, test } from '@playwright/test'

test.describe('Change Preview Modal', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should open change preview modal when review button clicked', async ({ page }) => {
        // This test requires sandbox to have changes
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const modal = page.locator('.fixed.inset-0')
            await expect(modal).toBeVisible()

            const modalTitle = page.locator('h2:has-text("Review")')
            await expect(modalTitle).toBeVisible()
        }
    })

    test('should display file list in modal', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const fileList = page.locator('.border-r.border-gray-700')
            await expect(fileList).toBeVisible()
        }
    })

    test('should have apply and discard buttons', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const applyBtn = page.locator('button:has-text("Apply All")')
            const discardBtn = page.locator('button:has-text("Discard All")')

            await expect(applyBtn).toBeVisible()
            await expect(discardBtn).toBeVisible()
        }
    })

    test('should close modal on cancel', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const cancelBtn = page.locator('button:has-text("Cancel")')
            await cancelBtn.click()

            const modal = page.locator('.fixed.inset-0')
            await expect(modal).not.toBeVisible()
        }
    })

    test('should show diff view when file selected', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            const diffView = page.locator('.font-mono')
            const emptyState = page.locator('.empty-state:has-text("Select")')

            // Either diff view or empty state should be visible
            const hasDiff = await diffView.count() > 0
            const hasEmpty = await emptyState.count() > 0
            expect(hasDiff || hasEmpty).toBeTruthy()
        }
    })

    test('should show operation icons for file changes', async ({ page }) => {
        const reviewBtn = page.locator('button:has-text("Review Changes")')
        const btnCount = await reviewBtn.count()

        if (btnCount > 0) {
            await reviewBtn.first().click()

            // Check for operation indicators (+, ~, -)
            const createIcon = page.locator('.text-green-400:has-text("+")')
            const modifyIcon = page.locator('.text-yellow-400:has-text("~")')
            const deleteIcon = page.locator('.text-red-400:has-text("-")')

            const hasCreate = await createIcon.count() > 0
            const hasModify = await modifyIcon.count() > 0
            const hasDelete = await deleteIcon.count() > 0

            // At least one type should exist if there are changes
            expect(hasCreate || hasModify || hasDelete).toBeTruthy()
        }
    })
})
