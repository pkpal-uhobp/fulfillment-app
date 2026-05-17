export const orderStatusOptions = [
  { value: '', label: 'Все статусы' },
  { value: 'created', label: 'Создана' },
  { value: 'waiting_pickup', label: 'Ожидает забора' },
  { value: 'waiting_delivery', label: 'Ожидает сдачи' },
  { value: 'received', label: 'Принята на склад' },
  { value: 'stored', label: 'На хранении' },
  { value: 'assigned_to_shipping', label: 'Назначена к отгрузке' },
  { value: 'shipped', label: 'Отгружена' },
  { value: 'delivered', label: 'Доставлена' },
  { value: 'cancelled', label: 'Отменена' },
]

export const cargoStatusOptions = [
  { value: '', label: 'Все статусы' },
  { value: 'accepted', label: 'Принято' },
  { value: 'stored', label: 'На хранении' },
  { value: 'ready_to_ship', label: 'Готово к отгрузке' },
  { value: 'shipped', label: 'Отгружено' },
  { value: 'lost', label: 'Потеряно' },
  { value: 'damaged', label: 'Повреждено' },
  { value: 'cancelled', label: 'Отменено' },
]

export const shipmentStatusOptions = [
  { value: '', label: 'Все статусы' },
  { value: 'planned', label: 'Запланирована' },
  { value: 'loading', label: 'Погрузка' },
  { value: 'shipped', label: 'Отправлена' },
  { value: 'completed', label: 'Завершена' },
  { value: 'cancelled', label: 'Отменена' },
]

export const handoverLabels = {
  pickup: 'Забор с адреса',
  self_delivery: 'Сдача на склад',
}

export const orderStatusLabels = Object.fromEntries(orderStatusOptions.filter((x) => x.value).map((x) => [x.value, x.label]))
export const cargoStatusLabels = Object.fromEntries(cargoStatusOptions.filter((x) => x.value).map((x) => [x.value, x.label]))
export const shipmentStatusLabels = Object.fromEntries(shipmentStatusOptions.filter((x) => x.value).map((x) => [x.value, x.label]))

export function labelFromMap(map, value) {
  return map[value] || value || '—'
}

export function formatDate(value) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return new Intl.DateTimeFormat('ru-RU', { day: '2-digit', month: 'long', year: 'numeric' }).format(date)
}

export function formatDateTime(value) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

export function formatIsoDate(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

export function addDays(date, days) {
  const next = new Date(date)
  next.setDate(next.getDate() + days)
  return next
}

export function monthRange(date) {
  const from = new Date(date.getFullYear(), date.getMonth(), 1)
  const to = new Date(date.getFullYear(), date.getMonth() + 1, 0)
  return { from: formatIsoDate(from), to: formatIsoDate(to) }
}

export function warehouseName(id, warehouses = []) {
  const item = warehouses.find((warehouse) => Number(warehouse.id) === Number(id))
  if (!item) return id ? `Склад #${id}` : '—'
  const market = item.marketplace ? `${item.marketplace} · ` : ''
  return `${market}${item.name || `Склад #${item.id}`}`
}

export function zoneName(id, zones = []) {
  const item = zones.find((zone) => Number(zone.id) === Number(id))
  return item ? item.name : id ? `Зона #${id}` : '—'
}

export function gateName(id, gates = []) {
  const item = gates.find((gate) => Number(gate.id) === Number(id))
  return item ? item.name : id ? `Гейт #${id}` : '—'
}

export function unwrapList(payload, key) {
  if (!payload) return []
  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload[key])) return payload[key]
  if (Array.isArray(payload.items)) return payload.items
  if (Array.isArray(payload.data)) return payload.data
  return []
}

export function unwrapOne(payload, key) {
  if (!payload) return null
  return payload[key] || payload.item || payload.data || payload
}

export function shortId(prefix, id) {
  return `${prefix}-${String(id || '').padStart(4, '0')}`
}

export function isTerminalStatus(value) {
  return ['cancelled', 'delivered', 'shipped', 'completed', 'lost', 'damaged'].includes(value)
}
