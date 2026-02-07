# Git Operations Caching & Optimization Analysis

## Executive Summary

**Current State:** Frontend has basic caching (5-min TTL), backend has NO caching
**Optimization Potential:** 70-85% reduction in Git command execution time
**Priority:** HIGH - Git operations are called frequently in UI

---

## 1. Current Git Operations

### Backend Git Commands (infrastructure/git/repository.go)

| Operation | Git Command | Frequency | Cacheable |
|-----------|-------------|-----------|-----------|
| `GetBranches` | `git branch -a` | Medium | ✅ Yes (2-5 min) |
| `GetCurrentBranch` | `git rev-parse --abbrev-ref HEAD` | High | ✅ Yes (1-2 min) |
| `GetCommitHistory` | `git log --pretty=format:...` | Medium | ✅ Yes (1-2 min) |
| `ListFilesAtRef` | `git ls-tree -r --name-only {ref}` | High | ✅ Yes (5-10 min) |
| `GetFileAtRef` | `git show {ref}:{file}` | Very High | ✅ Yes (10-15 min) |
| `GetUncommittedFiles` | `git status --porcelain` | Very High | ⚠️ Partial (30 sec) |
| `IsGitRepository` | `git rev-parse --git-dir` | Low | ✅ Yes (session) |
| `GetRichCommitHistory` | `git log --name-status` | Low | ✅ Yes (2-5 min) |

### Remote Operations (GitHub/GitLab API)

| Operation | API Call | Frequency | Cacheable |
|-----------|----------|-----------|-----------|
| `GitHubGetBranches` | GET /repos/{owner}/{repo}/branches | Low | ✅ Yes (5-10 min) |
| `GitHubListFiles` | GET /repos/{owner}/{repo}/git/trees/{ref}?recursive=1 | Medium | ✅ Yes (10-15 min) |
| `GitHubGetFileContent` | GET raw.githubusercontent.com/{owner}/{repo}/{ref}/{path} | High | ✅ Yes (15-30 min) |
| `GitLabGetBranches` | GET /projects/{id}/repository/branches | Low | ✅ Yes (5-10 min) |
| `GitLabListFiles` | GET /projects/{id}/repository/tree?recursive=true | Medium | ✅ Yes (10-15 min) |
| `GitLabGetFileContent` | GET /projects/{id}/repository/files/{path}/raw | High | ✅ Yes (15-30 min) |

---

## 2. Call Frequency Analysis

### Frontend Call Patterns

**GitLocalPanel.vue:**
- `loadGitInfo()` - Called on component mount + when project changes
  - Calls: `isGitRepository`, `getCurrentBranch`, `getBranches` (3 Git commands)
- `loadCommits()` - Called when switching to "Commits" tab
  - Calls: `getCommitHistory` (1 Git command)
- `loadFilesAtRef(ref)` - Called when selecting branch/commit
  - Calls: `listFilesAtRef` (1 Git command)

**useGitSource.ts:**
- `watch(() => projectStore.currentPath, checkGitRepo, { immediate: true })`
  - Triggers on EVERY project change
- File selection triggers `listFilesAtRef` for EACH ref selection

**Estimated Call Frequency (per session):**
- `isGitRepository`: 5-10 times
- `getCurrentBranch`: 10-20 times
- `getBranches`: 5-10 times
- `getCommitHistory`: 2-5 times
- `listFilesAtRef`: 10-30 times (depends on user interaction)
- `getFileAtRef`: 50-200 times (when building context from git ref)

### Performance Impact

**Without Caching:**
- Each `git` command: 50-200ms (local repo)
- Each GitHub API call: 200-1000ms (network latency)
- Building context from 50 files at ref: 50 × 100ms = **5 seconds**

**With Backend Caching:**
- Cache hit: <1ms (memory lookup)
- Building context from 50 files at ref: 50 × 1ms = **50ms**
- **100x faster** for cached operations

---

## 3. Current Caching Implementation

### Frontend Cache (useGitCache.ts)

