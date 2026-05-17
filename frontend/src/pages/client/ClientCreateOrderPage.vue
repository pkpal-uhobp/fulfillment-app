<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowRight, CalendarDays, PackagePlus, Plus, Trash2 } from '@lucide/vue'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import DatePicker from '@/shared/ui/DatePicker.vue'
import FieldMessage from '@/shared/ui/FieldMessage.vue'
import TimeSelect from '@/shared/ui/TimeSelect.vue'
import { apiFetch, getCurrentUser } from '@/shared/api/http'
import {
  calendarDayAvailability,
  calendarDayClosed,
  calendarDayMap,
  catalogOption,
  compactName,
  isDestinationWarehouse,
  isNonNegativeNumber,
  isPositiveInteger,
  isTerminalWarehouse,
  isValidPhone,
  normalizeCollection,
  timeToMinutes,
  toNumberOrUndefined,
  warehouseOption,
  warehouseTypeLabel,
} from './clientUtils'

const router = useRouter()
const loading = ref(false)
const loadingCatalogs = ref(true)
const error = ref('')
const success = ref('')
const warehouses = ref([])
const productTypes = ref([])
const cargoPlaceTypes = ref([])
const calendarDays = ref([])
const fieldErrors = ref({})
const user = ref(getCurrentUser())

const today = new Date().toISOString().slice(0, 10)
const maxDate = (() => {
  const date = new Date()
  date.setDate(date.getDate() + 60)
  return date.toISOString().slice(0, 10)
})()

const form = ref({
  receiving_warehouse_id: '',
  destination_warehouse_id: '',
  product_type_id: '',
  handover_type: 'self_delivery',
  self_delivery_date: today,
  self_delivery_time_from: '10:00',
  self_delivery_time_to: '18:00',
  comment: '',
  pickup: {
    pickup_address: '',
    pickup_date: today,
    pickup_time_from: '10:00',
    pickup_time_to: '18:00',
    contact_name: user.value?.full_name || '',
    contact_phone: user.value?.phone || '',
    comment: '',
  },
  cargo_places: [
    {
      cargo_place_type_id: '',
      quantity: '1',
      weight_per_place_kg: '',
      length_cm: '',
      width_cm: '',
      height_cm: '',
      comment: '',
    },
  ],
})

const receivingWarehouses = computed(() => warehouses.value.filter(isTerminalWarehouse))
const destinationWarehouses = computed(() => warehouses.value.filter(isDestinationWarehouse))
const selectedReceivingWarehouse = computed(() => warehouses.value.find((warehouse) => String(warehouse.id) === String(form.value.receiving_warehouse_id)))
const selectedDestinationWarehouse = computed(() => warehouses.value.find((warehouse) => String(warehouse.id) === String(form.value.destination_warehouse_id)))
const receivingWarehouseOptions = computed(() => receivingWarehouses.value.map(warehouseOption))
const destinationWarehouseOptions = computed(() => destinationWarehouses.value.map(warehouseOption))
const productTypeOptions = computed(() => productTypes.value.map(catalogOption))
const cargoPlaceTypeOptions = computed(() => cargoPlaceTypes.value.map(catalogOption))
const currentCalendarMap = computed(() => calendarDayMap(calendarDays.value))
const selectedHandoverDate = computed(() => form.value.handover_type === 'pickup' ? form.value.pickup.pickup_date : form.value.self_delivery_date)
const selectedCalendarDay = computed(() => currentCalendarMap.value[selectedHandoverDate.value])
const totalPlaces = computed(() => form.value.cargo_places.reduce((sum, place) => sum + Number(place.quantity || 0), 0))

function optionalString(value) {
  const trimmed = String(value || '').trim()
  return trimmed || undefined
}

function numericValue(value) {
  return toNumberOrUndefined(value)
}

function addCargoPlace() {
  form.value.cargo_places.push({
    cargo_place_type_id: cargoPlaceTypes.value[0]?.id || '',
    quantity: '1',
    weight_per_place_kg: '',
    length_cm: '',
    width_cm: '',
    height_cm: '',
    comment: '',
  })
}

