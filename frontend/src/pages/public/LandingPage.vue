<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import {
  ArrowRight,
  CheckCircle2,
  ChevronDown,
  Loader2,
  LogOut,
  Menu,
  PackageCheck,
  QrCode,
  ShieldCheck,
  Truck,
  User,
  Warehouse,
  X,
} from '@lucide/vue'
import QrStatusPanel from '@/features/scanQr/QrStatusPanel.vue'
import { apiFetch, clearAuth, getAccessToken, getCurrentUser, loadMe } from '@/shared/api/http'

const router = useRouter()

const navLinks = [
  { label: 'Услуги', href: '#services' },
  { label: 'Процесс', href: '#process' },
  { label: 'Склады', href: '#warehouses' },
  { label: 'Преимущества', href: '#benefits' },
  { label: 'FAQ', href: '#faq' },
]

const services = [
  {
    icon: PackageCheck,
    title: 'Приёмка товара',
    text: 'Проверяем количество мест, фиксируем параметры, принимаем товар и создаём понятную историю обработки.',
  },
  {
    icon: QrCode,
    title: 'QR-контроль мест',
    text: 'Каждое грузовое место получает QR-код, по которому можно быстро открыть карточку товара и его статус.',
  },
  {
    icon: Warehouse,
    title: 'Хранение по зонам',
    text: 'Логист назначает зону хранения, а складской работник видит, куда нужно переместить товар.',
  },
  {
    icon: Truck,
    title: 'Отправка на склад назначения',
    text: 'Грузовые места закрепляются за гейтом и передаются на нужный склад маркетплейса или клиента.',
  },
]

const processSteps = [
  ['01', 'Заявка', 'Клиент оформляет заявку, выбирает склад приёмки, склад назначения и способ передачи товара.'],
  ['02', 'Календарь', 'Система проверяет доступную дату и помогает не перегружать складскую приёмку.'],
  ['03', 'Приёмка', 'Работник принимает грузовые места, присваивает QR и фиксирует фактические параметры.'],
  ['04', 'Хранение', 'Логист распределяет товар по зонам хранения и назначает дальнейшие операции.'],
  ['05', 'Отправка', 'Груз перемещается к гейту и отправляется на склад назначения с прозрачной историей статусов.'],
]

const benefits = [
  'статус каждого грузового места доступен по QR-коду',
  'история действий хранится в системе',
  'клиенты видят свои заявки без звонков менеджеру',
  'логист управляет календарём, зонами и гейтами',
  'склад работает по понятным операциям',
  'администратор управляет пользователями и складами',
]

const faqItems = ref([
  {
    question: 'Можно ли проверить товар по QR-коду?',
    answer: 'Да. После входа в аккаунт пользователь вводит QR-код или сканирует его камерой. Система открывает отдельную карточку грузового места.',
    open: true,
  },
  {
    question: 'Почему перед сканированием нужен вход?',
    answer: 'QR-код связан с заявкой и грузовым местом. Backend проверяет роль пользователя и не показывает чужие данные.',
    open: false,
  },
  {
    question: 'Что увидит клиент?',
    answer: 'Клиент видит только свои грузовые места: QR-код, текущий статус, заявку и разрешённые складские данные.',
    open: false,
  },
  {
    question: 'Что делает логист?',
    answer: 'Логист управляет календарём приёмки, назначает зоны хранения, гейты и формирует отправку товара на склад назначения.',
    open: false,
  },
  {
    question: 'Что делает складской работник?',
    answer: 'Работник принимает товар, работает с QR-кодами, меняет статусы и выполняет операции по размещению или отправке.',
    open: false,
  },
  {
    question: 'Какие склады отображаются на сайте?',
    answer: 'Склады загружаются из базы данных через публичный endpoint. Универсальный склад может использоваться и для приёмки, и как склад назначения.',
    open: false,
  },
  {
    question: 'Можно ли закрыть дату приёмки?',
    answer: 'Да. Логист или администратор может закрывать даты и ограничивать количество заявок для конкретного склада.',
    open: false,
  },
])

