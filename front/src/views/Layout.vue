<template>
  <el-container class="layout-container">
    <el-aside width="220px">
      <el-menu
        router
        :default-active="$route.path"
        class="el-menu-vertical"
        background-color="#304156"
        text-color="#fff"
        active-text-color="#409EFF"
      >
        <div class="logo">HealthTrack</div>
        <el-menu-item index="/appointments">
          <el-icon><Calendar /></el-icon>
          <span>预约管理</span>
        </el-menu-item>
        <el-menu-item index="/providers">
          <el-icon><FirstAidKit /></el-icon>
          <span>寻找医生</span>
        </el-menu-item>
        <el-menu-item index="/challenges">
          <el-icon><Trophy /></el-icon>
          <span>健康挑战</span>
        </el-menu-item>
        <el-menu-item index="/reports">
          <el-icon><DataLine /></el-icon>
          <span>月度报告</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header class="header">
        <div class="user-info">
          <span>欢迎, {{ userStore.userInfo.name }}</span>
          <el-button link type="danger" @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout-container { height: 100vh; }
.el-menu-vertical { height: 100%; border-right: none; }
.logo { height: 60px; line-height: 60px; text-align: center; color: white; font-weight: bold; font-size: 20px; }
.header { background: #fff; border-bottom: 1px solid #dcdfe6; display: flex; align-items: center; justify-content: flex-end; padding-right: 20px;}
.user-info span { margin-right: 15px; font-weight: 500; }
</style>