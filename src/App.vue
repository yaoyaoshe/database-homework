<template>
  <el-config-provider>
    <div class="app-container">
      <el-container>
        <el-header v-if="!isLoginPage" class="main-header">
          <div class="brand">HealthTrack 健康平台</div>
          <div class="user-info">
            <span>欢迎, 张三</span>
            <el-button link @click="logout" class="logout-btn">退出登录</el-button>
          </div>
        </el-header>
        
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </div>
  </el-config-provider>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

const isLoginPage = computed(() => route.path === '/login');

const logout = () => {
  router.push('/login');
};
</script>

<style>
/* 1. 定义全局颜色变量：鲜艳的蓝青色渐变 */
:root {
  --vivid-blue-primary: #0072ff;
  --vivid-blue-gradient: linear-gradient(135deg, #00c6ff 0%, #0072ff 100%);
  --hover-shadow: 0 8px 20px rgba(0, 114, 255, 0.3);
  --bg-color: #f0f2f5;
}

body {
  margin: 0;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  background-color: var(--bg-color); /* 背景色稍微深一点，突出白色卡片 */
}

/* 2. 顶部导航栏美化 */
.main-header {
  background: var(--vivid-blue-gradient) !important; /* 使用渐变背景 */
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 2px 10px rgba(0, 114, 255, 0.2); /* 增加层次感阴影 */
  border-bottom: none;
  height: 60px;
}

.brand {
  font-size: 22px;
  font-weight: bold;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.1); /* 文字微立体效果 */
  letter-spacing: 1px;
}

.logout-btn {
  color: rgba(255, 255, 255, 0.9) !important;
  font-weight: 500;
}
.logout-btn:hover {
  color: #fff !important;
}

/* 3. 全局覆盖 Element Plus 主按钮样式 (让所有 primary 按钮都变渐变) */
.el-button--primary {
    background: var(--vivid-blue-gradient) !important;
    border: none !important;
    transition: all 0.3s;
    opacity: 0.9;
}
.el-button--primary:hover {
    opacity: 1;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 198, 255, 0.4);
}
</style>