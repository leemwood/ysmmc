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

const title = ref('')
const description = ref('')
const tags = ref('')
const isPublic = ref(true)
const file = ref<File | null>(null)
const imageFile = ref<File | null>(null)
const imagePreview = ref('')
const loading = ref(false)
const error = ref('')

function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    file.value = target.files[0]
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
  imagePreview.value = ''
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

  try {
    const modelUpload = await uploadApi.uploadModel(file.value)
    const modelData = modelUpload.data.data!

    let imageUrl = ''
    if (imageFile.value) {
      const imageUpload = await uploadApi.uploadImage(imageFile.value)
      imageUrl = imageUpload.data.data!.url
    }

    await modelApi.create({
      title: title.value,
      description: description.value,
      file_path: modelData.file_path,
      file_size: modelData.file_size,
      image_url: imageUrl || undefined,
      tags: tags.value.split(',').map(t => t.trim()).filter(Boolean),
      is_public: isPublic.value,
    })

    router.push('/profile')
  } catch (err: any) {
    error.value = err.response?.data?.message || '上传失败，请重试'
  } finally {
    loading.value = false
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
            <p v-if="file" class="text-sm text-muted-foreground">
              已选择: {{ file.name }}
            </p>
          </div>

          <div class="space-y-2">
            <Label>预览图</Label>
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
