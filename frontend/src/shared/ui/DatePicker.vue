<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { CalendarDays, ChevronLeft, ChevronRight, X } from '@lucide/vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  min: { type: String, default: '' },
  max: { type: String, default: '' },
  days: { type: Array, default: () => [] },
  placeholder: { type: String, default: 'Выберите дату' },
  error: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue', 'change', 'open', 'close'])
const root = ref(null)
const open = ref(false)
const visibleMonth = ref(new Date())
const popupStyle = ref({})
const previousOverflow = ref('')
const previousPaddingRight = ref('')

const monthNames = ['январь', 'февраль', 'март', 'апрель', 'май', 'июнь', 'июль', 'август', 'сентябрь', 'октябрь', 'ноябрь', 'декабрь']
const weekdayNames = ['пн', 'вт', 'ср', 'чт', 'пт', 'сб', 'вс']

const dayMap = computed(() => {
  const map = new Map()
  props.days.forEach((day) => {
    const key = day.date || day.pickup_date || day.day
    if (key) map.set(String(key).slice(0, 10), day)
  })
  return map
})

const displayValue = computed(() => props.modelValue ? formatHumanDate(props.modelValue) : props.placeholder)
const currentMonthLabel = computed(() => `${monthNames[visibleMonth.value.getMonth()]} ${visibleMonth.value.getFullYear()}`)

const cells = computed(() => {
  const year = visibleMonth.value.getFullYear()
  const month = visibleMonth.value.getMonth()
  const first = new Date(year, month, 1)
  const start = new Date(first)
  const day = first.getDay() || 7
  start.setDate(first.getDate() - day + 1)
  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(start)
    date.setDate(start.getDate() + index)
    const value = toISO(date)
    const info = dayMap.value.get(value)
    return {
      date,
      value,
      day: date.getDate(),
      inMonth: date.getMonth() === month,
      selected: value === props.modelValue,
      today: value === toISO(new Date()),
      info,
      disabled: !isSelectable(value, info),
      label: dayLimitLabel(info),
      closed: isClosed(info),
      full: isFull(info),
    }
  })
})

function toISO(date) {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function parseISO(value) {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(String(value || ''))) return null
  const [year, month, day] = value.split('-').map(Number)
  return new Date(year, month - 1, day)
}

function formatHumanDate(value) {
  const date = parseISO(value)
  if (!date) return value
  return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }).format(date)
}

function isClosed(info) {
  return Boolean(info?.is_closed || info?.closed)
}

function isFull(info) {
  const max = Number(info?.max_orders ?? info?.capacity ?? info?.limit ?? 0)
  const current = Number(info?.current_orders ?? info?.orders_count ?? info?.booked ?? 0)
  return max === 0 || (Number.isFinite(max) && max > 0 && current >= max)
}

function dayLimitLabel(info) {
  if (!info) return ''
  if (isClosed(info)) return 'закрыто'
  const max = info.max_orders ?? info.capacity ?? info.limit
  const current = info.current_orders ?? info.orders_count ?? info.booked ?? 0
  if (max === undefined || max === null) return 'доступно'
  return `${current}/${max}`
}

function isSelectable(value, info = dayMap.value.get(value)) {
  if (props.min && value < props.min) return false
  if (props.max && value > props.max) return false
  if (info && (isClosed(info) || isFull(info))) return false
  return true
}

function scrollbarWidth() {
  return window.innerWidth - document.documentElement.clientWidth
}

function lockScroll() {
  previousOverflow.value = document.body.style.overflow
  previousPaddingRight.value = document.body.style.paddingRight
  const width = scrollbarWidth()
  document.body.style.overflow = 'hidden'
  if (width > 0) document.body.style.paddingRight = `${width}px`
}

function unlockScroll() {
  document.body.style.overflow = previousOverflow.value
  document.body.style.paddingRight = previousPaddingRight.value
}

function updatePosition() {
  if (!root.value) return
  const rect = root.value.getBoundingClientRect()
  const margin = 16
  const gap = 8
  const width = Math.min(390, window.innerWidth - margin * 2)
  const left = Math.min(Math.max(rect.left, margin), window.innerWidth - width - margin)
  const estimatedHeight = 520
  const spaceBelow = window.innerHeight - rect.bottom - margin
  const spaceAbove = rect.top - margin
  const openAbove = spaceBelow < 460 && spaceAbove > spaceBelow
  const maxHeight = Math.max(360, Math.min(estimatedHeight, openAbove ? spaceAbove - gap : spaceBelow - gap))
  const top = openAbove ? Math.max(margin, rect.top - gap - maxHeight) : rect.bottom + gap
  popupStyle.value = { left: `${left}px`, top: `${top}px`, width: `${width}px`, maxHeight: `${maxHeight}px` }
}

async function setOpen(value) {
  open.value = value
  if (value) {
    const base = parseISO(props.modelValue || props.min) || new Date()
    visibleMonth.value = new Date(base.getFullYear(), base.getMonth(), 1)
    await nextTick()
    updatePosition()
    emit('open')
  } else {
    emit('close')
  }
}

function selectDate(cell) {
  if (cell.disabled) return
  emit('update:modelValue', cell.value)
  emit('change', cell.value)
  setOpen(false)
}

