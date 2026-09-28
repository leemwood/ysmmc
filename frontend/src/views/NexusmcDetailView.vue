<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Skeleton } from '@/components/ui/skeleton'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import {
  ArrowLeft,
  Download,
  Eye,
  ExternalLink,
  ImageOff,
  CalendarDays,
} from 'lucide-vue-next'
import {
  isNexusmcConfigError,
  nexusmcApi,
  resourceCover,
  resourceDescription,
  resourceDownloads,
  resourcePageUrl,
  resourceTags,
  resourceTitle,
  resourceViews,
  type NexusmcResource,
} from '@/lib/nexusmc'
import { NEXUSMC_SITE, useNexusmcStore } from '@/stores/nexusmc'
import { tiptapToHtml } from '@/utils/tiptap'

const route = useRoute()
const store = useNexusmcStore()

const loading = ref(true)
const error = ref('')
const notFound = ref(false)
const item = ref<NexusmcResource | null>(null)

const resourceId = computed(() => String(route.params.id || ''))
const pageUrl = computed(() => (item.value ? resourcePageUrl(item.value) : ''))
const updatedAt = computed(() => {
  if (!item.value) return ''
  const raw = item.value.updated_at || item.value.updatedAt || item.value.updated_at_time
  if (typeof raw !== 'string' || !raw) return ''
  const date = new Date(raw)
  return Number.isNaN(date.getTime()) ? '' : date.toLocaleDateString('zh-CN')
})
const displayName = computed(() => store.user?.username || '')
// TipTap 正文 → 受限 HTML；无正文时回退为简介
const contentHtml = computed(() => {
  if (!item.value) return ''
  const raw = (item.value as Record<string, unknown>).content ?? (item.value as Record<string, unknown>).body
  return tiptapToHtml(raw)
})

async function fetchDetail() {
  loading.value = true
  error.value = ''
  notFound.value = false
  try {
    const { data } = await nexusmcApi.resource(resourceId.value)
    const node = (data as Record<string, unknown>).resource ?? (data as Record<string, unknown>).item ?? data
    if (node && typeof node === 'object' && (resourceTitle(node as NexusmcResource) || (node as NexusmcResource).page_url)) {
      item.value = node as NexusmcResource
    } else {
      notFound.value = true
    }
  } catch (e) {
    if (isNexusmcConfigError(e)) {
      error.value = 'not_configured'
    } else {
      notFound.value = true
    }
  } finally {
    loading.value = false
  }
}

onMounted(fetchDetail)
</script>

