import { expect, test } from '@playwright/test'

/**
 * Performance test for folder selection
 * 
 * NOTE: This test is for Wails desktop app, not web.
 * To run manually:
 * 1. Build app: npm run build
 * 2. Start app: wails dev
 * 3. Open project: D:\Sources\StartUp\Muffin\frontend\muffin
 * 4. Click on 'lib' folder checkbox
 * 5. Check DevTools Performance tab
 */

test.describe('Folder Selection Performance', () => {
    test.skip('should select large folder quickly', async ({ page }) => {
        // This test is skipped because it requires Wails app to be running
        // Use it as a reference for manual testing

        const PROJECT_PATH = 'D:\\Sources\\StartUp\\Muffin\\frontend\\muffin'
        const FOLDER_TO_SELECT = 'lib'

        // Navigate to app (adjust URL based on your Wails dev server)
        await page.goto('http://localhost:34115')

        // Wait for app to load
        await page.waitForSelector('[data-tour="file-tree"]')

        // Open project (this would need to be implemented via Wails API)
        // await page.evaluate((path) => {
        //   window.runtime.EventsEmit('open-project', path)
        // }, PROJECT_PATH)

        // Wait for file tree to load
        await page.waitForSelector('.virtual-tree-scroller')

        // Find the folder
        const folderRow = page.locator(`.tree-row:has-text("${FOLDER_TO_SELECT}")`)
        await expect(folderRow).toBeVisible()

        // Start performance measurement
        await page.evaluate(() => performance.mark('folder-select-start'))

        // Click folder checkbox
        const checkbox = folderRow.locator('.tree-checkbox')
        await checkbox.click()

        // Wait for selection to complete
        await page.waitForFunction(() => {
            const badge = document.querySelector('.badge')
            return badge && parseInt(badge.textContent || '0') > 0
        })

        // End performance measurement
        const duration = await page.evaluate(() => {
            performance.mark('folder-select-end')
            performance.measure('folder-selection', 'folder-select-start', 'folder-select-end')
            const measure = performance.getEntriesByName('folder-selection')[0]
            return measure.duration
        })

        console.log(`Folder selection took: ${duration}ms`)

        // Assert performance (should be under 1 second)
        expect(duration).toBeLessThan(1000)
    })
})
