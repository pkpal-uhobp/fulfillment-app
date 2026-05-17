<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="toolbar">
      <div>
        <p>QR-контроль</p>
        <h2>Грузовые места</h2>
      </div>
      <div class="scan-box">
        <input v-model.trim="qrCode" placeholder="QR-TPRO-MSK-240001" @keyup.enter="scanQr" />
        <button type="button" @click="scanQr">Проверить QR</button>
      </div>
    </div>

    <div v-if="scanned" class="scan-result">
      <div>
        <p>Найдено грузовое место</p>
        <h3>{{ scanned.qr_code }}</h3>
        <span>{{ labelFromMap(cargoStatusLabels, scanned.status) }}</span>
      </div>
      <div class="result-meta">
        <b>Заявка #{{ scanned.order_id }}</b>
        <b>{{ zoneName(scanned.storage_zone_id, zones) }}</b>
        <b>{{ gateName(scanned.gate_id, gates) }}</b>
      </div>
    </div>

    <div class="filters">
      <select v-model="filters.status" @change="loadCargoItems">
        <option v-for="option in cargoStatusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
      </select>
      <select v-model="filters.zoneId" @change="loadCargoItems">
        <option value="">Все зоны</option>
        <option v-for="zone in zones" :key="zone.id" :value="zone.id">{{ zone.name }}</option>
      </select>
      <select v-model="filters.gateId" @change="loadCargoItems">
        <option value="">Все гейты</option>
        <option v-for="gate in gates" :key="gate.id" :value="gate.id">{{ gate.name }}</option>
      </select>
      <button type="button" @click="loadCargoItems">Обновить</button>
    </div>

    <div class="cargo-grid">
      <article v-for="cargo in cargoItems" :key="cargo.id" class="cargo-card">
        <div class="card-top">
          <div>
            <span>QR</span>
            <h3>{{ cargo.qr_code }}</h3>
          </div>
          <em>{{ labelFromMap(cargoStatusLabels, cargo.status) }}</em>
        </div>

        <div class="meta-grid">
          <div><small>Заявка</small><b>#{{ cargo.order_id }}</b></div>
          <div><small>Тип места</small><b>#{{ cargo.cargo_place_type_id }}</b></div>
          <div><small>Зона</small><b>{{ zoneName(cargo.storage_zone_id, zones) }}</b></div>
          <div><small>Гейт</small><b>{{ gateName(cargo.gate_id, gates) }}</b></div>
          <div><small>Принято</small><b>{{ formatDateTime(cargo.received_at) }}</b></div>
          <div><small>Отгружено</small><b>{{ formatDateTime(cargo.shipped_at) }}</b></div>
        </div>

        <div class="action-box">
          <label>
            Статус
            <select v-model="forms[cargo.id].status">
              <option v-for="option in cargoStatusOptions.filter((item) => item.value)" :key="option.value" :value="option.value">{{ option.label }}</option>
            </select>
          </label>
          <label>
            Зона
            <select v-model="forms[cargo.id].storageZoneId">
              <option value="">Не назначать</option>
              <option v-for="zone in zones" :key="zone.id" :value="zone.id">{{ zone.name }}</option>
            </select>
          </label>
          <label>
            Гейт
            <select v-model="forms[cargo.id].gateId">
              <option value="">Не назначать</option>
              <option v-for="gate in gates" :key="gate.id" :value="gate.id">{{ gate.name }}</option>
            </select>
          </label>
          <label>
            Комментарий
            <input v-model.trim="forms[cargo.id].comment" placeholder="Комментарий операции" />
          </label>
          <button type="button" :disabled="loadingId === cargo.id" @click="applyCargoChanges(cargo)">
            {{ loadingId === cargo.id ? 'Сохраняем...' : 'Сохранить изменения' }}
          </button>
        </div>

        <details class="history" @toggle="onHistoryToggle($event, cargo.id)">
          <summary>История груза</summary>
          <div v-if="historyLoadingId === cargo.id" class="muted">Загружаем историю...</div>
          <div v-else-if="history[cargo.id]?.length" class="history-list">
            <div v-for="item in history[cargo.id]" :key="item.id" class="history-item">
              <b>{{ labelFromMap(cargoStatusLabels, item.new_status) }}</b>
              <span>{{ formatDateTime(item.changed_at) }}</span>
              <small>{{ item.comment || 'Без комментария' }}</small>
            </div>
          </div>
          <div v-else class="muted">Истории пока нет</div>
        </details>
      </article>
    </div>

    <div v-if="!cargoItems.length && !error" class="empty">Грузовые места не найдены</div>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'
