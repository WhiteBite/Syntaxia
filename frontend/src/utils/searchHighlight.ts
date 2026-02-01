/**
 * Search Highlight Utility
 * Splits text into segments for highlighting matched characters
 */

import type { FuseResultMatch } from 'fuse.js'

export interface TextSegment {
    text: string
    isMatch: boolean
}

/**
 * Highlight matches in text based on Fuse.js match indices
 * @param text - The text to highlight
 * @param matches - Match data from Fuse.js (optional)
 * @returns Array of text segments with match flags
 */
export function highlightMatches(
    text: string,
    matches?: ReadonlyArray<FuseResultMatch>
): TextSegment[] {
    if (!matches || matches.length === 0) {
        return [{ text, isMatch: false }]
    }

    // Find matches for the 'name' key (file name)
    const nameMatch = matches.find((m) => m.key === 'name')
    if (!nameMatch || !nameMatch.indices || nameMatch.indices.length === 0) {
        return [{ text, isMatch: false }]
    }

    const segments: TextSegment[] = []
    const indices = nameMatch.indices

    // Sort indices by start position
    const sortedIndices = [...indices].sort((a, b) => a[0] - b[0])

    let lastIndex = 0

    for (const [start, end] of sortedIndices) {
        // Add non-matched text before this match
        if (start > lastIndex) {
            segments.push({
                text: text.substring(lastIndex, start),
                isMatch: false,
            })
        }

        // Add matched text
        segments.push({
            text: text.substring(start, end + 1),
            isMatch: true,
        })

        lastIndex = end + 1
    }

    // Add remaining non-matched text
    if (lastIndex < text.length) {
        segments.push({
            text: text.substring(lastIndex),
            isMatch: false,
        })
    }

    return segments
}

/**
 * Simple string-based highlighting for non-fuzzy search
 * @param text - The text to highlight
 * @param query - The search query
 * @returns Array of text segments with match flags
 */
export function highlightSimple(text: string, query: string): TextSegment[] {
    if (!query) {
        return [{ text, isMatch: false }]
    }

    const segments: TextSegment[] = []
    const lowerText = text.toLowerCase()
    const lowerQuery = query.toLowerCase()

    let lastIndex = 0
    let index = lowerText.indexOf(lowerQuery, lastIndex)

    while (index !== -1) {
        // Add non-matched text before this match
        if (index > lastIndex) {
            segments.push({
                text: text.substring(lastIndex, index),
                isMatch: false,
            })
        }

        // Add matched text
        segments.push({
            text: text.substring(index, index + query.length),
            isMatch: true,
        })

        lastIndex = index + query.length
        index = lowerText.indexOf(lowerQuery, lastIndex)
    }

    // Add remaining non-matched text
    if (lastIndex < text.length) {
        segments.push({
            text: text.substring(lastIndex),
            isMatch: false,
        })
    }

    return segments.length > 0 ? segments : [{ text, isMatch: false }]
}
