package runner

import (
	"context"
	"strings"
	"testing"
)

type fakeExecutor struct {
	request ExecutionRequest
	result  ExecutionResult
	err     error
}

func (f *fakeExecutor) Execute(_ context.Context, request ExecutionRequest) (ExecutionResult, error) {
	f.request = request
	return f.result, f.err
}

func TestServiceValidatesRequests(t *testing.T) {
	service := NewService(&fakeExecutor{}, 1)
	tests := []ExecutionRequest{
		{Language: "ruby", Source: "puts 1"},
		{Language: "cpp"},
		{Language: "cpp", Source: strings.Repeat("x", MaxSourceBytes+1)},
		{Language: "cpp", Source: "int main(){}", Stdin: strings.Repeat("x", MaxStdinBytes+1)},
	}
	for _, request := range tests {
		if _, err := service.Execute(context.Background(), request); err == nil {
			t.Fatalf("request %#v was accepted", request)
		}
	}
}

func TestServicePassesStructuredRequest(t *testing.T) {
	fake := &fakeExecutor{result: ExecutionResult{Status: StatusSuccess}}
	service := NewService(fake, 1)
	request := ExecutionRequest{Language: "python", Source: "print(input())", Stdin: "hello\n"}
	if _, err := service.Execute(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	if fake.request != request {
		t.Fatalf("executor received %#v, want %#v", fake.request, request)
	}
}
