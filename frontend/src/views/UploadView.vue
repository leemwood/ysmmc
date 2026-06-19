<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { modelApi, uploadApi, modelImageApi } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Accordion, AccordionItem, AccordionTrigger, AccordionContent } from '@/components/ui/accordion'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Separator } from '@/components/ui/separator'
import { Upload, X, ImageIcon, FileArchive, ChevronUp, ChevronDown, Trash2, AlertCircle } from 'lucide-vue-next'

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
const fieldErrors = ref<Record<string, boolean>>({})
const uploadStage = ref('')
const isDraggingModel = ref(false)
const isDraggingCard = ref(false)

const canAddGalleryImage = computed(() => galleryFiles.value.length < MAX_GALLERY_IMAGES)

function clearFieldErrors() {
  fieldErrors.value = {}
}

function setFieldError(field: string) {
  fieldErrors.value[field] = true
}

function handleModelFile(selectedFile: File) {
  if (selectedFile.size > MAX_FILE_SIZE) {
    error.value = '模型文件大小不能超过100MB'
    setFieldError('file')
    return
  }

  if (!selectedFile.name.match(/\.(ysm|zip)$/i)) {
    error.value = '仅支持 .ysm, .zip 格式的模型文件'
    setFieldError('file')
    return
  }

  file.value = selectedFile
  error.value = ''
  clearFieldErrors()
}

function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    handleModelFile(target.files[0])
    target.value = ''
  }
}

function handleModelDrop(e: DragEvent) {
  e.preventDefault()
  isDraggingModel.value = false
  const droppedFile = e.dataTransfer?.files?.[0]
  if (droppedFile) {
    handleModelFile(droppedFile)
  }
}

function handleModelDragOver(e: DragEvent) {
  e.preventDefault()
  isDraggingModel.value = true
}

function handleModelDragLeave(e: DragEvent) {
  e.preventDefault()
  isDraggingModel.value = false
}

function handleCardImage(selectedFile: File) {
  if (selectedFile.size > MAX_IMAGE_SIZE) {
    error.value = '图片大小不能超过5MB'
    setFieldError('cardImage')
    return
  }

  if (!selectedFile.type.startsWith('image/')) {
    error.value = '请选择有效的图片文件'
    setFieldError('cardImage')
    return
  }

  cardImageFile.value = selectedFile
  const reader = new FileReader()
  reader.onload = (e) => {
    cardImagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(cardImageFile.value)
  error.value = ''
  clearFieldErrors()
}

function handleCardImageChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    handleCardImage(target.files[0])
    target.value = ''
  }
}

function handleCardImageDrop(e: DragEvent) {
  e.preventDefault()
  isDraggingCard.value = false
  const droppedFile = e.dataTransfer?.files?.[0]
  if (droppedFile) {
    handleCardImage(droppedFile)
  }
}

function handleCardImageDragOver(e: DragEvent) {
  e.preventDefault()
  isDraggingCard.value = true
}

function handleCardImageDragLeave(e: DragEvent) {
  e.preventDefault()
  isDraggingCard.value = false
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

    target.value = ''
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
  clearFieldErrors()

  if (!title.value) {
    error.value = '请输入标题'
    setFieldError('title')
    return
  }

  if (!file.value) {
    error.value = '请选择模型文件'
    setFieldError('file')
    return
  }

  if (!cardImageFile.value) {
    error.value = '请上传卡片预览图'
    setFieldError('cardImage')
    return
  }

  loading.value = true
  uploadProgress.value = 0
  uploadStage.value = '准备上传...'

  try {
    uploadProgress.value = 10
    uploadStage.value = '上传模型文件...'
    const modelUpload = await uploadApi.uploadModel(file.value)
    const modelData = modelUpload.data.data!

    uploadProgress.value = 30
    uploadStage.value = '上传卡片预览图...'
    let cardImageId: string | undefined
    if (cardImageFile.value) {
      const imageUpload = await uploadApi.uploadImage(cardImageFile.value)
      cardImageId = imageUpload.data.data?.file_id
    }

    uploadProgress.value = 50
    uploadStage.value = '创建模型记录...'
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
      uploadStage.value = '上传展示图...'
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
    uploadStage.value = '完成'
    router.push('/profile')
  } catch (err: any) {
    error.value = err.response?.data?.message || '上传失败，请重试'
  } finally {
    loading.value = false
    uploadProgress.value = 0
    uploadStage.value = ''
  }
}
</script>

