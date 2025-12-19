<template>
  <div>
    <el-page-header @back="$router.go(-1)" content="月度健康摘要" class="mb-20" />

    <el-card class="mb-20">
      <span style="margin-right: 15px; font-weight: bold;">选择月份:</span>
      <el-date-picker v-model="selectedMonth" type="month" placeholder="选择月份" @change="fetchData" :clearable="false" />
    </el-card>

    <div v-loading="loading">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card header="步数历史 (本月)" class="mb-20">
            <div id="stepsChart" style="height: 300px;"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card header="体征数据 (体重 & 血压)" class="mb-20">
             <div id="vitalsChart" style="height: 300px;"></div>
          </el-card>
        </el-col>
      </el-row>

      <el-card header="本月预约记录" class="mb-20">
        <el-table :data="appointments" stripe>
          <el-table-column prop="date" label="日期" />
          <el-table-column prop="provider" label="医生" />
          <el-table-column prop="type" label="类型" />
          <el-table-column prop="status" label="状态">
            <template #default="scope">
              <el-tag :type="scope.row.status === '已完成' ? 'success' : 'primary'">{{ scope.row.status }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
      
      <el-alert :title="summaryText" type="info" show-icon :closable="false" class="mb-20" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import * as echarts from 'echarts';
import { getMonthlyData } from '@/api/stats'; // 使用新接口

const selectedMonth = ref(new Date());
const loading = ref(false);
const appointments = ref([]);
const summaryText = ref('');

const fetchData = async () => {
  loading.value = true;
  try {
    const data = await getMonthlyData(selectedMonth.value.getFullYear(), selectedMonth.value.getMonth() + 1);
    appointments.value = data.appointments;
    summaryText.value = data.summaryText;
    
    await nextTick();
    initCharts(data.metrics);
  } finally {
    loading.value = false;
  }
};

const initCharts = (metrics) => {
  // 销毁旧实例防止内存泄漏（简单处理）
  echarts.getInstanceByDom(document.getElementById('stepsChart'))?.dispose();
  echarts.getInstanceByDom(document.getElementById('vitalsChart'))?.dispose();

  const stepsChart = echarts.init(document.getElementById('stepsChart'));
  stepsChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['W1', 'W2', 'W3', 'W4'] },
    yAxis: { type: 'value' },
    series: [{ data: metrics.steps.slice(0, 4), type: 'line', smooth: true, itemStyle: { color: '#0072ff' } }]
  });

  const vitalsChart = echarts.init(document.getElementById('vitalsChart'));
  vitalsChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['体重', '收缩压'] },
    xAxis: { type: 'category', data: ['W1', 'W2', 'W3', 'W4'] },
    yAxis: { type: 'value' },
    series: [
      { name: '体重', type: 'bar', data: metrics.weight, itemStyle: { color: '#00c6ff' } },
      { name: '收缩压', type: 'line', data: metrics.bp_sys, itemStyle: { color: '#ff4d4f' } }
    ]
  });
  
  // 窗口大小改变时重绘
  window.addEventListener('resize', () => {
    stepsChart.resize();
    vitalsChart.resize();
  });
};

onMounted(fetchData);
</script>

<style scoped> .mb-20 { margin-bottom: 20px; } </style>