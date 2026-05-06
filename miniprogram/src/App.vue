<template>
  <view class="app">
    <offline-badge v-if="!store.isOnline" />
    <router-view />
  </view>
</template>

<script setup lang="ts">
import { onLaunch } from '@dcloudio/uni-app'
import { useAppStore } from './store'
import OfflineBadge from './components/offline-badge/index.vue'

const store = useAppStore()

onLaunch(() => {
  store.restoreSession()

  // Watch network status
  uni.onNetworkStatusChange((res) => {
    store.setOnlineStatus(res.isConnected)
  })
  uni.getNetworkType({ success: (res) => {
    store.setOnlineStatus(res.networkType !== 'none')
  }})
})
</script>

<style>
.app { min-height: 100vh; background: #f5f5f5; }
</style>
