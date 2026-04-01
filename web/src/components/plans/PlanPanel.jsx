import { useState } from 'preact/hooks'
import { ApproveDialog } from './ApproveDialog'

/**
 * PlanPanel - Displays pending Terraform plan output and approve button.
 *
 * @param {Object} props
 * @param {string} props.namespace - Resource namespace
 * @param {string} props.name - Resource name
 * @param {string} props.pendingPlanId - Pending plan ID from status
 * @param {boolean} props.isDestroyPlan - Whether the plan destroys resources
 * @param {boolean} props.isDriftDetectionPlan - Whether this is a drift detection plan
 * @param {string|null} props.planContent - Decoded plan output (null if not loaded)
 * @param {boolean} props.planLoading - Whether plan content is loading
 * @param {string|null} props.planError - Error message if plan fetch failed
 * @param {Function} props.onApprove - Callback after successful approval
 */
export function PlanPanel({
  namespace,
  name,
  pendingPlanId,
  isDestroyPlan,
  isDriftDetectionPlan,
  planContent,
  planLoading,
  planError,
  onApprove
}) {
  const [showDialog, setShowDialog] = useState(false)
  const [approving, setApproving] = useState(false)
  const [approveError, setApproveError] = useState(null)

  const handleApprove = async () => {
    setApproving(true)
    setApproveError(null)
    try {
      const resp = await fetch('/api/v1/resource/action', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action: 'approve', namespace, name })
      })
      if (!resp.ok) {
        const text = await resp.text()
        throw new Error(text || `HTTP ${resp.status}`)
      }
      setShowDialog(false)
      onApprove?.()
    } catch (err) {
      setApproveError(err.message)
    } finally {
      setApproving(false)
    }
  }

  return (
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-orange-200 dark:border-orange-800">
      {/* Plan header */}
      <div class="flex items-center justify-between px-4 py-3 border-b border-orange-200 dark:border-orange-800">
        <div class="flex items-center gap-2">
          <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01" />
            </svg>
            Pending Approval
          </span>
          {isDestroyPlan && (
            <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400">
              Destroy Plan
            </span>
          )}
          {isDriftDetectionPlan && (
            <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400">
              Drift Detection
            </span>
          )}
        </div>

        <button
          data-testid="approve-button"
          onClick={() => setShowDialog(true)}
          class={`px-3 py-1.5 text-sm font-medium text-white rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-tofu-purple transition-colors ${
            isDestroyPlan ? 'bg-red-600 hover:bg-red-700' : 'bg-tofu-purple hover:opacity-90'
          }`}
        >
          Approve Plan
        </button>
      </div>

      {/* Plan ID */}
      <div class="px-4 py-2 bg-gray-50 dark:bg-gray-900/30 border-b border-gray-200 dark:border-gray-700">
        <p class="text-xs text-gray-500 dark:text-gray-400">
          Plan ID: <span class="font-mono text-gray-900 dark:text-white">{pendingPlanId}</span>
        </p>
      </div>

      {/* Plan content */}
      <div class="p-4">
        {planLoading && (
          <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
            <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            Loading plan output...
          </div>
        )}

        {planError && (
          <p class="text-sm text-red-600 dark:text-red-400">Failed to load plan: {planError}</p>
        )}

        {!planLoading && !planError && planContent && (
          <pre
            data-testid="plan-content"
            class="text-xs font-mono text-gray-800 dark:text-gray-200 bg-gray-50 dark:bg-gray-900 rounded p-3 overflow-auto max-h-96 whitespace-pre-wrap break-all"
          >
            {planContent}
          </pre>
        )}
      </div>

      {/* Approve error */}
      {approveError && (
        <div class="px-4 pb-4">
          <p class="text-sm text-red-600 dark:text-red-400">Approval failed: {approveError}</p>
        </div>
      )}

      {/* ApproveDialog */}
      {showDialog && (
        <ApproveDialog
          planId={pendingPlanId}
          isDestroyPlan={isDestroyPlan}
          onConfirm={handleApprove}
          onCancel={() => setShowDialog(false)}
          loading={approving}
        />
      )}
    </div>
  )
}
