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
import { getPopularChallenges, joinChallenge, checkinChallenge, getDailyProgress } from '@/api/all'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const activeTab = ref('popular')
const popularList = ref([])
const progressRecords = ref([])

const checkinForm = ref({
  challenge_id: null,
  date: new Date().toISOString().split('T')[0],
  progress_value: 0,
  notes: ''
})

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

onMounted(loadPopular)
</script>

<style scoped>
.checkin-container { max-width: 500px; margin-top: 20px; }
</style>