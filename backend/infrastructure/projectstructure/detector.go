package projectstructure

import (
	"os"
	"path/filepath"
	"regexp"
	"syntaxia/domain"
	"strings"
)

// Category constants
const (
	categoryWeb = "web"
)

// Detector implements ProjectStructureDetector
type Detector struct {
	frameworkDetectors []frameworkDetector
	archDetectors      []architectureDetector
}

// NewDetector creates a new project structure detector
func NewDetector() *Detector {
	d := &Detector{}
	d.initFrameworkDetectors()
	d.initArchitectureDetectors()
	return d
}

// DetectStructure analyzes project and returns complete structure info
func (d *Detector) DetectStructure(projectPath string) (*domain.ProjectStructure, error) {
	arch, _ := d.DetectArchitecture(projectPath)
	conventions, _ := d.DetectConventions(projectPath)
	frameworks, _ := d.DetectFrameworks(projectPath)
	buildSystems := d.detectBuildSystems(projectPath)
	languages := d.detectLanguages(projectPath)
	projectType := d.detectProjectType(projectPath, frameworks, buildSystems)

	confidence := 0.5
	if arch != nil && arch.Confidence > 0 {
		confidence = arch.Confidence
	}

	return &domain.ProjectStructure{
		Architecture: arch,
		Conventions:  conventions,
		Frameworks:   frameworks,
		BuildSystems: buildSystems,
		Languages:    languages,
		ProjectType:  projectType,
		Confidence:   confidence,
	}, nil
}

// DetectArchitecture detects architecture pattern
func (d *Detector) DetectArchitecture(projectPath string) (*domain.ArchitectureInfo, error) {
	var bestMatch *domain.ArchitectureInfo
	var bestScore float64

	for _, detector := range d.archDetectors {
		info := detector.detect(projectPath)
		if info != nil && info.Confidence > bestScore {
			bestScore = info.Confidence
			bestMatch = info
		}
	}

	if bestMatch == nil {
		return &domain.ArchitectureInfo{
			Type:        domain.ArchUnknown,
			Confidence:  0.0,
			Description: "Could not determine architecture pattern",
		}, nil
	}

	// Detect layers
	bestMatch.Layers = d.detectLayers(projectPath, bestMatch.Type)
	return bestMatch, nil
}

// DetectConventions detects naming and code conventions
func (d *Detector) DetectConventions(projectPath string) (*domain.ConventionInfo, error) {
	conventions := &domain.ConventionInfo{
		NamingStyle:     d.detectNamingStyle(projectPath),
		FileNaming:      d.detectFileNaming(projectPath),
		FolderStructure: d.detectFolderStructure(projectPath),
		TestConventions: d.detectTestConventions(projectPath),
		ImportStyle:     d.detectImportStyle(projectPath),
		CodeStyle:       d.detectCodeStyle(projectPath),
	}
	return conventions, nil
}

// DetectFrameworks detects frameworks used in the project
func (d *Detector) DetectFrameworks(projectPath string) ([]domain.FrameworkInfo, error) {
	var frameworks []domain.FrameworkInfo

	for _, detector := range d.frameworkDetectors {
		if info := detector.detect(projectPath); info != nil {
			frameworks = append(frameworks, *info)
		}
	}

	return frameworks, nil
}

// GetRelatedLayers returns layers related to a file
func (d *Detector) GetRelatedLayers(projectPath, filePath string) ([]domain.LayerInfo, error) {
	arch, err := d.DetectArchitecture(projectPath)
	if err != nil || arch == nil {
		return nil, err
	}

	var related []domain.LayerInfo
	relPath, _ := filepath.Rel(projectPath, filePath)

	for _, layer := range arch.Layers {
		if strings.HasPrefix(relPath, layer.Path) {
			related = append(related, layer)
			// Add dependent layers
			for _, dep := range layer.Dependencies {
				for _, l := range arch.Layers {
					if l.Name == dep {
						related = append(related, l)
					}
				}
			}
		}
	}

	return related, nil
}

