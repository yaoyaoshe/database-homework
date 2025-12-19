import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue') },
  { path: '/', name: 'MainMenu', component: () => import('../views/MainMenu.vue') },
  { path: '/account-info', name: 'AccountInfo', component: () => import('../views/AccountInfo.vue') },
  { path: '/book-appointment', name: 'BookAppointment', component: () => import('../views/BookAppointment.vue') },
  { path: '/create-challenge', name: 'CreateChallenge', component: () => import('../views/CreateChallenge.vue') },
  { path: '/monthly-summary', name: 'MonthlySummary', component: () => import('../views/MonthlySummary.vue') },
  { path: '/search-records', name: 'SearchRecords', component: () => import('../views/SearchRecords.vue') },
  { path: '/account-functions', name: 'AccountFunctions', component: () => import('../views/AccountFunctions.vue') },
  { path: '/summary-functions', name: 'SummaryFunctions', component: () => import('../views/SummaryFunctions.vue') }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router