**Implementation:**
```typescript
const cache = new Map<string, CacheEntry<unknown>>()
const DEFAULT_TTL = 5 * 60 * 1000 // 5 minutes
```

**Cache Keys:**
- `git:branches:{projectPath}`
- `git:commits:{projectPath}:{ref}`
- `git:files:{source}:{ref}`
- `git:content:{source}:{filePath}:{ref}`

**TTL Configuration:**
- Branches: 2 min
- Commits: 1 min
- Files: 5 min
- File content: 10 min

**Limitations:**
1. ❌ Cache is per-tab (not shared across windows)
2. ❌ Cache is lost on app restart
3. ❌ No cache invalidation on file changes
4. ❌ Backend still executes Git commands on every request

### Backend Cache

**Current State:** ❌ **NO CACHING**

Every API call executes `exec.Command("git", ...)` directly.

---

## 4. Performance Measurements

### Estimated Git Command Times (Local Repository)

| Command | Cold (no cache) | Warm (OS cache) | Backend Cache |
|---------|-----------------|-----------------|---------------|
| `git branch -a` | 80-150ms | 40-80ms | <1ms |
| `git rev-parse --abbrev-ref HEAD` | 30-60ms | 15-30ms | <1ms |
| `git log -n50` | 100-200ms | 50-100ms | <1ms |
| `git ls-tree -r {ref}` | 150-300ms | 80-150ms | <1ms |
| `git show {ref}:{file}` | 80-150ms | 40-80ms | <1ms |
| `git status --porcelain` | 100-500ms | 50-200ms | N/A (must be fresh) |

### Real-World Scenario: Building Context from Git Ref

**Task:** Build context from 50 files at specific commit

**Without Backend Cache:**
- 50 × `git show {ref}:{file}` = 50 × 100ms = **5,000ms (5 seconds)**

**With Backend Cache (after first load):**
- 50 × cache lookup = 50 × 0.5ms = **25ms**
- **200x faster**

**With File Watcher Invalidation:**
- Cache remains valid until files actually change
- 99% cache hit rate for read-only operations

---

## 5. Caching Opportunities

### High Priority (Immediate Impact)

#### 1. **ListFilesAtRef** - CRITICAL
- **Why:** Called frequently when browsing git history
- **Cache Duration:** 10-15 minutes
- **Invalidation:** On branch switch, new commits
- **Expected Impact:** 90% reduction in calls

#### 2. **GetFileAtRef** - CRITICAL
- **Why:** Called 50-200 times when building context
- **Cache Duration:** 15-30 minutes (immutable for specific ref)
- **Invalidation:** Never (git objects are immutable)
- **Expected Impact:** 95% reduction in calls

#### 3. **GetBranches** - HIGH
- **Why:** Called on every project switch
- **Cache Duration:** 5 minutes
- **Invalidation:** On `git fetch`, `git checkout -b`
- **Expected Impact:** 80% reduction in calls

#### 4. **GetCurrentBranch** - HIGH
- **Why:** Called frequently in UI
- **Cache Duration:** 2 minutes
- **Invalidation:** On `git checkout`
- **Expected Impact:** 85% reduction in calls

### Medium Priority

#### 5. **GetCommitHistory** - MEDIUM
- **Why:** Called when opening commits tab
- **Cache Duration:** 2-5 minutes
- **Invalidation:** On new commits
- **Expected Impact:** 70% reduction in calls

#### 6. **GitHub/GitLab API Calls** - MEDIUM
- **Why:** Network latency is expensive
- **Cache Duration:** 10-30 minutes
- **Invalidation:** Manual refresh button
- **Expected Impact:** 90% reduction in API calls

### Low Priority (Special Cases)

#### 7. **GetUncommittedFiles** - LOW
- **Why:** Must be fresh for accurate status
- **Cache Duration:** 10-30 seconds (very short)
- **Invalidation:** On file watcher events
- **Expected Impact:** 50% reduction (only for rapid repeated calls)

---

## 6. Cache Invalidation Strategy

### File Watcher Integration

