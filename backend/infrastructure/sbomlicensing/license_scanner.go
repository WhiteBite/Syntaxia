package sbomlicensing

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"syntaxia/domain"
	"syntaxia/internal/executil"
)

// LicenseScanner представляет сканер лицензий
type LicenseScanner struct {
	log domain.Logger
}

// NewLicenseScanner создает новый сканер лицензий
func NewLicenseScanner(log domain.Logger) *LicenseScanner {
	return &LicenseScanner{
		log: log,
	}
}

// IsAvailable проверяет доступность инструментов сканирования лицензий
func (l *LicenseScanner) IsAvailable() bool {
	// Проверяем доступность различных инструментов
	tools := []string{"licensecheck", "licensee", "go-licenses"}
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err == nil {
			return true
		}
	}
	return false
}

// GetAvailableTools возвращает список доступных инструментов
func (l *LicenseScanner) GetAvailableTools() []string {
	var tools []string
	availableTools := []string{"licensecheck", "licensee", "go-licenses"}

	for _, tool := range availableTools {
		if _, err := exec.LookPath(tool); err == nil {
			tools = append(tools, tool)
		}
	}

	return tools
}

// ScanLicenses сканирует лицензии в проекте
func (l *LicenseScanner) ScanLicenses(ctx context.Context, projectPath string) (*domain.LicenseScanResult, error) {
	if !l.IsAvailable() {
		return &domain.LicenseScanResult{
			Success:     false,
			ProjectPath: projectPath,
			Error:       "no license scanning tools available",
		}, nil
	}

	l.log.Info(fmt.Sprintf("Scanning licenses for: %s", projectPath))

	// Определяем язык проекта
	language := l.detectLanguage(projectPath)

	// Выбираем подходящий инструмент
	tool := l.selectTool(language)

	// Сканируем лицензии
	var licenses []*domain.LicenseInfo
	var err error

	switch tool {
	case "licensecheck":
		licenses, err = l.scanWithLicensecheck(ctx, projectPath)
	case "licensee":
		licenses, err = l.scanWithLicensee(ctx, projectPath)
	case "go-licenses":
		licenses, err = l.scanWithGoLicenses(ctx, projectPath)
	default:
		return &domain.LicenseScanResult{
			Success:     false,
			ProjectPath: projectPath,
			Error:       fmt.Sprintf("no suitable tool found for language: %s", language),
		}, nil
	}

	if err != nil {
		return &domain.LicenseScanResult{
			Success:     false,
			ProjectPath: projectPath,
			Error:       fmt.Sprintf("license scan failed: %v", err),
		}, nil
	}

	// Рассчитываем сводку
	summary := l.calculateLicenseSummary(licenses)

	return &domain.LicenseScanResult{
		Success:     true,
		ProjectPath: projectPath,
		Licenses:    licenses,
		Summary:     summary,
		Metadata: map[string]interface{}{
			"tool":     tool,
			"language": language,
		},
	}, nil
}

// detectLanguage определяет язык проекта
func (l *LicenseScanner) detectLanguage(projectPath string) string {
	if _, err := exec.LookPath(filepath.Join(projectPath, "go.mod")); err == nil {
		return "go"
	}
	if _, err := exec.LookPath(filepath.Join(projectPath, "package.json")); err == nil {
		return "node"
	}
	if _, err := exec.LookPath(filepath.Join(projectPath, "pom.xml")); err == nil {
		return "java"
	}
	if _, err := exec.LookPath(filepath.Join(projectPath, "build.gradle")); err == nil {
		return "java"
	}
	return "unknown"
}

// selectTool выбирает подходящий инструмент для языка
func (l *LicenseScanner) selectTool(language string) string {
	switch language {
	case "go":
		if _, err := exec.LookPath("go-licenses"); err == nil {
			return "go-licenses"
		}
	case "node":
		if _, err := exec.LookPath("licensee"); err == nil {
			return "licensee"
		}
	}

	// Fallback к licensecheck
	if _, err := exec.LookPath("licensecheck"); err == nil {
		return "licensecheck"
	}

	return ""
}

// scanWithLicensecheck сканирует с помощью licensecheck
func (l *LicenseScanner) scanWithLicensecheck(ctx context.Context, projectPath string) ([]*domain.LicenseInfo, error) {
	l.log.Info("Scanning licenses with licensecheck")

	cmd := exec.CommandContext(ctx, "licensecheck", "-r", projectPath)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("licensecheck failed: %w, output: %s", err, string(output))
	}

	return l.parseLicensecheckOutput(output), nil
}

// scanWithLicensee сканирует с помощью licensee
func (l *LicenseScanner) scanWithLicensee(ctx context.Context, projectPath string) ([]*domain.LicenseInfo, error) {
	l.log.Info("Scanning licenses with licensee")

	cmd := exec.CommandContext(ctx, "licensee", "detect", projectPath)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("licensee failed: %w, output: %s", err, string(output))
	}

	return l.parseLicenseeOutput(output), nil
}

