<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { apiFetch } from '@/shared/api/http'

const warehouses = ref([])
const zones = ref([])
const gates = ref([])
const selectedWarehouseId = ref('')
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const notice = ref('')

const warehouseForm = reactive({
  name: '',
  warehouse_type: 'both',
  marketplace: '',
  city: '',
  address: '',
})

const zoneForm = reactive({
  name: '',
  description: '',
})

const gateForm = reactive({
  name: '',
})

const warehouseTypes = [
  { value: 'receiving', label: 'Склад приёмки' },
  { value: 'destination', label: 'Склад назначения' },
  { value: 'both', label: 'Приёмка и назначение' },
]

const typeLabels = Object.fromEntries(warehouseTypes.map((item) => [item.value, item.label]))

function collection(payload, keys) {
  if (Array.isArray(payload)) return payload

  for (const key of keys) {
    if (Array.isArray(payload?.[key])) return payload[key]
  }

  return []
}

const selectedWarehouse = computed(() => {
  return warehouses.value.find((warehouse) => String(warehouse.id) === String(selectedWarehouseId.value)) || null
})

const selectedZones = computed(() => {
  return zones.value.filter((zone) => String(zone.warehouse_id) === String(selectedWarehouseId.value))
})

const selectedGates = computed(() => {
  return gates.value.filter((gate) => String(gate.warehouse_id) === String(selectedWarehouseId.value))
})

const counters = computed(() => ({
  warehouses: warehouses.value.length,
  activeWarehouses: warehouses.value.filter((warehouse) => warehouse.is_active).length,
  zones: zones.value.length,
  gates: gates.value.length,
}))

function resetWarehouseForm() {
  warehouseForm.name = ''
  warehouseForm.warehouse_type = 'both'
  warehouseForm.marketplace = ''
  warehouseForm.city = ''
  warehouseForm.address = ''
}

function resetZoneForm() {
  zoneForm.name = ''
  zoneForm.description = ''
}

function resetGateForm() {
  gateForm.name = ''
}

async function loadData() {
  loading.value = true
  error.value = ''

  try {
    const [warehousesPayload, zonesPayload, gatesPayload] = await Promise.all([
      apiFetch('/warehouses', { auth: true }),
      apiFetch('/storage-zones', { auth: true }),
      apiFetch('/gates', { auth: true }),
    ])

    warehouses.value = collection(warehousesPayload, ['warehouses', 'items', 'data'])
    zones.value = collection(zonesPayload, ['storage_zones', 'storageZones', 'items', 'data'])
    gates.value = collection(gatesPayload, ['gates', 'items', 'data'])

    if (!selectedWarehouseId.value && warehouses.value.length) {
      selectedWarehouseId.value = String(warehouses.value[0].id)
    }
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить складскую структуру'
  } finally {
    loading.value = false
  }
}

async function createWarehouse() {
  saving.value = true
  error.value = ''
  notice.value = ''

  try {
    const payload = await apiFetch('/warehouses', {
      method: 'POST',
      auth: true,
      body: {
        name: warehouseForm.name,
        warehouse_type: warehouseForm.warehouse_type,
        marketplace: warehouseForm.marketplace || undefined,
        city: warehouseForm.city,
        address: warehouseForm.address,
      },
    })

    notice.value = 'Склад создан'
    resetWarehouseForm()
    await loadData()

    const created = payload?.warehouse
    if (created?.id) selectedWarehouseId.value = String(created.id)
  } catch (err) {
    error.value = err.message || 'Не удалось создать склад'
  } finally {
    saving.value = false
  }
}

async function patchWarehouse(warehouse, patch) {
  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/warehouses/${warehouse.id}`, {
      method: 'PATCH',
      auth: true,
      body: patch,
    })

    notice.value = 'Склад обновлён'
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось обновить склад'
  }
}

async function deleteWarehouse(warehouse) {
  if (!confirm(`Удалить склад «${warehouse.name}»?`)) return

  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/warehouses/${warehouse.id}`, {
      method: 'DELETE',
      auth: true,
    })

    notice.value = 'Склад удалён'
    selectedWarehouseId.value = ''
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось удалить склад'
  }
}

async function createZone() {
  if (!selectedWarehouse.value) return

  saving.value = true
  error.value = ''
  notice.value = ''

  try {
    await apiFetch('/storage-zones', {
      method: 'POST',
      auth: true,
      body: {
        warehouse_id: Number(selectedWarehouseId.value),
        name: zoneForm.name,
        description: zoneForm.description || undefined,
      },
    })

    notice.value = 'Зона хранения создана'
    resetZoneForm()
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось создать зону хранения'
  } finally {
    saving.value = false
  }
}

