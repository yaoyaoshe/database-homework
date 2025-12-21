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
            <el-form-item label="医生 ID">
              <el-input-number v-model="form.provider_id" :min="1" style="width: 100%" />
            </el-form-item>

            <el-form-item label="预约时间">
              <el-date-picker 
                v-model="form.appointment_date" 
                type="datetime" 
                placeholder="选择日期时间" 
                style="width: 100%"
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
                  <el-input-number v-model="form.duration_minutes" :step="15" style="width: 100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="症状描述 / 原因">
              <el-input v-model="form.reason" type="textarea" :rows="3" placeholder="请简要描述您的症状..." />
            </el-form-item>

            <el-button type="primary" @click="submitAppointment" style="width: 100%; margin-top: 10px;">
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
            <el-table :data="appointments" style="width: 100%" height="500">
              <el-table-column label="时间" min-width="160">
                <template #default="scope">
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <el-icon><Clock /></el-icon>
                    {{ formatTime(scope.row.appointment_date) }}
                  </div>
                </template>
              </el-table-column>
              
              <el-table-column prop="provider_id" label="医生ID" width="80" align="center" />
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
            <el-empty v-if="appointments.length === 0" description="暂无预约记录" />
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
                          style="margin-top:2px; font-size: 10px; width: 100%; border:none;"
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
const viewMode = ref('list') // 控制默认视图

const form = reactive({
  provider_id: 1,
  appointment_date: '',
  consultation_type: '线下就诊',
  reason: '',
  duration_minutes: 30
})

const fetchAppointments = async () => {
  try {
    const res = await request.get(`/users/${userStore.userId}/appointments`)
    appointments.value = res || [] // 确保是数组
  } catch(e) {
    appointments.value = []
  }
}

const submitAppointment = async () => {
  if (!form.appointment_date) return ElMessage.warning('请选择时间')
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
    fetchAppointments()
  } catch(e) {}
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

const formatTime = (t) => {
  if(!t) return ''
  return new Date(t).toLocaleString()
}

const getStatusType = (status) => {
  if (status === '已预约') return 'primary'
  if (status === '已完成') return 'success'
  if (status === '已取消') return 'info'
  return 'warning'
}

// 日历辅助函数：筛选特定日期的预约
const getApptOnDay = (dateStr) => {
  return appointments.value.filter(appt => {
    if (!appt.appointment_date) return false
    // 截取 YYYY-MM-DD
    return appt.appointment_date.substring(0, 10) === dateStr
  })
}

onMounted(fetchAppointments)
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
/* 隐藏滚动条但允许滚动 */
.appt-list-scroll::-webkit-scrollbar {
  display: none; 
}
</style>