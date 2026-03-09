<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { modelApi, uploadApi, modelImageApi, fileApi } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Model, ModelImage } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Save, X, Loader2, ArrowLeft, Plus } from 'lucide-vue-next'
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
    imageFile.value = target.files[0]
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(imageFile.value)
  }
}

function removeImage() {
  imageFile.value = null
  imagePreview.value = getModelImageUrl(model.value?.image_id, model.value?.image_url) || ''
}

function handleNewGalleryImagesChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files) {
    const remainingSlots = MAX_GALLERY_IMAGES - images.value.length - newImageFiles.value.length
    const filesToAdd = Array.from(target.files).slice(0, remainingSlots)
    
    for (const selectedFile of filesToAdd) {
      if (selectedFile.size > 5 * 1024 * 1024) continue
      if (!selectedFile.type.startsWith('image/')) continue
      
      const reader = new FileReader()
      reader.onload = (e) => {
        newImageFiles.value.push({
          file: selectedFile,
          preview: e.target?.result as string,
          id: crypto.randomUUID()
        })
      }
      reader.readAsDataURL(selectedFile)
    }
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

function getImageUrl(fileId: string) {
  return fileApi.getUrl(fileId)
}

async function handleSubmit() {
  if (!title.value) return

  loading.value = true

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
  } catch (error) {
    console.error('Failed to update model:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchModel)
</script>

<template>
  <div class="mx-auto max-w-2xl px-4 py-8">
    <Button variant="ghost" size="sm" class="mb-4" @click="router.back()">
      <ArrowLeft class="mr-2 h-4 w-4" />
      返回
    </Button>

    <Card v-if="fetching">
      <CardContent class="py-8 text-center">
        <Loader2 class="mx-auto h-8 w-8 animate-spin text-muted-foreground" />
      </CardContent>
    </Card>

    <Card v-else>
      <CardHeader>
        <CardTitle class="text-2xl">编辑模型</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div class="space-y-2">
            <Label for="title">标题</Label>
            <Input id="title" v-model="title" placeholder="请输入模型标题" required />
          </div>

          <div class="space-y-2">
            <Label for="description">描述</Label>
            <Textarea
              id="description"
              v-model="description"
              placeholder="请输入模型描述"
              :rows="4"
            />
          </div>

          <div class="space-y-2">
            <Label for="tags">标签</Label>
            <Input
              id="tags"
              v-model="tags"
              placeholder="多个标签用逗号分隔"
            />
          </div>

          <div class="space-y-2">
            <Label>卡片预览图 (建议 4:3 比例)</Label>
            <div v-if="imagePreview" class="relative inline-block">
              <img :src="imagePreview" class="aspect-[4/3] w-64 rounded-md object-cover" />
              <Button
                type="button"
                variant="destructive"
                size="icon"
                class="absolute right-2 top-2 h-6 w-6"
                @click="removeImage"
              >
                <X class="h-4 w-4" />
              </Button>
            </div>
            <Input v-else type="file" accept="image/*" @change="handleImageChange" />
          </div>

          <div class="space-y-2">
            <Label>展示图 (最多 {{ MAX_GALLERY_IMAGES }} 张)</Label>
            
            <div class="grid grid-cols-5 gap-2 mb-2">
              <div
                v-for="img in images"
                :key="img.id"
                class="relative group"
              >
                <img :src="getImageUrl(img.file_id)" class="aspect-square w-full rounded-md object-cover" />
                <Button
                  type="button"
                  variant="destructive"
                  size="icon"
                  class="absolute -right-1 -top-1 h-5 w-5 opacity-0 group-hover:opacity-100 transition-opacity"
                  @click="deleteExistingImage(img)"
                >
                  <X class="h-3 w-3" />
                </Button>
              </div>
              
              <div
                v-for="img in newImageFiles"
                :key="img.id"
                class="relative group"
              >
                <img :src="img.preview" class="aspect-square w-full rounded-md object-cover" />
                <Button
                  type="button"
                  variant="destructive"
                  size="icon"
                  class="absolute -right-1 -top-1 h-5 w-5 opacity-0 group-hover:opacity-100 transition-opacity"
                  @click="removeNewImage(img.id)"
                >
                  <X class="h-3 w-3" />
                </Button>
                <span class="absolute bottom-1 left-1 bg-green-500 text-white text-xs px-1 rounded">
                  新
                </span>
              </div>
            </div>
            
            <div v-if="canAddGalleryImage" class="flex items-center gap-2">
              <label class="flex-1 cursor-pointer">
                <div class="border-2 border-dashed border-muted-foreground/25 rounded-md p-3 text-center hover:border-primary transition-colors">
                  <Plus class="mx-auto h-6 w-6 text-muted-foreground" />
                  <p class="text-xs text-muted-foreground">添加展示图</p>
                </div>
                <Input type="file" accept="image/*" multiple class="hidden" @change="handleNewGalleryImagesChange" />
              </label>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <input
              id="isPublic"
              v-model="isPublic"
              type="checkbox"
              class="h-4 w-4 rounded border-gray-300"
            />
            <Label for="isPublic">公开显示</Label>
          </div>

          <div v-if="model?.update_status === 'pending_review'" class="rounded-md bg-yellow-500/10 p-3 text-sm text-yellow-600">
            当前有修改正在审核中
          </div>

          <Button type="submit" class="w-full" :disabled="loading">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            <Save v-else class="mr-2 h-4 w-4" />
            {{ loading ? '保存中...' : '保存修改' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
