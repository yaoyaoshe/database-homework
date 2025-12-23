<template>
  <div class="page-container">
    <div class="dashboard-header no-print">
      <div class="header-left">
        <el-button icon="ArrowLeft" circle @click="$router.push('/')" class="back-btn" />
        <div>
          <h2>月度健康报告</h2>
          <p class="subtitle">全面掌握您的健康趋势</p>
        </div>
      </div>
      
      <div class="header-right">
         <el-button type="success" icon="Plus" @click="showRecordDialog = true">记一笔</el-button>
        <el-date-picker 
          v-model="month" 
          type="month" 
          placeholder="选择月份" 
          value-format="YYYY-MM-01" 
          :clearable="false"
          @change="fetchData"
          class="month-picker"
        />
        <el-button type="primary" icon="Printer" @click="printReport">打印/导出PDF</el-button>
      </div>
    </div>

    <div class="print-header only-print">
      <h1>HealthTrack 月度健康报告</h1>
      <p>报告月份: {{ month }}</p>
      <p>生成时间: {{ new Date().toLocaleDateString() }}</p>
    </div>

    <div v-loading="loading" class="dashboard-content">
      <el-empty v-if="!reportData && !loading" description="请选择月份查看报告" />

      <template v-else-if="reportData">
        
        <div class="section-title">
          <el-icon><DataLine /></el-icon> 核心指标概览
        </div>
        <el-row :gutter="20" class="kpi-row">
          <el-col :xs="24" :sm="12" :md="6">
            <el-card shadow="hover" class="kpi-card">
              <div class="stat-value">{{ reportData.total_appointments }}</div>
              <div class="stat-label">本月预约总数</div>
              <div class="stat-footer text-gray">
                取消: {{ reportData.cancelled_appointments }} 次
              </div>
            </el-card>
          </el-col>

          <el-col :xs="24" :sm="12" :md="6">
            <el-card shadow="hover" class="kpi-card">
              <div class="chart-container">
                <el-progress type="circle" :percentage="appointmentRate" :width="80" status="success" />
              </div>
              <div class="stat-label center-text">预约完成率</div>
            </el-card>
          </el-col>

          <el-col :xs="24" :sm="12" :md="6">
            <el-card shadow="hover" class="kpi-card">
              <div class="stat-value text-blue">{{ reportData.total_challenges }}</div>
              <div class="stat-label">参与挑战数</div>
              <div class="stat-footer text-gray">
                已完成: {{ reportData.completed_challenges }} 个
              </div>
            </el-card>
          </el-col>

          <el-col :xs="24" :sm="12" :md="6">
            <el-card shadow="hover" class="kpi-card">
              <div class="chart-container">
                 <el-progress type="circle" :percentage="challengeRate" :width="80" color="#e6a23c" />
              </div>
              <div class="stat-label center-text">挑战达标率</div>
            </el-card>
          </el-col>
        </el-row>

        <div class="section-title" style="margin-top: 30px;">
          <el-icon><TrendCharts /></el-icon> 健康数据摘要
        </div>
        
        <el-row :gutter="20">
            <el-col :span="8" v-if="parsedHealthSummary?.weight_stats">
                <el-card class="metric-card" shadow="never">
                    <template #header><div class="card-head"><el-icon><ScaleToOriginal /></el-icon> 体重记录</div></template>
                    <div class="metric-grid">
                        <div class="m-item"><span>平均</span><strong>{{ parsedHealthSummary.weight_stats.avg }} kg</strong></div>
                        <div class="m-item"><span>最低</span><strong>{{ parsedHealthSummary.weight_stats.min }} kg</strong></div>
                        <div class="m-item"><span>最高</span><strong>{{ parsedHealthSummary.weight_stats.max }} kg</strong></div>
                    </div>
                </el-card>
            </el-col>
             <el-col :span="8" v-if="parsedHealthSummary?.blood_pressure_stats">
                <el-card class="metric-card" shadow="never">
                     <template #header><div class="card-head"><el-icon><Odometer /></el-icon> 血压平均值</div></template>
                    <div class="metric-grid">
                        <div class="m-item"><span>收缩压 (高压)</span><strong>{{ parsedHealthSummary.blood_pressure_stats.systolic_avg }} mmHg</strong></div>
                        <div class="m-item"><span>舒张压 (低压)</span><strong>{{ parsedHealthSummary.blood_pressure_stats.diastolic_avg }} mmHg</strong></div>
                    </div>
                </el-card>
            </el-col>
             <el-col :span="8" v-if="parsedHealthSummary?.step_stats">
                 <el-card class="metric-card" shadow="never">
                    <template #header><div class="card-head"><el-icon><Stopwatch /></el-icon> 运动步数</div></template>
                    <div class="metric-grid">
                        <div class="m-item"><span>本月总计</span><strong>{{ parsedHealthSummary.step_stats.total }} 步</strong></div>
                        <div class="m-item"><span>日均步数</span><strong>{{ parsedHealthSummary.step_stats.daily_avg }} 步</strong></div>
                    </div>
                </el-card>
            </el-col>
        </el-row>

        <div class="section-title" style="margin-top: 30px;">
          <el-icon><FirstAidKit /></el-icon> 综合建议
        </div>
        <el-card shadow="hover" class="advice-card">
          <div class="advice-content">
            <div class="advice-icon">
              <el-icon><ChatDotRound /></el-icon>
            </div>
            <div class="advice-text">
              <h3>HealthTrack 助手建议:</h3>
              <p>{{ reportData.recommendations || '根据您本月的数据，您的健康状况良好。建议继续保持当前的运动和饮食习惯。' }}</p>
            </div>
          </div>
        </el-card>

      </template>
    </div>

    <el-dialog v-model="showRecordDialog" title="录入健康数据" width="400px">
      <el-form label-position="top">
        <el-form-item label="日期时间">
          <el-date-picker v-model="recordForm.date" type="datetime" placeholder="选择日期时间" style="width: 100%" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
        <el-form-item label="指标类型">
          <el-select v-model="recordForm.metric_type" placeholder="请选择" style="width: 100%" @change="handleTypeChange">
             <el-option v-for="opt in metricOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="数值">
           <el-input v-model.number="recordForm.value" type="number" placeholder="请输入数值">
             <template #append>{{ recordForm.unit }}</template>
           </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRecordDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRecord" :loading="submitting">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const month = ref(new Date().toISOString().slice(0, 7) + '-01') // 默认当月1号