import {
  cargoStatusLabels,
  cargoStatusOptions,
  formatDateTime,
  gateName,
  labelFromMap,
  unwrapList,
  unwrapOne,
  zoneName,
} from './logistUtils'

const cargoItems = ref([])
const zones = ref([])
const gates = ref([])
const forms = reactive({})
const history = reactive({})
const filters = reactive({ status: '', zoneId: '', gateId: '' })
const qrCode = ref('QR-TPRO-MSK-240001')
const scanned = ref(null)
const error = ref('')
const success = ref('')
const loadingId = ref(null)
const historyLoadingId = ref(null)

function ensureForm(cargo) {
  if (!forms[cargo.id]) {
    forms[cargo.id] = {
      status: cargo.status || 'accepted',
      storageZoneId: cargo.storage_zone_id || '',
      gateId: cargo.gate_id || '',
      comment: '',
    }
  }
}

async function loadCatalogs() {
  const [zonesPayload, gatesPayload] = await Promise.all([
    apiFetch('/storage-zones', { auth: true }),
    apiFetch('/gates', { auth: true }),
  ])
  zones.value = unwrapList(zonesPayload, 'storage_zones')
  gates.value = unwrapList(gatesPayload, 'gates')
}

async function loadCargoItems() {
  error.value = ''
  success.value = ''
  const params = new URLSearchParams({ limit: '100' })
  if (filters.status) params.set('status', filters.status)
  if (filters.zoneId) params.set('storage_zone_id', filters.zoneId)
  if (filters.gateId) params.set('gate_id', filters.gateId)
  try {
    const payload = await apiFetch(`/cargo-items?${params}`, { auth: true })
    cargoItems.value = unwrapList(payload, 'cargo_items')
    cargoItems.value.forEach(ensureForm)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить грузовые места'
  }
}

async function scanQr() {
  if (!qrCode.value) {
    error.value = 'Введите QR-код'
    return
  }
  error.value = ''
  success.value = ''
  try {
    const payload = await apiFetch(`/cargo-items/scan?qr_code=${encodeURIComponent(qrCode.value)}`, { auth: true })
    scanned.value = unwrapOne(payload, 'cargo_item')
    if (scanned.value) {
      ensureForm(scanned.value)
      success.value = `QR ${scanned.value.qr_code} найден`
    }
  } catch (err) {
    error.value = err.message || 'Не удалось проверить QR-код'
  }
}

async function applyCargoChanges(cargo) {
  loadingId.value = cargo.id
  error.value = ''
  success.value = ''
  const form = forms[cargo.id]
  try {
    let updated = cargo
    if (form.storageZoneId && Number(form.storageZoneId) !== Number(cargo.storage_zone_id || 0)) {
      const payload = await apiFetch(`/cargo-items/${cargo.id}/assign-zone`, {
        auth: true,
        method: 'PATCH',
        body: { storage_zone_id: Number(form.storageZoneId), comment: form.comment || undefined },
      })
      updated = unwrapOne(payload, 'cargo_item')
    }
    if (form.gateId && Number(form.gateId) !== Number(updated.gate_id || 0)) {
      const payload = await apiFetch(`/cargo-items/${cargo.id}/assign-gate`, {
        auth: true,
        method: 'PATCH',
        body: { gate_id: Number(form.gateId), comment: form.comment || undefined },
      })
      updated = unwrapOne(payload, 'cargo_item')
    }
    if (form.status && form.status !== updated.status) {
      const payload = await apiFetch(`/cargo-items/${cargo.id}/status`, {
        auth: true,
        method: 'PATCH',
        body: { status: form.status, comment: form.comment || undefined },
      })
      updated = unwrapOne(payload, 'cargo_item')
    }
    const index = cargoItems.value.findIndex((item) => item.id === cargo.id)
    if (index !== -1) cargoItems.value[index] = updated
    forms[cargo.id] = {
      status: updated.status,
      storageZoneId: updated.storage_zone_id || '',
      gateId: updated.gate_id || '',
      comment: '',
    }
    success.value = `Груз ${updated.qr_code} обновлён`
    delete history[cargo.id]
  } catch (err) {
    error.value = err.message || 'Не удалось сохранить изменения'
  } finally {
    loadingId.value = null
  }
}

