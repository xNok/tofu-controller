// web/vitest.setup.js

// Import jest-dom matchers for DOM assertions
import '@testing-library/jest-dom/vitest'

// Mock localStorage globally before any modules are imported
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  clear: vi.fn(),
}

Object.defineProperty(global, 'localStorage', {
  value: localStorageMock,
  writable: true,
})

// Mock matchMedia globally
const matchMediaMock = vi.fn((_query) => ({
  matches: false, // Default to light theme for system preference
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
}))

Object.defineProperty(global, 'matchMedia', {
  value: matchMediaMock,
  writable: true,
})

// Mock crypto.randomUUID globally
const cryptoMock = {
  randomUUID: vi.fn(() => 'test-uuid-1234'),
}

Object.defineProperty(global, 'crypto', {
  value: cryptoMock,
  writable: true,
})

// Expose the mocks for individual test files to reset/inspect
global.localStorageMock = localStorageMock
global.matchMediaMock = matchMediaMock
global.cryptoMock = cryptoMock

// Suppress expected console messages during tests
const suppressedErrorPatterns = [
  /Failed to fetch/,
  /Failed to parse URL/,
  /Network error/,
  /Network connection failed/,
]

const originalConsoleError = console.error
console.error = (...args) => {
  const message = args.join(' ')
  const shouldSuppress = suppressedErrorPatterns.some(pattern => pattern.test(message))
  if (!shouldSuppress) {
    originalConsoleError(...args)
  }
}
