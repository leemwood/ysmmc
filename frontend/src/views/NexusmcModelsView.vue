<script setup lang="ts">
import axios from 'axios'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, ArrowUpRight, Box, CircleAlert, Clock3, Download, Eye, RefreshCw, Search, Upload, X } from 'lucide-vue-next'
import HeroArtwork from '@/components/HeroArtwork.vue'
import ResourceCover from '@/components/ResourceCover.vue'
import ResourcePagination from '@/components/ResourcePagination.vue'
import {
  isNexusmcConfigError,
  NEXUSMC_SITE,
  nexusmcApi,
  normalizeCatalogOptions,
  resourceCover,
  resourceDescription,
  resourceDetailPath,
  resourceDownloads,
  resourcePageUrl,
  resourceTags,
  resourceTitle,
  resourceViews,
  type CatalogOption,
  type NexusmcResource,
} from '@/lib/nexusmc'

const PAGE_SIZE = 20
const MAX_QUERY_LENGTH = 100
const SEARCH_HISTORY_KEY = 'ysmmc-recent-searches'
const sortOptions = { relevance: '相关度', latest: '最新发布', popular: '热度' } as const
const numberFormat = new Intl.NumberFormat('zh-CN')
const route = useRoute()
const router = useRouter()
const searchField = ref<HTMLInputElement>()
const resultsSection = ref<HTMLElement>()
const keyword = computed(() => typeof route.query.q === 'string' ? route.query.q.trim().slice(0, MAX_QUERY_LENGTH) : '')
const category = computed(() => typeof route.query.category === 'string' ? route.query.category : '')
const sort = computed(() => route.query.sort === 'latest' || route.query.sort === 'popular' ? route.query.sort : 'relevance')
const page = computed(() => {
  const value = typeof route.query.page === 'string' ? Number(route.query.page) : 1
  return Number.isSafeInteger(value) && value > 0 ? value : 1
})
const searchInput = ref(keyword.value)
const searchFocused = ref(false)
const recentSearches = ref<string[]>([])
const showHistory = computed(() => searchFocused.value && !searchInput.value.trim() && recentSearches.value.length > 0)
const compact = computed(() => Boolean(keyword.value || category.value || page.value > 1))
const loading = ref(true)
const error = ref<'offline' | 'search-unavailable' | 'load' | ''>('')
const items = ref<NexusmcResource[]>([])
const total = ref(0)
const totalPages = ref(1)
const categories = ref<CatalogOption[]>([])
let request: AbortController | undefined

const cards = computed(() => items.value.map((item) => {
  const downloads = resourceDownloads(item)
  const views = resourceViews(item)
  return {
    title: resourceTitle(item),
    cover: resourceCover(item),
    description: resourceDescription(item),
    tags: [...new Set(resourceTags(item))],
    path: resourceDetailPath(item),
    externalUrl: resourcePageUrl(item),
    downloads: downloads === undefined ? '' : numberFormat.format(downloads),
    views: views === undefined ? '' : numberFormat.format(views),
  }
}))

async function fetchCatalog() {
  try {
    const { data } = await nexusmcApi.catalog()
    categories.value = normalizeCatalogOptions(data)
  } catch {
    // 已加载的分类仍可用于筛选。
  }
}

async function fetchResources() {
  request?.abort()
  const currentRequest = new AbortController()
  request = currentRequest
  loading.value = true
  error.value = ''
  try {
    const params = { page: page.value, pageSize: PAGE_SIZE, category: category.value || undefined }
    const { data } = await (keyword.value
      ? nexusmcApi.search(keyword.value, { ...params, sort: sort.value }, currentRequest.signal)
      : nexusmcApi.resources(params, currentRequest.signal))
    if (currentRequest.signal.aborted) return

    const lastPage = Math.max(data.pagination?.totalPages ?? 1, 1)
    if (page.value > lastPage) {
      await router.replace({ query: { ...route.query, page: lastPage > 1 ? String(lastPage) : undefined } })
      return
    }
    items.value = data.items ?? []
    total.value = data.pagination?.total ?? items.value.length
    totalPages.value = lastPage
  } catch (cause) {
    if (currentRequest.signal.aborted) return
    const missingSearchScope = keyword.value && axios.isAxiosError(cause) &&
      cause.response?.status === 403 && cause.response.data?.code === 'SITE_API_KEY_SCOPE_REQUIRED'
    error.value = missingSearchScope ? 'search-unavailable' : isNexusmcConfigError(cause) ? 'offline' : 'load'
    items.value = []
    total.value = 0
    totalPages.value = 1
  } finally {
    if (!currentRequest.signal.aborted) loading.value = false
  }
}

