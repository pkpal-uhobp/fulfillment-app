<script setup>
import { Boxes, CheckCircle2, PackageCheck, QrCode, Truck, Warehouse } from '@lucide/vue'

defineProps({
  mode: {
    type: String,
    default: 'auth',
  },
})
</script>

<template>
  <div class="relative overflow-hidden rounded-[2rem] border border-white/10 bg-[#070f1f] p-8 shadow-2xl">
    <div class="pointer-events-none absolute -left-20 -top-24 h-72 w-72 rounded-full bg-red-500/25 blur-3xl" />
    <div class="pointer-events-none absolute -bottom-20 -right-24 h-72 w-72 rounded-full bg-cyan-500/20 blur-3xl" />

    <div class="relative">
      <div class="mb-8 inline-flex items-center gap-3 rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white">
        <span class="flex h-11 w-11 items-center justify-center rounded-xl bg-[#ff3f4b] shadow-lg shadow-red-500/30">
          <Truck class="h-6 w-6" />
        </span>
        <div>
          <div class="text-lg font-black leading-none">Fulfillment Transit</div>
          <div class="mt-1 text-[10px] font-black uppercase tracking-[0.45em] text-red-200">Marketplace logistics</div>
        </div>
      </div>

      <div class="rounded-[1.75rem] border border-white/10 bg-gradient-to-br from-white/10 to-white/5 p-6">
        <div class="relative min-h-[280px] rounded-[1.5rem] bg-[#081426] p-6">
          <div class="absolute left-8 top-8 h-24 w-44 rounded-3xl border border-white/10 bg-white/10 p-4">
            <div class="mb-4 h-10 w-10 rounded-xl bg-[#ff3f4b]" />
            <div class="mb-3 h-3 w-24 rounded-full bg-white/40" />
            <div class="h-3 w-32 rounded-full bg-white/20" />
          </div>

          <div class="absolute right-10 top-10 flex h-24 w-24 items-center justify-center rounded-3xl border border-white/10 bg-cyan-300/10 text-cyan-50">
            <QrCode class="h-12 w-12" />
          </div>

          <div class="absolute bottom-10 left-10 flex h-24 w-24 items-center justify-center rounded-3xl border border-white/10 bg-red-500/15 text-red-100">
            <Warehouse class="h-12 w-12" />
          </div>

          <div class="absolute bottom-10 right-10 h-28 w-52 rounded-3xl border border-red-300/30 bg-red-400/10 p-5">
            <PackageCheck class="mb-4 h-10 w-10 text-red-100" />
            <div class="mb-3 h-3 w-28 rounded-full bg-white/30" />
            <div class="h-3 w-20 rounded-full bg-white/15" />
          </div>

          <div class="absolute left-1/2 top-1/2 flex -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full border border-white/10 bg-white/10 p-6 text-white backdrop-blur">
            <Boxes class="h-14 w-14" />
          </div>
        </div>

        <div class="mt-5 grid grid-cols-3 gap-3">
          <div class="rounded-2xl bg-white/10 p-4 text-center text-white">
            <QrCode class="mx-auto mb-2 h-5 w-5 text-red-200" />
            <div class="text-xl font-black">QR</div>
            <div class="text-xs text-slate-300">контроль</div>
          </div>
          <div class="rounded-2xl bg-white/10 p-4 text-center text-white">
            <CheckCircle2 class="mx-auto mb-2 h-5 w-5 text-red-200" />
            <div class="text-xl font-black">24/7</div>
            <div class="text-xs text-slate-300">статусы</div>
          </div>
          <div class="rounded-2xl bg-white/10 p-4 text-center text-white">
            <Truck class="mx-auto mb-2 h-5 w-5 text-red-200" />
            <div class="text-xl font-black">API</div>
            <div class="text-xs text-slate-300">заявки</div>
          </div>
        </div>
      </div>

      <div class="mt-8">
        <div class="text-[13px] font-black uppercase tracking-[0.6em] text-red-200">Fulfillment Transit</div>
        <h2 v-if="mode === 'login'" class="mt-5 max-w-xl text-5xl font-black leading-tight text-white">Вход в операционную систему</h2>
        <h2 v-else-if="mode === 'profile'" class="mt-5 max-w-xl text-5xl font-black leading-tight text-white">Личный профиль клиента</h2>
        <h2 v-else class="mt-5 max-w-xl text-5xl font-black leading-tight text-white">Регистрация клиента</h2>
        <p v-if="mode === 'login'" class="mt-5 max-w-xl text-lg leading-8 text-slate-200">Клиент попадает в личный кабинет, логист — в операционную панель, администратор — в управление системой.</p>
        <p v-else-if="mode === 'profile'" class="mt-5 max-w-xl text-lg leading-8 text-slate-200">Здесь собраны данные аккаунта, которые используются для заявок, QR-проверки и истории статусов.</p>
        <p v-else class="mt-5 max-w-xl text-lg leading-8 text-slate-200">Создайте аккаунт, чтобы оформлять заявки, выбирать склады, отслеживать грузовые места и видеть историю статусов.</p>
      </div>
    </div>
  </div>
</template>
