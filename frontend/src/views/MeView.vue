<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ArrowUpRight, Bell, CheckCheck, CircleAlert, Heart, LogIn, LogOut, Package, UserRound } from 'lucide-vue-next'
import type { AxiosResponse } from 'axios'
import NexusmcAvatar from '@/components/NexusmcAvatar.vue'
import ResourceCover from '@/components/ResourceCover.vue'
import ResourcePagination from '@/components/ResourcePagination.vue'
import {
  isNexusmcAuthError,
  isNexusmcConfigError,
  NEXUSMC_REGISTER_URL,
  NEXUSMC_SITE,
  nexusmcMeApi,
  resourceCover,
  resourcePageUrl,
  resourceTitle,
  resourceUpdatedAt,
  type NexusmcMeListResponse,
  type NexusmcResource,
} from '@/lib/nexusmc'
import { useNexusmcStore } from '@/stores/nexusmc'

type LooseItem = Record<string, unknown>
type MeTab = 'content' | 'favorites' | 'notifications'

const PAGE_SIZE = 20
const tabs = [
  { key: 'content', label: '我的投稿', icon: Package, empty: '当前没有符合条件的投稿。' },
  { key: 'favorites', label: '我的收藏', icon: Heart, empty: '还没有收藏过模型。' },
  { key: 'notifications', label: '通知', icon: Bell, empty: '暂无通知。' },
] as const

const statusLabels: Record<string, string> = {
  approved: '已通过',
  pending: '审核中',
  rejected: '未通过',
}

const store = useNexusmcStore()
const activeTab = ref<MeTab>('content')
const contentStatus = ref('all')
const unread = ref(0)
const markingRead = ref(false)
const readError = ref('')
let unreadRequest = 0

function asRecord(value: unknown): LooseItem | null {
  return value && typeof value === 'object' && !Array.isArray(value) ? value as LooseItem : null
}

function unwrapResource(item: LooseItem): NexusmcResource {
  return asRecord(item.resource) ?? asRecord(item.item) ?? item
}

function looseString(item: LooseItem, keys: string[]): string {
  for (const key of keys) {
    const value = item[key]
    if (typeof value === 'string' && value) return value
  }
  return ''
}

function looseId(item: LooseItem): string {
  return looseString(item, ['id', 'slug', 'nid']) || (typeof item.id === 'number' ? String(item.id) : '')
}

function formatDate(raw: string): string {
  const date = new Date(raw)
  return Number.isNaN(date.getTime()) ? '' : date.toLocaleDateString('zh-CN')
}

function isUnread(item: LooseItem): boolean {
  return !(item.isRead ?? item.read ?? item.readAt ?? item.read_at)
}

function createListState(loader: (page: number, pageSize: number) => Promise<AxiosResponse<NexusmcMeListResponse>>) {
  const items = ref<LooseItem[]>([])
  const loading = ref(false)
  const error = ref('')
  const authError = ref(false)
  const page = ref(1)
  const totalPages = ref(1)
  const loaded = ref(false)
  let requestId = 0

  async function fetchList(targetPage = 1) {
    const currentRequest = ++requestId
    loading.value = true
    error.value = ''
    authError.value = false
    page.value = targetPage

    try {
      const { data } = await loader(targetPage, PAGE_SIZE)
      if (currentRequest !== requestId) return
      items.value = data.items ?? []
      totalPages.value = Math.max(data.pagination?.totalPages ?? 1, 1)
      loaded.value = true
    } catch (cause) {
      if (currentRequest !== requestId) return
      authError.value = isNexusmcAuthError(cause)
      if (isNexusmcConfigError(cause)) {
        error.value = 'NexusMC 服务暂不可用，请稍后重试。'
      } else if (authError.value) {
        error.value = '登录已过期，或当前授权缺少个人数据权限。'
      } else {
        error.value = '列表加载失败，请稍后重试。'
      }
      items.value = []
      loaded.value = false
    } finally {
      if (currentRequest === requestId) loading.value = false
    }
  }

  return { items, loading, error, authError, page, totalPages, loaded, fetchList }
}

const contentState = createListState((page, pageSize) =>
  nexusmcMeApi.content({
    type: 'resource',
    status: contentStatus.value === 'all' ? undefined : contentStatus.value,
    page,
    pageSize,
  }),
)
const favoritesState = createListState((page, pageSize) => nexusmcMeApi.favorites({ page, pageSize }))
const notificationsState = createListState((page, pageSize) => nexusmcMeApi.notifications({ page, pageSize }))
const states = reactive({ content: contentState, favorites: favoritesState, notifications: notificationsState })
const activeState = computed(() => states[activeTab.value])
const selectedTab = computed(() => tabs.find(tab => tab.key === activeTab.value) ?? tabs[0])

const resources = computed(() => activeState.value.items.map((item) => {
  const resource = unwrapResource(item)
  const status = looseString(item, ['status']) || looseString(resource, ['status'])
  return {
    id: looseId(item) || looseId(resource),
    title: resourceTitle(resource),
    cover: resourceCover(resource),
    updatedAt: formatDate(resourceUpdatedAt(resource)),
    url: resourcePageUrl(resource),
    status: statusLabels[status] ?? '',
  }
}))

