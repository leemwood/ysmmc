<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { userApi } from '@/lib/api'
import type { User as UserType, Model, PaginatedResponse } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { RouterLink } from 'vue-router'
import { User, Calendar, Package, Download } from 'lucide-vue-next'

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
  <div class="mx-auto max-w-2xl px-4 py-8">
    <div v-if="loading" class="space-y-4">
      <Skeleton class="h-32 w-full" />
      <Skeleton class="h-64 w-full" />
    </div>

    <template v-else-if="user">
      <Card class="mb-6">
        <CardContent class="flex items-center gap-4 p-6">
          <Avatar class="h-20 w-20">
            <AvatarImage v-if="user.avatar_url" :src="user.avatar_url">
              <User class="h-8 w-8 text-muted-foreground" />
            </AvatarImage>
          </Avatar>
          <div>
            <h1 class="text-2xl font-bold">{{ user.username }}</h1>
            <p v-if="user.bio" class="mt-1 text-muted-foreground">{{ user.bio }}</p>
            <div class="mt-2 flex items-center gap-4 text-sm text-muted-foreground">
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

      <Card>
        <CardContent class="p-6">
          <h2 class="mb-4 flex items-center gap-2 text-lg font-semibold">
            <Package class="h-5 w-5" />
            发布的模型
          </h2>

          <div v-if="loadingModels" class="grid gap-4 sm:grid-cols-2">
            <Skeleton v-for="i in 4" :key="i" class="h-24" />
          </div>

          <div v-else-if="models.length === 0" class="text-center py-8 text-muted-foreground">
            暂无公开模型
          </div>

          <div v-else class="grid gap-4 sm:grid-cols-2">
            <RouterLink
              v-for="model in models"
              :key="model.id"
              :to="`/model/${model.id}`"
              class="block rounded-lg border bg-card p-4 hover:bg-accent transition-colors"
            >
              <div class="flex items-start gap-3">
                <div class="h-14 w-14 flex-shrink-0 rounded bg-muted overflow-hidden">
                  <img v-if="model.image_url" :src="model.image_url" class="h-full w-full object-cover" />
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="font-medium line-clamp-1">{{ model.title }}</h3>
                  <p class="mt-1 text-xs text-muted-foreground line-clamp-2">{{ model.description || '暂无描述' }}</p>
                  <div class="mt-2 flex items-center gap-2 text-xs text-muted-foreground">
                    <Download class="h-3 w-3" />
                    {{ model.downloads }} 下载
                  </div>
                </div>
              </div>
            </RouterLink>
          </div>
        </CardContent>
      </Card>
    </template>

    <Card v-else>
      <CardContent class="py-12 text-center text-muted-foreground">
        用户不存在
      </CardContent>
    </Card>
  </div>
</template>
