<template>
  <div class="menu-container">
    <el-row :gutter="24">
      <el-col :span="8" :xs="24" :sm="12" :md="8" v-for="item in menuItems" :key="item.path" style="margin-bottom: 24px;">
        <el-card 
          shadow="hover" 
          class="menu-card" 
          :body-style="{ padding: '0px', height: '100%', display: 'flex', 'justify-content': 'center', 'align-items': 'center', 'flex-direction': 'column' }"
          @click="$router.push(item.path)"
        >
          <div class="card-content">
            <el-icon :size="48" class="icon"><component :is="item.icon" /></el-icon>
            <h3>{{ item.title }}</h3>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { User, Calendar, Trophy, PieChart, Search, Setting, DataAnalysis } from '@element-plus/icons-vue'

const menuItems = [
  { title: '账户信息', path: '/account-info', icon: 'User' },
  { title: '预约服务', path: '/book-appointment', icon: 'Calendar' },
  { title: '创建健康挑战', path: '/create-challenge', icon: 'Trophy' },
  { title: '月度健康摘要', path: '/monthly-summary', icon: 'PieChart' },
  { title: '搜索记录', path: '/search-records', icon: 'Search' },
  { title: '账户功能管理', path: '/account-functions', icon: 'Setting' },
  { title: '系统统计报表', path: '/summary-functions', icon: 'DataAnalysis' }
]
</script>

<style scoped>
.menu-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto; /* 居中显示整个菜单区域 */
}

/* 卡片基础样式 */
.menu-card {
  cursor: pointer;
  height: 180px; /* 增加高度，更大气 */
  border-radius: 16px; /* 更圆润的边角 */
  border: none;
  overflow: hidden; /* 关键！隐藏滚动条 */
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  background: #fff;
}

/* 内部内容容器 */
.card-content {
  text-align: center;
  transition: all 0.3s ease;
  z-index: 1;
}

/* 图标样式 */
.icon {
  color: var(--vivid-blue-primary); /* 默认使用鲜艳蓝 */
  margin-bottom: 16px;
  transition: all 0.3s ease;
}

/* 标题样式 */
.menu-card h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  transition: all 0.3s ease;
}

/* --- 悬停交互效果 --- */
.menu-card:hover {
  /* 悬停时背景变为渐变色 */
  background: var(--vivid-blue-gradient);
  /* 上浮效果 */
  transform: translateY(-8px);
  /* 炫酷的彩色阴影 */
  box-shadow: var(--hover-shadow) !important;
}

/* 悬停时，强制反转图标和文字颜色为白色 */
.menu-card:hover .icon,
.menu-card:hover h3 {
  color: #ffffff !important;
}
</style>