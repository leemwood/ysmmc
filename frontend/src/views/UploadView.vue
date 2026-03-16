<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { modelApi, uploadApi, modelImageApi } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Upload, X, Loader2, GripVertical, ImageIcon } from 'lucide-vue-next'

const router = useRouter()

const MAX_FILE_SIZE = 100 * 1024 * 1024
const MAX_IMAGE_SIZE = 5 * 1024 * 1024
const MAX_GALLERY_IMAGES = 10

const title = ref('')
const description = ref('')
const tags = ref('')
const isPublic = ref(true)
const file = ref<File | null>(null)
const cardImageFile = ref<File | null>(null)
const cardImagePreview = ref('')
const galleryFiles = ref<{ file: File; preview: string; id: string }[]>([])
const loading = ref(false)
const error = ref('')
const uploadProgress = ref(0)

const canAddGalleryImage = computed(() => galleryFiles.value.length < MAX_GALLERY_IMAGES)

function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const selectedFile = target.files[0]
    
    if (selectedFile.size > MAX_FILE_SIZE) {
      error.value = '模型文件大小不能超过100MB'
      return
    }
    
    file.value = selectedFile
    error.value = ''
  }
}

function handleCardImageChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const selectedFile = target.files[0]
    
    if (selectedFile.size > MAX_IMAGE_SIZE) {
      error.value = '图片大小不能超过5MB'
      return
    }
    
    if (!selectedFile.type.startsWith('image/')) {
      error.value = '请选择有效的图片文件'
      return
    }
    
    cardImageFile.value = selectedFile
    const reader = new FileReader()
    reader.onload = (e) => {
      cardImagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(cardImageFile.value)
    error.value = ''
  }
}

function removeCardImage() {
  cardImageFile.value = null
  cardImagePreview.value = ''
}

function handleGalleryImagesChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files) {
    const remainingSlots = MAX_GALLERY_IMAGES - galleryFiles.value.length
    const filesToAdd = Array.from(target.files).slice(0, remainingSlots)
    
    for (const selectedFile of filesToAdd) {
      if (selectedFile.size > MAX_IMAGE_SIZE) {
        error.value = `图片 ${selectedFile.name} 大小超过5MB，已跳过`
        continue
      }
      
      if (!selectedFile.type.startsWith('image/')) {
        error.value = `文件 ${selectedFile.name} 不是有效的图片，已跳过`
        continue
      }
      
      const reader = new FileReader()
      reader.onload = (e) => {
        galleryFiles.value.push({
          file: selectedFile,
          preview: e.target?.result as string,
          id: crypto.randomUUID()
        })
      }
      reader.readAsDataURL(selectedFile)
    }
    
    error.value = ''
  }
}

function removeGalleryImage(id: string) {
  galleryFiles.value = galleryFiles.value.filter(img => img.id !== id)
}

