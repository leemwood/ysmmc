<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { adminApi } from '@/lib/api'
import type { Model, PaginatedResponse } from '@/types'
import AdminLayout from '@/components/admin/AdminLayout.vue'
import ResponsiveTable from '@/components/admin/ResponsiveTable.vue'
import type { ColumnConfig } from '@/components/admin/ResponsiveTable.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
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
import { useToast } from '@/composables/useToast'

const router = useRouter()
const { toast } = useToast()

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

const modelColumns: ColumnConfig<Model>[] = [
  { key: 'model', title: '模型' },
  { key: 'status', title: '状态', width: '100px' },
  { key: 'author', title: '作者', width: '140px' },
  { key: 'created_at', title: '发布时间', width: '120px', formatter: (row) => new Date(row.created_at).toLocaleDateString('zh-CN') },
  { key: 'actions', title: '操作', align: 'right', width: '140px' },
]

async function fetchModels() {
  loading.value = true
  try {
    const status = statusFilter.value === 'all' ? '' : statusFilter.value
    const response = await adminApi.listAllModels(page.value, pageSize.value, status, search.value)
    const data = response.data.data as PaginatedResponse<Model>
    models.value = data.items || []
    total.value = data.total || 0
  } catch (error) {
    console.error('Failed to fetch models:', error)
    toast.error('获取模型列表失败')
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
    toast.success('模型已删除')
    deleteDialog.value = false
    deleteTarget.value = null
    fetchModels()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '删除失败')
  } finally {
    deleteLoading.value = false
  }
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
  <AdminLayout title="模型管理" description="管理所有用户上传的模型">
    <div class="flex flex-col gap-3 sm:flex-row sm:gap-4">
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          v-model="search"
          placeholder="搜索模型..."
          class="h-11 pl-10"
          @keyup.enter="handleSearch"
        />
      </div>
      <div class="flex gap-2 overflow-x-auto pb-1 sm:pb-0">
        <Button
          v-for="option in statusOptions"
          :key="option.value"
          :variant="statusFilter === option.value ? 'default' : 'outline'"
          size="sm"
          class="btn-press h-9 whitespace-nowrap"
          @click="handleStatusChange(option.value)"
        >
          {{ option.label }}
        </Button>
      </div>
      <Button class="btn-press h-9" @click="handleSearch">搜索</Button>
    </div>

    <div class="surface p-4 sm:p-6">
      <ResponsiveTable
        :columns="modelColumns"
        :rows="models"
        :loading="loading"
        :skeleton-rows="8"
      >
        <template #model="{ row }">
          <div class="flex items-center gap-3 sm:gap-4">
            <div class="h-14 w-14 flex-shrink-0 overflow-hidden rounded bg-muted sm:h-16 sm:w-16">
              <img
                v-if="row.image_id || row.image_url"
                :src="getModelImageUrl(row.image_id, row.image_url)"
                :alt="row.title"
                class="h-full w-full object-cover"
                loading="lazy"
                decoding="async"
              />
              <div v-else class="flex h-full items-center justify-center text-xs text-muted-foreground">
                无预览
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="font-medium line-clamp-1">{{ row.title }}</p>
              <p class="text-sm text-muted-foreground line-clamp-1">{{ row.description || '暂无描述' }}</p>
            </div>
          </div>
        </template>

        <template #status="{ row }">
          <Badge :variant="getStatusBadgeVariant(row.status)" class="text-xs">
            {{ getStatusLabel(row.status) }}
          </Badge>
        </template>

        <template #author="{ row }">
          <span class="text-sm text-muted-foreground">{{ row.user?.username || '未知' }}</span>
        </template>

        <template #actions="{ row }">
          <div class="flex flex-wrap items-center justify-end gap-2">
            <Button variant="outline" size="sm" class="btn-press" @click="viewModel(row.id)">
              <Eye class="mr-1 h-3 w-3" />
              查看
            </Button>
            <Button variant="destructive" size="sm" class="btn-press" @click="confirmDelete(row)">
              <Trash2 class="mr-1 h-3 w-3" />
              删除
            </Button>
          </div>
        </template>
      </ResponsiveTable>

      <div v-if="totalPages > 1" class="mt-4 flex items-center justify-center gap-2 sm:mt-6">
        <Button variant="outline" size="sm" class="btn-press" :disabled="page === 1" @click="prevPage">
          <ChevronLeft class="h-4 w-4 sm:mr-1" />
          <span class="hidden sm:inline">上一页</span>
        </Button>
        <span class="text-sm text-muted-foreground">
          <span class="hidden sm:inline">第 {{ page }} / {{ totalPages }} 页，共 {{ total }} 条</span>
          <span class="sm:hidden">{{ page }} / {{ totalPages }}</span>
        </span>
        <Button variant="outline" size="sm" class="btn-press" :disabled="page === totalPages" @click="nextPage">
          <span class="hidden sm:inline">下一页</span>
          <ChevronRight class="h-4 w-4 sm:ml-1" />
        </Button>
      </div>
    </div>

    <Dialog v-model:open="deleteDialog">
      <DialogContent class="max-w-[calc(100%-2rem)] sm:max-w-lg">
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
  </AdminLayout>
</template>
