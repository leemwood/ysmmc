<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Skeleton } from '@/components/ui/skeleton'
import NexusmcAvatar from '@/components/NexusmcAvatar.vue'
import ResourceCover from '@/components/ResourceCover.vue'
import {
  Bell,
  CheckCheck,
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  Heart,
  LogOut,
  Package,
  Sparkles,
} from 'lucide-vue-next'
import {
  isNexusmcAuthError,
  nexusmcMeApi,
  resourceCover,
  resourcePageUrl,
  resourceTitle,
  resourceUpdatedAt,
  type NexusmcMeListResponse,
  type NexusmcResource,
} from '@/lib/nexusmc'
import { useNexusmcStore } from '@/stores/nexusmc'
import type { AxiosResponse } from 'axios'

type LooseItem = Record<string, unknown>
type MeTab = 'content' | 'favorites' | 'notifications'

const PAGE_SIZE = 20

const store = useNexusmcStore()
const activeTab = ref<MeTab>('content')
const contentStatus = ref('all')
const unread = ref(0)

// --- 宽松字段工具：me/* 各端点条目结构未在文档中完整定义 ---------------

function asRecord(v: unknown): LooseItem | null {
  return v && typeof v === 'object' && !Array.isArray(v) ? (v as LooseItem) : null
}

// 投稿/收藏条目可能是裸资源，也可能包一层 { resource: ... } / { item: ... }
function unwrapResource(item: LooseItem): NexusmcResource {
  return (asRecord(item.resource) ?? asRecord(item.item) ?? item) as NexusmcResource
}

function looseString(obj: LooseItem, keys: string[]): string {
  for (const k of keys) {
    const v = obj[k]
    if (typeof v === 'string' && v) return v
  }
  return ''
}

function looseId(item: LooseItem): string {
  return looseString(item, ['id', 'slug', 'nid'])
}

const STATUS_LABELS: Record<string, { label: string; variant: 'secondary' | 'outline' | 'destructive' }> = {
  approved: { label: '已通过', variant: 'secondary' },
  pending: { label: '审核中', variant: 'outline' },
  rejected: { label: '未通过', variant: 'destructive' },
}
function contentStatusOf(item: LooseItem) {
  return STATUS_LABELS[looseString(item, ['status'])] ?? STATUS_LABELS.approved
}

function formatDate(raw: string): string {
  if (!raw) return ''
  const date = new Date(raw)
  return Number.isNaN(date.getTime()) ? '' : date.toLocaleDateString('zh-CN')
}

function notificationTitle(n: LooseItem): string {
  return looseString(n, ['title', 'subject']) || '通知'
}

function notificationBody(n: LooseItem): string {
  return looseString(n, ['content', 'message', 'body', 'excerpt'])
}

function notificationTime(n: LooseItem): string {
  return formatDate(looseString(n, ['createdAt', 'created_at', 'time']))
}

function isUnread(n: LooseItem): boolean {
  return !(n.isRead ?? n.read ?? n.readAt ?? n.read_at)
}

// --- 三个列表共享的分页加载逻辑 -----------------------------------------

function createListState(loader: (page: number, pageSize: number) => Promise<AxiosResponse<NexusmcMeListResponse>>) {
  const items = ref<LooseItem[]>([])
  const loading = ref(false)
  const error = ref('')
  const page = ref(1)
  const totalPages = ref(1)
  const loaded = ref(false)

  async function fetchList(targetPage = 1) {
    loading.value = true
    error.value = ''
    try {
      const { data } = await loader(targetPage, PAGE_SIZE)
      items.value = data.items ?? []
      totalPages.value = Math.max(data.pagination?.totalPages ?? 1, 1)
      page.value = targetPage
      loaded.value = true
    } catch (e) {
      error.value = isNexusmcAuthError(e)
        ? '登录已过期，或授权未包含个人数据权限，请重新登录。'
        : '加载失败，请稍后重试。'
      items.value = []
    } finally {
      loading.value = false
    }
  }

  return { items, loading, error, page, totalPages, loaded, fetchList }
}

const contentState = createListState((p, ps) =>
  nexusmcMeApi.content({
    type: 'resource',
    status: contentStatus.value === 'all' ? undefined : contentStatus.value,
    page: p,
    pageSize: ps,
  }),
)
const favoritesState = createListState((p, ps) => nexusmcMeApi.favorites({ page: p, pageSize: ps }))
const notificationsState = createListState((p, ps) => nexusmcMeApi.notifications({ page: p, pageSize: ps }))

const states = reactive<Record<MeTab, ReturnType<typeof createListState>>>({
  content: contentState,
  favorites: favoritesState,
  notifications: notificationsState,
})

const activeState = computed(() => states[activeTab.value])

function switchTab(tab: MeTab) {
  activeTab.value = tab
  const state = states[tab]
  if (!state.loaded) void state.fetchList(1)
}

