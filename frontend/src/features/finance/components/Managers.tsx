import { Drawer } from '@/components/ui/Drawer'
import {
  useCards,
  useCategories,
  useDeleteCard,
  useDeleteCategory,
  useSaveCard,
  useSaveCategory,
} from '../queries'
import { NamedListEditor } from './NamedListEditor'

export function CategoryManager() {
  const { data = [] } = useCategories()
  const save = useSaveCategory()
  const remove = useDeleteCategory()
  return (
    <NamedListEditor
      noun="category"
      items={data.map((c) => ({ id: c.id, name: c.description }))}
      onSave={(name, id) => save.mutateAsync({ id, description: name })}
      onDelete={(id) => remove.mutateAsync(id)}
    />
  )
}

export function CardManager() {
  const { data = [] } = useCards()
  const save = useSaveCard()
  const remove = useDeleteCard()
  return (
    <NamedListEditor
      noun="card"
      items={data.map((c) => ({ id: c.id, name: c.card_name }))}
      onSave={(name, id) => save.mutateAsync({ id, card_name: name })}
      onDelete={(id) => remove.mutateAsync(id)}
    />
  )
}

type ManageDrawerProps = { kind: 'category' | 'card' | null; onClose: () => void }

// Opened from the transaction and recurring forms so lists can be edited without leaving them.
export function ManageDrawer({ kind, onClose }: ManageDrawerProps) {
  return (
    <Drawer
      open={kind !== null}
      title={kind === 'card' ? 'Manage cards' : 'Manage categories'}
      onClose={onClose}
    >
      {kind === 'card' ? <CardManager /> : <CategoryManager />}
    </Drawer>
  )
}
