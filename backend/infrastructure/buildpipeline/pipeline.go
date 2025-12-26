package buildpipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"syntaxia/domain"
	"syntaxia/infrastructure/sandbox"
)

// Language constants
const (
	langGo         = "go"
	langTypeScript = "typescript"
	langJava       = "java"
	langRust       = "rust"
	langPython     = "python"
	langKotlin     = "kotlin"
	langCSharp     = "csharp"
	langDart       = "dart"
	langCpp        = "cpp"
	langSwift      = "swift"
	langPHP        = "php"
	langRuby       = "ruby"
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
	case langRust, "rs":
		result, err = p.buildRust(ctx, projectPath)
	case langPython, "py":
		result, err = p.buildPython(ctx, projectPath)
	case langCSharp, "cs":
		result, err = p.buildCSharp(ctx, projectPath)
	case langKotlin, "kt":
		result, err = p.buildKotlin(ctx, projectPath)
	case langDart:
		result, err = p.buildDart(ctx, projectPath)
	case langCpp, "c", "c++", "cc":
		result, err = p.buildCpp(ctx, projectPath)
	case langSwift:
		result, err = p.buildSwift(ctx, projectPath)
	case langPHP:
		result, err = p.buildPHP(ctx, projectPath)
	case langRuby, "rb":
		result, err = p.buildRuby(ctx, projectPath)
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
	case langRust, "rs":
		result, err = p.typeCheckRust(ctx, projectPath)
	case langPython, "py":
		result, err = p.typeCheckPython(ctx, projectPath)
	case langCSharp, "cs":
		result, err = p.typeCheckCSharp(ctx, projectPath)
	case langKotlin, "kt":
		result, err = p.typeCheckKotlin(ctx, projectPath)
	case langDart:
		result, err = p.typeCheckDart(ctx, projectPath)
	case langCpp, "c", "c++", "cc":
		result, err = p.typeCheckCpp(ctx, projectPath)
	case langSwift:
		result, err = p.typeCheckSwift(ctx, projectPath)
	case langPHP:
		result, err = p.typeCheckPHP(ctx, projectPath)
	case langRuby, "rb":
		result, err = p.typeCheckRuby(ctx, projectPath)
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
	return []string{"go", "typescript", "ts", "java", "rust", "rs", "python", "py", "csharp", "cs", "kotlin", "kt", "dart", "cpp", "c", "c++", "cc", "swift", "php", "ruby", "rb"}
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

// buildRust выполняет сборку Rust проекта
func (p *Impl) buildRust(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langRust,
		ProjectPath: projectPath,
	}

	// Проверяем наличие Cargo.toml
	cargoPath := filepath.Join(projectPath, "Cargo.toml")
	if _, err := os.Stat(cargoPath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "Cargo.toml not found"
		return result, nil
	}

	// Проверяем наличие cargo
	if _, err := exec.LookPath("cargo"); err != nil {
		result.Success = false
		result.Error = "cargo not found in PATH"
		result.Warnings = append(result.Warnings, "Rust toolchain not available")
		return result, nil
	}

	// Выполняем cargo build
	cmd := exec.CommandContext(ctx, "cargo", "build", "--message-format=short")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true

	// Ищем артефакты в target/debug или target/release
	debugPath := filepath.Join(projectPath, "target", "debug")
	if _, err := os.Stat(debugPath); err == nil {
		result.Artifacts = append(result.Artifacts, "target/debug/")
	}

	return result, nil
}

// typeCheckRust выполняет проверку типов Rust проекта через cargo check
func (p *Impl) typeCheckRust(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langRust,
		ProjectPath: projectPath,
	}

	// Проверяем наличие Cargo.toml
	cargoPath := filepath.Join(projectPath, "Cargo.toml")
	if _, err := os.Stat(cargoPath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "Cargo.toml not found"
		return result, nil
	}

	// Проверяем наличие cargo
	if _, err := exec.LookPath("cargo"); err != nil {
		result.Success = false
		result.Error = "cargo not found in PATH"
		return result, nil
	}

	// Выполняем cargo check
	cmd := exec.CommandContext(ctx, "cargo", "check", "--message-format=short")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parseRustIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// parseRustIssues парсит ошибки Rust компилятора
