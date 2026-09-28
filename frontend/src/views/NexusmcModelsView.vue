<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Skeleton } from '@/components/ui/skeleton'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import {
  Download,
  Eye,
  ExternalLink,
  ImageOff,
  LogOut,
  RefreshCw,
  Search,
  ChevronLeft,
  ChevronRight,
  UserPlus,
  Sparkles,
} from 'lucide-vue-next'
import {
  isNexusmcConfigError,
  nexusmcApi,
  resourceCover,
  resourceDescription,
  resourceDownloads,
  resourceTags,
  resourceTitle,
  resourceViews,
  type NexusmcListResponse,
  type NexusmcResource,
} from '@/lib/nexusmc'
import { NEXUSMC_REGISTER_URL, NEXUSMC_SITE, useNexusmcStore } from '@/stores/nexusmc'

interface CatalogOption {
  value: string
  label: string
}

const PAGE_SIZE = 20

const store = useNexusmcStore()
const loading = ref(false)
const error = ref('')
const notConfigured = ref(false)
const items = ref<NexusmcResource[]>([])
const total = ref(0)
const totalPages = ref(1)
const page = ref(1)
const keyword = ref('')
const searchInput = ref('')
const categories = ref<CatalogOption[]>([])
const activeCategory = ref('')

const totalPagesSafe = computed(() => Math.max(totalPages.value, 1))
const displayName = computed(() => store.user?.username || 'NexusMC 用户')

