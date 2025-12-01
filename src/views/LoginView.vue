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
// const username = ref('')
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
      if (error) throw error
    } else {
      // Registration: Email & Password only
      const { error } = await supabase.auth.signUp({
        email: email.value,
        password: password.value
      })
      if (error) throw error
      
      // No manual profile creation here, rely on trigger to set random username
      if (!error) {
        alert('注册成功！请前往您的邮箱确认验证邮件，然后登录。')
        isLogin.value = true
        return
      }
    }
    
    await userStore.fetchUser()
    router.push('/')
  } catch (e: any) {
    console.error('Auth error:', e)
    // Translate common Supabase errors to Chinese
    if (e.message.includes('Invalid login credentials')) {
      errorMsg.value = '邮箱或密码错误'
    } else if (e.message.includes('Email not confirmed')) {
      errorMsg.value = '请先前往邮箱确认您的注册邮件'
    } else if (e.message.includes('User already registered')) {
      errorMsg.value = '该邮箱已被注册'
    } else if (e.message.includes('Password should be at least')) {
      errorMsg.value = '密码长度不能少于6位'
    } else if (e.message.includes('Error sending confirmation email')) {
      errorMsg.value = '发送验证邮件失败。请联系管理员检查邮件服务配置。'
    } else {
      errorMsg.value = e.message || '操作失败，请稍后重试'
    }
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
          <div class="label-row">
            <label for="password">密码</label>
            <router-link v-if="isLogin" to="/reset-password" class="forgot-link">忘记密码？</router-link>
          </div>
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

  .label-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  label {
    font-weight: 500;
    font-size: 0.875rem;
    color: var(--color-text-main);
  }

  .forgot-link {
    font-size: 0.875rem;
    color: var(--color-primary);
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
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
