import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 从 localStorage 初始化，防止刷新丢失
  const userInfo = ref(JSON.parse(localStorage.getItem('ht_user') || '{}'))

  const setUser = (user) => {
    userInfo.value = user
    localStorage.setItem('ht_user', JSON.stringify(user))
  }

  const logout = () => {
    userInfo.value = {}
    localStorage.removeItem('ht_user')
  }

  return { userInfo, setUser, logout }
})