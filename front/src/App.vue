<template>
  <router-view v-slot="{ Component }">
    <transition name="fade" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>

<style>
/* 全局重置与字体 */
:root {
  --primary-color: #007bff; /* 医疗蓝 */
  --bg-color: #f4f6f9;      /* 浅灰背景 */
  --card-radius: 16px;
}

html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  background-color: var(--bg-color);
  font-family: 'Inter', 'Helvetica Neue', Helvetica, 'PingFang SC', 'Microsoft YaHei', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  color: #2c3e50;
}

/* 页面容器通用样式 */
.page-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  animation: slideUp 0.5s ease-out;
}

/* 路由切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Element Plus 覆盖样式 */
.el-card {
  border-radius: var(--card-radius) !important;
  border: none !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05) !important;
}

.el-button--primary {
  background: linear-gradient(135deg, #007bff 0%, #00d2ff 100%) !important;
  border: none !important;
}
.el-button--primary:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}
</style>