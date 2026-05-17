<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { CheckCircle2, Edit3, LogOut, Mail, Phone, Save, ShieldCheck, UserRound, X } from '@lucide/vue'
import { apiFetch, clearAuth, getCurrentUser, loadMe, updateCurrentUserLocal } from '@/shared/api/http'
import { isValidEmail, isValidPhone, roleLabel } from './clientUtils'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const editing = ref(false)
const error = ref('')
const success = ref('')
const user = ref(getCurrentUser())

const form = ref({
  full_name: '',
  email: '',
  phone: '',
})
const formErrors = ref({})

const displayName = computed(() => user.value?.full_name || user.value?.email || 'Клиент')
const initials = computed(() => displayName.value.split(' ').filter(Boolean).slice(0, 2).map((item) => item[0]?.toUpperCase()).join('') || 'К')

function fillForm() {
  form.value = {
    full_name: user.value?.full_name || '',
    email: user.value?.email || '',
    phone: user.value?.phone || '',
  }
}

async function refresh() {
  loading.value = true
  error.value = ''
  try {
    user.value = await loadMe()
    fillForm()
  } catch (err) {
    error.value = err?.message || 'Не удалось обновить профиль.'
    user.value = getCurrentUser()
    fillForm()
  } finally {
    loading.value = false
  }
}

function startEdit() {
  fillForm()
  editing.value = true
  error.value = ''
  success.value = ''
  formErrors.value = {}
}

function cancelEdit() {
  fillForm()
  editing.value = false
  error.value = ''
  formErrors.value = {}
}

function validateProfile() {
  const errors = {}
  if (!form.value.full_name.trim()) errors.full_name = 'Укажите ФИО.'
  if (!isValidEmail(form.value.email)) errors.email = 'Укажите корректный email в формате name@example.com.'
  if (!isValidPhone(form.value.phone)) errors.phone = 'Укажите телефон в формате +79991234567.'
  formErrors.value = errors
  return Object.keys(errors).length === 0
}

function clearProfileError(key) {
  if (formErrors.value[key]) formErrors.value = { ...formErrors.value, [key]: '' }
}

async function saveProfile() {
  error.value = ''
  success.value = ''
  if (!validateProfile()) {
    error.value = 'Проверьте выделенные поля.'
    return
  }

  saving.value = true

  const patch = {
    full_name: form.value.full_name.trim(),
    email: form.value.email.trim(),
    phone: form.value.phone.trim(),
  }

  try {
    // Если позже добавишь backend endpoint PATCH /auth/me, этот код сразу начнёт сохранять профиль на сервере.
    const payload = await apiFetch('/auth/me', {
      method: 'PATCH',
      auth: true,
      body: patch,
    })
    const updated = payload?.user || payload?.data?.user || payload
    user.value = updateCurrentUserLocal(updated && typeof updated === 'object' ? updated : patch)
    success.value = 'Профиль обновлён.'
  } catch (err) {
    if ([404, 405, 501].includes(err?.status)) {
      user.value = updateCurrentUserLocal(patch)
      success.value = 'Профиль обновлён в интерфейсе. Для сохранения в базе добавь backend endpoint PATCH /auth/me.'
    } else {
      error.value = err?.message || 'Не удалось сохранить профиль.'
      saving.value = false
      return
    }
  }

  editing.value = false
  saving.value = false
}

function logout() {
  clearAuth()
  router.push('/')
}

onMounted(refresh)
</script>

