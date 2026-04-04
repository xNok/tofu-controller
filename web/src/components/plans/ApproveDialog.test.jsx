import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/preact'
import { ApproveDialog } from './ApproveDialog'

const defaultProps = {
  planId: 'plan-abc123',
  isDestroyPlan: false,
  onConfirm: vi.fn(),
  onCancel: vi.fn(),
  loading: false
}

describe('ApproveDialog', () => {
  it('renders plan ID', () => {
    render(<ApproveDialog {...defaultProps} />)
    expect(screen.getByText('plan-abc123')).toBeTruthy()
  })

  it('renders Approve Plan title', () => {
    render(<ApproveDialog {...defaultProps} />)
    expect(screen.getByRole('heading', { name: 'Approve Plan' })).toBeTruthy()
  })

  it('shows destroy warning when isDestroyPlan is true', () => {
    render(<ApproveDialog {...defaultProps} isDestroyPlan={true} />)
    expect(screen.getByText(/This is a destroy plan/i)).toBeTruthy()
    expect(screen.getByText(/permanently destroy/i)).toBeTruthy()
  })

  it('does not show destroy warning for normal plan', () => {
    render(<ApproveDialog {...defaultProps} isDestroyPlan={false} />)
    expect(screen.queryByText(/destroy plan/i)).toBeNull()
  })

  it('calls onCancel when Cancel is clicked', () => {
    const onCancel = vi.fn()
    render(<ApproveDialog {...defaultProps} onCancel={onCancel} />)
    fireEvent.click(screen.getByText('Cancel'))
    expect(onCancel).toHaveBeenCalledOnce()
  })

  it('calls onConfirm when Approve Plan button is clicked', () => {
    const onConfirm = vi.fn()
    render(<ApproveDialog {...defaultProps} onConfirm={onConfirm} />)
    const buttons = screen.getAllByRole('button')
    const approveBtn = buttons.find(b => b.textContent === 'Approve Plan')
    fireEvent.click(approveBtn)
    expect(onConfirm).toHaveBeenCalledOnce()
  })

  it('shows Approving... when loading is true', () => {
    render(<ApproveDialog {...defaultProps} loading={true} />)
    expect(screen.getByText('Approving...')).toBeTruthy()
  })

  it('disables buttons when loading', () => {
    render(<ApproveDialog {...defaultProps} loading={true} />)
    const buttons = screen.getAllByRole('button')
    for (const btn of buttons) {
      expect(btn.disabled).toBe(true)
    }
  })
})