// SuggestRelatedFiles suggests related files based on architecture
func (d *Detector) SuggestRelatedFiles(projectPath, filePath string) ([]string, error) {
	relPath, _ := filepath.Rel(projectPath, filePath)
	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	coreName := extractCoreName(strings.TrimSuffix(baseName, ext))

	var suggestions []string
	suggestions = append(suggestions, d.findFilesInOtherLayers(projectPath, relPath, coreName)...)
	suggestions = append(suggestions, d.findTestFiles(projectPath, relPath, baseName, ext)...)

	return deduplicateStrings(suggestions, relPath), nil
}

// extractCoreName removes common suffixes to get base name
func extractCoreName(name string) string {
	suffixes := []string{"_test", ".test", ".spec", "_spec", "_handler", "_service", "_repository", "_controller", "_model", "_entity"}
	for _, suffix := range suffixes {
		name = strings.TrimSuffix(name, suffix)
	}
	return name
}

// findFilesInOtherLayers finds related files in architectural layers
func (d *Detector) findFilesInOtherLayers(projectPath, relPath, coreName string) []string {
	var suggestions []string
	arch, _ := d.DetectArchitecture(projectPath)
	if arch == nil {
		return suggestions
	}

	for _, layer := range arch.Layers {
		if strings.HasPrefix(relPath, layer.Path) {
			continue
		}
		layerPath := filepath.Join(projectPath, layer.Path)
		_ = filepath.Walk(layerPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(coreName)) {
				rel, _ := filepath.Rel(projectPath, path)
				suggestions = append(suggestions, rel)
			}
			return nil
		})
	}
	return suggestions
}

// findTestFiles finds test files for a given source file
func (d *Detector) findTestFiles(projectPath, relPath, baseName, ext string) []string {
	var suggestions []string
	testPatterns := []string{
		strings.Replace(relPath, ext, "_test"+ext, 1),
		strings.Replace(relPath, ext, ".test"+ext, 1),
		strings.Replace(relPath, ext, ".spec"+ext, 1),
		filepath.Join(filepath.Dir(relPath), "__tests__", baseName),
		filepath.Join("tests", relPath),
		filepath.Join("test", relPath),
	}

	for _, pattern := range testPatterns {
		if _, err := os.Stat(filepath.Join(projectPath, pattern)); err == nil {
			suggestions = append(suggestions, pattern)
		}
	}
	return suggestions
}

// deduplicateStrings removes duplicates and excludes a specific path
func deduplicateStrings(items []string, exclude string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, s := range items {
		if !seen[s] && s != exclude {
			seen[s] = true
			unique = append(unique, s)
		}
	}
	return unique
}

// detectLayers detects architectural layers based on architecture type
func (d *Detector) detectLayers(projectPath string, archType domain.ArchitectureType) []domain.LayerInfo {
	var layers []domain.LayerInfo

	// Common layer patterns to look for
	layerPatterns := map[string][]string{
		"domain":         {"domain", "entities", "models", "core"},
		"application":    {"application", "services", "usecases", "use_cases", "use-cases"},
		"infrastructure": {"infrastructure", "infra", "adapters", "external"},
		"presentation":   {"presentation", "ui", "views", "pages", "screens"},
		"handlers":       {"handlers", "controllers", "api", "routes", "endpoints"},
		"repository":     {"repository", "repositories", "repo", "data", "persistence"},
	}

	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return layers
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := strings.ToLower(entry.Name())
		for layerName, patterns := range layerPatterns {
			for _, pattern := range patterns {
				if name == pattern || strings.Contains(name, pattern) {
					layer := domain.LayerInfo{
						Name:        layerName,
						Path:        entry.Name(),
						Description: d.getLayerDescription(layerName),
						Patterns:    d.detectPatternsInDir(filepath.Join(projectPath, entry.Name())),
					}
					layer.Dependencies = d.getLayerDependencies(layerName, archType)
					layers = append(layers, layer)
					break
				}
			}
		}
	}

	return layers
}

func (d *Detector) getLayerDescription(layerName string) string {
	descriptions := map[string]string{
		"domain":         "Core business logic and entities",
		"application":    "Application services and use cases",
		"infrastructure": "External services, databases, APIs",
		"presentation":   "User interface components",
		"handlers":       "HTTP/API request handlers",
		"repository":     "Data access and persistence",
	}
	if desc, ok := descriptions[layerName]; ok {
		return desc
	}
	return ""
}