function moveGalleryImage(fromIndex: number, toIndex: number) {
  if (toIndex < 0 || toIndex >= galleryFiles.value.length) return
  const items = galleryFiles.value.splice(fromIndex, 1)
  if (items.length > 0) {
    galleryFiles.value.splice(toIndex, 0, items[0]!)
  }
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

async function handleSubmit() {
  error.value = ''

  if (!title.value) {
    error.value = '请输入标题'
    return
  }

  if (!file.value) {
    error.value = '请选择模型文件'
    return
  }

  if (!cardImageFile.value) {
    error.value = '请上传卡片预览图'
    return
  }

  loading.value = true
  uploadProgress.value = 0

  try {
    uploadProgress.value = 10
    const modelUpload = await uploadApi.uploadModel(file.value)
    const modelData = modelUpload.data.data!

    uploadProgress.value = 30
    let cardImageId: string | undefined
    if (cardImageFile.value) {
      const imageUpload = await uploadApi.uploadImage(cardImageFile.value)
      cardImageId = imageUpload.data.data?.file_id
    }

    uploadProgress.value = 50
    const modelResponse = await modelApi.create({
      title: title.value,
      description: description.value,
      file_path: modelData.file_path,
      file_size: modelData.file_size,
      image_id: cardImageId,
      tags: tags.value.split(',').map(t => t.trim()).filter(Boolean),
      is_public: isPublic.value,
    })
    const modelId = modelResponse.data.data!.id

    if (galleryFiles.value.length > 0) {
      uploadProgress.value = 60
      const totalGallery = galleryFiles.value.length
      for (let i = 0; i < totalGallery; i++) {
        const galleryItem = galleryFiles.value[i]
        if (galleryItem) {
          const galleryUpload = await uploadApi.uploadImage(galleryItem.file)
          const galleryImageId = galleryUpload.data.data?.file_id
          if (galleryImageId) {
            await modelImageApi.add(modelId, galleryImageId)
          }
        }
        uploadProgress.value = 60 + Math.floor((i + 1) / totalGallery * 30)
      }
    }

    uploadProgress.value = 100
    router.push('/profile')
  } catch (err: any) {
    error.value = err.response?.data?.message || '上传失败，请重试'
  } finally {
    loading.value = false
    uploadProgress.value = 0
  }
}
</script>

<template>
  <div class="mx-auto max-w-2xl px-4 py-6 sm:py-8">
    <Card>
      <CardHeader>
        <CardTitle class="text-xl sm:text-2xl">上传模型</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-5 sm:space-y-6">
          <div v-if="error" class="rounded-md bg-destructive/10 p-3 text-sm text-destructive animate-fade-in">
            {{ error }}
          </div>

          <div class="space-y-2">
            <Label for="title">标题 *</Label>
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
              placeholder="多个标签用逗号分隔，如：建筑,现代,简约"
            />
          </div>

          <div class="space-y-2">
            <Label for="file">模型文件 * (.ysm, .zip)</Label>
            <Input
              id="file"
              type="file"
              accept=".ysm,.zip"
              @change="handleFileChange"
              required
              class="text-sm"
            />
            <p class="text-xs sm:text-sm text-muted-foreground">
              支持 .ysm, .zip 格式，最大 100MB
            </p>
            <p v-if="file" class="text-xs sm:text-sm text-primary">
              已选择: {{ file.name }} ({{ formatFileSize(file.size) }})
            </p>
          </div>

          <div class="space-y-2">
            <Label>卡片预览图 * (建议 4:3 比例)</Label>
            <p class="text-xs sm:text-sm text-muted-foreground">
              用于模型卡片展示，建议使用 4:3 比例的图片
            </p>
            <div v-if="cardImagePreview" class="relative inline-block">
              <img :src="cardImagePreview" class="aspect-[4/3] w-48 sm:w-64 rounded-md object-cover" />
              <Button
                type="button"
                variant="destructive"
                size="icon"
                class="absolute right-2 top-2 h-6 w-6 sm:h-7 sm:w-7"
                @click="removeCardImage"
              >
                <X class="h-4 w-4" />
              </Button>
            </div>
            <Input v-else type="file" accept="image/*" @change="handleCardImageChange" class="text-sm" />
          </div>

          <div class="space-y-2">
            <Label>展示图 (可选，最多 {{ MAX_GALLERY_IMAGES }} 张)</Label>
            <p class="text-xs sm:text-sm text-muted-foreground">
              用于模型详情页展示，支持 jpg, png, gif 格式
            </p>
            
            <div v-if="galleryFiles.length > 0" class="grid grid-cols-4 sm:grid-cols-5 gap-2 mb-2">
              <div
                v-for="(img, index) in galleryFiles"
                :key="img.id"
                class="relative group aspect-square"
              >
                <img :src="img.preview" class="w-full h-full rounded-md object-cover" />
                <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity rounded-md flex items-center justify-center gap-1">
                  <Button
                    type="button"
                    variant="secondary"
                    size="icon"
                    class="h-5 w-5 sm:h-6 sm:w-6 hidden sm:flex"
                    :disabled="index === 0"
                    @click="moveGalleryImage(index, index - 1)"
                  >
                    <GripVertical class="h-3 w-3" />
                  </Button>
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    class="h-5 w-5 sm:h-6 sm:w-6"
                    @click="removeGalleryImage(img.id)"
                  >
                    <X class="h-3 w-3" />
                  </Button>
                </div>
                <span class="absolute bottom-1 left-1 bg-black/70 text-white text-[10px] sm:text-xs px-1 rounded">
                  {{ index + 1 }}
                </span>
              </div>
            </div>
            
            <div v-if="canAddGalleryImage" class="flex items-center gap-2">
              <label class="flex-1 cursor-pointer">
                <div class="border-2 border-dashed border-muted-foreground/25 rounded-md p-3 sm:p-4 text-center hover:border-primary transition-colors">
                  <ImageIcon class="mx-auto h-6 w-6 sm:h-8 sm:w-8 text-muted-foreground mb-2" />
                  <p class="text-xs sm:text-sm text-muted-foreground">点击添加展示图</p>
                  <p class="text-[10px] sm:text-xs text-muted-foreground">还可添加 {{ MAX_GALLERY_IMAGES - galleryFiles.length }} 张</p>
                </div>
                <Input type="file" accept="image/*" multiple class="hidden" @change="handleGalleryImagesChange" />
              </label>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <input
              id="isPublic"
              v-model="isPublic"
              type="checkbox"
              class="h-4 w-4 rounded border-gray-300 accent-primary"
            />
            <Label for="isPublic" class="cursor-pointer">公开显示</Label>
          </div>

          <div v-if="loading && uploadProgress > 0" class="space-y-2">
            <div class="flex justify-between text-sm">
              <span>上传进度</span>
              <span>{{ uploadProgress }}%</span>
            </div>
            <div class="h-2 bg-muted rounded-full overflow-hidden">
              <div 
                class="h-full bg-primary transition-all duration-300"
                :style="{ width: uploadProgress + '%' }"
              ></div>
            </div>
          </div>

          <Button type="submit" class="w-full btn-press" :disabled="loading">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            <Upload v-else class="mr-2 h-4 w-4" />
            {{ loading ? '上传中...' : '上传模型' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
