<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { User as UserIcon, Package, Edit2, Save, X } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model } from '../types'

const userStore = useUserStore()

useHead({
  title: '个人中心 - YSM 模型站',
  meta: [
    { name: 'robots', content: 'noindex, nofollow' }
  ]
})

const userModels = ref<Model[]>([])
const loading = ref(true)

// Edit Profile State
const isEditing = ref(false)
const editUsername = ref('')
const editBio = ref('')
const avatarFile = ref<File | null>(null)
const avatarPreview = ref<string | null>(null)
const updateLoading = ref(false)
const updateMsg = ref('')
const updateError = ref(false)

const currentAvatar = computed(() => {
  if (avatarPreview.value) return avatarPreview.value
  if (userStore.profile?.avatar_url) return userStore.profile.avatar_url
  return null
})

const fetchUserModels = async () => {
  if (!userStore.user) return
  
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*')
    .eq('user_id', userStore.user.id)
    .order('created_at', { ascending: false })

  if (error) {
    console.error(error)
  } else {
    userModels.value = data as any
  }
  loading.value = false
}

const startEditing = () => {
  if (userStore.profile) {
    editUsername.value = userStore.profile.username
    editBio.value = userStore.profile.bio || ''
    isEditing.value = true
    updateMsg.value = ''
  }
}

const cancelEditing = () => {
  isEditing.value = false
  avatarFile.value = null
  avatarPreview.value = null
}

const handleAvatarSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    if (file.size > 3 * 1024 * 1024) {
      alert('头像大小不能超过 3MB')
      return
    }
    avatarFile.value = file
    avatarPreview.value = URL.createObjectURL(file)
  }
}

const saveProfile = async () => {
  if (!userStore.user || !userStore.profile) return
  
  if (editUsername.value.length < 1 || editUsername.value.length > 16) {
    updateError.value = true
    updateMsg.value = '用户名长度必须在 1-16 个字符之间'
    return
  }
  if (editBio.value.length > 300) {
    updateError.value = true
    updateMsg.value = '简介不能超过 300 字'
    return
  }

  updateLoading.value = true
  updateMsg.value = ''
  updateError.value = false

  try {
    let avatarUrl = userStore.profile.avatar_url
    
    // Upload avatar if changed
    if (avatarFile.value) {
      const fileExt = avatarFile.value.name.split('.').pop()
      const fileName = `avatar_${Math.random()}.${fileExt}`
      const filePath = `${userStore.user.id}/${fileName}`
      
      const { error: uploadError } = await supabase.storage
        .from('images')
        .upload(filePath, avatarFile.value)
        
      if (uploadError) throw uploadError
      
      const { data } = supabase.storage.from('images').getPublicUrl(filePath)
      avatarUrl = data.publicUrl
    }

    // Construct pending changes
    const pendingChanges = {
      username: editUsername.value !== userStore.profile.username ? editUsername.value : undefined,
      bio: editBio.value !== userStore.profile.bio ? editBio.value : undefined,
      avatar_url: avatarUrl !== userStore.profile.avatar_url ? avatarUrl : undefined
    }

    // Only update if there are changes
    if (Object.values(pendingChanges).some(v => v !== undefined)) {
      const { error } = await supabase
        .from('profiles')
        .update({
          pending_changes: pendingChanges,
          profile_status: 'pending_review'
        })
        .eq('id', userStore.user.id)

      if (error) throw error
      
      updateMsg.value = '资料更新请求已提交，等待管理员审核。'
      isEditing.value = false
      
      // Refresh profile locally to show status
      await userStore.fetchUser()
    } else {
      isEditing.value = false
    }

  } catch (e: any) {
    console.error(e)
    updateError.value = true
    updateMsg.value = '更新失败: ' + e.message
  } finally {
    updateLoading.value = false
  }
}

onMounted(() => {
  fetchUserModels()
})
</script>

