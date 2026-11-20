<script setup lang="ts">
import axios from 'axios'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, ArrowUpRight, CalendarDays, CircleAlert, Download, Eye, LogIn } from 'lucide-vue-next'
import NexusmcAvatar from '@/components/NexusmcAvatar.vue'
import ResourceCover from '@/components/ResourceCover.vue'
import {
  isNexusmcConfigError,
  NEXUSMC_SITE,
  nexusmcApi,
  resourceCover,
  resourceDescription,
  resourceDownloads,
  resourcePageUrl,
  resourceTags,
  resourceTitle,
  resourceUpdatedAt,
  resourceViews,
  type NexusmcResource,
} from '@/lib/nexusmc'
import { useNexusmcStore } from '@/stores/nexusmc'
import { tiptapToHtml } from '@/utils/tiptap'

const route = useRoute()
const router = useRouter()
const store = useNexusmcStore()
const loading = ref(true)
const error = ref<'offline' | 'missing' | 'load' | ''>('')
const item = ref<NexusmcResource | null>(null)
const resourceId = computed(() => String(route.params.id ?? ''))
let requestId = 0

const catalogPath = computed(() => {
  const previous = router.options.history.state.back
  return typeof previous === 'string' && /^\/(?:[?#]|$)/.test(previous) ? previous : '/'
})

const model = computed(() => {
  if (!item.value) return null
  const resource = item.value
  const date = new Date(resourceUpdatedAt(resource))
  const format = new Intl.NumberFormat('zh-CN')
  const downloads = resourceDownloads(resource)
  const views = resourceViews(resource)
  return {
    title: resourceTitle(resource),
    description: resourceDescription(resource),
    cover: resourceCover(resource),
    tags: [...new Set(resourceTags(resource))],
    downloads: downloads === undefined ? '' : format.format(downloads),
    views: views === undefined ? '' : format.format(views),
    updatedAt: Number.isNaN(date.getTime()) ? '' : date.toLocaleDateString('zh-CN'),
    pageUrl: resourcePageUrl(resource) || `${NEXUSMC_SITE}/resources/${encodeURIComponent(resourceId.value)}`,
    contentHtml: tiptapToHtml(resource.content ?? resource.body),
  }
})

const errorMessage = computed(() => {
  if (error.value === 'offline') return 'NexusMC 资源服务暂不可用'
  if (error.value === 'missing') return '资源不存在或暂不可见'
  return '模型详情加载失败'
})

async function fetchDetail() {
  const currentRequest = ++requestId
  loading.value = true
  error.value = ''
  item.value = null

  try {
    const { data } = await nexusmcApi.resource(resourceId.value)
    if (currentRequest !== requestId) return
    const resource = data?.resource ?? data?.item ?? data
    if (resource && typeof resource === 'object' && (resource.title || resource.name || resource.page_url)) {
      item.value = resource
    } else {
      error.value = 'missing'
    }
  } catch (cause) {
    if (currentRequest !== requestId) return
    if (isNexusmcConfigError(cause)) {
      error.value = 'offline'
    } else if (axios.isAxiosError(cause) && [403, 404].includes(cause.response?.status ?? 0)) {
      error.value = 'missing'
    } else {
      error.value = 'load'
    }
  } finally {
    if (currentRequest === requestId) loading.value = false
  }
}

watch(resourceId, fetchDetail, { immediate: true })
</script>

<template>
  <div class="container-app py-6 sm:py-9">
    <RouterLink :to="catalogPath" class="inline-flex items-center gap-2 rounded-sm text-sm font-semibold text-muted-foreground hover:text-primary">
      <ArrowLeft class="size-4" aria-hidden="true" /> 返回模型库
    </RouterLink>

    <div v-if="loading" class="mt-6" role="status">
      <span class="sr-only">正在加载模型详情…</span>
      <div class="mb-6 h-10 max-w-80 animate-pulse rounded-md bg-muted" aria-hidden="true" />
      <div class="detail-grid" aria-hidden="true">
        <div class="aspect-[16/10] animate-pulse rounded-2xl bg-muted" />
        <div class="h-56 animate-pulse rounded-2xl bg-muted" />
      </div>
    </div>

    <section v-else-if="error" class="site-panel mt-6 flex flex-col items-start gap-4 p-6 sm:p-8" role="alert">
      <div class="status-icon"><CircleAlert class="size-6" aria-hidden="true" /></div>
      <div>
        <h1 class="text-xl font-bold">{{ errorMessage }}</h1>
        <p class="mt-1 text-sm text-muted-foreground">你也可以前往 NexusMC 查找这个资源。</p>
      </div>
      <div class="flex flex-wrap gap-2">
        <button v-if="error !== 'missing'" type="button" class="site-button site-button-primary" @click="fetchDetail">重试</button>
        <a :href="`${NEXUSMC_SITE}/resources`" target="_blank" rel="noopener noreferrer" class="site-button site-button-outline">
          浏览 NexusMC <ArrowUpRight class="size-4" aria-hidden="true" />
        </a>
      </div>
    </section>

    <article v-else-if="model" class="mt-5 pb-8">
      <header class="mb-6 max-w-3xl">
        <h1 class="font-display text-[clamp(2rem,4.5vw,3.5rem)] font-bold tracking-tight">{{ model.title }}</h1>
        <p v-if="model.description" class="mt-3 whitespace-pre-line text-muted-foreground">{{ model.description }}</p>
        <ul v-if="model.tags.length" class="mt-3 flex flex-wrap gap-1.5" aria-label="资源标签">
          <li v-for="tag in model.tags" :key="tag" class="max-w-full truncate rounded-md bg-accent px-2 py-0.5 text-sm font-medium text-primary">{{ tag }}</li>
        </ul>
      </header>

      <div class="detail-grid">
        <div class="site-panel aspect-[16/10] overflow-hidden bg-muted">
          <ResourceCover :cover="model.cover" :title="model.title" eager icon-class="size-12" />
        </div>

        <aside class="site-panel p-5 sm:p-6" aria-labelledby="download-heading">
          <div class="status-icon"><Download class="size-[1.375rem]" aria-hidden="true" /></div>
          <h2 id="download-heading" class="mt-4 text-xl font-bold">下载模型</h2>
          <p class="mt-1 text-sm text-muted-foreground">在 NexusMC 资源页获取模型文件和最新版本。</p>
          <a :href="model.pageUrl" target="_blank" rel="noopener noreferrer" class="site-button site-button-primary mt-4 w-full">
            前往 NexusMC 下载 <ArrowUpRight class="size-4" aria-hidden="true" />
          </a>

          <dl v-if="model.downloads || model.views || model.updatedAt" class="mt-5 grid gap-2.5 border-t border-border pt-4 text-sm">
            <div v-if="model.downloads" class="flex justify-between gap-3">
              <dt class="inline-flex items-center gap-2 text-muted-foreground"><Download class="size-4" aria-hidden="true" /> 下载次数</dt>
              <dd class="font-semibold tabular-nums">{{ model.downloads }}</dd>
            </div>
            <div v-if="model.views" class="flex justify-between gap-3">
              <dt class="inline-flex items-center gap-2 text-muted-foreground"><Eye class="size-4" aria-hidden="true" /> 浏览次数</dt>
              <dd class="font-semibold tabular-nums">{{ model.views }}</dd>
            </div>
            <div v-if="model.updatedAt" class="flex justify-between gap-3">
              <dt class="inline-flex items-center gap-2 text-muted-foreground"><CalendarDays class="size-4" aria-hidden="true" /> 最近更新</dt>
              <dd class="font-semibold">{{ model.updatedAt }}</dd>
            </div>
          </dl>

          <div class="mt-5 flex flex-wrap items-center gap-2.5 border-t border-border pt-4">
            <template v-if="store.isLoggedIn">
              <NexusmcAvatar size="xs" />
              <RouterLink to="/me" class="min-w-0 truncate rounded-sm text-sm font-semibold hover:text-primary">{{ store.displayName }}</RouterLink>
            </template>
            <template v-else>
              <p class="grow text-sm text-muted-foreground">登录后可在 NexusMC 收藏和评论。</p>
              <button type="button" class="site-button site-button-outline" @click="store.login(route.fullPath)">
                <LogIn class="size-4" aria-hidden="true" /> 登录
              </button>
            </template>
          </div>
        </aside>

        <section v-if="model.contentHtml" class="site-panel min-w-0 p-5 sm:p-7" aria-labelledby="description-heading">
          <h2 id="description-heading" class="mb-3 border-b border-border pb-3 text-xl font-bold">模型介绍</h2>
          <div class="nexusmc-article" v-html="model.contentHtml" />
        </section>
      </div>
    </article>
  </div>
</template>

<style scoped>
.detail-grid {
  display: grid;
  align-items: start;
  gap: 1.25rem;
}

.nexusmc-article {
  overflow-wrap: anywhere;
}

.nexusmc-article :deep(p),
.nexusmc-article :deep(ul),
.nexusmc-article :deep(ol),
.nexusmc-article :deep(blockquote),
.nexusmc-article :deep(pre),
.nexusmc-article :deep(img),
.nexusmc-article :deep(hr) {
  margin-block: 0.75rem;
}

.nexusmc-article :deep(h2),
.nexusmc-article :deep(h3),
.nexusmc-article :deep(h4) {
  margin-top: 1.25rem;
  font-weight: 700;
}

.nexusmc-article :deep(ul),
.nexusmc-article :deep(ol) {
  padding-inline-start: 1.5rem;
  list-style: revert;
}

.nexusmc-article :deep(blockquote) {
  border-inline-start: 3px solid var(--color-border);
  padding-inline-start: 1rem;
  color: var(--color-muted-foreground);
}

.nexusmc-article :deep(pre) {
  overflow-x: auto;
  border-radius: 0.5rem;
  background: var(--color-muted);
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
}

.nexusmc-article :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
  text-underline-offset: 2px;
}

@media (min-width: 850px) {
  .detail-grid {
    grid-template-columns: minmax(0, 2fr) minmax(18rem, 1fr);
    gap: 1.5rem;
  }
}
</style>
