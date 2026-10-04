import { Outlet } from 'react-router-dom'
import { BottomNav } from './BottomNav'
import { Sidebar } from './Sidebar'

export function AppShell() {
  return (
    <div className="flex min-h-dvh">
      <Sidebar />
      <main className="min-w-0 flex-1 px-10 py-8 max-lg:px-6 max-md:px-4 max-md:pb-24">
        <Outlet />
      </main>
      <BottomNav />
    </div>
  )
}
