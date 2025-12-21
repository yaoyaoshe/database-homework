<template>
  <div class="page-container">
    <h2>账户信息与功能</h2>
    
    <el-tabs type="border-card">
      <el-tab-pane label="个人详情">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="姓名">{{ userInfo?.name }}</el-descriptions-item>
          <el-descriptions-item label="Health ID">{{ userInfo?.health_id }}</el-descriptions-item>
          <el-descriptions-item label="性别">{{ userInfo?.gender }}</el-descriptions-item>
          <el-descriptions-item label="出生日期">{{ formatDate(userInfo?.date_of_birth) }}</el-descriptions-item>
        </el-descriptions>
      </el-tab-pane>

      <el-tab-pane label="邮箱管理">
        <div class="action-bar">
          <el-input v-model="newEmail" placeholder="输入新邮箱" style="width: 200px; margin-right: 10px;" />
          <el-button type="primary" @click="addEmail">添加邮箱</el-button>
        </div>
        <el-table :data="emailList" style="width: 100%; margin-top: 10px;">
          <el-table-column prop="email_address" label="邮箱地址" />
          <el-table-column prop="is_verified" label="验证状态">
             <template #default="scope">{{ scope.row.is_verified ? '已验证' : '未验证' }}</template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button type="danger" size="small" @click="deleteEmail(scope.row.email_id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="电话管理">
         <div class="action-bar">
          <el-input v-model="newPhone" placeholder="输入新电话" style="width: 200px; margin-right: 10px;" />
          <el-button type="primary" @click="addPhone">添加电话</el-button>
        </div>
        <el-table :data="phoneList" style="width: 100%; margin-top: 10px;">
          <el-table-column prop="phone_number" label="电话号码" />
          <el-table-column label="操作">
            <template #default="scope">
              <el-button type="danger" size="small" @click="deletePhone(scope.row.phone_id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="医疗提供方">
        <div class="action-bar">
          <el-input v-model="providerIdInput" placeholder="输入医生ID" style="width: 200px; margin-right: 10px;" />
          <el-button type="primary" @click="linkProvider">关联医生</el-button>
        </div>
         <el-alert title="注意：后端未提供获取已关联医生列表的接口，此处仅提供关联/解绑功能" type="warning" :closable="false" style="margin-top:10px"/>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo)

// 本地状态模拟列表（因为后端没有GET列表接口）
const emailList = ref([]) 
const phoneList = ref([])
const newEmail = ref('')
const newPhone = ref('')
const providerIdInput = ref('')

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}

const addEmail = async () => {
  if (!newEmail.value) return
  try {
    const res = await request.post(`/users/${userStore.userId}/emails`, { email: newEmail.value })
    emailList.value.push(res) // 添加到本地列表显示
    newEmail.value = ''
    ElMessage.success('邮箱添加成功')
  } catch(e){}
}

const deleteEmail = async (id) => {
  try {
    await request.delete(`/users/${userStore.userId}/emails/${id}`)
    emailList.value = emailList.value.filter(e => e.email_id !== id)
    ElMessage.success('删除成功')
  } catch(e){}
}

const addPhone = async () => {
  if (!newPhone.value) return
  try {
    const res = await request.post(`/users/${userStore.userId}/phones`, { phone: newPhone.value })
    phoneList.value.push(res)
    newPhone.value = ''
    ElMessage.success('电话添加成功')
  } catch(e){}
}

const deletePhone = async (id) => {
   try {
    await request.delete(`/users/${userStore.userId}/phones/${id}`)
    phoneList.value = phoneList.value.filter(p => p.phone_id !== id)
    ElMessage.success('删除成功')
  } catch(e){}
}

const linkProvider = async () => {
  if (!providerIdInput.value) return
  try {
    await request.post(`/users/${userStore.userId}/providers/${providerIdInput.value}`)
    ElMessage.success('关联医生成功')
  } catch(e){}
}

onMounted(() => {
  userStore.fetchUserInfo()
})
</script>

<style scoped>
.page-container { padding: 20px; max-width: 1000px; margin: 0 auto; }
.action-bar { display: flex; margin-bottom: 20px; }
</style>