func (d *Detector) getLayerDependencies(layerName string, archType domain.ArchitectureType) []string {
	// Clean Architecture dependencies
	if archType == domain.ArchCleanArchitecture || archType == domain.ArchHexagonal {
		deps := map[string][]string{
			"handlers":       {"application"},
			"application":    {"domain"},
			"infrastructure": {"domain", "application"},
			"presentation":   {"application"},
			"repository":     {"domain"},
		}
		if d, ok := deps[layerName]; ok {
			return d
		}
	}
	return nil
}

func (d *Detector) detectPatternsInDir(dirPath string) []string {
	var patterns []string
	patternIndicators := map[string][]string{
		"repository": {"Repository", "Repo", "Store", "DAO"},
		"service":    {"Service", "UseCase", "Interactor"},
		"handler":    {"Handler", "Controller", "Endpoint"},
		"factory":    {"Factory", "Builder", "Creator"},
		"observer":   {"Observer", "Listener", "Subscriber"},
	}

	_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		name := info.Name()
		for pattern, indicators := range patternIndicators {
			for _, ind := range indicators {
				if strings.Contains(name, ind) {
					patterns = append(patterns, pattern)
					break
				}
			}
		}
		return nil
	})

	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, p := range patterns {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}
	return unique
}

func (d *Detector) detectNamingStyle(projectPath string) domain.NamingStyle {
	var camelCount, snakeCount, pascalCount, kebabCount int

	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
			return filepath.SkipDir
		}

		name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		if regexp.MustCompile(`^[a-z]+[A-Z]`).MatchString(name) {
			camelCount++
		}
		if regexp.MustCompile(`^[a-z]+_[a-z]`).MatchString(name) {
			snakeCount++
		}
		if regexp.MustCompile(`^[A-Z][a-z]+[A-Z]`).MatchString(name) {
			pascalCount++
		}
		if regexp.MustCompile(`^[a-z]+-[a-z]`).MatchString(name) {
			kebabCount++
		}
		return nil
	})

	maxCount := camelCount
	style := domain.NamingCamelCase

	if snakeCount > maxCount {
		maxCount = snakeCount
		style = domain.NamingSnakeCase
	}
	if pascalCount > maxCount {
		maxCount = pascalCount
		style = domain.NamingPascalCase
	}
	if kebabCount > maxCount {
		style = domain.NamingKebabCase
	}

	return style
}

func (d *Detector) detectFileNaming(projectPath string) domain.FileNamingStyle {
	style := domain.FileNamingStyle{
		Style:    d.detectNamingStyle(projectPath),
		Suffixes: []string{},
		Prefixes: []string{},
		Examples: []string{},
	}

	suffixCounts, prefixCounts := countFileAffixes(projectPath)
	style.Suffixes = filterByMinCount(suffixCounts, 2)
	style.Prefixes = filterByMinCount(prefixCounts, 2)
	return style
}

// countFileAffixes counts suffix and prefix occurrences in file names
func countFileAffixes(projectPath string) (map[string]int, map[string]int) {
	suffixCounts := make(map[string]int)
	prefixCounts := make(map[string]int)
	suffixes := []string{"_test", ".test", ".spec", "_spec", "_handler", "_service", "_controller", "_model", "_repository"}
	prefixes := []string{"test_", "spec_", "I", "Abstract"}

	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
			return filepath.SkipDir
		}

		name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		for _, suffix := range suffixes {
			if strings.HasSuffix(name, suffix) {
				suffixCounts[suffix]++
			}
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				prefixCounts[prefix]++
			}
		}
		return nil
	})
	return suffixCounts, prefixCounts
}

// filterByMinCount returns keys with count > minCount
func filterByMinCount(counts map[string]int, minCount int) []string {
	var result []string
	for key, count := range counts {
		if count > minCount {
			result = append(result, key)
		}
	}
	return result
}

func (d *Detector) detectFolderStructure(projectPath string) domain.FolderStructure {
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return domain.FolderFlat
	}

	dirNames := collectDirNames(entries)
	if len(dirNames) < 3 {
		return domain.FolderFlat
	}

	layerNames := []string{"domain", "application", "infrastructure", "handlers", "services", "controllers", "models", "views"}
	typeNames := []string{"components", "utils", "helpers", "types", "interfaces", "constants"}
	featureIndicators := []string{"features", "modules", "pages", "screens"}

	if matchesAny(dirNames, featureIndicators) {
		return domain.FolderByFeature
	}
	if matchesAnyPartial(dirNames, layerNames) {
		return domain.FolderByLayer
	}
	if matchesAny(dirNames, typeNames) {
		return domain.FolderByType
	}
	return domain.FolderHybrid
}

