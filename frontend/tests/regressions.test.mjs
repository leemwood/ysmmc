import assert from 'node:assert/strict'
import { test } from 'node:test'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import { createPinia, setActivePinia } from 'pinia'
import { createRenderer, nextTick } from 'vue'
import { createMemoryHistory, createRouter } from 'vue-router'
import { AxiosError } from 'axios'
import { tiptapToHtml } from '../src/utils/tiptap.ts'
import { useNexusmcStore } from '../src/stores/nexusmc.ts'
import { isNexusmcAuthError, isNexusmcConfigError, nexusmcApi, nexusmcMeApi } from '../src/lib/nexusmc.ts'
import CallbackView from '../src/views/NexusmcCallbackView.vue'
import MeView from '../src/views/MeView.vue'
import ModelsView from '../src/views/NexusmcModelsView.vue'

const storage = new Map()
globalThis.localStorage = {
  getItem: (key) => storage.get(key) ?? null,
  setItem: (key, value) => storage.set(key, value),
  removeItem: (key) => storage.delete(key),
}

// Vue 的内存渲染器运行真实组件及生命周期，不需要浏览器 DOM。
function element(tag = '') {
  return {
    tag, children: [], text: '', parent: null,
    addEventListener() {}, removeEventListener() {},
    get options() { return this.children.filter((node) => node.tag === 'option') },
  }
}
const renderer = createRenderer({
  createElement: element,
  createText: (text) => ({ ...element(), text }),
  createComment: () => element(),
  setText: (node, text) => { node.text = text },
  setElementText: (node, text) => { node.text = text; node.children = [] },
  patchProp(node, key, _previous, value) { node[key] = value },
  insert(node, parent, anchor) {
    if (node.parent) node.parent.children.splice(node.parent.children.indexOf(node), 1)
    const index = anchor ? parent.children.indexOf(anchor) : -1
    parent.children.splice(index < 0 ? parent.children.length : index, 0, node)
    node.parent = parent
  },
  remove(node) {
    if (node.parent) node.parent.children.splice(node.parent.children.indexOf(node), 1)
    node.parent = null
  },
  parentNode: (node) => node.parent,
  nextSibling: (node) => node.parent?.children[node.parent.children.indexOf(node) + 1] ?? null,
})

function apiError(status, code) {
  return new AxiosError('NexusMC request failed', undefined, undefined, undefined, { status, data: { code } })
}

function texts(node) {
  return [node.text, ...node.children.flatMap(texts)]
}

function findNode(node, predicate) {
  if (predicate(node)) return node
  for (const child of node.children) {
    const found = findNode(child, predicate)
    if (found) return found
  }
}

test('站点密钥错误与用户登录过期分开处理', () => {
  for (const error of [apiError(401, 'SITE_API_KEY_INVALID'), apiError(403, 'SITE_API_KEY_SCOPE_REQUIRED'), apiError(503, 'not_configured')]) {
    assert.equal(isNexusmcConfigError(error), true)
    assert.equal(isNexusmcAuthError(error), false)
  }
  for (const error of [apiError(401, 'OAUTH_TOKEN_INVALID'), apiError(403, 'OAUTH_SCOPE_REQUIRED')]) {
    assert.equal(isNexusmcConfigError(error), false)
    assert.equal(isNexusmcAuthError(error), true)
  }
  assert.equal(isNexusmcConfigError(new Error('network unavailable')), false)
  assert.equal(isNexusmcConfigError(apiError(403, 'RESOURCE_FORBIDDEN')), false)
})

test('资源页在密钥失效时提示服务不可用，刷新成功后恢复列表和分类', async () => {
  const catalog = nexusmcApi.catalog
  const resources = nexusmcApi.resources
  let available = false
  nexusmcApi.catalog = async () => {
    if (!available) throw apiError(401, 'SITE_API_KEY_INVALID')
    return { data: ['ysm'] }
  }
  nexusmcApi.resources = async () => {
    if (!available) throw apiError(401, 'SITE_API_KEY_INVALID')
    return { data: { items: [], pagination: { total: 0, totalPages: 1 } } }
  }
  const router = createRouter({ history: createMemoryHistory(), routes: [{ path: '/', component: ModelsView }] })
  await router.push('/')
  const root = element()
  const app = renderer.createApp(ModelsView)
  try {
    app.use(router).mount(root)
    await new Promise((resolve) => setImmediate(resolve))
    await nextTick()
    assert.ok(texts(root).join('').includes('资源服务暂不可用'))
    assert.ok(!texts(root).join('').includes('暂无 YSM 模型资源'))
    available = true
    const refresh = findNode(root, (node) => node['aria-label'] === '刷新')
    assert.ok(refresh, 'refresh action must remain available')
    refresh.onClick()
    await new Promise((resolve) => setImmediate(resolve))
    await nextTick()
    assert.ok(!texts(root).join('').includes('资源服务暂不可用'))
    assert.ok(texts(root).join('').includes('暂无 YSM 模型资源'))
    assert.ok(findNode(root, (node) => node.tag === 'option' && node.text === 'ysm'))
  } finally {
    app.unmount()
    nexusmcApi.catalog = catalog
    nexusmcApi.resources = resources
  }
})

