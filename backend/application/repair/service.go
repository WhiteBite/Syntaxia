package repair

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"syntaxia/domain"
	"sort"
	"strings"
	"time"
)

// Service реализует RepairService
type Service struct {
	log           domain.Logger
	commandRunner domain.CommandRunner
}

// NewService создает новый сервис repair
func NewService(log domain.Logger, commandRunner domain.CommandRunner) domain.RepairService {
	return &Service{
		log:           log,
		commandRunner: commandRunner,
	}
}

// ExecuteRepair выполняет repair цикл
func (s *Service) ExecuteRepair(ctx context.Context, req domain.RepairRequest) (*domain.RepairResult, error) {
	startTime := time.Now()
	s.log.Info(fmt.Sprintf("Starting repair cycle for project: %s", req.ProjectPath))

	// Проверяем существование проекта
	if _, err := os.Stat(req.ProjectPath); os.IsNotExist(err) {
		return &domain.RepairResult{
			Success: false,
			Error:   fmt.Sprintf("project path does not exist: %s", req.ProjectPath),
		}, nil
	}

	// Если правила не указаны, получаем доступные для языка
	if len(req.Rules) == 0 {
		rules, err := s.GetAvailableRules(req.Language)
		if err != nil {
			return &domain.RepairResult{
				Success: false,
				Error:   fmt.Sprintf("failed to get rules: %v", err),
			}, err
		}
		req.Rules = rules
	}

	// Сортируем правила по приоритету
	sort.Slice(req.Rules, func(i, j int) bool {
		return req.Rules[i].Priority > req.Rules[j].Priority
	})

	// Выполняем repair цикл
	for attempt := 1; attempt <= req.MaxAttempts; attempt++ {
		s.log.Info(fmt.Sprintf("Repair attempt %d/%d", attempt, req.MaxAttempts))

		// Анализируем ошибки и применяем правила
		appliedRules := s.applyRepairRules(ctx, req.ProjectPath, req.ErrorOutput, req.Rules)

		if len(appliedRules) == 0 {
			s.log.Info("No applicable repair rules found")
			break
		}

		// Проверяем, исправились ли ошибки
		success, newErrors := s.verifyRepair(ctx, req.ProjectPath, req.Language)
		if success {
			duration := time.Since(startTime)
			return &domain.RepairResult{
				Success:    true,
				FixedFiles: appliedRules,
				Duration:   duration,
				Attempts:   attempt,
			}, nil
		}

		// Обновляем ошибки для следующей попытки
		req.ErrorOutput = newErrors
	}

	duration := time.Since(startTime)
	return &domain.RepairResult{
		Success:  false,
		Error:    "repair cycle completed but errors remain",
		Duration: duration,
		Attempts: req.MaxAttempts,
	}, nil
}

// GetAvailableRules возвращает доступные правила для языка
func (s *Service) GetAvailableRules(language string) ([]domain.RepairRule, error) {
	rules := s.getDefaultRules(language)
	return rules, nil
}

// AddRule добавляет новое правило
func (s *Service) AddRule(rule domain.RepairRule) error {
	// В простой реализации просто логируем
	s.log.Info(fmt.Sprintf("Adding repair rule: %s", rule.Name))
	return nil
}

// RemoveRule удаляет правило
func (s *Service) RemoveRule(ruleID string) error {
	s.log.Info(fmt.Sprintf("Removing repair rule: %s", ruleID))
	return nil
}

// ValidateRule проверяет корректность правила
func (s *Service) ValidateRule(rule domain.RepairRule) error {
	if rule.ID == "" {
		return fmt.Errorf("rule ID is required")
	}
	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}
	if rule.Pattern == "" {
		return fmt.Errorf("rule pattern is required")
	}

	// Проверяем, что pattern является валидным regex
	if _, err := regexp.Compile(rule.Pattern); err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	return nil
}

// applyRepairRules применяет правила к проекту
func (s *Service) applyRepairRules(ctx context.Context, projectPath, errorOutput string, rules []domain.RepairRule) []string {
	var fixedFiles []string

	for _, rule := range rules {
		// Проверяем, подходит ли правило к ошибкам
		if !s.matchesError(errorOutput, rule) {
			continue
		}

		// Применяем правило
		files, err := s.applyRule(ctx, projectPath, rule)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to apply rule %s: %v", rule.Name, err))
			continue
		}

		fixedFiles = append(fixedFiles, files...)
		s.log.Info(fmt.Sprintf("Applied rule %s to %d files", rule.Name, len(files)))
	}

	return fixedFiles
}

// matchesError проверяет, подходит ли правило к ошибкам
func (s *Service) matchesError(errorOutput string, rule domain.RepairRule) bool {
	re, err := regexp.Compile(rule.Pattern)
	if err != nil {
		s.log.Warning(fmt.Sprintf("Invalid regex in rule %s: %v", rule.Name, err))
		return false
	}
	return re.MatchString(errorOutput)
}

