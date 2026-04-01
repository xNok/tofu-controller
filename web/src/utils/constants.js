// Copyright 2024 The Flux authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Auto-refresh polling interval in milliseconds
export const POLL_INTERVAL_MS = 30000

// Fast polling interval when a plan is pending or applying (5 seconds)
export const FAST_POLL_INTERVAL_MS = 5000

// Timeout to revert to normal polling after plan completion (5 minutes)
export const FAST_POLL_TIMEOUT_MS = 300000

// Terraform resource status values
export const resourceStatuses = [
  { value: '', label: 'All' },
  { value: 'Ready', label: 'Ready' },
  { value: 'Failed', label: 'Failed' },
  { value: 'Planning', label: 'Planning' },
  { value: 'Applying', label: 'Applying' },
  { value: 'PendingApproval', label: 'Pending Approval' },
  { value: 'Suspended', label: 'Suspended' },
  { value: 'Unknown', label: 'Unknown' },
]

// Terraform source kinds
export const sourceKinds = [
  { value: '', label: 'All Sources' },
  { value: 'OCIRepository', label: 'OCIRepository' },
  { value: 'GitRepository', label: 'GitRepository' },
  { value: 'Bucket', label: 'Bucket' },
]

// Plan approval values
export const APPROVE_AUTO = 'auto'
export const APPROVE_DISABLE = 'disable'
