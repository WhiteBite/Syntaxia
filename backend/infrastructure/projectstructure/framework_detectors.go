package projectstructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

type frameworkDetector interface {
	detect(projectPath string) *domain.FrameworkInfo
}

func (d *Detector) initFrameworkDetectors() {
	d.frameworkDetectors = []frameworkDetector{
		// Frontend frameworks
		&vueDetector{},
		&reactDetector{},
		&angularDetector{},
		&svelteDetector{},
		&nextjsDetector{},
		&nuxtDetector{},

		// Backend frameworks - Go
		&ginDetector{},
		&echoDetector{},
		&fiberDetector{},
		&wailsDetector{},

		// Backend frameworks - Node
		&expressDetector{},
		&nestjsDetector{},
		&fastifyDetector{},

		// Backend frameworks - Python
		&djangoDetector{},
		&flaskDetector{},
		&fastapiDetector{},

		// Backend frameworks - Java/Kotlin
		&springDetector{},
		&ktorDetector{},

		// Mobile
		&flutterDetector{},
		&reactNativeDetector{},

		// Testing
		&jestDetector{},
		&vitestDetector{},
		&pytestDetector{},
	}
}

// Vue.js detector
type vueDetector struct{}

func (v *vueDetector) detect(projectPath string) *domain.FrameworkInfo {
	pkgPath := filepath.Join(projectPath, "package.json")
	content, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), `"vue"`) {
		return nil
	}

	version := extractDependencyVersion(content, "vue")

	var configFiles []string
	configs := []string{"vite.config.ts", "vite.config.js", "vue.config.js", "nuxt.config.ts"}
	for _, cfg := range configs {
		if _, err := os.Stat(filepath.Join(projectPath, cfg)); err == nil {
			configFiles = append(configFiles, cfg)
		}
	}

	return &domain.FrameworkInfo{
		Name:        "Vue.js",
		Version:     version,
		Category:    "web",
		Language:    "TypeScript",
		ConfigFiles: configFiles,
		Indicators:  []string{"vue dependency in package.json"},
		BestPractices: []string{
			"Use Composition API for complex components",
			"Organize components by feature in src/features/",
			"Use composables for reusable logic",
			"Keep components small and focused",
			"Use TypeScript for type safety",
		},
	}
}

// React detector
type reactDetector struct{}

func (r *reactDetector) detect(projectPath string) *domain.FrameworkInfo {
	pkgPath := filepath.Join(projectPath, "package.json")
	content, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), `"react"`) {
		return nil
	}

	version := extractDependencyVersion(content, "react")

	return &domain.FrameworkInfo{
		Name:       "React",
		Version:    version,
		Category:   "web",
		Language:   "TypeScript",
		Indicators: []string{"react dependency in package.json"},
		BestPractices: []string{
			"Use functional components with hooks",
			"Organize by feature, not by type",
			"Use custom hooks for reusable logic",
			"Implement proper error boundaries",
			"Use React.memo for performance optimization",
		},
	}
}

// Angular detector
type angularDetector struct{}

func (a *angularDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "angular.json")); err != nil {
		return nil
	}

	pkgPath := filepath.Join(projectPath, "package.json")
	content, _ := os.ReadFile(pkgPath)
	version := extractDependencyVersion(content, "@angular/core")

	return &domain.FrameworkInfo{
		Name:        "Angular",
		Version:     version,
		Category:    "web",
		Language:    "TypeScript",
		ConfigFiles: []string{"angular.json", "tsconfig.json"},
		Indicators:  []string{"angular.json config file"},
		BestPractices: []string{
			"Use standalone components",
			"Implement lazy loading for modules",
			"Use OnPush change detection",
			"Follow Angular style guide",
		},
	}
}

// Svelte detector
type svelteDetector struct{}

func (s *svelteDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "svelte.config.js")); err != nil {
		pkgPath := filepath.Join(projectPath, "package.json")
		content, err := os.ReadFile(pkgPath)
		if err != nil || !strings.Contains(string(content), `"svelte"`) {
			return nil
		}
	}

	return &domain.FrameworkInfo{
		Name:        "Svelte",
		Category:    "web",
		Language:    "TypeScript",
		ConfigFiles: []string{"svelte.config.js"},
		Indicators:  []string{"svelte.config.js or svelte dependency"},
		BestPractices: []string{
			"Use stores for shared state",
			"Keep components small",
			"Use TypeScript for type safety",
		},
	}
}

