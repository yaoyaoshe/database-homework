<template>
  <div class="login-wrapper">
    <div class="login-box">
      <div class="login-header">
        <el-icon class="logo-icon"><Monitor /></el-icon>
        <h1>HealthTrack</h1>
        <p>您的个人健康管理专家</p>
      </div>
      
      <el-card class="login-card">
        <el-tabs v-model="activeTab" stretch class="custom-tabs">
          <el-tab-pane label="用户登录" name="login">
            <el-form 
              ref="loginFormRef" 
              :model="loginForm" 
              :rules="loginRules" 
              size="large"
              class="auth-form"
            >
              <el-form-item prop="identifier">
                <el-input 
                  v-model="loginForm.identifier" 
                  placeholder="Health ID / 邮箱 / 手机号" 
                  prefix-icon="User"
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input 
                  v-model="loginForm.password" 
                  type="password" 
                  placeholder="请输入密码" 
                  prefix-icon="Lock"
                  show-password
                  @keyup.enter="handleLogin"
                />
              </el-form-item>
              <el-button type="primary" class="submit-btn" :loading="loading" @click="handleLogin" round>
                立即登录
              </el-button>
            </el-form>
          </el-tab-pane>

          <el-tab-pane label="注册账号" name="register">
            <el-form 
              ref="registerFormRef" 
              :model="registerForm" 
              :rules="registerRules" 
              size="large"
              class="auth-form"
            >
              <el-form-item prop="health_id">
                <el-input 
                  v-model="registerForm.health_id" 
                  placeholder="设置 Health ID (唯一标识)" 
                  prefix-icon="IdCard"
                />
              </el-form-item>
              <el-form-item prop="name">
                <el-input 
                  v-model="registerForm.name" 
                  placeholder="您的真实姓名" 
                  prefix-icon="User"
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input 
                  v-model="registerForm.password" 
                  type="password" 
                  placeholder="设置登录密码" 
                  prefix-icon="Lock"
                  show-password
                />
              </el-form-item>
              
              <div class="form-row">
                <el-form-item prop="gender" class="half-width">
                  <el-select v-model="registerForm.gender" placeholder="性别">
                    <el-option label="男" value="男" />
                    <el-option label="女" value="女" />
                  </el-select>
                </el-form-item>
                <el-form-item prop="date_of_birth" class="half-width">
                   <el-date-picker
                    v-model="registerForm.date_of_birth"
                    type="date"
                    placeholder="出生日期"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                    style="width: 100%"
                  />
                </el-form-item>
              </div>

              <el-button type="success" class="submit-btn" :loading="loading" @click="handleRegister" round>
                注册并加入
              </el-button>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('login')
const loading = ref(false)
const loginFormRef = ref(null)
const registerFormRef = ref(null)

// 登录数据
// 修改点 3: 键名改为 identifier，与后端 struct 定义一致
const loginForm = reactive({
  identifier: '', 
  password: ''
})

// 注册数据
const registerForm = reactive({
  health_id: '',
  name: '',
  password: '',
  gender: '',
  date_of_birth: ''
})

// 验证规则
// 修改点 4: 规则键名改为 identifier
const loginRules = {
  identifier: [{ required: true, message: '请输入账号标识', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const registerRules = {
  health_id: [{ required: true, message: '请设置Health ID', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  password: [{ required: true, message: '请设置密码', trigger: 'blur' }, { min: 6, message: '密码至少6位', trigger: 'blur' }]
}

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.login(loginForm)
        ElMessage.success('欢迎回来 ' + userStore.userInfo.name)
        router.push('/')
      } catch (e) {
        // 错误已在 request.js 拦截器或 store 中处理
      } finally {
        loading.value = false
      }
    }
  })
}

// 处理注册
const handleRegister = async () => {
  if (!registerFormRef.value) return
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.register(registerForm)
        // 注册成功后切换到登录 Tab
        activeTab.value = 'login'
        // 修改点 5: 自动填入账号时使用 identifier
        loginForm.identifier = registerForm.health_id 
        registerFormRef.value.resetFields()
      } catch (e) {
        // 错误处理
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #f0f7ff 0%, #e1effe 100%);
  position: relative;
  overflow: hidden;
}

/* 装饰背景 */
.login-wrapper::before, .login-wrapper::after {
  content: '';
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  z-index: 0;
}
.login-wrapper::before {
  top: -10%;
  right: -5%;
  width: 500px;
  height: 500px;
  background: rgba(0, 123, 255, 0.1);
}
.login-wrapper::after {
  bottom: -10%;
  left: -5%;
  width: 400px;
  height: 400px;
  background: rgba(0, 210, 255, 0.15);
}

.login-box {
  text-align: center;
  z-index: 1;
  width: 100%;
  max-width: 460px;
  padding: 20px;
}

.login-header {
  margin-bottom: 25px;
  animation: slideDown 0.6s ease-out;
}

.logo-icon {
  font-size: 42px;
  color: var(--primary-color);
  background: white;
  padding: 12px;
  border-radius: 16px;
  box-shadow: 0 8px 16px rgba(0, 123, 255, 0.15);
  margin-bottom: 12px;
}

.login-header h1 {
  margin: 10px 0 5px;
  font-size: 26px;
  color: #2c3e50;
  font-weight: 800;
  letter-spacing: -0.5px;
}

.login-header p {
  color: #64748b;
  margin: 0;
  font-size: 14px;
}

.login-card {
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(255, 255, 255, 0.5) !important;
}

.custom-tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: #e2e8f0;
}

.auth-form {
  padding: 10px 5px 0;
  text-align: left;
}

.submit-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  letter-spacing: 1px;
  margin-top: 15px;
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.2);
}

.form-row {
  display: flex;
  gap: 15px;
}

.half-width {
  flex: 1;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 移动端适配 */
@media (max-width: 480px) {
  .login-box {
    padding: 15px;
  }
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}
</style>