import { Plus } from 'lucide-react'
import { useState } from 'react'
import { Button } from '@/components/ui/Button'
import { Card } from '@/components/ui/Card'
import { cn } from '@/lib/cn'
import { errorMessage } from '@/lib/errorMessage'
import { formatMoney } from '@/lib/format'
import { typeLabels } from '../labels'
import { useLogRecurring, useRecurring, useSummary, useTypes } from '../queries'
import type { Recurring, TypeName } from '../types'
import { MonthPicker } from './MonthPicker'
import { RecurringDrawer } from './RecurringDrawer'

type RecurringTabProps = { month: string; onMonthChange: (month: string) => void }

export function RecurringTab({ month, onMonthChange }: RecurringTabProps) {
  const { data: items, isPending, error } = useRecurring(month)
  const { data: summary } = useSummary(month)
  const { data: types = [] } = useTypes()
  const log = useLogRecurring()
  const [editingId, setEditingId] = useState<number | null>(null)
  const [adding, setAdding] = useState(false)

  const typeName = (item: Recurring): TypeName | undefined =>
    types.find((t) => t.id === item.type_id)?.description

  const expected = (name: TypeName) =>
    (items ?? [])
      .filter((i) => i.active && typeName(i) === name)
      .reduce((sum, i) => sum + (i.amount ?? 0), 0)

  const editing = items?.find((i) => i.id === editingId) ?? null

  return (
    <>
      <div className="mb-4 flex flex-wrap items-center gap-3">
        <MonthPicker month={month} onChange={onMonthChange} />
        <Button className="ml-auto" onClick={() => setAdding(true)}>
          <Plus className="size-4" aria-hidden />
          Add recurring item
        </Button>
      </div>

      <div className="mb-4 grid grid-cols-2 gap-4 max-md:grid-cols-1">
        <Stat label="Income" expected={expected('income')} logged={summary?.current.income} />
        <Stat label="Fixed costs" expected={expected('fixed')} logged={summary?.current.fixed} />
      </div>

      <Card className="p-0 max-md:p-0">
        {isPending && <p className="text-muted p-6 text-sm">Loading…</p>}
        {error && <p className="text-negative p-6 text-sm">Could not load recurring items.</p>}
        {items && items.length === 0 && (
          <p className="text-muted p-6 text-sm">
            No recurring items yet. Add your salary and fixed charges to track what is expected each
            month.
          </p>
        )}
        <ul className="divide-line divide-y">
          {items?.map((item) => (
            <li
              key={item.id}
              className={cn(
                'flex items-center gap-4 p-4 max-md:flex-wrap',
                !item.active && 'opacity-50',
              )}
            >
              <button
                type="button"
                className="min-w-0 flex-1 text-left"
                onClick={() => setEditingId(item.id)}
              >
                <span className="block truncate text-sm font-medium">{item.description}</span>
                <span className="text-muted block text-xs">
                  {(item.type_id && typeLabels[typeName(item) ?? 'spend']) || 'No type'} · day{' '}
                  {item.day_of_month}
                </span>
              </button>
              <span className="text-sm tabular-nums">
                {item.amount !== null ? formatMoney(item.amount) : '—'}
              </span>
              <span className="w-24 text-right">
                {!item.active ? (
                  <span className="text-muted text-xs">Not active</span>
                ) : item.logged_transaction_id ? (
                  <span className="text-muted text-xs">Logged</span>
                ) : (
                  <Button
                    variant="secondary"
                    className="h-8 px-3"
                    disabled={log.isPending || item.amount === null}
                    onClick={() => log.mutate({ id: item.id, month })}
                  >
                    Log it
                  </Button>
                )}
              </span>
            </li>
          ))}
        </ul>
      </Card>
      {log.error && <p className="text-negative mt-2 text-sm">{errorMessage(log.error)}</p>}

      <RecurringDrawer
        open={adding || editing !== null}
        item={editing}
        defaultMonth={month}
        onClose={() => {
          setAdding(false)
          setEditingId(null)
        }}
      />
    </>
  )
}

function Stat({ label, expected, logged }: { label: string; expected: number; logged?: number }) {
  return (
    <Card>
      <p className="text-muted text-xs">{label}</p>
      <p className="mt-1 text-lg font-semibold tabular-nums">{formatMoney(expected)} expected</p>
      <p className="text-muted text-sm tabular-nums">
        {logged === undefined ? '—' : formatMoney(logged)} logged this month
      </p>
    </Card>
  )
}
