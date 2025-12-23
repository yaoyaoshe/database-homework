<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>预约服务</h2>
    </div>

    <el-row :gutter="24">
      <el-col :lg="8" :md="24" style="margin-bottom: 20px;">
        <el-card class="form-card">
          <template #header>
            <div class="card-header">
              <span><el-icon><EditPen /></el-icon> 新建预约</span>
            </div>
          </template>
          
          <el-form :model="form" label-position="top" size="large">
            
            <el-form-item label="选择医生" required>
              <el-select 
                v-model="form.provider_id" 
                placeholder="请选择您的医生" 
                style="width: 100%" 
                no-data-text="暂无关联医生"
              >
                <el-option 
                  v-for="doc in linkedProviders" 
                  :key="doc.provider_id" 
                  :label="doc.name + ' (' + (doc.relationship_type || '医生') + ')'" 
                  :value="doc.provider_id" 
                >
                  <span style="float: left">{{ doc.name }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ doc.specialty }}</span>
                </el-option>
              </el-select>
              
              <div class="form-helper-link">
                <span v-if="linkedProviders.length === 0" style="color: #e6a23c; margin-right: 5px;">
                  <el-icon><Warning /></el-icon> 您尚未关联任何医生
                </span>
                <router-link to="/doctor-list">去医生库查找/关联医生 &gt;</router-link>
              </div>
            </el-form-item>

            <el-form-item label="预约时间" required>
              <el-date-picker 
                v-model="form.appointment_date" 
                type="datetime" 
                placeholder="请选择就诊时间" 
                style="width: 100%"
                :disabled-date="disabledDate"
                format="YYYY-MM-DD HH:mm"
              />
            </el-form-item>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="预约类型">
                  <el-select v-model="form.consultation_type" placeholder="请选择" style="width: 100%">
                    <el-option label="线下就诊" value="线下就诊" />
                    <el-option label="线上咨询" value="线上咨询" />
                    <el-option label="电话咨询" value="电话咨询" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="预计时长 (分)">
                  <el-input-number v-model="form.duration_minutes" :step="15" :min="15" style="width: 100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="症状描述 / 原因">
              <el-input v-model="form.reason" type="textarea" :rows="3" placeholder="请简要描述您的症状..." />
            </el-form-item>

            <el-button type="primary" @click="submitAppointment" style="width: 100%; margin-top: 10px;" :loading="submitting">
              提交预约申请
            </el-button>
          </el-form>
        </el-card>
      </el-col>

      <el-col :lg="16" :md="24">
        <el-card>
          <template #header>
            <div class="card-header" style="justify-content: space-between;">
              <span style="display:flex; align-items:center; gap:8px;">
                <el-icon><List /></el-icon> 我的预约
              </span>
              <el-radio-group v-model="viewMode" size="small">
                <el-radio-button label="list">列表视图</el-radio-button>
                <el-radio-button label="calendar">日历视图</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          
          <div v-if="viewMode === 'list'">
            <el-table :data="appointments" style="width: 100%" height="500" v-loading="loading">
              <el-table-column label="时间" min-width="160">
                <template #default="scope">
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <el-icon><Clock /></el-icon>
                    {{ formatTime(scope.row.appointment_date) }}
                  </div>
                </template>
              </el-table-column>
              
              <el-table-column label="医生" width="120">
                <template #default="scope">
                  {{ getProviderName(scope.row.provider_id) }}
                </template>
              </el-table-column>

              <el-table-column prop="consultation_type" label="类型" width="100" />
              
              <el-table-column prop="status" label="状态" width="100">
                <template #default="scope">
                  <el-tag :type="getStatusType(scope.row.status)" effect="light" round>
                    {{ scope.row.status }}
                  </el-tag>
                </template>
              </el-table-column>
              
              <el-table-column label="操作" width="90" fixed="right">
                <template #default="scope">
                  <el-button 
                    v-if="scope.row.status !== '已取消'" 
                    type="danger" 
                    link
                    icon="CircleClose"
                    @click="cancelAppt(scope.row.appointment_id)"
                  >取消</el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="appointments.length === 0 && !loading" description="暂无预约记录" />
          </div>

          <div v-else class="calendar-view">
            <el-calendar>
              <template #date-cell="{ data }">
                <div class="date-cell-content" :class="{ 'has-appt': getApptOnDay(data.day).length > 0 }">
                  <span class="day-number">{{ data.day.split('-')[2] }}</span>
                  
                  <div class="appt-list-scroll">
                    <div v-for="appt in getApptOnDay(data.day)" :key="appt.appointment_id" class="appt-mark">
                      <el-tooltip :content="`${formatTime(appt.appointment_date)}: ${appt.reason || '无描述'}`" placement="top">
                        <el-tag 
                          size="small" 
                          :type="getStatusType(appt.status)" 
                          effect="dark" 
                          style="margin-top:2px; font-size: 10px; width: 100%; border:none; text-overflow: ellipsis; overflow: hidden;"
                        >
                          {{ appt.consultation_type }}
                        </el-tag>
                      </el-tooltip>
                    </div>
                  </div>
                </div>
              </template>
            </el-calendar>
          </div>
          
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const appointments = ref([])
const linkedProviders = ref([]) // 存储已关联的医生列表
const viewMode = ref('list') 
const loading = ref(false)
const submitting = ref(false)

