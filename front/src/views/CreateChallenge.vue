<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>健康挑战中心</h2>
    </div>

    <el-alert 
      v-if="pendingInvites.length > 0" 
      type="info" 
      show-icon 
      :closable="false"
      class="invite-alert"
    >
      <template #title>
        <span style="font-weight: bold; font-size: 15px;">您有 {{ pendingInvites.length }} 个待处理的挑战邀请</span>
      </template>
      <div class="invite-list">
        <div v-for="inv in pendingInvites" :key="inv.invitation_id" class="invite-item">
          <span>挑战ID: <strong>{{ inv.challenge_id }}</strong> (来自用户 {{ inv.sender_id }})</span>
          <div class="invite-msg" v-if="inv.message">留言: {{ inv.message }}</div>
          <el-button type="primary" size="small" @click="handleAcceptInvite(inv.invitation_id)">接受邀请</el-button>
        </div>
      </div>
    </el-alert>

    <el-row :gutter="24" style="margin-top: 20px;">
       <el-col :xs="24" :lg="8">
        <el-card class="box-card create-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span><el-icon><Plus /></el-icon> 发起新挑战</span>
            </div>
          </template>
          
          <el-form :model="form" ref="formRef" :rules="rules" label-position="top" size="large">
            <el-form-item label="挑战名称" prop="challenge_name">
              <el-input v-model="form.challenge_name" placeholder="例如：30天减脂计划" />
            </el-form-item>

            <el-form-item label="挑战描述" prop="description">
              <el-input v-model="form.description" type="textarea" :rows="2" />
            </el-form-item>

            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="类型" prop="challenge_type">
                  <el-select v-model="form.challenge_type" placeholder="请选择" style="width:100%">
                    <el-option label="减重" value="减重" />
                    <el-option label="运动" value="运动" />
                    <el-option label="饮食" value="饮食" />
                    <el-option label="睡眠" value="睡眠" />
                    <el-option label="健康习惯" value="健康习惯" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="目标指标" prop="target_metric">
                  <el-select v-model="form.target_metric" placeholder="请选择" style="width:100%">
                    <el-option label="体重" value="体重" />
                    <el-option label="步数" value="步数" />
                    <el-option label="睡眠时长" value="睡眠时长" />
                    <el-option label="运动时长" value="运动时长" />
                    <el-option label="卡路里消耗" value="卡路里消耗" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="12">
              <el-col :span="14">
                <el-form-item label="目标数值" prop="target_value">
                  <el-input-number v-model="form.target_value" :min="0.1" :precision="1" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="10">
                <el-form-item label="单位" prop="target_unit">
                  <el-input v-model="form.target_unit" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="挑战周期" prop="dateRange">
              <el-date-picker
                v-model="form.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始"
                end-placeholder="结束"
                value-format="YYYY-MM-DD"
                :disabled-date="disabledDate"
                style="width: 100%"
              />
            </el-form-item>
            
            <el-form-item label="公开挑战">
              <el-switch v-model="form.is_public" active-text="公开" inactive-text="私密" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" class="submit-btn" :loading="creating" @click="handleCreateChallenge">
                立即创建并加入
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
       
       <el-col :xs="24" :lg="16">
        <el-card class="box-card list-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span><el-icon><Trophy /></el-icon> 我参加的挑战</span>
              <el-button link type="primary" icon="Refresh" @click="fetchMyChallenges">刷新</el-button>
            </div>
          </template>

          <el-table :data="myChallenges" style="width: 100%" v-loading="loading" empty-text="暂无参与的挑战">
            <el-table-column prop="challenge_name" label="名称" min-width="120" show-overflow-tooltip />
            <el-table-column label="进度" min-width="150">
              <template #default="{ row }">
                <div class="progress-wrapper">
                  <el-progress 
                    :percentage="calculatePercentage(row)" 
                    :status="row.status === '已完成' ? 'success' : ''"
                  />
                  <div class="progress-info">
                    <span>当前: {{ row.current_progress }}</span>
                    <span>目标: {{ row.target_value }} {{ row.target_unit }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <span :class="getStatusClass(row.status)">{{ row.status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" align="right">
              <template #default="{ row }">
                <el-button type="success" link icon="EditPen" @click="openCheckinDialog(row)" :disabled="row.status==='已完成'">打卡</el-button>
                <el-button type="primary" link icon="Position" @click="openInviteDialog(row)">邀请</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="inviteVisible" title="邀请好友加入" width="400px" destroy-on-close>
      <el-form :model="inviteForm" label-position="top">
        <el-form-item label="发送方式">
          <el-radio-group v-model="inviteForm.recipient_type">
            <el-radio label="Health ID">Health ID</el-radio>
            <el-radio label="邮箱">邮箱</el-radio>
            <el-radio label="手机号">手机号</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item :label="inviteForm.recipient_type === '邮箱' ? '邮箱地址' : (inviteForm.recipient_type === 'Health ID' ? '好友 Health ID' : '手机号码')">
          <el-input v-model="inviteForm.recipient_value" placeholder="请输入..." />
        </el-form-item>
        
        <el-form-item label="邀请留言">
          <el-input v-model="inviteForm.message" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="inviteVisible = false">取消</el-button>
        <el-button type="primary" @click="sendInvite" :loading="inviting">发送</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="checkinVisible" title="每日打卡" width="400px" destroy-on-close>
        <p style="margin-bottom:15px;color:#666">挑战：<strong>{{ currentCheckinChallenge?.challenge_name }}</strong></p>
      <el-form :model="checkinForm" label-position="top">
        <el-form-item label="日期">
          <el-date-picker v-model="checkinForm.date" type="date" value-format="YYYY-MM-DD" style="width: 100%" :disabled-date="d => d > new Date()" />
        </el-form-item>
        <el-form-item label="今日完成量">
          <el-input-number v-model="checkinForm.progress_value" :min="0" style="width: 100%" />
          <span style="margin-left: 10px">{{ currentCheckinChallenge?.target_unit }}</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="checkinForm.notes" type="textarea" />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="checkinForm.is_completed">标记为今日已达标</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="checkinVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCheckin" :loading="checkingIn">提交记录</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, watch } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

// ... (其他逻辑保持不变) ...
const userStore = useUserStore()
const formRef = ref(null)

const loading = ref(false)
const creating = ref(false)
const inviting = ref(false)
const checkingIn = ref(false)

const inviteVisible = ref(false)
const checkinVisible = ref(false)

const myChallenges = ref([])
const pendingInvites = ref([])
const currentChallengeId = ref(null)
const currentCheckinChallenge = ref(null)

const form = reactive({
  challenge_name: '',
  description: '',
  challenge_type: '运动',
  target_metric: '步数',
  target_value: 10000,
  target_unit: '步',
  dateRange: [], 
  is_public: true
})

// 修改：默认选中 Health ID
const inviteForm = reactive({
  recipient_type: 'Health ID',
  recipient_value: '',
  message: '一起来参加这个健康挑战吧！'
})

const checkinForm = reactive({
  date: new Date().toISOString().substring(0, 10),
  progress_value: 0,
  is_completed: false,
  notes: ''
})

watch(() => form.target_metric, (val) => {
  const map = { '体重': 'kg', '步数': '步', '睡眠时长': '小时', '运动时长': '分钟', '卡路里消耗': '千卡' }
  if (map[val]) form.target_unit = map[val]
})

const rules = {
  challenge_name: [{ required: true, message: '请输入挑战名称', trigger: 'blur' }],
  challenge_type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  target_metric: [{ required: true, message: '请选择指标', trigger: 'change' }],
  target_value: [{ required: true, message: '请输入目标值', trigger: 'blur' }],
  dateRange: [{ required: true, message: '请选择起止时间', trigger: 'change' }]
}

const disabledDate = (time) => time.getTime() < Date.now() - 8.64e7

const fetchMyChallenges = async () => {
  if (!userStore.userId) return
  loading.value = true
  try {
    const [challRes, invRes] = await Promise.all([
       request.get(`/users/${userStore.userId}/challenges`),
       request.get(`/users/${userStore.userId}/pending-invites`)
    ])
    myChallenges.value = challRes || []
    pendingInvites.value = invRes || []
  } catch (error) {
    ElMessage.error("获取列表失败")
  } finally {
    loading.value = false
  }
}

const handleAcceptInvite = async (invId) => {
  try {
    await request.post('/invitations/accept', {
      invitation_id: invId,
      user_id: parseInt(userStore.userId)
    })
    ElMessage.success('已接受邀请')
    fetchMyChallenges()
  } catch(e) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  }
}

const handleCreateChallenge = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      creating.value = true
      try {
        const startDateStr = `${form.dateRange[0]}T00:00:00Z`
        const endDateStr = `${form.dateRange[1]}T23:59:59Z`

        const payload = {
          creator_id: parseInt(userStore.userId),
          challenge_name: form.challenge_name,
          description: form.description || form.challenge_name,
          challenge_type: form.challenge_type,
          target_metric: form.target_metric,
          target_value: parseFloat(form.target_value),
          target_unit: form.target_unit,
          start_date: startDateStr, 
          end_date: endDateStr,
          max_participants: 50,
          is_public: form.is_public
        }
        const res = await request.post('/challenges', payload)
        if (res && res.challenge_id) {
          ElMessage.success('挑战创建成功')
          await joinChallenge(res.challenge_id)
          formRef.value.resetFields()
          fetchMyChallenges()
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '创建失败')
      } finally {
        creating.value = false
      }
    }
  })
}

