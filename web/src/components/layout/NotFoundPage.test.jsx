
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/preact'
import { NotFoundPage } from './NotFoundPage'

describe('NotFoundPage', () => {
  it('renders a 404 heading', () => {
    render(<NotFoundPage />)
    expect(screen.getByText('Page Not Found')).toBeInTheDocument()
  })

  it('renders a link back to the dashboard', () => {
    render(<NotFoundPage />)
    const link = screen.getByRole('link', { name: /go to dashboard/i })
    expect(link).toBeInTheDocument()
    expect(link.getAttribute('href')).toBe('/')
  })

  it('has the correct test id', () => {
    render(<NotFoundPage />)
    expect(screen.getByTestId('not-found-page')).toBeInTheDocument()
  })
})
