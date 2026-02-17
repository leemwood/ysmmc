<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/lib/api'
import type { Model, User, PaginatedResponse } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { LayoutDashboard, Package, Users, Download, Check, X, Loader2, FileText } from 'lucide-vue-next'

const activeTab = ref('overview')
const loading = ref(true)
const stats = ref({
  total_users: 0,
  total_models: 0,
  pending_models: 0,
  total_downloads: 0,
})
const pendingModels = ref<Model[]>([])
const pendingProfiles = ref<User[]>([])

async function fetchStats() {
  try {
    const response = await adminApi.getStats()
    stats.value = response.data.data!
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

async function fetchPendingModels() {
  try {
    const response = await adminApi.listPendingModels(1, 20)
    pendingModels.value = (response.data.data as PaginatedResponse<Model>).items
  } catch (error) {
    console.error('Failed to fetch pending models:', error)
  }
}

async function fetchPendingProfiles() {
  try {
    const response = await adminApi.listPendingProfiles(1, 20)
    pendingProfiles.value = (response.data.data as PaginatedResponse<User>).items
  } catch (error) {
    console.error('Failed to fetch pending profiles:', error)
  }
}

async function approveModel(id: string) {
  try {
    await adminApi.approveModel(id)
    fetchPendingModels()
    fetchStats()
  } catch (error) {
    console.error('Failed to approve model:', error)
  }
}

async function rejectModel(id: string) {
  const reason = prompt('请输入拒绝原因:')
  if (!reason) return
  
  try {
    await adminApi.rejectModel(id, reason)
    fetchPendingModels()
    fetchStats()
  } catch (error) {
    console.error('Failed to reject model:', error)
  }
}

async function approveProfile(id: string) {
  try {
    await adminApi.approveProfile(id)
    fetchPendingProfiles()
  } catch (error) {
    console.error('Failed to approve profile:', error)
  }
}

async function rejectProfile(id: string) {
  try {
    await adminApi.rejectProfile(id)
    fetchPendingProfiles()
  } catch (error) {
    console.error('Failed to reject profile:', error)
  }
}

onMounted(async () => {
  await Promise.all([fetchStats(), fetchPendingModels(), fetchPendingProfiles()])
  loading.value = false
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 py-8">
    <h1 class="mb-6 text-2xl font-bold">管理后台</h1>

    <div class="mb-6 flex gap-4 border-b">
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'overview' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'overview'"
      >
        <LayoutDashboard class="h-4 w-4" />
        概览
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'models' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'models'"
      >
        <Package class="h-4 w-4" />
        模型审核
        <Badge v-if="stats.pending_models > 0" variant="destructive" class="ml-1">{{ stats.pending_models }}</Badge>
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'profiles' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'profiles'"
      >
        <Users class="h-4 w-4" />
        资料审核
        <Badge v-if="pendingProfiles.length > 0" variant="destructive" class="ml-1">{{ pendingProfiles.length }}</Badge>
      </button>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <div v-if="activeTab === 'overview'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent class="flex items-center gap-4 p-6">
            <div class="rounded-full bg-primary/10 p-3">
              <Users class="h-6 w-6 text-primary" />
            </div>
            <div>
              <p class="text-sm text-muted-foreground">总用户数</p>
              <p class="text-2xl font-bold">{{ stats.total_users }}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="flex items-center gap-4 p-6">
            <div class="rounded-full bg-primary/10 p-3">
              <Package class="h-6 w-6 text-primary" />
            </div>
            <div>
              <p class="text-sm text-muted-foreground">总模型数</p>
              <p class="text-2xl font-bold">{{ stats.total_models }}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="flex items-center gap-4 p-6">
            <div class="rounded-full bg-yellow-500/10 p-3">
              <FileText class="h-6 w-6 text-yellow-500" />
            </div>
            <div>
              <p class="text-sm text-muted-foreground">待审核</p>
              <p class="text-2xl font-bold">{{ stats.pending_models }}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="flex items-center gap-4 p-6">
            <div class="rounded-full bg-green-500/10 p-3">
              <Download class="h-6 w-6 text-green-500" />
            </div>
            <div>
              <p class="text-sm text-muted-foreground">总下载量</p>
              <p class="text-2xl font-bold">{{ stats.total_downloads }}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-else-if="activeTab === 'models'">
        <div v-if="pendingModels.length === 0" class="text-center py-12 text-muted-foreground">
          暂无待审核模型
        </div>

        <div v-else class="space-y-4">
          <Card v-for="model in pendingModels" :key="model.id">
            <CardContent class="flex items-center justify-between p-4">
              <div class="flex items-center gap-4">
                <div class="h-16 w-16 flex-shrink-0 rounded bg-muted">
                  <img v-if="model.image_url" :src="model.image_url" class="h-full w-full object-cover rounded" />
                </div>
                <div>
                  <RouterLink :to="`/model/${model.id}`" class="font-medium hover:underline">
                    {{ model.title }}
                  </RouterLink>
                  <p class="text-sm text-muted-foreground line-clamp-1">
                    {{ model.description || '暂无描述' }}
                  </p>
                  <p class="text-xs text-muted-foreground">
                    上传者: {{ model.user?.username }}
                  </p>
                </div>
              </div>
              <div class="flex gap-2">
                <Button size="sm" variant="default" @click="approveModel(model.id)">
                  <Check class="mr-1 h-4 w-4" />
                  通过
                </Button>
                <Button size="sm" variant="destructive" @click="rejectModel(model.id)">
                  <X class="mr-1 h-4 w-4" />
                  拒绝
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <div v-else-if="activeTab === 'profiles'">
        <div v-if="pendingProfiles.length === 0" class="text-center py-12 text-muted-foreground">
          暂无待审核资料
        </div>

        <div v-else class="space-y-4">
          <Card v-for="user in pendingProfiles" :key="user.id">
            <CardContent class="flex items-center justify-between p-4">
              <div class="flex items-center gap-4">
                <Avatar class="h-12 w-12">
                  <AvatarImage :src="user.avatar_url || undefined">
                    <Users class="h-6 w-6" />
                  </AvatarImage>
                </Avatar>
                <div>
                  <p class="font-medium">{{ user.username }}</p>
                  <p class="text-sm text-muted-foreground">
                    新用户名: {{ user.pending_changes?.username || '无变更' }}
                  </p>
                  <p class="text-sm text-muted-foreground">
                    新简介: {{ user.pending_changes?.bio || '无变更' }}
                  </p>
                </div>
              </div>
              <div class="flex gap-2">
                <Button size="sm" variant="default" @click="approveProfile(user.id)">
                  <Check class="mr-1 h-4 w-4" />
                  通过
                </Button>
                <Button size="sm" variant="destructive" @click="rejectProfile(user.id)">
                  <X class="mr-1 h-4 w-4" />
                  拒绝
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </template>
  </div>
</template>
