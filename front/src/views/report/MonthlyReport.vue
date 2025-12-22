<template>
  <div>
    <h2>健康月报</h2>
    <div class="search-box">
      <el-date-picker 
        v-model="month" 
        type="month" 
        placeholder="选择月份" 
        format="YYYY-MM"
        value-format="YYYY-MM-01"
      />
      <el-button type="primary" @click="fetchReport" style="margin-left: 10px">生成/查看报告</el-button>
    </div>

    <div v-if="report" class="report-content">
      <el-descriptions title="概览数据" border>
        <el-descriptions-item label="总预约">{{ report.total_appointments }}</el-descriptions-item>
        <el-descriptions-item label="已完成预约">{{ report.completed_appointments }}</el-descriptions-item>
        <el-descriptions-item label="总挑战">{{ report.total_challenges }}</el-descriptions-item>
        <el-descriptions-item label="完成挑战">{{ report.completed_challenges }}</el-descriptions-item>
      </el-descriptions>

      <el-card class="mt-20">
        <template #header><strong>健康建议</strong></template>
        <p>{{ report.recommendations }}</p>
      </el-card>

      <el-card class="mt-20" v-if="report.health_summary">
        <template #header><strong>详细指标 (JSON)</strong></template>
        <pre>{{ JSON.stringify(report.health_summary, null, 2) }}</pre>
      </el-card>
    </div>
    
    <el-empty v-else description="暂无数据，请选择月份生成" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { generateReport } from '@/api/all'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const month = ref('')
const report = ref(null)

const fetchReport = async () => {
  if(!month.value) return
  const res = await generateReport({
    user_id: userStore.userInfo.user_id,
    month: month.value
  })
  report.value = res
}
</script>

<style scoped>
.search-box { margin-bottom: 30px; }
.mt-20 { margin-top: 20px; }
pre { background: #f4f4f5; padding: 15px; border-radius: 4px; }
</style>