package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"code-runner/internal/runner"
)

type fakeExecutor struct{ result runner.ExecutionResult }

func (f fakeExecutor) Execute(context.Context, runner.ExecutionRequest) (runner.ExecutionResult, error) {
	return f.result, nil
}

func TestRunsRejectsInvalidRequests(t *testing.T) {
	handler := NewRuns(runner.NewService(fakeExecutor{}, 1))
	for _, body := range []string{
		`{"language":"go","source":"package main"}`,
		`{"language":"cpp","source":""}`,
		`not json`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/runs", bytes.NewBufferString(body))
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("status = %d for %s", response.Code, body)
		}
	}
}

func TestRunsReturnsExecutionResult(t *testing.T) {
	result := runner.ExecutionResult{Status: runner.StatusSuccess, Stdout: "15\n", ExecutionTimeMs: 12}
	handler := NewRuns(runner.NewService(fakeExecutor{result: result}, 1))
	req := httptest.NewRequest(http.MethodPost, "/api/runs", bytes.NewBufferString(`{"language":"cpp","source":"int main(){}"}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK || !bytes.Contains(response.Body.Bytes(), []byte(`"stdout":"15\n"`)) {
		t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
	}
}
