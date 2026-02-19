<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { modelApi, uploadApi } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Model } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Save, X, Loader2, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const model = ref<Model | null>(null)
const title = ref('')
const description = ref('')
const tags = ref('')
const isPublic = ref(true)
const imageFile = ref<File | null>(null)
const imagePreview = ref('')
const loading = ref(false)
const fetching = ref(true)

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
    imagePreview.value = model.value.image_url || ''
  } catch (error) {
    console.error('Failed to fetch model:', error)
    router.push('/')
  } finally {
    fetching.value = false
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
  imagePreview.value = model.value?.image_url || ''
}

async function handleSubmit() {
  if (!title.value) return

  loading.value = true

  try {
    let imageUrl: string | undefined = model.value?.image_url || undefined
    if (imageFile.value) {
      const imageUpload = await uploadApi.uploadImage(imageFile.value)
      imageUrl = imageUpload.data.data!.url
    }

    await modelApi.update(route.params.id as string, {
      title: title.value,
      description: description.value,
      tags: tags.value.split(',').map(t => t.trim()).filter(Boolean),
      is_public: isPublic.value,
      image_url: imageUrl,
    })

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
