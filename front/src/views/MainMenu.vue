<template>
  <div class="main-menu">
    <el-header class="header">
      <div class="logo">HealthTrack</div>
      <div class="user-info">欢迎, {{ userStore.userInfo?.name || 'User ' + userStore.userId }}
        <el-button link @click="logout">退出登录</el-button>
      </div>
    </el-header>

    <el-main>
      <el-row :gutter="20">
        <el-col :span="8" v-for="item in menuItems" :key="item.path">
          <el-card class="menu-card" shadow="hover" @click="$router.push(item.path)">
            <el-icon :size="40"><component :is="item.icon" /></el-icon>
            <h3>{{ item.title }}</h3>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const menuItems = [
  { title: '账户信息', path: '/account-info', icon: 'User' },
  { title: '预约服务', path: '/book-appointment', icon: 'Calendar' },
  { title: '创建挑战', path: '/create-challenge', icon: 'Trophy' },
  { title: '月度摘要', path: '/monthly-summary', icon: 'DataLine' },
  { title: '搜索记录', path: '/search-records', icon: 'Search' },
  { title: '统计功能', path: '/summary-functions', icon: 'PieChart' }
]

onMounted(() => userStore.fetchUserInfo())

const logout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.header { background: #409EFF; color: white; display: flex; justify-content: space-between; align-items: center; padding: 0 20px; }
.menu-card { text-align: center; cursor: pointer; margin-bottom: 20px; padding: 20px; transition: transform 0.2s; }
.menu-card:hover { transform: translateY(-5px); }
.logo { font-size: 24px; font-weight: bold; }
</style>