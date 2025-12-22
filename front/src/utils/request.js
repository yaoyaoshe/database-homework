import axios from 'axios'
import { ElMessage } from 'element-plus'

const service = axios.create({
  baseURL: 'http://localhost:8080/api', // 对应 Go 后端地址
  timeout: 5000
})

// 响应拦截器
service.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    console.error('Request Error:', error)
    const msg = error.response?.data?.error || '请求失败'
    ElMessage.error(msg)
    return Promise.reject(error)
  }
)

export default service