func (p *Impl) parseRustIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// Rust error format: "error[E0425]: cannot find value `x` in this scope"
	// followed by: " --> src/main.rs:5:13"
	lines := strings.Split(output, "\n")

	var currentMessage string
	var currentCode string
	var currentSeverity string

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Parse error/warning line
		if strings.HasPrefix(line, "error[") {
			currentSeverity = "error"
			// Extract error code like E0425
			if idx := strings.Index(line, "]"); idx > 6 {
				currentCode = line[6:idx]
			}
			if colonIdx := strings.Index(line, "]:"); colonIdx != -1 {
				currentMessage = strings.TrimSpace(line[colonIdx+2:])
			}
		} else if strings.HasPrefix(line, "warning[") {
			currentSeverity = "warning"
			if idx := strings.Index(line, "]"); idx > 8 {
				currentCode = line[8:idx]
			}
			if colonIdx := strings.Index(line, "]:"); colonIdx != -1 {
				currentMessage = strings.TrimSpace(line[colonIdx+2:])
			}
		} else if strings.HasPrefix(line, "error:") {
			currentSeverity = "error"
			currentCode = ""
			currentMessage = strings.TrimSpace(line[6:])
		} else if strings.HasPrefix(line, "warning:") {
			currentSeverity = "warning"
			currentCode = ""
			currentMessage = strings.TrimSpace(line[8:])
		}

		// Parse location line: " --> src/main.rs:5:13"
		if strings.HasPrefix(line, "-->") && currentMessage != "" {
			locPart := strings.TrimPrefix(line, "-->")
			locPart = strings.TrimSpace(locPart)

			// Parse "src/main.rs:5:13"
			parts := strings.Split(locPart, ":")
			if len(parts) >= 3 {
				issue := &domain.TypeIssue{
					File:     parts[0],
					Line:     p.parseInt(parts[1]),
					Column:   p.parseInt(parts[2]),
					Severity: currentSeverity,
					Message:  currentMessage,
					Code:     currentCode,
				}
				issues = append(issues, issue)
			}

			// Reset for next error
			currentMessage = ""
			currentCode = ""
			currentSeverity = ""
		}

		// Also handle short format: "src/main.rs:5:13: error[E0425]: message"
		if strings.Contains(line, ".rs:") && (strings.Contains(line, ": error") || strings.Contains(line, ": warning")) {
			re := regexp.MustCompile(`([^:]+\.rs):(\d+):(\d+):\s+(error|warning)(?:\[([^\]]+)\])?:\s+(.+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 7 {
				issue := &domain.TypeIssue{
					File:     matches[1],
					Line:     p.parseInt(matches[2]),
					Column:   p.parseInt(matches[3]),
					Severity: matches[4],
					Code:     matches[5],
					Message:  matches[6],
				}
				issues = append(issues, issue)
			}
		}

		_ = i // Suppress unused variable warning
	}

	return issues
}

// buildPython выполняет сборку Python проекта
func (p *Impl) buildPython(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langPython,
		ProjectPath: projectPath,
	}

	// Проверяем наличие python
	pythonCmd := p.findPythonCommand()
	if pythonCmd == "" {
		result.Success = false
		result.Error = "python not found in PATH"
		result.Warnings = append(result.Warnings, "Python interpreter not available")
		return result, nil
	}

	// Устанавливаем зависимости в зависимости от конфигурации проекта
	var installOutput string
	var installErr error

	// Проверяем pyproject.toml → poetry install
	if _, err := os.Stat(filepath.Join(projectPath, "pyproject.toml")); err == nil {
		if _, poetryErr := exec.LookPath("poetry"); poetryErr == nil {
			cmd := exec.CommandContext(ctx, "poetry", "install", "--no-interaction")
			cmd.Dir = projectPath
			output, err := cmd.CombinedOutput()
			installOutput = string(output)
			installErr = err
			if err == nil {
				result.Artifacts = append(result.Artifacts, "poetry.lock")
			}
		} else {
			// Fallback to pip install if poetry not available
			cmd := exec.CommandContext(ctx, pythonCmd, "-m", "pip", "install", "-e", ".")
			cmd.Dir = projectPath
			output, err := cmd.CombinedOutput()
			installOutput = string(output)
			installErr = err
		}
	} else if _, err := os.Stat(filepath.Join(projectPath, "requirements.txt")); err == nil {
		// requirements.txt → pip install -r requirements.txt
		cmd := exec.CommandContext(ctx, pythonCmd, "-m", "pip", "install", "-r", "requirements.txt")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		installOutput = string(output)
		installErr = err
	} else if _, err := os.Stat(filepath.Join(projectPath, "setup.py")); err == nil {
		// setup.py → pip install -e .
		cmd := exec.CommandContext(ctx, pythonCmd, "-m", "pip", "install", "-e", ".")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		installOutput = string(output)
		installErr = err
	}

	if installErr != nil {
		result.Output = installOutput
		result.Success = false
		result.Error = installErr.Error()
		result.Warnings = append(result.Warnings, "Dependency installation failed")
		return result, nil
	}

	// Выполняем синтаксическую проверку .py файлов
	syntaxOutput, syntaxErr := p.checkPythonSyntax(ctx, projectPath, pythonCmd)
	result.Output = installOutput + "\n" + syntaxOutput

	if syntaxErr != nil {
		result.Success = false
		result.Error = syntaxErr.Error()
		return result, nil
	}

	result.Success = true
	return result, nil
}

// typeCheckPython выполняет проверку типов Python проекта через mypy
func (p *Impl) typeCheckPython(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langPython,
		ProjectPath: projectPath,
	}

	// Проверяем наличие mypy
	if _, err := exec.LookPath("mypy"); err != nil {
		result.Success = false
		result.Error = "mypy not found in PATH"
		result.Issues = append(result.Issues, &domain.TypeIssue{
			Severity: "warning",
			Message:  "mypy not installed, type checking skipped. Install with: pip install mypy",
		})
		return result, nil
	}

	// Выполняем mypy
	cmd := exec.CommandContext(ctx, "mypy", ".", "--ignore-missing-imports", "--no-error-summary")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parsePythonIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// findPythonCommand находит доступную команду Python
func (p *Impl) findPythonCommand() string {
	// Пробуем python3 сначала
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	// Затем python
	if _, err := exec.LookPath("python"); err == nil {
		return "python"
	}
	return ""
}

// checkPythonSyntax проверяет синтаксис Python файлов
func (p *Impl) checkPythonSyntax(ctx context.Context, projectPath, pythonCmd string) (string, error) {
	var allOutput strings.Builder
	var syntaxErrors []string

	// Находим все .py файлы
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории и не-.py файлы
		if info.IsDir() {
			// Пропускаем виртуальные окружения и кэш
			if info.Name() == "venv" || info.Name() == ".venv" ||
				info.Name() == "__pycache__" || info.Name() == ".git" ||
				info.Name() == "node_modules" || info.Name() == ".mypy_cache" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".py") {
			return nil
		}

		// Проверяем синтаксис файла
		cmd := exec.CommandContext(ctx, pythonCmd, "-m", "py_compile", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			allOutput.WriteString(string(output))
			allOutput.WriteString("\n")
			syntaxErrors = append(syntaxErrors, path)
		}

		return nil
	})

	if err != nil {
		return allOutput.String(), fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(syntaxErrors) > 0 {
		return allOutput.String(), fmt.Errorf("syntax errors in %d file(s)", len(syntaxErrors))
	}

	return allOutput.String(), nil
}

// parsePythonIssues парсит ошибки mypy
func (p *Impl) parsePythonIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// mypy output format: "file.py:10: error: Message [error-code]"
	// or: "file.py:10:5: error: Message [error-code]"
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip summary lines
		if strings.HasPrefix(line, "Found ") || strings.HasPrefix(line, "Success:") {
			continue
		}

		// Parse mypy format with column: "file.py:10:5: error: Message [code]"
		reWithCol := regexp.MustCompile(`^([^:]+):(\d+):(\d+):\s+(error|warning|note):\s+(.+?)(?:\s+\[([^\]]+)\])?$`)
		matches := reWithCol.FindStringSubmatch(line)
		if len(matches) >= 6 {
			severity := matches[4]
			if severity == "note" {
				severity = "info"
			}
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[2]),
				Column:   p.parseInt(matches[3]),
				Severity: severity,
				Message:  matches[5],
			}
			if len(matches) >= 7 && matches[6] != "" {
				issue.Code = matches[6]
			}
			issues = append(issues, issue)
			continue
		}

		// Parse mypy format without column: "file.py:10: error: Message [code]"
		reNoCol := regexp.MustCompile(`^([^:]+):(\d+):\s+(error|warning|note):\s+(.+?)(?:\s+\[([^\]]+)\])?$`)
		matches = reNoCol.FindStringSubmatch(line)
		if len(matches) >= 5 {
			severity := matches[3]
			if severity == "note" {
				severity = "info"
			}
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[2]),
				Severity: severity,
				Message:  matches[4],
			}
			if len(matches) >= 6 && matches[5] != "" {
				issue.Code = matches[5]
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// buildCSharp выполняет сборку C# проекта
func (p *Impl) buildCSharp(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langCSharp,
		ProjectPath: projectPath,
	}

	// Проверяем наличие dotnet
	if _, err := exec.LookPath("dotnet"); err != nil {
		result.Success = false
		result.Error = "dotnet not found in PATH"
		result.Warnings = append(result.Warnings, ".NET SDK not available")
		return result, nil
	}

	// Проверяем наличие .csproj или .sln файла
	hasSolution := p.hasDotnetProject(projectPath)
	if !hasSolution {
		result.Success = false
		result.Error = "no .csproj or .sln file found"
		return result, nil
	}

	// Выполняем dotnet build
	cmd := exec.CommandContext(ctx, "dotnet", "build", "--nologo")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true

	// Ищем артефакты в bin/
	binPath := filepath.Join(projectPath, "bin")
	if _, err := os.Stat(binPath); err == nil {
		result.Artifacts = append(result.Artifacts, "bin/")
	}

	return result, nil
}

// typeCheckCSharp выполняет проверку типов C# проекта через dotnet build
func (p *Impl) typeCheckCSharp(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langCSharp,
		ProjectPath: projectPath,
	}

	// Проверяем наличие dotnet
	if _, err := exec.LookPath("dotnet"); err != nil {
		result.Success = false
		result.Error = "dotnet not found in PATH"
		return result, nil
	}

	// Проверяем наличие .csproj или .sln файла
	if !p.hasDotnetProject(projectPath) {
		result.Success = false
		result.Error = "no .csproj or .sln file found"
		return result, nil
	}

	// dotnet build включает проверку типов
	cmd := exec.CommandContext(ctx, "dotnet", "build", "--nologo", "--no-restore")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parseDotnetIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// hasDotnetProject проверяет наличие .csproj или .sln файла
func (p *Impl) hasDotnetProject(projectPath string) bool {
	// Проверяем .sln файлы
	slnFiles, _ := filepath.Glob(filepath.Join(projectPath, "*.sln"))
	if len(slnFiles) > 0 {
		return true
	}

	// Проверяем .csproj файлы
	csprojFiles, _ := filepath.Glob(filepath.Join(projectPath, "*.csproj"))
	if len(csprojFiles) > 0 {
		return true
	}

	// Проверяем .csproj в поддиректориях
	found := false
	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".csproj") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})

	return found
}

// parseDotnetIssues парсит ошибки dotnet build
func (p *Impl) parseDotnetIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// dotnet build output format: "File.cs(10,5): error CS0103: The name 'x' does not exist"
	// or: "path/File.cs(10,5): warning CS0168: The variable 'x' is declared but never used"
	lines := strings.Split(output, "\n")

	// Regex for dotnet error format
	re := regexp.MustCompile(`([^(]+)\((\d+),(\d+)\):\s+(error|warning)\s+(CS\d+):\s+(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 7 {
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[2]),
				Column:   p.parseInt(matches[3]),
				Severity: matches[4],
				Code:     matches[5],
				Message:  matches[6],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// buildKotlin выполняет сборку Kotlin проекта
func (p *Impl) buildKotlin(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langKotlin,
		ProjectPath: projectPath,
	}

	// Проверяем наличие build.gradle или build.gradle.kts
	hasGradle := p.hasGradleKotlinProject(projectPath)
	if !hasGradle {
		result.Success = false
		result.Error = "no build.gradle or build.gradle.kts found"
		return result, nil
	}

	// Определяем команду gradle
	gradleCmd := p.getGradleCommand(projectPath)
	if gradleCmd == "" {
		result.Success = false
		result.Error = "gradle not found in PATH"
		result.Warnings = append(result.Warnings, "Gradle build tool not available")
		return result, nil
	}

	// Выполняем gradle build
	cmd := exec.CommandContext(ctx, gradleCmd, "build", "--quiet")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true

	// Ищем артефакты в build/
	buildPath := filepath.Join(projectPath, "build")
	if _, err := os.Stat(buildPath); err == nil {
		result.Artifacts = append(result.Artifacts, "build/")
	}

	return result, nil
}

// typeCheckKotlin выполняет проверку типов Kotlin проекта через gradle compileKotlin
func (p *Impl) typeCheckKotlin(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langKotlin,
		ProjectPath: projectPath,
	}

	// Проверяем наличие build.gradle или build.gradle.kts
	if !p.hasGradleKotlinProject(projectPath) {
		result.Success = false
		result.Error = "no build.gradle or build.gradle.kts found"
		return result, nil
	}

	// Определяем команду gradle
	gradleCmd := p.getGradleCommand(projectPath)
	if gradleCmd == "" {
		result.Success = false
		result.Error = "gradle not found in PATH"
		return result, nil
	}

	// Выполняем gradle compileKotlin
	cmd := exec.CommandContext(ctx, gradleCmd, "compileKotlin", "--quiet")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parseKotlinIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// hasGradleKotlinProject проверяет наличие Gradle Kotlin проекта
func (p *Impl) hasGradleKotlinProject(projectPath string) bool {
	gradleFiles := []string{"build.gradle.kts", "build.gradle"}
	for _, gradleFile := range gradleFiles {
		if _, err := os.Stat(filepath.Join(projectPath, gradleFile)); err == nil {
			return true
		}
	}
	return false
}

// getGradleCommand возвращает команду для запуска Gradle
func (p *Impl) getGradleCommand(projectPath string) string {
	// Проверяем наличие gradlew wrapper
	if p.hasGradleWrapper(projectPath) {
		if os.PathSeparator == '\\' {
			return filepath.Join(projectPath, "gradlew.bat")
		}
		return filepath.Join(projectPath, "gradlew")
	}

	// Используем системный gradle
	if _, err := exec.LookPath("gradle"); err == nil {
		return "gradle"
	}

	return ""
}

// hasGradleWrapper проверяет наличие Gradle wrapper
func (p *Impl) hasGradleWrapper(projectPath string) bool {
	wrapperFiles := []string{"gradlew", "gradlew.bat"}
	for _, wrapper := range wrapperFiles {
		if _, err := os.Stat(filepath.Join(projectPath, wrapper)); err == nil {
			return true
		}
	}
	return false
}

// parseKotlinIssues парсит ошибки Kotlin компилятора
func (p *Impl) parseKotlinIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// Kotlin error format: "e: file.kt:10:5: Unresolved reference: x"
	// or: "w: file.kt:10:5: Parameter 'x' is never used"
	lines := strings.Split(output, "\n")

	// Regex for Kotlin error format
	re := regexp.MustCompile(`^([ew]):\s*([^:]+):(\d+):(\d+):\s+(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 6 {
			severity := "error"
			if matches[1] == "w" {
				severity = "warning"
			}
			issue := &domain.TypeIssue{
				File:     matches[2],
				Line:     p.parseInt(matches[3]),
				Column:   p.parseInt(matches[4]),
				Severity: severity,
				Message:  matches[5],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// buildDart выполняет сборку Dart/Flutter проекта
func (p *Impl) buildDart(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langDart,
		ProjectPath: projectPath,
	}

	// Проверяем наличие pubspec.yaml
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "pubspec.yaml not found"
		return result, nil
	}

	// Определяем, Flutter или чистый Dart проект
	isFlutter := p.isFlutterProject(projectPath)

	// Запускаем pub get для зависимостей
	var pubCmd *exec.Cmd
	if isFlutter {
		pubCmd = exec.CommandContext(ctx, "flutter", "pub", "get")
	} else {
		pubCmd = exec.CommandContext(ctx, "dart", "pub", "get")
	}
	pubCmd.Dir = projectPath

	pubOutput, err := pubCmd.CombinedOutput()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("pub get failed: %v", err)
		result.Output = string(pubOutput)
		return result, nil
	}

	// Выполняем сборку
	var buildCmd *exec.Cmd
	if isFlutter {
		buildCmd = exec.CommandContext(ctx, "flutter", "build", "apk", "--debug")
	} else {
		// Для чистого Dart пробуем скомпилировать
		buildCmd = exec.CommandContext(ctx, "dart", "compile", "exe", "bin/main.dart", "-o", "build/main")
	}
	buildCmd.Dir = projectPath

	buildOutput, err := buildCmd.CombinedOutput()
	result.Output = string(pubOutput) + "\n" + string(buildOutput)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	if isFlutter {
		result.Artifacts = append(result.Artifacts, "build/app/outputs/")
	} else {
		result.Artifacts = append(result.Artifacts, "build/")
	}

	return result, nil
}

// typeCheckDart выполняет проверку типов Dart проекта через dart analyze
func (p *Impl) typeCheckDart(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langDart,
		ProjectPath: projectPath,
	}

	// Проверяем наличие pubspec.yaml
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "pubspec.yaml not found"
		return result, nil
	}

	// Проверяем наличие dart
	if _, err := exec.LookPath("dart"); err != nil {
		result.Success = false
		result.Error = "dart not found in PATH"
		return result, nil
	}

	// Выполняем dart analyze
	cmd := exec.CommandContext(ctx, "dart", "analyze", "--format=machine")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parseDartIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// isFlutterProject проверяет, является ли проект Flutter проектом
func (p *Impl) isFlutterProject(projectPath string) bool {
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "flutter:")
}

// parseDartIssues парсит ошибки dart analyze
func (p *Impl) parseDartIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// Machine format: SEVERITY|TYPE|ERROR_CODE|FILE|LINE|COLUMN|LENGTH|MESSAGE
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`^(ERROR|WARNING|INFO)\|([^|]+)\|([^|]+)\|([^|]+)\|(\d+)\|(\d+)\|(\d+)\|(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 9 {
			issue := &domain.TypeIssue{
				File:     matches[4],
				Line:     p.parseInt(matches[5]),
				Column:   p.parseInt(matches[6]),
				Severity: strings.ToLower(matches[1]),
				Code:     matches[3],
				Message:  matches[8],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// buildCpp выполняет сборку C/C++ проекта
func (p *Impl) buildCpp(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	// Detect build system
	buildSystem := p.detectCppBuildSystem(projectPath)

	switch buildSystem {
	case "cmake":
		return p.buildCMakeProject(ctx, projectPath)
	case "make":
		return p.buildMakeProject(ctx, projectPath)
	case "meson":
		return p.buildMesonProject(ctx, projectPath)
	default:
		return p.buildCppDirect(ctx, projectPath)
	}
}

// typeCheckCpp выполняет проверку типов C/C++ проекта
func (p *Impl) typeCheckCpp(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langCpp,
		ProjectPath: projectPath,
	}

	// For C/C++, type checking is done during compilation
	// We use clang or gcc with -fsyntax-only flag
	compiler := p.findCppCompiler()
	if compiler == "" {
		result.Success = false
		result.Error = "no C++ compiler found (g++, clang++)"
		return result, nil
	}

	// Find source files
	sourceFiles, err := p.findCppSourceFiles(projectPath)
	if err != nil || len(sourceFiles) == 0 {
		result.Success = false
		result.Error = "no C/C++ source files found"
		return result, nil
	}

	// Run syntax check on each file
	var allOutput strings.Builder
	var issues []*domain.TypeIssue

	for _, srcFile := range sourceFiles {
		args := []string{"-fsyntax-only", "-std=c++17", srcFile}
		cmd := exec.CommandContext(ctx, compiler, args...)
		cmd.Dir = projectPath

		output, err := cmd.CombinedOutput()
		allOutput.WriteString(string(output))

		if err != nil {
			fileIssues := p.parseCppIssues(string(output))
			issues = append(issues, fileIssues...)
		}
	}

	result.Output = allOutput.String()
	result.Issues = issues
	result.Success = len(issues) == 0

	if !result.Success {
		result.Error = fmt.Sprintf("found %d type/syntax issues", len(issues))
	}

	return result, nil
}

// detectCppBuildSystem detects the build system used
func (p *Impl) detectCppBuildSystem(projectPath string) string {
	if _, err := os.Stat(filepath.Join(projectPath, "CMakeLists.txt")); err == nil {
		return "cmake"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "Makefile")); err == nil {
		return "make"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "meson.build")); err == nil {
		return "meson"
	}
	return "direct"
}

// buildCMakeProject builds a CMake project
func (p *Impl) buildCMakeProject(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langCpp,
		ProjectPath: projectPath,
	}

	// Check for cmake
	if _, err := exec.LookPath("cmake"); err != nil {
		result.Success = false
		result.Error = "cmake not found in PATH"
		result.Warnings = append(result.Warnings, "CMake build tool not available")
		return result, nil
	}

	buildDir := filepath.Join(projectPath, "build")

	// Create build directory
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create build directory: %v", err)
		return result, nil
	}

	// Configure with cmake
	configCmd := exec.CommandContext(ctx, "cmake", "-S", projectPath, "-B", buildDir)
	configCmd.Dir = projectPath

	configOutput, err := configCmd.CombinedOutput()
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Output = string(configOutput)
		return result, nil
	}

	// Build
	buildCmd := exec.CommandContext(ctx, "cmake", "--build", buildDir)
	buildCmd.Dir = projectPath

	buildOutput, err := buildCmd.CombinedOutput()
	result.Output = string(configOutput) + "\n" + string(buildOutput)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, "build/")

	return result, nil
}

