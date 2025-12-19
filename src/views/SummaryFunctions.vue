<template>
  <div>
    <el-page-header @back="$router.go(-1)" content="系统统计报表" class="mb-20" />

    <el-card class="mb-20">
      <span class="mr-10" style="margin-right: 10px; font-weight: bold;">统计范围:</span>
      <el-date-picker v-model="dateRange" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" @change="loadStats" />
      <el-button type="primary" style="margin-left: 20px;" @click="loadStats">生成报告</el-button>
      <el-button type="warning">导出 PDF</el-button>
    </el-card>

    <div v-if="stats" v-loading="loading">
      <el-row :gutter="20" class="mb-20">
        <el-col :span="8"><el-card shadow="hover"><el-statistic title="总预约数" :value="stats.totalAppts" /></el-card></el-col>
        <el-col :span="8"><el-card shadow="hover"><el-statistic title="平均体重 (kg)" :value="stats.metrics.weight.avg" precision="1" /></el-card></el-col>
        <el-col :span="8"><el-card shadow="hover"><el-statistic title="最大体重 (kg)" :value="stats.metrics.weight.max" precision="1" /></el-card></el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card header="热门健康挑战">
            <el-table :data="stats.topChallenges" stripe>
              <el-table-column type="index" label="排名" width="60" />
              <el-table-column prop="name" label="挑战名称" />
              <el-table-column prop="count" label="参与人数" align="right" />
            </el-table>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card header="最活跃用户">
            <el-table :data="stats.activeUsers" stripe>
              <el-table-column type="index" label="排名" width="60" />
              <el-table-column prop="name" label="用户" />
              <el-table-column prop="score" label="活跃积分" align="right" />
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { getGlobalStats } from '@/api/stats'; // 使用新接口

const dateRange = ref([]);
const stats = ref(null);
const loading = ref(false);

const loadStats = async () => {
  if(!dateRange.value) return;
  loading.value = true;
  try {
    stats.value = await getGlobalStats(dateRange.value);
  } finally {
    loading.value = false;
  }
};
</script>
<style scoped> .mb-20 { margin-bottom: 20px; } </style>