<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { apiFetch } from '@/shared/api/http'
import {
  cargoStatuses,
  cargoTitle,
  extractCargo,
  formatDateTime,
  gateTitle,
  normalizeCollection,
  routeTitle,
  safeFileName,
  statusLabel,
  statusTone,
  zoneTitle,
} from './workerUtils'

const route = useRoute()
const qrCode = ref(String(route.query.qr || 'QR-TPRO-MSK-240001'))
const cargo = ref(null)
const history = ref([])
const selectedStatus = ref('accepted')
const comment = ref('')
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const notice = ref('')
const scannerMode = ref('manual')
const cameraState = ref('idle')
const cameraError = ref('')
const readerId = 'worker-qr-reader'
let scanner = null
let lastScan = ''

const hasCargo = computed(() => Boolean(cargo.value?.id))
const currentStatus = computed(() => cargo.value?.status || 'accepted')
const terminalStatuses = ['shipped', 'lost', 'damaged', 'cancelled']

const nextStatusMap = {
  accepted: ['accepted', 'stored', 'lost', 'damaged', 'cancelled'],
  stored: ['stored', 'ready_to_ship', 'lost', 'damaged', 'cancelled'],
  ready_to_ship: ['ready_to_ship', 'shipped', 'lost', 'damaged', 'cancelled'],
  shipped: ['shipped'],
  lost: ['lost'],
  damaged: ['damaged'],
  cancelled: ['cancelled'],
}

const workerStatusOptions = computed(() => {
  const allowed = nextStatusMap[currentStatus.value] || ['accepted', 'stored', 'ready_to_ship', 'shipped', 'lost', 'damaged', 'cancelled']
  return cargoStatuses.filter((item) => allowed.includes(item.value))
})

const statusHint = computed(() => {
  if (!hasCargo.value) return ''
  if (terminalStatuses.includes(currentStatus.value)) return 'Груз находится в финальном статусе. Его нельзя перевести дальше.'
  if (currentStatus.value === 'accepted' && !cargo.value?.storage_zone_id) return 'Чтобы перевести груз “На хранении”, логист должен назначить зону хранения.'
  if (currentStatus.value === 'stored' && !cargo.value?.gate_id) return 'Чтобы подготовить груз к отгрузке, логист должен назначить гейт.'
  if (currentStatus.value === 'ready_to_ship' && !cargo.value?.gate_id) return 'Для отгрузки нужен назначенный гейт.'
  return 'Выберите следующий допустимый статус и сохраните операцию.'
})

function setMessage(kind, text) {
  if (kind === 'error') {
    error.value = text
    notice.value = ''
  } else {
    notice.value = text
    error.value = ''
  }
}

function applyCargo(payload, fallbackQr = '') {
  const item = extractCargo(payload)
  cargo.value = item
  qrCode.value = item?.qr_code || fallbackQr || qrCode.value
  selectedStatus.value = item?.status || 'accepted'
}

async function checkQr(code = qrCode.value) {
  const normalized = String(code || '').trim()
  if (!normalized) {
    setMessage('error', 'Введите или отсканируйте QR-код грузового места')
    return
  }
  loading.value = true
  error.value = ''
  notice.value = ''
  try {
    const payload = await apiFetch(`/cargo-items/scan?qr_code=${encodeURIComponent(normalized)}`, { auth: true })
    applyCargo(payload, normalized)
    setMessage('success', 'Грузовое место найдено')
    await loadHistory()
  } catch (err) {
    cargo.value = null
    history.value = []
    setMessage('error', err.message || 'Не удалось получить информацию по QR')
  } finally {
    loading.value = false
  }
}

async function loadHistory() {
  if (!cargo.value?.id) return
  try {
    const payload = await apiFetch(`/cargo-items/${cargo.value.id}/history`, { auth: true })
    history.value = normalizeCollection(payload, ['history', 'items', 'data'])
  } catch {
    history.value = []
  }
}

