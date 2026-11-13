<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useUIStore } from '../../store/ui'
  import {
    createCron,
    updateCron,
    getCron,
    type CronCreateRequest,
    type CronUpdateRequest,
  } from '../../api/cron'
  import { describeCron } from '../../utils/cron'

  const route = useRoute()
  const router = useRouter()
  const ui = useUIStore()

  const isEdit = computed(() => Boolean(route.params.id))
  const cronId = computed(() => route.params.id as string | undefined)

  const formName = ref('')
  const cronMinute = ref('0')
  const cronHour = ref('*')
  const cronDom = ref('*')
  const cronMonth = ref('*')
  const cronDow = ref('*')

  const formSchedule = computed(() => {
    return `${cronMinute.value} ${cronHour.value} ${cronDom.value} ${cronMonth.value} ${cronDow.value}`
  })

  const formTimezone = ref('Asia/Jakarta')
  const formOverlap = ref<'skip' | 'allow'>('skip')
  const formEnabled = ref(true)
  const formScript = ref(`async function scheduled(event, env) {
  console.log("Running scheduled task...");
  const res = await fetch("https://httpbin.org/json");
  const data = await res.json();
  console.log("Fetched data:", JSON.stringify(data));
}`)
  const envPairs = ref<Array<{ key: string; value: string; show?: boolean }>>([])
  
  const loading = ref(false)
  const submitting = ref(false)

  onMounted(async () => {
    if (isEdit.value && cronId.value) {
      loading.value = true
      try {
        const cron = await getCron(cronId.value)
        formName.value = cron.name
        const parts = cron.schedule.trim().split(/\s+/)
        if (parts.length === 5) {
          cronMinute.value = parts[0]
          cronHour.value = parts[1]
          cronDom.value = parts[2]
          cronMonth.value = parts[3]
          cronDow.value = parts[4]
        }
        formTimezone.value = cron.timezone || 'Asia/Jakarta'
        formOverlap.value = cron.overlap_policy || 'skip'
        formEnabled.value = cron.enabled
        formScript.value = cron.script
        envPairs.value = cron.env
          ? Object.entries(cron.env).map(([key, value]) => ({ key, value, show: false }))
          : []
      } catch (err: any) {
        ui.showToast(err.message || 'Failed to load cron job detail.', 'error')
      } finally {
        loading.value = false
      }
    }
  })

  const addEnvPair = () => {
    envPairs.value.push({ key: '', value: '', show: !isEdit.value })
  }

  const removeEnvPair = (index: number) => {
    envPairs.value.splice(index, 1)
  }

  const handleSave = async () => {
    if (!formName.value.trim() || !formSchedule.value.trim() || !formScript.value.trim()) {
      ui.showToast('Name, schedule, and script are required.', 'error')
      return
    }

    submitting.value = true

    const envMap: Record<string, string> = {}
    for (const pair of envPairs.value) {
      if (pair.key.trim()) {
        envMap[pair.key.trim()] = pair.value
      }
    }

    const payload: CronCreateRequest | CronUpdateRequest = {
      name: formName.value.trim(),
      schedule: formSchedule.value.trim(),
      timezone: formTimezone.value.trim() || 'Asia/Jakarta',
      script: formScript.value,
      env: Object.keys(envMap).length > 0 ? envMap : undefined,
      overlap_policy: formOverlap.value,
      enabled: formEnabled.value,
    }

    try {
      if (isEdit.value && cronId.value) {
        await updateCron(cronId.value, payload as CronUpdateRequest)
      } else {
        await createCron(payload as CronCreateRequest)
      }
      ui.showToast(isEdit.value ? 'Cron job updated successfully.' : 'Cron job created successfully.', 'success')
      router.push({ name: 'admin-cron' })
    } catch (err: any) {
      ui.showToast(err.message || 'Failed to save cron job.', 'error')
    } finally {
      submitting.value = false
    }
  }

  const applyPreset = (expr: string) => {
    const parts = expr.trim().split(/\s+/)
    if (parts.length === 5) {
      cronMinute.value = parts[0]
      cronHour.value = parts[1]
      cronDom.value = parts[2]
      cronMonth.value = parts[3]
      cronDow.value = parts[4]
    }
  }

  const cronDescription = computed(() => describeCron(formSchedule.value))
