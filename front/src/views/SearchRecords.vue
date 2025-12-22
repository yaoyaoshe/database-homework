<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>搜索记录</h2>
    </div>

    <el-card>
      <el-form :inline="true" class="search-form">
        <el-form-item label="搜索类型">
          <el-select v-model="searchType" style="width: 140px">
            <el-option label="预约记录" value="appointment" />
            <el-option label="健康数据" value="health" />
          </el-select>
        </el-form-item>

        <template v-if="searchType === 'appointment'">
          <el-form-item label="关键词">
            <el-input v-model="keyword" placeholder="搜医生/原因..." />
          </el-form-item>
        </template>

        <template v-else>
          <el-form-item label="指标类型">
            <el-select v-model="healthType" style="width: 120px">
              <el-option label="体重" value="Weight" />
              <el-option label="步数" value="Steps" />
            </el-select>
          </el-form-item>
          <el-form-item label="月份">
            <el-date-picker v-model="searchMonth" type="month" value-format="YYYY-MM" placeholder="不限" />
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div style="margin-top: 20px;">
      <div v-if="searchType === 'health' && healthResults">
        <el-alert type="success" :closable="false" style="margin-bottom: 15px;">
          <h3>总计数值: {{ healthResults.total_value }}</h3>
        </el-alert>
        <el-table :data="healthResults.records" border stripe>
          <el-table-column prop="recorded_at" label="记录时间" />
          <el-table-column prop="data_value" label="数值" />
          <el-table-column prop="unit" label="单位" />
        </el-table>
      </div>

      <el-table v-else-if="searchType === 'appointment'" :data="filteredAppts" border stripe>
        <el-table-column prop="appointment_date" label="时间" />
        <el-table-column prop="consultation_type" label="类型" />
        <el-table-column prop="status" label="状态" />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const searchType = ref('appointment')
// 预约搜索 state
const keyword = ref('')
const apptList = ref([])
// 健康搜索 state
const healthType = ref('Steps')
const searchMonth = ref('')
const healthResults = ref(null)

const handleSearch = async () => {
  if (searchType.value === 'appointment') {
    const res = await request.get(`/users/${userStore.userId}/appointments`)
    apptList.value = res || []
  } else {
    // 调用新的高级搜索接口
    const res = await request.get(`/search/health`, {
      params: {
        user_id: userStore.userId,
        type: healthType.value,
        month: searchMonth.value
      }
    })
    healthResults.value = res
  }
}

const filteredAppts = computed(() => {
  if (!keyword.value) return apptList.value
  return apptList.value.filter(i => JSON.stringify(i).toLowerCase().includes(keyword.value.toLowerCase()))
})
</script>

<style scoped>
.page-container { padding: 20px; }
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
</style>