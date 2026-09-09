<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { nexusmcApi } from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { ExternalLink, ChevronLeft, ChevronRight, Package } from 'lucide-vue-next'

// NexusMC 资源列表结构字段以实际 API 返回为准（文档：站点 API：资源），
// 这里对常见字段做防御性取值。
interface NexusResource {
  id?: string | number
  title?: string
  name?: string
  slug?: string
  description?: string
  summary?: string
  icon?: string
  avatar?: string
  downloads?: number
  likes?: number
  views?: number
  [key: string]: unknown
}

interface NexusListResponse {
  data?: NexusResource[] | { items?: NexusResource[]; list?: NexusResource[] }
  items?: NexusResource[]
  pagination?: { total?: number; page?: number; pageSize?: number }
  [key: string]: unknown
}

const resources = ref<NexusResource[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const platform = ref('')
const loading = ref(false)
const error = ref('')

function pickList(payload: NexusListResponse): NexusResource[] {
  const d = payload.data
  if (Array.isArray(d)) return d
  if (d && typeof d === 'object') {
    const obj = d as { items?: NexusResource[]; list?: NexusResource[] }
    if (Array.isArray(obj.items)) return obj.items
    if (Array.isArray(obj.list)) return obj.list
  }
  if (Array.isArray(payload.items)) return payload.items
  return []
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const params: Record<string, string | number> = { page: page.value, pageSize }
    if (platform.value) params.platform = platform.value
    const res = await nexusmcApi.resources(params)
    const payload = res.data as unknown as NexusListResponse
    resources.value = pickList(payload)
    total.value = payload.pagination?.total ?? resources.value.length
  } catch (err: any) {
    error.value = err.response?.data?.message || err.message || '加载 NexusMC 资源失败'
    resources.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(page, load)
watch(platform, () => {
  page.value = 1
  load()
})

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

function displayName(r: NexusResource) {
  return r.title || r.name || r.slug || String(r.id ?? '未命名资源')
}
function displaySummary(r: NexusResource) {
  const s = r.summary || r.description || ''
  return s.length > 120 ? s.slice(0, 120) + '…' : s
}
function iconUrl(r: NexusResource) {
  const p = r.icon || r.avatar
  if (!p) return ''
  return p.startsWith('http') ? p : `https://www.nexusmc.cn${p}`
}
</script>

<template>
  <div class="mx-auto max-w-5xl px-4 py-6 sm:py-8">
    <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-2xl font-bold">NexusMC 资源市场</h1>
        <p class="text-sm text-muted-foreground">浏览 NexusMC 公开资源（数据来自 NexusMC 站点 API）</p>
      </div>
      <div class="flex items-center gap-2">
        <select
          v-model="platform"
          class="h-9 rounded-md border border-input bg-background px-3 text-sm"
        >
          <option value="">全部平台</option>
          <option value="java">Java 版</option>
          <option value="bedrock">基岩版</option>
        </select>
      </div>
    </div>

    <Alert v-if="error" variant="destructive" class="mb-4">
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Skeleton v-for="i in 6" :key="i" class="h-36 rounded-xl" />
    </div>

    <template v-else-if="resources.length">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <Card v-for="r in resources" :key="String(r.id ?? r.slug)" class="transition-shadow hover:shadow-md">
          <CardContent class="flex h-full flex-col gap-3 p-4">
            <div class="flex items-center gap-3">
              <img
                v-if="iconUrl(r)"
                :src="iconUrl(r)"
                :alt="displayName(r)"
                class="h-10 w-10 rounded-md object-cover"
                loading="lazy"
              />
              <div v-else class="flex h-10 w-10 items-center justify-center rounded-md bg-muted">
                <Package class="h-5 w-5 text-muted-foreground" />
              </div>
              <h2 class="line-clamp-1 flex-1 font-semibold">{{ displayName(r) }}</h2>
            </div>
            <p class="line-clamp-3 flex-1 text-sm text-muted-foreground">
              {{ displaySummary(r) || '暂无简介' }}
            </p>
            <div class="flex items-center justify-between text-xs text-muted-foreground">
              <span v-if="r.downloads != null">下载 {{ r.downloads }}</span>
              <a
                v-if="r.slug"
                :href="`https://www.nexusmc.cn/resources/${r.slug}`"
                target="_blank"
                rel="noopener noreferrer"
                class="ml-auto inline-flex items-center gap-1 hover:text-foreground"
              >
                前往 NexusMC
                <ExternalLink class="h-3 w-3" />
              </a>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-if="totalPages > 1" class="mt-6 flex items-center justify-center gap-3">
        <Button variant="outline" size="sm" :disabled="page <= 1 || loading" @click="page--">
          <ChevronLeft class="h-4 w-4" />
          上一页
        </Button>
        <span class="text-sm text-muted-foreground">{{ page }} / {{ totalPages }}</span>
        <Button variant="outline" size="sm" :disabled="page >= totalPages || loading" @click="page++">
          下一页
          <ChevronRight class="h-4 w-4" />
        </Button>
      </div>
    </template>

    <Card v-else-if="!error">
      <CardContent class="py-12 text-center text-muted-foreground">
        暂无资源
      </CardContent>
    </Card>
  </div>
</template>
