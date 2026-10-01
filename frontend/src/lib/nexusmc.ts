import axios from 'axios'

// NexusMC 开放平台对接（经 EdgeOne Cloud Functions 同源代理，密钥只存在服务端）。
// 上游接口契约见 docs/reference/nexusmc-open-platform/。

export interface NexusmcUser {
  sub: string
  uid?: number
  username: string
  slug?: string
  avatar?: string
  email?: string
  email_verified?: boolean
}

// 上游资源条目字段名未在公开文档中完整定义，这里按常见命名做宽松映射。
export interface NexusmcResource {
  id?: string
  slug?: string
  page_url?: string
  title?: string
  name?: string
  description?: string
  excerpt?: string
  summary?: string
  cover?: string
  cover_image_url?: string
  image_url?: string
  thumbnail?: string
  downloads?: number
  download_count?: number
  views?: number
  view_count?: number
  tags?: (string | { name?: string; label?: string })[]
  updated_at?: string
  updated_at_time?: string
  updatedAt?: string
  [key: string]: unknown
}

export interface NexusmcPagination {
  page?: number
  pageSize?: number
  total?: number
  totalPages?: number
}

export interface NexusmcListResponse {
  items?: NexusmcResource[]
  pagination?: NexusmcPagination
  [key: string]: unknown
}

const api = axios.create({
  baseURL: '/api/nexusmc',
  timeout: 20000,
  headers: { Accept: 'application/json' },
})

export function isNexusmcConfigError(error: unknown): boolean {
  return axios.isAxiosError(error) && error.response?.status === 503
}

export const nexusmcApi = {
  config: () => api.get<{ configured: boolean; oauth_configured: boolean; auth_start: string }>('/config'),
  catalog: (platform?: string) => api.get('/catalog', { params: platform ? { platform } : {} }),
  resources: (params: Record<string, string | number | undefined>) => api.get<NexusmcListResponse>('/resources', { params }),
  resource: (id: string) => api.get('/resource', { params: { id } }),
  search: (q: string, params: Record<string, string | number | undefined> = {}) =>
    api.get<NexusmcListResponse>('/search', { params: { q, ...params } }),
}

// --- OAuth 用户数据（/me/*，需登录，token 由服务端 httpOnly cookie 持有） ----

export interface NexusmcMeListResponse {
  items?: Record<string, unknown>[]
  pagination?: NexusmcPagination
  [key: string]: unknown
}

export const nexusmcMeApi = {
  content: (params: Record<string, string | number | undefined> = {}) =>
    api.get<NexusmcMeListResponse>('/me/content', { params }),
  favorites: (params: Record<string, string | number | undefined> = {}) =>
    api.get<NexusmcMeListResponse>('/me/favorites', { params }),
  likes: (params: Record<string, string | number | undefined> = {}) =>
    api.get<NexusmcMeListResponse>('/me/likes', { params }),
  notifications: (params: Record<string, string | number | undefined> = {}) =>
    api.get<NexusmcMeListResponse>('/me/notifications', { params }),
  unreadCount: () => api.get<{ unread?: number; [key: string]: unknown }>('/me/notifications/unread-count'),
  markRead: (id: string | number) => api.post(`/me/notifications/${encodeURIComponent(String(id))}/read`),
  markAllRead: () => api.post('/me/notifications/read-all'),
}

// 401：token 缺失/过期，或授权未包含所需 scope，需要重新登录
export function isNexusmcAuthError(error: unknown): boolean {
  return axios.isAxiosError(error) && error.response?.status === 401
}

const USER_KEY = 'nexusmc_user'

export function getNexusmcUser(): NexusmcUser | null {
  try {
    const raw = localStorage.getItem(USER_KEY)
    return raw ? (JSON.parse(raw) as NexusmcUser) : null
  } catch {
    return null
  }
}

export function setNexusmcUser(user: NexusmcUser): void {
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export function clearNexusmcUser(): void {
  localStorage.removeItem(USER_KEY)
}

export const NEXUSMC_SITE = 'https://www.nexusmc.cn'
export const NEXUSMC_REGISTER_URL = `${NEXUSMC_SITE}/register`
export const NEXUSMC_LOGIN_URL = '/api/nexusmc/auth/start'

// 宽松取值工具：上游字段名存在多种可能。
export function pickString(item: NexusmcResource, keys: string[]): string {
  for (const k of keys) {
    const v = item[k]
    if (typeof v === 'string' && v) return v
  }
  return ''
}

export function pickNumber(item: NexusmcResource, keys: string[]): number | undefined {
  for (const k of keys) {
    const v = item[k]
    if (typeof v === 'number' && Number.isFinite(v)) return v
  }
  return undefined
}

export function resourceTags(item: NexusmcResource): string[] {
  const tags = item.tags
  if (!Array.isArray(tags)) return []
  return tags
    .map((t) => (typeof t === 'string' ? t : t?.name || t?.label || ''))
    .filter((t): t is string => Boolean(t))
    .slice(0, 4)
}

export function resourceTitle(item: NexusmcResource): string {
  return pickString(item, ['title', 'name']) || '未命名资源'
}

export function resourceDescription(item: NexusmcResource): string {
  return pickString(item, ['description', 'excerpt', 'summary'])
}

export function resourceCover(item: NexusmcResource): string {
  return pickString(item, ['cover', 'cover_image_url', 'coverImage', 'image_url', 'thumbnail'])
}

export function resourcePageUrl(item: NexusmcResource): string {
  return item.page_url || pickString(item, ['page_url', 'url']) || ''
}

export function resourceDownloads(item: NexusmcResource): number | undefined {
  return pickNumber(item, ['downloads', 'download_count'])
}

export function resourceViews(item: NexusmcResource): number | undefined {
  return pickNumber(item, ['views', 'view_count'])
}

export function resourceUpdatedAt(item: NexusmcResource): string {
  return pickString(item, ['updated_at', 'updatedAt', 'updated_at_time'])
}

// 详情页路由地址：优先 slug，退回 id；缺失时返回空串由调用方兜底。
export function resourceDetailPath(item: NexusmcResource): string {
  const id = item.slug || item.id
  return id ? `/nexusmc/resource/${encodeURIComponent(String(id))}` : ''
}

// --- 资源分类选项 -----------------------------------------------------------

export interface CatalogOption {
  value: string
  label: string
}

// 上游 catalog 返回结构不稳定：递归收集字符串/键值对形式的选项并去重。
export function normalizeCatalogOptions(data: unknown): CatalogOption[] {
  const out: CatalogOption[] = []
  const visit = (node: unknown) => {
    if (Array.isArray(node)) {
      for (const entry of node) {
        if (typeof entry === 'string') {
          out.push({ value: entry, label: entry })
        } else if (entry && typeof entry === 'object') {
          const obj = entry as Record<string, unknown>
          const value = (obj.value || obj.slug || obj.id || obj.key) as string | undefined
          const label = (obj.label || obj.name || obj.title || value) as string | undefined
          if (value && label) out.push({ value: String(value), label: String(label) })
        }
      }
      return
    }
    if (node && typeof node === 'object') {
      for (const value of Object.values(node as Record<string, unknown>)) {
        if (Array.isArray(value)) visit(value)
      }
    }
  }
  visit(data)
  const seen = new Set<string>()
  return out.filter((o) => (seen.has(o.value) ? false : (seen.add(o.value), true)))
}