**Existing Infrastructure:**
- ✅ `backend/infrastructure/fswatcher/watcher.go` already exists
- ✅ Uses `fsnotify` for file system events
- ✅ Connected to event bus

**Invalidation Triggers:**

| Event | Invalidate Cache |
|-------|------------------|
| File modified | `GetUncommittedFiles` |
| File created/deleted | `GetUncommittedFiles` |
| `.git/HEAD` changed | `GetCurrentBranch`, `GetBranches` |
| `.git/refs/heads/*` changed | `GetBranches`, `GetCommitHistory` |
| `.git/index` changed | `GetUncommittedFiles` |
| `git fetch` detected | Remote branches cache |

### Cache Key Strategy

```go
// Immutable caches (never invalidate)
"git:file:{projectPath}:{ref}:{filePath}"  // Git objects are immutable
"git:files:{projectPath}:{ref}"            // Tree at ref is immutable

// Mutable caches (invalidate on events)
"git:branches:{projectPath}"               // Invalidate on branch changes
"git:current-branch:{projectPath}"         // Invalidate on checkout
"git:status:{projectPath}"                 // Invalidate on file changes
"git:commits:{projectPath}:{branch}"       // Invalidate on new commits
```

### Manual Invalidation

**User Actions:**
- "Refresh" button in Git panel → Clear all caches for project
- Project switch → Clear caches for old project
- App restart → Persist cache to disk (optional)

---

## 7. Implementation Plan

### Phase 1: Backend Cache Infrastructure (Priority: HIGH)

**Goal:** Add caching layer to backend Git operations

**Tasks:**
1. Create `backend/infrastructure/git/cache.go`
   - In-memory cache with TTL
   - Thread-safe (sync.RWMutex)
   - LRU eviction policy
   - Cache statistics

2. Wrap Git operations with cache layer
   - `CachedGetBranches(projectPath string) ([]string, error)`
   - `CachedGetCurrentBranch(projectPath string) (string, error)`
   - `CachedListFilesAtRef(projectPath, ref string) ([]string, error)`
   - `CachedGetFileAtRef(projectPath, filePath, ref string) (string, error)`

3. Add cache configuration
   - TTL per operation type
   - Max cache size (MB)
   - Enable/disable per operation

**Estimated Time:** 4-6 hours

**Expected Impact:** 70-85% reduction in Git command execution

---

### Phase 2: File Watcher Integration (Priority: MEDIUM)

**Goal:** Automatic cache invalidation on file changes

**Tasks:**
1. Extend `fswatcher.Watcher` to watch `.git/` directory
   - Monitor `.git/HEAD` for branch changes
   - Monitor `.git/refs/heads/` for branch updates
   - Monitor `.git/index` for staging changes

2. Emit cache invalidation events
   - `git:branch-changed` → Invalidate branch caches
   - `git:files-changed` → Invalidate status cache
   - `git:commit-added` → Invalidate commit history

3. Connect cache to event bus
   - Subscribe to invalidation events
   - Clear relevant cache entries

**Estimated Time:** 3-4 hours

**Expected Impact:** 95% cache hit rate (vs 70% without invalidation)

---

### Phase 3: Remote API Caching (Priority: MEDIUM)

**Goal:** Cache GitHub/GitLab API responses

**Tasks:**
1. Add caching to `github_api.go` and `gitlab_api.go`
   - Cache tree responses (immutable)
   - Cache file content (immutable)
   - Cache branch lists (5-10 min TTL)

2. Add manual refresh mechanism
   - "Refresh" button in remote panel
   - Clear remote cache on demand

**Estimated Time:** 2-3 hours

**Expected Impact:** 90% reduction in API calls, faster remote browsing

---

### Phase 4: Persistent Cache (Priority: LOW)

**Goal:** Persist cache to disk for faster app startup

**Tasks:**
1. Serialize cache to disk on app shutdown
2. Load cache from disk on app startup
3. Validate cache entries (check if project still exists)

**Estimated Time:** 2-3 hours

**Expected Impact:** Instant Git info on app restart

