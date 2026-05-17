<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="toolbar">
      <div>
        <p>Календарь</p>
        <h2>Лимиты и закрытые даты</h2>
      </div>
      <select v-model="selectedWarehouseId" @change="loadCalendar">
        <option value="">Выберите склад приёмки</option>
        <option v-for="warehouse in receivingWarehouses" :key="warehouse.id" :value="warehouse.id">
          {{ warehouseName(warehouse.id, warehouses) }} · {{ warehouse.city }}
        </option>
      </select>
    </div>

    <div class="calendar-panel">
      <div class="month-head">
        <button type="button" @click="changeMonth(-1)">←</button>
        <div>
          <p>{{ monthTitle }}</p>
          <span>Выберите день, чтобы изменить лимит или закрыть дату</span>
        </div>
        <button type="button" @click="changeMonth(1)">→</button>
      </div>

      <div class="weekdays">
        <span v-for="day in ['Пн','Вт','Ср','Чт','Пт','Сб','Вс']" :key="day">{{ day }}</span>
      </div>

      <div class="days-grid">
        <button
          v-for="day in monthCells"
          :key="day.key"
          type="button"
          class="day-card"
          :class="{ muted: !day.currentMonth, closed: day.data?.is_closed, filled: isFilled(day.data), selected: selectedDate === day.iso }"
          @click="selectDay(day)"
        >
          <b>{{ day.date.getDate() }}</b>
          <span v-if="day.data?.is_closed">Закрыто</span>
          <span v-else-if="day.data">{{ day.data.current_orders }}/{{ day.data.max_orders }}</span>
          <span v-else>—</span>
        </button>
      </div>
    </div>

    <aside class="editor" v-if="selectedDate">
      <div>
        <p>Дата</p>
        <h3>{{ selectedDate }}</h3>
        <span v-if="selectedDay?.block?.reason">Причина закрытия: {{ selectedDay.block.reason }}</span>
      </div>

      <div class="form-grid">
        <label>
          Лимит заявок
          <input v-model.number="capacityForm.maxOrders" type="number" min="0" />
        </label>
        <label>
          Уже записано
          <input v-model.number="capacityForm.currentOrders" type="number" min="0" />
        </label>
        <label class="checkbox">
          <input v-model="capacityForm.isClosed" type="checkbox" />
          Закрыть день для приёмки
        </label>
        <label>
          Причина закрытия
          <input v-model.trim="blockReason" placeholder="Например: инвентаризация" />
        </label>
      </div>

      <div class="actions">
        <button type="button" @click="saveCapacity">Сохранить лимит</button>
        <button type="button" class="danger" @click="blockDate">Закрыть дату</button>
        <button v-if="selectedDay?.block?.id" type="button" class="ghost" @click="unblockDate">Открыть дату</button>
      </div>
    </aside>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'
import { formatIsoDate, monthRange, unwrapList, warehouseName } from './logistUtils'

const warehouses = ref([])
const days = ref([])
const selectedWarehouseId = ref('')
const currentMonth = ref(new Date())
const selectedDate = ref('')
const error = ref('')
const success = ref('')
const blockReason = ref('')
const capacityForm = reactive({ maxOrders: 18, currentOrders: 0, isClosed: false })

const receivingWarehouses = computed(() => warehouses.value.filter((warehouse) => ['receiving', 'both'].includes(warehouse.warehouse_type)))
const monthTitle = computed(() => new Intl.DateTimeFormat('ru-RU', { month: 'long', year: 'numeric' }).format(currentMonth.value))
const selectedDay = computed(() => days.value.find((day) => day.date === selectedDate.value))

const monthCells = computed(() => {
  const year = currentMonth.value.getFullYear()
  const month = currentMonth.value.getMonth()
  const first = new Date(year, month, 1)
  const startOffset = (first.getDay() + 6) % 7
  const start = new Date(year, month, 1 - startOffset)
  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(start)
    date.setDate(start.getDate() + index)
    const iso = formatIsoDate(date)
    return {
      key: iso,
      iso,
      date,
      currentMonth: date.getMonth() === month,
      data: days.value.find((day) => day.date === iso),
    }
  })
})

function isFilled(day) {
  return day && day.max_orders > 0 && day.current_orders >= day.max_orders
}

async function loadWarehouses() {
  const payload = await apiFetch('/warehouses')
  warehouses.value = unwrapList(payload, 'warehouses')
  if (!selectedWarehouseId.value && receivingWarehouses.value.length) {
    selectedWarehouseId.value = String(receivingWarehouses.value[0].id)
  }
}

async function loadCalendar() {
  if (!selectedWarehouseId.value) return
  error.value = ''
  success.value = ''
  const range = monthRange(currentMonth.value)
  try {
    const payload = await apiFetch(`/pickup-calendar?warehouse_id=${selectedWarehouseId.value}&date_from=${range.from}&date_to=${range.to}`, { auth: true })
    days.value = unwrapList(payload, 'days')
    if (selectedDate.value) syncForm()
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить календарь'
  }
}

function changeMonth(delta) {
  currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() + delta, 1)
  selectedDate.value = ''
  loadCalendar()
}

