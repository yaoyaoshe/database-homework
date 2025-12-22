<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <div class="card-header">
          <h2>注册 HealthTrack</h2>
        </div>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="90px">
        <el-form-item label="Health ID" prop="health_id">
          <el-input 
            v-model="form.health_id" 
            placeholder="设置您的唯一健康标识ID" 
          />
        </el-form-item>

        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="您的真实姓名" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="form.password" 
            type="password" 
            show-password 
            placeholder="设置登录密码" 
          />
        </el-form-item>

        <el-form-item label="性别" prop="gender">
          <el-select v-model="form.gender" placeholder="请选择性别" style="width: 100%">
            <el-option label="男" value="男" />
            <el-option label="女" value="女" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>

        <el-form-item label="出生日期" prop="date_of_birth">
          <el-date-picker
            v-model="form.date_of_birth"
            type="date"
            placeholder="选择出生日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleRegister" :loading="loading" class="w-100">
            立即注册
          </el-button>
        </el-form-item>

        <div class="links">
          <router-link to="/login">已有账号？返回登录</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { register } from '@/api/all' // 确保你已经在 api/all.js 中导出了 register 方法
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)

// 表单数据，字段名需与后端 handler.go 中的 CreateUser 结构体 JSON tag 对应
const form = reactive({
  health_id: '',
  name: '',
  password: '',
  gender: '',
  date_of_birth: ''
})

// 表单验证规则
const rules = {
  health_id: [
    { required: true, message: '请输入 Health ID', trigger: 'blur' },
    { min: 3, message: '长度至少 3 个字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // 调用后端注册接口 POST /api/users
        await register(form)
        ElMessage.success('注册成功，请登录')
        // 注册成功后跳转回登录页
        router.push('/login')
      } catch (error) {
        // 错误通常会在 request.js 拦截器中处理，这里不需要额外操作
        console.error(error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}

.register-card {
  width: 450px;
}

.card-header {
  text-align: center;
}

.w-100 {
  width: 100%;
}

.links {
  margin-top: 15px;
  text-align: center;
  font-size: 14px;
}
</style>