const reportData = ref(null)
const loading = ref(false)

// 录入相关
const showRecordDialog = ref(false)
const submitting = ref(false)
const recordForm = reactive({
  date: '',
  metric_type: '',
  value: '',
  unit: ''
})

const metricOptions = [
  { label: '体重', value: '体重', unit: 'kg' },
  { label: '收缩压 (高压)', value: '血压收缩压', unit: 'mmHg' },
  { label: '舒张压 (低压)', value: '血压舒张压', unit: 'mmHg' },
  { label: '步数', value: '步数', unit: '步' },
  { label: '心率', value: '心率', unit: 'bpm' },
  { label: '空腹血糖', value: '血糖空腹', unit: 'mmol/L' },
  { label: '餐后血糖', value: '血糖餐后', unit: 'mmol/L' },
]

const handleTypeChange = (val) => {
  const opt = metricOptions.find(o => o.value === val)
  if (opt) recordForm.unit = opt.unit
}

// 计算属性：预约完成率
const appointmentRate = computed(() => {
  if (!reportData.value || reportData.value.total_appointments === 0) return 0
  const rate = (reportData.value.completed_appointments / reportData.value.total_appointments) * 100
  return Math.round(rate)
})

// 计算属性：挑战达标率
const challengeRate = computed(() => {
  if (!reportData.value || reportData.value.total_challenges === 0) return 0
  const rate = (reportData.value.completed_challenges / reportData.value.total_challenges) * 100
  return Math.round(rate)
})

// 计算属性：解析 health_summary JSON
const parsedHealthSummary = computed(() => {
  if (!reportData.value || !reportData.value.health_summary) return {}
  try {
    const raw = reportData.value.health_summary
    // 如果后端返回已经是对象则直接用，如果是字符串则解析
    return typeof raw === 'string' ? JSON.parse(raw) : raw
  } catch (e) {
    console.error("JSON解析失败", e)
    return {}
  }
})

