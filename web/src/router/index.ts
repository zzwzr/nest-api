import { createRouter, createWebHistory } from 'vue-router'
import { fetchInstallStatus } from '@/api/install'
import { fetchSiteInfo } from '@/api/auth'
import { getAccessToken } from '@/utils/auth-storage'
import { useAuth } from '@/composables/useAuth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/install',
      name: 'install',
      component: () => import('@/views/install/InstallView.vue'),
      meta: { title: '系统安装', public: true, installOnly: true },
    },
    {
      path: '/',
      component: () => import('@/layouts/AuthLayout.vue'),
      meta: { public: true },
      children: [
        {
          path: 'login',
          name: 'login',
          component: () => import('@/views/auth/LoginView.vue'),
          meta: { title: '登录', guestOnly: true },
        },
        {
          path: 'register',
          name: 'register',
          component: () => import('@/views/auth/RegisterView.vue'),
          meta: { title: '注册', guestOnly: true },
        },
      ],
    },
    {
      path: '/',
      component: () => import('@/layouts/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/home',
        },
        {
          path: 'home',
          name: 'home',
          component: () => import('@/views/home/HomeView.vue'),
          meta: { title: '首页' },
        },
        {
          path: 'admin/settings',
          name: 'admin-settings',
          component: () => import('@/views/admin/SettingsView.vue'),
          meta: { title: '系统设置', requiresAdmin: true },
        },
      ],
    },
  ],
})

let installChecked = false
let installed = false

router.beforeEach(async (to, _from, next) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} - ApiNest` : 'ApiNest'

  if (!installChecked) {
    try {
      const status = await fetchInstallStatus()
      installed = status.installed
    } catch {
      try {
        const site = await fetchSiteInfo()
        installed = site.installed
      } catch {
        installed = false
      }
    }
    installChecked = true
  }

  if (!installed) {
    if (to.path === '/install') {
      next()
      return
    }
    next('/install')
    return
  }

  if (to.path === '/install') {
    next('/login')
    return
  }

  const token = getAccessToken()
  const { bootstrap, user } = useAuth()
  if (token) {
    await bootstrap()
  }

  if (to.meta.guestOnly && token && user.value) {
    next('/home')
    return
  }

  if (to.meta.requiresAuth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresAdmin && !user.value?.is_admin) {
    next('/home')
    return
  }

  next()
})

export default router
