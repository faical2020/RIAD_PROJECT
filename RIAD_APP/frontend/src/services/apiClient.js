const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1'

let isRefreshing = false
let refreshQueue = []

function getAccessToken() { return localStorage.getItem('token') }
function getRefreshToken() { return localStorage.getItem('refresh_token') }
function setTokens(access, refresh) {
  localStorage.setItem('token', access)
  if (refresh) localStorage.setItem('refresh_token', refresh)
}

async function refreshTokens() {
  const refresh = getRefreshToken()
  if (!refresh) throw new Error('no refresh token')

  const res = await fetch(`${API_BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refresh }),
  })

  if (!res.ok) {
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    localStorage.removeItem('role')
    window.location.href = '/login'
    throw new Error('refresh failed')
  }

  const data = await res.json()
  setTokens(data.access_token, data.refresh_token)
  return data.access_token
}

function onRefreshed(token) {
  refreshQueue.forEach(cb => cb(token))
  refreshQueue = []
}

export async function apiClient(path, options = {}) {
  const url = path.startsWith('http') ? path : `${API_BASE_URL}${path}`
  const token = getAccessToken()

  const headers = { 'Content-Type': 'application/json', ...options.headers }
  if (token) headers['Authorization'] = `Bearer ${token}`

  let res = await fetch(url, { ...options, headers })

  if (res.status === 401 && getRefreshToken()) {
    if (!isRefreshing) {
      isRefreshing = true
      try {
        const newToken = await refreshTokens()
        isRefreshing = false
        onRefreshed(newToken)
        headers['Authorization'] = `Bearer ${newToken}`
        res = await fetch(url, { ...options, headers })
      } catch {
        isRefreshing = false
        refreshQueue = []
        throw new Error('Session expirée, veuillez vous reconnecter')
      }
    } else {
      const newToken = await new Promise(resolve => refreshQueue.push(resolve))
      headers['Authorization'] = `Bearer ${newToken}`
      res = await fetch(url, { ...options, headers })
    }
  }

  return res
}

export async function apiGet(path) { return apiClient(path) }

export async function apiPost(path, body) {
  return apiClient(path, { method: 'POST', body: JSON.stringify(body) })
}

export async function apiPatch(path, body) {
  return apiClient(path, { method: 'PATCH', body: JSON.stringify(body) })
}

export async function apiDelete(path) {
  return apiClient(path, { method: 'DELETE' })
}
