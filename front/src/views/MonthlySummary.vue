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

    <div v-loading="loading" class="dashboard-content">
      <el-empty v-if="!reportData && !loading" description="请选择月份查看报告" />

      <template v-else-if="reportData">
        <el-row :gutter="20" class="kpi-row">
          <el-col :xs="24" :sm="12">
            <el-card shadow="hover" class="kpi-card appt-card">
              <div class="kpi-content">
                <div class="icon-wrapper"><el-icon><Calendar /></el-icon></div>
                <div class="stat-info">
                  <div class="label">本月预约总数</div>
                  <el-statistic :value="reportData.total_appointments" />
                  <div class="sub-stat">
                    已完成: {{ reportData.completed_appointments }} | 取消: {{ reportData.cancelled_appointments }}
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :xs="24" :sm="12">
            <el-card shadow="hover" class="kpi-card challenge-card">
              <div class="kpi-content">
                <div class="icon-wrapper"><el-icon><Trophy /></el-icon></div>
                <div class="stat-info">
                  <div class="label">参与挑战数</div>
                  <el-statistic :value="reportData.total_challenges" />
                  <div class="sub-stat">
                    已达标: {{ reportData.completed_challenges }}
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-card shadow="never" class="advice-card" style="margin-top: 24px;">
          <template #header>
            <div class="card-header">
              <span><el-icon><FirstAidKit /></el-icon> 系统生成建议</span>
            </div>
          </template>
          <div class="advice-content">
            <el-alert
              :title="reportData.recommendations || '暂无特别建议，请继续保持。'"
              type="success"
              :closable="false"
              show-icon
            />
            
            <div v-if="reportData.health_summary" style="margin-top: 15px; background: #f8f9fa; padding: 15px; border-radius: 8px;">
               <h4>健康数据摘要:</h4>
               <pre style="font-family: inherit; color: #666;">{{ parseSummary(reportData.health_summary) }}</pre>
            </div>
          </div>
        </el-card>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const month = ref(new Date().toISOString().slice(0, 7) + '-01')
const reportData = ref(null)
const loading = ref(false)

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
    ElMessage.error('获取报告失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const parseSummary = (raw) => {
  if (!raw) return '无数据'
  try {
    // 如果是 JSON 对象，格式化显示
    if (typeof raw === 'object') return JSON.stringify(raw, null, 2)
    // 如果是字符串，尝试解析
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

const printReport = () => window.print()

onMounted(fetchData)
</script>

<style scoped>
.page-container { padding: 24px; max-width: 1200px; margin: 0 auto; }
.dashboard-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
.header-left h2 { margin: 0; }
.subtitle { margin: 4px 0 0 0; color: #909399; font-size: 14px; }
.kpi-row { margin-bottom: 24px; }
.kpi-card { height: 100%; border-radius: 12px; }
.kpi-content { display: flex; align-items: center; padding: 10px 0; }
.icon-wrapper { width: 56px; height: 56px; border-radius: 16px; display: flex; justify-content: center; align-items: center; margin-right: 20px; font-size: 28px; }
.stat-info .label { font-size: 14px; color: #606266; margin-bottom: 4px; }
.sub-stat { font-size: 12px; color: #999; margin-top: 5px; }

.appt-card .icon-wrapper { background: rgba(64, 158, 255, 0.1); color: #409eff; }
.challenge-card .icon-wrapper { background: rgba(230, 162, 60, 0.1); color: #e6a23c; }

@media print {
  .dashboard-header .header-right, .back-btn { display: none; }
  .page-container { padding: 0; }
  .el-card { box-shadow: none !important; border: 1px solid #ddd !important; }
}
</style>