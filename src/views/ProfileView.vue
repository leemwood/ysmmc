<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { User as UserIcon, Package, Edit2, Save, X, Camera, Mail, CheckCircle, ShieldAlert } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model } from '../types'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ModelCard from '../components/ModelCard.vue'

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
  <div class="container py-xl">
    <div class="profile-header-card">
      <div class="profile-header-main">
        <div class="avatar-container">
           <img v-if="currentAvatar" :src="currentAvatar" :alt="`${userStore.profile?.username || '用户'} 的头像`" class="avatar-image">
           <div v-else class="avatar-placeholder" aria-hidden="true">
             <UserIcon :size="48" />
           </div>
           
           <label v-if="isEditing" class="avatar-edit-label" title="点击更换头像">
             <input type="file" accept="image/*" @change="handleAvatarSelect" class="sr-only">
             <Camera :size="20" aria-hidden="true" />
             <span>更换头像</span>
           </label>
        </div>
        
        <div class="profile-details" v-if="!isEditing">
          <div class="name-row">
             <h1>{{ userStore.profile?.username || '未命名用户' }}</h1>
             <span v-if="userStore.profile?.profile_status === 'pending_review'" class="status-badge pending" role="status">资料审核中</span>
          </div>
          <p class="bio-text">{{ userStore.profile?.bio || '这个人很懒，什么都没有留下。' }}</p>
          <div class="user-meta">
            <span class="meta-item"><Mail :size="16" aria-hidden="true" /> {{ userStore.user?.email }}</span>
          </div>
        </div>

        <div class="profile-edit-form" v-else>
          <div class="form-group">
            <label for="username">用户名 <span class="required" aria-hidden="true">*</span></label>
            <input 
              id="username" 
              v-model="editUsername" 
              type="text" 
              class="input" 
              maxlength="16"
              required
              aria-required="true"
            >
            <span class="input-hint">1-16 个字符</span>
          </div>
          <div class="form-group">
            <label for="bio">个人简介</label>
            <textarea 
              id="bio" 
              v-model="editBio" 
              class="input textarea" 
              rows="3" 
              maxlength="300"
              placeholder="介绍一下你自己..."
            ></textarea>
            <div class="char-counter" :class="{ 'at-limit': editBio.length >= 300 }" aria-hidden="true">
              {{ editBio.length }}/300
            </div>
          </div>
        </div>
      </div>

      <div class="profile-header-actions">
        <button v-if="!isEditing" @click="startEditing" class="btn btn--secondary" aria-label="编辑个人资料">
          <Edit2 :size="16" aria-hidden="true" /> 编辑资料
        </button>
        <div v-else class="edit-actions">
          <button @click="cancelEditing" class="btn btn--secondary" :disabled="updateLoading" aria-label="取消编辑">
            <X :size="16" aria-hidden="true" /> 取消
          </button>
          <button @click="saveProfile" class="btn btn--primary" :disabled="updateLoading" aria-label="保存修改">
            <template v-if="updateLoading">
              <LoadingSpinner :size="16" class="inline-spinner" /> 提交中...
            </template>
            <template v-else>
              <Save :size="16" aria-hidden="true" /> 保存修改
            </template>
          </button>
        </div>
      </div>
    </div>
    
    <Transition name="fade">
      <div v-if="updateMsg" :class="['alert', updateError ? 'alert-error' : 'alert-success']" role="alert" aria-live="polite">
        <ShieldAlert v-if="updateError" :size="18" aria-hidden="true" />
        <CheckCircle v-else :size="18" aria-hidden="true" />
        {{ updateMsg }}
      </div>
    </Transition>

    <div class="content-section">
      <div class="section-header">
        <div class="header-title">
          <Package :size="24" aria-hidden="true" />
          <h2>我的模型</h2>
        </div>
        <RouterLink to="/upload" class="btn btn--primary btn--sm">
          上传新模型
        </RouterLink>
      </div>

      <LoadingSpinner v-if="loading" message="正在加载您的模型..." />

      <div v-else-if="userModels.length === 0" class="empty-state">
        <div class="empty-icon" aria-hidden="true">
          <Package :size="48" />
        </div>
        <h3>还没有上传过模型</h3>
        <p>在这里分享你的创意，让更多人看到你的作品！</p>
        <RouterLink to="/upload" class="btn btn--primary">
          立即上传第一个模型
        </RouterLink>
      </div>

      <div v-else class="model-grid">
        <ModelCard 
          v-for="(model, index) in userModels" 
          :key="model.id" 
          :model="model"
          :index="index"
        />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.profile-header-card {
  background: white;
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-md;
  padding: $spacing-2xl;
  margin-bottom: $spacing-2xl;
  display: flex;
  flex-direction: column;
  gap: $spacing-xl;

  @media (min-width: 768px) {
    flex-direction: row;
    align-items: flex-start;
    justify-content: space-between;
  }
}

