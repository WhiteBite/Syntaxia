package sbomlicensing

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"syntaxia/domain"
	"syntaxia/internal/executil"
)

// syftJSONOutput представляет структуру JSON вывода Syft
type syftJSONOutput struct {
	Artifacts []syftArtifact `json:"artifacts"`
}

// syftArtifact представляет артефакт в выводе Syft
type syftArtifact struct {
	Name      string        `json:"name"`
	Version   string        `json:"version"`
	Type      string        `json:"type"`
	PURL      string        `json:"purl"`
	Licenses  []syftLicense `json:"licenses"`
	Locations []struct {
		Path string `json:"path"`
	} `json:"locations"`
	Metadata map[string]interface{} `json:"metadata"`
}

// syftLicense представляет лицензию в выводе Syft
type syftLicense struct {
	Value  string `json:"value"`
	SPDXID string `json:"spdxExpression"`
	Type   string `json:"type"`
}

// spdxDocument представляет SPDX документ
type spdxDocument struct {
	Packages []spdxPackage `json:"packages"`
}

// spdxPackage представляет пакет в SPDX формате
type spdxPackage struct {
	Name             string `json:"name"`
	VersionInfo      string `json:"versionInfo"`
	SPDXID           string `json:"SPDXID"`
	LicenseConcluded string `json:"licenseConcluded"`
	LicenseDeclared  string `json:"licenseDeclared"`
	ExternalRefs     []struct {
		ReferenceType     string `json:"referenceType"`
		ReferenceLocator  string `json:"referenceLocator"`
		ReferenceCategory string `json:"referenceCategory"`
	} `json:"externalRefs"`
}

// SyftGenerator представляет генератор SBOM с использованием Syft
type SyftGenerator struct {
	log domain.Logger
}

// NewSyftGenerator создает новый генератор Syft
func NewSyftGenerator(log domain.Logger) *SyftGenerator {
	return &SyftGenerator{
		log: log,
	}
}

// IsAvailable проверяет доступность Syft
func (s *SyftGenerator) IsAvailable() bool {
	_, err := exec.LookPath("syft")
	return err == nil
}

// GenerateSBOM генерирует SBOM для проекта
func (s *SyftGenerator) GenerateSBOM(ctx context.Context, projectPath string, format domain.SBOMFormat) (*domain.SBOMResult, error) {
	if !s.IsAvailable() {
		return &domain.SBOMResult{
			Success:     false,
			ProjectPath: projectPath,
			Format:      format,
			Error:       "syft is not available",
		}, nil
	}

	s.log.Info(fmt.Sprintf("Generating SBOM with Syft for: %s, format: %s", projectPath, format))

	// Определяем расширение файла
	var fileExt string
	switch format {
	case domain.SBOMFormatSPDX:
		fileExt = "spdx"
	case domain.SBOMFormatCycloneDX:
		fileExt = "cyclonedx"
	case domain.SBOMFormatJSON:
		fileExt = "json"
	default:
		fileExt = "json"
	}

	outputPath := filepath.Join(projectPath, fmt.Sprintf("sbom.%s", fileExt))

	// Запускаем Syft
	cmd := exec.CommandContext(ctx, "syft", projectPath, "-o", string(format), "-f", outputPath) //nolint:gosec // External tool command
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &domain.SBOMResult{
			Success:     false,
			ProjectPath: projectPath,
			Format:      format,
			Error:       fmt.Sprintf("syft failed: %v, output: %s", err, string(output)),
		}, nil
	}

	// Парсим результат
	components, err := s.parseSyftOutput(output, format)
	if err != nil {
		return &domain.SBOMResult{
			Success:     false,
			ProjectPath: projectPath,
			Format:      format,
			Error:       fmt.Sprintf("failed to parse syft output: %v", err),
		}, nil
	}

	return &domain.SBOMResult{
		Success:     true,
		ProjectPath: projectPath,
		Format:      format,
		OutputPath:  outputPath,
		Components:  components,
		Metadata: map[string]interface{}{
			"tool":    "syft",
			"version": s.getVersion(),
			"command": cmd.String(),
		},
	}, nil
}

// ValidateSBOM валидирует SBOM файл
func (s *SyftGenerator) ValidateSBOM(ctx context.Context, sbomPath string, format domain.SBOMFormat) error {
	if !s.IsAvailable() {
		return fmt.Errorf("syft is not available")
	}

	s.log.Info(fmt.Sprintf("Validating SBOM with Syft: %s", sbomPath))

	// Запускаем Syft для валидации
	cmd := exec.CommandContext(ctx, "syft", "validate", sbomPath)
	executil.HideWindow(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("syft validation failed: %w, output: %s", err, string(output))
	}

	return nil
}

// GetSupportedFormats возвращает поддерживаемые форматы
func (s *SyftGenerator) GetSupportedFormats() []domain.SBOMFormat {
	return []domain.SBOMFormat{
		domain.SBOMFormatSPDX,
		domain.SBOMFormatCycloneDX,
		domain.SBOMFormatJSON,
	}
}

