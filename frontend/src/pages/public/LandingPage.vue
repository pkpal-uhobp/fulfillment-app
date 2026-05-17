<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { clearAuth, getAccessToken, getCurrentUser, loadMe } from '@/shared/api/http'
import {
  ArrowRight,
  Building2,
  CalendarCheck,
  CheckCircle2,
  ChevronDown,
  ClipboardList,
  Clock3,
  MapPin,
  Menu,
  MoreHorizontal,
  PackageCheck,
  QrCode,
  ShieldCheck,
  Truck,
  X,
} from '@lucide/vue'

const router = useRouter()

const rawBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1'
const apiBaseUrl = import.meta.env.DEV && rawBaseUrl.includes(':8080')
  ? '/api/v1'
  : rawBaseUrl.replace(/\/$/, '')

const warehouses = ref([])
const warehousesLoading = ref(false)
const warehousesError = ref('')

const selectedReceivingId = ref(null)
const selectedDestinationId = ref(null)
const openedDropdown = ref(null)
const trackingCode = ref('QR-TPRO-MSK-240001')
const trackingResult = ref(null)
const trackingLoading = ref(false)
const trackingError = ref('')
const openedFaq = ref(0)

const mobileMenuOpen = ref(false)
const tabletMenuOpen = ref(false)
const headerTheme = ref('dark')

const currentUser = ref(getCurrentUser())

const isAuthenticated = computed(() => Boolean(currentUser.value || getAccessToken()))

const userDisplayName = computed(() => {
  const user = currentUser.value
  if (!user) return ''
  return user.full_name || user.fullName || user.email || 'Пользователь'
})

const userInitials = computed(() => {
  const name = userDisplayName.value.trim()
  if (!name) return 'FT'
  const parts = name.split(/\s+/).slice(0, 2)
  return parts.map((part) => part[0]).join('').toUpperCase()
})

function refreshAuthState() {
  currentUser.value = getCurrentUser()
}

async function syncCurrentUser() {
  if (!getAccessToken()) {
    currentUser.value = null
    return
  }

  if (currentUser.value) return

  try {
    currentUser.value = await loadMe()
  } catch (error) {
    clearAuth()
    currentUser.value = null
  }
}

function logout() {
  clearAuth()
  currentUser.value = null
  trackingResult.value = null
  closeMenus()
  router.push('/')
}

function warehouseLabel(warehouse) {
  if (!warehouse) return 'Выберите склад'

  if (warehouse.marketplace) {
    return `${warehouse.marketplace} · ${warehouse.name.replace(`${warehouse.marketplace} `, '')}`
  }

  return warehouse.name
    .replace('TransitPro ', '')
    .replace('Fulfillment Transit ', '')
}

function warehouseKindLabel(type) {
  const map = {
    both: 'универсальный склад',
    receiving: 'склад приёмки',
    destination: 'склад назначения',
  }

  return map[type] || 'склад'
}

function statusLabel(status) {
  const map = {
    accepted: 'Принято на склад',
    stored: 'На хранении',
    ready_to_ship: 'Готово к отправке',
    shipped: 'Отправлено',
    lost: 'Утеряно',
    damaged: 'Повреждено',
    cancelled: 'Отменено',
    created: 'Заявка создана',
    waiting_pickup: 'Ожидает забора',
    waiting_delivery: 'Ожидает сдачи на склад',
    received: 'Принято',
    assigned_to_shipping: 'Назначено к отправке',
    delivered: 'Доставлено',
  }

  return map[status] || status || 'Статус уточняется'
}

function normalizeWarehouses(payload) {
  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.warehouses)) return payload.warehouses
  if (Array.isArray(payload?.items)) return payload.items
  if (Array.isArray(payload?.data)) return payload.data
  return []
}