function moveMonth(direction) {
  visibleMonth.value = new Date(visibleMonth.value.getFullYear(), visibleMonth.value.getMonth() + direction, 1)
}

function onKeydown(event) {
  if (event.key === 'Escape') setOpen(false)
}

watch(open, (value) => {
  if (value) lockScroll()
  else unlockScroll()
})

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
  window.addEventListener('resize', updatePosition)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  window.removeEventListener('resize', updatePosition)
  if (open.value) unlockScroll()
})
</script>

<template>
  <div ref="root" class="relative">
    <button
      type="button"
      class="flex w-full items-center justify-between gap-3 rounded-[1.35rem] border bg-slate-50 px-5 py-4 text-left font-black outline-none transition"
      :class="error ? 'border-[#ff4248] bg-red-50/60 shadow-[0_0_0_4px_rgba(255,66,72,0.10)]' : open ? 'border-[#ff4248] bg-white shadow-[0_0_0_4px_rgba(255,66,72,0.10)]' : 'border-slate-200 hover:border-slate-300'"
      @click="setOpen(true)"
    >
      <span :class="modelValue ? 'text-[#07101f]' : 'text-slate-400'">{{ displayValue }}</span>
      <CalendarDays class="h-5 w-5 text-slate-500" />
    </button>
    <p v-if="error" class="mt-2 text-sm font-bold text-[#e11d48]">{{ error }}</p>

    <Teleport to="body">
      <div v-if="open" class="fixed inset-0 z-[999] bg-black/5 backdrop-blur-[1px]" @click.self="setOpen(false)"></div>
      <div
        v-if="open"
        class="fixed z-[1000] overflow-y-auto overscroll-contain rounded-[1.7rem] border border-slate-200 bg-white p-4 text-[#07101f] shadow-[0_30px_90px_rgba(15,23,42,0.24)]"
        :style="popupStyle"
      >
        <div class="flex items-center justify-between gap-3">
          <button type="button" class="grid h-11 w-11 place-items-center rounded-2xl bg-slate-100 text-slate-600 transition hover:bg-[#ff4248] hover:text-white" @click="moveMonth(-1)">
            <ChevronLeft class="h-5 w-5" />
          </button>
          <div class="text-center">
            <div class="text-lg font-black capitalize">{{ currentMonthLabel }}</div>
            <div class="mt-1 text-xs font-bold text-slate-400">дни с ограничениями логиста подсвечены</div>
          </div>
          <button type="button" class="grid h-11 w-11 place-items-center rounded-2xl bg-slate-100 text-slate-600 transition hover:bg-[#ff4248] hover:text-white" @click="moveMonth(1)">
            <ChevronRight class="h-5 w-5" />
          </button>
        </div>

        <div class="mt-4 grid grid-cols-7 gap-1 text-center text-[11px] font-black uppercase tracking-[0.16em] text-slate-400">
          <div v-for="name in weekdayNames" :key="name" class="py-2">{{ name }}</div>
        </div>

        <div class="grid grid-cols-7 gap-1">
          <button
            v-for="cell in cells"
            :key="cell.value"
            type="button"
            class="relative min-h-[54px] rounded-2xl p-1 text-left transition"
            :class="[
              cell.selected ? 'bg-[#ff4248] text-white shadow-[0_14px_34px_rgba(255,66,72,0.28)]' : cell.disabled ? 'cursor-not-allowed bg-slate-50 text-slate-300' : 'bg-white text-[#07101f] hover:bg-slate-100',
              !cell.inMonth && !cell.selected ? 'opacity-35' : '',
              cell.today && !cell.selected ? 'ring-2 ring-[#ff4248]/25' : '',
            ]"
            :disabled="cell.disabled"
            @click="selectDate(cell)"
          >
            <span class="block text-center text-sm font-black">{{ cell.day }}</span>
            <span
              v-if="cell.label"
              class="mt-1 block truncate rounded-full px-1.5 py-0.5 text-center text-[9px] font-black"
              :class="cell.selected ? 'bg-white/18 text-white' : cell.closed || cell.full ? 'bg-red-50 text-red-500' : 'bg-emerald-50 text-emerald-600'"
            >
              {{ cell.label }}
            </span>
          </button>
        </div>

        <div class="mt-4 rounded-2xl bg-slate-50 p-3 text-xs font-bold leading-5 text-slate-500">
          <div class="flex items-start gap-2"><span class="mt-1 h-2 w-2 rounded-full bg-emerald-500"></span>Доступные дни можно выбрать для сдачи или забора товара.</div>
          <div class="mt-1 flex items-start gap-2"><span class="mt-1 h-2 w-2 rounded-full bg-red-500"></span>Закрытые и заполненные дни заблокированы по ограничениям логиста.</div>
        </div>

        <button type="button" class="mt-3 inline-flex w-full items-center justify-center gap-2 rounded-2xl border border-slate-200 px-4 py-3 font-black text-slate-600 transition hover:bg-slate-50" @click="setOpen(false)">
          <X class="h-4 w-4" /> Закрыть
        </button>
      </div>
    </Teleport>
  </div>
</template>
