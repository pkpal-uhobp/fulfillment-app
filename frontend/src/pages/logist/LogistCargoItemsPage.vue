<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="toolbar panel">
      <div>
        <p>Грузы и QR</p>
        <h2>Грузовые места</h2>
        <span>Поиск, фильтрация по статусу, складу, зоне и гейту без отдельного QR-контроля.</span>
      </div>

      <button type="button" @click="loadCargoItems">Обновить</button>
    </div>

    <div class="filters panel">
      <label class="search-field">
        <span>Поиск</span>
        <input
          v-model.trim="filters.search"
          type="text"
          placeholder="QR, ID места, номер заявки, зона или гейт"
        />
      </label>

      <BaseSelect
        v-model="filters.status"
        :options="cargoStatusOptions"
        label="Статус"
        placeholder="Все статусы"
        @change="loadCargoItems"
      />

      <BaseSelect
        v-model="filters.warehouseId"
        :options="warehouseOptions"
        label="Склад"
        placeholder="Все склады"
      />

      <BaseSelect
        v-model="filters.zoneId"
        :options="zoneOptions"
        label="Зона"
        placeholder="Все зоны"
        @change="loadCargoItems"
      />

      <BaseSelect
        v-model="filters.gateId"
        :options="gateOptions"
        label="Гейт"
        placeholder="Все гейты"
        @change="loadCargoItems"
      />
    </div>

    <div class="summary-row">
      <article>
        <span>Всего</span>
        <strong>{{ filteredCargoItems.length }}</strong>
      </article>
      <article>
        <span>Принято</span>
        <strong>{{ counters.accepted }}</strong>
      </article>
      <article>
        <span>Хранение</span>
        <strong>{{ counters.stored }}</strong>
      </article>
      <article>
        <span>К отгрузке</span>
        <strong>{{ counters.ready }}</strong>
      </article>
    </div>

    <div class="cargo-grid">
      <article v-for="cargo in filteredCargoItems" :key="cargo.id" class="cargo-card">
        <div class="card-top">
          <div>
            <span>QR</span>
            <h3>{{ cargo.qr_code || `Место #${cargo.id}` }}</h3>
          </div>

          <em>{{ labelFromMap(cargoStatusLabels, cargo.status) }}</em>
        </div>

        <div class="meta-grid">
          <div><small>Заявка</small><b>#{{ cargo.order_id }}</b></div>
          <div><small>Тип места</small><b>#{{ cargo.cargo_place_type_id }}</b></div>
          <div><small>Склад</small><b>{{ cargoWarehouseName(cargo) }}</b></div>
          <div><small>Зона</small><b>{{ zoneName(cargo.storage_zone_id, zones) }}</b></div>
          <div><small>Гейт</small><b>{{ gateName(cargo.gate_id, gates) }}</b></div>
          <div><small>Обновлено</small><b>{{ formatDateTime(cargo.updated_at || cargo.received_at || cargo.created_at) }}</b></div>
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

    <div v-if="!filteredCargoItems.length && !error" class="empty">Грузовые места не найдены</div>
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
  warehouseName,
  zoneName,
} from './logistUtils'

const cargoItems = ref([])
const warehouses = ref([])
const zones = ref([])
const gates = ref([])
const forms = reactive({})
const history = reactive({})
const filters = reactive({ search: '', status: '', warehouseId: '', zoneId: '', gateId: '' })

const error = ref('')
const success = ref('')
const loadingId = ref(null)
const historyLoadingId = ref(null)

const warehouseOptions = computed(() => [
  { value: '', label: 'Все склады' },
  ...warehouses.value.map((warehouse) => ({
    value: String(warehouse.id),
    label: warehouseName(warehouse.id, warehouses.value),
    description: warehouse.address || warehouse.city || '',
  })),
])

const zoneOptions = computed(() => [
  { value: '', label: 'Все зоны' },
  ...zones.value.map((zone) => ({
    value: String(zone.id),
    label: zone.name,
    description: zone.warehouse_name || warehouseName(zone.warehouse_id, warehouses.value),
  })),
])

