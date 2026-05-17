<template>
  <div ref="root" class="date-picker">
    <label v-if="label" class="dp-label">{{ label }}</label>

    <button
      ref="button"
      type="button"
      class="dp-trigger"
      :class="{ open, error }"
      :disabled="disabled"
      @click="toggle"
    >
      <span :class="{ placeholder: !modelValue }">
        {{ modelValue ? human(modelValue) : placeholder }}
      </span>

      <span class="dp-icon" aria-hidden="true">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M7 3v3M17 3v3M4.5 9.2h15" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          <path d="M6.7 5h10.6c1.2 0 2.2 1 2.2 2.2v10.1c0 1.2-1 2.2-2.2 2.2H6.7c-1.2 0-2.2-1-2.2-2.2V7.2C4.5 6 5.5 5 6.7 5Z" stroke="currentColor" stroke-width="2" />
          <path d="M8 13h.01M12 13h.01M16 13h.01M8 16.5h.01M12 16.5h.01" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
        </svg>
      </span>
    </button>

    <p v-if="error" class="dp-error">{{ error }}</p>

    <Teleport to="body">
      <div v-show="open" ref="panel" class="dp-panel" :style="style" @click.stop>
        <header>
          <button type="button" class="nav-btn" aria-label="Предыдущий месяц" @click="prev">‹</button>
          <strong>{{ months[month] }} {{ year }} г.</strong>
          <button type="button" class="nav-btn" aria-label="Следующий месяц" @click="next">›</button>
        </header>

        <div class="week">
          <span v-for="d in week" :key="d">{{ d }}</span>
        </div>

        <div class="grid">
          <button
            v-for="day in calendarDays"
            :key="day.key"
            type="button"
            class="day"
            :class="{
              muted: !day.current,
              today: day.today,
              selected: day.iso === modelValue,
              disabled: day.disabled,
              closed: day.closed,
              available: day.available,
            }"
            :disabled="day.disabled"
            :title="day.title"
            @click="select(day)"
          >
            <b>{{ day.date.getDate() }}</b>
            <small v-if="day.meta && day.current">{{ day.closed ? 'закрыто' : free(day.meta) }}</small>
          </button>
        </div>

        <footer>
          <span class="legend-badge available-legend">доступно</span>
          <span class="legend-badge closed-legend">закрыто / лимит</span>
        </footer>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: '' },
  placeholder: { type: String, default: 'Выберите дату' },
  error: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  min: { type: String, default: '' },
  max: { type: String, default: '' },
  availabilityDays: { type: Array, default: () => [] },
  days: { type: Array, default: () => [] },
})

const emit = defineEmits(['update:modelValue', 'change'])

const root = ref(null)
const button = ref(null)
const panel = ref(null)
const open = ref(false)
const style = ref({})

const now = new Date()
const initial = props.modelValue ? parse(props.modelValue) : now
const year = ref(initial.getFullYear())
const month = ref(initial.getMonth())

const months = ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь']
const week = ['ПН', 'ВТ', 'СР', 'ЧТ', 'ПТ', 'СБ', 'ВС']

const availabilityMap = computed(() => {
  const m = new Map()
  const source = props.availabilityDays.length ? props.availabilityDays : props.days

  source.forEach((item) => {
    const iso = String(item.pickup_date || item.date || item.blocked_date || '').slice(0, 10)
    if (iso) m.set(iso, item)
  })

  return m
})

const calendarDays = computed(() => {
  const first = new Date(year.value, month.value, 1)
  const shift = (first.getDay() + 6) % 7
  const start = new Date(year.value, month.value, 1 - shift)

  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(start)
    date.setDate(start.getDate() + index)

    const iso = toIso(date)
    const meta = availabilityMap.value.get(iso)
    const closed = isClosed(meta)
    const disabled = isDisabled(iso, meta)
    const available = Boolean(meta && !closed && !disabled)

    return {
      key: `${iso}-${index}`,
      date,
      iso,
      meta,
      closed,
      available,
      disabled,
      current: date.getMonth() === month.value,
      today: iso === toIso(now),
      title: title(meta, disabled),
    }
  })
})

