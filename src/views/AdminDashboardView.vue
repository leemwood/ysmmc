<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { supabase } from '../supabase/client'
import { Check, X, User, Box } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model, Profile } from '../types'

useHead({
  title: '管理员控制台 - YSM 模型站',
  meta: [
    { name: 'robots', content: 'noindex, nofollow' }
  ]
})

const activeTab = ref<'models' | 'profiles'>('models')

const pendingModels = ref<Model[]>([])
const pendingProfiles = ref<Profile[]>([])
const loading = ref(true)

const fetchPendingModels = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*, profiles(username)')
    .or('status.eq.pending,update_status.eq.pending_review') // Fetch pending new OR pending updates
    .order('created_at', { ascending: false })

  if (error) console.error(error)
  else pendingModels.value = data as any
  loading.value = false
}

const fetchPendingProfiles = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('profiles')
    .select('*')
    .eq('profile_status', 'pending_review')
    .order('updated_at', { ascending: false })

  if (error) console.error(error)
  else pendingProfiles.value = data as any
  loading.value = false
}

const handleApproveModel = async (model: Model) => {
  // Case 1: New Model Approval
  if (model.status === 'pending') {
    const { error } = await supabase
      .from('models')
      .update({ status: 'approved' })
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  } 
  // Case 2: Model Update Approval
  else if (model.update_status === 'pending_review' && model.pending_changes) {
    const changes = model.pending_changes
    const updates: any = {
      update_status: 'idle',
      pending_changes: null,
      updated_at: new Date().toISOString()
    }
    
    // Apply changes
    if (changes.title) updates.title = changes.title
    if (changes.description) updates.description = changes.description
    if (changes.tags) updates.tags = changes.tags
    if (changes.file_path) updates.file_path = changes.file_path
    if (changes.image_url) updates.image_url = changes.image_url
    if (changes.is_public !== undefined) updates.is_public = changes.is_public

    const { error } = await supabase
      .from('models')
      .update(updates)
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  }
}

const handleRejectModel = async (model: Model) => {
  // Case 1: New Model Rejection
  if (model.status === 'pending') {
    const { error } = await supabase
      .from('models')
      .update({ status: 'rejected' })
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  }
  // Case 2: Model Update Rejection
  else if (model.update_status === 'pending_review') {
    const { error } = await supabase
      .from('models')
      .update({ 
        update_status: 'idle',
        pending_changes: null
      })
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  }
}

const handleApproveProfile = async (profile: Profile) => {
  if (!profile.pending_changes) return

  const changes = profile.pending_changes as any
  const updates: any = {
    profile_status: 'approved',
    pending_changes: null
  }
  
  if (changes.username) updates.username = changes.username
  if (changes.bio) updates.bio = changes.bio
  if (changes.avatar_url) updates.avatar_url = changes.avatar_url

  const { error } = await supabase
    .from('profiles')
    .update(updates)
    .eq('id', profile.id)

  if (!error) {
    pendingProfiles.value = pendingProfiles.value.filter(p => p.id !== profile.id)
  }
}

const handleRejectProfile = async (id: string) => {
  const { error } = await supabase
    .from('profiles')
    .update({ 
      profile_status: 'approved', // Revert status to approved but keep old data
      pending_changes: null 
    })
    .eq('id', id)

  if (!error) {
    pendingProfiles.value = pendingProfiles.value.filter(p => p.id !== id)
  }
}

const switchTab = (tab: 'models' | 'profiles') => {
  activeTab.value = tab
  if (tab === 'models') fetchPendingModels()
  else fetchPendingProfiles()
}

onMounted(() => {
  fetchPendingModels()
})
</script>

