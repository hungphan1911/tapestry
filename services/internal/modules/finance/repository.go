package finance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/hungphan1911/tapestry/services/internal/core"
	"strings"
)

type Repository interface {
	ListTypes(ctx context.Context) ([]Type, error)

	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, description string) (Category, error)
	UpdateCategory(ctx context.Context, id int, description string) (Category, error)
	DeleteCategory(ctx context.Context, id int) error

	ListCards(ctx context.Context) ([]Card, error)
	CreateCard(ctx context.Context, name string) (Card, error)
	UpdateCard(ctx context.Context, id int, name string) (Card, error)
	DeleteCard(ctx context.Context, id int) error

	ListTransactions(ctx context.Context, f TransactionFilter) ([]Transaction, error)
	GetTransaction(ctx context.Context, id int) (Transaction, error)
	CreateTransaction(ctx context.Context, t Transaction) (Transaction, error)
	UpdateTransaction(ctx context.Context, id int, t Transaction) (Transaction, error)
	DeleteTransaction(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

const transactionColumns = `id, date, description, amount, amount_expression, type_id, category_id, card_id`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTransaction(s rowScanner) (Transaction, error) {
	var (
		t                          Transaction
		desc, expr                 sql.NullString
		amount                     sql.NullFloat64
		typeID, categoryID, cardID sql.NullInt64
	)
	if err := s.Scan(&t.ID, &t.Date, &desc, &amount, &expr, &typeID, &categoryID, &cardID); err != nil {
		return Transaction{}, err
	}
	t.Description = desc.String
	t.AmountExpression = expr.String
	t.Amount = amount.Float64
	t.TypeID = nullIntPtr(typeID)
	t.CategoryID = nullIntPtr(categoryID)
	t.CardID = nullIntPtr(cardID)
	return t, nil
}

func nullIntPtr(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int64)
	return &v
}

func checkAffected(res sql.Result, op string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if n == 0 {
		return fmt.Errorf("%s: %w", op, core.ErrNotFound)
	}
	return nil
}

// Types