async function onHistoryToggle(event, cargoId) {
  if (!event.target.open || history[cargoId]) return
  historyLoadingId.value = cargoId
  try {
    const payload = await apiFetch(`/cargo-items/${cargoId}/history`, { auth: true })
    history[cargoId] = unwrapList(payload, 'history')
  } catch {
    history[cargoId] = []
  } finally {
    historyLoadingId.value = null
  }
}

onMounted(async () => {
  await Promise.all([loadCatalogs(), loadCargoItems()])
})
</script>

<style scoped>
.page { display:grid; gap:20px; }
.alert, .success { padding:16px 18px; border-radius:18px; font-weight:900; }
.alert { background:#fee2e2; color:#991b1b; }
.success { background:#d1fae5; color:#065f46; }
.toolbar, .filters, .scan-result { background:white; border-radius:32px; padding:24px; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.toolbar { display:flex; align-items:end; justify-content:space-between; gap:18px; }
.toolbar p { margin:0 0 8px; color:#ff3f4d; letter-spacing:.22em; text-transform:uppercase; font-size:12px; font-weight:900; }
.toolbar h2 { margin:0; font-size:34px; letter-spacing:-.04em; }
.scan-box, .filters { display:flex; gap:12px; flex-wrap:wrap; }
input, select { height:48px; border:1px solid #dbe3ef; border-radius:16px; padding:0 14px; background:#f6f8fb; color:#07101f; font-weight:800; }
button { height:48px; border:0; border-radius:16px; padding:0 18px; background:#ff3f4d; color:white; font-weight:900; cursor:pointer; }
button:disabled { opacity:.55; cursor:wait; }
.scan-result { display:flex; justify-content:space-between; gap:18px; background:#07101f; color:white; }
.scan-result p { margin:0 0 8px; color:#ff9da5; letter-spacing:.22em; text-transform:uppercase; font-weight:900; font-size:12px; }
.scan-result h3 { margin:0 0 10px; font-size:30px; }
.scan-result span { color:#7dd3fc; font-weight:900; }
.result-meta { display:grid; gap:8px; text-align:right; }
.cargo-grid { display:grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap:20px; }
.cargo-card { display:flex; flex-direction:column; gap:18px; padding:24px; border-radius:32px; background:white; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.card-top { display:flex; justify-content:space-between; gap:14px; align-items:start; }
.card-top span { color:#ff3f4d; letter-spacing:.22em; text-transform:uppercase; font-size:12px; font-weight:900; }
.card-top h3 { margin:8px 0 0; font-size:28px; letter-spacing:-.04em; }
.card-top em { font-style:normal; padding:10px 12px; border-radius:999px; background:#ffe6e8; color:#ff3f4d; font-weight:900; white-space:nowrap; }
.meta-grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:12px; }
.meta-grid div { padding:16px; border-radius:20px; background:#f6f8fb; display:grid; gap:5px; }
small, .muted { color:#6b7b91; font-weight:800; }
.action-box { margin-top:auto; display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:12px; padding:16px; border-radius:24px; background:#f6f8fb; }
.action-box label { display:grid; gap:8px; font-size:12px; text-transform:uppercase; letter-spacing:.14em; color:#8b9ab0; font-weight:900; }
.action-box button { grid-column:1 / -1; }
.history { border-top:1px solid #edf1f7; padding-top:14px; }
summary { cursor:pointer; font-weight:900; }
.history-list { display:grid; gap:10px; margin-top:12px; }
.history-item { display:grid; gap:4px; padding:14px; border-radius:18px; background:#f6f8fb; }
.empty { padding:26px; border-radius:26px; background:white; font-weight:900; color:#63738a; }
@media (max-width:1100px) { .cargo-grid { grid-template-columns:1fr; } .toolbar, .scan-result { flex-direction:column; align-items:stretch; } }
@media (max-width:640px) { .meta-grid, .action-box { grid-template-columns:1fr; } .card-top { flex-direction:column; } }
</style>
