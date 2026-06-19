<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { adminApi } from '@/lib/api'
import type { User, PaginatedResponse } from '@/types'
import { useAuthStore } from '@/stores/auth'
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
import { Loader2, Shield, ShieldCheck, ShieldAlert, Ban, Search, Users, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'

const authStore = useAuthStore()
const { toast } = useToast()

const users = ref<User[]>([])
const loading = ref(true)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const search = ref('')

const banDialog = ref(false)
const banUserId = ref<string>('')
const banReason = ref('')
const banning = ref(false)

const isSuperAdmin = computed(() => authStore.user?.role === 'super_admin')
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const userColumns: ColumnConfig<User>[] = [
  { key: 'user', title: '用户' },
  { key: 'email', title: '邮箱' },
  { key: 'role', title: '角色' },
  { key: 'status', title: '状态' },
  { key: 'created_at', title: '注册时间', formatter: (row) => new Date(row.created_at).toLocaleDateString() },
  { key: 'actions', title: '操作', align: 'right' },
]

async function fetchUsers() {
  loading.value = true
  try {
    const response = await adminApi.listUsers(page.value, pageSize.value)
    const data = response.data.data as PaginatedResponse<User>
    users.value = data.items
    total.value = data.total
  } catch (error) {
    console.error('Failed to fetch users:', error)
  } finally {
    loading.value = false
  }
}

async function setAdmin(userId: string) {
  try {
    await adminApi.setAdmin(userId)
    toast.success('已设置为管理员')
    await fetchUsers()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

async function removeAdmin(userId: string) {
  try {
    await adminApi.removeAdmin(userId)
    toast.success('已取消管理员')
    await fetchUsers()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

function openBanDialog(userId: string) {
  banUserId.value = userId
  banReason.value = ''
  banDialog.value = true
}

async function banUser() {
  if (!banReason.value.trim()) {
    toast.warning('请输入封禁原因')
    return
  }

  banning.value = true
  try {
    await adminApi.banUser(banUserId.value, banReason.value)
    toast.success('用户已封禁')
    banDialog.value = false
    await fetchUsers()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  } finally {
    banning.value = false
  }
}

async function unbanUser(userId: string) {
  try {
    await adminApi.unbanUser(userId)
    toast.success('用户已解封')
    await fetchUsers()
  } catch (error: any) {
    toast.error(error.response?.data?.message || '操作失败')
  }
}

function getRoleBadgeVariant(role: string): 'default' | 'secondary' | 'outline' {
  switch (role) {
    case 'super_admin':
      return 'default'
    case 'admin':
      return 'secondary'
    default:
      return 'outline'
  }
}

function getRoleIcon(role: string) {
  switch (role) {
    case 'super_admin':
      return ShieldAlert
    case 'admin':
      return ShieldCheck
    default:
      return Shield
  }
}

function getRoleLabel(role: string) {
  switch (role) {
    case 'super_admin':
      return '站长'
    case 'admin':
      return '管理员'
    default:
      return '用户'
  }
}

onMounted(fetchUsers)
</script>

<template>
  <AdminLayout title="用户管理" description="管理系统用户和权限">
    <template #actions>
      <div class="flex items-center gap-2">
        <Users class="h-5 w-5 text-muted-foreground" />
        <span class="text-sm text-muted-foreground">共 {{ total }} 个用户</span>
      </div>
    </template>

    <div class="surface p-4 sm:p-6">
      <div class="mb-4 sm:mb-6">
        <div class="relative max-w-full sm:max-w-sm">
          <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input 
            v-model="search" 
            placeholder="搜索用户..." 
            class="h-11 pl-10"
          />
        </div>
      </div>

      <ResponsiveTable
        :columns="userColumns"
        :rows="users"
        :loading="loading"
        :skeleton-rows="5"
      >
        <template #user="{ row }">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-muted">
              <span class="text-sm font-medium">{{ row.username?.charAt(0).toUpperCase() }}</span>
            </div>
            <div class="min-w-0">
              <p class="font-medium">{{ row.username }}</p>
              <p class="text-sm text-muted-foreground line-clamp-1">{{ row.bio || '暂无简介' }}</p>
            </div>
          </div>
        </template>

        <template #email="{ row }">
          <span class="line-clamp-1 text-sm">{{ row.email }}</span>
        </template>

        <template #role="{ row }">
          <Badge :variant="getRoleBadgeVariant(row.role)">
            <component :is="getRoleIcon(row.role)" class="mr-1 h-3 w-3" />
            {{ getRoleLabel(row.role) }}
          </Badge>
        </template>

        <template #status="{ row }">
          <Badge v-if="row.is_banned" variant="destructive">
            <Ban class="mr-1 h-3 w-3" />
            已封禁
          </Badge>
          <Badge v-else variant="success">
            正常
          </Badge>
        </template>

        <template #actions="{ row }">
          <div class="flex flex-wrap items-center justify-end gap-2">
            <template v-if="isSuperAdmin && row.role !== 'super_admin'">
              <Button 
                v-if="row.role === 'user'"
                size="sm" 
                variant="outline"
                class="btn-press"
                @click="setAdmin(row.id)"
              >
                设为管理员
              </Button>
              <Button 
                v-else-if="row.role === 'admin'"
                size="sm" 
                variant="outline"
                class="btn-press"
                @click="removeAdmin(row.id)"
              >
                取消管理员
              </Button>
            </template>
            
            <template v-if="!row.is_banned && row.role !== 'super_admin'">
              <Button 
                size="sm" 
                variant="destructive"
                class="btn-press"
                @click="openBanDialog(row.id)"
              >
                <Ban class="h-4 w-4 sm:mr-1" />
                <span class="hidden sm:inline">封禁</span>
              </Button>
            </template>
            
            <template v-if="row.is_banned">
              <Button 
                size="sm" 
                variant="outline"
                class="btn-press"
                @click="unbanUser(row.id)"
              >
                解封
              </Button>
            </template>
          </div>
        </template>
      </ResponsiveTable>
      
      <div v-if="totalPages > 1" class="mt-4 flex items-center justify-center gap-2 sm:mt-6">
        <Button 
          variant="outline" 
          size="sm"
          class="btn-press"
          :disabled="page === 1"
          @click="page--; fetchUsers()"
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
          @click="page++; fetchUsers()"
        >
          <span class="hidden sm:inline">下一页</span>
          <ChevronRight class="h-4 w-4 sm:ml-1" />
        </Button>
      </div>
    </div>

    <Dialog v-model:open="banDialog">
      <DialogContent class="max-w-[calc(100%-2rem)] sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>封禁用户</DialogTitle>
          <DialogDescription>
            请输入封禁原因，用户将无法登录和使用平台功能。
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div>
            <Label>封禁原因</Label>
            <Textarea 
              v-model="banReason" 
              placeholder="请输入封禁原因..."
              :rows="3"
              class="h-24"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="banDialog = false">取消</Button>
          <Button variant="destructive" @click="banUser" :disabled="banning">
            <Loader2 v-if="banning" class="mr-2 h-4 w-4 animate-spin" />
            确认封禁
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </AdminLayout>
</template>
