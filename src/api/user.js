import request from '@/utils/request'

// 获取用户详细信息 (适配接口文档 GET /users/profile)
// export function getUserProfile() {
//   return request({
//     url: '/api/users/profile',
//     method: 'get'
//   })
// }

export function getUserProfile(userId) {
  return request({
    url: `/api/users/${userId}`,
    method: 'get'
  })
}

// 登录接口 (适配接口文档 POST /users/login)
export function login(data) {
  return request({
    url: '/api/users/login',
    method: 'post',
    data
  })
}

// 更新用户信息 (适配接口文档 PUT /users/profile)
export function updateProfile(data) {
  return request({
    url: '/api/users/profile',
    method: 'put',
    data
  })
}





// import request from '@/utils/request'

// // 获取用户详细信息
// export function getUserProfile() {
//   // --- 真实对接时使用下方代码 ---
//   // return request({
//   //   url: '/user/profile',
//   //   method: 'get'
//   // })

//   // --- 暂时返回默认值 ---
//   return Promise.resolve({
//     name: "张三",
//     healthId: "HID-8839201",
//     dob: "1985-04-12",
//     emails: [
//       { id: 1, address: "zhangsan@example.com", verified: true },
//       { id: 2, address: "zhang.work@corp.com", verified: false }
//     ],
//     phones: [
//       { id: 1, number: "138-0000-0000", verified: true }
//     ],
//     providers: [
//       { id: 101, license: "MD-9921", name: "李医生", type: "全科医师", isPrimary: true },
//       { id: 102, license: "PT-4421", name: "王医生", type: "理疗专家", isPrimary: false }
//     ]
//   })
// }

// // 登录接口
// export function login(data) {
//   // return request({ url: '/user/login', method: 'post', data })
//   return Promise.resolve({ token: "fake-token-12345" })
// }

// // 更新用户信息
// export function updateProfile(data) {
//   // return request({ url: '/user/update', method: 'post', data })
//   return Promise.resolve({ code: 200, message: "更新成功" })
// }