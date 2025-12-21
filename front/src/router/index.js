import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import MainMenu from '../views/MainMenu.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login', name: 'Login', component: LoginView },
    { path: '/', name: 'MainMenu', component: MainMenu },
    { 
      path: '/account-info', 
      name: 'AccountInfo', 
      component: () => import('../views/AccountInfo.vue') 
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

// 简单路由守卫
router.beforeEach((to, from, next) => {
  const userId = localStorage.getItem('userId')
  if (to.name !== 'Login' && !userId) next({ name: 'Login' })
  else next()
})

export default router