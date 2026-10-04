import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as api from './api'
import type { BudgetInput, RecurringInput, TransactionFilter, TransactionInput } from './types'

const root = ['finance'] as const

export const keys = {
  types: [...root, 'types'] as const,
  categories: [...root, 'categories'] as const,
  cards: [...root, 'cards'] as const,
  transactions: (filter: TransactionFilter) => [...root, 'transactions', filter] as const,
  summary: (month: string) => [...root, 'summary', month] as const,
  budget: (month: string) => [...root, 'budget', month] as const,
  recurring: (month: string) => [...root, 'recurring', month] as const,
}

// Most writes touch several views (a transaction changes the summary and recurring status),
// so every mutation refreshes all finance queries.
function useRefresh() {
  const queryClient = useQueryClient()
  return () => queryClient.invalidateQueries({ queryKey: root })
}

export const useTypes = () =>
  useQuery({ queryKey: keys.types, queryFn: api.listTypes, staleTime: Infinity })
export const useCategories = () =>
  useQuery({ queryKey: keys.categories, queryFn: api.listCategories })
export const useCards = () => useQuery({ queryKey: keys.cards, queryFn: api.listCards })

export const useTransactions = (filter: TransactionFilter) =>
  useQuery({ queryKey: keys.transactions(filter), queryFn: () => api.listTransactions(filter) })
export const useSummary = (month: string) =>
  useQuery({ queryKey: keys.summary(month), queryFn: () => api.getSummary(month) })
export const useBudget = (month: string) =>
  useQuery({ queryKey: keys.budget(month), queryFn: () => api.getBudget(month) })
export const useRecurring = (month: string) =>
  useQuery({ queryKey: keys.recurring(month), queryFn: () => api.listRecurring(month) })

export function useSaveCategory() {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (v: { id?: number; description: string }) =>
      v.id ? api.updateCategory(v.id, v.description) : api.createCategory(v.description),
    onSuccess: refresh,
  })
}
export function useDeleteCategory() {
  const refresh = useRefresh()
  return useMutation({ mutationFn: api.deleteCategory, onSuccess: refresh })
}

export function useSaveCard() {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (v: { id?: number; card_name: string }) =>
      v.id ? api.updateCard(v.id, v.card_name) : api.createCard(v.card_name),
    onSuccess: refresh,
  })
}
export function useDeleteCard() {
  const refresh = useRefresh()
  return useMutation({ mutationFn: api.deleteCard, onSuccess: refresh })
}

export function useSaveTransaction() {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (v: { id?: number; input: TransactionInput }) =>
      v.id ? api.updateTransaction(v.id, v.input) : api.createTransaction(v.input),
    onSuccess: refresh,
  })
}
export function useDeleteTransaction() {
  const refresh = useRefresh()
  return useMutation({ mutationFn: api.deleteTransaction, onSuccess: refresh })
}

export function useSetBudget(month: string) {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (input: BudgetInput) => api.setBudget(month, input),
    onSuccess: refresh,
  })
}

export function useSaveRecurring() {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (v: { id?: number; input: RecurringInput }) =>
      v.id ? api.updateRecurring(v.id, v.input) : api.createRecurring(v.input),
    onSuccess: refresh,
  })
}
export function useDeleteRecurring() {
  const refresh = useRefresh()
  return useMutation({ mutationFn: api.deleteRecurring, onSuccess: refresh })
}
export function useSetRecurringAmount() {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (v: { id: number; month: string; expression: string }) =>
      api.setRecurringAmount(v.id, v.month, v.expression),
    onSuccess: refresh,
  })
}
export function useDeleteRecurringAmount() {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (v: { id: number; month: string }) => api.deleteRecurringAmount(v.id, v.month),
    onSuccess: refresh,
  })
}
export function useLogRecurring() {
  const refresh = useRefresh()
  return useMutation({
    mutationFn: (v: { id: number; month: string }) => api.logRecurring(v.id, v.month),
    onSuccess: refresh,
  })
}
