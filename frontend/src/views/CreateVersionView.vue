<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { modelApi, modelVersionApi, uploadApi } from '@/lib/api'
import type { Model } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { Upload, X, ArrowLeft, FileArchive, ImageOff } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const MAX_FILE_SIZE = 100 * 1024 * 1024
const MAX_IMAGE_SIZE = 5 * 1024 * 1024

const model = ref<Model | null>(null)
const loading = ref(true)
const versionNumber = ref('')
const description = ref('')
const changelog = ref('')
const file = ref<File | null>(null)
const imageFile = ref<File | null>(null)
const imagePreview = ref('')
const submitting = ref(false)
const error = ref('')
const uploadProgress = ref(0)
const fileDragOver = ref(false)
const imageDragOver = ref(false)

const modelId = computed(() => route.params.id as string)

const versionNumberError = computed(() => {
  if (!versionNumber.value) return ''
  const regex = /^\d+\.\d+\.\d+$/
  if (!regex.test(versionNumber.value)) {
    return '版本号格式不正确，应为 x.y.z 格式（如 1.0.1）'
  }
  return ''
})

async function fetchModel() {
  loading.value = true
  try {
    const response = await modelApi.getById(modelId.value)
    model.value = response.data.data!.model
  } catch (err) {
    console.error('Failed to fetch model:', err)
    router.push('/')
  } finally {
    loading.value = false
  }
}

function validateModelFile(selectedFile: File): boolean {
  if (selectedFile.size > MAX_FILE_SIZE) {
    error.value = '模型文件大小不能超过100MB'
    return false
  }
  return true
}

function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const selectedFile = target.files[0]
    if (validateModelFile(selectedFile)) {
      file.value = selectedFile
      error.value = ''
    }
  }
}

function handleFileDrop(e: DragEvent) {
  e.preventDefault()
  fileDragOver.value = false
  const files = e.dataTransfer?.files
  if (files && files[0]) {
    const selectedFile = files[0]
    if (validateModelFile(selectedFile)) {
      file.value = selectedFile
      error.value = ''
    }
  }
}

function removeFile() {
  file.value = null
}

function validateImageFile(selectedFile: File): boolean {
  if (selectedFile.size > MAX_IMAGE_SIZE) {
    error.value = '图片大小不能超过5MB'
    return false
  }
  if (!selectedFile.type.startsWith('image/')) {
    error.value = '请选择有效的图片文件'
    return false
  }
  return true
}

function handleImageChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const selectedFile = target.files[0]
    if (validateImageFile(selectedFile)) {
      readImageFile(selectedFile)
    }
  }
}

function handleImageDrop(e: DragEvent) {
  e.preventDefault()
  imageDragOver.value = false
  const files = e.dataTransfer?.files
  if (files && files[0]) {
    const selectedFile = files[0]
    if (validateImageFile(selectedFile)) {
      readImageFile(selectedFile)
    }
  }
}

function readImageFile(selectedFile: File) {
  imageFile.value = selectedFile
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(selectedFile)
  error.value = ''
}