// collectDirNames returns lowercase names of directories
func collectDirNames(entries []os.DirEntry) []string {
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, strings.ToLower(entry.Name()))
		}
	}
	return names
}

// matchesAny checks if any dir name exactly matches patterns
func matchesAny(dirNames, patterns []string) bool {
	for _, name := range dirNames {
		for _, p := range patterns {
			if name == p {
				return true
			}
		}
	}
	return false
}

// matchesAnyPartial checks if any dir name matches or contains patterns
func matchesAnyPartial(dirNames, patterns []string) bool {
	for _, name := range dirNames {
		for _, p := range patterns {
			if name == p || strings.Contains(name, p) {
				return true
			}
		}
	}
	return false
}

func (d *Detector) detectTestConventions(projectPath string) domain.TestConventions {
	conventions := domain.TestConventions{
		Location:   "same-dir",
		FileSuffix: "_test",
	}

	// Check for test directories
	testDirs := []string{"tests", "test", "__tests__", "spec"}
	for _, dir := range testDirs {
		if _, err := os.Stat(filepath.Join(projectPath, dir)); err == nil {
			conventions.Location = dir
			break
		}
	}

	// Detect test framework
	if _, err := os.Stat(filepath.Join(projectPath, "jest.config.js")); err == nil {
		conventions.Framework = "jest"
		conventions.FileSuffix = ".test"
	} else if _, err := os.Stat(filepath.Join(projectPath, "vitest.config.ts")); err == nil {
		conventions.Framework = "vitest"
		conventions.FileSuffix = ".test"
	} else if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		conventions.Framework = "go test"
		conventions.FileSuffix = "_test"
	} else if _, err := os.Stat(filepath.Join(projectPath, "pytest.ini")); err == nil {
		conventions.Framework = "pytest"
		conventions.FileSuffix = "_test"
	} else if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		conventions.Framework = "junit"
		conventions.FileSuffix = "Test"
	}

	return conventions
}

func (d *Detector) detectImportStyle(projectPath string) domain.ImportStyle {
	style := domain.ImportStyle{
		ImportOrder: []string{"stdlib", "external", "internal"},
	}

	// Check for tsconfig paths
	if content, err := os.ReadFile(filepath.Join(projectPath, "tsconfig.json")); err == nil {
		if strings.Contains(string(content), `"paths"`) {
			style.AliasedImports = true
		}
		if strings.Contains(string(content), `"baseUrl"`) {
			style.AbsoluteImports = true
		}
	}

	// Check for vite config
	if content, err := os.ReadFile(filepath.Join(projectPath, "vite.config.ts")); err == nil {
		if strings.Contains(string(content), "alias") {
			style.AliasedImports = true
		}
	}

	return style
}

func (d *Detector) detectCodeStyle(projectPath string) domain.CodeStyleInfo {
	style := domain.CodeStyleInfo{
		IndentStyle:   "spaces",
		IndentSize:    2,
		MaxLineLength: 120,
	}

	// Check for config files
	configFiles := []string{".prettierrc", ".prettierrc.json", ".editorconfig", ".eslintrc.json", ".eslintrc.js"}
	for _, cf := range configFiles {
		if _, err := os.Stat(filepath.Join(projectPath, cf)); err == nil {
			style.ConfigFile = cf
			break
		}
	}

	// Parse .editorconfig if exists
	if content, err := os.ReadFile(filepath.Join(projectPath, ".editorconfig")); err == nil {
		contentStr := string(content)
		if strings.Contains(contentStr, "indent_style = tab") {
			style.IndentStyle = "tabs"
		}
		if strings.Contains(contentStr, "indent_size = 4") {
			style.IndentSize = 4
		}
	}

	return style
}

