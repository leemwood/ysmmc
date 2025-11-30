import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useUserStore } from '../stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/upload',
      name: 'upload',
      component: () => import('../views/UploadView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/model/:id',
      name: 'model-detail',
      component: () => import('../views/ModelDetailView.vue')
    },
    {
      path: '/model/:id/edit',
      name: 'model-edit',
      component: () => import('../views/EditModelView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/AdminDashboardView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/user/:id',
      name: 'user-profile',
      component: () => import('../views/UserPublicProfileView.vue')
    }
  ]
})

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()
  if (userStore.loading) {
    await userStore.fetchUser()
  }

  if (to.meta.requiresAuth && !userStore.user) {
    next({ name: 'login' })
  } else if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router
