import { LayoutDashboard, Wallet } from 'lucide-react'
import { routes } from '@/config/routes'

export const navItems = [
  { to: routes.home, label: 'Home', icon: LayoutDashboard, end: true },
  { to: routes.finance, label: 'Finance', icon: Wallet, end: false },
]
