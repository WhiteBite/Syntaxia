package initmanager

import (
	"context"
	"syntaxia/application/ai"
	"syntaxia/application/diff"
	"syntaxia/application/export"
	"syntaxia/application/router"
	"syntaxia/application/sbom"
	"syntaxia/application/symbol"
	"syntaxia/application/verification"
	"syntaxia/domain"
	contextservice "syntaxia/internal/context"
)

// LazyServices holds lazy-initialized services
type LazyServices struct {
	// AI Services
	AIService *LazyService[*ai.Service]

	// Context Services (unified service replaces old separate services)
	ContextService  *LazyService[*contextservice.Service]
	ContextBuilder  *LazyService[domain.ContextBuilder]
	ContextAnalyzer *LazyService[domain.ContextAnalyzer]

	// Analysis Services
	TestService           *LazyService[domain.ITestService]
	StaticAnalyzerService *LazyService[domain.IStaticAnalyzerService]
	SBOMService           *LazyService[*sbom.Service]

	// Build and Verification Services
	BuildService                *LazyService[domain.IBuildService]
	VerificationPipelineService *LazyService[*verification.Service]

	// Export and Report Services
	ExportService *LazyService[*export.Service]
	ReportService *LazyService[*export.ReportService]

	// Task Protocol Services
	TaskProtocolService *LazyService[domain.TaskProtocolService]
	ErrorAnalyzer       *LazyService[domain.ErrorAnalyzer]
	CorrectionEngine    *LazyService[domain.CorrectionEngine]

	// Symbol Graph Service
	SymbolGraph *LazyService[*symbol.Service]

	// Repair Service
	RepairService *LazyService[domain.RepairService]

	// Guardrail Service
	GuardrailService *LazyService[domain.GuardrailService]

	// Taskflow Service
	TaskflowService *LazyService[domain.TaskflowService]

	// UX Metrics Service
	UXMetricsService *LazyService[domain.UXMetricsService]

	// Apply and Diff Services
	ApplyService *LazyService[*diff.ApplyService]
	DiffService  *LazyService[*diff.Service]

	// Router LLM Service
	RouterLLMService *LazyService[*router.LLMService]
}

// NewLazyServices creates a new LazyServices instance
func NewLazyServices() *LazyServices {
	return &LazyServices{
		AIService:                   NewLazyService(func(ctx context.Context) (*ai.Service, error) { return nil, nil }),
		ContextService:              NewLazyService(func(ctx context.Context) (*contextservice.Service, error) { return nil, nil }),
		ContextBuilder:              NewLazyService(func(ctx context.Context) (domain.ContextBuilder, error) { return nil, nil }),
		ContextAnalyzer:             NewLazyService(func(ctx context.Context) (domain.ContextAnalyzer, error) { return nil, nil }),
		TestService:                 NewLazyService(func(ctx context.Context) (domain.ITestService, error) { return nil, nil }),
		StaticAnalyzerService:       NewLazyService(func(ctx context.Context) (domain.IStaticAnalyzerService, error) { return nil, nil }),
		SBOMService:                 NewLazyService(func(ctx context.Context) (*sbom.Service, error) { return nil, nil }),
		BuildService:                NewLazyService(func(ctx context.Context) (domain.IBuildService, error) { return nil, nil }),
		VerificationPipelineService: NewLazyService(func(ctx context.Context) (*verification.Service, error) { return nil, nil }),
		ExportService:               NewLazyService(func(ctx context.Context) (*export.Service, error) { return nil, nil }),
		ReportService:               NewLazyService(func(ctx context.Context) (*export.ReportService, error) { return nil, nil }),
		TaskProtocolService:         NewLazyService(func(ctx context.Context) (domain.TaskProtocolService, error) { return nil, nil }),
		ErrorAnalyzer:               NewLazyService(func(ctx context.Context) (domain.ErrorAnalyzer, error) { return nil, nil }),
		CorrectionEngine:            NewLazyService(func(ctx context.Context) (domain.CorrectionEngine, error) { return nil, nil }),
		SymbolGraph:                 NewLazyService(func(ctx context.Context) (*symbol.Service, error) { return nil, nil }),
		RepairService:               NewLazyService(func(ctx context.Context) (domain.RepairService, error) { return nil, nil }),
		GuardrailService:            NewLazyService(func(ctx context.Context) (domain.GuardrailService, error) { return nil, nil }),
		TaskflowService:             NewLazyService(func(ctx context.Context) (domain.TaskflowService, error) { return nil, nil }),
		UXMetricsService:            NewLazyService(func(ctx context.Context) (domain.UXMetricsService, error) { return nil, nil }),
		ApplyService:                NewLazyService(func(ctx context.Context) (*diff.ApplyService, error) { return nil, nil }),
		DiffService:                 NewLazyService(func(ctx context.Context) (*diff.Service, error) { return nil, nil }),
		RouterLLMService:            NewLazyService(func(ctx context.Context) (*router.LLMService, error) { return nil, nil }),
	}
}
