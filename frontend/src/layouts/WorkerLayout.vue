<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { Boxes, Home, LogOut, Menu, PackageCheck, QrCode, UserRound, X } from '@lucide/vue'
import { clearAuth, getCurrentUser, loadMe } from '@/shared/api/http'

const router = useRouter()
const route = useRoute()
const user = ref(getCurrentUser())
const menuOpen = ref(false)

const navItems = [
  { label: 'Главная', to: '/worker', icon: Home },
  { label: 'QR-сканер', to: '/worker/scan', icon: QrCode },
  { label: 'Грузовые места', to: '/worker/cargo-items', icon: Boxes },
  { label: 'Профиль', to: '/worker/profile', icon: UserRound },
]

const displayName = computed(() => user.value?.full_name || user.value?.email || 'Рабочий склада')
const initials = computed(() => {
  const parts = displayName.value.split(' ').filter(Boolean)
  if (!parts.length) return 'Р'
  return parts.slice(0, 2).map((part) => part[0]?.toUpperCase()).join('')
})

function isActive(path) {
  if (path === '/worker') return route.path === '/worker'
  return route.path.startsWith(path)
}

async function refreshUser() {
  try { user.value = await loadMe() } catch { user.value = getCurrentUser() }
}
function logout() { clearAuth(); router.push('/') }
function onAuthChanged() { user.value = getCurrentUser() }

onMounted(() => { window.addEventListener('auth:changed', onAuthChanged); refreshUser() })
onBeforeUnmount(() => window.removeEventListener('auth:changed', onAuthChanged))
</script>

<template>
  <div class="min-h-screen bg-[#07101f] text-white selection:bg-[#ff4248] selection:text-white">
    <div class="fixed inset-0 -z-10 bg-[radial-gradient(circle_at_4%_4%,rgba(255,66,72,0.18),transparent_34%),radial-gradient(circle_at_92%_12%,rgba(0,166,214,0.22),transparent_32%),linear-gradient(135deg,#160711_0%,#07101f_52%,#073e4f_100%)]"></div>

    <header class="sticky top-0 z-40 border-b border-white/10 bg-[#07101f]/92 backdrop-blur-xl">
      <div class="mx-auto flex max-w-[1480px] items-center justify-between gap-4 px-4 py-4 sm:px-6 lg:px-8">
        <div class="flex min-w-0 items-center gap-3">
          <RouterLink to="/worker" class="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-[#ff4248] shadow-[0_18px_44px_rgba(255,66,72,0.28)]">
            <PackageCheck class="h-6 w-6 text-white" />
          </RouterLink>
          <div class="min-w-0">
            <div class="truncate text-lg font-black tracking-[-0.04em]">Fulfillment Transit</div>
            <div class="text-[10px] font-black uppercase tracking-[0.4em] text-[#ff8e92]">Панель рабочего</div>
          </div>
        </div>

        <nav class="hidden items-center gap-2 lg:flex">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-2 rounded-2xl px-4 py-3 text-sm font-black transition"
            :class="isActive(item.to) ? 'bg-white text-[#07101f]' : 'text-white/72 hover:bg-white/10 hover:text-white'"
          >
            <component :is="item.icon" class="h-4 w-4" />
            {{ item.label }}
          </RouterLink>
        </nav>

        <div class="hidden items-center gap-3 md:flex">
          <RouterLink to="/worker/profile" class="flex items-center gap-3 rounded-2xl border border-white/10 bg-white/[0.06] px-4 py-2 transition hover:border-[#ff4248]/60 hover:bg-white/10">
            <div class="grid h-9 w-9 place-items-center rounded-xl bg-[#ff4248] text-sm font-black">{{ initials }}</div>
            <div class="max-w-[190px] truncate text-sm font-black">{{ displayName }}</div>
          </RouterLink>
          <button class="inline-flex items-center gap-2 rounded-2xl border border-white/15 px-4 py-3 text-sm font-black text-white/80 transition hover:bg-white hover:text-[#07101f]" @click="logout">
            Выйти <LogOut class="h-4 w-4" />
          </button>
        </div>

        <button class="grid h-11 w-11 place-items-center rounded-2xl border border-white/15 bg-white/[0.06] lg:hidden" @click="menuOpen = true">
          <Menu class="h-5 w-5" />
        </button>
      </div>
    </header>

    <Teleport to="body">
      <div v-if="menuOpen" class="fixed inset-0 z-50 bg-black/70 p-4 backdrop-blur-sm lg:hidden" @click.self="menuOpen = false">
        <div class="ml-auto flex h-full max-w-sm flex-col rounded-[2rem] border border-white/10 bg-[#07101f] p-5 text-white shadow-2xl">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-lg font-black">Fulfillment Transit</div>
              <div class="text-[10px] font-black uppercase tracking-[0.35em] text-[#ff8e92]">Рабочий</div>
            </div>
            <button class="grid h-11 w-11 place-items-center rounded-2xl bg-white/10" @click="menuOpen = false"><X class="h-5 w-5" /></button>
          </div>
          <nav class="mt-6 grid gap-2">
            <RouterLink
              v-for="item in navItems"
              :key="item.to"
              :to="item.to"
              class="flex items-center gap-3 rounded-2xl px-4 py-4 font-black transition"
              :class="isActive(item.to) ? 'bg-white text-[#07101f]' : 'bg-white/[0.06] text-white/75'"
              @click="menuOpen = false"
            >
              <component :is="item.icon" class="h-5 w-5" />
              {{ item.label }}
            </RouterLink>
          </nav>
          <button class="mt-auto flex items-center justify-center gap-2 rounded-2xl bg-[#ff4248] px-5 py-4 font-black text-white" @click="logout">
            Выйти <LogOut class="h-5 w-5" />
          </button>
        </div>
      </div>
    </Teleport>

    <main class="mx-auto max-w-[1480px] px-4 py-6 sm:px-6 lg:px-8 lg:py-8">
      <RouterView />
    </main>
  </div>
</template>
