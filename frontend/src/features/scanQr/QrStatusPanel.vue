<script setup>
import { computed, nextTick, onBeforeUnmount, ref } from 'vue'
import {
  AlertTriangle,
  Camera,
  CheckCircle2,
  ClipboardList,
  Keyboard,
  Loader2,
  LogIn,
  PackageCheck,
  QrCode,
  Search,
  X,
} from '@lucide/vue'
import { apiFetch, getAccessToken } from '@/shared/api/http'

const manualCode = ref('QR-TPRO-MSK-240001')
const modalMode = ref(null) // null | 'camera' | 'manual'
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const cargoItem = ref(null)
const scanner = ref(null)
const scannerRunning = ref(false)
const cameraMessage = ref('Камера не запущена')
const scannerElementId = `qr-reader-${Math.random().toString(36).slice(2)}`

const isAuthenticated = computed(() => Boolean(getAccessToken()))

const modeTitle = computed(() => {
  if (modalMode.value === 'camera') return 'Сканирование QR-кода'
  if (modalMode.value === 'manual') return 'Ручная проверка'
  return ''
})

const statusLabel = computed(() => mapStatus(cargoItem.value?.status))

const infoRows = computed(() => {
  const item = cargoItem.value || {}
  const order = item.order || item.order_info || {}
  const zone = item.storage_zone || item.zone || {}
  const gate = item.gate || {}

  return [
    ['QR-код', item.qr_code || item.qrCode],
    ['Грузовое место', item.id ? `№ ${item.id}` : null],
    ['Заявка', item.order_id || item.orderId || order.id],
    ['Статус', statusLabel.value],
    ['Зона хранения', zone.name || item.storage_zone_name || item.storageZoneName || item.storage_zone_id],
    ['Гейт', gate.name || item.gate_name || item.gateName || item.gate_id],
    ['Принято', formatDate(item.received_at || item.receivedAt)],
    ['Отгружено', formatDate(item.shipped_at || item.shippedAt)],
    ['Создано', formatDate(item.created_at || item.createdAt)],
    ['Обновлено', formatDate(item.updated_at || item.updatedAt)],
  ].filter(([, value]) => value !== undefined && value !== null && value !== '')
})

function mapStatus(status) {
  const map = {
    created: 'Создано',
    pending: 'Ожидает обработки',
    accepted: 'Принято на складе',
    received: 'Принято на складе',
    stored: 'На хранении',
    assigned_to_gate: 'Назначено к гейту',
    ready_to_ship: 'Готово к отправке',
    shipped: 'Отгружено',
    delivered: 'Доставлено',
    cancelled: 'Отменено',
  }
  return map[status] || status || '—'
}

function formatDate(value) {
  if (!value) return null
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

function normalizeCargoResponse(data) {
  return data?.cargo_item || data?.cargoItem || data?.item || data?.cargo || data
}

function resetState({ keepResult = false } = {}) {
  errorMessage.value = ''
  successMessage.value = ''
  isLoading.value = false
  if (!keepResult) cargoItem.value = null
}

async function openManualModal() {
  await stopScanner()
  resetState()
  modalMode.value = 'manual'
}

async function openCameraModal() {
  resetState()
  modalMode.value = 'camera'
  await nextTick()
  await startScanner()
}

async function closeModal() {
  await stopScanner()
  modalMode.value = null
  resetState({ keepResult: true })
}

async function startScanner() {
  if (scannerRunning.value) return

  if (!navigator.mediaDevices?.getUserMedia) {
    errorMessage.value = 'Браузер не поддерживает доступ к камере. Введите QR-код вручную.'
    return
  }

  try {
    cameraMessage.value = 'Запрашиваем доступ к камере...'

    // Явно вызываем getUserMedia, чтобы браузер показал окно разрешения.
    const stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' },
      audio: false,
    })
    stream.getTracks().forEach((track) => track.stop())

    const { Html5Qrcode } = await import('html5-qrcode')
    const instance = new Html5Qrcode(scannerElementId)
    scanner.value = instance

    cameraMessage.value = 'Наведите камеру на QR-код'
    await instance.start(
      { facingMode: 'environment' },
      {
        fps: 10,
        qrbox: { width: 260, height: 260 },
        aspectRatio: 1,
      },
      async (decodedText) => {
        if (!decodedText) return
        await stopScanner()
        manualCode.value = decodedText.trim()
        await fetchCargoByQr(decodedText.trim())
      },
      () => {},
    )

    scannerRunning.value = true
  } catch (error) {
    const name = error?.name || ''
    if (name === 'NotAllowedError' || name === 'PermissionDeniedError') {
      errorMessage.value = 'Доступ к камере запрещён. Разрешите камеру в браузере или введите код вручную.'
    } else if (name === 'NotFoundError' || name === 'DevicesNotFoundError') {
      errorMessage.value = 'Камера не найдена. Введите QR-код вручную.'
    } else {
      errorMessage.value = error?.message || 'Не удалось открыть камеру. Введите QR-код вручную.'
    }
    cameraMessage.value = 'Камера не запущена'
  }
}

