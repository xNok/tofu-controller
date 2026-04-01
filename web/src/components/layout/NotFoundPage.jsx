// Copyright 2024 The Flux authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

/**
 * NotFoundPage - 404 error page for unmatched routes
 */
export function NotFoundPage() {
  return (
    <main data-testid="not-found-page" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-grow w-full">
      <div class="flex items-center justify-center min-h-[60vh]">
        <div class="text-center">
          <div class="mx-auto w-20 h-20 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center mb-6">
            <svg class="w-12 h-12 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>

          <h1 class="text-2xl sm:text-3xl font-semibold text-gray-900 dark:text-white mb-2">
            Page Not Found
          </h1>

          <p class="text-gray-500 dark:text-gray-400 mb-8 max-w-md mx-auto">
            The page you&#39;re looking for doesn&#39;t exist or has been moved.
          </p>

          <a
            href="/"
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-tofu-purple hover:opacity-90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-tofu-purple transition-colors"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            Go to Dashboard
          </a>
        </div>
      </div>
    </main>
  )
}
