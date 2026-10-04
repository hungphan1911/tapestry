// Mirrors services/internal/modules/finance/models.go. Months are "YYYY-MM".

export type TypeName = 'spend' | 'fixed' | 'income' | 'non-salary'

export type TxType = { id: number; description: TypeName }
export type Category = { id: number; description: string }
export type Card = { id: number; card_name: string }

export type Transaction = {
  id: number
  date: string
  description: string
  amount: number
  amount_expression: string
  type_id: number | null
  category_id: number | null
  card_id: number | null
  recurring_id: number | null
}

export type TransactionInput = {
  date: string
  description: string
  amount_expression: string
  type_id: number | null
  category_id: number | null
  card_id: number | null
}

export type TransactionFilter = {
  month: string
  categoryId?: number
  typeId?: number
  cardId?: number
}

export type BudgetCategory = { category_id: number; amount: number }
export type Budget = { month: string; total: number | null; categories: BudgetCategory[] }
export type BudgetInput = { total: number | null; categories: BudgetCategory[] }

export type RecurringAmount = {
  effective_month: string
  amount: number
  amount_expression: string
}

export type Recurring = {
  id: number
  description: string
  type_id: number | null
  category_id: number | null
  card_id: number | null
  day_of_month: number
  start_month: string
  end_month: string | null
  amounts: RecurringAmount[]
  active: boolean
  amount: number | null
  amount_expression: string
  logged_transaction_id: number | null
}

export type RecurringInput = {
  description: string
  type_id: number | null
  category_id: number | null
  card_id: number | null
  day_of_month: number
  start_month: string
  end_month: string | null
  amount_expression: string
}

export type Totals = {
  fixed: number
  spend: number
  non_salary: number
  income: number
  gross: number
  net: number
}

export type CategoryTotal = {
  category_id: number | null
  spend: number
  fixed: number
  budget: number | null
}

export type DailyTotal = { date: string; amount: number; cumulative: number }

export type Summary = {
  month: string
  current: Totals
  previous: Totals
  budget: number | null
  remaining: number | null
  categories: CategoryTotal[]
  daily: DailyTotal[]
}
