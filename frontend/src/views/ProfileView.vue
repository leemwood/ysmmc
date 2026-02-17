<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi, modelApi, favoriteApi } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Model, Favorite, PaginatedResponse } from '@/types'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { RouterLink } from 'vue-router'
import { User, Save, Loader2, Heart, Package } from 'lucide-vue-next'

const authStore = useAuthStore()

const activeTab = ref('profile')
const username = ref('')
const bio = ref('')
const loading = ref(false)
const models = ref<Model[]>([])
const favorites = ref<Favorite[]>([])
const loadingModels = ref(true)
const loadingFavorites = ref(true)
const message = ref('')

async function loadProfile() {
  if (authStore.user) {
    username.value = authStore.user.username
    bio.value = authStore.user.bio || ''
  }
}

async function loadModels() {
  loadingModels.value = true
  try {
    const response = await modelApi.list(1, 100)
    const allModels = (response.data.data as PaginatedResponse<Model>).items
    models.value = allModels.filter(m => m.user_id === authStore.user?.id)
  } catch (error) {
    console.error('Failed to load models:', error)
  } finally {
    loadingModels.value = false
  }
}

async function loadFavorites() {
  loadingFavorites.value = true
  try {
    const response = await favoriteApi.list(1, 100)
    favorites.value = (response.data.data as PaginatedResponse<Favorite>).items
  } catch (error) {
    console.error('Failed to load favorites:', error)
  } finally {
    loadingFavorites.value = false
  }
}

async function handleUpdateProfile() {
  loading.value = true
  message.value = ''

  try {
    await userApi.updateMe({
      username: username.value !== authStore.user?.username ? username.value : undefined,
      bio: bio.value,
    })
    await authStore.fetchUser()
    message.value = '资料更新成功'
    if (!authStore.isAdmin) {
      message.value = '资料更新已提交审核'
    }
  } catch (error: any) {
    message.value = error.response?.data?.message || '更新失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadProfile()
  loadModels()
  loadFavorites()
})
</script>

<template>
  <div class="mx-auto max-w-4xl px-4 py-8">
    <div class="mb-6 flex gap-4 border-b">
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'profile' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'profile'"
      >
        <User class="h-4 w-4" />
        个人资料
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'models' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'models'"
      >
        <Package class="h-4 w-4" />
        我的模型
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'favorites' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'favorites'"
      >
        <Heart class="h-4 w-4" />
        我的收藏
      </button>
    </div>

    <div v-if="activeTab === 'profile'">
      <Card>
        <CardHeader>
          <CardTitle>个人资料</CardTitle>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleUpdateProfile" class="space-y-4">
            <div v-if="message" class="rounded-md bg-primary/10 p-3 text-sm text-primary">
              {{ message }}
            </div>

            <div class="flex items-center gap-4">
              <Avatar class="h-20 w-20">
                <AvatarImage :src="authStore.user?.avatar_url || undefined">
                  <User class="h-8 w-8" />
                </AvatarImage>
              </Avatar>
              <div>
                <p class="font-medium">{{ authStore.user?.username }}</p>
                <p class="text-sm text-muted-foreground">{{ authStore.user?.email }}</p>
                <Badge v-if="authStore.isAdmin" variant="default" class="mt-1">管理员</Badge>
              </div>
            </div>

            <div class="space-y-2">
              <Label for="username">用户名</Label>
              <Input id="username" v-model="username" placeholder="请输入用户名" />
            </div>

            <div class="space-y-2">
              <Label for="bio">个人简介</Label>
              <Textarea id="bio" v-model="bio" placeholder="介绍一下自己..." :rows="3" />
            </div>

            <div v-if="authStore.user?.profile_status === 'pending_review'" class="rounded-md bg-yellow-500/10 p-3 text-sm text-yellow-600">
              资料修改正在审核中
            </div>

            <Button type="submit" :disabled="loading">
              <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
              <Save v-else class="mr-2 h-4 w-4" />
              保存修改
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>

    <div v-else-if="activeTab === 'models'">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-xl font-semibold">我的模型</h2>
        <RouterLink to="/upload">
          <Button size="sm">上传模型</Button>
        </RouterLink>
      </div>

      <div v-if="loadingModels" class="grid gap-4 sm:grid-cols-2">
        <Skeleton v-for="i in 4" :key="i" class="h-32" />
      </div>

      <div v-else-if="models.length === 0" class="text-center py-12 text-muted-foreground">
        暂无模型，<RouterLink to="/upload" class="text-primary hover:underline">去上传</RouterLink>
      </div>

      <div v-else class="grid gap-4 sm:grid-cols-2">
        <Card v-for="model in models" :key="model.id">
          <CardContent class="p-4">
            <div class="flex items-start gap-4">
              <div class="h-16 w-16 flex-shrink-0 rounded bg-muted">
                <img v-if="model.image_url" :src="model.image_url" class="h-full w-full object-cover rounded" />
              </div>
              <div class="flex-1 min-w-0">
                <RouterLink :to="`/model/${model.id}`" class="font-medium hover:underline line-clamp-1">
                  {{ model.title }}
                </RouterLink>
                <div class="mt-1 flex items-center gap-2">
                  <Badge :variant="model.status === 'approved' ? 'success' : model.status === 'rejected' ? 'destructive' : 'secondary'">
                    {{ model.status === 'approved' ? '已通过' : model.status === 'rejected' ? '已拒绝' : '审核中' }}
                  </Badge>
                  <span class="text-xs text-muted-foreground">{{ model.downloads }} 下载</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>

    <div v-else-if="activeTab === 'favorites'">
      <h2 class="mb-4 text-xl font-semibold">我的收藏</h2>

      <div v-if="loadingFavorites" class="grid gap-4 sm:grid-cols-2">
        <Skeleton v-for="i in 4" :key="i" class="h-32" />
      </div>

      <div v-else-if="favorites.length === 0" class="text-center py-12 text-muted-foreground">
        暂无收藏
      </div>

      <div v-else class="grid gap-4 sm:grid-cols-2">
        <Card v-for="favorite in favorites" :key="favorite.id">
          <CardContent class="p-4">
            <div class="flex items-start gap-4">
              <div class="h-16 w-16 flex-shrink-0 rounded bg-muted">
                <img v-if="favorite.model?.image_url" :src="favorite.model.image_url" class="h-full w-full object-cover rounded" />
              </div>
              <div class="flex-1 min-w-0">
                <RouterLink :to="`/model/${favorite.model_id}`" class="font-medium hover:underline line-clamp-1">
                  {{ favorite.model?.title }}
                </RouterLink>
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ favorite.model?.downloads || 0 }} 下载
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>
