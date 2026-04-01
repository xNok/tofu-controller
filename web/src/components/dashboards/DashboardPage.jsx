import { reportData, reportUpdatedAt } from '../../app'
import { formatTime } from '../../utils/time'

/**
 * SummaryCard - a single metric card for the dashboard.
 */
function SummaryCard({ title, count, color, borderColor, icon, href }) {
  const cardClass = `bg-white dark:bg-gray-800 rounded-lg shadow-sm border-l-4 ${borderColor} p-6 flex items-center justify-between transition-colors`

  const content = (
    <>
      <div>
        <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{title}</p>
        <p class={`text-3xl font-bold mt-1 ${color}`}>{count}</p>
      </div>
      <div class={`w-12 h-12 rounded-full flex items-center justify-center ${color.replace('text-', 'bg-').replace('600', '100').replace('500', '100')} dark:bg-opacity-20`}>
        {icon}
      </div>
    </>
  )

  if (href) {
    return (
      <a href={href} class={`${cardClass} hover:shadow-md cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-tofu-purple`}>
        {content}
      </a>
    )
  }

  return <div class={cardClass}>{content}</div>
}

/**
 * DashboardPage - Summary dashboard with Terraform resource counts.
 *
 * Shows 4 summary cards: Total, Pending Approval, Failed, Ready.
 * Reads from the global reportData signal populated by the polling loop in app.jsx.
 */
export function DashboardPage() {
  const report = reportData.value

  if (!report) {
    return (
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
        <p class="text-gray-500 dark:text-gray-400">Loading dashboard...</p>
      </main>
    )
  }

  const cards = [
    {
      title: 'Total Resources',
      count: report.total ?? 0,
      color: 'text-gray-700 dark:text-gray-200',
      borderColor: 'border-gray-400 dark:border-gray-500',
      href: '/resources',
      icon: (
        <svg class="w-6 h-6 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
        </svg>
      )
    },
    {
      title: 'Pending Approval',
      count: report.pendingApproval ?? 0,
      color: 'text-orange-600 dark:text-orange-400',
      borderColor: 'border-orange-500',
      href: '/resources?status=PendingApproval',
      icon: (
        <svg class="w-6 h-6 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      )
    },
    {
      title: 'Failed',
      count: report.failed ?? 0,
      color: 'text-red-600 dark:text-red-400',
      borderColor: 'border-red-500',
      href: '/resources?status=Failed',
      icon: (
        <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      )
    },
    {
      title: 'Ready',
      count: report.ready ?? 0,
      color: 'text-green-600 dark:text-green-400',
      borderColor: 'border-green-500',
      href: '/resources?status=Ready',
      icon: (
        <svg class="w-6 h-6 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      )
    }
  ]

  return (
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">Dashboard</h1>
        <span class="text-sm text-gray-500 dark:text-gray-400">
          Updated {formatTime(reportUpdatedAt.value)}
        </span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {cards.map(card => (
          <SummaryCard key={card.title} {...card} />
        ))}
      </div>

      {report.suspended > 0 && (
        <div class="mt-4">
          <a
            href="/resources?status=Suspended"
            class="inline-flex items-center gap-2 text-sm text-yellow-700 dark:text-yellow-400 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md px-3 py-2 hover:bg-yellow-100 dark:hover:bg-yellow-900/30 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {report.suspended} resource{report.suspended !== 1 ? 's' : ''} suspended
          </a>
        </div>
      )}
    </main>
  )
}
