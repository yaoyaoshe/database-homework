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
            <el-form-item label="姓名">
              <el-input v-model="editForm.name" />
            </el-form-item>
            <el-form-item label="性别">
              <el-select v-model="editForm.gender" style="width: 100%">
                <el-option label="男" value="男" />
                <el-option label="女" value="女" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>
            <el-form-item label="出生日期">
              <el-date-picker 
                v-model="editForm.date_of_birth" 
                type="date" 
                placeholder="选择日期" 
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item label="注册时间">
              <el-input :value="formatDate(userInfo?.created_at)" disabled />
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
          <el-alert title="注意：后端未提供邮箱列表查询接口，刷新后列表将重置。" type="info" show-icon :closable="false" style="margin-bottom: 10px;" />
          <el-table :data="emailList" border stripe>
            <el-table-column prop="email_address" label="邮箱地址" />
            <el-table-column prop="is_verified" label="状态" width="100">
               <template #default="scope">
                 <el-tag :type="scope.row.is_verified ? 'success' : 'info'">
                   {{ scope.row.is_verified ? '已验证' : '未验证' }}
                 </el-tag>
               </template>
            </el-table-column>
            <el-table-column label="操作" width="80" align="center">
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
          <el-alert title="注意：后端未提供电话列表查询接口，刷新后列表将重置。" type="info" show-icon :closable="false" style="margin-bottom: 10px;" />
           <el-table :data="phoneList" border stripe>
            <el-table-column prop="phone_number" label="电话号码" />
            <el-table-column prop="phone_type" label="类型" width="100" />
            <el-table-column label="操作" width="80" align="center">
              <template #default="scope">
                <el-button type="danger" link icon="Delete" @click="deletePhone(scope.row.phone_id)"></el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="医疗团队">
          <div class="action-bar">
            <el-input v-model="linkForm.providerId" placeholder="输入医生 ID" style="width: 150px; margin-right: 10px;" />
            
            <el-select v-model="linkForm.relation" placeholder="关系类型" style="width: 150px; margin-right: 10px;">
              <el-option label="家庭医生" value="家庭医生" />
              <el-option label="主治医生" value="主治医生" />
              <el-option label="专科顾问" value="专科顾问" />
            </el-select>

            <el-button type="success" @click="linkProvider">关联医生</el-button>
            <el-button link type="primary" @click="$router.push('/doctor-list')" style="margin-left: auto">
              去医生库查找 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
          
          <el-table :data="providerList" border stripe v-loading="loading">
            <el-table-column prop="name" label="医生姓名" width="120">
              <template #default="scope">
                <strong>{{ scope.row.name }}</strong>
              </template>
            </el-table-column>
            
            <el-table-column prop="specialty" label="专业领域" width="120" />
            
            <el-table-column prop="relationship_type" label="关系类型" width="120">
              <template #default="scope">
                <el-tag effect="plain">{{ scope.row.relationship_type || '未设置' }}</el-tag>
              </template>
            </el-table-column>

            <el-table-column label="科室" width="150">
              <template #default="scope">
                {{ getContactDetail(scope.row.contact_info, '科室') }}
              </template>
            </el-table-column>

            <el-table-column label="操作" width="80" align="center" fixed="right">
              <template #default="scope">
                <el-popconfirm title="确定要解除关联吗？" @confirm="unlinkProvider(scope.row.provider_id)">
                  <template #reference>
                    <el-button type="danger" size="small" link>解绑</el-button>
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
const allProvidersCache = ref([]) // 缓存所有医生列表，用于前端匹配

const newEmail = ref('')
const newPhone = ref('')

