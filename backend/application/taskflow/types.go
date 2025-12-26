package taskflow

import (
	"syntaxia/application/router"
)

// Re-export router types for backward compatibility
type TaskPipeline = router.TaskPipeline
type TaskPipelineStep = router.TaskPipelineStep
type PipelinePolicy = router.PipelinePolicy

// Re-export router constants
const (
	PipelineStatusPending   = router.PipelineStatusPending
	PipelineStatusRunning   = router.PipelineStatusRunning
	PipelineStatusCompleted = router.PipelineStatusCompleted
	PipelineStatusFailed    = router.PipelineStatusFailed
)

const (
	StepStatusPending = router.StepStatusPending
	StepStatusRunning = router.StepStatusRunning
	StepStatusDone    = router.StepStatusCompleted // Map to completed
	StepStatusFailed  = router.StepStatusFailed
)

// Taskflow-specific step types (different from router step types)
const (
	StepTypeAnalysis = "analysis"
	StepTypeRepair   = "repair"
	StepTypeBuild    = "build"
	StepTypeTest     = "test"
)

// Constants for estimated time calculation
const (
	DefaultEstimatedTimeSeconds = 300  // 5 minutes default
	MaxEstimatedTimeSeconds     = 3600 // Maximum 1 hour
)
