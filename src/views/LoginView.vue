<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { useHead } from '@vueuse/head'

const router = useRouter()
const userStore = useUserStore()

useHead({
  title: '登录 / 注册 - YSM 模型站',
  meta: [
    { name: 'description', content: '登录 YSM 模型站，上传和管理你的 Minecraft 模型。' }
  ]
})

const isLogin = ref(true)
const email = ref('')
const password = ref('')
const username = ref('')
const loading = ref(false)
const errorMsg = ref('')

const handleSubmit = async () => {
  loading.value = true
  errorMsg.value = ''

  try {
    if (isLogin.value) {
      const { error } = await supabase.auth.signInWithPassword({
        email: email.value,
        password: password.value
      })
      if (error) throw new Error('登录失败，请检查邮箱和密码')
    } else {
      // Registration: Email & Password only
      const { error } = await supabase.auth.signUp({
        email: email.value,
        password: password.value
      })
      if (error) throw new Error('注册失败，请稍后重试')
      
      // No manual profile creation here, rely on trigger to set random username
    }
    
    await userStore.fetchUser()
    router.push('/')
  } catch (e: any) {
    errorMsg.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container auth-container">
    <div class="card auth-card">
      <h1 class="auth-title">{{ isLogin ? '欢迎回来' : '创建账户' }}</h1>
      
      <form @submit.prevent="handleSubmit" class="auth-form">
        <!-- Removed Username Input -->

        <div class="form-group">
          <label for="email">邮箱</label>
          <input 
            id="email" 
            v-model="email" 
            type="email" 
            class="input" 
            required 
            placeholder="请输入邮箱地址"
          >
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input 
            id="password" 
            v-model="password" 
            type="password" 
            class="input" 
            required 
            minlength="6"
            placeholder="请输入密码"
          >
        </div>

        <div v-if="errorMsg" class="error-message">
          {{ errorMsg }}
        </div>

        <button type="submit" class="btn btn--primary" :disabled="loading">
          {{ loading ? '处理中...' : (isLogin ? '登录' : '注册') }}
        </button>
      </form>

      <div class="auth-footer">
        <p>
          {{ isLogin ? "还没有账户？" : "已有账户？" }}
          <a href="#" @click.prevent="isLogin = !isLogin">
            {{ isLogin ? '立即注册' : '立即登录' }}
          </a>
        </p>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: $spacing-xl;
}

.auth-title {
  text-align: center;
  margin-bottom: $spacing-xl;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-main);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: $spacing-lg;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;

  label {
    font-weight: 500;
    font-size: 0.875rem;
    color: var(--color-text-main);
  }
}

.error-message {
  color: var(--color-danger);
  font-size: 0.875rem;
  text-align: center;
}

.auth-footer {
  margin-top: $spacing-lg;
  text-align: center;
  font-size: 0.875rem;
  color: var(--color-text-muted);

  a {
    color: var(--color-primary);
    font-weight: 500;
    
    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
