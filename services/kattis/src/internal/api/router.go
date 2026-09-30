// src/internal/api/router.go
package api

import (
	"net/http"
	"kattis/src/internal/api/handler"
)

type Server struct {

}

func NewRouter() http.Handler {
	server := &Server{}
	mux := http.NewServeMux()

	server.registerHealthRoutes(mux)
	server.registerProblemsetRoutes(mux)
	server.registerStatsRoutes(mux)
	server.registerDashboardRoutes(mux)
	server.registerSubmissionRoutes(mux)
	
	return mux
}

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", handler.HealthCheck)
}

func (s *Server) registerProblemsetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /problemsets", handler.ProblemSets)
	mux.HandleFunc("GET /problemsets/completed", handler.CompletedProblemSets)
}

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /stats/summary", handler.Summary)
	mux.HandleFunc("GET /stats/difficulty-distribution", handler.DifficultyDistribution)
	mux.HandleFunc("GET /stats/activity", handler.Activity)
}

func (s *Server) registerDashboardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /dashboard", handler.Dashboard)
}

func (s *Server) registerSubmissionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /submissions", handler.Submissions)
}