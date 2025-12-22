<template>
  <div class="dashboard-container">
    <div class="nav-header">
      <div class="brand">
        <el-icon class="brand-icon"><FirstAidKit /></el-icon>
        <span>HealthTrack</span>
      </div>
      <div class="user-profile">
        <el-avatar :size="36" style="background: #007bff">{{ userStore.userInfo?.name?.[0] || 'U' }}</el-avatar>
        <span class="username">欢迎, {{ userStore.userInfo?.name || 'User' }}</span>
        <el-divider direction="vertical" />
        <el-button type="danger" link @click="logout">退出</el-button>
      </div>
    </div>

    <div class="welcome-banner">
      <div class="banner-content">
        <h2>今天感觉怎么样?</h2>
        <p>管理您的健康数据，开启活力每一天。</p>
      </div>
      <el-icon class="banner-icon"><Sunny /></el-icon>
    </div>

    <div class="menu-grid">
      <el-row :gutter="24">
        <el-col :xs="24" :sm="12" :md="8" v-for="(item, index) in menuItems" :key="item.path">
          <div class="menu-card" :class="`card-style-${index % 6}`" @click="$router.push(item.path)">
            <div class="card-icon-wrapper">
              <el-icon><component :is="item.icon" /></el-icon>
            </div>
            <div class="card-info">
              <h3>{{ item.title }}</h3>
              <p>{{ item.desc }}</p>
            </div>
            <el-icon class="arrow-icon"><ArrowRight /></el-icon>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const menuItems = [
  { title: '账户信息', desc: '管理个人档案与联系方式', path: '/account-info', icon: 'User' },
  { title: '预约服务', desc: '在线预约医生与专家', path: '/book-appointment', icon: 'Calendar' },
  { title: '健康挑战', desc: '查看参与的挑战或发起新挑战', path: '/create-challenge', icon: 'Trophy' },
  { title: '月度摘要', desc: '查看健康数据趋势分析', path: '/monthly-summary', icon: 'DataLine' },
  { title: '搜索记录', desc: '查询历史预约与数据', path: '/search-records', icon: 'Search' },
  { title: '统计功能', desc: '平台热门挑战与排行', path: '/summary-functions', icon: 'PieChart' }
]

onMounted(() => userStore.fetchUserInfo())

const logout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding: 0 24px 24px;
}

.nav-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  margin: 0 -24px 24px;
  padding: 0 40px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.03);
  position: sticky;
  top: 0;
  z-index: 100;
}

.brand {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: 800;
  color: #007bff;
  gap: 10px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
}
.username { font-weight: 500; }

.welcome-banner {
  background: linear-gradient(120deg, #007bff, #00d2ff);
  border-radius: 20px;
  padding: 40px;
  color: white;
  margin-bottom: 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 10px 25px rgba(0, 123, 255, 0.3);
}

.banner-content h2 { font-size: 28px; margin: 0 0 8px 0; }
.banner-content p { margin: 0; opacity: 0.9; font-size: 16px; }
.banner-icon { font-size: 80px; opacity: 0.2; }

.menu-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
  border: 1px solid rgba(0,0,0,0.03);
  position: relative;
  overflow: hidden;
}

.menu-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px rgba(0,0,0,0.08);
}

.card-icon-wrapper {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 26px;
  margin-right: 20px;
  flex-shrink: 0;
}

.card-info h3 { margin: 0 0 6px 0; font-size: 18px; color: #333; }
.card-info p { margin: 0; color: #888; font-size: 13px; line-height: 1.4; }

.arrow-icon {
  margin-left: auto;
  color: #ddd;
  transition: transform 0.3s;
}

.menu-card:hover .arrow-icon {
  transform: translateX(5px);
  color: #007bff;
}

/* 不同卡片的配色装饰 */
.card-style-0 .card-icon-wrapper { background: #e3f2fd; color: #1e88e5; }
.card-style-1 .card-icon-wrapper { background: #e8f5e9; color: #43a047; }
.card-style-2 .card-icon-wrapper { background: #fff3e0; color: #fb8c00; }
.card-style-3 .card-icon-wrapper { background: #f3e5f5; color: #8e24aa; }
.card-style-4 .card-icon-wrapper { background: #e0f2f1; color: #00897b; }
.card-style-5 .card-icon-wrapper { background: #ffebee; color: #e53935; }
</style>