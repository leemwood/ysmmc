<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { modelApi, uploadApi } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Upload, X, Loader2 } from 'lucide-vue-next'

const router = useRouter()

const MAX_FILE_SIZE = 100 * 1024 * 1024
const MAX_IMAGE_SIZE = 5 * 1024 * 1024

const title = ref('')
const description = ref('')
const tags = ref('')
const isPublic = ref(true)
const file = ref<File | null>(null)
const imageFile = ref<File | null>(null)
const imagePreview = ref('')
const loading = ref(false)
const error = ref('')
const uploadProgress = ref(0)

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

function handleImageChange(e: Event) {
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
    
    imageFile.value = selectedFile
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(imageFile.value)
    error.value = ''
  }
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

  loading.value = true
  uploadProgress.value = 0

  try {
    uploadProgress.value = 30
    const modelUpload = await uploadApi.uploadModel(file.value)
    const modelData = modelUpload.data.data!

    uploadProgress.value = 60
    let imageUrl = ''
    if (imageFile.value) {
      const imageUpload = await uploadApi.uploadImage(imageFile.value)
      imageUrl = imageUpload.data.data!.url
    }

    uploadProgress.value = 80
    await modelApi.create({
      title: title.value,
      description: description.value,
      file_path: modelData.file_path,
      file_size: modelData.file_size,
      image_url: imageUrl || undefined,
      tags: tags.value.split(',').map(t => t.trim()).filter(Boolean),
      is_public: isPublic.value,
    })

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
  <div class="mx-auto max-w-2xl px-4 py-8">
    <Card>
      <CardHeader>
        <CardTitle class="text-2xl">上传模型</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div v-if="error" class="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
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
            />
            <p class="text-sm text-muted-foreground">
              支持 .ysm, .zip 格式，最大 100MB
            </p>
            <p v-if="file" class="text-sm text-primary">
              已选择: {{ file.name }} ({{ formatFileSize(file.size) }})
            </p>
          </div>

          <div class="space-y-2">
            <Label>预览图</Label>
            <p class="text-sm text-muted-foreground">支持 jpg, png, gif 格式，最大 5MB</p>
            <div v-if="imagePreview" class="relative inline-block">
              <img :src="imagePreview" class="h-40 w-auto rounded-md object-cover" />
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

          <div class="flex items-center gap-2">
            <input
              id="isPublic"
              v-model="isPublic"
              type="checkbox"
              class="h-4 w-4 rounded border-gray-300"
            />
            <Label for="isPublic">公开显示</Label>
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

          <Button type="submit" class="w-full" :disabled="loading">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            <Upload v-else class="mr-2 h-4 w-4" />
            {{ loading ? '上传中...' : '上传模型' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
