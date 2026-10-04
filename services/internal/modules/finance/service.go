package finance

import (
	"context"
	"errors"
	"fmt"
	"github.com/hungphan1911/tapestry/services/internal/core"
	"strings"
	"time"

	"github.com/lib/pq"
)

const maxTextLength = 256

type Service interface {
	ListTypes(ctx context.Context) ([]Type, error)

	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, req CategoryRequest) (Category, error)
	UpdateCategory(ctx context.Context, id int, req CategoryRequest) (Category, error)
	DeleteCategory(ctx context.Context, id int) error

	ListCards(ctx context.Context) ([]Card, error)
	CreateCard(ctx context.Context, req CardRequest) (Card, error)
	UpdateCard(ctx context.Context, id int, req CardRequest) (Card, error)
	DeleteCard(ctx context.Context, id int) error

	ListTransactions(ctx context.Context, f TransactionFilter) ([]Transaction, error)
	GetTransaction(ctx context.Context, id int) (Transaction, error)
	CreateTransaction(ctx context.Context, req TransactionRequest) (Transaction, error)
	UpdateTransaction(ctx context.Context, id int, req TransactionRequest) (Transaction, error)
	DeleteTransaction(ctx context.Context, id int) error
}

type service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return &service{repo: repo, now: time.Now}
}

func validateText(field, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%w: %s is required", core.ErrInvalidInput, field)
	}
	if len(value) > maxTextLength {
		return "", fmt.Errorf("%w: %s exceeds %d characters", core.ErrInvalidInput, field, maxTextLength)
	}
	return value, nil
}

// mapDBError converts foreign-key violations (unknown type/category/card id)
// into ErrInvalidInput so callers see a 400 rather than a 500.
func mapDBError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23503" {
		return fmt.Errorf("%w: referenced type, category or card does not exist", core.ErrInvalidInput)
	}
	return err
}

func (s *service) ListTypes(ctx context.Context) ([]Type, error) {
	return s.repo.ListTypes(ctx)
}

// Categories

func (s *service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *service) CreateCategory(ctx context.Context, req CategoryRequest) (Category, error) {
	desc, err := validateText("description", req.Description)
	if err != nil {
		return Category{}, err
	}
	return s.repo.CreateCategory(ctx, desc)
}

func (s *service) UpdateCategory(ctx context.Context, id int, req CategoryRequest) (Category, error) {
	desc, err := validateText("description", req.Description)
	if err != nil {
		return Category{}, err
	}
	return s.repo.UpdateCategory(ctx, id, desc)
}

func (s *service) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

// Cards

func (s *service) ListCards(ctx context.Context) ([]Card, error) {
	return s.repo.ListCards(ctx)
}

func (s *service) CreateCard(ctx context.Context, req CardRequest) (Card, error) {
	name, err := validateText("card_name", req.CardName)
	if err != nil {
		return Card{}, err
	}
	return s.repo.CreateCard(ctx, name)
}

func (s *service) UpdateCard(ctx context.Context, id int, req CardRequest) (Card, error) {
	name, err := validateText("card_name", req.CardName)
	if err != nil {
		return Card{}, err
	}
	return s.repo.UpdateCard(ctx, id, name)
}

func (s *service) DeleteCard(ctx context.Context, id int) error {
	return s.repo.DeleteCard(ctx, id)
}

// Transactions

func (s *service) ListTransactions(ctx context.Context, f TransactionFilter) ([]Transaction, error) {
	if f.From != nil && f.To != nil && f.From.After(*f.To) {
		return nil, fmt.Errorf("%w: from must not be after to", core.ErrInvalidInput)
	}
	return s.repo.ListTransactions(ctx, f)
}

func (s *service) GetTransaction(ctx context.Context, id int) (Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *service) CreateTransaction(ctx context.Context, req TransactionRequest) (Transaction, error) {
	t, err := s.buildTransaction(req)
	if err != nil {
		return Transaction{}, err
	}
	created, err := s.repo.CreateTransaction(ctx, t)
	return created, mapDBError(err)
}

func (s *service) UpdateTransaction(ctx context.Context, id int, req TransactionRequest) (Transaction, error) {
	t, err := s.buildTransaction(req)
	if err != nil {
		return Transaction{}, err
	}
	updated, err := s.repo.UpdateTransaction(ctx, id, t)
	return updated, mapDBError(err)
}

func (s *service) DeleteTransaction(ctx context.Context, id int) error {
	return s.repo.DeleteTransaction(ctx, id)
}

// buildTransaction validates a request and evaluates its amount expression.
func (s *service) buildTransaction(req TransactionRequest) (Transaction, error) {
	desc := strings.TrimSpace(req.Description)
	if len(desc) > maxTextLength {
		return Transaction{}, fmt.Errorf("%w: description exceeds %d characters", core.ErrInvalidInput, maxTextLength)
	}
	expr := strings.TrimSpace(req.AmountExpression)
	if len(expr) > maxTextLength {
		return Transaction{}, fmt.Errorf("%w: amount_expression exceeds %d characters", core.ErrInvalidInput, maxTextLength)
	}
	amount, err := EvaluateExpression(expr)
	if err != nil {
		return Transaction{}, err
	}
	date := s.now()
	if req.Date != nil {
		date = *req.Date
	}
	return Transaction{
		Date:             date,
		Description:      desc,
		Amount:           amount,
		AmountExpression: expr,
		TypeID:           req.TypeID,
		CategoryID:       req.CategoryID,
		CardID:           req.CardID,
	}, nil
}

var _ Service = (*service)(nil)
