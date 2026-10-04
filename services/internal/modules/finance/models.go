package finance

import "time"

// Entities

type Type struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Category struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Card struct {
	ID       int    `json:"id"`
	CardName string `json:"card_name"`
}

type Transaction struct {
	ID               int       `json:"id"`
	Date             time.Time `json:"date"`
	Description      string    `json:"description"`
	Amount           float64   `json:"amount"`
	AmountExpression string    `json:"amount_expression"`
	TypeID           *int      `json:"type_id"`
	CategoryID       *int      `json:"category_id"`
	CardID           *int      `json:"card_id"`
	RecurringID      *int      `json:"recurring_id"`
}

// Budgets: months are "YYYY-MM"; a nil category is the month's total budget.

type BudgetCategory struct {
	CategoryID int     `json:"category_id"`
	Amount     float64 `json:"amount"`
}

type Budget struct {
	Month      string           `json:"month"`
	Total      *float64         `json:"total"`
	Categories []BudgetCategory `json:"categories"`
}

// Recurring items: anticipated salary and fixed charges. Not transactions.

type RecurringAmount struct {
	EffectiveMonth   string  `json:"effective_month"`
	Amount           float64 `json:"amount"`
	AmountExpression string  `json:"amount_expression"`
}

type Recurring struct {
	ID          int               `json:"id"`
	Description string            `json:"description"`
	TypeID      *int              `json:"type_id"`
	CategoryID  *int              `json:"category_id"`
	CardID      *int              `json:"card_id"`
	DayOfMonth  int               `json:"day_of_month"`
	StartMonth  string            `json:"start_month"`
	EndMonth    *string           `json:"end_month"`
	Amounts     []RecurringAmount `json:"amounts"`

	// Filled for a requested month by the service.
	Active              bool     `json:"active"`
	Amount              *float64 `json:"amount"`
	AmountExpression    string   `json:"amount_expression"`
	LoggedTransactionID *int     `json:"logged_transaction_id"`
}

// Summary

type Totals struct {
	Fixed     float64 `json:"fixed"`
	Spend     float64 `json:"spend"`
	NonSalary float64 `json:"non_salary"`
	Income    float64 `json:"income"`
	// Gross is fixed + spend; Net is gross minus non-salary (cashback, refunds).
	Gross float64 `json:"gross"`
	Net   float64 `json:"net"`
}

type CategoryTotal struct {
	CategoryID *int     `json:"category_id"`
	Spend      float64  `json:"spend"`
	Fixed      float64  `json:"fixed"`
	Budget     *float64 `json:"budget"`
}

type DailyTotal struct {
	Date       string  `json:"date"`
	Amount     float64 `json:"amount"`
	Cumulative float64 `json:"cumulative"`
}

type Summary struct {
	Month      string          `json:"month"`
	Current    Totals          `json:"current"`
	Previous   Totals          `json:"previous"`
	Budget     *float64        `json:"budget"`
	Remaining  *float64        `json:"remaining"`
	Categories []CategoryTotal `json:"categories"`
	Daily      []DailyTotal    `json:"daily"`
}

// Filters

type TransactionFilter struct {
	From       *time.Time
	To         *time.Time
	CategoryID *int
	TypeID     *int
	CardID     *int
}

// Request bodies

type CategoryRequest struct {
	Description string `json:"description"`
}

type CardRequest struct {
	CardName string `json:"card_name"`
}

type TransactionRequest struct {
	Date             *time.Time `json:"date"`
	Description      string     `json:"description"`
	AmountExpression string     `json:"amount_expression"`
	TypeID           *int       `json:"type_id"`
	CategoryID       *int       `json:"category_id"`
	CardID           *int       `json:"card_id"`
}

type BudgetRequest struct {
	Total      *float64         `json:"total"`
	Categories []BudgetCategory `json:"categories"`
}

// RecurringRequest's amount_expression sets the first amount on create and is ignored on update.
type RecurringRequest struct {
	Description      string  `json:"description"`
	TypeID           *int    `json:"type_id"`
	CategoryID       *int    `json:"category_id"`
	CardID           *int    `json:"card_id"`
	DayOfMonth       int     `json:"day_of_month"`
	StartMonth       string  `json:"start_month"`
	EndMonth         *string `json:"end_month"`
	AmountExpression string  `json:"amount_expression"`
}

type RecurringAmountRequest struct {
	AmountExpression string `json:"amount_expression"`
}

// LogRecurringRequest turns a recurring item into a transaction for Month.
// Date and AmountExpression override the item's defaults when set.
type LogRecurringRequest struct {
	Month            string     `json:"month"`
	Date             *time.Time `json:"date"`
	AmountExpression string     `json:"amount_expression"`
}
