<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="create-panel">
      <div>
        <p>Отгрузки</p>
        <h2>Создать отгрузку</h2>
      </div>

      <div class="create-grid">
        <label>
          <span>Склад назначения</span>
          <BaseSelect
            v-model="form.destinationWarehouseId"
            :options="destinationWarehouseOptions"
            placeholder="Выберите склад"
          />
        </label>

        <label>
          <span>Гейт</span>
          <BaseSelect v-model="form.gateId" :options="createGateOptions" placeholder="Выберите гейт" />
        </label>

        <label ref="datePickerRef" class="date-picker-field">
          <span>Дата и время отправки</span>

          <button type="button" class="date-button" @click.stop="toggleCalendar">
            <span :class="{ muted: !form.plannedDate }">{{ plannedDateLabel }}</span>
            <b>{{ form.plannedTime || '12:00' }}</b>
          </button>

          <div v-if="calendarOpen" class="calendar-popover" @click.stop>
            <div class="calendar-head">
              <button type="button" aria-label="Предыдущий месяц" @click="shiftMonth(-1)">‹</button>
              <strong>{{ monthLabel }}</strong>
              <button type="button" aria-label="Следующий месяц" @click="shiftMonth(1)">›</button>
            </div>

            <div class="calendar-weekdays">
              <span v-for="day in weekdays" :key="day">{{ day }}</span>
            </div>

            <div class="calendar-grid">
              <button
                v-for="day in calendarDays"
                :key="day.key"
                type="button"
                :class="{
                  outside: day.outside,
                  today: day.isToday,
                  selected: day.value === form.plannedDate,
                }"
                @click="pickDate(day.value)"
              >
                {{ day.label }}
              </button>
            </div>

            <div class="time-picker">
              <div>
                <span>Часы</span>
                <div class="time-list">
                  <button
                    v-for="hour in hourOptions"
                    :key="hour"
                    type="button"
                    :class="{ active: selectedHour === hour }"
                    @click="setHour(hour)"
                  >
                    {{ hour }}
                  </button>
                </div>
              </div>

              <div>
                <span>Минуты</span>
                <div class="time-list">
                  <button
                    v-for="minute in minuteOptions"
                    :key="minute"
                    type="button"
                    :class="{ active: selectedMinute === minute }"
                    @click="setMinute(minute)"
                  >
                    {{ minute }}
                  </button>
                </div>
              </div>
            </div>

            <div class="time-panel">
              <strong>{{ form.plannedTime || '12:00' }}</strong>
              <button type="button" @click="setNow">Сегодня</button>
              <button type="button" class="done-btn" @click="closeCalendar">Готово</button>
            </div>
          </div>
        </label>

        <button class="create-btn" type="button" @click="createShipment">Создать</button>
      </div>
    </div>

    <div class="toolbar panel">
      <div>
        <p>Список</p>
        <h2>Партии к отправке</h2>
      </div>

      <div class="filters filters--shipments">
        <BaseSelect
          v-model="filters.status"
          :options="shipmentStatusOptions"
          placeholder="Все статусы"
          @change="loadShipments"
        />

        <BaseSelect
          v-model="filters.destinationWarehouseId"
          :options="destinationWarehouseFilterOptions"
          placeholder="Все склады"
        />

        <BaseSelect
          v-model="filters.gateId"
          :options="gateFilterOptions"
          placeholder="Все гейты"
        />

        <label class="date-filter">
          <span>С даты</span>
          <input v-model="filters.fromDate" type="date" />
        </label>

        <label class="date-filter">
          <span>По дату</span>
          <input v-model="filters.toDate" type="date" />
        </label>

        <button type="button" @click="loadShipments">Обновить</button>
      </div>
    </div>

    <div class="shipment-grid">
      <article v-for="shipment in filteredShipments" :key="shipment.id" class="shipment-card">
        <div class="card-top">
          <div>
            <span>Отгрузка #{{ shipment.id }}</span>
            <h3>{{ labelFromMap(shipmentStatusLabels, shipment.status) }}</h3>
          </div>

          <RouterLink :to="`/logist/shipments/${shipment.id}`">Открыть</RouterLink>
        </div>

        <div class="meta-grid">
          <div><small>Назначение</small><b>{{ warehouseName(shipment.destination_warehouse_id, warehouses) }}</b></div>
          <div><small>Гейт</small><b>{{ gateName(shipment.gate_id, gates) }}</b></div>
          <div><small>План отправки</small><b>{{ formatDateTime(shipment.planned_departure_at) }}</b></div>
          <div><small>Мест</small><b>{{ shipment.items?.length || 0 }}</b></div>
        </div>

        <div class="action-box">
          <BaseSelect
            v-model="forms[shipment.id].status"
            :options="shipmentStatusOptions.filter((item) => item.value)"
            placeholder="Выберите статус"
          />

          <button type="button" :disabled="loadingId === shipment.id" @click="updateStatus(shipment)">
            Обновить статус
          </button>
        </div>
      </article>
    </div>

    <div v-if="!filteredShipments.length && !error" class="empty">Отгрузки не найдены</div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { apiFetch } from '@/shared/api/http'
