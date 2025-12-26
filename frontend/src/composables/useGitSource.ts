import { gitApi } from '@/features/git/api/git.api'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, ref, watch } from 'vue'

// Re-export CommitInfo from api.service for consistency
import { type CommitInfo } from '@/services/api.service'
export type { CommitInfo }

export interface RecentRepo {
    url: string
    name: string
    isGitHub: boolean
    lastUsed: number
}

export interface RemoteBranch {
    name: string
    commit: { sha: string }
}

// Constants
const RECENT_REPOS_KEY = 'git-recent-repos'
const MAX_RECENT_REPOS = 10
const SELECTION_KEY = 'git-file-selection'

// File type filters
export const fileTypeFilters = [
    { id: 'code', label: 'Code', icon: '💻', extensions: ['.ts', '.js', '.tsx', '.jsx', '.vue', '.go', '.py', '.java', '.rs', '.cpp', '.c', '.h', '.cs', '.rb', '.php', '.swift', '.kt'] },
    { id: 'styles', label: 'Styles', icon: '🎨', extensions: ['.css', '.scss', '.sass', '.less', '.styl'] },
    { id: 'config', label: 'Config', icon: '⚙️', extensions: ['.json', '.yaml', '.yml', '.toml', '.xml', '.ini', '.env', '.config'] },
    { id: 'docs', label: 'Docs', icon: '📄', extensions: ['.md', '.txt', '.rst', '.adoc'] },
]

