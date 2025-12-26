package taskflow

import (
	"syntaxia/domain"
	"testing"
	"time"
)

func TestCalculateEstimatedTimeRemaining_CompletedTask(t *testing.T) {
	service := &Service{}

	status := &domain.TaskStatus{
		State:    domain.TaskStateDone,
		Progress: 1.0,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result != 0 {
		t.Errorf("expected 0 for completed task, got %d", result)
	}
}

func TestCalculateEstimatedTimeRemaining_FailedTask(t *testing.T) {
	service := &Service{}

	status := &domain.TaskStatus{
		State:    domain.TaskStateFailed,
		Progress: 0.5,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result != 0 {
		t.Errorf("expected 0 for failed task, got %d", result)
	}
}

func TestCalculateEstimatedTimeRemaining_NoStartTime(t *testing.T) {
	service := &Service{}

	status := &domain.TaskStatus{
		State:     domain.TaskStateTodo,
		Progress:  0.5,
		StartedAt: nil,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result != DefaultEstimatedTimeSeconds {
		t.Errorf("expected default %d for task without start time, got %d", DefaultEstimatedTimeSeconds, result)
	}
}

func TestCalculateEstimatedTimeRemaining_ZeroProgress(t *testing.T) {
	service := &Service{}
	now := time.Now()

	status := &domain.TaskStatus{
		State:     domain.TaskStateTodo,
		Progress:  0.0,
		StartedAt: &now,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result != DefaultEstimatedTimeSeconds {
		t.Errorf("expected default %d for zero progress, got %d", DefaultEstimatedTimeSeconds, result)
	}
}

func TestCalculateEstimatedTimeRemaining_InProgress(t *testing.T) {
	service := &Service{}

	startTime := time.Now().Add(-60 * time.Second)

	status := &domain.TaskStatus{
		State:     domain.TaskStateTodo,
		Progress:  0.5,
		StartedAt: &startTime,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result < 55 || result > 65 {
		t.Errorf("expected ~60 seconds remaining, got %d", result)
	}
}

func TestCalculateEstimatedTimeRemaining_AlmostComplete(t *testing.T) {
	service := &Service{}

	startTime := time.Now().Add(-90 * time.Second)

	status := &domain.TaskStatus{
		State:     domain.TaskStateTodo,
		Progress:  0.9,
		StartedAt: &startTime,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result < 5 || result > 15 {
		t.Errorf("expected ~10 seconds remaining, got %d", result)
	}
}

func TestCalculateEstimatedTimeRemaining_FullProgress(t *testing.T) {
	service := &Service{}
	startTime := time.Now().Add(-60 * time.Second)

	status := &domain.TaskStatus{
		State:     domain.TaskStateTodo,
		Progress:  1.0,
		StartedAt: &startTime,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result != 0 {
		t.Errorf("expected 0 for 100%% progress, got %d", result)
	}
}

func TestCalculateEstimatedTimeRemaining_MaxCapped(t *testing.T) {
	service := &Service{}

	startTime := time.Now().Add(-10 * time.Second)

	status := &domain.TaskStatus{
		State:     domain.TaskStateTodo,
		Progress:  0.001,
		StartedAt: &startTime,
	}

	result := service.calculateEstimatedTimeRemaining(status)

	if result > MaxEstimatedTimeSeconds {
		t.Errorf("expected result capped at %d, got %d", MaxEstimatedTimeSeconds, result)
	}
}
