<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi, modelApi, favoriteApi, uploadApi, authApi } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import type { Model, Favorite, PaginatedResponse } from '@/types'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { RouterLink } from 'vue-router'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { User, Save, Loader2, Heart, Package, Camera, Key, Mail } from 'lucide-vue-next'

const authStore = useAuthStore()

const activeTab = ref('profile')
const username = ref('')
const bio = ref('')
const loading = ref(false)
const models = ref<Model[]>([])
const favorites = ref<Favorite[]>([])
const loadingModels = ref(true)
const loadingFavorites = ref(true)
const message = ref('')

const avatarFile = ref<File | null>(null)
const avatarPreview = ref<string | null>(null)
const uploadingAvatar = ref(false)

const passwordDialog = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const changingPassword = ref(false)
const passwordMessage = ref('')

const emailDialog = ref(false)
const newEmail = ref('')
const changingEmail = ref(false)
const emailMessage = ref('')

async function loadProfile() {
  if (authStore.user) {
    username.value = authStore.user.username
    bio.value = authStore.user.bio || ''
  }
}

async function loadModels() {
  loadingModels.value = true
  try {
    const response = await modelApi.list(1, 100)
    const allModels = (response.data.data as PaginatedResponse<Model>).items
    models.value = allModels.filter(m => m.user_id === authStore.user?.id)
  } catch (error) {
    console.error('Failed to load models:', error)
  } finally {
    loadingModels.value = false
  }
}

async function loadFavorites() {
  loadingFavorites.value = true
  try {
    const response = await favoriteApi.list(1, 100)
    favorites.value = (response.data.data as PaginatedResponse<Favorite>).items
  } catch (error) {
    console.error('Failed to load favorites:', error)
  } finally {
    loadingFavorites.value = false
  }
}

async function handleUpdateProfile() {
  loading.value = true
  message.value = ''

  try {
    await userApi.updateMe({
      username: username.value !== authStore.user?.username ? username.value : undefined,
      bio: bio.value,
    })
    await authStore.fetchUser()
    if (authStore.isAdmin) {
      message.value = '资料更新成功'
    } else {
      message.value = '资料更新已提交审核'
    }
  } catch (error: any) {
    message.value = error.response?.data?.message || '更新失败'
  } finally {
    loading.value = false
  }
}

const MAX_AVATAR_SIZE = 2 * 1024 * 1024

function handleAvatarChange(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    
    if (file.size > MAX_AVATAR_SIZE) {
      message.value = '图片大小不能超过2MB'
      return
    }
    
    if (!file.type.startsWith('image/')) {
      message.value = '请选择有效的图片文件'
      return
    }
    
    avatarFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      avatarPreview.value = e.target?.result as string
    }
    reader.readAsDataURL(avatarFile.value)
    message.value = ''
  }
}

async function uploadAvatar() {
  if (!avatarFile.value) return

  uploadingAvatar.value = true
  try {
    const uploadResponse = await uploadApi.uploadImage(avatarFile.value)
    const avatarUrl = uploadResponse.data.data?.url

    if (avatarUrl) {
      await userApi.updateMe({ avatar_url: avatarUrl })
      await authStore.fetchUser()
      avatarFile.value = null
      avatarPreview.value = null
      message.value = '头像更新成功'
    }
  } catch (error: any) {
    message.value = error.response?.data?.message || '头像上传失败'
  } finally {
    uploadingAvatar.value = false
  }
}

function openPasswordDialog() {
  oldPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  passwordMessage.value = ''
  passwordDialog.value = true
}

async function changePassword() {
  if (!oldPassword.value || !newPassword.value || !confirmPassword.value) {
    passwordMessage.value = '请填写所有字段'
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    passwordMessage.value = '两次输入的密码不一致'
    return
  }

  if (newPassword.value.length < 6) {
    passwordMessage.value = '密码长度至少6位'
    return
  }

  changingPassword.value = true
  try {
    await userApi.changePassword({
      old_password: oldPassword.value,
      new_password: newPassword.value,
    })
    passwordDialog.value = false
    message.value = '密码修改成功'
  } catch (error: any) {
    passwordMessage.value = error.response?.data?.message || '密码修改失败'
  } finally {
    changingPassword.value = false
  }
}

function openEmailDialog() {
  newEmail.value = ''
  emailMessage.value = ''
  emailDialog.value = true
}

