import { Card } from '@/components/ui/Card'
import { CardManager, CategoryManager } from './Managers'

export function SettingsTab() {
  return (
    <div className="grid grid-cols-2 gap-6 max-lg:grid-cols-1">
      <Card>
        <h2 className="mb-3 text-base font-semibold">Categories</h2>
        <CategoryManager />
      </Card>
      <Card>
        <h2 className="mb-3 text-base font-semibold">Cards</h2>
        <CardManager />
      </Card>
    </div>
  )
}
