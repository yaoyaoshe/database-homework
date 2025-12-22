<template>
  <div>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="热门挑战" name="popular">
        <el-table :data="popularList" stripe>
          <el-table-column prop="challenge_name" label="挑战名称" />
          <el-table-column prop="challenge_type" label="类型" width="100" />
          <el-table-column prop="participant_count" label="参与人数" width="100" align="center" />
          <el-table-column label="操作" width="120">
            <template #default="scope">
              <el-button size="small" type="success" @click="handleJoin(scope.row.challenge_id)">加入</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="发起挑战" name="create">
        <div class="create-container">
          <el-form :model="createForm" label-width="100px">
            <el-form-item label="挑战名称">
              <el-input v-model="createForm.challenge_name" placeholder="例如：30天减脂计划" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="createForm.description" type="textarea" />
            </el-form-item>
            <el-form-item label="类型">
              <el-select v-model="createForm.challenge_type">
                <el-option label="运动" value="运动" />
                <el-option label="减重" value="减重" />
                <el-option label="饮食" value="饮食" />
                <el-option label="睡眠" value="睡眠" />
              </el-select>
            </el-form-item>
            <el-form-item label="目标指标">
              <el-select v-model="createForm.target_metric">
                <el-option label="步数" value="步数" />
                <el-option label="体重" value="体重" />
                <el-option label="运动时长" value="运动时长" />
              </el-select>
            </el-form-item>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="目标值">
                  <el-input-number v-model="createForm.target_value" :min="1" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="单位">
                  <el-input v-model="createForm.target_unit" placeholder="如：步、kg" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="公开挑战">
              <el-switch v-model="createForm.is_public" />
            </el-form-item>
            <el-button type="primary" @click="submitChallenge">创建挑战</el-button>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane label="每日打卡" name="checkin">
        <div class="checkin-container">
          <el-form :model="checkinForm" label-width="80px">
            <el-form-item label="挑战ID">
              <el-input v-model.number="checkinForm.challenge_id" placeholder="输入已参加的挑战ID" />
            </el-form-item>
            <el-form-item label="日期">
              <el-date-picker v-model="checkinForm.date" value-format="YYYY-MM-DD" placeholder="打卡日期" />
            </el-form-item>
            <el-form-item label="进度值">
              <el-input-number v-model="checkinForm.progress_value" />
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model="checkinForm.notes" />
            </el-form-item>
            <el-button type="primary" @click="submitCheckin">提交打卡</el-button>
            <el-button @click="loadProgress">查看历史记录</el-button>
          </el-form>

          <el-divider v-if="progressRecords.length" content-position="left">历史记录</el-divider>
          <el-table v-if="progressRecords.length" :data="progressRecords" size="small">
            <el-table-column prop="progress_date" label="日期">
               <template #default="{row}">{{ row.progress_date.substring(0,10) }}</template>
            </el-table-column>
            <el-table-column prop="progress_value" label="数值" />
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
// 注意引入新加的 createChallenge
import { getPopularChallenges, joinChallenge, checkinChallenge, getDailyProgress, createChallenge } from '@/api/all'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const activeTab = ref('popular')
const popularList = ref([])
const progressRecords = ref([])
const dateRange = ref([])

// === 创建挑战表单数据 ===
const createForm = ref({
  challenge_name: '',
  description: '',
  challenge_type: '运动',
  target_metric: '步数',
  target_value: 10000,
  target_unit: '步',
  is_public: true,
  max_participants: 50
})

const checkinForm = ref({
  challenge_id: null,
  date: new Date().toISOString().split('T')[0],
  progress_value: 0,
  notes: ''
})

// ... (保留 loadPopular, handleJoin, submitCheckin, loadProgress 等原有方法) ...

const loadPopular = async () => {
  popularList.value = await getPopularChallenges(10)
}

const handleJoin = async (id) => {
  await joinChallenge(id, { user_id: userStore.userInfo.user_id })
  ElMessage.success('加入成功！快去打卡吧')
}

const submitCheckin = async () => {
  if(!checkinForm.value.challenge_id) return
  await checkinChallenge(checkinForm.value.challenge_id, {
    user_id: userStore.userInfo.user_id,
    date: checkinForm.value.date,
    progress_value: checkinForm.value.progress_value,
    notes: checkinForm.value.notes
  })
  ElMessage.success('打卡成功')
  loadProgress()
}

const loadProgress = async () => {
  if(!checkinForm.value.challenge_id) return
  progressRecords.value = await getDailyProgress(checkinForm.value.challenge_id, userStore.userInfo.user_id)
}

// === 新增：提交挑战 ===
const submitChallenge = async () => {
  if (!dateRange.value || dateRange.value.length < 2) return ElMessage.warning('请选择时间范围')
  
  const payload = {
    ...createForm.value,
    creator_id: userStore.userInfo.user_id,
    start_date: dateRange.value[0],
    end_date: dateRange.value[1]
  }
  
  await createChallenge(payload)
  ElMessage.success('挑战创建成功！')
  createForm.value = { ...createForm.value, challenge_name: '' } // reset simple fields
  activeTab.value = 'popular'
  loadPopular()
}

onMounted(loadPopular)
</script>

<style scoped>
.checkin-container, .create-container { max-width: 600px; margin-top: 20px; }
</style>