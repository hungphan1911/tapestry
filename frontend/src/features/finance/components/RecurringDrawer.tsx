import { Trash2 } from 'lucide-react'
import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/Button'
import { Drawer } from '@/components/ui/Drawer'
import { Field, Input, Select } from '@/components/ui/Input'
import { errorMessage } from '@/lib/errorMessage'
import { formatMoney } from '@/lib/format'
import { typeLabels } from '../labels'
import {
  useCards,
  useCategories,
  useDeleteRecurring,
  useDeleteRecurringAmount,
  useSaveRecurring,
  useSetRecurringAmount,
  useTypes,
} from '../queries'
import type { Recurring } from '../types'
import { ManageDrawer } from './Managers'

type RecurringDrawerProps = {
  open: boolean
  item: Recurring | null
  defaultMonth: string
  onClose: () => void
}

export function RecurringDrawer({ open, item, defaultMonth, onClose }: RecurringDrawerProps) {
  return (
    <Drawer
      open={open}
      title={item ? 'Edit recurring item' : 'Add recurring item'}
      onClose={onClose}
    >
      <RecurringForm item={item} defaultMonth={defaultMonth} onDone={onClose} />
    </Drawer>
  )
}

const toId = (v: string) => (v === '' ? null : Number(v))

type FormProps = { item: Recurring | null; defaultMonth: string; onDone: () => void }

function RecurringForm({ item, defaultMonth, onDone }: FormProps) {
  const { data: types = [] } = useTypes()
  const { data: categories = [] } = useCategories()
  const { data: cards = [] } = useCards()
  const save = useSaveRecurring()
  const remove = useDeleteRecurring()

  const [description, setDescription] = useState(item?.description ?? '')
  const [typeId, setTypeId] = useState(item?.type_id?.toString() ?? '')
  const [categoryId, setCategoryId] = useState(item?.category_id?.toString() ?? '')
  const [cardId, setCardId] = useState(item?.card_id?.toString() ?? '')
  const [day, setDay] = useState(String(item?.day_of_month ?? 1))
  const [startMonth, setStartMonth] = useState(item?.start_month ?? defaultMonth)
  const [endMonth, setEndMonth] = useState(item?.end_month ?? '')
  const [expression, setExpression] = useState('')
  const [managing, setManaging] = useState<'category' | 'card' | null>(null)

  const submit = (e: FormEvent) => {
    e.preventDefault()
    save.mutate(
      {
        id: item?.id,
        input: {
          description,
          type_id: toId(typeId),
          category_id: toId(categoryId),
          card_id: toId(cardId),
          day_of_month: Number(day),
          start_month: startMonth,
          end_month: endMonth || null,
          amount_expression: expression,
        },
      },
      { onSuccess: onDone },
    )
  }

  const manageLink = (kind: 'category' | 'card') => (
    <button
      type="button"
      className="text-muted hover:text-ink text-xs underline"
      onClick={() => setManaging(kind)}
    >
      Manage
    </button>
  )

  return (
    <div className="space-y-6">
      <form onSubmit={submit} className="space-y-4">
        <Field label="Description">
          <Input required value={description} onChange={(e) => setDescription(e.target.value)} />
        </Field>
        <Field label="Type">
          <Select value={typeId} onChange={(e) => setTypeId(e.target.value)}>
            <option value="">None</option>
            {types.map((t) => (
              <option key={t.id} value={t.id}>
                {typeLabels[t.description] ?? t.description}
              </option>
            ))}
          </Select>
        </Field>
        <Field label="Category" action={manageLink('category')}>
          <Select value={categoryId} onChange={(e) => setCategoryId(e.target.value)}>
            <option value="">None</option>
            {categories.map((c) => (
              <option key={c.id} value={c.id}>
                {c.description}
              </option>
            ))}
          </Select>
        </Field>
        <Field label="Card" action={manageLink('card')}>
          <Select value={cardId} onChange={(e) => setCardId(e.target.value)}>
            <option value="">None</option>
            {cards.map((c) => (
              <option key={c.id} value={c.id}>
                {c.card_name}
              </option>
            ))}
          </Select>
        </Field>
        <div className="grid grid-cols-2 gap-3">
          <Field label="Day of month" hint="Short months use their last day.">
            <Input
              type="number"
              min="1"
              max="31"
              required
              value={day}
              onChange={(e) => setDay(e.target.value)}
            />
          </Field>
          {!item && (
            <Field label="Amount" hint="Sums work, e.g. 1500+111.24">
              <Input
                required
                inputMode="decimal"
                className="tabular-nums"
                value={expression}
                onChange={(e) => setExpression(e.target.value)}
              />
            </Field>
          )}
        </div>
        <div className="grid grid-cols-2 gap-3">
          <Field label="Starts">
            <Input
              type="month"
              required
              value={startMonth}
              onChange={(e) => setStartMonth(e.target.value)}
            />
          </Field>
          <Field label="Ends (optional)">
            <Input type="month" value={endMonth} onChange={(e) => setEndMonth(e.target.value)} />
          </Field>
        </div>

        {(save.error ?? remove.error) && (
          <p className="text-negative text-sm">{errorMessage(save.error ?? remove.error)}</p>
        )}

        <div className="flex items-center justify-between pt-2">
          <Button type="submit" disabled={save.isPending}>
            {item ? 'Save changes' : 'Add item'}
          </Button>
          {item && (
            <Button
              variant="ghost"
              onClick={() => {
                if (window.confirm(`Delete "${item.description}"? Logged transactions are kept.`)) {
                  remove.mutate(item.id, { onSuccess: onDone })
                }
              }}
            >
              Delete
            </Button>
          )}
        </div>
      </form>

      {item && <AmountHistory item={item} defaultMonth={defaultMonth} />}
      <ManageDrawer kind={managing} onClose={() => setManaging(null)} />
    </div>
  )
}