function removeCargoPlace(index) {
  if (form.value.cargo_places.length === 1) return
  form.value.cargo_places.splice(index, 1)
}

function clearFieldError(key) {
  if (fieldErrors.value[key]) {
    fieldErrors.value = { ...fieldErrors.value, [key]: '' }
  }
}

function setHandoverType(value) {
  form.value.handover_type = value
  clearFieldError('handover_type')
  loadCalendar()
}

function dayRestrictionMessage(date) {
  if (!date) return 'Выберите дату.'
  if (date < today) return 'Дата не может быть раньше сегодняшнего дня.'
  if (date > maxDate) return 'Дата слишком далеко. Выберите день в ближайшие 60 дней.'
  const day = currentCalendarMap.value[date]
  if (day && calendarDayClosed(day)) return 'Эта дата закрыта логистом или лимит заявок исчерпан.'
  return ''
}

function timeRangeMessage(from, to) {
  const start = timeToMinutes(from)
  const end = timeToMinutes(to)
  if (start === null || end === null) return 'Выберите время.'
  if (start >= end) return 'Время окончания должно быть позже времени начала.'
  return ''
}

function validate() {
  const errors = {}
  if (!form.value.receiving_warehouse_id) errors.receiving_warehouse_id = 'Выберите склад приёмки.'
  if (!form.value.destination_warehouse_id) errors.destination_warehouse_id = 'Выберите склад назначения.'
  if (!form.value.product_type_id) errors.product_type_id = 'Выберите тип товара.'

  if (form.value.handover_type === 'pickup') {
    if (!form.value.pickup.pickup_address.trim()) errors.pickup_address = 'Укажите адрес забора.'
    if (!form.value.pickup.contact_name.trim()) errors.contact_name = 'Укажите контактное лицо.'
    if (!isValidPhone(form.value.pickup.contact_phone)) errors.contact_phone = 'Укажите телефон в формате +79991234567.'
    const dateError = dayRestrictionMessage(form.value.pickup.pickup_date)
    if (dateError) errors.pickup_date = dateError
    const timeError = timeRangeMessage(form.value.pickup.pickup_time_from, form.value.pickup.pickup_time_to)
    if (timeError) errors.pickup_time = timeError
  } else {
    const dateError = dayRestrictionMessage(form.value.self_delivery_date)
    if (dateError) errors.self_delivery_date = dateError
    const timeError = timeRangeMessage(form.value.self_delivery_time_from, form.value.self_delivery_time_to)
    if (timeError) errors.self_delivery_time = timeError
  }

  form.value.cargo_places.forEach((place, index) => {
    if (!place.cargo_place_type_id) errors[`cargo_${index}_type`] = 'Выберите тип места.'
    if (!isPositiveInteger(place.quantity)) errors[`cargo_${index}_quantity`] = 'Количество должно быть целым числом больше 0.'
    if (!isNonNegativeNumber(place.weight_per_place_kg)) errors[`cargo_${index}_weight`] = 'Вес должен быть числом не меньше 0.'
    if (!isNonNegativeNumber(place.length_cm)) errors[`cargo_${index}_length`] = 'Длина должна быть числом не меньше 0.'
    if (!isNonNegativeNumber(place.width_cm)) errors[`cargo_${index}_width`] = 'Ширина должна быть числом не меньше 0.'
    if (!isNonNegativeNumber(place.height_cm)) errors[`cargo_${index}_height`] = 'Высота должна быть числом не меньше 0.'
  })

  fieldErrors.value = errors
  return Object.keys(errors).length === 0
}

