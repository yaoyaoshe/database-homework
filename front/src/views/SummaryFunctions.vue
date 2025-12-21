<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>统计概览</h2>
    </div>

    <el-row :gutter="24">
      <el-col :lg="14" :md="24" style="margin-bottom: 20px;">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">
                <el-icon class="icon-trophy"><Trophy /></el-icon> 热门挑战排行榜
              </span>
              <el-tag type="danger" effect="dark" size="small">HOT</el-tag>
            </div>
          </template>
          
          <el-table :data="popularChallenges" stripe style="width: 100%">
            <el-table-column type="index" label="排名" width="60" align="center">
              <template #default="scope">
                <span :class="getRankClass(scope.$index)">{{ scope.$index + 1 }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="challenge_name" label="挑战名称" min-width="120" show-overflow-tooltip />
            
            <el-table-column prop="participant_count" label="热度" align="center" width="80">
               <template #default="scope">
                 <span style="color: #ff9800; font-weight: bold;">{{ scope.row.participant_count }}</span>
               </template>
            </el-table-column>
            
            <el-table-column prop="avg_progress" label="平均进度" width="140">
               <template #default="scope">
                 <el-progress 
                   :percentage="Math.min(Number(scope.row.avg_progress), 100)" 
                   :stroke-width="8" 
                   :color="getProgressColor(scope.row.avg_progress)"
                 />
               </template>
            </el-table-column>
            
            <el-table-column label="操作" align="right" width="90">
               <template #default="scope">
                 <el-button 
                   size="small" 
                   type="primary" 
                   plain 
                   round
                   @click="joinChallenge(scope.row)"
                 >
                   加入
                 </el-button>
               </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :lg="10" :md="24">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">
                <el-icon class="icon-user"><UserFilled /></el-icon> 社区活跃榜
              </span>
            </div>
          </template>

          <el-table :data="activeUsers" stripe style="width: 100%">
            <el-table-column prop="name" label="用户">
              <template #default="scope">
                <div style="display: flex; align-items: center; gap: 10px;">
                  <el-avatar :size="28" style="background: #e0e0e0; color: #666;">
                    {{ scope.row.name ? scope.row.name[0] : 'U' }}
                  </el-avatar>
                  {{ scope.row.name }}
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="health_id" label="Health ID" width="100" />
            <el-table-column prop="activity_score" label="活跃分" align="right" width="90">
              <template #default="scope">
                <el-tag effect="plain" type="success">{{ scope.row.activity_score }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const popularChallenges = ref([])
const activeUsers = ref([]) // 补全了之前代码中缺失的变量定义

// 加入挑战逻辑
const joinChallenge = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要加入"${row.challenge_name}"挑战吗?`, '加入确认', {
      confirmButtonText: '加入',
      cancelButtonText: '再想想',
      type: 'info'
    })
    
    // 接口文档：POST /challenges/{id}/join
    await request.post(`/challenges/${row.challenge_id}/join`, {
      user_id: parseInt(userStore.userId)
    })
    
    ElMessage.success('成功加入挑战！')
    // 简单模拟热度+1，避免刷新页面
    row.participant_count += 1
  } catch(e) {
    if(e !== 'cancel') ElMessage.error('加入失败或您已在挑战中')
  }
}

// 样式辅助函数
const getRankClass = (index) => {
  if (index === 0) return 'rank-1'
  if (index === 1) return 'rank-2'
  if (index === 2) return 'rank-3'
  return 'rank-common'
}

const getProgressColor = (val) => {
  if (val >= 100) return '#67c23a'
  if (val >= 60) return '#409eff'
  return '#e6a23c'
}

onMounted(async () => {
  try {
    // 获取热门挑战
    const res1 = await request.get('/stats/popular_challenges?limit=5')
    popularChallenges.value = res1 || []
    
    // 获取活跃用户
    const res2 = await request.get('/stats/active_users?limit=10')
    activeUsers.value = res2 || []
  } catch(e) {
    // 接口可能未通，保持空数组不报错
  }
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.page-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}
.page-header h2 { margin: 0; font-size: 24px; color: #333; }

.stat-card {
  border: none;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
  font-size: 16px;
  color: #333;
}

.icon-trophy { color: #ff9800; font-size: 20px; }
.icon-user { color: #007bff; font-size: 20px; }

/* 排名样式 */
.rank-1 { color: #f56c6c; font-weight: 900; font-size: 16px; }
.rank-2 { color: #ff9800; font-weight: 800; font-size: 15px; }
.rank-3 { color: #e6a23c; font-weight: 700; font-size: 14px; }
.rank-common { color: #909399; font-weight: 500; }
</style>