function selectDay(day) {
  if (!day.currentMonth) return
  selectedDate.value = day.iso
  syncForm()
}

function syncForm() {
  const day = selectedDay.value
  capacityForm.maxOrders = day?.max_orders ?? 18
  capacityForm.currentOrders = day?.current_orders ?? 0
  capacityForm.isClosed = Boolean(day?.is_closed)
  blockReason.value = day?.block?.reason || ''
}

async function saveCapacity() {
  if (!selectedWarehouseId.value || !selectedDate.value) return
  error.value = ''
  success.value = ''
  try {
    await apiFetch('/pickup-calendar/capacity', {
      auth: true,
      method: 'PATCH',
      body: {
        warehouse_id: Number(selectedWarehouseId.value),
        pickup_date: selectedDate.value,
        max_orders: Number(capacityForm.maxOrders || 0),
        current_orders: Number(capacityForm.currentOrders || 0),
        is_closed: Boolean(capacityForm.isClosed),
      },
    })
    success.value = 'Лимит календаря сохранён'
    await loadCalendar()
  } catch (err) {
    error.value = err.message || 'Не удалось сохранить лимит'
  }
}

async function blockDate() {
  if (!selectedWarehouseId.value || !selectedDate.value) return
  error.value = ''
  success.value = ''
  try {
    await apiFetch('/pickup-calendar/blocks', {
      auth: true,
      method: 'POST',
      body: {
        warehouse_id: Number(selectedWarehouseId.value),
        blocked_date: selectedDate.value,
        reason: blockReason.value || undefined,
      },
    })
    capacityForm.isClosed = true
    success.value = 'Дата закрыта для приёмки'
    await loadCalendar()
  } catch (err) {
    error.value = err.message || 'Не удалось закрыть дату'
  }
}

async function unblockDate() {
  const blockId = selectedDay.value?.block?.id
  if (!blockId) return
  error.value = ''
  success.value = ''
  try {
    await apiFetch(`/pickup-calendar/blocks/${blockId}`, { auth: true, method: 'DELETE' })
    success.value = 'Дата снова доступна'
    await loadCalendar()
  } catch (err) {
    error.value = err.message || 'Не удалось открыть дату'
  }
}

onMounted(async () => {
  await loadWarehouses()
  await loadCalendar()
})
</script>

<style scoped>
.page { display:grid; gap:20px; }
.alert, .success { padding:16px 18px; border-radius:18px; font-weight:900; }
.alert { background:#fee2e2; color:#991b1b; }
.success { background:#d1fae5; color:#065f46; }
.toolbar, .calendar-panel, .editor { background:white; border-radius:34px; padding:26px; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.toolbar { display:flex; justify-content:space-between; align-items:end; gap:20px; }
p { margin:0 0 8px; color:#ff3f4d; text-transform:uppercase; letter-spacing:.22em; font-size:12px; font-weight:900; }
h2, h3 { margin:0; letter-spacing:-.04em; }
h2 { font-size:34px; }
h3 { font-size:30px; }
select, input { height:52px; border:1px solid #dbe3ef; border-radius:18px; padding:0 16px; background:#f6f8fb; color:#07101f; font-weight:800; }
.month-head { display:flex; align-items:center; justify-content:space-between; gap:20px; margin-bottom:18px; }
.month-head button, .actions button { border:0; cursor:pointer; height:48px; border-radius:16px; padding:0 18px; font-weight:900; color:white; background:#07101f; }
.month-head span, .editor span { color:#6b7b91; font-weight:800; }
.weekdays, .days-grid { display:grid; grid-template-columns: repeat(7, minmax(0,1fr)); gap:10px; }
.weekdays span { color:#94a3b8; font-weight:900; text-align:center; }
.day-card { min-height:104px; border:1px solid #dbe3ef; border-radius:22px; background:#f6f8fb; display:grid; gap:6px; align-content:center; justify-items:center; cursor:pointer; color:#07101f; }
.day-card b { font-size:24px; }
.day-card span { color:#6b7b91; font-weight:900; font-size:12px; }
.day-card.muted { opacity:.35; }
.day-card.closed { background:#fee2e2; border-color:#fecaca; }
.day-card.filled { background:#fff7ed; border-color:#fed7aa; }
.day-card.selected { background:#07101f; color:white; border-color:#07101f; }
.day-card.selected span { color:#dbeafe; }
.form-grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:14px; margin-top:18px; }
label { display:grid; gap:8px; color:#8b9ab0; font-size:12px; text-transform:uppercase; letter-spacing:.14em; font-weight:900; }
.checkbox { grid-column:1 / -1; display:flex; flex-direction:row; align-items:center; }
.checkbox input { width:20px; height:20px; }
.actions { margin-top:18px; display:flex; gap:12px; flex-wrap:wrap; }
.actions .danger { background:#ff3f4d; }
.actions .ghost { background:#edf2f7; color:#07101f; }
@media (max-width:900px) { .toolbar { flex-direction:column; align-items:stretch; } .weekdays, .days-grid { gap:6px; } .day-card { min-height:72px; border-radius:16px; } .form-grid { grid-template-columns:1fr; } }
</style>
