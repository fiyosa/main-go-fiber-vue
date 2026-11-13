<script setup lang="ts">
  import { ref, onMounted, computed, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import {
    listCrons,
    deleteCron,
    toggleCron,
    triggerCron,
    listCronRuns,
    listAllRuns,
    getRunDetail,
    type CronEntity,
    type CronRunEntity,
  } from '../../api/cron'
  import { describeCron } from '../../utils/cron'
  import { useUIStore } from '../../store/ui'
  import AdminPagination from '../../components/admin/ui/AdminPagination.vue'

  const router = useRouter()
  const ui = useUIStore()
  const crons = ref<CronEntity[]>([])
  const allRuns = ref<CronRunEntity[]>([])
  const activeTab = ref<'crons' | 'runs'>('crons')
  const loading = ref(false)
  const searchQuery = ref('')

  const filteredCrons = computed(() => {
    const q = searchQuery.value.trim().toLowerCase()
    if (!q) return crons.value
    return crons.value.filter(
      (c) =>
        c.name.toLowerCase().includes(q) ||
        c.schedule.toLowerCase().includes(q) ||
        (c.timezone || '').toLowerCase().includes(q)
    )
  })

  const filteredAllRuns = computed(() => {
    const q = searchQuery.value.trim().toLowerCase()
    if (!q) return allRuns.value
    return allRuns.value.filter(
      (r) =>
        (r.status || '').toLowerCase().includes(q) ||
        (r.trigger_type || '').toLowerCase().includes(q) ||
        (r.error || '').toLowerCase().includes(q) ||
        (r.output || '').toLowerCase().includes(q)
    )
  })

  // Runs modal state
  const showRunsModal = ref(false)
  const selectedCronName = ref('')
  const cronRuns = ref<CronRunEntity[]>([])
  const loadingRuns = ref(false)

  // Run detail modal state
  const showDetailModal = ref(false)
  const selectedRun = ref<CronRunEntity | null>(null)

  // Pagination state for Cron Jobs
  const cronPage = ref(1)
  const cronLimit = ref(25)
  const cronTotalPages = computed(() => Math.max(1, Math.ceil(filteredCrons.value.length / cronLimit.value)))
  const paginatedCrons = computed(() => {
    const page = Math.min(cronPage.value, cronTotalPages.value)
    const start = (page - 1) * cronLimit.value
    return filteredCrons.value.slice(start, start + cronLimit.value)
  })

  // Pagination state for Execution History
  const runsPage = ref(1)
  const runsLimit = ref(25)
  const runsTotalPages = computed(() => Math.max(1, Math.ceil(filteredAllRuns.value.length / runsLimit.value)))
  const paginatedAllRuns = computed(() => {
    const page = Math.min(runsPage.value, runsTotalPages.value)
    const start = (page - 1) * runsLimit.value
    return filteredAllRuns.value.slice(start, start + runsLimit.value)
  })

  // Pagination state for Runs History Modal
  const modalRunsPage = ref(1)
  const modalRunsLimit = ref(25)
  const modalRunsTotalPages = computed(() => Math.max(1, Math.ceil(cronRuns.value.length / modalRunsLimit.value)))
  const paginatedCronRuns = computed(() => {
    const page = Math.min(modalRunsPage.value, modalRunsTotalPages.value)
    const start = (page - 1) * modalRunsLimit.value
    return cronRuns.value.slice(start, start + modalRunsLimit.value)
  })

  watch(searchQuery, () => {
    cronPage.value = 1
    runsPage.value = 1
  })

  const loadData = async () => {
    loading.value = true
    try {
      const [cronList, runsList] = await Promise.all([listCrons(), listAllRuns()])
      crons.value = cronList
      allRuns.value = runsList
    } catch (err: any) {
      ui.showToast(err.message || 'Failed to load cron jobs.', 'error')
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    loadData()
  })

  const openCreateModal = () => {
    router.push({ name: 'admin-cron-create' })
  }

  const openEditModal = (cron: CronEntity) => {
    router.push({ name: 'admin-cron-edit', params: { id: cron.id } })
  }

  const handleToggle = async (id: string) => {
    try {
      await toggleCron(id)
      await loadData()
      ui.showToast('Cron status updated.', 'success')
    } catch (err: any) {
      ui.showToast(err.message || 'Failed to toggle cron.', 'error')
    }
  }

  const handleTrigger = async (cron: CronEntity) => {
    const ok = await ui.confirm({
      title: 'Run Cron Job',
      message: `Are you sure you want to run "${cron.name}" now?`,
      confirmText: 'Run Now',
      cancelText: 'Cancel',
      type: 'info',
    })
    if (!ok) return

    try {
      ui.showToast('Triggering cron job...', 'success')
      const run = await triggerCron(cron.id)
      await loadData()
      ui.showToast(`Cron triggered successfully (${run.status}, ${run.duration_ms}ms).`, 'success')
    } catch (err: any) {
      ui.showToast(err.message || 'Failed to trigger cron job.', 'error')
    }
  }

  const handleDelete = async (id: string) => {
    const ok = await ui.confirm({
      title: 'Delete Cron Job',
      message: 'Are you sure you want to delete this cron job and its run history? This action cannot be undone.',
      confirmText: 'Delete',
      cancelText: 'Cancel',
      type: 'danger',
    })
    if (!ok) return

    try {
      await deleteCron(id)
      await loadData()
      ui.showToast('Cron job deleted.', 'success')
    } catch (err: any) {
      ui.showToast(err.message || 'Failed to delete cron job.', 'error')
    }
  }

  const viewCronRuns = async (cron: CronEntity) => {
    selectedCronName.value = cron.name
    modalRunsPage.value = 1
    loadingRuns.value = true
    showRunsModal.value = true
    try {
      cronRuns.value = await listCronRuns(cron.id)
    } catch (err: any) {
      ui.showToast(err.message || 'Failed to load execution runs.', 'error')
    } finally {
      loadingRuns.value = false
    }
  }

  const viewRunDetail = async (run: CronRunEntity) => {
    try {
      const detail = await getRunDetail(run.id)
      selectedRun.value = detail
      showDetailModal.value = true
    } catch (err: any) {
      ui.showToast(err.message || 'Failed to load run detail.', 'error')
    }
  }

  const formatDate = (isoStr?: string) => {
    if (!isoStr) return '-'
    const d = new Date(isoStr)
    if (isNaN(d.getTime())) return isoStr
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
  }
</script>

<template>
  <div class="flex flex-col gap-8">
    <!-- Header Banner -->
    <div
      class="flex flex-col justify-between gap-4 rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm sm:flex-row sm:items-center"
    >
      <div class="flex flex-col gap-1">
        <div class="flex items-center gap-2">
          <p class="text-xs font-semibold tracking-widest text-[var(--admin-muted)] uppercase">Automation</p>
          <span class="rounded-full bg-blue-500/10 px-2 py-0.5 text-[10px] font-bold text-blue-500">Goja Engine</span>
        </div>
        <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Cron Jobs Manager</h1>
        <p class="text-sm text-[var(--admin-muted)]">
          Schedule and manage JavaScript background tasks with automated fetch support and encrypted environment
          variables.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          type="button"
          class="inline-flex min-h-11 items-center justify-center rounded-xl bg-blue-600 px-4 text-xs font-semibold text-white shadow-sm transition-colors hover:bg-blue-700"
          @click="openCreateModal"
        >
          + Create
        </button>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="relative w-full max-w-md">
        <svg
          class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-[var(--admin-muted)]"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M21 21l-4.35-4.35m0 0A7.5 7.5 0 105.1 5.1a7.5 7.5 0 0011.6 11.6z"
          />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="
            activeTab === 'crons'
              ? 'Search cron jobs by name, schedule...'
              : 'Search execution runs by status, trigger...'
          "
          class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] py-2 pr-4 pl-9 text-xs transition-colors outline-none focus:border-blue-500"
        />
      </div>

      <button
        type="button"
        class="inline-flex min-h-11 shrink-0 items-center justify-center gap-2 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 text-xs font-semibold text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50 disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="loading"
        @click="loadData()"
      >
        <svg
          class="size-4"
          :class="{ 'animate-spin': loading }"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        Refresh
      </button>
    </div>

    <!-- Tabs -->
    <div class="flex border-b border-[var(--admin-border)]">
      <button
        type="button"
        class="border-b-2 px-5 py-3 text-xs font-bold tracking-wider uppercase transition-colors"
        :class="
          activeTab === 'crons'
            ? 'border-blue-600 text-blue-600 dark:text-blue-400'
            : 'border-transparent text-[var(--admin-muted)] hover:text-[var(--admin-fg)]'
        "
        @click="activeTab = 'crons'"
      >
        Cron Jobs ({{ filteredCrons.length
        }}<span v-if="filteredCrons.length !== crons.length"> / {{ crons.length }}</span
        >)
      </button>
      <button
        type="button"
        class="border-b-2 px-5 py-3 text-xs font-bold tracking-wider uppercase transition-colors"
        :class="
          activeTab === 'runs'
            ? 'border-blue-600 text-blue-600 dark:text-blue-400'
            : 'border-transparent text-[var(--admin-muted)] hover:text-[var(--admin-fg)]'
        "
        @click="activeTab = 'runs'"
      >
        Execution History ({{ filteredAllRuns.length
        }}<span v-if="filteredAllRuns.length !== allRuns.length"> / {{ allRuns.length }}</span
        >)
      </button>
    </div>

    <!-- Tab 1: Crons List -->
    <div v-if="activeTab === 'crons'" class="flex flex-col gap-4">
      <div v-if="loading" class="py-12 text-center text-xs text-[var(--admin-muted)]">Loading cron jobs...</div>
      <div
        v-else-if="filteredCrons.length === 0"
        class="flex flex-col items-center justify-center rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-12 text-center"
      >
        <p class="text-sm font-semibold text-[var(--admin-fg)]">No cron jobs found</p>
        <p class="mt-1 text-xs text-[var(--admin-muted)]">
          {{
            crons.length === 0
              ? 'Create your first scheduled JS background task to get started.'
              : 'Try adjusting your search query.'
          }}
        </p>
        <button
          v-if="crons.length === 0"
          type="button"
          class="mt-4 inline-flex min-h-10 items-center justify-center rounded-xl bg-blue-600 px-4 text-xs font-semibold text-white"
          @click="openCreateModal"
        >
          Create Cron Job
        </button>
      </div>

      <div v-else class="flex flex-col gap-5">
        <div
          v-for="cron in paginatedCrons"
          :key="cron.id"
          class="flex flex-col justify-between gap-4 rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm transition-all hover:shadow-md sm:flex-row sm:items-center"
        >
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-3">
              <h3 class="text-base font-bold tracking-tight text-[var(--admin-fg)]">{{ cron.name }}</h3>
              <span
                class="rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase"
                :class="
                  cron.enabled
                    ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                    : 'bg-slate-500/10 text-[var(--admin-muted)]'
                "
              >
                {{ cron.enabled ? 'Active' : 'Paused' }}
              </span>
              <span class="rounded-full bg-blue-500/10 px-2.5 py-0.5 text-[10px] font-semibold text-blue-500">
                {{ cron.overlap_policy }} overlap
              </span>
            </div>

            <div class="flex flex-wrap items-center gap-4 text-xs text-[var(--admin-muted)]">
              <div class="flex items-center gap-1.5 font-mono">
                <svg
                  class="size-4 text-blue-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <span>( {{ cron.schedule }} )</span>
              </div>
              <span>•</span>
              <span class="font-medium text-blue-500">({{ describeCron(cron.schedule) }})</span>
              <span>•</span>
              <span>TZ: {{ cron.timezone }}</span>
              <span>•</span>
              <span>Updated: {{ formatDate(cron.updated_at) }}</span>
            </div>
          </div>

          <div class="flex flex-wrap items-center gap-2">
            <button
              type="button"
              class="inline-flex min-h-9 min-w-9 items-center justify-center rounded-xl border border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 transition-colors hover:bg-emerald-500/20"
              @click="handleTrigger(cron)"
              title="Run Now"
              aria-label="Run Now"
            >
              <svg class="size-4" viewBox="0 0 24 24" fill="currentColor">
                <path d="M8 5v14l11-7z" />
              </svg>
            </button>
            <div class="h-6 w-px bg-[var(--admin-border)]" aria-hidden="true"></div>
            <button
              type="button"
              class="inline-flex min-h-9 min-w-9 items-center justify-center rounded-xl border border-blue-500/30 bg-blue-500/10 text-blue-600 dark:text-blue-400 transition-colors hover:bg-blue-500/20"
              @click="openEditModal(cron)"
              title="Edit"
              aria-label="Edit"
            >
              <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path
                  d="M4 20h4L19.5 8.5a2.121 2.121 0 00-3-3L5 17v3z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </button>
            <button
              type="button"
              class="inline-flex min-h-9 min-w-9 items-center justify-center rounded-xl border border-purple-500/30 bg-purple-500/10 text-purple-600 dark:text-purple-400 transition-colors hover:bg-purple-500/20"
              @click="viewCronRuns(cron)"
              title="Runs History"
              aria-label="Runs History"
            >
              <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="9" />
                <path d="M12 7v5l3 3" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </button>
            <div class="h-6 w-px bg-[var(--admin-border)]" aria-hidden="true"></div>
            <button
              type="button"
              class="inline-flex min-h-9 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3 text-xs font-semibold transition-colors"
              :class="
                cron.enabled ? 'text-amber-600 hover:bg-amber-500/10' : 'text-emerald-600 hover:bg-emerald-500/10'
              "
              @click="handleToggle(cron.id)"
            >
              {{ cron.enabled ? 'Pause' : 'Enable' }}
            </button>
            <button
              type="button"
              class="inline-flex min-h-9 items-center justify-center rounded-xl border border-red-500/30 bg-red-500/10 px-3 text-xs font-semibold text-red-500 transition-colors hover:bg-red-500/20"
              @click="handleDelete(cron.id)"
            >
              Delete
            </button>
          </div>
        </div>

        <!-- Cron Pagination Controls -->
        <AdminPagination
          v-if="filteredCrons.length > 0"
          v-model:page="cronPage"
          v-model:limit="cronLimit"
          :total="filteredCrons.length"
        />
      </div>
    </div>

    <!-- Tab 2: All Runs History -->
    <div v-if="activeTab === 'runs'" class="flex flex-col gap-4">
      <div v-if="loading" class="py-12 text-center text-xs text-[var(--admin-muted)]">Loading execution history...</div>
      <div
        v-else-if="filteredAllRuns.length === 0"
        class="flex flex-col items-center justify-center rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-12 text-center"
      >
        <p class="text-sm font-semibold text-[var(--admin-fg)]">No execution runs found</p>
        <p class="mt-1 text-xs text-[var(--admin-muted)]">
          {{
            allRuns.length === 0
              ? 'Cron execution logs and console outputs will appear here.'
              : 'Try adjusting your search query.'
          }}
        </p>
      </div>
      <div v-else class="flex flex-col gap-4">
        <div class="overflow-x-auto rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)]">
          <table class="w-full text-left text-xs">
            <thead
              class="border-b border-[var(--admin-border)] bg-[var(--admin-border)]/20 text-[var(--admin-muted)] uppercase"
            >
              <tr>
                <th class="px-4 py-3 font-semibold">Status</th>
                <th class="px-4 py-3 font-semibold">Trigger</th>
                <th class="px-4 py-3 font-semibold">Started At</th>
                <th class="px-4 py-3 font-semibold">Duration</th>
                <th class="px-4 py-3 font-semibold">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[var(--admin-border)]">
              <tr
                v-for="run in paginatedAllRuns"
                :key="run.id"
                class="transition-colors hover:bg-[var(--admin-border)]/10"
              >
                <td class="px-4 py-3.5">
                  <span
                    class="rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase"
                    :class="{
                      'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400': run.status === 'success',
                      'bg-rose-500/10 text-rose-600 dark:text-rose-400': run.status === 'failed',
                      'bg-amber-500/10 text-amber-600 dark:text-amber-400': run.status === 'running',
                    }"
                  >
                    {{ run.status }}
                  </span>
                </td>
                <td class="px-4 py-3.5 font-medium capitalize">{{ run.trigger_type }}</td>
                <td class="px-4 py-3.5 text-[var(--admin-muted)]">{{ formatDate(run.started_at) }}</td>
                <td class="px-4 py-3.5 font-mono">{{ run.duration_ms }} ms</td>
                <td class="px-4 py-3.5">
                  <button
                    type="button"
                    class="rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 py-1 text-xs font-semibold hover:bg-[var(--admin-border)]/50"
                    @click="viewRunDetail(run)"
                  >
                    View Output
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Runs Pagination Controls -->
        <AdminPagination
          v-if="filteredAllRuns.length > 0"
          v-model:page="runsPage"
          v-model:limit="runsLimit"
          :total="filteredAllRuns.length"
        />
      </div>
    </div>

    <!-- Runs History Modal for specific Cron -->
    <div
      v-if="showRunsModal"
      class="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-slate-950/60 p-4 backdrop-blur-xs"
    >
      <div
        class="relative my-8 w-full max-w-4xl rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-2xl"
      >
        <div class="flex items-center justify-between border-b border-[var(--admin-border)] pb-4">
          <h2 class="text-base font-bold tracking-tight">Execution Runs: {{ selectedCronName }}</h2>
          <button
            type="button"
            class="text-[var(--admin-muted)] hover:text-[var(--admin-fg)]"
            @click="showRunsModal = false"
          >
            ✕
          </button>
        </div>

        <div class="mt-4 flex flex-col gap-4 max-h-[65vh] overflow-y-auto">
          <div v-if="loadingRuns" class="py-8 text-center text-xs text-[var(--admin-muted)]">
            Loading runs history...
          </div>
          <div v-else-if="cronRuns.length === 0" class="py-8 text-center text-xs text-[var(--admin-muted)]">
            No execution runs found for this job.
          </div>
          <template v-else>
            <div class="overflow-x-auto">
              <table class="w-full text-left text-xs">
                <thead
                  class="border-b border-[var(--admin-border)] bg-[var(--admin-border)]/20 text-[var(--admin-muted)] uppercase"
                >
                  <tr>
                    <th class="px-3 py-2.5 font-semibold">Status</th>
                    <th class="px-3 py-2.5 font-semibold">Trigger</th>
                    <th class="px-3 py-2.5 font-semibold">Started At</th>
                    <th class="px-3 py-2.5 font-semibold">Duration</th>
                    <th class="px-3 py-2.5 font-semibold">Output / Error</th>
                    <th class="px-3 py-2.5 font-semibold">Action</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-[var(--admin-border)]">
                  <tr v-for="run in paginatedCronRuns" :key="run.id" class="hover:bg-[var(--admin-border)]/10">
                    <td class="px-3 py-3">
                      <span
                        class="rounded-full px-2 py-0.5 text-[10px] font-bold uppercase"
                        :class="{
                          'bg-emerald-500/10 text-emerald-600': run.status === 'success',
                          'bg-rose-500/10 text-rose-600': run.status === 'failed',
                          'bg-amber-500/10 text-amber-600': run.status === 'running',
                        }"
                      >
                        {{ run.status }}
                      </span>
                    </td>
                    <td class="px-3 py-3 capitalize">{{ run.trigger_type }}</td>
                    <td class="px-3 py-3 text-[var(--admin-muted)]">{{ formatDate(run.started_at) }}</td>
                    <td class="px-3 py-3 font-mono">{{ run.duration_ms }} ms</td>
                    <td class="max-w-xs truncate px-3 py-3 font-mono text-[11px] text-[var(--admin-muted)]">
                      {{ run.error || run.output || '-' }}
                    </td>
                    <td class="px-3 py-3">
                      <button
                        type="button"
                        class="rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 py-1 text-xs font-semibold hover:bg-[var(--admin-border)]/50"
                        @click="viewRunDetail(run)"
                      >
                        Details
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Modal Runs Pagination Controls -->
            <AdminPagination
              v-if="cronRuns.length > 0"
              v-model:page="modalRunsPage"
              v-model:limit="modalRunsLimit"
              :total="cronRuns.length"
            />
          </template>
        </div>
      </div>
    </div>

    <!-- Run Detail / Output Modal -->
    <div
      v-if="showDetailModal && selectedRun"
      class="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-slate-950/60 p-4 backdrop-blur-xs"
    >
      <div
        class="relative my-8 w-full max-w-2xl rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-2xl"
      >
        <div class="flex items-center justify-between border-b border-[var(--admin-border)] pb-4">
          <div class="flex items-center gap-3">
            <h2 class="text-base font-bold tracking-tight">Execution Run Details</h2>
            <span
              class="rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase"
              :class="{
                'bg-emerald-500/10 text-emerald-600': selectedRun.status === 'success',
                'bg-rose-500/10 text-rose-600': selectedRun.status === 'failed',
              }"
            >
              {{ selectedRun.status }}
            </span>
          </div>
          <button
            type="button"
            class="text-[var(--admin-muted)] hover:text-[var(--admin-fg)]"
            @click="showDetailModal = false"
          >
            ✕
          </button>
        </div>

        <div class="mt-4 flex flex-col gap-4 text-xs">
          <div
            class="grid grid-cols-2 gap-4 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-border)]/10 p-3"
          >
            <div>
              <span class="text-[var(--admin-muted)]">Trigger Type:</span>
              <strong class="capitalize">{{ selectedRun.trigger_type }}</strong>
            </div>
            <div>
              <span class="text-[var(--admin-muted)]">Duration:</span>
              <strong class="font-mono">{{ selectedRun.duration_ms }} ms</strong>
            </div>
            <div>
              <span class="text-[var(--admin-muted)]">Started At:</span> {{ formatDate(selectedRun.started_at) }}
            </div>
            <div>
              <span class="text-[var(--admin-muted)]">Finished At:</span> {{ formatDate(selectedRun.finished_at) }}
            </div>
          </div>

          <div v-if="selectedRun.error" class="flex flex-col gap-1.5">
            <label class="font-semibold text-rose-500">Error Message</label>
            <pre
              class="rounded-xl border border-rose-500/30 bg-rose-500/10 p-3 font-mono text-xs whitespace-pre-wrap text-rose-600 dark:text-rose-400"
              >{{ selectedRun.error }}</pre>
          </div>

          <div class="flex flex-col gap-1.5">
            <label class="font-semibold">Console Output (`console.log` capture)</label>
            <pre
              class="max-h-60 overflow-y-auto rounded-xl border border-[var(--admin-border)] bg-slate-950 p-4 font-mono text-xs whitespace-pre-wrap text-emerald-400"
              >{{ selectedRun.output || '(No console output)' }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
