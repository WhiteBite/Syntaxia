/**
 * highlightMatches - Helper for highlighting search matches in text
 * Used by VirtualTreeRow to highlight matched characters in file names
 */

import type { FuseResultMatch } from 'fuse.js'

export interface TextSegment {
    text: string
    isMatch: boolean
}

/**
 * Split text into segments based on match indices from Fuse.js
 * @param text - The text to highlight
 * @param matches - Match information from Fuse.js (optional)
 * @param key - The key to match against (e.g., 'name', 'path')
 * @returns Array of text segments with match flags
 */
export function highlightMatches(
    text: string,
    matches?: ReadonlyArray<FuseResultMatch>,
    key: string = 'name'
): TextSegment[] {
    if (!matches || matches.length === 0) {
        return [{ text, isMatch: false }]
    }

    // Find the match for the specified key
    const match = matches.find(m => m.key === key)
    if (!match || !match.indices || match.indices.length === 0) {
        return [{ text, isMatch: false }]
    }

    const segments: TextSegment[] = []
    let lastIndex = 0

    // Sort indices by start position
    const sortedIndices = [...match.indices].sort((a, b) => a[0] - b[0])

    for (const [start, end] of sortedIndices) {
        // Add unmatched text before this match
        if (start > lastIndex) {
            segments.push({
                text: text.substring(lastIndex, start),
                isMatch: false
            })
        }

        // Add matched text
        segments.push({
            text: text.substring(start, end + 1),
            isMatch: true
        })

        lastIndex = end + 1
    }

    // Add remaining unmatched text
    if (lastIndex < text.length) {
        segments.push({
            text: text.substring(lastIndex),
            isMatch: false
        })
    }

    return segments
}
