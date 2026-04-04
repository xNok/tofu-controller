import { mockResources } from './resources'

// A fully detailed resource with spec and status for the detail view.
// Used when the endpoint matches; falls back to first resource with pendingPlan.
const detailedResources = mockResources.map(r => ({
  ...r,
  spec: {
    approvePlan: '',
    suspend: r.suspended,
    sourceRef: { kind: r.sourceKind, name: r.sourceName, namespace: r.namespace }
  },
  terraformStatus: {
    plan: {
      pending: r.pendingPlan || '',
      lastApplied: r.pendingPlan ? '' : 'plan-20240114-180000',
      isDestroyPlan: false,
      isDriftDetectionPlan: false
    },
    conditions: r.message ? [
      {
        type: 'Ready',
        status: r.status === 'Ready' ? 'True' : r.status === 'Failed' ? 'False' : 'Unknown',
        message: r.message,
        lastTransitionTime: r.lastUpdated
      }
    ] : []
  }
}))

// Override one resource to have a destroy plan for demo purposes
const destroyIdx = detailedResources.findIndex(r => r.pendingPlan)
if (destroyIdx >= 0) {
  detailedResources[destroyIdx].terraformStatus.plan.isDestroyPlan = false
}

/**
 * Mock for GET /api/v1/resource?namespace=...&name=...
 * Parses namespace/name from the endpoint query string.
 */
export function getMockResource(endpoint) {
  const url = new URL(endpoint, 'http://localhost')
  const ns = url.searchParams.get('namespace')
  const name = url.searchParams.get('name')

  const found = detailedResources.find(r => r.namespace === ns && r.name === name)
  if (found) return found

  // Default: return first resource with pending plan, or just first resource
  return detailedResources.find(r => r.pendingPlan) || detailedResources[0]
}
