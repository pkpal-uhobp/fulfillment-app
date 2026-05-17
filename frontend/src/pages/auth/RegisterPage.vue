<script setup>
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { ArrowLeft, ArrowRight, LockKeyhole, Mail, Phone, UserRound } from '@lucide/vue'
import { apiFetch, saveAuth } from '@/shared/api/http'

const router = useRouter()

const fullName = ref('')
const email = ref('')
const phone = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

async function submit() {
  error.value = ''
  success.value = ''

  if (!fullName.value.trim() || !email.value.trim() || !password.value) {
    error.value = 'Заполните ФИО, email и пароль.'
    return
  }

  loading.value = true

  try {
    const payload = await apiFetch('/auth/register', {
      method: 'POST',
      body: {
        full_name: fullName.value.trim(),
        email: email.value.trim(),
        phone: phone.value.trim() || undefined,
        password: password.value,
      },
    })

    const { user } = saveAuth(payload)
    success.value = `Аккаунт создан: ${user?.full_name || user?.email || 'пользователь'}.`

    setTimeout(() => router.push('/'), 450)
  } catch (err) {
    error.value = err?.message || 'Не удалось зарегистрироваться. Проверьте данные.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#07101f] text-white selection:bg-[#ff4248] selection:text-white">
    <section class="relative min-h-screen overflow-hidden px-5 py-8 sm:px-8 lg:px-12">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_10%_10%,rgba(255,66,72,0.28),transparent_34%),radial-gradient(circle_at_90%_15%,rgba(0,166,214,0.28),transparent_34%),linear-gradient(135deg,#1d0714_0%,#07101f_48%,#06394a_100%)]"></div>

      <div class="relative mx-auto flex min-h-[calc(100vh-4rem)] max-w-[1220px] items-center justify-center">
        <div class="grid w-full overflow-hidden rounded-[2.2rem] border border-white/10 bg-white/5 shadow-[0_35px_120px_rgba(0,0,0,0.35)] backdrop-blur lg:grid-cols-[0.9fr_1.1fr]">
          <div class="hidden p-10 lg:flex lg:flex-col lg:justify-between">
            <RouterLink to="/" class="inline-flex items-center gap-3 text-sm font-black text-white/80 transition hover:text-white">
              <ArrowLeft class="h-5 w-5" /> На главную
            </RouterLink>

            <div class="my-10 rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-[0_30px_90px_rgba(0,0,0,0.25)] backdrop-blur">
              <div class="relative h-[270px] overflow-hidden rounded-[1.5rem] bg-[#0b1527]">
                <div class="absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(255,66,72,0.35),transparent_30%),radial-gradient(circle_at_85%_35%,rgba(0,166,214,0.28),transparent_32%)]"></div>
                <svg class="absolute inset-0 h-full w-full" viewBox="0 0 520 300" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                  <rect x="58" y="82" width="190" height="136" rx="22" fill="rgba(255,255,255,0.08)" stroke="rgba(255,255,255,0.18)"/>
                  <rect x="88" y="113" width="48" height="48" rx="10" fill="#ff4248"/>
                  <rect x="150" y="113" width="68" height="14" rx="7" fill="rgba(255,255,255,0.55)"/>
                  <rect x="150" y="143" width="82" height="14" rx="7" fill="rgba(255,255,255,0.25)"/>
                  <rect x="88" y="181" width="132" height="14" rx="7" fill="rgba(255,255,255,0.18)"/>
                  <path d="M320 198h-54v-58h112l35 36v22h-30" fill="rgba(255,66,72,0.18)" stroke="#ff767b" stroke-width="6" stroke-linejoin="round"/>
                  <path d="M378 140v38h35" stroke="#ff767b" stroke-width="6" stroke-linecap="round" stroke-linejoin="round"/>
                  <circle cx="302" cy="205" r="18" fill="#07101f" stroke="white" stroke-width="6"/>
                  <circle cx="383" cy="205" r="18" fill="#07101f" stroke="white" stroke-width="6"/>
                  <path d="M94 64h88M334 82h86M263 238h166" stroke="rgba(255,255,255,0.25)" stroke-width="8" stroke-linecap="round"/>
                  <rect x="340" y="52" width="74" height="74" rx="22" fill="rgba(255,255,255,0.08)" stroke="rgba(255,255,255,0.16)"/>
                  <path d="M362 74h12v12h-12V74Zm24 0h12v12h-12V74Zm-24 24h12v12h-12V98Zm24 24h12v12h-12v-12Zm-12-24h12v12h-12V98Zm-12 24h12v12h-12v-12Z" fill="#fff" fill-opacity="0.9"/>
                </svg>
              </div>
              <div class="mt-5 grid grid-cols-3 gap-3 text-center">
                <div class="rounded-2xl bg-white/10 p-4">
                  <div class="text-2xl font-black">QR</div>
                  <div class="mt-1 text-xs text-white/55">контроль</div>
                </div>
                <div class="rounded-2xl bg-white/10 p-4">
                  <div class="text-2xl font-black">24/7</div>
                  <div class="mt-1 text-xs text-white/55">статусы</div>
                </div>
                <div class="rounded-2xl bg-white/10 p-4">
                  <div class="text-2xl font-black">API</div>
                  <div class="mt-1 text-xs text-white/55">заявки</div>
                </div>
              </div>
            </div>
            <div>
              <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff9ca0]">Fulfillment Transit</p>
              <h1 class="mt-6 text-5xl font-black leading-tight tracking-[-0.05em]">Регистрация клиента</h1>
              <p class="mt-6 text-lg leading-8 text-white/72">Создайте аккаунт, чтобы оформлять заявки, выбирать склады, отслеживать грузовые места и видеть историю статусов.</p>
            </div>
          </div>

          <div class="bg-white p-6 text-[#07101f] sm:p-10 lg:p-12">
            <RouterLink to="/" class="mb-8 inline-flex items-center gap-2 text-sm font-black text-slate-500 transition hover:text-[#ff4248] lg:hidden">
              <ArrowLeft class="h-5 w-5" /> На главную
            </RouterLink>

            <h2 class="text-4xl font-black tracking-[-0.05em]">Регистрация</h2>
            <p class="mt-3 text-slate-600">По умолчанию создаётся клиентский аккаунт для оформления заявок.</p>

            <form class="mt-8 space-y-5" @submit.prevent="submit">
              <label class="block">
                <span class="mb-2 block text-xs font-black uppercase tracking-[0.28em] text-slate-500">ФИО</span>
                <span class="flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-5 py-4 focus-within:border-[#ff4248]">
                  <UserRound class="h-5 w-5 text-slate-400" />
                  <input v-model="fullName" class="w-full bg-transparent font-bold outline-none" placeholder="Иванов Иван" autocomplete="name" />
                </span>
              </label>

              <label class="block">
                <span class="mb-2 block text-xs font-black uppercase tracking-[0.28em] text-slate-500">Email</span>
                <span class="flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-5 py-4 focus-within:border-[#ff4248]">
                  <Mail class="h-5 w-5 text-slate-400" />
                  <input v-model="email" class="w-full bg-transparent font-bold outline-none" type="email" placeholder="client@example.com" autocomplete="email" />
                </span>
              </label>

              <label class="block">
                <span class="mb-2 block text-xs font-black uppercase tracking-[0.28em] text-slate-500">Телефон</span>
                <span class="flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-5 py-4 focus-within:border-[#ff4248]">
                  <Phone class="h-5 w-5 text-slate-400" />
                  <input v-model="phone" class="w-full bg-transparent font-bold outline-none" placeholder="+79991234567" autocomplete="tel" />
                </span>
              </label>

              <label class="block">
                <span class="mb-2 block text-xs font-black uppercase tracking-[0.28em] text-slate-500">Пароль</span>
                <span class="flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-5 py-4 focus-within:border-[#ff4248]">
                  <LockKeyhole class="h-5 w-5 text-slate-400" />
                  <input v-model="password" class="w-full bg-transparent font-bold outline-none" type="password" placeholder="Минимум 6 символов" autocomplete="new-password" />
                </span>
              </label>

              <p v-if="error" class="rounded-2xl border border-red-200 bg-red-50 px-5 py-4 text-sm font-bold text-red-700">{{ error }}</p>
              <p v-if="success" class="rounded-2xl border border-emerald-200 bg-emerald-50 px-5 py-4 text-sm font-bold text-emerald-700">{{ success }}</p>

              <button class="flex w-full items-center justify-center gap-3 rounded-2xl bg-[#ff4248] px-6 py-5 font-black text-white shadow-[0_18px_45px_rgba(255,66,72,0.24)] transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60" :disabled="loading" type="submit">
                {{ loading ? 'Создаём...' : 'Создать аккаунт' }} <ArrowRight class="h-5 w-5" />
              </button>
            </form>

            <p class="mt-8 text-center text-sm text-slate-500">
              Уже есть аккаунт?
              <RouterLink to="/login" class="font-black text-[#ff4248] hover:underline">Войти</RouterLink>
            </p>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
