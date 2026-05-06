<template>
  <view class="history-page">
    <view class="list" v-if="reports.length">
      <view class="report-card" v-for="item in reports" :key="item.id">
        <view class="card-header">
          <text class="month">{{ item.report_month }}</text>
          <text class="status" :class="item.status">{{ item.status === 'submitted' ? '已上报' : item.status }}</text>
        </view>
        <view class="card-body">
          <text class="total-label">营业额</text>
          <text class="total-value">¥{{ item.total_revenue?.toFixed(2) || '0.00' }}</text>
        </view>
        <view class="card-footer">
          <text class="source">{{ item.source_type }}</text>
          <text class="photo-link" v-if="item.photo_urls" @click="previewPhoto(item.photo_urls)">查看附件</text>
        </view>
      </view>
    </view>
    <view class="empty" v-else>
      <text>暂无上报记录</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getReportHistory, SERVER_HOST } from '../../utils/api'

const reports = ref<any[]>([])

async function loadHistory() {
  try {
    const res = await getReportHistory()
    reports.value = res.list || []
  } catch {}
}

function previewPhoto(urls: string) {
  const photos = urls.split(',').filter(Boolean)
  if (photos.length === 0) return
  const fullUrls = photos.map(p => p.startsWith('http') ? p : SERVER_HOST + p)
  uni.previewImage({
    urls: fullUrls,
    current: fullUrls[0]
  })
}

onShow(() => loadHistory())
</script>

<style scoped>
.history-page { min-height: 100vh; background: #f5f5f5; padding: 30rpx; }
.list { display: flex; flex-direction: column; gap: 20rpx; }
.report-card { background: #fff; border-radius: 16rpx; padding: 24rpx; }
.card-header { display: flex; justify-content: space-between; margin-bottom: 16rpx; }
.month { font-size: 30rpx; font-weight: bold; color: #333; }
.status { font-size: 22rpx; padding: 4rpx 16rpx; border-radius: 20rpx; }
.status.submitted { background: #e8f5e9; color: #4caf50; }
.card-body { display: flex; justify-content: space-between; align-items: center; padding: 10rpx 0; }
.total-label { font-size: 26rpx; color: #999; }
.total-value { font-size: 36rpx; font-weight: bold; color: #333; }
.card-footer { display: flex; justify-content: space-between; margin-top: 16rpx; padding-top: 16rpx; border-top: 1rpx solid #f0f0f0; }
.source { font-size: 22rpx; color: #999; }
.photo-link { font-size: 22rpx; color: #2979FF; }
.empty { text-align: center; padding: 200rpx 0; color: #999; font-size: 28rpx; }
</style>