<template>
  <div class="container">
    <div class="admin-header">
      <h1>管理员控制台</h1>
      <p>审核用户上传的模型和资料更新</p>
    </div>

    <div class="tabs">
      <button 
        @click="switchTab('models')" 
        :class="['tab-btn', activeTab === 'models' ? 'active' : '']"
      >
        <Box :size="18" /> 模型审核
        <span v-if="pendingModels.length" class="badge">{{ pendingModels.length }}</span>
      </button>
      <button 
        @click="switchTab('profiles')" 
        :class="['tab-btn', activeTab === 'profiles' ? 'active' : '']"
      >
        <User :size="18" /> 资料审核
        <span v-if="pendingProfiles.length" class="badge">{{ pendingProfiles.length }}</span>
      </button>
    </div>

    <div v-if="loading" class="loading">
      正在加载...
    </div>

    <!-- Models List -->
    <div v-else-if="activeTab === 'models'">
      <div v-if="pendingModels.length === 0" class="empty-state">
        <div class="empty-icon">
          <Check :size="48" />
        </div>
        <p>太棒了！所有模型都已审核完毕。</p>
      </div>

      <div v-else class="audit-list">
        <div v-for="model in pendingModels" :key="model.id" class="audit-item">
          <div class="item-preview">
             <img :src="model.image_url || 'https://via.placeholder.com/150x100?text=暂无预览'" :alt="model.title">
          </div>
          
          <div class="item-info">
            <h3>
              {{ model.title }} 
              <span class="arrow" v-if="model.update_status === 'pending_review' && model.pending_changes?.title">→ {{ model.pending_changes.title }}</span>
              <span v-if="model.update_status === 'pending_review'" class="badge update-badge">更新申请</span>
              <span v-else class="badge new-badge">新上传</span>
            </h3>
            
            <template v-if="model.update_status === 'pending_review' && model.pending_changes">
               <div v-if="model.pending_changes.description" class="change-block">
                 <p class="label">描述变更</p>
                 <p class="new-val">{{ model.pending_changes.description }}</p>
               </div>
               <div v-if="model.pending_changes.file_path" class="change-block">
                 <p class="new-val">申请更新模型文件</p>
               </div>
               <div v-if="model.pending_changes.image_url" class="change-block">
                 <p class="new-val">申请更新预览图</p>
               </div>
            </template>

            <p class="item-meta">
              <span>上传者: {{ model.profiles?.username || '未知' }}</span>
              <span>•</span>
              <span>{{ new Date(model.created_at).toLocaleDateString() }}</span>
            </p>
            <p class="item-desc">{{ model.description || '无描述' }}</p>
            <div class="item-tags">
              <span v-for="tag in model.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
          </div>

          <div class="item-actions">
            <button @click="handleApproveModel(model)" class="btn btn--success btn--sm" title="通过">
              <Check :size="18" /> 通过
            </button>
            <button @click="handleRejectModel(model)" class="btn btn--danger btn--sm" title="拒绝">
              <X :size="18" /> 拒绝
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Profiles List -->
    <div v-else-if="activeTab === 'profiles'">
      <div v-if="pendingProfiles.length === 0" class="empty-state">
        <div class="empty-icon">
          <Check :size="48" />
        </div>
        <p>没有待审核的资料更新。</p>
      </div>

      <div v-else class="audit-list">
        <div v-for="profile in pendingProfiles" :key="profile.id" class="audit-item">
          <div class="item-preview profile-preview">
             <div class="old-new-avatar">
               <div class="avatar-box">
                 <span>当前</span>
                 <img :src="profile.avatar_url || 'https://via.placeholder.com/50'" class="avatar-mini">
               </div>
               <div v-if="profile.pending_changes?.avatar_url" class="avatar-box">
                 <span>新</span>
                 <img :src="profile.pending_changes.avatar_url" class="avatar-mini">
               </div>
             </div>
          </div>
          
          <div class="item-info">
            <h3>{{ profile.username }} <span class="arrow" v-if="profile.pending_changes?.username">→ {{ profile.pending_changes.username }}</span></h3>
            
            <div v-if="profile.pending_changes?.bio" class="change-block">
              <p class="label">简介变更:</p>
              <p class="old-val">{{ profile.bio || '(无)' }}</p>
              <p class="new-val">{{ profile.pending_changes.bio }}</p>
            </div>
            
            <div v-if="!profile.pending_changes?.username && !profile.pending_changes?.bio && profile.pending_changes?.avatar_url" class="change-block">
              <p>仅申请更新头像</p>
            </div>
          </div>

          <div class="item-actions">
            <button @click="handleApproveProfile(profile)" class="btn btn--success btn--sm" title="通过">
              <Check :size="18" /> 通过
            </button>
            <button @click="handleRejectProfile(profile.id)" class="btn btn--danger btn--sm" title="拒绝">
              <X :size="18" /> 拒绝
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.admin-header {
  margin-bottom: $spacing-lg;
  
  h1 {
    font-size: 2rem;
    font-weight: 700;
    color: var(--color-text-main);
  }
  
  p {
    color: var(--color-text-muted);
  }
}

