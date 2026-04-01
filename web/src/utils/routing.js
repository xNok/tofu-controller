// Copyright 2024 The Flux authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

import { useEffect } from 'preact/hooks'
import { useLocation } from 'preact-iso'

/**
 * Serialize filter signals to URL query string.
 *
 * @param {Object} filters - Object mapping filter names to signal values
 * @returns {string} URL query string (without leading ?)
 */
export function serializeFilters(filters) {
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(filters)) {
    if (value) {
      params.set(key, value)
    }
  }
  return params.toString()
}

/**
 * Hook to restore filter signals from URL query parameters on mount.
 *
 * @param {Object} signals - Object mapping filter names to writable signals
 */
export function useRestoreFiltersFromUrl(signals) {
  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    for (const [key, signal] of Object.entries(signals)) {
      const value = params.get(key)
      if (value !== null) {
        signal.value = value
      }
    }
  }, [])
}

/**
 * Hook to sync filter signals to URL query parameters (debounced).
 *
 * @param {Object} signals - Object mapping filter names to readable signals
 * @param {number} delay - Debounce delay in ms (default: 300)
 */
export function useSyncFiltersToUrl(signals, delay = 300) {
  const location = useLocation()

  useEffect(() => {
    const timer = setTimeout(() => {
      const filterValues = {}
      for (const [key, signal] of Object.entries(signals)) {
        filterValues[key] = signal.value
      }
      const query = serializeFilters(filterValues)
      const newUrl = query ? `${location.path}?${query}` : location.path
      window.history.replaceState(null, '', newUrl)
    }, delay)

    return () => clearTimeout(timer)
  }, Object.values(signals).map(s => s.value))
}
