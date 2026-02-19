<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { adminApi, userApi } from '@/lib/api'
import type { User, PaginatedResponse } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
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
import { Loader2, Shield, ShieldCheck, ShieldAlert, Ban, Search, Users } from 'lucide-vue-next'

const authStore = useAuthStore()

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

const deleteDialog = ref(false)
const deleteUserId = ref<string>('')
const deleting = ref(false)

const isSuperAdmin = computed(() => authStore.user?.role === 'super_admin')

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

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
    await fetchUsers()
  } catch (error: any) {
    alert(error.response?.data?.message || '操作失败')
  }
}

async function removeAdmin(userId: string) {
  try {
    await adminApi.removeAdmin(userId)
    await fetchUsers()
  } catch (error: any) {
    alert(error.response?.data?.message || '操作失败')
  }
}

function openBanDialog(userId: string) {
  banUserId.value = userId
  banReason.value = ''
  banDialog.value = true
}

async function banUser() {
  if (!banReason.value.trim()) {
    alert('请输入封禁原因')
    return
  }
  
  banning.value = true
  try {
    await adminApi.banUser(banUserId.value, banReason.value)
    banDialog.value = false
    await fetchUsers()
  } catch (error: any) {
    alert(error.response?.data?.message || '操作失败')
  } finally {
    banning.value = false
  }
}

async function unbanUser(userId: string) {
  try {
    await adminApi.unbanUser(userId)
    await fetchUsers()
  } catch (error: any) {
    alert(error.response?.data?.message || '操作失败')
  }
}

function openDeleteDialog(userId: string) {
  deleteUserId.value = userId
  deleteDialog.value = true
}

async function deleteUser() {
  deleting.value = true
  try {
    await userApi.delete(deleteUserId.value)
    deleteDialog.value = false
    await fetchUsers()
  } catch (error: any) {
    alert(error.response?.data?.message || '操作失败')
  } finally {
    deleting.value = false
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
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">用户管理</h1>
        <p class="text-muted-foreground">管理系统用户和权限</p>
      </div>
      <div class="flex items-center gap-2">
        <Users class="h-5 w-5 text-muted-foreground" />
        <span class="text-sm text-muted-foreground">共 {{ total }} 个用户</span>
      </div>
    </div>

    <Card>
      <CardHeader>
        <div class="flex items-center gap-4">
          <div class="relative flex-1 max-w-sm">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input 
              v-model="search" 
              placeholder="搜索用户..." 
              class="pl-10"
            />
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex justify-center py-8">
          <Loader2 class="h-8 w-8 animate-spin text-primary" />
        </div>
        
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead>用户</TableHead>
              <TableHead>邮箱</TableHead>
              <TableHead>角色</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>注册时间</TableHead>
              <TableHead class="text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="user in users" :key="user.id">
              <TableCell>
                <div class="flex items-center gap-3">
                  <div class="h-10 w-10 rounded-full bg-muted flex items-center justify-center">
                    <span class="text-sm font-medium">{{ user.username?.charAt(0).toUpperCase() }}</span>
                  </div>
                  <div>
                    <p class="font-medium">{{ user.username }}</p>
                    <p class="text-sm text-muted-foreground">{{ user.bio || '暂无简介' }}</p>
                  </div>
                </div>
              </TableCell>
              <TableCell>{{ user.email }}</TableCell>
              <TableCell>
                <Badge :variant="getRoleBadgeVariant(user.role)">
                  <component :is="getRoleIcon(user.role)" class="h-3 w-3 mr-1" />
                  {{ getRoleLabel(user.role) }}
                </Badge>
              </TableCell>
              <TableCell>
                <Badge v-if="user.is_banned" variant="destructive">
                  <Ban class="h-3 w-3 mr-1" />
                  已封禁
                </Badge>
                <Badge v-else variant="default" class="bg-green-500">
                  正常
                </Badge>
              </TableCell>
              <TableCell>
                {{ new Date(user.created_at).toLocaleDateString() }}
              </TableCell>
              <TableCell class="text-right">
                <div class="flex items-center justify-end gap-2">
                  <template v-if="isSuperAdmin && user.role !== 'super_admin'">
                    <Button 
                      v-if="user.role === 'user'"
                      size="sm" 
                      variant="outline"
                      @click="setAdmin(user.id)"
                    >
                      设为管理员
                    </Button>
                    <Button 
                      v-else-if="user.role === 'admin'"
                      size="sm" 
                      variant="outline"
                      @click="removeAdmin(user.id)"
                    >
                      取消管理员
                    </Button>
                  </template>
                  
                  <template v-if="!user.is_banned && user.role !== 'super_admin'">
                    <Button 
                      size="sm" 
                      variant="destructive"
                      @click="openBanDialog(user.id)"
                    >
                      <Ban class="h-4 w-4 mr-1" />
                      封禁
                    </Button>
                  </template>
                  
                  <template v-if="user.is_banned">
                    <Button
                      size="sm"
                      variant="outline"
                      @click="unbanUser(user.id)"
                    >
                      解封
                    </Button>
                  </template>

                  <Button
                    size="sm"
                    variant="destructive"
                    @click="openDeleteDialog(user.id)"
                  >
                    删除
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
        
        <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-4">
          <Button 
            variant="outline" 
            size="sm"
            :disabled="page === 1"
            @click="page--; fetchUsers()"
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
            @click="page++; fetchUsers()"
          >
            下一页
          </Button>
        </div>
      </CardContent>
    </Card>

    <Dialog v-model:open="banDialog">
      <DialogContent>
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
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="banDialog = false">取消</Button>
          <Button variant="destructive" @click="banUser" :disabled="banning">
            <Loader2 v-if="banning" class="h-4 w-4 mr-2 animate-spin" />
            确认封禁
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="deleteDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>删除用户</DialogTitle>
          <DialogDescription>
            确定要删除该用户吗？此操作不可逆，用户的所有数据将被删除。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="deleteDialog = false">取消</Button>
          <Button variant="destructive" @click="deleteUser" :disabled="deleting">
            <Loader2 v-if="deleting" class="h-4 w-4 mr-2 animate-spin" />
            确认删除
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
