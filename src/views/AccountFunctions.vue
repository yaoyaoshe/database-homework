<template>
  <div>
    <el-page-header @back="$router.go(-1)" content="账户功能管理" class="mb-20" />

    <el-tabs type="border-card" class="function-tabs">
      <el-tab-pane label="修改个人信息">
        <el-form label-width="100px" style="max-width: 500px; padding-top: 20px;">
          <el-form-item label="姓名"><el-input v-model="form.name" /></el-form-item>
          <el-form-item label="健康ID"><el-input v-model="form.id" disabled /></el-form-item>
          <el-button type="primary" style="margin-left: 100px;">保存更改</el-button>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="管理邮箱">
        <div class="add-box">
          <el-input v-model="newEmail" placeholder="输入新邮箱地址" style="width: 300px; margin-right: 10px;" />
          <el-button type="success" @click="mockAdd('邮箱')">添加</el-button>
        </div>
        <el-table :data="emails" border>
          <el-table-column prop="address" label="邮箱地址" />
          <el-table-column label="操作" width="120">
             <template #default> <el-button type="danger" size="small" plain>删除</el-button> </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="管理电话">
         <div class="add-box">
          <el-input v-model="newPhone" placeholder="输入新电话号码" style="width: 300px; margin-right: 10px;" />
          <el-button type="success" @click="mockAdd('电话')">添加</el-button>
        </div>
         <el-table :data="phones" border>
          <el-table-column prop="number" label="电话号码" />
           <el-table-column label="操作" width="120">
             <template #default> <el-button type="danger" size="small" plain>删除</el-button> </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="管理提供方">
        <div class="add-box">
          <el-input v-model="newProvider" placeholder="输入提供方执照号" style="width: 300px; margin-right: 10px;" />
          <el-button type="success" @click="mockAdd('提供方')">关联提供方</el-button>
        </div>
        <el-table :data="providers" border>
          <el-table-column prop="name" label="姓名" />
          <el-table-column prop="license" label="执照号" />
           <el-table-column label="操作" width="120">
             <template #default> <el-button type="danger" size="small" plain>取消关联</el-button> </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getUserProfile } from '@/api/user';
import { ElMessage } from 'element-plus';

const form = ref({ name: '', id: '' });
const emails = ref([]);
const phones = ref([]);
const providers = ref([]);

const newEmail = ref('');
const newPhone = ref('');
const newProvider = ref('');

onMounted(async () => {
  const data = await getUserProfile();
  form.value = { name: data.name, id: data.healthId };
  emails.value = data.emails;
  phones.value = data.phones;
  providers.value = data.providers;
});

// 优化后的添加逻辑：不只弹窗，还更新列表（前端模拟）
const mockAdd = (type) => {
  if (type === '邮箱' && newEmail.value) {
    emails.value.push({ address: newEmail.value, verified: false });
    newEmail.value = '';
    ElMessage.success('已添加新邮箱');
  } else if (type === '电话' && newPhone.value) {
    phones.value.push({ number: newPhone.value, verified: false });
    newPhone.value = '';
    ElMessage.success('已添加新号码');
  } else if (type === '提供方' && newProvider.value) {
    providers.value.push({ name: '新关联医生', license: newProvider.value, type: '未验证' });
    newProvider.value = '';
    ElMessage.success('已关联新提供方');
  } else {
    ElMessage.warning('请输入内容');
  }
};
</script>