async function updateQuery(
  q: string,
  nextCategory: string,
  nextPage = 1,
  scroll = false,
  nextSort: keyof typeof sortOptions = sort.value,
) {
  await router.push({
    path: '/',
    query: {
      ...route.query,
      q: q || undefined,
      category: nextCategory || undefined,
      page: nextPage > 1 ? String(nextPage) : undefined,
      sort: q && nextSort !== 'relevance' ? nextSort : undefined,
    },
  })
  if (scroll) resultsSection.value?.scrollIntoView({ block: 'start' })
}

function setHistory(entries: string[]) {
  recentSearches.value = entries
  try {
    window.localStorage.setItem(SEARCH_HISTORY_KEY, JSON.stringify(entries))
  } catch {
    // 当前会话仍可使用搜索记录。
  }
}

function submitSearch(value = searchInput.value) {
  const q = value.trim().slice(0, MAX_QUERY_LENGTH)
  searchInput.value = q
  if (q) setHistory([q, ...recentSearches.value.filter(term => term !== q)].slice(0, 5))
  searchField.value?.blur?.()
  searchFocused.value = false
  void updateQuery(q, category.value, 1, true, q === keyword.value ? sort.value : 'relevance')
}

function closeSearchOnExit(event: FocusEvent) {
  if (!(event.currentTarget as HTMLElement).contains(event.relatedTarget as Node | null)) {
    searchFocused.value = false
  }
}

function changeCategory(event: Event) {
  const select = event.target as HTMLSelectElement
  void updateQuery(keyword.value, select.value)
}

function changeSort(event: Event) {
  const select = event.target as HTMLSelectElement
  void updateQuery(keyword.value, category.value, 1, false, select.value as keyof typeof sortOptions)
}

function paginationLink(nextPage: number) {
  return {
    path: '/',
    query: { ...route.query, page: nextPage > 1 ? String(nextPage) : undefined },
    hash: '#resources',
  }
}

function refresh() {
  void fetchCatalog()
  void fetchResources()
}

function focusSearch() {
  if (route.hash === '#search') void nextTick(() => searchField.value?.focus({ preventScroll: true }))
}

watch([keyword, category, page, sort], fetchResources)
watch(keyword, (q) => { searchInput.value = q })
watch(() => route.hash, focusSearch)
onMounted(() => {
  try {
    const stored: unknown = JSON.parse(window.localStorage.getItem(SEARCH_HISTORY_KEY) ?? '[]')
    if (Array.isArray(stored)) {
      recentSearches.value = stored.filter((term): term is string =>
        typeof term === 'string' && Boolean(term.trim()) && term.length <= MAX_QUERY_LENGTH,
      ).slice(0, 5)
    }
  } catch {
    // 存储不可用或记录损坏时从空列表开始。
  }
  refresh()
  focusSearch()
})
onUnmounted(() => request?.abort())
</script>

