<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { modelApi, uploadApi, modelImageApi } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Model, ModelImage } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Save, X, ArrowLeft, Plus, ImageOff } from 'lucide-vue-next'
import { getModelImageUrl } from '@/utils/image'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const MAX_GALLERY_IMAGES = 10

const model = ref<Model | null>(null)
const title = ref('')
const description = ref('')
const tags = ref('')
const isPublic = ref(true)
const imageFile = ref<File | null>(null)
const imagePreview = ref('')
const loading = ref(false)
const fetching = ref(true)
const images = ref<ModelImage[]>([])
const imagesLoading = ref(false)
const newImageFiles = ref<{ file: File; preview: string; id: string }[]>([])
const error = ref('')
const dragOverPreview = ref(false)
const dragOverGallery = ref(false)

const canAddGalleryImage = computed(() => images.value.length + newImageFiles.value.length < MAX_GALLERY_IMAGES)

async function fetchModel() {
  try {
    const response = await modelApi.getById(route.params.id as string)
    model.value = response.data.data!.model

    if (model.value.user_id !== authStore.user?.id && !authStore.isAdmin) {
      router.push('/')
      return
    }

    title.value = model.value.title
    description.value = model.value.description || ''
    tags.value = model.value.tags?.join(', ') || ''
    isPublic.value = model.value.is_public
    imagePreview.value = getModelImageUrl(model.value.image_id, model.value.image_url) || ''

    fetchImages()
  } catch (error) {
    console.error('Failed to fetch model:', error)
    router.push('/')
  } finally {
    fetching.value = false
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

function handleImageChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    readImageFile(target.files[0])
  }
}

function readImageFile(file: File) {
  if (!file.type.startsWith('image/')) return
  imageFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
  error.value = ''
}

function removeImage() {
  imageFile.value = null
  imagePreview.value = getModelImageUrl(model.value?.image_id, model.value?.image_url) || ''
}

function handleNewGalleryImagesChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files) {
    addGalleryFiles(Array.from(target.files))
  }
}

function addGalleryFiles(files: File[]) {
  const remainingSlots = MAX_GALLERY_IMAGES - images.value.length - newImageFiles.value.length
  const filesToAdd = files.slice(0, remainingSlots)

  for (const selectedFile of filesToAdd) {
    if (selectedFile.size > 5 * 1024 * 1024) continue
    if (!selectedFile.type.startsWith('image/')) continue

    const reader = new FileReader()
    reader.onload = (e) => {
      newImageFiles.value.push({
        file: selectedFile,
        preview: e.target?.result as string,
        id: crypto.randomUUID(),
      })
    }
    reader.readAsDataURL(selectedFile)
  }
}

function removeNewImage(id: string) {
  newImageFiles.value = newImageFiles.value.filter(img => img.id !== id)
}

async function deleteExistingImage(image: ModelImage) {
  if (!model.value || !confirm('确定要删除这张展示图吗？')) return

  try {
    await modelImageApi.delete(model.value.id, image.file_id)
    images.value = images.value.filter(img => img.id !== image.id)
  } catch (error) {
    console.error('Failed to delete image:', error)
  }
}

function handlePreviewDragOver(e: DragEvent) {
  e.preventDefault()
  dragOverPreview.value = true
}

function handlePreviewDragLeave(e: DragEvent) {
  e.preventDefault()
  dragOverPreview.value = false
}

function handlePreviewDrop(e: DragEvent) {
  e.preventDefault()
  dragOverPreview.value = false
  const files = e.dataTransfer?.files
  if (files && files[0]) {
    readImageFile(files[0])
  }
}

function handleGalleryDragOver(e: DragEvent) {
  e.preventDefault()
  dragOverGallery.value = true
}

function handleGalleryDragLeave(e: DragEvent) {
  e.preventDefault()
  dragOverGallery.value = false
}

function handleGalleryDrop(e: DragEvent) {
  e.preventDefault()
  dragOverGallery.value = false
  const files = e.dataTransfer?.files
  if (files) {
    addGalleryFiles(Array.from(files))
  }
}

async function handleSubmit() {
  if (!title.value) {
    error.value = '请输入模型标题'
    return
  }

  loading.value = true
  error.value = ''

  try {
    let imageId: string | undefined = model.value?.image_id || undefined
    if (imageFile.value) {
      const imageUpload = await uploadApi.uploadImage(imageFile.value)
      imageId = imageUpload.data.data?.file_id
    }

    await modelApi.update(route.params.id as string, {
      title: title.value,
      description: description.value,
      tags: tags.value.split(',').map(t => t.trim()).filter(Boolean),
      is_public: isPublic.value,
      image_id: imageId,
    })

    for (const newImage of newImageFiles.value) {
      const upload = await uploadApi.uploadImage(newImage.file)
      const newImageId = upload.data.data?.file_id
      if (newImageId && model.value) {
        await modelImageApi.add(model.value.id, newImageId)
      }
    }

    router.push(`/model/${route.params.id}`)
  } catch (err: any) {
    error.value = err.response?.data?.message || '保存失败，请重试'
  } finally {
    loading.value = false
  }
}

onMounted(fetchModel)
</script>