function parse(value) {
  const [y, m, d] = String(value).slice(0, 10).split('-').map(Number)
  return new Date(y || now.getFullYear(), (m || 1) - 1, d || 1)
}

function toIso(date) {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function human(value) {
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(parse(value))
}

function free(day) {
  const max = Number(day.max_orders || day.capacity || 0)
  const current = Number(day.current_orders || day.orders_count || 0)
  return max ? `${Math.max(0, max - current)}/${max}` : 'свободно'
}

function isClosed(day) {
  if (!day) return false
  if (day.is_closed || day.closed || day.is_blocked) return true

  const max = Number(day.max_orders || day.capacity || 0)
  const current = Number(day.current_orders || day.orders_count || 0)
  return max > 0 && current >= max
}

function isDisabled(iso, day) {
  if (props.min && iso < props.min) return true
  if (props.max && iso > props.max) return true
  return isClosed(day)
}

function title(day, disabled) {
  if (!day) return disabled ? 'Дата недоступна' : 'Дата доступна'
  if (day.is_closed || day.closed || day.is_blocked) return day.reason || 'Дата закрыта логистом'

  const max = Number(day.max_orders || day.capacity || 0)
  const current = Number(day.current_orders || day.orders_count || 0)
  if (max > 0 && current >= max) return 'Лимит заявок исчерпан'
  return max ? `Доступно ${Math.max(0, max - current)} из ${max}` : 'Дата доступна'
}

function position() {
  const r = button.value?.getBoundingClientRect()
  if (!r) return

  const width = Math.max(560, r.width)
  const height = 620
  const up = r.bottom + height + 12 > window.innerHeight
  const left = Math.min(Math.max(12, r.left), window.innerWidth - Math.min(width, window.innerWidth - 24) - 12)

  style.value = {
    left: `${left}px`,
    top: `${up ? Math.max(12, r.top - height - 10) : r.bottom + 10}px`,
    width: `${Math.min(width, window.innerWidth - 24)}px`,
  }
}

async function show() {
  if (props.disabled || open.value) return
  open.value = true
  await nextTick()
  position()
  window.addEventListener('resize', position)
  window.addEventListener('scroll', position, true)
  window.addEventListener('click', outside, true)
  window.addEventListener('keydown', esc)
}

function hide() {
  if (!open.value) return
  open.value = false
  window.removeEventListener('resize', position)
  window.removeEventListener('scroll', position, true)
  window.removeEventListener('click', outside, true)
  window.removeEventListener('keydown', esc)
}

function toggle() {
  open.value ? hide() : show()
}

function outside(event) {
  if (!root.value?.contains(event.target) && !panel.value?.contains(event.target)) hide()
}

function esc(event) {
  if (event.key === 'Escape') hide()
}

function prev() {
  if (month.value === 0) {
    month.value = 11
    year.value -= 1
  } else {
    month.value -= 1
  }
}

function next() {
  if (month.value === 11) {
    month.value = 0
    year.value += 1
  } else {
    month.value += 1
  }
}

function select(day) {
  if (day.disabled) return
  emit('update:modelValue', day.iso)
  emit('change', day)
  hide()
}

watch(() => props.modelValue, (value) => {
  if (!value) return
  const date = parse(value)
  year.value = date.getFullYear()
  month.value = date.getMonth()
})

onBeforeUnmount(hide)
</script>

<style scoped>
.date-picker {
  position: relative;
}

.dp-label {
  display: block;
  margin: 0 0 10px;
  color: #94a3b8;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .28em;
  text-transform: uppercase;
}

.dp-trigger {
  width: 100%;
  min-height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding: 0 18px 0 22px;
  border: 1px solid #dce5f0;
  border-radius: 22px;
  background: #f7faff;
  color: #07101f;
  font: inherit;
  font-weight: 950;
  cursor: pointer;
  text-align: left;
  transition: border-color .18s ease, background .18s ease, box-shadow .18s ease;
}

.dp-trigger:hover,
.dp-trigger.open {
  border-color: #ff3f4c;
  background: #fff;
  box-shadow: 0 18px 42px rgba(255, 63, 76, .14);
}

.dp-trigger.error {
  border-color: #ff3f4c;
  background: #fff5f6;
}

.placeholder {
  color: #8da0ba;
}

.dp-icon {
  width: 42px;
  height: 42px;
  flex: 0 0 auto;
  display: grid;
  place-items: center;
  border-radius: 14px;
  background: #eef4ff;
  color: #ff3f4c;
}

.dp-icon svg {
  width: 22px;
  height: 22px;
}

.dp-error {
  margin: 9px 0 0;
  color: #ff3f4c;
  font-weight: 850;
  font-size: 13px;
}

.dp-panel {
  position: fixed;
  z-index: 9999;
  padding: 28px;
  border: 1px solid rgba(220, 229, 240, .95);
  border-radius: 34px;
  background: rgba(255, 255, 255, .98);
  box-shadow: 0 38px 90px rgba(7, 16, 31, .24);
  backdrop-filter: blur(18px);
}

.dp-panel header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 22px;
}

