import assert from 'node:assert/strict'
import { test } from 'node:test'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'

test('构建产物对个人 API 使用 NetworkOnly 并加载旧缓存清理脚本', () => {
  const routes = []
  const imported = []
  class NetworkOnly {}
  class NavigationRoute {}
  class CacheFirst {}
  class ExpirationPlugin {}
  class CacheableResponsePlugin {}
  const workbox = {
    NetworkOnly, NavigationRoute, CacheFirst, ExpirationPlugin, CacheableResponsePlugin,
    clientsClaim() {}, precacheAndRoute() {}, cleanupOutdatedCaches() {},
    createHandlerBoundToURL() { return () => {} },
    registerRoute(...args) { routes.push(args) },
  }
  const context = {
    self: { define() {}, skipWaiting() {} },
    define: (_dependencies, factory) => factory(workbox),
    importScripts: (path) => imported.push(path),
  }
  vm.runInNewContext(readFileSync(new URL('../dist/sw.js', import.meta.url), 'utf8'), context)
  for (const path of ['/me/content', '/me/favorites', '/me/notifications', '/me/notifications/unread-count', '/auth/start']) {
    const route = routes.find(([matcher]) => typeof matcher?.test === 'function' && matcher.test('https://ysmmc.example/api/nexusmc' + path))
    assert.ok(route, `missing API route: ${path}`)
    assert.ok(route[1] instanceof NetworkOnly, `API is cacheable: ${path}`)
  }
  assert.ok(imported.includes('/clear-legacy-api-cache.js'))
})