async function updateStatus() {
  if (!cargo.value?.id) return
  if (!selectedStatus.value) {
    setMessage('error', 'Выберите статус грузового места')
    return
  }
  saving.value = true
  error.value = ''
  notice.value = ''
  try {
    const payload = await apiFetch(`/cargo-items/${cargo.value.id}/status`, {
      method: 'PATCH',
      auth: true,
      body: { status: selectedStatus.value, comment: comment.value.trim() || undefined },
    })
    applyCargo(payload)
    comment.value = ''
    setMessage('success', 'Статус грузового места обновлён')
    await loadHistory()
  } catch (err) {
    const message = err.message || 'Не удалось обновить статус'
    if (message.includes('storage zone')) {
      setMessage('error', 'Нужно назначить зону хранения у логиста, потом можно поставить статус “На хранении”.')
    } else if (message.includes('gate')) {
      setMessage('error', 'Нужно назначить гейт у логиста, потом можно готовить груз к отгрузке.')
    } else if (message.includes('transition')) {
      setMessage('error', 'Недопустимый переход статуса для текущего этапа обработки.')
    } else {
      setMessage('error', message)
    }
  } finally {
    saving.value = false
  }
}

async function startCamera() {
  scannerMode.value = 'camera'
  cameraError.value = ''
  error.value = ''
  notice.value = ''
  cameraState.value = 'starting'
  await nextTick()

  try {
    if (!window.isSecureContext && !['localhost', '127.0.0.1'].includes(window.location.hostname)) {
      throw new Error('Камера работает только на HTTPS, localhost или 127.0.0.1')
    }
    const { Html5Qrcode } = await import('html5-qrcode')
    if (scanner) {
      await stopCamera()
      await nextTick()
    }
    scanner = new Html5Qrcode(readerId, { verbose: false })
    lastScan = ''
    await scanner.start(
      { facingMode: 'environment' },
      { fps: 10, qrbox: { width: 280, height: 280 }, aspectRatio: 1 },
      async (decodedText) => {
        const value = String(decodedText || '').trim()
        if (!value || value === lastScan) return
        lastScan = value
        qrCode.value = value
        await stopCamera()
        await checkQr(value)
      },
      () => {},
    )
    cameraState.value = 'running'
  } catch (err) {
    cameraState.value = 'error'
    cameraError.value = err?.message || 'Не удалось включить камеру. Проверьте разрешение браузера.'
    try { await stopCamera() } catch {}
  }
}

async function stopCamera() {
  if (!scanner) {
    cameraState.value = 'idle'
    return
  }
  try {
    if (scanner.isScanning) await scanner.stop()
    await scanner.clear()
  } catch {
    // scanner may already be stopped
  } finally {
    scanner = null
    cameraState.value = 'idle'
  }
}

async function switchToManual() {
  await stopCamera()
  scannerMode.value = 'manual'
}

function selectStatus(value) {
  selectedStatus.value = value
}

async function downloadCargoQr() {
  if (!cargo.value?.qr_code) return
  const qrModule = await import('qrcode')
  const QRCode = qrModule.default || qrModule
  const dataUrl = await QRCode.toDataURL(cargo.value.qr_code, {
    width: 720,
    margin: 3,
    errorCorrectionLevel: 'M',
  })
  const link = document.createElement('a')
  link.href = dataUrl
  link.download = `${safeFileName(cargo.value.qr_code)}.png`
  link.click()
}

onMounted(() => {
  if (route.query.qr) checkQr(route.query.qr)
})

onBeforeUnmount(stopCamera)
</script>

