const rawBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1'

// В dev-режиме frontend ходит через Vite proxy, чтобы не ловить CORS от Go backend.
export const API_BASE_URL = import.meta.env.DEV && rawBaseUrl.includes(':8080')
  ? '/api/v1'
  : rawBaseUrl.replace(/\/$/, '')

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
  } catch (error) {
    return null
  }
}

export function saveAuth(payload) {
  const source = payload?.data || payload || {}
  const tokens = source.tokens || source.Tokens || {
    access_token: source.access_token || source.accessToken,
    refresh_token: source.refresh_token || source.refreshToken,
  }
  const user = source.user || source.User || null

  if (tokens?.access_token) localStorage.setItem('access_token', tokens.access_token)
  if (tokens?.refresh_token) localStorage.setItem('refresh_token', tokens.refresh_token)
  if (user) localStorage.setItem('current_user', JSON.stringify(user))

  return { tokens, user }
}

export function clearAuth() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('current_user')
}

export async function loadMe() {
  const payload = await apiFetch('/auth/me', { auth: true })
  const user = payload?.user || payload?.User || null
  if (user) localStorage.setItem('current_user', JSON.stringify(user))
  return user
}

export async function apiFetch(path, options = {}) {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  const token = options.auth ? getAccessToken() : ''
  const headers = {
    Accept: 'application/json',
    ...(options.body ? { 'Content-Type': 'application/json' } : {}),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...(options.headers || {}),
  }

  const response = await fetch(`${API_BASE_URL}${normalizedPath}`, {
    ...options,
    headers,
  })

  const text = await response.text()
  let data = null

  try {
    data = text ? JSON.parse(text) : null
  } catch (error) {
    data = { message: text }
  }

  if (!response.ok) {
    const message = data?.error || data?.message || data?.detail || `HTTP ${response.status}`
    throw new Error(message)
  }

  return data
}
