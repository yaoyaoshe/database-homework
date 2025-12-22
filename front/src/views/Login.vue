<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header><h2>HealthTrack 登录</h2></template>
      <el-form :model="form" label-width="0">
        <el-form-item>
          <el-input v-model="form.identifier" placeholder="Health ID / 邮箱 / 手机号" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码" show-password />
        </el-form-item>
        <el-button type="primary" class="w-100" @click="handleLogin" :loading="loading">登录</el-button>
        <div class="links">
          <router-link to="/register">没有账号? 去注册</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { login } from '@/api/all'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const form = ref({ identifier: '', password: '' })
const loading = ref(false)
const userStore = useUserStore()
const router = useRouter()

const handleLogin = async () => {
  if(!form.value.identifier || !form.value.password) return ElMessage.warning('请输入完整信息')
  
  loading.value = true
  try {
    const res = await login(form.value)
    // res 包含 { user_id, health_id, name }
    userStore.setUser(res)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (e) {
    // 错误在 request.js 已处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container { display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
.login-card { width: 400px; }
.w-100 { width: 100%; }
.links { margin-top: 15px; text-align: center; }
</style>