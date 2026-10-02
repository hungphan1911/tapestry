//go:build integration

package executor

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"code-runner/internal/runner"
)

func integrationExecutor() *Docker {
	limits := DefaultLimits()
	limits.RunTimeout = 500 * time.Millisecond
	limits.CompileTimeout = 10 * time.Second
	limits.OutputBytes = 1024
	return NewDocker(limits)
}

func TestDockerCPP(t *testing.T) {
	executor := integrationExecutor()
	tests := []struct {
		name   string
		source string
		stdin  string
		status runner.Status
		stdout string
		stderr string
	}{
		{"valid and stdin", `#include <iostream>
int main(){int x; std::cin>>x; std::cout << x * 3 << '\n';}`, "5\n", runner.StatusSuccess, "15\n", ""},
		{"compile error", `int main( {`, "", runner.StatusCompilationError, "", ""},
		{"runtime error", `#include <iostream>
int main(){std::cerr << "bad"; return 7;}`, "", runner.StatusRuntimeError, "", "bad"},
		{"timeout", `int main(){for(;;){}}`, "", runner.StatusTimeLimitExceeded, "", ""},
		{"output limit", `#include <iostream>
int main(){for(int i=0;i<2000;i++)std::cout << 'x';}`, "", runner.StatusOutputLimitExceeded, "", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := executor.Execute(context.Background(), runner.ExecutionRequest{Language: "cpp", Source: test.source, Stdin: test.stdin})
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != test.status {
				t.Fatalf("status = %s, output=%q, compile=%q", result.Status, result.Stderr, result.CompileOutput)
			}
			if test.stdout != "" && result.Stdout != test.stdout {
				t.Fatalf("stdout = %q", result.Stdout)
			}
			if test.stderr != "" && result.Stderr != test.stderr {
				t.Fatalf("stderr = %q", result.Stderr)
			}
		})
	}
}

func TestDockerPyPy(t *testing.T) {
	executor := integrationExecutor()
	tests := []struct {
		name, source, stdin string
		status              runner.Status
		stdout, stderr      string
	}{
		{"valid and stdin", `import sys
print(sys.stdin.read().strip()[::-1])`, "abc\n", runner.StatusSuccess, "cba\n", ""},
		{"syntax error", `def broken(:`, "", runner.StatusRuntimeError, "", ""},
		{"runtime stderr", `import sys
print("bad", file=sys.stderr)
raise RuntimeError("boom")`, "", runner.StatusRuntimeError, "", "bad"},
		{"timeout", `while True: pass`, "", runner.StatusTimeLimitExceeded, "", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := executor.Execute(context.Background(), runner.ExecutionRequest{Language: "python", Source: test.source, Stdin: test.stdin})
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != test.status {
				t.Fatalf("status = %s, stderr=%q", result.Status, result.Stderr)
			}
			if test.stdout != "" && result.Stdout != test.stdout {
				t.Fatalf("stdout = %q", result.Stdout)
			}
			if test.stderr != "" && !strings.Contains(result.Stderr, test.stderr) {
				t.Fatalf("stderr = %q", result.Stderr)
			}
		})
	}
}

func TestDockerIsolation(t *testing.T) {
	executor := integrationExecutor()
	marker, err := os.CreateTemp("", "code-runner-host-secret-")
	if err != nil {
		t.Fatal(err)
	}
	markerPath := marker.Name()
	if _, err := marker.WriteString("host-secret"); err != nil {
		t.Fatal(err)
	}
	if err := marker.Close(); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(markerPath)
	source := fmt.Sprintf(`import socket
try:
    socket.create_connection(("1.1.1.1", 53), 0.1)
    print("network-accessible")
except OSError:
    print("network-blocked")
try:
    open(%q).read()
    print("host-accessible")
except OSError:
    print("host-blocked")
`, markerPath)
	result, err := executor.Execute(context.Background(), runner.ExecutionRequest{Language: "python", Source: source})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != runner.StatusSuccess || result.Stdout != "network-blocked\nhost-blocked\n" {
		t.Fatalf("isolation failed: %#v", result)
	}
}
