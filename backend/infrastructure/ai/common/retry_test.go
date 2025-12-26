package common

import (
	"context"
	"errors"
	"syntaxia/domain"
	"testing"
	"time"
)

func TestWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   100 * time.Millisecond,
	}

	callCount := 0
	result, err := WithRetry(ctx, config, func() (string, error) {
		callCount++
		return "success", nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("expected 'success', got %s", result)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestWithRetry_SuccessAfterRetries(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   100 * time.Millisecond,
	}

	callCount := 0
	result, err := WithRetry(ctx, config, func() (string, error) {
		callCount++
		if callCount < 3 {
			return "", errors.New("transient error")
		}
		return "success", nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("expected 'success', got %s", result)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestWithRetry_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries: 2,
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   100 * time.Millisecond,
	}

	callCount := 0
	expectedErr := errors.New("persistent error")
	_, err := WithRetry(ctx, config, func() (string, error) {
		callCount++
		return "", expectedErr
	})

	if err == nil {
		t.Error("expected error, got nil")
	}
	// MaxRetries=2 means initial attempt + 2 retries = 3 total calls
	if callCount != 3 {
		t.Errorf("expected 3 calls (initial + 2 retries), got %d", callCount)
	}
}

func TestWithRetry_NonRecoverableError(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   100 * time.Millisecond,
	}

	callCount := 0
	nonRecoverableErr := &domain.DomainError{
		Code:        domain.ErrCodeInvalidAPIKey,
		Message:     "Invalid API key",
		Recoverable: false,
	}

	_, err := WithRetry(ctx, config, func() (string, error) {
		callCount++
		return "", nonRecoverableErr
	})

	if err == nil {
		t.Error("expected error, got nil")
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retries for non-recoverable), got %d", callCount)
	}
}

func TestWithRetry_RecoverableError(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries: 2,
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   100 * time.Millisecond,
	}

	callCount := 0
	recoverableErr := &domain.DomainError{
		Code:        domain.ErrCodeRateLimitExceeded,
		Message:     "Rate limit exceeded",
		Recoverable: true,
	}

	result, err := WithRetry(ctx, config, func() (string, error) {
		callCount++
		if callCount < 2 {
			return "", recoverableErr
		}
		return "success", nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("expected 'success', got %s", result)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestWithRetry_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := RetryConfig{
		MaxRetries: 5,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   1 * time.Second,
	}

	callCount := 0
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := WithRetry(ctx, config, func() (string, error) {
		callCount++
		return "", errors.New("transient error")
	})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestWithRetry_ContextCanceledError_NoRetry(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   100 * time.Millisecond,
	}

	callCount := 0
	_, err := WithRetry(ctx, config, func() (string, error) {
		callCount++
		return "", context.Canceled
	})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retries for context.Canceled), got %d", callCount)
	}
}

func TestCalculateBackoff(t *testing.T) {
	config := RetryConfig{
		MaxRetries: 5,
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},  // 1 * 2^0 = 1s
		{1, 2 * time.Second},  // 1 * 2^1 = 2s
		{2, 4 * time.Second},  // 1 * 2^2 = 4s
		{3, 8 * time.Second},  // 1 * 2^3 = 8s
		{4, 16 * time.Second}, // 1 * 2^4 = 16s
		{5, 30 * time.Second}, // 1 * 2^5 = 32s, capped at 30s
		{6, 30 * time.Second}, // capped at max
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := calculateBackoff(tt.attempt, config)
			if result != tt.expected {
				t.Errorf("attempt %d: expected %v, got %v", tt.attempt, tt.expected, result)
			}
		})
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != DefaultMaxRetries {
		t.Errorf("expected MaxRetries=%d, got %d", DefaultMaxRetries, config.MaxRetries)
	}
	if config.BaseDelay != DefaultBaseDelay {
		t.Errorf("expected BaseDelay=%v, got %v", DefaultBaseDelay, config.BaseDelay)
	}
	if config.MaxDelay != DefaultMaxDelay {
		t.Errorf("expected MaxDelay=%v, got %v", DefaultMaxDelay, config.MaxDelay)
	}
}

func TestWithRetryNoResult(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries: 2,
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   100 * time.Millisecond,
	}

	callCount := 0
	err := WithRetryNoResult(ctx, config, func() error {
		callCount++
		if callCount < 2 {
			return errors.New("transient error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}