const fetchData = async () => {
  if (!userStore.userId) return
  
  loading.value = true
  try {
    const res = await request.post('/reports/generate', {
      user_id: parseInt(userStore.userId), 
      month: month.value
    })
    reportData.value = res
  } catch (error) {
    console.error(error)
    ElMessage.error(error.response?.data?.error || '获取报告失败')
    reportData.value = null
  } finally {
    loading.value = false
  }
}

const submitRecord = async () => {
  if(!recordForm.metric_type || !recordForm.value) return ElMessage.warning('请填写完整信息')
  
  submitting.value = true
  try {
    await request.post(`/users/${userStore.userId}/metrics`, {
      metric_type: recordForm.metric_type,
      value: parseFloat(recordForm.value),
      unit: recordForm.unit,
      date: recordForm.date || undefined
    })
    ElMessage.success('记录成功')
    showRecordDialog.value = false
    // 清空表单
    recordForm.value = ''
    recordForm.metric_type = ''
    // 重新获取报告以更新数据
    fetchData()
  } catch(e) {
    ElMessage.error('保存失败')
  } finally {
    submitting.value = false
  }
}

const printReport = () => {
  window.print()
}

onMounted(() => {
    // 设置默认时间为当前
    const now = new Date()
    const isoString = new Date(now.getTime() - (now.getTimezoneOffset() * 60000)).toISOString()
    recordForm.date = isoString.slice(0, 19).replace('T', ' ')
    fetchData()
})
</script>

<style scoped>
.page-container { padding: 24px; max-width: 1200px; margin: 0 auto; background-color: #f5f7fa; min-height: 100vh; }

/* 头部样式 */
.dashboard-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; background: white; padding: 20px; border-radius: 12px; box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05); }
.header-left { display: flex; align-items: center; gap: 15px; }
.header-left h2 { margin: 0; color: #303133; }
.subtitle { margin: 4px 0 0 0; color: #909399; font-size: 14px; }
.header-right { display: flex; gap: 10px; }

/* 标题样式 */
.section-title { font-size: 18px; font-weight: 600; color: #303133; margin-bottom: 15px; display: flex; align-items: center; gap: 8px; }

/* KPI 卡片 */
.kpi-row { margin-bottom: 10px; }
.kpi-card { height: 160px; display: flex; flex-direction: column; justify-content: center; position: relative; border: none; border-radius: 12px; }
.stat-value { font-size: 36px; font-weight: bold; color: #303133; margin-bottom: 5px; text-align: center; }
.stat-label { font-size: 14px; color: #909399; text-align: center; }
.stat-footer { font-size: 12px; margin-top: 10px; text-align: center; }
.text-gray { color: #909399; }
.text-blue { color: #409eff; }
.chart-container { display: flex; justify-content: center; margin-bottom: 10px; }
.center-text { text-align: center; }

/* 详情数据卡片 */
.metric-card { border-radius: 8px; height: 100%; }
.card-head { display: flex; align-items: center; gap: 8px; font-weight: bold; color: #606266; }
.metric-grid { display: flex; flex-direction: column; gap: 12px; padding: 10px 0; }
.m-item { display: flex; justify-content: space-between; border-bottom: 1px dashed #eee; padding-bottom: 8px; }
.m-item span { color: #909399; font-size: 14px; }
.m-item strong { color: #303133; font-size: 16px; }

/* 建议卡片 */
.advice-card { border-radius: 12px; background: linear-gradient(135deg, #f0f9eb 0%, #e1f3d8 100%); border: 1px solid #c2e7b0; }
.advice-content { display: flex; gap: 20px; align-items: flex-start; }
.advice-icon { font-size: 40px; color: #67c23a; margin-top: 5px; }
.advice-text h3 { margin: 0 0 10px 0; color: #2c3e50; font-size: 16px; }
.advice-text p { margin: 0; color: #5e6d82; line-height: 1.6; }

/* 打印样式控制 */
.only-print { display: none; }

@media print {
  .no-print { display: none !important; }
  .only-print { display: block; margin-bottom: 30px; text-align: center; }
  .page-container { background: white; padding: 0; }
  .el-card { box-shadow: none !important; border: 1px solid #ccc !important; page-break-inside: avoid; }
  .kpi-card { height: auto; padding: 15px; border: 1px solid #eee; }
}
</style>