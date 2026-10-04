import { useSearchParams } from 'react-router-dom'
import { PageHeader } from '@/components/layout/PageHeader'
import { Tabs } from '@/components/ui/Tabs'
import { BudgetTab } from '@/features/finance/components/BudgetTab'
import { OverviewTab } from '@/features/finance/components/OverviewTab'
import { RecurringTab } from '@/features/finance/components/RecurringTab'
import { SettingsTab } from '@/features/finance/components/SettingsTab'
import { TransactionsTab } from '@/features/finance/components/TransactionsTab'
import { tabLabels } from '@/features/finance/labels'
import { currentMonth } from '@/lib/format'

type TabId = keyof typeof tabLabels

const tabItems = (Object.keys(tabLabels) as TabId[]).map((id) => ({ id, label: tabLabels[id] }))
const isTab = (value: string | null): value is TabId => value !== null && value in tabLabels

export function FinancePage() {
  const [params, setParams] = useSearchParams()
  const tabParam = params.get('tab')
  const tab: TabId = isTab(tabParam) ? tabParam : 'overview'
  const month = params.get('month') ?? currentMonth()

  const update = (changes: Record<string, string>) => {
    const next = new URLSearchParams(params)
    Object.entries(changes).forEach(([key, value]) => next.set(key, value))
    setParams(next, { replace: true })
  }

  return (
    <>
      <PageHeader title="Personal finance" />
      <Tabs items={tabItems} value={tab} onChange={(id) => update({ tab: id })} />
      {tab === 'overview' && (
        <OverviewTab month={month} onMonthChange={(m) => update({ month: m })} />
      )}
      {tab === 'transactions' && (
        <TransactionsTab month={month} onMonthChange={(m) => update({ month: m })} />
      )}
      {tab === 'budget' && <BudgetTab month={month} onMonthChange={(m) => update({ month: m })} />}
      {tab === 'recurring' && (
        <RecurringTab month={month} onMonthChange={(m) => update({ month: m })} />
      )}
      {tab === 'settings' && <SettingsTab />}
    </>
  )
}