<template>
  <section class="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
    <div class="relative overflow-hidden rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 backdrop-blur sm:p-8">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_15%_10%,rgba(255,66,72,0.24),transparent_32%),radial-gradient(circle_at_80%_15%,rgba(20,184,166,0.16),transparent_30%)]"></div>
      <div class="relative">
        <p class="text-xs font-black uppercase tracking-[0.45em] text-[#ff9ca0]">Профиль клиента</p>
        <h1 class="mt-4 text-5xl font-black leading-none tracking-[-0.06em]">Личный кабинет</h1>
        <p class="mt-4 max-w-2xl text-white/65">Данные профиля используются для заявок, уведомлений, QR-проверки и ограничения доступа к вашим грузовым местам.</p>

        <div class="mt-8 flex items-center gap-5 rounded-[2rem] border border-white/10 bg-[#0b1527]/80 p-5">
          <div class="grid h-20 w-20 shrink-0 place-items-center rounded-[1.6rem] bg-[#ff4248] text-3xl font-black shadow-[0_20px_60px_rgba(255,66,72,0.25)]">{{ initials }}</div>
          <div class="min-w-0">
            <div class="truncate text-2xl font-black">{{ displayName }}</div>
            <div class="mt-2 inline-flex rounded-full border border-white/10 bg-white/[0.06] px-3 py-1 text-xs font-black text-white/60">{{ roleLabel(user?.role || 'client') }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-[2rem] bg-white p-6 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.12)] sm:p-8">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <p class="text-xs font-black uppercase tracking-[0.35em] text-[#ff4248]">Аккаунт</p>
          <h2 class="mt-2 text-3xl font-black tracking-[-0.05em]">Данные пользователя</h2>
        </div>
        <div class="flex gap-2">
          <button v-if="!editing" class="inline-flex items-center justify-center gap-2 rounded-2xl bg-slate-100 px-4 py-3 text-sm font-black text-slate-700 transition hover:bg-[#07101f] hover:text-white" @click="refresh">
            {{ loading ? '...' : 'Обновить' }}
          </button>
          <button v-if="!editing" class="inline-flex items-center justify-center gap-2 rounded-2xl bg-[#ff4248] px-4 py-3 text-sm font-black text-white transition hover:-translate-y-0.5" @click="startEdit">
            <Edit3 class="h-4 w-4" /> Редактировать
          </button>
        </div>
      </div>

      <p v-if="error" class="mt-5 rounded-2xl border border-red-200 bg-red-50 px-5 py-4 text-sm font-bold text-red-700">{{ error }}</p>
      <p v-if="success" class="mt-5 rounded-2xl border border-emerald-200 bg-emerald-50 px-5 py-4 text-sm font-bold text-emerald-700">{{ success }}</p>

      <div v-if="!editing" class="mt-7 grid gap-4">
        <div class="flex items-center gap-4 rounded-3xl bg-slate-50 p-5">
          <UserRound class="h-6 w-6 text-[#ff4248]" />
          <div>
            <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">ФИО</div>
            <div class="mt-1 font-black">{{ user?.full_name || '—' }}</div>
          </div>
        </div>
        <div class="flex items-center gap-4 rounded-3xl bg-slate-50 p-5">
          <Mail class="h-6 w-6 text-[#ff4248]" />
          <div>
            <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Email</div>
            <div class="mt-1 font-black">{{ user?.email || '—' }}</div>
          </div>
        </div>
        <div class="flex items-center gap-4 rounded-3xl bg-slate-50 p-5">
          <Phone class="h-6 w-6 text-[#ff4248]" />
          <div>
            <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Телефон</div>
            <div class="mt-1 font-black">{{ user?.phone || '—' }}</div>
          </div>
        </div>
        <div class="flex items-center gap-4 rounded-3xl bg-slate-50 p-5">
          <ShieldCheck class="h-6 w-6 text-[#ff4248]" />
          <div>
            <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Роль</div>
            <div class="mt-1 font-black">{{ roleLabel(user?.role || 'client') }}</div>
          </div>
        </div>
      </div>

      <form v-else class="mt-7 grid gap-4" novalidate @submit.prevent="saveProfile">
        <label class="block">
          <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">ФИО</span>
          <input v-model="form.full_name" class="w-full rounded-3xl border bg-slate-50 px-5 py-4 font-black outline-none transition focus:bg-white" :class="formErrors.full_name ? 'border-[#ff4248]' : 'border-slate-200 focus:border-[#ff4248]'" placeholder="Иванов Иван" @input="clearProfileError('full_name')" />
          <p v-if="formErrors.full_name" class="mt-2 text-sm font-bold text-[#e11d48]">{{ formErrors.full_name }}</p>
        </label>
        <label class="block">
          <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Email</span>
          <input v-model="form.email" class="w-full rounded-3xl border bg-slate-50 px-5 py-4 font-black outline-none transition focus:bg-white" :class="formErrors.email ? 'border-[#ff4248]' : 'border-slate-200 focus:border-[#ff4248]'" inputmode="email" placeholder="client@example.com" @input="clearProfileError('email')" />
          <p v-if="formErrors.email" class="mt-2 text-sm font-bold text-[#e11d48]">{{ formErrors.email }}</p>
        </label>
        <label class="block">
          <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Телефон</span>
          <input v-model="form.phone" class="w-full rounded-3xl border bg-slate-50 px-5 py-4 font-black outline-none transition focus:bg-white" :class="formErrors.phone ? 'border-[#ff4248]' : 'border-slate-200 focus:border-[#ff4248]'" inputmode="tel" placeholder="+79991234567" @input="clearProfileError('phone')" />
          <p v-if="formErrors.phone" class="mt-2 text-sm font-bold text-[#e11d48]">{{ formErrors.phone }}</p>
        </label>
        <div class="grid gap-3 sm:grid-cols-2">
          <button class="inline-flex items-center justify-center gap-2 rounded-2xl bg-[#ff4248] px-5 py-4 font-black text-white disabled:opacity-60" :disabled="saving" type="submit">
            <Save class="h-5 w-5" /> {{ saving ? 'Сохраняем...' : 'Сохранить' }}
          </button>
          <button class="inline-flex items-center justify-center gap-2 rounded-2xl border border-slate-200 px-5 py-4 font-black text-slate-700 hover:bg-slate-50" type="button" @click="cancelEdit">
            <X class="h-5 w-5" /> Отмена
          </button>
        </div>
      </form>

      <button class="mt-7 flex w-full items-center justify-center gap-3 rounded-2xl bg-[#07101f] px-6 py-5 font-black text-white transition hover:bg-[#ff4248]" @click="logout">
        Выйти из аккаунта <LogOut class="h-5 w-5" />
      </button>
    </div>
  </section>
</template>
