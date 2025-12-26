package sbom

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"strings"
)

// Service предоставляет высокоуровневый API для работы с SBOM и лицензиями
type Service struct {
	log              domain.Logger
	sbomGenerator    domain.SBOMGenerator
	vulnScanner      domain.VulnerabilityScanner
	licenseScanner   domain.LicenseScanner
	fileStatProvider domain.FileStatProvider
}

// NewService создает новый сервис SBOM
func NewService(log domain.Logger, sbomGenerator domain.SBOMGenerator, vulnScanner domain.VulnerabilityScanner, licenseScanner domain.LicenseScanner, fileStatProvider domain.FileStatProvider) *Service {
	return &Service{
		log:              log,
		sbomGenerator:    sbomGenerator,
		vulnScanner:      vulnScanner,
		licenseScanner:   licenseScanner,
		fileStatProvider: fileStatProvider,
	}
}

// GenerateSBOM генерирует SBOM для проекта
func (s *Service) GenerateSBOM(ctx context.Context, projectPath string, format domain.SBOMFormat) (*domain.SBOMResult, error) {
	s.log.Info(fmt.Sprintf("Generating SBOM for project: %s, format: %s", projectPath, format))

	// Проверяем существование проекта
	if _, err := s.fileStatProvider.Stat(projectPath); err != nil {
		return &domain.SBOMResult{
			Success:     false,
			ProjectPath: projectPath,
			Format:      format,
			Error:       fmt.Sprintf("project path does not exist: %s", projectPath),
		}, nil
	}

	// Проверяем доступность Syft
	if !s.sbomGenerator.IsAvailable() {
		return &domain.SBOMResult{
			Success:     false,
			ProjectPath: projectPath,
			Format:      format,
			Error:       "SBOM generator is not available",
		}, nil
	}

	// Генерируем SBOM с помощью Syft
	result, err := s.sbomGenerator.GenerateSBOM(ctx, projectPath, format)
	if err != nil {
		return &domain.SBOMResult{
			Success:     false,
			ProjectPath: projectPath,
			Format:      format,
			Error:       err.Error(),
		}, nil
	}

	return result, nil
}

// checkProjectAndScanner validates project exists and scanner is available
func (s *Service) checkProjectAndScanner(projectPath string, isAvailable func() bool, scannerName string) (bool, string) {
	if _, err := s.fileStatProvider.Stat(projectPath); err != nil {
		return false, fmt.Sprintf("project path does not exist: %s", projectPath)
	}
	if !isAvailable() {
		return false, scannerName + " is not available"
	}
	return true, ""
}

// ScanVulnerabilities сканирует уязвимости в проекте
func (s *Service) ScanVulnerabilities(ctx context.Context, projectPath string) (*domain.VulnerabilityScanResult, error) {
	s.log.Info(fmt.Sprintf("Scanning vulnerabilities for project: %s", projectPath))

	if ok, errMsg := s.checkProjectAndScanner(projectPath, s.vulnScanner.IsAvailable, "Vulnerability scanner"); !ok {
		return &domain.VulnerabilityScanResult{Success: false, ProjectPath: projectPath, Error: errMsg}, nil
	}

	result, err := s.vulnScanner.ScanVulnerabilities(ctx, projectPath)
	if err != nil {
		return &domain.VulnerabilityScanResult{Success: false, ProjectPath: projectPath, Error: err.Error()}, nil
	}
	return result, nil
}

// ScanLicenses сканирует лицензии в проекте
func (s *Service) ScanLicenses(ctx context.Context, projectPath string) (*domain.LicenseScanResult, error) {
	s.log.Info(fmt.Sprintf("Scanning licenses for project: %s", projectPath))

	if ok, errMsg := s.checkProjectAndScanner(projectPath, s.licenseScanner.IsAvailable, "License scanner"); !ok {
		return &domain.LicenseScanResult{Success: false, ProjectPath: projectPath, Error: errMsg}, nil
	}

	result, err := s.licenseScanner.ScanLicenses(ctx, projectPath)
	if err != nil {
		return &domain.LicenseScanResult{Success: false, ProjectPath: projectPath, Error: err.Error()}, nil
	}
	return result, nil
}

// GenerateComplianceReport генерирует отчет о соответствии
func (s *Service) GenerateComplianceReport(ctx context.Context, projectPath string, requirements *domain.ComplianceRequirements) (*domain.ComplianceReport, error) {
	s.log.Info(fmt.Sprintf("Generating compliance report for project: %s", projectPath))

	// Проверяем существование проекта
	if _, err := s.fileStatProvider.Stat(projectPath); err != nil {
		return &domain.ComplianceReport{
			Success:     false,
			ProjectPath: projectPath,
			Compliant:   false,
		}, nil
	}

	// Генерируем SBOM
	sbomResult, err := s.GenerateSBOM(ctx, projectPath, domain.SBOMFormatSPDX)
	if err != nil {
		return &domain.ComplianceReport{
			Success:     false,
			ProjectPath: projectPath,
			Compliant:   false,
		}, err
	}

	// Сканируем уязвимости
	vulnResult, err := s.ScanVulnerabilities(ctx, projectPath)
	if err != nil {
		return &domain.ComplianceReport{
			Success:     false,
			ProjectPath: projectPath,
			Compliant:   false,
		}, err
	}

	// Сканируем лицензии
	licenseResult, err := s.ScanLicenses(ctx, projectPath)
	if err != nil {
		return &domain.ComplianceReport{
			Success:     false,
			ProjectPath: projectPath,
			Compliant:   false,
		}, err
	}

	// Анализируем соответствие
	report := s.analyzeCompliance(projectPath, requirements, sbomResult, vulnResult, licenseResult)

	return report, nil
}

