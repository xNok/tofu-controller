import { describe, it, expect } from 'vitest'
import { mockResources } from './resources'

const VALID_STATUSES = ['Ready', 'Failed', 'Suspended', 'PendingApproval', 'Progressing', 'Planning', 'Applying']
const VALID_SOURCE_KINDS = ['GitRepository', 'OCIRepository', 'Bucket']

describe('mockResources', () => {
  it('is a non-empty array', () => {
    expect(Array.isArray(mockResources)).toBe(true)
    expect(mockResources.length).toBeGreaterThan(0)
  })

  it('each item has required fields', () => {
    const requiredFields = ['namespace', 'name', 'status', 'suspended', 'sourceKind', 'sourceName']
    for (const item of mockResources) {
      for (const field of requiredFields) {
        expect(item, `item ${item.name} missing field ${field}`).toHaveProperty(field)
      }
    }
  })

  it('each item has a valid status', () => {
    for (const item of mockResources) {
      expect(VALID_STATUSES, `unexpected status "${item.status}"`).toContain(item.status)
    }
  })

  it('each item has a valid sourceKind', () => {
    for (const item of mockResources) {
      expect(VALID_SOURCE_KINDS, `unexpected sourceKind "${item.sourceKind}"`).toContain(item.sourceKind)
    }
  })

  it('suspended items have status Suspended', () => {
    for (const item of mockResources) {
      if (item.suspended) {
        expect(item.status).toBe('Suspended')
      }
    }
  })

  it('PendingApproval items have a pendingPlan', () => {
    for (const item of mockResources) {
      if (item.status === 'PendingApproval') {
        expect(item.pendingPlan).toBeTruthy()
      }
    }
  })
})
