package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hungphan1911/tapestry/services/internal/core"
	"github.com/hungphan1911/tapestry/services/internal/modules/finance"
)

func main() {
	// Load configuration
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Loading .env failed: %v", err)
	}
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	// Initialize layer architecture
	financeRepo := finance.NewRepository(db) 
	financeSvc := finance.NewService(financeRepo)   
	financeHandler := finance.NewHandler(financeSvc)

	// 4. Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		financeHandler.RegisterRoutes(r)
	})

	// 5. Start HTTP Server
	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server listening on %s\n", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
