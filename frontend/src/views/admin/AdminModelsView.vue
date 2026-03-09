<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { adminApi } from '@/lib/api'
import type { Model } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardContent } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Search, Trash2, Eye, ChevronLeft, ChevronRight, Loader2 } from 'lucide-vue-next'
import { getModelImageUrl } from '@/utils/image'

const router = useRouter()

const models = ref<Model[]>([])
const loading = ref(true)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)
const search = ref('')
const statusFilter = ref('all')
const deleteDialog = ref(false)
const deleteTarget = ref<Model | null>(null)
const deleteLoading = ref(false)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const statusOptions = [
  { value: 'all', label: '全部' },
  { value: 'pending', label: '待审核' },
  { value: 'approved', label: '已通过' },
  { value: 'rejected', label: '已拒绝' },
]

async function fetchModels() {
  loading.value = true
  try {
    const status = statusFilter.value === 'all' ? '' : statusFilter.value
    const response = await adminApi.listAllModels(page.value, pageSize.value, status, search.value)
    const respData = response.data as unknown as { data: Model[]; total: number }
    models.value = respData.data || []
    total.value = respData.total || 0
  } catch (error) {
    console.error('Failed to fetch models:', error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  fetchModels()
}

function handleStatusChange(value: string) {
  statusFilter.value = value
  page.value = 1
  fetchModels()
}

function prevPage() {
  if (page.value > 1) {
    page.value--
    fetchModels()
  }
}

function nextPage() {
  if (page.value < totalPages.value) {
    page.value++
    fetchModels()
  }
}

function viewModel(id: string) {
  router.push(`/model/${id}`)
}

function confirmDelete(model: Model) {
  deleteTarget.value = model
  deleteDialog.value = true
}

async function deleteModel() {
  if (!deleteTarget.value) return
  
  deleteLoading.value = true
  try {
    await adminApi.deleteModel(deleteTarget.value.id)
    deleteDialog.value = false
    deleteTarget.value = null
    fetchModels()
  } catch (error) {
    console.error('Failed to delete model:', error)
  } finally {
    deleteLoading.value = false
  }
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}

function getStatusBadgeVariant(status: string) {
  switch (status) {
    case 'approved':
      return 'default'
    case 'pending':
      return 'secondary'
    case 'rejected':
      return 'destructive'
    default:
      return 'outline'
  }
}

function getStatusLabel(status: string) {
  switch (status) {
    case 'approved':
      return '已通过'
    case 'pending':
      return '待审核'
    case 'rejected':
      return '已拒绝'
    default:
      return status
  }
}

onMounted(fetchModels)
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6">模型管理</h1>

    <div class="flex flex-wrap gap-4 mb-6">
      <div class="flex-1 min-w-[200px]">
        <div class="relative">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            v-model="search"
            placeholder="搜索模型..."
            class="pl-10"
            @keyup.enter="handleSearch"
          />
        </div>
      </div>
      <div class="flex gap-2">
        <Button
          v-for="option in statusOptions"
          :key="option.value"
          :variant="statusFilter === option.value ? 'default' : 'outline'"
          size="sm"
          @click="handleStatusChange(option.value)"
        >
          {{ option.label }}
        </Button>
      </div>
      <Button @click="handleSearch">搜索</Button>
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <Card v-for="i in 8" :key="i">
        <CardContent class="p-4">
          <Skeleton class="aspect-[4/3] w-full mb-3" />
          <Skeleton class="h-5 w-3/4 mb-2" />
          <Skeleton class="h-4 w-1/2" />
        </CardContent>
      </Card>
    </div>

    <div v-else-if="models.length === 0" class="text-center py-12 text-muted-foreground">
      暂无模型数据
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <Card v-for="m in models" :key="m.id" class="overflow-hidden">
        <div class="aspect-[4/3] w-full bg-muted">
          <img
            v-if="m.image_id || m.image_url"
            :src="getModelImageUrl(m.image_id, m.image_url)"
            :alt="m.title"
            class="h-full w-full object-cover"
          />
          <div v-else class="flex h-full items-center justify-center text-muted-foreground">
            无预览图
          </div>
        </div>
        <CardContent class="p-4">
          <div class="flex items-start justify-between gap-2">
            <h3 class="font-semibold line-clamp-1">{{ m.title }}</h3>
            <Badge :variant="getStatusBadgeVariant(m.status)">
              {{ getStatusLabel(m.status) }}
            </Badge>
          </div>
          <p class="mt-1 text-sm text-muted-foreground line-clamp-2">
            {{ m.description || '暂无描述' }}
          </p>
          <div class="mt-2 text-xs text-muted-foreground">
            <span>作者: {{ m.user?.username || '未知' }}</span>
            <span class="mx-2">|</span>
            <span>{{ formatDate(m.created_at) }}</span>
          </div>
          <div class="mt-3 flex gap-2">
            <Button variant="outline" size="sm" @click="viewModel(m.id)">
              <Eye class="mr-1 h-3 w-3" />
              查看
            </Button>
            <Button variant="destructive" size="sm" @click="confirmDelete(m)">
              <Trash2 class="mr-1 h-3 w-3" />
              删除
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>

    <div v-if="totalPages > 1" class="mt-6 flex items-center justify-center gap-4">
      <Button variant="outline" size="sm" :disabled="page === 1" @click="prevPage">
        <ChevronLeft class="h-4 w-4" />
        上一页
      </Button>
      <span class="text-sm text-muted-foreground">
        第 {{ page }} / {{ totalPages }} 页，共 {{ total }} 条
      </span>
      <Button variant="outline" size="sm" :disabled="page === totalPages" @click="nextPage">
        下一页
        <ChevronRight class="h-4 w-4" />
      </Button>
    </div>

    <Dialog v-model:open="deleteDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>确认删除</DialogTitle>
          <DialogDescription>
            确定要删除模型 "{{ deleteTarget?.title }}" 吗？此操作不可撤销。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="deleteDialog = false">取消</Button>
          <Button variant="destructive" :disabled="deleteLoading" @click="deleteModel">
            <Loader2 v-if="deleteLoading" class="mr-2 h-4 w-4 animate-spin" />
            删除
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
