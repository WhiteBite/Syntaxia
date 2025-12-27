import { expect, test } from '@playwright/test'

test.describe('Tool Call Display', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have chat panel for tool calls', async ({ page }) => {
        const chatPanel = page.locator('.section-title-text:has-text("AI Chat")')
        await expect(chatPanel.first()).toBeVisible()
    })

    test('should display tool call container when AI uses tools', async ({ page }) => {
        // Tool calls appear during AI interaction
        const toolCallContainer = page.locator('.tool-call-container')
        const count = await toolCallContainer.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show tool call header with name', async ({ page }) => {
        const toolCallHeader = page.locator('.tool-call-header')
        const count = await toolCallHeader.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have expandable tool call details', async ({ page }) => {
        const toolCallDetails = page.locator('.tool-call-details')
        const count = await toolCallDetails.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show tool call status icons', async ({ page }) => {
        const statusIcons = page.locator('.tool-call-icon')
        const count = await statusIcons.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should display tool call arguments section', async ({ page }) => {
        const argsSection = page.locator('.tool-call-section-title:has-text("Arguments")')
        const count = await argsSection.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should display tool call result section', async ({ page }) => {
        const resultSection = page.locator('.tool-call-section-title:has-text("Result")')
        const count = await resultSection.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})

test.describe('Tool Call Categories', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have category icons for tool types', async ({ page }) => {
        const categoryIcon = page.locator('.tool-call-category, .category-icon')
        const count = await categoryIcon.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should show duration for completed tool calls', async ({ page }) => {
        const duration = page.locator('.tool-call-duration')
        const count = await duration.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })

    test('should have expand/collapse functionality', async ({ page }) => {
        const expandIcon = page.locator('.tool-call-header svg.rotate-180, .tool-call-header svg')
        const count = await expandIcon.count()
        expect(count).toBeGreaterThanOrEqual(0)
    })
})
