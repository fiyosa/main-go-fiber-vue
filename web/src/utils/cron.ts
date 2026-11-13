/**
 * Builds a human-readable explanation for a cron schedule expression.
 * Supports both 5-field (minute hour dom month dow) and 6-field
 * (second minute hour dom month dow) expressions.
 */
export function describeCron(schedule: string): string {
  const raw = schedule.trim()
  if (!raw) return 'Custom schedule.'

  const parts = raw.split(/\s+/)
  if (parts.length === 6) {
    parts.shift()
  }

  const min = parts[0] || '*'
  const hr = parts[1] || '*'
  const dom = parts[2] || '*'
  const mon = parts[3] || '*'
  const dow = parts[4] || '*'

  if (min === '*' && hr === '*' && dom === '*' && mon === '*' && dow === '*') {
    return 'Runs every minute of every day.'
  }
  if (min.startsWith('*/') && hr === '*' && dom === '*' && mon === '*' && dow === '*') {
    return `Runs every ${min.replace('*/', '')} minutes.`
  }
  if (min === '0' && hr === '*' && dom === '*' && mon === '*' && dow === '*') {
    return 'Runs every hour, at minute 0.'
  }
  if (min === '0' && hr === '0' && dom === '*' && mon === '*' && dow === '*') {
    return 'Runs every day at 00:00 (midnight).'
  }
  if (/^\d+$/.test(min) && hr === '*' && dom === '*' && mon === '*' && dow === '*') {
    return `Runs every hour, at minute ${min}.`
  }
  if (/^\d+$/.test(min) && /^\d+$/.test(hr) && dom === '*' && mon === '*' && dow === '*') {
    return `Runs every day at ${hr.padStart(2, '0')}:${min.padStart(2, '0')}.`
  }
  return `Custom schedule: ${raw}`
}