// 关联表单
const linkForm = reactive({
  providerId: '',
  relation: '家庭医生'
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
    // 修改：仅调用后端实际存在的接口，避免 404
    // 1. 获取用户信息
    const userRes = await request.get(`/users/${userStore.userId}`)
    userStore.userInfo = userRes
    editForm.name = userRes.name
    editForm.gender = userRes.gender
    editForm.date_of_birth = userRes.date_of_birth ? userRes.date_of_birth.substring(0, 10) : ''

    // 2. 初始化列表为空 (因为后端没有提供查询接口)
    emailList.value = []
    phoneList.value = []
    providerList.value = []

    // 3. 尝试获取所有医生列表，并根据用户的 primary_provider_id 模拟显示主治医生
    try {
      const allProviders = await request.get('/providers')
      allProvidersCache.value = allProviders // 缓存以备后用
      
      if (userRes.primary_provider_id) {
        const primary = allProviders.find(p => p.provider_id === userRes.primary_provider_id)
        if (primary) {
          // 手动添加关系类型，模拟关联数据
          providerList.value.push({
            ...primary,
            relationship_type: '主治医生 (Primary)'
          })
        }
      }
    } catch (err) {
      console.error('无法获取医生列表', err)
    }

  } catch(e) {
    console.error(e)
    ElMessage.error('加载用户信息失败')
  } finally {
    loading.value = false
  }
}

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
    // 调用 POST 接口 (存在)
    const res = await request.post(`/users/${userStore.userId}/emails`, { email: newEmail.value })
    ElMessage.success('邮箱添加成功')
    // 手动添加到本地列表显示，因为无法重新拉取
    emailList.value.push(res)
    newEmail.value = ''
  } catch(e){}
}
const deleteEmail = async (id) => {
  try {
    await request.delete(`/users/${userStore.userId}/emails/${id}`)
    ElMessage.success('删除成功')
    // 手动从本地列表移除
    emailList.value = emailList.value.filter(e => e.email_id !== id)
  } catch(e){}
}

const addPhone = async () => {
  if (!newPhone.value) return
  try {
    // 调用 POST 接口 (存在)
    const res = await request.post(`/users/${userStore.userId}/phones`, { phone: newPhone.value })
    ElMessage.success('电话添加成功')
    // 手动添加到本地列表显示
    phoneList.value.push(res)
    newPhone.value = ''
  } catch(e){}
}
const deletePhone = async (id) => {
   try {
    await request.delete(`/users/${userStore.userId}/phones/${id}`)
    ElMessage.success('删除成功')
    // 手动从本地列表移除
    phoneList.value = phoneList.value.filter(p => p.phone_id !== id)
  } catch(e){}
}

const linkProvider = async () => {
  if (!linkForm.providerId) return ElMessage.warning('请输入医生ID')
  try {
    const pid = parseInt(linkForm.providerId)
    // 1. 调用后端关联接口
    await request.post(`/users/${userStore.userId}/providers`, {
      provider_id: pid,
      relationship_type: linkForm.relation
    })
    ElMessage.success('关联医生成功')
    
    // 2. 尝试从缓存的所有医生中找到详细信息并添加到列表
    const providerDetail = allProvidersCache.value.find(p => p.provider_id === pid)
    if (providerDetail) {
      // 检查是否已存在
      const exists = providerList.value.find(p => p.provider_id === pid)
      if (!exists) {
        providerList.value.push({
          ...providerDetail,
          relationship_type: linkForm.relation
        })
      }
    } else {
      // 如果找不到详情，刷新页面(虽然只会显示主治医生)
      fetchAllData()
    }
    
    linkForm.providerId = ''
  } catch(e) {
    ElMessage.error(e.response?.data?.error || '关联失败')
  }
}

const unlinkProvider = async (pid) => {
  try {
    await request.delete(`/users/${userStore.userId}/providers/${pid}`)
    ElMessage.success('已解除关联')
    providerList.value = providerList.value.filter(p => p.provider_id !== pid)
  } catch(e) {}
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}

const getContactDetail = (jsonStr, key) => {
  if (!jsonStr) return '-'
  try {
    const obj = JSON.parse(jsonStr)
    return obj[key] || '-'
  } catch (e) {
    return '-'
  }
}

onMounted(() => {
  fetchAllData()
})
</script>

<style scoped>
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
</style>