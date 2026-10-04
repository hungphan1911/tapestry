import type { LucideIcon } from 'lucide-react'
import { ArrowRight } from 'lucide-react'
import { Link } from 'react-router-dom'
import { Card } from '@/components/ui/Card'

type ModuleCardProps = { to: string; title: string; description: string; icon: LucideIcon }

export function ModuleCard({ to, title, description, icon: Icon }: ModuleCardProps) {
  return (
    <Link
      to={to}
      className="group rounded-card focus-visible:outline-primary block focus-visible:outline-2 focus-visible:outline-offset-2"
    >
      <Card className="group-hover:border-primary/40 h-full transition-colors">
        <div className="bg-subtle mb-4 flex size-10 items-center justify-center rounded-lg">
          <Icon className="size-5" aria-hidden />
        </div>
        <h2 className="text-base font-semibold">{title}</h2>
        <p className="text-muted mt-1 text-sm">{description}</p>
        <span className="mt-4 inline-flex items-center gap-1 text-sm font-medium">
          Open{' '}
          <ArrowRight
            className="size-4 transition-transform group-hover:translate-x-0.5"
            aria-hidden
          />
        </span>
      </Card>
    </Link>
  )
}
