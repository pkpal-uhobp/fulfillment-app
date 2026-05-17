<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="toolbar">
      <div>
        <p>Календарь</p>
        <h2>Выбранный склад</h2>
      </div>

      <BaseSelect
        v-model="selectedWarehouseId"
        class="warehouse-select"
        :options="warehouseOptions"
        placeholder="Выберите склад приёмки"
        searchable
        empty-text="Нет складов приёмки"
        @change="onWarehouseChange"
      />
    </div>

    <div class="calendar-panel">
      <div class="month-head">
        <button type="button" aria-label="Предыдущий месяц" @click="changeMonth(-1)">←</button>
        <div>
          <p>{{ monthTitle }}</p>
          <span>Выберите день, чтобы изменить лимит или закрыть дату</span>
        </div>
        <button type="button" aria-label="Следующий месяц" @click="changeMonth(1)">→</button>
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
          :class="{
            muted: !day.currentMonth,
            closed: day.data?.is_closed,
            filled: isFilled(day.data),
            selected: selectedDate === day.iso,
          }"
          @click="selectDay(day)"
        >
          <b>{{ day.date.getDate() }}</b>
          <span v-if="day.data?.is_closed">Закрыто</span>
          <span v-else-if="day.data">{{ day.data.current_orders ?? 0 }}/{{ day.data.max_orders ?? 0 }}</span>
          <span v-else>—</span>
        </button>
      </div>
    </div>

    <aside v-if="selectedDate" class="editor">
      <div class="editor-head">
        <div>
          <p>Выбранный день</p>
          <h3>{{ selectedDateLabel }}</h3>
          <span v-if="selectedWarehouse" class="muted-text">{{ warehouseName(selectedWarehouse.id, warehouses) }} · {{ selectedWarehouse.city }}</span>
          <span v-if="selectedDay?.block?.reason" class="closed-reason">Причина закрытия: {{ selectedDay.block.reason }}</span>
        </div>

        <div class="day-summary">
          <div>
            <small>Записано</small>
            <strong>{{ selectedBookedCount }}</strong>
          </div>
          <div>
            <small>Лимит</small>
            <strong>{{ Number(capacityForm.maxOrders || 0) }}</strong>
          </div>
          <div>
            <small>Свободно</small>
            <strong>{{ freeSlots }}</strong>
          </div>
        </div>
      </div>

      <div class="editor-grid">
        <label class="field-card">
          <span>Лимит заявок</span>
          <input
            v-model.number="capacityForm.maxOrders"
            type="number"
            min="0"
            inputmode="numeric"
            @focus="lockScroll"
            @blur="unlockScroll"
          />
          <small>Сколько заявок можно принять в этот день.</small>
        </label>

        <label class="field-card readonly-card">
          <span>Уже записано</span>
          <input :value="selectedBookedCount" type="number" readonly tabindex="-1" />
          <small>Считается автоматически по заявкам выбранного склада на эту дату.</small>
        </label>

        <label class="checkbox-card">
          <input v-model="capacityForm.isClosed" type="checkbox" />
          <span>Закрыть день для приёмки</span>
        </label>

        <label class="field-card reason-card">
          <span>Причина закрытия</span>
          <input v-model.trim="blockReason" placeholder="Например: инвентаризация" />
          <small>Заполняется, если дату нужно закрыть для клиентов.</small>
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
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { formatIsoDate, monthRange, unwrapList, warehouseName } from './logistUtils'

const warehouses = ref([])
const days = ref([])
const selectedWarehouseId = ref('')
const currentMonth = ref(new Date())
const selectedDate = ref('')
const error = ref('')
const success = ref('')
const blockReason = ref('')
const capacityForm = reactive({ maxOrders: 18, isClosed: false })

const receivingWarehouses = computed(() => warehouses.value.filter((warehouse) => ['receiving', 'both'].includes(warehouse.warehouse_type)))
const warehouseOptions = computed(() => receivingWarehouses.value.map((warehouse) => ({
  value: String(warehouse.id),
  label: `${warehouseName(warehouse.id, warehouses.value)} · ${warehouse.city}`,
  description: warehouse.address || 'Склад приёмки',
})))

const selectedWarehouse = computed(() => warehouses.value.find((warehouse) => String(warehouse.id) === String(selectedWarehouseId.value)))
const monthTitle = computed(() => new Intl.DateTimeFormat('ru-RU', { month: 'long', year: 'numeric' }).format(currentMonth.value))
const selectedDay = computed(() => days.value.find((day) => day.date === selectedDate.value))
const selectedBookedCount = computed(() => Number(selectedDay.value?.current_orders ?? 0))
const freeSlots = computed(() => Math.max(0, Number(capacityForm.maxOrders || 0) - selectedBookedCount.value))
const selectedDateLabel = computed(() => {
  if (!selectedDate.value) return ''
  return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }).format(new Date(`${selectedDate.value}T00:00:00`))
})

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
  return day && Number(day.max_orders || 0) > 0 && Number(day.current_orders || 0) >= Number(day.max_orders || 0)
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

