<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>账户中心</h2>
    </div>
    
    <div class="profile-header-card">
      <el-avatar :size="80" style="background: #007bff; font-size: 32px;">
        {{ userInfo?.name ? userInfo.name[0] : 'U' }}
      </el-avatar>
      <div class="profile-info">
        <h3>{{ userInfo?.name || '加载中...' }}</h3>
        <p>Health ID: <strong>{{ userInfo?.health_id || '-' }}</strong></p>
      </div>
    </div>

    <el-card class="main-content-card">
      <el-tabs type="card" class="custom-tabs">
        <el-tab-pane label="个人档案">
          <div class="info-grid">
            <div class="info-item">
              <label>姓名</label>
              <div class="value">{{ userInfo?.name || '-' }}</div>
            </div>
            <div class="info-item">
              <label>性别</label>
              <div class="value">{{ userInfo?.gender || '-' }}</div>
            </div>
             <div class="info-item">
              <label>出生日期</label>
              <div class="value">{{ formatDate(userInfo?.date_of_birth) }}</div>
            </div>
             <div class="info-item">
              <label>注册时间</label>
              <div class="value">{{ formatDate(userInfo?.created_at) }}</div>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="邮箱管理">
          <div class="action-bar">
            <el-input v-model="newEmail" placeholder="example@mail.com" style="width: 240px; margin-right: 10px;">
              <template #prefix><el-icon><Message /></el-icon></template>
            </el-input>
            <el-button type="primary" @click="addEmail">添加邮箱</el-button>
          </div>
          <el-table :data="emailList" border stripe>
            <el-table-column prop="email_address" label="邮箱地址" />
            <el-table-column prop="is_verified" label="状态" width="120">
               <template #default="scope">
                 <el-tag :type="scope.row.is_verified ? 'success' : 'warning'">
                   {{ scope.row.is_verified ? '已验证' : '未验证' }}
                 </el-tag>
               </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="center">
              <template #default="scope">
                <el-button type="danger" link icon="Delete" @click="deleteEmail(scope.row.email_id)"></el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="电话管理">
           <div class="action-bar">
            <el-input v-model="newPhone" placeholder="输入电话号码" style="width: 240px; margin-right: 10px;">
              <template #prefix><el-icon><Iphone /></el-icon></template>
            </el-input>
            <el-button type="primary" @click="addPhone">添加电话</el-button>
          </div>
           <el-table :data="phoneList" border stripe>
            <el-table-column prop="phone_number" label="电话号码" />
            <el-table-column label="操作" width="100" align="center">
              <template #default="scope">
                <el-button type="danger" link icon="Delete" @click="deletePhone(scope.row.phone_id)"></el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="医疗团队">
          <div class="action-bar">
            <el-input v-model="providerIdInput" placeholder="输入医生执照号/ID" style="width: 240px; margin-right: 10px;">
              <template #prefix><el-icon><FirstAidKit /></el-icon></template>
            </el-input>
            <el-button type="success" @click="linkProvider">关联医生</el-button>
          </div>
           <el-alert title="提示：此处仅提供关联功能，已关联列表需调用额外接口获取" type="info" show-icon :closable="false" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo || {}) // 提供默认空对象防止报错
const emailList = ref([]) 
const phoneList = ref([])
const newEmail = ref('')
const newPhone = ref('')
const providerIdInput = ref('')

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}

const addEmail = async () => {
  if (!newEmail.value) return
  try {
    const res = await request.post(`/users/${userStore.userId}/emails`, { email: newEmail.value })
    emailList.value.push(res)
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
/* 统一样式 */
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 24px; color: #333; }

.profile-header-card {
  background: white;
  padding: 30px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
}

.profile-info h3 { margin: 0 0 8px 0; font-size: 24px; }
.profile-info p { margin: 0; color: #666; font-size: 16px; }

.main-content-card { min-height: 400px; }
.action-bar { display: flex; margin-bottom: 20px; align-items: center; background: #f8f9fa; padding: 15px; border-radius: 8px; }
.info-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); gap: 24px; padding: 10px; }
.info-item label { display: block; color: #888; margin-bottom: 5px; font-size: 13px; }
.info-item .value { font-size: 16px; font-weight: 500; color: #333; }
</style>