function removeImage() {
  imageFile.value = null
  imagePreview.value = ''
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function suggestNextVersion() {
  if (model.value?.current_version?.version_number) {
    const parts = model.value.current_version.version_number.split('.').map(p => parseInt(p, 10) || 0)
    parts[2] = (parts[2] || 0) + 1
    versionNumber.value = parts.join('.')
  } else {
    versionNumber.value = '1.0.0'
  }
}

async function handleSubmit() {
  error.value = ''

  if (!versionNumber.value) {
    error.value = '请输入版本号'
    return
  }

  if (versionNumberError.value) {
    error.value = versionNumberError.value
    return
  }

  if (!file.value) {
    error.value = '请选择模型文件'
    return
  }

  submitting.value = true
  uploadProgress.value = 0

  try {
    uploadProgress.value = 30
    const modelUpload = await uploadApi.uploadModel(file.value)
    const modelData = modelUpload.data.data!

    uploadProgress.value = 60
    let imageId: string | undefined
    if (imageFile.value) {
      const imageUpload = await uploadApi.uploadImage(imageFile.value)
      imageId = imageUpload.data.data?.file_id
    }

    uploadProgress.value = 80
    await modelVersionApi.create(modelId.value, {
      version_number: versionNumber.value,
      description: description.value || undefined,
      file_path: modelData.file_path,
      file_size: modelData.file_size,
      image_id: imageId,
      changelog: changelog.value || undefined,
    })

    uploadProgress.value = 100
    router.push(`/model/${modelId.value}`)
  } catch (err: any) {
    error.value = err.response?.data?.message || '上传失败，请重试'
  } finally {
    submitting.value = false
    uploadProgress.value = 0
  }
}

onMounted(() => {
  fetchModel()
  suggestNextVersion()
})
</script>

<template>
  <div class="mx-auto max-w-2xl px-4 py-6 sm:py-8">
    <Button variant="ghost" size="sm" class="focus-ring btn-press mb-4" @click="router.back()">
      <ArrowLeft class="mr-2 h-4 w-4" />
      返回
    </Button>

    <div v-if="loading" class="space-y-4">
      <Skeleton class="h-8 w-2/3" />
      <Skeleton class="h-64 w-full" />
    </div>

    <template v-else-if="model">
      <Card class="card-hover mb-6">
        <CardHeader>
          <CardTitle class="text-xl">{{ model.title }}</CardTitle>
          <p class="text-sm text-muted-foreground">
            当前版本: {{ model.current_version?.version_number || '无' }}
          </p>
        </CardHeader>
      </Card>

      <Card class="card-hover">
        <CardHeader>
          <CardTitle class="text-xl sm:text-2xl">上传新版本</CardTitle>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleSubmit" class="space-y-6">
            <Alert v-if="error" variant="destructive" class="animate-shake">
              <AlertDescription>{{ error }}</AlertDescription>
            </Alert>

            <div class="space-y-2">
              <Label for="version">版本号 *</Label>
              <div class="flex gap-2">
                <Input
                  id="version"
                  v-model="versionNumber"
                  placeholder="如: 1.0.1"
                  :error="!!versionNumberError"
                  class="h-11"
                />
                <Button type="button" variant="outline" class="focus-ring btn-press" @click="suggestNextVersion">
                  自动递增
                </Button>
              </div>
              <p v-if="versionNumberError" class="text-sm text-destructive">
                {{ versionNumberError }}
              </p>
              <p class="text-xs text-muted-foreground">
                版本号格式: x.y.z (如 1.0.1, 2.0.0)
              </p>
            </div>

            <div class="space-y-2">
              <Label for="description">版本描述</Label>
              <Textarea
                id="description"
                v-model="description"
                placeholder="简短描述这个版本的更新内容"
                :rows="2"
              />
            </div>

            <div class="space-y-2">
              <Label for="changelog">更新日志</Label>
              <Textarea
                id="changelog"
                v-model="changelog"
                placeholder="详细的更新内容，如：&#10;- 修复了xxx问题&#10;- 新增了xxx功能&#10;- 优化了xxx性能"
                :rows="4"
                class="min-h-[100px]"
              />
            </div>

            <div class="space-y-2">
              <Label>模型文件 * (.ysm, .zip)</Label>
              <div v-if="file" class="surface flex items-center justify-between p-3">
                <div class="flex items-center gap-3 min-w-0">
                  <FileArchive class="h-8 w-8 flex-shrink-0 text-primary" />
                  <div class="min-w-0">
                    <p class="truncate text-sm font-medium">{{ file.name }}</p>
                    <p class="text-xs text-muted-foreground">{{ formatFileSize(file.size) }}</p>
                  </div>
                </div>
                <Button type="button" variant="ghost" size="icon" class="focus-ring btn-press h-8 w-8" aria-label="移除文件" @click="removeFile">
                  <X class="h-4 w-4" />
                </Button>
              </div>
              <label
                v-else
                class="flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-6 text-center transition-colors"
                :class="fileDragOver ? 'border-primary bg-primary/5' : 'border-muted-foreground/25 hover:border-primary hover:bg-primary/5'"
                @dragover.prevent="fileDragOver = true"
                @dragleave.prevent="fileDragOver = false"
                @drop="handleFileDrop"
              >
                <Upload class="h-8 w-8 text-muted-foreground" />
                <p class="mt-2 text-sm text-muted-foreground">点击或拖拽上传模型文件</p>
                <p class="text-xs text-muted-foreground">支持 .ysm, .zip 格式，最大 100MB</p>
                <Input id="file" type="file" accept=".ysm,.zip" class="hidden" @change="handleFileChange" />
              </label>
            </div>

            <div class="space-y-2">
              <Label>预览图（可选）</Label>
              <p class="text-xs text-muted-foreground">支持 jpg, png, gif 格式，最大 5MB</p>
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
                :class="imageDragOver ? 'border-primary bg-primary/5' : 'border-muted-foreground/25 hover:border-primary hover:bg-primary/5'"
                @dragover.prevent="imageDragOver = true"
                @dragleave.prevent="imageDragOver = false"
                @drop="handleImageDrop"
              >
                <ImageOff class="h-8 w-8 text-muted-foreground" />
                <p class="mt-2 text-sm text-muted-foreground">点击或拖拽上传预览图</p>
                <p class="text-xs text-muted-foreground">支持 JPG、PNG、GIF、WEBP</p>
                <Input type="file" accept="image/*" class="hidden" @change="handleImageChange" />
              </label>
            </div>

            <Progress
              v-if="submitting && uploadProgress > 0"
              :model-value="uploadProgress"
              :max="100"
              show-value
              size="lg"
            >
              <span>上传进度</span>
            </Progress>

            <Button type="submit" class="focus-ring btn-press w-full sm:w-auto" :loading="submitting">
              <Upload v-if="!submitting" class="mr-2 h-4 w-4" />
              {{ submitting ? '上传中...' : '上传新版本' }}
            </Button>
          </form>
        </CardContent>
      </Card>
    </template>
  </div>
</template>
