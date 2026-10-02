package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code-runner/internal/api"
	"code-runner/internal/executor"
	"code-runner/internal/runner"
)

func main() {
	limits := executor.DefaultLimits()
	service := runner.NewService(executor.NewDocker(limits), 2)
	server := &http.Server{Addr: ":8080", Handler: api.NewRouter(service), ReadHeaderTimeout: 5 * time.Second}

	go func() {
		slog.Info("server started", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
}
