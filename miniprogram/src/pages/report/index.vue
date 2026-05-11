<template>
  <view class="report-page">
    <view class="photo-preview" v-if="photoUrl">
      <image :src="photoUrl" mode="widthFix" class="receipt-img" />
    </view>

    <view class="form-section">
      <text class="section-title">本月营业额申报</text>

      <view class="form-item">
        <text class="item-label">营业额 (元)</text>
        <input class="item-input" v-model.number="totalRevenue" type="digit" placeholder="请输入营业额" />
      </view>

      <view class="photo-info" v-if="!photoUrl">
        <text class="photo-hint">无附件照片</text>
      </view>
    </view>

    <view class="submit-section">
      <button class="submit-btn" @click="handleSubmit" :loading="submitting">确认上报</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useAppStore } from '../../store'
import { submitReport } from '../../utils/api'
import { enqueueReport } from '../../utils/offline-sync'

const store = useAppStore()
const submitting = ref(false)
const photoUrl = ref('')
const totalRevenue = ref(0)
const categories = ref<Record<string, number>>({})

onLoad((options: any) => {
  if (options?.data) {
    const data = JSON.parse(decodeURIComponent(options.data))
    totalRevenue.value = data.total_revenue || 0
    photoUrl.value = data.photo_url || data.photo_local || ''
    categories.value = {
      restaurant_revenue: data.restaurant_revenue || 0,
      retail_revenue: data.retail_revenue || 0,
      accommodation_revenue: data.accommodation_revenue || 0,
      tobacco_alcohol_revenue: data.tobacco_alcohol_revenue || 0,
      other_revenue: data.other_revenue || 0,
    }
    // Save raw text locally for reference (not submitted)
    if (data.raw_text) {
      uni.setStorageSync('last_ocr_text', data.raw_text)
    }
  }
})

async function handleSubmit() {
  if (totalRevenue.value <= 0) {
    uni.showToast({ title: '请输入营业额', icon: 'none' })
    return
  }

  submitting.value = true

  const reportData = {
    business_id: store.businessId,
    total_revenue: totalRevenue.value,
    restaurant_revenue: categories.value.restaurant_revenue || 0,
    retail_revenue: categories.value.retail_revenue || 0,
    accommodation_revenue: categories.value.accommodation_revenue || 0,
    tobacco_alcohol_revenue: categories.value.tobacco_alcohol_revenue || 0,
    other_revenue: categories.value.other_revenue || 0,
    source_type: 'OCR',
    photo_urls: photoUrl.value ? [photoUrl.value] : []
  }

  try {
    if (!store.isOnline) {
      enqueueReport({
        id: Date.now().toString(),
        businessId: store.businessId,
        data: reportData,
        photoUrls: photoUrl.value ? [photoUrl.value] : [],
        createdAt: Date.now(),
        synced: false
      })
      uni.showToast({ title: '已暂存，联网后自动上报', icon: 'none' })
    } else {
      await submitReport(reportData)
      uni.showToast({ title: '上报成功', icon: 'success' })
    }

    setTimeout(() => {
      uni.reLaunch({ url: '/pages/home/index' })
    }, 1500)
  } catch (err: any) {
    uni.showToast({ title: err.data?.error || '上报失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.report-page { min-height: 100vh; background: #f5f5f5; padding: 30rpx; }
.photo-preview { margin-bottom: 20rpx; }
.receipt-img { width: 100%; border-radius: 12rpx; }
.form-section { background: #fff; border-radius: 16rpx; padding: 30rpx; margin-bottom: 30rpx; }
.section-title { font-size: 30rpx; font-weight: bold; color: #333; margin-bottom: 24rpx; display: block; }
.form-item { display: flex; align-items: center; justify-content: space-between; padding: 20rpx 0; }
.item-label { font-size: 28rpx; color: #333; flex-shrink: 0; }
.item-input { text-align: right; font-size: 32rpx; color: #333; width: 400rpx; border-bottom: 2rpx solid #e0e0e0; padding: 10rpx 0; }
.photo-info { margin-top: 30rpx; padding: 30rpx; background: #f9f9f9; border-radius: 12rpx; text-align: center; }
.photo-hint { font-size: 26rpx; color: #999; }
.submit-section { text-align: center; padding: 40rpx 0; }
.submit-btn { width: 100%; height: 90rpx; background: #2979FF; color: #fff; border-radius: 50rpx; border: none; font-size: 32rpx; }
.warning { color: #f44336; font-size: 24rpx; margin-bottom: 20rpx; display: block; }
</style>