const warehouses = ref([])
const warehousesLoading = ref(false)
const warehousesError = ref('')
const mobileMenuOpen = ref(false)
const desktopMoreOpen = ref(false)
const headerTheme = ref('dark')
const currentUser = ref(getCurrentUser())
const authVersion = ref(0)

const isAuthorized = computed(() => {
  authVersion.value
  return Boolean(currentUser.value || getAccessToken())
})
const userDisplayName = computed(() => {
  const user = currentUser.value
  if (!user) return ''
  return user.full_name || user.fullName || user.name || user.email || 'Пользователь'
})

const activeWarehouses = computed(() => warehouses.value.filter((item) => item.is_active !== false))
const warehouseCards = computed(() => activeWarehouses.value.slice(0, 6))
const warehousesCount = computed(() => activeWarehouses.value.length || '—')

const headerClass = computed(() => {
  return headerTheme.value === 'light'
    ? 'border-slate-200 bg-white/95 text-slate-950 shadow-[0_12px_50px_rgba(15,23,42,0.10)] backdrop-blur-xl'
    : 'border-white/10 bg-[#07101f]/96 text-white shadow-[0_14px_55px_rgba(0,0,0,0.30)] backdrop-blur-xl'
})

const navClass = computed(() => {
  return headerTheme.value === 'light'
    ? 'text-slate-700 hover:text-slate-950'
    : 'text-white/75 hover:text-white'
})

function shortWarehouseName(warehouse) {
  if (!warehouse) return '—'
  return String(warehouse.name || '')
    .replace(/^TransitPro\s+/i, '')
    .replace(/^Fulfillment Transit\s+/i, '')
    .replace(/\s+/g, ' ')
    .trim() || warehouse.city || 'Склад'
}

function warehouseTypeLabel(type) {
  const labels = {
    receiving: 'Склад приёмки',
    destination: 'Склад назначения',
    both: 'Универсальный склад',
  }

  return labels[type] || 'Склад'
}

async function fetchWarehouses() {
  warehousesLoading.value = true
  warehousesError.value = ''

  try {
    const data = await apiFetch('/warehouses')
    warehouses.value = Array.isArray(data?.warehouses)
      ? data.warehouses
      : Array.isArray(data)
        ? data
        : []
  } catch (error) {
    warehousesError.value = 'Не удалось загрузить склады из базы данных. Проверьте backend и Vite proxy.'
    warehouses.value = []
  } finally {
    warehousesLoading.value = false
  }
}

async function refreshUser() {
  if (!getAccessToken()) {
    currentUser.value = null
    return
  }

  try {
    currentUser.value = await loadMe()
  } catch {
    currentUser.value = getCurrentUser()
  }
}

function logout() {
  clearAuth()
  currentUser.value = null
  authVersion.value += 1
  router.push('/')
}

function toggleFaq(index) {
  faqItems.value[index].open = !faqItems.value[index].open
}

function scrollToStatus() {
  document.querySelector('#status')?.scrollIntoView({ behavior: 'smooth' })
}

function updateHeaderTheme() {
  const probeY = 112
  const sections = Array.from(document.querySelectorAll('[data-header-theme]'))
  const activeSection = sections.find((section) => {
    const rect = section.getBoundingClientRect()
    return rect.top <= probeY && rect.bottom > probeY
  })

  headerTheme.value = activeSection?.dataset.headerTheme || 'dark'
}

function onAuthChanged() {
  authVersion.value += 1
  currentUser.value = getCurrentUser()
  refreshUser()
}

onMounted(() => {
  fetchWarehouses()
  refreshUser()
  updateHeaderTheme()
  window.addEventListener('scroll', updateHeaderTheme, { passive: true })
  window.addEventListener('resize', updateHeaderTheme)
  window.addEventListener('auth:changed', onAuthChanged)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateHeaderTheme)
  window.removeEventListener('resize', updateHeaderTheme)
  window.removeEventListener('auth:changed', onAuthChanged)
})
</script>

