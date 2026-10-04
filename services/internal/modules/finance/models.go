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