function AmountHistory({ item, defaultMonth }: { item: Recurring; defaultMonth: string }) {
  const setAmount = useSetRecurringAmount()
  const removeAmount = useDeleteRecurringAmount()
  const [month, setMonth] = useState(defaultMonth)
  const [expression, setExpression] = useState('')

  const submit = (e: FormEvent) => {
    e.preventDefault()
    setAmount.mutate({ id: item.id, month, expression }, { onSuccess: () => setExpression('') })
  }

  const error = setAmount.error ?? removeAmount.error

  return (
    <section className="border-line border-t pt-4">
      <h3 className="mb-2 text-sm font-semibold">Amount history</h3>
      <ul className="divide-line divide-y">
        {item.amounts.map((a) => (
          <li key={a.effective_month} className="flex items-center justify-between py-2 text-sm">
            <span>From {a.effective_month}</span>
            <span className="flex items-center gap-2">
              <span className="tabular-nums">{formatMoney(a.amount)}</span>
              {item.amounts.length > 1 && (
                <button
                  type="button"
                  aria-label={`Remove change from ${a.effective_month}`}
                  className="hover:bg-subtle rounded p-1"
                  onClick={() => removeAmount.mutate({ id: item.id, month: a.effective_month })}
                >
                  <Trash2 className="size-4" aria-hidden />
                </button>
              )}
            </span>
          </li>
        ))}
      </ul>

      <form onSubmit={submit} className="mt-3 flex items-end gap-2">
        <Field label="Change from">
          <Input type="month" required value={month} onChange={(e) => setMonth(e.target.value)} />
        </Field>
        <Field label="New amount">
          <Input
            required
            inputMode="decimal"
            className="tabular-nums"
            value={expression}
            onChange={(e) => setExpression(e.target.value)}
          />
        </Field>
        <Button type="submit" variant="secondary" disabled={setAmount.isPending}>
          Set
        </Button>
      </form>
      {error && <p className="text-negative mt-2 text-sm">{errorMessage(error)}</p>}
    </section>
  )
}