// applyRule применяет конкретное правило
func (s *Service) applyRule(ctx context.Context, projectPath string, rule domain.RepairRule) ([]string, error) {
	var fixedFiles []string

	switch rule.Category {
	case "format":
		fixedFiles = s.applyFormatRule(ctx, projectPath, rule)
	case "import":
		fixedFiles = s.applyImportRule(ctx, projectPath, rule)
	case "syntax":
		fixedFiles = s.applySyntaxRule(ctx, projectPath, rule)
	default:
		s.log.Warning(fmt.Sprintf("Unknown rule category: %s", rule.Category))
	}

	return fixedFiles, nil
}

// applyFormatRule применяет правило форматирования
func (s *Service) applyFormatRule(ctx context.Context, projectPath string, rule domain.RepairRule) []string {
	var fixedFiles []string

	// Определяем язык и применяем соответствующий форматтер
	if strings.Contains(rule.Language, "go") {
		// gofmt
		_, err := s.commandRunner.RunCommandInDir(ctx, projectPath, "gofmt", "-w", ".")
		if err == nil {
			fixedFiles = append(fixedFiles, "*.go")
		}

		// goimports
		_, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "goimports", "-w", ".")
		if err == nil {
			fixedFiles = append(fixedFiles, "*.go")
		}
	} else if strings.Contains(rule.Language, "typescript") || strings.Contains(rule.Language, "javascript") {
		// prettier
		_, err := s.commandRunner.RunCommandInDir(ctx, projectPath, "npx", "prettier", "--write", ".")
		if err == nil {
			fixedFiles = append(fixedFiles, "*.ts", "*.js", "*.vue")
		}
	}

	return fixedFiles
}

// applyImportRule применяет правило импортов
func (s *Service) applyImportRule(ctx context.Context, projectPath string, rule domain.RepairRule) []string {
	var fixedFiles []string

	if strings.Contains(rule.Language, "go") {
		// go mod tidy
		_, err := s.commandRunner.RunCommandInDir(ctx, projectPath, "go", "mod", "tidy")
		if err == nil {
			fixedFiles = append(fixedFiles, "go.mod", "go.sum")
		}
	} else if strings.Contains(rule.Language, "typescript") || strings.Contains(rule.Language, "javascript") {
		// npm install
		_, err := s.commandRunner.RunCommandInDir(ctx, projectPath, "npm", "install")
		if err == nil {
			fixedFiles = append(fixedFiles, "package.json", "package-lock.json")
		}
	}

	return fixedFiles
}

// applySyntaxRule применяет правило синтаксиса
func (s *Service) applySyntaxRule(ctx context.Context, projectPath string, rule domain.RepairRule) []string {
	// В простой реализации возвращаем пустой список
	// В реальной реализации здесь была бы логика исправления синтаксических ошибок
	return []string{}
}

// verifyRepair проверяет, исправились ли ошибки
func (s *Service) verifyRepair(ctx context.Context, projectPath, language string) (bool, string) {
	var output []byte
	var err error

	switch language {
	case langGo:
		output, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "go", "build", "./...")
	case langTypeScript, langJavaScript:
		output, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "npx", "tsc", "--noEmit")
	default:
		return false, "unsupported language for verification"
	}

	if err == nil {
		return true, ""
	}

	return false, string(output)
}

// getDefaultRules возвращает правила по умолчанию для языка
func (s *Service) getDefaultRules(language string) []domain.RepairRule {
	switch language {
	case "go":
		return []domain.RepairRule{
			{
				ID:          "go-format",
				Name:        "Go Format",
				Description: "Format Go code with gofmt and goimports",
				Pattern:     `gofmt|goimports`,
				Fix:         "format",
				Priority:    100,
				Language:    "go",
				Category:    "format",
			},
			{
				ID:          "go-imports",
				Name:        "Go Imports",
				Description: "Fix Go imports with go mod tidy",
				Pattern:     `undefined|imported and not used`,
				Fix:         "imports",
				Priority:    90,
				Language:    "go",
				Category:    "import",
			},
		}
	case "typescript", "javascript":
		return []domain.RepairRule{
			{
				ID:          "ts-format",
				Name:        "TypeScript Format",
				Description: "Format TypeScript code with prettier",
				Pattern:     `prettier|format`,
				Fix:         "format",
				Priority:    100,
				Language:    "typescript",
				Category:    "format",
			},
			{
				ID:          "ts-imports",
				Name:        "TypeScript Imports",
				Description: "Fix TypeScript imports with npm install",
				Pattern:     `Cannot find module|Module not found`,
				Fix:         "imports",
				Priority:    90,
				Language:    "typescript",
				Category:    "import",
			},
		}
	default:
		return []domain.RepairRule{}
	}
}
