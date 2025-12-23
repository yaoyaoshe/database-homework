<template>
  <div class="login-wrapper">
    <div class="login-box">
      <div class="login-header">
        <el-icon class="logo-icon"><Plus /></el-icon>
        <h1>创建账号</h1>
        <p>加入 HealthTrack，开启健康生活</p>
      </div>
      
      <el-card class="login-card">
        <el-form ref="formRef" :model="form" :rules="rules" size="large" @submit.prevent="handleRegister">
          <el-form-item prop="health_id">
            <el-input 
              v-model="form.health_id" 
              placeholder="设置您的 Health ID (唯一标识)" 
              prefix-icon="postcard"
            />
          </el-form-item>
          
          <el-form-item prop="name">
            <el-input 
              v-model="form.name" 
              placeholder="您的姓名" 
              prefix-icon="User"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input 
              v-model="form.password" 
              type="password"
              placeholder="设置密码" 
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item prop="confirmPassword">
            <el-input 
              v-model="form.confirmPassword" 
              type="password"
              placeholder="确认密码" 
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-row :gutter="10">
            <el-col :span="12">
              <el-form-item prop="gender">
                <el-select v-model="form.gender" placeholder="性别">
                  <el-option label="男" value="男" />
                  <el-option label="女" value="女" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item prop="date_of_birth">
                <el-date-picker 
                  v-model="form.date_of_birth" 
                  type="date" 
                  placeholder="出生日期" 
                  value-format="YYYY-MM-DD"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-button type="primary" class="login-btn" :loading="loading" @click="handleRegister" round>
            立即注册
          </el-button>
        </el-form>
        <div class="login-footer">
          <span>已有账号? <router-link to="/login">去登录</router-link></span>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const formRef = ref(null)

const form = reactive({
  health_id: '',
  name: '',
  password: '',
  confirmPassword: '',
  gender: '',
  date_of_birth: ''
})

const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const rules = {
  health_id: [{ required: true, message: '请输入Health ID', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  confirmPassword: [{ validator: validatePass2, trigger: 'blur' }]
}

// SHA-256 加密工具函数 (保持与登录一致)
const sha256 = async (message) => {
  const msgBuffer = new TextEncoder().encode(message);
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  return hashHex;
}

const handleRegister = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // 为了安全性，注册时也对密码进行哈希，保持与LoginView一致
        const hashedPassword = await sha256(form.password)

        const payload = {
          health_id: form.health_id,
          name: form.name,
          password: hashedPassword,
          gender: form.gender,
          date_of_birth: form.date_of_birth
        }

        await request.post('/users', payload)
        
        ElMessage.success('注册成功，请登录')
        router.push('/login')
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '注册失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
/* 复用 LoginView 的样式 */
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #e0f7fa 0%, #e1effe 100%);
  position: relative;
  overflow: hidden;
}
.login-wrapper::before {
  content: '';
  position: absolute;
  top: -10%;
  right: -5%;
  width: 500px;
  height: 500px;
  background: linear-gradient(45deg, #007bff, #00d2ff);
  border-radius: 50%;
  opacity: 0.1;
  filter: blur(80px);
}
.login-box {
  text-align: center;
  z-index: 1;
  width: 100%;
  max-width: 450px; /* 稍微宽一点以容纳更多字段 */
  padding: 20px;
}
.login-header { margin-bottom: 20px; }
.logo-icon {
  font-size: 40px;
  color: #007bff;
  background: white;
  padding: 10px;
  border-radius: 12px;
  box-shadow: 0 10px 20px rgba(0,123,255,0.2);
  margin-bottom: 15px;
}
.login-header h1 { margin: 10px 0 5px; font-size: 26px; color: #333; font-weight: 700; }
.login-header p { color: #666; margin: 0; }
.login-btn { width: 100%; font-weight: bold; letter-spacing: 1px; margin-top: 10px; height: 45px; }
.login-footer { margin-top: 20px; font-size: 14px; color: #999; }
.login-footer a { color: #007bff; text-decoration: none; font-weight: 500; }
</style>