import {
  daysInMonth,
  formatDate,
  formatMoney,
  monthBounds,
  monthLabel,
  shiftMonth,
  toApiDate,
} from './format'

describe('format helpers', () => {
  it('formats money as USD', () => {
    expect(formatMoney(1611.2)).toBe('$1,611.20')
    expect(formatMoney(-42.14)).toBe('-$42.14')
  })

  it('shifts months across year boundaries', () => {
    expect(shiftMonth('2026-01', -1)).toBe('2025-12')
    expect(shiftMonth('2026-12', 1)).toBe('2027-01')
    expect(shiftMonth('2026-09', 0)).toBe('2026-09')
  })

  it('labels months and counts days', () => {
    expect(monthLabel('2026-09')).toBe('September 2026')
    expect(daysInMonth('2026-02')).toBe(28)
    expect(daysInMonth('2028-02')).toBe(29)
  })

  it('builds month bounds and converts dates', () => {
    expect(monthBounds('2026-09')).toEqual({
      from: '2026-09-01T00:00:00Z',
      to: '2026-09-30T23:59:59Z',
    })
    expect(formatDate('2026-09-01T00:00:00Z')).toBe('01/09/2026')
    expect(toApiDate('2026-09-01')).toBe('2026-09-01T00:00:00Z')
  })
})
