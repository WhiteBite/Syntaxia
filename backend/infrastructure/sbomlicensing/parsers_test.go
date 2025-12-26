package sbomlicensing

import (
	"testing"

	"syntaxia/domain"
)

// mockLogger implements domain.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(msg string)                                       {}
func (m *mockLogger) Info(msg string)                                        {}
func (m *mockLogger) Warning(msg string)                                     {}
func (m *mockLogger) Error(msg string)                                       {}
func (m *mockLogger) Fatal(msg string)                                       {}
func (m *mockLogger) Debugf(format string, args ...interface{})              {}
func (m *mockLogger) Infof(format string, args ...interface{})               {}
func (m *mockLogger) Warningf(format string, args ...interface{})            {}
func (m *mockLogger) Errorf(format string, args ...interface{})              {}
func (m *mockLogger) Fatalf(format string, args ...interface{})              {}
func (m *mockLogger) WithField(key string, value interface{}) domain.Logger  { return m }
func (m *mockLogger) WithFields(fields map[string]interface{}) domain.Logger { return m }

func TestGrypeScanner_parseGrypeOutput(t *testing.T) {
	log := &mockLogger{}
	scanner := NewGrypeScanner(log)

	tests := []struct {
		name           string
		input          string
		expectedCount  int
		expectedFirst  *domain.Vulnerability
		wantErr        bool
	}{
		{
			name:          "empty output",
			input:         "",
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "single vulnerability",
			input: `{
				"matches": [
					{
						"vulnerability": {
							"id": "CVE-2021-1234",
							"severity": "HIGH",
							"description": "Test vulnerability",
							"fix": {
								"versions": ["1.2.3"],
								"state": "fixed"
							},
							"cvss": [{"version": "3.0", "metrics": 7.5}]
						},
						"artifact": {
							"name": "test-package",
							"version": "1.0.0"
						}
					}
				]
			}`,
			expectedCount: 1,
			expectedFirst: &domain.Vulnerability{
				ID:          "CVE-2021-1234",
				Severity:    "HIGH",
				Description: "Test vulnerability",
				FixedIn:     "1.2.3",
				CVSS:        7.5,
			},
			wantErr: false,
		},
		{
			name: "multiple vulnerabilities",
			input: `{
				"matches": [
					{
						"vulnerability": {"id": "CVE-2021-0001", "severity": "CRITICAL", "description": "Critical vuln"},
						"artifact": {"name": "pkg1", "version": "1.0.0"}
					},
					{
						"vulnerability": {"id": "CVE-2021-0002", "severity": "LOW", "description": "Low vuln"},
						"artifact": {"name": "pkg2", "version": "2.0.0"}
					}
				]
			}`,
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name:          "invalid json",
			input:         "not valid json",
			expectedCount: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := scanner.parseGrypeOutput([]byte(tt.input))

			if (err != nil) != tt.wantErr {
				t.Errorf("parseGrypeOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("parseGrypeOutput() count = %d, expected %d", len(result), tt.expectedCount)
			}

			if tt.expectedFirst != nil && len(result) > 0 {
				if result[0].ID != tt.expectedFirst.ID {
					t.Errorf("ID = %s, expected %s", result[0].ID, tt.expectedFirst.ID)
				}
				if result[0].Severity != tt.expectedFirst.Severity {
					t.Errorf("Severity = %s, expected %s", result[0].Severity, tt.expectedFirst.Severity)
				}
				if result[0].FixedIn != tt.expectedFirst.FixedIn {
					t.Errorf("FixedIn = %s, expected %s", result[0].FixedIn, tt.expectedFirst.FixedIn)
				}
			}
		})
	}
}

