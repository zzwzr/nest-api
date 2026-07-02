import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/install',
    },
    {
      path: '/install',
      name: 'install',
      component: () => import('@/views/install/InstallView.vue'),
      meta: { title: '系统安装' },
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} - ApiNest` : 'ApiNest'
  next()
})

export default router
