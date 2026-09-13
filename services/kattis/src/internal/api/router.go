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

	return mux
}

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", handler.HealthCheck)
}

func (s *Server) registerProblemsetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /problemsets", nil)
	mux.HandleFunc("GET /problemsets/completed", nil)
}

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /stats/summary", nil)
	mux.HandleFunc("GET /stats/difficulty-distribution", nil)
	mux.HandleFunc("GET /stats/activity", nil)
}

func (s *Server) registerDashboardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /dashboard", nil)
}