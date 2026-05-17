export const cargoStatuses = [
  { value: 'accepted', label: 'Принято', description: 'Груз принят и ожидает размещения' },
  { value: 'stored', label: 'На хранении', description: 'Груз размещён в назначенной зоне' },
  { value: 'ready_to_ship', label: 'Готово к отгрузке', description: 'Груз подготовлен к передаче на гейт' },
  { value: 'shipped', label: 'Отгружено', description: 'Груз передан на отгрузку' },
  { value: 'lost', label: 'Потеряно', description: 'Груз не найден, нужна проверка' },
  { value: 'damaged', label: 'Повреждено', description: 'Зафиксировано повреждение груза' },
  { value: 'cancelled', label: 'Отменено', description: 'Работа с грузом отменена' },
]

const statusMap = {
  created: 'Создано',
  pending_pickup: 'Ожидает забора',
  pending_self_delivery: 'Ожидает сдачи',
  accepted: 'Принято',
  stored: 'На хранении',
  ready_to_ship: 'Готово к отгрузке',
  assigned_to_shipment: 'Назначено к отгрузке',
  shipped: 'Отгружено',
  delivered: 'Доставлено',
  lost: 'Потеряно',
  damaged: 'Повреждено',
  cancelled: 'Отменено',
}

const toneMap = {
  created: 'gray',
  pending_pickup: 'amber',
  pending_self_delivery: 'amber',
  accepted: 'blue',
  stored: 'green',
  ready_to_ship: 'violet',
  assigned_to_shipment: 'violet',
  shipped: 'dark',
  delivered: 'green',
  lost: 'red',
  damaged: 'red',
  cancelled: 'gray',
}

export function statusLabel(status) {
  return statusMap[status] || status || 'Неизвестно'
}

export function statusTone(status) {
  return toneMap[status] || 'gray'
}

export function normalizeCollection(payload, keys = ['cargo_items', 'cargoItems', 'items', 'data', 'orders', 'shipments', 'history']) {
  if (Array.isArray(payload)) return payload
  for (const key of keys) {
    if (Array.isArray(payload?.[key])) return payload[key]
  }
  return []
}

export function extractCargo(payload) {
  return payload?.cargo_item || payload?.cargoItem || payload?.item || payload?.data || payload
}

export function formatDateTime(value) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

export function cargoTitle(item) {
  return item?.qr_code || item?.qrCode || `Грузовое место #${item?.id || '—'}`
}

export function routeTitle(item) {
  const order = item?.order || item?.order_info || {}
  const from = item?.receiving_warehouse?.name || item?.receivingWarehouse?.name || order?.receiving_warehouse?.name || order?.receivingWarehouse?.name || item?.receiving_warehouse_name || item?.receivingWarehouseName || 'Склад приёмки'
  const to = item?.destination_warehouse?.name || item?.destinationWarehouse?.name || order?.destination_warehouse?.name || order?.destinationWarehouse?.name || item?.destination_warehouse_name || item?.destinationWarehouseName || 'Склад назначения'
  return `${from} → ${to}`
}

export function zoneTitle(item) {
  return item?.storage_zone?.name || item?.storageZone?.name || item?.storage_zone_name || item?.storageZoneName || (item?.storage_zone_id ? `Зона #${item.storage_zone_id}` : 'Не назначена')
}

export function gateTitle(item) {
  return item?.gate?.name || item?.gate_name || item?.gateName || (item?.gate_id ? `Гейт #${item.gate_id}` : 'Не назначен')
}

export function safeFileName(value) {
  return String(value || 'qr').replace(/[^A-Za-zА-Яа-я0-9_.-]+/g, '_')
}