.nav-btn {
  width: 58px;
  height: 58px;
  border: 0;
  border-radius: 20px;
  background: #f1f5fb;
  color: #07101f;
  font-size: 30px;
  line-height: 1;
  font-weight: 950;
  cursor: pointer;
}

.nav-btn:hover {
  background: #e8eef7;
}

.dp-panel header strong {
  font-size: 26px;
  font-weight: 950;
}

.week,
.grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 12px;
}

.week {
  margin-bottom: 12px;
  color: #8da0ba;
  font-size: 12px;
  font-weight: 950;
  text-align: center;
}

.day {
  min-height: 72px;
  display: grid;
  place-items: center;
  gap: 2px;
  border: 0;
  border-radius: 22px;
  border: 1px solid transparent;
  background: #f6f9fd;
  color: #07101f;
  cursor: pointer;
  position: relative;
}


.day:hover:not(:disabled) {
  background: #eaf1fb;
  transform: translateY(-1px);
}

.day b {
  font-size: 18px;
  font-weight: 950;
}

.day small {
  font-size: 11px;
  font-weight: 850;
  color: #64748b;
}

.day.muted {
  opacity: .42;
}

.day.today {
  outline: 2px solid rgba(255, 63, 76, .28);
}

.day.selected {
  background: #ff3f4c;
  color: #fff;
  box-shadow: 0 18px 36px rgba(255, 63, 76, .24);
}

.day.selected small {
  color: rgba(255, 255, 255, .82);
}

.day.available {
  background: #ecfdf7;
  border-color: rgba(16, 185, 129, .18);
  color: #07101f;
}

.day.available small {
  color: #047857;
}

.day.closed {
  background: #fff1f2;
  border-color: rgba(255, 63, 76, .22);
  color: #ff3f4c;
}

.day.disabled {
  cursor: not-allowed;
  color: #98a8bc;
  background: #f6f8fb;
}

.day.disabled.closed {
  background: #fff1f2;
  color: #ff3f4c;
}

.day.selected {
  background: #ff3f4c;
  border-color: #ff3f4c;
  color: #fff;
  box-shadow: 0 18px 36px rgba(255, 63, 76, .24);
}

.day.selected small {
  color: rgba(255, 255, 255, .88);
}


footer {
  display: flex;
  flex-wrap: wrap;
  gap: 14px 18px;
  margin-top: 22px;
  color: #64748b;
  font-size: 12px;
  font-weight: 850;
}

footer span {
  display: flex;
  align-items: center;
  gap: 8px;
}

.legend-badge {
  min-height: 30px;
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 0 12px;
  font-size: 12px;
  font-weight: 950;
}

.available-legend {
  background: #ecfdf7;
  color: #047857;
  border: 1px solid rgba(16, 185, 129, .18);
}

.closed-legend {
  background: #fff1f2;
  color: #ff3f4c;
  border: 1px solid rgba(255, 63, 76, .22);
}


@media (max-width: 640px) {
  .dp-panel {
    padding: 20px;
  }

  .week,
  .grid {
    gap: 8px;
  }

  .day {
    min-height: 58px;
    border-radius: 18px;
  }

  .day b {
    font-size: 16px;
  }
}
</style>
