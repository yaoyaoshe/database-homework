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
          <el-table-column label="进度" width="180">
             <template #default="scope">
               <el-progress :percentage="Math.min(Number(scope.row.current_progress||0), 100)" />
               <div class="small-text">目标: {{ scope.row.target_value }} {{ scope.row.target_unit }}</div>
             </template>
          </el-table-column>
          <el-table-column label="操作" width="120" align="right">
            <template #default="scope">
              <el-button type="primary" size="small" plain @click="openInviteDialog(scope.row)">邀请好友</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="创建新挑战" name="create">
        <el-form :model="form" label-width="100px" style="max-width: 600px; margin: 20px auto;">
          <el-form-item label="挑战名称"><el-input v-model="form.challenge_name" /></el-form-item>
          <el-form-item label="类型">
             <el-select v-model="form.challenge_type" style="width:100%">
               <el-option label="减重" value="减重" /><el-option label="运动" value="运动" />
             </el-select>
          </el-form-item>
          <el-form-item label="目标数值"><el-input-number v-model="form.target_value" style="width:100%" /></el-form-item>
          <el-form-item label="单位"><el-input v-model="form.target_unit" placeholder="例如: kg, 步" /></el-form-item>
          <el-form-item>
             <el-button type="primary" @click="createChallenge">创建</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="inviteVisible" title="邀请好友加入挑战" width="400px">
      <el-form :model="inviteForm">
        <el-form-item label="邀请方式">
          <el-radio-group v-model="inviteForm.recipient_type">
            <el-radio label="Email">邮箱</el-radio>
            <el-radio label="Phone">电话</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="联系方式">
          <el-input v-model="inviteForm.recipient_value" placeholder="输入邮箱或电话" />
        </el-form-item>
        <el-form-item label="留言">
          <el-input v-model="inviteForm.message" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="inviteVisible = false">取消</el-button>
        <el-button type="primary" @click="sendInvite">发送邀请</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const activeTab = ref('my_challenges')
const myChallenges = ref([])
const loading = ref(false)
const inviteVisible = ref(false)
const currentChallengeId = ref(null)

const form = reactive({ challenge_name: '', challenge_type: '运动', target_value: 0, target_unit: '步' })
const inviteForm = reactive({ recipient_type: 'Email', recipient_value: '', message: '来和我一起挑战吧！' })

const fetchMyChallenges = async () => {
  loading.value = true
  try {
    const res = await request.get(`/users/${userStore.userId}/challenges`)
    myChallenges.value = res || []
  } catch(e){} finally { loading.value = false }
}

const createChallenge = async () => {
  try {
    await request.post('/challenges', { creator_id: parseInt(userStore.userId), ...form })
    ElMessage.success('创建成功'); activeTab.value = 'my_challenges'; fetchMyChallenges()
  } catch(e){}
}

const openInviteDialog = (row) => {
  currentChallengeId.value = row.challenge_id
  inviteVisible.value = true
}

const sendInvite = async () => {
  try {
    await request.post(`/challenges/${currentChallengeId.value}/invite`, {
      sender_id: parseInt(userStore.userId),
      ...inviteForm
    })
    ElMessage.success('邀请已发送'); inviteVisible.value = false
  } catch(e){}
}

onMounted(fetchMyChallenges)
</script>
<style scoped>
.page-container { padding: 20px; }
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.small-text { font-size: 12px; color: #999; margin-top: 4px; }
</style>