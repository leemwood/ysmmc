<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { supabase } from '../supabase/client'
import { Check, X, User, Box, Search, Trash2, Shield, ShieldAlert, Edit, Users, FolderTree } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model, Profile } from '../types'

useHead({
  title: '管理员控制台 - YSM 模型站',
  meta: [
    { name: 'robots', content: 'noindex, nofollow' }
  ]
})

const activeTab = ref<'audit_models' | 'audit_profiles' | 'manage_models' | 'manage_users'>('audit_models')

const pendingModels = ref<Model[]>([])
const pendingProfiles = ref<Profile[]>([])
const allModels = ref<Model[]>([])
const allProfiles = ref<Profile[]>([])
const loading = ref(true)

// Search filters
const modelSearch = ref('')
const userSearch = ref('')

const filteredModels = computed(() => {
  if (!modelSearch.value) return allModels.value
  const query = modelSearch.value.toLowerCase()
  return allModels.value.filter(m => 
    m.title.toLowerCase().includes(query) || 
    m.profiles?.username?.toLowerCase().includes(query)
  )
})

const filteredProfiles = computed(() => {
  if (!userSearch.value) return allProfiles.value
  const query = userSearch.value.toLowerCase()
  return allProfiles.value.filter(p => 
    p.username.toLowerCase().includes(query) || 
    p.id.toLowerCase().includes(query)
  )
})

const fetchPendingModels = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*, profiles(username)')
    .or('status.eq.pending,update_status.eq.pending_review')
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

const fetchAllModels = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*, profiles(username)')
    .order('created_at', { ascending: false })

  if (error) console.error(error)
  else allModels.value = data as any
  loading.value = false
}

const fetchAllProfiles = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('profiles')
    .select('*')
    .order('created_at', { ascending: false })

  if (error) console.error(error)
  else allProfiles.value = data as any
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

const handleDeleteModel = async (id: string) => {
  if (!confirm('确定要永久删除此模型吗？此操作不可撤销。')) return

  const { error } = await supabase
    .from('models')
    .delete()
    .eq('id', id)

  if (error) {
    alert('删除失败: ' + error.message)
  } else {
    allModels.value = allModels.value.filter(m => m.id !== id)
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
      profile_status: 'approved',
      pending_changes: null 
    })
    .eq('id', id)

  if (!error) {
    pendingProfiles.value = pendingProfiles.value.filter(p => p.id !== id)
  }
}

const handleUpdateUserRole = async (id: string, newRole: 'user' | 'admin') => {
  const { error } = await supabase
    .from('profiles')
    .update({ role: newRole })
    .eq('id', id)

  if (error) {
    alert('角色更新失败: ' + error.message)
  } else {
    const user = allProfiles.value.find(p => p.id === id)
    if (user) user.role = newRole
  }
}

const handleDeleteUser = async (id: string) => {
  if (!confirm('确定要删除此用户资料吗？(注意：这不会删除其 Auth 账号，仅删除资料表记录)')) return

  const { error } = await supabase
    .from('profiles')
    .delete()
    .eq('id', id)

  if (error) {
    alert('删除失败: ' + error.message)
  } else {
    allProfiles.value = allProfiles.value.filter(p => p.id !== id)
  }
}

