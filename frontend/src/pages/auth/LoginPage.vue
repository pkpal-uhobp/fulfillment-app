<script setup>
import { computed, reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'

import { ArrowLeft, ArrowRight, LockKeyhole, Mail } from '@lucide/vue'

import AuthIllustration from '@/shared/ui/AuthIllustration.vue'
import { apiFetch, saveAuth } from '@/shared/api/http'
import { humanizeApiError, normalizeEmail, validateEmail, validatePassword } from './authValidation'

const router = useRouter()

const form = reactive({ email: '', password: '' })
const errors = reactive({ email: '', password: '', common: '' })
const touched = reactive({ email: false, password: false })
const isSubmitting = ref(false)

const isValid = computed(() => !validateEmail(form.email) && !validatePassword(form.password))

function validateField(field) {
  if (field === 'email') errors.email = validateEmail(form.email)
  if (field === 'password') errors.password = validatePassword(form.password)
}

function validateForm() {
  touched.email = true
  touched.password = true
  validateField('email')
  validateField('password')
  errors.common = ''

  return !errors.email && !errors.password
}

function redirectByRole(user) {
  const role = String(user?.role || '').toLowerCase()

  if (role === 'admin') return router.push('/admin')
  if (role === 'logist' || role === 'logistician') return router.push('/logist')
  if (role === 'worker' || role === 'warehouse_worker') return router.push('/worker')

  return router.push('/client')
}

async function submit() {
  if (!validateForm()) return

  isSubmitting.value = true
  errors.common = ''

  try {
    const payload = await apiFetch('/auth/login', {
      method: 'POST',
      body: {
        email: normalizeEmail(form.email),
        password: form.password,
      },
    })

    const auth = saveAuth(payload)
    await redirectByRole(auth.user || payload?.user)
  } catch (error) {
    errors.common = humanizeApiError(error, 'Не удалось войти в аккаунт.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="min-h-[100dvh] overflow-hidden bg-[#07101f] p-0 text-[#07101f] lg:p-4">
    <div class="mx-auto grid min-h-[100dvh] bg-white shadow-2xl lg:h-[calc(100dvh-2rem)] lg:min-h-0 lg:max-w-[1500px] lg:overflow-hidden lg:rounded-[2rem] lg:grid-cols-[0.9fr_1.1fr]">
      <section class="hidden min-h-0 flex-col overflow-hidden bg-gradient-to-br from-[#3a1220] via-[#111827] to-[#073b46] p-6 lg:flex xl:p-8">
        <RouterLink to="/" class="inline-flex shrink-0 items-center gap-3 rounded-2xl px-1 py-2 text-sm font-black text-white/90 hover:text-white xl:text-base">
          <ArrowLeft class="h-5 w-5" />
          На главную
        </RouterLink>

        <div class="mt-5 min-h-0 flex-1 xl:mt-7">
          <AuthIllustration mode="login" compact />
        </div>
      </section>

      <section class="flex min-h-0 items-center justify-center overflow-y-auto px-5 py-6 sm:px-8 lg:px-10 lg:py-6 xl:px-14">
        <form novalidate class="w-full max-w-2xl" @submit.prevent="submit">
          <p class="text-[12px] font-black uppercase tracking-[0.42em] text-[#ff3f4b] sm:text-[13px]">Аккаунт</p>
          <h1 class="mt-3 text-4xl font-black tracking-tight sm:text-5xl xl:text-6xl">Войти</h1>

          <p class="mt-3 max-w-xl text-base font-medium leading-7 text-slate-600 xl:text-lg">
            Введите email и пароль, чтобы открыть личный кабинет или операционную панель.
          </p>

          <div class="mt-6 space-y-4 xl:mt-7 xl:space-y-5">
            <label class="block">
              <span class="mb-2 block text-[12px] font-black uppercase tracking-[0.35em] text-slate-500 sm:text-[13px]">Email</span>
              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-5 py-3.5 transition xl:px-6 xl:py-4', errors.email && touched.email ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <Mail class="h-5 w-5 shrink-0 text-slate-400 xl:h-6 xl:w-6" />
                <input
                  v-model="form.email"
                  type="text"
                  autocomplete="email"
                  placeholder="client@example.com"
                  class="w-full min-w-0 bg-transparent text-lg font-black outline-none placeholder:text-slate-400 xl:text-xl"
                  @blur="touched.email = true; validateField('email')"
                  @input="validateField('email')"
                />
              </div>
              <p v-if="errors.email && touched.email" class="mt-2 rounded-2xl bg-red-50 px-4 py-2 text-sm font-bold text-red-600">{{ errors.email }}</p>
            </label>

            <label class="block">
              <span class="mb-2 block text-[12px] font-black uppercase tracking-[0.35em] text-slate-500 sm:text-[13px]">Пароль</span>
              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-5 py-3.5 transition xl:px-6 xl:py-4', errors.password && touched.password ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <LockKeyhole class="h-5 w-5 shrink-0 text-slate-400 xl:h-6 xl:w-6" />
                <input
                  v-model="form.password"
                  type="password"
                  autocomplete="current-password"
                  placeholder="Введите пароль"
                  class="w-full min-w-0 bg-transparent text-lg font-black outline-none placeholder:text-slate-400 xl:text-xl"
                  @blur="touched.password = true; validateField('password')"
                  @input="validateField('password')"
                />
              </div>
              <p v-if="errors.password && touched.password" class="mt-2 rounded-2xl bg-red-50 px-4 py-2 text-sm font-bold text-red-600">{{ errors.password }}</p>
            </label>
          </div>

          <div v-if="errors.common" class="mt-4 rounded-3xl border border-red-200 bg-red-50 px-5 py-4 text-sm font-bold leading-6 text-red-700 xl:mt-5 xl:text-base xl:leading-7">{{ errors.common }}</div>

          <button
            type="submit"
            :disabled="isSubmitting || !isValid"
            class="mt-6 inline-flex min-h-[64px] w-full items-center justify-center gap-3 rounded-3xl bg-[#ff3f4b] px-6 text-lg font-black text-white shadow-[0_20px_45px_rgba(255,63,75,0.28)] transition hover:-translate-y-0.5 hover:bg-[#f12f3c] disabled:cursor-not-allowed disabled:opacity-60 xl:mt-7 xl:min-h-[72px] xl:text-xl"
          >
            {{ isSubmitting ? 'Входим…' : 'Войти' }}
            <ArrowRight class="h-5 w-5 xl:h-6 xl:w-6" />
          </button>

          <p class="mt-5 text-center text-sm font-bold text-slate-500 xl:text-base">
            Нет аккаунта?
            <RouterLink to="/register" class="text-[#ff3f4b] hover:underline">Зарегистрироваться</RouterLink>
          </p>
        </form>
      </section>
    </div>
  </main>
</template>
