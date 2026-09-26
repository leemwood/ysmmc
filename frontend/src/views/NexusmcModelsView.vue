<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
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
  Settings2,
  ChevronLeft,
  ChevronRight,
} from 'lucide-vue-next'
import {
  NEXUSMC_LOGIN_URL,
  clearNexusmcUser,
  getNexusmcUser,
  isNexusmcConfigError,
  nexusmcApi,
  resourceCover,
  resourceDescription,
  resourceDownloads,
  resourcePageUrl,
  resourceTags,
  resourceTitle,
  resourceViews,
  type NexusmcListResponse,
  type NexusmcResource,
  type NexusmcUser,
} from '@/lib/nexusmc'

interface CatalogOption {
  value: string
  label: string
}

const PAGE_SIZE = 20

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
const user = ref<NexusmcUser | null>(null)

const totalPagesSafe = computed(() => Math.max(totalPages.value, 1))

function extractCategories(data: unknown): CatalogOption[] {
  // 上游目录结构未在文档中给出完整字段，做多形状兼容。
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

function handleLogout() {
  clearNexusmcUser()
  user.value = null
}

onMounted(async () => {
  user.value = getNexusmcUser()
  fetchCatalog()
  fetchResources()
})
</script>

<template>
  <div class="container-app py-6 sm:py-8">
    <!-- Header -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight">NexusMC 模型资源</h1>
        <p class="mt-1 text-sm text-muted-foreground">
          来自 NexusMC 社区的公开模型资源，点击下载将跳转到 NexusMC 资源页面。
        </p>
      </div>

      <!-- NexusMC login -->
      <div class="flex items-center gap-3">
        <template v-if="user">
          <Avatar class="h-9 w-9 border">
            <AvatarImage v-if="user.avatar" :src="user.avatar" :alt="user.username" />
            <span v-else class="flex h-full w-full items-center justify-center rounded-full bg-muted text-sm font-medium">
              {{ user.username?.slice(0, 1).toUpperCase() || 'N' }}
            </span>
          </Avatar>
          <div class="text-sm">
            <div class="font-medium leading-tight">{{ user.username }}</div>
            <div class="text-xs text-muted-foreground">NexusMC 已登录</div>
          </div>
          <Button variant="ghost" size="sm" class="btn-press focus-ring" @click="handleLogout">
            <LogOut class="h-4 w-4" />
          </Button>
        </template>
        <template v-else>
          <a :href="NEXUSMC_LOGIN_URL">
            <Button size="sm" class="btn-press focus-ring" :disabled="notConfigured">
              通过 NexusMC 登录
            </Button>
          </a>
        </template>
      </div>
    </div>

    <!-- Not configured banner -->
    <Alert v-if="notConfigured" class="mt-6">
      <Settings2 class="h-4 w-4" />
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
    <div class="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center">
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
        <div class="aspect-[4/3] w-full bg-muted overflow-hidden">
          <img
            v-if="resourceCover(item)"
            :src="resourceCover(item)"
            :alt="resourceTitle(item)"
            loading="lazy"
            decoding="async"
            class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
          <div v-else class="flex h-full flex-col items-center justify-center gap-2 text-muted-foreground">
            <ImageOff class="h-10 w-10 opacity-50" />
            <span class="text-sm">无预览图</span>
          </div>
        </div>
        <CardContent class="p-4">
          <h3 class="font-semibold line-clamp-1 text-balance">{{ resourceTitle(item) }}</h3>
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
            <a
              v-if="resourcePageUrl(item)"
              :href="resourcePageUrl(item)"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Button size="sm" class="btn-press focus-ring">
                <ExternalLink class="mr-1.5 h-3.5 w-3.5" />
                前往下载
              </Button>
            </a>
          </div>
        </CardContent>
      </Card>
    </div>

    <div v-else-if="!notConfigured && !error" class="mt-16 text-center text-muted-foreground">
      没有找到匹配的资源。
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
  </div>
</template>
