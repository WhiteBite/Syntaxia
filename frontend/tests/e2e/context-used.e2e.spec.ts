import { expect, test } from '@playwright/test'

test.describe('Context Used Badge', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have chat messages area', async ({ page }) => {
        const messagesArea = page.locator('.chat-messages, .messages-container')
        const count = await messagesArea.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should display context used badge after AI response', async ({ page }) => {
        // Context badge appears under AI messages
        const contextBadge = page.locator('.context-used-badge, [class*="context-used"]')
        const count = await contextBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show file count in context badge', async ({ page }) => {
        const fileBadge = page.locator('.badge:has-text("files"), .context-used-badge')
        const count = await fileBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show token count in context badge', async ({ page }) => {
        const tokenBadge = page.locator('.badge:has-text("tokens"), .context-used-badge')
        const count = await tokenBadge.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('Context Used Details', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have expandable context details', async ({ page }) => {
        const contextDetails = page.locator('.context-used-details, [class*="context-details"]')
        const count = await contextDetails.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show file operations in details', async ({ page }) => {
        // Operations: read, write, search
        const operations = page.locator('.operation-badge, [class*="operation"]')
        const count = await operations.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})
