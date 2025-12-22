import { defineStore } from 'pinia'
import request from '@/utils/request'

export const useUserStore = defineStore('user', {
  state: () => ({
    userId: localStorage.getItem('userId') || null,
    userName: localStorage.getItem('userName') || '',
    healthId: localStorage.getItem('healthId') || '',
    userInfo: null
  }),
  actions: {
    // 登录动作
    async login(identifier, password) {
      try {
        const res = await request.post('/auth/login', {
          identifier,
          password
        })
        
        // 保存登录态
        this.userId = res.user_id
        this.userName = res.name
        this.healthId = res.health_id
        
        localStorage.setItem('userId', res.user_id)
        localStorage.setItem('userName', res.name)
        localStorage.setItem('healthId', res.health_id)
        
        return true
      } catch (error) {
        console.error("登录失败", error)
        throw error
      }
    },
    
    // 获取详细信息
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
      this.userName = ''
      this.healthId = ''
      
      localStorage.removeItem('userId')
      localStorage.removeItem('userName')
      localStorage.removeItem('healthId')
    }
  }
})