---

### Phase 5: Cache Statistics & Monitoring (Priority: LOW)

**Goal:** Monitor cache effectiveness

**Tasks:**
1. Add cache hit/miss counters
2. Add cache size monitoring
3. Add cache eviction statistics
4. Expose metrics in UI (optional)

**Estimated Time:** 1-2 hours

**Expected Impact:** Visibility into cache performance

---

## 8. Expected Impact

### Performance Improvements

| Scenario | Before | After | Improvement |
|----------|--------|-------|-------------|
| Load Git info (3 commands) | 200-400ms | 10-20ms | **20x faster** |
| Load commit history | 100-200ms | 1-5ms | **50x faster** |
| List files at ref | 150-300ms | 1-5ms | **100x faster** |
| Build context (50 files) | 5,000ms | 50ms | **100x faster** |
| Browse remote repo | 2,000-5,000ms | 100-200ms | **20x faster** |

### User Experience

**Before:**
- ❌ Noticeable lag when switching branches
- ❌ Slow context building from git refs
- ❌ Sluggish remote repository browsing
- ❌ Repeated Git commands on every action

**After:**
- ✅ Instant branch switching
- ✅ Fast context building (100x faster)
- ✅ Smooth remote repository browsing
- ✅ Minimal Git command execution

### Resource Usage

**Memory:**
- Estimated cache size: 10-50 MB (for typical project)
- LRU eviction prevents unbounded growth
- Configurable max cache size

**CPU:**
- Reduced CPU usage from fewer `exec.Command` calls
- Minimal overhead for cache lookups (<1ms)

**Disk I/O:**
- Reduced disk I/O from fewer Git operations
- Optional persistent cache adds minimal I/O

---

## 9. Implementation Code Sketch

### Backend Cache Layer

```go
// backend/infrastructure/git/cache.go
package git

import (
    "sync"
    "time"
)

type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

type GitCache struct {
    mu    sync.RWMutex
    cache map[string]*CacheEntry
    stats CacheStats
}

type CacheStats struct {
    Hits   int64
    Misses int64
    Size   int64
}

func NewGitCache() *GitCache {
    return &GitCache{
        cache: make(map[string]*CacheEntry),
    }
}

func (c *GitCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    entry, exists := c.cache[key]
    if !exists || time.Now().After(entry.ExpiresAt) {
        c.stats.Misses++
        return nil, false
    }
    
    c.stats.Hits++
    return entry.Data, true
}

func (c *GitCache) Set(key string, data interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[key] = &CacheEntry{
        Data:      data,
        ExpiresAt: time.Now().Add(ttl),
    }
}

func (c *GitCache) Invalidate(pattern string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    for key := range c.cache {
        if strings.Contains(key, pattern) {
            delete(c.cache, key)
        }
    }
}
```

### Cached Repository Methods

```go
// backend/infrastructure/git/repository.go
type Repository struct {
    log   domain.Logger
    cache *GitCache
}

func (r *Repository) GetBranches(projectRoot string) ([]string, error) {
    cacheKey := fmt.Sprintf("branches:%s", projectRoot)
    
    // Check cache
    if cached, ok := r.cache.Get(cacheKey); ok {
        return cached.([]string), nil
    }
    
    // Execute git command
    branches, err := r.executeBranchesCommand(projectRoot)
    if err != nil {
        return nil, err
    }
    
    // Cache result (5 min TTL)
    r.cache.Set(cacheKey, branches, 5*time.Minute)
    
    return branches, nil
}

func (r *Repository) GetFileAtRef(projectPath, filePath, ref string) (string, error) {
    // Git objects are immutable - cache forever
    cacheKey := fmt.Sprintf("file:%s:%s:%s", projectPath, ref, filePath)
    
    if cached, ok := r.cache.Get(cacheKey); ok {
        return cached.(string), nil
    }
    
    content, err := r.executeShowCommand(projectPath, filePath, ref)
    if err != nil {
        return "", err
    }
    
    // Cache with long TTL (30 min) - git objects don't change
    r.cache.Set(cacheKey, content, 30*time.Minute)
    
    return content, nil
}
```

