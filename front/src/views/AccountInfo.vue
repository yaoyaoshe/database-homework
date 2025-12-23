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
      <div style="margin-left: auto;">
        <el-button type="primary" @click="handleUpdateProfile" :loading="saving">保存修改</el-button>
      </div>
    </div>

    <el-card class="main-content-card">
      <el-tabs type="card" class="custom-tabs">
        <el-tab-pane label="个人档案">
           <el-form label-position="top" class="info-grid">
            <el-form-item label="姓名"><el-input v-model="editForm.name" /></el-form-item>
            <el-form-item label="性别">
              <el-select v-model="editForm.gender" style="width: 100%">
                <el-option label="男" value="男" />
                <el-option label="女" value="女" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>
            <el-form-item label="出生日期">
              <el-date-picker v-model="editForm.date_of_birth" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-form>
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
                <el-table-column prop="is_verified" label="状态" width="100">
                <template #default="scope">
                    <el-tag :type="scope.row.is_verified ? 'success' : 'info'">{{ scope.row.is_verified ? '已验证' : '未验证' }}</el-tag>
                </template>
                </el-table-column>
                <el-table-column label="操作" width="80" align="center">
                <template #default="scope"><el-button type="danger" link icon="Delete" @click="deleteEmail(scope.row.email_id)"></el-button></template>
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
                <el-table-column prop="phone_type" label="类型" width="100" />
                <el-table-column label="操作" width="80" align="center">
                <template #default="scope"><el-button type="danger" link icon="Delete" @click="deletePhone(scope.row.phone_id)"></el-button></template>
                </el-table-column>
            </el-table>
        </el-tab-pane>

        <el-tab-pane label="医疗团队">
            <div class="action-bar">
                <span>您的专属医疗团队</span>
                <el-button link type="primary" @click="$router.push('/doctor-list')" style="margin-left: auto">
                去寻找/添加医生 <el-icon><ArrowRight /></el-icon>
                </el-button>
            </div>
            <el-table :data="providerList" border stripe empty-text="暂无关联医生">
                <el-table-column prop="name" label="医生姓名" width="120" />
                <el-table-column prop="specialty" label="专业领域" width="120" />
                <el-table-column prop="relationship_type" label="关系类型" width="120" />
                <el-table-column label="操作" width="80" align="center" fixed="right">
                <template #default="scope">
                    <el-popconfirm title="确定要解除关联吗？" @confirm="unlinkProvider(scope.row.provider_id)">
                    <template #reference><el-button type="danger" size="small" link>解绑</el-button></template>
                    </el-popconfirm>
                </template>
                </el-table-column>
            </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="家庭成员">
          <div v-if="pendingRequests.length > 0" class="pending-section">
            <h4 style="margin: 0 0 10px 0; color: #E6A23C"><el-icon><Bell /></el-icon> 待处理的关联请求</h4>
            <div v-for="req in pendingRequests" :key="req.initiator_id" class="request-item">
              <span>用户 <strong>{{ req.initiator_name }}</strong> 想要添加您为 <strong>{{ req.relationship }}</strong></span>
              <div>
                <el-button type="success" size="small" @click="handleVerify(req.initiator_id, true)">同意</el-button>
                <el-button type="danger" size="small" @click="handleVerify(req.initiator_id, false)">拒绝</el-button>
              </div>
            </div>
            <el-divider />
          </div>

          <div class="action-bar">
            <el-input v-model="familyForm.targetHealthId" placeholder="成员 Health ID" style="width: 200px; margin-right: 10px;">
              <template #prefix><el-icon><Postcard /></el-icon></template>
            </el-input>
            
            <el-select v-model="familyForm.relation" placeholder="关系" style="width: 140px; margin-right: 10px;">
              <el-option label="父母" value="父母" />
              <el-option label="子女" value="子女" />
              <el-option label="配偶" value="配偶" />
              <el-option label="兄弟姐妹" value="兄弟姐妹" />
              <el-option label="亲戚" value="亲戚" />
              <el-option label="朋友" value="朋友" />
              <el-option label="监护人" value="监护人" />
            </el-select>
            <el-button type="primary" @click="addFamilyMember">添加成员</el-button>
          </div>

          <el-table :data="familyList" border stripe v-loading="loading" empty-text="暂无家庭成员">
            <el-table-column prop="name" label="姓名" />
            <el-table-column prop="relationship" label="关系" width="120" />
            <el-table-column prop="is_verified" label="验证状态" width="120">
               <template #default="scope">
                 <el-tag :type="scope.row.is_verified ? 'success' : 'warning'">
                   {{ scope.row.is_verified ? '已验证' : '待验证' }}
                 </el-tag>
               </template>
            </el-table-column>
            <el-table-column label="操作" width="80" align="center">
              <template #default="scope">
                <el-popconfirm title="确定删除该成员？" @confirm="deleteFamilyMember(scope.row.related_user_id)">
                  <template #reference>
                    <el-button type="danger" link icon="Delete"></el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, toRefs } from 'vue'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const { userInfo } = toRefs(userStore)