async function stopScanner() {
  if (!scanner.value) return

  try {
    if (scannerRunning.value) {
      await scanner.value.stop()
    }
  } catch {
    // scanner мог быть уже остановлен библиотекой
  }

  try {
    await scanner.value.clear()
  } catch {
    // clear может упасть, если контейнер уже очищен
  }

  scanner.value = null
  scannerRunning.value = false
  cameraMessage.value = 'Камера не запущена'
}

async function checkManualCode() {
  const code = manualCode.value.trim()
  if (!code) {
    errorMessage.value = 'Введите QR-код или номер грузового места.'
    return
  }
  await fetchCargoByQr(code)
}

async function fetchCargoByQr(code) {
  resetState({ keepResult: false })

  if (!isAuthenticated.value) {
    errorMessage.value = 'Сначала войдите в аккаунт. Запрос к QR защищён авторизацией.'
    return
  }

  isLoading.value = true

  try {
    const data = await apiFetch(`/cargo-items/scan?qr_code=${encodeURIComponent(code)}`, {
      auth: true,
    })

    cargoItem.value = normalizeCargoResponse(data)
    successMessage.value = 'Грузовое место найдено'
  } catch (error) {
    if (error.status === 401) {
      errorMessage.value = 'Сессия истекла или токен не найден. Войдите в аккаунт заново.'
    } else if (error.status === 403) {
      errorMessage.value = 'Недостаточно прав для просмотра этого QR-кода.'
    } else if (error.status === 404) {
      errorMessage.value = 'QR-код не найден.'
    } else {
      errorMessage.value = error?.message || 'Не удалось получить информацию по QR-коду.'
    }
  } finally {
    isLoading.value = false
  }
}

onBeforeUnmount(async () => {
  await stopScanner()
})
</script>

