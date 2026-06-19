<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
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
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { RouterLink } from 'vue-router'
import ModelCard from '@/components/ModelCard.vue'
import {
  User,
  Save,
  Loader2,
  Heart,
  Package,
  Camera,
  Key,
  Mail,
  UploadCloud,
  Inbox,
  AlertCircle,
  Pencil,
} from 'lucide-vue-next'
import { getAvatarUrl } from '@/utils/image'

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
const messageType = ref<'success' | 'error'>('success')

const avatarFile = ref<File | null>(null)
const avatarPreview = ref<string | null>(null)
const uploadingAvatar = ref(false)
const uploadProgress = ref(0)
const isDragging = ref(false)

const passwordSheet = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const changingPassword = ref(false)
const passwordMessage = ref('')

const emailSheet = ref(false)
const newEmail = ref('')
const changingEmail = ref(false)
const emailMessage = ref('')

const isMobile = ref(false)
let resizeHandler: (() => void) | null = null

const sheetSide = computed(() => (isMobile.value ? 'bottom' : 'right'))

function checkMobile() {
  isMobile.value = window.innerWidth < 640
}

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

function setMessage(text: string, type: 'success' | 'error' = 'success') {
  message.value = text
  messageType.value = type
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
    setMessage(authStore.isAdmin ? '资料更新成功' : '资料更新已提交审核')
  } catch (error: any) {
    setMessage(error.response?.data?.message || '更新失败', 'error')
  } finally {
    loading.value = false
  }
}

const MAX_AVATAR_SIZE = 2 * 1024 * 1024

function validateAvatarFile(file: File): boolean {
  if (file.size > MAX_AVATAR_SIZE) {
    setMessage('图片大小不能超过2MB', 'error')
    return false
  }

  if (!file.type.startsWith('image/')) {
    setMessage('请选择有效的图片文件', 'error')
    return false
  }

  return true
}

function readAvatarFile(file: File) {
  avatarFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    avatarPreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
  message.value = ''
}

function handleAvatarChange(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    if (validateAvatarFile(file)) {
      readAvatarFile(file)
    }
    target.value = ''
  }
}

function handleDragOver(event: DragEvent) {
  event.preventDefault()
  isDragging.value = true
}

function handleDragLeave(event: DragEvent) {
  event.preventDefault()
  isDragging.value = false
}

function handleDrop(event: DragEvent) {
  event.preventDefault()
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (files && files[0]) {
    const file = files[0]
    if (validateAvatarFile(file)) {
      readAvatarFile(file)
    }
  }
}

async function uploadAvatar() {
  if (!avatarFile.value) return

  uploadingAvatar.value = true
  uploadProgress.value = 0
  const progressInterval = window.setInterval(() => {
    if (uploadProgress.value < 90) {
      uploadProgress.value += Math.floor(Math.random() * 10) + 5
      if (uploadProgress.value > 90) uploadProgress.value = 90
    }
  }, 200)

  try {
    const uploadResponse = await uploadApi.uploadImage(avatarFile.value)
    const avatarId = uploadResponse.data.data?.file_id

    if (avatarId) {
      await userApi.updateMe({ avatar_id: avatarId })
      await authStore.fetchUser()
      avatarFile.value = null
      avatarPreview.value = null
      setMessage('头像更新成功')
    }
  } catch (error: any) {
    setMessage(error.response?.data?.message || '头像上传失败', 'error')
  } finally {
    window.clearInterval(progressInterval)
    uploadProgress.value = 100
    uploadingAvatar.value = false
    window.setTimeout(() => {
      uploadProgress.value = 0
    }, 500)
  }
}

function cancelAvatarUpload() {
  avatarPreview.value = null
  avatarFile.value = null
  uploadProgress.value = 0
}

function openPasswordSheet() {
  oldPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  passwordMessage.value = ''
  passwordSheet.value = true
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
    passwordSheet.value = false
    setMessage('密码修改成功')
  } catch (error: any) {
    passwordMessage.value = error.response?.data?.message || '密码修改失败'
  } finally {
    changingPassword.value = false
  }
}

