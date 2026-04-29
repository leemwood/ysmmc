<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/lib/api'
import type { Model, User, PaginatedResponse } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { LayoutDashboard, Package, Users, Download, Check, X, Loader2, FileText, Megaphone } from 'lucide-vue-next'
import { getAvatarUrl, getModelImageUrl } from '@/utils/image'

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

const rejectDialog = ref(false)
const rejectModelId = ref('')
const rejectReason = ref('')
const rejecting = ref(false)

const successMessage = ref('')
const errorMessage = ref('')

function showSuccess(msg: string) {
  successMessage.value = msg
  setTimeout(() => { successMessage.value = '' }, 3000)
}

function showError(msg: string) {
  errorMessage.value = msg
  setTimeout(() => { errorMessage.value = '' }, 3000)
}

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
    showSuccess('模型已通过审核')
    fetchPendingModels()
    fetchStats()
  } catch (error: any) {
    showError(error.response?.data?.message || '操作失败')
  }
}

function openRejectDialog(id: string) {
  rejectModelId.value = id
  rejectReason.value = ''
  rejectDialog.value = true
}

async function rejectModel() {
  if (!rejectReason.value.trim()) {
    showError('请输入拒绝原因')
    return
  }
  
  rejecting.value = true
  try {
    await adminApi.rejectModel(rejectModelId.value, rejectReason.value)
    showSuccess('模型已拒绝')
    rejectDialog.value = false
    fetchPendingModels()
    fetchStats()
  } catch (error: any) {
    showError(error.response?.data?.message || '操作失败')
  } finally {
    rejecting.value = false
  }
}

async function approveProfile(id: string) {
  try {
    await adminApi.approveProfile(id)
    showSuccess('资料已通过审核')
    fetchPendingProfiles()
  } catch (error: any) {
    showError(error.response?.data?.message || '操作失败')
  }
}

async function rejectProfile(id: string) {
  try {
    await adminApi.rejectProfile(id)
    showSuccess('资料变更已拒绝')
    fetchPendingProfiles()
  } catch (error: any) {
    showError(error.response?.data?.message || '操作失败')
  }
}

