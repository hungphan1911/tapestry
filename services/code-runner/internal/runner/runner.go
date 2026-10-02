package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code-runner/internal/language"
)

const (
	MaxSourceBytes = 512 * 1024
	MaxStdinBytes  = 1024 * 1024
)

type Status string

const (
	StatusSuccess             Status = "success"
	StatusCompilationError    Status = "compilation_error"
	StatusRuntimeError        Status = "runtime_error"
	StatusTimeLimitExceeded   Status = "time_limit_exceeded"
	StatusOutputLimitExceeded Status = "output_limit_exceeded"
	StatusInternalError       Status = "internal_error"
)

type ExecutionRequest struct {
	Language string `json:"language"`
	Source   string `json:"source"`
	Stdin    string `json:"stdin"`
}

type ExecutionResult struct {
	Status          Status `json:"status"`
	Stdout          string `json:"stdout"`
	Stderr          string `json:"stderr"`
	CompileOutput   string `json:"compileOutput"`
	ExitCode        *int   `json:"exitCode"`
	ExecutionTimeMs int64  `json:"executionTimeMs"`
}

type Executor interface {
	Execute(context.Context, ExecutionRequest) (ExecutionResult, error)
}

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

type Service struct {
	executor Executor
	slots    chan struct{}
}

func NewService(executor Executor, concurrency int) *Service {
	return &Service{executor: executor, slots: make(chan struct{}, concurrency)}
}

func (s *Service) Execute(ctx context.Context, request ExecutionRequest) (ExecutionResult, error) {
	if _, err := language.Lookup(request.Language); err != nil {
		return ExecutionResult{}, &ValidationError{Message: err.Error()}
	}
	if request.Source == "" {
		return ExecutionResult{}, &ValidationError{Message: "source must not be empty"}
	}
	if len(request.Source) > MaxSourceBytes {
		return ExecutionResult{}, &ValidationError{Message: fmt.Sprintf("source exceeds %d byte limit", MaxSourceBytes)}
	}
	if len(request.Stdin) > MaxStdinBytes {
		return ExecutionResult{}, &ValidationError{Message: fmt.Sprintf("stdin exceeds %d byte limit", MaxStdinBytes)}
	}
	select {
	case s.slots <- struct{}{}:
		defer func() { <-s.slots }()
	case <-ctx.Done():
		return ExecutionResult{}, ctx.Err()
	}

	result, err := s.executor.Execute(ctx, request)
	if err != nil {
		return ExecutionResult{Status: StatusInternalError}, err
	}
	return result, nil
}

func ExitCode(err error) *int {
	var exitErr interface{ ExitCode() int }
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()
		return &code
	}
	return nil
}

func DurationMillis(start time.Time) int64 { return time.Since(start).Milliseconds() }
