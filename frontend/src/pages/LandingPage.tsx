import { Wallet } from 'lucide-react'
import { PageHeader } from '@/components/layout/PageHeader'
import { routes } from '@/config/routes'
import { ModuleCard } from '@/features/landing/ModuleCard'

const modules = [
  {
    to: routes.finance,
    title: 'Personal finance',
    description: 'Track transactions, categories and cards.',
    icon: Wallet,
  },
]

export function LandingPage() {
  return (
    <>
      <PageHeader title="Tapestry" description="Personal home lab." />
      <div className="grid grid-cols-3 gap-4 max-lg:grid-cols-2 max-md:grid-cols-1">
        {modules.map((m) => (
          <ModuleCard key={m.to} {...m} />
        ))}
      </div>
    </>
  )
}