const loading = ref(false)
const saving = ref(false)
const emailList = ref([])
const phoneList = ref([])
const providerList = ref([])
const familyList = ref([])
const pendingRequests = ref([]) // 新增：存储待处理请求

const newEmail = ref('')
const newPhone = ref('')

const familyForm = reactive({
  targetHealthId: '',
  relation: '亲戚' // 默认值修正为有效的枚举值
})

const editForm = reactive({
  name: '',
  gender: '',
  date_of_birth: ''
})

const fetchAllData = async () => {
  if (!userStore.userId) return
  loading.value = true
  try {
    const userRes = await request.get(`/users/${userStore.userId}`)
    userStore.userInfo = userRes
    editForm.name = userRes.name
    editForm.gender = userRes.gender
    editForm.date_of_birth = userRes.date_of_birth ? userRes.date_of_birth.substring(0, 10) : ''

    const [emailRes, phoneRes, providerRes, familyRes, reqRes] = await Promise.all([
      request.get(`/users/${userStore.userId}/emails`),
      request.get(`/users/${userStore.userId}/phones`),
      request.get(`/users/${userStore.userId}/providers`),
      request.get(`/users/${userStore.userId}/family`),
      request.get(`/users/${userStore.userId}/family/requests`) // 获取待处理请求
    ])
    
    emailList.value = emailRes || []
    phoneList.value = phoneRes || []
    providerList.value = providerRes || []
    familyList.value = familyRes || []
    pendingRequests.value = reqRes || []

  } catch(e) {
    console.error(e)
    ElMessage.error('加载信息失败')
  } finally {
    loading.value = false
  }
}

// ... (handleUpdateProfile, addEmail, deleteEmail, addPhone, deletePhone, unlinkProvider 保持不变) ...

const handleUpdateProfile = async () => {
  saving.value = true
  try {
    await request.put(`/users/${userStore.userId}`, editForm)
    ElMessage.success('个人资料已更新')
    fetchAllData()
  } catch(e) {} finally {
    saving.value = false
  }
}
const addEmail = async () => {
  if (!newEmail.value) return
  try {
    await request.post(`/users/${userStore.userId}/emails`, { email: newEmail.value })
    ElMessage.success('邮箱添加成功')
    newEmail.value = ''
    fetchAllData()
  } catch(e){}
}
const deleteEmail = async (id) => {
  try {
    await request.delete(`/users/${userStore.userId}/emails/${id}`)
    ElMessage.success('删除成功')
    fetchAllData()
  } catch(e){}
}
const addPhone = async () => {
  if (!newPhone.value) return
  try {
    await request.post(`/users/${userStore.userId}/phones`, { phone: newPhone.value })
    ElMessage.success('电话添加成功')
    newPhone.value = ''
    fetchAllData()
  } catch(e){}
}
const deletePhone = async (id) => {
   try {
    await request.delete(`/users/${userStore.userId}/phones/${id}`)
    ElMessage.success('删除成功')
    fetchAllData()
  } catch(e){}
}
const unlinkProvider = async (pid) => {
  try {
    await request.delete(`/users/${userStore.userId}/providers/${pid}`)
    ElMessage.success('已解除关联')
    fetchAllData()
  } catch(e) {}
}

const addFamilyMember = async () => {
  if(!familyForm.targetHealthId) return ElMessage.warning('请输入成员 Health ID')
  try {
    await request.post(`/users/${userStore.userId}/family`, {
      target_health_id: familyForm.targetHealthId,
      relationship: familyForm.relation
    })
    ElMessage.success('已发送请求，等待对方验证')
    familyForm.targetHealthId = ''
    fetchAllData()
  } catch(e) {
     ElMessage.error(e.response?.data?.error || '添加失败，请检查Health ID是否正确')
  }
}

const deleteFamilyMember = async (rid) => {
  try {
    await request.delete(`/users/${userStore.userId}/family/${rid}`)
    ElMessage.success('已移除成员')
    fetchAllData()
  } catch(e) {}
}

// 新增：处理验证请求
const handleVerify = async (initiatorId, accept) => {
  try {
    await request.post(`/users/${userStore.userId}/family/verify`, {
      initiator_id: initiatorId,
      accept: accept
    })
    ElMessage.success(accept ? '已同意关联' : '已拒绝请求')
    fetchAllData()
  } catch(e) {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  fetchAllData()
})
</script>

<style scoped>
/* 原有样式保持不变 */
.page-container { padding: 20px; }
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
.main-content-card { min-height: 500px; border-radius: 16px; }
.action-bar { display: flex; margin-bottom: 20px; align-items: center; background: #f8f9fa; padding: 15px; border-radius: 8px; }
.info-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; padding: 10px; }

/* 新增样式 */
.pending-section { background: #fdf6ec; padding: 15px; border-radius: 8px; margin-bottom: 20px; border: 1px solid #faecd8; }
.request-item { display: flex; justify-content: space-between; align-items: center; padding: 8px 0; border-bottom: 1px dashed #eee; }
.request-item:last-child { border-bottom: none; }
</style>