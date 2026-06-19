<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { userApi } from '@/lib/api'
import type { User as UserType, Model, PaginatedResponse } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import ModelCard from '@/components/ModelCard.vue'
import { User, Calendar, Package, ArrowLeft } from 'lucide-vue-next'
import { getAvatarUrl } from '@/utils/image'

const route = useRoute()

const user = ref<UserType | null>(null)
const models = ref<Model[]>([])
const loading = ref(true)
const loadingModels = ref(true)

async function fetchUser() {
  loading.value = true
  try {
    const response = await userApi.getById(route.params.id as string)
    user.value = response.data.data!
    await fetchUserModels()
  } catch (error) {
    console.error('Failed to fetch user:', error)
  } finally {
    loading.value = false
  }
}

async function fetchUserModels() {
  if (!user.value) return
  loadingModels.value = true
  try {
    const response = await userApi.getUserModels(user.value.id, 1, 20)
    const data = response.data.data as PaginatedResponse<Model>
    models.value = data.items.filter(m => m.is_public && m.status === 'approved')
  } catch (error) {
    console.error('Failed to fetch user models:', error)
  } finally {
    loadingModels.value = false
  }
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(fetchUser)
</script>

<template>
  <div class="mx-auto max-w-4xl px-4 py-6 sm:py-8">
    <RouterLink to="/">
      <Button variant="ghost" size="sm" class="focus-ring btn-press mb-4">
        <ArrowLeft class="mr-2 h-4 w-4" />
        返回
      </Button>
    </RouterLink>

    <div v-if="loading" class="space-y-6">
      <div class="surface flex flex-col items-center gap-4 p-6 sm:flex-row sm:items-start">
        <Skeleton variant="shimmer" class="h-20 w-20 rounded-full" />
        <div class="w-full space-y-3 text-center sm:text-left">
          <Skeleton variant="shimmer" class="mx-auto h-7 w-40 sm:mx-0" />
          <Skeleton variant="shimmer" class="mx-auto h-4 w-2/3 sm:mx-0" />
          <Skeleton variant="shimmer" class="mx-auto h-4 w-32 sm:mx-0" />
        </div>
      </div>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <Skeleton v-for="i in 4" :key="i" variant="shimmer" class="h-48" />
      </div>
    </div>

    <template v-else-if="user">
      <Card class="card-hover mb-6">
        <CardContent class="flex flex-col items-center gap-4 p-6 sm:flex-row sm:items-start">
          <Avatar class="h-20 w-20 sm:h-24 sm:w-24">
            <AvatarImage :src="getAvatarUrl(user.avatar_id, user.avatar_url) || undefined">
              <User class="h-8 w-8 text-muted-foreground" />
            </AvatarImage>
          </Avatar>
          <div class="text-center sm:text-left">
            <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">{{ user.username }}</h1>
            <p v-if="user.bio" class="mt-1 max-w-lg text-sm text-muted-foreground sm:text-base">
              {{ user.bio }}
            </p>
            <div class="mt-3 flex flex-wrap items-center justify-center gap-3 text-sm text-muted-foreground sm:justify-start">
              <span class="flex items-center gap-1">
                <Calendar class="h-4 w-4" />
                {{ formatDate(user.created_at) }} 加入
              </span>
              <Badge v-if="user.role === 'super_admin'" variant="default">站长</Badge>
              <Badge v-else-if="user.role === 'admin'" variant="secondary">管理员</Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card class="card-hover">
        <CardContent class="p-6">
          <h2 class="mb-4 flex items-center gap-2 text-lg font-semibold sm:text-xl">
            <Package class="h-5 w-5" />
            发布的模型
          </h2>

          <div v-if="loadingModels" class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Skeleton v-for="i in 4" :key="i" variant="shimmer" class="h-48" />
          </div>

          <div
            v-else-if="models.length === 0"
            class="surface flex flex-col items-center justify-center py-16 text-center animate-fade-in"
          >
            <Package class="h-12 w-12 text-muted-foreground opacity-50" />
            <p class="mt-4 text-muted-foreground">暂无公开模型</p>
            <RouterLink to="/">
              <Button variant="link" class="btn-press mt-2 h-auto p-0">去发现模型</Button>
            </RouterLink>
          </div>

          <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <ModelCard v-for="model in models" :key="model.id" :model="model" />
          </div>
        </CardContent>
      </Card>
    </template>

    <Card v-else class="card-hover">
      <CardContent class="flex flex-col items-center justify-center py-16 text-center text-muted-foreground">
        <User class="h-12 w-12 opacity-50" />
        <p class="mt-4">用户不存在或已被移除</p>
        <RouterLink to="/">
          <Button class="btn-press focus-ring mt-4">返回首页</Button>
        </RouterLink>
      </CardContent>
    </Card>
  </div>
</template>