func (d *Detector) detectBuildSystems(projectPath string) []domain.BuildSystemInfo {
	var systems []domain.BuildSystemInfo

	buildFiles := map[string]domain.BuildSystemInfo{
		"package.json":   {Name: "npm", ConfigFile: "package.json"},
		"Makefile":       {Name: "make", ConfigFile: "Makefile"},
		"go.mod":         {Name: "go", ConfigFile: "go.mod"},
		"Cargo.toml":     {Name: "cargo", ConfigFile: "Cargo.toml"},
		"pom.xml":        {Name: "maven", ConfigFile: "pom.xml"},
		"build.gradle":   {Name: "gradle", ConfigFile: "build.gradle"},
		"CMakeLists.txt": {Name: "cmake", ConfigFile: "CMakeLists.txt"},
	}

	for file, info := range buildFiles {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			// Extract scripts if package.json
			if file == "package.json" {
				if content, err := os.ReadFile(filepath.Join(projectPath, file)); err == nil {
					info.Scripts = extractNpmScripts(string(content))
				}
			}
			systems = append(systems, info)
		}
	}

	return systems
}

func extractNpmScripts(content string) []string {
	var scripts []string
	// Simple extraction - look for "scripts" section
	if idx := strings.Index(content, `"scripts"`); idx != -1 {
		// Find the scripts object
		start := strings.Index(content[idx:], "{")
		if start != -1 {
			end := strings.Index(content[idx+start:], "}")
			if end != -1 {
				scriptsSection := content[idx+start : idx+start+end+1]
				// Extract script names
				re := regexp.MustCompile(`"([^"]+)":\s*"`)
				matches := re.FindAllStringSubmatch(scriptsSection, -1)
				for _, m := range matches {
					if len(m) > 1 {
						scripts = append(scripts, m[1])
					}
				}
			}
		}
	}
	return scripts
}

func (d *Detector) detectLanguages(projectPath string) []domain.LanguageInfo {
	langCounts := make(map[string]int)
	totalFiles := 0

	extToLang := map[string]string{
		".go":    "Go",
		".ts":    "TypeScript",
		".tsx":   "TypeScript",
		".js":    "JavaScript",
		".jsx":   "JavaScript",
		".vue":   "Vue",
		".py":    "Python",
		".java":  "Java",
		".kt":    "Kotlin",
		".rs":    "Rust",
		".cs":    "C#",
		".cpp":   "C++",
		".c":     "C",
		".rb":    "Ruby",
		".php":   "PHP",
		".swift": "Swift",
		".dart":  "Dart",
	}

	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "node_modules") || strings.Contains(path, ".git") || strings.Contains(path, "vendor") {
			return filepath.SkipDir
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if lang, ok := extToLang[ext]; ok {
			langCounts[lang]++
			totalFiles++
		}
		return nil
	})

	languages := make([]domain.LanguageInfo, 0, len(langCounts))
	var maxCount int
	var primaryLang string

	for lang, count := range langCounts {
		if count > maxCount {
			maxCount = count
			primaryLang = lang
		}
	}

	for lang, count := range langCounts {
		percentage := 0.0
		if totalFiles > 0 {
			percentage = float64(count) / float64(totalFiles) * 100
		}
		languages = append(languages, domain.LanguageInfo{
			Name:       lang,
			FileCount:  count,
			Percentage: percentage,
			Primary:    lang == primaryLang,
		})
	}

	return languages
}

func (d *Detector) detectProjectType(projectPath string, frameworks []domain.FrameworkInfo, buildSystems []domain.BuildSystemInfo) string {
	// Check for specific indicators
	if _, err := os.Stat(filepath.Join(projectPath, "main.go")); err == nil {
		if _, err := os.Stat(filepath.Join(projectPath, "cmd")); err == nil {
			return "cli"
		}
	}

	for _, fw := range frameworks {
		if fw.Category == categoryWeb {
			return categoryWeb
		}
	}

	// Check for library indicators
	for _, bs := range buildSystems {
		if bs.Name == "npm" {
			if content, err := os.ReadFile(filepath.Join(projectPath, "package.json")); err == nil {
				if strings.Contains(string(content), `"main"`) && !strings.Contains(string(content), `"bin"`) {
					return "library"
				}
			}
		}
	}

	// Check for monorepo
	if _, err := os.Stat(filepath.Join(projectPath, "packages")); err == nil {
		return "monorepo"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "apps")); err == nil {
		return "monorepo"
	}

	return "service"
}
