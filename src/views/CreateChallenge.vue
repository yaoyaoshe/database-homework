<template>
  <div>
    <el-page-header @back="$router.go(-1)" content="创建健康挑战" class="mb-20" />
    
    <el-card>
      <el-form label-position="top" size="large">
        <el-form-item label="挑战目标">
          <el-input v-model="form.goal" placeholder="例如：一个月内步行 100 公里" />
        </el-form-item>

        <el-form-item label="持续时间">
          <el-date-picker v-model="form.dates" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 100%" />
        </el-form-item>

        <el-form-item label="邀请参与者">
          <div style="display: flex; gap: 10px; margin-bottom: 15px;">
            <el-input v-model="inviteInput" placeholder="输入邮箱或电话号码" style="max-width: 400px;" @keyup.enter="addInvitee" />
            <el-button @click="addInvitee">添加</el-button>
          </div>
          <div class="tags-area">
             <el-tag v-for="(user, index) in form.invitees" :key="index" closable @close="form.invitees.splice(index, 1)" size="large" style="margin-right: 10px;">
              {{ user }}
            </el-tag>
          </div>
        </el-form-item>

        <el-button type="primary" @click="create" size="large" style="margin-top: 20px;">创建挑战</el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { ElMessage } from 'element-plus';
import { createChallenge } from '@/api/challenge'; // 使用新接口

const form = ref({ goal: '', dates: [], invitees: [] });
const inviteInput = ref('');

const addInvitee = () => {
  if (inviteInput.value) {
    form.value.invitees.push(inviteInput.value);
    inviteInput.value = '';
  }
};

const create = async () => {
  if(!form.value.goal) return ElMessage.warning("请输入挑战目标");
  await createChallenge(form.value);
  ElMessage.success('挑战创建成功，已发送邀请！');
};
</script>

<style scoped> .mb-20 { margin-bottom: 20px; } </style>