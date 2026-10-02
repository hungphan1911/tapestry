package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
	"sync"
	"time"

	"code-runner/internal/language"
	"code-runner/internal/runner"
)

const outputLimitError = "output limit exceeded"

type Limits struct {
	MemoryBytes    int64
	CPUs           string
	PIDs           int
	WorkspaceBytes int64
	OutputBytes    int
	CompileTimeout time.Duration
	RunTimeout     time.Duration
}

func DefaultLimits() Limits {
	return Limits{
		MemoryBytes: 256 * 1024 * 1024, CPUs: "1", PIDs: 64, WorkspaceBytes: 64 * 1024 * 1024,
		OutputBytes: 1024 * 1024, CompileTimeout: 10 * time.Second, RunTimeout: 2 * time.Second,
	}
}

type Docker struct{ limits Limits }

func NewDocker(limits Limits) *Docker { return &Docker{limits: limits} }

func (d *Docker) Execute(ctx context.Context, request runner.ExecutionRequest) (result runner.ExecutionResult, err error) {
	lang, err := language.Lookup(request.Language)
	if err != nil {
		return result, err
	}
	containerID, err := d.create(ctx, lang)
	if err != nil {
		return result, err
	}
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if cleanupErr := docker(cleanupCtx, nil, nil, "rm", "-f", containerID); cleanupErr != nil {
			slog.Error("container cleanup failed", "container", containerID, "error", cleanupErr)
		}
	}()
	if err := docker(ctx, nil, nil, "start", containerID); err != nil {
		return result, fmt.Errorf("start container: %w", err)
	}
	if err := d.writeSource(ctx, containerID, lang.SourceFilename, request.Source); err != nil {
		return result, err
	}

	if len(lang.CompileCommand) > 0 {
		slog.Info("run compile phase", "language", lang.ID)
		compileCtx, cancel := context.WithTimeout(ctx, d.limits.CompileTimeout)
		compile := d.exec(compileCtx, containerID, nil, lang.CompileCommand)
		cancel()
		result.CompileOutput = compile.stdout.String() + compile.stderr.String()
		if compile.outputLimit {
			result.Status = runner.StatusOutputLimitExceeded
			return result, nil
		}
		if compile.timedOut {
			result.Status = runner.StatusTimeLimitExceeded
			return result, nil
		}
		if compile.err != nil {
			result.Status = runner.StatusCompilationError
			result.ExitCode = runner.ExitCode(compile.err)
			return result, nil
		}
	}

	slog.Info("run execute phase", "language", lang.ID)
	start := time.Now()
	runCtx, cancel := context.WithTimeout(ctx, d.limits.RunTimeout)
	executed := d.exec(runCtx, containerID, strings.NewReader(request.Stdin), lang.RunCommand)
	cancel()
	result.ExecutionTimeMs = time.Since(start).Milliseconds()
	result.Stdout = executed.stdout.String()
	result.Stderr = executed.stderr.String()
	result.ExitCode = runner.ExitCode(executed.err)
	if executed.err == nil {
		code := 0
		result.ExitCode = &code
	}
	switch {
	case executed.outputLimit:
		result.Status = runner.StatusOutputLimitExceeded
	case executed.timedOut:
		result.Status = runner.StatusTimeLimitExceeded
	case executed.err != nil:
		result.Status = runner.StatusRuntimeError
	default:
		result.Status = runner.StatusSuccess
	}
	return result, nil
}

func (d *Docker) writeSource(ctx context.Context, containerID, filename, source string) error {
	// tee receives source only through stdin; filename is server-controlled by the language registry.
	if err := docker(ctx, strings.NewReader(source), nil, "exec", "-i", "--workdir", "/workspace", containerID, "tee", filename); err != nil {
		return fmt.Errorf("write source into container: %w", err)
	}
	return nil
}

func (d *Docker) create(ctx context.Context, lang language.Language) (string, error) {
	// The compiled C++ binary must execute from this otherwise isolated tmpfs.
	tmpfs := fmt.Sprintf("/workspace:rw,exec,nosuid,nodev,size=%d,uid=10001,gid=10001,mode=700", d.limits.WorkspaceBytes)
	args := []string{
		"create", "--network", "none", "--user", "runner", "--read-only", "--workdir", "/workspace",
		"--tmpfs", tmpfs, "--tmpfs", "/tmp:rw,noexec,nosuid,nodev,size=16m,uid=10001,gid=10001,mode=700",
		"--memory", fmt.Sprintf("%d", d.limits.MemoryBytes), "--memory-swap", fmt.Sprintf("%d", d.limits.MemoryBytes),
		"--cpus", d.limits.CPUs, "--pids-limit", fmt.Sprintf("%d", d.limits.PIDs),
		"--cap-drop", "ALL", "--security-opt", "no-new-privileges", lang.Image, "sleep", "infinity",
	}
	var output bytes.Buffer
	if err := docker(ctx, nil, &output, args...); err != nil {
		return "", fmt.Errorf("create container: %w: %s", err, strings.TrimSpace(output.String()))
	}
	return strings.TrimSpace(output.String()), nil
}

type commandResult struct {
	stdout      *limitedBuffer
	stderr      *limitedBuffer
	err         error
	timedOut    bool
	outputLimit bool
}

func (d *Docker) exec(parent context.Context, containerID string, stdin io.Reader, command []string) commandResult {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	args := append([]string{"exec", "-i", "--workdir", "/workspace", containerID}, command...)
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin = stdin
	stdout, stderr := &limitedBuffer{limit: d.limits.OutputBytes}, &limitedBuffer{limit: d.limits.OutputBytes}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	done := make(chan struct{})
	go func() {
		for !stdout.Limited() && !stderr.Limited() {
			select {
			case <-done:
				return
			default:
			}
			time.Sleep(time.Millisecond)
		}
		cancel()
	}()
	err := cmd.Run()
	close(done)
	return commandResult{stdout: stdout, stderr: stderr, err: err, timedOut: errors.Is(parent.Err(), context.DeadlineExceeded), outputLimit: stdout.Limited() || stderr.Limited()}
}

type limitedBuffer struct {
	mu      sync.Mutex
	buffer  bytes.Buffer
	limit   int
	limited bool
}

func (b *limitedBuffer) Write(data []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	remaining := b.limit - b.buffer.Len()
	if remaining <= 0 {
		b.limited = true
		return 0, errors.New(outputLimitError)
	}
	if len(data) > remaining {
		_, _ = b.buffer.Write(data[:remaining])
		b.limited = true
		return remaining, errors.New(outputLimitError)
	}
	return b.buffer.Write(data)
}

func (b *limitedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.String()
}

func (b *limitedBuffer) Limited() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.limited
}

func docker(ctx context.Context, stdin io.Reader, stdout io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin, cmd.Stdout = stdin, stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker %s: %w: %s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
