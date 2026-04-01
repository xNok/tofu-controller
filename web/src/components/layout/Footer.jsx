
/**
 * Footer component - Application footer with links and license information
 */
export function Footer() {
  return (
    <footer class="bg-white dark:bg-gray-800 sm:border-t border-gray-200 dark:border-gray-700 transition-colors mt-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 sm:py-4">
        {/* Mobile: centered, minimal footer */}
        <div class="flex sm:hidden items-center justify-center gap-4 text-xs text-gray-500 dark:text-gray-400">
          <a
            href="https://github.com/flux-iac/tofu-controller"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-1.5 hover:text-gray-900 dark:hover:text-white transition-colors"
          >
            <svg class="w-3.5 h-3.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
            <span>Source code</span>
          </a>
          <span class="text-gray-300 dark:text-gray-600">•</span>
          <a
            href="https://flux-iac.github.io/tofu-controller/"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-1.5 hover:text-gray-900 dark:hover:text-white transition-colors"
          >
            <svg class="w-3.5 h-3.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>Documentation</span>
          </a>
        </div>

        {/* Desktop: full footer */}
        <div class="hidden sm:flex flex-row items-center justify-between gap-4">
          <div class="flex flex-row items-center gap-6 text-sm">
            <a
              href="https://github.com/flux-iac/tofu-controller"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              <svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              <span>Source code</span>
            </a>
            <span class="text-gray-300 dark:text-gray-600">•</span>
            <a
              href="https://flux-iac.github.io/tofu-controller/"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              <svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span>Documentation</span>
            </a>
          </div>
          <div class="text-sm text-gray-600 dark:text-gray-400 text-right">
            <p>Apache-2.0 Licensed</p>
          </div>
        </div>
      </div>
    </footer>
  )
}