const joinChallenge = async (cid) => {
  try {
    await request.post(`/challenges/${cid}/join`, { user_id: parseInt(userStore.userId) })
  } catch (e) { console.error(e) }
}

const openInviteDialog = (row) => {
  currentChallengeId.value = row.challenge_id
  inviteForm.recipient_value = ''
  inviteVisible.value = true
}

const sendInvite = async () => {
  if (!inviteForm.recipient_value) return ElMessage.warning('请输入联系方式')
  inviting.value = true
  try {
    await request.post(`/challenges/${currentChallengeId.value}/invite`, {
      sender_id: parseInt(userStore.userId),
      recipient_type: inviteForm.recipient_type,
      recipient_value: inviteForm.recipient_value,
      message: inviteForm.message
    })
    ElMessage.success('邀请已发送')
    inviteVisible.value = false
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '发送邀请失败，请检查Health ID')
  } finally {
    inviting.value = false
  }
}

const openCheckinDialog = (row) => {
  currentCheckinChallenge.value = row
  checkinForm.progress_value = 0
  checkinForm.notes = ''
  checkinForm.is_completed = false
  checkinVisible.value = true
}

const submitCheckin = async () => {
  if (!checkinForm.progress_value && checkinForm.progress_value !== 0) return
  checkingIn.value = true
  try {
    await request.post(`/challenges/${currentCheckinChallenge.value.challenge_id}/checkin`, {
      user_id: parseInt(userStore.userId),
      date: checkinForm.date,
      progress_value: parseFloat(checkinForm.progress_value),
      progress_unit: currentCheckinChallenge.value.target_unit,
      is_completed: checkinForm.is_completed,
      notes: checkinForm.notes
    })
    ElMessage.success('打卡成功')
    checkinVisible.value = false
    fetchMyChallenges()
  } catch(e) {
    ElMessage.error(e.response?.data?.error || '打卡失败')
  } finally {
    checkingIn.value = false
  }
}

