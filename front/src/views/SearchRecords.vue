<template>
  <div class="page-container">
    <h2>搜索记录</h2>
    <el-card>
      <el-form :inline="true">
        <el-form-item label="搜索类型">
          <el-select v-model="searchType">
            <el-option label="预约记录" value="appointment" />
            <el-option label="健康数据" value="health" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="keyword" placeholder="输入关键词" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <div style="margin-top:20px; text-align:center; color:#999">
      暂无搜索结果 (后端暂未开放高级搜索接口)
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const searchType = ref('appointment')
const keyword = ref('')

const handleSearch = () => {
  // 模拟搜索
}
</script>

<style scoped>
.page-container { padding: 20px; }
</style>