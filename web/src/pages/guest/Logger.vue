<script setup lang="ts">
  import { ref, computed, watch } from 'vue'
  import { getLogs, getLogDetail, deleteLog } from '../../api/logger'
  import { axiosLib } from '../../lib'

  const selectedFile = ref<string | null>(null)
  const searchQuery = ref('')
  const selectedLevels = ref<Record<string, boolean>>({
    info: true,
    error: true,
  })
  const currentPage = ref(1)
  const expandedIndex = ref<number | null>(null)

  const filesQuery = getLogs()

  const levelsParam = computed(() =>
    Object.keys(selectedLevels.value)
      .filter((k) => selectedLevels.value[k])
      .join(',')
  )

  const detailQuery = getLogDetail({
    props: computed(() => ({
      param: { file_name: selectedFile.value ?? '' },
      query: {
        search: searchQuery.value,
        page: String(currentPage.value),
        limit: '50',
        levels: levelsParam.value,
      },
    })),
  })

  const deleteMutation = deleteLog({
    onSuccess: (_, fileName) => {
      if (selectedFile.value === fileName) selectedFile.value = null
      filesQuery.refetch()
    },
  })

  let searchTimeout: ReturnType<typeof setTimeout> | null = null
  watch(searchQuery, () => {
    if (searchTimeout) clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
      currentPage.value = 1
    }, 300)
  })

  function levelBadge(level: string): string {
    if (level === 'error' || level === 'emergency' || level === 'alert' || level === 'critical')
      return 'bg-red-600 text-white'
    if (level === 'warning') return 'bg-yellow-600 text-white'
    if (level === 'info' || level === 'notice') return 'bg-blue-600 text-white'
    if (level === 'sql') return 'bg-indigo-600 text-white'
    return 'bg-gray-600 text-white'
  }

  function formatSize(bytes: number) {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
  }

  function toggleLevel(level: string) {
    selectedLevels.value = { ...selectedLevels.value, [level]: !selectedLevels.value[level] }
    currentPage.value = 1
  }

  function clearSearch() {
    searchQuery.value = ''
    currentPage.value = 1
  }

  function selectFile(name: string) {
    selectedFile.value = name
    currentPage.value = 1
    expandedIndex.value = null
  }

  function toggleExpand(index: number) {
    expandedIndex.value = expandedIndex.value === index ? null : index
  }

  const totalPages = computed(() => {
    if (!detailQuery.data.value?.meta) return 1
    return Math.ceil(detailQuery.data.value.meta.total / detailQuery.data.value.meta.limit) || 1
  })

  const pageNumbers = computed((): (number | '...')[] => {
    const tp = totalPages.value
    if (tp <= 7) return Array.from({ length: tp }, (_, i) => i + 1)
    const pages: (number | '...')[] = [1]
    if (currentPage.value > 3) pages.push('...')
    const start = Math.max(2, currentPage.value - 1)
    const end = Math.min(tp - 1, currentPage.value + 1)
    for (let i = start; i <= end; i++) pages.push(i)
    if (currentPage.value < tp - 2) pages.push('...')
    if (tp > 1) pages.push(tp)
    return pages
  })

  function handleDownload(name: string) {
    const url = axiosLib.instance.defaults.baseURL + '/log/' + encodeURIComponent(name) + '/download'
    const a = document.createElement('a')
    a.href = url
    a.download = name
    a.click()
  }

  function handleDelete(name: string) {
    if (!confirm(`Delete ${name}?`)) return
    deleteMutation.mutate(name)
  }
</script>

