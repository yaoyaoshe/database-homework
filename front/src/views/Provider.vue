<template>
  <div>
    <h2>医疗服务提供者</h2>
    <el-table :data="providers" style="width: 100%">
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column prop="specialty" label="专科" width="120" />
      <el-table-column prop="qualification" label="资质" />
      <el-table-column prop="contact_info" label="联系方式" show-overflow-tooltip />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" type="success" @click="linkProvider(row.provider_id)">关联医生</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { useUserStore } from '../stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const providers = ref([])
const userStore = useUserStore()

onMounted(async () => {
  providers.value = await request.get('/providers')
})

const linkProvider = async (pid) => {
  ElMessageBox.prompt('请输入关系类型 (如: 主治医生, 家庭医生)', '关联医生', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async ({ value }) => {
    await request.post(`/users/${userStore.userId}/providers`, {
      provider_id: pid,
      relationship_type: value || '咨询医生'
    })
    ElMessage.success('关联成功')
  })
}
</script>