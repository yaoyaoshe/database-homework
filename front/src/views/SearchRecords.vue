<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>搜索记录</h2>
    </div>

    <el-card>
      <el-form :inline="true" class="search-form">
        <el-form-item label="搜索类型">
          <el-select v-model="searchType" style="width: 140px">
            <el-option label="预约记录" value="appointment" />
            <el-option label="健康统计" value="health" />
          </el-select>
        </el-form-item>

        <template v-if="searchType === 'appointment'">
          <el-form-item label="状态">
             <el-select v-model="apptFilters.status" placeholder="全部" clearable style="width: 120px">
               <el-option label="已预约" value="已预约" />
               <el-option label="已完成" value="已完成" />
               <el-option label="已取消" value="已取消" />
             </el-select>
          </el-form-item>
          <el-form-item label="日期范围">
             <el-date-picker 
               v-model="apptFilters.dateRange" 
               type="daterange" 
               range-separator="至"
               start-placeholder="开始" 
               end-placeholder="结束"
               value-format="YYYY-MM-DD"
             />
          </el-form-item>
        </template>

        <template v-else>
          <el-form-item label="指标类型">
            <el-select v-model="healthType" style="width: 120px">
              <el-option label="体重" value="Weight" />
              <el-option label="步数" value="Steps" />
              <el-option label="心率" value="HeartRate" />
              <el-option label="血压" value="BloodPressure" />
            </el-select>
          </el-form-item>
          <el-form-item label="月份">
            <el-date-picker v-model="searchMonth" type="month" value-format="YYYY-MM" placeholder="选择月份" />
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div style="margin-top: 20px;">
      <div v-if="searchType === 'health' && healthResults">
         <el-row :gutter="20">
            <el-col :span="12">
              <el-card shadow="hover" style="text-align: center; background: #f0f9eb;">
                 <h3>累计数值</h3>
                 <div style="font-size: 24px; font-weight: bold; color: #67c23a">{{ healthResults.total || 0 }}</div>
              </el-card>
            </el-col>
            <el-col :span="12">
               <el-card shadow="hover" style="text-align: center; background: #e6f7ff;">
                 <h3>平均数值</h3>
                 <div style="font-size: 24px; font-weight: bold; color: #409eff">{{ healthResults.avg ? healthResults.avg.toFixed(2) : 0 }}</div>
              </el-card>
            </el-col>
         </el-row>
      </div>

      <el-table v-else-if="searchType === 'appointment'" :data="apptList" border stripe>
        <el-table-column prop="appointment_date" label="时间">
           <template #default="scope">{{ formatTime(scope.row.appointment_date) }}</template>
        </el-table-column>
        <el-table-column prop="consultation_type" label="类型" />
        <el-table-column prop="status" label="状态">
           <template #default="scope">
             <el-tag :type="scope.row.status === '已取消' ? 'info' : 'success'">{{ scope.row.status }}</el-tag>
           </template>
        </el-table-column>
        <el-table-column prop="reason" label="原因" show-overflow-tooltip />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const searchType = ref('appointment')

// 预约搜索 state
const apptFilters = reactive({
  status: '',
  dateRange: []
})
const apptList = ref([])

// 健康搜索 state
const healthType = ref('Steps')
const searchMonth = ref(new Date().toISOString().slice(0, 7)) // 默认当月
const healthResults = ref(null)

const handleSearch = async () => {
  if (searchType.value === 'appointment') {
    try {
       const params = {}
       if (apptFilters.status) params.status = apptFilters.status
       if (apptFilters.dateRange && apptFilters.dateRange.length === 2) {
          params.start_date = apptFilters.dateRange[0]
          params.end_date = apptFilters.dateRange[1]
       }
       // 对应后端接口：func SearchAppointments
       const res = await request.get(`/users/${userStore.userId}/appointments/search`, { params })
       apptList.value = res || []
    } catch(e) {
       ElMessage.error('搜索失败')
    }
  } else {
    try {
      if (!searchMonth.value) { return ElMessage.warning('请选择月份') }
      // 对应后端接口：func MetricSummary
      const res = await request.get(`/users/${userStore.userId}/metrics/summary`, {
        params: {
          metric_type: healthType.value,
          month: searchMonth.value
        }
      })
      healthResults.value = res // { total, avg }
    } catch(e) {
      ElMessage.error('获取统计失败')
    }
  }
}

const formatTime = (t) => {
  return t ? new Date(t).toLocaleString() : '-'
}
</script>

<style scoped>
.page-container { padding: 20px; }
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
</style>