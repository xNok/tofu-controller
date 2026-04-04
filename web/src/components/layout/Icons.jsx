
/**
 * TofuControllerIcon - OpenTofu-inspired cube icon
 * Uses currentColor so it adapts to text color
 */
export function TofuControllerIcon({ className }) {
  return (
    <svg class={className} viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M32 8L56 20V44L32 56L8 44V20L32 8Z" stroke="currentColor" stroke-width="3" stroke-linejoin="round" fill="none"/>
      <path d="M32 8V56" stroke="currentColor" stroke-width="3"/>
      <path d="M8 20L56 20" stroke="currentColor" stroke-width="3"/>
      <path d="M8 44L32 32L56 44" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>
      <path d="M32 32L8 20" stroke="currentColor" stroke-width="3"/>
      <path d="M32 32L56 20" stroke="currentColor" stroke-width="3"/>
    </svg>
  )
}
