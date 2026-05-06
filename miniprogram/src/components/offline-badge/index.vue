<template>
  <view class="offline-badge" v-if="!isOnline">
    <text class="icon">!</text>
    <text>当前离线，数据将暂存本地</text>
  </view>
</template>

<script setup lang="ts">
import { useAppStore } from '../../store'
import { computed } from 'vue'

const store = useAppStore()
const isOnline = computed(() => store.isOnline)
</script>

<style scoped>
.offline-badge {
  position: fixed; top: 0; left: 0; right: 0; z-index: 999;
  background: #ff9800; color: #fff; padding: 8rpx 20rpx;
  font-size: 24rpx; text-align: center; display: flex;
  align-items: center; justify-content: center; gap: 8rpx;
}
.icon {
  display: inline-flex; width: 32rpx; height: 32rpx;
  background: #fff; color: #ff9800; border-radius: 50%;
  align-items: center; justify-content: center; font-weight: bold;
}
</style>