async function patchZone(zone, patch) {
  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/storage-zones/${zone.id}`, {
      method: 'PATCH',
      auth: true,
      body: patch,
    })

    notice.value = 'Зона обновлена'
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось обновить зону'
  }
}

async function deleteZone(zone) {
  if (!confirm(`Удалить зону «${zone.name}»?`)) return

  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/storage-zones/${zone.id}`, {
      method: 'DELETE',
      auth: true,
    })

    notice.value = 'Зона удалена'
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось удалить зону'
  }
}

async function createGate() {
  if (!selectedWarehouse.value) return

  saving.value = true
  error.value = ''
  notice.value = ''

  try {
    await apiFetch('/gates', {
      method: 'POST',
      auth: true,
      body: {
        warehouse_id: Number(selectedWarehouseId.value),
        name: gateForm.name,
      },
    })

    notice.value = 'Гейт создан'
    resetGateForm()
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось создать гейт'
  } finally {
    saving.value = false
  }
}

async function patchGate(gate, patch) {
  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/gates/${gate.id}`, {
      method: 'PATCH',
      auth: true,
      body: patch,
    })

    notice.value = 'Гейт обновлён'
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось обновить гейт'
  }
}

async function deleteGate(gate) {
  if (!confirm(`Удалить гейт «${gate.name}»?`)) return

  error.value = ''
  notice.value = ''

  try {
    await apiFetch(`/gates/${gate.id}`, {
      method: 'DELETE',
      auth: true,
    })

    notice.value = 'Гейт удалён'
    await loadData()
  } catch (err) {
    error.value = err.message || 'Не удалось удалить гейт'
  }
}

watch(warehouses, () => {
  if (selectedWarehouseId.value && !selectedWarehouse.value && warehouses.value.length) {
    selectedWarehouseId.value = String(warehouses.value[0].id)
  }
})

onMounted(loadData)
</script>

<template>
  <section class="admin-page">
    <header class="hero-card">
      <div>
        <p class="eyebrow">Склады</p>
        <h1>Складская структура</h1>
        <span>Управляйте складами, зонами хранения и гейтами отгрузки в одном месте.</span>
      </div>

      <button class="dark-btn" type="button" :disabled="loading" @click="loadData">
        {{ loading ? 'Загрузка…' : 'Обновить' }}
      </button>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>
    <div v-if="notice" class="alert success">{{ notice }}</div>

    <section class="stats-grid">
      <article>
        <span>Складов</span>
        <strong>{{ counters.warehouses }}</strong>
      </article>

      <article>
        <span>Активных</span>
        <strong>{{ counters.activeWarehouses }}</strong>
      </article>

      <article>
        <span>Зон хранения</span>
        <strong>{{ counters.zones }}</strong>
      </article>

      <article>
        <span>Гейтов</span>
        <strong>{{ counters.gates }}</strong>
      </article>
    </section>

    <section class="workspace-grid">
      <aside class="panel-card">
        <p class="eyebrow">Новый склад</p>
        <h2>Создать склад</h2>

        <form class="form-grid" @submit.prevent="createWarehouse">
          <label>
            <span>Название</span>
            <input v-model.trim="warehouseForm.name" required type="text" placeholder="Склад приёмки" />
          </label>

          <label>
            <span>Тип</span>
            <select v-model="warehouseForm.warehouse_type">
              <option v-for="type in warehouseTypes" :key="type.value" :value="type.value">
                {{ type.label }}
              </option>
            </select>
          </label>

          <label>
            <span>Маркетплейс</span>
            <input v-model.trim="warehouseForm.marketplace" type="text" placeholder="Ozon, WB или свой" />
          </label>

          <label>
            <span>Город</span>
            <input v-model.trim="warehouseForm.city" required type="text" placeholder="Москва" />
          </label>

          <label>
            <span>Адрес</span>
            <input v-model.trim="warehouseForm.address" required type="text" placeholder="Улица, дом" />
          </label>

          <button class="red-btn" type="submit" :disabled="saving">
            {{ saving ? 'Сохраняем…' : 'Создать склад' }}
          </button>
        </form>
      </aside>

      <section class="panel-card warehouses-card">
        <div class="section-head">
          <div>
            <p class="eyebrow">Список</p>
            <h2>Склады</h2>
          </div>
        </div>

        <div v-if="!warehouses.length" class="empty">Складов пока нет.</div>

        <div v-else class="warehouse-list">
          <article
            v-for="warehouse in warehouses"
            :key="warehouse.id"
            class="warehouse-row"
            :class="{ active: String(selectedWarehouseId) === String(warehouse.id) }"
            @click="selectedWarehouseId = String(warehouse.id)"
          >
            <div>
              <strong>{{ warehouse.name }}</strong>
              <span>{{ warehouse.city }} · {{ warehouse.address }}</span>
              <small>{{ typeLabels[warehouse.warehouse_type] || warehouse.warehouse_type }}</small>
            </div>

            <em :class="{ off: !warehouse.is_active }">
              {{ warehouse.is_active ? 'Активен' : 'Отключён' }}
            </em>

            <button
              type="button"
              class="small-btn"
              @click.stop="patchWarehouse(warehouse, { is_active: !warehouse.is_active })"
            >
              {{ warehouse.is_active ? 'Отключить' : 'Включить' }}
            </button>

            <button type="button" class="small-btn danger" @click.stop="deleteWarehouse(warehouse)">
              Удалить
            </button>
          </article>
        </div>
      </section>
    </section>

    <section class="structure-grid">
      <article class="panel-card">
        <p class="eyebrow">Зоны хранения</p>
        <h2>{{ selectedWarehouse ? selectedWarehouse.name : 'Выберите склад' }}</h2>

        <form class="inline-form" @submit.prevent="createZone">
          <input v-model.trim="zoneForm.name" required type="text" placeholder="Название зоны" :disabled="!selectedWarehouse" />
          <input v-model.trim="zoneForm.description" type="text" placeholder="Описание" :disabled="!selectedWarehouse" />
          <button class="red-btn" type="submit" :disabled="!selectedWarehouse || saving">Добавить</button>
        </form>

        <div v-if="!selectedZones.length" class="empty">У выбранного склада пока нет зон хранения.</div>

        <div v-else class="mini-list">
          <article v-for="zone in selectedZones" :key="zone.id" class="mini-row">
            <div>
              <strong>{{ zone.name }}</strong>
              <span>{{ zone.description || 'Без описания' }}</span>
            </div>

            <em :class="{ off: !zone.is_active }">{{ zone.is_active ? 'Активна' : 'Отключена' }}</em>

            <button class="small-btn" type="button" @click="patchZone(zone, { is_active: !zone.is_active })">
              {{ zone.is_active ? 'Отключить' : 'Включить' }}
            </button>

            <button class="small-btn danger" type="button" @click="deleteZone(zone)">Удалить</button>
          </article>
        </div>
      </article>

      <article class="panel-card">
        <p class="eyebrow">Гейты</p>
        <h2>{{ selectedWarehouse ? selectedWarehouse.name : 'Выберите склад' }}</h2>

        <form class="inline-form two" @submit.prevent="createGate">
          <input v-model.trim="gateForm.name" required type="text" placeholder="Название гейта" :disabled="!selectedWarehouse" />
          <button class="red-btn" type="submit" :disabled="!selectedWarehouse || saving">Добавить</button>
        </form>

        <div v-if="!selectedGates.length" class="empty">У выбранного склада пока нет гейтов.</div>

        <div v-else class="mini-list">
          <article v-for="gate in selectedGates" :key="gate.id" class="mini-row gate-row">
            <div>
              <strong>{{ gate.name }}</strong>
              <span>Склад #{{ gate.warehouse_id }}</span>
            </div>

            <em :class="{ off: !gate.is_active }">{{ gate.is_active ? 'Активен' : 'Отключён' }}</em>

            <button class="small-btn" type="button" @click="patchGate(gate, { is_active: !gate.is_active })">
              {{ gate.is_active ? 'Отключить' : 'Включить' }}
            </button>

            <button class="small-btn danger" type="button" @click="deleteGate(gate)">Удалить</button>
          </article>
        </div>
      </article>
    </section>
  </section>
</template>

<style scoped>
.admin-page {
  display: grid;
  gap: 26px;
}

.hero-card,
.panel-card,
.stats-grid article {
  background: #fff;
  border-radius: 34px;
  box-shadow: 0 18px 62px rgba(15, 23, 42, .08);
}

.hero-card {
  padding: 34px;
  display: flex;
  justify-content: space-between;
  gap: 24px;
}

.eyebrow {
  margin: 0 0 12px;
  color: #ff3f4d;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .28em;
  text-transform: uppercase;
}

h1,
h2 {
  margin: 0;
  color: #061126;
  font-weight: 950;
  letter-spacing: -.06em;
}

h1 {
  font-size: clamp(48px, 7vw, 84px);
  line-height: .9;
}

h2 {
  font-size: clamp(28px, 3vw, 42px);
  line-height: 1;
}

.hero-card span {
  display: block;
  margin-top: 14px;
  color: #5d6d83;
  font-size: 18px;
  font-weight: 800;
  line-height: 1.5;
}

.dark-btn,
.red-btn,
.small-btn {
  border: 0;
  border-radius: 20px;
  font-weight: 950;
  cursor: pointer;
}

.dark-btn,
.red-btn {
  min-height: 58px;
  padding: 0 24px;
  color: #fff;
  font-size: 16px;
}

.dark-btn {
  background: #061126;
}

.red-btn {
  background: #ff3f4d;
  box-shadow: 0 18px 42px rgba(255, 63, 77, .2);
}

.alert,
.empty {
  padding: 18px 22px;
  border-radius: 22px;
  font-weight: 900;
}

.alert.error {
  background: #fff0f1;
  color: #be123c;
}

.alert.success {
  background: #e8fff5;
  color: #047857;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
}

.stats-grid article {
  min-height: 130px;
  padding: 24px;
  display: grid;
  align-content: center;
}

.stats-grid span,
label span {
  color: #97a5bb;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .18em;
  text-transform: uppercase;
}

.stats-grid strong {
  margin-top: 10px;
  font-size: 48px;
  line-height: 1;
  font-weight: 950;
}

.workspace-grid,
.structure-grid {
  display: grid;
  grid-template-columns: minmax(330px, .66fr) minmax(0, 1.34fr);
  gap: 22px;
  align-items: start;
}

.structure-grid {
  grid-template-columns: 1fr;
}

.panel-card {
  padding: 30px;
}

.form-grid,
.inline-form {
  margin-top: 22px;
  display: grid;
  gap: 14px;
}

.inline-form {
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr) auto;
  align-items: end;
}

.inline-form.two {
  grid-template-columns: minmax(0, 1fr) auto;
}

label {
  display: grid;
  gap: 8px;
}

input,
select {
  width: 100%;
  min-height: 54px;
  border: 1px solid #dbe4ef;
  border-radius: 18px;
  background: #f8fbff;
  color: #061126;
  padding: 0 16px;
  font-size: 16px;
  font-weight: 850;
  outline: none;
  box-sizing: border-box;
}

input:focus,
select:focus {
  border-color: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .12);
  background: #fff;
}

.section-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.warehouse-list,
.mini-list {
  display: grid;
  gap: 12px;
  max-height: 560px;
  overflow-y: auto;
  padding-right: 6px;
}

.warehouse-row,
.mini-row {
  border: 1px solid #dbe4ef;
  border-radius: 24px;
  background: #f8fbff;
  padding: 16px;
  display: grid;
  grid-template-columns: minmax(220px, 1fr) auto auto auto;
  gap: 12px;
  align-items: center;
  cursor: pointer;
  min-width: 0;
}

.warehouse-row.active {
  background: #061126;
  color: #fff;
  border-color: #061126;
}

.warehouse-row strong,
.mini-row strong {
  display: block;
  font-size: 18px;
  font-weight: 950;
}

.warehouse-row span,
.warehouse-row small,
.mini-row span {
  display: block;
  margin-top: 5px;
  color: #66758a;
  font-weight: 800;
}

.warehouse-row.active span,
.warehouse-row.active small {
  color: #cbd5e1;
}

em {
  border-radius: 999px;
  padding: 10px 14px;
  background: #dcfce7;
  color: #047857;
  font-style: normal;
  font-weight: 950;
  white-space: nowrap;
}

em.off {
  background: #e2e8f0;
  color: #64748b;
}

.small-btn {
  min-height: 46px;
  padding: 0 14px;
  background: #061126;
  color: #fff;
}

.small-btn.danger {
  background: #ffe4e6;
  color: #be123c;
}

.empty {
  margin-top: 18px;
  background: #f6f9fd;
  color: #64748b;
}

.mini-list {
  margin-top: 18px;
}

.mini-row {
  cursor: default;
}

.gate-row {
  grid-template-columns: minmax(220px, 1fr) auto auto auto;
}

@media (max-width: 1280px) {
  .workspace-grid,
  .structure-grid,
  .stats-grid,
  .inline-form,
  .inline-form.two {
    grid-template-columns: 1fr;
  }

  .warehouse-row,
  .mini-row,
  .gate-row {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 760px) {
  .hero-card {
    flex-direction: column;
    align-items: stretch;
  }

  .warehouse-row,
  .mini-row,
  .gate-row {
    grid-template-columns: 1fr;
  }
}
</style>