// scanWithGoLicenses сканирует с помощью go-licenses
func (l *LicenseScanner) scanWithGoLicenses(ctx context.Context, projectPath string) ([]*domain.LicenseInfo, error) {
	l.log.Info("Scanning licenses with go-licenses")

	cmd := exec.CommandContext(ctx, "go-licenses", "csv", ".")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("go-licenses failed: %w, output: %s", err, string(output))
	}

	return l.parseGoLicensesOutput(output), nil
}

// parseLicensecheckOutput парсит вывод licensecheck
// Формат: "filename: LICENSE_TYPE"
func (l *LicenseScanner) parseLicensecheckOutput(output []byte) []*domain.LicenseInfo {
	l.log.Info(fmt.Sprintf("Parsing licensecheck output, size: %d bytes", len(output)))

	if len(output) == 0 {
		return []*domain.LicenseInfo{}
	}

	var licenses []*domain.LicenseInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Парсим формат "filename: LICENSE_TYPE"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		filename := strings.TrimSpace(parts[0])
		licenseType := strings.TrimSpace(parts[1])

		if licenseType == "" || licenseType == "UNKNOWN" {
			continue
		}

		license := &domain.LicenseInfo{
			Name:       licenseType,
			SPDXID:     l.mapToSPDXID(licenseType),
			Type:       l.classifyLicenseType(licenseType),
			Files:      []string{filename},
			Confidence: 1.0,
		}

		licenses = append(licenses, license)
	}

	l.log.Info(fmt.Sprintf("Parsed %d licenses from licensecheck output", len(licenses)))

	return licenses
}

