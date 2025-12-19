<template>
  <div>
    <el-page-header @back="$router.go(-1)" content="账户信息" class="mb-20" />
    
    <div v-loading="loading">
      <el-card class="mb-20">
        <template #header><span class="card-title">个人详细信息</span></template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="姓名">{{ profile.name }}</el-descriptions-item>
          <el-descriptions-item label="健康ID">{{ profile.healthId }}</el-descriptions-item>
          <el-descriptions-item label="出生日期">{{ profile.dob }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="mb-20" style="height: 100%;">
            <template #header><span class="card-title">联系方式</span></template>
            <h4>电子邮箱</h4>
            <div v-for="email in profile.emails" :key="email.id" class="item-row">
              <el-tag :type="email.verified ? 'success' : 'warning'">
                {{ email.address }} {{ email.verified ? '(已验证)' : '(未验证)' }}
              </el-tag>
            </div>
            <h4>电话号码</h4>
            <div v-for="phone in profile.phones" :key="phone.id" class="item-row">
              <el-tag type="success">{{ phone.number }}</el-tag>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="12">
          <el-card class="mb-20" style="height: 100%;">
            <template #header><span class="card-title">关联医疗提供方</span></template>
            <el-table :data="profile.providers" stripe style="width: 100%">
              <el-table-column prop="name" label="姓名" />
              <el-table-column prop="type" label="专科类型" />
              <el-table-column label="状态" width="100">
                <template #default="scope">
                  <el-tag v-if="scope.row.isPrimary" effect="dark">主管医生</el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUserProfile } from '@/api/user'

const loading = ref(true)
const profile = ref({
  emails: [],
  phones: [],
  providers: []
})

onMounted(async () => {
  try {
    // 从本地存储获取用户ID
    const userId = localStorage.getItem('userId')
    profile.value = await getUserProfile(userId)
  } finally {
    loading.value = false
  }
})
</script>


<!-- <script setup>
import { ref, onMounted } from 'vue';
import { getUserProfile } from '@/api/user'; // 使用新接口

const loading = ref(true);
const profile = ref({ emails: [], phones: [], providers: [] });

onMounted(async () => {
  try {
    profile.value = await getUserProfile();
  } finally {
    loading.value = false;
  }
});
</script> -->

<style scoped>
.mb-20 { margin-bottom: 20px; }
.card-title { font-weight: bold; font-size: 16px; }
.item-row { margin-bottom: 8px; }
</style>