// Next.js detector
type nextjsDetector struct{}

func (n *nextjsDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "next.config.js")); err != nil {
		if _, err := os.Stat(filepath.Join(projectPath, "next.config.mjs")); err != nil {
			return nil
		}
	}

	return &domain.FrameworkInfo{
		Name:        "Next.js",
		Category:    "web",
		Language:    "TypeScript",
		ConfigFiles: []string{"next.config.js", "next.config.mjs"},
		Indicators:  []string{"next.config.js config file"},
		BestPractices: []string{
			"Use App Router for new projects",
			"Implement proper data fetching patterns",
			"Use Server Components where possible",
			"Optimize images with next/image",
		},
	}
}

// Nuxt detector
type nuxtDetector struct{}

func (n *nuxtDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "nuxt.config.ts")); err != nil {
		if _, err := os.Stat(filepath.Join(projectPath, "nuxt.config.js")); err != nil {
			return nil
		}
	}

	return &domain.FrameworkInfo{
		Name:        "Nuxt",
		Category:    "web",
		Language:    "TypeScript",
		ConfigFiles: []string{"nuxt.config.ts", "nuxt.config.js"},
		Indicators:  []string{"nuxt.config.ts config file"},
		BestPractices: []string{
			"Use auto-imports feature",
			"Organize pages by route structure",
			"Use composables for shared logic",
			"Implement proper SEO with useHead",
		},
	}
}

// Gin (Go) detector
type ginDetector struct{}

func (g *ginDetector) detect(projectPath string) *domain.FrameworkInfo {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), "github.com/gin-gonic/gin") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Gin",
		Category:    "web",
		Language:    "Go",
		ConfigFiles: []string{"go.mod"},
		Indicators:  []string{"gin-gonic/gin in go.mod"},
		BestPractices: []string{
			"Use middleware for cross-cutting concerns",
			"Implement proper error handling",
			"Use gin.Context for request/response",
			"Group routes logically",
			"Use binding for request validation",
		},
	}
}

// Echo (Go) detector
type echoDetector struct{}

func (e *echoDetector) detect(projectPath string) *domain.FrameworkInfo {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), "github.com/labstack/echo") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Echo",
		Category:    "web",
		Language:    "Go",
		ConfigFiles: []string{"go.mod"},
		Indicators:  []string{"labstack/echo in go.mod"},
		BestPractices: []string{
			"Use middleware for logging, auth, etc.",
			"Implement custom error handler",
			"Use context for request data",
			"Group routes with prefixes",
		},
	}
}

// Fiber (Go) detector
type fiberDetector struct{}

func (f *fiberDetector) detect(projectPath string) *domain.FrameworkInfo {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), "github.com/gofiber/fiber") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Fiber",
		Category:    "web",
		Language:    "Go",
		ConfigFiles: []string{"go.mod"},
		Indicators:  []string{"gofiber/fiber in go.mod"},
		BestPractices: []string{
			"Use built-in middleware",
			"Implement proper error handling",
			"Use fiber.Ctx for request handling",
		},
	}
}

// Wails detector
type wailsDetector struct{}

func (w *wailsDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "wails.json")); err != nil {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Wails",
		Category:    "desktop",
		Language:    "Go",
		ConfigFiles: []string{"wails.json"},
		Indicators:  []string{"wails.json config file"},
		BestPractices: []string{
			"Use Go for backend logic",
			"Use frontend framework for UI",
			"Implement proper error handling between Go and JS",
			"Use events for async communication",
		},
	}
}

// Express detector
type expressDetector struct{}

func (e *expressDetector) detect(projectPath string) *domain.FrameworkInfo {
	pkgPath := filepath.Join(projectPath, "package.json")
	content, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), `"express"`) {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:       "Express",
		Category:   "web",
		Language:   "JavaScript",
		Indicators: []string{"express dependency in package.json"},
		BestPractices: []string{
			"Use middleware for cross-cutting concerns",
			"Implement proper error handling middleware",
			"Use router for route organization",
			"Validate input with express-validator",
		},
	}
}

// NestJS detector
type nestjsDetector struct{}

