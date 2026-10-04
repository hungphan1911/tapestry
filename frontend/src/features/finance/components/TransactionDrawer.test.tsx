import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TransactionDrawer } from './TransactionDrawer'

const calls: { url: string; method: string }[] = []

beforeEach(() => {
  calls.length = 0
  vi.stubGlobal(
    'fetch',
    vi.fn(async (url: string, init?: RequestInit) => {
      calls.push({ url, method: init?.method ?? 'GET' })
      const created = init?.method === 'POST' && url.includes('/categories')
      return new Response(JSON.stringify(created ? { id: 1, description: 'Pets' } : []), {
        status: created ? 201 : 200,
      })
    }),
  )
})

afterEach(() => vi.unstubAllGlobals())

describe('TransactionDrawer', () => {
  it('adds a category from the manage drawer without saving the transaction', async () => {
    const onClose = vi.fn()
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
    render(
      <QueryClientProvider client={client}>
        <TransactionDrawer open transaction={null} defaultDate="2026-10-03" onClose={onClose} />
      </QueryClientProvider>,
    )

    await userEvent.click(screen.getAllByText('Manage')[0])
    await userEvent.type(screen.getByLabelText('New category'), 'Pets')
    await userEvent.click(screen.getByRole('button', { name: 'Add category' }))

    expect(calls.some((c) => c.method === 'POST' && c.url.includes('/categories'))).toBe(true)
    expect(calls.some((c) => c.method === 'POST' && c.url.includes('/transactions'))).toBe(false)
    expect(onClose).not.toHaveBeenCalled()
  })
})
