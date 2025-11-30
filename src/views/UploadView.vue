<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { Upload, X } from 'lucide-vue-next'
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
  <div class="container">
    <div class="upload-header">
      <h1>上传新模型</h1>
      <p class="subtitle">与社区分享你的 YSM 创作</p>
    </div>

    <div class="card upload-form">
      <form @submit.prevent="handleSubmit">
        <!-- File Upload Section -->
        <div class="form-section">
          <h2>模型文件</h2>
          
          <div class="form-group">
            <label>模型文件 (.ysm, .zip)</label>
            <div class="file-input-wrapper">
              <input 
                type="file" 
                accept=".ysm,.zip,.rar,.7z" 
                @change="handleModelSelect" 
                class="file-input"
                required
              >
              <div class="file-placeholder" :class="{ 'has-file': modelFile }">
                <Upload :size="24" />
                <span>{{ modelFile ? modelFile.name : '点击选择或拖拽文件到此处' }}</span>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label>预览图片</label>
            <div class="file-input-wrapper">
              <input 
                type="file" 
                accept="image/*" 
                @change="handleImageSelect" 
                class="file-input"
              >
              <div class="file-placeholder" :class="{ 'has-file': imageFile }">
                <img v-if="imagePreview" :src="imagePreview" class="image-preview" />
                <div v-else class="placeholder-content">
                  <Upload :size="24" />
                  <span>选择预览图片</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Details Section -->
        <div class="form-section">
          <h2>模型详情</h2>
          
          <div class="form-group">
            <label for="title">标题</label>
            <input 
              id="title" 
              v-model="title" 
              type="text" 
              class="input" 
              required 
              placeholder="例如：酷炫的 Steve 变体"
            >
          </div>

          <div class="form-group">
            <label for="description">描述 (支持 Markdown)</label>
            <textarea 
              id="description" 
              v-model="description" 
              class="input textarea" 
              rows="8"
              placeholder="描述你的模型... 支持 Markdown 语法"
            ></textarea>
          </div>

          <div class="form-group">
            <label>标签</label>
            <div class="tags-input">
              <div v-for="tag in tags" :key="tag" class="tag">
                {{ tag }}
                <button type="button" @click="removeTag(tag)"><X :size="14" /></button>
              </div>
              <input 
                v-model="tagInput" 
                @keydown.enter.prevent="addTag" 
                type="text" 
                class="input-ghost" 
                placeholder="添加标签 (按回车)"
              >
            </div>
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="isPublic">
              <span>公开此模型</span>
            </label>
          </div>
        </div>

        <div v-if="errorMsg" class="error-message">
          {{ errorMsg }}
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--secondary" @click="router.back()">取消</button>
          <button type="submit" class="btn btn--primary" :disabled="loading">
            {{ loading ? '上传中...' : '发布模型' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.upload-header {
  margin-bottom: $spacing-xl;
  
  h1 {
    font-size: 2rem;
    font-weight: 700;
    color: var(--color-text-main);
  }

  .subtitle {
    color: var(--color-text-muted);
  }
}

.upload-form {
  padding: $spacing-xl;
}

.form-section {
  margin-bottom: $spacing-xl;
  padding-bottom: $spacing-xl;
  border-bottom: 1px solid var(--color-border);

  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
    margin-bottom: 0;
  }

  h2 {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: $spacing-md;
    color: var(--color-text-main);
  }
}

.form-group {
  margin-bottom: $spacing-md;

  label {
    display: block;
    font-weight: 500;
    margin-bottom: $spacing-sm;
    color: var(--color-text-main);
  }
}

.file-input-wrapper {
  position: relative;
  height: 120px;
  border: 2px dashed var(--color-border);
  border-radius: $radius-md;
  overflow: hidden;
  transition: $transition-base;

  &:hover {
    border-color: var(--color-primary);
    background-color: #f9fafb;
  }
}

.file-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
  z-index: 10;
}

.file-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-muted);
  gap: $spacing-sm;

  &.has-file {
    color: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.image-preview {
  height: 100%;
  width: 100%;
  object-fit: cover;
}

.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-sm;
}

.textarea {
  resize: vertical;
}

.tags-input {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
  padding: $spacing-sm;
  border: 1px solid var(--color-border);
  border-radius: $radius-md;
  background-color: white;

  &:focus-within {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px rgba($color-primary, 0.2);
  }
}

.tag {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  background-color: #e0e7ff;
  color: var(--color-primary);
  padding: 2px 8px;
  border-radius: $radius-full;
  font-size: 0.875rem;

  button {
    background: none;
    border: none;
    cursor: pointer;
    color: inherit;
    display: flex;
  }
}

.input-ghost {
  border: none;
  outline: none;
  flex: 1;
  min-width: 120px;
  font-size: 0.875rem;
}

.checkbox-group {
  margin-top: $spacing-md;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  cursor: pointer;
  user-select: none;

  input {
    width: 1.25rem;
    height: 1.25rem;
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: $spacing-md;
  margin-top: $spacing-xl;
}

.error-message {
  color: var(--color-danger);
  margin-top: $spacing-md;
  text-align: center;
}
</style>