async function onWarehouseChange() {
  selectedDate.value = ''
  await loadCalendar()
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
  capacityForm.maxOrders = Number(day?.max_orders ?? 18)
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
        current_orders: selectedBookedCount.value,
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

function lockScroll() {
  document.body.classList.add('no-scroll')
}

function unlockScroll() {
  document.body.classList.remove('no-scroll')
}

onBeforeUnmount(() => {
  document.body.classList.remove('no-scroll')
})

onMounted(async () => {
  await loadWarehouses()
  await loadCalendar()
})
</script>

<style scoped>
* { box-sizing: border-box; }
.page { display:grid; gap:22px; width:100%; min-width:0; }
.alert, .success { padding:16px 18px; border-radius:18px; font-weight:900; }
.alert { background:#fee2e2; color:#991b1b; }
.success { background:#d1fae5; color:#065f46; }
.toolbar, .calendar-panel, .editor { background:#fff; border-radius:34px; padding:26px; box-shadow:0 18px 42px rgba(7,16,31,.08); min-width:0; }
.toolbar { display:flex; justify-content:space-between; align-items:end; gap:20px; }
.warehouse-select { width:min(520px, 100%); }
p { margin:0 0 8px; color:#ff3f4d; text-transform:uppercase; letter-spacing:.22em; font-size:12px; font-weight:950; }
h2, h3 { margin:0; letter-spacing:-.04em; color:#07101f; }
h2 { font-size:clamp(24px, 3.2vw, 42px); line-height:1.05; }
h3 { font-size:clamp(28px, 3vw, 44px); line-height:1; }
.month-head { display:flex; align-items:center; justify-content:space-between; gap:20px; margin-bottom:18px; }
.month-head button, .actions button { border:0; cursor:pointer; height:58px; border-radius:20px; padding:0 24px; font-weight:950; color:white; background:#07101f; box-shadow:0 18px 34px rgba(7,16,31,.12); }
.month-head button { width:58px; padding:0; font-size:22px; }
.month-head span, .muted-text, .closed-reason, .field-card small { color:#6b7b91; font-weight:850; }
.weekdays, .days-grid { display:grid; grid-template-columns: repeat(7, minmax(0,1fr)); gap:10px; }
.weekdays span { color:#94a3b8; font-weight:950; text-align:center; }
.day-card { min-height:104px; border:1px solid #dbe3ef; border-radius:22px; background:#f6f8fb; display:grid; gap:6px; align-content:center; justify-items:center; cursor:pointer; color:#07101f; transition:.18s; }
.day-card:hover { transform:translateY(-2px); border-color:#ff3f4d; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.day-card b { font-size:24px; }
.day-card span { color:#6b7b91; font-weight:950; font-size:12px; }
.day-card.muted { opacity:.35; }
.day-card.closed { background:#fee2e2; border-color:#fecaca; }
.day-card.filled { background:#fff7ed; border-color:#fed7aa; }
.day-card.selected { background:#07101f; color:white; border-color:#07101f; box-shadow:0 24px 54px rgba(7,16,31,.22); }
.day-card.selected span { color:#dbeafe; }
.editor { overflow:hidden; }
.editor-head { display:grid; grid-template-columns:minmax(0,1fr) auto; gap:22px; align-items:start; padding-bottom:22px; border-bottom:1px solid #e7eef7; }
.editor-head .closed-reason { display:block; margin-top:10px; color:#ff3f4d; }
.day-summary { display:grid; grid-template-columns:repeat(3, minmax(100px, 1fr)); gap:10px; min-width:360px; }
.day-summary > div { padding:16px; border-radius:22px; background:#f5f8fc; border:1px solid #e1e9f3; }
.day-summary small { display:block; margin-bottom:6px; color:#8b9ab0; font-size:11px; font-weight:950; letter-spacing:.18em; text-transform:uppercase; }
.day-summary strong { display:block; color:#07101f; font-size:28px; line-height:1; }
.editor-grid { display:grid; grid-template-columns:repeat(2, minmax(0,1fr)); gap:18px; margin-top:22px; align-items:start; }
.field-card, .checkbox-card { min-width:0; border:1px solid #dbe3ef; border-radius:24px; background:#f8fbff; padding:20px 22px; }
.field-card { display:grid; gap:12px; }
.field-card > span, .checkbox-card span { color:#8b9ab0; font-size:12px; text-transform:uppercase; letter-spacing:.18em; font-weight:950; }
.field-card input { width:100%; max-width:100%; height:64px; border:1px solid #dbe3ef; border-radius:20px; padding:0 18px; background:white; color:#07101f; outline:none; font:inherit; font-size:18px; font-weight:950; }
.field-card input:focus { border-color:#ff3f4d; box-shadow:0 0 0 4px rgba(255,63,77,.12); }
.readonly-card input { background:#eef4fb; color:#07101f; cursor:default; }
.checkbox-card { grid-column:1 / -1; display:flex; align-items:center; gap:14px; min-height:76px; }
.checkbox-card input { width:24px; height:24px; accent-color:#ff3f4d; }
.reason-card { grid-column:1 / -1; }
.actions { margin-top:22px; display:flex; gap:14px; flex-wrap:wrap; align-items:center; }
.actions .danger { background:#ff3f4d; box-shadow:0 20px 40px rgba(255,63,77,.24); }
.actions .ghost { background:#edf2f7; color:#07101f; }
@media (max-width:1100px) {
  .toolbar { flex-direction:column; align-items:stretch; }
  .warehouse-select { width:100%; }
  .editor-head { grid-template-columns:1fr; }
  .day-summary { min-width:0; }
}
@media (max-width:900px) {
  .toolbar, .calendar-panel, .editor { border-radius:26px; padding:18px; }
  .weekdays, .days-grid { gap:6px; }
  .day-card { min-height:72px; border-radius:16px; }
  .day-card b { font-size:20px; }
  .editor-grid, .day-summary { grid-template-columns:1fr; }
  .actions { display:grid; grid-template-columns:1fr; }
  .actions button { width:100%; }
}
:global(body.no-scroll) { overflow: hidden; }
</style>