<template>
  <div class="container-app">
    <section class="hero" :class="{ 'hero-compact': compact }" aria-labelledby="hero-title">
      <div class="hero-content" :class="{ 'enter-content': !compact }">
        <template v-if="!compact">
          <p class="hero-kicker"><Box class="size-4" aria-hidden="true" /> Minecraft · Yes Steve Model</p>
          <h1 id="hero-title">给方块世界<br /><span class="text-primary">换个新造型。</span></h1>
          <p class="hero-description">探索玩家创作的 YSM 角色模型，前往 NexusMC 下载并带进你的 Minecraft 世界。</p>
        </template>
        <h1 v-else id="hero-title" class="font-display text-2xl font-bold">YSM 模型库</h1>

        <form
          id="search"
          role="search"
          class="hero-search"
          @submit.prevent="submitSearch()"
          @focusout="closeSearchOnExit"
          @keydown.esc.stop="searchFocused = false"
        >
          <label for="model-search" class="sr-only">搜索模型名称</label>
          <Search class="ml-3 size-5 shrink-0 text-muted-foreground sm:ml-4" aria-hidden="true" />
          <input
            id="model-search"
            ref="searchField"
            v-model="searchInput"
            type="search"
            enterkeyhint="search"
            autocomplete="off"
            :maxlength="MAX_QUERY_LENGTH"
            placeholder="搜索模型名称"
            class="min-w-0 flex-1 self-stretch bg-transparent px-2 text-base sm:px-3"
            @focus="searchFocused = true"
            @input="searchFocused = true"
          />
          <button type="submit" class="site-button site-button-primary shrink-0">搜索</button>
          <div v-if="showHistory" class="site-panel enter-content search-history p-2 shadow-lg">
            <div class="flex items-center justify-between px-2 text-sm">
              <p class="font-semibold">最近搜索</p>
              <button type="button" class="site-button site-button-ghost text-primary" @click="setHistory([]); searchField?.focus()">清除记录</button>
            </div>
            <ul>
              <li v-for="term in recentSearches" :key="term">
                <button type="button" class="site-button site-button-ghost w-full justify-start text-left" @click="submitSearch(term)">
                  <Clock3 class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
                  <span class="truncate">{{ term }}</span>
                </button>
              </li>
            </ul>
          </div>
        </form>
      </div>
      <HeroArtwork v-if="!compact" class="enter-content" />
    </section>

    <section id="resources" ref="resultsSection" class="catalog" aria-labelledby="catalog-title">
      <div class="catalog-header">
        <div class="min-w-0">
          <h2 id="catalog-title" class="font-display text-[1.75rem] font-bold tracking-tight sm:text-[2rem]">模型资源</h2>
          <div class="mt-1.5 flex flex-wrap items-center gap-x-2.5 gap-y-1 text-sm text-muted-foreground" role="status">
            <span v-if="loading">正在加载模型…</span>
            <span v-else-if="error">查看 NexusMC 社区发布的 YSM 模型</span>
            <span v-else-if="keyword" class="break-all">「{{ keyword }}」的搜索结果 · {{ numberFormat.format(total) }} 个</span>
            <span v-else>共 {{ numberFormat.format(total) }} 个模型</span>
            <button
              v-if="keyword || category"
              type="button"
              class="inline-flex min-h-8 items-center gap-1 rounded-sm font-semibold text-primary hover:underline max-sm:min-h-11"
              @click="updateQuery('', '')"
            ><X class="size-3.5" aria-hidden="true" /> 清除筛选</button>
          </div>
        </div>

        <div class="catalog-controls flex flex-wrap items-center gap-2">
          <select
            v-if="categories.length || category"
            :value="category"
            class="site-select max-w-44"
            aria-label="资源分类"
            @change="changeCategory"
          >
            <option value="">全部分类</option>
            <option v-if="category && !categories.some(option => option.value === category)" :value="category">{{ category }}</option>
            <option v-for="option in categories" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
          <select v-if="keyword" :value="sort" class="site-select" aria-label="搜索排序" @change="changeSort">
            <option v-for="(label, value) in sortOptions" :key="value" :value="value">{{ label }}</option>
          </select>
          <button type="button" class="site-button site-button-outline site-button-icon" aria-label="刷新" :disabled="loading" @click="refresh">
            <RefreshCw class="size-4" :class="{ 'animate-spin': loading }" aria-hidden="true" />
          </button>
          <a
            :href="NEXUSMC_SITE + '/resources/new'"
            target="_blank"
            rel="noopener noreferrer"
            class="site-button site-button-outline"
            :class="{ 'max-[379px]:hidden!': keyword }"
            aria-label="我要投稿"
          >
            <Upload class="size-4" aria-hidden="true" /> 我要投稿
          </a>
        </div>
      </div>

      <div :aria-busy="loading">
        <div v-if="loading" class="resource-grid" aria-hidden="true">
          <div v-for="i in 4" :key="i" class="site-panel overflow-hidden">
            <div class="aspect-[16/10] animate-pulse bg-muted" />
            <div class="space-y-3 p-4">
              <div class="h-4 w-2/3 animate-pulse rounded-sm bg-muted" />
              <div class="h-3 w-full animate-pulse rounded-sm bg-muted" />
              <div class="h-3 w-2/5 animate-pulse rounded-sm bg-muted" />
            </div>
          </div>
        </div>

        <div v-else-if="error" class="site-panel flex flex-wrap items-center gap-4 p-5 sm:p-6" role="alert">
          <div class="status-icon"><CircleAlert class="size-[1.375rem]" aria-hidden="true" /></div>
          <div class="min-w-52 flex-1">
            <h3 class="font-bold">{{ error === 'offline' ? 'NexusMC 资源服务暂不可用' : error === 'search-unavailable' ? '站内搜索暂不可用' : '模型加载失败' }}</h3>
            <p class="mt-1 text-sm text-muted-foreground">{{ error === 'search-unavailable' ? '可以清除筛选继续浏览，或前往 NexusMC 搜索。' : '请稍后重试，或直接前往 NexusMC 浏览资源。' }}</p>
          </div>
          <a
            :href="keyword ? NEXUSMC_SITE + '/search?q=' + encodeURIComponent(keyword) : NEXUSMC_SITE + '/resources'"
            target="_blank"
            rel="noopener noreferrer"
            class="site-button site-button-outline"
          >
            {{ keyword ? '去 NexusMC 搜索' : '浏览 NexusMC' }} <ArrowUpRight class="size-4" aria-hidden="true" />
          </a>
        </div>

        <ul v-else-if="cards.length" class="resource-grid enter-content" aria-label="模型列表">
          <li v-for="(card, index) in cards" :key="card.path || card.externalUrl || index" class="min-w-0">
            <component
              :is="card.path ? 'RouterLink' : card.externalUrl ? 'a' : 'article'"
              :to="card.path || undefined"
              :href="!card.path ? card.externalUrl || undefined : undefined"
              :target="!card.path && card.externalUrl ? '_blank' : undefined"
              :rel="!card.path && card.externalUrl ? 'noopener noreferrer' : undefined"
              class="resource-card site-panel group flex h-full flex-col overflow-hidden"
            >
              <div class="aspect-[16/10] overflow-hidden border-b border-border bg-muted">
                <ResourceCover :cover="card.cover" :title="card.title" />
              </div>
              <div class="flex flex-1 flex-col p-4">
                <h3 class="line-clamp-2 text-[1.0625rem] font-bold group-hover:text-primary">{{ card.title }}</h3>
                <p v-if="card.description" class="mt-1.5 line-clamp-2 text-sm text-muted-foreground">{{ card.description }}</p>
                <div v-if="card.tags.length" class="mt-3 flex flex-wrap gap-1.5">
                  <span v-for="tag in card.tags" :key="tag" class="max-w-full truncate rounded-md bg-accent px-2 py-0.5 text-xs font-medium text-primary">{{ tag }}</span>
                </div>
                <div v-if="card.downloads || card.views || card.path || card.externalUrl" class="mt-auto flex flex-wrap items-center gap-x-3 gap-y-1 pt-4 text-xs text-muted-foreground">
                  <span v-if="card.downloads" class="inline-flex items-center gap-1"><Download class="size-3.5" aria-label="下载次数" />{{ card.downloads }}</span>
                  <span v-if="card.views" class="inline-flex items-center gap-1"><Eye class="size-3.5" aria-label="浏览次数" />{{ card.views }}</span>
                  <ArrowRight v-if="card.path" class="ml-auto size-4 text-primary" aria-hidden="true" />
                  <ArrowUpRight v-else-if="card.externalUrl" class="ml-auto size-4 text-primary" aria-hidden="true" />
                </div>
              </div>
            </component>
          </li>
        </ul>

        <div v-else class="site-panel flex flex-col items-center gap-3 px-5 py-9 text-center">
          <div class="status-icon"><Box class="size-6" aria-hidden="true" /></div>
          <h3 class="font-bold">{{ keyword || category ? '没有找到符合条件的模型' : '暂无 YSM 模型资源' }}</h3>
          <p class="text-sm text-muted-foreground">{{ keyword || category ? '试试其他关键词，或清除筛选后重新浏览。' : 'NexusMC 审核通过的新模型会展示在这里。' }}</p>
        </div>

        <ResourcePagination
          v-if="!loading && !error"
          :page="page"
          :total-pages="totalPages"
          :to="paginationLink"
        />
      </div>
    </section>
  </div>
