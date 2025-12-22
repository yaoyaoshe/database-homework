<template>
  <div>
    <h2>健康挑战</h2>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="参与中的挑战" name="active">
        <el-table :data="myChallenges" border>
          <el-table-column prop="challenge_name" label="挑战名称" />
          <el-table-column prop="current_progress" label="当前进度" />
          <el-table-column prop="progress_unit" label="单位" width="80" />
          <el-table-column prop="status" label="状态" width="100" />
          <el-table-column label="操作">
            <template #default="{ row }">
              <el-button size="small" type="primary" @click="openCheckin(row)">打卡</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="checkinVisible" title="每日打卡">
      <el-form :model="checkinForm">
        <el-form-item label="日期">
           <el-date-picker v-model="checkinForm.date" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" />
        </el-form-item>
        <el-form-item label="进度值">
          <el-input-number v-model="checkinForm.progress_value" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="checkinForm.notes" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="submitCheckin">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const activeTab = ref('active')
const myChallenges = ref([])
const checkinVisible = ref(false)
const currentChallengeId = ref(0)
const checkinForm = ref({ date: '', progress_value: 0, notes: '' })

const fetchMyChallenges = async () => {
  // 注意：此处使用了补充文档中的 summary 接口，id传0或任意值因为后端似乎未严格校验path id, 主要是query param
  // 假设后端接口逻辑是 GET /challenges/:id/summary?user_id=...
  // 我们这里随便传个id=1，重点是 query
  try {
     const res = await request.get(`/challenges/0/summary?user_id=${userStore.userId}`)
     // 后端如果返回单个对象，转为数组；如果返回数组直接用
     myChallenges.value = Array.isArray(res) ? res : [res]
  } catch (e) {
      console.log(e)
  }
}

const openCheckin = (row) => {
  currentChallengeId.value = row.challenge_id
  checkinForm.value = { 
    date: new Date().toISOString().split('T')[0], 
    progress_value: 0, 
    notes: '',
    progress_unit: row.progress_unit
  }
  checkinVisible.value = true
}

const submitCheckin = async () => {
  await request.post(`/challenges/${currentChallengeId.value}/checkin`, {
    user_id: parseInt(userStore.userId),
    date: checkinForm.value.date,
    progress_value: Number(checkinForm.value.progress_value),
    progress_unit: checkinForm.value.progress_unit,
    notes: checkinForm.value.notes
  })
  ElMessage.success('打卡成功')
  checkinVisible.value = false
  fetchMyChallenges() // 刷新进度
}

onMounted(fetchMyChallenges)
</script>