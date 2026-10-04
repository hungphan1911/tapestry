package finance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hungphan1911/tapestry/services/internal/core"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const maxBodyBytes = 1 << 20

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/finance", func(r chi.Router) {
		r.Get("/types", h.ListTypes)

		r.Route("/categories", func(r chi.Router) {
			r.Get("/", h.ListCategories)
			r.Post("/", h.CreateCategory)
			r.Put("/{id}", h.UpdateCategory)
			r.Delete("/{id}", h.DeleteCategory)
		})

		r.Route("/cards", func(r chi.Router) {
			r.Get("/", h.ListCards)
			r.Post("/", h.CreateCard)
			r.Put("/{id}", h.UpdateCard)
			r.Delete("/{id}", h.DeleteCard)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", h.ListTransactions)
			r.Post("/", h.CreateTransaction)
			r.Get("/{id}", h.GetTransaction)
			r.Put("/{id}", h.UpdateTransaction)
			r.Delete("/{id}", h.DeleteTransaction)
		})
	})
}

// Types

func (h *Handler) ListTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.service.ListTypes(r.Context())
	respond(w, http.StatusOK, types, err)
}

// Categories

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	respond(w, http.StatusOK, categories, err)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	category, err := h.service.CreateCategory(r.Context(), req)
	respond(w, http.StatusCreated, category, err)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req CategoryRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	category, err := h.service.UpdateCategory(r.Context(), id, req)
	respond(w, http.StatusOK, category, err)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	respondNoContent(w, h.service.DeleteCategory(r.Context(), id))
}

// Cards

func (h *Handler) ListCards(w http.ResponseWriter, r *http.Request) {
	cards, err := h.service.ListCards(r.Context())
	respond(w, http.StatusOK, cards, err)
}

func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req CardRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	card, err := h.service.CreateCard(r.Context(), req)
	respond(w, http.StatusCreated, card, err)
}

func (h *Handler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req CardRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	card, err := h.service.UpdateCard(r.Context(), id, req)
	respond(w, http.StatusOK, card, err)
}

func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	respondNoContent(w, h.service.DeleteCard(r.Context(), id))
}

// Transactions

func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	filter, err := parseTransactionFilter(r)
	if err != nil {
		writeError(w, err)
		return
	}
	transactions, err := h.service.ListTransactions(r.Context(), filter)
	respond(w, http.StatusOK, transactions, err)
}

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	transaction, err := h.service.GetTransaction(r.Context(), id)
	respond(w, http.StatusOK, transaction, err)
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	transaction, err := h.service.CreateTransaction(r.Context(), req)
	respond(w, http.StatusCreated, transaction, err)
}

func (h *Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req TransactionRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	transaction, err := h.service.UpdateTransaction(r.Context(), id, req)
	respond(w, http.StatusOK, transaction, err)
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	respondNoContent(w, h.service.DeleteTransaction(r.Context(), id))
}

// Helpers

func parseTransactionFilter(r *http.Request) (TransactionFilter, error) {
	q := r.URL.Query()
	var f TransactionFilter

	for _, p := range []struct {
		name string
		dst  **time.Time
	}{{"from", &f.From}, {"to", &f.To}} {
		if v := q.Get(p.name); v != "" {
			t, err := parseTime(v)
			if err != nil {
				return f, fmt.Errorf("%w: %s must be RFC3339 or YYYY-MM-DD", core.ErrInvalidInput, p.name)
			}
			*p.dst = &t
		}
	}

	for _, p := range []struct {
		name string
		dst  **int
	}{{"category_id", &f.CategoryID}, {"type_id", &f.TypeID}, {"card_id", &f.CardID}} {
		if v := q.Get(p.name); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return f, fmt.Errorf("%w: %s must be an integer", core.ErrInvalidInput, p.name)
			}
			*p.dst = &n
		}
	}
	return f, nil
}

func parseTime(v string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", v)
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		writeError(w, fmt.Errorf("%w: id must be a positive integer", core.ErrInvalidInput))
		return 0, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, fmt.Errorf("%w: malformed request body: %v", core.ErrInvalidInput, err))
		return false
	}
	return true
}

func respond(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, status, body)
}

func respondNoContent(w http.ResponseWriter, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"
	switch {
	case errors.Is(err, core.ErrNotFound):
		status, msg = http.StatusNotFound, "resource not found"
	case errors.Is(err, core.ErrInvalidInput):
		status, msg = http.StatusBadRequest, err.Error()
	default:
		log.Printf("internal error: %v", err)
	}
	writeJSON(w, status, map[string]string{"error": msg})
}