function extractCategories(data: unknown): CatalogOption[] {
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

async function fetchCatalog() {
  try {
    const { data } = await nexusmcApi.catalog()
    categories.value = extractCategories(data)
  } catch {
    categories.value = []
  }
}

async function fetchResources() {
  loading.value = true
  error.value = ''
  try {
    let data: NexusmcListResponse
    if (keyword.value) {
      data = (
        await nexusmcApi.search(keyword.value, { page: page.value, pageSize: PAGE_SIZE, category: activeCategory.value || undefined })
      ).data
    } else {
      data = (
        await nexusmcApi.resources({
          page: page.value,
          pageSize: PAGE_SIZE,
          category: activeCategory.value || undefined,
        })
      ).data
    }
    items.value = data.items || []
    total.value = data.pagination?.total ?? items.value.length
    totalPages.value = data.pagination?.totalPages ?? 1
  } catch (e) {
    if (isNexusmcConfigError(e)) {
      notConfigured.value = true
    } else {
      error.value = '获取 NexusMC 资源失败，请稍后重试。'
    }
    items.value = []
  } finally {
    loading.value = false
  }
}

function submitSearch() {
  keyword.value = searchInput.value.trim()
  page.value = 1
  fetchResources()
}

function clearSearch() {
  searchInput.value = ''
  keyword.value = ''
  page.value = 1
  fetchResources()
}

function changeCategory() {
  page.value = 1
  fetchResources()
}

function goPage(delta: number) {
  const next = page.value + delta
  if (next < 1 || next > totalPagesSafe.value) return
  page.value = next
  fetchResources()
  if (typeof window !== 'undefined') window.scrollTo({ top: 0, behavior: 'smooth' })
}

function detailTo(item: NexusmcResource) {
  const id = item.slug || item.id
  return id ? `/nexusmc/resource/${encodeURIComponent(String(id))}` : ''
}

onMounted(() => {
  fetchCatalog()
  fetchResources()
})
</script>

<template>
  <div class="container-app py-6 sm:py-8">
    <!-- Hero：站点定位与引流 -->
    <section class="rounded-2xl border bg-gradient-to-br from-primary/10 via-background to-background p-6 sm:p-10">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
        <div class="max-w-2xl">
          <p class="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary">
            <Sparkles class="h-3.5 w-3.5" />
            NexusMC 资源精选站
          </p>
          <h1 class="mt-3 text-2xl font-bold tracking-tight sm:text-3xl">
            {{ store.isLoggedIn ? `欢迎回来，${displayName}` : 'YSM 模型资源精选' }}
          </h1>
          <p class="mt-2 text-sm text-muted-foreground sm:text-base">
            这里汇集了 NexusMC 社区的公开模型资源。浏览、挑选，一键前往 NexusMC 完成下载——
            登录、收藏与互动都在 NexusMC 进行。
          </p>
          <div class="mt-5 flex flex-wrap items-center gap-3">
            <a :href="NEXUSMC_SITE" target="_blank" rel="noopener noreferrer">
              <Button class="btn-press focus-ring">
                <ExternalLink class="mr-1.5 h-4 w-4" />
                前往 NexusMC
              </Button>
            </a>
            <template v-if="!store.isLoggedIn">
              <Button variant="outline" class="btn-press focus-ring" @click="store.login()">
                通过 NexusMC 登录
              </Button>
              <a :href="NEXUSMC_REGISTER_URL" target="_blank" rel="noopener noreferrer">
                <Button variant="ghost" class="btn-press focus-ring">
                  <UserPlus class="mr-1.5 h-4 w-4" />
                  注册 NexusMC
                </Button>
              </a>
            </template>
          </div>
        </div>

        <!-- 用户信息卡片 -->
        <div v-if="store.user" class="flex items-center gap-4 rounded-xl border bg-background/80 p-4 lg:w-80">
          <Avatar class="h-12 w-12 border">
            <AvatarImage v-if="store.user.avatar" :src="store.user.avatar" :alt="store.user.username" />
            <span v-else class="flex h-full w-full items-center justify-center rounded-full bg-muted text-lg font-medium">
              {{ displayName.slice(0, 1).toUpperCase() }}
            </span>
          </Avatar>
          <div class="min-w-0 flex-1">
            <div class="truncate font-semibold">{{ displayName }}</div>
            <div class="truncate text-xs text-muted-foreground">
              {{ store.user.email || `NexusMC UID ${store.user.uid ?? '—'}` }}
            </div>
          </div>
          <Button variant="ghost" size="icon" class="btn-press focus-ring shrink-0" aria-label="退出 NexusMC 登录" @click="store.logout()">
            <LogOut class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </section>

    <!-- Not configured banner -->
    <Alert v-if="notConfigured" class="mt-6">
      <AlertDescription>
        <span class="font-medium">NexusMC 凭据尚未配置。</span>
        部署环境缺少 NEXUSMC_SITE_API_KEY 等 EdgeOne 环境变量，资源展示与登录暂不可用。
        配置步骤见仓库 docs/edgeone-deploy.md。
      </AlertDescription>
    </Alert>

    <Alert v-else-if="error" class="mt-6" variant="destructive">
      <AlertDescription>
        <span class="font-medium">加载失败。</span>
        {{ error }}
      </AlertDescription>
    </Alert>

    <!-- Toolbar -->
    <div class="mt-8 flex flex-col gap-3 sm:flex-row sm:items-center">
      <form class="flex flex-1 gap-2" @submit.prevent="submitSearch">
        <div class="relative flex-1 max-w-md">
          <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            v-model="searchInput"
            placeholder="搜索 NexusMC 资源…"
            class="pl-9"
            maxlength="100"
          />
        </div>
        <Button type="submit" variant="secondary" size="sm" class="btn-press focus-ring">搜索</Button>
        <Button
          v-if="keyword"
          type="button"
          variant="ghost"
          size="sm"
          class="btn-press focus-ring"
          @click="clearSearch"
        >
          清除
        </Button>
      </form>

      <div class="flex items-center gap-2">
        <select
          v-if="categories.length"
          v-model="activeCategory"
          class="h-9 rounded-md border border-input bg-background px-3 text-sm focus-ring"
          @change="changeCategory"
        >
          <option value="">全部分类</option>
          <option v-for="opt in categories" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <Button
          variant="ghost"
          size="icon"
          class="btn-press focus-ring"
          aria-label="刷新"
          :disabled="loading"
          @click="fetchResources"
        >
          <RefreshCw class="h-4 w-4" :class="{ 'animate-spin': loading }" />
        </Button>
      </div>
    </div>

    <!-- Grid -->
    <div
      v-if="loading"
      class="mt-6 grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
    >
      <Skeleton v-for="i in 8" :key="i" class="h-64 rounded-xl" />
    </div>

    <div
      v-else-if="items.length"
      class="mt-6 grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
    >
      <Card
        v-for="item in items"
        :key="item.id || item.slug || resourceTitle(item)"
        class="overflow-hidden card-hover group"
      >
        <RouterLink :to="detailTo(item)" class="block focus-ring">
          <div class="aspect-[4/3] w-full bg-muted overflow-hidden">
            <img
              v-if="resourceCover(item)"
              :src="resourceCover(item)"
              :alt="resourceTitle(item)"
              loading="lazy"
              decoding="async"
              style="height: 100%; width: 100%; object-fit: contain"
              class="transition-transform duration-300 group-hover:scale-105"
            />
            <div v-else class="flex h-full flex-col items-center justify-center gap-2 text-muted-foreground">
              <ImageOff class="h-10 w-10 opacity-50" />
              <span class="text-sm">无预览图</span>
            </div>
          </div>
          <CardContent class="p-4">
            <h3 class="font-semibold line-clamp-1 text-balance group-hover:text-primary transition-colors">
              {{ resourceTitle(item) }}
            </h3>
            <p class="mt-1 text-sm text-muted-foreground line-clamp-2 min-h-[2.5rem]">
              {{ resourceDescription(item) || '暂无描述' }}
            </p>
            <div v-if="resourceTags(item).length" class="mt-3 flex flex-wrap gap-1.5">
              <Badge v-for="tag in resourceTags(item)" :key="tag" variant="secondary" class="text-xs">
                {{ tag }}
              </Badge>
            </div>
            <div class="mt-3 flex items-center justify-between">
              <div class="flex items-center gap-3 text-xs text-muted-foreground">
                <span v-if="resourceDownloads(item) !== undefined" class="inline-flex items-center gap-1">
                  <Download class="h-3.5 w-3.5" />
                  {{ resourceDownloads(item) }}
                </span>
                <span v-if="resourceViews(item) !== undefined" class="inline-flex items-center gap-1">
                  <Eye class="h-3.5 w-3.5" />
                  {{ resourceViews(item) }}
                </span>
              </div>
              <span class="text-xs font-medium text-primary">查看详情 →</span>
            </div>
          </CardContent>
        </RouterLink>
      </Card>
    </div>

    <div v-else-if="!notConfigured && !error" class="mt-16 text-center">
      <p class="text-muted-foreground">暂无 YSM 模型资源。</p>
      <p class="mt-2 text-sm text-muted-foreground/70">
        新模型在 NexusMC 审核通过后会自动出现在这里，
        <a :href="NEXUSMC_SITE + '/resources/new'" target="_blank" rel="noopener noreferrer" class="text-primary hover:underline">欢迎前往 NexusMC 投稿 →</a>
      </p>
    </div>

    <!-- Pagination -->
    <div
      v-if="!loading && items.length && totalPagesSafe > 1"
      class="mt-8 flex items-center justify-center gap-4"
    >
      <Button
        variant="outline"
        size="sm"
        class="btn-press focus-ring"
        :disabled="page <= 1"
        @click="goPage(-1)"
      >
        <ChevronLeft class="mr-1 h-4 w-4" />
        上一页
      </Button>
      <span class="text-sm text-muted-foreground">第 {{ page }} / {{ totalPagesSafe }} 页 · 共 {{ total }} 个资源</span>
      <Button
        variant="outline"
        size="sm"
        class="btn-press focus-ring"
        :disabled="page >= totalPagesSafe"
        @click="goPage(1)"
      >
        下一页
        <ChevronRight class="ml-1 h-4 w-4" />
      </Button>
    </div>

    <!-- 底部引流 -->
    <section class="mt-12 rounded-2xl border bg-muted/30 p-6 text-center sm:p-8">
      <h2 class="text-lg font-semibold">在 NexusMC 获得完整体验</h2>
      <p class="mx-auto mt-2 max-w-xl text-sm text-muted-foreground">
        登录收藏、社区互动、投稿分享——注册 NexusMC 账号，发现更多资源与玩法。
      </p>
      <div class="mt-4 flex flex-wrap items-center justify-center gap-3">
        <a :href="NEXUSMC_SITE" target="_blank" rel="noopener noreferrer">
          <Button class="btn-press focus-ring">打开 NexusMC</Button>
        </a>
        <a :href="NEXUSMC_REGISTER_URL" target="_blank" rel="noopener noreferrer">
          <Button variant="outline" class="btn-press focus-ring">免费注册</Button>
        </a>
      </div>
    </section>
  </div>
</template>
