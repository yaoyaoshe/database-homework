<template>
  <div class="profile-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>基本信息</span>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="姓名">{{ userInfo.name }}</el-descriptions-item>
        <el-descriptions-item label="Health ID">{{ userInfo.health_id }}</el-descriptions-item>
        <el-descriptions-item label="性别">{{ userInfo.gender || '未知' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="box-card mt-20">
      <template #header>
        <div class="card-header">
          <span>联系方式管理</span>
          <div>
            <el-button size="small" @click="openDialog('email')">添加邮箱</el-button>
            <el-button size="small" @click="openDialog('phone')">添加电话</el-button>
          </div>
        </div>
      </template>
      
      <el-empty description="暂无法查看联系方式列表 (需后端接口支持)，但可以添加新联系方式" />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogType === 'email' ? '添加邮箱' : '添加电话'" width="400px">
      <el-form :model="form" label-width="80px">
        <el-form-item :label="dialogType === 'email' ? '邮箱' : '电话'">
          <el-input v-model="form.value" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getUser, addEmail, addPhone } from '@/api/all'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const userInfo = ref({})
const dialogVisible = ref(false)
const dialogType = ref('email') // 'email' or 'phone'
const form = ref({ value: '' })

const loadData = async () => {
  // 调用 GET /users/{id}
  userInfo.value = await getUser(userStore.userInfo.user_id)
}

const openDialog = (type) => {
  dialogType.value = type
  form.value.value = ''
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!form.value.value) return
  
  try {
    if (dialogType.value === 'email') {
      // POST /users/{id}/emails
      await addEmail(userStore.userInfo.user_id, { email: form.value.value })
    } else {
      // POST /users/{id}/phones
      await addPhone(userStore.userInfo.user_id, { phone: form.value.value })
    }
    ElMessage.success('添加成功')
    dialogVisible.value = false
  } catch (error) {
    // 错误处理已在 request.js 中
  }
}

onMounted(loadData)
</script>

<style scoped>
.mt-20 { margin-top: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>