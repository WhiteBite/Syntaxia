/**
 * Smart context naming utilities
 * Generates human-readable names based on file patterns
 */

const TYPE_LABELS: Record<string, string> = {
    'vue': 'Vue компоненты',
    'ts': 'TypeScript',
    'tsx': 'React TSX',
    'js': 'JavaScript',
    'jsx': 'React JSX',
    'go': 'Go код',
    'py': 'Python',
    'java': 'Java',
    'css': 'Стили',
    'scss': 'SCSS стили',
    'json': 'Конфигурация',
    'yaml': 'YAML конфиг',
    'yml': 'YAML конфиг',
    'md': 'Документация',
    'sql': 'SQL запросы',
    'html': 'HTML шаблоны',
    'test': 'Тесты',
    'spec': 'Тесты'
}

interface FileAnalysis {
    extensions: Map<string, number>
    folders: Map<string, number>
    fileNames: string[]
    hasTests: boolean
}

function analyzeFiles(files: string[]): FileAnalysis {
    const extensions = new Map<string, number>()
    const folders = new Map<string, number>()
    const fileNames: string[] = []

    for (const file of files) {
        const parts = file.split('/')
        const fileName = parts[parts.length - 1]
        fileNames.push(fileName)

        // Count extensions
        const ext = fileName.includes('.') ? fileName.split('.').pop()?.toLowerCase() : ''
        if (ext) {
            extensions.set(ext, (extensions.get(ext) || 0) + 1)
        }

        // Count top-level folders
        if (parts.length > 1) {
            const topFolder = parts[0]
            folders.set(topFolder, (folders.get(topFolder) || 0) + 1)
        }
    }

    const hasTests = fileNames.some(f =>
        f.includes('test') || f.includes('spec') || f.includes('Test')
    )

    return { extensions, folders, fileNames, hasTests }
}

/**
 * Generate a smart name for context based on file patterns
 */
export function generateSmartName(files: string[]): string {
    if (!files || files.length === 0) return 'Пустой контекст'

    const { extensions, folders, fileNames, hasTests } = analyzeFiles(files)

    // Sort by frequency
    const sortedExts = [...extensions.entries()].sort((a, b) => b[1] - a[1])
    const sortedFolders = [...folders.entries()].sort((a, b) => b[1] - a[1])

    // Single file - use its name
    if (files.length === 1) {
        const fileName = fileNames[0]
        return fileName.length > 30 ? fileName.substring(0, 27) + '...' : fileName
    }

    let name = ''

    // All files from one folder
    if (sortedFolders.length === 1 && sortedFolders[0][1] === files.length) {
        name = sortedFolders[0][0]
    }
    // Dominant file type (>60%)
    else if (sortedExts.length > 0 && sortedExts[0][1] / files.length > 0.6) {
        const ext = sortedExts[0][0]
        name = TYPE_LABELS[ext] || `.${ext} файлы`
    }
    // Clear folder leader (>50%)
    else if (sortedFolders.length > 0 && sortedFolders[0][1] / files.length > 0.5) {
        name = sortedFolders[0][0]
    }
    // Mixed context
    else {
        const topTypes = sortedExts.slice(0, 2).map(([ext]) => TYPE_LABELS[ext] || ext)
        name = topTypes.join(' + ') || 'Смешанный'
    }

    // Add test label
    if (hasTests && !name.toLowerCase().includes('тест')) {
        name += ' (с тестами)'
    }

    // Add file count
    name += ` [${files.length}]`

    return name
}
