export const orderStatusLabels = {
  created: 'Создана',
  new: 'Создана',
  pending: 'Ожидает обработки',
  draft: 'Черновик',
  waiting_pickup: 'Ожидает забора',
  waiting_delivery: 'Ожидает сдачи на склад',
  pickup_scheduled: 'Забор назначен',
  accepted: 'Принята',
  received: 'Принята на склад',
  warehouse_received: 'Принята на склад',
  stored: 'На хранении',
  in_storage: 'На хранении',
  processing: 'В обработке',
  assigned_to_shipping: 'Назначена к отгрузке',
  ready_to_ship: 'Готова к отгрузке',
  shipped: 'Отгружена',
  in_transit: 'В пути',
  delivered: 'Доставлена',
  completed: 'Завершена',
  closed: 'Закрыта',
  cancelled: 'Отменена',
  canceled: 'Отменена',
  rejected: 'Отклонена',
}

export const cargoStatusLabels = {
  created: 'Создано',
  new: 'Создано',
  pending: 'Ожидает обработки',
  accepted: 'Принято',
  received: 'Принято на склад',
  warehouse_received: 'Принято на склад',
  stored: 'На хранении',
  in_storage: 'На хранении',
  processing: 'В обработке',
  assigned_to_shipping: 'Назначено к отгрузке',
  ready_to_ship: 'Готово к отгрузке',
  shipped: 'Отгружено',
  in_transit: 'В пути',
  delivered: 'Доставлено',
  completed: 'Завершено',
  damaged: 'Повреждено',
  lost: 'Утеряно',
  cancelled: 'Отменено',
  canceled: 'Отменено',
  rejected: 'Отклонено',
}

export const handoverLabels = {
  self_delivery: 'Сдача на склад',
  pickup: 'Забор с адреса',
}

export const roleLabels = {
  client: 'Клиент',
  worker: 'Работник склада',
  logist: 'Логист',
  admin: 'Администратор',
}

const unknownStatusLabels = {
  order_created: 'Заявка создана',
  order_received: 'Заявка принята на склад',
  cargo_created: 'Грузовое место создано',
  cargo_received: 'Грузовое место принято',
}

export function statusLabel(status, type = 'order') {
  const value = String(status || '').trim()
  const dict = type === 'cargo' ? cargoStatusLabels : orderStatusLabels
  return dict[value] || unknownStatusLabels[value] || translateUnknown(value) || '—'
}

export function translateUnknown(value) {
  if (!value) return ''
  return String(value)
    .replace(/_/g, ' ')
    .replace(/\b\w/g, (char) => char.toUpperCase())
}

export function roleLabel(value) {
  return roleLabels[value] || translateUnknown(value) || '—'
}

export function handoverLabel(value) {
  return handoverLabels[value] || translateUnknown(value) || '—'
}

export function formatDate(value) {
  if (!value) return '—'
  try {
    return new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(new Date(value))
  } catch {
    return String(value)
  }
}

export function formatDateTime(value) {
  if (!value) return '—'
  try {
    return new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
  } catch {
    return String(value)
  }
}

export function formatHumanDate(value) {
  if (!value) return '—'
  try {
    const [y, m, d] = String(value).slice(0, 10).split('-').map(Number)
    return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }).format(new Date(y, m - 1, d))
  } catch {
    return String(value)
  }
}

export function normalizeCollection(payload, key) {
  if (Array.isArray(payload)) return payload
  return payload?.[key] || payload?.data?.[key] || payload?.items || []
}

export function byId(items) {
  return Object.fromEntries((items || []).map((item) => [String(item.id), item]))
}

export function compactName(value = '') {
  return value
    .replace(/^TransitPro\s+/i, '')
    .replace(/^Fulfillment Transit\s+/i, '')
    .replace(/^WB\s+/i, 'Wildberries · ')
    .replace(/^Ozon\s+/i, 'Ozon · ')
    .replace(/\s+/g, ' ')
    .trim()
}

export function warehouseTypeLabel(value) {
  if (value === 'both') return 'универсальный склад'
  if (value === 'receiving') return 'склад приёмки'
  if (value === 'destination') return 'склад назначения'
  return value || 'склад'
}

export function isTerminalWarehouse(warehouse) {
  return warehouse?.warehouse_type === 'receiving' || warehouse?.warehouse_type === 'both'
}

export function isDestinationWarehouse(warehouse) {
  return warehouse?.warehouse_type === 'destination' || warehouse?.warehouse_type === 'both'
}

export function warehouseOption(warehouse) {
  return {
    value: warehouse.id,
    label: `${compactName(warehouse.name)} · ${warehouse.city || 'город не указан'}`,
    description: `${warehouseTypeLabel(warehouse.warehouse_type)}${warehouse.marketplace ? ` · ${warehouse.marketplace}` : ''}`,
  }
}

export function catalogOption(item) {
  return {
    value: item.id,
    label: item.name || item.title || `#${item.id}`,
    description: item.description || '',
  }
}

export function toNumberOrUndefined(value) {
  if (value === '' || value === null || value === undefined) return undefined
  const normalized = String(value).replace(',', '.').trim()
  const number = Number(normalized)
  return Number.isFinite(number) ? number : undefined
}

export function isPositiveInteger(value) {
  const number = Number(value)
  return Number.isInteger(number) && number > 0
}

export function isNonNegativeNumber(value) {
  if (value === '' || value === null || value === undefined) return true
  const number = Number(String(value).replace(',', '.'))
  return Number.isFinite(number) && number >= 0
}

export function timeToMinutes(value) {
  const match = /^(\d{2}):(\d{2})$/.exec(String(value || ''))
  if (!match) return null
  return Number(match[1]) * 60 + Number(match[2])
}

export function isValidEmail(value) {
  if (!value) return false
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(value).trim())
}

export function isValidPhone(value) {
  if (!value) return true
  return /^\+?[0-9\s()\-]{10,20}$/.test(String(value).trim())
}

export function calendarDayKey(day) {
  return String(day?.date || day?.pickup_date || day?.day || '').slice(0, 10)
}

export function calendarDayMap(days = []) {
  return Object.fromEntries(days.map((day) => [calendarDayKey(day), day]).filter(([key]) => key))
}

export function calendarDayClosed(day) {
  if (!day) return false
  const max = Number(day.max_orders ?? day.capacity ?? day.limit ?? 0)
  const current = Number(day.current_orders ?? day.orders_count ?? day.booked ?? 0)
  return Boolean(day.is_closed || day.closed) || max === 0 || (max > 0 && current >= max)
}

export function calendarDayAvailability(day) {
  if (!day) return 'нет данных календаря'
  if (calendarDayClosed(day)) return 'день закрыт или лимит исчерпан'
  const max = day.max_orders ?? day.capacity ?? day.limit
  const current = day.current_orders ?? day.orders_count ?? day.booked ?? 0
  if (max === undefined || max === null) return 'доступно'
  return `доступно ${current}/${max}`
}

export function historyCommentLabel(comment = '') {
  const value = String(comment || '')
  const normalized = value.toLowerCase()
  if (!value) return ''
  if (normalized.includes('seed:') && normalized.includes('created')) return 'Заявка создана из тестовых данных.'
  if (normalized.includes('seed:') && normalized.includes('received')) return 'Заявка принята на склад по тестовому сценарию.'
  if (normalized.includes('seed:') && normalized.includes('shipped')) return 'Заявка отгружена по тестовому сценарию.'
  if (normalized.includes('created')) return 'Заявка создана.'
  if (normalized.includes('received')) return 'Принято на склад.'
  if (normalized.includes('shipped')) return 'Отгружено.'
  return value
}