</script>

<template>
  <div class="flex flex-col gap-6 max-w-5xl mx-auto pb-12">
    <!-- Top Nav / Header -->
    <div class="flex flex-col gap-2">
      <button
        type="button"
        class="inline-flex items-center gap-1.5 text-xs font-semibold text-[var(--admin-muted)] hover:text-[var(--admin-fg)] w-fit"
        @click="router.push({ name: 'admin-cron' })"
      >
        ← Back to Cron Jobs
      </button>
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold tracking-tight">{{ isEdit ? 'Edit Cron Job' : 'Create Cron Job' }}</h1>
          <p class="text-sm text-[var(--admin-muted)]">
            Configure schedule, timezone, encrypted secrets, and Goja JavaScript worker code.
          </p>
        </div>
        <div class="flex items-center gap-3">
          <button
            type="button"
            class="inline-flex min-h-10 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 text-xs font-semibold text-[var(--admin-fg)] hover:bg-[var(--admin-border)]/50"
            @click="router.push({ name: 'admin-cron' })"
          >
            Cancel
          </button>
          <button
            type="button"
            class="inline-flex min-h-10 items-center justify-center rounded-xl bg-blue-600 px-5 text-xs font-semibold text-white shadow-sm hover:bg-blue-700 disabled:opacity-60"
            :disabled="submitting || loading"
            @click="handleSave"
          >
            {{ submitting ? 'Saving...' : (isEdit ? 'Save Changes' : 'Create Cron Job') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="py-16 text-center text-xs text-[var(--admin-muted)]">Loading cron configuration...</div>

    <div v-else class="flex flex-col gap-6">
      <!-- Basic Configuration Card -->
      <div class="rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm flex flex-col gap-5">
        <h2 class="text-base font-bold tracking-tight border-b border-[var(--admin-border)] pb-3">1. General Configuration</h2>

        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-semibold">Job Name <span class="text-red-500">*</span></label>
            <input
              v-model="formName"
              type="text"
              placeholder="e.g. Hourly Exchange Rates Sync"
              class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3.5 text-xs text-[var(--admin-fg)] outline-none focus:border-blue-500"
            />
          </div>

          <div class="flex flex-col gap-1.5 sm:col-span-2">
            <div class="flex items-center justify-between">
              <label class="text-xs font-semibold">Schedule Expression (5-field Cron) <span class="text-red-500">*</span></label>
              <div class="flex items-center gap-1.5 text-[11px]">
                <span class="text-[var(--admin-muted)]">Presets:</span>
                <button type="button" class="text-blue-500 hover:underline" @click="applyPreset('0 * * * *')">Hourly</button>
                <span>•</span>
                <button type="button" class="text-blue-500 hover:underline" @click="applyPreset('0 0 * * *')">Daily</button>
                <span>•</span>
                <button type="button" class="text-blue-500 hover:underline" @click="applyPreset('*/15 * * * *')">Every 15m</button>
              </div>
            </div>
            
            <div class="grid grid-cols-5 gap-2">
              <div class="flex flex-col gap-1">
                <input
                  v-model="cronMinute"
                  type="text"
                  placeholder="*"
                  class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 text-center text-xs font-mono text-[var(--admin-fg)] outline-none focus:border-blue-500"
                />
                <span class="text-[10px] text-center text-[var(--admin-muted)]">Minute (0-59)</span>
              </div>
              <div class="flex flex-col gap-1">
                <input
                  v-model="cronHour"
                  type="text"
                  placeholder="*"
                  class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 text-center text-xs font-mono text-[var(--admin-fg)] outline-none focus:border-blue-500"
                />
                <span class="text-[10px] text-center text-[var(--admin-muted)]">Hour (0-23)</span>
              </div>
              <div class="flex flex-col gap-1">
                <input
                  v-model="cronDom"
                  type="text"
                  placeholder="*"
                  class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 text-center text-xs font-mono text-[var(--admin-fg)] outline-none focus:border-blue-500"
                />
                <span class="text-[10px] text-center text-[var(--admin-muted)]">Day (1-31)</span>
              </div>
              <div class="flex flex-col gap-1">
                <input
                  v-model="cronMonth"
                  type="text"
                  placeholder="*"
                  class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 text-center text-xs font-mono text-[var(--admin-fg)] outline-none focus:border-blue-500"
                />
                <span class="text-[10px] text-center text-[var(--admin-muted)]">Month (1-12)</span>
              </div>
              <div class="flex flex-col gap-1">
                <input
                  v-model="cronDow"
                  type="text"
                  placeholder="*"
                  class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 text-center text-xs font-mono text-[var(--admin-fg)] outline-none focus:border-blue-500"
                />
                <span class="text-[10px] text-center text-[var(--admin-muted)]">Day of Week</span>
              </div>
            </div>
            <div class="mt-2 flex flex-col sm:flex-row sm:items-center sm:justify-between px-1 text-[11px] gap-1">
              <span class="text-[var(--admin-muted)]">Combined Cron: <code class="font-mono font-semibold text-[var(--admin-fg)]">{{ formSchedule }}</code></span>
              <span class="font-medium text-blue-500 bg-blue-500/10 px-2 py-0.5 rounded-md w-fit">💡 {{ cronDescription }}</span>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-5 sm:grid-cols-3">
          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-semibold">Timezone</label>
            <select
              v-model="formTimezone"
              class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3.5 text-xs text-[var(--admin-fg)] outline-none focus:border-blue-500"
            >
              <option value="Asia/Jakarta">Asia/Jakarta (WIB)</option>
              <option value="Asia/Singapore">Asia/Singapore (SGT)</option>
              <option value="Asia/Tokyo">Asia/Tokyo (JST)</option>
              <option value="UTC">UTC</option>
              <option value="America/New_York">America/New_York (EST/EDT)</option>
              <option value="Europe/London">Europe/London (GMT/BST)</option>
              <option value="Australia/Sydney">Australia/Sydney (AEST/AEDT)</option>
            </select>
          </div>

          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-semibold">Overlap Policy</label>
            <select
              v-model="formOverlap"
              class="min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3.5 text-xs text-[var(--admin-fg)] outline-none focus:border-blue-500"
            >
              <option value="skip">skip (prevent concurrent runs)</option>
              <option value="allow">allow concurrent runs</option>
            </select>
          </div>

          <div class="flex flex-col justify-end">
            <label class="flex items-center gap-2.5 cursor-pointer min-h-11 py-2.5">
              <input v-model="formEnabled" type="checkbox" class="size-4 rounded border-[var(--admin-border)] text-blue-600 focus:ring-blue-500" />
              <span class="text-xs font-semibold">Enabled (Active Scheduler)</span>
            </label>
          </div>
        </div>
      </div>

      <!-- Environment Variables Card -->
      <div class="rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-[var(--admin-border)] pb-3">
          <div>
            <h2 class="text-base font-bold tracking-tight">2. Encrypted Environment Variables (env)</h2>
            <p class="text-xs text-[var(--admin-muted)]">Secrets are encrypted via AES-256-GCM before storing in database.</p>
          </div>
          <button type="button" class="inline-flex items-center gap-1 rounded-xl bg-blue-500/10 px-3 py-1.5 text-xs font-semibold text-blue-500 hover:bg-blue-500/20" @click="addEnvPair">
            + Add Variable
          </button>
        </div>

        <div v-if="envPairs.length === 0" class="text-xs text-[var(--admin-muted)] italic py-2">No environment variables defined. Accessible in script via `env.KEY`.</div>
        
        <div v-for="(pair, idx) in envPairs" :key="idx" class="flex items-center gap-3">
          <input
            v-model="pair.key"
            type="text"
            placeholder="KEY (e.g. API_KEY)"
            class="min-h-10 w-1/3 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3 text-xs font-mono text-[var(--admin-fg)] outline-none"
          />
          <div class="flex-1 relative">
            <input
              v-if="!pair.show"
              :value="'*****'"
              type="text"
              disabled
              class="min-h-10 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-border)]/30 px-3 pr-16 text-xs font-mono text-[var(--admin-muted)] outline-none cursor-not-allowed"
            />
            <input
              v-else
              v-model="pair.value"
              type="text"
              placeholder="VALUE (Encrypted on save)"
              class="min-h-10 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3 pr-16 text-xs font-mono text-[var(--admin-fg)] outline-none focus:border-blue-500"
            />
            <button
              type="button"
              class="absolute right-1.5 top-1/2 -translate-y-1/2 px-2.5 py-1.5 rounded-lg text-[10px] font-semibold text-[var(--admin-muted)] hover:text-[var(--admin-fg)] hover:bg-[var(--admin-border)]/50"
              @click="pair.show = !pair.show"
            >
              {{ pair.show ? 'Hide' : 'Show' }}
            </button>
          </div>
          <button type="button" class="text-xs font-semibold text-red-500 px-3 py-2 rounded-xl hover:bg-red-500/10" @click="removeEnvPair(idx)">Remove</button>
        </div>
      </div>

      <!-- JavaScript Script Editor Card -->
      <div class="rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm flex flex-col gap-4">
        <div class="flex items-center justify-between border-b border-[var(--admin-border)] pb-3">
          <div>
            <h2 class="text-base font-bold tracking-tight">3. JavaScript Worker Script</h2>
            <p class="text-xs text-[var(--admin-muted)]">Executed via Goja embedded engine. Supports `async function scheduled(event, env)` with `await fetch` and `console.log`.</p>
          </div>
          <span class="rounded-full bg-emerald-500/10 px-3 py-1 text-xs font-semibold text-emerald-600 dark:text-emerald-400">ES2020 + Fetch</span>
        </div>

        <div class="overflow-hidden rounded-xl border border-[var(--admin-border)] shadow-sm">
          <!-- Editor Header -->
          <div class="flex items-center justify-between bg-slate-100 dark:bg-slate-900 px-4 py-2.5 border-b border-[var(--admin-border)]">
            <div class="flex items-center gap-2">
              <div class="flex items-center gap-1.5">
                <span class="size-2.5 rounded-full bg-red-400"></span>
                <span class="size-2.5 rounded-full bg-amber-400"></span>
                <span class="size-2.5 rounded-full bg-emerald-400"></span>
              </div>
              <span class="ml-3 text-xs font-semibold text-[var(--admin-muted)]">worker.js</span>
            </div>
            <span class="text-[10px] font-medium text-[var(--admin-muted)] uppercase tracking-wider">JavaScript / ES2020</span>
          </div>

          <!-- Editor Body -->
          <textarea
            v-model="formScript"
            rows="18"
            spellcheck="false"
            class="w-full bg-slate-50 dark:bg-slate-900 p-4 font-mono text-xs text-slate-800 dark:text-slate-300 outline-none leading-relaxed resize-y"
            placeholder="async function scheduled(event, env) { ... }"
          ></textarea>
        </div>
      </div>

      <!-- Bottom Actions -->
      <div class="flex items-center justify-end gap-3 pt-2">
        <button
          type="button"
          class="inline-flex min-h-11 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-5 text-xs font-semibold text-[var(--admin-fg)] hover:bg-[var(--admin-border)]/50"
          @click="router.push({ name: 'admin-cron' })"
        >
          Cancel
        </button>
        <button
          type="button"
          class="inline-flex min-h-11 items-center justify-center rounded-xl bg-blue-600 px-6 text-xs font-semibold text-white shadow-sm hover:bg-blue-700 disabled:opacity-60"
          :disabled="submitting || loading"
          @click="handleSave"
        >
          {{ submitting ? 'Saving...' : (isEdit ? 'Save Changes' : 'Create Cron Job') }}
        </button>
      </div>
    </div>
  </div>
</template>
