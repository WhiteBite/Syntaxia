import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { GitBranch, GitCommit, GitProvider, RecentRepo } from './types'

const RECENT_REPOS_KEY = 'git-recent-repos'
const SELECTION_KEY = 'git-file-selection'
const MAX_RECENT_REPOS = 10
const SELECTION_TTL = 24 * 60 * 60 * 1000 // 24 hours

export const useGitStore = defineStore('git', () => {
    // Local repo state
    const isGitRepo = ref(false)
    const currentBranch = ref('')
    const branches = ref<string[]>([])
    const commits = ref<GitCommit[]>([])
    const selectedRef = ref<string | null>(null)
    const filesAtRef = ref<string[]>([])
    const selectedFiles = ref<Set<string>>(new Set())

    // Remote repo state
    const remoteUrl = ref('')
    const remoteProvider = ref<GitProvider>('unknown')
    const remoteLoading = ref(false)
    const remoteLoaded = ref(false)
    const remoteBranches = ref<GitBranch[]>([])
    const remoteSelectedBranch = ref('')
    const remoteFiles = ref<string[]>([])
    const remoteSelectedFiles = ref<Set<string>>(new Set())

    // UI state
    const sourceType = ref<'local' | 'remote'>('local')
    const refType = ref<'branches' | 'commits'>('branches')
    const isLoading = ref(false)
    const loadingMessage = ref('')

    // Lazy loading flags
    const branchesLoaded = ref(false)
    const commitsLoaded = ref(false)

    // Recent repos
    const recentRepos = ref<RecentRepo[]>([])

    // Computed
    const hasSelectedFiles = computed(() =>
        sourceType.value === 'local'
            ? selectedFiles.value.size > 0
            : remoteSelectedFiles.value.size > 0
    )

    const selectedFilesCount = computed(() =>
        sourceType.value === 'local'
            ? selectedFiles.value.size
            : remoteSelectedFiles.value.size
    )

    // Actions
    function setLocalRepo(isRepo: boolean, branch: string) {
        isGitRepo.value = isRepo
        currentBranch.value = branch
        branchesLoaded.value = false
        commitsLoaded.value = false
        branches.value = []
        commits.value = []
    }

    function setBranches(branchList: string[]) {
        branches.value = branchList
        branchesLoaded.value = true
    }

    function setCommits(commitList: GitCommit[]) {
        commits.value = commitList
        commitsLoaded.value = true
    }

    function selectRef(ref: string) {
        selectedRef.value = ref
        selectedFiles.value = new Set()
        filesAtRef.value = []
    }

    function clearSelectedRef() {
        selectedRef.value = null
        filesAtRef.value = []
        selectedFiles.value = new Set()
    }

    function setFilesAtRef(files: string[]) {
        filesAtRef.value = files
    }

    function toggleFileSelection(path: string) {
        const files = sourceType.value === 'local' ? selectedFiles : remoteSelectedFiles
        if (files.value.has(path)) {
            files.value.delete(path)
        } else {
            files.value.add(path)
        }
        files.value = new Set(files.value)
    }

    function selectFolderFiles(paths: string[]) {
        const files = sourceType.value === 'local' ? selectedFiles : remoteSelectedFiles
        const allSelected = paths.every(p => files.value.has(p))

        if (allSelected) {
            paths.forEach(p => files.value.delete(p))
        } else {
            paths.forEach(p => files.value.add(p))
        }
        files.value = new Set(files.value)
    }

    function selectAllFiles() {
        const source = sourceType.value === 'local' ? filesAtRef : remoteFiles
        const target = sourceType.value === 'local' ? selectedFiles : remoteSelectedFiles
        target.value = new Set(source.value)
    }

    function clearFileSelection() {
        const target = sourceType.value === 'local' ? selectedFiles : remoteSelectedFiles
        target.value = new Set()
    }

    // Remote repo actions
    function setRemoteRepo(url: string, provider: GitProvider) {
        remoteUrl.value = url
        remoteProvider.value = provider
        remoteLoaded.value = false
        remoteBranches.value = []
        remoteFiles.value = []
        remoteSelectedFiles.value = new Set()
    }

    function setRemoteBranches(branchList: GitBranch[], defaultBranch: string) {
        remoteBranches.value = branchList
        remoteSelectedBranch.value = defaultBranch
        remoteLoaded.value = true
    }

    function setRemoteFiles(files: string[]) {
        remoteFiles.value = files
    }

    // Recent repos persistence
    function loadRecentRepos() {
        try {
            const stored = localStorage.getItem(RECENT_REPOS_KEY)
            if (stored) {
                recentRepos.value = JSON.parse(stored)
            }
        } catch (e) {
            console.error('[GitStore] Failed to load recent repos:', e)
        }
    }

    function addRecentRepo(url: string, name: string, provider: GitProvider) {
        const existing = recentRepos.value.findIndex(r => r.url === url)
        if (existing >= 0) {
            recentRepos.value.splice(existing, 1)
        }

        recentRepos.value.unshift({
            url,
            name,
            provider,
            lastUsed: Date.now()
        })

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
    function saveSelection(key: string, files: Set<string>) {
        try {
            const data = { files: Array.from(files), timestamp: Date.now() }
            localStorage.setItem(`${SELECTION_KEY}-${key}`, JSON.stringify(data))
        } catch (e) {
            console.error('[GitStore] Failed to save selection:', e)
        }
    }

    function loadSelection(key: string): Set<string> {
        try {
            const stored = localStorage.getItem(`${SELECTION_KEY}-${key}`)
            if (stored) {
                const data = JSON.parse(stored)
                if (Date.now() - data.timestamp < SELECTION_TTL) {
                    return new Set(data.files)
                }
            }
        } catch (e) {
            console.error('[GitStore] Failed to load selection:', e)
        }
        return new Set()
    }

    function reset() {
        isGitRepo.value = false
        currentBranch.value = ''
        branches.value = []
        commits.value = []
        selectedRef.value = null
        filesAtRef.value = []
        selectedFiles.value = new Set()
        branchesLoaded.value = false
        commitsLoaded.value = false
    }

    return {
        // State
        isGitRepo,
        currentBranch,
        branches,
        commits,
        selectedRef,
        filesAtRef,
        selectedFiles,
        remoteUrl,
        remoteProvider,
        remoteLoading,
        remoteLoaded,
        remoteBranches,
        remoteSelectedBranch,
        remoteFiles,
        remoteSelectedFiles,
        sourceType,
        refType,
        isLoading,
        loadingMessage,
        branchesLoaded,
        commitsLoaded,
        recentRepos,

        // Computed
        hasSelectedFiles,
        selectedFilesCount,

        // Actions
        setLocalRepo,
        setBranches,
        setCommits,
        selectRef,
        clearSelectedRef,
        setFilesAtRef,
        toggleFileSelection,
        selectFolderFiles,
        selectAllFiles,
        clearFileSelection,
        setRemoteRepo,
        setRemoteBranches,
        setRemoteFiles,
        loadRecentRepos,
        addRecentRepo,
        clearRecentRepos,
        saveSelection,
        loadSelection,
        reset,
    }
})
