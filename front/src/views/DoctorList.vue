<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>寻找医生</h2>
      <el-button type="success" plain style="margin-left: auto" @click="$router.push('/account-info')">
        <el-icon style="margin-right: 5px"><UserFilled /></el-icon> 我的医疗团队
      </el-button>
    </div>

    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="专业领域">
          <el-select v-model="filters.specialty" placeholder="全部" clearable style="width: 180px">
            <el-option label="全科" value="全科" />
            <el-option label="内科" value="内科" />
            <el-option label="外科" value="外科" />
            <el-option label="儿科" value="儿科" />
            <el-option label="心理咨询" value="心理咨询" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="filters.verified" label="只看已认证" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div class="provider-grid" v-loading="loading">
      <el-empty v-if="providers.length === 0 && !loading" description="未找到匹配的医生" />
      
      <el-card v-for="doctor in providers" :key="doctor.provider_id" class="doctor-card" shadow="hover">
        <div class="doctor-header">
          <el-avatar :size="60" style="background: #e6f7ff; color: #1890ff; font-size: 24px">
            {{ doctor.name[0] }}
          </el-avatar>
          <div class="doctor-info">
            <h3>
              {{ doctor.name }}
              <el-icon v-if="doctor.is_verified" color="#67C23A" title="已认证"><CircleCheckFilled /></el-icon>
            </h3>
            <el-tag size="small" effect="plain">{{ doctor.specialty || '全科医生' }}</el-tag>
          </div>
        </div>
        
        <div class="doctor-details">
          <p class="detail-item"><el-icon><Medal /></el-icon> {{ doctor.qualification || '执业医师' }}</p>
          <p class="detail-item"><el-icon><OfficeBuilding /></el-icon> {{ getContactInfo(doctor.contact_info, '科室') }}</p>
        </div>

        <div class="doctor-actions">
          <el-button 
            v-if="isLinked(doctor.provider_id)" 
            type="success" 
            plain 
            block 
            disabled
          >
            <el-icon><Check /></el-icon> 已关联
          </el-button>
          
          <el-button 
            v-else 
            type="primary" 
            plain 
            block 
            @click="openLinkDialog(doctor)"
          >
            建立关联
          </el-button>
        </div>
      </el-card>
    </div>

    <el-dialog v-model="dialogVisible" title="关联医生" width="400px">
      <p>您正在请求关联医生：<strong>{{ currentDoctor?.name }}</strong></p>
      <el-form :model="linkForm" label-position="top" style="margin-top: 20px">
        <el-form-item label="请选择关系类型">
          <el-select v-model="linkForm.relation" style="width: 100%">
            <el-option label="家庭医生" value="家庭医生" />
            <el-option label="主治医生" value="主治医生" />
            </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmLink" :loading="linking">确认关联</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const providers = ref([])
const myProviders = ref([]) // 存储已关联的医生 ID
const loading = ref(false)
const linking = ref(false)
const dialogVisible = ref(false)
const currentDoctor = ref(null)

const filters = reactive({
  specialty: '',
  verified: true
})

const linkForm = reactive({
  relation: '家庭医生'
})

// 同时获取所有医生和已关联医生
const fetchData = async () => {
  loading.value = true
  try {
    const params = {}
    if (filters.specialty) params.specialty = filters.specialty
    if (filters.verified) params.verified = 'true'
    
    // 并行请求
    const [allRes, myRes] = await Promise.all([
      request.get('/providers', { params }),
      userStore.userId ? request.get(`/users/${userStore.userId}/providers`) : Promise.resolve([])
    ])

    providers.value = allRes || []
    myProviders.value = myRes || [] // myRes 是包含 provider_id 的对象数组
  } catch (e) {
    ElMessage.error('数据加载失败')
  } finally {
    loading.value = false
  }
}

// 检查是否已关联
const isLinked = (pid) => {
  return myProviders.value.some(p => p.provider_id === pid)
}

const openLinkDialog = (doctor) => {
  currentDoctor.value = doctor
  linkForm.relation = '家庭医生'
  dialogVisible.value = true
}

const confirmLink = async () => {
  if (!userStore.userId) return ElMessage.warning('请先登录')
  linking.value = true
  try {
    await request.post(`/users/${userStore.userId}/providers`, {
      provider_id: currentDoctor.value.provider_id,
      relationship_type: linkForm.relation
    })
    ElMessage.success('关联成功')
    dialogVisible.value = false
    fetchData() // 刷新列表状态
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '关联失败')
  } finally {
    linking.value = false
  }
}

const getContactInfo = (jsonStr, key) => {
  try {
    const obj = JSON.parse(jsonStr)
    return obj[key] || '未知'
  } catch { return '-' }
}

onMounted(fetchData)
</script>

<style scoped>
/* 样式保持不变 */
.page-container { padding: 24px; max-width: 1200px; margin: 0 auto; }
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.filter-card { margin-bottom: 20px; }
.provider-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}
.doctor-header { display: flex; align-items: center; gap: 15px; margin-bottom: 15px; }
.doctor-info h3 { margin: 0 0 5px 0; display: flex; align-items: center; gap: 5px; }
.doctor-details { color: #666; font-size: 14px; margin-bottom: 20px; min-height: 50px; }
.detail-item { margin: 5px 0; display: flex; align-items: center; gap: 8px; }
.doctor-actions .el-button { width: 100%; }
</style>