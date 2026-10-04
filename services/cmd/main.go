package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/hungphan1911/tapestry/services/internal/core"
	"github.com/hungphan1911/tapestry/services/internal/modules"
)

func main() {
	// Load configuration
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Loading .env failed: %v", err)
	}
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Build every module; each opens its own database
	mods, err := modules.Build(cfg.Postgres)
	if err != nil {
		log.Fatalf("Building modules failed: %v", err)
	}
	defer modules.CloseAll(mods)

	// Setup router, mounting each module under /api/v1/<name>
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		for _, m := range mods {
			r.Route("/"+m.Name(), m.RegisterRoutes)
		}
	})

	// Start HTTP Server
	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server listening on %s\n", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
	}
}
