import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/',
      component: () => import('../views/Layout.vue'),
      redirect: '/dashboard',
      children: [
        { path: 'dashboard', component: () => import('../views/Dashboard.vue'), name: 'Dashboard' },
        { path: 'appointments', component: () => import('../views/Appointment.vue'), name: 'Appointments' },
        { path: 'providers', component: () => import('../views/Provider.vue'), name: 'Providers' },
        { path: 'challenges', component: () => import('../views/Challenge.vue'), name: 'Challenges' },
        { path: 'reports', component: () => import('../views/Report.vue'), name: 'Reports' }
      ]
    }
  ]
})

// 路由守卫：未登录跳转到登录页
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.name !== 'login' && !userStore.userId) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router