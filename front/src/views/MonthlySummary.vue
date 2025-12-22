<template>
  <div class="page-container">
    <div class="dashboard-header">
      <div class="header-left">
        <el-button icon="ArrowLeft" circle @click="$router.push('/')" class="back-btn" />
        <div>
          <h2>月度健康报告</h2>
          <p class="subtitle">全面掌握您的健康趋势</p>
        </div>
      </div>
      
      <div class="header-right">
        <el-date-picker 
          v-model="month" 
          type="month" 
          placeholder="选择月份" 
          value-format="YYYY-MM-01" 
          :clearable="false"
          @change="fetchData"
          class="month-picker"
        />
        <el-button type="primary" icon="Download" @click="printReport">导出报告</el-button>
      </div>
    </div>

    <div v-loading="loading" element-loading-text="正在分析健康数据..." class="dashboard-content">
      
      <el-empty v-if="!reportData && !loading" description="请选择月份查看报告" />

      <template v-else-if="reportData">
        <el-row :gutter="20" class="kpi-row">
          <el-col :xs="24" :sm="8">
            <el-card shadow="hover" class="kpi-card steps-card">
              <div class="kpi-content">
                <div class="icon-wrapper">
                  <el-icon><collection-tag /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="label">本月总步数</div>
                  <el-statistic :value="reportData.total_steps" group-separator="," />
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :xs="24" :sm="8">
            <el-card shadow="hover" class="kpi-card weight-card">
              <div class="kpi-content">
                <div class="icon-wrapper">
                  <el-icon><odometer /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="label">平均体重</div>
                  <el-statistic :value="reportData.weight_stats?.avg" :precision="1" suffix="kg" />
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :xs="24" :sm="8">
            <el-card shadow="hover" class="kpi-card range-card">
              <div class="kpi-content">
                <div class="icon-wrapper">
                  <el-icon><trend-charts /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="label">体重波动范围</div>
                  <div class="range-value">
                    {{ reportData.weight_stats?.min || '-' }} ~ {{ reportData.weight_stats?.max || '-' }} <span class="unit">kg</span>
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-row :gutter="20" style="margin-top: 24px;">
          <el-col :span="24">
            <el-card shadow="never" class="chart-card">
              <template #header>
                <div class="card-header">
                  <span><el-icon><Histogram /></el-icon> 体重变化趋势</span>
                </div>
              </template>
              <div id="chart-weight" style="height: 400px; width: 100%;"></div>
            </el-card>
          </el-col>
        </el-row>

        <el-card shadow="never" class="advice-card" style="margin-top: 24px;">
          <template #header>
            <div class="card-header">
              <span><el-icon><FirstAidKit /></el-icon> 健康小贴士</span>
            </div>
          </template>
          <div class="advice-content">
            <el-alert
              title="保持活跃"
              type="success"
              :closable="false"
              description="您本月的步数表现不错，建议继续保持每日散步的习惯，有助于心肺功能提升。"
              show-icon
              style="margin-bottom: 12px;"
            />
            <el-alert
              title="体重管理"
              type="info"
              :closable="false"
              description="定期监测体重有助于了解身体状况。建议固定在早晨空腹时测量，数据更准确。"
              show-icon
            />
          </div>
        </el-card>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onUnmounted } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { 
  ArrowLeft, Download, CollectionTag, Odometer, 
  TrendCharts, Histogram, FirstAidKit 
} from '@element-plus/icons-vue'

const userStore = useUserStore()
const month = ref(new Date().toISOString().slice(0, 7) + '-01') // 默认当月
const reportData = ref(null)
const loading = ref(false)
let myChart = null

const fetchData = async () => {
  if (!userStore.userId) {
    ElMessage.warning('用户信息未加载')
    return
  }
  
  loading.value = true
  try {
    // 1. 获取统计摘要
    const res = await request.post('/reports/generate', {
      user_id: parseInt(userStore.userId), 
      month: month.value
    })
    reportData.value = res

    // 2. 获取详细数据用于绘图
    const rawData = await request.get(`/health-data/${userStore.userId}?type=Weight`)
    
    // 等待 DOM 更新后渲染图表
    await nextTick()
    renderChart(rawData)
  } catch (error) {
    console.error(error)
    ElMessage.error('获取报告数据失败')
  } finally {
    loading.value = false
  }
}

const renderChart = (data) => {
  const chartDom = document.getElementById('chart-weight')
  if (!chartDom) return

  // 销毁旧实例防止内存泄漏
  if (myChart) myChart.dispose()
  
  myChart = echarts.init(chartDom)
  
  // 数据处理：按时间正序排列
  const sorted = (data || [])
    .filter(item => item.recorded_at.startsWith(month.value.slice(0, 7))) // 仅筛选选中月份的数据
    .sort((a, b) => new Date(a.recorded_at) - new Date(b.recorded_at))
    
  const dates = sorted.map(i => i.recorded_at.substring(8, 10) + '日') // 显示 "01日"
  const values = sorted.map(i => i.data_value)

  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      borderColor: '#eee',
      textStyle: { color: '#333' }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates.length ? dates : ['无数据'],
      axisLine: { lineStyle: { color: '#ccc' } },
      axisLabel: { color: '#666' }
    },
    yAxis: {
      type: 'value',
      scale: true, // 不从0开始，突显变化
      splitLine: { lineStyle: { type: 'dashed', color: '#eee' } },
      name: '体重 (kg)',
      nameTextStyle: { align: 'right' }
    },
    series: [{
      name: '体重',
      data: values.length ? values : [0],
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      itemStyle: { color: '#409EFF' },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
          { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
        ])
      }
    }]
  }

  myChart.setOption(option)
}

const printReport = () => {
  window.print()
}

// 监听窗口大小变化调整图表
const handleResize = () => myChart && myChart.resize()

onMounted(() => {
  fetchData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (myChart) myChart.dispose()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
}

/* 头部样式 */
.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.header-left h2 {
  margin: 0;
  font-size: 24px;
  color: #1a1a1a;
}

.subtitle {
  margin: 4px 0 0 0;
  color: #909399;
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* KPI 卡片样式 */
.kpi-row {
  margin-bottom: 24px;
}

.kpi-card {
  height: 100%;
  border: none;
  border-radius: 12px;
  transition: transform 0.3s;
}

.kpi-card:hover {
  transform: translateY(-4px);
}

.kpi-content {
  display: flex;
  align-items: center;
  padding: 10px 0;
}

.icon-wrapper {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 20px;
  font-size: 28px;
}

.stat-info {
  flex: 1;
}

.stat-info .label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 4px;
}

/* 特定卡片配色 */
.steps-card .icon-wrapper {
  background: rgba(103, 194, 58, 0.1);
  color: #67c23a;
}

.weight-card .icon-wrapper {
  background: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.range-card .icon-wrapper {
  background: rgba(144, 147, 153, 0.1);
  color: #909399;
}

.range-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.range-value .unit {
  font-size: 14px;
  font-weight: normal;
  color: #909399;
  margin-left: 4px;
}

/* 图表与通用卡片 */
.chart-card, .advice-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.card-header {
  display: flex;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
}

.card-header .el-icon {
  margin-right: 8px;
  color: #409EFF;
}

/* 打印样式适配 */
@media print {
  .dashboard-header .header-right, .back-btn {
    display: none;
  }
  .page-container {
    padding: 0;
  }
  .el-card {
    box-shadow: none !important;
    border: 1px solid #ddd !important;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .header-right {
    width: 100%;
    justify-content: space-between;
  }
  .kpi-row .el-col {
    margin-bottom: 16px;
  }
}
</style>