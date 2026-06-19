<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { adminApi, announcementApi } from '@/lib/api'
import type { Announcement, PaginatedResponse } from '@/types'
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
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Loader2, Plus, Pencil, Trash2, Megaphone, Power, PowerOff, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'

const { toast } = useToast()

const announcements = ref<Announcement[]>([])
const loading = ref(true)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const editDialog = ref(false)
const deleteDialog = ref(false)
const editingId = ref<string>('')
const editTitle = ref('')
const editContent = ref('')
const saving = ref(false)
const deleting = ref(false)
const isCreate = ref(false)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const announcementColumns: ColumnConfig<Announcement>[] = [
  { key: 'title', title: '标题', width: '200px' },
  { key: 'content', title: '内容' },
  { key: 'status', title: '状态', width: '100px' },
  { key: 'created_at', title: '创建时间', width: '120px', formatter: (row) => new Date(row.created_at).toLocaleDateString() },
  { key: 'actions', title: '操作', align: 'right', width: '220px' },
]

async function fetchAnnouncements() {
  loading.value = true
  try {
    const response = await announcementApi.listAll(page.value, pageSize.value)
    const data = response.data.data as PaginatedResponse<Announcement>
    announcements.value = data.items
    total.value = data.total
  } catch (error) {
    console.error('Failed to fetch announcements:', error)
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  isCreate.value = true
  editingId.value = ''
  editTitle.value = ''
  editContent.value = ''
  editDialog.value = true
}

function openEditDialog(announcement: Announcement) {
  isCreate.value = false
  editingId.value = announcement.id
  editTitle.value = announcement.title
  editContent.value = announcement.content
  editDialog.value = true
}

function openDeleteDialog(id: string) {
  editingId.value = id
  deleteDialog.value = true
}

async function saveAnnouncement() {
  if (!editTitle.value.trim() || !editContent.value.trim()) {
    toast.warning('请填写标题和内容')
    return
  }

  saving.value = true
  try {
    if (isCreate.value) {
      await adminApi.createAnnouncement({
        title: editTitle.value,
        content: editContent.value
      })
      toast.success('公告已创建')
    } else {
      await adminApi.updateAnnouncement(editingId.value, {
        title: editTitle.value,
        content: editContent.value
      })
      toast.success('公告已更新')
    }
    editDialog.value = false
    await fetchAnnouncements()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

async function deleteAnnouncement() {
  deleting.value = true
  try {
    await adminApi.deleteAnnouncement(editingId.value)
    toast.success('公告已删除')
    deleteDialog.value = false
    await fetchAnnouncements()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '删除失败')
  } finally {
    deleting.value = false
  }
}

async function toggleActive(announcement: Announcement) {
  try {
    await adminApi.updateAnnouncement(announcement.id, {
      is_active: !announcement.is_active
    })
    announcement.is_active = !announcement.is_active
    toast.success(announcement.is_active ? '公告已启用' : '公告已禁用')
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

onMounted(fetchAnnouncements)
</script>

<template>
  <AdminLayout title="公告管理" description="管理系统公告内容">
    <template #actions>
      <div class="flex items-center gap-3 sm:gap-4">
        <div class="flex items-center gap-2">
          <Megaphone class="h-5 w-5 text-muted-foreground" />
          <span class="text-sm text-muted-foreground">共 {{ total }} 条公告</span>
        </div>
        <Button class="btn-press" @click="openCreateDialog">
          <Plus class="mr-2 h-4 w-4" />
          新建
        </Button>
      </div>
    </template>

    <div class="surface p-4 sm:p-6">
      <ResponsiveTable
        :columns="announcementColumns"
        :rows="announcements"
        :loading="loading"
        :skeleton-rows="5"
      >
        <template #title="{ row }">
          <span class="font-medium line-clamp-1">{{ row.title }}</span>
        </template>

        <template #content="{ row }">
          <p class="line-clamp-2 text-sm text-muted-foreground">
            {{ row.content }}
          </p>
        </template>

        <template #status="{ row }">
          <Badge :variant="row.is_active ? 'success' : 'secondary'">
            {{ row.is_active ? '已启用' : '已禁用' }}
          </Badge>
        </template>

        <template #actions="{ row }">
          <div class="flex flex-wrap items-center justify-end gap-2">
            <Button 
              size="sm" 
              :variant="row.is_active ? 'outline' : 'default'"
              class="btn-press"
              @click="toggleActive(row)"
            >
              <component :is="row.is_active ? PowerOff : Power" class="h-4 w-4 sm:mr-1" />
              <span class="hidden sm:inline">{{ row.is_active ? '禁用' : '启用' }}</span>
            </Button>
            <Button 
              size="sm" 
              variant="outline"
              class="btn-press"
              @click="openEditDialog(row)"
            >
              <Pencil class="h-4 w-4 sm:mr-1" />
              <span class="hidden sm:inline">编辑</span>
            </Button>
            <Button 
              size="sm" 
              variant="destructive"
              class="btn-press"
              @click="openDeleteDialog(row.id)"
            >
              <Trash2 class="h-4 w-4 sm:mr-1" />
              <span class="hidden sm:inline">删除</span>
            </Button>
          </div>
        </template>
      </ResponsiveTable>
      
      <div v-if="announcements.length === 0 && !loading" class="py-8 text-center text-muted-foreground">
        暂无公告
      </div>
      
      <div v-if="totalPages > 1" class="mt-4 flex items-center justify-center gap-2 sm:mt-6">
        <Button 
          variant="outline" 
          size="sm"
          class="btn-press"
          :disabled="page === 1"
          @click="page--; fetchAnnouncements()"
        >
          <ChevronLeft class="h-4 w-4 sm:mr-1" />
          <span class="hidden sm:inline">上一页</span>
        </Button>
        <span class="px-2 text-sm text-muted-foreground">
          {{ page }} / {{ totalPages }}
        </span>
        <Button 
          variant="outline" 
          size="sm"
          class="btn-press"
          :disabled="page === totalPages"
          @click="page++; fetchAnnouncements()"
        >
          <span class="hidden sm:inline">下一页</span>
          <ChevronRight class="h-4 w-4 sm:ml-1" />
        </Button>
      </div>
    </div>

    <Dialog v-model:open="editDialog">
      <DialogContent class="max-w-[calc(100%-2rem)] sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{{ isCreate ? '新建公告' : '编辑公告' }}</DialogTitle>
          <DialogDescription>
            {{ isCreate ? '创建一条新的系统公告' : '修改公告内容' }}
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div>
            <Label>标题</Label>
            <Input 
              v-model="editTitle" 
              placeholder="请输入公告标题..."
              class="h-11"
            />
          </div>
          <div>
            <Label>内容</Label>
            <Textarea 
              v-model="editContent" 
              placeholder="请输入公告内容..."
              :rows="5"
              class="min-h-[120px]"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="editDialog = false">取消</Button>
          <Button @click="saveAnnouncement" :disabled="saving">
            <Loader2 v-if="saving" class="mr-2 h-4 w-4 animate-spin" />
            {{ isCreate ? '创建' : '保存' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="deleteDialog">
      <DialogContent class="max-w-[calc(100%-2rem)] sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>删除公告</DialogTitle>
          <DialogDescription>
            确定要删除这条公告吗？此操作不可撤销。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="deleteDialog = false">取消</Button>
          <Button variant="destructive" @click="deleteAnnouncement" :disabled="deleting">
            <Loader2 v-if="deleting" class="mr-2 h-4 w-4 animate-spin" />
            确认删除
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </AdminLayout>
</template>
