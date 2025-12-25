package buildpipeline

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syntaxia/domain"
	"syntaxia/infrastructure/sandbox"
	"strings"
	"time"
)

// Language constants
const (
	langGo         = "go"
	langTypeScript = "typescript"
	langJava       = "java"
)

// Impl реализует BuildPipeline
type Impl struct {
	log           domain.Logger
	sandboxRunner domain.SandboxRunner
}

// NewBuildPipeline создает новый build pipeline
func NewBuildPipeline(log domain.Logger) *Impl {
	return &Impl{
		log:           log,
		sandboxRunner: sandbox.NewSandboxRunner(log),
	}
}

// Build выполняет сборку проекта
func (p *Impl) Build(ctx context.Context, projectPath, language string) (*domain.BuildResult, error) {
	p.log.Info(fmt.Sprintf("Building %s project at %s", language, projectPath))
	startTime := time.Now()

	var result *domain.BuildResult
	var err error
	switch language {
	case langGo:
		result, err = p.buildGo(ctx, projectPath)
	case langTypeScript, "ts":
		result, err = p.buildTypeScript(ctx, projectPath)
	case langJava:
		result, err = p.buildJava(ctx, projectPath)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	if err != nil {
		return nil, err
	}
	result.Duration = time.Since(startTime).Seconds()
	return result, nil
}

// TypeCheck выполняет проверку типов
func (p *Impl) TypeCheck(ctx context.Context, projectPath, language string) (*domain.TypeCheckResult, error) {
	p.log.Info(fmt.Sprintf("Type checking %s project at %s", language, projectPath))
	startTime := time.Now()

	var result *domain.TypeCheckResult
	var err error
	switch language {
	case langGo:
		result, err = p.typeCheckGo(ctx, projectPath)
	case langTypeScript, "ts":
		result, err = p.typeCheckTypeScript(ctx, projectPath)
	case langJava:
		result, err = p.typeCheckJava(ctx, projectPath)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	if err != nil {
		return nil, err
	}
	result.Duration = time.Since(startTime).Seconds()
	return result, nil
}

// BuildAndTypeCheck выполняет сборку и проверку типов
func (p *Impl) BuildAndTypeCheck(ctx context.Context, projectPath, language string) (*domain.BuildResult, *domain.TypeCheckResult, error) {
	p.log.Info(fmt.Sprintf("Building and type checking %s project at %s", language, projectPath))

	// Сначала выполняем проверку типов
	typeCheckResult, err := p.TypeCheck(ctx, projectPath, language)
	if err != nil {
		return nil, nil, fmt.Errorf("type check failed: %w", err)
	}

	// Затем выполняем сборку
	buildResult, err := p.Build(ctx, projectPath, language)
	if err != nil {
		return nil, typeCheckResult, fmt.Errorf("build failed: %w", err)
	}

	return buildResult, typeCheckResult, nil
}

// BuildInSandbox выполняет сборку в песочнице
func (p *Impl) BuildInSandbox(ctx context.Context, projectPath, language string, sandboxConfig domain.SandboxConfig) (*domain.BuildResult, error) {
	p.log.Info(fmt.Sprintf("Building %s project in sandbox at %s", language, projectPath))

	// Проверяем доступность песочницы
	if !p.sandboxRunner.IsAvailable(ctx) {
		p.log.Warning("Sandbox not available, falling back to local build")
		return p.Build(ctx, projectPath, language)
	}

	// Определяем команду для сборки
	var command []string
	switch language {
	case "go":
		command = []string{"go", "build", "./..."}
	case "java":
		if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
			command = []string{"mvn", "compile"}
		} else if _, err := os.Stat(filepath.Join(projectPath, "build.gradle")); err == nil {
			command = []string{"gradle", "build"}
		} else {
			return nil, fmt.Errorf("no build configuration found for Java project")
		}
	case "typescript", "ts":
		command = []string{"npm", "run", "build"}
	default:
		return nil, fmt.Errorf("unsupported language for sandbox build: %s", language)
	}

	// Настраиваем монтирование проекта
	if len(sandboxConfig.Mounts) == 0 {
		sandboxConfig.Mounts = []domain.SandboxMount{
			{
				Source:   projectPath,
				Target:   "/workspace",
				ReadOnly: false,
				Type:     "bind",
			},
		}
	}

	// Настраиваем рабочую директорию
	if sandboxConfig.WorkingDir == "" {
		sandboxConfig.WorkingDir = "/workspace"
	}

	// Настраиваем таймаут
	if sandboxConfig.Timeout == 0 {
		sandboxConfig.Timeout = 300 // 5 минут по умолчанию
	}

	// Запускаем команду в песочнице
	result, err := p.sandboxRunner.Run(ctx, sandboxConfig, command)
	if err != nil {
		return nil, fmt.Errorf("sandbox execution failed: %w", err)
	}

	// Конвертируем результат
	buildResult := &domain.BuildResult{
		Language:    language,
		ProjectPath: projectPath,
		Success:     result.Success,
		Output:      result.Output,
		Error:       result.Error,
		Duration:    result.Duration,
		Metadata: map[string]interface{}{
			"sandbox": true,
			"engine":  sandboxConfig.Engine,
			"image":   sandboxConfig.Image,
		},
	}

	if !result.Success {
		buildResult.Warnings = append(buildResult.Warnings, "Build executed in sandbox but failed")
	}

	return buildResult, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (p *Impl) GetSupportedLanguages() []string {
	return []string{"go", "typescript", "ts", "java"}
}

// buildGo выполняет сборку Go проекта
func (p *Impl) buildGo(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    "go",
		ProjectPath: projectPath,
	}

	// Проверяем наличие go.mod
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); os.IsNotExist(err) {
		result.Success = false
		result.Error = "go.mod not found"
		return result, nil
	}

	// Выполняем go build
	cmd := exec.CommandContext(ctx, "go", "build", "-o", "syntaxia.exe", ".")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true

	// Ищем артефакты
	if _, err := os.Stat(filepath.Join(projectPath, "syntaxia.exe")); err == nil {
		result.Artifacts = append(result.Artifacts, "syntaxia.exe")
	}

	return result, nil
}

