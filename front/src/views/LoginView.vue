<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2>HealthTrack 登录</h2>
      <el-form @submit.prevent="handleLogin">
        <el-form-item label="User ID">
          <el-input v-model="userIdInput" placeholder="请输入用户ID (例如: 1)" />
        </el-form-item>
        <el-button type="primary" @click="handleLogin" block>进入系统</el-button>
      </el-form>
    </el-card>
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
.login-container { display: flex; justify-content: center; align-items: center; height: 100vh; background: #f0f2f5; }
.login-card { width: 400px; }
</style>