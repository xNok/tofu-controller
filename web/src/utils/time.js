
/**
 * Format a timestamp into a human-readable relative time string.
 *
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Formatted time string (e.g., "just now", "5m ago", "3h ago", or absolute date)
 */
export const formatTimestamp = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)

  if (diffMins < 1) return 'just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffMins < 1440) return `${Math.floor(diffMins / 60)}h ago`

  const includeYear = date.getFullYear() !== now.getFullYear()
  const options = {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }
  if (includeYear) {
    options.year = 'numeric'
  }
  return date.toLocaleString(undefined, options)
}

/**
 * Format a date into an absolute time string (HH:MM:SS format).
 *
 * @param {Date|null} date - The date to format
 * @returns {string} Formatted time string or "Never" if date is null
 */
export const formatTime = (date) => {
  if (!date) return 'Never'
  return new Intl.DateTimeFormat(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(date)
}