<template>
  <div class="container-app py-6 sm:py-8">
    <Button variant="ghost" size="sm" class="btn-press focus-ring -ml-2" @click="$router.back()">
      <ArrowLeft class="mr-1.5 h-4 w-4" />
      返回资源列表
    </Button>

    <!-- Loading -->
    <div v-if="loading" class="mt-6 grid gap-6 lg:grid-cols-[minmax(0,2fr)_minmax(0,1fr)]">
      <Skeleton class="aspect-[16/9] rounded-xl" />
      <div class="space-y-4">
        <Skeleton class="h-8 w-3/4" />
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-5/6" />
        <Skeleton class="h-10 w-40" />
      </div>
    </div>

    <!-- Error / fallback：详情缺失时仍给出 NexusMC 出口，保证引流不中断 -->
    <template v-else-if="notFound || error">
      <Alert class="mt-6" :variant="error === 'not_configured' ? 'default' : 'destructive'">
        <AlertDescription>
          <template v-if="error === 'not_configured'">
            <span class="font-medium">NexusMC 凭据尚未配置。</span>
            详情展示暂不可用，可直接前往 NexusMC 浏览该资源。
          </template>
          <template v-else>
            <span class="font-medium">未能获取资源详情。</span>
            该资源可能已下架或暂不可见，你可以前往 NexusMC 搜索查看。
          </template>
        </AlertDescription>
      </Alert>
      <div class="mt-6 flex flex-wrap gap-3">
        <a :href="`${NEXUSMC_SITE}/resources`" target="_blank" rel="noopener noreferrer">
          <Button class="btn-press focus-ring">
            <ExternalLink class="mr-1.5 h-4 w-4" />
            前往 NexusMC 资源库
          </Button>
        </a>
        <Button variant="outline" class="btn-press focus-ring" @click="$router.push('/')">
          返回资源列表
        </Button>
      </div>
    </template>

    <!-- Detail -->
    <div v-else-if="item" class="mt-4">
      <div class="flex flex-col gap-6 lg:flex-row">
        <!-- 封面 -->
        <div class="lg:w-3/5">
          <div class="aspect-[16/9] w-full overflow-hidden rounded-xl border bg-muted">
            <img
              v-if="resourceCover(item)"
              :src="resourceCover(item)"
              :alt="resourceTitle(item)"
              class="h-full w-full object-cover"
            />
            <div v-else class="flex h-full flex-col items-center justify-center gap-2 text-muted-foreground">
              <ImageOff class="h-12 w-12 opacity-50" />
              <span class="text-sm">无预览图</span>
            </div>
          </div>
        </div>

        <!-- 信息与 CTA -->
        <div class="flex flex-col lg:w-2/5">
          <h1 class="text-xl font-bold tracking-tight sm:text-2xl">{{ resourceTitle(item) }}</h1>
          <p v-if="resourceDescription(item)" class="mt-3 whitespace-pre-line text-sm leading-relaxed text-muted-foreground">
            {{ resourceDescription(item) }}
          </p>

          <div v-if="resourceTags(item).length" class="mt-4 flex flex-wrap gap-1.5">
            <Badge v-for="tag in resourceTags(item)" :key="tag" variant="secondary">{{ tag }}</Badge>
          </div>

          <div class="mt-4 flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            <span v-if="resourceDownloads(item) !== undefined" class="inline-flex items-center gap-1.5">
              <Download class="h-4 w-4" />
              {{ resourceDownloads(item) }} 次下载
            </span>
            <span v-if="resourceViews(item) !== undefined" class="inline-flex items-center gap-1.5">
              <Eye class="h-4 w-4" />
              {{ resourceViews(item) }} 次浏览
            </span>
            <span v-if="updatedAt" class="inline-flex items-center gap-1.5">
              <CalendarDays class="h-4 w-4" />
              更新于 {{ updatedAt }}
            </span>
          </div>

          <!-- 核心引流 CTA -->
          <div class="mt-6 rounded-xl border bg-primary/5 p-4">
            <p class="text-sm font-medium">下载与更新都在 NexusMC 进行</p>
            <p class="mt-1 text-xs text-muted-foreground">
              点击前往 NexusMC 资源页获取最新版本；登录、收藏与评论同样在 NexusMC 完成。
            </p>
            <a v-if="pageUrl" :href="pageUrl" target="_blank" rel="noopener noreferrer" class="mt-3 block">
              <Button size="lg" class="btn-press focus-ring w-full sm:w-auto">
                <ExternalLink class="mr-2 h-4 w-4" />
                前往 NexusMC 下载
              </Button>
            </a>
          </div>

          <!-- 登录状态 -->
          <div class="mt-4 flex items-center gap-3 rounded-xl border p-3">
            <template v-if="store.user">
              <Avatar class="h-8 w-8 border">
                <AvatarImage v-if="store.user.avatar" :src="store.user.avatar" :alt="store.user.username" />
                <span v-else class="flex h-full w-full items-center justify-center rounded-full bg-muted text-xs font-medium">
                  {{ displayName.slice(0, 1).toUpperCase() }}
                </span>
              </Avatar>
              <div class="min-w-0 text-sm">
                已以 <span class="font-medium">{{ displayName }}</span> 身份浏览，
                <a :href="NEXUSMC_SITE" target="_blank" rel="noopener noreferrer" class="text-primary hover:underline">去 NexusMC 互动 →</a>
              </div>
            </template>
            <template v-else>
              <div class="min-w-0 flex-1 text-sm text-muted-foreground">
                登录 NexusMC 后可收藏资源、参与评论。
              </div>
              <Button size="sm" variant="outline" class="btn-press focus-ring shrink-0" @click="store.login()">
                登录
              </Button>
            </template>
          </div>
        </div>
      </div>

      <!-- 详情正文（TipTap JSON 渲染为受限 HTML；纯文本原样展示） -->
      <Card v-if="contentHtml" class="mt-8">
        <CardContent class="nexusmc-article p-6" v-html="contentHtml" />
      </Card>
    </div>
  </div>
</template>

<style scoped>
.nexusmc-article :deep(p) {
  margin: 0.75rem 0;
  line-height: 1.8;
}
.nexusmc-article :deep(h2),
.nexusmc-article :deep(h3),
.nexusmc-article :deep(h4) {
  margin: 1.25rem 0 0.5rem;
  font-weight: 600;
}
.nexusmc-article :deep(ul),
.nexusmc-article :deep(ol) {
  margin: 0.75rem 0;
  padding-left: 1.5rem;
  list-style: revert;
}
.nexusmc-article :deep(li) {
  margin: 0.25rem 0;
  line-height: 1.7;
}
.nexusmc-article :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 0.75rem 0;
}
.nexusmc-article :deep(blockquote) {
  border-left: 3px solid var(--border);
  padding-left: 1rem;
  margin: 0.75rem 0;
  color: var(--muted-foreground);
}
.nexusmc-article :deep(pre) {
  background: var(--muted);
  border-radius: 0.5rem;
  padding: 0.75rem 1rem;
  overflow-x: auto;
  font-size: 0.875rem;
}
.nexusmc-article :deep(code) {
  font-size: 0.875em;
}
.nexusmc-article :deep(a) {
  color: var(--primary);
  text-decoration: underline;
  text-underline-offset: 2px;
  overflow-wrap: anywhere;
}
.nexusmc-article :deep(a:hover) {
  opacity: 0.85;
}
.nexusmc-article :deep(hr) {
  border-color: var(--border);
  margin: 1rem 0;
}
</style>