// buildMakeProject builds a Make project
func (p *Impl) buildMakeProject(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langCpp,
		ProjectPath: projectPath,
	}

	// Check for make
	if _, err := exec.LookPath("make"); err != nil {
		result.Success = false
		result.Error = "make not found in PATH"
		result.Warnings = append(result.Warnings, "Make build tool not available")
		return result, nil
	}

	cmd := exec.CommandContext(ctx, "make")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	return result, nil
}

// buildMesonProject builds a Meson project
func (p *Impl) buildMesonProject(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langCpp,
		ProjectPath: projectPath,
	}

	// Check for meson
	if _, err := exec.LookPath("meson"); err != nil {
		result.Success = false
		result.Error = "meson not found in PATH"
		result.Warnings = append(result.Warnings, "Meson build tool not available")
		return result, nil
	}

	buildDir := filepath.Join(projectPath, "builddir")

	// Setup if needed
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		setupCmd := exec.CommandContext(ctx, "meson", "setup", buildDir)
		setupCmd.Dir = projectPath

		setupOutput, err := setupCmd.CombinedOutput()
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			result.Output = string(setupOutput)
			return result, nil
		}
	}

	// Compile
	compileCmd := exec.CommandContext(ctx, "meson", "compile", "-C", buildDir)
	compileCmd.Dir = projectPath

	output, err := compileCmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, "builddir/")

	return result, nil
}

