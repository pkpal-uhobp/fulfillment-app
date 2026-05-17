const DEFAULT_API_BASE_URL = '/api/v1'

function normalizeBaseUrl(value) {
  const raw = value || DEFAULT_API_BASE_URL

  return raw.includes('localhost:8080') || raw.includes('127.0.0.1:8080')
    ? DEFAULT_API_BASE_URL
    : raw.replace(/\/+$/, '')
}

export const API_BASE_URL = normalizeBaseUrl(import.meta.env.VITE_API_BASE_URL)

function buildUrl(path) {
  if (/^https?:\/\//i.test(path)) return path

  return `${API_BASE_URL}${path?.startsWith('/') ? path : `/${path}`}`
}

function notifyAuthChanged() {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent('auth:changed'))
  }
}

function parseJwtPayload(token) {
  try {
    const [, payload] = String(token).split('.')

    if (!payload) return null

    return JSON.parse(
      decodeURIComponent(
        atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
          .split('')
          .map((char) => `%${char.charCodeAt(0).toString(16).padStart(2, '0')}`)
          .join(''),
      ),
    )
  } catch {
    return null
  }
}

export function getAccessToken() {
  return localStorage.getItem('access_token') || ''
}

export function getRefreshToken() {
  return localStorage.getItem('refresh_token') || ''
}

export function getCurrentUser() {
  try {
    const raw = localStorage.getItem('current_user')

    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export function getCurrentRole() {
  return String(getCurrentUser()?.role || '').toLowerCase()
}

export function clearAuth() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('current_user')
  notifyAuthChanged()
}

function normalizeUser(payload, token) {
  const user = payload?.user || payload?.me || payload?.profile || payload || null

  if (user && typeof user === 'object' && (user.email || user.full_name || user.role)) {
    return {
      ...user,
      role: String(user.role || '').toLowerCase(),
    }
  }

  const jwtPayload = token ? parseJwtPayload(token) : null

  return jwtPayload
    ? {
        id: jwtPayload.user_id || jwtPayload.sub || jwtPayload.id,
        email: jwtPayload.email || '',
        full_name: jwtPayload.full_name || jwtPayload.name || jwtPayload.email || 'Пользователь',
        role: String(jwtPayload.role || '').toLowerCase(),
      }
    : null
}

export function saveAuth(payload = {}) {
  const accessToken =
    payload.access_token ||
    payload.accessToken ||
    payload.token ||
    payload.tokens?.access_token ||
    payload.tokens?.accessToken ||
    ''

  const refreshToken =
    payload.refresh_token ||
    payload.refreshToken ||
    payload.tokens?.refresh_token ||
    payload.tokens?.refreshToken ||
    ''

  if (accessToken) localStorage.setItem('access_token', accessToken)
  if (refreshToken) localStorage.setItem('refresh_token', refreshToken)

  const user = normalizeUser(payload, accessToken)

  if (user) localStorage.setItem('current_user', JSON.stringify(user))

  notifyAuthChanged()

  return {
    accessToken,
    refreshToken,
    user,
  }
}

export function cabinetPathByRole(role) {
  const currentRole = String(role || '').toLowerCase()

  if (currentRole === 'admin') return '/admin'
  if (currentRole === 'worker' || currentRole === 'warehouse_worker') return '/worker'
  if (currentRole === 'logist' || currentRole === 'logistician') return '/logist'

  return '/client'
}

export async function apiFetch(path, options = {}) {
  const { method = 'GET', body, auth = false, headers = {} } = options
  const requestHeaders = new Headers(headers)
  let requestBody = body

  if (body !== undefined && body !== null && !(body instanceof FormData)) {
    if (!requestHeaders.has('Content-Type')) requestHeaders.set('Content-Type', 'application/json')
    requestBody = typeof body === 'string' ? body : JSON.stringify(body)
  }

  if (auth) {
    const token = getAccessToken()

    if (token) requestHeaders.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(buildUrl(path), {
    method,
    headers: requestHeaders,
    body: requestBody,
  })

  const text = await response.text()
  let data = null

  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = { message: text }
    }
  }

  if (!response.ok) {
    const error = new Error(data?.error || data?.message || data?.detail || `HTTP ${response.status}`)
    error.status = response.status
    error.data = data
    throw error
  }

  return data
}

export async function loadMe() {
  const data = await apiFetch('/auth/me', { auth: true })
  const user = normalizeUser(data, getAccessToken())

  if (user) {
    localStorage.setItem('current_user', JSON.stringify(user))
    notifyAuthChanged()
  }

  return user
}
