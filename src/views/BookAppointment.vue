<template>
  <div>
    <el-page-header @back="$router.go(-1)" content="预约服务" class="mb-20" />
    
    <el-card style="max-width: 800px; margin: 0 auto; padding: 20px;">
      <el-form :model="form" label-width="120px" size="large">
        <el-form-item label="医生/提供方" required>
          <el-select v-model="form.provider" placeholder="请选择医疗提供方" style="width: 100%">
            <el-option v-for="p in providers" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="日期与时间" required>
          <el-date-picker v-model="form.date" type="datetime" placeholder="选择日期和时间" style="width: 100%" />
        </el-form-item>

        <el-form-item label="就诊类型">
          <el-radio-group v-model="form.type">
            <el-radio label="线下门诊" border />
            <el-radio label="线上问诊" border />
            <el-radio label="电话咨询" border />
          </el-radio-group>
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="form.memo" type="textarea" rows="4" placeholder="请输入就诊原因或症状..." />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="onSubmit" :loading="submitting" style="width: 200px;">确认预约</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getProviders, submitAppointment } from '@/api/appointment'; // 使用新接口

const providers = ref([]);
const submitting = ref(false);
const form = ref({ provider: '', date: '', type: '线下门诊', memo: '' });

onMounted(async () => {
  providers.value = await getProviders();
});

const onSubmit = async () => {
  if(!form.value.provider || !form.value.date) return ElMessage.warning('请填写必填项');
  
  submitting.value = true;
  try {
    await submitAppointment(form.value);
    ElMessage.success('预约成功！');
    form.value.memo = ''; // 重置备注
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped> .mb-20 { margin-bottom: 20px; } </style>