<template>
  <div>
    <div class="action-bar">
      <h2>我的预约</h2>
      <el-button type="primary" @click="dialogVisible = true">新建预约</el-button>
    </div>

    <el-table :data="appointments" style="width: 100%" v-loading="loading">
      <el-table-column prop="appointment_date" label="时间">
        <template #default="scope">
          {{ formatDate(scope.row.appointment_date) }}
        </template>
      </el-table-column>
      <el-table-column prop="consultation_type" label="类型" />
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="原因" show-overflow-tooltip />
      <el-table-column label="操作" width="120">
        <template #default="scope">
          <el-button 
            v-if="scope.row.status === '已预约' && canCancel(scope.row.appointment_date)"
            size="small" 
            type="danger" 
            @click="handleCancel(scope.row)"
          >取消</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="创建预约" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="医生ID">
          <el-input v-model.number="form.provider_id" placeholder="输入 Provider ID" />
        </el-form-item>
        <el-form-item label="日期时间">
          <el-date-picker 
            v-model="form.appointment_date" 
            type="datetime" 
            placeholder="选择时间" 
            format="YYYY-MM-DD HH:mm"
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.consultation_type">
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
        <el-button type="primary" @click="submitAppointment">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAppointments, createAppointment, cancelAppointment } from '@/api/all'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import { ElMessageBox, ElMessage } from 'element-plus'

const userStore = useUserStore()
const appointments = ref([])
const loading = ref(false)
const dialogVisible = ref(false)

const form = ref({
  provider_id: null,
  appointment_date: '',
  consultation_type: '线下就诊',
  reason: '',
  duration_minutes: 30
})

const fetchList = async () => {
  loading.value = true
  try {
    appointments.value = await getAppointments(userStore.userInfo.user_id)
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusType = (status) => {
  const map = { '已预约': 'primary', '已取消': 'info', '已完成': 'success', '未到诊': 'danger' }
  return map[status] || 'default'
}

// 24小时前才可取消
const canCancel = (dateStr) => {
  return dayjs(dateStr).diff(dayjs(), 'hour') >= 24
}

const submitAppointment = async () => {
  if (!form.value.appointment_date || !form.value.provider_id) return
  
  const payload = {
    ...form.value,
    user_id: userStore.userInfo.user_id,
    // 转换为 Go RFC3339 格式
    appointment_date: dayjs(form.value.appointment_date).format() 
  }

  await createAppointment(payload)
  ElMessage.success('预约成功')
  dialogVisible.value = false
  fetchList()
}

const handleCancel = (row) => {
  ElMessageBox.prompt('请输入取消原因', '取消预约', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async ({ value }) => {
    await cancelAppointment(row.appointment_id, { reason: value })
    ElMessage.success('已取消')
    fetchList()
  })
}

onMounted(fetchList)
</script>

<style scoped>
.action-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
</style>