function buildPayload() {
  const cargo_places = form.value.cargo_places.map((place) => ({
    cargo_place_type_id: Number(place.cargo_place_type_id),
    quantity: Number(place.quantity || 1),
    weight_per_place_kg: numericValue(place.weight_per_place_kg),
    length_cm: numericValue(place.length_cm),
    width_cm: numericValue(place.width_cm),
    height_cm: numericValue(place.height_cm),
    comment: optionalString(place.comment),
  }))

  const payload = {
    receiving_warehouse_id: Number(form.value.receiving_warehouse_id),
    destination_warehouse_id: Number(form.value.destination_warehouse_id),
    product_type_id: Number(form.value.product_type_id),
    handover_type: form.value.handover_type,
    comment: optionalString(form.value.comment),
    cargo_places,
  }

  if (form.value.handover_type === 'pickup') {
    payload.pickup = {
      pickup_address: form.value.pickup.pickup_address.trim(),
      pickup_date: form.value.pickup.pickup_date,
      pickup_time_from: optionalString(form.value.pickup.pickup_time_from),
      pickup_time_to: optionalString(form.value.pickup.pickup_time_to),
      contact_name: optionalString(form.value.pickup.contact_name),
      contact_phone: optionalString(form.value.pickup.contact_phone),
      comment: optionalString(form.value.pickup.comment),
    }
  } else {
    payload.self_delivery_date = form.value.self_delivery_date
    payload.self_delivery_time_from = optionalString(form.value.self_delivery_time_from)
    payload.self_delivery_time_to = optionalString(form.value.self_delivery_time_to)
  }

  return payload
}

async function loadCatalogs() {
  loadingCatalogs.value = true
  error.value = ''
  try {
    const [warehousesPayload, productTypesPayload, cargoPlaceTypesPayload] = await Promise.all([
      apiFetch('/warehouses'),
      apiFetch('/product-types'),
      apiFetch('/cargo-place-types'),
    ])

    warehouses.value = normalizeCollection(warehousesPayload, 'warehouses').filter((warehouse) => warehouse.is_active !== false)
    productTypes.value = normalizeCollection(productTypesPayload, 'product_types').filter((item) => item.is_active !== false)
    cargoPlaceTypes.value = normalizeCollection(cargoPlaceTypesPayload, 'cargo_place_types').filter((item) => item.is_active !== false)

    if (!form.value.receiving_warehouse_id) form.value.receiving_warehouse_id = receivingWarehouses.value[0]?.id || ''
    if (!form.value.destination_warehouse_id) form.value.destination_warehouse_id = destinationWarehouses.value[0]?.id || ''
    if (!form.value.product_type_id) form.value.product_type_id = productTypes.value[0]?.id || ''
    form.value.cargo_places.forEach((place) => {
      if (!place.cargo_place_type_id) place.cargo_place_type_id = cargoPlaceTypes.value[0]?.id || ''
    })
  } catch (err) {
    error.value = err?.message || 'Не удалось загрузить справочники.'
  } finally {
    loadingCatalogs.value = false
  }
}

async function loadCalendar() {
  if (!form.value.receiving_warehouse_id) return
  try {
    const dateFrom = today
    const toDate = new Date()
    toDate.setDate(toDate.getDate() + 60)
    const payload = await apiFetch(`/pickup-calendar?warehouse_id=${form.value.receiving_warehouse_id}&date_from=${dateFrom}&date_to=${toDate.toISOString().slice(0, 10)}`, { auth: true })
    calendarDays.value = normalizeCollection(payload, 'days')
  } catch {
    calendarDays.value = []
  }
}

async function submit() {
  error.value = ''
  success.value = ''
  if (!validate()) {
    error.value = 'Проверьте выделенные поля. Заявка не будет создана, пока данные не исправлены.'
    return
  }

  loading.value = true
  try {
    const payload = await apiFetch('/orders', {
      method: 'POST',
      auth: true,
      body: buildPayload(),
    })
    const order = payload?.order || payload?.data?.order || payload
    success.value = `Заявка #${order?.id || ''} создана.`
    setTimeout(() => router.push(order?.id ? `/client/orders/${order.id}` : '/client/orders'), 500)
  } catch (err) {
    error.value = err?.message || 'Не удалось создать заявку.'
  } finally {
    loading.value = false
  }
}