<template>
  <div class="mx-auto max-w-2xl px-4 py-6 sm:py-8 pb-28 md:pb-8">
    <Card>
      <CardHeader class="pb-4">
        <CardTitle class="text-xl sm:text-2xl">上传模型</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-5 sm:space-y-6">
          <Alert v-if="error" variant="destructive" class="animate-fade-in">
            <AlertCircle class="h-4 w-4" />
            <AlertDescription>{{ error }}</AlertDescription>
          </Alert>

          <Accordion type="multiple" :default-value="['basic', 'model', 'preview', 'gallery']" class="w-full">
            <AccordionItem value="basic" class="border-border">
              <AccordionTrigger class="text-base sm:text-lg font-semibold py-3">
                基本信息
              </AccordionTrigger>
              <AccordionContent>
                <div class="space-y-4">
                  <div class="space-y-2">
                    <Label for="title">标题 <span class="text-destructive">*</span></Label>
                    <Input
                      id="title"
                      v-model="title"
                      placeholder="请输入模型标题"
                      :error="fieldErrors.title"
                      required
                    />
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

                  <div class="flex items-center gap-2">
                    <input
                      id="isPublic"
                      v-model="isPublic"
                      type="checkbox"
                      class="h-4 w-4 rounded border-gray-300 accent-primary"
                    />
                    <Label for="isPublic" class="cursor-pointer">公开显示</Label>
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="model" class="border-border">
              <AccordionTrigger class="text-base sm:text-lg font-semibold py-3">
                模型文件
              </AccordionTrigger>
              <AccordionContent>
                <div class="space-y-3">
                  <Label for="file-input">模型文件 <span class="text-destructive">*</span></Label>
                  <label
                    for="file-input"
                    class="relative flex min-h-[140px] sm:min-h-[160px] cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed p-4 text-center transition-all"
                    :class="[
                      fieldErrors.file
                        ? 'border-destructive bg-destructive/5'
                        : isDraggingModel
                          ? 'border-primary bg-primary/5'
                          : 'border-muted-foreground/25 hover:border-primary hover:bg-primary/5'
                    ]"
                    @dragover="handleModelDragOver"
                    @drop="handleModelDrop"
                    @dragleave="handleModelDragLeave"
                  >
                    <input
                      id="file-input"
                      ref="modelFileInput"
                      type="file"
                      accept=".ysm,.zip"
                      class="sr-only"
                      @change="handleFileChange"
                    />
                    <template v-if="file">
                      <FileArchive class="mb-2 h-8 w-8 text-primary sm:h-10 sm:w-10" />
                      <p class="text-sm font-medium text-foreground break-all px-2">{{ file.name }}</p>
                      <p class="text-xs text-muted-foreground">{{ formatFileSize(file.size) }}</p>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        class="mt-3"
                        @click.stop="file = null"
                      >
                        重新选择
                      </Button>
                    </template>
                    <template v-else>
                      <Upload class="mb-2 h-8 w-8 text-muted-foreground sm:h-10 sm:w-10" />
                      <p class="text-sm font-medium text-foreground">点击或拖拽上传模型文件</p>
                      <p class="mt-1 text-xs text-muted-foreground">支持 .ysm, .zip 格式，最大 100MB</p>
                    </template>
                  </label>
                </div>
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="preview" class="border-border">
              <AccordionTrigger class="text-base sm:text-lg font-semibold py-3">
                卡片预览图
              </AccordionTrigger>
              <AccordionContent>
                <div class="space-y-3">
                  <Label for="card-image-input">卡片预览图 <span class="text-destructive">*</span> <span class="text-muted-foreground font-normal">(建议 4:3 比例)</span></Label>
                  <div v-if="cardImagePreview" class="relative inline-block w-full max-w-[16rem]">
                    <img :src="cardImagePreview" class="aspect-[4/3] w-full rounded-xl object-cover" alt="卡片预览图" />
                    <Button
                      type="button"
                      variant="destructive"
                      size="icon"
                      class="absolute right-2 top-2 h-8 w-8"
                      @click="removeCardImage"
                    >
                      <X class="h-4 w-4" />
                    </Button>
                  </div>
                  <label
                    v-else
                    for="card-image-input"
                    class="relative flex aspect-[4/3] w-full max-w-[16rem] cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed p-4 text-center transition-all"
                    :class="[
                      fieldErrors.cardImage
                        ? 'border-destructive bg-destructive/5'
                        : isDraggingCard
                          ? 'border-primary bg-primary/5'
                          : 'border-muted-foreground/25 hover:border-primary hover:bg-primary/5'
                    ]"
                    @dragover="handleCardImageDragOver"
                    @drop="handleCardImageDrop"
                    @dragleave="handleCardImageDragLeave"
                  >
                    <input
                      id="card-image-input"
                      ref="cardImageInput"
                      type="file"
                      accept="image/*"
                      class="sr-only"
                      @change="handleCardImageChange"
                    />
                    <ImageIcon class="mb-2 h-8 w-8 text-muted-foreground sm:h-10 sm:w-10" />
                    <p class="text-sm font-medium text-foreground">点击或拖拽上传预览图</p>
                    <p class="mt-1 text-xs text-muted-foreground">建议 4:3 比例，最大 5MB</p>
                  </label>
                </div>
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="gallery" class="border-border">
              <AccordionTrigger class="text-base sm:text-lg font-semibold py-3">
                展示图 <span class="ml-1 text-xs font-normal text-muted-foreground">(可选)</span>
              </AccordionTrigger>
              <AccordionContent>
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <Label>展示图 (可选，最多 {{ MAX_GALLERY_IMAGES }} 张)</Label>
                    <span class="text-xs text-muted-foreground">{{ galleryFiles.length }}/{{ MAX_GALLERY_IMAGES }}</span>
                  </div>
                  <p class="text-xs text-muted-foreground">
                    用于模型详情页展示，支持 jpg, png, gif 格式
                  </p>

                  <div v-if="galleryFiles.length > 0" class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-2">
                    <div
                      v-for="(img, index) in galleryFiles"
                      :key="img.id"
                      class="relative aspect-square rounded-xl overflow-hidden group"
                    >
                      <img :src="img.preview" class="w-full h-full object-cover" loading="lazy" decoding="async" />
                      <div class="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-1.5 p-1 sm:flex-row sm:gap-1">
                        <Button
                          type="button"
                          variant="secondary"
                          size="icon"
                          class="h-8 w-8 sm:h-6 sm:w-6"
                          :disabled="index === 0"
                          @click="moveGalleryImage(index, index - 1)"
                        >
                          <ChevronUp class="h-4 w-4 sm:h-3 sm:w-3" />
                        </Button>
                        <Button
                          type="button"
                          variant="secondary"
                          size="icon"
                          class="h-8 w-8 sm:h-6 sm:w-6"
                          :disabled="index === galleryFiles.length - 1"
                          @click="moveGalleryImage(index, index + 1)"
                        >
                          <ChevronDown class="h-4 w-4 sm:h-3 sm:w-3" />
                        </Button>
                        <Button
                          type="button"
                          variant="destructive"
                          size="icon"
                          class="h-8 w-8 sm:h-6 sm:w-6"
                          @click="removeGalleryImage(img.id)"
                        >
                          <Trash2 class="h-4 w-4 sm:h-3 sm:w-3" />
                        </Button>
                      </div>
                      <span class="absolute bottom-1 left-1 bg-black/70 text-white text-[10px] sm:text-xs px-1 rounded">
                        {{ index + 1 }}
                      </span>
                    </div>
                  </div>

                  <div v-if="canAddGalleryImage" class="flex items-center gap-2">
                    <label class="flex-1 cursor-pointer">
                      <div class="border-2 border-dashed border-muted-foreground/25 rounded-xl p-3 sm:p-4 text-center hover:border-primary hover:bg-primary/5 transition-colors">
                        <ImageIcon class="mx-auto h-6 w-6 sm:h-8 sm:w-8 text-muted-foreground mb-2" />
                        <p class="text-xs sm:text-sm text-muted-foreground">点击添加展示图</p>
                        <p class="text-[10px] sm:text-xs text-muted-foreground">还可添加 {{ MAX_GALLERY_IMAGES - galleryFiles.length }} 张</p>
                      </div>
                      <Input type="file" accept="image/*" multiple class="hidden" @change="handleGalleryImagesChange" />
                    </label>
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>
          </Accordion>

          <Separator />

          <div v-if="loading && uploadProgress > 0" class="space-y-2">
            <Progress v-model="uploadProgress" :max="100" show-value>
              <span class="text-sm text-muted-foreground">{{ uploadStage }}</span>
            </Progress>
          </div>

          <div class="hidden md:block">
            <Button type="submit" class="w-full btn-press" size="lg" :loading="loading">
              <Upload v-if="!loading" class="mr-2 h-4 w-4" />
              {{ loading ? '上传中...' : '上传模型' }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <div class="fixed bottom-0 left-0 right-0 z-40 border-t bg-background/95 p-4 backdrop-blur md:hidden pb-[max(1rem,env(safe-area-inset-bottom))]">
      <Button type="button" class="w-full btn-press" size="lg" :loading="loading" @click="handleSubmit">
        <Upload v-if="!loading" class="mr-2 h-4 w-4" />
        {{ loading ? '上传中...' : '上传模型' }}
      </Button>
    </div>
  </div>
</template>