// buildCppDirect compiles C/C++ files directly
func (p *Impl) buildCppDirect(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langCpp,
		ProjectPath: projectPath,
	}

	compiler := p.findCppCompiler()
	if compiler == "" {
		result.Success = false
		result.Error = "no C++ compiler found (g++, clang++)"
		result.Warnings = append(result.Warnings, "C++ compiler not available")
		return result, nil
	}

	// Find source files
	sourceFiles, err := p.findCppSourceFiles(projectPath)
	if err != nil || len(sourceFiles) == 0 {
		result.Success = false
		result.Error = "no C/C++ source files found"
		return result, nil
	}

	// Compile all source files
	outputPath := filepath.Join(projectPath, "a.out")
	args := append([]string{"-o", outputPath, "-std=c++17"}, sourceFiles...)

	cmd := exec.CommandContext(ctx, compiler, args...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, "a.out")

	return result, nil
}

// findCppCompiler finds an available C++ compiler
func (p *Impl) findCppCompiler() string {
	compilers := []string{"g++", "clang++", "c++"}
	for _, compiler := range compilers {
		if _, err := exec.LookPath(compiler); err == nil {
			return compiler
		}
	}
	return ""
}

// findCppSourceFiles finds C/C++ source files in project
func (p *Impl) findCppSourceFiles(projectPath string) ([]string, error) {
	var sourceFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip build directories
			if info.Name() == "build" || info.Name() == "builddir" ||
				info.Name() == "cmake-build-debug" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".cpp" || ext == ".cc" || ext == ".cxx" || ext == ".c" {
			relPath, err := filepath.Rel(projectPath, path)
			if err != nil {
				return nil
			}
			sourceFiles = append(sourceFiles, relPath)
		}

		return nil
	})

	return sourceFiles, err
}

