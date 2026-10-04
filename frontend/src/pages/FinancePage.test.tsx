import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { FinancePage } from './FinancePage'

const typeRows = [
  { id: 1, description: 'spend' },
  { id: 2, description: 'fixed' },
]

beforeEach(() => {
  vi.stubGlobal(
    'fetch',
    vi.fn(async (url: string) => {
      const body = url.includes('/types')
        ? typeRows
        : url.includes('/transactions/')
          ? [
              {
                id: 1,
                date: '2026-09-01T00:00:00Z',
                description: 'Tiền nhà',
                amount: 1611.24,
                amount_expression: '1611.24',
                type_id: 2,
                category_id: null,
                card_id: null,
                recurring_id: null,
              },
            ]
          : url.includes('/budgets/')
            ? { month: '2026-09', total: null, categories: [] }
            : []
      return new Response(JSON.stringify(body), { status: 200 })
    }),
  )
})

afterEach(() => vi.unstubAllGlobals())

function renderPage(path: string) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter initialEntries={[path]}>
        <FinancePage />
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

describe('FinancePage', () => {
  it('shows the requested tab and month from the URL', async () => {
    renderPage('/finance?tab=transactions&month=2026-09')
    expect(screen.getByRole('tab', { name: 'Transactions' })).toHaveAttribute(
      'aria-selected',
      'true',
    )
    expect(screen.getByText('September 2026')).toBeInTheDocument()
    expect((await screen.findAllByText('Tiền nhà')).length).toBeGreaterThan(0)
    expect((await screen.findAllByText('$1,611.24')).length).toBeGreaterThan(0)
  })

  it('switches tabs', async () => {
    renderPage('/finance?tab=settings')
    expect(screen.getByRole('heading', { name: /categories/i })).toBeInTheDocument()
    await userEvent.click(screen.getByRole('tab', { name: 'Budget' }))
    expect(await screen.findByText(/total monthly budget/i)).toBeInTheDocument()
  })
})
