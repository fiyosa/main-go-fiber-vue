/// Formats an ISO timestamp as `YYYY-MM-DD HH:mm:ss` in local time.
export const formatDateTime = (iso?: string): string => {
  if (!iso) return '-'

  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso

  const pad = (value: number) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

/// Renders a short relative label such as `30 min ago` or `2 days ago`.
export const formatTimeAgo = (iso: string, now: number = Date.now()): string => {
  const then = new Date(iso).getTime()
  if (Number.isNaN(then)) return iso

  const minutes = Math.floor(Math.max(0, now - then) / 60_000)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes} min ago`

  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} hr ago`

  const days = Math.floor(hours / 24)
  return `${days} day${days === 1 ? '' : 's'} ago`
}

/// Renders a millisecond duration in the most readable unit.
export const formatDuration = (ms: number): string => {
  if (ms < 1000) return `${ms} ms`

  const seconds = ms / 1000
  if (seconds < 60) return `${seconds.toFixed(1)} s`

  const minutes = Math.floor(seconds / 60)
  const rest = Math.round(seconds % 60)
  return `${minutes}m ${rest}s`
}
