import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Legend,
  Line,
  LineChart,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { Card } from '@/components/ui/Card'
import { cn } from '@/lib/cn'
import { daysInMonth, formatMoney } from '@/lib/format'
import { useCategories, useSummary } from '../queries'
import type { Summary } from '../types'
import { MonthPicker } from './MonthPicker'

type OverviewTabProps = { month: string; onMonthChange: (month: string) => void }

const palette = [1, 2, 3, 4, 5, 6].map((n) => `var(--color-chart-${n})`)
const axisProps = { stroke: 'var(--color-muted)', fontSize: 12, tickLine: false } as const
const money = (value: unknown) => formatMoney(Number(value))

export function OverviewTab({ month, onMonthChange }: OverviewTabProps) {
  const { data, isPending, error } = useSummary(month)
  const { data: categories = [] } = useCategories()
  const categoryName = (id: number | null) =>
    (id && categories.find((c) => c.id === id)?.description) || 'Uncategorised'

  return (
    <>
      <div className="mb-4">
        <MonthPicker month={month} onChange={onMonthChange} />
      </div>
      {isPending && <p className="text-muted text-sm">Loading…</p>}
      {error && <p className="text-negative text-sm">Could not load the summary.</p>}
      {data && (
        <div className="grid grid-cols-2 gap-6 max-lg:grid-cols-1">
          <SummaryCard summary={data} />
          <SpendDonut summary={data} categoryName={categoryName} />
          <PaceChart summary={data} month={month} />
          <MonthComparison summary={data} />
          <BudgetProgress summary={data} categoryName={categoryName} />
        </div>
      )}
    </>
  )
}

function Row({
  label,
  value,
  strong,
  tone,
}: {
  label: string
  value: number
  strong?: boolean
  tone?: boolean
}) {
  return (
    <div className={cn('flex justify-between py-2 text-sm', strong && 'font-semibold')}>
      <span>{label}</span>
      <span className={cn('tabular-nums', tone && (value < 0 ? 'text-negative' : 'text-positive'))}>
        {formatMoney(value)}
      </span>
    </div>
  )
}

function SummaryCard({ summary }: { summary: Summary }) {
  const { current, budget, remaining } = summary
  return (
    <Card>
      <h2 className="mb-2 text-base font-semibold">Monthly summary</h2>
      <div className="divide-line divide-y">
        <Row label="Fixed costs" value={current.fixed} />
        <Row label="Other spending" value={current.spend} />
        <Row label="Refunds and cashback" value={current.non_salary} />
        <Row label="Monthly budget" value={budget ?? 0} strong />
        <Row label="Total spent" value={current.gross} strong />
        <Row label="Total spent minus cashback" value={current.net} strong />
        <Row label="Remaining" value={remaining ?? 0} strong tone={remaining !== null} />
      </div>
      {budget === null && (
        <p className="text-muted mt-2 text-xs">
          No budget set for this month. Add one in the Budget tab.
        </p>
      )}
    </Card>
  )
}

function SpendDonut({
  summary,
  categoryName,
}: {
  summary: Summary
  categoryName: (id: number | null) => string
}) {
  const slices = summary.categories
    .filter((c) => c.spend > 0)
    .map((c) => ({ name: categoryName(c.category_id), value: c.spend }))

  return (
    <Card>
      <h2 className="mb-2 text-base font-semibold">Spend by category</h2>
      {slices.length === 0 ? (
        <p className="text-muted text-sm">No spending yet this month.</p>
      ) : (
        <div className="h-72">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie data={slices} dataKey="value" nameKey="name" innerRadius="55%" outerRadius="85%">
                {slices.map((s, i) => (
                  <Cell key={s.name} fill={palette[i % palette.length]} />
                ))}
              </Pie>
              <Tooltip formatter={money} />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
      )}
    </Card>
  )
}

