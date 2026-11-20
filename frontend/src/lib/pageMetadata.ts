import type { RouteLocationNormalized } from 'vue-router'

const SITE_NAME = 'YSM 模型站'
const SITE_ORIGIN = 'https://ysmmc.cn'
const HOME_DESCRIPTION = '浏览和搜索 Minecraft Yes Steve Model（YSM）角色模型，查看预览、作者介绍和资源版本，并前往 NexusMC 下载。'

function updateTag(selector: string, attributes: Record<string, string>) {
  const tag = document.head.querySelector(selector) ?? document.createElement(selector.startsWith('link') ? 'link' : 'meta')
  for (const [name, value] of Object.entries(attributes)) tag.setAttribute(name, value)
  if (!tag.isConnected) document.head.append(tag)
}

export function updatePageMetadata(
  route: RouteLocationNormalized,
  resource?: { title: string; description: string; cover: string },
  missing = false,
) {
  const isHome = route.name === 'home'
  const isResource = route.name === 'nexusmc-resource'
  const hasSearchFilters = isHome && Boolean(route.query.q || route.query.category)
  const indexable = (isHome || isResource) && !hasSearchFilters && !missing
  const page = typeof route.query.page === 'string' ? Number(route.query.page) : 1
  const isPaginated = isHome && Number.isSafeInteger(page) && page > 1

  const pageTitle = resource?.title || (isHome
    ? hasSearchFilters ? '模型搜索' : isPaginated ? `Minecraft YSM 模型下载 · 第 ${page} 页` : 'Minecraft YSM 模型下载'
    : isResource ? '模型详情'
    : route.name === 'me' ? '个人主页'
    : route.name === 'nexusmc-callback' ? '登录授权'
    : '页面未找到')
  const title = `${pageTitle} · ${SITE_NAME}`
  const description = resource
    ? resource.description.trim().replace(/\s+/g, ' ').slice(0, 160) || `查看 ${resource.title} 的模型预览，并前往 NexusMC 下载。`
    : HOME_DESCRIPTION
  const canonical = new URL(isHome ? '/' : route.path, SITE_ORIGIN)
  if (isPaginated && !hasSearchFilters) canonical.searchParams.set('page', String(page))

  let image = new URL('/icons/icon-512x512.png', SITE_ORIGIN).href
  if (resource?.cover) {
    try {
      const cover = new URL(resource.cover, SITE_ORIGIN)
      if (cover.protocol === 'https:' || cover.protocol === 'http:') image = cover.href
    } catch {
      // 使用站点图标。
    }
  }

  document.title = title
  updateTag('meta[name="description"]', { name: 'description', content: description })
  updateTag('meta[name="robots"]', { name: 'robots', content: indexable ? 'index, follow' : 'noindex, follow' })
  const canonicalTag = document.head.querySelector('link[rel="canonical"]')
  if (indexable) updateTag('link[rel="canonical"]', { rel: 'canonical', href: canonical.href })
  else canonicalTag?.remove()
  for (const [key, content] of Object.entries({ 'og:title': title, 'og:description': description, 'og:url': canonical.href, 'og:image': image })) {
    updateTag(`meta[property="${key}"]`, { property: key, content })
  }
  for (const [key, content] of Object.entries({ 'twitter:title': title, 'twitter:description': description, 'twitter:image': image })) {
    updateTag(`meta[name="${key}"]`, { name: key, content })
  }
}