</template>

<style scoped>
.hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  align-items: center;
  gap: clamp(1.75rem, 5vw, 4.25rem);
  padding-block: clamp(1.75rem, 5vw, 3.5rem);
}

.hero-kicker {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--color-primary);
  font-size: 0.8125rem;
  font-weight: 700;
}

.hero-content {
  position: relative;
  z-index: 1;
}

.hero h1 {
  margin-top: 0.75rem;
  font-family: var(--font-display);
  font-size: clamp(2.65rem, 5.3vw, 4.75rem);
  font-weight: 700;
  line-height: 1.12;
  letter-spacing: -0.06em;
}

.hero-description {
  max-width: 27rem;
  margin-top: 1rem;
  color: var(--color-muted-foreground);
  font-size: clamp(0.9375rem, 1.45vw, 1.0625rem);
  text-wrap: pretty;
}

.hero-search {
  position: relative;
  display: flex;
  max-width: 29rem;
  min-height: 3.375rem;
  align-items: center;
  margin-top: 1.5rem;
  padding: 0.3125rem;
  border: 1px solid var(--color-border);
  border-radius: 0.875rem;
  background: var(--color-card);
}

.search-history {
  position: absolute;
  z-index: 10;
  inset: calc(100% + 0.375rem) 0 auto;
}