func (n *nestjsDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "nest-cli.json")); err != nil {
		pkgPath := filepath.Join(projectPath, "package.json")
		content, err := os.ReadFile(pkgPath)
		if err != nil || !strings.Contains(string(content), `"@nestjs/core"`) {
			return nil
		}
	}

	return &domain.FrameworkInfo{
		Name:        "NestJS",
		Category:    "web",
		Language:    "TypeScript",
		ConfigFiles: []string{"nest-cli.json"},
		Indicators:  []string{"nest-cli.json or @nestjs/core dependency"},
		BestPractices: []string{
			"Use modules for feature organization",
			"Implement DTOs for validation",
			"Use dependency injection",
			"Follow NestJS conventions",
		},
	}
}

// Fastify detector
type fastifyDetector struct{}

func (f *fastifyDetector) detect(projectPath string) *domain.FrameworkInfo {
	pkgPath := filepath.Join(projectPath, "package.json")
	content, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), `"fastify"`) {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:       "Fastify",
		Category:   "web",
		Language:   "TypeScript",
		Indicators: []string{"fastify dependency in package.json"},
		BestPractices: []string{
			"Use plugins for modularity",
			"Implement JSON Schema validation",
			"Use hooks for lifecycle events",
		},
	}
}

// Django detector
type djangoDetector struct{}

func (d *djangoDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "manage.py")); err != nil {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Django",
		Category:    "web",
		Language:    "Python",
		ConfigFiles: []string{"manage.py", "settings.py"},
		Indicators:  []string{"manage.py file"},
		BestPractices: []string{
			"Use Django REST Framework for APIs",
			"Implement proper model design",
			"Use class-based views",
			"Follow Django conventions",
		},
	}
}

// Flask detector
type flaskDetector struct{}

func (f *flaskDetector) detect(projectPath string) *domain.FrameworkInfo {
	reqPath := filepath.Join(projectPath, "requirements.txt")
	content, err := os.ReadFile(reqPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(strings.ToLower(string(content)), "flask") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:       "Flask",
		Category:   "web",
		Language:   "Python",
		Indicators: []string{"flask in requirements.txt"},
		BestPractices: []string{
			"Use blueprints for modularity",
			"Implement proper error handling",
			"Use Flask-SQLAlchemy for ORM",
		},
	}
}

// FastAPI detector
type fastapiDetector struct{}

func (f *fastapiDetector) detect(projectPath string) *domain.FrameworkInfo {
	reqPath := filepath.Join(projectPath, "requirements.txt")
	content, err := os.ReadFile(reqPath)
	if err != nil {
		// Check pyproject.toml
		pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
		content, err = os.ReadFile(pyprojectPath)
		if err != nil {
			return nil
		}
	}

	if !strings.Contains(strings.ToLower(string(content)), "fastapi") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:       "FastAPI",
		Category:   "web",
		Language:   "Python",
		Indicators: []string{"fastapi in requirements.txt or pyproject.toml"},
		BestPractices: []string{
			"Use Pydantic models for validation",
			"Implement dependency injection",
			"Use async/await for I/O operations",
			"Document with OpenAPI",
		},
	}
}

// Spring detector
type springDetector struct{}

func (s *springDetector) detect(projectPath string) *domain.FrameworkInfo {
	pomPath := filepath.Join(projectPath, "pom.xml")
	content, err := os.ReadFile(pomPath)
	if err != nil {
		gradlePath := filepath.Join(projectPath, "build.gradle")
		content, err = os.ReadFile(gradlePath)
		if err != nil {
			return nil
		}
	}

	if !strings.Contains(string(content), "spring") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Spring Boot",
		Category:    "web",
		Language:    "Java",
		ConfigFiles: []string{"pom.xml", "build.gradle", "application.properties"},
		Indicators:  []string{"spring dependency in pom.xml or build.gradle"},
		BestPractices: []string{
			"Use Spring Boot starters",
			"Implement proper layered architecture",
			"Use Spring Data for repositories",
			"Configure with application.properties/yaml",
		},
	}
}

// Ktor detector
type ktorDetector struct{}

func (k *ktorDetector) detect(projectPath string) *domain.FrameworkInfo {
	gradlePath := filepath.Join(projectPath, "build.gradle.kts")
	content, err := os.ReadFile(gradlePath)
	if err != nil {
		gradlePath = filepath.Join(projectPath, "build.gradle")
		content, err = os.ReadFile(gradlePath)
		if err != nil {
			return nil
		}
	}

	if !strings.Contains(string(content), "ktor") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Ktor",
		Category:    "web",
		Language:    "Kotlin",
		ConfigFiles: []string{"build.gradle.kts", "application.conf"},
		Indicators:  []string{"ktor dependency in build.gradle"},
		BestPractices: []string{
			"Use features for modularity",
			"Implement proper routing",
			"Use coroutines for async operations",
		},
	}
}

