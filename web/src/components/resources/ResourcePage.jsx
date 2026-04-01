import { useEffect, useState } from 'preact/hooks'
import { useRoute } from 'preact-iso'
import { getStatusBadgeClass, getStatusLabel, getSourceKindBadgeClass } from '../../utils/status'
import { formatTimestamp } from '../../utils/time'
import { FAST_POLL_INTERVAL_MS, POLL_INTERVAL_MS } from '../../utils/constants'
import { fetchWithMock } from '../../utils/fetch'
import { PlanPanel } from '../plans/PlanPanel'

/**
 * ResourcePage - Full resource detail with plan approval.
 *
 * Route: /resource/:namespace/:name
 */
export function ResourcePage() {
  const { params } = useRoute()
  const { namespace, name } = params

  const [resource, setResource] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [planContent, setPlanContent] = useState(null)
  const [planLoading, setPlanLoading] = useState(false)
  const [planError, setPlanError] = useState(null)

  async function fetchResource() {
    const endpoint = `/api/v1/resource?namespace=${encodeURIComponent(namespace)}&name=${encodeURIComponent(name)}`
    try {
      const data = await fetchWithMock({
        endpoint,
        mockPath: '../../mock/resource',
        mockExport: 'getMockResource'
      })
      setResource(data)
      setError(null)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function fetchPlan(ns, n) {
    const endpoint = `/api/v1/plan?namespace=${encodeURIComponent(ns)}&name=${encodeURIComponent(n)}`
    setPlanLoading(true)
    setPlanError(null)
    try {
      const data = await fetchWithMock({
        endpoint,
        mockPath: '../../mock/plan',
        mockExport: 'getMockPlan'
      })
      setPlanContent(data.content)
    } catch (err) {
      setPlanError(err.message)
    } finally {
      setPlanLoading(false)
    }
  }

  useEffect(() => {
    fetchResource()
  }, [namespace, name])

  // Fetch plan content when pending plan is detected
  useEffect(() => {
    if (resource?.pendingPlan) {
      fetchPlan(namespace, name)
    } else {
      setPlanContent(null)
      setPlanError(null)
    }
  }, [resource?.pendingPlan, namespace, name])

  // Fast poll while pending approval or progressing
  useEffect(() => {
    const isPendingOrActive = resource?.status === 'PendingApproval' ||
      resource?.status === 'Planning' ||
      resource?.status === 'Applying' ||
      resource?.status === 'Progressing'

    const interval = isPendingOrActive ? FAST_POLL_INTERVAL_MS : POLL_INTERVAL_MS
    const timer = setInterval(fetchResource, interval)
    return () => clearInterval(timer)
  }, [resource?.status, namespace, name])

  if (loading) {
    return (
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
        <div class="animate-pulse space-y-4">
          <div class="h-6 bg-gray-200 dark:bg-gray-700 rounded w-1/3"></div>
          <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
        </div>
      </main>
    )
  }

  if (error) {
    return (
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <p class="text-sm text-red-700 dark:text-red-400">Failed to load resource: {error}</p>
        </div>
      </main>
    )
  }

  if (!resource) return null

  return (
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
      {/* Breadcrumb */}
      <nav class="text-sm text-gray-500 dark:text-gray-400 mb-4">
        <a href="/resources" class="hover:text-gray-700 dark:hover:text-gray-200 transition-colors">
          Resources
        </a>
        <span class="mx-2">/</span>
        <span class="text-gray-900 dark:text-white font-medium">{namespace}/{name}</span>
      </nav>

      {/* Header */}
      <div class="flex items-start justify-between gap-3 mb-6">
        <div>
          <div class="flex items-center gap-2 mb-1 flex-wrap">
            <span class={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${getStatusBadgeClass(resource.status)}`}>
              {getStatusLabel(resource.status)}
            </span>
            <span class={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${getSourceKindBadgeClass(resource.sourceKind)}`}>
              {resource.sourceKind}
            </span>
            {resource.suspended && (
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400">
                Suspended
              </span>
            )}
          </div>
          <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
            {namespace}/{name}
          </h1>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
            Source: {resource.sourceName}
            {resource.lastUpdated && ` · Updated ${formatTimestamp(resource.lastUpdated)}`}
          </p>
        </div>
      </div>

      {/* Pending plan panel */}
      {resource.pendingPlan && (
        <div class="mb-6">
          <PlanPanel
            namespace={namespace}
            name={name}
            pendingPlanId={resource.pendingPlan}
            isDestroyPlan={resource.terraformStatus?.plan?.isDestroyPlan ?? false}
            isDriftDetectionPlan={resource.terraformStatus?.plan?.isDriftDetectionPlan ?? false}
            planContent={planContent}
            planLoading={planLoading}
            planError={planError}
            onApprove={fetchResource}
          />
        </div>
      )}

      {/* Status message */}
      {resource.message && (
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 mb-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Status Message</p>
          <p class="text-sm font-mono text-gray-800 dark:text-gray-200 break-all">{resource.message}</p>
        </div>
      )}
    </main>
  )
}
