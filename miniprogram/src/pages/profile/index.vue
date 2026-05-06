<template>
  <view class="profile-page">
    <view class="user-card">
      <image class="avatar" src="/static/avatar-default.png" mode="aspectFill" />
      <text class="name">{{ store.realName || '用户' }}</text>
      <text class="biz">{{ store.businessName || '未绑定企业' }}</text>
    </view>

    <view class="info-section">
      <view class="section-title">企业信息</view>
      <view class="info-item">
        <text class="label">企业名称</text>
        <text class="value">{{ store.businessName || '-' }}</text>
      </view>
      <view class="info-item">
        <text class="label">绑定状态</text>
        <text class="value status-active" v-if="store.isBound">已绑定</text>
        <text class="value status-pending" v-else>未绑定</text>
      </view>
    </view>

    <view class="menu-section">
      <view class="menu-item" @click="syncOfflineData">
        <text>离线数据同步</text>
        <text class="arrow">></text>
      </view>
      <view class="menu-item" @click="handleLogout">
        <text class="logout-text">退出登录</text>
        <text class="arrow">></text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { useAppStore } from '../../store'
import { getPendingQueue, syncPending } from '../../utils/offline-sync'
import { submitReport } from '../../utils/api'

const store = useAppStore()

async function syncOfflineData() {
  const pending = getPendingQueue()
  if (!pending.length) {
    uni.showToast({ title: '无待同步数据', icon: 'none' })
    return
  }

  uni.showLoading({ title: '同步中...' })
  const synced = await syncPending(async (report) => {
    return submitReport(report.data)
  })
  uni.hideLoading()
  uni.showToast({ title: `已同步 ${synced} 条记录`, icon: 'success' })
}

function handleLogout() {
  uni.showModal({
    title: '提示',
    content: '确定退出登录？',
    success: (res) => {
      if (res.confirm) {
        store.logout()
        uni.reLaunch({ url: '/pages/login/index' })
      }
    }
  })
}
</script>

<style scoped>
.profile-page { min-height: 100vh; background: #f5f5f5; }
.user-card { background: linear-gradient(135deg, #2979FF, #1565C0); padding: 60rpx 30rpx 40rpx; display: flex; flex-direction: column; align-items: center; }
.avatar { width: 120rpx; height: 120rpx; border-radius: 50%; border: 4rpx solid rgba(255,255,255,0.3); margin-bottom: 20rpx; }
.name { font-size: 36rpx; color: #fff; font-weight: bold; }
.biz { font-size: 24rpx; color: rgba(255,255,255,0.8); margin-top: 8rpx; }
.info-section { background: #fff; margin: 20rpx; border-radius: 16rpx; padding: 30rpx; }
.section-title { font-size: 28rpx; color: #999; margin-bottom: 20rpx; }
.info-item { display: flex; justify-content: space-between; padding: 16rpx 0; }
.label { font-size: 28rpx; color: #333; }
.value { font-size: 28rpx; color: #666; }
.status-active { color: #4caf50; }
.status-pending { color: #ff9800; }
.menu-section { background: #fff; margin: 20rpx; border-radius: 16rpx; overflow: hidden; }
.menu-item { display: flex; justify-content: space-between; padding: 30rpx; font-size: 28rpx; color: #333; border-bottom: 1rpx solid #f5f5f5; }
.arrow { color: #ccc; }
.logout-text { color: #f44336; }
</style>
