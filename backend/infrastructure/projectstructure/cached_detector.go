package projectstructure

import (
	"syntaxia/domain"
	"sync"
	"time"
)

const (
	// DefaultCacheTTL is the default time-to-live for cached results
	DefaultCacheTTL = 5 * time.Minute
)

// CachedDetector wraps Detector with caching capabilities
type CachedDetector struct {
	detector *Detector
	mu       sync.RWMutex

	// Cache storage
	structureCache    map[string]*cachedStructure
	architectureCache map[string]*cachedArchitecture
	frameworksCache   map[string]*cachedFrameworks
	conventionsCache  map[string]*cachedConventions

	// Configuration
	ttl time.Duration
}

type cachedStructure struct {
	result    *domain.ProjectStructure
	timestamp time.Time
}

type cachedArchitecture struct {
	result    *domain.ArchitectureInfo
	timestamp time.Time
}

type cachedFrameworks struct {
	result    []domain.FrameworkInfo
	timestamp time.Time
}

type cachedConventions struct {
	result    *domain.ConventionInfo
	timestamp time.Time
}

// NewCachedDetector creates a new cached detector with default TTL
func NewCachedDetector() *CachedDetector {
	return NewCachedDetectorWithTTL(DefaultCacheTTL)
}

// NewCachedDetectorWithTTL creates a new cached detector with custom TTL
func NewCachedDetectorWithTTL(ttl time.Duration) *CachedDetector {
	return &CachedDetector{
		detector:          NewDetector(),
		structureCache:    make(map[string]*cachedStructure),
		architectureCache: make(map[string]*cachedArchitecture),
		frameworksCache:   make(map[string]*cachedFrameworks),
		conventionsCache:  make(map[string]*cachedConventions),
		ttl:               ttl,
	}
}

// DetectStructure returns cached or fresh project structure
func (cd *CachedDetector) DetectStructure(projectPath string) (*domain.ProjectStructure, error) {
	cd.mu.RLock()
	if cached, ok := cd.structureCache[projectPath]; ok && cd.isValid(cached.timestamp) {
		cd.mu.RUnlock()
		return cached.result, nil
	}
	cd.mu.RUnlock()

	result, err := cd.detector.DetectStructure(projectPath)
	if err != nil {
		return nil, err
	}

	cd.mu.Lock()
	cd.structureCache[projectPath] = &cachedStructure{
		result:    result,
		timestamp: time.Now(),
	}
	cd.mu.Unlock()

	return result, nil
}

// DetectArchitecture returns cached or fresh architecture info
func (cd *CachedDetector) DetectArchitecture(projectPath string) (*domain.ArchitectureInfo, error) {
	cd.mu.RLock()
	if cached, ok := cd.architectureCache[projectPath]; ok && cd.isValid(cached.timestamp) {
		cd.mu.RUnlock()
		return cached.result, nil
	}
	cd.mu.RUnlock()

	result, err := cd.detector.DetectArchitecture(projectPath)
	if err != nil {
		return nil, err
	}

	cd.mu.Lock()
	cd.architectureCache[projectPath] = &cachedArchitecture{
		result:    result,
		timestamp: time.Now(),
	}
	cd.mu.Unlock()

	return result, nil
}

// DetectFrameworks returns cached or fresh frameworks info
func (cd *CachedDetector) DetectFrameworks(projectPath string) ([]domain.FrameworkInfo, error) {
	cd.mu.RLock()
	if cached, ok := cd.frameworksCache[projectPath]; ok && cd.isValid(cached.timestamp) {
		cd.mu.RUnlock()
		return cached.result, nil
	}
	cd.mu.RUnlock()

	result, err := cd.detector.DetectFrameworks(projectPath)
	if err != nil {
		return nil, err
	}

	cd.mu.Lock()
	cd.frameworksCache[projectPath] = &cachedFrameworks{
		result:    result,
		timestamp: time.Now(),
	}
	cd.mu.Unlock()

	return result, nil
}

// DetectConventions returns cached or fresh conventions info
func (cd *CachedDetector) DetectConventions(projectPath string) (*domain.ConventionInfo, error) {
	cd.mu.RLock()
	if cached, ok := cd.conventionsCache[projectPath]; ok && cd.isValid(cached.timestamp) {
		cd.mu.RUnlock()
		return cached.result, nil
	}
	cd.mu.RUnlock()

	result, err := cd.detector.DetectConventions(projectPath)
	if err != nil {
		return nil, err
	}

	cd.mu.Lock()
	cd.conventionsCache[projectPath] = &cachedConventions{
		result:    result,
		timestamp: time.Now(),
	}
	cd.mu.Unlock()

	return result, nil
}

// GetRelatedLayers delegates to underlying detector (no caching needed)
func (cd *CachedDetector) GetRelatedLayers(projectPath, filePath string) ([]domain.LayerInfo, error) {
	return cd.detector.GetRelatedLayers(projectPath, filePath)
}

// SuggestRelatedFiles delegates to underlying detector (no caching needed)
func (cd *CachedDetector) SuggestRelatedFiles(projectPath, filePath string) ([]string, error) {
	return cd.detector.SuggestRelatedFiles(projectPath, filePath)
}

// Invalidate clears all caches for a project
func (cd *CachedDetector) Invalidate(projectPath string) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	delete(cd.structureCache, projectPath)
	delete(cd.architectureCache, projectPath)
	delete(cd.frameworksCache, projectPath)
	delete(cd.conventionsCache, projectPath)
}

// InvalidateAll clears all caches
func (cd *CachedDetector) InvalidateAll() {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.structureCache = make(map[string]*cachedStructure)
	cd.architectureCache = make(map[string]*cachedArchitecture)
	cd.frameworksCache = make(map[string]*cachedFrameworks)
	cd.conventionsCache = make(map[string]*cachedConventions)
}

// Stats returns cache statistics
func (cd *CachedDetector) Stats() map[string]int {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return map[string]int{
		"structure_entries":    len(cd.structureCache),
		"architecture_entries": len(cd.architectureCache),
		"frameworks_entries":   len(cd.frameworksCache),
		"conventions_entries":  len(cd.conventionsCache),
	}
}

// isValid checks if cached entry is still valid
func (cd *CachedDetector) isValid(timestamp time.Time) bool {
	return time.Since(timestamp) < cd.ttl
}

// SetTTL updates the cache TTL
func (cd *CachedDetector) SetTTL(ttl time.Duration) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.ttl = ttl
}

// GetTTL returns current cache TTL
func (cd *CachedDetector) GetTTL() time.Duration {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return cd.ttl
}
