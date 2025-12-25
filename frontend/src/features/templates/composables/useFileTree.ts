// File tree generation for templates

interface TreeNode {
    name: string
    isDir: boolean
    children: TreeNode[]
}

/**
 * Generate ASCII tree representation from file paths
 */
export function generateFileTree(filePaths: string[], projectName: string = ''): string {
    if (!filePaths || filePaths.length === 0) return ''

    // Build tree structure
    const root: TreeNode = { name: projectName || 'project', isDir: true, children: [] }

    for (const path of filePaths) {
        const parts = path.split(/[/\\]/).filter(Boolean)
        let current = root

        for (let i = 0; i < parts.length; i++) {
            const part = parts[i]
            const isLast = i === parts.length - 1

            let child = current.children.find(c => c.name === part)
            if (!child) {
                child = { name: part, isDir: !isLast, children: [] }
                current.children.push(child)
            }
            current = child
        }
    }

    // Sort children: directories first, then alphabetically
    sortTree(root)

    // Generate ASCII representation
    const lines: string[] = []
    lines.push(root.name + '/')
    renderTree(root.children, '', lines)

    return lines.join('\n')
}

function sortTree(node: TreeNode): void {
    node.children.sort((a, b) => {
        if (a.isDir !== b.isDir) return a.isDir ? -1 : 1
        return a.name.localeCompare(b.name)
    })
    for (const child of node.children) {
        sortTree(child)
    }
}

function renderTree(nodes: TreeNode[], prefix: string, lines: string[]): void {
    for (let i = 0; i < nodes.length; i++) {
        const node = nodes[i]
        const isLast = i === nodes.length - 1
        const connector = isLast ? '└── ' : '├── '
        const suffix = node.isDir ? '/' : ''

        lines.push(prefix + connector + node.name + suffix)

        if (node.children.length > 0) {
            const childPrefix = prefix + (isLast ? '    ' : '│   ')
            renderTree(node.children, childPrefix, lines)
        }
    }
}

/**
 * Detect programming languages from file extensions
 */
export function detectLanguages(filePaths: string[]): string[] {
    const extToLang: Record<string, string> = {
        '.ts': 'TypeScript',
        '.tsx': 'TypeScript',
        '.js': 'JavaScript',
        '.jsx': 'JavaScript',
        '.vue': 'Vue',
        '.go': 'Go',
        '.py': 'Python',
        '.java': 'Java',
        '.kt': 'Kotlin',
        '.rs': 'Rust',
        '.cpp': 'C++',
        '.c': 'C',
        '.cs': 'C#',
        '.rb': 'Ruby',
        '.php': 'PHP',
        '.swift': 'Swift',
        '.scala': 'Scala',
        '.html': 'HTML',
        '.css': 'CSS',
        '.scss': 'SCSS',
        '.sql': 'SQL',
        '.sh': 'Shell',
        '.yaml': 'YAML',
        '.yml': 'YAML',
        '.json': 'JSON',
        '.md': 'Markdown',
    }

    const languages = new Set<string>()

    for (const path of filePaths) {
        const ext = '.' + path.split('.').pop()?.toLowerCase()
        if (ext && extToLang[ext]) {
            languages.add(extToLang[ext])
        }
    }

    return Array.from(languages).sort()
}