const switchTab = (tab: typeof activeTab.value) => {
  activeTab.value = tab
  if (tab === 'audit_models') fetchPendingModels()
  else if (tab === 'audit_profiles') fetchPendingProfiles()
  else if (tab === 'manage_models') fetchAllModels()
  else if (tab === 'manage_users') fetchAllProfiles()
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
        @click="switchTab('audit_models')" 
        :class="['tab-btn', activeTab === 'audit_models' ? 'active' : '']"
      >
        <FolderTree :size="18" /> 模型审核
        <span v-if="pendingModels.length" class="badge">{{ pendingModels.length }}</span>
      </button>
      <button 
        @click="switchTab('audit_profiles')" 
        :class="['tab-btn', activeTab === 'audit_profiles' ? 'active' : '']"
      >
        <User :size="18" /> 资料审核
        <span v-if="pendingProfiles.length" class="badge">{{ pendingProfiles.length }}</span>
      </button>
      <button 
        @click="switchTab('manage_models')" 
        :class="['tab-btn', activeTab === 'manage_models' ? 'active' : '']"
      >
        <Box :size="18" /> 模型管理
      </button>
      <button 
        @click="switchTab('manage_users')" 
        :class="['tab-btn', activeTab === 'manage_users' ? 'active' : '']"
      >
        <Users :size="18" /> 用户管理
      </button>
    </div>

    <div v-if="loading" class="loading">
      正在加载...
    </div>

    <!-- Models Audit List -->
    <div v-else-if="activeTab === 'audit_models'">
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

    <!-- Profiles Audit List -->
    <div v-else-if="activeTab === 'audit_profiles'">
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

    <!-- Model Management List -->
    <div v-else-if="activeTab === 'manage_models'">
      <div class="management-toolbar">
        <div class="search-box">
          <Search :size="18" />
          <input v-model="modelSearch" type="text" placeholder="搜索模型标题或作者...">
        </div>
      </div>

      <div class="management-table-wrapper">
        <table class="management-table">
          <thead>
            <tr>
              <th>预览</th>
              <th>标题</th>
              <th>作者</th>
              <th>状态</th>
              <th>下载</th>
              <th>上传日期</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="model in filteredModels" :key="model.id">
              <td>
                <img :src="model.image_url || 'https://via.placeholder.com/40x30'" class="table-preview">
              </td>
              <td class="td-title">{{ model.title }}</td>
              <td>{{ model.profiles?.username || '未知' }}</td>
              <td>
                <span :class="['status-badge', model.status]">
                  {{ model.status === 'approved' ? '已通过' : model.status === 'pending' ? '待审核' : '已拒绝' }}
                </span>
              </td>
              <td>{{ model.downloads }}</td>
              <td>{{ new Date(model.created_at).toLocaleDateString() }}</td>
              <td>
                <div class="table-actions">
                  <router-link :to="`/model/${model.id}/edit`" class="action-btn" title="编辑">
                    <Edit :size="18" />
                  </router-link>
                  <button @click="handleDeleteModel(model.id)" class="action-btn danger" title="删除">
                    <Trash2 :size="18" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- User Management List -->
    <div v-else-if="activeTab === 'manage_users'">
      <div class="management-toolbar">
        <div class="search-box">
          <Search :size="18" />
          <input v-model="userSearch" type="text" placeholder="搜索用户名或 ID...">
        </div>
      </div>

      <div class="management-table-wrapper">
        <table class="management-table">
          <thead>
            <tr>
              <th>头像</th>
              <th>用户名</th>
              <th>角色</th>
              <th>状态</th>
              <th>注册日期</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="profile in filteredProfiles" :key="profile.id">
              <td>
                <img :src="profile.avatar_url || 'https://via.placeholder.com/32'" class="table-avatar">
              </td>
              <td>
                <div class="user-info-cell">
                  <span class="username">{{ profile.username }}</span>
                  <span class="user-id">{{ profile.id }}</span>
                </div>
              </td>
              <td>
                <span :class="['role-badge', profile.role]">
                  {{ profile.role === 'admin' ? '管理员' : '普通用户' }}
                </span>
              </td>
              <td>
                <span :class="['status-badge', profile.profile_status]">
                  {{ profile.profile_status === 'approved' ? '正常' : '待审核' }}
                </span>
              </td>
              <td>{{ new Date(profile.created_at).toLocaleDateString() }}</td>
              <td>
                <div class="table-actions">
                  <button 
                    v-if="profile.role === 'user'" 
                    @click="handleUpdateUserRole(profile.id, 'admin')" 
                    class="action-btn success" 
                    title="设为管理员"
                  >
                    <Shield :size="18" />
                  </button>
                  <button 
                    v-else 
                    @click="handleUpdateUserRole(profile.id, 'user')" 
                    class="action-btn warning" 
                    title="取消管理员"
                  >
                    <ShieldAlert :size="18" />
                  </button>
                  <button @click="handleDeleteUser(profile.id)" class="action-btn danger" title="删除资料">
                    <Trash2 :size="18" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
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

// Management styles
.management-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: $spacing-lg;
}

.search-box {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  background-color: white;
  border: 1px solid var(--color-border);
  border-radius: $radius-md;
  padding: 0 $spacing-md;
  width: 300px;
  
  &:focus-within {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px rgba($color-primary, 0.1);
  }
  
  input {
    border: none;
    outline: none;
    padding: $spacing-sm 0;
    width: 100%;
    font-size: 0.875rem;
  }
  
  color: var(--color-text-muted);
}

.management-table-wrapper {
  background-color: white;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  overflow: hidden;
  overflow-x: auto;
}

.management-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
  
  th {
    background-color: #f9fafb;
    padding: $spacing-md;
    font-weight: 600;
    color: var(--color-text-muted);
    border-bottom: 1px solid var(--color-border);
  }
  
  td {
    padding: $spacing-md;
    border-bottom: 1px solid var(--color-border);
    vertical-align: middle;
  }
  
  tr:last-child td {
    border-bottom: none;
  }
  
  .table-preview {
    width: 40px;
    height: 30px;
    object-fit: cover;
    border-radius: $radius-sm;
  }
  
  .table-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
  }
  
  .td-title {
    font-weight: 500;
    color: var(--color-text-main);
  }
}

.user-info-cell {
  display: flex;
  flex-direction: column;
  
  .username {
    font-weight: 500;
    color: var(--color-text-main);
  }
  
  .user-id {
    font-size: 0.75rem;
    color: var(--color-text-muted);
  }
}

.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: $radius-full;
  font-size: 0.75rem;
  font-weight: 500;
  
  &.approved {
    background-color: #ecfdf5;
    color: #065f46;
  }
  
  &.pending, &.pending_review {
    background-color: #fffbeb;
    color: #92400e;
  }
  
  &.rejected {
    background-color: #fef2f2;
    color: #991b1b;
  }
}

.role-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: $radius-full;
  font-size: 0.75rem;
  font-weight: 500;
  
  &.admin {
    background-color: #eef2ff;
    color: #3730a3;
  }
  
  &.user {
    background-color: #f3f4f6;
    color: #374151;
  }
}

.table-actions {
  display: flex;
  gap: $spacing-xs;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: $radius-md;
  border: 1px solid var(--color-border);
  background-color: white;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: $transition-base;
  text-decoration: none;
  
  &:hover {
    color: var(--color-primary);
    border-color: var(--color-primary);
    background-color: #f5f3ff;
  }
  
  &.danger:hover {
    color: var(--color-danger);
    border-color: var(--color-danger);
    background-color: #fef2f2;
  }
  
  &.success:hover {
    color: var(--color-secondary);
    border-color: var(--color-secondary);
    background-color: #ecfdf5;
  }
  
  &.warning:hover {
    color: #d97706;
    border-color: #d97706;
    background-color: #fffbeb;
  }
}
</style>