import {
  formatDateTime,
  gateName,
  labelFromMap,
  shipmentStatusLabels,
  shipmentStatusOptions,
  unwrapList,
  unwrapOne,
  warehouseName,
} from './logistUtils'

const shipments = ref([])
const warehouses = ref([])
const gates = ref([])
const forms = reactive({})
const filters = reactive({
  status: '',
  destinationWarehouseId: '',
  gateId: '',
  fromDate: '',
  toDate: '',
})
const form = reactive({ destinationWarehouseId: '', gateId: '', plannedDate: '', plannedTime: '12:00' })

const error = ref('')
const success = ref('')
const loadingId = ref(null)

const datePickerRef = ref(null)
const calendarOpen = ref(false)
const shownMonth = ref(startOfMonth(new Date()))

const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
const hourOptions = Array.from({ length: 24 }, (_, index) => String(index).padStart(2, '0'))
const minuteOptions = Array.from({ length: 12 }, (_, index) => String(index * 5).padStart(2, '0'))

const selectedHour = computed(() => (form.plannedTime || '12:00').split(':')[0])
const selectedMinute = computed(() => (form.plannedTime || '12:00').split(':')[1])

const destinationWarehouses = computed(() =>
  warehouses.value.filter((warehouse) => ['destination', 'both'].includes(warehouse.warehouse_type)),
)

const destinationWarehouseOptions = computed(() =>
  destinationWarehouses.value.map((warehouse) => ({
    value: String(warehouse.id),
    label: `${warehouseName(warehouse.id, warehouses.value)} · ${warehouse.city || '—'}`,
    description: warehouse.address,
  })),
)

const destinationWarehouseFilterOptions = computed(() => [
  { value: '', label: 'Все склады' },
  ...destinationWarehouses.value.map((warehouse) => ({
    value: String(warehouse.id),
    label: warehouseName(warehouse.id, warehouses.value),
    description: warehouse.address || warehouse.city || '',
  })),
])

const createGateOptions = computed(() => {
  const selectedWarehouseId = Number(form.destinationWarehouseId)

  return gates.value
    .filter((gate) => !selectedWarehouseId || Number(gate.warehouse_id) === selectedWarehouseId)
    .map((gate) => ({
      value: String(gate.id),
      label: gate.name,
      description: gate.warehouse_name || warehouseName(gate.warehouse_id, warehouses.value),
    }))
})

const gateFilterOptions = computed(() => [
  { value: '', label: 'Все гейты' },
  ...gates.value
    .filter((gate) => !filters.destinationWarehouseId || Number(gate.warehouse_id) === Number(filters.destinationWarehouseId))
    .map((gate) => ({
      value: String(gate.id),
      label: gate.name,
      description: gate.warehouse_name || warehouseName(gate.warehouse_id, warehouses.value),
    })),
])

