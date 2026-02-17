<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { modelApi } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Model } from '@/types'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Download, Heart, Edit, Trash2, User, Calendar, Tag, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const model = ref<Model | null>(null)
const loading = ref(true)
const isFavorited = ref(false)
const favoriteCount = ref(0)

const isOwner = computed(() => {
  return authStore.user?.id === model.value?.user_id
})

async function fetchModel() {
  loading.value = true
  try {
    const response = await modelApi.getById(route.params.id as string)
    const data = response.data.data!
    model.value = data.model
    isFavorited.value = data.is_favorited
    favoriteCount.value = data.favorite_count
  } catch (error) {
    console.error('Failed to fetch model:', error)
    router.push('/')
  } finally {
    loading.value = false
  }
}

async function handleDownload() {
  if (!model.value) return
  
  try {
    await modelApi.download(model.value.id)
    window.open(`/uploads/${model.value.file_path.split('/').pop()}`, '_blank')
    fetchModel()
  } catch (error) {
    console.error('Failed to download:', error)
  }
}

async function toggleFavorite() {
  if (!model.value || !authStore.isAuthenticated) return

  try {
    if (isFavorited.value) {
      await modelApi.removeFavorite(model.value.id)
      favoriteCount.value--
    } else {
      await modelApi.addFavorite(model.value.id)
      favoriteCount.value++
    }
    isFavorited.value = !isFavorited.value
  } catch (error) {
    console.error('Failed to toggle favorite:', error)
  }
}

async function handleDelete() {
  if (!model.value || !confirm('确定要删除这个模型吗？')) return

  try {
    await modelApi.delete(model.value.id)
    router.push('/')
  } catch (error) {
    console.error('Failed to delete:', error)
  }
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}

function formatFileSize(bytes: number) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(fetchModel)
</script>

<template>
  <div class="mx-auto max-w-4xl px-4 py-8">
    <Button variant="ghost" size="sm" class="mb-4" @click="router.back()">
      <ArrowLeft class="mr-2 h-4 w-4" />
      返回
    </Button>

    <div v-if="loading" class="space-y-4">
      <Skeleton class="h-8 w-2/3" />
      <Skeleton class="aspect-video w-full" />
      <Skeleton class="h-32 w-full" />
    </div>

    <template v-else-if="model">
      <div class="mb-6 flex items-start justify-between">
        <div>
          <h1 class="text-3xl font-bold">{{ model.title }}</h1>
          <div class="mt-2 flex items-center gap-4 text-sm text-muted-foreground">
            <RouterLink :to="`/user/${model.user?.id}`" class="flex items-center gap-2 hover:text-foreground">
              <Avatar class="h-6 w-6">
                <AvatarImage :src="model.user?.avatar_url || undefined">
                  <User class="h-3 w-3" />
                </AvatarImage>
              </Avatar>
              {{ model.user?.username }}
            </RouterLink>
            <span class="flex items-center gap-1">
              <Calendar class="h-4 w-4" />
              {{ formatDate(model.created_at) }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <Button v-if="isOwner || authStore.isAdmin" variant="outline" size="sm" @click="router.push(`/model/${model.id}/edit`)">
            <Edit class="mr-2 h-4 w-4" />
            编辑
          </Button>
          <Button v-if="isOwner || authStore.isAdmin" variant="destructive" size="sm" @click="handleDelete">
            <Trash2 class="mr-2 h-4 w-4" />
            删除
          </Button>
        </div>
      </div>

      <Card class="mb-6">
        <div class="aspect-video w-full bg-muted">
          <img
            v-if="model.image_url"
            :src="model.image_url"
            :alt="model.title"
            class="h-full w-full object-contain"
          />
          <div v-else class="flex h-full items-center justify-center text-muted-foreground">
            无预览图
          </div>
        </div>
      </Card>

      <div class="mb-6 flex flex-wrap items-center gap-4">
        <Button @click="handleDownload">
          <Download class="mr-2 h-4 w-4" />
          下载 ({{ model.downloads }})
        </Button>
        <Button
          :variant="isFavorited ? 'default' : 'outline'"
          @click="toggleFavorite"
          :disabled="!authStore.isAuthenticated"
        >
          <Heart class="mr-2 h-4 w-4" :class="{ 'fill-current': isFavorited }" />
          {{ isFavorited ? '已收藏' : '收藏' }} ({{ favoriteCount }})
        </Button>
        <span class="text-sm text-muted-foreground">
          文件大小: {{ formatFileSize(model.file_size) }}
        </span>
      </div>

      <div class="mb-6 flex flex-wrap gap-2">
        <Badge v-for="tag in model.tags" :key="tag" variant="secondary">
          <Tag class="mr-1 h-3 w-3" />
          {{ tag }}
        </Badge>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>描述</CardTitle>
        </CardHeader>
        <CardContent>
          <p class="whitespace-pre-wrap">{{ model.description || '暂无描述' }}</p>
        </CardContent>
      </Card>

      <Card v-if="model.status === 'rejected'" class="mt-6 border-destructive">
        <CardHeader>
          <CardTitle class="text-destructive">审核未通过</CardTitle>
        </CardHeader>
        <CardContent>
          <p>{{ model.rejection_reason || '未提供原因' }}</p>
        </CardContent>
      </Card>
    </template>
  </div>
</template>
