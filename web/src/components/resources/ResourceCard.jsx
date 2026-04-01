import { useState } from 'preact/hooks'
import { getStatusBadgeClass, getStatusLabel, getSourceKindBadgeClass } from '../../utils/status'
import { formatTimestamp } from '../../utils/time'

/**
 * ResourceCard - Displays a single Terraform resource with status badge and source type.
 *
 * @param {Object} props
 * @param {Object} props.resource - ResourceSummary from API
 */
export function ResourceCard({ resource }) {
  const [isExpanded, setIsExpanded] = useState(false)

  const isLongMessage = resource.message && (resource.message.length > 150 || resource.message.includes('\n'))
  const shouldTruncate = isLongMessage && !isExpanded

  const displayMessage = shouldTruncate
    ? resource.message.split('\n')[0].substring(0, 150) + '...'
    : (resource.message || '')

  return (
    <div
      data-testid="resource-card"
      class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 hover:shadow-md transition-shadow"
    >
      <div class="flex items-start justify-between gap-2">
        <div class="flex-grow min-w-0">
          <div class="flex items-center gap-2 flex-wrap mb-1">
            {/* Status badge */}
            <span class={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${getStatusBadgeClass(resource.status)}`}>
              {resource.status === 'PendingApproval' && (
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01" />
                </svg>
              )}
              {getStatusLabel(resource.status)}
            </span>

            {/* Source kind badge */}
            <span class={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${getSourceKindBadgeClass(resource.sourceKind)}`}>
              {resource.sourceKind}
            </span>
          </div>

          {/* Resource name and namespace */}
          <a
            href={`/resource/${encodeURIComponent(resource.namespace)}/${encodeURIComponent(resource.name)}`}
            class="text-sm font-semibold text-gray-900 dark:text-white hover:text-tofu-purple dark:hover:text-purple-400 transition-colors"
          >
            {resource.namespace}/{resource.name}
          </a>

          {/* Source name */}
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
            Source: {resource.sourceName}
          </p>
        </div>

        {/* Last updated */}
        {resource.lastUpdated && (
          <span class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap flex-shrink-0">
            {formatTimestamp(resource.lastUpdated)}
          </span>
        )}
      </div>

      {/* Status message */}
      {resource.message && (
        <div class="mt-2">
          <p class="text-xs text-gray-600 dark:text-gray-400 font-mono break-all">
            {displayMessage}
          </p>
          {isLongMessage && (
            <button
              onClick={() => setIsExpanded(!isExpanded)}
              class="text-xs text-tofu-purple dark:text-purple-400 hover:underline mt-1"
            >
              {isExpanded ? 'Show less' : 'Show more'}
            </button>
          )}
        </div>
      )}
    </div>
  )
}
