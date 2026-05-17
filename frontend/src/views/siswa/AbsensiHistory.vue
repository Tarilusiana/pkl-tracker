<template>
  <div>
    <div class="flex items-center gap-3 mb-5">
      <router-link to="/siswa/absensi" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
        <ArrowLeftIcon :size="20" />
      </router-link>
      <div>
        <h2 class="text-lg font-bold text-gray-800">Riwayat Absensi</h2>
        <p class="text-xs text-gray-500">Bulan {{ currentMonth }}</p>
      </div>
    </div>

    <!-- Month selector -->
    <div class="flex items-center gap-2 mb-4">
      <button @click="prevMonth" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
        <ChevronLeftIcon :size="18" />
      </button>
      <span class="flex-1 text-center text-sm font-semibold text-gray-800">{{ currentMonth }}</span>
      <button @click="nextMonth" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
        <ChevronRightIcon :size="18" />
      </button>
    </div>

    <!-- Stats -->
    <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <div class="grid grid-cols-4 gap-2 text-center">
        <div>
          <p class="text-xl font-bold text-accent">{{ summary.hadir }}</p>
          <p class="text-[10px] text-gray-500">Hadir</p>
        </div>
        <div>
          <p class="text-xl font-bold text-warning">{{ summary.terlambat }}</p>
          <p class="text-[10px] text-gray-500">Terlambat</p>
        </div>
        <div>
          <p class="text-xl font-bold text-info">{{ summary.izin }}</p>
          <p class="text-[10px] text-gray-500">Izin</p>
        </div>
        <div>
          <p class="text-xl font-bold text-danger">{{ summary.sakit }}</p>
          <p class="text-[10px] text-gray-500">Sakit</p>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat...</div>

    <!-- Empty -->
    <div v-else-if="filteredHistory.length === 0" class="text-center py-8">
      <ClockIcon :size="40" class="mx-auto mb-2 text-gray-300" />
      <p class="text-sm text-gray-400">Belum ada riwayat absensi</p>
    </div>

    <!-- List -->
    <div v-else class="space-y-2">
      <div
        v-for="a in filteredHistory"
        :key="a.ID || a.id"
        class="bg-white rounded-xl p-4 border border-gray-100 flex items-center gap-3"
      >
        <div :class="['w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0', statusIconBg(a.Status || a.status)]">
          <component :is="statusIcon(a.Status || a.status)" :size="18" :class="statusIconColor(a.Status || a.status)" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-800">{{ formatDate(a.Timestamp || a.timestamp) }}</p>
          <p class="text-xs text-gray-500">{{ formatTime(a.Timestamp || a.timestamp) }} - {{ (a.Latitude ?? a.latitude)?.toFixed(4) }}, {{ (a.Longitude ?? a.longitude)?.toFixed(4) }}</p>
        </div>
        <span :class="['inline-flex px-2.5 py-0.5 rounded-full text-[10px] font-medium', statusBadgeStyle(a.Status || a.status)]">
          {{ statusLabel(a.Status || a.status) }}
        </span>
        <ChevronRightIcon :size="16" class="text-gray-300" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  ArrowLeftIcon, ChevronLeftIcon, ChevronRightIcon,
  CheckCircleIcon, ClockIcon, FileTextIcon, HeartIcon, AlertCircleIcon
} from 'lucide-vue-next'

const monthOffset = ref(0)
const history = ref([])
const loading = ref(true)

const currentMonth = computed(() => {
  const d = new Date()
  d.setMonth(d.getMonth() - monthOffset.value)
  return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
})

const filteredHistory = computed(() => {
  const d = new Date()
  d.setMonth(d.getMonth() - monthOffset.value)
  const targetMonth = d.getMonth()
  const targetYear = d.getFullYear()
  return history.value.filter((a) => {
    const t = new Date(a.Timestamp || a.timestamp)
    return t.getMonth() === targetMonth && t.getFullYear() === targetYear
  })
})

const summary = computed(() => {
  const items = filteredHistory.value
  return {
    hadir: items.filter((a) => (a.Status || a.status) === 'hadir').length,
    terlambat: items.filter((a) => (a.Status || a.status) === 'terlambat').length,
    izin: items.filter((a) => (a.Status || a.status) === 'izin').length,
    sakit: items.filter((a) => (a.Status || a.status) === 'sakit').length,
  }
})

function formatDate(ts) {
  const t = new Date(ts)
  return t.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

function formatTime(ts) {
  const t = new Date(ts)
  return t.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) + ' WIB'
}

function prevMonth() { monthOffset.value++ }
function nextMonth() { if (monthOffset.value > 0) monthOffset.value-- }

async function fetchHistory() {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/absensi/history', {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Gagal memuat riwayat')
    const json = await res.json()
    history.value = json.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchHistory)

function statusLabel(s) {
  return { hadir: 'Hadir', terlambat: 'Terlambat', izin: 'Izin', sakit: 'Sakit' }[s]
}

function statusBadgeStyle(s) {
  return {
    hadir: 'bg-accent/10 text-accent',
    terlambat: 'bg-warning/10 text-warning',
    izin: 'bg-info/10 text-info',
    sakit: 'bg-purple-100 text-purple-600'
  }[s]
}

function statusIcon(s) {
  return {
    hadir: CheckCircleIcon,
    terlambat: ClockIcon,
    izin: FileTextIcon,
    sakit: HeartIcon
  }[s]
}

function statusIconBg(s) {
  return {
    hadir: 'bg-accent/10',
    terlambat: 'bg-warning/10',
    izin: 'bg-info/10',
    sakit: 'bg-purple-100'
  }[s]
}

function statusIconColor(s) {
  return {
    hadir: 'text-accent',
    terlambat: 'text-warning',
    izin: 'text-info',
    sakit: 'text-purple-600'
  }[s]
}
</script>
