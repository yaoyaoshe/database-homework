import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 从 localStorage 初始化，防止刷新丢失
  const userInfo = ref(JSON.parse(localStorage.getItem('health_user') || '{}'))
  const userId = ref(localStorage.getItem('health_uid') || '')

  function login(data) {
    userInfo.value = data
    userId.value = data.user_id
    localStorage.setItem('health_user', JSON.stringify(data))
    localStorage.setItem('health_uid', data.user_id)
  }

  function logout() {
    userInfo.value = {}
    userId.value = ''
    localStorage.removeItem('health_user')
    localStorage.removeItem('health_uid')
  }

  return { userInfo, userId, login, logout }
})