onMounted(async () => {
  await Promise.all([fetchStats(), fetchPendingModels(), fetchPendingProfiles()])
  loading.value = false
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 py-6 sm:py-8">
    <h1 class="mb-6 text-xl sm:text-2xl font-bold">管理后台</h1>

    <div v-if="successMessage" class="mb-4 rounded-md bg-green-500/10 p-3 text-sm text-green-600 animate-fade-in">
      {{ successMessage }}
    </div>
    <div v-if="errorMessage" class="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive animate-fade-in">
      {{ errorMessage }}
    </div>

    <div class="mb-6 -mx-4 px-4 overflow-x-auto">
      <div class="flex gap-1 sm:gap-4 border-b min-w-max">
        <button
          class="flex items-center gap-2 px-3 sm:px-4 py-2 text-sm font-medium border-b-2 transition-colors whitespace-nowrap"
          :class="activeTab === 'overview' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
          @click="activeTab = 'overview'"
        >
          <LayoutDashboard class="h-4 w-4" />
          <span class="hidden sm:inline">概览</span>
        </button>
        <button
          class="flex items-center gap-2 px-3 sm:px-4 py-2 text-sm font-medium border-b-2 transition-colors whitespace-nowrap"
          :class="activeTab === 'models' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
          @click="activeTab = 'models'"
        >
          <Package class="h-4 w-4" />
          <span class="hidden sm:inline">模型审核</span>
          <Badge v-if="stats.pending_models > 0" variant="destructive" class="ml-1 text-xs">{{ stats.pending_models }}</Badge>
        </button>
        <button
          class="flex items-center gap-2 px-3 sm:px-4 py-2 text-sm font-medium border-b-2 transition-colors whitespace-nowrap"
          :class="activeTab === 'profiles' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
          @click="activeTab = 'profiles'"
        >
          <Users class="h-4 w-4" />
          <span class="hidden sm:inline">资料审核</span>
          <Badge v-if="pendingProfiles.length > 0" variant="destructive" class="ml-1 text-xs">{{ pendingProfiles.length }}</Badge>
        </button>
        <RouterLink
          to="/admin/users"
          class="flex items-center gap-2 px-3 sm:px-4 py-2 text-sm font-medium border-b-2 border-transparent text-muted-foreground hover:text-foreground transition-colors whitespace-nowrap"
        >
          <Users class="h-4 w-4" />
          <span class="hidden sm:inline">用户管理</span>
        </RouterLink>
        <RouterLink
          to="/admin/models"
          class="flex items-center gap-2 px-3 sm:px-4 py-2 text-sm font-medium border-b-2 border-transparent text-muted-foreground hover:text-foreground transition-colors whitespace-nowrap"
        >
          <Package class="h-4 w-4" />
          <span class="hidden sm:inline">模型管理</span>
        </RouterLink>
        <RouterLink
          to="/admin/announcements"
          class="flex items-center gap-2 px-3 sm:px-4 py-2 text-sm font-medium border-b-2 border-transparent text-muted-foreground hover:text-foreground transition-colors whitespace-nowrap"
        >
          <Megaphone class="h-4 w-4" />
          <span class="hidden sm:inline">公告管理</span>
        </RouterLink>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <div v-if="activeTab === 'overview'" class="grid gap-3 sm:gap-4 grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent class="flex items-center gap-3 sm:gap-4 p-4 sm:p-6">
            <div class="rounded-full bg-primary/10 p-2 sm:p-3">
              <Users class="h-5 w-5 sm:h-6 sm:w-6 text-primary" />
            </div>
            <div>
              <p class="text-xs sm:text-sm text-muted-foreground">总用户数</p>
              <p class="text-xl sm:text-2xl font-bold">{{ stats.total_users }}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="flex items-center gap-3 sm:gap-4 p-4 sm:p-6">
            <div class="rounded-full bg-primary/10 p-2 sm:p-3">
              <Package class="h-5 w-5 sm:h-6 sm:w-6 text-primary" />
            </div>
            <div>
              <p class="text-xs sm:text-sm text-muted-foreground">总模型数</p>
              <p class="text-xl sm:text-2xl font-bold">{{ stats.total_models }}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="flex items-center gap-3 sm:gap-4 p-4 sm:p-6">
            <div class="rounded-full bg-yellow-500/10 p-2 sm:p-3">
              <FileText class="h-5 w-5 sm:h-6 sm:w-6 text-yellow-500" />
            </div>
            <div>
              <p class="text-xs sm:text-sm text-muted-foreground">待审核</p>
              <p class="text-xl sm:text-2xl font-bold">{{ stats.pending_models }}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="flex items-center gap-3 sm:gap-4 p-4 sm:p-6">
            <div class="rounded-full bg-green-500/10 p-2 sm:p-3">
              <Download class="h-5 w-5 sm:h-6 sm:w-6 text-green-500" />
            </div>
            <div>
              <p class="text-xs sm:text-sm text-muted-foreground">总下载量</p>
              <p class="text-xl sm:text-2xl font-bold">{{ stats.total_downloads }}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-else-if="activeTab === 'models'">
        <div v-if="pendingModels.length === 0" class="text-center py-12 text-muted-foreground">
          暂无待审核模型
        </div>

        <div v-else class="space-y-3 sm:space-y-4">
          <Card v-for="model in pendingModels" :key="model.id" class="overflow-hidden">
            <CardContent class="p-4">
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <div class="flex items-center gap-3 sm:gap-4">
                  <div class="h-14 w-14 sm:h-16 sm:w-16 flex-shrink-0 rounded bg-muted overflow-hidden">
                    <img v-if="model.image_id || model.image_url" :src="getModelImageUrl(model.image_id, model.image_url)" class="h-full w-full object-cover rounded" />
                  </div>
                  <div class="min-w-0 flex-1">
                    <RouterLink :to="`/model/${model.id}`" class="font-medium hover:underline line-clamp-1">
                      {{ model.title }}
                    </RouterLink>
                    <p class="text-sm text-muted-foreground line-clamp-1 mt-0.5">
                      {{ model.description || '暂无描述' }}
                    </p>
                    <p class="text-xs text-muted-foreground mt-0.5">
                      上传者: {{ model.user?.username }}
                    </p>
                  </div>
                </div>
                <div class="flex gap-2 justify-end sm:flex-shrink-0">
                  <Button size="sm" variant="default" class="btn-press" @click="approveModel(model.id)">
                    <Check class="mr-1 h-4 w-4" />
                    通过
                  </Button>
                  <Button size="sm" variant="destructive" class="btn-press" @click="openRejectDialog(model.id)">
                    <X class="mr-1 h-4 w-4" />
                    拒绝
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <div v-else-if="activeTab === 'profiles'">
        <div v-if="pendingProfiles.length === 0" class="text-center py-12 text-muted-foreground">
          暂无待审核资料
        </div>

        <div v-else class="space-y-3 sm:space-y-4">
          <Card v-for="user in pendingProfiles" :key="user.id" class="overflow-hidden">
            <CardContent class="p-4">
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <div class="flex items-center gap-3 sm:gap-4">
                  <Avatar class="h-10 w-10 sm:h-12 sm:w-12">
                    <AvatarImage :src="getAvatarUrl(user.avatar_id, user.avatar_url) || undefined">
                      <Users class="h-5 w-5 sm:h-6 sm:w-6" />
                    </AvatarImage>
                  </Avatar>
                  <div class="min-w-0 flex-1">
                    <p class="font-medium">{{ user.username }}</p>
                    <p class="text-sm text-muted-foreground line-clamp-1">
                      新用户名: {{ user.pending_changes?.username || '无变更' }}
                    </p>
                    <p class="text-sm text-muted-foreground line-clamp-1">
                      新简介: {{ user.pending_changes?.bio || '无变更' }}
                    </p>
                  </div>
                </div>
                <div class="flex gap-2 justify-end sm:flex-shrink-0">
                  <Button size="sm" variant="default" class="btn-press" @click="approveProfile(user.id)">
                    <Check class="mr-1 h-4 w-4" />
                    通过
                  </Button>
                  <Button size="sm" variant="destructive" class="btn-press" @click="rejectProfile(user.id)">
                    <X class="mr-1 h-4 w-4" />
                    拒绝
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </template>

    <Dialog v-model:open="rejectDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>拒绝模型</DialogTitle>
          <DialogDescription>请输入拒绝原因，将通知上传者</DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div class="space-y-2">
            <Label>拒绝原因</Label>
            <Textarea 
              v-model="rejectReason" 
              placeholder="请输入拒绝原因..."
              :rows="3"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="rejectDialog = false">取消</Button>
          <Button variant="destructive" @click="rejectModel" :disabled="rejecting">
            <Loader2 v-if="rejecting" class="mr-2 h-4 w-4 animate-spin" />
            确认拒绝
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