// parseCppIssues parses C/C++ compiler error output
func (p *Impl) parseCppIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	lines := strings.Split(output, "\n")

	// gcc/clang format: file.cpp:10:5: error: message
	re := regexp.MustCompile(`([^:\s]+\.(c|cpp|cc|cxx|h|hpp)):(\d+):(\d+):\s*(error|warning|note):\s*(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 7 {
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[3]),
				Column:   p.parseInt(matches[4]),
				Severity: matches[5],
				Message:  matches[6],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// buildSwift выполняет сборку Swift проекта
func (p *Impl) buildSwift(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langSwift,
		ProjectPath: projectPath,
	}

	// Проверяем наличие swift
	if _, err := exec.LookPath("swift"); err != nil {
		result.Success = false
		result.Error = "swift not found in PATH"
		result.Warnings = append(result.Warnings, "Swift toolchain not available")
		return result, nil
	}

	// Проверяем наличие Package.swift (Swift Package Manager)
	packagePath := filepath.Join(projectPath, "Package.swift")
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "Package.swift not found"
		return result, nil
	}

	// Выполняем swift build
	cmd := exec.CommandContext(ctx, "swift", "build")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, ".build/")

	return result, nil
}

// typeCheckSwift выполняет проверку типов Swift проекта
func (p *Impl) typeCheckSwift(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langSwift,
		ProjectPath: projectPath,
	}

	// Проверяем наличие swift
	if _, err := exec.LookPath("swift"); err != nil {
		result.Success = false
		result.Error = "swift not found in PATH"
		return result, nil
	}

	// Проверяем наличие Package.swift
	packagePath := filepath.Join(projectPath, "Package.swift")
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "Package.swift not found"
		return result, nil
	}

	// swift build включает проверку типов
	cmd := exec.CommandContext(ctx, "swift", "build")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parseSwiftIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// parseSwiftIssues парсит ошибки Swift компилятора
func (p *Impl) parseSwiftIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	lines := strings.Split(output, "\n")

	// Swift error format: file.swift:10:5: error: message
	re := regexp.MustCompile(`([^:\s]+\.swift):(\d+):(\d+):\s*(error|warning|note):\s*(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 6 {
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[2]),
				Column:   p.parseInt(matches[3]),
				Severity: matches[4],
				Message:  matches[5],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// buildPHP выполняет сборку PHP проекта (composer install)
func (p *Impl) buildPHP(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langPHP,
		ProjectPath: projectPath,
	}

	// Проверяем наличие composer.json
	composerPath := filepath.Join(projectPath, "composer.json")
	if _, err := os.Stat(composerPath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "composer.json not found"
		return result, nil
	}

	// Проверяем наличие composer
	composerCmd := p.findComposerCommand()
	if composerCmd == "" {
		result.Success = false
		result.Error = "composer not found in PATH"
		result.Warnings = append(result.Warnings, "Composer not available")
		return result, nil
	}

	// Выполняем composer install
	cmd := exec.CommandContext(ctx, composerCmd, "install", "--no-interaction", "--prefer-dist")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, "vendor/")

	// Проверяем синтаксис PHP файлов
	syntaxOutput, syntaxErr := p.checkPHPSyntax(ctx, projectPath)
	if syntaxErr != nil {
		result.Output += "\n" + syntaxOutput
		result.Warnings = append(result.Warnings, "PHP syntax check failed: "+syntaxErr.Error())
	}

	return result, nil
}

