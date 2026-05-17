export function normalizeEmail(value) {
  return String(value || '').trim().toLowerCase()
}

export function normalizePhone(value) {
  return String(value || '').trim().replace(/[\s()-]/g, '')
}

export function validateEmail(value) {
  const email = normalizeEmail(value)
  if (!email) return 'Введите email.'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return 'Введите корректный email, например client@example.com.'
  return ''
}

export function validatePassword(value, { min = 6 } = {}) {
  const password = String(value || '')
  if (!password) return 'Введите пароль.'
  if (password.length < min) return `Пароль должен быть не короче ${min} символов.`
  return ''
}

export function validateFullName(value) {
  const fullName = String(value || '').trim().replace(/\s+/g, ' ')
  if (!fullName) return 'Введите ФИО.'
  if (fullName.length < 2) return 'ФИО должно содержать минимум 2 символа.'
  if (!/^[А-Яа-яЁёA-Za-z\s-]+$/.test(fullName)) return 'ФИО может содержать только буквы, пробелы и дефис.'
  return ''
}

export function validatePhone(value) {
  const phone = normalizePhone(value)
  if (!phone) return ''
  if (!/^\+?\d{10,16}$/.test(phone)) return 'Телефон должен содержать от 10 до 16 цифр.'
  return ''
}

export function humanizeApiError(error, fallback = 'Не удалось выполнить запрос.') {
  const raw = error?.message || error?.payload?.error || error?.payload?.message || ''
  if (!raw) return fallback
  if (raw.includes('invalid credentials') || raw.includes('unauthorized')) return 'Неверный email или пароль.'
  if (raw.includes('already exists') || raw.includes('duplicate')) return 'Пользователь с таким email уже существует.'
  if (raw.includes('decode json')) return 'Форма отправлена в неверном формате. Обновите страницу и попробуйте снова.'
  return raw
}

export function roleLabel(role) {
  const labels = {
    admin: 'Администратор',
    logist: 'Логист',
    worker: 'Работник склада',
    client: 'Клиент',
  }
  return labels[String(role || '').toLowerCase()] || role || 'Пользователь'
}
