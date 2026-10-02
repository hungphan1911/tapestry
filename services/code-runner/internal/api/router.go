package api

import (
	"net/http"

	"code-runner/internal/api/handler"
	"code-runner/internal/runner"
)

func NewRouter(service *runner.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.Handle("POST /api/runs", handler.NewRuns(service))
	return mux
}