// parseLicenseeOutput парсит вывод licensee
// Формат: "filename: LICENSE_TYPE (confidence: XX%)"
func (l *LicenseScanner) parseLicenseeOutput(output []byte) []*domain.LicenseInfo {
	l.log.Info(fmt.Sprintf("Parsing licensee output, size: %d bytes", len(output)))

	if len(output) == 0 {
		return []*domain.LicenseInfo{}
	}

	var licenses []*domain.LicenseInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Регулярное выражение для парсинга формата с confidence
	confidenceRegex := regexp.MustCompile(`^(.+?):\s*(.+?)\s*(?:\(confidence:\s*(\d+(?:\.\d+)?)%?\))?$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := confidenceRegex.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		filename := strings.TrimSpace(matches[1])
		licenseType := strings.TrimSpace(matches[2])

		if licenseType == "" || licenseType == "UNKNOWN" || licenseType == "NOASSERTION" {
			continue
		}

		confidence := 1.0
		if len(matches) >= 4 && matches[3] != "" {
			if parsed, err := strconv.ParseFloat(matches[3], 64); err == nil {
				confidence = parsed / 100.0
			}
		}

		license := &domain.LicenseInfo{
			Name:       licenseType,
			SPDXID:     l.mapToSPDXID(licenseType),
			Type:       l.classifyLicenseType(licenseType),
			Files:      []string{filename},
			Confidence: confidence,
		}

		licenses = append(licenses, license)
	}

	l.log.Info(fmt.Sprintf("Parsed %d licenses from licensee output", len(licenses)))

	return licenses
}

// parseGoLicensesOutput парсит вывод go-licenses
// Формат CSV: "package,license,path"
func (l *LicenseScanner) parseGoLicensesOutput(output []byte) []*domain.LicenseInfo {
	l.log.Info(fmt.Sprintf("Parsing go-licenses output, size: %d bytes", len(output)))

	if len(output) == 0 {
		return []*domain.LicenseInfo{}
	}

	var licenses []*domain.LicenseInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Парсим CSV формат: package,license,path
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		packageName := strings.TrimSpace(parts[0])
		licenseType := strings.TrimSpace(parts[1])

		if licenseType == "" || licenseType == "Unknown" {
			continue
		}

		var licensePath string
		if len(parts) >= 3 {
			licensePath = strings.TrimSpace(parts[2])
		}

		files := []string{packageName}
		if licensePath != "" {
			files = append(files, licensePath)
		}

		license := &domain.LicenseInfo{
			Name:       licenseType,
			SPDXID:     l.mapToSPDXID(licenseType),
			Type:       l.classifyLicenseType(licenseType),
			Files:      files,
			Confidence: 1.0,
		}

		licenses = append(licenses, license)
	}

	l.log.Info(fmt.Sprintf("Parsed %d licenses from go-licenses output", len(licenses)))

	return licenses
}

// mapToSPDXID преобразует название лицензии в SPDX ID
func (l *LicenseScanner) mapToSPDXID(licenseName string) string {
	spdxMap := map[string]string{
		"MIT":              "MIT",
		"MIT License":      "MIT",
		"Apache-2.0":       "Apache-2.0",
		"Apache 2.0":       "Apache-2.0",
		"Apache License 2": "Apache-2.0",
		"GPL-3.0":          "GPL-3.0-only",
		"GPLv3":            "GPL-3.0-only",
		"GPL-2.0":          "GPL-2.0-only",
		"GPLv2":            "GPL-2.0-only",
		"LGPL-3.0":         "LGPL-3.0-only",
		"LGPLv3":           "LGPL-3.0-only",
		"BSD-3-Clause":     "BSD-3-Clause",
		"BSD 3-Clause":     "BSD-3-Clause",
		"BSD-2-Clause":     "BSD-2-Clause",
		"BSD 2-Clause":     "BSD-2-Clause",
		"ISC":              "ISC",
		"MPL-2.0":          "MPL-2.0",
		"CC0-1.0":          "CC0-1.0",
		"Unlicense":        "Unlicense",
	}

	if spdxID, ok := spdxMap[licenseName]; ok {
		return spdxID
	}

	// Пробуем найти частичное совпадение
	lowerName := strings.ToLower(licenseName)
	for name, spdxID := range spdxMap {
		if strings.Contains(lowerName, strings.ToLower(name)) {
			return spdxID
		}
	}

	return licenseName
}

// classifyLicenseType классифицирует тип лицензии
func (l *LicenseScanner) classifyLicenseType(licenseName string) string {
	lowerName := strings.ToLower(licenseName)

	// Copyleft лицензии
	copyleftLicenses := []string{"gpl", "lgpl", "agpl", "cc-by-sa", "mpl"}
	for _, cl := range copyleftLicenses {
		if strings.Contains(lowerName, cl) {
			return "copyleft"
		}
	}

	// Permissive лицензии
	permissiveLicenses := []string{"mit", "apache", "bsd", "isc", "cc0", "unlicense", "wtfpl", "zlib"}
	for _, pl := range permissiveLicenses {
		if strings.Contains(lowerName, pl) {
			return "permissive"
		}
	}

	// Проприетарные
	proprietaryLicenses := []string{"proprietary", "commercial", "private"}
	for _, prop := range proprietaryLicenses {
		if strings.Contains(lowerName, prop) {
			return "proprietary"
		}
	}

	return "unknown"
}

// calculateLicenseSummary рассчитывает сводку лицензий
func (l *LicenseScanner) calculateLicenseSummary(licenses []*domain.LicenseInfo) *domain.LicenseSummary {
	summary := &domain.LicenseSummary{
		TotalLicenses: len(licenses),
		ByType:        make(map[string]int),
		ByLicense:     make(map[string]int),
		Conflicts:     []*domain.LicenseConflict{},
	}

	// Подсчитываем лицензии по типам и названиям
	for _, license := range licenses {
		summary.ByType[license.Type]++
		summary.ByLicense[license.Name]++
	}

	// Проверяем конфликты лицензий
	summary.Conflicts = l.detectLicenseConflicts(licenses)

	return summary
}

// detectLicenseConflicts обнаруживает конфликты лицензий
func (l *LicenseScanner) detectLicenseConflicts(licenses []*domain.LicenseInfo) []*domain.LicenseConflict {
	var conflicts []*domain.LicenseConflict

	// Упрощенная логика обнаружения конфликтов
	// В реальном приложении здесь будет более сложная логика

	// Пример: GPL и проприетарные лицензии
	gplLicenses := make(map[string]bool)
	proprietaryLicenses := make(map[string]bool)

	for _, license := range licenses {
		if strings.Contains(strings.ToLower(license.Name), "gpl") {
			gplLicenses[license.Name] = true
		}
		if strings.Contains(strings.ToLower(license.Type), "proprietary") {
			proprietaryLicenses[license.Name] = true
		}
	}

	// Проверяем конфликты
	if len(gplLicenses) > 0 && len(proprietaryLicenses) > 0 {
		for gplLicense := range gplLicenses {
			for propLicense := range proprietaryLicenses {
				conflicts = append(conflicts, &domain.LicenseConflict{
					License1:    gplLicense,
					License2:    propLicense,
					Description: "GPL license conflicts with proprietary license",
					Severity:    "high",
				})
			}
		}
	}

	return conflicts
}

// ValidateLicense проверяет валидность лицензии
func (l *LicenseScanner) ValidateLicense(licenseName string) (bool, string) {
	// Список известных лицензий
	knownLicenses := map[string]bool{
		"MIT":          true,
		"Apache-2.0":   true,
		"GPL-3.0":      true,
		"GPL-2.0":      true,
		"LGPL-3.0":     true,
		"BSD-3-Clause": true,
		"BSD-2-Clause": true,
		"ISC":          true,
		"MPL-2.0":      true,
		"CC0-1.0":      true,
		"Unlicense":    true,
	}

	if knownLicenses[licenseName] {
		return true, "Valid license"
	}

	return false, "Unknown or invalid license"
}
