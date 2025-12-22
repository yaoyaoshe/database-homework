<template>
  <div class="login-wrapper">
    <div class="login-box">
      <div class="login-header">
        <el-icon class="logo-icon"><Monitor /></el-icon>
        <h1>HealthTrack</h1>
        <p>您的个人健康管理专家</p>
      </div>
      
      <el-card class="login-card">
        <h2 class="form-title">欢迎回来</h2>
        <el-form @submit.prevent="handleLogin" size="large">
          <el-form-item label="用户 ID">
            <el-input 
              v-model="userIdInput" 
              placeholder="请输入您的 User ID (例如: 1)" 
              prefix-icon="User"
            />
          </el-form-item>
          <el-button type="primary" class="login-btn" @click="handleLogin" round>
            立即登录
          </el-button>
        </el-form>
        <div class="login-footer">
          <span>还没有账号? <a href="#">立即注册</a></span>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const userIdInput = ref('1')
const router = useRouter()
const userStore = useUserStore()

const handleLogin = () => {
  if (userIdInput.value) {
    userStore.setUserId(userIdInput.value)
    router.push('/')
  }
}
</script>

<style scoped>
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #e0f7fa 0%, #e1effe 100%);
  position: relative;
  overflow: hidden;
}

/* 装饰背景圆 */
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
  max-width: 420px;
  padding: 20px;
}

.login-header {
  margin-bottom: 30px;
}

.logo-icon {
  font-size: 48px;
  color: #007bff;
  background: white;
  padding: 10px;
  border-radius: 12px;
  box-shadow: 0 10px 20px rgba(0,123,255,0.2);
  margin-bottom: 15px;
}

.login-header h1 {
  margin: 10px 0 5px;
  font-size: 28px;
  color: #333;
  font-weight: 700;
}

.login-header p {
  color: #666;
  margin: 0;
}

.form-title {
  text-align: left;
  margin-bottom: 25px;
  font-weight: 600;
  color: #2c3e50;
}

.login-btn {
  width: 100%;
  font-weight: bold;
  letter-spacing: 1px;
  margin-top: 10px;
  height: 45px;
}

.login-footer {
  margin-top: 20px;
  font-size: 14px;
  color: #999;
}

.login-footer a {
  color: #007bff;
  text-decoration: none;
  font-weight: 500;
}
</style>