<template>
  <div class="flex h-screen min-h-0">
    <aside class="flex w-60 shrink-0 flex-col border-r border-gray-800 bg-gray-950">
      <div class="border-b border-gray-800 p-3">
        <h2 class="text-xs font-semibold tracking-wider text-gray-400 uppercase">Log Files</h2>
      </div>
      <div class="flex-1 overflow-y-auto">
        <div v-if="filesQuery.isLoading.value" class="p-4 text-sm text-gray-600">Loading...</div>
        <div v-else-if="filesQuery.data.value && filesQuery.data.value.length > 0">
          <button
            v-for="file in filesQuery.data.value"
            :key="file.name"
            :class="[
              'w-full border-b border-gray-900 px-3 py-2 text-left text-sm transition hover:bg-gray-800/50',
              selectedFile === file.name ? 'bg-gray-800 text-white' : 'text-gray-400',
            ]"
            @click="selectFile(file.name)"
          >
            <div class="flex items-center justify-between gap-2">
              <span class="truncate">{{ file.name }}</span>
              <span class="shrink-0 text-xs text-gray-600">{{ formatSize(file.size) }}</span>
            </div>
          </button>
        </div>
        <div v-else class="p-4 text-sm text-gray-600">No log files found</div>
      </div>
    </aside>

    <main class="flex min-w-0 flex-1 flex-col bg-gray-950">
      <template v-if="selectedFile">
        <div class="flex flex-wrap items-center gap-2 border-b border-gray-800 bg-gray-950 p-3">
          <div class="relative min-w-[200px] flex-1">
            <input
              type="text"
              placeholder="Search logs..."
              class="w-full rounded border border-gray-700 bg-gray-900 px-3 py-1.5 text-sm text-gray-200 placeholder-gray-600 focus:border-blue-500 focus:outline-none"
              v-model="searchQuery"
            />
            <button
              v-if="searchQuery"
              class="absolute top-1/2 right-2 -translate-y-1/2 text-gray-500 hover:text-gray-300"
              @click="clearSearch"
            >
              ✕
            </button>
          </div>
          <div class="flex flex-wrap items-center gap-1">
            <button
              :class="[
                'mr-1.5 rounded bg-blue-700 px-3 py-1.5 text-xs font-semibold text-white uppercase transition',
                selectedLevels['info'] ? 'ring-2 ring-white/40' : 'opacity-60 ring-0',
              ]"
              @click="toggleLevel('info')"
            >
              info
            </button>
            <button
              :class="[
                'rounded bg-red-700 px-3 py-1.5 text-xs font-semibold text-white uppercase transition',
                selectedLevels['error'] ? 'ring-2 ring-white/40' : 'opacity-60 ring-0',
              ]"
              @click="toggleLevel('error')"
            >
              error
            </button>
          </div>
          <div class="ml-auto flex items-center gap-1">
            <button
              class="rounded bg-gray-800 px-3 py-1.5 text-xs text-gray-300 transition hover:bg-gray-700"
              @click="detailQuery.refetch()"
            >
              ⟳
            </button>
            <button
              class="rounded bg-gray-800 px-3 py-1.5 text-xs text-gray-300 transition hover:bg-gray-700"
              @click="handleDownload(selectedFile!)"
            >
              Download
            </button>
            <button
              class="rounded bg-red-800 px-3 py-1.5 text-xs text-white transition hover:bg-red-700"
              @click="handleDelete(selectedFile!)"
            >
              Delete
            </button>
          </div>
        </div>

        <div class="flex-1 overflow-y-auto">
          <div v-if="detailQuery.isLoading.value" class="p-6 text-sm text-gray-500">Loading...</div>
          <div v-else-if="detailQuery.isError.value" class="p-6 text-sm text-red-400">Failed to load log file</div>
          <div v-else-if="(detailQuery.data.value?.data ?? []).length > 0">
            <div v-for="(entry, i) in detailQuery.data.value?.data ?? []" :key="i" class="border-b border-gray-800">
              <button
                class="flex w-full items-start gap-2 px-4 py-2.5 text-left transition hover:bg-gray-800/30"
                @click="toggleExpand(i)"
              >
                <span class="mt-0.5 w-3 shrink-0 text-xs text-gray-700">{{ expandedIndex === i ? '▼' : '▶' }}</span>
                <span
                  :class="[
                    'shrink-0 rounded px-2 py-0.5 font-mono text-xs font-semibold uppercase',
                    levelBadge(entry.level),
                  ]"
                  >{{ entry.level }}</span
                >
                <span class="w-32 shrink-0 font-mono text-xs text-gray-500">{{ entry.time }}</span>
                <span class="line-clamp-2 text-sm break-all text-gray-300">{{ entry.message }}</span>
              </button>
              <div v-if="expandedIndex === i" class="bg-gray-900/30 px-12 pt-1 pb-3 font-mono text-xs text-gray-400">
                <pre class="border-l border-gray-700 pl-3 break-all whitespace-pre-wrap text-gray-400">{{
                  entry.message
                }}</pre>
              </div>
            </div>
          </div>
          <div v-else class="p-6 text-sm text-gray-500">No log entries match your filters</div>
        </div>

        <div
          v-if="detailQuery.data.value?.meta && totalPages > 1"
          class="flex items-center justify-between border-t border-gray-800 bg-gray-950 px-4 py-2"
        >
          <span class="text-xs text-gray-500">{{ detailQuery.data.value.meta.total }} entries</span>
          <div class="flex items-center gap-1">
            <button
              class="rounded bg-gray-800 px-3 py-1 text-xs text-gray-300 transition hover:bg-gray-700 disabled:cursor-not-allowed disabled:opacity-40"
              :disabled="currentPage <= 1"
              @click="currentPage = Math.max(1, currentPage - 1)"
            >
              Prev
            </button>
            <template v-for="p in pageNumbers" :key="typeof p === 'string' ? p : p">
              <span v-if="p === '...'" class="px-1 text-xs text-gray-600">...</span>
              <button
                v-else
                :class="[
                  'rounded px-3 py-1 text-xs transition',
                  currentPage === p ? 'bg-gray-700 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700',
                ]"
                @click="currentPage = p"
              >
                {{ p }}
              </button>
            </template>
            <button
              class="rounded bg-gray-800 px-3 py-1 text-xs text-gray-300 transition hover:bg-gray-700 disabled:cursor-not-allowed disabled:opacity-40"
              :disabled="currentPage >= totalPages"
              @click="currentPage = Math.min(totalPages, currentPage + 1)"
            >
              Next
            </button>
          </div>
        </div>
      </template>
      <div v-else class="flex flex-1 items-center justify-center text-sm text-gray-600">Select a log file to view</div>
    </main>
  </div>
</template>
