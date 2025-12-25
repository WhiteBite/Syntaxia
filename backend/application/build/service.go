package build

import (
	"context"
	"fmt"

	"syntaxia/domain"
)

// Language constants
const (
	langGo         = "go"
	langTypeScript = "typescript"
	langTS         = "ts"
	langJava       = "java"
)

// Service предоставляет высокоуровневый API для работы с build pipeline
type Service struct {
	log      domain.Logger
	pipeline domain.BuildPipeline
}

// NewService создает новый сервис сборки
func NewService(log domain.Logger, pipeline domain.BuildPipeline) *Service {
	return &Service{
		log:      log,
		pipeline: pipeline,
	}
}

// Build выполняет сборку проекта
func (s *Service) Build(ctx context.Context, projectPath, language string) (*domain.BuildResult, error) {
	s.log.Info(fmt.Sprintf("Building %s project at %s", language, projectPath))
	return s.pipeline.Build(ctx, projectPath, language)
}

// TypeCheck выполняет проверку типов
func (s *Service) TypeCheck(ctx context.Context, projectPath, language string) (*domain.TypeCheckResult, error) {
	s.log.Info(fmt.Sprintf("Type checking %s project at %s", language, projectPath))
	return s.pipeline.TypeCheck(ctx, projectPath, language)
}

// BuildAndTypeCheck выполняет сборку и проверку типов
func (s *Service) BuildAndTypeCheck(ctx context.Context, projectPath, language string) (*domain.BuildResult, *domain.TypeCheckResult, error) {
	s.log.Info(fmt.Sprintf("Building and type checking %s project at %s", language, projectPath))
	return s.pipeline.BuildAndTypeCheck(ctx, projectPath, language)
}

// BuildMultiLanguage выполняет сборку для нескольких языков
func (s *Service) BuildMultiLanguage(ctx context.Context, projectPath string, languages []string) (map[string]*domain.BuildResult, error) {
	s.log.Info(fmt.Sprintf("Building multi-language project at %s for languages: %v", projectPath, languages))

	results := make(map[string]*domain.BuildResult)
	for _, language := range languages {
		result, err := s.Build(ctx, projectPath, language)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to build %s: %v", language, err))
			continue
		}
		results[language] = result
	}
	return results, nil
}

// TypeCheckMultiLanguage выполняет проверку типов для нескольких языков
func (s *Service) TypeCheckMultiLanguage(ctx context.Context, projectPath string, languages []string) (map[string]*domain.TypeCheckResult, error) {
	s.log.Info(fmt.Sprintf("Type checking multi-language project at %s for languages: %v", projectPath, languages))

	results := make(map[string]*domain.TypeCheckResult)
	for _, language := range languages {
		result, err := s.TypeCheck(ctx, projectPath, language)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to type check %s: %v", language, err))
			continue
		}
		results[language] = result
	}
	return results, nil
}

// ValidateProject выполняет полную валидацию проекта
func (s *Service) ValidateProject(ctx context.Context, projectPath string, languages []string) (*domain.ProjectValidationResult, error) {
	s.log.Info(fmt.Sprintf("Validating project at %s for languages: %v", projectPath, languages))

	validation := &domain.ProjectValidationResult{
		ProjectPath: projectPath,
		Languages:   languages,
		Results:     make(map[string]*domain.LanguageValidationResult),
	}

	for _, language := range languages {
		langResult := &domain.LanguageValidationResult{Language: language}

		typeCheckResult, err := s.TypeCheck(ctx, projectPath, language)
		if err != nil {
			langResult.TypeCheckError = err.Error()
		} else {
			langResult.TypeCheckResult = typeCheckResult
		}

		buildResult, err := s.Build(ctx, projectPath, language)
		if err != nil {
			langResult.BuildError = err.Error()
		} else {
			langResult.BuildResult = buildResult
		}

		langResult.Success = (langResult.TypeCheckResult != nil && langResult.TypeCheckResult.Success) &&
			(langResult.BuildResult != nil && langResult.BuildResult.Success)
		validation.Results[language] = langResult
	}

	validation.Success = true
	for _, result := range validation.Results {
		if !result.Success {
			validation.Success = false
			break
		}
	}
	return validation, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (s *Service) GetSupportedLanguages() []string {
	return s.pipeline.GetSupportedLanguages()
}

// DetectLanguages определяет языки в проекте
func (s *Service) DetectLanguages(ctx context.Context, projectPath string) ([]string, error) {
	s.log.Info(fmt.Sprintf("Detecting languages in project at %s", projectPath))

	var detectedLanguages []string
	supportedLanguages := s.GetSupportedLanguages()

	for _, language := range supportedLanguages {
		switch language {
		case langGo:
			if s.hasFile(projectPath, "go.mod") {
				detectedLanguages = append(detectedLanguages, language)
			}
		case langTypeScript, langTS:
			if s.hasFile(projectPath, "package.json") || s.hasFile(projectPath, "tsconfig.json") {
				detectedLanguages = append(detectedLanguages, language)
			}
		case langJava:
			if s.hasFile(projectPath, "pom.xml") || s.hasFile(projectPath, "build.gradle") {
				detectedLanguages = append(detectedLanguages, language)
			}
		}
	}

	s.log.Info(fmt.Sprintf("Detected languages: %v", detectedLanguages))
	return detectedLanguages, nil
}

func (s *Service) hasFile(projectPath, filename string) bool {
	return true // Placeholder - real implementation should use os.Stat
}
