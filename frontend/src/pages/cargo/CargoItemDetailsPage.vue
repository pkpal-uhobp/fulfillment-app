<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import {
  ArrowLeft,
  Box,
  CalendarClock,
  CheckCircle2,
  ClipboardList,
  Loader2,
  LogIn,
  MapPin,
  QrCode,
  ShieldCheck,
  Warehouse,
} from '@lucide/vue'
import { apiFetch, getAccessToken } from '@/shared/api/http'

const route = useRoute()

const loading = ref(false)
const error = ref('')
const cargoItem = ref(null)
const history = ref([])

const qrCode = computed(() => decodeURIComponent(String(route.params.qrCode || '')))
const isAuthorized = computed(() => Boolean(getAccessToken()))

function statusLabel(status) {
  const labels = {
    created: 'Создано',
    accepted: 'Принято на склад',
    received: 'Принято на склад',
    stored: 'На хранении',
    in_storage: 'На хранении',
    assigned_to_gate: 'Назначено к гейту',
    ready_to_ship: 'Готово к отправке',
    shipped: 'Отправлено',
    delivered: 'Доставлено',
    lost: 'Утеряно',
    damaged: 'Повреждено',
    cancelled: 'Отменено',
  }

  return labels[status] || status || 'Статус не указан'
}

