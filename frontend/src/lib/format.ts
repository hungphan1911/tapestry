const usd = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })

export function formatMoney(amount: number): string {
  return usd.format(amount)
}

// Months are "YYYY-MM" strings and dates are UTC calendar dates, matching the backend.

export function currentMonth(): string {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
}

export function todayInput(): string {
  const now = new Date()
  return `${currentMonth()}-${String(now.getDate()).padStart(2, '0')}`
}

export function shiftMonth(month: string, delta: number): string {
  const [year, m] = month.split('-').map(Number)
  const d = new Date(Date.UTC(year, m - 1 + delta, 1))
  return `${d.getUTCFullYear()}-${String(d.getUTCMonth() + 1).padStart(2, '0')}`
}

export function monthLabel(month: string): string {
  const [year, m] = month.split('-').map(Number)
  return new Intl.DateTimeFormat('en-US', {
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(Date.UTC(year, m - 1, 1)))
}

export function daysInMonth(month: string): number {
  const [year, m] = month.split('-').map(Number)
  return new Date(Date.UTC(year, m, 0)).getUTCDate()
}

export function monthBounds(month: string): { from: string; to: string } {
  return { from: `${month}-01T00:00:00Z`, to: `${month}-${daysInMonth(month)}T23:59:59Z` }
}

// "2026-09-01T00:00:00Z" -> "01/09/2026"
export function formatDate(iso: string): string {
  const [year, month, day] = iso.slice(0, 10).split('-')
  return `${day}/${month}/${year}`
}

export function toDateInput(iso: string): string {
  return iso.slice(0, 10)
}

export function toApiDate(input: string): string {
  return `${input}T00:00:00Z`
}
