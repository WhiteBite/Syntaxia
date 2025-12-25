package projectstructure

import (
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

type architectureDetector interface {
	detect(projectPath string) *domain.ArchitectureInfo
}

func (d *Detector) initArchitectureDetectors() {
	d.archDetectors = []architectureDetector{
		&cleanArchDetector{},
		&hexagonalDetector{},
		&mvcDetector{},
		&mvvmDetector{},
		&layeredDetector{},
		&dddDetector{},
	}
}

// Helper functions for architecture detection

// scanDirsForWeights scans directories and returns score and indicators
func scanDirsForWeights(searchPath, projectPath string, weights map[string]float64) (float64, []string) {
	entries, err := os.ReadDir(searchPath)
	if err != nil {
		return 0, nil
	}
	var score float64
	var indicators []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if weight, ok := weights[name]; ok {
			score += weight
			rel, _ := filepath.Rel(projectPath, filepath.Join(searchPath, entry.Name()))
			indicators = append(indicators, "Found "+rel+" directory")
		}
	}
	return score, indicators
}

// capScore caps score at 0.95
func capScore(score float64) float64 {
	if score > 0.95 {
		return 0.95
	}
	return score
}

// buildArchInfo creates ArchitectureInfo if score meets threshold
func buildArchInfo(archType domain.ArchitectureType, score, threshold float64, desc string, indicators []string) *domain.ArchitectureInfo {
	if score < threshold {
		return nil
	}
	return &domain.ArchitectureInfo{
		Type:        archType,
		Confidence:  capScore(score),
		Description: desc,
		Indicators:  indicators,
	}
}

// Clean Architecture detector
type cleanArchDetector struct{}

func (c *cleanArchDetector) detect(projectPath string) *domain.ArchitectureInfo {
	cleanDirs := map[string]float64{
		"domain": 0.2, "application": 0.2, "infrastructure": 0.2, "interfaces": 0.1,
		"usecases": 0.15, "use_cases": 0.15, "use-cases": 0.15, "entities": 0.15,
	}

	score, indicators := scanDirsForWeights(projectPath, projectPath, cleanDirs)

	// Check nested structures
	for _, subdir := range []string{"backend", "internal"} {
		subPath := filepath.Join(projectPath, subdir)
		s, ind := scanDirsForWeights(subPath, projectPath, cleanDirs)
		if subdir == "internal" {
			s *= 0.8
		}
		score += s
		indicators = append(indicators, ind...)
	}

	return buildArchInfo(domain.ArchCleanArchitecture, score, 0.3,
		"Clean Architecture with separation of domain, application, and infrastructure layers", indicators)
}

// Hexagonal Architecture detector
type hexagonalDetector struct{}

func (h *hexagonalDetector) detect(projectPath string) *domain.ArchitectureInfo {
	hexDirs := map[string]float64{
		"ports": 0.25, "adapters": 0.25, "core": 0.2,
		"inbound": 0.15, "outbound": 0.15, "driven": 0.15, "driving": 0.15,
	}

	score, indicators := scanDirsForWeights(projectPath, projectPath, hexDirs)
	s, ind := scanDirsForWeights(filepath.Join(projectPath, "src"), projectPath, hexDirs)
	score += s
	indicators = append(indicators, ind...)

	return buildArchInfo(domain.ArchHexagonal, score, 0.3,
		"Hexagonal Architecture (Ports & Adapters) with clear separation of core logic and external interfaces", indicators)
}

// MVC detector
type mvcDetector struct{}

func (m *mvcDetector) detect(projectPath string) *domain.ArchitectureInfo {
	mvcDirs := map[string]float64{
		"models": 0.3, "views": 0.3, "controllers": 0.3,
		"model": 0.25, "view": 0.25, "controller": 0.25,
	}

	var score float64
	var indicators []string
	for _, subdir := range []string{"", "app", "src"} {
		searchPath := filepath.Join(projectPath, subdir)
		s, ind := scanDirsForWeights(searchPath, projectPath, mvcDirs)
		score += s
		indicators = append(indicators, ind...)
	}

	// Check for Rails-style structure
	if _, err := os.Stat(filepath.Join(projectPath, "app", "models")); err == nil {
		if _, err := os.Stat(filepath.Join(projectPath, "app", "controllers")); err == nil {
			score += 0.3
			indicators = append(indicators, "Rails-style app/models and app/controllers structure")
		}
	}

	return buildArchInfo(domain.ArchMVC, score, 0.5,
		"Model-View-Controller architecture with separation of data, presentation, and control logic", indicators)
}