const filteredShipments = computed(() => {
  const from = filters.fromDate ? new Date(`${filters.fromDate}T00:00:00`).getTime() : null
  const to = filters.toDate ? new Date(`${filters.toDate}T23:59:59`).getTime() : null

  return shipments.value.filter((shipment) => {
    if (filters.destinationWarehouseId && Number(shipment.destination_warehouse_id) !== Number(filters.destinationWarehouseId)) {
      return false
    }

    if (filters.gateId && Number(shipment.gate_id) !== Number(filters.gateId)) {
      return false
    }

    const planned = shipment.planned_departure_at ? new Date(shipment.planned_departure_at).getTime() : null

    if (from && (!planned || planned < from)) return false
    if (to && (!planned || planned > to)) return false

    return true
  })
})

const monthLabel = computed(() =>
  new Intl.DateTimeFormat('ru-RU', { month: 'long', year: 'numeric' }).format(shownMonth.value),
)

const plannedDateLabel = computed(() => {
  if (!form.plannedDate) return 'Выберите дату'

  const date = new Date(`${form.plannedDate}T00:00:00`)
  return new Intl.DateTimeFormat('ru-RU', { day: '2-digit', month: 'long', year: 'numeric' }).format(date)
})

const calendarDays = computed(() => {
  const first = startOfMonth(shownMonth.value)
  const start = new Date(first)
  const shift = (first.getDay() + 6) % 7
  start.setDate(first.getDate() - shift)

  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(start)
    date.setDate(start.getDate() + index)
    const value = toInputDate(date)

    return {
      key: `${value}-${index}`,
      value,
      label: date.getDate(),
      outside: date.getMonth() !== shownMonth.value.getMonth(),
      isToday: value === toInputDate(new Date()),
    }
  })
})

function startOfMonth(date) {
  return new Date(date.getFullYear(), date.getMonth(), 1)
}

