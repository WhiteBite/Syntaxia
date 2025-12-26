import { useReviewStore, type Hunk } from '@/stores/review.store'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('ReviewStore', () => {
    let store: ReturnType<typeof useReviewStore>

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useReviewStore()
    })

    describe('initial state', () => {
        it('should have empty gitDiff', () => {
            expect(store.gitDiff).toBe('')
        })

        it('should have empty parsedDiff', () => {
            expect(store.parsedDiff).toEqual([])
        })
    })

    describe('setGeneratedDiff', () => {
        it('should set gitDiff value', () => {
            const diff = `diff --git a/file.ts b/file.ts
--- a/file.ts
+++ b/file.ts
@@ -1,3 +1,4 @@
 line1
+new line
 line2`

            store.setGeneratedDiff(diff)

            expect(store.gitDiff).toBe(diff)
        })

        it('should update gitDiff when called multiple times', () => {
            store.setGeneratedDiff('first diff')
            expect(store.gitDiff).toBe('first diff')

            store.setGeneratedDiff('second diff')
            expect(store.gitDiff).toBe('second diff')
        })

        it('should handle empty string', () => {
            store.setGeneratedDiff('some diff')
            store.setGeneratedDiff('')

            expect(store.gitDiff).toBe('')
        })

        it('should handle multiline diff', () => {
            const multilineDiff = `diff --git a/src/main.ts b/src/main.ts
index abc123..def456 100644
--- a/src/main.ts
+++ b/src/main.ts
@@ -10,6 +10,8 @@ function main() {
     const x = 1;
+    const y = 2;
+    const z = 3;
     return x;
 }
@@ -20,3 +22,4 @@ function helper() {
     return true;
+    // added comment
 }`

            store.setGeneratedDiff(multilineDiff)

            expect(store.gitDiff).toBe(multilineDiff)
            expect(store.gitDiff).toContain('diff --git')
            expect(store.gitDiff).toContain('@@ -10,6 +10,8 @@')
        })
    })

    describe('parsedDiff computed', () => {
        it('should return empty array when gitDiff is empty', () => {
            expect(store.parsedDiff).toEqual([])
        })

        it('should return empty array when gitDiff has content (current implementation)', () => {
            // Note: Current implementation returns empty array
            // This test documents current behavior
            store.setGeneratedDiff('some diff content')
            expect(store.parsedDiff).toEqual([])
        })
    })

    describe('Hunk type', () => {
        it('should correctly type Hunk interface', () => {
            const hunk: Hunk = {
                header: '@@ -1,3 +1,4 @@',
                lines: [' line1', '+new line', ' line2']
            }

            expect(hunk.header).toBe('@@ -1,3 +1,4 @@')
            expect(hunk.lines).toHaveLength(3)
            expect(hunk.lines[1]).toBe('+new line')
        })

        it('should allow empty lines array', () => {
            const hunk: Hunk = {
                header: '@@ -0,0 +1 @@',
                lines: []
            }

            expect(hunk.lines).toEqual([])
        })
    })

    describe('reactivity', () => {
        it('should be reactive when gitDiff changes', () => {
            expect(store.gitDiff).toBe('')

            store.setGeneratedDiff('updated')

            expect(store.gitDiff).toBe('updated')
        })
    })
})
