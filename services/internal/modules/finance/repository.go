package finance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/hungphan1911/tapestry/services/internal/core"
	"strings"
	"time"
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

	GetBudget(ctx context.Context, month string) (Budget, error)
	ReplaceBudget(ctx context.Context, month string, req BudgetRequest) error

	ListRecurring(ctx context.Context) ([]Recurring, error)
	GetRecurring(ctx context.Context, id int) (Recurring, error)
	CreateRecurring(ctx context.Context, r Recurring) (Recurring, error)
	UpdateRecurring(ctx context.Context, id int, r Recurring) error
	DeleteRecurring(ctx context.Context, id int) error
	UpsertRecurringAmount(ctx context.Context, id int, a RecurringAmount) error
	DeleteRecurringAmount(ctx context.Context, id int, month string) error
	// LoggedRecurring maps recurring id to the first transaction logged for it in [from, to).
	LoggedRecurring(ctx context.Context, from, to time.Time) (map[int]int, error)

	MonthTotals(ctx context.Context, from, to time.Time) (Totals, error)
	CategoryTotals(ctx context.Context, from, to time.Time) ([]CategoryTotal, error)
	DailyTotals(ctx context.Context, from, to time.Time) ([]DailyTotal, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

const transactionColumns = `id, date, description, amount, amount_expression, type_id, category_id, card_id, recurring_id`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTransaction(s rowScanner) (Transaction, error) {
	var (
		t                                       Transaction
		desc, expr                              sql.NullString
		amount                                  sql.NullFloat64
		typeID, categoryID, cardID, recurringID sql.NullInt64
	)
	if err := s.Scan(&t.ID, &t.Date, &desc, &amount, &expr, &typeID, &categoryID, &cardID, &recurringID); err != nil {
		return Transaction{}, err
	}
	t.Description = desc.String
	t.AmountExpression = expr.String
	t.Amount = amount.Float64
	t.TypeID = nullIntPtr(typeID)
	t.CategoryID = nullIntPtr(categoryID)
	t.CardID = nullIntPtr(cardID)
	t.RecurringID = nullIntPtr(recurringID)
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
		`insert into pft_transaction (date, description, amount, amount_expression, type_id, category_id, card_id, recurring_id)
		 values ($1, $2, $3, $4, $5, $6, $7, $8)
		 returning `+transactionColumns,
		t.Date, t.Description, t.Amount, t.AmountExpression, t.TypeID, t.CategoryID, t.CardID, t.RecurringID)
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

// Budgets

func (r *repository) GetBudget(ctx context.Context, month string) (Budget, error) {
	rows, err := r.db.QueryContext(ctx,
		`select category_id, amount from pft_budget where month = $1::date order by category_id`,
		month+"-01")
	if err != nil {
		return Budget{}, fmt.Errorf("get budget: %w", err)
	}
	defer rows.Close()

	b := Budget{Month: month, Categories: []BudgetCategory{}}
	for rows.Next() {
		var (
			categoryID sql.NullInt64
			amount     float64
		)
		if err := rows.Scan(&categoryID, &amount); err != nil {
			return Budget{}, fmt.Errorf("get budget: %w", err)
		}
		if categoryID.Valid {
			b.Categories = append(b.Categories, BudgetCategory{CategoryID: int(categoryID.Int64), Amount: amount})
		} else {
			total := amount
			b.Total = &total
		}
	}
	if err := rows.Err(); err != nil {
		return Budget{}, fmt.Errorf("get budget: %w", err)
	}
	return b, nil
}

func (r *repository) ReplaceBudget(ctx context.Context, month string, req BudgetRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("replace budget: %w", err)
	}
	defer tx.Rollback()

	date := month + "-01"
	if _, err := tx.ExecContext(ctx, `delete from pft_budget where month = $1::date`, date); err != nil {
		return fmt.Errorf("replace budget: %w", err)
	}
	if req.Total != nil {
		if _, err := tx.ExecContext(ctx,
			`insert into pft_budget (month, amount) values ($1::date, $2)`, date, *req.Total); err != nil {
			return fmt.Errorf("replace budget: %w", err)
		}
	}
	for _, c := range req.Categories {
		if _, err := tx.ExecContext(ctx,
			`insert into pft_budget (month, category_id, amount) values ($1::date, $2, $3)`,
			date, c.CategoryID, c.Amount); err != nil {
			return fmt.Errorf("replace budget: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("replace budget: %w", err)
	}
	return nil
}

// Recurring

const recurringColumns = `id, description, type_id, category_id, card_id, day_of_month,
	to_char(start_month, 'YYYY-MM'), to_char(end_month, 'YYYY-MM')`

func (r *repository) queryRecurring(ctx context.Context, id *int) ([]Recurring, error) {
	query := `select ` + recurringColumns + ` from pft_recurring`
	var args []any
	if id != nil {
		query += ` where id = $1`
		args = append(args, *id)
	}
	query += ` order by id`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Recurring{}
	for rows.Next() {
		var (
			it                         Recurring
			typeID, categoryID, cardID sql.NullInt64
			endMonth                   sql.NullString
		)
		if err := rows.Scan(&it.ID, &it.Description, &typeID, &categoryID, &cardID,
			&it.DayOfMonth, &it.StartMonth, &endMonth); err != nil {
			return nil, err
		}
		it.TypeID = nullIntPtr(typeID)
		it.CategoryID = nullIntPtr(categoryID)
		it.CardID = nullIntPtr(cardID)
		if endMonth.Valid {
			it.EndMonth = &endMonth.String
		}
		it.Amounts = []RecurringAmount{}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	amountRows, err := r.db.QueryContext(ctx,
		`select recurring_id, to_char(effective_month, 'YYYY-MM'), amount, amount_expression
		 from pft_recurring_amount order by recurring_id, effective_month`)
	if err != nil {
		return nil, err
	}
	defer amountRows.Close()

	index := make(map[int]*Recurring, len(items))
	for i := range items {
		index[items[i].ID] = &items[i]
	}
	for amountRows.Next() {
		var (
			recurringID int
			a           RecurringAmount
			expr        sql.NullString
		)
		if err := amountRows.Scan(&recurringID, &a.EffectiveMonth, &a.Amount, &expr); err != nil {
			return nil, err
		}
		a.AmountExpression = expr.String
		if it, ok := index[recurringID]; ok {
			it.Amounts = append(it.Amounts, a)
		}
	}
	return items, amountRows.Err()
}

func (r *repository) ListRecurring(ctx context.Context) ([]Recurring, error) {
	items, err := r.queryRecurring(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("list recurring: %w", err)
	}
	return items, nil
}

func (r *repository) GetRecurring(ctx context.Context, id int) (Recurring, error) {
	items, err := r.queryRecurring(ctx, &id)
	if err != nil {
		return Recurring{}, fmt.Errorf("get recurring: %w", err)
	}
	if len(items) == 0 {
		return Recurring{}, fmt.Errorf("get recurring: %w", core.ErrNotFound)
	}
	return items[0], nil
}

func (r *repository) CreateRecurring(ctx context.Context, it Recurring) (Recurring, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Recurring{}, fmt.Errorf("create recurring: %w", err)
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRowContext(ctx,
		`insert into pft_recurring (description, type_id, category_id, card_id, day_of_month, start_month, end_month)
		 values ($1, $2, $3, $4, $5, $6::date, $7::date)
		 returning id`,
		it.Description, it.TypeID, it.CategoryID, it.CardID, it.DayOfMonth,
		it.StartMonth+"-01", monthDatePtr(it.EndMonth)).Scan(&id)
	if err != nil {
		return Recurring{}, fmt.Errorf("create recurring: %w", err)
	}
	for _, a := range it.Amounts {
		if _, err := tx.ExecContext(ctx,
			`insert into pft_recurring_amount (recurring_id, effective_month, amount, amount_expression)
			 values ($1, $2::date, $3, $4)`,
			id, a.EffectiveMonth+"-01", a.Amount, a.AmountExpression); err != nil {
			return Recurring{}, fmt.Errorf("create recurring: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return Recurring{}, fmt.Errorf("create recurring: %w", err)
	}
	return r.GetRecurring(ctx, id)
}

func (r *repository) UpdateRecurring(ctx context.Context, id int, it Recurring) error {
	res, err := r.db.ExecContext(ctx,
		`update pft_recurring
		 set description = $1, type_id = $2, category_id = $3, card_id = $4,
		     day_of_month = $5, start_month = $6::date, end_month = $7::date
		 where id = $8`,
		it.Description, it.TypeID, it.CategoryID, it.CardID, it.DayOfMonth,
		it.StartMonth+"-01", monthDatePtr(it.EndMonth), id)
	if err != nil {
		return fmt.Errorf("update recurring: %w", err)
	}
	return checkAffected(res, "update recurring")
}

func (r *repository) DeleteRecurring(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `delete from pft_recurring where id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete recurring: %w", err)
	}
	return checkAffected(res, "delete recurring")
}

func (r *repository) UpsertRecurringAmount(ctx context.Context, id int, a RecurringAmount) error {
	_, err := r.db.ExecContext(ctx,
		`insert into pft_recurring_amount (recurring_id, effective_month, amount, amount_expression)
		 values ($1, $2::date, $3, $4)
		 on conflict (recurring_id, effective_month)
		 do update set amount = excluded.amount, amount_expression = excluded.amount_expression`,
		id, a.EffectiveMonth+"-01", a.Amount, a.AmountExpression)
	if err != nil {
		return fmt.Errorf("upsert recurring amount: %w", err)
	}
	return nil
}

func (r *repository) DeleteRecurringAmount(ctx context.Context, id int, month string) error {
	res, err := r.db.ExecContext(ctx,
		`delete from pft_recurring_amount where recurring_id = $1 and effective_month = $2::date`,
		id, month+"-01")
	if err != nil {
		return fmt.Errorf("delete recurring amount: %w", err)
	}
	return checkAffected(res, "delete recurring amount")
}

func (r *repository) LoggedRecurring(ctx context.Context, from, to time.Time) (map[int]int, error) {
	rows, err := r.db.QueryContext(ctx,
		`select recurring_id, min(id) from pft_transaction
		 where recurring_id is not null and date >= $1 and date < $2
		 group by recurring_id`, from, to)
	if err != nil {
		return nil, fmt.Errorf("logged recurring: %w", err)
	}
	defer rows.Close()

	logged := map[int]int{}
	for rows.Next() {
		var recurringID, transactionID int
		if err := rows.Scan(&recurringID, &transactionID); err != nil {
			return nil, fmt.Errorf("logged recurring: %w", err)
		}
		logged[recurringID] = transactionID
	}
	return logged, rows.Err()
}

// monthDatePtr turns an optional "YYYY-MM" into an optional date string.
func monthDatePtr(month *string) any {
	if month == nil {
		return nil
	}
	return *month + "-01"
}

// Summary

func (r *repository) MonthTotals(ctx context.Context, from, to time.Time) (Totals, error) {
	var t Totals
	err := r.db.QueryRowContext(ctx,
		`select
		   coalesce(sum(t.amount) filter (where ty.description = 'fixed'), 0),
		   coalesce(sum(t.amount) filter (where ty.description = 'spend'), 0),
		   coalesce(sum(t.amount) filter (where ty.description = 'non-salary'), 0),
		   coalesce(sum(t.amount) filter (where ty.description = 'income'), 0)
		 from pft_transaction t
		 join pft_type ty on ty.id = t.type_id
		 where t.date >= $1 and t.date < $2`, from, to).Scan(&t.Fixed, &t.Spend, &t.NonSalary, &t.Income)
	if err != nil {
		return Totals{}, fmt.Errorf("month totals: %w", err)
	}
	return t, nil
}

func (r *repository) CategoryTotals(ctx context.Context, from, to time.Time) ([]CategoryTotal, error) {
	rows, err := r.db.QueryContext(ctx,
		`select t.category_id,
		   coalesce(sum(t.amount) filter (where ty.description = 'spend'), 0),
		   coalesce(sum(t.amount) filter (where ty.description = 'fixed'), 0)
		 from pft_transaction t
		 join pft_type ty on ty.id = t.type_id
		 where t.date >= $1 and t.date < $2 and ty.description in ('spend', 'fixed')
		 group by t.category_id
		 order by 2 desc`, from, to)
	if err != nil {
		return nil, fmt.Errorf("category totals: %w", err)
	}
	defer rows.Close()

	totals := []CategoryTotal{}
	for rows.Next() {
		var (
			c          CategoryTotal
			categoryID sql.NullInt64
		)
		if err := rows.Scan(&categoryID, &c.Spend, &c.Fixed); err != nil {
			return nil, fmt.Errorf("category totals: %w", err)
		}
		c.CategoryID = nullIntPtr(categoryID)
		totals = append(totals, c)
	}
	return totals, rows.Err()
}

func (r *repository) DailyTotals(ctx context.Context, from, to time.Time) ([]DailyTotal, error) {
	rows, err := r.db.QueryContext(ctx,
		`select to_char(t.date at time zone 'UTC', 'YYYY-MM-DD') as day, sum(t.amount)
		 from pft_transaction t
		 join pft_type ty on ty.id = t.type_id
		 where t.date >= $1 and t.date < $2 and ty.description in ('spend', 'fixed')
		 group by day
		 order by day`, from, to)
	if err != nil {
		return nil, fmt.Errorf("daily totals: %w", err)
	}
	defer rows.Close()

	totals := []DailyTotal{}
	for rows.Next() {
		var d DailyTotal
		if err := rows.Scan(&d.Date, &d.Amount); err != nil {
			return nil, fmt.Errorf("daily totals: %w", err)
		}
		totals = append(totals, d)
	}
	return totals, rows.Err()
}

var _ Repository = (*repository)(nil)