// GetSupportedSBOMFormats возвращает поддерживаемые форматы SBOM
func (s *Service) GetSupportedSBOMFormats() []domain.SBOMFormat {
	return []domain.SBOMFormat{
		domain.SBOMFormatSPDX,
		domain.SBOMFormatCycloneDX,
		domain.SBOMFormatJSON,
	}
}

// ValidateSBOM валидирует SBOM
func (s *Service) ValidateSBOM(ctx context.Context, sbomPath string, format domain.SBOMFormat) error {
	s.log.Info(fmt.Sprintf("Validating SBOM: %s, format: %s", sbomPath, format))

	// Проверяем существование файла
	if _, err := s.fileStatProvider.Stat(sbomPath); err != nil {
		return fmt.Errorf("SBOM file does not exist: %s", sbomPath)
	}

	// Проверяем доступность Syft
	if !s.sbomGenerator.IsAvailable() {
		return fmt.Errorf("SBOM generator is not available for validation")
	}

	// Валидируем SBOM с помощью Syft
	return s.sbomGenerator.ValidateSBOM(ctx, sbomPath, format)
}

// checkSBOMCompliance checks SBOM requirements
func checkSBOMCompliance(requirements *domain.ComplianceRequirements, sbomResult *domain.SBOMResult) []*domain.ComplianceIssue {
	if requirements.RequireSBOM && !sbomResult.Success {
		return []*domain.ComplianceIssue{{Type: "sbom", Severity: "high", Description: "SBOM generation failed"}}
	}
	return nil
}

// checkVulnerabilityCompliance checks vulnerability requirements
func checkVulnerabilityCompliance(requirements *domain.ComplianceRequirements, vulnResult *domain.VulnerabilityScanResult) []*domain.ComplianceIssue {
	if !vulnResult.Success || vulnResult.Summary == nil {
		return nil
	}
	var issues []*domain.ComplianceIssue
	if vulnResult.Summary.Total > requirements.MaxVulnerabilities {
		issues = append(issues, &domain.ComplianceIssue{
			Type: "vulnerability", Severity: "high",
			Description: fmt.Sprintf("Too many vulnerabilities: %d (max: %d)", vulnResult.Summary.Total, requirements.MaxVulnerabilities),
		})
	}
	if vulnResult.Summary.Critical > 0 {
		issues = append(issues, &domain.ComplianceIssue{
			Type: "vulnerability", Severity: "critical",
			Description: fmt.Sprintf("Critical vulnerabilities found: %d", vulnResult.Summary.Critical),
		})
	}
	return issues
}

// checkLicenseCompliance checks license requirements
func checkLicenseCompliance(requirements *domain.ComplianceRequirements, licenseResult *domain.LicenseScanResult) []*domain.ComplianceIssue {
	if !licenseResult.Success || licenseResult.Summary == nil {
		return nil
	}
	var issues []*domain.ComplianceIssue
	for _, license := range licenseResult.Licenses {
		for _, forbidden := range requirements.ForbiddenLicenses {
			if strings.Contains(strings.ToLower(license.Name), strings.ToLower(forbidden)) {
				issues = append(issues, &domain.ComplianceIssue{
					Type: "license", Severity: "high",
					Description: fmt.Sprintf("Forbidden license found: %s", license.Name), Component: license.Name,
				})
			}
		}
	}
	return issues
}

// analyzeCompliance анализирует соответствие требованиям
func (s *Service) analyzeCompliance(projectPath string, requirements *domain.ComplianceRequirements, sbomResult *domain.SBOMResult, vulnResult *domain.VulnerabilityScanResult, licenseResult *domain.LicenseScanResult) *domain.ComplianceReport {
	var issues []*domain.ComplianceIssue
	issues = append(issues, checkSBOMCompliance(requirements, sbomResult)...)
	issues = append(issues, checkVulnerabilityCompliance(requirements, vulnResult)...)
	issues = append(issues, checkLicenseCompliance(requirements, licenseResult)...)

	// Рассчитываем сводку
	summary := &domain.ComplianceSummary{
		TotalIssues: len(issues),
	}

	for _, issue := range issues {
		switch issue.Severity {
		case "critical":
			summary.Critical++
		case "high":
			summary.High++
		case "medium":
			summary.Medium++
		case "low":
			summary.Low++
		}
	}

	// Определяем общее соответствие
	compliant := len(issues) == 0

	return &domain.ComplianceReport{
		Success:     true,
		ProjectPath: projectPath,
		Compliant:   compliant,
		Issues:      issues,
		Summary:     summary,
	}
}