.profile-header-main {
  display: flex;
  flex-direction: column;
  gap: $spacing-xl;
  flex: 1;

  @media (min-width: 640px) {
    flex-direction: row;
    align-items: flex-start;
  }
}

.avatar-container {
  position: relative;
  width: 120px;
  height: 120px;
  flex-shrink: 0;

  .avatar-image, .avatar-placeholder {
    width: 100%;
    height: 100%;
    border-radius: $radius-xl;
    object-fit: cover;
    border: 4px solid white;
    box-shadow: $shadow-sm;
  }

  .avatar-placeholder {
    background-color: var(--color-bg-light);
    color: var(--color-text-muted);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--color-border);
  }

  .avatar-edit-label {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    color: white;
    border-radius: $radius-xl;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: $spacing-xs;
    cursor: pointer;
    opacity: 0;
    transition: $transition-base;
    font-size: 0.75rem;
    font-weight: 600;

    &:hover {
      opacity: 1;
    }
  }
}

.profile-details {
  .name-row {
    display: flex;
    align-items: center;
    gap: $spacing-md;
    margin-bottom: $spacing-sm;

    h1 {
      font-size: 1.75rem;
      font-weight: 800;
      color: var(--color-text-main);
      letter-spacing: -0.025em;
    }

    .status-badge {
      padding: 2px 8px;
      border-radius: $radius-full;
      font-size: 0.75rem;
      font-weight: 600;

      &.pending {
        background-color: #fef3c7;
        color: #92400e;
      }
    }
  }

  .bio-text {
    color: var(--color-text-main);
    line-height: 1.6;
    margin-bottom: $spacing-md;
    max-width: 600px;
  }

  .user-meta {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-md;
    color: var(--color-text-muted);
    font-size: 0.875rem;

    .meta-item {
      display: flex;
      align-items: center;
      gap: $spacing-xs;
    }
  }
}

.profile-edit-form {
  flex: 1;
  max-width: 500px;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;

  .form-group {
    label {
      display: block;
      font-size: 0.875rem;
      font-weight: 600;
      color: var(--color-text-main);
      margin-bottom: $spacing-xs;

      .required {
        color: var(--color-danger);
      }
    }

    .input-hint {
      display: block;
      font-size: 0.75rem;
      color: var(--color-text-muted);
      margin-top: $spacing-xs;
    }

    .char-counter {
      text-align: right;
      font-size: 0.75rem;
      color: var(--color-text-muted);
      margin-top: $spacing-xs;

      &.at-limit {
        color: var(--color-danger);
      }
    }
  }

  .textarea {
    resize: vertical;
    min-height: 100px;
  }
}

.profile-header-actions {
  .edit-actions {
    display: flex;
    gap: $spacing-sm;
  }
}

.inline-spinner {
  margin-right: $spacing-xs;
}

.alert {
  padding: $spacing-md $spacing-lg;
  border-radius: $radius-lg;
  margin-bottom: $spacing-xl;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: $spacing-sm;

  &-success {
    background-color: #ecfdf5;
    color: #065f46;
    border: 1px solid #a7f3d0;
  }

  &-error {
    background-color: #fef2f2;
    color: #991b1b;
    border: 1px solid #fecaca;
  }
}

.content-section {
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: $spacing-xl;

    .header-title {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      color: var(--color-primary);

      h2 {
        font-size: 1.5rem;
        font-weight: 700;
        color: var(--color-text-main);
      }
    }
  }
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: $spacing-xl;
}

.empty-state {
  text-align: center;
  padding: $spacing-3xl 0;
  background: white;
  border-radius: $radius-xl;
  border: 1px dashed var(--color-border);

  .empty-icon {
    font-size: 3rem;
    margin-bottom: $spacing-md;
  }

  h3 {
    font-size: 1.25rem;
    font-weight: 700;
    margin-bottom: $spacing-sm;
  }

  p {
    color: var(--color-text-muted);
    margin-bottom: $spacing-xl;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>