export function useGitSource() {
    const projectStore = useProjectStore()
    const uiStore = useUIStore()

    // Local git state
    const isGitRepo = ref(false)
    const currentBranch = ref('')
    const branches = ref<string[]>([])
    const commits = ref<CommitInfo[]>([])
    const commitsLoaded = ref(false)
    const selectedRef = ref<string | null>(null)
    const filesAtRef = ref<string[]>([])
    const selectedFiles = ref<Set<string>>(new Set())

    // Remote state
    const remoteUrl = ref('')
    const isCloning = ref(false)
    const clonedPath = ref<string | null>(null)
    const isGitHubRepo = ref(false)
    const isGitLabRepo = ref(false)
    const isLoadingRemote = ref(false)
    const remoteRepoLoaded = ref(false)
    const remoteBranches = ref<RemoteBranch[]>([])
    const remoteSelectedBranch = ref('')
    const remoteFiles = ref<string[]>([])
    const remoteSelectedFiles = ref<Set<string>>(new Set())

    // Loading state
    const isLoading = ref(false)
    const loadingMessage = ref('')
    const isBuilding = ref(false)

    // Search & Filter state
    const localSearchQuery = ref('')
    const remoteSearchQuery = ref('')
    const localActiveFilters = ref<Set<string>>(new Set())
    const remoteActiveFilters = ref<Set<string>>(new Set())

    // Recent repos
    const recentRepos = ref<RecentRepo[]>([])

    const projectPath = computed(() => projectStore.currentPath || '')

    // Filter helpers
    function getFilteredFiles(files: string[], searchQuery: string, activeFilters: Set<string>): string[] {
        let result = files

        if (searchQuery) {
            const query = searchQuery.toLowerCase()
            result = result.filter(f => f.toLowerCase().includes(query))
        }

        if (activeFilters.size > 0) {
            const activeExtensions = new Set<string>()
            activeFilters.forEach(filterId => {
                const filter = fileTypeFilters.find(f => f.id === filterId)
                if (filter) filter.extensions.forEach(ext => activeExtensions.add(ext))
            })
            result = result.filter(f => {
                const ext = '.' + f.split('.').pop()?.toLowerCase()
                return activeExtensions.has(ext)
            })
        }

        return result
    }

    const filteredLocalFiles = computed(() =>
        getFilteredFiles(filesAtRef.value || [], localSearchQuery.value, localActiveFilters.value)
    )

    const filteredRemoteFiles = computed(() =>
        getFilteredFiles(remoteFiles.value || [], remoteSearchQuery.value, remoteActiveFilters.value)
    )

    function getFilterCount(files: string[], filterId: string): number {
        const filter = fileTypeFilters.find(f => f.id === filterId)
        if (!filter || !files) return 0
        return files.filter(f => {
            const ext = '.' + f.split('.').pop()?.toLowerCase()
            return filter.extensions.includes(ext)
        }).length
    }

    function toggleFilter(activeFilters: Set<string>, filterId: string): Set<string> {
        const newFilters = new Set(activeFilters)
        if (newFilters.has(filterId)) {
            newFilters.delete(filterId)
        } else {
            newFilters.add(filterId)
        }
        return newFilters
    }

    // Recent repos management
    function loadRecentReposFromStorage() {
        try {
            const stored = localStorage.getItem(RECENT_REPOS_KEY)
            if (stored) {
                recentRepos.value = JSON.parse(stored)
            }
        } catch {
            recentRepos.value = []
        }
    }

    function saveRecentRepo(url: string, name: string, isGitHub: boolean) {
        const existing = recentRepos.value.findIndex(r => r.url === url)
        if (existing >= 0) {
            recentRepos.value.splice(existing, 1)
        }

        recentRepos.value.unshift({ url, name, isGitHub, lastUsed: Date.now() })

        if (recentRepos.value.length > MAX_RECENT_REPOS) {
            recentRepos.value = recentRepos.value.slice(0, MAX_RECENT_REPOS)
        }

        localStorage.setItem(RECENT_REPOS_KEY, JSON.stringify(recentRepos.value))
    }

    function clearRecentRepos() {
        recentRepos.value = []
        localStorage.removeItem(RECENT_REPOS_KEY)
    }

    // Selection persistence
    function saveSelectionToStorage(key: string, files: Set<string>) {
        try {
            const data = { files: Array.from(files), timestamp: Date.now() }
            localStorage.setItem(`${SELECTION_KEY}-${key}`, JSON.stringify(data))
        } catch {
            // Ignore storage errors
        }
    }

    function loadSelectionFromStorage(key: string): Set<string> {
        try {
            const stored = localStorage.getItem(`${SELECTION_KEY}-${key}`)
            if (stored) {
                const data = JSON.parse(stored)
                if (Date.now() - data.timestamp < 24 * 60 * 60 * 1000) {
                    return new Set(data.files)
                }
            }
        } catch {
            // Ignore storage errors
        }
        return new Set()
    }

    // Git operations
    async function loadGitInfo() {
        if (!projectPath.value) return

        isLoading.value = true
        loadingMessage.value = 'Loading git info...'

        try {
            const [isRepo, branch, branchList] = await Promise.all([
                gitApi.isGitRepository(projectPath.value),
                gitApi.getCurrentBranch(projectPath.value).catch(() => ''),
                gitApi.getBranches(projectPath.value).catch(() => [] as string[]),
            ])
            isGitRepo.value = isRepo
            currentBranch.value = branch
            branches.value = branchList
        } catch {
            isGitRepo.value = false
        } finally {
            isLoading.value = false
        }
    }

    async function loadCommits() {
        if (!projectPath.value || commitsLoaded.value) return

        isLoading.value = true
        loadingMessage.value = 'Loading commits...'

        try {
            const gitCommits = await gitApi.getCommits(projectPath.value, 50)
            commits.value = gitCommits.map(c => ({
                hash: c.hash,
                subject: c.message || '',
                author: c.author,
                date: c.date,
            }))
            commitsLoaded.value = true
        } catch {
            commits.value = []
        } finally {
            isLoading.value = false
        }
    }

    async function loadFilesAtRef(ref: string) {
        if (!projectPath.value) return

        isLoading.value = true
        loadingMessage.value = `Loading files at ${ref.slice(0, 7)}...`

        try {
            filesAtRef.value = await gitApi.listFilesAtRef(projectPath.value, ref)
            const savedSelection = loadSelectionFromStorage(`local-${projectPath.value}-${ref}`)
            if (savedSelection.size > 0) {
                selectedFiles.value = savedSelection
            }
        } catch {
            filesAtRef.value = []
            uiStore.addToast('Failed to load files', 'error')
        } finally {
            isLoading.value = false
        }
    }

    function selectRef(ref: string) {
        selectedRef.value = ref
        loadFilesAtRef(ref)
    }

    function clearSelectedRef() {
        selectedRef.value = null
        filesAtRef.value = []
        selectedFiles.value = new Set()
    }

    // File selection
    function toggleFileSelection(path: string) {
        const newSelection = new Set(selectedFiles.value)
        if (newSelection.has(path)) {
            newSelection.delete(path)
        } else {
            newSelection.add(path)
        }
        selectedFiles.value = newSelection

        if (selectedRef.value && projectPath.value) {
            saveSelectionToStorage(`local-${projectPath.value}-${selectedRef.value}`, newSelection)
        }
    }

    function selectAllFiles() {
        selectedFiles.value = new Set(filteredLocalFiles.value)
    }

    function clearFileSelection() {
        selectedFiles.value = new Set()
    }

    function selectFolderFiles(folderPath: string) {
        const newSelection = new Set(selectedFiles.value)
        filteredLocalFiles.value
            .filter(f => f.startsWith(folderPath + '/'))
            .forEach(f => newSelection.add(f))
        selectedFiles.value = newSelection
    }

    // Remote operations
    function checkRemoteType(url: string) {
        isGitHubRepo.value = url.includes('github.com')
        isGitLabRepo.value = url.includes('gitlab.com')
    }

    async function loadRemoteRepo() {
        if (!remoteUrl.value) return

        checkRemoteType(remoteUrl.value)
        isLoadingRemote.value = true

        try {
            const provider = isGitHubRepo.value ? 'github' : isGitLabRepo.value ? 'gitlab' : 'unknown'
            if (provider === 'unknown') {
                uiStore.addToast('Unsupported repository provider', 'error')
                return
            }

            const [branchList, defaultBranch] = await Promise.all([
                gitApi.getRemoteBranches(remoteUrl.value, provider),
                gitApi.getRemoteDefaultBranch(remoteUrl.value, provider),
            ])

            remoteBranches.value = branchList.map(b => ({ name: b.name, commit: { sha: b.commit.sha } }))
            remoteSelectedBranch.value = defaultBranch
            remoteRepoLoaded.value = true
            await loadRemoteFiles()

            const repoName = remoteUrl.value.split('/').slice(-2).join('/')
            saveRecentRepo(remoteUrl.value, repoName, isGitHubRepo.value)
        } catch {
            uiStore.addToast('Failed to load repository', 'error')
        } finally {
            isLoadingRemote.value = false
        }
    }

    async function loadRemoteFiles() {
        if (!remoteSelectedBranch.value) return

        isLoading.value = true
        loadingMessage.value = 'Loading files...'

        try {
            const provider = isGitHubRepo.value ? 'github' : isGitLabRepo.value ? 'gitlab' : 'unknown'
            if (provider !== 'unknown') {
                remoteFiles.value = await gitApi.listRemoteFiles(remoteUrl.value, remoteSelectedBranch.value, provider)
            }
        } catch {
            remoteFiles.value = []
        } finally {
            isLoading.value = false
        }
    }

    // Remote file selection
    function toggleRemoteFileSelection(path: string) {
        const newSelection = new Set(remoteSelectedFiles.value)
        if (newSelection.has(path)) {
            newSelection.delete(path)
        } else {
            newSelection.add(path)
        }
        remoteSelectedFiles.value = newSelection
    }

    function selectAllRemoteFiles() {
        remoteSelectedFiles.value = new Set(filteredRemoteFiles.value)
    }

    function clearRemoteFileSelection() {
        remoteSelectedFiles.value = new Set()
    }

    function selectRemoteFolderFiles(folderPath: string) {
        const newSelection = new Set(remoteSelectedFiles.value)
        filteredRemoteFiles.value
            .filter(f => f.startsWith(folderPath + '/'))
            .forEach(f => newSelection.add(f))
        remoteSelectedFiles.value = newSelection
    }

    // Watch for URL changes
    watch(remoteUrl, (url) => {
        if (url) {
            checkRemoteType(url)
        }
    })

    return {
        // Local state
        isGitRepo,
        currentBranch,
        branches,
        commits,
        commitsLoaded,
        selectedRef,
        filesAtRef,
        selectedFiles,
        filteredLocalFiles,
        localSearchQuery,
        localActiveFilters,

        // Remote state
        remoteUrl,
        isCloning,
        clonedPath,
        isGitHubRepo,
        isGitLabRepo,
        isLoadingRemote,
        remoteRepoLoaded,
        remoteBranches,
        remoteSelectedBranch,
        remoteFiles,
        remoteSelectedFiles,
        filteredRemoteFiles,
        remoteSearchQuery,
        remoteActiveFilters,

        // Loading state
        isLoading,
        loadingMessage,
        isBuilding,

        // Recent repos
        recentRepos,

        // Computed
        projectPath,

        // Methods
        loadGitInfo,
        loadCommits,
        selectRef,
        clearSelectedRef,
        toggleFileSelection,
        selectAllFiles,
        clearFileSelection,
        selectFolderFiles,
        loadRemoteRepo,
        loadRemoteFiles,
        toggleRemoteFileSelection,
        selectAllRemoteFiles,
        clearRemoteFileSelection,
        selectRemoteFolderFiles,
        loadRecentReposFromStorage,
        saveRecentRepo,
        clearRecentRepos,
        getFilterCount,
        toggleFilter,
        fileTypeFilters,
    }
}
