/**
 * Mock for GET /api/v1/resources?namespace=...&name=...&status=...&sourceKind=...
 * Applies the same filters the real API would apply.
 */
export function getMockResources(endpoint) {
  const url = new URL(endpoint, 'http://localhost')
  const ns = url.searchParams.get('namespace') || ''
  const name = url.searchParams.get('name') || ''
  const status = url.searchParams.get('status') || ''
  const sourceKind = url.searchParams.get('sourceKind') || ''

  return mockResources.filter(r => {
    if (ns && r.namespace !== ns) return false
    if (name && !r.name.includes(name)) return false
    if (status && r.status !== status) return false
    if (sourceKind && r.sourceKind !== sourceKind) return false
    return true
  })
}

export const mockResources = [
  {
    namespace: 'production',
    name: 'aws-infra',
    status: 'Ready',
    suspended: false,
    sourceKind: 'OCIRepository',
    sourceName: 'aws-modules',
    message: 'Applied revision: sha256:abc123',
    lastUpdated: '2024-01-15T10:30:00Z'
  },
  {
    namespace: 'production',
    name: 'gcp-infra',
    status: 'PendingApproval',
    suspended: false,
    sourceKind: 'OCIRepository',
    sourceName: 'gcp-modules',
    pendingPlan: 'plan-20240115-103045',
    message: 'Planned with changes, awaiting approval',
    lastUpdated: '2024-01-15T10:30:45Z'
  },
  {
    namespace: 'staging',
    name: 'staging-db',
    status: 'Failed',
    suspended: false,
    sourceKind: 'GitRepository',
    sourceName: 'infra-repo',
    message: 'TFExecApplyFailed: Error creating resource',
    lastUpdated: '2024-01-15T09:00:00Z'
  },
  {
    namespace: 'staging',
    name: 'staging-vpc',
    status: 'Ready',
    suspended: false,
    sourceKind: 'GitRepository',
    sourceName: 'infra-repo',
    message: 'Applied revision: main/abc456',
    lastUpdated: '2024-01-14T18:00:00Z'
  },
  {
    namespace: 'development',
    name: 'dev-cluster',
    status: 'Suspended',
    suspended: true,
    sourceKind: 'Bucket',
    sourceName: 'tf-modules',
    message: 'Reconciliation is suspended',
    lastUpdated: '2024-01-10T12:00:00Z'
  },
  {
    namespace: 'development',
    name: 'dev-network',
    status: 'Progressing',
    suspended: false,
    sourceKind: 'OCIRepository',
    sourceName: 'network-modules',
    message: 'Progressing: initializing terraform',
    lastUpdated: '2024-01-15T10:31:00Z'
  }
]