// typeCheckPHP выполняет проверку типов PHP проекта через PHPStan или Psalm
func (p *Impl) typeCheckPHP(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langPHP,
		ProjectPath: projectPath,
	}

	// Проверяем наличие composer.json
	composerPath := filepath.Join(projectPath, "composer.json")
	if _, err := os.Stat(composerPath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "composer.json not found"
		return result, nil
	}

	// Пробуем PHPStan сначала
	phpstanPath := p.findPHPStanPath(projectPath)
	if phpstanPath != "" {
		return p.runPHPStan(ctx, projectPath, phpstanPath)
	}

	// Пробуем Psalm
	psalmPath := p.findPsalmPath(projectPath)
	if psalmPath != "" {
		return p.runPsalm(ctx, projectPath, psalmPath)
	}

	// Fallback: проверка синтаксиса через php -l
	result.Success = true
	result.Output = "No static analyzer found (PHPStan/Psalm). Running syntax check only."
	result.Issues = append(result.Issues, &domain.TypeIssue{
		Severity: "warning",
		Message:  "Install PHPStan or Psalm for type checking: composer require --dev phpstan/phpstan",
	})

	syntaxOutput, syntaxErr := p.checkPHPSyntax(ctx, projectPath)
	result.Output += "\n" + syntaxOutput

	if syntaxErr != nil {
		result.Success = false
		result.Error = syntaxErr.Error()
	}

	return result, nil
}

// findComposerCommand находит доступную команду Composer
func (p *Impl) findComposerCommand() string {
	// Пробуем composer сначала
	if _, err := exec.LookPath("composer"); err == nil {
		return "composer"
	}
	// Затем composer.phar
	if _, err := exec.LookPath("composer.phar"); err == nil {
		return "composer.phar"
	}
	return ""
}

