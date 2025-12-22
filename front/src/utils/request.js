import axios from 'axios'
import { ElMessage } from 'element-plus'

const service = axios.create({
  baseURL: '/api',
  timeout: 5000
})

service.interceptors.response.use(
  response => response.data,
  error => {
    ElMessage.error(error.response?.data?.error || '请求失败')
    return Promise.reject(error)
  }
)

// 在 request.js 的 service.interceptors.request.use 中
service.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      // 这里的 Header Key 需与后端中间件验证的 Key 一致，通常是 Authorization
      // 假设后端格式为 "Bearer <token>" 或直接 "<token>"
      // 根据提供的Go后端，并没有明显的中间件代码展示，但通常JWT放在Authorization
      // 补充文档 Source 297 显示返回了 JWT
      config.headers['Authorization'] = `Bearer ${token}` 
    }
    return config
  },
  error => Promise.reject(error)
)

export default service