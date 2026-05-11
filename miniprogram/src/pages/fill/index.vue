<template>
  <view class="fill-page">
    <view class="header">
      <text class="header-title">手动填报</text>
      <text class="header-sub">直接输入经营数据，无需扫描单据</text>
    </view>

    <!-- Photo (optional) -->
    <view class="photo-section" @click="choosePhoto">
      <image v-if="photoUrl" :src="photoUrl" class="photo-preview" mode="aspectFit" />
      <view v-else class="photo-placeholder">
        <text class="photo-icon">+</text>
        <text class="photo-text">上传单据照片（选填）</text>
      </view>
      <view v-if="photoUrl" class="photo-remove" @click.stop="photoUrl = ''">
        <text>×</text>
      </view>
    </view>

    <!-- Category inputs -->
    <view class="section">
      <text class="section-title">收入分类</text>
      <view class="form-card">
        <view class="form-row" v-for="row in rows" :key="row.key">
          <text class="form-label">{{ row.label }}</text>
          <view class="form-input-wrap">
            <text class="yen-sign">¥</text>
            <input
              class="form-input"
              type="digit"
              v-model.number="row.value"
              placeholder="0.00"
              :focus="row.key === 'restaurant_revenue'"
            />
          </view>
        </view>
      </view>
    </view>

    <!-- Auto-calculated total -->
    <view class="total-bar">
      <text class="total-label">合计金额</text>
      <text class="total-value">¥{{ totalRevenue.toFixed(2) }}</text>
    </view>

    <!-- Submit -->
    <button class="submit-btn" :disabled="totalRevenue <= 0 || submitting" @click="handleSubmit" :loading="submitting">
      {{ submitting ? '提交中...' : '确认上报' }}
    </button>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useAppStore } from '../../store'
import { submitReport, uploadPhoto } from '../../utils/api'
import { enqueueReport } from '../../utils/offline-sync'

const store = useAppStore()
const submitting = ref(false)
const photoUrl = ref('')
const photoLocal = ref('')
const rows = reactive([
  { key: 'restaurant_revenue', label: '餐饮收入', value: 0 },
  { key: 'retail_revenue', label: '零售收入', value: 0 },
  { key: 'accommodation_revenue', label: '住宿收入', value: 0 },
  { key: 'tobacco_alcohol_revenue', label: '烟酒收入', value: 0 },
  { key: 'other_revenue', label: '其他收入', value: 0 },
])

const totalRevenue = computed(() => {
  const sum = rows.reduce((s, r) => s + (r.value || 0), 0)
  return Math.round(sum * 100) / 100
})

function categoryValue(key: string): number {
  const row = rows.find(r => r.key === key)
  return row ? (row.value || 0) : 0
}

async function choosePhoto() {
  const res = await uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album'],
  })
  const path = res.tempFilePaths[0]
  photoLocal.value = path

  try {
    const uploadRes = await uploadPhoto(path)
    if ((uploadRes as any).url) {
      photoUrl.value = (uploadRes as any).url
    }
  } catch {
    uni.showToast({ title: '图片上传失败，可稍后重试', icon: 'none' })
  }
}

async function handleSubmit() {
  if (totalRevenue.value <= 0) {
    uni.showToast({ title: '请输入至少一项金额', icon: 'none' })
    return
  }

  submitting.value = true
  try {
    const data: any = {
      business_id: store.businessId,
      total_revenue: totalRevenue.value,
      restaurant_revenue: categoryValue('restaurant_revenue'),
      retail_revenue: categoryValue('retail_revenue'),
      accommodation_revenue: categoryValue('accommodation_revenue'),
      tobacco_alcohol_revenue: categoryValue('tobacco_alcohol_revenue'),
      other_revenue: categoryValue('other_revenue'),
      source_type: 'manual',
      photo_urls: photoUrl.value ? [photoUrl.value] : [],
    }

    if (!store.isOnline) {
      await enqueueReport(data)
      uni.showToast({ title: '已保存到本地，联网后自动上报', icon: 'none' })
      setTimeout(() => uni.switchTab({ url: '/pages/home/index' }), 1200)
      return
    }

    await submitReport(data)
    uni.showToast({ title: '上报成功', icon: 'success' })
    setTimeout(() => uni.switchTab({ url: '/pages/home/index' }), 1200)
  } catch (err: any) {
    uni.showToast({ title: err.errMsg || '提交失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.fill-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 30rpx;
}
.header {
  padding: 20rpx 0 30rpx;
}
.header-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  display: block;
}
.header-sub {
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
}
.photo-section {
  position: relative;
  margin-bottom: 30rpx;
}
.photo-placeholder {
  background: #fff;
  border: 2rpx dashed #ccc;
  border-radius: 16rpx;
  padding: 50rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.photo-icon {
  font-size: 60rpx;
  color: #ccc;
}
.photo-text {
  font-size: 24rpx;
  color: #999;
  margin-top: 12rpx;
}
.photo-preview {
  width: 100%;
  height: 300rpx;
  border-radius: 16rpx;
}
.photo-remove {
  position: absolute;
  top: -12rpx;
  right: -12rpx;
  width: 48rpx;
  height: 48rpx;
  background: #ff4444;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 32rpx;
}
.section {
  margin-bottom: 30rpx;
}
.section-title {
  font-size: 28rpx;
  color: #666;
  margin-bottom: 16rpx;
  display: block;
}
.form-card {
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}
.form-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}
.form-row:last-child {
  border-bottom: none;
}
.form-label {
  font-size: 30rpx;
  color: #333;
  flex-shrink: 0;
}
.form-input-wrap {
  display: flex;
  align-items: center;
}
.yen-sign {
  font-size: 32rpx;
  color: #999;
  margin-right: 8rpx;
}
.form-input {
  text-align: right;
  font-size: 32rpx;
  width: 200rpx;
}
.total-bar {
  background: linear-gradient(135deg, #2979FF, #1565C0);
  border-radius: 16rpx;
  padding: 36rpx 30rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40rpx;
}
.total-label {
  font-size: 30rpx;
  color: rgba(255,255,255,0.9);
}
.total-value {
  font-size: 44rpx;
  font-weight: bold;
  color: #fff;
}
.submit-btn {
  width: 100%;
  height: 100rpx;
  background: linear-gradient(135deg, #2979FF, #1565C0);
  color: #fff;
  font-size: 32rpx;
  border-radius: 50rpx;
  border: none;
}
.submit-btn[disabled] {
  background: #ccc;
  color: #999;
}
</style>
