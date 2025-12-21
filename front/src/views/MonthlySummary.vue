<template>
  <div class="page-container">
    <div class="header">
      <h2>月度健康摘要</h2>
      <div class="controls">
        <el-date-picker v-model="month" type="month" placeholder="选择月份" value-format="YYYY-MM-01" />
        <el-button type="primary" @click="generateReport">生成/刷新报告</el-button>
      </div>
    </div>

    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="12">
        <el-card>
          <div id="chart-weight" style="height: 300px;"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <div id="chart-steps" style="height: 300px;"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-alert 
      title="数据说明" 
      type="info" 
      description="后端接口仅支持生成报告，暂无获取报告JSON详情的API。上方图表为演示数据。"
      show-icon 
      style="margin-top:20px"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const month = ref('2024-03-01')

const generateReport = async () => {
  if(!month.value) return
  try {
    const res = await request.post('/reports/generate', {
      user_id: parseInt(userStore.userId),
      month: month.value
    })
    ElMessage.success(res.message)
    // 这里应该调用GET接口刷新数据，但后端未提供，故不操作
    initCharts() // 重绘演示图表
  } catch(e) {}
}

// 模拟图表初始化
const initCharts = () => {
  const chartWeight = echarts.init(document.getElementById('chart-weight'))
  chartWeight.setOption({
    title: { text: '体重变化 (kg)' },
    xAxis: { type: 'category', data: ['W1', 'W2', 'W3', 'W4'] },
    yAxis: { type: 'value', min: 60 },
    series: [{ data: [72.5, 71.8, 71.2, 70.8], type: 'line' }]
  })

  const chartSteps = echarts.init(document.getElementById('chart-steps'))
  chartSteps.setOption({
    title: { text: '每日步数' },
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
    yAxis: { type: 'value' },
    series: [{ data: [8000, 9200, 10500, 7800, 11000, 12500, 9000], type: 'bar' }]
  })
}

onMounted(() => {
  initCharts()
})
</script>

<style scoped>
.page-container { padding: 20px; }
.header { display: flex; justify-content: space-between; align-items: center; }
</style>