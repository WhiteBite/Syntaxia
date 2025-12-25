import type { FileNode } from '../model/file.store'

/**
 * Get compact path for nested folders (e.g., "src/main/java")
 */
export function getCompactPath(node: FileNode): string {
    if (!node.isDir || !node.children) return node.name

    const path: string[] = [node.name]
    let current = node

    // Traverse down while folder has only one subfolder and no files
    while (current.children && current.children.length === 1 && current.children[0].isDir) {
        const child = current.children[0]

        // Check if this child has any files (not just folders)
        const hasFiles = child.children?.some(n => !n.isDir) || false

        // Stop if we found files or if there are multiple children
        if (hasFiles) break

        path.push(child.name)
        current = child
    }

    return path.join('/')
}

/**
 * Check if path matches gitignore pattern
 */
export function isIgnoredByGitignore(path: string, gitignoreRules: string): boolean {
    if (!gitignoreRules) return false

    const rules = gitignoreRules.split('\n')
        .map(line => line.trim())
        .filter(line => line && !line.startsWith('#'))

    return rules.some(rule => {
        // Simple pattern matching (can be enhanced with proper gitignore library)
        const pattern = rule.replace(/\*/g, '.*').replace(/\?/g, '.')
        const regex = new RegExp(`^${pattern}$`)
        return regex.test(path) || path.includes(rule.replace(/\*/g, ''))
    })
}

/**
 * Check if path matches custom ignore rules
 */
export function isIgnoredByCustomRules(path: string, customRules: string): boolean {
    return isIgnoredByGitignore(path, customRules)
}

/**
 * Copy text to clipboard
 */
export async function copyToClipboard(text: string): Promise<void> {
    try {
        if (navigator.clipboard && window.isSecureContext) {
            await navigator.clipboard.writeText(text)
        } else {
            // Fallback for older browsers
            const textArea = document.createElement('textarea')
            textArea.value = text
            textArea.style.position = 'fixed'
            textArea.style.left = '-999999px'
            document.body.appendChild(textArea)
            textArea.select()
            document.execCommand('copy')
            document.body.removeChild(textArea)
        }
    } catch (err) {
        console.error('Failed to copy to clipboard:', err)
        throw err
    }
}

/**
 * Get relative path from base path
 */
export function getRelativePath(fullPath: string, basePath: string): string {
    if (!fullPath.startsWith(basePath)) return fullPath

    const relative = fullPath.slice(basePath.length)
    return relative.replace(/^[\\/]+/, '')
}

/**
 * Normalize path separators to forward slashes
 */
export function normalizePathSeparators(path: string): string {
    return path.replace(/\\/g, '/')
}

/**
 * Validate ignore pattern syntax
 */
export function validateIgnorePattern(pattern: string): { valid: boolean; error?: string } {
    if (!pattern.trim()) {
        return { valid: false, error: 'Pattern cannot be empty' }
    }

    // Check for invalid characters
    const invalidChars = /[<>:"|]/
    if (invalidChars.test(pattern)) {
        return { valid: false, error: 'Pattern contains invalid characters' }
    }

    // Check for valid glob patterns
    try {
        // Basic validation - can be enhanced
        const hasValidGlob = /^[a-zA-Z0-9_\-.*/\\]+$/.test(pattern)
        if (!hasValidGlob && !pattern.startsWith('#')) {
            return { valid: false, error: 'Invalid glob pattern' }
        }
    } catch (err) {
        return { valid: false, error: 'Invalid pattern syntax' }
    }

    return { valid: true }
}

/**
 * Parse ignore rules from string
 */
export function parseIgnoreRules(rules: string): string[] {
    return rules
        .split('\n')
        .map(line => line.trim())
        .filter(line => line && !line.startsWith('#'))
}

/**
 * Format file size
 */
export function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 B'

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

/**
 * Get file extension
 */
export function getFileExtension(filePath: string): string {
    const filename = getFileName(filePath)
    const lastDot = filename.lastIndexOf('.')
    if (lastDot === -1 || lastDot === 0) return ''
    return filename.slice(lastDot)
}

/**
 * Get file name from path
 */
export function getFileName(filePath: string): string {
    const normalized = normalizePathSeparators(filePath)
    const parts = normalized.split('/')
    return parts[parts.length - 1] || ''
}

/**
 * Check if path is a directory
 */
export function isDirectory(path: string): boolean {
    return !path.includes('.') || path.endsWith('/')
}
