import axios from 'axios'
import { ElMessage } from 'element-plus'

const service = axios.create({
  baseURL: 'http://localhost:8080/api', // 对应 main.go 中的端口
  timeout: 5000
})

// 响应拦截器
service.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    console.error('API Error:', error)
    ElMessage.error(error.response?.data?.error || '请求失败，请检查后端服务')
    return Promise.reject(error)
  }
)

export default service