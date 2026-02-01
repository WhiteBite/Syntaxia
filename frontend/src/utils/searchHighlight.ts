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
 * @param key - The key to match against (default: 'name')
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

    // Find matches for the specified key (file name or path)
    const match = matches.find((m) => m.key === key)
    if (!match || !match.indices || match.indices.length === 0) {
        return [{ text, isMatch: false }]
    }

    const segments: TextSegment[] = []
    const indices = match.indices

    // Sort indices by start position to handle overlapping matches
    const sortedIndices = [...indices].sort((a, b) => a[0] - b[0])

    let lastIndex = 0

    for (const [start, end] of sortedIndices) {
        // Skip if this match overlaps with previous one
        if (start < lastIndex) {
            continue
        }

        // Add non-matched text before this match
        if (start > lastIndex) {
            segments.push({
                text: text.substring(lastIndex, start),
                isMatch: false,
            })
        }

        // Add matched text (end is inclusive in Fuse.js indices)
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
 * Highlights all occurrences of the query string (case-insensitive)
 * @param text - The text to highlight
 * @param query - The search query
 * @returns Array of text segments with match flags
 */
export function highlightSimple(text: string, query: string): TextSegment[] {
    if (!query || !query.trim()) {
        return [{ text, isMatch: false }]
    }

    const segments: TextSegment[] = []
    const lowerText = text.toLowerCase()
    const lowerQuery = query.toLowerCase().trim()

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
            text: text.substring(index, index + lowerQuery.length),
            isMatch: true,
        })

        lastIndex = index + lowerQuery.length
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

/**
 * Fuzzy character-by-character highlighting
 * Highlights individual characters that match the query in sequence
 * Useful for fuzzy search where characters don't need to be consecutive
 * @param text - The text to highlight
 * @param query - The search query
 * @returns Array of text segments with match flags
 */
export function highlightFuzzy(text: string, query: string): TextSegment[] {
    if (!query || !query.trim()) {
        return [{ text, isMatch: false }]
    }

    const segments: TextSegment[] = []
    const lowerText = text.toLowerCase()
    const lowerQuery = query.toLowerCase().trim()

    let textIndex = 0
    let queryIndex = 0
    let currentSegment = ''
    let isCurrentMatch = false

    while (textIndex < text.length) {
        const char = text[textIndex]
        const lowerChar = lowerText[textIndex]

        // Check if current character matches the next query character
        const isMatch = queryIndex < lowerQuery.length && lowerChar === lowerQuery[queryIndex]

        if (isMatch) {
            // Character matches query
            if (!isCurrentMatch && currentSegment) {
                // Save previous non-match segment
                segments.push({ text: currentSegment, isMatch: false })
                currentSegment = ''
            }
            currentSegment += char
            isCurrentMatch = true
            queryIndex++
        } else {
            // Character doesn't match
            if (isCurrentMatch && currentSegment) {
                // Save previous match segment
                segments.push({ text: currentSegment, isMatch: true })
                currentSegment = ''
            }
            currentSegment += char
            isCurrentMatch = false
        }

        textIndex++
    }

    // Add final segment
    if (currentSegment) {
        segments.push({ text: currentSegment, isMatch: isCurrentMatch })
    }

    return segments.length > 0 ? segments : [{ text, isMatch: false }]
}