<template>
  <section class="scan-page">
    <header class="scan-hero">
      <div>
        <p class="eyebrow">QR-контроль</p>
        <h1>Работа с грузом</h1>
        <span>Сканируйте QR камерой или введите код вручную. После проверки можно скачать QR и изменить статус.</span>
      </div>
      <div class="mode-tabs">
        <button type="button" :class="{ active: scannerMode === 'manual' }" @click="switchToManual">Ввести код</button>
        <button type="button" :class="{ active: scannerMode === 'camera' }" @click="startCamera">Сканировать</button>
      </div>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>
    <div v-if="notice" class="alert success">{{ notice }}</div>

    <div class="scan-layout">
      <section class="panel action-panel">
        <div v-if="scannerMode === 'manual'" class="manual-card">
          <p class="eyebrow">Ручной ввод</p>
          <h2>Введите QR-код</h2>
          <label class="input-field">
            <span>Код грузового места</span>
            <input v-model.trim="qrCode" type="text" placeholder="QR-TPRO-MSK-240001" @keyup.enter="checkQr()" />
          </label>
          <button type="button" class="red-btn" :disabled="loading" @click="checkQr()">
            {{ loading ? 'Проверяем…' : 'Получить информацию' }}
          </button>
        </div>

        <div v-else class="camera-card">
          <div class="camera-head">
            <div>
              <p class="eyebrow">Камера</p>
              <h2>Наведите на QR</h2>
            </div>
            <button v-if="cameraState === 'running'" type="button" class="mini-btn" @click="stopCamera">Остановить</button>
            <button v-else type="button" class="mini-btn red" @click="startCamera">Включить</button>
          </div>
          <div class="camera-frame">
            <div :id="readerId" class="reader"></div>
            <div v-if="cameraState !== 'running'" class="camera-overlay">
              <strong>{{ cameraState === 'starting' ? 'Запрашиваем доступ…' : 'Камера не активна' }}</strong>
              <span>Нажмите “Включить” и разрешите доступ к камере.</span>
            </div>
          </div>
          <p v-if="cameraError" class="camera-message">{{ cameraError }}</p>
        </div>
      </section>

      <section class="panel details-panel" :class="{ empty: !hasCargo }">
        <template v-if="hasCargo">
          <div class="details-head">
            <div>
              <p class="eyebrow">Карточка груза</p>
              <h2>{{ cargoTitle(cargo) }}</h2>
            </div>
            <em :class="statusTone(cargo.status)">{{ statusLabel(cargo.status) }}</em>
          </div>
          <div class="info-grid">
            <div><span>ID места</span><strong>#{{ cargo.id }}</strong></div>
            <div><span>Заявка</span><strong>#{{ cargo.order_id || '—' }}</strong></div>
            <div><span>Зона</span><strong>{{ zoneTitle(cargo) }}</strong></div>
            <div><span>Гейт</span><strong>{{ gateTitle(cargo) }}</strong></div>
            <div class="wide"><span>Маршрут</span><strong>{{ routeTitle(cargo) }}</strong></div>
            <div class="wide"><span>Обновлено</span><strong>{{ formatDateTime(cargo.updated_at || cargo.created_at) }}</strong></div>
          </div>
          <button type="button" class="light-btn" @click="downloadCargoQr">Скачать QR этого места</button>
        </template>
        <template v-else>
          <div class="empty-state">
            <strong>Сканируйте QR</strong>
            <span>После проверки здесь появятся данные грузового места, маршрут, зона, гейт и текущий статус.</span>
          </div>
        </template>
      </section>
    </div>

    <section v-if="hasCargo" class="panel status-panel">
      <div class="panel-head">
        <div>
          <p class="eyebrow">Операция</p>
          <h2>Обновить статус</h2>
          <span class="hint">{{ statusHint }}</span>
        </div>
        <button type="button" class="dark-btn" :disabled="saving" @click="updateStatus">
          {{ saving ? 'Сохраняем…' : 'Сохранить статус' }}
        </button>
      </div>

      <div class="status-cards">
        <button
          v-for="item in workerStatusOptions"
          :key="item.value"
          type="button"
          class="status-card"
          :class="{ active: selectedStatus === item.value }"
          @click="selectStatus(item.value)"
        >
          <strong>{{ item.label }}</strong>
          <span>{{ item.description }}</span>
        </button>
      </div>

      <label class="input-field comment-field">
        <span>Комментарий</span>
        <input v-model.trim="comment" type="text" placeholder="Например: принято без повреждений" />
      </label>
    </section>

    <section v-if="hasCargo" class="panel history-panel">
      <p class="eyebrow">История груза</p>
      <div v-if="!history.length" class="empty-note">История пока недоступна.</div>
      <div v-else class="history-list">
        <article v-for="event in history" :key="event.id || `${event.status}-${event.created_at}`">
          <strong>{{ statusLabel(event.status || event.new_status) }}</strong>
          <span>{{ formatDateTime(event.created_at || event.changed_at) }}</span>
          <p v-if="event.comment">{{ event.comment }}</p>
        </article>
      </div>
    </section>
  </section>
