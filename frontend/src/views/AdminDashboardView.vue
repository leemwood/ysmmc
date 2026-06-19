<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/lib/api'
import type { Model, User, PaginatedResponse } from '@/types'
import AdminLayout from '@/components/admin/AdminLayout.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
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
import { LayoutDashboard, Package, Users, Download, Check, X, Loader2, FileText } from 'lucide-vue-next'
import { getModelImageUrl } from '@/utils/image'
import { useToast } from '@/composables/useToast'

const { toast } = useToast()

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
    toast.success('模型已通过审核')
    fetchPendingModels()
    fetchStats()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

function openRejectDialog(id: string) {
  rejectModelId.value = id
  rejectReason.value = ''
  rejectDialog.value = true
}

async function rejectModel() {
  if (!rejectReason.value.trim()) {
    toast.warning('请输入拒绝原因')
    return
  }

  rejecting.value = true
  try {
    await adminApi.rejectModel(rejectModelId.value, rejectReason.value)
    toast.success('模型已拒绝')
    rejectDialog.value = false
    fetchPendingModels()
    fetchStats()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  } finally {
    rejecting.value = false
  }
}

async function approveProfile(id: string) {
  try {
    await adminApi.approveProfile(id)
    toast.success('资料已通过审核')
    fetchPendingProfiles()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

async function rejectProfile(id: string) {
  try {
    await adminApi.rejectProfile(id)
    toast.success('资料变更已拒绝')
    fetchPendingProfiles()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

onMounted(async () => {
  await Promise.all([fetchStats(), fetchPendingModels(), fetchPendingProfiles()])
  loading.value = false
})
</script>

<template>
  <AdminLayout>
    <div class="flex flex-wrap gap-2">
      <Button
        :variant="activeTab === 'overview' ? 'default' : 'outline'"
        size="sm"
        class="btn-press"
        @click="activeTab = 'overview'"
      >
        <LayoutDashboard class="mr-1 h-4 w-4" />
        统计概览
      </Button>
      <Button
        :variant="activeTab === 'models' ? 'default' : 'outline'"
        size="sm"
        class="btn-press"
        @click="activeTab = 'models'"
      >
        <Package class="mr-1 h-4 w-4" />
        待审核模型
        <Badge v-if="stats.pending_models > 0" variant="destructive" class="ml-2 text-xs">{{ stats.pending_models }}</Badge>
      </Button>
      <Button
        :variant="activeTab === 'profiles' ? 'default' : 'outline'"
        size="sm"
        class="btn-press"
        @click="activeTab = 'profiles'"
      >
        <Users class="mr-1 h-4 w-4" />
        待审核资料
        <Badge v-if="pendingProfiles.length > 0" variant="destructive" class="ml-2 text-xs">{{ pendingProfiles.length }}</Badge>
      </Button>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <div v-if="activeTab === 'overview'" class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4">
        <Card class="card-hover">
          <CardContent class="flex flex-col items-center gap-2 p-4 text-center sm:flex-row sm:gap-4 sm:p-6 sm:text-left">
            <div class="rounded-full bg-primary/10 p-2 sm:p-3">
              <Users class="h-5 w-5 text-primary sm:h-6 sm:w-6" />
            </div>
            <div>
              <p class="text-xs text-muted-foreground sm:text-sm">总用户数</p>
              <p class="text-xl font-bold sm:text-2xl">{{ stats.total_users }}</p>
            </div>
          </CardContent>
        </Card>

        <Card class="card-hover">
          <CardContent class="flex flex-col items-center gap-2 p-4 text-center sm:flex-row sm:gap-4 sm:p-6 sm:text-left">
            <div class="rounded-full bg-primary/10 p-2 sm:p-3">
              <Package class="h-5 w-5 text-primary sm:h-6 sm:w-6" />
            </div>
            <div>
              <p class="text-xs text-muted-foreground sm:text-sm">总模型数</p>
              <p class="text-xl font-bold sm:text-2xl">{{ stats.total_models }}</p>
            </div>
          </CardContent>
        </Card>

        <Card class="card-hover">
          <CardContent class="flex flex-col items-center gap-2 p-4 text-center sm:flex-row sm:gap-4 sm:p-6 sm:text-left">
            <div class="rounded-full bg-yellow-500/10 p-2 sm:p-3">
              <FileText class="h-5 w-5 text-yellow-500 sm:h-6 sm:w-6" />
            </div>
            <div>
              <p class="text-xs text-muted-foreground sm:text-sm">待审核</p>
              <p class="text-xl font-bold sm:text-2xl">{{ stats.pending_models }}</p>
            </div>
          </CardContent>
        </Card>

        <Card class="card-hover">
          <CardContent class="flex flex-col items-center gap-2 p-4 text-center sm:flex-row sm:gap-4 sm:p-6 sm:text-left">
            <div class="rounded-full bg-green-500/10 p-2 sm:p-3">
              <Download class="h-5 w-5 text-green-500 sm:h-6 sm:w-6" />
            </div>
            <div>
              <p class="text-xs text-muted-foreground sm:text-sm">总下载量</p>
              <p class="text-xl font-bold sm:text-2xl">{{ stats.total_downloads }}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-else-if="activeTab === 'models'">
        <div v-if="pendingModels.length === 0" class="text-center py-12 text-muted-foreground">
          暂无待审核模型
        </div>

        <div v-else class="space-y-3 sm:space-y-4">
          <Card v-for="model in pendingModels" :key="model.id" class="overflow-hidden card-hover">
            <CardContent class="p-4">
              <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div class="flex items-center gap-3 sm:gap-4">
                  <div class="h-14 w-14 flex-shrink-0 overflow-hidden rounded bg-muted sm:h-16 sm:w-16">
                    <img v-if="model.image_id || model.image_url" :src="getModelImageUrl(model.image_id, model.image_url)" class="h-full w-full rounded object-cover" />
                  </div>
                  <div class="min-w-0 flex-1">
                    <RouterLink :to="`/model/${model.id}`" class="font-medium hover:underline line-clamp-1">
                      {{ model.title }}
                    </RouterLink>
                    <p class="mt-0.5 line-clamp-1 text-sm text-muted-foreground">
                      {{ model.description || '暂无描述' }}
                    </p>
                    <p class="mt-0.5 text-xs text-muted-foreground">
                      上传者: {{ model.user?.username }}
                    </p>
                  </div>
                </div>
                <div class="flex flex-col gap-2 sm:flex-row sm:flex-shrink-0">
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
          <Card v-for="user in pendingProfiles" :key="user.id" class="overflow-hidden card-hover">
            <CardContent class="p-4">
              <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div class="flex items-center gap-3 sm:gap-4">
                  <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-muted text-sm font-medium sm:h-12 sm:w-12 sm:text-base">
                    {{ user.username?.charAt(0).toUpperCase() }}
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="font-medium">{{ user.username }}</p>
                    <p class="line-clamp-1 text-sm text-muted-foreground">
                      新用户名: {{ user.pending_changes?.username || '无变更' }}
                    </p>
                    <p class="line-clamp-1 text-sm text-muted-foreground">
                      新简介: {{ user.pending_changes?.bio || '无变更' }}
                    </p>
                  </div>
                </div>
                <div class="flex flex-col gap-2 sm:flex-row sm:flex-shrink-0">
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
      <DialogContent class="max-w-[calc(100%-2rem)] sm:max-w-lg">
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
  </AdminLayout>
</template>
