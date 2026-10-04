package finance

import (
	"context"
	"errors"
	"fmt"
	"github.com/hungphan1911/tapestry/services/internal/core"
	"sort"
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

	GetBudget(ctx context.Context, month string) (Budget, error)
	SetBudget(ctx context.Context, month string, req BudgetRequest) (Budget, error)

	ListRecurring(ctx context.Context, month string) ([]Recurring, error)
	CreateRecurring(ctx context.Context, req RecurringRequest) (Recurring, error)
	UpdateRecurring(ctx context.Context, id int, req RecurringRequest) (Recurring, error)
	DeleteRecurring(ctx context.Context, id int) error
	SetRecurringAmount(ctx context.Context, id int, month string, req RecurringAmountRequest) (Recurring, error)
	DeleteRecurringAmount(ctx context.Context, id int, month string) (Recurring, error)
	LogRecurring(ctx context.Context, id int, req LogRecurringRequest) (Transaction, error)

	Summary(ctx context.Context, month string) (Summary, error)
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

// Months are "YYYY-MM" strings in the API, treated as UTC calendar months.

const monthLayout = "2006-01"

func parseMonth(field, value string) (time.Time, error) {
	t, err := time.Parse(monthLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s must be YYYY-MM", core.ErrInvalidInput, field)
	}
	return t, nil
}

func (s *service) currentMonth() string {
	return s.now().UTC().Format(monthLayout)
}

func monthRange(month time.Time) (from, to time.Time) {
	return month, month.AddDate(0, 1, 0)
}

// Budgets

func (s *service) GetBudget(ctx context.Context, month string) (Budget, error) {
	if _, err := parseMonth("month", month); err != nil {
		return Budget{}, err
	}
	return s.repo.GetBudget(ctx, month)
}

func (s *service) SetBudget(ctx context.Context, month string, req BudgetRequest) (Budget, error) {
	if _, err := parseMonth("month", month); err != nil {
		return Budget{}, err
	}
	if req.Total != nil && *req.Total < 0 {
		return Budget{}, fmt.Errorf("%w: total must not be negative", core.ErrInvalidInput)
	}
	seen := map[int]bool{}
	for _, c := range req.Categories {
		if c.Amount < 0 {
			return Budget{}, fmt.Errorf("%w: category amounts must not be negative", core.ErrInvalidInput)
		}
		if seen[c.CategoryID] {
			return Budget{}, fmt.Errorf("%w: duplicate category in budget", core.ErrInvalidInput)
		}
		seen[c.CategoryID] = true
	}
	if err := s.repo.ReplaceBudget(ctx, month, req); err != nil {
		return Budget{}, mapDBError(err)
	}
	return s.repo.GetBudget(ctx, month)
}

// Recurring

// amountForMonth returns the latest amount whose effective month is on or before month.
func amountForMonth(amounts []RecurringAmount, month string) (RecurringAmount, bool) {
	var (
		best  RecurringAmount
		found bool
	)
	for _, a := range amounts {
		if a.EffectiveMonth <= month && (!found || a.EffectiveMonth > best.EffectiveMonth) {
			best, found = a, true
		}
	}
	return best, found
}

func activeInMonth(r Recurring, month string) bool {
	return r.StartMonth <= month && (r.EndMonth == nil || *r.EndMonth >= month)
}

// clampDay returns the given day of month, capped to the month's last day.
func clampDay(month time.Time, day int) time.Time {
	last := month.AddDate(0, 1, -1).Day()
	if day > last {
		day = last
	}
	return time.Date(month.Year(), month.Month(), day, 0, 0, 0, 0, time.UTC)
}

func applyMonth(r Recurring, month string, logged map[int]int) Recurring {
	r.Active = activeInMonth(r, month)
	if a, ok := amountForMonth(r.Amounts, month); ok && r.Active {
		amount := a.Amount
		r.Amount = &amount
		r.AmountExpression = a.AmountExpression
	}
	if id, ok := logged[r.ID]; ok && r.Active {
		r.LoggedTransactionID = &id
	}
	return r
}

func (s *service) ListRecurring(ctx context.Context, month string) ([]Recurring, error) {
	if month == "" {
		month = s.currentMonth()
	}
	m, err := parseMonth("month", month)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListRecurring(ctx)
	if err != nil {
		return nil, err
	}
	from, to := monthRange(m)
	logged, err := s.repo.LoggedRecurring(ctx, from, to)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i] = applyMonth(items[i], month, logged)
	}
	return items, nil
}

