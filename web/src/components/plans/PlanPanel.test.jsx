import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/preact'
import { PlanPanel } from './PlanPanel'

const defaultProps = {
  namespace: 'default',
  name: 'my-app',
  pendingPlanId: 'plan-abc123',
  isDestroyPlan: false,
  isDriftDetectionPlan: false,
  planContent: 'Terraform will perform the following actions...',
  planLoading: false,
  planError: null,
  onApprove: vi.fn()
}

describe('PlanPanel', () => {
  it('shows pending plan ID', () => {
    render(<PlanPanel {...defaultProps} />)
    expect(screen.getByText('plan-abc123')).toBeTruthy()
  })

  it('shows plan content', () => {
    render(<PlanPanel {...defaultProps} />)
    expect(screen.getByTestId('plan-content')).toBeTruthy()
    expect(screen.getByText(/Terraform will perform/)).toBeTruthy()
  })

  it('shows Approve Plan button', () => {
    render(<PlanPanel {...defaultProps} />)
    expect(screen.getByTestId('approve-button')).toBeTruthy()
  })

  it('shows Destroy Plan badge when isDestroyPlan is true', () => {
    render(<PlanPanel {...defaultProps} isDestroyPlan={true} />)
    expect(screen.getByText('Destroy Plan')).toBeTruthy()
  })

  it('does not show Destroy Plan badge for normal plan', () => {
    render(<PlanPanel {...defaultProps} isDestroyPlan={false} />)
    expect(screen.queryByText('Destroy Plan')).toBeNull()
  })

  it('shows Drift Detection badge when isDriftDetectionPlan is true', () => {
    render(<PlanPanel {...defaultProps} isDriftDetectionPlan={true} />)
    expect(screen.getByText('Drift Detection')).toBeTruthy()
  })

  it('shows loading spinner when planLoading is true', () => {
    render(<PlanPanel {...defaultProps} planContent={null} planLoading={true} />)
    expect(screen.getByText('Loading plan output...')).toBeTruthy()
  })

  it('shows error message when planError is set', () => {
    render(<PlanPanel {...defaultProps} planContent={null} planError="Connection refused" />)
    expect(screen.getByText(/Connection refused/)).toBeTruthy()
  })

  it('shows ApproveDialog when Approve Plan button is clicked', async () => {
    render(<PlanPanel {...defaultProps} />)
    const button = screen.getByTestId('approve-button')
    fireEvent.click(button)
    // Dialog should appear
    expect(screen.getByRole('dialog')).toBeTruthy()
  })
})
