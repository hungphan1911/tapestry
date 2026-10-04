import type { TypeName } from './types'

export const typeLabels: Record<TypeName, string> = {
  spend: 'Spend',
  fixed: 'Fixed',
  income: 'Income',
  'non-salary': 'Non-salary',
}

export const tabLabels = {
  overview: 'Overview',
  transactions: 'Transactions',
  budget: 'Budget',
  recurring: 'Recurring',
  settings: 'Settings',
} as const