<template>
  <div class="container">
    <div class="profile-header">
      <div class="profile-header-content">
        <div class="avatar-wrapper">
           <img v-if="currentAvatar" :src="currentAvatar" alt="Avatar" class="avatar-img">
           <div v-else class="profile-avatar">
             <UserIcon :size="48" />
           </div>
           
           <label v-if="isEditing" class="avatar-upload-overlay">
             <input type="file" accept="image/*" @change="handleAvatarSelect" hidden>
             <span>更换头像</span>
           </label>
        </div>
        
        <div class="profile-info" v-if="!isEditing">
          <div class="info-top">
             <h1>{{ userStore.profile?.username || '未命名用户' }}</h1>
             <span v-if="userStore.profile?.profile_status === 'pending_review'" class="status-badge pending">审核中</span>
          </div>
          <p class="bio">{{ userStore.profile?.bio || '暂无简介' }}</p>
          <p class="email">{{ userStore.user?.email }}</p>
        </div>

        <div class="profile-edit-form" v-else>
          <div class="form-group">
            <label>用户名 (1-16字符)</label>
            <input v-model="editUsername" type="text" class="input" maxlength="16">
          </div>
          <div class="form-group">
            <label>简介 (300字以内)</label>
            <textarea v-model="editBio" class="input textarea" rows="3" maxlength="300"></textarea>
            <div class="char-count">{{ editBio.length }}/300</div>
          </div>
        </div>
      </div>

      <div class="profile-actions">
        <button v-if="!isEditing" @click="startEditing" class="btn btn--secondary">
          <Edit2 :size="16" /> 编辑资料
        </button>
        <template v-else>
          <button @click="cancelEditing" class="btn btn--secondary" :disabled="updateLoading">
            <X :size="16" /> 取消
          </button>
          <button @click="saveProfile" class="btn btn--primary" :disabled="updateLoading">
            <Save :size="16" /> {{ updateLoading ? '提交中...' : '保存修改' }}
          </button>
        </template>
      </div>
    </div>
    
    <div v-if="updateMsg" :class="['alert', updateError ? 'alert-error' : 'alert-success']">
      {{ updateMsg }}
    </div>

    <div class="content-section">
      <div class="section-header">
        <h2>我的模型</h2>
        <RouterLink to="/upload" class="btn btn--primary btn--sm">上传新模型</RouterLink>
      </div>

      <div v-if="loading" class="loading">
        正在加载模型...
      </div>

      <div v-else-if="userModels.length === 0" class="empty-state">
        <Package :size="48" />
        <p>你还没有上传任何模型。</p>
        <RouterLink to="/upload" class="btn btn--primary">开始上传</RouterLink>
      </div>

      <div v-else class="model-list">
        <div v-for="model in userModels" :key="model.id" class="list-item">
          <div class="item-image">
             <img :src="model.image_url || 'https://via.placeholder.com/100x75?text=暂无预览'" :alt="model.title">
          </div>
          <div class="item-content">
            <h3>{{ model.title }}</h3>
            <p class="item-meta">
              <span class="status" :class="{ 'public': model.is_public }">
                {{ model.is_public ? '公开' : '私有' }}
              </span>
              <span>•</span>
              <span class="status" :class="model.status">
                {{ model.status === 'approved' ? '已通过' : (model.status === 'rejected' ? '已拒绝' : '审核中') }}
              </span>
              <span>•</span>
              <span>{{ model.downloads }} 次下载</span>
            </p>
          </div>
          <div class="item-actions">
            <RouterLink :to="`/model/${model.id}`" class="btn btn--secondary btn--sm">查看</RouterLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.profile-header {
  display: flex;
  flex-direction: column;
  gap: $spacing-lg;
  padding: $spacing-xl;
  background-color: white;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  margin-bottom: $spacing-xl;
  
  @media (min-width: 768px) {
    flex-direction: row;
    align-items: flex-start;
    justify-content: space-between;
  }
}

.profile-header-content {
  display: flex;
  gap: $spacing-lg;
  flex: 1;
  
  @media (max-width: 600px) {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
}

.avatar-wrapper {
  position: relative;
  width: 100px;
  height: 100px;
  flex-shrink: 0;
}

.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.profile-avatar {
  width: 100%;
  height: 100%;
  background-color: #e0e7ff;
  color: var(--color-primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s;
  font-size: 0.875rem;
  
  &:hover {
    opacity: 1;
  }
}

.profile-info {
  flex: 1;
  
  .info-top {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      margin-bottom: $spacing-xs;
      
      @media (max-width: 600px) {
        justify-content: center;
      }
  }

  h1 {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-text-main);
  }
  
  .bio {
    color: var(--color-text-main);
    margin-bottom: $spacing-sm;
    white-space: pre-wrap;
  }
  
  .email {
    color: var(--color-text-muted);
    font-size: 0.875rem;
  }
}

.profile-edit-form {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
  max-width: 400px;
}

.char-count {
  text-align: right;
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.profile-actions {
  display: flex;
  gap: $spacing-sm;
  align-self: flex-start;
}

.status-badge {
  font-size: 0.75rem;
  padding: 2px 8px;
  border-radius: $radius-full;
  background-color: #f3f4f6;
  color: var(--color-text-muted);
  
  &.pending {
    background-color: #fff7ed;
    color: #c2410c;
    border: 1px solid #ffedd5;
  }
}

.alert {
  padding: $spacing-md;
  border-radius: $radius-md;
  margin-bottom: $spacing-xl;
  text-align: center;
}

.alert-success {
  background-color: #ecfdf5;
  color: #047857;
  border: 1px solid #d1fae5;
}

.alert-error {
  background-color: #fef2f2;
  color: #b91c1c;
  border: 1px solid #fee2e2;
}

/* Existing styles... */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-lg;
  
  h2 {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--color-text-main);
  }
}

.model-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.list-item {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  background-color: white;
  padding: $spacing-md;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  transition: $transition-base;
  
  &:hover {
    border-color: var(--color-primary);
  }
}

.item-image {
  width: 100px;
  height: 75px;
  border-radius: $radius-md;
  overflow: hidden;
  background-color: #f3f4f6;
  flex-shrink: 0;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.item-content {
  flex: 1;
  
  h3 {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text-main);
    margin-bottom: $spacing-xs;
  }
}

.item-meta {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.status {
  font-weight: 500;
  
  &.public {
    color: var(--color-secondary);
  }
  
  &.approved { color: var(--color-secondary); }
  &.rejected { color: var(--color-danger); }
  &.pending { color: #f59e0b; }
}

.empty-state {
  text-align: center;
  padding: $spacing-2xl;
  color: var(--color-text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-md;
}

.loading {
  text-align: center;
  padding: $spacing-xl;
  color: var(--color-text-muted);
}
</style>
