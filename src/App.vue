<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import AppNavbar from './components/AppNavbar.vue'
import AppFooter from './components/AppFooter.vue'
import { useUserStore } from './stores/user'

const userStore = useUserStore()

onMounted(() => {
  userStore.fetchUser()
})
</script>

<template>
  <div class="app-layout">
    <AppNavbar />
    <main class="main-content">
      <RouterView />
    </main>
    <AppFooter />
  </div>
</template>

<style lang="scss" scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding-top: 2rem;
  padding-bottom: 2rem;
}
</style>