func (s *service) buildRecurring(req RecurringRequest) (Recurring, error) {
	desc, err := validateText("description", req.Description)
	if err != nil {
		return Recurring{}, err
	}
	if req.DayOfMonth < 1 || req.DayOfMonth > 31 {
		return Recurring{}, fmt.Errorf("%w: day_of_month must be between 1 and 31", core.ErrInvalidInput)
	}
	if _, err := parseMonth("start_month", req.StartMonth); err != nil {
		return Recurring{}, err
	}
	if req.EndMonth != nil {
		if _, err := parseMonth("end_month", *req.EndMonth); err != nil {
			return Recurring{}, err
		}
		if *req.EndMonth < req.StartMonth {
			return Recurring{}, fmt.Errorf("%w: end_month must not be before start_month", core.ErrInvalidInput)
		}
	}
	return Recurring{
		Description: desc,
		TypeID:      req.TypeID,
		CategoryID:  req.CategoryID,
		CardID:      req.CardID,
		DayOfMonth:  req.DayOfMonth,
		StartMonth:  req.StartMonth,
		EndMonth:    req.EndMonth,
	}, nil
}

func buildRecurringAmount(month, expression string) (RecurringAmount, error) {
	expr := strings.TrimSpace(expression)
	if len(expr) > maxTextLength {
		return RecurringAmount{}, fmt.Errorf("%w: amount_expression exceeds %d characters", core.ErrInvalidInput, maxTextLength)
	}
	amount, err := EvaluateExpression(expr)
	if err != nil {
		return RecurringAmount{}, err
	}
	return RecurringAmount{EffectiveMonth: month, Amount: amount, AmountExpression: expr}, nil
}

func (s *service) CreateRecurring(ctx context.Context, req RecurringRequest) (Recurring, error) {
	r, err := s.buildRecurring(req)
	if err != nil {
		return Recurring{}, err
	}
	first, err := buildRecurringAmount(r.StartMonth, req.AmountExpression)
	if err != nil {
		return Recurring{}, err
	}
	r.Amounts = []RecurringAmount{first}
	created, err := s.repo.CreateRecurring(ctx, r)
	if err != nil {
		return Recurring{}, mapDBError(err)
	}
	return applyMonth(created, s.currentMonth(), nil), nil
}

func (s *service) UpdateRecurring(ctx context.Context, id int, req RecurringRequest) (Recurring, error) {
	r, err := s.buildRecurring(req)
	if err != nil {
		return Recurring{}, err
	}
	if err := s.repo.UpdateRecurring(ctx, id, r); err != nil {
		return Recurring{}, mapDBError(err)
	}
	updated, err := s.repo.GetRecurring(ctx, id)
	if err != nil {
		return Recurring{}, err
	}
	return applyMonth(updated, s.currentMonth(), nil), nil
}

func (s *service) DeleteRecurring(ctx context.Context, id int) error {
	return s.repo.DeleteRecurring(ctx, id)
}

func (s *service) SetRecurringAmount(ctx context.Context, id int, month string, req RecurringAmountRequest) (Recurring, error) {
	if _, err := parseMonth("month", month); err != nil {
		return Recurring{}, err
	}
	a, err := buildRecurringAmount(month, req.AmountExpression)
	if err != nil {
		return Recurring{}, err
	}
	if _, err := s.repo.GetRecurring(ctx, id); err != nil {
		return Recurring{}, err
	}
	if err := s.repo.UpsertRecurringAmount(ctx, id, a); err != nil {
		return Recurring{}, err
	}
	updated, err := s.repo.GetRecurring(ctx, id)
	if err != nil {
		return Recurring{}, err
	}
	return applyMonth(updated, s.currentMonth(), nil), nil
}

func (s *service) DeleteRecurringAmount(ctx context.Context, id int, month string) (Recurring, error) {
	if _, err := parseMonth("month", month); err != nil {
		return Recurring{}, err
	}
	r, err := s.repo.GetRecurring(ctx, id)
	if err != nil {
		return Recurring{}, err
	}
	if len(r.Amounts) <= 1 {
		return Recurring{}, fmt.Errorf("%w: a recurring item needs at least one amount", core.ErrInvalidInput)
	}
	if err := s.repo.DeleteRecurringAmount(ctx, id, month); err != nil {
		return Recurring{}, err
	}
	updated, err := s.repo.GetRecurring(ctx, id)
	if err != nil {
		return Recurring{}, err
	}
	return applyMonth(updated, s.currentMonth(), nil), nil
}