async function changeEmail() {
  if (!newEmail.value) {
    emailMessage.value = '请输入新邮箱地址'
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(newEmail.value)) {
    emailMessage.value = '请输入有效的邮箱地址'
    return
  }

  changingEmail.value = true
  try {
    await authApi.changeEmail(newEmail.value)
    emailDialog.value = false
    message.value = '验证邮件已发送到新邮箱，请查收'
  } catch (error: any) {
    emailMessage.value = error.response?.data?.message || '发送失败'
  } finally {
    changingEmail.value = false
  }
}

onMounted(() => {
  loadProfile()
  loadModels()
  loadFavorites()
})
</script>

<template>
  <div class="mx-auto max-w-4xl px-4 py-8">
    <div class="mb-6 flex gap-4 border-b">
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'profile' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'profile'"
      >
        <User class="h-4 w-4" />
        个人资料
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'models' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'models'"
      >
        <Package class="h-4 w-4" />
        我的模型
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="activeTab === 'favorites' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        @click="activeTab = 'favorites'"
      >
        <Heart class="h-4 w-4" />
        我的收藏
      </button>
    </div>

    <div v-if="activeTab === 'profile'">
      <Card>
        <CardHeader>
          <CardTitle>个人资料</CardTitle>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleUpdateProfile" class="space-y-4">
            <div v-if="message" class="rounded-md bg-primary/10 p-3 text-sm text-primary">
              {{ message }}
            </div>

            <div class="flex items-start gap-6">
              <div class="relative group">
                <Avatar class="h-24 w-24">
                  <AvatarImage v-if="avatarPreview || authStore.user?.avatar_url" :src="avatarPreview || authStore.user?.avatar_url || undefined">
                    <User class="h-12 w-12 text-muted-foreground" />
                  </AvatarImage>
                </Avatar>
                <label class="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full opacity-0 group-hover:opacity-100 cursor-pointer transition-opacity">
                  <Camera class="h-6 w-6 text-white" />
                  <input type="file" accept="image/*" class="hidden" @change="handleAvatarChange" />
                </label>
              </div>
              <div class="flex-1 space-y-3">
                <div v-if="avatarPreview" class="flex gap-2">
                  <Button size="sm" @click="uploadAvatar" :disabled="uploadingAvatar">
                    <Loader2 v-if="uploadingAvatar" class="mr-2 h-4 w-4 animate-spin" />
                    保存头像
                  </Button>
                  <Button size="sm" variant="outline" @click="avatarPreview = null; avatarFile = null">
                    取消
                  </Button>
                </div>
                <p class="text-sm text-muted-foreground">点击头像更换图片</p>
              </div>
            </div>

            <div class="flex items-center gap-4 pt-2">
              <div>
                <p class="font-medium">{{ authStore.user?.username }}</p>
                <p class="text-sm text-muted-foreground">{{ authStore.user?.email }}</p>
                <div class="flex items-center gap-2 mt-1">
                  <Badge v-if="authStore.user?.role === 'super_admin'" variant="default">站长</Badge>
                  <Badge v-else-if="authStore.isAdmin" variant="secondary">管理员</Badge>
                  <Badge v-if="authStore.user?.email_verified" variant="outline" class="text-green-600">已验证</Badge>
                </div>
              </div>
            </div>

            <div class="space-y-2">
              <Label for="username">用户名</Label>
              <Input id="username" v-model="username" placeholder="请输入用户名" />
            </div>

            <div class="space-y-2">
              <Label for="bio">个人简介</Label>
              <Textarea id="bio" v-model="bio" placeholder="介绍一下自己..." :rows="3" />
            </div>

            <div v-if="authStore.user?.profile_status === 'pending_review'" class="rounded-md bg-yellow-500/10 p-3 text-sm text-yellow-600">
              资料修改正在审核中
            </div>

            <div class="flex gap-2">
              <Button type="submit" :disabled="loading">
                <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
                <Save v-else class="mr-2 h-4 w-4" />
                保存修改
              </Button>
              <Button type="button" variant="outline" @click="openPasswordDialog">
                <Key class="mr-2 h-4 w-4" />
                修改密码
              </Button>
              <Button type="button" variant="outline" @click="openEmailDialog">
                <Mail class="mr-2 h-4 w-4" />
                修改邮箱
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>

    <div v-else-if="activeTab === 'models'">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-xl font-semibold">我的模型</h2>
        <RouterLink to="/upload">
          <Button size="sm">上传模型</Button>
        </RouterLink>
      </div>

      <div v-if="loadingModels" class="grid gap-4 sm:grid-cols-2">
        <Skeleton v-for="i in 4" :key="i" class="h-32" />
      </div>

      <div v-else-if="models.length === 0" class="text-center py-12 text-muted-foreground">
        暂无模型，<RouterLink to="/upload" class="text-primary hover:underline">去上传</RouterLink>
      </div>

      <div v-else class="grid gap-4 sm:grid-cols-2">
        <Card v-for="model in models" :key="model.id">
          <CardContent class="p-4">
            <div class="flex items-start gap-4">
              <div class="h-16 w-16 flex-shrink-0 rounded bg-muted">
                <img v-if="model.image_url" :src="model.image_url" class="h-full w-full object-cover rounded" />
              </div>
              <div class="flex-1 min-w-0">
                <RouterLink :to="`/model/${model.id}`" class="font-medium hover:underline line-clamp-1">
                  {{ model.title }}
                </RouterLink>
                <div class="mt-1 flex items-center gap-2">
                  <Badge :variant="model.status === 'approved' ? 'default' : model.status === 'rejected' ? 'destructive' : 'secondary'">
                    {{ model.status === 'approved' ? '已通过' : model.status === 'rejected' ? '已拒绝' : '审核中' }}
                  </Badge>
                  <span class="text-xs text-muted-foreground">{{ model.downloads }} 下载</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>

    <div v-else-if="activeTab === 'favorites'">
      <h2 class="mb-4 text-xl font-semibold">我的收藏</h2>

      <div v-if="loadingFavorites" class="grid gap-4 sm:grid-cols-2">
        <Skeleton v-for="i in 4" :key="i" class="h-32" />
      </div>

      <div v-else-if="favorites.length === 0" class="text-center py-12 text-muted-foreground">
        暂无收藏
      </div>

      <div v-else class="grid gap-4 sm:grid-cols-2">
        <Card v-for="favorite in favorites" :key="favorite.id">
          <CardContent class="p-4">
            <div class="flex items-start gap-4">
              <div class="h-16 w-16 flex-shrink-0 rounded bg-muted">
                <img v-if="favorite.model?.image_url" :src="favorite.model.image_url" class="h-full w-full object-cover rounded" />
              </div>
              <div class="flex-1 min-w-0">
                <RouterLink :to="`/model/${favorite.model_id}`" class="font-medium hover:underline line-clamp-1">
                  {{ favorite.model?.title }}
                </RouterLink>
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ favorite.model?.downloads || 0 }} 下载
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>

    <Dialog v-model:open="passwordDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>修改密码</DialogTitle>
          <DialogDescription>请输入旧密码和新密码</DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div class="space-y-2">
            <Label>旧密码</Label>
            <Input v-model="oldPassword" type="password" placeholder="请输入旧密码" />
          </div>
          <div class="space-y-2">
            <Label>新密码</Label>
            <Input v-model="newPassword" type="password" placeholder="请输入新密码（至少6位）" />
          </div>
          <div class="space-y-2">
            <Label>确认密码</Label>
            <Input v-model="confirmPassword" type="password" placeholder="请再次输入新密码" />
          </div>
          <div v-if="passwordMessage" class="text-sm text-destructive">{{ passwordMessage }}</div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="passwordDialog = false">取消</Button>
          <Button @click="changePassword" :disabled="changingPassword">
            <Loader2 v-if="changingPassword" class="mr-2 h-4 w-4 animate-spin" />
            确认修改
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="emailDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>修改邮箱</DialogTitle>
          <DialogDescription>新邮箱需要验证后才能生效</DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div class="space-y-2">
            <Label>当前邮箱</Label>
            <Input :model-value="authStore.user?.email" disabled />
          </div>
          <div class="space-y-2">
            <Label>新邮箱</Label>
            <Input v-model="newEmail" type="email" placeholder="请输入新邮箱地址" />
          </div>
          <div v-if="emailMessage" class="text-sm text-destructive">{{ emailMessage }}</div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="emailDialog = false">取消</Button>
          <Button @click="changeEmail" :disabled="changingEmail">
            <Loader2 v-if="changingEmail" class="mr-2 h-4 w-4 animate-spin" />
            发送验证邮件
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
