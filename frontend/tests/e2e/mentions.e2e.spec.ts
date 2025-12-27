import { expect, test } from '@playwright/test'

test.describe('Mention Dropdown', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should have chat input textarea', async ({ page }) => {
        const textarea = page.locator('textarea')
        await expect(textarea.first()).toBeVisible()
    })

    test('should accept @ character in input', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@')

        const value = await textarea.inputValue()
        expect(value).toBe('@')
    })

    test('should type @file mention syntax', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@file:test.ts')

        const value = await textarea.inputValue()
        expect(value).toContain('@file')
    })

    test('should type @folder mention syntax', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@folder:src/')

        const value = await textarea.inputValue()
        expect(value).toContain('@folder')
    })

    test('should type @symbol mention syntax', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@symbol:MyClass')

        const value = await textarea.inputValue()
        expect(value).toContain('@symbol')
    })

    test('should type @git mention syntax', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@git:diff')

        const value = await textarea.inputValue()
        expect(value).toContain('@git')
    })

    test('should type @docs mention syntax', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@docs:README')

        const value = await textarea.inputValue()
        expect(value).toContain('@docs')
    })

    test('should enable send button when message has content', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        const sendBtn = page.locator('.action-btn-accent')

        await textarea.fill('Test message with @file:test.ts')
        await expect(sendBtn).toBeEnabled()
    })
})

test.describe('Mention Categories', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/')
    })

    test('should support file mentions', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@file:package.json add dependency')

        const value = await textarea.inputValue()
        expect(value).toContain('@file:package.json')
    })

    test('should support multiple mentions in one message', async ({ page }) => {
        const textarea = page.locator('textarea').first()
        await textarea.fill('@file:a.ts @file:b.ts compare these files')

        const value = await textarea.inputValue()
        expect(value).toContain('@file:a.ts')
        expect(value).toContain('@file:b.ts')
    })
})
