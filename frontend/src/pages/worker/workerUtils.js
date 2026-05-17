export const workerCargoStatusOptions = [
  { value: 'accepted', label: 'Принято', description: 'Грузовое место принято' },
  { value: 'received', label: 'Принято на склад', description: 'Проверено при приёмке' },
  { value: 'stored', label: 'На хранении', description: 'Размещено в зоне' },
  { value: 'ready_to_ship', label: 'Готово к отгрузке', description: 'Ожидает отправки' },
  { value: 'shipped', label: 'Отгружено', description: 'Передано в отгрузку' },
  { value: 'damaged', label: 'Повреждено', description: 'Есть повреждение' },
  { value: 'lost', label: 'Потеряно', description: 'Нужно разбирательство' },
]
const map = Object.fromEntries(workerCargoStatusOptions.map((i) => [i.value, i.label]))
export function cargoStatusLabel(v) { return map[v] || v || 'Не указан' }
export function normalizeCollection(data, keys) { if (Array.isArray(data)) return data; for (const k of keys) if (Array.isArray(data?.[k])) return data[k]; return [] }
export function formatDateTime(value) { if (!value) return '—'; try { return new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) } catch { return value } }
