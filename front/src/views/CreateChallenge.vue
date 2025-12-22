<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>健康挑战中心</h2>
    </div>

    <el-row :gutter="24">
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
              <el-input v-model="form.description" type="textarea" :rows="2" placeholder="简单描述挑战的目标..." />
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
                    <el-option label="综合" value="综合" />
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
                  <el-input v-model="form.target_unit" placeholder="如: kg" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="挑战周期" prop="dateRange">
              <el-date-picker
                v-model="form.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
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
            <el-table-column prop="challenge_name" label="名称" min-width="140" />
            <el-table-column label="类型" width="100">
              <template #default="{ row }">
                <el-tag size="small" :type="getTagType(row.challenge_type)">{{ row.challenge_type }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="进度" min-width="180">
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
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <span :class="getStatusClass(row.status)">{{ row.status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="right">
              <template #default="{ row }">
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
            <el-radio label="邮箱">邮箱</el-radio>
            <el-radio label="手机号">手机号</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="inviteForm.recipient_type === '邮箱' ? '邮箱地址' : '手机号码'">
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
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, watch } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const formRef = ref(null)

// 状态控制
const loading = ref(false)
const creating = ref(false)
const inviting = ref(false)
const inviteVisible = ref(false)
const myChallenges = ref([])
const currentChallengeId = ref(null)

// 表单数据初始化
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

const inviteForm = reactive({
  recipient_type: '邮箱',
  recipient_value: '',
  message: '一起来参加这个健康挑战吧！'
})

// 根据指标自动关联单位
watch(() => form.target_metric, (val) => {
  const map = {
    '体重': 'kg',
    '步数': '步',
    '睡眠时长': '小时',
    '运动时长': '分钟',
    '卡路里消耗': '千卡'
  }
  if (map[val]) {
    form.target_unit = map[val]
  }
})

// 校验规则
const rules = {
  challenge_name: [{ required: true, message: '请输入挑战名称', trigger: 'blur' }],
  challenge_type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  target_metric: [{ required: true, message: '请选择指标', trigger: 'change' }],
  target_value: [{ required: true, message: '请输入目标值', trigger: 'blur' }],
  target_unit: [{ required: true, message: '请输入单位', trigger: 'blur' }],
  dateRange: [{ required: true, message: '请选择起止时间', trigger: 'change' }]
}

// 禁止选择过去的日期
const disabledDate = (time) => {
  return time.getTime() < Date.now() - 8.64e7
}

const fetchMyChallenges = async () => {
  if (!userStore.userId) return
  loading.value = true
  try {
    const res = await request.get(`/users/${userStore.userId}/challenges`)
    myChallenges.value = res || []
  } catch (error) {
    console.error("获取挑战失败", error)
    ElMessage.error("获取列表失败")
  } finally {
    loading.value = false
  }
}

const handleCreateChallenge = async () => {
  if (!formRef.value) return
  if (!userStore.userId) {
    ElMessage.error("用户信息丢失，请重新登录")
    return
  }
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      creating.value = true
      try {
        // 关键修复：手动构建 Go 语言 time.Time 能接受的 RFC3339 格式
        // 避免使用 new Date().toISOString() 导致的时区偏移问题
        const startDateStr = `${form.dateRange[0]}T00:00:00Z`
        const endDateStr = `${form.dateRange[1]}T23:59:59Z`

        const payload = {
          creator_id: parseInt(userStore.userId), // 确保是整数
          challenge_name: form.challenge_name,
          description: form.description || form.challenge_name,
          challenge_type: form.challenge_type,
          target_metric: form.target_metric,
          target_value: parseFloat(form.target_value), // 确保是浮点数
          target_unit: form.target_unit,
          start_date: startDateStr, 
          end_date: endDateStr,
          max_participants: 50, // 数据库默认值逻辑在后端未体现，前端传参更安全
          is_public: form.is_public
        }

        console.log("发送创建请求:", payload) // 用于调试
        const res = await request.post('/challenges', payload)
        
        // 检查后端返回的 ID，防止因后端静默失败导致的无效 ID
        if (res && res.challenge_id && res.challenge_id > 0) {
          ElMessage.success('挑战创建成功，正在加入...')
          
          // 自动加入挑战
          await joinChallenge(res.challenge_id)
          
          resetForm()
          fetchMyChallenges()
        } else {
          throw new Error("创建返回了无效的 ID，请检查输入格式")
        }
      } catch (error) {
        console.error(error)
        ElMessage.error(error.message || '创建失败，请检查网络或联系管理员')
      } finally {
        creating.value = false
      }
    }
  })
}

// 抽离加入逻辑
const joinChallenge = async (challengeId) => {
  try {
    await request.post(`/challenges/${challengeId}/join`, {
      user_id: parseInt(userStore.userId)
    })
    ElMessage.success('已成功加入挑战！')
  } catch (error) {
    console.error("自动加入失败", error)
    ElMessage.warning('挑战创建成功，但自动加入失败，请手动加入')
  }
}

const resetForm = () => {
  formRef.value.resetFields()
  form.description = ''
  form.challenge_type = '运动'
  form.target_metric = '步数'
  form.target_value = 10000
}

const openInviteDialog = (row) => {
  currentChallengeId.value = row.challenge_id
  inviteForm.recipient_value = ''
  inviteVisible.value = true
}

const sendInvite = async () => {
  if (!inviteForm.recipient_value) {
    ElMessage.warning('请输入联系方式')
    return
  }
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
    ElMessage.error('发送邀请失败')
  } finally {
    inviting.value = false
  }
}

// 工具函数
const calculatePercentage = (row) => {
  if (!row.target_value || row.target_value === 0) return 0
  const p = (row.current_progress / row.target_value) * 100
  return Math.min(parseFloat(p.toFixed(1)), 100)
}

const getTagType = (type) => {
  const map = { '减重': 'danger', '运动': 'primary', '健康习惯': 'success', '饮食': 'warning', '睡眠': 'info' }
  return map[type] || 'info'
}

const getStatusClass = (status) => {
  if (status === '已完成') return 'text-success'
  if (status === '进行中' || status === '参与中') return 'text-primary'
  return 'text-gray'
}

onMounted(() => {
  fetchMyChallenges()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}
.page-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
}
.create-card {
  border-top: 4px solid var(--el-color-primary);
  margin-bottom: 20px;
}
.list-card {
  min-height: 500px;
}
.submit-btn {
  width: 100%;
  margin-top: 10px;
  font-weight: bold;
  height: 40px;
}
.progress-wrapper {
  padding-right: 15px;
}
.progress-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
.text-success { color: #67c23a; font-weight: bold; }
.text-primary { color: #409eff; font-weight: bold; }
.text-gray { color: #909399; }
</style>