// checkPHPSyntax проверяет синтаксис PHP файлов
func (p *Impl) checkPHPSyntax(ctx context.Context, projectPath string) (string, error) {
	// Проверяем наличие php
	if _, err := exec.LookPath("php"); err != nil {
		return "", fmt.Errorf("php not found in PATH")
	}

	var allOutput strings.Builder
	var syntaxErrors []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Пропускаем vendor и cache директории
			if info.Name() == "vendor" || info.Name() == ".git" ||
				info.Name() == "node_modules" || info.Name() == "cache" ||
				info.Name() == ".phpunit.cache" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".php") {
			return nil
		}

		// Проверяем синтаксис файла
		cmd := exec.CommandContext(ctx, "php", "-l", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			allOutput.WriteString(string(output))
			allOutput.WriteString("\n")
			syntaxErrors = append(syntaxErrors, path)
		}

		return nil
	})

	if err != nil {
		return allOutput.String(), fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(syntaxErrors) > 0 {
		return allOutput.String(), fmt.Errorf("syntax errors in %d file(s)", len(syntaxErrors))
	}

	return allOutput.String(), nil
}

// findPHPStanPath находит путь к PHPStan
func (p *Impl) findPHPStanPath(projectPath string) string {
	// Проверяем vendor/bin/phpstan
	vendorPath := filepath.Join(projectPath, "vendor", "bin", "phpstan")
	if _, err := os.Stat(vendorPath); err == nil {
		return vendorPath
	}

	// Проверяем глобальный phpstan
	if _, err := exec.LookPath("phpstan"); err == nil {
		return "phpstan"
	}

	return ""
}

// findPsalmPath находит путь к Psalm
func (p *Impl) findPsalmPath(projectPath string) string {
	// Проверяем vendor/bin/psalm
	vendorPath := filepath.Join(projectPath, "vendor", "bin", "psalm")
	if _, err := os.Stat(vendorPath); err == nil {
		return vendorPath
	}

	// Проверяем глобальный psalm
	if _, err := exec.LookPath("psalm"); err == nil {
		return "psalm"
	}

	return ""
}

// runPHPStan запускает PHPStan анализ
func (p *Impl) runPHPStan(ctx context.Context, projectPath, phpstanPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langPHP,
		ProjectPath: projectPath,
	}

	cmd := exec.CommandContext(ctx, phpstanPath, "analyse", "--error-format=json", "--no-progress")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parsePHPStanIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// runPsalm запускает Psalm анализ
func (p *Impl) runPsalm(ctx context.Context, projectPath, psalmPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langPHP,
		ProjectPath: projectPath,
	}

	cmd := exec.CommandContext(ctx, psalmPath, "--output-format=json", "--no-progress")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parsePsalmIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// parsePHPStanIssues парсит ошибки PHPStan
func (p *Impl) parsePHPStanIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// Пробуем парсить JSON формат
	var phpstanResult struct {
		Totals struct {
			Errors   int `json:"errors"`
			FileErrs int `json:"file_errors"`
		} `json:"totals"`
		Files map[string]struct {
			Errors   int `json:"errors"`
			Messages []struct {
				Message   string `json:"message"`
				Line      int    `json:"line"`
				Ignorable bool   `json:"ignorable"`
			} `json:"messages"`
		} `json:"files"`
	}

	if err := json.Unmarshal([]byte(output), &phpstanResult); err == nil {
		for filePath, fileData := range phpstanResult.Files {
			for _, msg := range fileData.Messages {
				issue := &domain.TypeIssue{
					File:     filePath,
					Line:     msg.Line,
					Severity: "error",
					Message:  msg.Message,
				}
				issues = append(issues, issue)
			}
		}
		return issues
	}

	// Fallback: парсим текстовый формат
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`([^:]+):(\d+):\s*(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 4 {
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[2]),
				Severity: "error",
				Message:  matches[3],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// parsePsalmIssues парсит ошибки Psalm
func (p *Impl) parsePsalmIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// Пробуем парсить JSON формат
	var psalmResults []struct {
		Severity    string `json:"severity"`
		Line        int    `json:"line_from"`
		Column      int    `json:"column_from"`
		Message     string `json:"message"`
		FileName    string `json:"file_name"`
		Type        string `json:"type"`
		ShortCode   int    `json:"shortcode"`
		ErrorLevel  int    `json:"error_level"`
	}

	if err := json.Unmarshal([]byte(output), &psalmResults); err == nil {
		for _, r := range psalmResults {
			issue := &domain.TypeIssue{
				File:     r.FileName,
				Line:     r.Line,
				Column:   r.Column,
				Severity: r.Severity,
				Message:  r.Message,
				Code:     r.Type,
			}
			issues = append(issues, issue)
		}
		return issues
	}

	// Fallback: парсим текстовый формат
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`([^:]+):(\d+):(\d+):\s*(error|warning):\s*(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 6 {
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[2]),
				Column:   p.parseInt(matches[3]),
				Severity: matches[4],
				Message:  matches[5],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// buildRuby выполняет сборку Ruby проекта (bundle install)
func (p *Impl) buildRuby(ctx context.Context, projectPath string) (*domain.BuildResult, error) {
	result := &domain.BuildResult{
		Language:    langRuby,
		ProjectPath: projectPath,
	}

	// Check for Gemfile
	gemfilePath := filepath.Join(projectPath, "Gemfile")
	if _, err := os.Stat(gemfilePath); os.IsNotExist(err) {
		result.Success = false
		result.Error = "Gemfile not found"
		return result, nil
	}

	// Check for bundler
	if _, err := exec.LookPath("bundle"); err != nil {
		result.Success = false
		result.Error = "bundler not found in PATH"
		result.Warnings = append(result.Warnings, "Install bundler with: gem install bundler")
		return result, nil
	}

	// Run bundle install
	cmd := exec.CommandContext(ctx, "bundle", "install")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	result.Artifacts = append(result.Artifacts, "Gemfile.lock")

	// Run syntax check on Ruby files
	syntaxOutput, syntaxErr := p.checkRubySyntax(ctx, projectPath)
	if syntaxErr != nil {
		result.Output += "\n" + syntaxOutput
		result.Warnings = append(result.Warnings, "Ruby syntax check found issues")
	}

	return result, nil
}

// typeCheckRuby выполняет проверку типов Ruby проекта через Sorbet или RuboCop
func (p *Impl) typeCheckRuby(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langRuby,
		ProjectPath: projectPath,
	}

	// Try Sorbet first (if available)
	if p.hasSorbet(projectPath) {
		return p.typeCheckWithSorbet(ctx, projectPath)
	}

	// Fallback to RuboCop
	if p.hasRuboCop(projectPath) {
		return p.typeCheckWithRuboCop(ctx, projectPath)
	}

	// Basic syntax check
	syntaxOutput, syntaxErr := p.checkRubySyntax(ctx, projectPath)
	result.Output = syntaxOutput

	if syntaxErr != nil {
		result.Success = false
		result.Error = syntaxErr.Error()
		return result, nil
	}

	result.Success = true
	return result, nil
}