<template>
  <main class="min-h-screen bg-slate-100 font-sans text-slate-950">
    <header
      :class="[
        'fixed left-0 right-0 top-0 z-50 border-b px-4 py-4 transition-all duration-300 lg:px-8',
        headerClass,
      ]"
    >
      <div class="mx-auto flex w-[min(1500px,calc(100%-16px))] items-center justify-between gap-4">
        <RouterLink to="/" class="flex items-center gap-4">
          <span class="flex h-14 w-14 items-center justify-center rounded-2xl bg-[#ff3f49] shadow-[0_18px_45px_rgba(255,63,73,0.35)]">
            <Truck class="h-7 w-7 text-white" />
          </span>
          <span>
            <span class="block text-xl font-black leading-none tracking-[-0.03em]">Fulfillment Transit</span>
            <span class="mt-2 block text-[10px] font-black uppercase tracking-[0.55em] text-[#ff5962]">marketplace logistics</span>
          </span>
        </RouterLink>

        <nav class="hidden items-center gap-8 text-sm font-black xl:flex">
          <a v-for="item in navLinks" :key="item.href" :href="item.href" :class="navClass">
            {{ item.label }}
          </a>
        </nav>

        <div class="hidden items-center gap-3 lg:flex">
          <template v-if="isAuthorized">
            <div class="flex items-center gap-3 rounded-2xl border px-4 py-3" :class="headerTheme === 'light' ? 'border-slate-200 bg-slate-50' : 'border-white/10 bg-white/5'">
              <User class="h-4 w-4 text-[#ff5962]" />
              <span class="max-w-[190px] truncate text-sm font-black">{{ userDisplayName }}</span>
            </div>
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-2xl border px-5 py-3 text-sm font-black transition hover:-translate-y-0.5"
              :class="headerTheme === 'light' ? 'border-slate-200 bg-white text-slate-950' : 'border-white/15 bg-white/5 text-white'"
              @click="logout"
            >
              Выйти
              <LogOut class="h-4 w-4" />
            </button>
          </template>

          <template v-else>
            <RouterLink :class="['rounded-2xl px-5 py-3 text-sm font-black transition hover:-translate-y-0.5', navClass]" to="/login">
              Войти
            </RouterLink>
            <RouterLink
              to="/register"
              class="rounded-2xl bg-[#ff3f49] px-7 py-4 text-sm font-black text-white shadow-[0_18px_45px_rgba(255,63,73,0.35)] transition hover:-translate-y-0.5"
            >
              Регистрация
            </RouterLink>
          </template>
        </div>

        <div class="relative hidden md:block xl:hidden">
          <button
            type="button"
            class="rounded-2xl border px-4 py-3 font-black"
            :class="headerTheme === 'light' ? 'border-slate-200 bg-white text-slate-950' : 'border-white/15 bg-white/5 text-white'"
            @click="desktopMoreOpen = !desktopMoreOpen"
          >
            •••
          </button>

          <div
            v-if="desktopMoreOpen"
            class="absolute right-0 top-14 w-72 rounded-3xl border border-white/10 bg-[#07101f] p-4 text-white shadow-2xl"
          >
            <a
              v-for="item in navLinks"
              :key="item.href"
              :href="item.href"
              class="block rounded-2xl px-4 py-3 text-sm font-black text-white/80 hover:bg-white/10 hover:text-white"
              @click="desktopMoreOpen = false"
            >
              {{ item.label }}
            </a>

            <div class="mt-3 border-t border-white/10 pt-3">
              <template v-if="isAuthorized">
                <div class="rounded-2xl bg-white/5 px-4 py-3 text-sm font-black">{{ userDisplayName }}</div>
                <button class="mt-2 w-full rounded-2xl bg-white px-4 py-3 text-sm font-black text-slate-950" @click="logout">Выйти</button>
              </template>
              <template v-else>
                <RouterLink class="block rounded-2xl px-4 py-3 text-sm font-black text-white/80 hover:bg-white/10" to="/login">Войти</RouterLink>
                <RouterLink class="mt-2 block rounded-2xl bg-[#ff3f49] px-4 py-3 text-center text-sm font-black text-white" to="/register">Регистрация</RouterLink>
              </template>
            </div>
          </div>
        </div>

        <button
          type="button"
          class="rounded-2xl border p-3 md:hidden"
          :class="headerTheme === 'light' ? 'border-slate-200 bg-white text-slate-950' : 'border-white/15 bg-white/5 text-white'"
          @click="mobileMenuOpen = true"
        >
          <Menu class="h-6 w-6" />
        </button>
      </div>
    </header>

    <div v-if="mobileMenuOpen" class="fixed inset-0 z-[60] bg-[#07101f]/90 p-4 backdrop-blur-xl md:hidden">
      <div class="rounded-[32px] border border-white/10 bg-[#0b1324] p-5 text-white">
        <div class="flex items-center justify-between">
          <span class="text-lg font-black">Fulfillment Transit</span>
          <button class="rounded-2xl bg-white/10 p-3" @click="mobileMenuOpen = false">
            <X class="h-5 w-5" />
          </button>
        </div>

        <div class="mt-6 grid gap-2">
          <a
            v-for="item in navLinks"
            :key="item.href"
            :href="item.href"
            class="rounded-2xl px-4 py-3 text-base font-black text-white/80 hover:bg-white/10"
            @click="mobileMenuOpen = false"
          >
            {{ item.label }}
          </a>

          <template v-if="isAuthorized">
            <div class="mt-3 rounded-2xl bg-white/5 px-4 py-3 font-black">{{ userDisplayName }}</div>
            <button class="rounded-2xl bg-white px-4 py-4 font-black text-slate-950" @click="logout">Выйти</button>
          </template>

          <template v-else>
            <RouterLink class="mt-3 rounded-2xl border border-white/15 px-4 py-4 text-center font-black text-white" to="/login">Войти</RouterLink>
            <RouterLink class="rounded-2xl bg-[#ff3f49] px-4 py-4 text-center font-black text-white" to="/register">Регистрация</RouterLink>
          </template>
        </div>
      </div>
    </div>

    <section data-header-theme="dark" class="relative overflow-hidden bg-[#07101f] pt-32 text-white">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_8%_18%,rgba(255,63,73,0.35),transparent_28%),radial-gradient(circle_at_92%_10%,rgba(20,184,166,0.28),transparent_30%),linear-gradient(135deg,#120814_0%,#07101f_50%,#06394b_100%)]"></div>

      <div class="relative mx-auto grid w-[min(1500px,calc(100%-32px))] gap-10 py-14 lg:grid-cols-[0.96fr_1.04fr] lg:py-24">
        <div class="flex flex-col justify-center">
          <p class="mb-7 text-sm font-black uppercase tracking-[0.45em] text-[#ff9aa0]">цифровой фулфилмент</p>
          <h1 class="max-w-[760px] text-5xl font-black leading-[0.95] tracking-[-0.07em] sm:text-6xl lg:text-7xl">
            Приёмка, хранение и отправка товара под контролем
          </h1>
          <p class="mt-8 max-w-[760px] text-xl leading-10 text-white/88">
            Организуем полный путь товара: заявка, календарь приёмки, QR-контроль каждого места,
            хранение по зонам и отправка на склад назначения.
          </p>

          <div class="mt-10 flex flex-col gap-4 sm:flex-row">
            <RouterLink
              to="/register"
              class="inline-flex items-center justify-center gap-3 rounded-[28px] bg-[#ff3f49] px-8 py-5 text-base font-black text-white shadow-[0_22px_65px_rgba(255,63,73,0.38)] transition hover:-translate-y-1"
            >
              Оформить заявку
              <ArrowRight class="h-5 w-5" />
            </RouterLink>
            <button
              type="button"
              class="inline-flex items-center justify-center gap-3 rounded-[28px] border border-white/25 px-8 py-5 text-base font-black text-white transition hover:-translate-y-1 hover:bg-white/10"
              @click="scrollToStatus"
            >
              Проверить QR
            </button>
          </div>

          <div class="mt-12 grid max-w-[720px] grid-cols-3 gap-4">
            <div class="rounded-3xl border border-white/25 px-5 py-6 text-center">
              <div class="text-3xl font-black">24/7</div>
              <div class="mt-2 text-xs text-white/60">доступ к статусам</div>
            </div>
            <div class="rounded-3xl border border-white/25 px-5 py-6 text-center">
              <div class="text-3xl font-black">QR</div>
              <div class="mt-2 text-xs text-white/60">каждое место</div>
            </div>
            <div class="rounded-3xl border border-white/25 px-5 py-6 text-center">
              <div class="text-3xl font-black">{{ warehousesCount }}</div>
              <div class="mt-2 text-xs text-white/60">складов в системе</div>
            </div>
          </div>
        </div>

        <div id="status" class="rounded-[40px] border border-white/10 bg-white/5 p-5 shadow-2xl backdrop-blur">
          <QrStatusPanel />
        </div>
      </div>
    </section>

    <section id="services" data-header-theme="light" class="bg-slate-100 py-24">
      <div class="mx-auto w-[min(1500px,calc(100%-32px))]">
        <div class="mb-12 flex flex-col justify-between gap-5 lg:flex-row lg:items-end">
          <div>
            <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff3f49]">услуги</p>
            <h2 class="mt-5 max-w-3xl text-5xl font-black leading-none tracking-[-0.06em]">Что берём на себя</h2>
          </div>
          <p class="max-w-xl text-lg leading-8 text-slate-600">
            Сервис закрывает полный цикл обработки товара для продавцов маркетплейсов.
          </p>
        </div>

        <div class="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
          <article
            v-for="service in services"
            :key="service.title"
            class="flex min-h-[320px] flex-col rounded-[34px] border border-slate-200 bg-white p-8 shadow-[0_18px_50px_rgba(15,23,42,0.06)]"
          >
            <component :is="service.icon" class="h-10 w-10 text-[#ff3f49]" />
            <h3 class="mt-14 text-3xl font-black leading-tight tracking-[-0.04em]">{{ service.title }}</h3>
            <p class="mt-5 leading-8 text-slate-600">{{ service.text }}</p>
          </article>
        </div>
      </div>
    </section>

    <section id="process" data-header-theme="light" class="bg-slate-100 py-24">
      <div class="mx-auto grid w-[min(1500px,calc(100%-32px))] gap-12 lg:grid-cols-[0.8fr_1.2fr]">
        <div>
          <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff3f49]">процесс</p>
          <h2 class="mt-5 text-5xl font-black leading-none tracking-[-0.06em]">От заявки до отправки</h2>
          <p class="mt-8 text-lg leading-9 text-slate-600">
            Все этапы фиксируются в системе: клиент видит статус, склад работает по QR-кодам,
            логист управляет календарём, зонами и гейтами.
          </p>
        </div>

        <div class="grid gap-5">
          <article
            v-for="[number, title, text] in processSteps"
            :key="number"
            class="grid gap-4 rounded-[34px] border border-slate-200 bg-white p-8 shadow-sm sm:grid-cols-[120px_1fr]"
          >
            <div class="text-4xl font-black text-slate-200">{{ number }}</div>
            <div>
              <h3 class="text-2xl font-black">{{ title }}</h3>
              <p class="mt-3 text-lg leading-8 text-slate-600">{{ text }}</p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <section id="warehouses" data-header-theme="light" class="bg-white py-24">
      <div class="mx-auto w-[min(1500px,calc(100%-32px))]">
        <div class="mb-12 flex flex-col justify-between gap-5 lg:flex-row lg:items-end">
          <div>
            <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff3f49]">география</p>
            <h2 class="mt-5 text-5xl font-black leading-none tracking-[-0.06em]">Склады из базы данных</h2>
          </div>
          <p class="max-w-xl text-lg leading-8 text-slate-600">
            Карточки строятся из публичного endpoint склада. Универсальный склад используется и для приёмки, и как конечное место.
          </p>
        </div>

        <div v-if="warehousesLoading" class="flex items-center gap-3 rounded-[32px] border border-slate-200 bg-slate-50 p-8 text-slate-600">
          <Loader2 class="h-5 w-5 animate-spin" />
          Загружаем склады...
        </div>

        <div v-else-if="warehousesError" class="rounded-[32px] border border-amber-200 bg-amber-50 p-8 text-amber-900">
          {{ warehousesError }}
        </div>

        <div v-else class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          <article
            v-for="warehouse in warehouseCards"
            :key="warehouse.id"
            class="flex min-h-[300px] flex-col rounded-[34px] border border-slate-200 bg-slate-50 p-8"
          >
            <Warehouse class="h-10 w-10 text-[#ff3f49]" />
            <h3 class="mt-10 text-3xl font-black leading-tight tracking-[-0.04em]">{{ shortWarehouseName(warehouse) }}</h3>
            <p class="mt-3 text-lg font-black text-slate-700">{{ warehouse.city }}</p>
            <p class="mt-4 leading-7 text-slate-600">{{ warehouse.address }}</p>
            <div class="mt-auto pt-8">
              <span class="inline-flex rounded-2xl bg-white px-4 py-3 text-sm font-black text-slate-900 ring-1 ring-slate-200">
                {{ warehouse.marketplace || warehouseTypeLabel(warehouse.warehouse_type) }}
              </span>
            </div>
          </article>
        </div>
      </div>
    </section>

    <section id="benefits" data-header-theme="dark" class="bg-[#07101f] py-24 text-white">
      <div class="mx-auto grid w-[min(1500px,calc(100%-32px))] gap-12 lg:grid-cols-[0.9fr_1.1fr]">
        <div>
          <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff9aa0]">преимущества</p>
          <h2 class="mt-5 text-5xl font-black leading-none tracking-[-0.06em]">Контроль без таблиц и ручных сверок</h2>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div v-for="item in benefits" :key="item" class="rounded-[28px] border border-white/10 bg-white/[0.05] p-6">
            <ShieldCheck class="h-7 w-7 text-[#ff5962]" />
            <p class="mt-5 text-lg font-black leading-7">{{ item }}</p>
          </div>
        </div>
      </div>
    </section>

    <section id="faq" data-header-theme="light" class="bg-slate-100 py-24">
      <div class="mx-auto w-[min(1050px,calc(100%-32px))]">
        <p class="text-center text-sm font-black uppercase tracking-[0.45em] text-[#ff3f49]">faq</p>
        <h2 class="mt-5 text-center text-5xl font-black leading-none tracking-[-0.06em]">Частые вопросы</h2>

        <div class="mt-12 grid gap-5">
          <article v-for="(item, index) in faqItems" :key="item.question" class="overflow-hidden rounded-[32px] border border-slate-200 bg-white shadow-sm">
            <button type="button" class="flex w-full items-center justify-between gap-5 px-7 py-6 text-left" @click="toggleFaq(index)">
              <span class="text-xl font-black">{{ item.question }}</span>
              <ChevronDown class="h-5 w-5 shrink-0 text-[#ff3f49] transition" :class="{ 'rotate-180': item.open }" />
            </button>
            <div v-if="item.open" class="px-7 pb-7 text-lg leading-8 text-slate-600">
              {{ item.answer }}
            </div>
          </article>
        </div>
      </div>
    </section>

    <footer data-header-theme="dark" class="bg-[#07101f] px-4 py-12 text-white">
      <div class="mx-auto w-[min(1500px,calc(100%-32px))]">
        <div class="rounded-[40px] bg-[#ff3f49] p-10 shadow-[0_20px_70px_rgba(255,63,73,0.28)] lg:p-16">
          <div class="flex flex-col justify-between gap-8 lg:flex-row lg:items-center">
            <div>
              <h2 class="max-w-3xl text-4xl font-black leading-none tracking-[-0.05em] lg:text-5xl">
                Готовы запустить управляемый фулфилмент?
              </h2>
              <p class="mt-6 max-w-2xl text-lg leading-8 text-white/90">
                Создайте аккаунт, оформите заявку и отслеживайте каждый этап обработки товара.
              </p>
            </div>
            <RouterLink to="/register" class="inline-flex shrink-0 items-center justify-center gap-3 rounded-[28px] bg-white px-8 py-5 text-base font-black text-[#bd111b]">
              Начать работу
              <ArrowRight class="h-5 w-5" />
            </RouterLink>
          </div>
        </div>

        <div class="mt-12 flex flex-col justify-between gap-6 border-t border-white/10 pt-8 text-sm text-white/60 md:flex-row">
          <p>© 2026 Fulfillment Transit. Все права защищены.</p>
          <p>Москва · Казань · Санкт-Петербург · маркетплейсы</p>
        </div>
      </div>
    </footer>
  </main>
</template>
