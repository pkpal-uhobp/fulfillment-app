<script setup>
import { computed, reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { ArrowLeft, ArrowRight, LockKeyhole, Mail, Phone, User } from '@lucide/vue'

import AuthIllustration from '@/shared/ui/AuthIllustration.vue'
import { apiFetch, saveAuth } from '@/shared/api/http'
import { humanizeApiError, normalizeEmail, normalizePhone, validateEmail, validateFullName, validatePassword, validatePhone } from './authValidation'

const router = useRouter()

const form = reactive({ fullName: '', email: '', phone: '', password: '' })
const errors = reactive({ fullName: '', email: '', phone: '', password: '', common: '' })
const touched = reactive({ fullName: false, email: false, phone: false, password: false })
const isSubmitting = ref(false)

const isValid = computed(() => !validateFullName(form.fullName) && !validateEmail(form.email) && !validatePhone(form.phone) && !validatePassword(form.password))

function validateField(field) {
  if (field === 'fullName') errors.fullName = validateFullName(form.fullName)
  if (field === 'email') errors.email = validateEmail(form.email)
  if (field === 'phone') errors.phone = validatePhone(form.phone)
  if (field === 'password') errors.password = validatePassword(form.password)
}

function validateForm() {
  Object.keys(touched).forEach((key) => { touched[key] = true })

  validateField('fullName')
  validateField('email')
  validateField('phone')
  validateField('password')
  errors.common = ''

  return !errors.fullName && !errors.email && !errors.phone && !errors.password
}

async function submit() {
  if (!validateForm()) return

  isSubmitting.value = true
  errors.common = ''

  try {
    const payload = await apiFetch('/auth/register', {
      method: 'POST',
      body: {
        full_name: form.fullName.trim().replace(/\s+/g, ' '),
        email: normalizeEmail(form.email),
        phone: normalizePhone(form.phone),
        password: form.password,
      },
    })

    saveAuth(payload)
    await router.push('/client')
  } catch (error) {
    errors.common = humanizeApiError(error, 'Не удалось создать аккаунт.')
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
          <AuthIllustration mode="register" compact />
        </div>
      </section>

      <section class="flex min-h-0 items-center justify-center overflow-y-auto px-5 py-5 sm:px-8 lg:px-10 lg:py-5 xl:px-14">
        <form novalidate class="w-full max-w-2xl" @submit.prevent="submit">
          <p class="text-[12px] font-black uppercase tracking-[0.42em] text-[#ff3f4b] sm:text-[13px]">Аккаунт</p>

          <h1 class="mt-2 text-4xl font-black tracking-tight sm:text-5xl xl:text-6xl">Регистрация</h1>

          <p class="mt-2 max-w-xl text-base font-medium leading-7 text-slate-600 xl:text-lg">
            По умолчанию создаётся клиентский аккаунт для оформления заявок.
          </p>

          <div class="mt-5 space-y-3 xl:mt-6 xl:space-y-4">
            <label class="block">
              <span class="mb-2 block text-[12px] font-black uppercase tracking-[0.35em] text-slate-500 sm:text-[13px]">ФИО</span>

              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-5 py-3 transition xl:px-6 xl:py-3.5', errors.fullName && touched.fullName ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <User class="h-5 w-5 shrink-0 text-slate-400 xl:h-6 xl:w-6" />
                <input
                  v-model="form.fullName"
                  type="text"
                  autocomplete="name"
                  placeholder="Иванов Иван"
                  class="w-full min-w-0 bg-transparent text-lg font-black outline-none placeholder:text-slate-400 xl:text-xl"
                  @blur="touched.fullName = true; validateField('fullName')"
                  @input="validateField('fullName')"
                />
              </div>

              <p v-if="errors.fullName && touched.fullName" class="mt-2 rounded-2xl bg-red-50 px-4 py-2 text-sm font-bold text-red-600">{{ errors.fullName }}</p>
            </label>

            <label class="block">
              <span class="mb-2 block text-[12px] font-black uppercase tracking-[0.35em] text-slate-500 sm:text-[13px]">Email</span>

              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-5 py-3 transition xl:px-6 xl:py-3.5', errors.email && touched.email ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
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
              <span class="mb-2 block text-[12px] font-black uppercase tracking-[0.35em] text-slate-500 sm:text-[13px]">Телефон</span>

              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-5 py-3 transition xl:px-6 xl:py-3.5', errors.phone && touched.phone ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <Phone class="h-5 w-5 shrink-0 text-slate-400 xl:h-6 xl:w-6" />
                <input
                  v-model="form.phone"
                  type="text"
                  autocomplete="tel"
                  placeholder="+79991234567"
                  class="w-full min-w-0 bg-transparent text-lg font-black outline-none placeholder:text-slate-400 xl:text-xl"
                  @blur="touched.phone = true; validateField('phone')"
                  @input="validateField('phone')"
                />
              </div>

              <p v-if="errors.phone && touched.phone" class="mt-2 rounded-2xl bg-red-50 px-4 py-2 text-sm font-bold text-red-600">{{ errors.phone }}</p>
            </label>

            <label class="block">
              <span class="mb-2 block text-[12px] font-black uppercase tracking-[0.35em] text-slate-500 sm:text-[13px]">Пароль</span>

              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-5 py-3 transition xl:px-6 xl:py-3.5', errors.password && touched.password ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <LockKeyhole class="h-5 w-5 shrink-0 text-slate-400 xl:h-6 xl:w-6" />
                <input
                  v-model="form.password"
                  type="password"
                  autocomplete="new-password"
                  placeholder="Минимум 6 символов"
                  class="w-full min-w-0 bg-transparent text-lg font-black outline-none placeholder:text-slate-400 xl:text-xl"
                  @blur="touched.password = true; validateField('password')"
                  @input="validateField('password')"
                />
              </div>

              <p v-if="errors.password && touched.password" class="mt-2 rounded-2xl bg-red-50 px-4 py-2 text-sm font-bold text-red-600">{{ errors.password }}</p>
            </label>
          </div>

          <div v-if="errors.common" class="mt-4 rounded-3xl border border-red-200 bg-red-50 px-5 py-4 text-sm font-bold leading-6 text-red-700 xl:text-base xl:leading-7">{{ errors.common }}</div>

          <button type="submit" :disabled="isSubmitting" class="mt-5 flex w-full items-center justify-center gap-4 rounded-3xl bg-[#ff3f4b] px-7 py-4 text-lg font-black text-white shadow-2xl shadow-red-500/25 transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 xl:mt-6 xl:px-8 xl:py-4 xl:text-xl">
            {{ isSubmitting ? 'Создаём...' : 'Создать аккаунт' }}
            <ArrowRight class="h-6 w-6" />
          </button>

          <p class="mt-4 text-center text-base font-semibold text-slate-500 xl:mt-5">
            Уже есть аккаунт?
            <RouterLink to="/login" class="font-black text-[#ff3f4b] hover:underline">Войти</RouterLink>
          </p>
        </form>
      </section>
    </div>
  </main>
</template>