// Flutter detector
type flutterDetector struct{}

func (f *flutterDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "pubspec.yaml")); err != nil {
		return nil
	}

	content, _ := os.ReadFile(filepath.Join(projectPath, "pubspec.yaml"))
	if !strings.Contains(string(content), "flutter:") {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:        "Flutter",
		Category:    "mobile",
		Language:    "Dart",
		ConfigFiles: []string{"pubspec.yaml"},
		Indicators:  []string{"flutter in pubspec.yaml"},
		BestPractices: []string{
			"Use BLoC or Provider for state management",
			"Organize by feature",
			"Use const constructors",
			"Implement proper widget composition",
		},
	}
}

// React Native detector
type reactNativeDetector struct{}

func (r *reactNativeDetector) detect(projectPath string) *domain.FrameworkInfo {
	pkgPath := filepath.Join(projectPath, "package.json")
	content, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	if !strings.Contains(string(content), `"react-native"`) {
		return nil
	}

	return &domain.FrameworkInfo{
		Name:       "React Native",
		Category:   "mobile",
		Language:   "TypeScript",
		Indicators: []string{"react-native dependency in package.json"},
		BestPractices: []string{
			"Use TypeScript for type safety",
			"Implement proper navigation",
			"Use React Query for data fetching",
			"Optimize with memo and useMemo",
		},
	}
}

// Jest detector
type jestDetector struct{}

func (j *jestDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "jest.config.js")); err != nil {
		if _, err := os.Stat(filepath.Join(projectPath, "jest.config.ts")); err != nil {
			return nil
		}
	}

	return &domain.FrameworkInfo{
		Name:        "Jest",
		Category:    "testing",
		Language:    "TypeScript",
		ConfigFiles: []string{"jest.config.js", "jest.config.ts"},
		Indicators:  []string{"jest.config.js config file"},
		BestPractices: []string{
			"Use describe/it for test organization",
			"Mock external dependencies",
			"Use snapshot testing sparingly",
			"Aim for high coverage",
		},
	}
}

// Vitest detector
type vitestDetector struct{}

func (v *vitestDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "vitest.config.ts")); err != nil {
		if _, err := os.Stat(filepath.Join(projectPath, "vitest.config.js")); err != nil {
			return nil
		}
	}

	return &domain.FrameworkInfo{
		Name:        "Vitest",
		Category:    "testing",
		Language:    "TypeScript",
		ConfigFiles: []string{"vitest.config.ts", "vitest.config.js"},
		Indicators:  []string{"vitest.config.ts config file"},
		BestPractices: []string{
			"Use describe/it for test organization",
			"Use vi.mock for mocking",
			"Run in watch mode during development",
		},
	}
}

// Pytest detector
type pytestDetector struct{}

func (p *pytestDetector) detect(projectPath string) *domain.FrameworkInfo {
	if _, err := os.Stat(filepath.Join(projectPath, "pytest.ini")); err != nil {
		if _, err := os.Stat(filepath.Join(projectPath, "pyproject.toml")); err != nil {
			return nil
		}
		content, _ := os.ReadFile(filepath.Join(projectPath, "pyproject.toml"))
		if !strings.Contains(string(content), "[tool.pytest") {
			return nil
		}
	}

	return &domain.FrameworkInfo{
		Name:        "pytest",
		Category:    "testing",
		Language:    "Python",
		ConfigFiles: []string{"pytest.ini", "pyproject.toml"},
		Indicators:  []string{"pytest.ini or pytest config in pyproject.toml"},
		BestPractices: []string{
			"Use fixtures for setup/teardown",
			"Use parametrize for test variations",
			"Organize tests by feature",
		},
	}
}

// Helper function to extract dependency version from package.json
func extractDependencyVersion(content []byte, depName string) string {
	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return ""
	}

	// Check dependencies
	if deps, ok := pkg["dependencies"].(map[string]interface{}); ok {
		if version, ok := deps[depName].(string); ok {
			return strings.TrimPrefix(version, "^")
		}
	}

	// Check devDependencies
	if deps, ok := pkg["devDependencies"].(map[string]interface{}); ok {
		if version, ok := deps[depName].(string); ok {
			return strings.TrimPrefix(version, "^")
		}
	}

	return ""
}