// MVVM detector
type mvvmDetector struct{}

func (m *mvvmDetector) detect(projectPath string) *domain.ArchitectureInfo {
	mvvmDirs := map[string]float64{
		"viewmodels": 0.35, "viewmodel": 0.35, "view-models": 0.35, "models": 0.2, "views": 0.2,
	}

	var score float64
	var indicators []string
	for _, subdir := range []string{"", "src", "lib"} {
		searchPath := filepath.Join(projectPath, subdir)
		s, ind := scanDirsForWeights(searchPath, projectPath, mvvmDirs)
		score += s
		indicators = append(indicators, ind...)
	}

	if count := countViewModelFiles(projectPath); count > 3 {
		score += 0.2
		indicators = append(indicators, "Found multiple ViewModel files")
	}

	return buildArchInfo(domain.ArchMVVM, score, 0.4,
		"Model-View-ViewModel architecture with data binding between View and ViewModel", indicators)
}

// countViewModelFiles counts files with viewmodel in name
func countViewModelFiles(projectPath string) int {
	count := 0
	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
			return filepath.SkipDir
		}
		name := strings.ToLower(info.Name())
		if strings.Contains(name, "viewmodel") || strings.Contains(name, "view_model") {
			count++
		}
		return nil
	})
	return count
}

// Layered Architecture detector
type layeredDetector struct{}

func (l *layeredDetector) detect(projectPath string) *domain.ArchitectureInfo {
	layeredDirs := map[string]float64{
		"presentation": 0.2, "business": 0.2, "data": 0.2, "services": 0.15,
		"handlers": 0.15, "repositories": 0.15, "repository": 0.15,
		"api": 0.1, "dal": 0.15, "bll": 0.15,
	}

	var score float64
	var indicators []string
	for _, subdir := range []string{"", "backend", "src", "app"} {
		searchPath := filepath.Join(projectPath, subdir)
		s, ind := scanDirsForWeights(searchPath, projectPath, layeredDirs)
		score += s
		indicators = append(indicators, ind...)
	}

	return buildArchInfo(domain.ArchLayered, score, 0.3,
		"Layered architecture with separation of presentation, business, and data access layers", indicators)
}

// DDD (Domain-Driven Design) detector
type dddDetector struct{}

func (d *dddDetector) detect(projectPath string) *domain.ArchitectureInfo {
	dddDirs := map[string]float64{
		"aggregates": 0.25, "entities": 0.2, "value_objects": 0.2, "valueobjects": 0.2,
		"value-objects": 0.2, "repositories": 0.15, "services": 0.1,
		"events": 0.15, "commands": 0.15, "queries": 0.15,
	}

	// Search in domain directory first
	score, indicators := scanDirsForWeights(filepath.Join(projectPath, "domain"), projectPath, dddDirs)

	// Check root level with reduced weight
	s, ind := scanDirsForWeights(projectPath, projectPath, dddDirs)
	score += s * 0.8
	indicators = append(indicators, ind...)

	// Check for CQRS pattern
	if hasCQRSPattern(projectPath) {
		score += 0.2
		indicators = append(indicators, "CQRS pattern detected (commands/queries separation)")
	}

	return buildArchInfo(domain.ArchDDD, score, 0.3,
		"Domain-Driven Design with rich domain model, aggregates, and bounded contexts", indicators)
}

// hasCQRSPattern checks for commands/queries separation
func hasCQRSPattern(projectPath string) bool {
	hasCommands, hasQueries := false, false
	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
			return filepath.SkipDir
		}
		if info.IsDir() {
			name := strings.ToLower(info.Name())
			if name == "commands" {
				hasCommands = true
			}
			if name == "queries" {
				hasQueries = true
			}
		}
		return nil
	})
	return hasCommands && hasQueries
}
