import { Plus, Repeat } from 'lucide-react'
import { useState } from 'react'
import { Button } from '@/components/ui/Button'
import { Card } from '@/components/ui/Card'
import { Input, Select } from '@/components/ui/Input'
import { cn } from '@/lib/cn'
import { currentMonth, formatDate, formatMoney, todayInput } from '@/lib/format'
import { typeLabels } from '../labels'
import { useCards, useCategories, useTransactions, useTypes } from '../queries'
import type { Transaction, TypeName } from '../types'
import { MonthPicker } from './MonthPicker'
import { TransactionDrawer } from './TransactionDrawer'

type TransactionsTabProps = { month: string; onMonthChange: (month: string) => void }

const inflows: TypeName[] = ['income', 'non-salary']

export function TransactionsTab({ month, onMonthChange }: TransactionsTabProps) {
  const [categoryId, setCategoryId] = useState('')
  const [typeId, setTypeId] = useState('')
  const [cardId, setCardId] = useState('')
  const [search, setSearch] = useState('')
  const [editing, setEditing] = useState<Transaction | null>(null)
  const [adding, setAdding] = useState(false)

  const { data: types = [] } = useTypes()
  const { data: categories = [] } = useCategories()
  const { data: cards = [] } = useCards()
  const { data, isPending, error } = useTransactions({
    month,
    categoryId: categoryId ? Number(categoryId) : undefined,
    typeId: typeId ? Number(typeId) : undefined,
    cardId: cardId ? Number(cardId) : undefined,
  })

  const typeById = new Map(types.map((t) => [t.id, t.description]))
  const categoryById = new Map(categories.map((c) => [c.id, c.description]))
  const cardById = new Map(cards.map((c) => [c.id, c.card_name]))

  const needle = search.trim().toLowerCase()
  const rows = (data ?? []).filter((t) => !needle || t.description.toLowerCase().includes(needle))

  const defaultDate = month === currentMonth() ? todayInput() : `${month}-01`
  const typeName = (t: Transaction) => (t.type_id ? typeById.get(t.type_id) : undefined)
  const amountClass = (t: Transaction) => {
    const name = typeName(t)
    return name && inflows.includes(name) ? 'text-positive' : ''
  }

  return (
    <>
      <div className="mb-4 flex flex-wrap items-center gap-3">
        <MonthPicker month={month} onChange={onMonthChange} />
        <Input
          className="w-48 max-md:w-full"
          placeholder="Search description"
          aria-label="Search description"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Select
          className="w-40"
          aria-label="Filter by type"
          value={typeId}
          onChange={(e) => setTypeId(e.target.value)}
        >
          <option value="">All types</option>
          {types.map((t) => (
            <option key={t.id} value={t.id}>
              {typeLabels[t.description] ?? t.description}
            </option>
          ))}
        </Select>
        <Select
          className="w-40"
          aria-label="Filter by category"
          value={categoryId}
          onChange={(e) => setCategoryId(e.target.value)}
        >
          <option value="">All categories</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.description}
            </option>
          ))}
        </Select>
        <Select
          className="w-40"
          aria-label="Filter by card"
          value={cardId}
          onChange={(e) => setCardId(e.target.value)}
        >
          <option value="">All cards</option>
          {cards.map((c) => (
            <option key={c.id} value={c.id}>
              {c.card_name}
            </option>
          ))}
        </Select>
        <Button className="ml-auto" onClick={() => setAdding(true)}>
          <Plus className="size-4" aria-hidden />
          Add transaction
        </Button>
      </div>

      <Card className="p-0 max-md:p-0">
        {isPending && <p className="text-muted p-6 text-sm">Loading…</p>}
        {error && <p className="text-negative p-6 text-sm">Could not load transactions.</p>}
        {data && rows.length === 0 && (
          <p className="text-muted p-6 text-sm">No transactions for this month.</p>
        )}

        {rows.length > 0 && (
          <table className="w-full text-sm max-md:hidden">
            <thead>
              <tr className="border-line text-muted border-b text-left text-xs">
                <th className="px-4 py-3 font-medium">Date</th>
                <th className="px-4 py-3 font-medium">Description</th>
                <th className="px-4 py-3 text-right font-medium">Amount</th>
                <th className="px-4 py-3 font-medium">Type</th>
                <th className="px-4 py-3 font-medium">Card</th>
                <th className="px-4 py-3 font-medium">Category</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((t) => (
                <tr
                  key={t.id}
                  className="border-line hover:bg-subtle/60 cursor-pointer border-b last:border-b-0"
                  onClick={() => setEditing(t)}
                >
                  <td className="px-4 py-2.5 whitespace-nowrap tabular-nums">
                    {formatDate(t.date)}
                  </td>
                  <td className="px-4 py-2.5">
                    <span className="inline-flex items-center gap-1.5">
                      {t.description}
                      {t.recurring_id && (
                        <Repeat
                          className="text-muted size-3.5"
                          aria-label="From a recurring item"
                        />
                      )}
                    </span>
                  </td>
                  <td className={cn('px-4 py-2.5 text-right tabular-nums', amountClass(t))}>
                    {formatMoney(t.amount)}
                  </td>
                  <td className="px-4 py-2.5">{typeName(t) ?? '—'}</td>
                  <td className="px-4 py-2.5">{(t.card_id && cardById.get(t.card_id)) || '—'}</td>
                  <td className="px-4 py-2.5">
                    {(t.category_id && categoryById.get(t.category_id)) || '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}

        <ul className="divide-line hidden divide-y max-md:block">
          {rows.map((t) => (
            <li key={t.id}>
              <button
                type="button"
                className="flex w-full items-start justify-between gap-3 p-4 text-left"
                onClick={() => setEditing(t)}
              >
                <span>
                  <span className="block text-sm font-medium">{t.description || '—'}</span>
                  <span className="text-muted block text-xs">
                    {formatDate(t.date)} ·{' '}
                    {(t.category_id && categoryById.get(t.category_id)) || 'No category'}
                  </span>
                </span>
                <span className={cn('text-sm tabular-nums', amountClass(t))}>
                  {formatMoney(t.amount)}
                </span>
              </button>
            </li>
          ))}
        </ul>
      </Card>

      <TransactionDrawer
        open={adding || editing !== null}
        transaction={editing}
        defaultDate={defaultDate}
        onClose={() => {
          setAdding(false)
          setEditing(null)
        }}
      />
    </>
  )
}
