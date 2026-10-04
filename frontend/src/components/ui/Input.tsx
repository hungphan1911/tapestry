import type { InputHTMLAttributes, ReactNode, SelectHTMLAttributes } from 'react'
import { cn } from '@/lib/cn'

const control =
  'h-9 w-full rounded-lg border border-line bg-surface px-3 text-sm text-ink placeholder:text-muted focus-visible:outline-primary focus-visible:outline-2 focus-visible:outline-offset-1 disabled:opacity-50'

export function Input({ className, ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return <input className={cn(control, className)} {...props} />
}

export function Select({ className, children, ...props }: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select className={cn(control, className)} {...props}>
      {children}
    </select>
  )
}

type FieldProps = { label: string; hint?: string; action?: ReactNode; children: ReactNode }

export function Field({ label, hint, action, children }: FieldProps) {
  return (
    <label className="block">
      <span className="mb-1 flex items-center justify-between text-xs font-medium">
        <span>{label}</span>
        {action}
      </span>
      {children}
      {hint && <span className="text-muted mt-1 block text-xs">{hint}</span>}
    </label>
  )
}
