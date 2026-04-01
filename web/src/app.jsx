
import { useEffect } from 'preact/hooks'
import { signal } from '@preact/signals'
import { LocationProvider, Router, Route, useLocation } from 'preact-iso'
import { fetchWithMock } from './utils/fetch'
import { POLL_INTERVAL_MS } from './utils/constants'
import './utils/theme'
import { ConnectionStatus } from './components/layout/ConnectionStatus'
import { Header } from './components/layout/Header'
import { Footer } from './components/layout/Footer'
import { NotFoundPage } from './components/layout/NotFoundPage'
import { TofuControllerIcon } from './components/layout/Icons'
import { DashboardPage } from './components/dashboards/DashboardPage'

// Global signals for report data and application state

// Summary report data from API (null until first successful fetch)
export const reportData = signal(null)

// Timestamp of last successful data fetch
export const reportUpdatedAt = signal(null)

// Initial loading state (prevents flash of error state on first load)
export const reportLoading = signal(true)

// Report fetch error message (null if no error)
export const reportError = signal(null)

// Connection status: 'loading' | 'connected' | 'disconnected'
export const connectionStatus = signal('loading')

/**
 * Fetches the report summary from the API or mock data.
 * Called on mount and every 30 seconds via auto-refresh.
 */
export async function fetchReport() {
  if (connectionStatus.value !== 'loading') {
    connectionStatus.value = 'loading'
  }
  if (!reportData.value) {
    reportLoading.value = true
  }

  try {
    reportData.value = await fetchWithMock({
      endpoint: '/api/v1/report',
      mockPath: '../mock/report',
      mockExport: 'mockReport'
    })

    connectionStatus.value = 'connected'
    reportUpdatedAt.value = new Date()
    reportLoading.value = false
    reportError.value = null
  } catch (error) {
    connectionStatus.value = 'disconnected'
    reportUpdatedAt.value = new Date()
    reportLoading.value = false
    reportError.value = error.message
  }
}

/**
 * TabNavigation - tabs for switching between Dashboard and Resources views
 */
function TabNavigation() {
  const location = useLocation()
  const currentPath = location.path

  return (
    <div class="border-b border-gray-200 dark:border-gray-700 transition-colors">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <nav class="-mb-px flex space-x-8">
          <a
            href="/"
            class={`${
              currentPath === '/'
                ? 'border-tofu-purple text-tofu-purple dark:text-purple-400'
                : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm focus:outline-none transition-colors`}
          >
            Dashboard
          </a>
          <a
            href="/resources"
            class={`${
              currentPath === '/resources'
                ? 'border-tofu-purple text-tofu-purple dark:text-purple-400'
                : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm focus:outline-none transition-colors`}
          >
            Resources
          </a>
        </nav>
      </div>
    </div>
  )
}

/**
 * App - Root application component
 *
 * Manages:
 * - Initial data fetch and auto-refresh (30s interval)
 * - Loading / error / normal states
 * - Global layout and routing
 */
export function App() {
  useEffect(() => {
    fetchReport()
    const interval = setInterval(fetchReport, POLL_INTERVAL_MS)
    return () => clearInterval(interval)
  }, [])

  // LOADING STATE
  if (reportLoading.value && !reportData.value && !reportError.value) {
    return (
      <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors flex flex-col">
        <ConnectionStatus />
        <div class="flex items-center justify-center flex-1">
          <div class="text-center">
            <TofuControllerIcon className="animate-spin h-12 w-12 text-tofu-purple mx-auto" />
            <p class="mt-4 text-gray-600 dark:text-gray-400">Loading Tofu Controller...</p>
          </div>
        </div>
      </div>
    )
  }

  // ERROR STATE: failed to load and no cached data available
  if (!reportData.value && reportError.value) {
    return (
      <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors flex flex-col items-center justify-center px-4 py-12 sm:py-0">
        <ConnectionStatus />
        <div class="w-full max-w-md">
          <div class="flex items-center justify-center gap-3 mb-8">
            <TofuControllerIcon className="w-12 h-12 text-tofu-purple" />
            <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
              Tofu Controller
            </h1>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-8">
            <div class="flex justify-center mb-6">
              <div class="w-16 h-16 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center">
                <svg class="w-8 h-8 text-red-500 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
              </div>
            </div>

            <div class="text-center mb-6">
              <h2 class="text-xl font-medium text-gray-900 dark:text-white mb-3">
                API Server Unavailable
              </h2>
              <p class="text-sm text-gray-600 dark:text-gray-400 leading-relaxed">
                Unable to connect to the server
              </p>
            </div>

            <div class="flex items-center justify-center gap-2 text-gray-500 dark:text-gray-400">
              <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              <span class="text-sm">Retrying automatically...</span>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // NORMAL STATE: render full application layout with routing
  return (
    <LocationProvider>
      <AppContent />
    </LocationProvider>
  )
}

/**
 * AppContent - Main app content with routing
 * Separated to allow useLocation() hook usage
 */
function AppContent() {
  const location = useLocation()
  const currentPath = location.path
  const isTabView = currentPath === '/' || currentPath === '/resources'

  return (
    <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors flex flex-col">
      <ConnectionStatus />
      <Header />
      {isTabView && <TabNavigation />}

      <Router>
        <Route path="/" component={DashboardPage} />
        <Route path="/resources" component={ResourcesPlaceholder} />
        <Route path="/resource/:namespace/:name" component={ResourcePlaceholder} />
        <Route default component={NotFoundPage} />
      </Router>

      <Footer />
    </div>
  )
}

// Placeholder components — replaced in later steps
function ResourcesPlaceholder() {
  return (
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
      <p class="text-gray-500 dark:text-gray-400">Resources list coming soon...</p>
    </main>
  )
}

function ResourcePlaceholder() {
  return (
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
      <p class="text-gray-500 dark:text-gray-400">Resource detail coming soon...</p>
    </main>
  )
}
