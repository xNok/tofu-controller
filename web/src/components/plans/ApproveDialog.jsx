/**
 * ApproveDialog - Confirmation modal for approving a Terraform plan.
 *
 * @param {Object} props
 * @param {string} props.planId - ID of the plan to approve
 * @param {boolean} props.isDestroyPlan - Whether this plan will destroy resources
 * @param {Function} props.onConfirm - Called when user confirms approval
 * @param {Function} props.onCancel - Called when user cancels
 * @param {boolean} props.loading - Whether the approval request is in progress
 */
export function ApproveDialog({ planId, isDestroyPlan, onConfirm, onCancel, loading }) {
  return (
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="approve-dialog-title"
    >
      {/* Backdrop */}
      <div
        class="absolute inset-0 bg-black bg-opacity-50"
        onClick={onCancel}
      />

      {/* Dialog */}
      <div class="relative bg-white dark:bg-gray-800 rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 w-full max-w-md p-6">
        {/* Header */}
        <div class="flex items-center gap-3 mb-4">
          {isDestroyPlan ? (
            <div class="w-10 h-10 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
          ) : (
            <div class="w-10 h-10 rounded-full bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-orange-600 dark:text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          )}
          <div>
            <h2 id="approve-dialog-title" class="text-lg font-semibold text-gray-900 dark:text-white">
              Approve Plan
            </h2>
            <p class="text-sm text-gray-500 dark:text-gray-400 font-mono">{planId}</p>
          </div>
        </div>

        {/* Destroy warning */}
        {isDestroyPlan && (
          <div class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md">
            <p class="text-sm font-medium text-red-700 dark:text-red-400">
              ⚠ This is a destroy plan
            </p>
            <p class="text-sm text-red-600 dark:text-red-300 mt-1">
              Approving will permanently destroy the managed infrastructure.
            </p>
          </div>
        )}

        <p class="text-sm text-gray-600 dark:text-gray-300 mb-6">
          This will set <code class="bg-gray-100 dark:bg-gray-700 px-1 rounded">spec.approvePlan</code> to the
          pending plan ID, triggering the controller to apply the plan.
        </p>

        {/* Actions */}
        <div class="flex justify-end gap-3">
          <button
            onClick={onCancel}
            disabled={loading}
            class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-tofu-purple disabled:opacity-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={onConfirm}
            disabled={loading}
            class={`px-4 py-2 text-sm font-medium text-white rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-tofu-purple disabled:opacity-50 transition-colors ${
              isDestroyPlan
                ? 'bg-red-600 hover:bg-red-700'
                : 'bg-tofu-purple hover:opacity-90'
            }`}
          >
            {loading ? 'Approving...' : 'Approve Plan'}
          </button>
        </div>
      </div>
    </div>
  )
}
