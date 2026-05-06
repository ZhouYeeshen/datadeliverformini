<template>
  <el-container class="app-layout">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <h2>经营数据上报系统</h2>
        <span>商务局管理端</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>综合态势</span>
        </el-menu-item>
        <el-menu-item index="/report">
          <el-icon><Document /></el-icon>
          <span>申报记录</span>
        </el-menu-item>
        <el-menu-item index="/business">
          <el-icon><OfficeBuilding /></el-icon>
          <span>商户管理</span>
        </el-menu-item>
        <el-menu-item index="/alert">
          <el-icon><Bell />
            <el-badge v-if="unreadCount" :value="unreadCount" class="badge" />
          </el-icon>
          <span>预警专栏</span>
        </el-menu-item>
        <el-menu-item index="/statistics">
          <el-icon><TrendCharts /></el-icon>
          <span>统计分析</span>
        </el-menu-item>
        <el-menu-item index="/pos-interface">
          <el-icon><Connection /></el-icon>
          <span>接口管理</span>
        </el-menu-item>
        <el-menu-item index="/system">
          <el-icon><Setting /></el-icon>
          <span>系统设置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="topbar">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentTitle">{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <span class="user">{{ store.username }}</span>
          <el-button type="danger" text @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from './store'
import { getUnreadAlertCount } from './api'

const route = useRoute()
const router = useRouter()
const store = useAdminStore()
const unreadCount = ref(0)

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta?.title as string || '')

async function loadUnreadCount() {
  try {
    const res = await getUnreadAlertCount()
    unreadCount.value = res.data.unread_count
  } catch {}
}

function handleLogout() {
  store.logout()
  router.push('/login')
}

onMounted(() => {
  loadUnreadCount()
  setInterval(loadUnreadCount, 30000) // refresh every 30s
})
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }

