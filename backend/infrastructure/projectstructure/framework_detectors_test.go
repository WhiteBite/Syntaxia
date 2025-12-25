package projectstructure

import (
	"testing"
)

const categoryWebTest = "web"

func TestVueDetector(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Vue project
	createTestFile(t, tmpDir, "package.json", `{
		"name": "vue-app",
		"dependencies": {
			"vue": "^3.4.0"
		}
	}`)
	createTestFile(t, tmpDir, "vite.config.ts", "export default {}")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Vue.js" {
			found = true
			if fw.Category != categoryWebTest {
				t.Errorf("Expected category 'web', got '%s'", fw.Category)
			}
			if fw.Version == "" {
				t.Error("Version not detected")
			}
			if len(fw.BestPractices) == 0 {
				t.Error("No best practices provided")
			}
		}
	}
	if !found {
		t.Error("Vue.js not detected")
	}
}

func TestReactDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "package.json", `{
		"name": "react-app",
		"dependencies": {
			"react": "^18.2.0",
			"react-dom": "^18.2.0"
		}
	}`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "React" {
			found = true
			if fw.Category != categoryWebTest {
				t.Errorf("Expected category 'web', got '%s'", fw.Category)
			}
		}
	}
	if !found {
		t.Error("React not detected")
	}
}

func TestAngularDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "angular.json", `{"version": 1}`)
	createTestFile(t, tmpDir, "package.json", `{
		"dependencies": {
			"@angular/core": "^17.0.0"
		}
	}`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Angular" {
			found = true
		}
	}
	if !found {
		t.Error("Angular not detected")
	}
}

func TestGinDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "go.mod", `module myapp

go 1.21

require github.com/gin-gonic/gin v1.9.1
`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Gin" {
			found = true
			if fw.Language != "Go" {
				t.Errorf("Expected language 'Go', got '%s'", fw.Language)
			}
			if fw.Category != categoryWebTest {
				t.Errorf("Expected category 'web', got '%s'", fw.Category)
			}
		}
	}
	if !found {
		t.Error("Gin not detected")
	}
}

func TestEchoDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "go.mod", `module myapp

go 1.21

require github.com/labstack/echo/v4 v4.11.0
`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Echo" {
			found = true
		}
	}
	if !found {
		t.Error("Echo not detected")
	}
}

func TestWailsDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "wails.json", `{"name": "myapp"}`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Wails" {
			found = true
			if fw.Category != "desktop" {
				t.Errorf("Expected category 'desktop', got '%s'", fw.Category)
			}
		}
	}
	if !found {
		t.Error("Wails not detected")
	}
}

func TestExpressDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "package.json", `{
		"dependencies": {
			"express": "^4.18.0"
		}
	}`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Express" {
			found = true
		}
	}
	if !found {
		t.Error("Express not detected")
	}
}

func TestNestJSDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "nest-cli.json", `{}`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "NestJS" {
			found = true
		}
	}
	if !found {
		t.Error("NestJS not detected")
	}
}

func TestDjangoDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "manage.py", "#!/usr/bin/env python")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Django" {
			found = true
			if fw.Language != "Python" {
				t.Errorf("Expected language 'Python', got '%s'", fw.Language)
			}
		}
	}
	if !found {
		t.Error("Django not detected")
	}
}

func TestFlaskDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "requirements.txt", "Flask==2.3.0\nrequests==2.31.0")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Flask" {
			found = true
		}
	}
	if !found {
		t.Error("Flask not detected")
	}
}

func TestFastAPIDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "requirements.txt", "fastapi==0.100.0\nuvicorn==0.23.0")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "FastAPI" {
			found = true
		}
	}
	if !found {
		t.Error("FastAPI not detected")
	}
}

func TestSpringDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "pom.xml", `<?xml version="1.0"?>
<project>
	<dependencies>
		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter</artifactId>
		</dependency>
	</dependencies>
</project>`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Spring Boot" {
			found = true
			if fw.Language != "Java" {
				t.Errorf("Expected language 'Java', got '%s'", fw.Language)
			}
		}
	}
	if !found {
		t.Error("Spring Boot not detected")
	}
}

func TestFlutterDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "pubspec.yaml", `name: myapp
flutter:
  sdk: flutter
`)

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Flutter" {
			found = true
			if fw.Category != "mobile" {
				t.Errorf("Expected category 'mobile', got '%s'", fw.Category)
			}
			if fw.Language != "Dart" {
				t.Errorf("Expected language 'Dart', got '%s'", fw.Language)
			}
		}
	}
	if !found {
		t.Error("Flutter not detected")
	}
}

func TestJestDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "jest.config.js", "module.exports = {}")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Jest" {
			found = true
			if fw.Category != "testing" {
				t.Errorf("Expected category 'testing', got '%s'", fw.Category)
			}
		}
	}
	if !found {
		t.Error("Jest not detected")
	}
}

func TestVitestDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "vitest.config.ts", "export default {}")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Vitest" {
			found = true
		}
	}
	if !found {
		t.Error("Vitest not detected")
	}
}

func TestNextJSDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "next.config.js", "module.exports = {}")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Next.js" {
			found = true
		}
	}
	if !found {
		t.Error("Next.js not detected")
	}
}

func TestNuxtDetector(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "nuxt.config.ts", "export default defineNuxtConfig({})")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Nuxt" {
			found = true
		}
	}
	if !found {
		t.Error("Nuxt not detected")
	}
}

func TestMultipleFrameworks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create project with multiple frameworks
	createTestFile(t, tmpDir, "package.json", `{
		"dependencies": {
			"vue": "^3.4.0"
		}
	}`)
	createTestFile(t, tmpDir, "vitest.config.ts", "export default {}")
	createTestFile(t, tmpDir, "vite.config.ts", "export default {}")

	d := NewDetector()
	frameworks, _ := d.DetectFrameworks(tmpDir)

	if len(frameworks) < 2 {
		t.Errorf("Expected at least 2 frameworks, got %d", len(frameworks))
	}

	hasVue := false
	hasVitest := false
	for _, fw := range frameworks {
		if fw.Name == "Vue.js" {
			hasVue = true
		}
		if fw.Name == "Vitest" {
			hasVitest = true
		}
	}

	if !hasVue {
		t.Error("Vue.js not detected in multi-framework project")
	}
	if !hasVitest {
		t.Error("Vitest not detected in multi-framework project")
	}
}

func TestExtractDependencyVersion(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		dep      string
		expected string
	}{
		{
			name:     "caret version",
			content:  `{"dependencies": {"vue": "^3.4.0"}}`,
			dep:      "vue",
			expected: "3.4.0",
		},
		{
			name:     "exact version",
			content:  `{"dependencies": {"react": "18.2.0"}}`,
			dep:      "react",
			expected: "18.2.0",
		},
		{
			name:     "devDependencies",
			content:  `{"devDependencies": {"jest": "^29.0.0"}}`,
			dep:      "jest",
			expected: "29.0.0",
		},
		{
			name:     "not found",
			content:  `{"dependencies": {"vue": "^3.4.0"}}`,
			dep:      "react",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := extractDependencyVersion([]byte(tt.content), tt.dep)
			if version != tt.expected {
				t.Errorf("Expected version '%s', got '%s'", tt.expected, version)
			}
		})
	}
}

func TestNoFrameworksDetected(t *testing.T) {
	tmpDir := t.TempDir()

	// Create empty project
	createTestFile(t, tmpDir, "main.go", "package main\n\nfunc main() {}")

	d := NewDetector()
	frameworks, err := d.DetectFrameworks(tmpDir)
	if err != nil {
		t.Fatalf("DetectFrameworks failed: %v", err)
	}

	// Should return empty or nil slice (both are acceptable)
	if len(frameworks) > 0 {
		t.Errorf("Expected no frameworks, got %d", len(frameworks))
	}
}
