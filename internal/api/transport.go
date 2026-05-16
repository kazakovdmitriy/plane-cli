package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	DefaultMaxRetries          = 4
	DefaultTimeout             = 30 * time.Second
	DefaultCircuitBreakerLimit = 5
)

type RetryConfig struct {
	MaxRetries          int
	TotalTimeout        time.Duration
	CircuitBreakerLimit int
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:          DefaultMaxRetries,
		TotalTimeout:        DefaultTimeout,
		CircuitBreakerLimit: DefaultCircuitBreakerLimit,
	}
}

type CircuitBreaker struct {
	mu               sync.Mutex
	consecutiveFails int
	limit            int
	open             bool
}

func NewCircuitBreaker(limit int) *CircuitBreaker {
	return &CircuitBreaker{limit: limit}
}

func (cb *CircuitBreaker) IsOpen() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.open
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.consecutiveFails++
	if cb.consecutiveFails >= cb.limit {
		cb.open = true
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.consecutiveFails = 0
	cb.open = false
}

var ErrCircuitBreakerOpen = fmt.Errorf("circuit breaker open")

func newRetryableHTTP(cfg *RetryConfig) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = cfg.MaxRetries
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 8 * time.Second
	client.HTTPClient.Timeout = cfg.TotalTimeout
	client.Logger = nil
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			return true, nil
		}
		status := resp.StatusCode
		if status == 429 || status >= 500 {
			return true, nil
		}
		return false, nil
	}
	return client
}
