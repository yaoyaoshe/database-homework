<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>健康挑战中心</h2>
    </div>

    <el-tabs v-model="activeTab" type="border-card" class="challenge-tabs">
      
      <el-tab-pane label="我参加的挑战" name="my_challenges">
        <el-table :data="myChallenges" style="width: 100%" stripe v-loading="loading">
          <el-table-column prop="challenge_name" label="挑战名称" />
          <el-table-column prop="challenge_type" label="类型" width="100" />
          <el-table-column label="进度" width="200">
             <template #default="scope">
               <el-progress 
                 :percentage="Math.min(Number(scope.row.current_progress || 0), 100)" 
                 :status="scope.row.status === '已完成' ? 'success' : ''"
               />
               <div class="progress-text">
                  目标: {{ scope.row.target_value }} {{ scope.row.target_unit }}
               </div>
             </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
             <template #default="scope">
               <el-tag :type="getStatusColor(scope.row.status)">{{ scope.row.status }}</el-tag>
             </template>
          </el-table-column>
          <el-table-column label="加入时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.joined_at) }}</template>
          </el-table-column>
        </el-table>
        
        <el-empty v-if="!loading && myChallenges.length === 0" description="您还没有参加任何挑战，快去创建或加入一个吧！" />
      </el-tab-pane>

      <el-tab-pane label="创建新挑战" name="create">
        <div class="form-wrapper">
          <el-form :model="form" label-width="100px" size="large">
            <el-form-item label="挑战名称">
              <el-input v-model="form.challenge_name" placeholder="例如：30天减重计划" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="form.description" type="textarea" :rows="3" placeholder="描述挑战的目标和规则" />
            </el-form-item>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="类型">
                  <el-select v-model="form.challenge_type" style="width: 100%">
                    <el-option label="减重" value="减重" />
                    <el-option label="运动" value="运动" />
                    <el-option label="睡眠" value="睡眠" />
                    <el-option label="健康习惯" value="健康习惯" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="时间范围">
                  <el-date-picker
                    v-model="dateRange"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始"
                    end-placeholder="结束"
                    style="width: 100%"
                    value-format="YYYY-MM-DD"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="8">
                <el-form-item label="目标数值">
                  <el-input-number v-model="form.target_value" :min="1" style="width: 100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="单位">
                  <el-input v-model="form.target_unit" placeholder="如：kg" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                 <el-form-item label="最大人数">
                  <el-input-number v-model="form.max_participants" :min="1" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <div style="margin-top: 20px; text-align: right;">
               <el-button @click="activeTab = 'my_challenges'">取消</el-button>
               <el-button type="primary" @click="createChallenge">立即创建挑战</el-button>
            </div>
          </el-form>
        </div>
      </el-tab-pane>

    </el-tabs>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const activeTab = ref('my_challenges') // 默认显示我的挑战
const myChallenges = ref([])
const loading = ref(false)

const dateRange = ref([])
const form = reactive({
  challenge_name: '',
  description: '',
  challenge_type: '运动',
  target_metric: '',
  target_value: 0,
  target_unit: '',
  max_participants: 100,
  is_public: true
})

// 获取我参与的挑战
const fetchMyChallenges = async () => {
  loading.value = true
  try {
    // 假设后端新增了 GET /users/{id}/challenges 接口
    // 如果后端尚未实现，这里可能会报错，建议同步更新后端 handler.go
    const res = await request.get(`/users/${userStore.userId}/challenges`)
    myChallenges.value = res || []
  } catch(e) {
    // 降级处理：如果不通，暂显示空
    console.warn("获取挑战列表失败，可能后端接口未就绪")
  } finally {
    loading.value = false
  }
}

const createChallenge = async () => {
  if (!dateRange.value || dateRange.value.length < 2) return ElMessage.warning('请选择日期')
  if (!form.challenge_name) return ElMessage.warning('请输入挑战名称')
  
  try {
    await request.post('/challenges', {
      creator_id: parseInt(userStore.userId),
      ...form,
      start_date: dateRange.value[0],
      end_date: dateRange.value[1]
    })
    ElMessage.success('挑战创建成功！')
    
    // 重置表单并切换回列表
    form.challenge_name = ''
    form.description = ''
    dateRange.value = []
    activeTab.value = 'my_challenges'
    fetchMyChallenges() // 刷新列表
    
  } catch(e) {}
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}

const getStatusColor = (status) => {
  if (status === '进行中') return 'success'
  if (status === '已结束') return 'info'
  return 'primary'
}

onMounted(() => {
  fetchMyChallenges()
})
</script>

<style scoped>
.page-container { padding: 20px; }
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 24px; color: #333; }

.challenge-tabs {
  min-height: 500px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05);
}

.form-wrapper {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.progress-text {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>