package sbomlicensing

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"syntaxia/domain"
	"syntaxia/internal/executil"
)

// grypeOutput представляет структуру JSON вывода Grype
type grypeOutput struct {
	Matches []grypeMatch `json:"matches"`
}

// grypeMatch представляет одно совпадение уязвимости
type grypeMatch struct {
	Vulnerability grypeVulnerability `json:"vulnerability"`
	Artifact      grypeArtifact      `json:"artifact"`
}

// grypeVulnerability представляет информацию об уязвимости
type grypeVulnerability struct {
	ID          string   `json:"id"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	Fix         grypeFix `json:"fix"`
	CVSS        []cvss   `json:"cvss"`
}

// grypeFix представляет информацию об исправлении
type grypeFix struct {
	Versions []string `json:"versions"`
	State    string   `json:"state"`
}

// cvss представляет CVSS оценку
type cvss struct {
	Version string  `json:"version"`
	Vector  string  `json:"vector"`
	Score   float64 `json:"metrics"`
}

// grypeArtifact представляет артефакт (пакет)
type grypeArtifact struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
	PURL    string `json:"purl"`
}

const versionUnknown = "unknown"

// GrypeScanner представляет сканер уязвимостей с использованием Grype
type GrypeScanner struct {
	log domain.Logger
}

// NewGrypeScanner создает новый сканер Grype
func NewGrypeScanner(log domain.Logger) *GrypeScanner {
	return &GrypeScanner{
		log: log,
	}
}

// IsAvailable проверяет доступность Grype
func (g *GrypeScanner) IsAvailable() bool {
	_, err := exec.LookPath("grype")
	return err == nil
}

// ScanVulnerabilities сканирует уязвимости в проекте
func (g *GrypeScanner) ScanVulnerabilities(ctx context.Context, projectPath string) (*domain.VulnerabilityScanResult, error) {
	if !g.IsAvailable() {
		return &domain.VulnerabilityScanResult{
			Success:     false,
			ProjectPath: projectPath,
			Error:       "grype is not available",
		}, nil
	}

	g.log.Info(fmt.Sprintf("Scanning vulnerabilities with Grype for: %s", projectPath))

	// Запускаем Grype
	cmd := exec.CommandContext(ctx, "grype", projectPath, "-o", "json")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &domain.VulnerabilityScanResult{
			Success:     false,
			ProjectPath: projectPath,
			Error:       fmt.Sprintf("grype failed: %v, output: %s", err, string(output)),
		}, nil
	}

	// Парсим результат
	vulnerabilities, err := g.parseGrypeOutput(output)
	if err != nil {
		return &domain.VulnerabilityScanResult{
			Success:     false,
			ProjectPath: projectPath,
			Error:       fmt.Sprintf("failed to parse grype output: %v", err),
		}, nil
	}

	// Рассчитываем сводку
	summary := g.calculateVulnerabilitySummary(vulnerabilities)

	return &domain.VulnerabilityScanResult{
		Success:         true,
		ProjectPath:     projectPath,
		Vulnerabilities: vulnerabilities,
		Summary:         summary,
		Metadata: map[string]interface{}{
			"tool":    "grype",
			"version": g.getVersion(),
			"command": cmd.String(),
		},
	}, nil
}

// ScanSBOM сканирует уязвимости в SBOM файле
func (g *GrypeScanner) ScanSBOM(ctx context.Context, sbomPath string) (*domain.VulnerabilityScanResult, error) {
	if !g.IsAvailable() {
		return &domain.VulnerabilityScanResult{
			Success: false,
			Error:   "grype is not available",
		}, nil
	}

	g.log.Info(fmt.Sprintf("Scanning SBOM with Grype: %s", sbomPath))

	// Запускаем Grype для SBOM
	cmd := exec.CommandContext(ctx, "grype", "sbom:"+sbomPath, "-o", "json") //nolint:gosec // External tool command
	executil.HideWindow(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &domain.VulnerabilityScanResult{
			Success: false,
			Error:   fmt.Sprintf("grype sbom scan failed: %v, output: %s", err, string(output)),
		}, nil
	}

	// Парсим результат
	vulnerabilities, err := g.parseGrypeOutput(output)
	if err != nil {
		return &domain.VulnerabilityScanResult{
			Success: false,
			Error:   fmt.Sprintf("failed to parse grype output: %v", err),
		}, nil
	}

	// Рассчитываем сводку
	summary := g.calculateVulnerabilitySummary(vulnerabilities)

	return &domain.VulnerabilityScanResult{
		Success:         true,
		Vulnerabilities: vulnerabilities,
		Summary:         summary,
		Metadata: map[string]interface{}{
			"tool":      "grype",
			"version":   g.getVersion(),
			"command":   cmd.String(),
			"sbom_path": sbomPath,
		},
	}, nil
}

// UpdateDatabase обновляет базу данных уязвимостей
func (g *GrypeScanner) UpdateDatabase(ctx context.Context) error {
	if !g.IsAvailable() {
		return fmt.Errorf("grype is not available")
	}

	g.log.Info("Updating Grype vulnerability database")

	cmd := exec.CommandContext(ctx, "grype", "db", "update")
	executil.HideWindow(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("grype db update failed: %w, output: %s", err, string(output))
	}

	return nil
}

// GetDatabaseInfo возвращает информацию о базе данных
func (g *GrypeScanner) GetDatabaseInfo(ctx context.Context) (map[string]interface{}, error) {
	if !g.IsAvailable() {
		return nil, fmt.Errorf("grype is not available")
	}

	g.log.Info("Getting Grype database info")

	cmd := exec.CommandContext(ctx, "grype", "db", "status")
	executil.HideWindow(cmd)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("grype db status failed: %w", err)
	}

	return map[string]interface{}{
		"status": string(output),
		"tool":   "grype",
	}, nil
}

// getVersion возвращает версию Grype
func (g *GrypeScanner) getVersion() string {
	cmd := exec.Command("grype", "--version")
	executil.HideWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return versionUnknown
	}
	return string(output)
}

// parseGrypeOutput парсит вывод Grype
func (g *GrypeScanner) parseGrypeOutput(output []byte) ([]*domain.Vulnerability, error) {
	g.log.Info(fmt.Sprintf("Parsing Grype output, size: %d bytes", len(output)))

	if len(output) == 0 {
		return []*domain.Vulnerability{}, nil
	}

	var grypeResult grypeOutput
	if err := json.Unmarshal(output, &grypeResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grype output: %w", err)
	}

	vulnerabilities := make([]*domain.Vulnerability, 0, len(grypeResult.Matches))

	for _, match := range grypeResult.Matches {
		vuln := &domain.Vulnerability{
			ID:          match.Vulnerability.ID,
			Severity:    strings.ToUpper(match.Vulnerability.Severity),
			Description: match.Vulnerability.Description,
		}

		// Извлекаем CVSS score (берём первый доступный)
		if len(match.Vulnerability.CVSS) > 0 {
			vuln.CVSS = match.Vulnerability.CVSS[0].Score
		}

		// Извлекаем версию с исправлением
		if len(match.Vulnerability.Fix.Versions) > 0 {
			vuln.FixedIn = match.Vulnerability.Fix.Versions[0]
		}

		vulnerabilities = append(vulnerabilities, vuln)
	}

	g.log.Info(fmt.Sprintf("Parsed %d vulnerabilities from Grype output", len(vulnerabilities)))

	return vulnerabilities, nil
}

// calculateVulnerabilitySummary рассчитывает сводку уязвимостей
func (g *GrypeScanner) calculateVulnerabilitySummary(vulnerabilities []*domain.Vulnerability) *domain.VulnerabilitySummary {
	summary := &domain.VulnerabilitySummary{
		Total: len(vulnerabilities),
	}

	for _, vuln := range vulnerabilities {
		switch vuln.Severity {
		case "CRITICAL":
			summary.Critical++
		case "HIGH":
			summary.High++
		case "MEDIUM":
			summary.Medium++
		case "LOW":
			summary.Low++
		}
	}

	return summary
}
