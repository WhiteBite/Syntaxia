// Package common provides shared utilities for AI providers
package common

import (
	"context"
	"errors"
	"math"
	"syntaxia/domain"
	"time"
)

const (
	// DefaultMaxRetries is the default number of retry attempts
	DefaultMaxRetries = 3
	// DefaultBaseDelay is the initial delay between retries
	DefaultBaseDelay = 1 * time.Second
	// DefaultMaxDelay is the maximum delay between retries
	DefaultMaxDelay = 30 * time.Second
)

// RetryConfig configures retry behavior for AI provider calls
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: DefaultMaxRetries,
		BaseDelay:  DefaultBaseDelay,
		MaxDelay:   DefaultMaxDelay,
	}
}

// WithRetry executes a function with exponential backoff retry logic
// It will retry on recoverable errors and respect context cancellation
func WithRetry[T any](ctx context.Context, config RetryConfig, fn func() (T, error)) (T, error) {
	var lastErr error
	var zero T

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		// Don't retry if error is not recoverable
		if !isRecoverable(err) {
			return zero, err
		}

		lastErr = err

		if attempt < config.MaxRetries {
			delay := calculateBackoff(attempt, config)
			select {
			case <-ctx.Done():
				return zero, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return zero, lastErr
}

// WithRetryNoResult executes a function with no return value with exponential backoff
func WithRetryNoResult(ctx context.Context, config RetryConfig, fn func() error) error {
	_, err := WithRetry(ctx, config, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	return err
}

// isRecoverable checks if an error is recoverable and should be retried
func isRecoverable(err error) bool {
	if err == nil {
		return false
	}

	// Check for context cancellation - not recoverable
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Check for DomainError with Recoverable flag
	var domainErr *domain.DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Recoverable
	}

	// Default: assume network/transient errors are recoverable
	return true
}

// calculateBackoff calculates the delay for a given retry attempt using exponential backoff
func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	delay := config.BaseDelay * time.Duration(math.Pow(2, float64(attempt)))
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}
	return delay
}
