<template>
  <div class="page-container">
    <h2>平台摘要统计</h2>
    
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card header="热门挑战 (Top 5)">
          <el-table :data="popularChallenges" style="width: 100%">
            <el-table-column prop="challenge_name" label="挑战名" />
            <el-table-column prop="participant_count" label="参与人数" />
            <el-table-column prop="avg_progress" label="平均进度">
               <template #default="scope">{{ scope.row.avg_progress }}%</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card header="最活跃用户 (Top 10)">
          <el-table :data="activeUsers" style="width: 100%">
            <el-table-column prop="name" label="姓名" />
            <el-table-column prop="activity_score" label="活跃分" />
            <el-table-column prop="health_record_count" label="记录数" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const popularChallenges = ref([])
const activeUsers = ref([])

onMounted(async () => {
  try {
    const res1 = await request.get('/stats/popular_challenges?limit=5')
    popularChallenges.value = res1 || []
    
    const res2 = await request.get('/stats/active_users?limit=10')
    activeUsers.value = res2 || []
  } catch(e) {}
})
</script>

<style scoped>
.page-container { padding: 20px; }
</style>