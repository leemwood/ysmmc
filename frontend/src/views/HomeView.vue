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
import { Search, ChevronLeft, ChevronRight, Megaphone, X } from 'lucide-vue-next'

const route = useRoute()

const models = ref<Model[]>([])
const announcements = ref<Announcement[]>([])
const loading = ref(true)
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
  <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8 sm:py-8">
    <div v-if="announcements.some(a => isAnnouncementVisible(a.id))" class="mb-6 space-y-3">
      <Card 
        v-for="announcement in announcements" 
        :key="announcement.id"
        v-show="isAnnouncementVisible(announcement.id)"
        class="bg-primary/5 border-primary/20 relative overflow-hidden"
      >
        <CardContent class="flex items-start gap-3 p-4">
          <Megaphone class="h-5 w-5 text-primary flex-shrink-0 mt-0.5" />
          <div class="flex-1 min-w-0">
            <h4 class="font-semibold text-sm sm:text-base">{{ announcement.title }}</h4>
            <p class="text-xs sm:text-sm text-muted-foreground mt-1">{{ announcement.content }}</p>
          </div>
          <Button 
            variant="ghost" 
            size="icon" 
            class="h-6 w-6 flex-shrink-0"
            @click="dismissAnnouncement(announcement.id)"
          >
            <X class="h-4 w-4" />
          </Button>
        </CardContent>
      </Card>
    </div>

    <div class="mb-6 sm:mb-8">
      <form class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3" @submit.prevent="handleSearch">
        <div class="relative flex-1 max-w-full sm:max-w-md">
          <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            v-model="search"
            placeholder="搜索模型..."
            class="pl-10"
          />
        </div>
        <Button type="submit" class="btn-press">搜索</Button>
      </form>
    </div>

    <div v-if="loading" class="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <Card v-for="i in 8" :key="i">
        <Skeleton class="aspect-video w-full" />
        <CardContent class="p-4 space-y-3">
          <Skeleton class="h-5 w-3/4" />
          <Skeleton class="h-4 w-full" />
          <Skeleton class="h-4 w-2/3" />
        </CardContent>
      </Card>
    </div>

    <div v-else-if="models.length === 0" class="text-center py-12 sm:py-16">
      <p class="text-muted-foreground">暂无模型</p>
    </div>

    <div v-else class="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <ModelCard v-for="model in models" :key="model.id" :model="model" />
    </div>

    <div v-if="totalPages > 1" class="mt-6 sm:mt-8 flex items-center justify-center gap-2">
      <Button variant="outline" size="sm" :disabled="page === 1" @click="prevPage" class="btn-press">
        <ChevronLeft class="h-4 w-4 sm:mr-1" />
        <span class="hidden sm:inline">上一页</span>
      </Button>
      <span class="text-sm text-muted-foreground px-2">
        {{ page }} / {{ totalPages }}
      </span>
      <Button variant="outline" size="sm" :disabled="page === totalPages" @click="nextPage" class="btn-press">
        <span class="hidden sm:inline">下一页</span>
        <ChevronRight class="h-4 w-4 sm:ml-1" />
      </Button>
    </div>
  </div>
</template>
