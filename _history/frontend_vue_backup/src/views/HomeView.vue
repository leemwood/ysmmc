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
import { Search, ChevronLeft, ChevronRight, Megaphone } from 'lucide-vue-next'

const route = useRoute()

const models = ref<Model[]>([])
const announcements = ref<Announcement[]>([])
const loading = ref(true)
const search = ref('')
const page = ref(1)
const pageSize = 12
const total = ref(0)
const totalPages = ref(0)

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
  <div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
    <div v-if="announcements.length > 0" class="mb-8">
      <Card v-for="announcement in announcements" :key="announcement.id" class="bg-primary/5 border-primary/20">
        <CardContent class="flex items-center gap-3 p-4">
          <Megaphone class="h-5 w-5 text-primary" />
          <div>
            <h4 class="font-semibold">{{ announcement.title }}</h4>
            <p class="text-sm text-muted-foreground">{{ announcement.content }}</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="mb-8 flex items-center gap-4">
      <div class="relative flex-1 max-w-md">
        <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          v-model="search"
          placeholder="搜索模型..."
          class="pl-10"
          @keyup.enter="handleSearch"
        />
      </div>
      <Button @click="handleSearch">搜索</Button>
    </div>

    <div v-if="loading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <Card v-for="i in 8" :key="i">
        <Skeleton class="aspect-video w-full" />
        <CardContent class="p-4 space-y-3">
          <Skeleton class="h-5 w-3/4" />
          <Skeleton class="h-4 w-full" />
          <Skeleton class="h-4 w-2/3" />
        </CardContent>
      </Card>
    </div>

    <div v-else-if="models.length === 0" class="text-center py-12">
      <p class="text-muted-foreground">暂无模型</p>
    </div>

    <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <ModelCard v-for="model in models" :key="model.id" :model="model" />
    </div>

    <div v-if="totalPages > 1" class="mt-8 flex items-center justify-center gap-2">
      <Button variant="outline" size="sm" :disabled="page === 1" @click="prevPage">
        <ChevronLeft class="h-4 w-4" />
        上一页
      </Button>
      <span class="text-sm text-muted-foreground">
        {{ page }} / {{ totalPages }}
      </span>
      <Button variant="outline" size="sm" :disabled="page === totalPages" @click="nextPage">
        下一页
        <ChevronRight class="h-4 w-4" />
      </Button>
    </div>
  </div>
</template>