function goPage(delta: number) {
  const state = activeState.value
  const next = state.page + delta
  if (next < 1 || next > state.totalPages) return
  void state.fetchList(next)
}

function changeStatus() {
  void contentState.fetchList(1)
}

async function refreshUnread() {
  try {
    const { data } = await nexusmcMeApi.unreadCount()
    unread.value = typeof data.unread === 'number' ? data.unread : 0
  } catch {
    unread.value = 0
  }
}

async function readAllNotifications() {
  try {
    await nexusmcMeApi.markAllRead()
    unread.value = 0
    for (const n of notificationsState.items.value) {
      n.isRead = true
    }
  } catch {
    // 静默失败，下次刷新时恢复
  }
}

async function readNotification(n: LooseItem) {
  if (!isUnread(n)) return
  const id = looseId(n)
  if (!id) {
    n.isRead = true
    unread.value = Math.max(unread.value - 1, 0)
    return
  }
  try {
    await nexusmcMeApi.markRead(id)
    n.isRead = true
    unread.value = Math.max(unread.value - 1, 0)
  } catch {
    // 保持未读，不打扰用户
  }
}

onMounted(() => {
  if (!store.isLoggedIn) return
  void contentState.fetchList(1)
  void refreshUnread()
})
</script>

<template>
  <div class="container-app py-6 sm:py-8">
    <!-- 未登录：引导登录（登录后自动回本页） -->
    <Card v-if="!store.isLoggedIn" class="mx-auto mt-10 w-full max-w-md">
      <CardContent class="flex flex-col items-center gap-4 p-8 text-center">
        <p class="font-semibold">登录后查看个人主页</p>
        <p class="text-sm text-muted-foreground">
          个人主页展示你在 NexusMC 的投稿、收藏与通知。首次登录需在 NexusMC 授权页
          同意「个人数据读取」权限。
        </p>
        <Button class="btn-press focus-ring" @click="store.login('/me')">
          <Sparkles class="mr-1.5 h-4 w-4" />
          通过 NexusMC 登录
        </Button>
        <RouterLink to="/" class="text-sm text-primary hover:underline">返回资源首页</RouterLink>
      </CardContent>
    </Card>

    <template v-else>
      <!-- 资料卡 -->
      <section class="flex items-center gap-4 rounded-2xl border bg-background/80 p-5 sm:p-6">
        <NexusmcAvatar size="lg" class="border" />
        <div class="min-w-0 flex-1">
          <h1 class="truncate text-xl font-bold tracking-tight">{{ store.displayName }}</h1>
          <p class="truncate text-sm text-muted-foreground">
            {{ store.user?.email || `NexusMC UID ${store.user?.uid ?? '—'}` }}
          </p>
        </div>
        <Button variant="ghost" size="icon" class="btn-press focus-ring shrink-0" aria-label="退出登录" @click="store.logout()">
          <LogOut class="h-4 w-4" />
        </Button>
      </section>

      <!-- 标签页 -->
      <div class="mt-6 flex items-center gap-1 border-b">
        <button
          v-for="tab in [
            { key: 'content', label: '我的投稿', icon: Package },
            { key: 'favorites', label: '我的收藏', icon: Heart },
            { key: 'notifications', label: '通知', icon: Bell },
          ] as const"
          :key="tab.key"
          class="relative inline-flex items-center gap-1.5 rounded-t-md px-3 py-2 text-sm font-medium transition-colors focus-ring"
          :class="activeTab === tab.key ? 'text-primary border-b-2 border-primary -mb-px' : 'text-muted-foreground hover:text-foreground'"
          @click="switchTab(tab.key)"
        >
          <component :is="tab.icon" class="h-4 w-4" />
          {{ tab.label }}
          <span
            v-if="tab.key === 'notifications' && unread > 0"
            class="ml-0.5 rounded-full bg-destructive px-1.5 text-xs font-medium text-destructive-foreground"
          >
            {{ unread }}
          </span>
        </button>
      </div>

      <!-- 投稿状态筛选 -->
      <div v-if="activeTab === 'content'" class="mt-4 flex items-center gap-2">
        <select
          v-model="contentStatus"
          class="h-9 rounded-md border border-input bg-background px-3 text-sm focus-ring"
          @change="changeStatus"
        >
          <option value="all">全部状态</option>
          <option value="approved">已通过</option>
          <option value="pending">审核中</option>
          <option value="rejected">未通过</option>
        </select>
      </div>

      <!-- 通知：全部已读 -->
      <div v-if="activeTab === 'notifications' && unread > 0" class="mt-4 flex justify-end">
        <Button variant="outline" size="sm" class="btn-press focus-ring" @click="readAllNotifications">
          <CheckCheck class="mr-1.5 h-4 w-4" />
          全部标为已读
        </Button>
      </div>

      <!-- 错误提示 -->
      <Alert v-if="activeState.error" class="mt-4" :variant="activeState.error.includes('登录') ? 'default' : 'destructive'">
        <AlertDescription class="flex flex-wrap items-center gap-3">
          <span>{{ activeState.error }}</span>
          <Button v-if="activeState.error.includes('登录')" size="sm" variant="outline" class="btn-press focus-ring" @click="store.login('/me')">
            重新登录
          </Button>
        </AlertDescription>
      </Alert>

      <!-- 加载骨架 -->
      <div v-if="activeState.loading" class="mt-4 space-y-3">
        <Skeleton v-for="i in 4" :key="i" class="h-20 rounded-xl" />
      </div>

      <!-- 投稿 / 收藏列表 -->
      <template v-else-if="activeTab !== 'notifications' && activeState.items.length">
        <div class="mt-4 space-y-3">
          <Card v-for="item in activeState.items" :key="looseId(item) || resourceTitle(unwrapResource(item))" class="card-hover">
            <CardContent class="flex items-center gap-4 p-3 sm:p-4">
              <div class="hidden h-16 w-24 shrink-0 overflow-hidden rounded-md bg-muted sm:block">
                <ResourceCover
                  :cover="resourceCover(unwrapResource(item))"
                  :title="resourceTitle(unwrapResource(item))"
                  icon-class="h-6 w-6"
                />
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <span class="truncate font-medium">{{ resourceTitle(unwrapResource(item)) }}</span>
                  <Badge v-if="activeTab === 'content'" :variant="contentStatusOf(item).variant" class="text-xs">
                    {{ contentStatusOf(item).label }}
                  </Badge>
                </div>
                <p class="mt-1 text-xs text-muted-foreground">
                  <template v-if="formatDate(resourceUpdatedAt(unwrapResource(item)))">
                    更新于 {{ formatDate(resourceUpdatedAt(unwrapResource(item))) }}
                  </template>
                  <template v-else>更新日期待定</template>
                </p>
              </div>
              <a
                v-if="resourcePageUrl(unwrapResource(item))"
                :href="resourcePageUrl(unwrapResource(item))"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button variant="outline" size="sm" class="btn-press focus-ring shrink-0">
                  <ExternalLink class="mr-1.5 h-3.5 w-3.5" />
                  NexusMC
                </Button>
              </a>
            </CardContent>
          </Card>
        </div>
      </template>

      <!-- 通知列表 -->
      <template v-else-if="activeTab === 'notifications' && activeState.items.length">
        <div class="mt-4 space-y-2">
          <button
            v-for="n in activeState.items"
            :key="looseId(n) || notificationTitle(n)"
            class="w-full rounded-xl border p-4 text-left transition-colors hover:bg-accent/50 focus-ring"
            :class="isUnread(n) ? 'bg-background' : 'opacity-70'"
            @click="readNotification(n)"
          >
            <div class="flex items-start gap-3">
              <span v-if="isUnread(n)" class="mt-1.5 h-2 w-2 shrink-0 rounded-full bg-primary" />
              <span v-else class="mt-1.5 h-2 w-2 shrink-0" />
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium" :class="{ 'font-semibold': isUnread(n) }">{{ notificationTitle(n) }}</p>
                <p v-if="notificationBody(n)" class="mt-1 text-sm text-muted-foreground line-clamp-2">
                  {{ notificationBody(n) }}
                </p>
                <p v-if="notificationTime(n)" class="mt-1 text-xs text-muted-foreground/70">{{ notificationTime(n) }}</p>
              </div>
            </div>
          </button>
        </div>
      </template>

      <!-- 空状态 -->
      <div v-else-if="!activeState.error" class="mt-12 text-center">
        <p class="text-muted-foreground">
          {{ activeTab === 'content' ? '还没有投稿过模型。' : activeTab === 'favorites' ? '还没有收藏过资源。' : '暂无通知。' }}
        </p>
        <p v-if="activeTab === 'content'" class="mt-2 text-sm text-muted-foreground/70">
          投稿与更新都在 NexusMC 进行，审核通过后会出现在这里。
        </p>
      </div>

      <!-- 分页 -->
      <div
        v-if="!activeState.loading && activeState.items.length && activeState.totalPages > 1"
        class="mt-6 flex items-center justify-center gap-4"
      >
        <Button variant="outline" size="sm" class="btn-press focus-ring" :disabled="activeState.page <= 1" @click="goPage(-1)">
          <ChevronLeft class="mr-1 h-4 w-4" />
          上一页
        </Button>
        <span class="text-sm text-muted-foreground">第 {{ activeState.page }} / {{ activeState.totalPages }} 页</span>
        <Button
          variant="outline"
          size="sm"
          class="btn-press focus-ring"
          :disabled="activeState.page >= activeState.totalPages"
          @click="goPage(1)"
        >
          下一页
          <ChevronRight class="ml-1 h-4 w-4" />
        </Button>
      </div>
    </template>
  </div>
</template>