// getVersion возвращает версию Syft
func (s *SyftGenerator) getVersion() string {
	cmd := exec.Command("syft", "--version")
	executil.HideWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return string(output)
}

// parseSyftOutput парсит вывод Syft
func (s *SyftGenerator) parseSyftOutput(output []byte, format domain.SBOMFormat) ([]*domain.SBOMComponent, error) {
	s.log.Info(fmt.Sprintf("Parsing Syft output, format: %s, size: %d bytes", format, len(output)))

	if len(output) == 0 {
		return []*domain.SBOMComponent{}, nil
	}

	switch format {
	case domain.SBOMFormatJSON:
		return s.parseSyftJSON(output)
	case domain.SBOMFormatSPDX:
		return s.parseSyftSPDX(output)
	case domain.SBOMFormatCycloneDX:
		return s.parseSyftCycloneDX(output)
	default:
		return s.parseSyftJSON(output)
	}
}

// parseSyftJSON парсит JSON формат Syft
func (s *SyftGenerator) parseSyftJSON(output []byte) ([]*domain.SBOMComponent, error) {
	var syftResult syftJSONOutput
	if err := json.Unmarshal(output, &syftResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal syft JSON output: %w", err)
	}

	components := make([]*domain.SBOMComponent, 0, len(syftResult.Artifacts))

	for _, artifact := range syftResult.Artifacts {
		component := &domain.SBOMComponent{
			Name:    artifact.Name,
			Version: artifact.Version,
			Type:    artifact.Type,
			PURL:    artifact.PURL,
		}

		// Извлекаем лицензию (берём первую доступную)
		if len(artifact.Licenses) > 0 {
			licenses := make([]string, 0, len(artifact.Licenses))
			for _, lic := range artifact.Licenses {
				if lic.SPDXID != "" {
					licenses = append(licenses, lic.SPDXID)
				} else if lic.Value != "" {
					licenses = append(licenses, lic.Value)
				}
			}
			component.License = strings.Join(licenses, ", ")
		}

		// Добавляем метаданные
		if artifact.Metadata != nil {
			component.Metadata = make(map[string]string)
			for k, v := range artifact.Metadata {
				if str, ok := v.(string); ok {
					component.Metadata[k] = str
				}
			}
		}

		components = append(components, component)
	}

	s.log.Info(fmt.Sprintf("Parsed %d components from Syft JSON output", len(components)))

	return components, nil
}

// parseSyftSPDX парсит SPDX формат
func (s *SyftGenerator) parseSyftSPDX(output []byte) ([]*domain.SBOMComponent, error) {
	var spdxDoc spdxDocument
	if err := json.Unmarshal(output, &spdxDoc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SPDX output: %w", err)
	}

	components := make([]*domain.SBOMComponent, 0, len(spdxDoc.Packages))

	for _, pkg := range spdxDoc.Packages {
		component := &domain.SBOMComponent{
			Name:    pkg.Name,
			Version: pkg.VersionInfo,
			Type:    "package",
		}

		// Определяем лицензию
		if pkg.LicenseConcluded != "" && pkg.LicenseConcluded != "NOASSERTION" {
			component.License = pkg.LicenseConcluded
		} else if pkg.LicenseDeclared != "" && pkg.LicenseDeclared != "NOASSERTION" {
			component.License = pkg.LicenseDeclared
		}

		// Извлекаем PURL из external refs
		for _, ref := range pkg.ExternalRefs {
			if ref.ReferenceType == "purl" {
				component.PURL = ref.ReferenceLocator
				break
			}
		}

		components = append(components, component)
	}

	s.log.Info(fmt.Sprintf("Parsed %d components from SPDX output", len(components)))

	return components, nil
}

// parseSyftCycloneDX парсит CycloneDX формат
func (s *SyftGenerator) parseSyftCycloneDX(output []byte) ([]*domain.SBOMComponent, error) {
	// CycloneDX JSON структура
	var cycloneDX struct {
		Components []struct {
			Name     string `json:"name"`
			Version  string `json:"version"`
			Type     string `json:"type"`
			PURL     string `json:"purl"`
			Licenses []struct {
				License struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"license"`
			} `json:"licenses"`
		} `json:"components"`
	}

	if err := json.Unmarshal(output, &cycloneDX); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CycloneDX output: %w", err)
	}

	components := make([]*domain.SBOMComponent, 0, len(cycloneDX.Components))

	for _, comp := range cycloneDX.Components {
		component := &domain.SBOMComponent{
			Name:    comp.Name,
			Version: comp.Version,
			Type:    comp.Type,
			PURL:    comp.PURL,
		}

		// Извлекаем лицензии
		if len(comp.Licenses) > 0 {
			licenses := make([]string, 0, len(comp.Licenses))
			for _, lic := range comp.Licenses {
				if lic.License.ID != "" {
					licenses = append(licenses, lic.License.ID)
				} else if lic.License.Name != "" {
					licenses = append(licenses, lic.License.Name)
				}
			}
			component.License = strings.Join(licenses, ", ")
		}

		components = append(components, component)
	}

	s.log.Info(fmt.Sprintf("Parsed %d components from CycloneDX output", len(components)))

	return components, nil
}