// runTypeCheck is a helper for running type check commands
func (p *Impl) runTypeCheck(ctx context.Context, projectPath, language string, cmdName string, cmdArgs []string, parseIssues func(string) []*domain.TypeIssue) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    language,
		ProjectPath: projectPath,
	}

	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		if parseIssues != nil {
			result.Issues = parseIssues(string(output))
		}
		return result, nil
	}

	result.Success = true
	return result, nil
}

// typeCheckGo выполняет проверку типов Go проекта
func (p *Impl) typeCheckGo(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	return p.runTypeCheck(ctx, projectPath, langGo, "go", []string{"vet", "./..."}, p.parseGoVetIssues)
}

// buildTypeScript выполняет сборку TypeScript проекта
func (p *Impl) buildTypeScript(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langTypeScript,
		ProjectPath: projectPath,
	}

	// Проверяем наличие package.json
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); os.IsNotExist(err) {
		result.Success = false
		result.Error = "package.json not found"
		return result, nil
	}

	// Выполняем npm run build или tsc
	cmd := exec.CommandContext(ctx, "npm", "run", "build")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		// Пробуем tsc напрямую
		cmd = exec.CommandContext(ctx, "npx", "tsc")
		cmd.Dir = projectPath

		output, err = cmd.CombinedOutput()
		result.Output = string(output)

		if err != nil {
			result.Success = false
			result.Error = err.Error()
			return result, nil
		}
	}

	result.Success = true

	// Ищем артефакты в dist/
	if distPath := filepath.Join(projectPath, "dist"); distPath != "" {
		if _, err := os.Stat(distPath); err == nil {
			result.Artifacts = append(result.Artifacts, "dist/")
		}
	}

	return result, nil
}

