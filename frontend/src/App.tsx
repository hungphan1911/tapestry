import { QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { queryClient } from '@/api/queryClient'
import { AppShell } from '@/components/layout/AppShell'
import { routes } from '@/config/routes'
import { FinancePage } from '@/pages/FinancePage'
import { LandingPage } from '@/pages/LandingPage'

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route element={<AppShell />}>
            <Route path={routes.home} element={<LandingPage />} />
            <Route path={routes.finance} element={<FinancePage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
