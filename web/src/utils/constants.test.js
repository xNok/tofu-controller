// Copyright 2024 The Flux authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

import { describe, it, expect } from 'vitest'
import { POLL_INTERVAL_MS, FAST_POLL_INTERVAL_MS, resourceStatuses, sourceKinds } from './constants'

describe('constants', () => {
  it('POLL_INTERVAL_MS is 30 seconds', () => {
    expect(POLL_INTERVAL_MS).toBe(30000)
  })

  it('FAST_POLL_INTERVAL_MS is 5 seconds', () => {
    expect(FAST_POLL_INTERVAL_MS).toBe(5000)
  })

  it('resourceStatuses includes all expected statuses', () => {
    const values = resourceStatuses.map(s => s.value)
    expect(values).toContain('Ready')
    expect(values).toContain('Failed')
    expect(values).toContain('PendingApproval')
    expect(values).toContain('Planning')
    expect(values).toContain('Applying')
    expect(values).toContain('Suspended')
  })

  it('resourceStatuses first item is "All" (empty value)', () => {
    expect(resourceStatuses[0].value).toBe('')
    expect(resourceStatuses[0].label).toBe('All')
  })

  it('sourceKinds includes OCIRepository, GitRepository, Bucket', () => {
    const values = sourceKinds.map(s => s.value)
    expect(values).toContain('OCIRepository')
    expect(values).toContain('GitRepository')
    expect(values).toContain('Bucket')
  })
})