watch(() => form.value.receiving_warehouse_id, () => {
  clearFieldError('receiving_warehouse_id')
  loadCalendar()
})
watch(() => form.value.destination_warehouse_id, () => clearFieldError('destination_warehouse_id'))
watch(() => form.value.product_type_id, () => clearFieldError('product_type_id'))
watch(selectedHandoverDate, () => {
  clearFieldError('self_delivery_date')
  clearFieldError('pickup_date')
})

onMounted(async () => {
  await loadCatalogs()
  await loadCalendar()
})
</script>

<template>
  <section class="space-y-6">
    <div class="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 backdrop-blur sm:p-8">
      <p class="text-xs font-black uppercase tracking-[0.45em] text-[#ff9ca0]">Новая заявка</p>
      <h1 class="mt-4 text-5xl font-black tracking-[-0.06em]">Оформить передачу товара</h1>
      <p class="mt-4 max-w-3xl text-white/65">Выберите маршрут, дату и состав грузовых мест. Календарь учитывает закрытые дни и лимиты, заданные логистом.</p>
    </div>

    <div v-if="error" class="rounded-3xl border border-red-400/30 bg-red-500/10 px-6 py-5 font-bold text-red-100">{{ error }}</div>
    <div v-if="success" class="rounded-3xl border border-emerald-400/30 bg-emerald-500/10 px-6 py-5 font-bold text-emerald-100">{{ success }}</div>

    <form class="grid gap-6 xl:grid-cols-[1fr_380px]" novalidate @submit.prevent="submit">
      <div class="space-y-6">
        <section class="rounded-[2rem] bg-white p-6 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.12)]">
          <div class="flex items-center gap-3">
            <div class="grid h-12 w-12 place-items-center rounded-2xl bg-[#ff4248]/10 text-[#ff4248]"><PackagePlus class="h-6 w-6" /></div>
            <div>
              <h2 class="text-2xl font-black tracking-[-0.04em]">Маршрут и тип товара</h2>
              <p class="text-sm font-bold text-slate-500">Склады и справочники загружаются из backend.</p>
            </div>
          </div>

          <div class="mt-6 grid gap-4 md:grid-cols-2">
            <label class="block">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Склад приёмки</span>
              <BaseSelect v-model="form.receiving_warehouse_id" :options="receivingWarehouseOptions" :error="fieldErrors.receiving_warehouse_id" placeholder="Выберите склад приёмки" />
            </label>
            <label class="block">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Склад назначения</span>
              <BaseSelect v-model="form.destination_warehouse_id" :options="destinationWarehouseOptions" :error="fieldErrors.destination_warehouse_id" placeholder="Выберите склад назначения" />
            </label>
            <label class="block md:col-span-2">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Тип товара</span>
              <BaseSelect v-model="form.product_type_id" :options="productTypeOptions" :error="fieldErrors.product_type_id" placeholder="Выберите тип товара" />
            </label>
          </div>
        </section>

        <section class="rounded-[2rem] bg-white p-6 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.12)]">
          <h2 class="text-2xl font-black tracking-[-0.04em]">Способ передачи</h2>
          <div class="mt-5 grid gap-3 sm:grid-cols-2">
            <button type="button" class="min-h-[140px] rounded-3xl border p-5 text-left transition" :class="form.handover_type === 'self_delivery' ? 'border-[#ff4248] bg-[#ff4248]/5 shadow-[0_0_0_4px_rgba(255,66,72,0.08)]' : 'border-slate-200 bg-slate-50 hover:border-slate-300'" @click="setHandoverType('self_delivery')">
              <div class="font-black">Сдача на склад</div>
              <p class="mt-2 text-sm leading-6 text-slate-500">Вы сами привозите товар на выбранный терминал.</p>
            </button>
            <button type="button" class="min-h-[140px] rounded-3xl border p-5 text-left transition" :class="form.handover_type === 'pickup' ? 'border-[#ff4248] bg-[#ff4248]/5 shadow-[0_0_0_4px_rgba(255,66,72,0.08)]' : 'border-slate-200 bg-slate-50 hover:border-slate-300'" @click="setHandoverType('pickup')">
              <div class="font-black">Забор с адреса</div>
              <p class="mt-2 text-sm leading-6 text-slate-500">Логист назначает забор, рабочий принимает груз, клиент видит движение в кабинете.</p>
            </button>
          </div>

          <div v-if="form.handover_type === 'self_delivery'" class="mt-6 grid gap-4 md:grid-cols-3">
            <label class="block md:col-span-1">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Дата сдачи</span>
              <DatePicker v-model="form.self_delivery_date" :min="today" :max="maxDate" :days="calendarDays" :error="fieldErrors.self_delivery_date" placeholder="Выберите дату" />
            </label>
            <label class="block">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">С</span>
              <TimeSelect v-model="form.self_delivery_time_from" :from="8" :to="22" :error="fieldErrors.self_delivery_time" placeholder="С" />
            </label>
            <label class="block">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">До</span>
              <TimeSelect v-model="form.self_delivery_time_to" :from="8" :to="23" placeholder="До" />
            </label>
          </div>

          <div v-else class="mt-6 grid gap-4 md:grid-cols-2">
            <label class="block md:col-span-2">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Адрес забора</span>
              <input v-model="form.pickup.pickup_address" class="w-full rounded-2xl border bg-slate-50 px-4 py-4 font-bold outline-none transition focus:bg-white" :class="fieldErrors.pickup_address ? 'border-[#ff4248]' : 'border-slate-200 focus:border-[#ff4248]'" placeholder="Москва, ул. Складская, 1" @input="clearFieldError('pickup_address')" />
              <FieldMessage :message="fieldErrors.pickup_address" />
            </label>
            <label class="block">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Дата забора</span>
              <DatePicker v-model="form.pickup.pickup_date" :min="today" :max="maxDate" :days="calendarDays" :error="fieldErrors.pickup_date" placeholder="Выберите дату" />
            </label>
            <div class="grid grid-cols-2 gap-3">
              <label class="block">
                <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">С</span>
                <TimeSelect v-model="form.pickup.pickup_time_from" :from="8" :to="22" :error="fieldErrors.pickup_time" placeholder="С" />
              </label>
              <label class="block">
                <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">До</span>
                <TimeSelect v-model="form.pickup.pickup_time_to" :from="8" :to="23" placeholder="До" />
              </label>
            </div>
            <label class="block">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Контакт</span>
              <input v-model="form.pickup.contact_name" class="w-full rounded-2xl border bg-slate-50 px-4 py-4 font-bold outline-none transition focus:bg-white" :class="fieldErrors.contact_name ? 'border-[#ff4248]' : 'border-slate-200 focus:border-[#ff4248]'" placeholder="Иванов Иван" @input="clearFieldError('contact_name')" />
              <FieldMessage :message="fieldErrors.contact_name" />
            </label>
            <label class="block">
              <span class="mb-2 block text-xs font-black uppercase tracking-[0.25em] text-slate-400">Телефон</span>
              <input v-model="form.pickup.contact_phone" class="w-full rounded-2xl border bg-slate-50 px-4 py-4 font-bold outline-none transition focus:bg-white" :class="fieldErrors.contact_phone ? 'border-[#ff4248]' : 'border-slate-200 focus:border-[#ff4248]'" inputmode="tel" placeholder="+79991234567" @input="clearFieldError('contact_phone')" />
              <FieldMessage :message="fieldErrors.contact_phone" />
            </label>
          </div>
        </section>

        <section class="rounded-[2rem] bg-white p-6 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.12)]">
          <div class="flex flex-wrap items-center justify-between gap-4">
            <div>
              <h2 class="text-2xl font-black tracking-[-0.04em]">Грузовые места</h2>
              <p class="mt-1 text-sm font-bold text-slate-500">Укажите количество, вес и габариты. Подсказки покажут ошибки до отправки.</p>
            </div>
            <button type="button" class="inline-flex items-center gap-2 rounded-2xl bg-[#07101f] px-4 py-3 text-sm font-black text-white" @click="addCargoPlace">
              <Plus class="h-4 w-4" /> Добавить
            </button>
          </div>

          <div class="mt-5 grid gap-4">
            <div v-for="(place, index) in form.cargo_places" :key="index" class="rounded-3xl border border-slate-200 bg-slate-50 p-4">
              <div class="flex items-center justify-between gap-4">
                <div class="font-black">Место {{ index + 1 }}</div>
                <button v-if="form.cargo_places.length > 1" type="button" class="rounded-xl p-2 text-red-500 hover:bg-red-50" @click="removeCargoPlace(index)">
                  <Trash2 class="h-5 w-5" />
                </button>
              </div>
              <div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
                <label class="block md:col-span-2">
                  <span class="mb-2 block text-xs font-black uppercase tracking-[0.2em] text-slate-400">Тип места</span>
                  <BaseSelect v-model="place.cargo_place_type_id" :options="cargoPlaceTypeOptions" :error="fieldErrors[`cargo_${index}_type`]" placeholder="Выберите тип места" compact />
                </label>
                <label class="block">
                  <span class="mb-2 block text-xs font-black uppercase tracking-[0.2em] text-slate-400">Количество</span>
                  <input v-model="place.quantity" class="w-full rounded-2xl border bg-white px-4 py-3 font-bold outline-none focus:border-[#ff4248]" :class="fieldErrors[`cargo_${index}_quantity`] ? 'border-[#ff4248]' : 'border-slate-200'" inputmode="numeric" placeholder="1" @input="clearFieldError(`cargo_${index}_quantity`)" />
                  <FieldMessage :message="fieldErrors[`cargo_${index}_quantity`]" />
                </label>
                <label class="block">
                  <span class="mb-2 block text-xs font-black uppercase tracking-[0.2em] text-slate-400">Вес, кг</span>
                  <input v-model="place.weight_per_place_kg" class="w-full rounded-2xl border bg-white px-4 py-3 font-bold outline-none focus:border-[#ff4248]" :class="fieldErrors[`cargo_${index}_weight`] ? 'border-[#ff4248]' : 'border-slate-200'" inputmode="decimal" placeholder="0" @input="clearFieldError(`cargo_${index}_weight`)" />
                  <FieldMessage :message="fieldErrors[`cargo_${index}_weight`]" />
                </label>
                <label class="block">
                  <span class="mb-2 block text-xs font-black uppercase tracking-[0.2em] text-slate-400">Длина, см</span>
                  <input v-model="place.length_cm" class="w-full rounded-2xl border bg-white px-4 py-3 font-bold outline-none focus:border-[#ff4248]" :class="fieldErrors[`cargo_${index}_length`] ? 'border-[#ff4248]' : 'border-slate-200'" inputmode="decimal" placeholder="0" @input="clearFieldError(`cargo_${index}_length`)" />
                  <FieldMessage :message="fieldErrors[`cargo_${index}_length`]" />
                </label>
                <label class="block">
                  <span class="mb-2 block text-xs font-black uppercase tracking-[0.2em] text-slate-400">Ширина, см</span>
                  <input v-model="place.width_cm" class="w-full rounded-2xl border bg-white px-4 py-3 font-bold outline-none focus:border-[#ff4248]" :class="fieldErrors[`cargo_${index}_width`] ? 'border-[#ff4248]' : 'border-slate-200'" inputmode="decimal" placeholder="0" @input="clearFieldError(`cargo_${index}_width`)" />
                  <FieldMessage :message="fieldErrors[`cargo_${index}_width`]" />
                </label>
                <label class="block">
                  <span class="mb-2 block text-xs font-black uppercase tracking-[0.2em] text-slate-400">Высота, см</span>
                  <input v-model="place.height_cm" class="w-full rounded-2xl border bg-white px-4 py-3 font-bold outline-none focus:border-[#ff4248]" :class="fieldErrors[`cargo_${index}_height`] ? 'border-[#ff4248]' : 'border-slate-200'" inputmode="decimal" placeholder="0" @input="clearFieldError(`cargo_${index}_height`)" />
                  <FieldMessage :message="fieldErrors[`cargo_${index}_height`]" />
                </label>
                <label class="block xl:col-span-1">
                  <span class="mb-2 block text-xs font-black uppercase tracking-[0.2em] text-slate-400">Комментарий</span>
                  <input v-model="place.comment" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 font-bold outline-none focus:border-[#ff4248]" placeholder="Опционально" />
                </label>
              </div>
            </div>
          </div>
        </section>
      </div>

      <aside class="space-y-6">
        <div class="rounded-[2rem] border border-white/10 bg-[#0b1527] p-6">
          <p class="text-xs font-black uppercase tracking-[0.35em] text-[#ff9ca0]">Итог</p>
          <h2 class="mt-3 text-3xl font-black tracking-[-0.05em]">Проверьте заявку</h2>
          <div class="mt-6 grid gap-4 text-sm">
            <div class="rounded-2xl bg-white/[0.06] p-4">
              <div class="text-white/45">Приёмка</div>
              <div class="mt-1 font-black">{{ compactName(selectedReceivingWarehouse?.name || '—') }}</div>
              <div class="mt-1 text-white/45">{{ warehouseTypeLabel(selectedReceivingWarehouse?.warehouse_type) }}</div>
            </div>
            <div class="rounded-2xl bg-white/[0.06] p-4">
              <div class="text-white/45">Назначение</div>
              <div class="mt-1 font-black">{{ compactName(selectedDestinationWarehouse?.name || '—') }}</div>
              <div class="mt-1 text-white/45">{{ warehouseTypeLabel(selectedDestinationWarehouse?.warehouse_type) }}</div>
            </div>
            <div class="rounded-2xl bg-white/[0.06] p-4">
              <div class="text-white/45">Дата</div>
              <div class="mt-1 font-black">{{ selectedHandoverDate }}</div>
              <div class="mt-1 text-white/45">{{ calendarDayAvailability(selectedCalendarDay) }}</div>
            </div>
            <div class="rounded-2xl bg-white/[0.06] p-4">
              <div class="text-white/45">Грузовые места</div>
              <div class="mt-1 font-black">{{ totalPlaces }} шт.</div>
            </div>
          </div>

          <button class="mt-6 flex w-full items-center justify-center gap-3 rounded-2xl bg-[#ff4248] px-6 py-5 font-black text-white shadow-[0_20px_60px_rgba(255,66,72,0.28)] disabled:cursor-not-allowed disabled:opacity-60" :disabled="loading || loadingCatalogs" type="submit">
            {{ loading ? 'Создаём...' : 'Создать заявку' }} <ArrowRight class="h-5 w-5" />
          </button>
        </div>

        <div class="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6">
          <div class="flex items-center gap-3">
            <CalendarDays class="h-6 w-6 text-[#ff9ca0]" />
            <h3 class="font-black">Календарь приёмки</h3>
          </div>
          <p class="mt-3 text-sm leading-6 text-white/55">Закрытые дни и заполненные лимиты недоступны для выбора в календаре.</p>
          <div class="mt-4 grid gap-2">
            <div v-for="day in calendarDays.slice(0, 8)" :key="day.date || day.pickup_date" class="rounded-2xl bg-white/[0.06] px-4 py-3 text-sm">
              <div class="font-black">{{ day.date || day.pickup_date || 'Дата' }}</div>
              <div class="mt-1" :class="calendarDayClosed(day) ? 'text-red-200' : 'text-emerald-200'">{{ calendarDayAvailability(day) }}</div>
            </div>
            <div v-if="!calendarDays.length" class="rounded-2xl bg-white/[0.06] px-4 py-3 text-sm text-white/45">Календарь не загружен или пока пуст.</div>
          </div>
        </div>
      </aside>
    </form>
  </section>
</template>