func (r *repository) ListTypes(ctx context.Context) ([]Type, error) {
	rows, err := r.db.QueryContext(ctx, `select id, description from pft_type order by id`)
	if err != nil {
		return nil, fmt.Errorf("list types: %w", err)
	}
	defer rows.Close()

	types := []Type{}
	for rows.Next() {
		var t Type
		if err := rows.Scan(&t.ID, &t.Description); err != nil {
			return nil, fmt.Errorf("list types: %w", err)
		}
		types = append(types, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list types: %w", err)
	}
	return types, nil
}

// Categories

func (r *repository) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, `select id, description from pft_category order by id`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Description); err != nil {
			return nil, fmt.Errorf("list categories: %w", err)
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

func (r *repository) CreateCategory(ctx context.Context, description string) (Category, error) {
	var c Category
	err := r.db.QueryRowContext(ctx,
		`insert into pft_category (description) values ($1) returning id, description`,
		description).Scan(&c.ID, &c.Description)
	if err != nil {
		return Category{}, fmt.Errorf("create category: %w", err)
	}
	return c, nil
}

func (r *repository) UpdateCategory(ctx context.Context, id int, description string) (Category, error) {
	var c Category
	err := r.db.QueryRowContext(ctx,
		`update pft_category set description = $1 where id = $2 returning id, description`,
		description, id).Scan(&c.ID, &c.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return Category{}, fmt.Errorf("update category: %w", core.ErrNotFound)
	}
	if err != nil {
		return Category{}, fmt.Errorf("update category: %w", err)
	}
	return c, nil
}

func (r *repository) DeleteCategory(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `delete from pft_category where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return checkAffected(res, "delete category")
}

// Cards

func (r *repository) ListCards(ctx context.Context) ([]Card, error) {
	rows, err := r.db.QueryContext(ctx, `select id, card_name from pft_card order by id`)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var c Card
		if err := rows.Scan(&c.ID, &c.CardName); err != nil {
			return nil, fmt.Errorf("list cards: %w", err)
		}
		cards = append(cards, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	return cards, nil
}

func (r *repository) CreateCard(ctx context.Context, name string) (Card, error) {
	var c Card
	err := r.db.QueryRowContext(ctx,
		`insert into pft_card (card_name) values ($1) returning id, card_name`,
		name).Scan(&c.ID, &c.CardName)
	if err != nil {
		return Card{}, fmt.Errorf("create card: %w", err)
	}
	return c, nil
}

func (r *repository) UpdateCard(ctx context.Context, id int, name string) (Card, error) {
	var c Card
	err := r.db.QueryRowContext(ctx,
		`update pft_card set card_name = $1 where id = $2 returning id, card_name`,
		name, id).Scan(&c.ID, &c.CardName)
	if errors.Is(err, sql.ErrNoRows) {
		return Card{}, fmt.Errorf("update card: %w", core.ErrNotFound)
	}
	if err != nil {
		return Card{}, fmt.Errorf("update card: %w", err)
	}
	return c, nil
}

func (r *repository) DeleteCard(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `delete from pft_card where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	return checkAffected(res, "delete card")
}

// Transactions

func (r *repository) ListTransactions(ctx context.Context, f TransactionFilter) ([]Transaction, error) {
	var (
		conds []string
		args  []any
	)
	add := func(cond string, v any) {
		args = append(args, v)
		conds = append(conds, fmt.Sprintf(cond, len(args)))
	}
	if f.From != nil {
		add("date >= $%d", *f.From)
	}
	if f.To != nil {
		add("date <= $%d", *f.To)
	}
	if f.CategoryID != nil {
		add("category_id = $%d", *f.CategoryID)
	}
	if f.TypeID != nil {
		add("type_id = $%d", *f.TypeID)
	}
	if f.CardID != nil {
		add("card_id = $%d", *f.CardID)
	}

	query := `select ` + transactionColumns + ` from pft_transaction`
	if len(conds) > 0 {
		query += ` where ` + strings.Join(conds, " AND ")
	}
	query += ` order by date desc, id desc`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("list transactions: %w", err)
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	return transactions, nil
}

func (r *repository) GetTransaction(ctx context.Context, id int) (Transaction, error) {
	row := r.db.QueryRowContext(ctx,
		`select `+transactionColumns+` from pft_transaction where id = $1`, id)
	t, err := scanTransaction(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Transaction{}, fmt.Errorf("get transaction: %w", core.ErrNotFound)
	}
	if err != nil {
		return Transaction{}, fmt.Errorf("get transaction: %w", err)
	}
	return t, nil
}

func (r *repository) CreateTransaction(ctx context.Context, t Transaction) (Transaction, error) {
	row := r.db.QueryRowContext(ctx,
		`insert into pft_transaction (date, description, amount, amount_expression, type_id, category_id, card_id)
		 values ($1, $2, $3, $4, $5, $6, $7)
		 returning `+transactionColumns,
		t.Date, t.Description, t.Amount, t.AmountExpression, t.TypeID, t.CategoryID, t.CardID)
	created, err := scanTransaction(row)
	if err != nil {
		return Transaction{}, fmt.Errorf("create transaction: %w", err)
	}
	return created, nil
}

func (r *repository) UpdateTransaction(ctx context.Context, id int, t Transaction) (Transaction, error) {
	row := r.db.QueryRowContext(ctx,
		`update pft_transaction
		 set date = $1, description = $2, amount = $3, amount_expression = $4,
		     type_id = $5, category_id = $6, card_id = $7
		 where id = $8
		 returning `+transactionColumns,
		t.Date, t.Description, t.Amount, t.AmountExpression, t.TypeID, t.CategoryID, t.CardID, id)
	updated, err := scanTransaction(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Transaction{}, fmt.Errorf("update transaction: %w", core.ErrNotFound)
	}
	if err != nil {
		return Transaction{}, fmt.Errorf("update transaction: %w", err)
	}
	return updated, nil
}

func (r *repository) DeleteTransaction(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `delete from pft_transaction where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}
	return checkAffected(res, "delete transaction")
}

var _ Repository = (*repository)(nil)
