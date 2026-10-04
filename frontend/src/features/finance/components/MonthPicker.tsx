import { ChevronLeft, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/Button'
import { monthLabel, shiftMonth } from '@/lib/format'

type MonthPickerProps = { month: string; onChange: (month: string) => void }

export function MonthPicker({ month, onChange }: MonthPickerProps) {
  return (
    <div className="flex items-center gap-1">
      <Button
        variant="secondary"
        className="w-9 px-0"
        aria-label="Previous month"
        onClick={() => onChange(shiftMonth(month, -1))}
      >
        <ChevronLeft className="size-4" aria-hidden />
      </Button>
      <span className="min-w-36 text-center text-sm font-medium">{monthLabel(month)}</span>
      <Button
        variant="secondary"
        className="w-9 px-0"
        aria-label="Next month"
        onClick={() => onChange(shiftMonth(month, 1))}
      >
        <ChevronRight className="size-4" aria-hidden />
      </Button>
    </div>
  )
}
