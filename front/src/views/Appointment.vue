<template>
  <div>
    <div class="toolbar">
      <h2>我的预约</h2>
      <el-button type="primary" @click="showBookDialog">新建预约</el-button>
    </div>

    <el-table :data="appointments" stripe style="width: 100%">
      <el-table-column prop="appointment_date" label="时间" :formatter="formatDate" />
      <el-table-column prop="consultation_type" label="类型" />
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="原因" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button size="small" type="danger" v-if="row.status === '已预约'" @click="cancelAppt(row.appointment_id)">取消</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="预约挂号" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="医生ID">
          <el-input v-model.number="form.provider_id" type="number" placeholder="请输入医生ID" />
        </el-form-item>
        <el-form-item label="时间">
           <el-date-picker v-model="form.appointment_date" type="datetime" placeholder="选择日期时间" style="width: 100%" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.consultation_type" placeholder="选择类型">
            <el-option label="线下就诊" value="线下就诊" />
            <el-option label="线上咨询" value="线上咨询" />
          </el-select>
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="form.reason" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAppointment">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const appointments = ref([])
const dialogVisible = ref(false)
const form = ref({ provider_id: '', appointment_date: '', consultation_type: '线下就诊', reason: '' })

const fetchList = async () => {
  appointments.value = await request.get(`/users/${userStore.userId}/appointments`)
}

const formatDate = (row) => {
  return new Date(row.appointment_date).toLocaleString()
}

const getStatusType = (status) => {
  if (status === '已预约') return 'primary'
  if (status === '已取消') return 'info'
  return 'success'
}

const showBookDialog = () => {
  dialogVisible.value = true
  form.value = { provider_id: '', appointment_date: '', consultation_type: '线下就诊', reason: '' }
}

const submitAppointment = async () => {
  if (!form.value.appointment_date) return
  // 转换为 RFC3339 格式
  const dateStr = new Date(form.value.appointment_date).toISOString()
  
  await request.post('/appointments', {
    user_id: parseInt(userStore.userId),
    provider_id: parseInt(form.value.provider_id),
    appointment_date: dateStr,
    duration_minutes: 30,
    consultation_type: form.value.consultation_type,
    reason: form.value.reason
  })
  ElMessage.success('预约成功')
  dialogVisible.value = false
  fetchList()
}

const cancelAppt = async (id) => {
  await request.put(`/appointments/${id}/cancel`, { reason: '用户主动取消' })
  ElMessage.success('已取消')
  fetchList()
}

onMounted(fetchList)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
</style>