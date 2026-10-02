package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"code-runner/internal/runner"
)

type Runs struct{ service *runner.Service }

func NewRuns(service *runner.Service) *Runs { return &Runs{service: service} }

func (h *Runs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, runner.MaxSourceBytes+runner.MaxStdinBytes+1024))
	decoder.DisallowUnknownFields()
	var request runner.ExecutionRequest
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	start := time.Now()
	slog.Info("run started", "language", request.Language)
	result, err := h.service.Execute(r.Context(), request)
	if err != nil {
		var validation *runner.ValidationError
		if errors.As(err, &validation) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": validation.Message})
			return
		}
		slog.Error("run failed", "language", request.Language, "duration_ms", time.Since(start).Milliseconds(), "error", err)
		writeJSON(w, http.StatusInternalServerError, result)
		return
	}
	slog.Info("run finished", "language", request.Language, "status", result.Status, "duration_ms", time.Since(start).Milliseconds())
	writeJSON(w, http.StatusOK, result)
}

func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
