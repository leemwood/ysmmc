import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

// 本站定位：NexusMC 资源展示与引流站。登录注册、用户体系均由 NexusMC 承载，
// 本站不设本地账号，因此无需路由守卫。
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/NexusmcModelsView.vue'),
  },
  {
    path: '/nexusmc/resource/:id',
    name: 'nexusmc-resource',
    component: () => import('@/views/NexusmcDetailView.vue'),
  },
  {
    path: '/nexusmc/callback',
    name: 'nexusmc-callback',
    component: () => import('@/views/NexusmcCallbackView.vue'),
  },
  {
    path: '/me',
    name: 'me',
    component: () => import('@/views/MeView.vue'),
  },
  {
    // 兼容旧地址
    path: '/nexusmc',
    redirect: '/',
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
