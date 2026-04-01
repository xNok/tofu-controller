// Copyright 2024 The Flux authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

/**
 * Get badge CSS classes for Terraform resource statuses.
 *
 * @param {string} status - Status value
 * @returns {string} Tailwind CSS classes for the badge
 */
export function getStatusBadgeClass(status) {
  switch (status) {
  case 'Ready':
    return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
  case 'Failed':
    return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
  case 'Planning':
  case 'Applying':
    return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
  case 'PendingApproval':
    return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400'
  case 'Suspended':
    return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
  default:
    return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400'
  }
}

/**
 * Get the human-readable display label for a status value.
 *
 * @param {string} status - Status value
 * @returns {string} Display label
 */
export function getStatusLabel(status) {
  switch (status) {
  case 'PendingApproval':
    return 'Pending Approval'
  default:
    return status || 'Unknown'
  }
}

/**
 * Get border color CSS class for resource status.
 *
 * @param {string} status - Status value
 * @returns {string} Tailwind CSS classes for the border
 */
export function getStatusBorderClass(status) {
  switch (status) {
  case 'Ready':
    return 'border-success'
  case 'Failed':
    return 'border-danger'
  case 'Suspended':
    return 'border-warning'
  case 'Planning':
  case 'Applying':
    return 'border-info'
  case 'PendingApproval':
    return 'border-orange-500 dark:border-orange-400'
  default:
    return 'border-gray-300 dark:border-gray-600'
  }
}

/**
 * Get solid background color class for status bars/charts.
 *
 * @param {string} status - Status value
 * @returns {string} Tailwind CSS classes for the background
 */
export function getStatusBarColor(status) {
  switch (status) {
  case 'Ready':
    return 'bg-green-500 dark:bg-green-600'
  case 'Failed':
    return 'bg-red-500 dark:bg-red-600'
  case 'Planning':
  case 'Applying':
    return 'bg-blue-500 dark:bg-blue-600'
  case 'PendingApproval':
    return 'bg-orange-500 dark:bg-orange-600'
  case 'Suspended':
    return 'bg-yellow-500 dark:bg-yellow-600'
  default:
    return 'bg-gray-200 dark:bg-gray-700'
  }
}

/**
 * Get badge CSS classes for source kind badges.
 *
 * @param {string} kind - Source kind (OCIRepository, GitRepository, Bucket)
 * @returns {string} Tailwind CSS classes
 */
export function getSourceKindBadgeClass(kind) {
  switch (kind) {
  case 'OCIRepository':
    return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400'
  case 'GitRepository':
    return 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-400'
  case 'Bucket':
    return 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-400'
  default:
    return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400'
  }
}
