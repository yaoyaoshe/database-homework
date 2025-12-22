<template>
  <div>
    <h2>寻找医生</h2>
    <div class="filter-bar">
      <el-switch v-model="verifiedOnly" active-text="仅显示已认证" @change="fetchProviders" />
    </div>

    <el-row :gutter="20">
      <el-col :span="8" v-for="p in providers" :key="p.provider_id" style="margin-bottom: 20px;">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ p.name }}</span>
              <el-tag v-if="p.is_verified" type="success" size="small">已认证</el-tag>
            </div>
          </template>
          <p><strong>专业:</strong> {{ p.specialty }}</p>
          <p><strong>资质:</strong> {{ p.qualification }}</p>
          <el-button type="primary" size="small" class="w-100" @click="openLinkDialog(p)">关联为我的医生</el-button>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="linkDialogVisible" title="关联医生">
      <p>正在关联医生: {{ selectedProvider?.name }}</p>
      <el-select v-model="relationType" placeholder="选择关系类型" class="w-100">
        <el-option label="主治医生" value="主治医生" />
        <el-option label="家庭医生" value="家庭医生" />
        <el-option label="专科医生" value="专科医生" />
      </el-select>
      <template #footer>
        <el-button type="primary" @click="confirmLink">确认关联</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProviders, linkProvider } from '@/api/all'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const providers = ref([])
const verifiedOnly = ref(false)
const linkDialogVisible = ref(false)
const selectedProvider = ref(null)
const relationType = ref('主治医生')

const fetchProviders = async () => {
  const params = {}
  if(verifiedOnly.value) params.verified = 'true'
  providers.value = await getProviders(params)
}

const openLinkDialog = (provider) => {
  selectedProvider.value = provider
  linkDialogVisible.value = true
}

const confirmLink = async () => {
  await linkProvider(userStore.userInfo.user_id, {
    provider_id: selectedProvider.value.provider_id,
    relationship_type: relationType.value
  })
  ElMessage.success('关联成功')
  linkDialogVisible.value = false
}

onMounted(fetchProviders)
</script>

<style scoped>
.filter-bar { margin-bottom: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
.w-100 { width: 100%; margin-top: 10px; }
</style>