// ... import existing ...
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Layout from '@/views/Layout.vue'
import Login from '@/views/Login.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // ... Login, Register routes ...
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue')
    },
    {
      path: '/',
      component: Layout,
      redirect: '/appointments',
      children: [
        // ... existing children ...
        {
            path: 'appointments',
            component: () => import('@/views/appointment/AppointmentManager.vue')
        },
        {
            path: 'providers',
            component: () => import('@/views/provider/ProviderList.vue')
        },
        {
            path: 'challenges',
            component: () => import('@/views/challenge/ChallengeCenter.vue')
        },
        {
            path: 'reports',
            component: () => import('@/views/report/MonthlyReport.vue')
        },
        // === 新增：个人中心路由 ===
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/profile/UserProfile.vue')
        }
      ]
    }
  ]
})
// ... guard ...
// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.name !== 'Login' && to.name !== 'Register' && !userStore.userInfo.user_id) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router