const gateOptions = computed(() => [
  { value: '', label: 'Все гейты' },
  ...gates.value.map((gate) => ({
    value: String(gate.id),
    label: gate.name,
    description: gate.warehouse_name || warehouseName(gate.warehouse_id, warehouses.value),
  })),
])

const zoneAssignOptions = computed(() => [
  { value: '', label: 'Не назначать' },
  ...zones.value.map((zone) => ({
    value: String(zone.id),
    label: zone.name,
    description: zone.warehouse_name || warehouseName(zone.warehouse_id, warehouses.value),
  })),
])

const gateAssignOptions = computed(() => [
  { value: '', label: 'Не назначать' },
  ...gates.value.map((gate) => ({
    value: String(gate.id),
    label: gate.name,
    description: gate.warehouse_name || warehouseName(gate.warehouse_id, warehouses.value),
  })),
])

const filteredCargoItems = computed(() => {
  const q = filters.search.trim().toLowerCase()

  return cargoItems.value.filter((cargo) => {
    const warehouseOk = !filters.warehouseId || cargoWarehouseIds(cargo).includes(String(filters.warehouseId))
    const searchText = [
      cargo.id,
      cargo.qr_code,
      cargo.order_id,
      cargo.cargo_place_type_id,
      cargo.status,
      labelFromMap(cargoStatusLabels, cargo.status),
      zoneName(cargo.storage_zone_id, zones.value),
      gateName(cargo.gate_id, gates.value),
      cargoWarehouseName(cargo),
    ]
      .join(' ')
      .toLowerCase()

    return warehouseOk && (!q || searchText.includes(q))
  })
})

const counters = computed(() => ({
  accepted: filteredCargoItems.value.filter((item) => item.status === 'accepted').length,
  stored: filteredCargoItems.value.filter((item) => item.status === 'stored').length,
  ready: filteredCargoItems.value.filter((item) => ['ready_to_ship', 'shipped'].includes(item.status)).length,
}))

function findZone(id) {
  return zones.value.find((zone) => Number(zone.id) === Number(id))
}

function findGate(id) {
  return gates.value.find((gate) => Number(gate.id) === Number(id))
}

function cargoWarehouseIds(cargo) {
  const ids = [
    cargo.warehouse_id,
    cargo.current_warehouse_id,
    cargo.receiving_warehouse_id,
    cargo.destination_warehouse_id,
    findZone(cargo.storage_zone_id)?.warehouse_id,
    findGate(cargo.gate_id)?.warehouse_id,
  ]
    .filter(Boolean)
    .map((id) => String(id))

  return [...new Set(ids)]
}

function cargoWarehouseName(cargo) {
  const ids = cargoWarehouseIds(cargo)

  if (!ids.length) return '—'

  return ids.map((id) => warehouseName(id, warehouses.value)).join(' / ')
}

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
  const [warehousesPayload, zonesPayload, gatesPayload] = await Promise.all([
    apiFetch('/warehouses'),
    apiFetch('/storage-zones', { auth: true }),
    apiFetch('/gates', { auth: true }),
  ])

  warehouses.value = unwrapList(warehousesPayload, 'warehouses')
  zones.value = unwrapList(zonesPayload, 'storage_zones')
  gates.value = unwrapList(gatesPayload, 'gates')
}

