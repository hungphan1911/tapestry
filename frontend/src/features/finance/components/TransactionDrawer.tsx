import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/Button'
import { Drawer } from '@/components/ui/Drawer'
import { Field, Input, Select } from '@/components/ui/Input'
import { errorMessage } from '@/lib/errorMessage'
import { toApiDate, toDateInput } from '@/lib/format'
import { typeLabels } from '../labels'
import {
  useCards,
  useCategories,
  useDeleteTransaction,
  useSaveTransaction,
  useTypes,
} from '../queries'
import type { Transaction } from '../types'
import { ManageDrawer } from './Managers'

type TransactionDrawerProps = {
  open: boolean
  transaction: Transaction | null
  defaultDate: string
  onClose: () => void
}

export function TransactionDrawer({
  open,
  transaction,
  defaultDate,
  onClose,
}: TransactionDrawerProps) {
  return (
    <Drawer
      open={open}
      title={transaction ? 'Edit transaction' : 'Add transaction'}
      onClose={onClose}
    >
      <TransactionForm transaction={transaction} defaultDate={defaultDate} onDone={onClose} />
    </Drawer>
  )
}

const toId = (v: string) => (v === '' ? null : Number(v))

type FormProps = { transaction: Transaction | null; defaultDate: string; onDone: () => void }

function TransactionForm({ transaction, defaultDate, onDone }: FormProps) {
  const { data: types = [] } = useTypes()
  const { data: categories = [] } = useCategories()
  const { data: cards = [] } = useCards()
  const save = useSaveTransaction()
  const remove = useDeleteTransaction()

  const [date, setDate] = useState(transaction ? toDateInput(transaction.date) : defaultDate)
  const [description, setDescription] = useState(transaction?.description ?? '')
  const [expression, setExpression] = useState(
    transaction ? transaction.amount_expression || String(transaction.amount) : '',
  )
  const [typeId, setTypeId] = useState(transaction?.type_id?.toString() ?? '')
  const [categoryId, setCategoryId] = useState(transaction?.category_id?.toString() ?? '')
  const [cardId, setCardId] = useState(transaction?.card_id?.toString() ?? '')
  const [managing, setManaging] = useState<'category' | 'card' | null>(null)

  // New transactions default to "spend", the most common type.
  const effectiveTypeId =
    typeId ||
    (transaction ? '' : (types.find((t) => t.description === 'spend')?.id.toString() ?? ''))

  const submit = (e: FormEvent) => {
    e.preventDefault()
    save.mutate(
      {
        id: transaction?.id,
        input: {
          date: toApiDate(date),
          description,
          amount_expression: expression,
          type_id: toId(effectiveTypeId),
          category_id: toId(categoryId),
          card_id: toId(cardId),
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

  const error = save.error ?? remove.error

  return (
    <form onSubmit={submit} className="space-y-4">
      <Field label="Date">
        <Input type="date" required value={date} onChange={(e) => setDate(e.target.value)} />
      </Field>
      <Field label="Description">
        <Input value={description} onChange={(e) => setDescription(e.target.value)} />
      </Field>
      <Field label="Amount" hint="Sums work too, e.g. 50+6.7. The formula is kept for editing.">
        <Input
          required
          inputMode="decimal"
          className="tabular-nums"
          value={expression}
          onChange={(e) => setExpression(e.target.value)}
        />
      </Field>
      <Field label="Type">
        <Select value={effectiveTypeId} onChange={(e) => setTypeId(e.target.value)}>
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

      {error && <p className="text-negative text-sm">{errorMessage(error)}</p>}

      <div className="flex items-center justify-between pt-2">
        <Button type="submit" disabled={save.isPending}>
          {transaction ? 'Save changes' : 'Add transaction'}
        </Button>
        {transaction && (
          <Button
            variant="ghost"
            onClick={() => {
              if (window.confirm('Delete this transaction?')) {
                remove.mutate(transaction.id, { onSuccess: onDone })
              }
            }}
          >
            Delete
          </Button>
        )}
      </div>

      <ManageDrawer kind={managing} onClose={() => setManaging(null)} />
    </form>
  )
}
