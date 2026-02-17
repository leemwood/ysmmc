<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { userApi } from '@/lib/api'
import type { User as UserType } from '@/types'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { User, Calendar, Package } from 'lucide-vue-next'

const route = useRoute()

const user = ref<UserType | null>(null)
const loading = ref(true)

async function fetchUser() {
  try {
    const response = await userApi.getById(route.params.id as string)
    user.value = response.data.data!
  } catch (error) {
    console.error('Failed to fetch user:', error)
  } finally {
    loading.value = false
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
            <AvatarImage :src="user.avatar_url || undefined">
              <User class="h-8 w-8" />
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
              <Badge v-if="user.role === 'admin'" variant="default">管理员</Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <Package class="h-5 w-5" />
            发布的模型
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-center py-8 text-muted-foreground">
            暂无公开模型
          </p>
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
