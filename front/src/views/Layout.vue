<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="aside">
      <div class="logo">HealthTrack</div>
      <el-menu router :default-active="$route.path" background-color="#001529" text-color="#fff" active-text-color="#409EFF">
        <el-menu-item index="/dashboard"><el-icon><Odometer /></el-icon>概览</el-menu-item>
        <el-menu-item index="/appointments"><el-icon><Calendar /></el-icon>预约挂号</el-menu-item>
        <el-menu-item index="/providers"><el-icon><FirstAidKit /></el-icon>我的医生</el-menu-item>
        <el-menu-item index="/challenges"><el-icon><Trophy /></el-icon>健康挑战</el-menu-item>
        <el-menu-item index="/reports"><el-icon><DataLine /></el-icon>健康报表</el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span>欢迎, {{ userStore.userInfo.name }}</span>
        <el-button link type="danger" @click="handleLogout">退出</el-button>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useUserStore } from '../stores/user'
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
.aside { background-color: #001529; color: white; }
.logo { height: 60px; line-height: 60px; text-align: center; font-size: 20px; font-weight: bold; background: #002140; }
.header { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #eee; background: #fff; }
</style>