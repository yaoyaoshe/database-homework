import { defineStore } from 'pinia'
import request from '@/utils/request'

export const useUserStore = defineStore('user', {
  state: () => ({
    userId: localStorage.getItem('userId') || null, // 模拟登录状态
    userInfo: null
  }),
  actions: {
    setUserId(id) {
      this.userId = id
      localStorage.setItem('userId', id)
    },
    async fetchUserInfo() {
      if (!this.userId) return
      try {
        const data = await request.get(`/users/${this.userId}`)
        this.userInfo = data
      } catch (error) {
        console.error(error)
      }
    },
    logout() {
      this.userId = null
      this.userInfo = null
      localStorage.removeItem('userId')
    }
  }
})