<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>创建健康挑战</h2>
    </div>

    <el-card class="form-card">
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
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-divider content-position="left">目标设定</el-divider>
        
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="目标指标">
              <el-input v-model="form.target_metric" placeholder="如：体重" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
             <el-form-item label="目标值">
              <el-input-number v-model="form.target_value" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="单位">
              <el-input v-model="form.target_unit" placeholder="如：kg" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <div style="margin-top: 20px; text-align: right;">
           <el-button @click="$router.push('/')">取消</el-button>
           <el-button type="primary" @click="createChallenge">立即创建挑战</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
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

const createChallenge = async () => {
  if (!dateRange.value || dateRange.value.length < 2) return ElMessage.warning('请选择日期')
  
  try {
    await request.post('/challenges', {
      creator_id: parseInt(userStore.userId),
      ...form,
      start_date: dateRange.value[0],
      end_date: dateRange.value[1]
    })
    ElMessage.success('挑战创建成功')
  } catch(e) {}
}
</script>

<style scoped>
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 24px; color: #333; }
.form-card { padding: 20px; }
</style>