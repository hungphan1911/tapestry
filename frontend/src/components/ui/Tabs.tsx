import { cn } from '@/lib/cn'

type TabsProps<T extends string> = {
  items: { id: T; label: string }[]
  value: T
  onChange: (id: T) => void
}

export function Tabs<T extends string>({ items, value, onChange }: TabsProps<T>) {
  return (
    <div role="tablist" className="border-line mb-6 flex gap-1 overflow-x-auto border-b">
      {items.map((item) => (
        <button
          key={item.id}
          type="button"
          role="tab"
          aria-selected={item.id === value}
          onClick={() => onChange(item.id)}
          className={cn(
            'shrink-0 border-b-2 px-4 py-2 text-sm font-medium transition-colors',
            item.id === value
              ? 'border-primary text-ink'
              : 'text-muted hover:text-ink border-transparent',
          )}
        >
          {item.label}
        </button>
      ))}
    </div>
  )
}
