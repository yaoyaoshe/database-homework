import request from '@/utils/request'

// 创建挑战 (适配接口文档 POST /challenges)
export function createChallenge(data) {
  return request({
    url: '/api/challenges',
    method: 'post',
    data
  })
}

// import request from '@/utils/request'

// // 创建挑战
// export function createChallenge(data) {
//   // return request({ url: '/challenge/create', method: 'post', data })
  
//   console.log("模拟创建挑战:", data)
//   return Promise.resolve({ success: true })
// }

// // 获取挑战列表 (如果需要)
// export function getChallenges() {
//   return Promise.resolve([])
// }