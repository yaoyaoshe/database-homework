import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const service = axios.create({
  // 未来这里填入真实的后端接口地址，例如 'http://localhost:8080/api'
  baseURL: import.meta.env.VITE_APP_BASE_API || '/api', 
  timeout: 5000 // 请求超时时间
})

// request 拦截器
service.interceptors.request.use(
  config => {
    //在这里可以添加 token，例如：
    // if (store.getters.token) {
    //   config.headers['Authorization'] = getToken()
    // }
    return config
  },
  error => {
    console.log(error)
    return Promise.reject(error)
  }
)

// response 拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    // 这里可以根据后端约定的状态码判断请求是否成功
    // 假设 200 是成功
    if (res.code !== 200) {
      // ElMessage.error(res.message || 'Error')
      // return Promise.reject(new Error(res.message || 'Error'))
      
      // 暂时直接返回数据，方便调试
      return res
    } else {
      return res
    }
  },
  error => {
    console.log('err' + error)
    ElMessage.error(error.message)
    return Promise.reject(error)
  }
)

export default service