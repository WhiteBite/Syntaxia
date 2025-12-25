package guardrails

import (
	"path/filepath"
	"regexp"
	"syntaxia/domain"
	"strings"
	"time"
)

// matchesPath проверяет, соответствует ли путь паттерну
func (s *ServiceImpl) matchesPath(path, pattern string) bool {
	if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			s.log.Warning("Invalid glob pattern: " + pattern)
			return false
		}
		return matched
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		s.log.Warning("Invalid regex pattern: " + pattern)
		return false
	}
	return re.MatchString(path)
}

// isCriticalPath проверяет, является ли путь критическим
func (s *ServiceImpl) isCriticalPath(path string) bool {
	criticalPatterns := []string{
		"go\\.mod",
		"package\\.json",
		"package-lock\\.json",
		"yarn\\.lock",
		"pnpm-lock\\.yaml",
	}

	for _, pattern := range criticalPatterns {
		if s.matchesPath(path, pattern) {
			return true
		}
	}

	return false
}

// isEphemeralExpired проверяет, истек ли ephemeral mode
func (s *ServiceImpl) isEphemeralExpired() bool {
	return time.Now().After(s.ephemeralEnd)
}

// disableEphemeralModeInternal отключает ephemeral mode (internal, no lock)
func (s *ServiceImpl) disableEphemeralModeInternal() {
	s.ephemeralMode = false
	s.ephemeralEnd = time.Time{}
	s.log.Info("Ephemeral mode expired and disabled")
}

// initializeDefaultPolicies инициализирует политики по умолчанию
func (s *ServiceImpl) initializeDefaultPolicies() {
	forbiddenPathsPolicy := domain.GuardrailPolicy{
		ID:          "forbidden-paths",
		Name:        "Forbidden Paths",
		Description: "Запрещенные пути и файлы",
		Type:        domain.GuardrailTypeForbiddenPath,
		Severity:    domain.GuardrailSeverityBlock,
		Enabled:     true,
		Rules: []domain.GuardrailRule{
			{
				ID:          "go-mod",
				Pattern:     "go\\.mod",
				Description: "Запрещено изменять go.mod",
				Action:      domain.GuardrailActionBlock,
				Message:     "Изменение go.mod запрещено",
			},
			{
				ID:          "package-json",
				Pattern:     "package\\.json|package-lock\\.json",
				Description: "Запрещено изменять package.json и package-lock.json",
				Action:      domain.GuardrailActionBlock,
				Message:     "Изменение package.json запрещено",
			},
			{
				ID:          "node-modules",
				Pattern:     "node_modules/.*",
				Description: "Запрещено изменять node_modules",
				Action:      domain.GuardrailActionBlock,
				Message:     "Изменение node_modules запрещено",
			},
			{
				ID:          "secrets",
				Pattern:     ".*\\.(key|pem|p12|pfx|env|secret)",
				Description: "Запрещено изменять файлы с секретами",
				Action:      domain.GuardrailActionBlock,
				Message:     "Изменение файлов с секретами запрещено",
			},
			{
				ID:          "binary-files",
				Pattern:     ".*\\.(exe|dll|so|dylib|bin|jar|war|ear)",
				Description: "Запрещено изменять бинарные файлы",
				Action:      domain.GuardrailActionBlock,
				Message:     "Изменение бинарных файлов запрещено",
			},
			{
				ID:          "temp-dirs",
				Pattern:     "(tmp|temp|cache)/.*",
				Description: "Запрещено изменять временные каталоги",
				Action:      domain.GuardrailActionBlock,
				Message:     "Изменение временных каталогов запрещено",
			},
			{
				ID:          "cursor-rules",
				Pattern:     "\\.cursor/rules/cutc\\.mdc",
				Description: "Запрещено изменять внутренний протокол",
				Action:      domain.GuardrailActionBlock,
				Message:     "Изменение .cursor/rules/cutc.mdc запрещено",
			},
		},
	}

	fileBudgetPolicy := domain.BudgetPolicy{
		ID:          "max-files",
		Name:        "Maximum Files",
		Description: "Максимальное количество изменяемых файлов",
		Type:        domain.BudgetTypeFiles,
		Limit:       150,
		Unit:        domain.BudgetUnitCount,
		TimeWindow:  time.Hour,
		Enabled:     true,
	}

	lineBudgetPolicy := domain.BudgetPolicy{
		ID:          "max-lines",
		Name:        "Maximum Lines",
		Description: "Максимальное количество изменяемых строк",
		Type:        domain.BudgetTypeLines,
		Limit:       1500,
		Unit:        domain.BudgetUnitCount,
		TimeWindow:  time.Hour,
		Enabled:     true,
	}

	tokenBudgetPolicy := domain.BudgetPolicy{
		ID:          "max-tokens",
		Name:        "Maximum Tokens",
		Description: "Максимальное количество токенов в запросе",
		Type:        domain.BudgetTypeTokens,
		Limit:       10000,
		Unit:        domain.BudgetUnitCount,
		TimeWindow:  time.Hour,
		Enabled:     true,
	}

	s.policies = append(s.policies, forbiddenPathsPolicy)
	s.budgets = append(s.budgets, fileBudgetPolicy, lineBudgetPolicy, tokenBudgetPolicy)

	s.log.Info("Initialized default guardrail policies")
}
