<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2 class="login-title">HealthTrack 登录</h2>
      </template>
      <el-form :model="form" label-width="0">
        <el-form-item>
          <el-input v-model="form.identifier" placeholder="请输入Health ID / 邮箱 / 电话" prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="请输入密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-button type="primary" class="login-btn" @click="handleLogin" :loading="loading">登录</el-button>
      </el-form>
      <div class="tips">测试账号: zhangsan123 (密码) / HealthID: BJ...1001</div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import request from '../api/request'
import { useUserStore } from '../stores/user'
import { useRouter } from 'vue-router'
import { sha256 } from 'js-sha256' // 需要 npm install js-sha256

const form = ref({ identifier: '', password: '' })
const loading = ref(false)
const userStore = useUserStore()
const router = useRouter()

const handleLogin = async () => {
  if (!form.value.identifier || !form.value.password) return
  loading.value = true
  try {
    // 适配后端：后端直接对比哈希值，所以前端需要先哈希
    // 注意：真实生产环境不应在前端哈希，这里是为了适配你提供的后端逻辑
    const hashedPassword = sha256(form.value.password)
    
    const res = await request.post('/auth/login', {
      identifier: form.value.identifier,
      password: hashedPassword 
    })
    userStore.login(res)
    router.push('/')
  } catch (e) {
    // error handled in interceptor
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container { height: 100vh; display: flex; justify-content: center; align-items: center; background: #f0f2f5; }
.login-card { width: 400px; }
.login-title { text-align: center; color: #409EFF; }
.login-btn { width: 100%; }
.tips { margin-top: 10px; font-size: 12px; color: #999; text-align: center; }
</style>