function switchTab(tab: MeTab) {
  activeTab.value = tab
  const state = states[tab]
  if (!state.loaded && !state.loading) void state.fetchList(state.page)
}

async function refreshUnread() {
  const currentRequest = ++unreadRequest
  try {
    const { data } = await nexusmcMeApi.unreadCount()
    if (currentRequest === unreadRequest) unread.value = typeof data.count === 'number' ? data.count : 0
  } catch {
    // 保留最近一次已知的未读数。
  }
}

async function readNotification(item?: LooseItem) {
  if (markingRead.value || (item && !isUnread(item))) return
  const id = item ? looseId(item) : ''
  markingRead.value = true
  readError.value = ''
  try {
    if (!item) {
      await nexusmcMeApi.markAllRead()
      for (const notification of notificationsState.items.value) notification.isRead = true
      unread.value = 0
    } else {
      if (id) await nexusmcMeApi.markRead(id)
      item.isRead = true
      unread.value = Math.max(unread.value - 1, 0)
    }
    ++unreadRequest
    if (!item || id) void refreshUnread()
  } catch {
    readError.value = '标记已读失败，请重试。'
  } finally {
    markingRead.value = false
  }
}

onMounted(() => {
  if (!store.isLoggedIn) return
  void contentState.fetchList()
  void refreshUnread()
})
</script>