test('OAuth 回调更新已初始化的 store，并返回来源页', async () => {
  storage.clear()
  const pinia = createPinia()
  setActivePinia(pinia)
  const store = useNexusmcStore()
  assert.equal(store.isLoggedIn, false)
  globalThis.window = { sessionStorage: {
    getItem: () => '/me', removeItem() {},
  } }
  const profile = { sub: 'user-1', username: 'Review' }
  const router = createRouter({ history: createMemoryHistory(), routes: [
    { path: '/nexusmc/callback', component: CallbackView },
    { path: '/me', component: { render: () => null } },
  ] })
  await router.push({ path: '/nexusmc/callback', query: { user: JSON.stringify(profile) } })
  const app = renderer.createApp(CallbackView)
  app.use(pinia).use(router)
  app.mount(element())
  await new Promise((resolve) => setImmediate(resolve))
  assert.equal(store.isLoggedIn, true)
  assert.equal(store.user.username, 'Review')
  assert.deepEqual(JSON.parse(storage.get('nexusmc_user')), profile)
  assert.equal(router.currentRoute.value.path, '/me')
  app.unmount()
  delete globalThis.window
})

test('个人主页使用上游 count 显示未读角标', async () => {
  storage.clear()
  const pinia = createPinia()
  setActivePinia(pinia)
  useNexusmcStore().completeLogin({ sub: 'user-2', username: 'Reader' })
  const content = nexusmcMeApi.content
  const unreadCount = nexusmcMeApi.unreadCount
  nexusmcMeApi.content = async () => ({ data: { items: [] } })
  nexusmcMeApi.unreadCount = async () => ({ data: { count: 3 } })
  const root = element()
  const router = createRouter({ history: createMemoryHistory(), routes: [{ path: '/me', component: MeView }] })
  await router.push('/me')
  const app = renderer.createApp(MeView)
  try {
    app.use(pinia).use(router).mount(root)
    await new Promise((resolve) => setImmediate(resolve))
    await nextTick()
    const texts = (node) => [node.text, ...node.children.flatMap(texts)]
    assert.ok(texts(root).some((text) => text.trim() === '3'))
  } finally {
    app.unmount()
    nexusmcMeApi.content = content
    nexusmcMeApi.unreadCount = unreadCount
  }
})

test('表格、任务列表和未知容器保留内容且不会递归自身', () => {
  const paragraph = { type: 'paragraph', content: [{ type: 'text', text: '<script>alert(1)</script>' }] }
  for (const node of [
    { type: 'table', content: [{ type: 'tableRow', content: [{ type: 'tableCell', content: [paragraph] }] }] },
    { type: 'taskList', content: [{ type: 'taskItem', content: [paragraph] }] },
    { type: 'futureNode', content: [paragraph] },
  ]) {
    const html = tiptapToHtml({ type: 'doc', content: [node] })
    assert.ok(html.includes('&lt;script&gt;'))
    assert.ok(!html.includes('<script>'))
  }
  assert.ok(tiptapToHtml({ type: 'doc', content: [{ type: 'image', attrs: { src: 'javascript:alert(1)' } }] }) === '')
})

test('Service Worker 升级清理旧的个人数据缓存', async () => {
  let listener
  const deleted = []
  vm.runInNewContext(readFileSync(new URL('../public/clear-legacy-api-cache.js', import.meta.url), 'utf8'), {
    self: { addEventListener: (type, callback) => { assert.equal(type, 'activate'); listener = callback } },
    caches: { delete: async (name) => { deleted.push(name); return true } },
  })
  let work
  listener({ waitUntil: (promise) => { work = promise } })
  await work
  assert.deepEqual(deleted, ['api-cache'])
})
