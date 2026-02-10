<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { Upload, X, FileCode, ImageIcon, Tag as TagIcon, Check, Loader2, ArrowLeft, ShieldAlert } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'

const router = useRouter()
const userStore = useUserStore()

useHead({
  title: '上传模型 - YSM 模型站',
  meta: [
    { name: 'description', content: '上传并分享你的 Minecraft Yes Steve Model 作品。' }
  ]
})

const title = ref('')
const description = ref('')
const tags = ref<string[]>([])
const tagInput = ref('')
const isPublic = ref(true)
const loading = ref(false)
const errorMsg = ref('')

const modelFile = ref<File | null>(null)
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)

const handleImageSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    imageFile.value = file
    imagePreview.value = URL.createObjectURL(file)
  }
}

const handleModelSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    modelFile.value = target.files[0]
  }
}

const addTag = () => {
  if (tagInput.value.trim() && !tags.value.includes(tagInput.value.trim())) {
    tags.value.push(tagInput.value.trim())
    tagInput.value = ''
  }
}

const removeTag = (tag: string) => {
  tags.value = tags.value.filter(t => t !== tag)
}

const uploadFile = async (file: File, bucket: string) => {
  const fileExt = file.name.split('.').pop()
  const fileName = `${Math.random()}.${fileExt}`
  const filePath = `${userStore.user?.id}/${fileName}`

  const { error } = await supabase.storage.from(bucket).upload(filePath, file)
  if (error) throw error
  return filePath
}

const getPublicUrl = (bucket: string, path: string) => {
    const { data } = supabase.storage.from(bucket).getPublicUrl(path)
    return data.publicUrl
}

