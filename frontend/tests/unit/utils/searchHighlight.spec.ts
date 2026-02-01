/**
 * Tests for searchHighlight utility
 * Validates highlighting logic for search results
 */

import { highlightFuzzy, highlightMatches, highlightSimple } from '@/utils/searchHighlight'
import type { FuseResultMatch } from 'fuse.js'
import { describe, expect, it } from 'vitest'

describe('searchHighlight', () => {
    describe('highlightMatches (Fuse.js based)', () => {
        it('should return single non-match segment when no matches provided', () => {
            const result = highlightMatches('test.ts')

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })

        it('should return single non-match segment when matches array is empty', () => {
            const result = highlightMatches('test.ts', [])

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })

        it('should highlight single match at start', () => {
            const matches: FuseResultMatch[] = [{
                indices: [[0, 3]],
                value: 'test.ts',
                key: 'name',
                refIndex: 0
            }]

            const result = highlightMatches('test.ts', matches)

            expect(result).toEqual([
                { text: 'test', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should highlight single match in middle', () => {
            const matches: FuseResultMatch[] = [{
                indices: [[5, 8]],
                value: 'file-test.ts',
                key: 'name',
                refIndex: 0
            }]

            const result = highlightMatches('file-test.ts', matches)

            expect(result).toEqual([
                { text: 'file-', isMatch: false },
                { text: 'test', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should highlight single match at end', () => {
            const matches: FuseResultMatch[] = [{
                indices: [[5, 6]],
                value: 'test.ts',
                key: 'name',
                refIndex: 0
            }]

            const result = highlightMatches('test.ts', matches)

            expect(result).toEqual([
                { text: 'test.', isMatch: false },
                { text: 'ts', isMatch: true }
            ])
        })

        it('should highlight multiple non-overlapping matches', () => {
            const matches: FuseResultMatch[] = [{
                indices: [[0, 3], [5, 8]],
                value: 'test-file.ts',
                key: 'name',
                refIndex: 0
            }]

            const result = highlightMatches('test-file.ts', matches)

            expect(result).toEqual([
                { text: 'test', isMatch: true },
                { text: '-', isMatch: false },
                { text: 'file', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should handle overlapping matches by skipping overlaps', () => {
            const matches: FuseResultMatch[] = [{
                indices: [[0, 5], [3, 8]],
                value: 'testfile.ts',
                key: 'name',
                refIndex: 0
            }]

            const result = highlightMatches('testfile.ts', matches)

            expect(result).toEqual([
                { text: 'testfi', isMatch: true },
                { text: 'le.ts', isMatch: false }
            ])
        })

        it('should use specified key parameter', () => {
            const matches: FuseResultMatch[] = [
                {
                    indices: [[0, 3]],
                    value: 'test.ts',
                    key: 'name',
                    refIndex: 0
                },
                {
                    indices: [[4, 7]],
                    value: 'src/test.ts',
                    key: 'path',
                    refIndex: 0
                }
            ]

            const result = highlightMatches('src/test.ts', matches, 'path')

            expect(result).toEqual([
                { text: 'src/', isMatch: false },
                { text: 'test', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should return non-match when key not found', () => {
            const matches: FuseResultMatch[] = [{
                indices: [[0, 3]],
                value: 'test.ts',
                key: 'name',
                refIndex: 0
            }]

            const result = highlightMatches('test.ts', matches, 'path')

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })
    })

    describe('highlightSimple (string-based)', () => {
        it('should return single non-match segment when query is empty', () => {
            const result = highlightSimple('test.ts', '')

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })

        it('should return single non-match segment when query is whitespace', () => {
            const result = highlightSimple('test.ts', '   ')

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })

        it('should highlight exact match (case-insensitive)', () => {
            const result = highlightSimple('test.ts', 'test')

            expect(result).toEqual([
                { text: 'test', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should highlight match with different case', () => {
            const result = highlightSimple('TestFile.ts', 'test')

            expect(result).toEqual([
                { text: 'Test', isMatch: true },
                { text: 'File.ts', isMatch: false }
            ])
        })

        it('should highlight multiple occurrences', () => {
            const result = highlightSimple('test-test.ts', 'test')

            expect(result).toEqual([
                { text: 'test', isMatch: true },
                { text: '-', isMatch: false },
                { text: 'test', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should highlight partial match in middle', () => {
            const result = highlightSimple('myTestFile.ts', 'test')

            expect(result).toEqual([
                { text: 'my', isMatch: false },
                { text: 'Test', isMatch: true },
                { text: 'File.ts', isMatch: false }
            ])
        })

        it('should return non-match when query not found', () => {
            const result = highlightSimple('test.ts', 'xyz')

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })

        it('should trim query before matching', () => {
            const result = highlightSimple('test.ts', '  test  ')

            expect(result).toEqual([
                { text: 'test', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })
    })

    describe('highlightFuzzy (character-by-character)', () => {
        it('should return single non-match segment when query is empty', () => {
            const result = highlightFuzzy('test.ts', '')

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })

        it('should highlight consecutive matching characters', () => {
            const result = highlightFuzzy('test.ts', 'test')

            expect(result).toEqual([
                { text: 'test', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should highlight non-consecutive matching characters', () => {
            const result = highlightFuzzy('TestFile.ts', 'tf')

            expect(result).toEqual([
                { text: 'T', isMatch: true },
                { text: 'est', isMatch: false },
                { text: 'F', isMatch: true },
                { text: 'ile.ts', isMatch: false }
            ])
        })

        it('should highlight fuzzy match across word boundaries', () => {
            const result = highlightFuzzy('my-test-file.ts', 'mtf')

            expect(result).toEqual([
                { text: 'm', isMatch: true },
                { text: 'y-', isMatch: false },
                { text: 't', isMatch: true },
                { text: 'est-', isMatch: false },
                { text: 'f', isMatch: true },
                { text: 'ile.ts', isMatch: false }
            ])
        })

        it('should be case-insensitive', () => {
            const result = highlightFuzzy('TestFile.ts', 'TF')

            expect(result).toEqual([
                { text: 'T', isMatch: true },
                { text: 'est', isMatch: false },
                { text: 'F', isMatch: true },
                { text: 'ile.ts', isMatch: false }
            ])
        })

        it('should return non-match when no characters match', () => {
            const result = highlightFuzzy('test.ts', 'xyz')

            expect(result).toEqual([
                { text: 'test.ts', isMatch: false }
            ])
        })

        it('should handle partial fuzzy match', () => {
            const result = highlightFuzzy('component.vue', 'cmpv')

            expect(result).toEqual([
                { text: 'c', isMatch: true },
                { text: 'o', isMatch: false },
                { text: 'mp', isMatch: true },
                { text: 'onent.', isMatch: false },
                { text: 'v', isMatch: true },
                { text: 'ue', isMatch: false }
            ])
        })

        it('should trim query before matching', () => {
            const result = highlightFuzzy('test.ts', '  ts  ')

            expect(result).toEqual([
                { text: 't', isMatch: true },
                { text: 'e', isMatch: false },
                { text: 's', isMatch: true },
                { text: 't.ts', isMatch: false }
            ])
        })

        it('should handle query longer than text gracefully', () => {
            const result = highlightFuzzy('ab', 'abcdef')

            expect(result).toEqual([
                { text: 'ab', isMatch: true }
            ])
        })
    })

    describe('edge cases', () => {
        it('should handle empty text', () => {
            expect(highlightMatches('')).toEqual([{ text: '', isMatch: false }])
            expect(highlightSimple('', 'test')).toEqual([{ text: '', isMatch: false }])
            expect(highlightFuzzy('', 'test')).toEqual([{ text: '', isMatch: false }])
        })

        it('should handle single character text', () => {
            expect(highlightSimple('a', 'a')).toEqual([{ text: 'a', isMatch: true }])
            expect(highlightFuzzy('a', 'a')).toEqual([{ text: 'a', isMatch: true }])
        })

        it('should handle special characters', () => {
            const result = highlightSimple('test-file_v2.ts', 'file_v2')

            expect(result).toEqual([
                { text: 'test-', isMatch: false },
                { text: 'file_v2', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })

        it('should handle unicode characters', () => {
            const result = highlightSimple('тест.ts', 'тест')

            expect(result).toEqual([
                { text: 'тест', isMatch: true },
                { text: '.ts', isMatch: false }
            ])
        })
    })
})
