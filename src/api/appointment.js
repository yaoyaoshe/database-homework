import request from '@/utils/request'

// 获取医生列表 (适配接口文档 GET /providers)
export function getProviders() {
  return request({
    url: '/api/providers',
    method: 'get'
  })
}

// 提交预约 (适配接口文档 POST /appointments)
export function submitAppointment(data) {
  return request({
    url: '/api/appointments',
    method: 'post',
    data
  })
}

// import request from '@/utils/request'

// // 获取医生列表
// export function getProviders() {
//   // return request({ url: '/appointment/providers', method: 'get' })
  
//   return Promise.resolve([
//     { id: 101, name: "李医生 (主治医师)", license: "MD-9921" },
//     { id: 102, name: "王医生 (专家)", license: "PT-4421" },
//     { id: 103, name: "赵医生 (诊断科)", license: "DG-1111" }
//   ])
// }

// // 提交预约
// export function submitAppointment(data) {
//   // return request({ url: '/appointment/submit', method: 'post', data })
  
//   console.log("模拟提交预约数据:", data)
//   return Promise.resolve({ success: true, id: "APT-" + Date.now() })
// }