async function loadWarehouses() {
  warehousesLoading.value = true
  warehousesError.value = ''

  try {
    const response = await fetch(`${apiBaseUrl}/warehouses`, {
      headers: { Accept: 'application/json' },
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }

    const payload = await response.json()
    const items = normalizeWarehouses(payload)
      .filter((warehouse) => warehouse && warehouse.is_active !== false)
      .sort((a, b) => Number(a.id || 0) - Number(b.id || 0))

    warehouses.value = items

    const firstReceiving = items.find((warehouse) => ['receiving', 'both'].includes(warehouse.warehouse_type))
    const firstDestination = items.find((warehouse) => ['destination', 'both'].includes(warehouse.warehouse_type))

    selectedReceivingId.value = firstReceiving?.id ?? null
    selectedDestinationId.value = firstDestination?.id ?? null
  } catch (error) {
    warehouses.value = []
    warehousesError.value = 'Не удалось загрузить склады из базы данных. Проверьте, что backend запущен и frontend использует Vite proxy.'
  } finally {
    warehousesLoading.value = false
  }
}

const receivingWarehouses = computed(() => warehouses.value.filter((warehouse) => (
  warehouse.warehouse_type === 'receiving' || warehouse.warehouse_type === 'both'
)))

const destinationWarehouses = computed(() => warehouses.value.filter((warehouse) => (
  warehouse.warehouse_type === 'destination' || warehouse.warehouse_type === 'both'
)))

const selectedReceivingWarehouse = computed(() => (
  warehouses.value.find((warehouse) => warehouse.id === selectedReceivingId.value) || null
))

const selectedDestinationWarehouse = computed(() => (
  warehouses.value.find((warehouse) => warehouse.id === selectedDestinationId.value) || null
))

const warehouseCards = computed(() => warehouses.value)

const stats = computed(() => [
  { value: '24/7', label: 'доступ к статусам' },
  { value: 'QR', label: 'учёт каждого места' },
  { value: warehouses.value.length ? `${warehouses.value.length}` : '—', label: 'складов в сети' },
])

const services = [
  {
    icon: ClipboardList,
    title: 'Приёмка товара',
    text: 'Проверяем количество мест, маркировку, вес и габариты при сдаче на терминал или после забора груза.',
  },
  {
    icon: QrCode,
    title: 'QR-контроль мест',
    text: 'Каждое грузовое место получает QR-код и проходит через понятную историю статусов.',
  },
  {
    icon: PackageCheck,
    title: 'Хранение и сортировка',
    text: 'Распределяем товар по зонам хранения, контролируем перемещение и подготовку к отправке.',
  },
  {
    icon: Truck,
    title: 'Отправка на склады',
    text: 'Формируем партии, закрепляем места за гейтами и отправляем груз на склад назначения.',
  },
]

const processSteps = [
  ['01', 'Заявка', 'Клиент выбирает склад приёмки, склад назначения, дату и формат передачи товара.'],
  ['02', 'Календарь', 'Система проверяет доступные даты и лимиты приёмки, чтобы не перегружать склад.'],
  ['03', 'Приёмка', 'Работник принимает грузовые места, присваивает QR и фиксирует фактические параметры.'],
  ['04', 'Хранение', 'Логист распределяет места по зонам хранения и назначает гейт для отправки.'],
  ['05', 'Отправка', 'Груз добавляется в партию, отправляется на склад назначения и закрывается по статусам.'],
]

const benefits = [
  'Прозрачный путь каждого грузового места от клиента до склада назначения.',
  'Единая система вместо таблиц, чатов и ручных согласований.',
  'Контроль дат приёмки, зон хранения, гейтов и истории статусов.',
]

const faqs = [
  {
    q: 'Как начать работу с Fulfillment Transit?',
    a: 'Зарегистрируйтесь, войдите в аккаунт и создайте заявку: выберите склад приёмки, склад назначения, дату передачи товара и количество грузовых мест.',
  },
  {
    q: 'Можно ли заранее закрыть дату приёмки?',
    a: 'Да. Логист или администратор может закрыть дату либо ограничить количество заявок для конкретного склада, чтобы не перегружать приёмку.',
  },
  {
    q: 'Как клиент отслеживает товар?',
    a: 'После входа в аккаунт клиент видит свои заявки, грузовые места, QR-коды, текущие статусы и историю изменений.',
  },
  {
    q: 'Что показывает QR-код?',
    a: 'QR-код открывает карточку грузового места: номер заявки, текущий статус, склад, зону хранения, назначенный гейт и историю операций.',
  },
  {
    q: 'Какие склады отображаются в заявке?',
    a: 'В списке складов приёмки показываются склады приёмки и универсальные склады. В списке назначения — склады назначения и универсальные склады.',
  },
  {
    q: 'Кто меняет статусы грузовых мест?',
    a: 'Складской работник, логист или администратор меняют статусы в зависимости от операции: приёмка, хранение, подготовка к отправке или отправка.',
  },
  {
    q: 'Можно ли посмотреть статус без входа?',
    a: 'Публичная форма показывает выбранный маршрут и предлагает войти. Полная карточка QR-кода защищена и доступна после авторизации с подходящей ролью.',
  },
]

function toggleFaq(index) {
  openedFaq.value = openedFaq.value === index ? null : index
}

function toggleDropdown(name) {
  openedDropdown.value = openedDropdown.value === name ? null : name
}

function selectWarehouse(type, warehouse) {
  if (type === 'receiving') selectedReceivingId.value = warehouse.id
  if (type === 'destination') selectedDestinationId.value = warehouse.id
  openedDropdown.value = null
}

function closeMenus() {
  mobileMenuOpen.value = false
  tabletMenuOpen.value = false
}

async function checkTracking() {
  trackingError.value = ''
  trackingResult.value = null

  const code = trackingCode.value.trim()
  if (!code) {
    trackingError.value = 'Введите номер заявки или QR-код.'
    return
  }

  const token = getAccessToken()

  if (!token) {
    trackingResult.value = {
      title: 'Маршрут подготовлен',
      subtitle: 'Для просмотра полной карточки войдите в аккаунт',
      description: 'Публичная форма показывает выбранный маршрут. Подробный статус, история операций и карточка грузового места доступны после входа.',
      code,
      status: 'Требуется вход',
    }
    return
  }

  trackingLoading.value = true

  try {
    const response = await fetch(`${apiBaseUrl}/cargo-items/scan?qr_code=${encodeURIComponent(code)}`, {
      headers: {
        Accept: 'application/json',
        Authorization: `Bearer ${token}`,
      },
    })

    if (!response.ok) throw new Error(`HTTP ${response.status}`)

    const payload = await response.json()
    trackingResult.value = {
      title: 'Грузовое место найдено',
      subtitle: statusLabel(payload.status || payload.cargo_item?.status),
      description: 'QR-код распознан. Полная карточка места доступна в личном кабинете.',
      code,
      status: statusLabel(payload.status || payload.cargo_item?.status),
    }
  } catch (error) {
    trackingResult.value = {
      title: 'Нужен вход с подходящей ролью',
      subtitle: 'Проверка QR защищена',
      description: 'Сканирование QR доступно сотрудникам склада, логисту и администратору. Войдите в аккаунт для полной проверки.',
      code,
      status: 'Ограниченный доступ',
    }
  } finally {
    trackingLoading.value = false
  }
}

function scrollToBlock(id) {
  closeMenus()
  document.querySelector(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function updateHeaderTheme() {
  const pointX = Math.max(40, Math.floor(window.innerWidth / 2))
  const pointY = 96
  const element = document.elementFromPoint(pointX, pointY)
  const section = element?.closest?.('[data-header-theme]')
  headerTheme.value = section?.dataset?.headerTheme || 'dark'
}

onMounted(() => {
  refreshAuthState()
  syncCurrentUser()
  loadWarehouses()
  checkTracking()
  updateHeaderTheme()
  window.addEventListener('scroll', updateHeaderTheme, { passive: true })
  window.addEventListener('resize', updateHeaderTheme)
  window.addEventListener('focus', refreshAuthState)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateHeaderTheme)
  window.removeEventListener('resize', updateHeaderTheme)
  window.removeEventListener('focus', refreshAuthState)
})
</script>

<template>
  <main class="min-h-screen bg-[#f3f6fb] text-[#07101f] selection:bg-[#ff4248] selection:text-white">
    <header
      class="fixed inset-x-0 top-0 z-50 border-b transition-colors duration-300"
      :class="headerTheme === 'dark'
        ? 'border-white/10 bg-[#07101f] text-white shadow-[0_24px_80px_rgba(0,0,0,0.25)]'
        : 'border-slate-200 bg-white text-[#07101f] shadow-[0_18px_50px_rgba(15,23,42,0.08)]'"
    >
      <div class="mx-auto flex h-[92px] max-w-[1500px] items-center justify-between px-5 sm:px-8 lg:px-12">
        <button class="flex items-center gap-4 text-left" @click="router.push('/')">
          <span class="grid h-14 w-14 shrink-0 place-items-center rounded-2xl bg-[#ff4248] text-white shadow-[0_18px_45px_rgba(255,66,72,0.35)]">
            <Truck class="h-7 w-7" />
          </span>
          <span>
            <span class="block text-xl font-black leading-none tracking-tight sm:text-2xl">Fulfillment Transit</span>
            <span class="mt-2 block text-[10px] font-black uppercase tracking-[0.45em] text-[#ff6d72]">Marketplace logistics</span>
          </span>
        </button>

        <nav class="hidden items-center gap-9 text-sm font-black lg:flex">
          <button class="hover:text-[#ff4248]" @click="scrollToBlock('#services')">Услуги</button>
          <button class="hover:text-[#ff4248]" @click="scrollToBlock('#process')">Как работаем</button>
          <button class="hover:text-[#ff4248]" @click="scrollToBlock('#warehouses')">Склады</button>
          <button class="hover:text-[#ff4248]" @click="scrollToBlock('#benefits')">Преимущества</button>
          <button class="hover:text-[#ff4248]" @click="scrollToBlock('#faq')">FAQ</button>
        </nav>

        <div v-if="isAuthenticated" class="hidden items-center gap-3 lg:flex">
          <div class="flex items-center gap-3 rounded-2xl border px-4 py-3" :class="headerTheme === 'dark' ? 'border-white/15 bg-white/5' : 'border-slate-200 bg-slate-50'">
            <span class="grid h-10 w-10 place-items-center rounded-xl bg-[#ff4248] text-sm font-black text-white">{{ userInitials }}</span>
            <span class="max-w-[220px] truncate text-sm font-black">{{ userDisplayName }}</span>
          </div>
          <button class="rounded-2xl bg-[#ff4248] px-7 py-4 text-sm font-black text-white shadow-[0_16px_40px_rgba(255,66,72,0.28)] transition hover:-translate-y-0.5 hover:bg-[#e7353b]" @click="logout">
            Выйти
          </button>
        </div>

        <div v-else class="hidden items-center gap-3 lg:flex">
          <RouterLink to="/login" class="rounded-2xl px-6 py-4 text-sm font-black transition hover:text-[#ff4248]">
            Войти
          </RouterLink>
          <RouterLink to="/register" class="rounded-2xl bg-[#ff4248] px-7 py-4 text-sm font-black text-white shadow-[0_16px_40px_rgba(255,66,72,0.28)] transition hover:-translate-y-0.5 hover:bg-[#e7353b]">
            Регистрация
          </RouterLink>
        </div>

        <button class="hidden rounded-2xl border p-3 md:grid lg:hidden" :class="headerTheme === 'dark' ? 'border-white/20' : 'border-slate-200'" @click="tabletMenuOpen = !tabletMenuOpen">
          <MoreHorizontal class="h-6 w-6" />
        </button>

        <button class="grid rounded-2xl border p-3 md:hidden" :class="headerTheme === 'dark' ? 'border-white/20' : 'border-slate-200'" @click="mobileMenuOpen = !mobileMenuOpen">
          <X v-if="mobileMenuOpen" class="h-6 w-6" />
          <Menu v-else class="h-6 w-6" />
        </button>
      </div>

      <div v-if="tabletMenuOpen || mobileMenuOpen" class="border-t border-white/10 bg-[#07101f] px-5 py-5 text-white shadow-2xl lg:hidden">
        <div class="mx-auto grid max-w-[1500px] gap-3 sm:grid-cols-2 md:grid-cols-4">
          <button class="rounded-2xl bg-white/5 px-5 py-4 text-left font-black" @click="scrollToBlock('#services')">Услуги</button>
          <button class="rounded-2xl bg-white/5 px-5 py-4 text-left font-black" @click="scrollToBlock('#process')">Как работаем</button>
          <button class="rounded-2xl bg-white/5 px-5 py-4 text-left font-black" @click="scrollToBlock('#warehouses')">Склады</button>
          <button class="rounded-2xl bg-white/5 px-5 py-4 text-left font-black" @click="scrollToBlock('#faq')">FAQ</button>
          <template v-if="isAuthenticated">
            <div class="rounded-2xl border border-white/15 px-5 py-4 font-black">
              {{ userDisplayName }}
            </div>
            <button class="rounded-2xl bg-[#ff4248] px-5 py-4 text-center font-black text-white" @click="logout">Выйти</button>
          </template>
          <template v-else>
            <RouterLink to="/login" class="rounded-2xl border border-white/15 px-5 py-4 text-center font-black" @click="closeMenus">Войти</RouterLink>
            <RouterLink to="/register" class="rounded-2xl bg-[#ff4248] px-5 py-4 text-center font-black text-white" @click="closeMenus">Регистрация</RouterLink>
          </template>
        </div>
      </div>
    </header>

    <section data-header-theme="dark" class="relative overflow-hidden bg-[#07101f] pt-[128px] text-white">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_10%_15%,rgba(255,66,72,0.34),transparent_30%),radial-gradient(circle_at_90%_5%,rgba(0,166,214,0.35),transparent_35%),linear-gradient(135deg,#1d0714_0%,#07101f_48%,#06394a_100%)]"></div>
      <div class="absolute inset-x-0 top-[92px] h-px bg-white/10"></div>
      <div class="relative mx-auto grid max-w-[1500px] gap-10 px-5 pb-20 pt-16 sm:px-8 lg:grid-cols-[1.05fr_0.95fr] lg:px-12 lg:pb-24">
        <div class="flex flex-col justify-center">
          <p class="mb-6 text-sm font-black uppercase tracking-[0.45em] text-[#ff8b8f]">fulfillment transit</p>
          <h1 class="max-w-[760px] text-5xl font-black leading-[0.94] tracking-[-0.06em] sm:text-6xl lg:text-7xl xl:text-8xl">
            Приёмка, хранение и отправка товара под контролем
          </h1>
          <p class="mt-8 max-w-[760px] text-xl leading-9 text-white/86">
            Организуем полный путь товара: заявка, календарь приёмки, QR-контроль каждого места, хранение по зонам и отправка на склад назначения.
          </p>
          <div class="mt-10 flex flex-col gap-4 sm:flex-row">
            <RouterLink to="/register" class="inline-flex items-center justify-center gap-3 rounded-3xl bg-[#ff4248] px-9 py-5 text-base font-black text-white shadow-[0_22px_55px_rgba(255,66,72,0.35)] transition hover:-translate-y-0.5">
              Оформить заявку <ArrowRight class="h-5 w-5" />
            </RouterLink>
            <button class="inline-flex items-center justify-center gap-3 rounded-3xl border border-white/70 px-9 py-5 text-base font-black text-white transition hover:bg-white hover:text-[#07101f]" @click="scrollToBlock('#tracking')">
              Проверить статус
            </button>
          </div>

          <div class="mt-14 grid max-w-[720px] gap-4 sm:grid-cols-3">
            <div v-for="item in stats" :key="item.label" class="rounded-3xl border border-white/50 bg-white/5 p-7 text-center backdrop-blur">
              <div class="text-4xl font-black">{{ item.value }}</div>
              <div class="mt-2 text-sm text-white/65">{{ item.label }}</div>
            </div>
          </div>
        </div>

        <div id="tracking" class="rounded-[2.25rem] border border-white/10 bg-white/5 p-5 shadow-[0_30px_100px_rgba(0,0,0,0.35)] backdrop-blur">
          <div class="rounded-[1.8rem] bg-[#07101f] p-6 sm:p-8 lg:p-10">
            <div class="mb-8 flex items-start justify-between gap-4">
              <div>
                <p class="text-xs font-black uppercase tracking-[0.45em] text-[#ff9ca0]">быстрый статус</p>
                <h2 class="mt-4 text-3xl font-black tracking-[-0.04em] sm:text-4xl">Проверка заявки</h2>
              </div>
              <div class="grid h-16 w-16 shrink-0 place-items-center rounded-3xl bg-[#ff4248]">
                <QrCode class="h-8 w-8" />
              </div>
            </div>

            <div class="space-y-6">
              <div>
                <label class="mb-3 block text-xs font-black uppercase tracking-[0.45em] text-white/70">Склад приёмки</label>
                <div class="relative">
                  <button class="flex w-full items-center justify-between rounded-2xl border border-white/10 bg-white/10 px-5 py-5 text-left font-black" @click="toggleDropdown('receiving')">
                    <span>{{ warehouseLabel(selectedReceivingWarehouse) }}</span>
                    <ChevronDown class="h-5 w-5 transition" :class="openedDropdown === 'receiving' ? 'rotate-180' : ''" />
                  </button>
                  <div v-if="openedDropdown === 'receiving'" class="absolute z-20 mt-3 w-full overflow-hidden rounded-2xl border border-white/10 bg-[#0d1628] p-2 shadow-2xl">
                    <button
                      v-for="warehouse in receivingWarehouses"
                      :key="warehouse.id"
                      class="flex w-full items-center justify-between rounded-xl px-4 py-4 text-left font-black text-white hover:bg-[#ff4248]"
                      :class="warehouse.id === selectedReceivingId ? 'bg-[#ff4248]' : ''"
                      @click="selectWarehouse('receiving', warehouse)"
                    >
                      <span>{{ warehouseLabel(warehouse) }}</span>
                      <span v-if="warehouse.warehouse_type === 'both'" class="rounded-full bg-white/15 px-3 py-1 text-[10px] uppercase tracking-widest">универсальный</span>
                    </button>
                    <p v-if="!receivingWarehouses.length" class="px-4 py-4 text-sm text-white/55">Нет складов приёмки</p>
                  </div>
                </div>
              </div>

              <div>
                <label class="mb-3 block text-xs font-black uppercase tracking-[0.45em] text-white/70">Склад назначения</label>
                <div class="relative">
                  <button class="flex w-full items-center justify-between rounded-2xl border border-white/10 bg-white/10 px-5 py-5 text-left font-black" @click="toggleDropdown('destination')">
                    <span>{{ warehouseLabel(selectedDestinationWarehouse) }}</span>
                    <ChevronDown class="h-5 w-5 transition" :class="openedDropdown === 'destination' ? 'rotate-180' : ''" />
                  </button>
                  <div v-if="openedDropdown === 'destination'" class="absolute z-20 mt-3 w-full overflow-hidden rounded-2xl border border-white/10 bg-[#0d1628] p-2 shadow-2xl">
                    <button
                      v-for="warehouse in destinationWarehouses"
                      :key="warehouse.id"
                      class="flex w-full items-center justify-between rounded-xl px-4 py-4 text-left font-black text-white hover:bg-[#ff4248]"
                      :class="warehouse.id === selectedDestinationId ? 'bg-[#ff4248]' : ''"
                      @click="selectWarehouse('destination', warehouse)"
                    >
                      <span>{{ warehouseLabel(warehouse) }}</span>
                      <span v-if="warehouse.warehouse_type === 'both'" class="rounded-full bg-white/15 px-3 py-1 text-[10px] uppercase tracking-widest">универсальный</span>
                    </button>
                    <p v-if="!destinationWarehouses.length" class="px-4 py-4 text-sm text-white/55">Нет складов назначения</p>
                  </div>
                </div>
              </div>

              <div>
                <label class="mb-3 block text-xs font-black uppercase tracking-[0.45em] text-white/70">Номер заявки или QR</label>
                <input v-model="trackingCode" class="w-full rounded-2xl border border-white/10 bg-white/10 px-5 py-5 font-black text-white outline-none placeholder:text-white/40 focus:border-[#ff8b8f]" placeholder="Например: QR-TPRO-MSK-240001" />
              </div>

              <button class="flex w-full items-center justify-center gap-3 rounded-2xl bg-[#ff4248] px-6 py-5 font-black text-white shadow-[0_24px_55px_rgba(255,66,72,0.28)] transition hover:-translate-y-0.5" @click="checkTracking">
                {{ trackingLoading ? 'Проверяем...' : 'Проверить статус' }} <ArrowRight class="h-5 w-5" />
              </button>

              <div v-if="trackingResult" class="rounded-2xl border border-emerald-400/25 bg-emerald-400/10 p-6">
                <div class="flex gap-3">
                  <CheckCircle2 class="mt-1 h-5 w-5 shrink-0 text-emerald-300" />
                  <div>
                    <h3 class="font-black text-white">{{ trackingResult.title }}</h3>
                    <p class="mt-1 font-bold text-emerald-200">{{ trackingResult.subtitle }}</p>
                    <p class="mt-4 leading-7 text-white/75">{{ trackingResult.description }}</p>
                    <p class="mt-4 text-sm text-white/55">
                      {{ warehouseLabel(selectedReceivingWarehouse) }} → {{ warehouseLabel(selectedDestinationWarehouse) }} · {{ trackingResult.code }}
                    </p>
                  </div>
                </div>
              </div>

              <div v-if="warehousesError" class="rounded-2xl border border-amber-300/35 bg-amber-300/10 p-6 text-sm leading-7 text-amber-100">
                {{ warehousesError }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="services" data-header-theme="light" class="py-20 sm:py-24">
      <div class="mx-auto max-w-[1500px] px-5 sm:px-8 lg:px-12">
        <div class="mb-12 flex flex-col justify-between gap-6 lg:flex-row lg:items-end">
          <div>
            <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff4248]">услуги</p>
            <h2 class="mt-5 max-w-[760px] text-4xl font-black leading-tight tracking-[-0.05em] sm:text-6xl">Всё для работы с товаром маркетплейсов</h2>
          </div>
          <p class="max-w-[560px] text-lg leading-8 text-slate-600">Берём на себя операционные этапы: приёмку, QR-контроль, хранение, сортировку и отправку на склады назначения.</p>
        </div>

        <div class="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
          <article v-for="service in services" :key="service.title" class="flex min-h-[330px] flex-col rounded-[2rem] border border-slate-200 bg-white p-8 shadow-[0_18px_45px_rgba(15,23,42,0.06)]">
            <component :is="service.icon" class="h-10 w-10 text-[#ff4248]" />
            <h3 class="mt-12 text-2xl font-black leading-tight tracking-[-0.03em]">{{ service.title }}</h3>
            <p class="mt-6 flex-1 text-lg leading-8 text-slate-600">{{ service.text }}</p>
            <button class="mt-8 inline-flex items-center gap-2 font-black text-[#ff4248]" @click="router.push('/register')">
              Подключить услугу <ArrowRight class="h-4 w-4" />
            </button>
          </article>
        </div>
      </div>
    </section>

    <section id="process" data-header-theme="light" class="bg-[#edf2f8] py-20 sm:py-24">
      <div class="mx-auto grid max-w-[1500px] gap-10 px-5 sm:px-8 lg:grid-cols-[0.75fr_1.25fr] lg:px-12">
        <div>
          <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff4248]">процесс</p>
          <h2 class="mt-5 text-4xl font-black leading-tight tracking-[-0.05em] sm:text-6xl">От заявки до отправки</h2>
          <p class="mt-7 text-lg leading-9 text-slate-600">Каждый этап фиксируется в системе, чтобы клиент и сотрудники видели актуальный статус без таблиц и ручных уточнений.</p>
        </div>
        <div class="space-y-5">
          <article v-for="step in processSteps" :key="step[0]" class="grid gap-5 rounded-[2rem] border border-slate-200 bg-white p-8 shadow-sm sm:grid-cols-[90px_1fr]">
            <div class="text-4xl font-black text-slate-200">{{ step[0] }}</div>
            <div>
              <h3 class="text-2xl font-black">{{ step[1] }}</h3>
              <p class="mt-3 text-lg leading-8 text-slate-600">{{ step[2] }}</p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <section id="warehouses" data-header-theme="light" class="py-20 sm:py-24">
      <div class="mx-auto max-w-[1500px] px-5 sm:px-8 lg:px-12">
        <div class="mb-12 flex flex-col justify-between gap-6 lg:flex-row lg:items-end">
          <div>
            <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff4248]">география</p>
            <h2 class="mt-5 text-4xl font-black leading-tight tracking-[-0.05em] sm:text-6xl">Склады из базы данных</h2>
          </div>
          <p class="max-w-[560px] text-lg leading-8 text-slate-600">Карточки строятся из публичного эндпоинта складов. Универсальные склады используются и для приёмки, и как конечное место.</p>
        </div>

        <div v-if="warehousesLoading" class="rounded-[2rem] border border-slate-200 bg-white p-10 text-center font-black text-slate-500">Загружаем склады...</div>
        <div v-else-if="!warehouseCards.length" class="rounded-[2rem] border border-slate-200 bg-white p-10 text-center font-black text-slate-500">Склады пока не загружены</div>
        <div v-else class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          <article v-for="warehouse in warehouseCards" :key="warehouse.id" class="flex min-h-[340px] flex-col rounded-[2rem] border border-slate-200 bg-white p-8 shadow-[0_18px_45px_rgba(15,23,42,0.06)]">
            <MapPin class="h-11 w-11 text-[#ff4248]" />
            <h3 class="mt-12 text-3xl font-black leading-tight tracking-[-0.04em]">{{ warehouse.city }}</h3>
            <p class="mt-4 text-xl font-black text-slate-800">{{ warehouseLabel(warehouse) }}</p>
            <p class="mt-4 flex-1 text-base leading-7 text-slate-600">{{ warehouse.address }}</p>
            <div class="mt-8 rounded-2xl border border-slate-200 bg-slate-50 px-5 py-4 text-sm font-black uppercase tracking-wider text-[#ff4248]">
              {{ warehouseKindLabel(warehouse.warehouse_type) }}
            </div>
          </article>
        </div>
      </div>
    </section>

    <section id="benefits" data-header-theme="dark" class="bg-[#07101f] py-20 text-white sm:py-24">
      <div class="mx-auto grid max-w-[1500px] gap-10 px-5 sm:px-8 lg:grid-cols-[0.85fr_1.15fr] lg:px-12">
        <div>
          <p class="text-sm font-black uppercase tracking-[0.45em] text-[#ff8b8f]">преимущества</p>
          <h2 class="mt-5 text-4xl font-black leading-tight tracking-[-0.05em] sm:text-6xl">Операционный контроль без хаоса</h2>
        </div>
        <div class="grid gap-5">
          <div v-for="benefit in benefits" :key="benefit" class="flex gap-4 rounded-[2rem] border border-white/10 bg-white/5 p-7">
            <ShieldCheck class="mt-1 h-6 w-6 shrink-0 text-[#ff8b8f]" />
            <p class="text-lg leading-8 text-white/82">{{ benefit }}</p>
          </div>
        </div>
      </div>
    </section>

    <section id="faq" data-header-theme="light" class="py-20 sm:py-24">
      <div class="mx-auto max-w-[980px] px-5 sm:px-8 lg:px-12">
        <p class="text-center text-sm font-black uppercase tracking-[0.45em] text-[#ff4248]">faq</p>
        <h2 class="mt-5 text-center text-4xl font-black tracking-[-0.05em] sm:text-6xl">Частые вопросы</h2>
        <div class="mt-12 space-y-5">
          <article
            v-for="(item, index) in faqs"
            :key="item.q"
            class="rounded-[2rem] border border-slate-200 bg-white p-0 shadow-sm transition hover:-translate-y-0.5 hover:shadow-[0_18px_45px_rgba(15,23,42,0.07)]"
          >
            <button class="flex w-full items-center justify-between gap-4 px-7 py-7 text-left text-xl font-black" @click="toggleFaq(index)">
              <span>{{ item.q }}</span>
              <ChevronDown class="h-5 w-5 shrink-0 text-[#ff4248] transition" :class="openedFaq === index ? 'rotate-180' : ''" />
            </button>
            <div v-show="openedFaq === index" class="px-7 pb-7">
              <p class="text-lg leading-8 text-slate-600">{{ item.a }}</p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <footer data-header-theme="dark" class="bg-[#07101f] px-5 py-12 text-white sm:px-8 lg:px-12">
      <div class="mx-auto max-w-[1500px]">
        <div class="rounded-[2rem] bg-[#ff4248] p-8 sm:p-12 lg:flex lg:items-center lg:justify-between">
          <div>
            <h2 class="max-w-[760px] text-4xl font-black leading-tight tracking-[-0.05em] sm:text-5xl">Готовы запустить управляемый fulfillment?</h2>
            <p class="mt-5 max-w-[760px] text-lg leading-8 text-white/90">Создайте заявку, выберите склад приёмки и отслеживайте путь товара по статусам.</p>
          </div>
          <RouterLink to="/register" class="mt-8 inline-flex items-center justify-center gap-3 rounded-3xl bg-white px-8 py-5 font-black text-[#c81620] lg:mt-0">
            Начать работу <ArrowRight class="h-5 w-5" />
          </RouterLink>
        </div>

        <div class="mt-10 flex flex-col justify-between gap-6 border-t border-white/10 pt-8 text-sm text-white/60 lg:flex-row">
          <p>© 2026 Fulfillment Transit. Фулфилмент для продавцов маркетплейсов.</p>
          <p>Москва · Казань · Санкт-Петербург · Екатеринбург</p>
        </div>
      </div>
    </footer>
  </main>
</template>
