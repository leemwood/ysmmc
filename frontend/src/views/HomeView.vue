<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { modelApi, announcementApi } from '@/lib/api'
import type { Model, Announcement, PaginatedResponse } from '@/types'
import ModelCard from '@/components/ModelCard.vue'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Search, ChevronLeft, ChevronRight, Megaphone, X, ImageOff, SlidersHorizontal } from 'lucide-vue-next'

const route = useRoute()

const models = ref<Model[]>([])
const announcements = ref<Announcement[]>([])
const loading = ref(true)
const loadingMore = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = 12
const total = ref(0)
const totalPages = ref(0)
const dismissedAnnouncements = ref<Set<string>>(new Set())

async function fetchModels() {
  loading.value = true
  try {
    const response = await modelApi.list(page.value, pageSize, search.value)
    const data = response.data.data as PaginatedResponse<Model>
    models.value = data.items
    total.value = data.total
    totalPages.value = data.total_pages
  } catch (error) {
    console.error('Failed to fetch models:', error)
  } finally {
    loading.value = false
  }
}

async function fetchAnnouncements() {
  try {
    const response = await announcementApi.list()
    announcements.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch announcements:', error)
  }
}

function handleSearch() {
  page.value = 1
  fetchModels()
}

function prevPage() {
  if (page.value > 1) {
    page.value--
    fetchModels()
  }
}

function nextPage() {
  if (page.value < totalPages.value) {
    page.value++
    fetchModels()
  }
}

async function loadMore() {
  if (page.value >= totalPages.value || loadingMore.value) return
  loadingMore.value = true
  try {
    const nextPage = page.value + 1
    const response = await modelApi.list(nextPage, pageSize, search.value)
    const data = response.data.data as PaginatedResponse<Model>
    models.value.push(...data.items)
    page.value = nextPage
  } catch (error) {
    console.error('Failed to load more models:', error)
  } finally {
    loadingMore.value = false
  }
}

function clearFilters() {
  search.value = ''
  page.value = 1
  fetchModels()
}

function dismissAnnouncement(id: string) {
  dismissedAnnouncements.value.add(id)
}

function isAnnouncementVisible(id: string) {
  return !dismissedAnnouncements.value.has(id)
}

watch(() => route.query.search, (newSearch) => {
  if (newSearch) {
    search.value = newSearch as string
    fetchModels()
  }
})

onMounted(() => {
  fetchModels()
  fetchAnnouncements()
})
</script>