.hero-search:focus-within {
  border-color: var(--color-primary);
  outline: 1px solid var(--color-primary);
}

.hero-search input:focus-visible {
  outline: none;
}

.hero-search .site-button {
  min-width: 4.875rem;
}

.catalog {
  border-top: 1px solid var(--color-border);
  padding-block: 1.75rem 3.75rem;
}

.catalog-header {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 0.875rem 1.25rem;
  margin-bottom: 1.25rem;
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 16rem), 1fr));
  gap: 1rem;
}

.resource-card:is(a):hover {
  border-color: var(--color-primary);
}

.resource-card {
  transition: border-color 180ms, transform 180ms, box-shadow 180ms;
}

.resource-card:active {
  transform: scale(0.985);
}

@media (hover: hover) {
  .resource-card:is(a):hover {
    transform: translateY(-3px);
    box-shadow: 0 6px 18px hsl(var(--foreground) / 0.08);
  }
}

@media (min-width: 800px) {
  .hero {
    grid-template-columns: minmax(0, 0.88fr) minmax(0, 1.12fr);
  }
}

@media (max-width: 479px) {
  .hero {
    gap: 1rem;
    padding-block: 1.125rem;
  }

  .hero :deep(.hero-art) {
    max-height: 10rem;
  }

  .hero :deep(canvas) {
    height: 100%;
    max-height: 10rem;
    object-fit: contain;
  }

  .hero h1 {
    font-size: clamp(2rem, 9vw, 2.75rem);
  }

  .hero-description {
    margin-top: 0.625rem;
  }

  .hero-search {
    margin-top: 0.875rem;
  }

  .catalog {
    padding-top: 1.25rem;
    padding-bottom: 2.25rem;
  }

  .resource-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0.625rem;
  }

  .resource-card > div:last-child {
    padding: 0.75rem;
  }

  .catalog-controls {
    width: 100%;
  }

  .catalog-controls .site-select {
    min-width: 0;
    flex: 1 1 6rem;
  }
}

@media (max-width: 359px) {
  .resource-grid {
    grid-template-columns: 1fr;
  }
}

.hero.hero-compact {
  grid-template-columns: minmax(0, 1fr);
  padding-block: 1.125rem;
}

.hero-compact h1 {
  margin-top: 0;
  font-size: 1.5rem;
  letter-spacing: -0.02em;
}

.hero-compact .hero-search {
  margin-top: 0.625rem;
}
</style>
