// Copyright 2024 The Flux authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

import { useLocation } from 'preact-iso'
import { fetchReport } from '../../app'
import { ThemeToggle } from './ThemeToggle'
import { TofuControllerIcon } from './Icons'

/**
 * Header component - Main application header with navigation and controls
 *
 * Features:
 * - Tofu Controller branding
 * - Refresh button (click logo when on home)
 * - Theme toggle
 * - Responsive design
 */
export function Header() {
  const location = useLocation()
  const currentPath = location.path

  const handleLogoClick = (e) => {
    if (currentPath === '/') {
      e.preventDefault()
      e.stopPropagation()
      fetchReport()
    }
  }

  return (
    <header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm transition-colors">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-2 sm:py-3">
        <div class="flex items-center justify-between">
          {/* Left: Logo and title */}
          <div class="flex items-center gap-3 sm:gap-4">
            <a
              href="/"
              onClick={handleLogoClick}
              class="flex items-center space-x-2 sm:space-x-3 hover:opacity-80 transition-opacity focus:outline-none"
              aria-label="Tofu Controller"
            >
              <TofuControllerIcon className="w-7 h-7 sm:w-8 sm:h-8 text-tofu-purple" />
              <h1 class="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white">Tofu Controller</h1>
            </a>
          </div>

          {/* Right: Theme toggle */}
          <div class="flex items-center space-x-2">
            <ThemeToggle />
          </div>
        </div>
      </div>
    </header>
  )
}
