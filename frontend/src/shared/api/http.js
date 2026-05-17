const DEFAULT_API_BASE_URL = '/api/v1'

function normalizeBaseUrl(value) {
  const raw = (value || DEFAULT_API_BASE_URL).trim().replace(/\/$/, '')

  // В dev через Vite proxy нельзя ходить напрямую на :8080 из браузера,
  // иначе снова будет CORS. Оставляем относительный /api/v1.
  if (import.meta.env.DEV && (raw.includes('localhost:8080') || raw.includes('127.0.0.1:8080'))) {
    return DEFAULT_API_BASE_URL
  }

  return raw || DEFAULT_API_BASE_URL
}

export const API_BASE_URL = normalizeBaseUrl(import.meta.env.VITE_API_BASE_URL)

export function getAccessToken() {
  return (
    localStorage.getItem('access_token') ||
    localStorage.getItem('accessToken') ||
    localStorage.getItem('token') ||
    localStorage.getItem('auth_token') ||
    ''
  )
}

export function getRefreshToken() {
  return (
    localStorage.getItem('refresh_token') ||
    localStorage.getItem('refreshToken') ||
    ''
  )
}

export function getCurrentUser() {
  const raw = localStorage.getItem('current_user') || localStorage.getItem('user')
  if (!raw) return null

  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

function pickAuthSource(payload = {}) {
  return payload.data || payload.result || payload || {}
}

function pickTokens(source = {}) {
  const tokens = source.tokens || source.Tokens || source.auth || source.Auth || {}

  return {
    accessToken:
      source.access_token ||
      source.accessToken ||
      source.token ||
      source.jwt ||
      tokens.access_token ||
      tokens.accessToken ||
      tokens.token ||
      '',
    refreshToken:
      source.refresh_token ||
      source.refreshToken ||
      tokens.refresh_token ||
      tokens.refreshToken ||
      '',
  }
}

function pickUser(source = {}) {
  return (
    source.user ||
    source.User ||
    source.current_user ||
    source.currentUser ||
    source.me ||
    source.profile ||
    null
  )
}

export function saveAuth(payload = {}) {
  const source = pickAuthSource(payload)
  const { accessToken, refreshToken } = pickTokens(source)
  const user = pickUser(source)

  if (accessToken) localStorage.setItem('access_token', accessToken)
  if (refreshToken) localStorage.setItem('refresh_token', refreshToken)
  if (user) localStorage.setItem('current_user', JSON.stringify(user))

  window.dispatchEvent(new CustomEvent('auth:changed'))

  return { accessToken, refreshToken, user }
}

export function clearAuth() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('accessToken')
  localStorage.removeItem('token')
  localStorage.removeItem('auth_token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('refreshToken')
  localStorage.removeItem('current_user')
  localStorage.removeItem('user')

  window.dispatchEvent(new CustomEvent('auth:changed'))
}

export async function loadMe() {
  const payload = await apiFetch('/auth/me', { auth: true })
  const source = pickAuthSource(payload)
  const user = pickUser(source) || source

  if (user && typeof user === 'object') {
    localStorage.setItem('current_user', JSON.stringify(user))
    window.dispatchEvent(new CustomEvent('auth:changed'))
  }

  return user
}

function makeRequestBody(body) {
  if (body === undefined || body === null) return undefined
  if (body instanceof FormData) return body
  if (typeof body === 'string') return body
  return JSON.stringify(body)
}

export async function apiFetch(path, options = {}) {
  const {
    auth = false,
    headers = {},
    body,
    method,
    ...rest
  } = options

  const url = path.startsWith('http')
    ? path
    : `${API_BASE_URL}${path.startsWith('/') ? path : `/${path}`}`

  const requestHeaders = {
    Accept: 'application/json',
    ...headers,
  }

  if (body !== undefined && body !== null && !(body instanceof FormData)) {
    requestHeaders['Content-Type'] = requestHeaders['Content-Type'] || 'application/json'
  }

  if (auth) {
    const token = getAccessToken()
    if (!token) {
      const error = new Error('Сначала войдите в аккаунт')
      error.status = 401
      throw error
    }
    requestHeaders.Authorization = `Bearer ${token}`
  }

  const response = await fetch(url, {
    method: method || (body !== undefined ? 'POST' : 'GET'),
    headers: requestHeaders,
    body: makeRequestBody(body),
    ...rest,
  })

  const contentType = response.headers.get('content-type') || ''
  const payload = contentType.includes('application/json')
    ? await response.json().catch(() => null)
    : await response.text().catch(() => '')

  if (!response.ok) {
    const message =
      payload?.error ||
      payload?.message ||
      payload?.detail ||
      (typeof payload === 'string' && payload) ||
      `HTTP ${response.status}`

    const error = new Error(message)
    error.status = response.status
    error.payload = payload
    throw error
  }

  return payload
}

export function updateCurrentUserLocal(patch = {}) {
  const current = getCurrentUser() || {}
  const next = { ...current, ...patch }
  localStorage.setItem('current_user', JSON.stringify(next))
  window.dispatchEvent(new CustomEvent('auth:changed'))
  return next
}