/* Sci-fi tech theme variables */
:root {
  --tech-dark: #0a1628;
  --tech-primary: #00d4ff;
  --tech-accent: #7c3aed;
  --tech-gradient: linear-gradient(135deg, #0a1628 0%, #1a1a3e 50%, #0d1b3e 100%);
  --glow-blue: 0 0 20px rgba(0, 212, 255, 0.15);
  --glow-purple: 0 0 20px rgba(124, 58, 237, 0.15);
}

body {
  background: #0a0f1a;
}

.app-layout { height: 100vh; }

/* Sidebar — deep space gradient */
.sidebar {
  background: var(--tech-gradient) !important;
  overflow-y: auto;
  border-right: 1px solid rgba(0, 212, 255, 0.08);
  position: relative;
}
.sidebar::after {
  content: '';
  position: absolute;
  top: 0; right: 0; bottom: 0; left: 0;
  background:
    radial-gradient(1px 1px at 20% 30%, rgba(0,212,255,0.3), transparent),
    radial-gradient(1px 1px at 40% 70%, rgba(124,58,237,0.3), transparent),
    radial-gradient(1px 1px at 80% 20%, rgba(0,212,255,0.25), transparent),
    radial-gradient(1px 1px at 60% 85%, rgba(124,58,237,0.25), transparent);
  pointer-events: none;
}

/* Logo with glow */
.logo {
  padding: 24px 16px;
  text-align: center;
  border-bottom: 1px solid rgba(0, 212, 255, 0.12);
  position: relative;
}
.logo h2 {
  background: linear-gradient(135deg, #00d4ff, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 2px;
}
.logo span {
  color: rgba(167, 139, 250, 0.6);
  font-size: 11px;
  letter-spacing: 3px;
  text-transform: uppercase;
}

/* Menu items */
.el-menu {
  border-right: none !important;
  background: transparent !important;
}
.el-menu-item {
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
}
.el-menu-item:hover {
  background: rgba(0, 212, 255, 0.08) !important;
  box-shadow: var(--glow-blue);
}
.el-menu-item.is-active {
  background: linear-gradient(135deg, rgba(0, 212, 255, 0.15), rgba(124, 58, 237, 0.15)) !important;
  box-shadow: var(--glow-blue);
  color: #00d4ff !important;
  border-left: 3px solid #00d4ff;
}

.badge { margin-left: -8px; }

/* Top header bar */
.topbar {
  background: linear-gradient(180deg, rgba(10, 22, 40, 0.95), rgba(13, 27, 62, 0.9));
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(0, 212, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  color: #c8d6e5;
}
.header-left { display: flex; align-items: center; }
.header-right { display: flex; align-items: center; gap: 16px; }
.user { color: rgba(167, 139, 250, 0.8); font-size: 14px; }

/* Main content area */
.main-content {
  background: #0a0f1a;
  background-image:
    radial-gradient(circle at 25% 25%, rgba(0, 212, 255, 0.03) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(124, 58, 237, 0.03) 0%, transparent 50%);
  padding: 24px;
  min-height: calc(100vh - 60px);
}

/* Global card override for tech look */
:deep(.el-card) {
  background: rgba(15, 23, 42, 0.8) !important;
  border: 1px solid rgba(0, 212, 255, 0.08) !important;
  border-radius: 12px !important;
  backdrop-filter: blur(8px);
  color: #c8d6e5 !important;
}
:deep(.el-card__header) {
  border-bottom: 1px solid rgba(0, 212, 255, 0.08) !important;
  color: #e2e8f0 !important;
}
:deep(.el-card__body) {
  color: #c8d6e5 !important;
}

/* Table override */
:deep(.el-table) {
  --el-table-bg-color: transparent;
  --el-table-tr-bg-color: transparent;
  --el-table-header-bg-color: rgba(0, 212, 255, 0.05);
  --el-table-border-color: rgba(0, 212, 255, 0.08);
  --el-table-text-color: #c8d6e5;
  --el-table-header-text-color: #94a3b8;
  --el-table-row-hover-bg-color: rgba(0, 212, 255, 0.05);
}
:deep(.el-table th) {
  background: rgba(0, 212, 255, 0.06) !important;
}
:deep(.el-table td) {
  border-bottom-color: rgba(0, 212, 255, 0.05) !important;
}
:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background: rgba(0, 212, 255, 0.02) !important;
}

/* Form input overrides */
:deep(.el-input__wrapper) {
  background: rgba(0, 212, 255, 0.04) !important;
  border: 1px solid rgba(0, 212, 255, 0.12) !important;
  box-shadow: none !important;
}
:deep(.el-input__inner) { color: #c8d6e5 !important; }
:deep(.el-select .el-input__wrapper) { background: rgba(0, 212, 255, 0.04) !important; }

/* Button overrides */
:deep(.el-button--primary) {
  background: linear-gradient(135deg, #00d4ff, #7c3aed) !important;
  border: none !important;
}
:deep(.el-button--primary:hover) {
  box-shadow: 0 0 20px rgba(0, 212, 255, 0.3) !important;
}

/* Dialog override */
:deep(.el-dialog) {
  background: rgba(15, 23, 42, 0.95) !important;
  border: 1px solid rgba(0, 212, 255, 0.1) !important;
  border-radius: 12px !important;
  backdrop-filter: blur(16px);
}
:deep(.el-dialog__title) { color: #e2e8f0 !important; }
:deep(.el-dialog__body) { color: #c8d6e5 !important; }

/* Pagination */
:deep(.el-pagination .btn-prev, .el-pagination .btn-next) {
  background: rgba(0, 212, 255, 0.05) !important;
}
:deep(.el-pager li) { color: #94a3b8 !important; }
:deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, #00d4ff, #7c3aed) !important;
}

/* Breadcrumb */
:deep(.el-breadcrumb__inner) { color: #94a3b8 !important; }
:deep(.el-breadcrumb__inner.is-link) { color: #00d4ff !important; }

/* Tag colors */
:deep(.el-tag--success) {
  background: rgba(0, 212, 255, 0.12) !important;
  border-color: rgba(0, 212, 255, 0.3) !important;
  color: #00d4ff !important;
}
</style>
