<template>
  <div>
    <div class="header-bar">
      <h2>健康月报</h2>
      <div class="controls">
        <el-date-picker v-model="month" type="month" placeholder="选择月份" value-format="YYYY-MM-01" />
        <el-button type="primary" @click="generateReport" style="margin-left: 10px;">生成/查看报告</el-button>
      </div>
    </div>

    <div v-if="reportData" class="report-content">
      <el-alert :title="reportData.recommendations" type="success" :closable="false" show-icon style="margin-bottom: 20px;" />
      
      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic title="总预约数" :value="reportData.total_appointments" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="已完成预约" :value="reportData.completed_appointments" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="参与挑战" :value="reportData.total_challenges" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="完成挑战" :value="reportData.completed_challenges" />
        </el-col>
      </el-row>

      <el-divider>健康数据摘要</el-divider>
      <pre class="json-box">{{ formatSummary(reportData.health_summary) }}</pre>
    </div>
    
    <el-empty v-else description="请选择月份并点击生成查看报告" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import request from '../api/request'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const month = ref('')
const reportData = ref(null)

const generateReport = async () => {
  if (!month.value) return
  const res = await request.post('/reports/generate', {
    user_id: parseInt(userStore.userId),
    month: month.value
  })
  reportData.value = res
}

const formatSummary = (json) => {
  try {
    return JSON.stringify(json, null, 2)
  } catch (e) {
    return json
  }
}
</script>

<style scoped>
.header-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.json-box { background: #f4f4f5; padding: 15px; border-radius: 4px; font-family: monospace; }
</style>