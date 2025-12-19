import request from '@/utils/request'

// 获取月度健康摘要 (适配接口文档 POST /reports/generate)
export function getMonthlyData(year, month) {
  return request({
    url: '/api/reports/generate',
    method: 'post',
    data: { year, month }
  })
}

// 搜索记录 (适配接口文档 GET /stats/search)
export function searchRecords(params) {
  return request({
    url: '/api/stats/search',
    method: 'get',
    params
  })
}

// 获取全局统计报表 (适配接口文档 GET /stats/active_users)
export function getGlobalStats(dateRange) {
  return request({
    url: '/api/stats/active_users',
    method: 'get',
    params: { limit: 10 }
  })
}


// import request from '@/utils/request'

// // 获取月度健康摘要
// export function getMonthlyData(year, month) {
//   // return request({ 
//   //   url: '/stats/monthly', 
//   //   method: 'get', 
//   //   params: { year, month } 
//   // })

//   return Promise.resolve({
//     metrics: {
//       steps: [6000, 8000, 12000, 7000, 9500, 11000, 8000],
//       weight: [70.5, 70.4, 70.2, 70.1],
//       bp_sys: [120, 118, 122, 119],
//       bp_dia: [80, 78, 82, 79]
//     },
//     appointments: [
//       { date: `${year}-${month}-05`, provider: "李医生", type: "线下门诊", status: "已完成" },
//       { date: `${year}-${month}-20`, provider: "王医生", type: "线上问诊", status: "待就诊" }
//     ],
//     summaryText: `${year}年${month}月总结：您本月平均每日步数保持在 8,500 步。`
//   })
// }

// // 搜索记录
// export function searchRecords(params) {
//   // return request({ url: '/stats/search', method: 'get', params })

//   return Promise.resolve([
//     { id: 1, date: "2024-11-01", type: "预约", detail: "年度体检 - 李医生" },
//     { id: 2, date: "2024-11-15", type: "挑战", detail: "目标：步行100公里" },
//     { id: 3, date: "2024-12-01", type: "指标", detail: "记录体重: 70kg" }
//   ])
// }

// // 获取全局统计报表
// export function getGlobalStats(dateRange) {
//   // return request({ url: '/stats/global', method: 'post', data: dateRange })

//   return Promise.resolve({
//     totalAppts: 42,
//     metrics: { weight: { avg: 70.2, min: 69.5, max: 71.0 } },
//     topChallenges: [
//       { name: "夏季长跑挑战", count: 150 },
//       { name: "早起瑜伽周", count: 89 }
//     ],
//     activeUsers: [
//       { name: "用户_A", score: 500 },
//       { name: "用户_B", score: 450 }
//     ]
//   })
// }