async function loadCargoItems() {
  error.value = ''
  success.value = ''

  const params = new URLSearchParams({ limit: '200' })

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

    success.value = `Груз ${updated.qr_code || `#${updated.id}`} обновлён`
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
.page {
  display: grid;
  gap: 20px;
}

.alert,
.success {
  padding: 16px 18px;
  border-radius: 18px;
  font-weight: 950;
}

.alert {
  background: #fee2e2;
  color: #991b1b;
}

.success {
  background: #d1fae5;
  color: #065f46;
}

.panel,
.cargo-card,
.empty,
.summary-row article {
  background: white;
  border-radius: 32px;
  padding: 24px;
  box-shadow: 0 18px 42px rgba(7, 16, 31, .08);
}

.toolbar {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 18px;
}

.toolbar p {
  margin: 0 0 8px;
  color: #ff3f4d;
  letter-spacing: .22em;
  text-transform: uppercase;
  font-size: 12px;
  font-weight: 950;
}

.toolbar h2 {
  margin: 0;
  font-size: 42px;
  line-height: 1;
  letter-spacing: -.05em;
}

.toolbar span {
  display: block;
  margin-top: 12px;
  max-width: 720px;
  color: #64748b;
  font-weight: 850;
  line-height: 1.5;
}

.toolbar button,
.action-box button {
  min-height: 58px;
  border: 0;
  border-radius: 18px;
  padding: 0 20px;
  background: #ff3f4d;
  color: white;
  font-weight: 950;
  cursor: pointer;
  font-family: inherit;
}

.toolbar button {
  background: #07101f;
}

button:disabled {
  opacity: .55;
  cursor: wait;
}

.filters {
  display: grid;
  grid-template-columns: minmax(260px, 1.35fr) repeat(4, minmax(180px, 1fr));
  gap: 14px;
  align-items: end;
}

.search-field {
  display: grid;
  gap: 10px;
}

.search-field span,
.summary-row span,
.card-top span,
.history summary,
.action-box label > span {
  color: #94a3b8;
  letter-spacing: .22em;
  text-transform: uppercase;
  font-size: 12px;
  font-weight: 950;
}

input {
  min-height: 64px;
  border: 1px solid #dbe3ef;
  border-radius: 22px;
  padding: 0 20px;
  background: #f6f8fb;
  color: #07101f;
  font-size: 17px;
  font-weight: 950;
  font-family: inherit;
  outline: none;
}

input:focus {
  border-color: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .10);
  background: #fff;
}

.summary-row {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.summary-row article {
  display: grid;
  gap: 8px;
  min-height: 110px;
  align-content: center;
}

.summary-row strong {
  font-size: 42px;
  line-height: 1;
  letter-spacing: -.05em;
  font-weight: 950;
}

.cargo-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.cargo-card {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.card-top {
  display: flex;
  justify-content: space-between;
  gap: 14px;
  align-items: start;
}

.card-top h3 {
  margin: 8px 0 0;
  font-size: 26px;
  letter-spacing: -.04em;
  overflow-wrap: anywhere;
}

.card-top em {
  display: inline-flex;
  align-items: center;
  min-height: 36px;
  padding: 0 12px;
  border-radius: 999px;
  background: #eef2ff;
  color: #3730a3;
  font-style: normal;
  font-weight: 950;
  white-space: nowrap;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.meta-grid div,
.action-box {
  background: #f6f8fb;
  border-radius: 22px;
  padding: 16px;
}

.meta-grid small {
  display: block;
  color: #64748b;
  font-weight: 900;
  margin-bottom: 6px;
}

.meta-grid b {
  display: block;
  font-size: 16px;
  overflow-wrap: anywhere;
}

.action-box {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  align-items: end;
}

.action-box label {
  display: grid;
  gap: 8px;
}

.action-box button {
  grid-column: 1 / -1;
}

.history {
  background: #f8fafc;
  border-radius: 22px;
  padding: 16px;
}

.history summary {
  cursor: pointer;
  color: #07101f;
}

.history-list {
  display: grid;
  gap: 10px;
  margin-top: 14px;
}

.history-item {
  display: grid;
  gap: 4px;
  padding: 12px;
  border-radius: 16px;
  background: white;
}

.history-item span,
.history-item small,
.muted {
  color: #64748b;
  font-weight: 800;
}

.empty {
  text-align: center;
  color: #64748b;
  font-weight: 950;
}

@media (max-width: 1280px) {
  .filters {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 1180px) {
  .cargo-grid {
    grid-template-columns: 1fr;
  }

  .toolbar {
    display: grid;
  }
}

@media (max-width: 760px) {
  .filters,
  .summary-row,
  .action-box,
  .meta-grid {
    grid-template-columns: 1fr;
  }

  .toolbar button {
    width: 100%;
  }
}
</style>
