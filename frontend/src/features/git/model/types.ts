// Git Feature Types

export type GitProvider = 'github' | 'gitlab' | 'local' | 'unknown'

export interface GitBranch {
    name: string
    commit: {
        sha: string
    }
    isDefault?: boolean
    isCurrent?: boolean
}

export interface GitCommit {
    hash: string
    shortHash: string
    subject: string
    message?: string
    author: string
    date: string
}

export interface RecentRepo {
    url: string
    name: string
    provider: GitProvider
    lastUsed: number
}

export interface GitFileNode {
    path: string
    name: string
    isDir: boolean
    depth: number
}

export interface GitDiffFile {
    path: string
    status: 'added' | 'deleted' | 'modified' | 'renamed'
    additions?: number
    deletions?: number
}

export interface GitRepoState {
    isGitRepo: boolean
    currentBranch: string
    branches: string[]
    commits: GitCommit[]
    selectedRef: string | null
    filesAtRef: string[]
    selectedFiles: Set<string>
}

export interface RemoteRepoState {
    url: string
    provider: GitProvider
    isLoading: boolean
    isLoaded: boolean
    branches: GitBranch[]
    selectedBranch: string
    files: string[]
    selectedFiles: Set<string>
}

export interface FileFilter {
    id: string
    label: string
    extensions: string[]
}

export const FILE_TYPE_FILTERS: FileFilter[] = [
    {
        id: 'code',
        label: 'Code',
        extensions: ['.ts', '.js', '.tsx', '.jsx', '.vue', '.go', '.py', '.java', '.rs', '.cpp', '.c', '.h', '.cs', '.rb', '.php', '.swift', '.kt']
    },
    {
        id: 'styles',
        label: 'Styles',
        extensions: ['.css', '.scss', '.sass', '.less', '.styl']
    },
    {
        id: 'config',
        label: 'Config',
        extensions: ['.json', '.yaml', '.yml', '.toml', '.xml', '.ini', '.env', '.config']
    },
    {
        id: 'docs',
        label: 'Docs',
        extensions: ['.md', '.txt', '.rst', '.adoc']
    },
]