function formatDate(value) {
  if (!value) return '—'
  try {
    return new Intl.DateTimeFormat('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(new Date(value))
  } catch {
    return value
  }
}

function field(...keys) {
  const item = cargoItem.value || {}
  for (const key of keys) {
    if (item[key] !== undefined && item[key] !== null && item[key] !== '') return item[key]
  }
  return '—'
}

async function loadCargoItem() {
  if (!isAuthorized.value) {
    error.value = 'Для просмотра карточки грузового места нужно войти в аккаунт.'
    return
  }

  loading.value = true
  error.value = ''
  cargoItem.value = null
  history.value = []

  try {
    const payload = await apiFetch(`/cargo-items/scan?qr_code=${encodeURIComponent(qrCode.value)}`, { auth: true })
    const item = payload?.cargo_item || payload?.cargoItem || payload?.item || payload
    cargoItem.value = item

    const id = item?.id || item?.cargo_item_id
    if (id) {
      try {
        const historyPayload = await apiFetch(`/cargo-items/${id}/history`, { auth: true })
        history.value = historyPayload?.history || historyPayload?.items || historyPayload?.cargo_status_history || historyPayload || []
        if (!Array.isArray(history.value)) history.value = []
      } catch {
        history.value = []
      }
    }
  } catch (err) {
    if (err.status === 401) {
      error.value = 'Сессия истекла. Войдите в аккаунт ещё раз.'
    } else if (err.status === 403) {
      error.value = 'У вашего аккаунта нет доступа к этому грузовому месту.'
    } else if (err.status === 404) {
      error.value = 'Грузовое место с таким QR-кодом не найдено.'
    } else {
      error.value = err.message || 'Не удалось загрузить информацию по QR-коду.'
    }
  } finally {
    loading.value = false
  }
}

onMounted(loadCargoItem)
</script>

<template>
  <main class="min-h-screen bg-[#07101f] text-white">
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_10%_0%,rgba(255,63,73,0.26),transparent_30%),radial-gradient(circle_at_95%_10%,rgba(20,184,166,0.24),transparent_28%),linear-gradient(135deg,#120814_0%,#07101f_55%,#06394b_100%)]"></div>

    <section class="relative mx-auto w-[min(1320px,calc(100%-32px))] py-10 lg:py-16">
      <div class="mb-8 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <RouterLink to="/" class="inline-flex items-center gap-3 rounded-2xl border border-white/12 bg-white/5 px-5 py-3 text-sm font-black text-white/80 transition hover:bg-white/10 hover:text-white">
          <ArrowLeft class="h-4 w-4" />
          На главную
        </RouterLink>

        <RouterLink v-if="!isAuthorized" to="/login" class="inline-flex items-center gap-3 rounded-2xl bg-[#ff3f49] px-5 py-3 text-sm font-black text-white">
          <LogIn class="h-4 w-4" />
          Войти
        </RouterLink>
      </div>

      <div class="rounded-[42px] border border-white/10 bg-white/[0.055] p-5 shadow-2xl backdrop-blur sm:p-8 lg:p-10">
        <div class="rounded-[34px] bg-[#07101f] p-7 sm:p-10">
          <div class="flex flex-col justify-between gap-8 lg:flex-row lg:items-start">
            <div>
              <p class="text-xs font-black uppercase tracking-[0.45em] text-[#ff9aa0]">карточка грузового места</p>
              <h1 class="mt-5 text-4xl font-black tracking-[-0.06em] sm:text-6xl">Информация по QR</h1>
              <p class="mt-5 max-w-2xl text-lg leading-8 text-white/70">
                Данные загружаются из backend по QR-коду. Клиент видит только свои грузовые места.
              </p>
            </div>

            <div class="rounded-[28px] border border-white/10 bg-white/[0.055] p-5">
              <div class="flex items-center gap-3 text-[#ff5962]">
                <QrCode class="h-7 w-7" />
                <span class="text-xl font-black">{{ qrCode }}</span>
              </div>
            </div>
          </div>

          <div v-if="loading" class="mt-10 flex items-center gap-3 rounded-[28px] border border-white/10 bg-white/[0.05] p-6 text-white/75">
            <Loader2 class="h-6 w-6 animate-spin text-[#ff5962]" />
            Загружаем карточку грузового места...
          </div>

          <div v-else-if="error" class="mt-10 rounded-[28px] border border-amber-400/30 bg-amber-400/10 p-6 text-amber-100">
            <p class="text-lg font-black">{{ error }}</p>
            <RouterLink v-if="!isAuthorized" to="/login" class="mt-5 inline-flex rounded-2xl bg-white px-5 py-3 text-sm font-black text-slate-950">
              Войти в аккаунт
            </RouterLink>
          </div>

          <template v-else-if="cargoItem">
            <div class="mt-10 grid gap-5 lg:grid-cols-[1fr_0.72fr]">
              <section class="rounded-[32px] border border-emerald-400/25 bg-emerald-400/10 p-6">
                <div class="flex flex-col justify-between gap-5 sm:flex-row sm:items-center">
                  <div class="flex items-center gap-4">
                    <span class="flex h-14 w-14 items-center justify-center rounded-2xl bg-emerald-400/15 text-emerald-200">
                      <CheckCircle2 class="h-7 w-7" />
                    </span>
                    <div>
                      <p class="text-sm font-black uppercase tracking-[0.25em] text-emerald-200/70">текущий статус</p>
                      <h2 class="mt-2 text-3xl font-black text-emerald-50">{{ statusLabel(cargoItem.status) }}</h2>
                    </div>
                  </div>
                </div>

                <dl class="mt-8 grid gap-4 sm:grid-cols-2">
                  <div class="rounded-2xl bg-white/[0.06] p-5">
                    <dt class="text-xs uppercase tracking-[0.22em] text-white/45">QR-код</dt>
                    <dd class="mt-2 font-black">{{ field('qr_code', 'qrCode') }}</dd>
                  </div>
                  <div class="rounded-2xl bg-white/[0.06] p-5">
                    <dt class="text-xs uppercase tracking-[0.22em] text-white/45">Грузовое место</dt>
                    <dd class="mt-2 font-black">№ {{ field('id', 'cargo_item_id') }}</dd>
                  </div>
                  <div class="rounded-2xl bg-white/[0.06] p-5">
                    <dt class="text-xs uppercase tracking-[0.22em] text-white/45">Заявка</dt>
                    <dd class="mt-2 font-black">№ {{ field('order_id', 'orderId') }}</dd>
                  </div>
                  <div class="rounded-2xl bg-white/[0.06] p-5">
                    <dt class="text-xs uppercase tracking-[0.22em] text-white/45">Статус</dt>
                    <dd class="mt-2 font-black">{{ statusLabel(cargoItem.status) }}</dd>
                  </div>
                </dl>
              </section>

              <aside class="grid gap-4">
                <div class="rounded-[28px] border border-white/10 bg-white/[0.055] p-5">
                  <div class="flex items-center gap-3 text-white">
                    <Warehouse class="h-5 w-5 text-[#ff5962]" />
                    <p class="font-black">Складские данные</p>
                  </div>
                  <dl class="mt-5 grid gap-4 text-sm">
                    <div>
                      <dt class="text-white/45">Зона хранения</dt>
                      <dd class="mt-1 font-black">{{ field('storage_zone_name', 'storage_zone_id', 'storageZoneID') }}</dd>
                    </div>
                    <div>
                      <dt class="text-white/45">Гейт</dt>
                      <dd class="mt-1 font-black">{{ field('gate_name', 'gate_id', 'gateID') }}</dd>
                    </div>
                  </dl>
                </div>

                <div class="rounded-[28px] border border-white/10 bg-white/[0.055] p-5">
                  <div class="flex items-center gap-3 text-white">
                    <CalendarClock class="h-5 w-5 text-[#ff5962]" />
                    <p class="font-black">Даты</p>
                  </div>
                  <dl class="mt-5 grid gap-4 text-sm">
                    <div>
                      <dt class="text-white/45">Создано</dt>
                      <dd class="mt-1 font-black">{{ formatDate(field('created_at', 'createdAt')) }}</dd>
                    </div>
                    <div>
                      <dt class="text-white/45">Обновлено</dt>
                      <dd class="mt-1 font-black">{{ formatDate(field('updated_at', 'updatedAt')) }}</dd>
                    </div>
                    <div>
                      <dt class="text-white/45">Принято</dt>
                      <dd class="mt-1 font-black">{{ formatDate(field('received_at', 'receivedAt')) }}</dd>
                    </div>
                    <div>
                      <dt class="text-white/45">Отправлено</dt>
                      <dd class="mt-1 font-black">{{ formatDate(field('shipped_at', 'shippedAt')) }}</dd>
                    </div>
                  </dl>
                </div>
              </aside>
            </div>

            <section class="mt-5 rounded-[32px] border border-white/10 bg-white/[0.045] p-6">
              <div class="flex items-center gap-3">
                <ClipboardList class="h-6 w-6 text-[#ff5962]" />
                <h2 class="text-2xl font-black">История статусов</h2>
              </div>

              <div v-if="history.length" class="mt-6 grid gap-3">
                <article v-for="(item, index) in history" :key="item.id || index" class="rounded-2xl bg-white/[0.055] p-5">
                  <div class="flex flex-col justify-between gap-2 sm:flex-row sm:items-center">
                    <p class="font-black">{{ statusLabel(item.new_status || item.status || item.to_status) }}</p>
                    <p class="text-sm text-white/50">{{ formatDate(item.created_at || item.changed_at) }}</p>
                  </div>
                  <p v-if="item.comment" class="mt-2 text-sm leading-6 text-white/60">{{ item.comment }}</p>
                </article>
              </div>

              <p v-else class="mt-5 rounded-2xl bg-white/[0.05] p-5 text-white/55">
                История пока не загружена или отсутствует.
              </p>
            </section>
          </template>
        </div>
      </div>
    </section>
  </main>
</template>