### File Watcher Integration

```go
// backend/infrastructure/fswatcher/watcher.go
func (w *Watcher) watchGitDirectory(projectPath string) {
    gitDir := filepath.Join(projectPath, ".git")
    
    // Watch .git/HEAD for branch changes
    w.fsWatcher.Add(filepath.Join(gitDir, "HEAD"))
    
    // Watch .git/refs/heads for branch updates
    w.fsWatcher.Add(filepath.Join(gitDir, "refs", "heads"))
    
    // Watch .git/index for staging changes
    w.fsWatcher.Add(filepath.Join(gitDir, "index"))
}

func (w *Watcher) handleGitEvent(event fsnotify.Event) {
    switch {
    case strings.Contains(event.Name, ".git/HEAD"):
        w.bus.Emit("git:branch-changed", event.Name)
    case strings.Contains(event.Name, ".git/refs/heads"):
        w.bus.Emit("git:refs-changed", event.Name)
    case strings.Contains(event.Name, ".git/index"):
        w.bus.Emit("git:index-changed", event.Name)
    }
}
```

---

## 10. Risks & Mitigation

### Risk 1: Stale Cache Data

**Risk:** Cache shows outdated information after external git operations

**Mitigation:**
- File watcher detects `.git/` changes
- Manual "Refresh" button in UI
- Short TTL for mutable data (1-5 min)
- Long TTL for immutable data (git objects)

### Risk 2: Memory Usage

**Risk:** Cache grows unbounded and consumes too much memory

**Mitigation:**
- LRU eviction policy
- Configurable max cache size (default: 100 MB)
- Periodic cleanup of expired entries
- Cache statistics monitoring

### Risk 3: Cache Invalidation Bugs

**Risk:** Cache not invalidated when it should be

**Mitigation:**
- Comprehensive test coverage
- Conservative TTL values
- Manual refresh option
- Cache statistics to detect issues

### Risk 4: Concurrency Issues

**Risk:** Race conditions in cache access

**Mitigation:**
- Thread-safe cache with `sync.RWMutex`
- Atomic cache operations
- Proper locking in all cache methods

---

## 11. Success Metrics

### Performance Metrics

- [ ] Git command execution time reduced by 70-85%
- [ ] Context building time reduced by 90-95%
- [ ] Cache hit rate > 80% after warmup
- [ ] Memory usage < 100 MB for cache

### User Experience Metrics

- [ ] Git panel loads instantly (<50ms)
- [ ] Branch switching feels instant
- [ ] Context building from git ref is fast (<500ms for 50 files)
- [ ] No noticeable lag in UI

### Code Quality Metrics

- [ ] Cache code has >80% test coverage
- [ ] No memory leaks detected
- [ ] No race conditions detected
- [ ] Cache statistics available for monitoring

---

## 12. Conclusion

**Current State:**
- Frontend has basic caching (5-min TTL, per-tab)
- Backend has NO caching (executes git commands on every request)
- Git operations are called frequently (10-200 times per session)

**Optimization Potential:**
- **70-85% reduction** in Git command execution time
- **100x faster** context building from git refs
- **20x faster** Git panel loading
- **90% reduction** in GitHub/GitLab API calls

**Recommended Approach:**
1. **Phase 1 (HIGH):** Implement backend cache infrastructure (4-6 hours)
2. **Phase 2 (MEDIUM):** Integrate file watcher for auto-invalidation (3-4 hours)
3. **Phase 3 (MEDIUM):** Cache remote API calls (2-3 hours)
4. **Phase 4 (LOW):** Add persistent cache (2-3 hours)
5. **Phase 5 (LOW):** Add cache monitoring (1-2 hours)

**Total Estimated Time:** 12-18 hours

**Expected Impact:** Massive performance improvement for Git operations, especially when building context from git history. Users will experience instant Git panel loading and fast context building.

**Priority:** **HIGH** - Git operations are a bottleneck in current implementation.
