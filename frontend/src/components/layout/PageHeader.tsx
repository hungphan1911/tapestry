import type { ReactNode } from 'react'

type PageHeaderProps = { title: string; description?: string; actions?: ReactNode }

export function PageHeader({ title, description, actions }: PageHeaderProps) {
  return (
    <header className="mb-8 flex items-end justify-between gap-4 max-md:mb-6 max-md:flex-col max-md:items-start">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight max-md:text-xl">{title}</h1>
        {description && <p className="text-muted mt-1 text-sm">{description}</p>}
      </div>
      {actions}
    </header>
  )
}
