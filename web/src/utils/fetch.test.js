
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shouldUseMockData, fetchWithMock } from './fetch'

describe('shouldUseMockData', () => {
  it('returns false in production mode', () => {
    expect(shouldUseMockData({ MODE: 'production', VITE_USE_MOCK_DATA: 'true' })).toBe(false)
  })

  it('returns false when VITE_USE_MOCK_DATA is not set', () => {
    expect(shouldUseMockData({ MODE: 'development', VITE_USE_MOCK_DATA: 'false' })).toBe(false)
  })

  it('returns true in development mode with mock flag', () => {
    expect(shouldUseMockData({ MODE: 'development', VITE_USE_MOCK_DATA: 'true' })).toBe(true)
  })

  it('returns true in test mode with mock flag', () => {
    expect(shouldUseMockData({ MODE: 'test', VITE_USE_MOCK_DATA: 'true' })).toBe(true)
  })
})

describe('fetchWithMock (real API path)', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  it('calls fetch with the correct endpoint', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ total: 5 }),
    })

    const result = await fetchWithMock({
      endpoint: '/api/v1/report',
      mockPath: '',
      mockExport: '',
      env: { MODE: 'production' },
    })

    expect(fetch).toHaveBeenCalledWith('/api/v1/report', { method: 'GET' })
    expect(result).toEqual({ total: 5 })
  })

  it('throws on non-ok HTTP response', async () => {
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      text: async () => 'Internal Server Error',
    })

    await expect(
      fetchWithMock({
        endpoint: '/api/v1/report',
        mockPath: '',
        mockExport: '',
        env: { MODE: 'production' },
      })
    ).rejects.toThrow('HTTP error! status: 500')
  })

  it('sends POST with JSON body', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ ok: true }),
    })

    await fetchWithMock({
      endpoint: '/api/v1/resource/action',
      mockPath: '',
      mockExport: '',
      env: { MODE: 'production' },
      method: 'POST',
      body: { action: 'approve', name: 'my-tf', namespace: 'default' },
    })

    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/resource/action',
      expect.objectContaining({
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action: 'approve', name: 'my-tf', namespace: 'default' }),
      })
    )
  })
})