<template>
  <div class="container-app py-6 sm:py-8">
    <!-- 公告横幅 -->
    <div v-if="announcements.some(a => isAnnouncementVisible(a.id))" class="mb-6">
      <div
        class="flex gap-3 overflow-x-auto pb-2 snap-x snap-mandatory [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
      >
        <Card
          v-for="announcement in announcements"
          :key="announcement.id"
          v-show="isAnnouncementVisible(announcement.id)"
          class="bg-primary/5 border-primary/20 relative overflow-hidden snap-start shrink-0 w-[85vw] sm:w-[60vw] lg:w-[45vw] max-w-md"
        >
          <CardContent class="flex items-start gap-3 p-4">
            <Megaphone class="h-5 w-5 text-primary flex-shrink-0 mt-0.5" />
            <div class="flex-1 min-w-0">
              <h4 class="font-semibold text-sm sm:text-base line-clamp-1">
                {{ announcement.title }}
              </h4>
              <p class="text-xs sm:text-sm text-muted-foreground mt-1 line-clamp-2">
                {{ announcement.content }}
              </p>
            </div>
            <Button
              variant="ghost"
              size="icon"
              class="h-6 w-6 flex-shrink-0"
              @click.stop="dismissAnnouncement(announcement.id)"
            >
              <X class="h-4 w-4" />
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="mb-6 sm:mb-8">
      <form
        class="flex items-center gap-2 p-1.5 pr-2 rounded-full border bg-card shadow-sm sm:bg-transparent sm:shadow-none sm:border-0 sm:rounded-lg sm:p-0"
        @submit.prevent="handleSearch"
      >
        <Input
          v-model="search"
          placeholder="搜索模型..."
          class="h-9 w-full border-0 bg-transparent shadow-none focus-visible:ring-0 sm:border sm:bg-background sm:focus-visible:ring-1"
        >
          <template #prefix>
            <Search class="h-4 w-4 text-muted-foreground" />
          </template>
        </Input>
        <Button
          type="submit"
          class="btn-press h-8 w-8 rounded-full sm:h-9 sm:w-auto sm:rounded-md shrink-0"
        >
          <Search class="h-4 w-4 sm:hidden" />
          <span class="hidden sm:inline">搜索</span>
        </Button>
      </form>
    </div>

    <!-- 加载骨架屏 -->
    <div
      v-if="loading"
      class="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
    >
      <Card v-for="i in 8" :key="i" class="overflow-hidden">
        <Skeleton variant="shimmer" class="aspect-[4/3] w-full rounded-none" />
        <CardContent class="p-4 space-y-3">
          <Skeleton variant="shimmer" class="h-5 w-3/4" />
          <Skeleton variant="shimmer" class="h-4 w-full" />
          <Skeleton variant="shimmer" class="h-4 w-2/3" />
          <div class="flex gap-1.5 pt-1">
            <Skeleton variant="shimmer" class="h-5 w-12 rounded-full" />
            <Skeleton variant="shimmer" class="h-5 w-12 rounded-full" />
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- 空状态 -->
    <div
      v-else-if="models.length === 0"
      class="flex flex-col items-center justify-center py-12 sm:py-16 text-center animate-fade-in"
    >
      <div class="surface mb-4 flex h-28 w-28 items-center justify-center rounded-2xl">
        <ImageOff class="h-14 w-14 text-muted-foreground opacity-50" />
      </div>
      <h3 class="text-lg font-semibold">暂无模型</h3>
      <p class="mt-1 text-sm text-muted-foreground max-w-xs">
        没有找到符合条件的模型，尝试清空筛选条件重新加载
      </p>
      <Button class="btn-press mt-5" @click="clearFilters">
        <SlidersHorizontal class="h-4 w-4 mr-2" />
        清空筛选
      </Button>
    </div>

    <!-- 模型网格 -->
    <div
      v-else
      class="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 animate-fade-in"
    >
      <ModelCard v-for="model in models" :key="model.id" :model="model" />
    </div>

    <!-- 桌面端分页 -->
    <div
      v-if="totalPages > 1"
      class="mt-6 sm:mt-8 hidden sm:flex items-center justify-center gap-2"
    >
      <Button
        variant="outline"
        size="sm"
        :disabled="page === 1"
        class="btn-press"
        @click="prevPage"
      >
        <ChevronLeft class="h-4 w-4 mr-1" />
        上一页
      </Button>
      <span class="text-sm text-muted-foreground px-2">
        {{ page }} / {{ totalPages }}
      </span>
      <Button
        variant="outline"
        size="sm"
        :disabled="page === totalPages"
        class="btn-press"
        @click="nextPage"
      >
        下一页
        <ChevronRight class="h-4 w-4 ml-1" />
      </Button>
    </div>

    <!-- 移动端加载更多 -->
    <div
      v-if="totalPages > 1"
      class="mt-6 sm:hidden flex flex-col items-center gap-3"
    >
      <Button
        variant="outline"
        class="btn-press w-full max-w-xs"
        :disabled="page >= totalPages || loadingMore"
        :loading="loadingMore"
        @click="loadMore"
      >
        {{ page >= totalPages ? '没有更多了' : '加载更多' }}
      </Button>
      <span class="text-xs text-muted-foreground">
        {{ page }} / {{ totalPages }}
      </span>
    </div>
  </div>
</template>
