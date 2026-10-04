import { NavLink } from 'react-router-dom'
import { cn } from '@/lib/cn'
import { navItems } from './navItems'

// Phone only (hidden at md and up via the max-md variant on display).
export function BottomNav() {
  return (
    <nav
      aria-label="Primary"
      className="border-line bg-surface fixed inset-x-0 bottom-0 z-10 hidden border-t pb-[env(safe-area-inset-bottom)] max-md:flex"
    >
      {navItems.map(({ to, label, icon: Icon, end }) => (
        <NavLink
          key={to}
          to={to}
          end={end}
          className={({ isActive }) =>
            cn(
              'text-muted flex flex-1 flex-col items-center gap-1 py-2 text-xs font-medium',
              isActive && 'text-ink',
            )
          }
        >
          <Icon className="size-5" aria-hidden />
          {label}
        </NavLink>
      ))}
    </nav>
  )
}
