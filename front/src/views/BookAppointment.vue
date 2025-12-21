<template>
  <div class="page-container">
    <h2>预约服务</h2>
    <el-card>
      <el-form :model="form" label-width="120px">
        <el-form-item label="医生ID">
           <el-input-number v-model="form.provider_id" :min="1" />
          <span style="margin-left:10px; color:#999; font-size:12px">（请输入已知的医生ID）</span>
        </el-form-item>

        <el-form-item label="日期时间">
          <el-date-picker 
            v-model="form.appointment_date" 
            type="datetime" 
            placeholder="选择日期时间" 
          />
        </el-form-item>

        <el-form-item label="预约类型">
          <el-select v-model="form.consultation_type" placeholder="请选择">
            <el-option label="线下就诊" value="线下就诊" />
            <el-option label="线上咨询" value="线上咨询" />
            <el-option label="电话咨询" value="电话咨询" />
          </el-select>
        </el-form-item>

        <el-form-item label="预约原因">
          <el-input v-model="form.reason" type="textarea" />
        </el-form-item>

        <el-form-item label="时长(分钟)">
          <el-input-number v-model="form.duration_minutes" :step="15" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitAppointment">提交预约</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <h3 style="margin-top:30px">我的预约记录</h3>
    <el-table :data="appointments" stripe style="width: 100%">
      <el-table-column prop="appointment_date" label="时间">
        <template #default="scope">{{ formatTime(scope.row.appointment_date) }}</template>
      </el-table-column>
      <el-table-column prop="provider_id" label="医生ID" />
      <el-table-column prop="consultation_type" label="类型" />
      <el-table-column prop="status" label="状态" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button v-if="scope.row.status !== '已取消'" type="danger" size="small" @click="cancelAppt(scope.row.appointment_id)">取消</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const appointments = ref([])

const form = reactive({
  provider_id: 1,
  appointment_date: '',
  consultation_type: '线下就诊',
  reason: '',
  duration_minutes: 30
})

const fetchAppointments = async () => {
  const res = await request.get(`/users/${userStore.userId}/appointments`)
  appointments.value = res
}

const submitAppointment = async () => {
  if (!form.appointment_date) return ElMessage.warning('请选择时间')
  
  // RFC3339 格式化
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
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async ({ value }) => {
    try {
      await request.put(`/appointments/${id}/cancel`, { reason: value || '用户取消' })
      ElMessage.success('已取消')
      fetchAppointments()
    } catch(e) {}
  })
}

const formatTime = (t) => new Date(t).toLocaleString()

onMounted(fetchAppointments)
</script>

<style scoped>
.page-container { padding: 20px; max-width: 800px; margin: 0 auto; }
</style>