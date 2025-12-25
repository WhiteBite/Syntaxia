package commands

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"time"
)

// VerifyCommand represents the verification command
type VerifyCommand struct {
	container *CLIContainer
}

// NewVerifyCommand creates a new verification command
func NewVerifyCommand(container *CLIContainer) *VerifyCommand {
	return &VerifyCommand{
		container: container,
	}
}

// Execute executes the verification command
func (c *VerifyCommand) Execute(ctx context.Context, args []string) error {
	// Create flags for the command
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	var (
		projectPath = fs.String("project", ".", "Project path to verify")
		languages   = fs.String("languages", "", "Comma-separated list of languages to verify (default: auto-detect)")
		output      = fs.String("output", "", "Output file for verification report (JSON)")
		verbose     = fs.Bool("verbose", false, "Verbose output")
		help        = fs.Bool("help", false, "Show help")
	)

	// Parse arguments
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Show help if requested
	if *help {
		c.printHelp()
		return nil
	}

	// Check if project exists
	if _, err := os.Stat(*projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", *projectPath)
	}

	// Get absolute path
	absPath, err := filepath.Abs(*projectPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if *verbose {
		fmt.Printf("Verifying project: %s\n", absPath)
	}

	// Parse languages
	var languageList []string
	if *languages != "" {
		languageList = []string{*languages}
	} else {
		// Auto-detect languages
		languageList, err = c.container.VerificationService.DetectLanguages(ctx, absPath)
		if err != nil {
			return fmt.Errorf("failed to detect languages: %w", err)
		}
		if *verbose {
			fmt.Printf("Detected languages: %v\n", languageList)
		}
	}

	// Create verification config
	config := &domain.VerificationConfig{
		ProjectPath: absPath,
		Languages:   languageList,
		Timeout:     300, // 5 minutes
		Verbose:     *verbose,
	}

	// Run verification pipeline
	result, err := c.container.VerificationService.RunVerificationPipeline(ctx, config)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	// Create verification result
	verifyResult := &VerifyResult{
		ProjectPath: absPath,
		Languages:   languageList,
		Success:     result.Success,
		Steps:       result.Steps,
		Timestamp:   time.Now(),
	}

	// Output result
	if *output != "" {
		// Save to file
		data, err := json.MarshalIndent(verifyResult, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal verification result: %w", err)
		}

		if err := os.WriteFile(*output, data, 0o644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Verification report saved to: %s\n", *output)
	} else {
		// Print to stdout
		if result.Success {
			fmt.Println("✅ Verification completed successfully!")
		} else {
			fmt.Println("❌ Verification failed!")
		}

		// Print step results
		for _, step := range result.Steps {
			status := "✅"
			if !step.Success {
				status = "❌"
			}
			fmt.Printf("%s %s\n", status, step.Name)
		}

		// Print detailed results in verbose mode
		if *verbose {
			data, err := json.MarshalIndent(verifyResult, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal verification result: %w", err)
			}
			fmt.Println("\nDetailed Results:")
			fmt.Println(string(data))
		}
	}

	return nil
}

// printHelp prints help for the command
func (c *VerifyCommand) printHelp() {
	fmt.Printf(`ark verify - Verify project quality and health

Usage: ark verify [options]

Options:
  -project string
        Project path to verify (default ".")
  -languages string
        Comma-separated list of languages to verify (default: auto-detect)
  -output string
        Output file for verification report (JSON)
  -verbose
        Verbose output
  -help
        Show this help message

Examples:
  ark verify --project ./my-project
  ark verify --project ./my-project --languages go,typescript
  ark verify --project ./my-project --output report.json --verbose
`)
}

// VerifyResult represents the result of verification
type VerifyResult struct {
	ProjectPath string                     `json:"project_path"`
	Languages   []string                   `json:"languages"`
	Success     bool                       `json:"success"`
	Steps       []*domain.VerificationStep `json:"steps"`
	Timestamp   time.Time                  `json:"timestamp"`
}