const form = reactive({
  provider_id: null,
  appointment_date: '',
  consultation_type: '线下就诊',
  reason: '',
  duration_minutes: 30
})

// 初始化数据
const initData = async () => {
  if (!userStore.userId) return
  loading.value = true
  await Promise.all([fetchAppointments(), fetchLinkedProviders()])
  loading.value = false
}

// 获取我的预约
const fetchAppointments = async () => {
  try {
    const res = await request.get(`/users/${userStore.userId}/appointments`)
    appointments.value = res || []
  } catch(e) {
    appointments.value = []
  }
}

// 🆕 获取已关联的医生
const fetchLinkedProviders = async () => {
  try {
    const res = await request.get(`/users/${userStore.userId}/providers`)
    linkedProviders.value = res || []
    
    // 如果有关联医生且当前表单未选，默认选中第一个
    if (linkedProviders.value.length > 0 && !form.provider_id) {
      form.provider_id = linkedProviders.value[0].provider_id
    }
  } catch(e) {
    console.error('获取关联医生失败', e)
  }
}

const submitAppointment = async () => {
  if (!form.provider_id) return ElMessage.warning('请选择医生')
  if (!form.appointment_date) return ElMessage.warning('请选择时间')
  
  submitting.value = true
  const rfcDate = new Date(form.appointment_date).toISOString()
  try {
    await request.post('/appointments', {
      user_id: parseInt(userStore.userId),
      provider_id: form.provider_id,
      appointment_date: rfcDate,
      duration_minutes: form.duration_minutes,
      consultation_type: form.consultation_type,
      reason: form.reason
    })
    ElMessage.success('预约成功')
    // 重置部分表单
    form.reason = ''
    form.appointment_date = ''
    fetchAppointments()
  } catch(e) {
    ElMessage.error(e.response?.data?.error || '预约失败')
  } finally {
    submitting.value = false
  }
}

const cancelAppt = (id) => {
  ElMessageBox.prompt('请输入取消原因', '取消预约', {
    confirmButtonText: '确定', cancelButtonText: '取消',
  }).then(async ({ value }) => {
    try {
      await request.put(`/appointments/${id}/cancel`, { reason: value || '用户取消' })
      ElMessage.success('已取消')
      fetchAppointments()
    } catch(e) {}
  })
}

// 辅助函数：通过ID反查医生姓名
const getProviderName = (pid) => {
  const doc = linkedProviders.value.find(p => p.provider_id === pid)
  return doc ? doc.name : `ID: ${pid}`
}

// 修改为显示年月日 (YYYY-MM-DD)
const formatTime = (t) => {
  if(!t) return ''
  const d = new Date(t)
  return `${d.getFullYear()}-${(d.getMonth() + 1).toString().padStart(2, '0')}-${d.getDate().toString().padStart(2, '0')}`
}

const getStatusType = (status) => {
  if (status === '已预约') return 'primary'
  if (status === '已完成') return 'success'
  if (status === '已取消') return 'info'
  return 'warning'
}

const disabledDate = (time) => {
  return time.getTime() < Date.now() - 8.64e7 // 禁止选择今天之前的日期
}

const getApptOnDay = (dateStr) => {
  return appointments.value.filter(appt => {
    if (!appt.appointment_date) return false
    return appt.appointment_date.substring(0, 10) === dateStr
  })
}

onMounted(initData)
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.page-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}
.page-header h2 { margin: 0; font-size: 24px; color: #333; }

.card-header { 
  font-weight: bold; 
  font-size: 16px; 
  color: #007bff; 
  display: flex; 
  align-items: center; 
  gap: 8px; 
}
.form-card { border-top: 4px solid #007bff !important; }

/* 引导链接样式 */
.form-helper-link {
  margin-top: 6px;
  font-size: 13px;
  text-align: right;
}
.form-helper-link a {
  color: #409eff;
  text-decoration: none;
  transition: color 0.2s;
}
.form-helper-link a:hover {
  color: #66b1ff;
  text-decoration: underline;
}

/* 日历样式优化 */
.calendar-view :deep(.el-calendar-table .el-calendar-day) {
  height: 85px;
  padding: 4px;
}
.date-cell-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.has-appt {
  background-color: #f0f9eb;
  border-radius: 4px;
}
.day-number {
  font-size: 12px;
  color: #999;
  margin-bottom: 2px;
}
.appt-list-scroll {
  flex: 1;
  overflow-y: auto;
}
.appt-list-scroll::-webkit-scrollbar {
  display: none; 
}
</style>