function PaceChart({ summary, month }: { summary: Summary; month: string }) {
  const days = daysInMonth(month)
  const byDay = new Map(summary.daily.map((d) => [Number(d.date.slice(8, 10)), d.cumulative]))
  const data = Array.from({ length: days }, (_, i) => ({
    day: i + 1,
    spent: byDay.get(i + 1),
    pace: summary.budget === null ? undefined : (summary.budget * (i + 1)) / days,
  }))

  return (
    <Card>
      <h2 className="mb-1 text-base font-semibold">Cumulative spending vs. budget</h2>
      <p className="text-muted mb-2 text-xs">
        Fixed + spend, against an even pace through the month.
      </p>
      <div className="h-72">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data}>
            <CartesianGrid stroke="var(--color-line)" vertical={false} />
            <XAxis dataKey="day" {...axisProps} />
            <YAxis {...axisProps} width={48} />
            <Tooltip formatter={money} labelFormatter={(day) => `Day ${day}`} />
            <Legend />
            <Line
              name="Spent"
              dataKey="spent"
              stroke="var(--color-chart-1)"
              strokeWidth={2}
              dot={false}
              connectNulls={false}
            />
            {summary.budget !== null && (
              <Line
                name="Budget pace"
                dataKey="pace"
                stroke="var(--color-chart-4)"
                strokeDasharray="5 4"
                strokeWidth={2}
                dot={false}
              />
            )}
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  )
}

function MonthComparison({ summary }: { summary: Summary }) {
  const { current, previous } = summary
  const data = [
    { name: 'Fixed', previous: previous.fixed, current: current.fixed },
    { name: 'Spend', previous: previous.spend, current: current.spend },
    { name: 'Net', previous: previous.net, current: current.net },
  ]
  const delta = current.net - previous.net

  return (
    <Card>
      <h2 className="mb-1 text-base font-semibold">Compared to last month</h2>
      <p className="text-muted mb-2 text-xs">
        Net spending is {formatMoney(Math.abs(delta))} {delta > 0 ? 'higher' : 'lower'} than last
        month.
      </p>
      <div className="h-64">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data}>
            <CartesianGrid stroke="var(--color-line)" vertical={false} />
            <XAxis dataKey="name" {...axisProps} />
            <YAxis {...axisProps} width={48} />
            <Tooltip formatter={money} />
            <Legend />
            <Bar
              name="Last month"
              dataKey="previous"
              fill="var(--color-chart-2)"
              radius={[4, 4, 0, 0]}
            />
            <Bar
              name="This month"
              dataKey="current"
              fill="var(--color-chart-1)"
              radius={[4, 4, 0, 0]}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </Card>
  )
}

function BudgetProgress({
  summary,
  categoryName,
}: {
  summary: Summary
  categoryName: (id: number | null) => string
}) {
  const rows = summary.categories.filter((c) => c.budget !== null)
  return (
    <Card>
      <h2 className="mb-3 text-base font-semibold">Category budgets</h2>
      {rows.length === 0 ? (
        <p className="text-muted text-sm">No category budgets set for this month.</p>
      ) : (
        <ul className="space-y-4">
          {rows.map((c) => {
            const spent = c.spend + c.fixed
            const budget = c.budget ?? 0
            const pct = budget > 0 ? Math.min(100, (spent / budget) * 100) : spent > 0 ? 100 : 0
            const over = spent - budget
            return (
              <li key={c.category_id ?? 'none'}>
                <div className="mb-1 flex justify-between text-sm">
                  <span>{categoryName(c.category_id)}</span>
                  <span className="tabular-nums">
                    {formatMoney(spent)} / {formatMoney(budget)}
                  </span>
                </div>
                <div className="bg-subtle h-2 rounded-full">
                  <div className="bg-primary h-2 rounded-full" style={{ width: `${pct}%` }} />
                </div>
                {over > 0 && (
                  <p className="text-negative mt-1 text-xs tabular-nums">
                    Over by {formatMoney(over)}
                  </p>
                )}
              </li>
            )
          })}
        </ul>
      )}
    </Card>
  )
}
