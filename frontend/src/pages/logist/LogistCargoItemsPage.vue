<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="toolbar panel">
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

    <div class="filters panel">
      <BaseSelect v-model="filters.status" :options="cargoStatusOptions" placeholder="Все статусы" @change="loadCargoItems" />
      <BaseSelect v-model="filters.zoneId" :options="zoneOptions" placeholder="Все зоны" @change="loadCargoItems" />
      <BaseSelect v-model="filters.gateId" :options="gateOptions" placeholder="Все гейты" @change="loadCargoItems" />
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
            <span>Статус</span>
            <BaseSelect
              v-model="forms[cargo.id].status"
              :options="cargoStatusOptions.filter((item) => item.value)"
              placeholder="Выберите статус"
            />
          </label>
          <label>
            <span>Зона</span>
            <BaseSelect v-model="forms[cargo.id].storageZoneId" :options="zoneAssignOptions" placeholder="Не назначать" />
          </label>
          <label>
            <span>Гейт</span>
            <BaseSelect v-model="forms[cargo.id].gateId" :options="gateAssignOptions" placeholder="Не назначать" />
          </label>
          <label>
            <span>Комментарий</span>
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
import { computed, onMounted, reactive, ref } from 'vue'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
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

const zoneOptions = computed(() => [{ value: '', label: 'Все зоны' }, ...zones.value.map((zone) => ({ value: String(zone.id), label: zone.name, description: zone.warehouse_name || '' }))])
const gateOptions = computed(() => [{ value: '', label: 'Все гейты' }, ...gates.value.map((gate) => ({ value: String(gate.id), label: gate.name, description: gate.warehouse_name || '' }))])
const zoneAssignOptions = computed(() => [{ value: '', label: 'Не назначать' }, ...zones.value.map((zone) => ({ value: String(zone.id), label: zone.name, description: zone.warehouse_name || '' }))])
const gateAssignOptions = computed(() => [{ value: '', label: 'Не назначать' }, ...gates.value.map((gate) => ({ value: String(gate.id), label: gate.name, description: gate.warehouse_name || '' }))])

function ensureForm(cargo) {
  if (!forms[cargo.id]) {
    forms[cargo.id] = {
      status: cargo.status || 'accepted',
      storageZoneId: cargo.storage_zone_id ? String(cargo.storage_zone_id) : '',
      gateId: cargo.gate_id ? String(cargo.gate_id) : '',
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
      updated = unwrapOne(payload, 'cargo_item') || updated
    }
    if (form.gateId && Number(form.gateId) !== Number(updated.gate_id || 0)) {
      const payload = await apiFetch(`/cargo-items/${cargo.id}/assign-gate`, {
        auth: true,
        method: 'PATCH',
        body: { gate_id: Number(form.gateId), comment: form.comment || undefined },
      })
      updated = unwrapOne(payload, 'cargo_item') || updated
    }
    if (form.status && form.status !== updated.status) {
      const payload = await apiFetch(`/cargo-items/${cargo.id}/status`, {
        auth: true,
        method: 'PATCH',
        body: { status: form.status, comment: form.comment || undefined },
      })
      updated = unwrapOne(payload, 'cargo_item') || { ...updated, status: form.status }
    }
    const index = cargoItems.value.findIndex((item) => item.id === cargo.id)
    if (index !== -1) cargoItems.value[index] = updated
    forms[cargo.id] = {
      status: updated.status,
      storageZoneId: updated.storage_zone_id ? String(updated.storage_zone_id) : '',
      gateId: updated.gate_id ? String(updated.gate_id) : '',
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
.alert, .success { padding:16px 18px; border-radius:18px; font-weight:950; }
.alert { background:#fee2e2; color:#991b1b; }
.success { background:#d1fae5; color:#065f46; }
.panel, .scan-result, .cargo-card, .empty { background:white; border-radius:32px; padding:24px; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.toolbar { display:flex; align-items:end; justify-content:space-between; gap:18px; }
.toolbar p, .scan-result p { margin:0 0 8px; color:#ff3f4d; letter-spacing:.22em; text-transform:uppercase; font-size:12px; font-weight:950; }
.toolbar h2 { margin:0; font-size:34px; letter-spacing:-.04em; }
.scan-box, .filters { display:grid; grid-template-columns: minmax(220px, 1fr) auto; gap:12px; align-items:end; }
.filters { grid-template-columns: repeat(3, minmax(180px, 1fr)) auto; }
input { min-height:58px; border:1px solid #dbe3ef; border-radius:18px; padding:0 18px; background:#f6f8fb; color:#07101f; font-weight:950; font-family:inherit; }
button { min-height:58px; border:0; border-radius:18px; padding:0 20px; background:#ff3f4d; color:white; font-weight:950; cursor:pointer; font-family:inherit; }
button:disabled { opacity:.55; cursor:wait; }
.filters button { background:#07101f; }
.scan-result { display:flex; justify-content:space-between; gap:18px; background:#07101f; color:white; }
.scan-result h3 { margin:0 0 10px; font-size:30px; }
.scan-result span { color:#7dd3fc; font-weight:950; }
.result-meta { display:grid; gap:8px; text-align:right; }
.cargo-grid { display:grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap:20px; }
.cargo-card { display:flex; flex-direction:column; gap:18px; }
.card-top { display:flex; justify-content:space-between; gap:14px; align-items:start; }
.card-top span, .history summary, .action-box label > span { color:#94a3b8; letter-spacing:.22em; text-transform:uppercase; font-size:12px; font-weight:950; }
.card-top h3 { margin:8px 0 0; font-size:26px; letter-spacing:-.04em; overflow-wrap:anywhere; }
.card-top em { display:inline-flex; align-items:center; min-height:36px; padding:0 12px; border-radius:999px; background:#eef2ff; color:#3730a3; font-style:normal; font-weight:950; }
.meta-grid { display:grid; grid-template-columns:repeat(2, minmax(0,1fr)); gap:12px; }
.meta-grid div, .action-box { background:#f6f8fb; border-radius:22px; padding:16px; }
.meta-grid small { display:block; color:#64748b; font-weight:900; margin-bottom:6px; }
.meta-grid b { display:block; font-size:16px; overflow-wrap:anywhere; }
.action-box { display:grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap:14px; align-items:end; }
.action-box label { display:grid; gap:8px; }
.action-box button { grid-column:1 / -1; }
.history { background:#f8fafc; border-radius:22px; padding:16px; }
.history summary { cursor:pointer; color:#07101f; }
.history-list { display:grid; gap:10px; margin-top:14px; }
.history-item { display:grid; gap:4px; padding:12px; border-radius:16px; background:white; }
.history-item span, .history-item small, .muted { color:#64748b; font-weight:800; }
.empty { text-align:center; color:#64748b; font-weight:950; }
@media (max-width: 1180px) {
  .cargo-grid { grid-template-columns:1fr; }
  .toolbar { display:grid; }
  .filters { grid-template-columns:1fr 1fr; }
}
@media (max-width: 760px) {
  .scan-box, .filters, .action-box, .meta-grid { grid-template-columns:1fr; }
  .scan-result { display:grid; }
  .result-meta { text-align:left; }
}
</style>
