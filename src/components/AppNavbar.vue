<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Upload, User as UserIcon, LogOut, Menu, X, Shield } from 'lucide-vue-next'
import { ref } from 'vue'

const userStore = useUserStore()
const router = useRouter()
const isMenuOpen = ref(false)

const handleSignOut = async () => {
  await userStore.signOut()
  router.push('/login')
}
</script>

<template>
  <nav class="navbar">
    <div class="container navbar__container">
      <RouterLink to="/" class="navbar__logo">
        YSM 管理器
      </RouterLink>

      <div class="navbar__desktop">
        <template v-if="userStore.user">
          <RouterLink to="/upload" class="btn btn--primary">
            <Upload :size="18" />
            上传模型
          </RouterLink>
          <div class="navbar__user">
            <RouterLink v-if="userStore.isAdmin" to="/admin" class="navbar__link" title="管理员控制台">
              <Shield :size="20" />
            </RouterLink>
            <RouterLink to="/profile" class="navbar__link">
              <UserIcon :size="20" />
              个人中心
            </RouterLink>
            <button @click="handleSignOut" class="btn btn--secondary btn--sm">
              <LogOut :size="18" />
              退出登录
            </button>
          </div>
        </template>
        <template v-else>
          <RouterLink to="/login" class="btn btn--primary">
            登录 / 注册
          </RouterLink>
        </template>
      </div>

      <button class="navbar__toggle" @click="isMenuOpen = !isMenuOpen">
        <Menu v-if="!isMenuOpen" />
        <X v-else />
      </button>
    </div>

    <div v-if="isMenuOpen" class="navbar__mobile">
      <template v-if="userStore.user">
        <RouterLink to="/upload" class="navbar__mobile-link" @click="isMenuOpen = false">
          <Upload :size="18" />
          上传模型
        </RouterLink>
        <RouterLink v-if="userStore.isAdmin" to="/admin" class="navbar__mobile-link" @click="isMenuOpen = false">
          <Shield :size="18" />
          管理员控制台
        </RouterLink>
        <RouterLink to="/profile" class="navbar__mobile-link" @click="isMenuOpen = false">
          <UserIcon :size="18" />
          个人中心
        </RouterLink>
        <button @click="handleSignOut; isMenuOpen = false" class="navbar__mobile-link">
          <LogOut :size="18" />
          退出登录
        </button>
      </template>
      <template v-else>
        <RouterLink to="/login" class="navbar__mobile-link" @click="isMenuOpen = false">
          登录 / 注册
        </RouterLink>
      </template>
    </div>
  </nav>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.navbar {
  background-color: white;
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 50;

  &__container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 4rem;
  }

  &__logo {
    font-weight: 700;
    font-size: 1.25rem;
    color: var(--color-primary);
  }

  &__desktop {
    display: none;
    align-items: center;
    gap: $spacing-md;

    @media (min-width: 768px) {
      display: flex;
    }
  }

  &__user {
    display: flex;
    align-items: center;
    gap: $spacing-md;
  }

  &__link {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    color: var(--color-text-muted);
    font-weight: 500;
    
    &:hover {
      color: var(--color-text-main);
    }
  }

  &__toggle {
    background: none;
    border: none;
    cursor: pointer;
    color: var(--color-text-main);
    
    @media (min-width: 768px) {
      display: none;
    }
  }

  &__mobile {
    border-top: 1px solid var(--color-border);
    padding: $spacing-md;
    display: flex;
    flex-direction: column;
    gap: $spacing-sm;
    background-color: white;

    &-link {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      padding: $spacing-sm;
      color: var(--color-text-main);
      border-radius: $radius-md;

      &:hover {
        background-color: var(--color-background);
      }
    }
  }
}

.btn--sm {
  padding: $spacing-xs $spacing-sm;
  font-size: 0.875rem;
}
</style>