<template>
  <div class="container-app py-7 sm:py-11">
    <section v-if="!store.isLoggedIn" class="site-panel mx-auto my-6 max-w-lg p-6 text-center sm:my-10 sm:p-10">
      <div class="status-icon mx-auto"><UserRound class="size-6" aria-hidden="true" /></div>
      <h1 class="mt-4 text-2xl font-bold">登录后查看个人主页</h1>
      <p class="mt-2 text-sm text-muted-foreground">使用 NexusMC 账号查看你的投稿、收藏和通知。在授权时开启「个人数据读取」权限即可访问。</p>
      <div class="mt-6 flex flex-wrap justify-center gap-2">
        <button type="button" class="site-button site-button-primary" @click="store.login('/me')">
          <LogIn class="size-4" aria-hidden="true" /> 通过 NexusMC 登录
        </button>
        <a :href="NEXUSMC_REGISTER_URL" target="_blank" rel="noopener noreferrer" class="site-button site-button-outline">
          注册账号 <ArrowUpRight class="size-4" aria-hidden="true" />
        </a>
      </div>
      <RouterLink to="/" class="mt-5 inline-block rounded-sm text-sm font-semibold text-primary hover:underline">返回模型库</RouterLink>
    </section>

    <template v-else>
      <h1 class="mb-5 font-display text-[clamp(1.875rem,4vw,2.75rem)] font-bold tracking-tight">个人主页</h1>
      <div class="account-layout">
        <aside class="site-panel overflow-hidden">
          <div class="flex items-center gap-3 px-4 py-4 sm:px-5">
            <NexusmcAvatar size="lg" class="shrink-0" />
            <div class="min-w-0 flex-1">
              <p class="truncate font-bold">{{ store.displayName }}</p>
              <p v-if="store.user?.email || store.user?.uid != null" class="truncate text-xs text-muted-foreground">
                {{ store.user?.email || `NexusMC UID ${store.user?.uid}` }}
              </p>
            </div>
            <button type="button" class="site-button site-button-ghost site-button-icon shrink-0" aria-label="退出登录" @click="store.logout()">
              <LogOut class="size-4 text-muted-foreground" aria-hidden="true" />
            </button>
          </div>

          <nav class="account-nav border-t border-border p-2" aria-label="个人主页导航">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              type="button"
              :aria-current="activeTab === tab.key ? 'true' : undefined"
              class="site-button site-button-ghost min-w-0 justify-start px-2 sm:px-3"
              :class="activeTab === tab.key ? 'bg-accent text-primary' : 'text-muted-foreground'"
              @click="switchTab(tab.key)"
            >
              <component :is="tab.icon" class="size-4 shrink-0" aria-hidden="true" />
              {{ tab.label }}
              <span
                v-if="tab.key === 'notifications' && unread > 0"
                class="inline-flex min-w-5 items-center justify-center rounded-full bg-primary px-1 text-xs text-primary-foreground"
                :aria-label="`${unread} 条未读`"
              >{{ unread }}</span>
            </button>
          </nav>
        </aside>

        <section class="min-w-0" aria-labelledby="me-list-title" :aria-busy="activeState.loading">
          <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
            <h2 id="me-list-title" class="text-xl font-bold">{{ selectedTab.label }}</h2>
            <select
              v-if="activeTab === 'content'"
              v-model="contentStatus"
              class="site-select"
              aria-label="投稿状态"
              @change="contentState.fetchList()"
            >
              <option value="all">全部状态</option>
              <option v-for="(label, value) in statusLabels" :key="value" :value="value">{{ label }}</option>
            </select>
            <button
              v-if="activeTab === 'notifications' && unread > 0"
              type="button"
              class="site-button site-button-outline"
              :disabled="markingRead || activeState.loading"
              @click="readNotification()"
            ><CheckCheck class="size-4" aria-hidden="true" /> 全部标为已读</button>
          </div>

          <p v-if="activeTab === 'notifications' && readError" role="alert" class="mb-3 text-sm text-destructive">{{ readError }}</p>

          <div v-if="activeState.loading" class="grid gap-3" role="status">
            <span class="sr-only">正在加载{{ selectedTab.label }}…</span>
            <div v-for="i in 3" :key="i" class="h-20 animate-pulse rounded-xl bg-muted" aria-hidden="true" />
          </div>

          <div v-else-if="activeState.error" class="site-panel flex flex-wrap items-center gap-3 p-4 sm:p-5" role="alert">
            <div class="status-icon"><CircleAlert class="size-5" aria-hidden="true" /></div>
            <p class="min-w-44 flex-1 text-sm text-muted-foreground">{{ activeState.error }}</p>
            <button v-if="activeState.authError" type="button" class="site-button site-button-outline" @click="store.login('/me')">重新登录</button>
            <button v-else type="button" class="site-button site-button-outline" @click="activeState.fetchList(activeState.page)">重试</button>
          </div>

          <ul v-else-if="activeTab === 'notifications' && activeState.items.length" class="grid gap-2">
            <li v-for="(item, index) in activeState.items" :key="looseId(item) || index">
              <article class="site-panel flex flex-wrap items-start gap-3 p-4">
                <span
                  class="mt-2 size-2 shrink-0 rounded-full"
                  :class="isUnread(item) ? 'bg-primary' : 'bg-muted'"
                  :aria-label="isUnread(item) ? '未读' : '已读'"
                />
                <div class="min-w-36 flex-1 wrap-anywhere">
                  <h3 class="text-sm font-bold">{{ looseString(item, ['title', 'subject']) || '通知' }}</h3>
                  <p class="mt-1 whitespace-pre-line text-sm text-muted-foreground">{{ looseString(item, ['content', 'message', 'body', 'excerpt']) }}</p>
                  <p class="mt-1 text-xs text-muted-foreground">{{ formatDate(looseString(item, ['createdAt', 'created_at', 'time'])) }}</p>
                </div>
                <button
                  v-if="isUnread(item)"
                  type="button"
                  class="site-button site-button-ghost text-primary"
                  :disabled="markingRead"
                  @click="readNotification(item)"
                >标为已读</button>
              </article>
            </li>
          </ul>

          <ul v-else-if="resources.length" class="grid gap-2">
            <li v-for="(resource, index) in resources" :key="resource.id || index" class="site-panel flex flex-wrap items-center gap-3 p-3 sm:p-4">
              <div class="hidden h-16 w-24 shrink-0 overflow-hidden rounded-md bg-muted sm:block">
                <ResourceCover :cover="resource.cover" :title="resource.title" icon-class="size-5" />
              </div>
              <div class="min-w-28 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <h3 class="font-bold">{{ resource.title }}</h3>
                  <span
                    v-if="activeTab === 'content' && resource.status"
                    class="rounded-md bg-accent px-2 py-0.5 text-xs font-medium text-primary"
                  >{{ resource.status }}</span>
                </div>
                <p v-if="resource.updatedAt" class="mt-1 text-xs text-muted-foreground">更新于 {{ resource.updatedAt }}</p>
              </div>
              <a v-if="resource.url" :href="resource.url" target="_blank" rel="noopener noreferrer" class="site-button site-button-outline">
                NexusMC <ArrowUpRight class="size-4" aria-hidden="true" />
              </a>
            </li>
          </ul>

          <div v-else class="site-panel grid justify-items-center gap-3 p-6 text-center sm:p-10">
            <div class="status-icon"><component :is="selectedTab.icon" class="size-5" aria-hidden="true" /></div>
            <p class="text-sm text-muted-foreground">{{ selectedTab.empty }}</p>
            <a v-if="activeTab === 'content'" :href="NEXUSMC_SITE + '/resources/new'" target="_blank" rel="noopener noreferrer" class="site-button site-button-outline">
              前往投稿 <ArrowUpRight class="size-4" aria-hidden="true" />
            </a>
            <RouterLink v-else-if="activeTab === 'favorites'" to="/" class="site-button site-button-outline">浏览模型库</RouterLink>
          </div>

          <ResourcePagination
            v-if="!activeState.loading && !activeState.error"
            :page="activeState.page"
            :total-pages="activeState.totalPages"
            @change="activeState.fetchList($event)"
          />
        </section>
      </div>
    </template>
  </div>
</template>

<style scoped>
.account-layout {
  display: grid;
  align-items: start;
  gap: 1.25rem;
}

.account-nav {
  display: flex;
  flex-wrap: wrap;
  gap: 0.125rem;
}

@media (min-width: 800px) {
  .account-layout {
    grid-template-columns: minmax(15rem, 1fr) minmax(0, 3fr);
    gap: 1.5rem;
  }

  .account-nav {
    display: grid;
  }
}
</style>
