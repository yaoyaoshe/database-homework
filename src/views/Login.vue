<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header><h2 class="login-title">HealthTrack 登录</h2></template>
      <el-form>
        <el-form-item><el-input v-model="form.username" placeholder="邮箱 / 用户ID" size="large" /></el-form-item>
        <el-form-item><el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password /></el-form-item>
        <el-button type="primary" class="login-btn" @click="handleLogin" :loading="loading">立即登录</el-button>
      </el-form>
    </el-card>
  </div>
</template>


<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '@/api/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const form = ref({ username: '', password: '' })
const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  try {
    const response = await login(form.value)
    // 保存用户ID到本地存储
    localStorage.setItem('userId', response.user_id)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    // 错误处理
  } finally {
    loading.value = false
  }
}
</script>

<!-- <script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { login } from '@/api/user'; // 引入新接口
import { ElMessage } from 'element-plus';

const router = useRouter();
const form = ref({ username: '', password: '' });
const loading = ref(false);

const handleLogin = async () => {
  loading.value = true;
  try {
    await login(form.value);
    ElMessage.success('登录成功');
    router.push('/');
  } catch (error) {
    // 错误处理已在 request.js 拦截器中处理，这里可留空或做额外逻辑
  } finally {
    loading.value = false;
  }
};
</script> -->

<style scoped>
.login-container { display: flex; justify-content: center; align-items: center; height: 100vh; background: #eef1f6; }
.login-card { width: 400px; border-radius: 12px; }
.login-title { text-align: center; color: #303133; margin: 0; }
.login-btn { width: 100%; height: 40px; font-size: 16px; margin-top: 10px; }
</style>