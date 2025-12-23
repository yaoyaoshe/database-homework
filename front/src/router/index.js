import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import MainMenu from '../views/MainMenu.vue'
// 1. 引入组件 (或者使用下方路由懒加载)

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login', name: 'Login', component: LoginView },
    // 2. 添加注册路由
    { 
      path: '/register', 
      name: 'Register', 
      component: () => import('../views/RegisterView.vue') 
    },
    { path: '/', name: 'MainMenu', component: MainMenu },
    // ... 其他路由保持不变 ...
    { 
      path: '/account-info', 
      name: 'AccountInfo', 
      component: () => import('../views/AccountInfo.vue') 
    },
    { 
      path: '/doctor-list', 
      name: 'DoctorList', 
      component: () => import('../views/DoctorList.vue') 
    },
    { 
      path: '/book-appointment', 
      name: 'BookAppointment', 
      component: () => import('../views/BookAppointment.vue') 
    },
    { 
      path: '/create-challenge', 
      name: 'CreateChallenge', 
      component: () => import('../views/CreateChallenge.vue') 
    },
    { 
      path: '/monthly-summary', 
      name: 'MonthlySummary', 
      component: () => import('../views/MonthlySummary.vue') 
    },
    { 
      path: '/search-records', 
      name: 'SearchRecords', 
      component: () => import('../views/SearchRecords.vue') 
    },
    { 
      path: '/summary-functions', 
      name: 'SummaryFunctions', 
      component: () => import('../views/SummaryFunctions.vue') 
    }
  ]
})

router.beforeEach((to, from, next) => {
  const userId = localStorage.getItem('userId')
  // 3. 修改路由守卫：允许访问注册页面 ('Register') 且无需登录
  if (to.name !== 'Login' && to.name !== 'Register' && !userId) next({ name: 'Login' })
  else next()
})

export default router