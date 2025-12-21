<template>
  <div class="page-container">
    <div class="page-header">
      <el-button icon="ArrowLeft" circle @click="$router.push('/')" style="margin-right: 15px" />
      <h2>搜索中心</h2>
    </div>

    <el-card shadow="hover" style="margin-bottom: 20px;">
      <el-form :inline="true" class="search-form">
        <el-form-item label="数据类型">
          <el-select v-model="searchType" style="width: 140px">
            <el-option label="预约记录" value="appointment" />
            </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input 
            v-model="keyword" 
            placeholder="搜医生ID/类型/说明..." 
            prefix-icon="Search" 
            clearable
            @input="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态筛选">
           <el-select v-model="filterStatus" placeholder="全部状态" clearable @change="handleSearch" style="width: 120px">
            <el-option label="已预约" value="已预约" />
            <el-option label="已完成" value="已完成" />
            <el-option label="已取消" value="已取消" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-loading="loading">
      <template #header>
        <div style="display: flex; justify-content: space-between;">
           <span>搜索结果 ({{ filteredList.length }})</span>
           <el-button v-if="filteredList.length > 0" size="small" type="success" @click="exportData">导出Excel</el-button>
        </div>
      </template>

      <el-table :data="filteredList" stripe style="width: 100%" height="400">
        <el-table-column prop="appointment_id" label="ID" width="80" />
        <el-table-column label="时间" width="180">
          <template #default="scope">{{ formatTime(scope.row.appointment_date) }}</template>
        </el-table-column>
        <el-table-column prop="provider_id" label="医生ID" width="100" />
        <el-table-column prop="consultation_type" label="类型" />
        <el-table-column prop="reason" label="原因/症状" show-overflow-tooltip />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
             <el-tag :type="scope.row.status === '已预约' ? 'primary' : 'info'">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <el-empty v-if="!loading && filteredList.length === 0" description="未找到匹配记录" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import * as XLSX from 'xlsx' // 建议安装: npm install xlsx

const userStore = useUserStore()
const searchType = ref('appointment')
const keyword = ref('')
const filterStatus = ref('')
const rawData = ref([])
const loading = ref(false)

// 获取原始数据
const fetchData = async () => {
  loading.value = true
  try {
    // 接口文档：GET /users/{id}/appointments
    const res = await request.get(`/users/${userStore.userId}/appointments`)
    rawData.value = res || []
  } catch(e) {
    rawData.value = []
  } finally {
    loading.value = false
  }
}

// 前端过滤逻辑
const filteredList = computed(() => {
  return rawData.value.filter(item => {
    // 1. 状态筛选
    if (filterStatus.value && item.status !== filterStatus.value) return false
    
    // 2. 关键词模糊匹配 (匹配 医生ID, 类型, 原因)
    if (!keyword.value) return true
    const searchStr = keyword.value.toLowerCase()
    
    return (
      String(item.provider_id).includes(searchStr) ||
      (item.consultation_type && item.consultation_type.includes(searchStr)) ||
      (item.reason && item.reason.includes(searchStr))
    )
  })
})

const handleSearch = () => {
  // 触发 computed 重新计算，实际是自动的，这里主要用于防抖扩展
}

const formatTime = (t) => new Date(t).toLocaleString()

// 纯前端导出功能
const exportData = () => {
  const ws = XLSX.utils.json_to_sheet(filteredList.value)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, "SearchResults")
  XLSX.writeFile(wb, "search_export.xlsx")
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 24px; color: #333; }
.search-form { display: flex; flex-wrap: wrap; }
</style>