const calculatePercentage = (row) => {
  if (!row.target_value) return 0
  const p = (row.current_progress / row.target_value) * 100
  return Math.min(parseFloat(p.toFixed(1)), 100)
}
const getStatusClass = (status) => {
  if (status === '已完成') return 'text-success'
  return 'text-primary'
}

onMounted(fetchMyChallenges)
</script>

<style scoped>
/* 样式保持不变 */
.page-container { padding: 24px; max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: 600; }
.create-card { border-top: 4px solid var(--el-color-primary); margin-bottom: 20px; }
.list-card { min-height: 500px; }
.submit-btn { width: 100%; margin-top: 10px; font-weight: bold; }
.progress-wrapper { padding-right: 15px; }
.progress-info { display: flex; justify-content: space-between; font-size: 12px; color: #909399; margin-top: 4px; }
.text-success { color: #67c23a; font-weight: bold; }
.text-primary { color: #409eff; }
.invite-alert { margin-bottom: 20px; }
.invite-list { display: flex; flex-direction: column; gap: 10px; margin-top: 10px; }
.invite-item { display: flex; justify-content: space-between; align-items: center; background: rgba(255,255,255,0.6); padding: 8px 12px; border-radius: 4px; }
.invite-msg { font-size: 12px; color: #666; margin-left: 10px; font-style: italic; }
</style>