func (s *service) LogRecurring(ctx context.Context, id int, req LogRecurringRequest) (Transaction, error) {
	m, err := parseMonth("month", req.Month)
	if err != nil {
		return Transaction{}, err
	}
	r, err := s.repo.GetRecurring(ctx, id)
	if err != nil {
		return Transaction{}, err
	}
	if !activeInMonth(r, req.Month) {
		return Transaction{}, fmt.Errorf("%w: item is not active in %s", core.ErrInvalidInput, req.Month)
	}
	from, to := monthRange(m)
	logged, err := s.repo.LoggedRecurring(ctx, from, to)
	if err != nil {
		return Transaction{}, err
	}
	if _, ok := logged[id]; ok {
		return Transaction{}, fmt.Errorf("%w: already logged for %s", core.ErrInvalidInput, req.Month)
	}

	amount, ok := amountForMonth(r.Amounts, req.Month)
	if !ok {
		return Transaction{}, fmt.Errorf("%w: no amount set for %s", core.ErrInvalidInput, req.Month)
	}
	if strings.TrimSpace(req.AmountExpression) != "" {
		if amount, err = buildRecurringAmount(req.Month, req.AmountExpression); err != nil {
			return Transaction{}, err
		}
	}
	date := clampDay(m, r.DayOfMonth)
	if req.Date != nil {
		date = *req.Date
	}
	created, err := s.repo.CreateTransaction(ctx, Transaction{
		Date:             date,
		Description:      r.Description,
		Amount:           amount.Amount,
		AmountExpression: amount.AmountExpression,
		TypeID:           r.TypeID,
		CategoryID:       r.CategoryID,
		CardID:           r.CardID,
		RecurringID:      &r.ID,
	})
	return created, mapDBError(err)
}

// Summary

func withDerived(t Totals) Totals {
	t.Gross = t.Fixed + t.Spend
	t.Net = t.Gross - t.NonSalary
	return t
}

// fillDaily returns one point per day from the first of the month through upto,
// with a running total of the spend recorded so far.
func fillDaily(month, upto time.Time, daily []DailyTotal) []DailyTotal {
	byDate := make(map[string]float64, len(daily))
	for _, d := range daily {
		byDate[d.Date] = d.Amount
	}
	out := []DailyTotal{}
	var running float64
	for day := month; !day.After(upto) && day.Month() == month.Month(); day = day.AddDate(0, 0, 1) {
		key := day.Format("2006-01-02")
		running += byDate[key]
		out = append(out, DailyTotal{Date: key, Amount: byDate[key], Cumulative: running})
	}
	return out
}

func (s *service) Summary(ctx context.Context, month string) (Summary, error) {
	if month == "" {
		month = s.currentMonth()
	}
	m, err := parseMonth("month", month)
	if err != nil {
		return Summary{}, err
	}
	from, to := monthRange(m)
	prevFrom := m.AddDate(0, -1, 0)

	current, err := s.repo.MonthTotals(ctx, from, to)
	if err != nil {
		return Summary{}, err
	}
	previous, err := s.repo.MonthTotals(ctx, prevFrom, from)
	if err != nil {
		return Summary{}, err
	}
	categories, err := s.repo.CategoryTotals(ctx, from, to)
	if err != nil {
		return Summary{}, err
	}
	daily, err := s.repo.DailyTotals(ctx, from, to)
	if err != nil {
		return Summary{}, err
	}
	budget, err := s.repo.GetBudget(ctx, month)
	if err != nil {
		return Summary{}, err
	}

	byCategory := make(map[int]float64, len(budget.Categories))
	for _, b := range budget.Categories {
		byCategory[b.CategoryID] = b.Amount
	}
	seen := map[int]bool{}
	for i := range categories {
		if id := categories[i].CategoryID; id != nil {
			seen[*id] = true
			if amount, ok := byCategory[*id]; ok {
				categories[i].Budget = &amount
			}
		}
	}
	// Categories with a budget but no spending still get a row so the UI can show an empty bar.
	for _, b := range budget.Categories {
		if !seen[b.CategoryID] {
			id, amount := b.CategoryID, b.Amount
			categories = append(categories, CategoryTotal{CategoryID: &id, Budget: &amount})
		}
	}
	sort.SliceStable(categories, func(i, j int) bool {
		return categories[i].Spend+categories[i].Fixed > categories[j].Spend+categories[j].Fixed
	})

	upto := to.AddDate(0, 0, -1)
	if today := s.now().UTC(); today.Before(to) {
		upto = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	}

	sum := Summary{
		Month:      month,
		Current:    withDerived(current),
		Previous:   withDerived(previous),
		Budget:     budget.Total,
		Categories: categories,
		Daily:      fillDaily(m, upto, daily),
	}
	if budget.Total != nil {
		remaining := *budget.Total - sum.Current.Net
		sum.Remaining = &remaining
	}
	return sum, nil
}

var _ Service = (*service)(nil)