// typeCheckTypeScript выполняет проверку типов TypeScript проекта
func (p *Impl) typeCheckTypeScript(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	return p.runTypeCheck(ctx, projectPath, langTypeScript, "npx", []string{"tsc", "--noEmit"}, p.parseTypeScriptIssues)
}

// buildJava выполняет сборку Java проекта
func (p *Impl) buildJava(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langJava,
		ProjectPath: projectPath,
	}

	// Проверяем наличие Java среды
	if !p.checkJavaEnvironment() {
		err := fmt.Errorf("Java environment not available (java, javac, mvn, or gradle not found)")
		result.Success = false
		result.Error = err.Error()
		result.Warnings = append(result.Warnings, "Java build tools not available, skipping build")
		return result, err
	}

	// Проверяем наличие pom.xml или build.gradle
	if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		// Maven проект
		return p.buildMavenProject(ctx, projectPath)
	}

	if _, err := os.Stat(filepath.Join(projectPath, "build.gradle")); err == nil {
		// Gradle проект
		return p.buildGradleProject(ctx, projectPath)
	}

	err := fmt.Errorf("neither pom.xml nor build.gradle found")
	result.Success = false
	result.Error = err.Error()
	result.Warnings = append(result.Warnings, "No Java build configuration found")
	return result, err
}

// typeCheckJava выполняет проверку типов Java проекта
func (p *Impl) typeCheckJava(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langJava,
		ProjectPath: projectPath,
	}

	// Проверяем наличие Java среды
	if !p.checkJavaEnvironment() {
		result.Success = false
		result.Error = "Java environment not available"
		return result, fmt.Errorf("%s", result.Error)
	}

	// Проверяем наличие pom.xml или build.gradle
	if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		// Maven проект - компиляция вкл��чает проверку типов
		cmd := exec.CommandContext(ctx, "mvn", "compile", "-q")
		cmd.Dir = projectPath

		output, err := cmd.CombinedOutput()
		result.Output = string(output)

		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.Issues = p.parseMavenIssues(string(output))
			return result, err
		}

		result.Success = true

		// Дополнительно запускаем ErrorProne, если доступен
		if p.hasErrorPronePlugin(projectPath) {
			errorProneResult, _ := p.runErrorProne(ctx, projectPath, "mvn")
			if !errorProneResult.Success {
				result.Issues = append(result.Issues, p.parseErrorProneIssues(errorProneResult.Output)...)
			}
		}
	} else if _, err := os.Stat(filepath.Join(projectPath, "build.gradle")); err == nil {
		// Gradle проект - компиляция включает проверку типов
		cmd := exec.CommandContext(ctx, "gradle", "compileJava", "--quiet")
		cmd.Dir = projectPath

		output, err := cmd.CombinedOutput()
		result.Output = string(output)

		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.Issues = p.parseGradleIssues(string(output))
			return result, err
		}

		result.Success = true
	} else {
		result.Success = false
		result.Error = "neither pom.xml nor build.gradle found"
		return result, fmt.Errorf("%s", result.Error)
	}

	return result, nil
}

