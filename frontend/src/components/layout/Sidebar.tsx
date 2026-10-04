import { NavLink } from 'react-router-dom'
import { cn } from '@/lib/cn'
import { navItems } from './navItems'

// Laptop: icon + label. Tablet (max-lg): icon rail. Phone (max-md): hidden, BottomNav takes over.
export function Sidebar() {
  return (
    <aside className="border-line bg-surface sticky top-0 flex h-dvh w-60 shrink-0 flex-col gap-6 border-r px-3 py-5 max-lg:w-16 max-lg:items-center max-lg:px-2 max-md:hidden">
      <div className="px-3 text-lg font-semibold tracking-tight max-lg:px-0 max-lg:text-base">
        <span className="max-lg:hidden">Tapestry</span>
        <span className="hidden max-lg:inline">T</span>
      </div>
      <nav aria-label="Primary" className="flex flex-col gap-1">
        {navItems.map(({ to, label, icon: Icon, end }) => (
          <NavLink
            key={to}
            to={to}
            end={end}
            title={label}
            className={({ isActive }) =>
              cn(
                'text-muted hover:bg-subtle hover:text-ink flex h-10 items-center gap-3 rounded-lg px-3 text-sm font-medium transition-colors max-lg:w-10 max-lg:justify-center max-lg:px-0',
                isActive && 'bg-subtle text-ink',
              )
            }
          >
            <Icon className="size-5 shrink-0" aria-hidden />
            <span className="max-lg:sr-only">{label}</span>
          </NavLink>
        ))}
      </nav>
    </aside>
  )
}