function toInputDate(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function ensureForm(shipment) {
  if (!forms[shipment.id]) forms[shipment.id] = { status: shipment.status || 'planned' }
}

function setTime(hour = selectedHour.value, minute = selectedMinute.value) {
  form.plannedTime = `${hour}:${minute}`
}

function setHour(hour) {
  setTime(hour, selectedMinute.value)
}

function setMinute(minute) {
  setTime(selectedHour.value, minute)
}

function toggleCalendar() {
  calendarOpen.value = !calendarOpen.value

  if (form.plannedDate) {
    const selected = new Date(`${form.plannedDate}T00:00:00`)
    shownMonth.value = startOfMonth(selected)
  }
}

function closeCalendar() {
  calendarOpen.value = false
}

function onDocumentClick(event) {
  if (!datePickerRef.value?.contains(event.target)) closeCalendar()
}

function shiftMonth(delta) {
  shownMonth.value = new Date(shownMonth.value.getFullYear(), shownMonth.value.getMonth() + delta, 1)
}

function pickDate(value) {
  form.plannedDate = value
}

function setNow() {
  const now = new Date()
  form.plannedDate = toInputDate(now)
  form.plannedTime = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
  shownMonth.value = startOfMonth(now)
}

function plannedIso() {
  if (!form.plannedDate || !form.plannedTime) return ''

  const date = new Date(`${form.plannedDate}T${form.plannedTime}:00`)
  return Number.isNaN(date.getTime()) ? '' : date.toISOString()
}

async function loadCatalogs() {
  const [warehousesPayload, gatesPayload] = await Promise.all([
    apiFetch('/warehouses'),
    apiFetch('/gates', { auth: true }),
  ])

  warehouses.value = unwrapList(warehousesPayload, 'warehouses')
  gates.value = unwrapList(gatesPayload, 'gates')
}

async function loadShipments() {
  error.value = ''
  success.value = ''

  const params = new URLSearchParams()
  if (filters.status) params.set('status', filters.status)

  try {
    const payload = await apiFetch(`/shipments${params.toString() ? `?${params}` : ''}`, { auth: true })
    shipments.value = unwrapList(payload, 'shipments')
    shipments.value.forEach(ensureForm)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить отгрузки'
  }
}

async function createShipment() {
  error.value = ''
  success.value = ''

  const departureAt = plannedIso()

  if (!form.destinationWarehouseId || !form.gateId || !departureAt) {
    error.value = 'Выберите склад, гейт, дату и время отправки'
    return
  }

  try {
    await apiFetch('/shipments', {
      auth: true,
      method: 'POST',
      body: {
        destination_warehouse_id: Number(form.destinationWarehouseId),
        gate_id: Number(form.gateId),
        planned_departure_at: departureAt,
      },
    })

    success.value = 'Отгрузка создана'
    form.destinationWarehouseId = ''
    form.gateId = ''
    form.plannedDate = ''
    form.plannedTime = '12:00'
    closeCalendar()
    await loadShipments()
  } catch (err) {
    error.value = err.message || 'Не удалось создать отгрузку'
  }
}

async function updateStatus(shipment) {
  error.value = ''
  success.value = ''
  loadingId.value = shipment.id

  try {
    const payload = await apiFetch(`/shipments/${shipment.id}/status`, {
      auth: true,
      method: 'PATCH',
      body: { status: forms[shipment.id].status },
    })

    const updated = unwrapOne(payload, 'shipment') || { ...shipment, status: forms[shipment.id].status }
    const index = shipments.value.findIndex((item) => item.id === shipment.id)

    if (index !== -1) shipments.value[index] = updated

    ensureForm(updated)
    success.value = `Статус отгрузки #${shipment.id} обновлён`
  } catch (err) {
    error.value = err.message || 'Не удалось обновить статус отгрузки'
  } finally {
    loadingId.value = null
  }
}

watch(() => form.destinationWarehouseId, () => {
  if (!form.gateId) return
  const selectedGate = gates.value.find((gate) => Number(gate.id) === Number(form.gateId))

  if (selectedGate && Number(selectedGate.warehouse_id) !== Number(form.destinationWarehouseId)) {
    form.gateId = ''
  }
})

watch(() => filters.destinationWarehouseId, () => {
  if (!filters.gateId) return
  const selectedGate = gates.value.find((gate) => Number(gate.id) === Number(filters.gateId))

  if (selectedGate && Number(selectedGate.warehouse_id) !== Number(filters.destinationWarehouseId)) {
    filters.gateId = ''
  }
})

onMounted(async () => {
  document.addEventListener('click', onDocumentClick)
  await Promise.all([loadCatalogs(), loadShipments()])
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick)
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
.create-panel,
.shipment-card,
.empty {
  background: white;
  border-radius: 32px;
  padding: 24px;
  box-shadow: 0 18px 42px rgba(7, 16, 31, .08);
}

.create-panel {
  display: grid;
  gap: 18px;
}

p {
  margin: 0 0 8px;
  color: #ff3f4d;
  letter-spacing: .22em;
  text-transform: uppercase;
  font-size: 12px;
  font-weight: 950;
}

h2,
h3 {
  margin: 0;
  letter-spacing: -.04em;
}

.toolbar {
  display: grid;
  gap: 20px;
}

.toolbar h2,
.create-panel h2 {
  font-size: 34px;
}

.create-grid {
  display: grid;
  grid-template-columns: 1.35fr 1fr 1.32fr auto;
  gap: 14px;
  align-items: end;
}

.create-grid label,
.search-field,
.date-filter {
  display: grid;
  gap: 8px;
}

.create-grid label > span,
.search-field span,
.date-filter span,
.time-picker span {
  color: #94a3b8;
  letter-spacing: .22em;
  text-transform: uppercase;
  font-size: 12px;
  font-weight: 950;
}

.date-picker-field {
  position: relative;
  min-width: 0;
}

.date-button,
input {
  min-height: 58px;
  border: 1px solid #dbe3ef;
  border-radius: 18px;
  padding: 0 18px;
  background: #f6f8fb;
  color: #07101f;
  font-weight: 950;
  font-family: inherit;
}

.date-button {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  text-align: left;
  cursor: pointer;
}

.date-button .muted {
  color: #94a3b8;
}

.date-button b {
  padding: 8px 10px;
  border-radius: 12px;
  background: #eaf1fb;
}

.calendar-popover {
  position: absolute;
  z-index: 80;
  top: calc(100% + 12px);
  right: 0;
  width: min(460px, 92vw);
  border-radius: 30px;
  border: 1px solid #dbe3ef;
  background: #fff;
  box-shadow: 0 28px 70px rgba(7, 16, 31, .20);
  padding: 18px;
}

.calendar-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.calendar-head strong {
  font-size: 20px;
  font-weight: 950;
  text-transform: capitalize;
}

.calendar-head button,
.time-panel button {
  min-height: 48px;
  width: 48px;
  border-radius: 16px;
  border: 0;
  background: #eef3f9;
  color: #07101f;
  font-weight: 950;
  cursor: pointer;
}

.calendar-weekdays,
.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 6px;
}

.calendar-weekdays span {
  text-align: center;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 950;
  text-transform: uppercase;
}

.calendar-grid {
  margin-top: 8px;
}

.calendar-grid button {
  min-height: 44px;
  border: 0;
  border-radius: 16px;
  background: #f6f8fb;
  color: #07101f;
  font-weight: 950;
  cursor: pointer;
}

.calendar-grid button.outside {
  opacity: .45;
}

.calendar-grid button.today {
  box-shadow: inset 0 0 0 2px #dbeafe;
}

.calendar-grid button.selected {
  background: #ff3f4d;
  color: #fff;
  box-shadow: 0 12px 24px rgba(255, 63, 77, .22);
}

.time-picker {
  margin-top: 16px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.time-picker > div {
  display: grid;
  gap: 8px;
}

.time-list {
  max-height: 132px;
  overflow-y: auto;
  overscroll-behavior: contain;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 6px;
  padding-right: 4px;
}

.time-list button {
  min-height: 40px;
  border: 0;
  border-radius: 14px;
  background: #f6f8fb;
  color: #07101f;
  font-weight: 950;
  cursor: pointer;
}

.time-list button.active {
  background: #ff3f4d;
  color: #fff;
}

.time-panel {
  margin-top: 14px;
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 10px;
  align-items: center;
}

.time-panel strong {
  min-height: 56px;
  border-radius: 18px;
  background: #f6f8fb;
  display: flex;
  align-items: center;
  padding: 0 18px;
  font-size: 20px;
  font-weight: 950;
}

.time-panel button {
  width: auto;
  padding: 0 16px;
  background: #07101f;
  color: #fff;
}

.time-panel .done-btn {
  background: #ff3f4d;
}

button,
.card-top a {
  min-height: 58px;
  border: 0;
  border-radius: 18px;
  padding: 0 20px;
  background: #ff3f4d;
  color: white;
  font-weight: 950;
  cursor: pointer;
  font-family: inherit;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.create-btn,
.filters button {
  background: #07101f;
}

.filters {
  display: grid;
  gap: 12px;
  align-items: end;
}

.filters--shipments {
  grid-template-columns: minmax(180px, 230px) minmax(190px, 260px) minmax(170px, 220px) minmax(220px, 1fr) minmax(150px, 180px) minmax(150px, 180px) auto;
}

.search-field input,
.date-filter input {
  width: 100%;
  box-sizing: border-box;
}

.shipment-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.shipment-card {
  display: grid;
  gap: 18px;
}

.card-top {
  display: flex;
  justify-content: space-between;
  gap: 14px;
}

.card-top span {
  color: #94a3b8;
  letter-spacing: .22em;
  text-transform: uppercase;
  font-size: 12px;
  font-weight: 950;
}

.card-top h3 {
  margin-top: 8px;
  font-size: 28px;
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

.action-box {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
}

.action-box button {
  min-height: 58px;
}

.empty {
  text-align: center;
  color: #64748b;
  font-weight: 950;
}

@media (max-width: 1440px) {
  .filters--shipments {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 1180px) {
  .shipment-grid {
    grid-template-columns: 1fr;
  }

  .create-grid,
  .filters--shipments {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 680px) {
  .meta-grid,
  .action-box,
  .time-picker,
  .time-panel {
    grid-template-columns: 1fr;
  }

  .calendar-popover {
    left: 0;
    right: auto;
  }
}
</style>