// parseGoVetIssues парсит ошибки go vet
func (p *Impl) parseGoVetIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 3)
			if len(parts) >= 3 {
				issue := &domain.TypeIssue{
					File:     parts[0],
					Severity: "warning",
					Message:  strings.TrimSpace(parts[2]),
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// parseTypeScriptIssues парсит ошибки TypeScript
func (p *Impl) parseTypeScriptIssues(output string) []*domain.TypeIssue {
	re := regexp.MustCompile(`([^(]+)\((\d+),(\d+)\):\s+error\s+(TS\d+):\s+(.+)`)
	return p.parseIssuesWithRegexCode(output, func(line string) bool {
		return strings.Contains(line, ".ts") && strings.Contains(line, "error TS")
	}, re)
}

// parseIssuesWithRegex is a helper for parsing build output with regex (4 groups: file, line, col, message)
func (p *Impl) parseIssuesWithRegex(output string, lineFilter func(string) bool, re *regexp.Regexp) []*domain.TypeIssue {
	var issues []*domain.TypeIssue
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if lineFilter(line) {
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 5 {
				issues = append(issues, &domain.TypeIssue{
					File:     matches[1],
					Line:     p.parseInt(matches[2]),
					Column:   p.parseInt(matches[3]),
					Message:  matches[4],
					Severity: "error",
				})
			}
		}
	}
	return issues
}

// parseIssuesWithRegexCode is a helper for parsing build output with regex (5 groups: file, line, col, code, message)
func (p *Impl) parseIssuesWithRegexCode(output string, lineFilter func(string) bool, re *regexp.Regexp) []*domain.TypeIssue {
	var issues []*domain.TypeIssue
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if lineFilter(line) {
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 6 {
				issues = append(issues, &domain.TypeIssue{
					File:     matches[1],
					Line:     p.parseInt(matches[2]),
					Column:   p.parseInt(matches[3]),
					Code:     matches[4],
					Message:  matches[5],
					Severity: "error",
				})
			}
		}
	}
	return issues
}

// parseMavenIssues парсит ошибки Maven
func (p *Impl) parseMavenIssues(output string) []*domain.TypeIssue {
	re := regexp.MustCompile(`\[ERROR\]\s+([^:]+):\[(\d+),(\d+)\]\s+(.+)`)
	return p.parseIssuesWithRegex(output, func(line string) bool {
		return strings.Contains(line, "[ERROR]") && strings.Contains(line, ".java")
	}, re)
}

// parseGradleIssues парсит ошибки Gradle
func (p *Impl) parseGradleIssues(output string) []*domain.TypeIssue {
	re := regexp.MustCompile(`([^:]+):(\d+):(\d+):\s+error:\s+(.+)`)
	return p.parseIssuesWithRegex(output, func(line string) bool {
		return strings.Contains(line, "error:") && strings.Contains(line, ".java")
	}, re)
}

// parseInt парсит строку в int
func (p *Impl) parseInt(s string) int {
	var i int
	_, _ = fmt.Sscanf(s, "%d", &i)
	return i
}

// checkJavaEnvironment проверяет наличие Java среды
func (p *Impl) checkJavaEnvironment() bool {
	// Проверяем наличие java
	if _, err := exec.LookPath("java"); err != nil {
		p.log.Warning("Java runtime not found")
		return false
	}

	// Проверяем наличие javac
	if _, err := exec.LookPath("javac"); err != nil {
		p.log.Warning("Java compiler not found")
		return false
	}

	// Проверяем наличие mvn или gradle
	if _, err := exec.LookPath("mvn"); err != nil {
		if _, err := exec.LookPath("gradle"); err != nil {
			p.log.Warning("Neither Maven nor Gradle found")
			return false
		}
	}

	return true
}

// buildJavaProject is a helper for building Java projects with Maven or Gradle
func (p *Impl) buildJavaProject(ctx context.Context, projectPath, toolName string, cmdArgs []string, artifactDir string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langJava,
		ProjectPath: projectPath,
	}

	if _, err := exec.LookPath(toolName); err != nil {
		result.Success = false
		result.Error = toolName + " not found in PATH"
		result.Warnings = append(result.Warnings, toolName+" build tool not available")
		return result, err
	}

	cmd := exec.CommandContext(ctx, toolName, cmdArgs...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Warnings = append(result.Warnings, toolName+" build failed")
		return result, err
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, artifactDir)

	if p.hasJUnitTests(projectPath) {
		testResult, testErr := p.runJUnitTests(ctx, projectPath, toolName)
		if testErr != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("JUnit test execution failed with error: %v", testErr))
		} else if !testResult.Success {
			result.Warnings = append(result.Warnings, "JUnit tests failed: "+testResult.Error)
		}
	}

	return result, nil
}

