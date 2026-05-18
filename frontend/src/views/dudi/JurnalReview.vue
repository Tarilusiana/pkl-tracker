<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Review Jurnal</h2>
      <p class="text-xs text-gray-500 mt-0.5">Baca dan beri komentar pada jurnal siswa</p>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat...</div>

    <template v-else>
      <div class="bg-white rounded-xl p-3 border border-gray-100 mb-4">
        <select
          v-model="selectedStudent"
          class="w-full px-4 py-2.5 rounded-xl border border-gray-200 text-sm bg-white focus:border-primary outline-none"
        >
          <option value="">Semua Siswa</option>
          <option v-for="s in students" :key="s.id" :value="s.id">{{ s.name }}</option>
        </select>
      </div>

      <div v-if="filteredJournals.length === 0" class="text-center py-8 text-gray-400 text-sm">
        Belum ada jurnal
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="j in filteredJournals"
          :key="j.ID || j.id"
          class="bg-white rounded-xl p-4 border border-gray-100"
        >
          <div class="flex items-start justify-between mb-2">
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
                {{ (j.Student?.FullName || '').charAt(0) }}
              </div>
              <div>
                <p class="text-sm font-medium text-gray-800">{{ j.Student?.FullName }}</p>
                <p class="text-[10px] text-gray-400">{{ formatDate(j.Date) }}</p>
              </div>
            </div>
          </div>

          <p class="text-sm text-gray-700 mb-3">{{ j.Activity }}</p>

          <div v-if="j.Reflection" class="mb-2 ml-2 pl-3 border-l-2 border-gray-200">
            <p class="text-xs text-gray-500 italic">{{ j.Reflection }}</p>
          </div>

          <div v-if="j.TeacherComment" class="mb-2 ml-2 pl-3 border-l-2 border-info/30">
            <p class="text-[10px] text-gray-400">Komentar Guru:</p>
            <p class="text-xs text-gray-700">{{ j.TeacherComment }}</p>
          </div>

          <div v-if="j.DudiComment" class="mb-2 ml-2 pl-3 border-l-2 border-warning/30">
            <p class="text-[10px] text-gray-400">Komentar Anda:</p>
            <p class="text-xs text-gray-700">{{ j.DudiComment }}</p>
          </div>

          <div v-if="!j.DudiComment" class="mt-2">
            <div class="flex items-center gap-2">
              <input
                v-model="commentInputs[j.ID || j.id]"
                type="text"
                placeholder="Tulis komentar..."
                class="flex-1 px-3 py-2 rounded-lg border border-gray-200 text-xs focus:border-primary outline-none"
                @keyup.enter="addComment(j)"
              />
              <button
                @click="addComment(j)"
                :disabled="!commentInputs[j.ID || j.id]?.trim() || commenting"
                class="px-3 py-2 bg-primary text-white rounded-lg text-xs font-medium hover:bg-primary-light disabled:opacity-50 transition-colors"
              >
                Kirim
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { get, post } from '../../api'

const students = ref([])
const journals = ref([])
const selectedStudent = ref('')
const loading = ref(true)
const commenting = ref(false)
const commentInputs = reactive({})

const filteredJournals = computed(() => {
  if (!selectedStudent.value) return journals.value
  return journals.value.filter(j => (j.StudentID || j.student_id) === selectedStudent.value)
})

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

async function fetchData() {
  loading.value = true
  try {
    const [dashRes, jrnRes] = await Promise.all([
      get('/dudi/dashboard'),
      get('/jurnal')
    ])
    students.value = dashRes.students || []
    journals.value = jrnRes.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function addComment(j) {
  const jid = j.ID || j.id
  const comment = commentInputs[jid]?.trim()
  if (!comment) return
  commenting.value = true
  try {
    await post('/jurnal/comment', { jurnal_id: jid, comment })
    j.DudiComment = comment
    commentInputs[jid] = ''
  } catch (e) {
    alert(e.message)
  } finally {
    commenting.value = false
  }
}

onMounted(fetchData)
</script>