</template>

<style scoped>
.scan-page { display: grid; gap: 26px; color: #061126; }
.scan-hero,
.panel { background: #fff; border-radius: 34px; padding: 32px; box-shadow: 0 18px 62px rgba(15, 23, 42, .08); }
.scan-hero { display: flex; align-items: flex-start; justify-content: space-between; gap: 24px; }
.eyebrow { margin: 0 0 10px; color: #ff3f4d; font-size: 13px; font-weight: 950; letter-spacing: .28em; text-transform: uppercase; }
h1 { margin: 0; font-size: clamp(44px, 6vw, 82px); line-height: .88; font-weight: 950; letter-spacing: -.06em; }
h2 { margin: 0; font-size: clamp(28px, 3vw, 42px); line-height: 1; font-weight: 950; letter-spacing: -.04em; }
.scan-hero span { display: block; margin-top: 16px; max-width: 760px; color: #5d6d83; font-size: 20px; line-height: 1.55; }
.mode-tabs { display: inline-flex; gap: 8px; padding: 8px; border-radius: 22px; background: #eef3f9; }
.mode-tabs button { min-height: 50px; border: 0; border-radius: 16px; padding: 0 18px; background: transparent; color: #5b6a80; font-weight: 950; cursor: pointer; }
.mode-tabs button.active { background: #061126; color: #fff; }
.scan-layout { display: grid; grid-template-columns: minmax(340px, .9fr) minmax(0, 1.1fr); gap: 22px; align-items: stretch; }
.action-panel,
.manual-card,
.camera-card { display: grid; gap: 20px; align-content: start; }
.input-field { display: grid; gap: 10px; }
.input-field span { color: #97a5bb; font-size: 13px; font-weight: 950; letter-spacing: .22em; text-transform: uppercase; }
.input-field input { width: 100%; min-height: 66px; border: 1px solid #dbe4ef; border-radius: 22px; background: #f8fbff; padding: 0 22px; color: #061126; font-size: 20px; font-weight: 950; outline: none; box-sizing: border-box; }
.input-field input:focus { border-color: #ff3f4d; box-shadow: 0 0 0 5px rgba(255,63,77,.12); background: #fff; }
.red-btn,
.dark-btn,
.light-btn,
.mini-btn { min-height: 62px; border: 0; border-radius: 20px; padding: 0 26px; color: #fff; font-size: 17px; font-weight: 950; cursor: pointer; }
.red-btn,
.mini-btn.red { background: #ff3f4d; box-shadow: 0 18px 42px rgba(255,63,77,.22); }
.dark-btn,
.mini-btn { background: #061126; }
.light-btn { margin-top: 18px; width: 100%; background: #eef3f9; color: #061126; }
.red-btn:disabled,
.dark-btn:disabled { opacity: .6; cursor: wait; }
.camera-head,
.details-head,
.panel-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 18px; }
.camera-frame { position: relative; min-height: 420px; border-radius: 28px; overflow: hidden; background: #071222; display: grid; place-items: center; }
.reader { width: 100%; min-height: 420px; }
.reader :deep(video) { width: 100% !important; height: 420px !important; object-fit: cover !important; border-radius: 28px; }
.reader :deep(img),
.reader :deep(button),
.reader :deep(select) { display: none !important; }
.camera-overlay { position: absolute; inset: 0; display: grid; place-items: center; text-align: center; padding: 28px; color: #fff; background: linear-gradient(135deg, rgba(7,18,34,.96), rgba(9,27,43,.86)); }
.camera-overlay strong { display: block; font-size: 28px; font-weight: 950; letter-spacing: -.04em; }
.camera-overlay span { display: block; margin-top: 10px; color: #a9b8ca; font-weight: 850; line-height: 1.45; }
.camera-message { margin: 0; padding: 14px 16px; border-radius: 18px; background: #fff7ed; color: #b45309; font-weight: 850; line-height: 1.45; }
.info-grid { margin-top: 22px; display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 14px; }
.info-grid div { border-radius: 22px; background: #f6f9fd; padding: 18px; display: grid; gap: 7px; }
.info-grid .wide { grid-column: 1 / -1; }
.info-grid span { color: #97a5bb; font-weight: 950; letter-spacing: .18em; text-transform: uppercase; font-size: 12px; }
.info-grid strong { font-size: 18px; font-weight: 950; overflow-wrap: anywhere; }
em { font-style: normal; padding: 10px 14px; border-radius: 999px; font-weight: 950; white-space: nowrap; text-align: center; }
em.green { background: #dcfce7; color: #047857; } em.red { background: #ffe4e6; color: #be123c; } em.blue { background: #dbeafe; color: #1d4ed8; } em.amber { background: #fef3c7; color: #b45309; } em.violet { background: #ede9fe; color: #6d28d9; } em.dark { background: #061126; color: #fff; } em.gray { background: #e2e8f0; color: #475569; }
.details-panel.empty { display: grid; place-items: center; min-height: 420px; }
.empty-state { max-width: 420px; text-align: center; display: grid; gap: 10px; }
.empty-state strong { font-size: 32px; font-weight: 950; letter-spacing: -.04em; }
.empty-state span { color: #66758a; font-size: 17px; line-height: 1.5; font-weight: 800; }
.status-panel { display: grid; gap: 22px; overflow: visible; }
.hint { display: block; margin-top: 12px; color: #64748b; font-weight: 850; line-height: 1.45; }
.status-cards { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 14px; }
.status-card { min-height: 110px; border: 1px solid #dbe4ef; border-radius: 24px; background: #f8fbff; padding: 18px; text-align: left; cursor: pointer; color: #061126; transition: transform .18s ease, border-color .18s ease, background .18s ease, box-shadow .18s ease; }
.status-card strong { display: block; font-size: 18px; font-weight: 950; line-height: 1.15; }
.status-card span { display: block; margin-top: 8px; color: #66758a; font-weight: 800; line-height: 1.35; }
.status-card:hover { transform: translateY(-2px); border-color: #ff9aa3; }
.status-card.active { background: #ff3f4d; border-color: #ff3f4d; color: #fff; box-shadow: 0 18px 42px rgba(255,63,77,.22); }
.status-card.active span { color: rgba(255,255,255,.86); }
.comment-field { max-width: 760px; }
.history-list { display: grid; gap: 12px; }
.history-list article { border-radius: 22px; background: #f6f9fd; padding: 16px 18px; }
.history-list strong { display: block; font-size: 18px; font-weight: 950; }
.history-list span { display: block; margin-top: 6px; color: #64748b; font-weight: 800; }
.history-list p { margin: 10px 0 0; color: #5c6c83; }
.alert,
.empty-note { padding: 18px 22px; border-radius: 22px; font-weight: 900; }
.alert.error { background: #fff0f1; color: #be123c; }
.alert.success { background: #ecfdf5; color: #047857; }
.empty-note { background: #f6f9fd; color: #64748b; }
@media (max-width: 1180px) { .scan-layout, .info-grid { grid-template-columns: 1fr; } .status-cards { grid-template-columns: repeat(2, 1fr); } .info-grid .wide { grid-column: auto; } }
@media (max-width: 760px) { .scan-hero, .camera-head, .details-head, .panel-head { flex-direction: column; } .mode-tabs { width: 100%; display: grid; grid-template-columns: 1fr 1fr; } .red-btn, .dark-btn, .light-btn, .mini-btn { width: 100%; } .status-cards { grid-template-columns: 1fr; } }
</style>
