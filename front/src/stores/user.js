import { defineStore } from 'pinia'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', {
  state: () => ({
    userId: localStorage.getItem('userId') || null,
    token: localStorage.getItem('token') || null,
    userInfo: JSON.parse(localStorage.getItem('userInfo') || 'null')
  }),
  
  actions: {
    // 登录动作
    async login(loginForm) {
      try {
        // 根据《接口文档补充》POST /api/auth/login
        // 入参: identity, password
        const data = await request.post('/auth/login', loginForm)
        
        // 保存状态
        this.token = data.token
        this.userId = data.user_id
        this.userInfo = {
          name: data.name,
          health_id: data.health_id
        }

        // 持久化存储
        localStorage.setItem('token', data.token)
        localStorage.setItem('userId', data.user_id)
        localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
        
        return true
      } catch (error) {
        console.error('登录失败:', error)
        throw error
      }
    },

    // 注册动作
    async register(registerForm) {
      try {
        // 根据《接口文档》POST /users
        // 入参: health_id, name, password, gender, date_of_birth
        await request.post('/users', registerForm)
        ElMessage.success('注册成功，请登录')
        return true
      } catch (error) {
        console.error('注册失败:', error)
        throw error
      }
    },

    // 登出
    logout() {
      this.userId = null
      this.token = null
      this.userInfo = null
      localStorage.removeItem('userId')
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  }
})