<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { modelApi, modelVersionApi, modelImageApi, fileApi } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Model, ModelVersion, ModelImage } from '@/types'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Download, Heart, Edit, Trash2, User, Calendar, Tag, ArrowLeft, History, Plus, Check, ChevronLeft, ChevronRight, X } from 'lucide-vue-next'
import { getModelImageUrl, getAvatarUrl } from '@/utils/image'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const model = ref<Model | null>(null)
const loading = ref(true)
const isFavorited = ref(false)
const favoriteCount = ref(0)
const versions = ref<ModelVersion[]>([])
const versionsLoading = ref(false)
const showVersions = ref(false)
const images = ref<ModelImage[]>([])
const imagesLoading = ref(false)
const currentImageIndex = ref(0)

const loginPromptDialog = ref(false)
const actionMessage = ref('')
const actionMessageType = ref<'success' | 'error'>('success')

const isOwner = computed(() => {
  return authStore.user?.id === model.value?.user_id
})

const hasGalleryImages = computed(() => images.value.length > 0)

function showActionMessage(msg: string, type: 'success' | 'error' = 'success') {
  actionMessage.value = msg
  actionMessageType.value = type
  setTimeout(() => { actionMessage.value = '' }, 3000)
}

async function fetchModel() {
  loading.value = true
  try {
    const response = await modelApi.getById(route.params.id as string)
    const data = response.data.data!
    model.value = data.model
    isFavorited.value = data.is_favorited
    favoriteCount.value = data.favorite_count
    fetchImages()
  } catch (error) {
    console.error('Failed to fetch model:', error)
    router.push('/')
  } finally {
    loading.value = false
  }
}

async function fetchImages() {
  if (!model.value) return
  imagesLoading.value = true
  try {
    const response = await modelImageApi.list(model.value.id)
    images.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch images:', error)
  } finally {
    imagesLoading.value = false
  }
}

async function fetchVersions() {
  if (!model.value) return
  versionsLoading.value = true
  try {
    const response = await modelVersionApi.list(model.value.id)
    versions.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch versions:', error)
  } finally {
    versionsLoading.value = false
  }
}

async function handleDownload() {
  if (!model.value) return
  
  try {
    const res = await modelApi.download(model.value.id)
    const data = res.data.data
    if (data) {
      const downloadUrl = (data as { download_url?: string; file_path?: string }).download_url
      if (downloadUrl) {
        const isDev = import.meta.env.DEV
        const baseUrl = isDev ? '' : 'https://api.ysmmc.cn'
        window.open(`${baseUrl}${downloadUrl}`, '_blank')
        fetchModel()
      }
    }
  } catch (error) {
    console.error('Failed to download:', error)
  }
}

async function handleVersionDownload(version: ModelVersion) {
  if (!model.value) return
  
  try {
    const res = await modelVersionApi.download(model.value.id, version.id)
    const data = res.data.data
    if (data) {
      const downloadUrl = (data as { download_url?: string }).download_url
      if (downloadUrl) {
        const isDev = import.meta.env.DEV
        const baseUrl = isDev ? '' : 'https://api.ysmmc.cn'
        window.open(`${baseUrl}${downloadUrl}`, '_blank')
        fetchVersions()
      }
    }
  } catch (error) {
    console.error('Failed to download version:', error)
  }
}

