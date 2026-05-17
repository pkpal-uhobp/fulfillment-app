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
  <main class="min-h-screen bg-[#07101f] p-5 text-[#07101f] lg:p-8">
    <div class="mx-auto grid min-h-[calc(100vh-2.5rem)] max-w-[1500px] overflow-hidden rounded-[2rem] bg-white shadow-2xl lg:grid-cols-[0.95fr_1.05fr]">
      <section class="bg-gradient-to-br from-[#3a1220] via-[#111827] to-[#073b46] p-8 lg:p-12">
        <RouterLink to="/" class="inline-flex items-center gap-3 rounded-2xl px-1 py-2 text-base font-black text-white/90 hover:text-white">
          <ArrowLeft class="h-5 w-5" />
          На главную
        </RouterLink>
        <div class="mt-12 lg:mt-20">
          <AuthIllustration mode="register" />
        </div>
      </section>

      <section class="flex items-center justify-center px-6 py-12 lg:px-16">
        <form novalidate class="w-full max-w-3xl" @submit.prevent="submit">
          <p class="text-[13px] font-black uppercase tracking-[0.5em] text-[#ff3f4b]">Аккаунт</p>
          <h1 class="mt-4 text-5xl font-black tracking-tight lg:text-6xl">Регистрация</h1>
          <p class="mt-4 max-w-xl text-lg leading-8 text-slate-600">По умолчанию создаётся клиентский аккаунт для оформления заявок.</p>

          <div class="mt-10 space-y-5">
            <label class="block">
              <span class="mb-3 block text-[13px] font-black uppercase tracking-[0.4em] text-slate-500">ФИО</span>
              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-6 py-5 transition', errors.fullName && touched.fullName ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <User class="h-6 w-6 text-slate-400" />
                <input v-model="form.fullName" type="text" autocomplete="name" placeholder="Иванов Иван" class="w-full bg-transparent text-xl font-black outline-none placeholder:text-slate-400" @blur="touched.fullName = true; validateField('fullName')" @input="validateField('fullName')" />
              </div>
              <p v-if="errors.fullName && touched.fullName" class="mt-2 rounded-2xl bg-red-50 px-4 py-3 text-sm font-bold text-red-600">{{ errors.fullName }}</p>
            </label>

            <label class="block">
              <span class="mb-3 block text-[13px] font-black uppercase tracking-[0.4em] text-slate-500">Email</span>
              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-6 py-5 transition', errors.email && touched.email ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <Mail class="h-6 w-6 text-slate-400" />
                <input v-model="form.email" type="text" autocomplete="email" placeholder="client@example.com" class="w-full bg-transparent text-xl font-black outline-none placeholder:text-slate-400" @blur="touched.email = true; validateField('email')" @input="validateField('email')" />
              </div>
              <p v-if="errors.email && touched.email" class="mt-2 rounded-2xl bg-red-50 px-4 py-3 text-sm font-bold text-red-600">{{ errors.email }}</p>
            </label>

            <label class="block">
              <span class="mb-3 block text-[13px] font-black uppercase tracking-[0.4em] text-slate-500">Телефон</span>
              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-6 py-5 transition', errors.phone && touched.phone ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <Phone class="h-6 w-6 text-slate-400" />
                <input v-model="form.phone" type="text" autocomplete="tel" placeholder="+79991234567" class="w-full bg-transparent text-xl font-black outline-none placeholder:text-slate-400" @blur="touched.phone = true; validateField('phone')" @input="validateField('phone')" />
              </div>
              <p v-if="errors.phone && touched.phone" class="mt-2 rounded-2xl bg-red-50 px-4 py-3 text-sm font-bold text-red-600">{{ errors.phone }}</p>
            </label>

            <label class="block">
              <span class="mb-3 block text-[13px] font-black uppercase tracking-[0.4em] text-slate-500">Пароль</span>
              <div :class="['flex items-center gap-4 rounded-3xl border bg-slate-50 px-6 py-5 transition', errors.password && touched.password ? 'border-[#ff3f4b] ring-4 ring-red-100' : 'border-slate-200 focus-within:border-[#ff3f4b] focus-within:ring-4 focus-within:ring-red-100']">
                <LockKeyhole class="h-6 w-6 text-slate-400" />
                <input v-model="form.password" type="password" autocomplete="new-password" placeholder="Минимум 6 символов" class="w-full bg-transparent text-xl font-black outline-none placeholder:text-slate-400" @blur="touched.password = true; validateField('password')" @input="validateField('password')" />
              </div>
              <p v-if="errors.password && touched.password" class="mt-2 rounded-2xl bg-red-50 px-4 py-3 text-sm font-bold text-red-600">{{ errors.password }}</p>
            </label>
          </div>

          <div v-if="errors.common" class="mt-6 rounded-3xl border border-red-200 bg-red-50 px-6 py-5 text-base font-bold leading-7 text-red-700">{{ errors.common }}</div>

          <button type="submit" :disabled="isSubmitting" class="mt-8 flex w-full items-center justify-center gap-4 rounded-3xl bg-[#ff3f4b] px-8 py-6 text-xl font-black text-white shadow-2xl shadow-red-500/25 transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60">
            {{ isSubmitting ? 'Создаём...' : 'Создать аккаунт' }}
            <ArrowRight class="h-6 w-6" />
          </button>

          <p class="mt-8 text-center text-base font-semibold text-slate-500">
            Уже есть аккаунт?
            <RouterLink to="/login" class="font-black text-[#ff3f4b] hover:underline">Войти</RouterLink>
          </p>
        </form>
      </section>
    </div>
  </main>
</template>
