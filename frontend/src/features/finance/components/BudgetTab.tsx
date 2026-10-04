import { useState } from 'react'
import { Button } from '@/components/ui/Button'
import { Card } from '@/components/ui/Card'
import { Field, Input } from '@/components/ui/Input'
import { errorMessage } from '@/lib/errorMessage'
import { formatMoney, shiftMonth } from '@/lib/format'
import { useBudget, useCategories, useSetBudget } from '../queries'
import type { Budget, Category } from '../types'
import { MonthPicker } from './MonthPicker'

type BudgetTabProps = { month: string; onMonthChange: (month: string) => void }

export function BudgetTab({ month, onMonthChange }: BudgetTabProps) {
  const { data: budget, isPending, error } = useBudget(month)
  const { data: previous } = useBudget(shiftMonth(month, -1))
  const { data: categories = [] } = useCategories()

  return (
    <>
      <div className="mb-4">
        <MonthPicker month={month} onChange={onMonthChange} />
      </div>
      {isPending && <p className="text-muted text-sm">Loading…</p>}
      {error && <p className="text-negative text-sm">Could not load the budget.</p>}
      {budget && (
        <BudgetForm
          key={month}
          month={month}
          budget={budget}
          previous={previous}
          categories={categories}
        />
      )}
    </>
  )
}

type FormProps = { month: string; budget: Budget; previous?: Budget; categories: Category[] }

const toText = (n: number | null | undefined) => (n === null || n === undefined ? '' : String(n))

function initialAmounts(budget: Budget) {
  return Object.fromEntries(budget.categories.map((c) => [c.category_id, String(c.amount)]))
}

function BudgetForm({ month, budget, previous, categories }: FormProps) {
  const save = useSetBudget(month)
  const [total, setTotal] = useState(toText(budget.total))
  const [amounts, setAmounts] = useState<Record<number, string>>(initialAmounts(budget))
  const [saved, setSaved] = useState(false)

  const categorySum = Object.values(amounts).reduce((sum, v) => sum + (Number(v) || 0), 0)

  const copyPrevious = () => {
    if (!previous) return
    setTotal(toText(previous.total))
    setAmounts(initialAmounts(previous))
    setSaved(false)
  }

  const submit = () => {
    setSaved(false)
    save.mutate(
      {
        total: total === '' ? null : Number(total),
        categories: Object.entries(amounts)
          .filter(([, v]) => v !== '')
          .map(([id, v]) => ({ category_id: Number(id), amount: Number(v) })),
      },
      { onSuccess: () => setSaved(true) },
    )
  }

  const change = (setter: () => void) => {
    setter()
    setSaved(false)
  }

  return (
    <Card className="max-w-xl">
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault()
          submit()
        }}
      >
        <Field
          label="Total monthly budget"
          hint="Compared against fixed + spend, minus non-salary."
        >
          <Input
            type="number"
            min="0"
            step="0.01"
            className="tabular-nums"
            value={total}
            onChange={(e) => change(() => setTotal(e.target.value))}
          />
        </Field>

        <div>
          <h2 className="mb-2 text-sm font-semibold">Per category</h2>
          {categories.length === 0 && (
            <p className="text-muted text-sm">Add categories in Settings to budget them.</p>
          )}
          <div className="space-y-2">
            {categories.map((c) => (
              <div key={c.id} className="flex items-center gap-3">
                <span className="flex-1 text-sm">{c.description}</span>
                <Input
                  type="number"
                  min="0"
                  step="0.01"
                  aria-label={`${c.description} budget`}
                  placeholder="No budget"
                  className="w-36 tabular-nums"
                  value={amounts[c.id] ?? ''}
                  onChange={(e) => change(() => setAmounts({ ...amounts, [c.id]: e.target.value }))}
                />
              </div>
            ))}
          </div>
          <p className="text-muted mt-2 text-xs tabular-nums">
            Categories add up to {formatMoney(categorySum)}
            {total !== '' && ` of ${formatMoney(Number(total))}`}.
          </p>
        </div>

        {save.error && <p className="text-negative text-sm">{errorMessage(save.error)}</p>}

        <div className="flex items-center gap-3">
          <Button type="submit" disabled={save.isPending}>
            Save budget
          </Button>
          <Button variant="secondary" disabled={!previous} onClick={copyPrevious}>
            Copy last month
          </Button>
          {saved && <span className="text-muted text-sm">Saved.</span>}
        </div>
      </form>
    </Card>
  )
}
