<template>
  <div>
    <el-page-header @back="$router.go(-1)" content="搜索记录" class="mb-20" />
    
    <el-card class="mb-20">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item label="类型">
          <el-select v-model="filters.type" placeholder="记录类型" style="width: 150px;">
            <el-option label="全部" value="all" />
            <el-option label="预约记录" value="apt" />
            <el-option label="健康指标" value="metric" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>

        <el-form-item label="关键字">
          <el-input v-model="filters.keyword" placeholder="健康ID / 医生姓名" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleSearch">搜索</el-button>
          <el-button icon="Download">导出 CSV</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="results" border v-loading="loading">
      <el-table-column prop="id" label="ID" width="100" />
      <el-table-column prop="date" label="日期" width="180" />
      <el-table-column prop="type" label="类型" width="150">
          <template #default="scope">
              <el-tag>{{ scope.row.type }}</el-tag>
          </template>
      </el-table-column>
      <el-table-column prop="detail" label="详细信息" />
    </el-table>
  </div>
</template>

<script setup>
import { ref } from 'vue';
// 修正：确保这里引用的是 api/stats 而不是 mockService，解决截图中的报错
import { searchRecords } from '@/api/stats'; 

const filters = ref({ 
  type: 'all', 
  keyword: '', 
  dateRange: [] // 新增日期范围字段
});
const results = ref([]);
const loading = ref(false);

const handleSearch = async () => {
  loading.value = true;
  try {
    // 将 filters 传给后端接口
    results.value = await searchRecords(filters.value);
  } finally {
    loading.value = false;
  }
};
</script>
<style scoped> .mb-20 { margin-bottom: 20px; } </style>