const handleSubmit = async () => {
  if (!modelFile.value) {
    errorMsg.value = '请选择模型文件'
    return
  }
  if (!userStore.user) return

  loading.value = true
  errorMsg.value = ''

  try {
    // 1. Upload Model
    const modelPath = await uploadFile(modelFile.value, 'models')
    
    // 2. Upload Image (optional)
    let imageUrl = null
    if (imageFile.value) {
        const imagePath = await uploadFile(imageFile.value, 'images')
        imageUrl = getPublicUrl('images', imagePath)
    }

    // 3. Create Database Record
    const { error } = await supabase.from('models').insert({
      user_id: userStore.user.id,
      title: title.value,
      description: description.value,
      file_path: modelPath,
      image_url: imageUrl,
      tags: tags.value,
      is_public: isPublic.value
    })

    if (error) throw error

    router.push('/')
  } catch (e: any) {
    console.error(e)
    errorMsg.value = e.message || '上传过程中发生错误'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container py-xl">
    <div class="upload-wrapper">
      <header class="upload-header">
        <button @click="router.back()" class="back-btn" aria-label="返回上一页">
          <ArrowLeft :size="20" aria-hidden="true" />
        </button>
        <div class="header-text">
          <h1>发布新模型</h1>
          <p>分享您的 Minecraft 创意作品，让更多人发现</p>
        </div>
      </header>

      <form @submit.prevent="handleSubmit" class="upload-form">
        <div class="form-grid">
          <!-- Left Column: Main Info -->
          <div class="form-main">
            <div class="form-card">
              <div class="form-group">
                <label for="title">模型标题 <span class="required" aria-hidden="true">*</span></label>
                <input 
                  id="title" 
                  v-model="title" 
                  type="text" 
                  class="input input--lg" 
                  required 
                  placeholder="给您的模型起一个响亮的名字"
                  aria-required="true"
                >
              </div>

              <div class="form-group">
                <label for="description">描述详情</label>
                <textarea 
                  id="description" 
                  v-model="description" 
                  class="input textarea" 
                  rows="10"
                  placeholder="详细介绍您的模型特点、使用方法等 (支持 Markdown)"
                ></textarea>
              </div>
            </div>

            <div class="form-card">
              <h3 class="card-title"><TagIcon :size="18" aria-hidden="true" /> 标签分类</h3>
              <div class="tag-input-wrapper">
                <input 
                  v-model="tagInput" 
                  type="text" 
                  class="input" 
                  placeholder="按回车添加标签"
                  @keydown.enter.prevent="addTag"
                  aria-label="输入标签并按回车添加"
                >
                <button type="button" @click="addTag" class="btn btn--secondary btn--sm">添加</button>
              </div>
              <div class="tags-display" role="list" aria-label="已添加的标签">
                <span v-for="tag in tags" :key="tag" class="tag-chip" role="listitem">
                  {{ tag }}
                  <button type="button" @click="removeTag(tag)" class="remove-tag" :aria-label="`删除标签 ${tag}`"><X :size="12" aria-hidden="true" /></button>
                </span>
                <p v-if="tags.length === 0" class="empty-tags">暂无标签</p>
              </div>
            </div>
          </div>

          <!-- Right Column: Files & Settings -->
          <aside class="form-side">
            <div class="form-card file-card">
              <h3 class="card-title"><FileCode :size="18" aria-hidden="true" /> 模型文件 <span class="required" aria-hidden="true">*</span></h3>
              <div class="file-upload-zone" :class="{ 'has-file': modelFile }">
                <input 
                  type="file" 
                  id="modelFile" 
                  class="file-input" 
                  @change="handleModelSelect"
                  accept=".ysm,.zip"
                  required
                  aria-required="true"
                >
                <label for="modelFile" class="file-label">
                  <template v-if="!modelFile">
                    <Upload :size="32" aria-hidden="true" />
                    <span>选择或拖拽模型文件</span>
                    <small>支持 .ysm, .zip 格式</small>
                  </template>
                  <template v-else>
                    <Check :size="32" class="success-icon" aria-hidden="true" />
                    <span class="file-name">{{ modelFile.name }}</span>
                    <button type="button" @click.stop.prevent="modelFile = null" class="btn btn--sm btn--secondary">更改文件</button>
                  </template>
                </label>
              </div>
            </div>

            <div class="form-card image-card">
              <h3 class="card-title"><ImageIcon :size="18" aria-hidden="true" /> 预览图片</h3>
              <div class="image-upload-zone" :class="{ 'has-preview': imagePreview }">
                <input 
                  type="file" 
                  id="imageFile" 
                  class="file-input" 
                  @change="handleImageSelect"
                  accept="image/*"
                >
                <label for="imageFile" class="image-label">
                  <img v-if="imagePreview" :src="imagePreview" :alt="title ? `${title} 的预览图` : '模型预览图'">
                  <template v-else>
                    <Upload :size="24" aria-hidden="true" />
                    <span>上传封面图</span>
                  </template>
                </label>
                <button v-if="imagePreview" type="button" @click.stop.prevent="imagePreview = null; imageFile = null" class="remove-image" aria-label="移除图片">
                  <X :size="16" aria-hidden="true" />
                </button>
              </div>
            </div>

            <div v-if="errorMsg" class="error-msg" role="alert" aria-live="assertive">
              <ShieldAlert :size="18" aria-hidden="true" /> {{ errorMsg }}
            </div>

            <button type="submit" class="btn btn--primary submit-btn" :disabled="loading" :aria-busy="loading">
              <template v-if="loading">
                <Loader2 class="animate-spin" :size="20" aria-hidden="true" /> 正在发布...
              </template>
              <template v-else>
                发布模型
              </template>
            </button>
          </aside>
        </div>
      </form>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.upload-wrapper {
  max-width: 1000px;
  margin: 0 auto;
}

.upload-header {
  display: flex;
  align-items: flex-start;
  gap: $spacing-lg;
  margin-bottom: $spacing-2xl;

  .back-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: $radius-full;
    border: 1px solid var(--color-border);
    background: white;
    color: var(--color-text-main);
    cursor: pointer;
    transition: $transition-base;
    margin-top: 4px;

    &:hover {
      background: var(--color-bg-alt);
      border-color: var(--color-text-muted);
      transform: translateX(-4px);
    }
  }

  .header-text {
    h1 {
      font-size: 2rem;
      font-weight: 800;
      color: var(--color-text-main);
      margin-bottom: $spacing-xs;
      letter-spacing: -0.025em;
    }

    p {
      color: var(--color-text-muted);
      font-size: 1.125rem;
    }
  }
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: $spacing-xl;
  align-items: flex-start;

  @media (max-width: 900px) {
    grid-template-columns: 1fr;
  }
}

.form-main, .form-side {
  display: flex;
  flex-direction: column;
  gap: $spacing-xl;
}

.form-card {
  background: white;
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  padding: $spacing-xl;
  box-shadow: $shadow-sm;

  .card-title {
    font-size: 1.125rem;
    font-weight: 700;
    color: var(--color-text-main);
    margin-bottom: $spacing-lg;
    display: flex;
    align-items: center;
    gap: $spacing-sm;
  }
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
  margin-bottom: $spacing-lg;

  &:last-child {
    margin-bottom: 0;
  }

  label {
    font-size: 0.9375rem;
    font-weight: 600;
    color: var(--color-text-main);

    .required {
      color: #ef4444;
      margin-left: 2px;
    }
  }

  .input--lg {
    font-size: 1.125rem;
    font-weight: 500;
  }

  .textarea {
    resize: vertical;
    line-height: 1.6;
  }
}

.tag-input-wrapper {
  display: flex;
  gap: $spacing-sm;
  margin-bottom: $spacing-md;
}

.tags-display {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;

  .tag-chip {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    padding: $spacing-xs $spacing-sm;
    background: var(--color-bg-alt);
    border: 1px solid var(--color-border);
    border-radius: $radius-full;
    font-size: 0.875rem;
    color: var(--color-text-main);

    .remove-tag {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      background: none;
      border: none;
      color: var(--color-text-muted);
      cursor: pointer;
      border-radius: $radius-full;

      &:hover {
        color: #ef4444;
        background: rgba(#ef4444, 0.1);
      }
    }
  }

  .empty-tags {
    font-size: 0.875rem;
    color: var(--color-text-muted);
    font-style: italic;
  }
}

.file-upload-zone {
  position: relative;
  border: 2px dashed var(--color-border);
  border-radius: $radius-lg;
  transition: $transition-base;
  
  &:hover {
    border-color: var(--color-primary);
    background: rgba($color-primary, 0.02);
  }

  &.has-file {
    border-style: solid;
    border-color: #10b981;
    background: rgba(#10b981, 0.02);
  }

  .file-input {
    position: absolute;
    inset: 0;
    opacity: 0;
    cursor: pointer;
    z-index: 1;
  }

  .file-label {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: $spacing-2xl $spacing-md;
    text-align: center;
    gap: $spacing-sm;
    color: var(--color-text-muted);

    span {
      font-weight: 600;
      color: var(--color-text-main);
    }

    small {
      font-size: 0.75rem;
    }

    .success-icon {
      color: #10b981;
    }

    .file-name {
      font-size: 0.875rem;
      word-break: break-all;
    }

    .btn {
      position: relative;
      z-index: 2;
    }
  }
}

.image-upload-zone {
  position: relative;
  border-radius: $radius-lg;
  overflow: hidden;
  background: var(--color-bg-alt);
  aspect-ratio: 16 / 9;
  border: 1px solid var(--color-border);

  .file-input {
    position: absolute;
    inset: 0;
    opacity: 0;
    cursor: pointer;
    z-index: 1;
  }

  .image-label {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: $spacing-xs;
    color: var(--color-text-muted);
    font-size: 0.875rem;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .remove-image {
    position: absolute;
    top: $spacing-sm;
    right: $spacing-sm;
    width: 32px;    height: 32px;
    border-radius: $radius-full;
    background: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    z-index: 2;
    backdrop-filter: blur(4px);
    transition: $transition-base;

    &:hover {
      background: rgba(0, 0, 0, 0.8);
      transform: scale(1.1);
    }
  }
}

.error-msg {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #ef4444;
  padding: $spacing-md;
  border-radius: $radius-lg;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.submit-btn {
  width: 100%;
  height: 3.5rem;
  font-size: 1.125rem;
  font-weight: 700;
  gap: $spacing-sm;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
