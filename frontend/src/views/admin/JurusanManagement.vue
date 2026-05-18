<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Data Jurusan</h2>
        <p class="text-sm text-gray-500 mt-0.5">Kelola daftar jurusan yang tersedia</p>
      </div>
      <button @click="openCreateModal" class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark">
        <PlusIcon :size="16" /> Tambah Jurusan
      </button>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <div v-else class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-100 bg-gray-50/50 text-left">
              <th class="px-4 py-3 font-medium text-gray-500">Kode</th>
              <th class="px-4 py-3 font-medium text-gray-500">Nama Jurusan</th>
              <th class="px-4 py-3 font-medium text-gray-500 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-for="j in jurusanList" :key="j.id" class="hover:bg-gray-50/50">
              <td class="px-4 py-3">
                <span class="inline-flex px-2 py-0.5 rounded-full text-[11px] font-mono font-medium bg-gray-100 text-gray-700">{{ j.kode }}</span>
              </td>
              <td class="px-4 py-3 font-medium text-gray-800">{{ j.nama }}</td>
              <td class="px-4 py-3 text-center">
                <div class="flex items-center justify-center gap-1">
                  <button @click="openEditModal(j)" class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10"><PencilIcon :size="14" /></button>
                  <button @click="confirmDelete(j)" class="p-1.5 rounded-lg text-gray-400 hover:text-danger hover:bg-danger/10"><TrashIcon :size="14" /></button>
                </div>
              </td>
            </tr>
            <tr v-if="jurusanList.length === 0">
              <td colspan="3" class="px-4 py-8 text-center text-gray-400">Belum ada data jurusan</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="closeModal">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full">
        <h3 class="text-lg font-bold text-gray-800 mb-4">{{ editing ? 'Edit Jurusan' : 'Tambah Jurusan' }}</h3>
        <form @submit.prevent="saveJurusan" class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Kode Jurusan</label>
            <input v-model="form.kode" placeholder="RPL" required maxlength="20" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none uppercase" />
            <p class="text-xs text-gray-400 mt-1">Singkatan, contoh: RPL, TKJ, MM</p>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Nama Jurusan</label>
            <input v-model="form.nama" placeholder="Rekayasa Perangkat Lunak" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div class="flex gap-2 pt-2">
            <button type="button" @click="closeModal" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
            <button type="submit" :disabled="saving" class="flex-1 py-2.5 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark disabled:opacity-50">{{ saving ? 'Menyimpan...' : 'Simpan' }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirm -->
    <div v-if="showDelete" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="showDelete = null">
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full text-center">
        <div class="w-12 h-12 bg-danger/10 rounded-full flex items-center justify-center mx-auto mb-3">
          <AlertCircleIcon :size="24" class="text-danger" />
        </div>
        <h3 class="text-lg font-bold text-gray-800 mb-1">Hapus Jurusan?</h3>
        <p class="text-sm text-gray-500 mb-4">{{ deleteTarget?.kode }} - {{ deleteTarget?.nama }}</p>
        <div class="flex gap-2">
          <button @click="showDelete = null" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
          <button @click="deleteJurusan" :disabled="saving" class="flex-1 py-2.5 bg-danger text-white rounded-xl text-sm font-medium hover:bg-red-600 disabled:opacity-50">{{ saving ? 'Menghapus...' : 'Hapus' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { PlusIcon, PencilIcon, TrashIcon, AlertCircleIcon } from 'lucide-vue-next'
import { get, post, put, del } from '@/api'

const jurusanList = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const editing = ref(null)
const showDelete = ref(null)
const deleteTarget = ref(null)

const form = reactive({
  kode: '', nama: ''
})

async function fetchJurusan() {
  loading.value = true
  try {
    const res = await get('/admin/jurusan')
    jurusanList.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editing.value = null
  Object.assign(form, { kode: '', nama: '' })
  showModal.value = true
}

function openEditModal(j) {
  editing.value = j
  Object.assign(form, { kode: j.kode, nama: j.nama })
  showModal.value = true
}

function closeModal() { showModal.value = false; editing.value = null }

async function saveJurusan() {
  saving.value = true
  try {
    if (editing.value) {
      await put('/admin/jurusan/' + editing.value.id, { ...form })
    } else {
      await post('/admin/jurusan', { ...form })
    }
    closeModal()
    fetchJurusan()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

function confirmDelete(j) { deleteTarget.value = j; showDelete.value = true }

async function deleteJurusan() {
  saving.value = true
  try {
    await del('/admin/jurusan/' + deleteTarget.value.id)
    showDelete.value = null
    fetchJurusan()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(fetchJurusan)
</script>