function openEmailSheet() {
  newEmail.value = ''
  emailMessage.value = ''
  emailSheet.value = true
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
    emailSheet.value = false
    setMessage('验证邮件已发送到新邮箱，请查收')
  } catch (error: any) {
    emailMessage.value = error.response?.data?.message || '发送失败'
  } finally {
    changingEmail.value = false
  }
}

onMounted(() => {
  checkMobile()
  resizeHandler = checkMobile
  window.addEventListener('resize', resizeHandler)
  loadProfile()
  loadModels()
  loadFavorites()
})

onUnmounted(() => {
  if (resizeHandler) {
    window.removeEventListener('resize', resizeHandler)
  }
})
</script>

<template>
  <div class="mx-auto max-w-4xl px-4 py-6 sm:py-8">
    <Tabs v-model="activeTab" class="mb-6">
      <TabsList class="grid h-auto w-full grid-cols-3 p-1 sm:h-10 sm:w-auto sm:inline-grid">
        <TabsTrigger value="profile" class="flex items-center justify-center gap-2 py-2.5 sm:py-1.5">
          <User class="h-4 w-4" />
          <span class="hidden sm:inline">个人资料</span>
          <span class="sm:hidden">资料</span>
        </TabsTrigger>
        <TabsTrigger value="models" class="flex items-center justify-center gap-2 py-2.5 sm:py-1.5">
          <Package class="h-4 w-4" />
          <span class="hidden sm:inline">我的模型</span>
          <span class="sm:hidden">模型</span>
        </TabsTrigger>
        <TabsTrigger value="favorites" class="flex items-center justify-center gap-2 py-2.5 sm:py-1.5">
          <Heart class="h-4 w-4" />
          <span class="hidden sm:inline">我的收藏</span>
          <span class="sm:hidden">收藏</span>
        </TabsTrigger>
      </TabsList>

      <TabsContent value="profile" class="animate-fade-in focus-visible:outline-none">
        <Card>
          <CardHeader>
            <CardTitle>个人资料</CardTitle>
          </CardHeader>
          <CardContent>
            <form @submit.prevent="handleUpdateProfile" class="space-y-5">
              <Alert v-if="message" :variant="messageType === 'error' ? 'destructive' : 'default'" class="animate-fade-in">
                <AlertCircle class="h-4 w-4" />
                <AlertDescription>{{ message }}</AlertDescription>
              </Alert>

              <div class="flex flex-col items-start gap-4 sm:flex-row sm:gap-6">
                <div
                  class="group relative self-center sm:self-start"
                  :class="[
                    'rounded-full border-2 border-dashed p-1 transition-colors',
                    isDragging ? 'border-primary bg-primary/5' : 'border-transparent',
                  ]"
                  @dragover="handleDragOver"
                  @dragleave="handleDragLeave"
                  @drop="handleDrop"
                >
                  <Avatar class="h-24 w-24 sm:h-28 sm:w-28">
                    <AvatarImage v-if="avatarPreview || authStore.user?.avatar_id || authStore.user?.avatar_url" :src="avatarPreview || getAvatarUrl(authStore.user?.avatar_id, authStore.user?.avatar_url) || undefined">
                      <User class="h-12 w-12 text-muted-foreground" />
                    </AvatarImage>
                    <div v-else class="flex h-full w-full items-center justify-center bg-muted">
                      <User class="h-12 w-12 text-muted-foreground" />
                    </div>
                  </Avatar>
                  <label
                    class="absolute inset-0 m-1 flex cursor-pointer items-center justify-center rounded-full bg-black/50 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100"
                  >
                    <Camera class="h-6 w-6 text-white" />
                    <input type="file" accept="image/*" class="hidden" @change="handleAvatarChange" />
                  </label>
                </div>
                <div class="flex-1 space-y-3 text-center sm:text-left">
                  <div class="text-sm text-muted-foreground">
                    <p class="flex items-center justify-center gap-1 sm:justify-start">
                      <UploadCloud class="h-4 w-4" />
                      点击头像或拖拽图片到此处上传
                    </p>
                    <p class="text-xs">支持 JPG、PNG、GIF、WEBP，最大 2MB</p>
                  </div>
                  <div v-if="avatarPreview" class="flex flex-col gap-2">
                    <div class="flex gap-2 justify-center sm:justify-start">
                      <Button size="sm" @click="uploadAvatar" :disabled="uploadingAvatar" class="h-9">
                        <Loader2 v-if="uploadingAvatar" class="mr-2 h-4 w-4 animate-spin" />
                        保存头像
                      </Button>
                      <Button size="sm" variant="outline" @click="cancelAvatarUpload" :disabled="uploadingAvatar" class="h-9">
                        取消
                      </Button>
                    </div>
                    <Progress v-if="uploadingAvatar" :model-value="uploadProgress" :max="100" show-value class="max-w-xs" />
                  </div>
                </div>
              </div>

              <div class="flex flex-col items-start gap-2 pt-2 sm:flex-row sm:items-center sm:gap-4">
                <div class="text-center sm:text-left">
                  <p class="font-medium">{{ authStore.user?.username }}</p>
                  <p class="text-sm text-muted-foreground">{{ authStore.user?.email }}</p>
                  <div class="mt-1 flex flex-wrap items-center justify-center gap-2 sm:justify-start">
                    <Badge v-if="authStore.user?.role === 'super_admin'" variant="default">站长</Badge>
                    <Badge v-else-if="authStore.isAdmin" variant="secondary">管理员</Badge>
                    <Badge v-if="authStore.user?.email_verified" variant="outline" class="text-green-600">已验证</Badge>
                  </div>
                </div>
              </div>

              <Separator />

              <div class="space-y-2">
                <Label for="username">用户名</Label>
                <Input id="username" v-model="username" placeholder="请输入用户名" class="h-11" />
              </div>

              <div class="space-y-2">
                <Label for="bio">个人简介</Label>
                <Textarea id="bio" v-model="bio" placeholder="介绍一下自己..." :rows="3" class="min-h-[80px]" />
              </div>

              <Alert v-if="authStore.user?.profile_status === 'pending_review'" variant="default" class="animate-fade-in border-yellow-500/50 text-yellow-600 dark:text-yellow-400 [&>svg]:text-yellow-600 dark:[&>svg]:text-yellow-400">
                <AlertCircle class="h-4 w-4" />
                <AlertDescription>资料修改正在审核中，审核通过前对外展示的信息不会更新。</AlertDescription>
              </Alert>

              <div class="flex flex-col gap-2 sm:flex-row">
                <Button type="submit" :disabled="loading" :loading="loading" class="btn-press h-10 w-full sm:flex-1">
                  <Save v-if="!loading" class="mr-2 h-4 w-4" />
                  保存修改
                </Button>
                <Button type="button" variant="outline" @click="openPasswordSheet" class="btn-press h-10 w-full sm:flex-1">
                  <Key class="mr-2 h-4 w-4" />
                  修改密码
                </Button>
                <Button type="button" variant="outline" @click="openEmailSheet" class="btn-press h-10 w-full sm:flex-1">
                  <Mail class="mr-2 h-4 w-4" />
                  修改邮箱
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="models" class="animate-fade-in focus-visible:outline-none">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-semibold sm:text-xl">我的模型</h2>
          <RouterLink to="/upload">
            <Button size="sm" class="btn-press h-9">上传模型</Button>
          </RouterLink>
        </div>

        <div v-if="loadingModels" class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Skeleton v-for="i in 4" :key="i" class="h-48" />
        </div>

        <div v-else-if="models.length === 0" class="surface flex flex-col items-center justify-center py-16 text-center animate-fade-in">
          <Inbox class="h-12 w-12 text-muted-foreground opacity-50" />
          <p class="mt-4 text-muted-foreground">暂无模型</p>
          <RouterLink to="/upload" class="mt-2 text-primary hover:underline">去上传第一个模型</RouterLink>
        </div>

        <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div v-for="model in models" :key="model.id" class="relative group">
            <ModelCard :model="model" />
            <RouterLink
              :to="`/model/${model.id}/edit`"
              class="absolute right-2 top-2 flex h-8 w-8 items-center justify-center rounded-md bg-background/90 text-foreground opacity-0 shadow-sm backdrop-blur-sm transition-opacity group-hover:opacity-100 focus:opacity-100 focus-ring"
              :aria-label="`编辑 ${model.title}`"
            >
              <Pencil class="h-4 w-4" />
            </RouterLink>
          </div>
        </div>
      </TabsContent>

      <TabsContent value="favorites" class="animate-fade-in focus-visible:outline-none">
        <h2 class="mb-4 text-lg font-semibold sm:text-xl">我的收藏</h2>

        <div v-if="loadingFavorites" class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Skeleton v-for="i in 4" :key="i" class="h-48" />
        </div>

        <div v-else-if="favorites.length === 0" class="surface flex flex-col items-center justify-center py-16 text-center animate-fade-in">
          <Heart class="h-12 w-12 text-muted-foreground opacity-50" />
          <p class="mt-4 text-muted-foreground">暂无收藏</p>
          <RouterLink to="/" class="mt-2 text-primary hover:underline">去发现模型</RouterLink>
        </div>

        <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <template v-for="favorite in favorites" :key="favorite.id">
            <ModelCard v-if="favorite.model" :model="favorite.model" />
            <Card v-else class="flex items-center justify-center p-4 text-center">
              <CardContent>
                <p class="text-sm text-muted-foreground">模型信息不可用</p>
              </CardContent>
            </Card>
          </template>
        </div>
      </TabsContent>
    </Tabs>

    <Sheet v-model:open="passwordSheet">
      <SheetContent :side="sheetSide" class="sm:max-w-sm">
        <SheetHeader>
          <SheetTitle>修改密码</SheetTitle>
          <SheetDescription>请输入旧密码和新密码</SheetDescription>
        </SheetHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="old-password">旧密码</Label>
            <Input id="old-password" v-model="oldPassword" type="password" placeholder="请输入旧密码" class="h-11" />
          </div>
          <div class="space-y-2">
            <Label for="new-password">新密码</Label>
            <Input id="new-password" v-model="newPassword" type="password" placeholder="请输入新密码（至少6位）" class="h-11" />
          </div>
          <div class="space-y-2">
            <Label for="confirm-password">确认密码</Label>
            <Input id="confirm-password" v-model="confirmPassword" type="password" placeholder="请再次输入新密码" class="h-11" />
          </div>
          <Alert v-if="passwordMessage" variant="destructive" class="animate-fade-in">
            <AlertCircle class="h-4 w-4" />
            <AlertDescription>{{ passwordMessage }}</AlertDescription>
          </Alert>
        </div>
        <SheetFooter class="flex-col-reverse sm:flex-row">
          <Button variant="outline" @click="passwordSheet = false" class="h-10 w-full sm:w-auto">取消</Button>
          <Button @click="changePassword" :loading="changingPassword" class="h-10 w-full sm:w-auto">
            确认修改
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>

    <Sheet v-model:open="emailSheet">
      <SheetContent :side="sheetSide" class="sm:max-w-sm">
        <SheetHeader>
          <SheetTitle>修改邮箱</SheetTitle>
          <SheetDescription>新邮箱需要验证后才能生效</SheetDescription>
        </SheetHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="current-email">当前邮箱</Label>
            <Input id="current-email" :model-value="authStore.user?.email" disabled class="h-11" />
          </div>
          <div class="space-y-2">
            <Label for="new-email">新邮箱</Label>
            <Input id="new-email" v-model="newEmail" type="email" placeholder="请输入新邮箱地址" class="h-11" />
          </div>
          <Alert v-if="emailMessage" variant="destructive" class="animate-fade-in">
            <AlertCircle class="h-4 w-4" />
            <AlertDescription>{{ emailMessage }}</AlertDescription>
          </Alert>
        </div>
        <SheetFooter class="flex-col-reverse sm:flex-row">
          <Button variant="outline" @click="emailSheet = false" class="h-10 w-full sm:w-auto">取消</Button>
          <Button @click="changeEmail" :loading="changingEmail" class="h-10 w-full sm:w-auto">
            发送验证邮件
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  </div>
</template>
