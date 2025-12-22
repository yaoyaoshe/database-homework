<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>个人信息</template>
          <p><strong>姓名:</strong> {{ info.name }}</p>
          <p><strong>Health ID:</strong> {{ info.health_id }}</p>
          <p><strong>性别:</strong> {{ info.gender }}</p>
        </el-card>
      </el-col>
      <el-col :span="16">
         <el-card shadow="hover">
           <template #header>平台热门挑战</template>
           <el-table :data="popularChallenges" style="width: 100%">
             <el-table-column prop="challenge_name" label="名称" />
             <el-table-column prop="participant_count" label="参与人数" />
             <el-table-column prop="avg_progress" label="平均进度" />
           </el-table>
         </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const info = ref({})
const popularChallenges = ref([])

onMounted(async () => {
  // 获取用户详情
  info.value = await request.get(`/users/${userStore.userId}`)
  // 获取热门统计
  popularChallenges.value = await request.get('/stats/popular_challenges?limit=5')
})
</script>