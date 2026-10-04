import { X } from 'lucide-react'
import { useEffect, type ReactNode } from 'react'

type DrawerProps = { open: boolean; title: string; onClose: () => void; children: ReactNode }

// Side panel on laptop, full screen on phone. A drawer rendered later in the tree stacks above an earlier one.
export function Drawer({ open, title, onClose, children }: DrawerProps) {
  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose()
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, onClose])

  if (!open) return null
  return (
    <div
      className="fixed inset-0 z-50 flex justify-end"
      role="dialog"
      aria-modal
      aria-label={title}
    >
      <button
        type="button"
        aria-label="Close"
        className="bg-ink/30 absolute inset-0"
        onClick={onClose}
      />
      <div className="bg-surface relative flex h-full w-full max-w-md flex-col shadow-xl max-md:max-w-none">
        <div className="border-line flex items-center justify-between border-b px-5 py-4">
          <h2 className="text-base font-semibold">{title}</h2>
          <button
            type="button"
            aria-label="Close"
            onClick={onClose}
            className="hover:bg-subtle rounded-lg p-1"
          >
            <X className="size-5" aria-hidden />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto p-5">{children}</div>
      </div>
    </div>
  )
}
