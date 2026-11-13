import { ref, computed } from 'vue'
import { listCrons, listAllRuns, type CronEntity, type CronRunEntity } from '../../api/cron'
import { getHealth, getReady } from '../../api/health'
import { getMe, type UserMe } from '../../api/auth'
import { summarizeRuns, countEnabledCrons, sortRunsByStartedAtDesc } from '../../utils/dashboard'

export const RUN_WINDOW_MS = 24 * 60 * 60 * 1000

export function useDashboard() {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const healthStatus = ref<string | null>(null)
  const readyStatus = ref<string | null>(null)
  const crons = ref<CronEntity[]>([])
  const runs = ref<CronRunEntity[]>([])
  const currentUser = ref<UserMe | null>(null)
  const lastUpdated = ref<Date | null>(null)

  const cronStats = computed(() => ({
    total: crons.value.length,
    enabled: countEnabledCrons(crons.value),
    disabled: crons.value.length - countEnabledCrons(crons.value),
  }))

  const runsStats = computed(() => summarizeRuns(runs.value, RUN_WINDOW_MS))
  const sortedRuns = computed(() => sortRunsByStartedAtDesc(runs.value))
  const failedRuns = computed(() => sortedRuns.value.filter((r) => r.status === 'failed'))
  const recentRuns = computed(() => sortedRuns.value.slice(0, 10))

  const loadDashboard = async () => {
    loading.value = true
    error.value = null

    try {
      const [health, ready, cronList, runList, me] = await Promise.all([
        getHealth(),
        getReady(),
        listCrons(),
        listAllRuns(),
        getMe(),
      ])
      healthStatus.value = health.status
      readyStatus.value = ready.status
      crons.value = cronList
      runs.value = runList
      currentUser.value = me
      lastUpdated.value = new Date()
    } catch (err: any) {
      error.value = err?.message || 'Failed to load dashboard data.'
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    healthStatus,
    readyStatus,
    crons,
    runs,
    currentUser,
    lastUpdated,
    cronStats,
    runsStats,
    recentRuns,
    failedRuns,
    loadDashboard,
  }
}