// buildMavenProject выполняет сборку Maven проекта
func (p *Impl) buildMavenProject(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	return p.buildJavaProject(ctx, projectPath, "mvn", []string{"compile", "-q"}, "target/")
}

// buildGradleProject выполняет сборку Gradle проекта
func (p *Impl) buildGradleProject(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	return p.buildJavaProject(ctx, projectPath, "gradle", []string{"build", "--quiet"}, "build/")
}

// hasJUnitTests проверяет наличие JUnit тестов
func (p *Impl) hasJUnitTests(projectPath string) bool {
	// Проверяем наличие тестовых директорий
	testDirs := []string{
		filepath.Join(projectPath, "src", "test", "java"),
		filepath.Join(projectPath, "test"),
		filepath.Join(projectPath, "tests"),
	}

	for _, dir := range testDirs {
		if _, err := os.Stat(dir); err == nil {
			return true
		}
	}

	return false
}

// runJUnitTests запускает JUnit тесты
func (p *Impl) runJUnitTests(ctx context.Context, projectPath, buildTool string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    "java",
		ProjectPath: projectPath,
	}

	var cmd *exec.Cmd
	if buildTool == "mvn" {
		cmd = exec.CommandContext(ctx, "mvn", "test", "-q")
	} else {
		cmd = exec.CommandContext(ctx, "gradle", "test", "--quiet")
	}
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Warnings = append(result.Warnings, "JUnit test execution failed")
		return result, err
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, "test-results/")
	return result, nil
}

// runErrorProne запускает ErrorProne анализ
func (p *Impl) runErrorProne(ctx context.Context, projectPath, buildTool string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    "java",
		ProjectPath: projectPath,
	}

	// ErrorProne доступен только для Maven
	if buildTool != "mvn" {
		err := fmt.Errorf("ErrorProne only supported with Maven")
		result.Success = false
		result.Error = err.Error()
		result.Warnings = append(result.Warnings, "ErrorProne analysis skipped (Gradle not supported)")
		return result, err
	}

	// Проверяем наличие ErrorProne plugin в pom.xml
	if !p.hasErrorPronePlugin(projectPath) {
		err := fmt.Errorf("ErrorProne plugin not configured")
		result.Success = false
		result.Error = err.Error()
		result.Warnings = append(result.Warnings, "ErrorProne analysis skipped (plugin not configured)")
		return result, err
	}

	// Запускаем ErrorProne анализ
	cmd := exec.CommandContext(ctx, "mvn", "compile", "-Perror-prone")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Warnings = append(result.Warnings, "ErrorProne analysis failed")
		return result, err
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, "error-prone-reports/")
	return result, nil
}

// hasErrorPronePlugin проверяет наличие ErrorProne plugin в pom.xml
func (p *Impl) hasErrorPronePlugin(projectPath string) bool {
	pomPath := filepath.Join(projectPath, "pom.xml")
	content, err := os.ReadFile(pomPath)
	if err != nil {
		p.log.Warning(fmt.Sprintf("could not read pom.xml to check for ErrorProne plugin: %v", err))
		return false
	}

	return strings.Contains(string(content), "error-prone") ||
		strings.Contains(string(content), "errorprone")
}

// parseErrorProneIssues парсит ошибки ErrorProne
func (p *Impl) parseErrorProneIssues(output string) []*domain.TypeIssue {
	re := regexp.MustCompile(`\[ERROR\]\s+([^:]+):\[(\d+),(\d+)\]\s+error:\s+\[([^\]]+)\]\s+(.+)`)
	return p.parseIssuesWithRegexCode(output, func(line string) bool {
		return strings.Contains(line, "[ERROR]") && strings.Contains(line, "error-prone")
	}, re)
}
