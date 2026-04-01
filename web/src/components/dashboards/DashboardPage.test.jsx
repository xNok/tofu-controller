import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/preact'
import { DashboardPage } from './DashboardPage'
import { reportData, reportUpdatedAt } from '../../app'
import { mockReport } from '../../mock/report'

describe('DashboardPage', () => {
  beforeEach(() => {
    reportData.value = null
    reportUpdatedAt.value = null
  })

  it('shows loading message when no data', () => {
    render(<DashboardPage />)
    expect(screen.getByText(/Loading dashboard/i)).toBeTruthy()
  })

  it('shows summary cards when data is loaded', () => {
    reportData.value = mockReport
    reportUpdatedAt.value = new Date()

    render(<DashboardPage />)

    expect(screen.getByText('Total Resources')).toBeTruthy()
    expect(screen.getByText('Pending Approval')).toBeTruthy()
    expect(screen.getByText('Failed')).toBeTruthy()
    expect(screen.getByText('Ready')).toBeTruthy()
  })

  it('displays correct counts from report data', () => {
    reportData.value = { total: 10, ready: 7, failed: 2, suspended: 0, pendingApproval: 1, progressing: 0 }
    reportUpdatedAt.value = new Date()

    render(<DashboardPage />)

    expect(screen.getByText('10')).toBeTruthy()
    expect(screen.getByText('7')).toBeTruthy()
    expect(screen.getByText('2')).toBeTruthy()
    expect(screen.getByText('1')).toBeTruthy()
  })

  it('shows suspended banner when suspended > 0', () => {
    reportData.value = { total: 5, ready: 3, failed: 0, suspended: 2, pendingApproval: 0, progressing: 0 }
    reportUpdatedAt.value = new Date()

    render(<DashboardPage />)

    expect(screen.getByText(/2 resources suspended/i)).toBeTruthy()
  })

  it('hides suspended banner when suspended == 0', () => {
    reportData.value = { total: 5, ready: 5, failed: 0, suspended: 0, pendingApproval: 0, progressing: 0 }
    reportUpdatedAt.value = new Date()

    render(<DashboardPage />)

    expect(screen.queryByText(/suspended/i)).toBeNull()
  })
})