func TestSyftGenerator_parseSyftJSON(t *testing.T) {
	log := &mockLogger{}
	generator := NewSyftGenerator(log)

	tests := []struct {
		name          string
		input         string
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "empty output",
			input:         "",
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "single artifact",
			input: `{
				"artifacts": [
					{
						"name": "test-package",
						"version": "1.0.0",
						"type": "go-module",
						"purl": "pkg:golang/test-package@1.0.0",
						"licenses": [{"value": "MIT", "spdxExpression": "MIT"}]
					}
				]
			}`,
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "multiple artifacts",
			input: `{
				"artifacts": [
					{"name": "pkg1", "version": "1.0.0", "type": "npm"},
					{"name": "pkg2", "version": "2.0.0", "type": "npm"}
				]
			}`,
			expectedCount: 2,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.parseSyftOutput([]byte(tt.input), domain.SBOMFormatJSON)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseSyftOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("parseSyftOutput() count = %d, expected %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestSyftGenerator_parseSyftSPDX(t *testing.T) {
	log := &mockLogger{}
	generator := NewSyftGenerator(log)

	input := `{
		"packages": [
			{
				"name": "test-package",
				"versionInfo": "1.0.0",
				"SPDXID": "SPDXRef-Package-test",
				"licenseConcluded": "MIT",
				"licenseDeclared": "MIT",
				"externalRefs": [
					{"referenceType": "purl", "referenceLocator": "pkg:npm/test@1.0.0"}
				]
			}
		]
	}`

	result, err := generator.parseSyftOutput([]byte(input), domain.SBOMFormatSPDX)
	if err != nil {
		t.Errorf("parseSyftOutput(SPDX) error = %v", err)
		return
	}

	if len(result) != 1 {
		t.Errorf("parseSyftOutput(SPDX) count = %d, expected 1", len(result))
		return
	}

	if result[0].Name != "test-package" {
		t.Errorf("Name = %s, expected test-package", result[0].Name)
	}
	if result[0].License != "MIT" {
		t.Errorf("License = %s, expected MIT", result[0].License)
	}
	if result[0].PURL != "pkg:npm/test@1.0.0" {
		t.Errorf("PURL = %s, expected pkg:npm/test@1.0.0", result[0].PURL)
	}
}

func TestLicenseScanner_parseLicensecheckOutput(t *testing.T) {
	log := &mockLogger{}
	scanner := NewLicenseScanner(log)

	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name:          "empty output",
			input:         "",
			expectedCount: 0,
		},
		{
			name:          "single license",
			input:         "LICENSE: MIT",
			expectedCount: 1,
		},
		{
			name:          "multiple licenses",
			input:         "LICENSE: MIT\nCOPYING: GPL-3.0\nNOTICE: Apache-2.0",
			expectedCount: 3,
		},
		{
			name:          "skip unknown",
			input:         "LICENSE: MIT\nREADME: UNKNOWN",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.parseLicensecheckOutput([]byte(tt.input))

			if len(result) != tt.expectedCount {
				t.Errorf("parseLicensecheckOutput() count = %d, expected %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestLicenseScanner_parseLicenseeOutput(t *testing.T) {
	log := &mockLogger{}
	scanner := NewLicenseScanner(log)

	tests := []struct {
		name               string
		input              string
		expectedCount      int
		expectedConfidence float64
	}{
		{
			name:          "empty output",
			input:         "",
			expectedCount: 0,
		},
		{
			name:               "with confidence",
			input:              "LICENSE: MIT (confidence: 95%)",
			expectedCount:      1,
			expectedConfidence: 0.95,
		},
		{
			name:               "without confidence",
			input:              "LICENSE: Apache-2.0",
			expectedCount:      1,
			expectedConfidence: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.parseLicenseeOutput([]byte(tt.input))

			if len(result) != tt.expectedCount {
				t.Errorf("parseLicenseeOutput() count = %d, expected %d", len(result), tt.expectedCount)
			}

			if tt.expectedCount > 0 && len(result) > 0 {
				if result[0].Confidence != tt.expectedConfidence {
					t.Errorf("Confidence = %f, expected %f", result[0].Confidence, tt.expectedConfidence)
				}
			}
		})
	}
}

func TestLicenseScanner_parseGoLicensesOutput(t *testing.T) {
	log := &mockLogger{}
	scanner := NewLicenseScanner(log)

	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name:          "empty output",
			input:         "",
			expectedCount: 0,
		},
		{
			name:          "single package",
			input:         "github.com/pkg/errors,MIT,/path/to/LICENSE",
			expectedCount: 1,
		},
		{
			name:          "multiple packages",
			input:         "github.com/pkg/errors,MIT,/path1\ngithub.com/sirupsen/logrus,MIT,/path2",
			expectedCount: 2,
		},
		{
			name:          "skip unknown",
			input:         "github.com/pkg/errors,MIT,/path\ngithub.com/unknown/pkg,Unknown,/path2",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.parseGoLicensesOutput([]byte(tt.input))

			if len(result) != tt.expectedCount {
				t.Errorf("parseGoLicensesOutput() count = %d, expected %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestLicenseScanner_classifyLicenseType(t *testing.T) {
	log := &mockLogger{}
	scanner := NewLicenseScanner(log)

	tests := []struct {
		licenseName  string
		expectedType string
	}{
		{"MIT", "permissive"},
		{"Apache-2.0", "permissive"},
		{"BSD-3-Clause", "permissive"},
		{"ISC", "permissive"},
		{"GPL-3.0", "copyleft"},
		{"LGPL-2.1", "copyleft"},
		{"AGPL-3.0", "copyleft"},
		{"MPL-2.0", "copyleft"},
		{"Proprietary", "proprietary"},
		{"Commercial License", "proprietary"},
		{"Custom License", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.licenseName, func(t *testing.T) {
			result := scanner.classifyLicenseType(tt.licenseName)
			if result != tt.expectedType {
				t.Errorf("classifyLicenseType(%s) = %s, expected %s", tt.licenseName, result, tt.expectedType)
			}
		})
	}
}

func TestLicenseScanner_mapToSPDXID(t *testing.T) {
	log := &mockLogger{}
	scanner := NewLicenseScanner(log)

	tests := []struct {
		licenseName    string
		expectedSPDXID string
	}{
		{"MIT", "MIT"},
		{"MIT License", "MIT"},
		{"Apache-2.0", "Apache-2.0"},
		{"Apache 2.0", "Apache-2.0"},
		{"GPL-3.0", "GPL-3.0-only"},
		{"GPLv3", "GPL-3.0-only"},
		{"BSD-3-Clause", "BSD-3-Clause"},
		{"Unknown License", "Unknown License"},
	}

	for _, tt := range tests {
		t.Run(tt.licenseName, func(t *testing.T) {
			result := scanner.mapToSPDXID(tt.licenseName)
			if result != tt.expectedSPDXID {
				t.Errorf("mapToSPDXID(%s) = %s, expected %s", tt.licenseName, result, tt.expectedSPDXID)
			}
		})
	}
}