async function toggleFavorite() {
  if (!model.value) return

  if (!authStore.isAuthenticated) {
    loginPromptDialog.value = true
    return
  }

  try {
    if (isFavorited.value) {
      await modelApi.removeFavorite(model.value.id)
      favoriteCount.value--
      showActionMessage('已取消收藏')
    } else {
      await modelApi.addFavorite(model.value.id)
      favoriteCount.value++
      showActionMessage('收藏成功')
    }
    isFavorited.value = !isFavorited.value
  } catch (error: any) {
    showActionMessage(error.response?.data?.message || '操作失败', 'error')
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

async function setCurrentVersion(version: ModelVersion) {
  if (!model.value) return
  
  try {
    await modelVersionApi.setCurrent(model.value.id, version.id)
    showActionMessage('已设置为当前版本')
    await fetchModel()
    await fetchVersions()
  } catch (error: any) {
    showActionMessage(error.response?.data?.message || '操作失败', 'error')
  }
}

async function deleteVersion(version: ModelVersion) {
  if (!model.value || !confirm('确定要删除这个版本吗？')) return
  
  try {
    await modelVersionApi.delete(model.value.id, version.id)
    showActionMessage('版本已删除')
    await fetchVersions()
  } catch (error: any) {
    showActionMessage(error.response?.data?.message || '删除失败', 'error')
  }
}

async function deleteImage(image: ModelImage) {
  if (!model.value || !confirm('确定要删除这张展示图吗？')) return
  
  try {
    await modelImageApi.delete(model.value.id, image.file_id)
    showActionMessage('展示图已删除')
    await fetchImages()
  } catch (error: any) {
    showActionMessage(error.response?.data?.message || '删除失败', 'error')
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

function toggleVersions() {
  showVersions.value = !showVersions.value
  if (showVersions.value && versions.value.length === 0) {
    fetchVersions()
  }
}

function prevImage() {
  if (currentImageIndex.value > 0) {
    currentImageIndex.value--
  } else {
    currentImageIndex.value = images.value.length - 1
  }
}

function nextImage() {
  if (currentImageIndex.value < images.value.length - 1) {
    currentImageIndex.value++
  } else {
    currentImageIndex.value = 0
  }
}

function getImageUrl(fileId: string) {
  return fileApi.getUrl(fileId)
}

onMounted(fetchModel)
</script>

<template>
  <div class="mx-auto max-w-4xl px-4 py-8">
    <Button variant="ghost" size="sm" class="mb-4" @click="router.back()">
      <ArrowLeft class="mr-2 h-4 w-4" />
      返回
    </Button>

    <div v-if="actionMessage" 
      class="mb-4 rounded-md p-3 text-sm"
      :class="actionMessageType === 'success' ? 'bg-green-500/10 text-green-600' : 'bg-destructive/10 text-destructive'">
      {{ actionMessage }}
    </div>

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
                <AvatarImage :src="getAvatarUrl(model.user?.avatar_id, model.user?.avatar_url) || undefined">
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
            v-if="model.image_id || model.image_url"
            :src="getModelImageUrl(model.image_id, model.image_url)"
            :alt="model.title"
            class="h-full w-full object-contain"
          />
          <div v-else class="flex h-full items-center justify-center text-muted-foreground">
            无预览图
          </div>
        </div>
      </Card>

      <Card v-if="hasGalleryImages" class="mb-6">
        <CardHeader>
          <CardTitle>展示图</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="relative">
            <div class="aspect-video w-full bg-muted rounded-md overflow-hidden">
              <img
                :src="getImageUrl(images[currentImageIndex]!.file_id)"
                :alt="`展示图 ${currentImageIndex + 1}`"
                class="h-full w-full object-contain"
              />
            </div>
            <Button
              v-if="images.length > 1"
              variant="outline"
              size="icon"
              class="absolute left-2 top-1/2 -translate-y-1/2"
              @click="prevImage"
            >
              <ChevronLeft class="h-4 w-4" />
            </Button>
            <Button
              v-if="images.length > 1"
              variant="outline"
              size="icon"
              class="absolute right-2 top-1/2 -translate-y-1/2"
              @click="nextImage"
            >
              <ChevronRight class="h-4 w-4" />
            </Button>
          </div>
          <div class="mt-4 flex gap-2 overflow-x-auto pb-2">
            <div
              v-for="(img, index) in images"
              :key="img.id"
              class="relative group flex-shrink-0"
            >
              <img
                :src="getImageUrl(img.file_id)"
                class="h-16 w-16 rounded-md object-cover cursor-pointer transition-opacity"
                :class="index === currentImageIndex ? 'ring-2 ring-primary' : 'opacity-60 hover:opacity-100'"
                @click="currentImageIndex = index"
              />
              <Button
                v-if="isOwner || authStore.isAdmin"
                variant="destructive"
                size="icon"
                class="absolute -right-1 -top-1 h-5 w-5 opacity-0 group-hover:opacity-100 transition-opacity"
                @click="deleteImage(img)"
              >
                <X class="h-3 w-3" />
              </Button>
            </div>
          </div>
          <p class="text-sm text-muted-foreground text-center">
            {{ currentImageIndex + 1 }} / {{ images.length }}
          </p>
        </CardContent>
      </Card>

      <div class="mb-6 flex flex-wrap items-center gap-4">
        <Button @click="handleDownload">
          <Download class="mr-2 h-4 w-4" />
          下载 ({{ model.downloads }})
        </Button>
        <Button
          :variant="isFavorited ? 'default' : 'outline'"
          @click="toggleFavorite"
        >
          <Heart class="mr-2 h-4 w-4" :class="{ 'fill-current': isFavorited }" />
          {{ isFavorited ? '已收藏' : '收藏' }} ({{ favoriteCount }})
        </Button>
        <span class="text-sm text-muted-foreground">
          文件大小: {{ formatFileSize(model.file_size) }}
        </span>
        <span v-if="model.current_version" class="text-sm text-muted-foreground">
          当前版本: {{ model.current_version.version_number }}
        </span>
        <span v-if="model.version_count > 1" class="text-sm text-muted-foreground">
          共 {{ model.version_count }} 个版本
        </span>
      </div>

      <div class="mb-6 flex flex-wrap gap-2">
        <Badge v-for="tag in model.tags" :key="tag" variant="secondary">
          <Tag class="mr-1 h-3 w-3" />
          {{ tag }}
        </Badge>
      </div>

      <Card class="mb-6">
        <CardHeader>
          <CardTitle>描述</CardTitle>
        </CardHeader>
        <CardContent>
          <p class="whitespace-pre-wrap">{{ model.description || '暂无描述' }}</p>
        </CardContent>
      </Card>

      <Card class="mb-6">
        <CardHeader class="cursor-pointer" @click="toggleVersions">
          <div class="flex items-center justify-between">
            <CardTitle class="flex items-center gap-2">
              <History class="h-5 w-5" />
              版本历史
            </CardTitle>
            <div class="flex items-center gap-2">
              <Button
                v-if="isOwner || authStore.isAdmin"
                variant="outline"
                size="sm"
                @click.stop="router.push(`/model/${model.id}/versions/new`)"
              >
                <Plus class="mr-1 h-4 w-4" />
                新版本
              </Button>
              <span class="text-sm text-muted-foreground">{{ showVersions ? '收起' : '展开' }}</span>
            </div>
          </div>
        </CardHeader>
        <CardContent v-if="showVersions">
          <div v-if="versionsLoading" class="py-4 text-center text-muted-foreground">
            加载中...
          </div>
          <div v-else-if="versions.length === 0" class="py-4 text-center text-muted-foreground">
            暂无版本记录
          </div>
          <div v-else class="space-y-4">
            <div
              v-for="version in versions"
              :key="version.id"
              class="rounded-lg border p-4"
              :class="{ 'border-primary bg-primary/5': version.is_current }"
            >
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center gap-2">
                    <span class="font-semibold">{{ version.version_number }}</span>
                    <Badge v-if="version.is_current" variant="default" class="text-xs">
                      当前版本
                    </Badge>
                  </div>
                  <p v-if="version.description" class="mt-1 text-sm text-muted-foreground">
                    {{ version.description }}
                  </p>
                  <div class="mt-2 flex items-center gap-4 text-xs text-muted-foreground">
                    <span>{{ formatDate(version.created_at) }}</span>
                    <span>{{ formatFileSize(version.file_size) }}</span>
                    <span>{{ version.downloads }} 次下载</span>
                  </div>
                  <p v-if="version.changelog" class="mt-2 whitespace-pre-wrap text-sm">
                    {{ version.changelog }}
                  </p>
                </div>
                <div class="flex items-center gap-2">
                  <Button variant="outline" size="sm" @click="handleVersionDownload(version)">
                    <Download class="mr-1 h-3 w-3" />
                    下载
                  </Button>
                  <Button
                    v-if="!version.is_current && (isOwner || authStore.isAdmin)"
                    variant="outline"
                    size="sm"
                    @click="setCurrentVersion(version)"
                  >
                    <Check class="mr-1 h-3 w-3" />
                    设为当前
                  </Button>
                  <Button
                    v-if="!version.is_current && (isOwner || authStore.isAdmin)"
                    variant="destructive"
                    size="sm"
                    @click="deleteVersion(version)"
                  >
                    <Trash2 class="h-3 w-3" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
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

    <Dialog v-model:open="loginPromptDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>请先登录</DialogTitle>
          <DialogDescription>收藏功能需要登录后才能使用</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="loginPromptDialog = false">取消</Button>
          <Button @click="router.push('/login')">去登录</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