// checkRubySyntax checks Ruby syntax using ruby -c
func (p *Impl) checkRubySyntax(ctx context.Context, projectPath string) (string, error) {
	var allOutput strings.Builder
	var syntaxErrors []string

	// Find Ruby command
	rubyCmd := p.findRubyCommand()
	if rubyCmd == "" {
		return "", fmt.Errorf("ruby not found in PATH")
	}

	// Find all .rb files
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip vendor and cache directories
			if info.Name() == "vendor" || info.Name() == ".bundle" ||
				info.Name() == ".git" || info.Name() == "node_modules" ||
				info.Name() == "tmp" || info.Name() == "log" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".rb") {
			return nil
		}

		// Check syntax
		cmd := exec.CommandContext(ctx, rubyCmd, "-c", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			allOutput.WriteString(string(output))
			allOutput.WriteString("\n")
			syntaxErrors = append(syntaxErrors, path)
		}

		return nil
	})

	if err != nil {
		return allOutput.String(), fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(syntaxErrors) > 0 {
		return allOutput.String(), fmt.Errorf("syntax errors in %d file(s)", len(syntaxErrors))
	}

	return allOutput.String(), nil
}

// findRubyCommand finds the Ruby executable
func (p *Impl) findRubyCommand() string {
	if _, err := exec.LookPath("ruby"); err == nil {
		return "ruby"
	}
	return ""
}

// hasSorbet checks if Sorbet is available in the project
func (p *Impl) hasSorbet(projectPath string) bool {
	// Check for sorbet directory
	sorbetDir := filepath.Join(projectPath, "sorbet")
	if _, err := os.Stat(sorbetDir); err == nil {
		return true
	}

	// Check Gemfile for sorbet
	gemfilePath := filepath.Join(projectPath, "Gemfile")
	content, err := os.ReadFile(gemfilePath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "sorbet")
}

// hasRuboCop checks if RuboCop is available
func (p *Impl) hasRuboCop(projectPath string) bool {
	// Check for .rubocop.yml
	rubocopConfig := filepath.Join(projectPath, ".rubocop.yml")
	if _, err := os.Stat(rubocopConfig); err == nil {
		return true
	}

	// Check Gemfile for rubocop
	gemfilePath := filepath.Join(projectPath, "Gemfile")
	content, err := os.ReadFile(gemfilePath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "rubocop")
}

// typeCheckWithSorbet runs Sorbet type checker
func (p *Impl) typeCheckWithSorbet(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langRuby,
		ProjectPath: projectPath,
	}

	// Determine command
	cmdName := "srb"
	args := []string{"tc"}

	// Use bundle exec if Gemfile exists
	if _, err := os.Stat(filepath.Join(projectPath, "Gemfile")); err == nil {
		cmdName = "bundle"
		args = []string{"exec", "srb", "tc"}
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parseSorbetIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// typeCheckWithRuboCop runs RuboCop for type-like checks
func (p *Impl) typeCheckWithRuboCop(ctx context.Context, projectPath string) (*domain.TypeCheckResult, error) {
	result := &domain.TypeCheckResult{
		Language:    langRuby,
		ProjectPath: projectPath,
	}

	// Determine command
	cmdName := "rubocop"
	args := []string{"--format", "json"}

	// Use bundle exec if Gemfile exists
	if _, err := os.Stat(filepath.Join(projectPath, "Gemfile")); err == nil {
		cmdName = "bundle"
		args = []string{"exec", "rubocop", "--format", "json"}
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Issues = p.parseRuboCopIssues(string(output))
		return result, nil
	}

	result.Success = true
	return result, nil
}

// parseSorbetIssues parses Sorbet output
func (p *Impl) parseSorbetIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// Sorbet format: "path/file.rb:10: error: Message"
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`^([^:]+):(\d+):\s*(error|warning):\s*(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 5 {
			issue := &domain.TypeIssue{
				File:     matches[1],
				Line:     p.parseInt(matches[2]),
				Severity: matches[3],
				Message:  matches[4],
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

// parseRuboCopIssues parses RuboCop JSON output
func (p *Impl) parseRuboCopIssues(output string) []*domain.TypeIssue {
	var issues []*domain.TypeIssue

	// Find JSON in output
	jsonStart := strings.Index(output, "{")
	if jsonStart == -1 {
		return issues
	}

	jsonStr := output[jsonStart:]

	var report struct {
		Files []struct {
			Path     string `json:"path"`
			Offenses []struct {
				Severity string `json:"severity"`
				Message  string `json:"message"`
				CopName  string `json:"cop_name"`
				Location struct {
					Line   int `json:"line"`
					Column int `json:"column"`
				} `json:"location"`
			} `json:"offenses"`
		} `json:"files"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &report); err != nil {
		return issues
	}

	for _, file := range report.Files {
		for _, offense := range file.Offenses {
			issue := &domain.TypeIssue{
				File:     file.Path,
				Line:     offense.Location.Line,
				Column:   offense.Location.Column,
				Severity: offense.Severity,
				Code:     offense.CopName,
				Message:  offense.Message,
			}
			issues = append(issues, issue)
		}
	}

	return issues
}
