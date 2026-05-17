<script setup>
import { RouterLink } from 'vue-router'
import { ArrowLeft, Home, PackageSearch, Radar, ShieldCheck } from '@lucide/vue'
</script>

<template>
  <main class="min-h-screen overflow-hidden bg-[#07101f] text-white">
    <section class="relative min-h-screen bg-[radial-gradient(circle_at_0%_0%,rgba(255,63,73,0.35),transparent_34%),radial-gradient(circle_at_100%_10%,rgba(0,118,150,0.35),transparent_36%),linear-gradient(135deg,#07101f_0%,#101827_48%,#073548_100%)] px-5 py-8 sm:px-8 lg:px-14">
      <div class="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-red-300/50 to-transparent"></div>

      <RouterLink
        to="/"
        class="inline-flex items-center gap-3 rounded-2xl border border-white/15 bg-white/10 px-5 py-3 text-sm font-black text-white transition hover:border-red-300 hover:bg-white/15"
      >
        <ArrowLeft class="h-5 w-5" />
        На главную
      </RouterLink>

      <div class="mx-auto grid min-h-[calc(100vh-7rem)] max-w-7xl items-center gap-12 lg:grid-cols-[0.9fr_1.1fr]">
        <div>
          <p class="mb-5 text-xs font-black uppercase tracking-[0.55em] text-red-300">Ошибка 404</p>
          <h1 class="text-7xl font-black leading-none sm:text-8xl lg:text-9xl">Маршрут не найден</h1>
          <p class="mt-7 max-w-2xl text-xl leading-9 text-slate-200">
            Такой страницы нет или адрес был изменён. Вернитесь на главную страницу Fulfillment Transit или проверьте QR-код ещё раз.
          </p>

          <div class="mt-10 flex flex-col gap-4 sm:flex-row">
            <RouterLink
              to="/"
              class="inline-flex items-center justify-center gap-3 rounded-2xl bg-red-500 px-7 py-5 text-base font-black text-white shadow-2xl shadow-red-500/25 transition hover:-translate-y-0.5 hover:bg-red-400"
            >
              <Home class="h-5 w-5" />
              Вернуться на главную
            </RouterLink>
            <RouterLink
              to="/#qr"
              class="inline-flex items-center justify-center gap-3 rounded-2xl border border-white/25 px-7 py-5 text-base font-black text-white transition hover:-translate-y-0.5 hover:border-red-300 hover:bg-white/10"
            >
              <PackageSearch class="h-5 w-5" />
              Проверить QR
            </RouterLink>
          </div>
        </div>

        <div class="relative">
          <div class="absolute -inset-6 rounded-[3rem] bg-red-500/20 blur-3xl"></div>
          <div class="relative rounded-[3rem] border border-white/15 bg-white/10 p-6 shadow-2xl backdrop-blur-xl sm:p-8">
            <div class="rounded-[2.5rem] bg-[#07101f] p-8">
              <div class="flex items-center justify-between gap-5">
                <div class="flex h-20 w-20 items-center justify-center rounded-[1.75rem] bg-red-500 shadow-xl shadow-red-500/30">
                  <Radar class="h-10 w-10" />
                </div>
                <span class="rounded-full border border-red-300/30 bg-red-500/10 px-4 py-2 text-sm font-black text-red-100">404</span>
              </div>

              <div class="mt-10 space-y-4">
                <div class="rounded-3xl border border-white/10 bg-white/5 p-5">
                  <p class="text-sm font-black uppercase tracking-[0.3em] text-slate-400">Система</p>
                  <p class="mt-3 text-2xl font-black">Fulfillment Transit</p>
                </div>
                <div class="grid gap-4 sm:grid-cols-2">
                  <div class="rounded-3xl border border-white/10 bg-white/5 p-5">
                    <PackageSearch class="mb-5 h-8 w-8 text-red-300" />
                    <p class="font-black">Заявки и QR</p>
                    <p class="mt-2 text-sm leading-6 text-slate-400">Проверьте код на главной странице.</p>
                  </div>
                  <div class="rounded-3xl border border-white/10 bg-white/5 p-5">
                    <ShieldCheck class="mb-5 h-8 w-8 text-red-300" />
                    <p class="font-black">Доступ защищён</p>
                    <p class="mt-2 text-sm leading-6 text-slate-400">Данные доступны после входа.</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
