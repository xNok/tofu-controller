import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/preact'
import { ResourceCard } from './ResourceCard'

const baseResource = {
  namespace: 'default',
  name: 'my-app',
  status: 'Ready',
  suspended: false,
  sourceKind: 'OCIRepository',
  sourceName: 'my-oci-repo',
  message: 'Applied revision: sha256:abc123',
  lastUpdated: '2024-01-15T10:30:00Z'
}

describe('ResourceCard', () => {
  it('renders namespace/name link', () => {
    render(<ResourceCard resource={baseResource} />)
    const link = screen.getByText('default/my-app')
    expect(link).toBeTruthy()
    expect(link.getAttribute('href')).toBe('/resource/default/my-app')
  })

  it('shows status badge', () => {
    render(<ResourceCard resource={baseResource} />)
    expect(screen.getByText('Ready')).toBeTruthy()
  })

  it('shows sourceKind badge', () => {
    render(<ResourceCard resource={baseResource} />)
    expect(screen.getByText('OCIRepository')).toBeTruthy()
  })

  it('shows source name', () => {
    render(<ResourceCard resource={baseResource} />)
    expect(screen.getByText('Source: my-oci-repo')).toBeTruthy()
  })

  it('shows status message', () => {
    render(<ResourceCard resource={baseResource} />)
    expect(screen.getByText('Applied revision: sha256:abc123')).toBeTruthy()
  })

  it('shows Pending Approval label for PendingApproval status', () => {
    const resource = { ...baseResource, status: 'PendingApproval' }
    render(<ResourceCard resource={resource} />)
    expect(screen.getByText('Pending Approval')).toBeTruthy()
  })

  it('renders card with testid', () => {
    render(<ResourceCard resource={baseResource} />)
    expect(screen.getByTestId('resource-card')).toBeTruthy()
  })

  it('renders failed status badge with correct class', () => {
    const resource = { ...baseResource, status: 'Failed' }
    const { container } = render(<ResourceCard resource={resource} />)
    const badge = container.querySelector('.bg-red-100')
    expect(badge).toBeTruthy()
  })
})