.tabs {
  display: flex;
  gap: $spacing-md;
  margin-bottom: $spacing-xl;
  border-bottom: 1px solid var(--color-border);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  padding: $spacing-md;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  font-weight: 500;
  color: var(--color-text-muted);
  
  &.active {
    color: var(--color-primary);
    border-bottom-color: var(--color-primary);
  }
  
  &:hover:not(.active) {
    color: var(--color-text-main);
  }
}

.badge {
  background-color: var(--color-danger);
  color: white;
  font-size: 0.75rem;
  padding: 0 6px;
  border-radius: $radius-full;
}

.update-badge {
  background-color: var(--color-primary);
  margin-left: 0.5rem;
}

.new-badge {
  background-color: #10b981;
  margin-left: 0.5rem;
}

.arrow {
  color: var(--color-primary);
  font-weight: 400;
  font-size: 1rem;
}

.audit-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-lg;
}

.audit-item {
  display: flex;
  gap: $spacing-lg;
  background-color: white;
  padding: $spacing-lg;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  
  @media (max-width: 768px) {
    flex-direction: column;
  }
}

.item-preview {
  width: 200px;
  flex-shrink: 0;
  background-color: #f3f4f6;
  border-radius: $radius-md;
  overflow: hidden;
  
  &.profile-preview {
    background: none;
    width: auto;
    min-width: 150px;
  }
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
}

.old-new-avatar {
  display: flex;
  gap: $spacing-md;
  
  .avatar-box {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    font-size: 0.75rem;
    color: var(--color-text-muted);
    
    .avatar-mini {
      width: 64px;
      height: 64px;
      border-radius: 50%;
      object-fit: cover;
    }
  }
}

.item-info {
  flex: 1;
  
  h3 {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: $spacing-xs;
    color: var(--color-text-main);
    
    .arrow {
      color: var(--color-primary);
      font-weight: 400;
    }
  }
}

.change-block {
  margin-top: $spacing-md;
  font-size: 0.875rem;
  
  .label {
    font-weight: 600;
    color: var(--color-text-muted);
    margin-bottom: 4px;
  }
  
  .old-val {
    text-decoration: line-through;
    color: var(--color-text-muted);
    margin-bottom: 2px;
  }
  
  .new-val {
    color: var(--color-success, #10b981);
    background-color: #ecfdf5;
    padding: 4px 8px;
    border-radius: $radius-sm;
    display: inline-block;
  }
}

.item-meta {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin-bottom: $spacing-sm;
  display: flex;
  gap: $spacing-sm;
}

.item-desc {
  color: var(--color-text-main);
  margin-bottom: $spacing-md;
  font-size: 0.9375rem;
}

.item-tags {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-xs;
  margin-bottom: $spacing-sm;
}

.tag {
  background-color: #f3f4f6;
  padding: 2px 8px;
  border-radius: $radius-full;
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.item-actions {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
  justify-content: center;
  
  @media (max-width: 768px) {
    flex-direction: row;
    justify-content: flex-end;
  }
}

.btn--success {
  background-color: var(--color-secondary);
  color: white;
  border: none;
  
  &:hover {
    opacity: 0.9;
  }
}

.loading, .empty-state {
  text-align: center;
  padding: $spacing-2xl;
  color: var(--color-text-muted);
}

.empty-icon {
  color: var(--color-secondary);
  margin-bottom: $spacing-md;
}
</style>
