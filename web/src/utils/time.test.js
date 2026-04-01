// Copyright 2024 The Flux authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

import { describe, it, expect } from 'vitest'
import { formatTimestamp, formatTime } from './time'

describe('formatTimestamp', () => {
  it('returns "just now" for timestamps less than 1 minute ago', () => {
    const now = new Date()
    expect(formatTimestamp(now.toISOString())).toBe('just now')
  })

  it('returns minutes ago for timestamps less than 1 hour ago', () => {
    const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000)
    expect(formatTimestamp(fiveMinutesAgo.toISOString())).toBe('5m ago')
  })

  it('returns hours ago for timestamps less than 24 hours ago', () => {
    const threeHoursAgo = new Date(Date.now() - 3 * 60 * 60 * 1000)
    expect(formatTimestamp(threeHoursAgo.toISOString())).toBe('3h ago')
  })

  it('returns a formatted date string for older timestamps', () => {
    const oldDate = new Date(Date.now() - 2 * 24 * 60 * 60 * 1000)
    const result = formatTimestamp(oldDate.toISOString())
    // Should be a non-relative string (contains date info)
    expect(result).not.toMatch(/ago/)
    expect(result.length).toBeGreaterThan(0)
  })
})

describe('formatTime', () => {
  it('returns "Never" for null input', () => {
    expect(formatTime(null)).toBe('Never')
  })

  it('returns "Never" for undefined input', () => {
    expect(formatTime(undefined)).toBe('Never')
  })

  it('returns a time string for a valid date', () => {
    const date = new Date('2024-01-15T14:30:45Z')
    const result = formatTime(date)
    expect(result).toBeTruthy()
    expect(typeof result).toBe('string')
  })
})