<template>
  <div class="mx-auto max-w-2xl px-4 py-6 sm:py-8">
    <Button variant="ghost" size="sm" class="focus-ring btn-press mb-4" @click="router.back()">
      <ArrowLeft class="mr-2 h-4 w-4" />
      返回
    </Button>

    <div v-if="fetching" class="space-y-4">
      <Skeleton class="h-8 w-2/3" />
      <Skeleton class="h-64 w-full" />
      <Skeleton class="h-32 w-full" />
    </div>

    <template v-else-if="model">
      <Card class="card-hover">
        <CardHeader>
          <CardTitle class="text-xl sm:text-2xl">编辑模型</CardTitle>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleSubmit" class="space-y-6">
            <Alert v-if="error" variant="destructive" class="animate-shake">
              <AlertDescription>{{ error }}</AlertDescription>
            </Alert>

            <div class="space-y-2">
              <Label for="title">标题</Label>
              <Input id="title" v-model="title" placeholder="请输入模型标题" required class="h-11" />
            </div>

            <div class="space-y-2">
              <Label for="description">描述</Label>
              <Textarea
                id="description"
                v-model="description"
                placeholder="请输入模型描述"
                :rows="4"
                class="min-h-[100px]"
              />
            </div>

            <div class="space-y-2">
              <Label for="tags">标签</Label>
              <Input
                id="tags"
                v-model="tags"
                placeholder="多个标签用逗号分隔"
                class="h-11"
              />
              <p class="text-xs text-muted-foreground">使用逗号分隔多个标签，例如：建筑, 红石, 生存</p>
            </div>

            <div class="space-y-2">
              <Label>卡片预览图 (建议 4:3 比例，最大 5MB)</Label>
              <div v-if="imagePreview" class="relative inline-block">
                <img :src="imagePreview" class="aspect-[4/3] w-64 rounded-lg object-cover" />
                <Button
                  type="button"
                  variant="destructive"
                  size="icon"
                  class="focus-ring btn-press absolute right-2 top-2 h-7 w-7"
                  aria-label="移除预览图"
                  @click="removeImage"
                >
                  <X class="h-4 w-4" />
                </Button>
              </div>
              <label
                v-else
                class="flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-6 text-center transition-colors"
                :class="dragOverPreview ? 'border-primary bg-primary/5' : 'border-muted-foreground/25 hover:border-primary hover:bg-primary/5'"
                @dragover="handlePreviewDragOver"
                @dragleave="handlePreviewDragLeave"
                @drop="handlePreviewDrop"
              >
                <ImageOff class="h-8 w-8 text-muted-foreground" />
                <p class="mt-2 text-sm text-muted-foreground">点击或拖拽上传预览图</p>
                <p class="text-xs text-muted-foreground">支持 JPG、PNG、GIF、WEBP</p>
                <Input type="file" accept="image/*" class="hidden" @change="handleImageChange" />
              </label>
            </div>

            <div class="space-y-2">
              <Label>展示图 (最多 {{ MAX_GALLERY_IMAGES }} 张，单张最大 5MB)</Label>

              <div class="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-5">
                <div
                  v-for="img in images"
                  :key="img.id"
                  class="group relative aspect-square"
                >
                  <img :src="getModelImageUrl(img.file_id)" class="h-full w-full rounded-lg object-cover" />
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    class="focus-ring btn-press absolute -right-1 -top-1 h-6 w-6 opacity-100 transition-opacity sm:opacity-0 sm:group-hover:opacity-100"
                    aria-label="删除展示图"
                    @click="deleteExistingImage(img)"
                  >
                    <X class="h-3 w-3" />
                  </Button>
                </div>

                <div
                  v-for="img in newImageFiles"
                  :key="img.id"
                  class="group relative aspect-square"
                >
                  <img :src="img.preview" class="h-full w-full rounded-lg object-cover" />
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    class="focus-ring btn-press absolute -right-1 -top-1 h-6 w-6 opacity-100 transition-opacity sm:opacity-0 sm:group-hover:opacity-100"
                    aria-label="移除新展示图"
                    @click="removeNewImage(img.id)"
                  >
                    <X class="h-3 w-3" />
                  </Button>
                  <span class="absolute bottom-1 left-1 rounded bg-green-500 px-1 text-xs text-white">新</span>
                </div>
              </div>

              <label
                v-if="canAddGalleryImage"
                class="flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-4 text-center transition-colors"
                :class="dragOverGallery ? 'border-primary bg-primary/5' : 'border-muted-foreground/25 hover:border-primary hover:bg-primary/5'"
                @dragover="handleGalleryDragOver"
                @dragleave="handleGalleryDragLeave"
                @drop="handleGalleryDrop"
              >
                <Plus class="mx-auto h-6 w-6 text-muted-foreground" />
                <p class="text-xs text-muted-foreground">添加展示图</p>
                <Input type="file" accept="image/*" multiple class="hidden" @change="handleNewGalleryImagesChange" />
              </label>
            </div>

            <div class="flex items-center gap-2">
              <input
                id="isPublic"
                v-model="isPublic"
                type="checkbox"
                class="h-4 w-4 rounded border-border text-primary focus:ring-2 focus:ring-ring"
              />
              <Label for="isPublic">公开显示</Label>
            </div>

            <Alert v-if="model?.update_status === 'pending_review'" variant="default" class="border-yellow-500/50 text-yellow-600 dark:text-yellow-400 [&>svg]:text-yellow-600 dark:[&>svg]:text-yellow-400">
              <AlertDescription>当前有修改正在审核中</AlertDescription>
            </Alert>

            <Button type="submit" class="focus-ring btn-press w-full sm:w-auto" :loading="loading">
              <Save v-if="!loading" class="mr-2 h-4 w-4" />
              {{ loading ? '保存中...' : '保存修改' }}
            </Button>
          </form>
        </CardContent>
      </Card>
    </template>
  </div>
</template>
