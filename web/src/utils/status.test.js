
import { describe, it, expect } from 'vitest'
import {
  getStatusBadgeClass,
  getStatusLabel,
  getStatusBorderClass,
  getStatusBarColor,
  getSourceKindBadgeClass,
} from './status'

describe('getStatusBadgeClass', () => {
  it('returns green classes for Ready', () => {
    expect(getStatusBadgeClass('Ready')).toContain('green')
  })

  it('returns red classes for Failed', () => {
    expect(getStatusBadgeClass('Failed')).toContain('red')
  })

  it('returns blue classes for Planning', () => {
    expect(getStatusBadgeClass('Planning')).toContain('blue')
  })

  it('returns blue classes for Applying', () => {
    expect(getStatusBadgeClass('Applying')).toContain('blue')
  })

  it('returns orange classes for PendingApproval', () => {
    expect(getStatusBadgeClass('PendingApproval')).toContain('orange')
  })

  it('returns yellow classes for Suspended', () => {
    expect(getStatusBadgeClass('Suspended')).toContain('yellow')
  })

  it('returns gray classes for unknown status', () => {
    expect(getStatusBadgeClass('Unknown')).toContain('gray')
  })

  it('returns gray classes for empty string', () => {
    expect(getStatusBadgeClass('')).toContain('gray')
  })
})

describe('getStatusLabel', () => {
  it('returns "Pending Approval" for PendingApproval', () => {
    expect(getStatusLabel('PendingApproval')).toBe('Pending Approval')
  })

  it('returns the status as-is for known statuses', () => {
    expect(getStatusLabel('Ready')).toBe('Ready')
    expect(getStatusLabel('Failed')).toBe('Failed')
    expect(getStatusLabel('Suspended')).toBe('Suspended')
  })

  it('returns "Unknown" for empty/null values', () => {
    expect(getStatusLabel('')).toBe('Unknown')
    expect(getStatusLabel(null)).toBe('Unknown')
    expect(getStatusLabel(undefined)).toBe('Unknown')
  })
})

describe('getStatusBorderClass', () => {
  it('returns success border for Ready', () => {
    expect(getStatusBorderClass('Ready')).toBe('border-success')
  })

  it('returns danger border for Failed', () => {
    expect(getStatusBorderClass('Failed')).toBe('border-danger')
  })

  it('returns warning border for Suspended', () => {
    expect(getStatusBorderClass('Suspended')).toBe('border-warning')
  })

  it('returns info border for Planning', () => {
    expect(getStatusBorderClass('Planning')).toBe('border-info')
  })

  it('returns orange border for PendingApproval', () => {
    expect(getStatusBorderClass('PendingApproval')).toContain('orange')
  })
})

describe('getStatusBarColor', () => {
  it('returns green for Ready', () => {
    expect(getStatusBarColor('Ready')).toContain('green')
  })

  it('returns red for Failed', () => {
    expect(getStatusBarColor('Failed')).toContain('red')
  })

  it('returns orange for PendingApproval', () => {
    expect(getStatusBarColor('PendingApproval')).toContain('orange')
  })
})

describe('getSourceKindBadgeClass', () => {
  it('returns purple classes for OCIRepository', () => {
    expect(getSourceKindBadgeClass('OCIRepository')).toContain('purple')
  })

  it('returns indigo classes for GitRepository', () => {
    expect(getSourceKindBadgeClass('GitRepository')).toContain('indigo')
  })

  it('returns cyan classes for Bucket', () => {
    expect(getSourceKindBadgeClass('Bucket')).toContain('cyan')
  })

  it('returns gray classes for unknown kind', () => {
    expect(getSourceKindBadgeClass('Unknown')).toContain('gray')
  })
})
