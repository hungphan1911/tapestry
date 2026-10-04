import { Check, Pencil, Plus, Trash2, X } from 'lucide-react'
import { useState } from 'react'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { errorMessage } from '@/lib/errorMessage'

type NamedListEditorProps = {
  items: { id: number; name: string }[]
  noun: string
  onSave: (name: string, id?: number) => Promise<unknown>
  onDelete: (id: number) => Promise<unknown>
}

export function NamedListEditor({ items, noun, onSave, onDelete }: NamedListEditorProps) {
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editName, setEditName] = useState('')
  const [newName, setNewName] = useState('')
  const [error, setError] = useState<string | null>(null)

  const run = async (action: () => Promise<unknown>, done?: () => void) => {
    setError(null)
    try {
      await action()
      done?.()
    } catch (err) {
      setError(errorMessage(err))
    }
  }

  const remove = (id: number, name: string) => {
    if (window.confirm(`Delete "${name}"? Transactions using it will become uncategorised.`)) {
      void run(() => onDelete(id))
    }
  }

  return (
    <div>
      <ul className="divide-line divide-y">
        {items.length === 0 && <li className="text-muted py-2 text-sm">No {noun}s yet.</li>}
        {items.map((item) => (
          <li key={item.id} className="flex items-center gap-2 py-2">
            {editingId === item.id ? (
              <>
                <Input
                  autoFocus
                  value={editName}
                  aria-label={`${noun} name`}
                  onChange={(e) => setEditName(e.target.value)}
                />
                <Button
                  variant="ghost"
                  className="px-2"
                  aria-label="Save"
                  onClick={() =>
                    run(
                      () => onSave(editName, item.id),
                      () => setEditingId(null),
                    )
                  }
                >
                  <Check className="size-4" aria-hidden />
                </Button>
                <Button
                  variant="ghost"
                  className="px-2"
                  aria-label="Cancel"
                  onClick={() => setEditingId(null)}
                >
                  <X className="size-4" aria-hidden />
                </Button>
              </>
            ) : (
              <>
                <span className="flex-1 text-sm">{item.name}</span>
                <Button
                  variant="ghost"
                  className="px-2"
                  aria-label={`Rename ${item.name}`}
                  onClick={() => {
                    setEditingId(item.id)
                    setEditName(item.name)
                  }}
                >
                  <Pencil className="size-4" aria-hidden />
                </Button>
                <Button
                  variant="ghost"
                  className="px-2"
                  aria-label={`Delete ${item.name}`}
                  onClick={() => remove(item.id, item.name)}
                >
                  <Trash2 className="size-4" aria-hidden />
                </Button>
              </>
            )}
          </li>
        ))}
      </ul>

      <form
        className="mt-3 flex gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          if (newName.trim())
            void run(
              () => onSave(newName),
              () => setNewName(''),
            )
        }}
      >
        <Input
          value={newName}
          placeholder={`New ${noun}`}
          aria-label={`New ${noun}`}
          onChange={(e) => setNewName(e.target.value)}
        />
        <Button type="submit" variant="secondary" aria-label={`Add ${noun}`}>
          <Plus className="size-4" aria-hidden />
          Add
        </Button>
      </form>
      {error && <p className="text-negative mt-2 text-sm">{error}</p>}
    </div>
  )
}
