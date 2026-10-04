import { apiFetch } from '@/api/client'
import { monthBounds } from '@/lib/format'
import type {
  Budget,
  BudgetInput,
  Card,
  Category,
  Recurring,
  RecurringInput,
  Summary,
  Transaction,
  TransactionFilter,
  TransactionInput,
  TxType,
} from './types'

const base = '/finance'

const json = (method: string, body?: unknown): RequestInit => ({
  method,
  body: body === undefined ? undefined : JSON.stringify(body),
})

export const listTypes = () => apiFetch<TxType[]>(`${base}/types`)

export const listCategories = () => apiFetch<Category[]>(`${base}/categories/`)
export const createCategory = (description: string) =>
  apiFetch<Category>(`${base}/categories/`, json('POST', { description }))
export const updateCategory = (id: number, description: string) =>
  apiFetch<Category>(`${base}/categories/${id}`, json('PUT', { description }))
export const deleteCategory = (id: number) =>
  apiFetch<void>(`${base}/categories/${id}`, json('DELETE'))

export const listCards = () => apiFetch<Card[]>(`${base}/cards/`)
export const createCard = (card_name: string) =>
  apiFetch<Card>(`${base}/cards/`, json('POST', { card_name }))
export const updateCard = (id: number, card_name: string) =>
  apiFetch<Card>(`${base}/cards/${id}`, json('PUT', { card_name }))
export const deleteCard = (id: number) => apiFetch<void>(`${base}/cards/${id}`, json('DELETE'))

export function listTransactions(filter: TransactionFilter) {
  const { from, to } = monthBounds(filter.month)
  const params = new URLSearchParams({ from, to })
  if (filter.categoryId) params.set('category_id', String(filter.categoryId))
  if (filter.typeId) params.set('type_id', String(filter.typeId))
  if (filter.cardId) params.set('card_id', String(filter.cardId))
  return apiFetch<Transaction[]>(`${base}/transactions/?${params}`)
}
export const createTransaction = (input: TransactionInput) =>
  apiFetch<Transaction>(`${base}/transactions/`, json('POST', input))
export const updateTransaction = (id: number, input: TransactionInput) =>
  apiFetch<Transaction>(`${base}/transactions/${id}`, json('PUT', input))
export const deleteTransaction = (id: number) =>
  apiFetch<void>(`${base}/transactions/${id}`, json('DELETE'))

export const getSummary = (month: string) => apiFetch<Summary>(`${base}/summary?month=${month}`)

export const getBudget = (month: string) => apiFetch<Budget>(`${base}/budgets/${month}/`)
export const setBudget = (month: string, input: BudgetInput) =>
  apiFetch<Budget>(`${base}/budgets/${month}/`, json('PUT', input))

export const listRecurring = (month: string) =>
  apiFetch<Recurring[]>(`${base}/recurring/?month=${month}`)
export const createRecurring = (input: RecurringInput) =>
  apiFetch<Recurring>(`${base}/recurring/`, json('POST', input))
export const updateRecurring = (id: number, input: RecurringInput) =>
  apiFetch<Recurring>(`${base}/recurring/${id}`, json('PUT', input))
export const deleteRecurring = (id: number) =>
  apiFetch<void>(`${base}/recurring/${id}`, json('DELETE'))
export const setRecurringAmount = (id: number, month: string, amount_expression: string) =>
  apiFetch<Recurring>(
    `${base}/recurring/${id}/amounts/${month}`,
    json('PUT', { amount_expression }),
  )
export const deleteRecurringAmount = (id: number, month: string) =>
  apiFetch<Recurring>(`${base}/recurring/${id}/amounts/${month}`, json('DELETE'))
export const logRecurring = (id: number, month: string) =>
  apiFetch<Transaction>(`${base}/recurring/${id}/log`, json('POST', { month }))
