import { signal } from '@preact/signals'
import { useEffect } from 'preact/hooks'
import { fetchWithMock } from '../../utils/fetch'
import { resourceStatuses, sourceKinds } from '../../utils/constants'
import { useRestoreFiltersFromUrl, useSyncFiltersToUrl } from '../../utils/routing'
import { ResourceCard } from './ResourceCard'

// Filter signals — updated by filter form and URL query params
export const filterNamespace = signal('')
export const filterName = signal('')
export const filterStatus = signal('')
export const filterSourceKind = signal('')

// Resources data signals
export const resourcesData = signal([])
export const resourcesLoading = signal(false)
export const resourcesError = signal(null)

/**
 * Fetch resources from API with current filter values.
 */
export async function fetchResources() {
  resourcesLoading.value = true
  resourcesError.value = null

  const params = new URLSearchParams()
  if (filterNamespace.value) params.set('namespace', filterNamespace.value)
  if (filterName.value) params.set('name', filterName.value)
  if (filterStatus.value) params.set('status', filterStatus.value)
  if (filterSourceKind.value) params.set('sourceKind', filterSourceKind.value)

  try {
    const data = await fetchWithMock({
      endpoint: `/api/v1/resources?${params.toString()}`,
      mockPath: '../../mock/resources',
      mockExport: 'mockResources'
    })
    resourcesData.value = Array.isArray(data) ? data : []
  } catch (error) {
    resourcesError.value = error.message
    resourcesData.value = []
  } finally {
    resourcesLoading.value = false
  }
}

/**
 * ResourceList - Filterable list of all Terraform resources.
 *
 * Supports filtering by namespace, name (contains), status, and sourceKind.
 * Filters are persisted to/restored from URL query parameters.
 */
export function ResourceList() {
  const filterSignals = { namespace: filterNamespace, name: filterName, status: filterStatus, sourceKind: filterSourceKind }

  useRestoreFiltersFromUrl(filterSignals)
  useSyncFiltersToUrl(filterSignals)

  // Fetch on mount and whenever filters change
  useEffect(() => {
    fetchResources()
  }, [filterNamespace.value, filterName.value, filterStatus.value, filterSourceKind.value])

  const resources = resourcesData.value
  const isLoading = resourcesLoading.value
  const error = resourcesError.value

  return (
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">Resources</h1>
        <span class="text-sm text-gray-500 dark:text-gray-400">
          {isLoading ? 'Loading...' : `${resources.length} resource${resources.length !== 1 ? 's' : ''}`}
        </span>
      </div>

      {/* Filter form */}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 mb-6">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">Namespace</label>
            <input
              type="text"
              placeholder="All namespaces"
              value={filterNamespace.value}
              onInput={e => { filterNamespace.value = e.target.value }}
              class="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-tofu-purple"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">Name</label>
            <input
              type="text"
              placeholder="Filter by name"
              value={filterName.value}
              onInput={e => { filterName.value = e.target.value }}
              class="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-tofu-purple"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">Status</label>
            <select
              value={filterStatus.value}
              onChange={e => { filterStatus.value = e.target.value }}
              class="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-tofu-purple"
            >
              {resourceStatuses.map(s => (
                <option key={s.value} value={s.value}>{s.label}</option>
              ))}
            </select>
          </div>

          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">Source Kind</label>
            <select
              value={filterSourceKind.value}
              onChange={e => { filterSourceKind.value = e.target.value }}
              class="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-tofu-purple"
            >
              {sourceKinds.map(k => (
                <option key={k.value} value={k.value}>{k.label}</option>
              ))}
            </select>
          </div>
        </div>

        {/* Clear filters */}
        {(filterNamespace.value || filterName.value || filterStatus.value || filterSourceKind.value) && (
          <button
            onClick={() => {
              filterNamespace.value = ''
              filterName.value = ''
              filterStatus.value = ''
              filterSourceKind.value = ''
            }}
            class="mt-3 text-xs text-tofu-purple dark:text-purple-400 hover:underline"
          >
            Clear filters
          </button>
        )}
      </div>

      {/* Error state */}
      {error && (
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-4">
          <p class="text-sm text-red-700 dark:text-red-400">{error}</p>
        </div>
      )}

      {/* Loading skeleton */}
      {isLoading && resources.length === 0 && (
        <div class="space-y-3">
          {[1, 2, 3].map(i => (
            <div key={i} class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 animate-pulse">
              <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/3 mb-2"></div>
              <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
            </div>
          ))}
        </div>
      )}

      {/* Empty state */}
      {!isLoading && resources.length === 0 && !error && (
        <div class="text-center py-12">
          <svg class="mx-auto w-12 h-12 text-gray-400 dark:text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p class="text-gray-500 dark:text-gray-400">No resources found</p>
        </div>
      )}

      {/* Resource cards */}
      {resources.length > 0 && (
        <div class="space-y-3">
          {resources.map(r => (
            <ResourceCard key={`${r.namespace}/${r.name}`} resource={r} />
          ))}
        </div>
      )}
    </main>
  )
}
