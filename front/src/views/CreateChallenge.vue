<template>
  <div class="page-container">
    <h2>创建健康挑战</h2>
    <el-card>
      <el-form :model="form" label-width="100px">
        <el-form-item label="挑战名称">
          <el-input v-model="form.challenge_name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.challenge_type">
            <el-option label="减重" value="减重" />
            <el-option label="运动" value="运动" />
            <el-option label="睡眠" value="睡眠" />
          </el-select>
        </el-form-item>
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
        <el-form-item label="目标指标">
          <el-input v-model="form.target_metric" placeholder="如：体重" />
        </el-form-item>
        <el-form-item label="目标值">
          <el-input-number v-model="form.target_value" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="form.target_unit" placeholder="如：kg" />
        </el-form-item>
        
        <el-button type="primary" @click="createChallenge">创建挑战</el-button>
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
.page-container { padding: 20px; max-width: 800px; margin: 0 auto; }
</style>