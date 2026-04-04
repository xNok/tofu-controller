import { describe, it, expect } from 'vitest'
import { mockReport } from './report'

describe('mockReport', () => {
  it('has all required fields', () => {
    const requiredFields = ['total', 'ready', 'failed', 'suspended', 'pendingApproval', 'progressing']
    for (const field of requiredFields) {
      expect(mockReport).toHaveProperty(field)
      expect(typeof mockReport[field]).toBe('number')
    }
  })

  it('counts are consistent with total', () => {
    const { total, ready, failed, suspended, pendingApproval, progressing } = mockReport
    expect(ready + failed + suspended + pendingApproval + progressing).toBeLessThanOrEqual(total)
  })
})
