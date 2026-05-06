import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/system/login.vue'),
      meta: { title: '登录', public: true }
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index.vue'),
      meta: { title: '综合态势' }
    },
    {
      path: '/business',
      name: 'Business',
      component: () => import('@/views/business/index.vue'),
      meta: { title: '商户管理' }
    },
    {
      path: '/business/:id',
      name: 'BusinessDetail',
      component: () => import('@/views/business/detail.vue'),
      meta: { title: '商户台账' }
    },
    {
      path: '/report',
      name: 'Report',
      component: () => import('@/views/report/index.vue'),
      meta: { title: '申报记录' }
    },
    {
      path: '/alert',
      name: 'Alert',
      component: () => import('@/views/alert/index.vue'),
      meta: { title: '预警专栏' }
    },
    {
      path: '/statistics',
      name: 'Statistics',
      component: () => import('@/views/statistics/index.vue'),
      meta: { title: '统计分析' }
    },
    {
      path: '/pos-interface',
      name: 'POSInterface',
      component: () => import('@/views/pos-interface/index.vue'),
      meta: { title: '接口管理' }
    },
    {
      path: '/system',
      name: 'System',
      component: () => import('@/views/system/index.vue'),
      meta: { title: '系统设置' }
    },
    { path: '/', redirect: '/dashboard' }
  ]
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('admin_token')
  if (!to.meta.public && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
