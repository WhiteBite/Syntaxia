package domain

import (
	"context"
	"time"
)

// SBOMFormat определяет формат SBOM
type SBOMFormat string

const (
	SBOMFormatSPDX      SBOMFormat = "spdx"
	SBOMFormatCycloneDX SBOMFormat = "cyclonedx"
	SBOMFormatJSON      SBOMFormat = "json"
)

// SBOMResult представляет результат генерации SBOM
type SBOMResult struct {
	Success     bool                   `json:"success"`
	ProjectPath string                 `json:"projectPath"`
	Format      SBOMFormat             `json:"format"`
	OutputPath  string                 `json:"outputPath"`
	Components  []*SBOMComponent       `json:"components"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// SBOMComponent представляет компонент в SBOM
type SBOMComponent struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Type            string            `json:"type"`
	PURL            string            `json:"purl,omitempty"`
	License         string            `json:"license,omitempty"`
	Vulnerabilities []*Vulnerability  `json:"vulnerabilities,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// Vulnerability представляет уязвимость
type Vulnerability struct {
	ID          string  `json:"id"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	CVSS        float64 `json:"cvss,omitempty"`
	FixedIn     string  `json:"fixedIn,omitempty"`
}

// VulnerabilityScanResult представляет результат сканирования уязвимостей
type VulnerabilityScanResult struct {
	Success         bool                   `json:"success"`
	ProjectPath     string                 `json:"projectPath"`
	Vulnerabilities []*Vulnerability       `json:"vulnerabilities"`
	Summary         *VulnerabilitySummary  `json:"summary"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Error           string                 `json:"error,omitempty"`
}

// VulnerabilitySummary представляет сводку уязвимостей
type VulnerabilitySummary struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Fixed    int `json:"fixed"`
}

// LicenseScanResult представляет результат сканирования лицензий
type LicenseScanResult struct {
	Success     bool                   `json:"success"`
	ProjectPath string                 `json:"projectPath"`
	Licenses    []*LicenseInfo         `json:"licenses"`
	Summary     *LicenseSummary        `json:"summary"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// LicenseInfo представляет информацию о лицензии
type LicenseInfo struct {
	Name        string   `json:"name"`
	SPDXID      string   `json:"spdxId"`
	Type        string   `json:"type"` // "permissive", "copyleft", "proprietary"
	Files       []string `json:"files"`
	Confidence  float64  `json:"confidence"`
	Description string   `json:"description,omitempty"`
}

// LicenseSummary представляет сводку лицензий
type LicenseSummary struct {
	TotalLicenses int                `json:"totalLicenses"`
	ByType        map[string]int     `json:"byType"`
	ByLicense     map[string]int     `json:"byLicense"`
	Conflicts     []*LicenseConflict `json:"conflicts,omitempty"`
}

// LicenseConflict представляет конфликт лицензий
type LicenseConflict struct {
	License1    string `json:"license1"`
	License2    string `json:"license2"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// ComplianceRequirements представляет требования соответствия
type ComplianceRequirements struct {
	AllowedLicenses            []string `json:"allowedLicenses"`
	ForbiddenLicenses          []string `json:"forbiddenLicenses"`
	MaxVulnerabilities         int      `json:"maxVulnerabilities"`
	MaxCriticalVulnerabilities int      `json:"maxCriticalVulnerabilities"`
	MaxHighVulnerabilities     int      `json:"maxHighVulnerabilities"`
	MaxCVSS                    float64  `json:"maxCVSS"`
	RequireSBOM                bool     `json:"requireSBOM"`
}

// ComplianceReport представляет отчет о соответствии
type ComplianceReport struct {
	Success             bool                     `json:"success"`
	ProjectPath         string                   `json:"projectPath"`
	Compliant           bool                     `json:"compliant"`
	Issues              []*ComplianceIssue       `json:"issues"`
	Summary             *ComplianceSummary       `json:"summary"`
	Metadata            map[string]interface{}   `json:"metadata,omitempty"`
	GeneratedAt         time.Time                `json:"generatedAt"`
	SBOMResult          *SBOMResult              `json:"sbomResult,omitempty"`
	VulnerabilityResult *VulnerabilityScanResult `json:"vulnerabilityResult,omitempty"`
	LicenseResult       *LicenseScanResult       `json:"licenseResult,omitempty"`
}

// ComplianceIssue представляет проблему соответствия
type ComplianceIssue struct {
	Type           string `json:"type"` // "license", "vulnerability", "sbom"
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Component      string `json:"component,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

// ComplianceSummary представляет сводку соответствия
type ComplianceSummary struct {
	TotalIssues int `json:"totalIssues"`
	Critical    int `json:"critical"`
	High        int `json:"high"`
	Medium      int `json:"medium"`
	Low         int `json:"low"`
}

// SBOMService определяет интерфейс для работы с SBOM и лицензиями
type SBOMService interface {
	// GenerateSBOM генерирует SBOM для проекта
	GenerateSBOM(ctx context.Context, projectPath string, format SBOMFormat) (*SBOMResult, error)

	// ScanVulnerabilities сканирует уязвимости в проекте
	ScanVulnerabilities(ctx context.Context, projectPath string) (*VulnerabilityScanResult, error)

	// ScanLicenses сканирует лицензии в проекте
	ScanLicenses(ctx context.Context, projectPath string) (*LicenseScanResult, error)

	// GenerateComplianceReport генерирует отчет о соответствии
	GenerateComplianceReport(ctx context.Context, projectPath string, requirements *ComplianceRequirements) (*ComplianceReport, error)

	// GetSupportedSBOMFormats возвращает поддерживаемые форматы SBOM
	GetSupportedSBOMFormats() []SBOMFormat

	// ValidateSBOM валидирует SBOM
	ValidateSBOM(ctx context.Context, sbomPath string, format SBOMFormat) error
}
