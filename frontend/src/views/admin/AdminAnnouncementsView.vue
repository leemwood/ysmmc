<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { adminApi, announcementApi } from '@/lib/api'
import type { Announcement, PaginatedResponse } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
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
import { Loader2, Plus, Pencil, Trash2, Megaphone, Power, PowerOff } from 'lucide-vue-next'

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
    alert('请填写标题和内容')
    return
  }
  
  saving.value = true
  try {
    if (isCreate.value) {
      await adminApi.createAnnouncement({
        title: editTitle.value,
        content: editContent.value
      })
    } else {
      await adminApi.updateAnnouncement(editingId.value, {
        title: editTitle.value,
        content: editContent.value
      })
    }
    editDialog.value = false
    await fetchAnnouncements()
  } catch (error: any) {
    alert(error.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

async function deleteAnnouncement() {
  deleting.value = true
  try {
    await adminApi.deleteAnnouncement(editingId.value)
    deleteDialog.value = false
    await fetchAnnouncements()
  } catch (error: any) {
    alert(error.response?.data?.message || '删除失败')
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
  } catch (error: any) {
    alert(error.response?.data?.message || '操作失败')
  }
}

onMounted(fetchAnnouncements)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">公告管理</h1>
        <p class="text-muted-foreground">管理系统公告内容</p>
      </div>
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
          <Megaphone class="h-5 w-5 text-muted-foreground" />
          <span class="text-sm text-muted-foreground">共 {{ total }} 条公告</span>
        </div>
        <Button @click="openCreateDialog">
          <Plus class="h-4 w-4 mr-2" />
          新建公告
        </Button>
      </div>
    </div>

    <Card>
      <CardContent class="pt-6">
        <div v-if="loading" class="flex justify-center py-8">
          <Loader2 class="h-8 w-8 animate-spin text-primary" />
        </div>
        
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead class="w-[200px]">标题</TableHead>
              <TableHead>内容</TableHead>
              <TableHead class="w-[100px]">状态</TableHead>
              <TableHead class="w-[150px]">创建时间</TableHead>
              <TableHead class="w-[200px] text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="announcement in announcements" :key="announcement.id">
              <TableCell class="font-medium">{{ announcement.title }}</TableCell>
              <TableCell>
                <p class="line-clamp-2 text-sm text-muted-foreground">
                  {{ announcement.content }}
                </p>
              </TableCell>
              <TableCell>
                <Badge :variant="announcement.is_active ? 'default' : 'secondary'">
                  {{ announcement.is_active ? '已启用' : '已禁用' }}
                </Badge>
              </TableCell>
              <TableCell>
                {{ new Date(announcement.created_at).toLocaleDateString() }}
              </TableCell>
              <TableCell class="text-right">
                <div class="flex items-center justify-end gap-2">
                  <Button 
                    size="sm" 
                    :variant="announcement.is_active ? 'outline' : 'default'"
                    @click="toggleActive(announcement)"
                  >
                    <component :is="announcement.is_active ? PowerOff : Power" class="h-4 w-4 mr-1" />
                    {{ announcement.is_active ? '禁用' : '启用' }}
                  </Button>
                  <Button 
                    size="sm" 
                    variant="outline"
                    @click="openEditDialog(announcement)"
                  >
                    <Pencil class="h-4 w-4 mr-1" />
                    编辑
                  </Button>
                  <Button 
                    size="sm" 
                    variant="destructive"
                    @click="openDeleteDialog(announcement.id)"
                  >
                    <Trash2 class="h-4 w-4 mr-1" />
                    删除
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
        
        <div v-if="announcements.length === 0 && !loading" class="text-center py-8 text-muted-foreground">
          暂无公告
        </div>
        
        <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-4">
          <Button 
            variant="outline" 
            size="sm"
            :disabled="page === 1"
            @click="page--; fetchAnnouncements()"
          >
            上一页
          </Button>
          <span class="text-sm text-muted-foreground">
            {{ page }} / {{ totalPages }}
          </span>
          <Button 
            variant="outline" 
            size="sm"
            :disabled="page === totalPages"
            @click="page++; fetchAnnouncements()"
          >
            下一页
          </Button>
        </div>
      </CardContent>
    </Card>

    <Dialog v-model:open="editDialog">
      <DialogContent>
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
            />
          </div>
          <div>
            <Label>内容</Label>
            <Textarea 
              v-model="editContent" 
              placeholder="请输入公告内容..."
              :rows="5"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="editDialog = false">取消</Button>
          <Button @click="saveAnnouncement" :disabled="saving">
            <Loader2 v-if="saving" class="h-4 w-4 mr-2 animate-spin" />
            {{ isCreate ? '创建' : '保存' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="deleteDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>删除公告</DialogTitle>
          <DialogDescription>
            确定要删除这条公告吗？此操作不可撤销。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="deleteDialog = false">取消</Button>
          <Button variant="destructive" @click="deleteAnnouncement" :disabled="deleting">
            <Loader2 v-if="deleting" class="h-4 w-4 mr-2 animate-spin" />
            确认删除
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
