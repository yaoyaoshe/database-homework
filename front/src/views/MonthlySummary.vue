<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>月度健康摘要</h2>
    </div>
    
    <div class="controls">
      <el-date-picker v-model="month" type="month" placeholder="选择月份" value-format="YYYY-MM-01" />
      <el-button type="primary" @click="fetchData" style="margin-left: 10px;">生成报告</el-button>
    </div>

    <el-row :gutter="20" style="margin-top: 20px;" v-if="reportData">
      <el-col :span="8">
        <el-card shadow="hover" class="stat-box">
          <h3>本月总步数</h3>
          <div class="big-num">{{ reportData.total_steps || 0 }}</div>
          <div class="unit">步</div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="hover" class="stat-box">
          <h3>体重统计 (kg)</h3>
          <el-row style="text-align: center; width: 100%;">
            <el-col :span="8"><div>平均: {{ reportData.weight_stats?.avg?.toFixed(1) || '-' }}</div></el-col>
            <el-col :span="8"><div>最低: {{ reportData.weight_stats?.min || '-' }}</div></el-col>
            <el-col :span="8"><div>最高: {{ reportData.weight_stats?.max || '-' }}</div></el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="hover" style="margin-top: 20px;">
      <div id="chart-weight" style="height: 300px;"></div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import * as echarts from 'echarts'

const userStore = useUserStore()
const month = ref('2024-03-01') // 默认值
const reportData = ref(null)

const fetchData = async () => {
  if(!month.value) return
  try {
    // 1. 获取统计摘要 (Aggregated Stats)
    const res = await request.post('/reports/generate', {
      user_id: parseInt(userStore.userId), month: month.value
    })
    reportData.value = res

    // 2. 获取详细数据用于绘图 (Raw Data)
    const rawData = await request.get(`/health-data/${userStore.userId}?type=Weight`)
    renderChart(rawData)
  } catch(e) {}
}

const renderChart = (data) => {
  const chartDom = document.getElementById('chart-weight')
  if(!chartDom) return
  const myChart = echarts.init(chartDom)
  
  // 简单处理数据，倒序排列
  const sorted = (data || []).reverse()
  const dates = sorted.map(i => i.recorded_at.substring(0, 10))
  const values = sorted.map(i => i.data_value)

  myChart.setOption({
    title: { text: '体重趋势图' },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: dates },
    yAxis: { type: 'value', scale: true },
    series: [{ data: values, type: 'line', smooth: true }]
  })
}

onMounted(fetchData)
</script>

<style scoped>
.page-container { padding: 20px; }
.page-header { display: flex; align-items: center; margin-bottom: 20px; }
.controls { background: #fff; padding: 15px; border-radius: 8px; }
.stat-box { text-align: center; height: 120px; display: flex; flex-direction: column; justify-content: center; }
.big-num { font-size: 32px; font-weight: bold; color: #007bff; margin: 10px 0; }
</style>