<template>
  <section class="rounded-[2rem] border border-white/10 bg-[#07101f] p-6 shadow-2xl sm:p-8 lg:p-10">
    <div class="flex flex-col gap-8 lg:flex-row lg:items-start lg:justify-between">
      <div class="max-w-xl">
        <p class="mb-4 text-xs font-black uppercase tracking-[0.45em] text-red-300">Быстрый статус</p>
        <h2 class="text-4xl font-black leading-tight text-white sm:text-5xl">Проверка по QR-коду</h2>
        <p class="mt-5 text-lg leading-8 text-slate-300">
          Сканируйте QR камерой или введите код вручную. После проверки откроется карточка грузового места.
        </p>
      </div>

      <div class="flex h-20 w-20 items-center justify-center rounded-[1.5rem] bg-red-500 text-white shadow-xl shadow-red-500/30">
        <QrCode class="h-10 w-10" />
      </div>
    </div>

    <div class="mt-10 grid gap-5 md:grid-cols-2">
      <button
        type="button"
        class="group rounded-[2rem] border border-white/20 bg-white/5 p-6 text-left transition hover:-translate-y-1 hover:border-red-300 hover:bg-white/10"
        @click="openCameraModal"
      >
        <span class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-red-500 text-white shadow-lg shadow-red-500/30">
          <Camera class="h-8 w-8" />
        </span>
        <span class="block text-3xl font-black text-white">Сканировать камерой</span>
        <span class="mt-4 block text-base leading-7 text-slate-300">Откроется окно сканирования и браузер запросит доступ к камере.</span>
      </button>

      <button
        type="button"
        class="group rounded-[2rem] border border-white/20 bg-white/5 p-6 text-left transition hover:-translate-y-1 hover:border-red-300 hover:bg-white/10"
        @click="openManualModal"
      >
        <span class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-red-500 text-white shadow-lg shadow-red-500/30">
          <Keyboard class="h-8 w-8" />
        </span>
        <span class="block text-3xl font-black text-white">Ввести код вручную</span>
        <span class="mt-4 block text-base leading-7 text-slate-300">Подходит для проверки кода вида QR-TPRO-MSK-240001.</span>
      </button>
    </div>
  </section>

  <Teleport to="body">
    <div
      v-if="modalMode"
      class="fixed inset-0 z-[100] flex items-center justify-center bg-black/75 p-4 backdrop-blur-sm"
      @click.self="closeModal"
    >
      <div class="relative max-h-[92vh] w-full max-w-6xl overflow-y-auto rounded-[2rem] border border-white/20 bg-[#07101f] p-5 text-white shadow-2xl sm:p-8">
        <button
          type="button"
          class="absolute right-5 top-5 flex h-12 w-12 items-center justify-center rounded-2xl bg-white/10 text-white transition hover:bg-white/20"
          @click="closeModal"
        >
          <X class="h-6 w-6" />
        </button>

        <div class="pr-14">
          <p class="text-xs font-black uppercase tracking-[0.45em] text-red-300">
            {{ modalMode === 'camera' ? 'Сканирование' : 'Ручная проверка' }}
          </p>
          <h3 class="mt-3 text-4xl font-black sm:text-5xl">{{ modeTitle }}</h3>
        </div>

        <div class="mt-8 grid gap-6 lg:grid-cols-[1fr_0.9fr]">
          <div class="rounded-[1.75rem] border border-white/20 bg-white/5 p-5">
            <template v-if="modalMode === 'camera'">
              <div class="mb-4 flex items-center justify-between gap-4">
                <div>
                  <h4 class="text-2xl font-black">Камера</h4>
                  <p class="mt-1 text-sm text-slate-300">Наведите камеру на QR-код грузового места.</p>
                </div>
                <button
                  type="button"
                  class="rounded-2xl bg-red-500 px-5 py-3 text-sm font-black text-white transition hover:bg-red-400 disabled:opacity-60"
                  :disabled="scannerRunning"
                  @click="startScanner"
                >
                  Запустить
                </button>
              </div>

              <div class="relative overflow-hidden rounded-[1.5rem] bg-[#030914] p-3">
                <div :id="scannerElementId" class="min-h-[320px] overflow-hidden rounded-[1.25rem]"></div>
                <div
                  v-if="!scannerRunning"
                  class="pointer-events-none absolute inset-3 flex items-center justify-center rounded-[1.25rem] border border-dashed border-white/15 bg-[#07101f] text-center text-slate-300"
                >
                  {{ cameraMessage }}
                </div>
              </div>
            </template>

            <template v-else>
              <h4 class="text-2xl font-black">Введите код</h4>
              <p class="mt-2 text-sm text-slate-300">Например: QR-TPRO-MSK-240001</p>
              <div class="mt-6 rounded-[1.5rem] bg-[#202838] p-4">
                <div class="flex items-center gap-3 rounded-[1.25rem] border border-white/10 bg-white/5 px-4 py-3">
                  <Keyboard class="h-5 w-5 text-slate-400" />
                  <input
                    v-model="manualCode"
                    class="w-full bg-transparent text-lg font-black text-white outline-none placeholder:text-slate-500"
                    placeholder="QR-TPRO-MSK-240001"
                    @keydown.enter="checkManualCode"
                  />
                </div>
              </div>
              <button
                type="button"
                class="mt-5 flex w-full items-center justify-center gap-3 rounded-[1.5rem] bg-red-500 px-6 py-5 text-lg font-black text-white shadow-xl shadow-red-500/25 transition hover:bg-red-400 disabled:cursor-not-allowed disabled:opacity-70"
                :disabled="isLoading"
                @click="checkManualCode"
              >
                <Loader2 v-if="isLoading" class="h-5 w-5 animate-spin" />
                <Search v-else class="h-5 w-5" />
                Получить информацию
              </button>
            </template>
          </div>

          <div class="rounded-[1.75rem] border border-white/20 bg-white/5 p-5">
            <h4 class="flex items-center gap-3 text-2xl font-black">
              <PackageCheck class="h-6 w-6 text-red-300" />
              Информация
            </h4>

            <div v-if="isLoading" class="mt-8 flex items-center gap-3 rounded-2xl bg-white/5 p-5 text-slate-200">
              <Loader2 class="h-5 w-5 animate-spin" />
              Получаем данные...
            </div>

            <div v-else-if="errorMessage" class="mt-8 rounded-2xl border border-red-400/30 bg-red-500/10 p-5 text-red-100">
              <div class="flex gap-3">
                <AlertTriangle class="mt-1 h-5 w-5 shrink-0" />
                <p class="leading-7">{{ errorMessage }}</p>
              </div>
            </div>

            <div v-else-if="cargoItem" class="mt-8 space-y-5">
              <div class="rounded-2xl border border-emerald-400/30 bg-emerald-500/10 p-5 text-emerald-100">
                <div class="flex gap-3">
                  <CheckCircle2 class="mt-1 h-5 w-5 shrink-0" />
                  <div>
                    <p class="font-black">{{ successMessage || 'Данные получены' }}</p>
                    <p class="mt-1 text-sm text-emerald-100/80">Карточка грузового места найдена.</p>
                  </div>
                </div>
              </div>

              <dl class="grid gap-3">
                <div
                  v-for="([label, value]) in infoRows"
                  :key="label"
                  class="flex items-start justify-between gap-4 rounded-2xl bg-[#111a2a] px-4 py-3"
                >
                  <dt class="text-sm text-slate-400">{{ label }}</dt>
                  <dd class="text-right font-black text-white">{{ value }}</dd>
                </div>
              </dl>
            </div>

            <div v-else class="mt-8 rounded-2xl bg-[#111a2a] p-5 text-slate-300">
              <ClipboardList class="mb-4 h-8 w-8 text-slate-500" />
              Здесь появится информация после сканирования QR или ручной проверки.
              <div v-if="!isAuthenticated" class="mt-4 flex items-center gap-2 text-sm text-red-200">
                <LogIn class="h-4 w-4" />
                Для запроса нужен вход в аккаунт.
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
