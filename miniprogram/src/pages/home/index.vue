<template>
  <view class="home-page">
    <view class="header-card">
      <text class="biz-name">{{ store.businessName || '未绑定企业' }}</text>
      <text class="greeting">欢迎回来，{{ store.realName || '用户' }}</text>
    </view>

    <view class="status-card">
      <view class="status-header">
        <text class="month">{{ currentMonth }}月申报状态</text>
      </view>
      <view class="status-body">
        <template v-if="reported">
          <view class="status-tag reported">已申报</view>
          <text class="total">¥{{ reportData?.total_revenue?.toFixed(2) || '0.00' }}</text>
          <text class="source">来源: {{ reportData?.source_type || '-' }}</text>
        </template>
        <template v-else>
          <view class="status-tag pending">待申报</view>
          <text class="hint">本月尚未申报，请尽快上报</text>
          <navigator url="/pages/scan/index" class="report-btn">立即申报</navigator>
        </template>
      </view>
    </view>

    <view class="quick-actions">
      <navigator url="/pages/scan/index" class="action-item">
        <view class="action-icon scan">S</view>
        <text>扫描上报</text>
      </navigator>
      <navigator url="/pages/history/index" class="action-item">
        <view class="action-icon history">H</view>
        <text>历史记录</text>
      </navigator>
      <navigator url="/pages/profile/index" class="action-item">
        <view class="action-icon profile">P</view>
        <text>企业信息</text>
      </navigator>
    </view>

    <view class="recent-section" v-if="recentReports.length">
      <text class="section-title">最近上报</text>
      <view class="report-item" v-for="r in recentReports" :key="r.id">
        <view class="report-info">
          <text class="report-month">{{ r.report_month }}</text>
          <text class="report-total">¥{{ r.total_revenue?.toFixed(2) }}</text>
        </view>
        <text class="report-status">{{ r.status === 'submitted' ? '已上报' : r.status }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useAppStore } from '../../store'
import { getCurrentMonthReport, getReportHistory } from '../../utils/api'

const store = useAppStore()
const reported = ref(false)
const reportData = ref<any>(null)
const recentReports = ref<any[]>([])

const currentMonth = new Date().getMonth() + 1

async function loadData() {
  try {
    const res = await getCurrentMonthReport()
    reported.value = res.reported
    reportData.value = res.report
  } catch {}

  try {
    const hRes = await getReportHistory()
    recentReports.value = (hRes.list || []).slice(0, 5)
  } catch {}
}

onShow(() => loadData())
</script>

<style scoped>
.home-page { min-height: 100vh; background: #f5f5f5; padding: 30rpx; }
.header-card { background: linear-gradient(135deg, #2979FF, #1565C0); border-radius: 20rpx; padding: 40rpx; margin-bottom: 30rpx; }
.biz-name { font-size: 36rpx; color: #fff; font-weight: bold; display: block; }
.greeting { font-size: 24rpx; color: rgba(255,255,255,0.8); margin-top: 8rpx; }
.status-card { background: #fff; border-radius: 20rpx; padding: 30rpx; margin-bottom: 30rpx; }
.status-header { margin-bottom: 20rpx; }
.month { font-size: 30rpx; font-weight: bold; color: #333; }
.status-body { display: flex; flex-direction: column; align-items: center; padding: 20rpx 0; }
.status-tag { padding: 8rpx 30rpx; border-radius: 30rpx; font-size: 24rpx; margin-bottom: 16rpx; }
.status-tag.reported { background: #e8f5e9; color: #4caf50; }
.status-tag.pending { background: #fff3e0; color: #ff9800; }
.total { font-size: 48rpx; font-weight: bold; color: #333; }
.source { font-size: 22rpx; color: #999; margin-top: 8rpx; }
.hint { font-size: 26rpx; color: #999; margin-bottom: 20rpx; }
.report-btn { background: #2979FF; color: #fff; padding: 16rpx 60rpx; border-radius: 40rpx; font-size: 28rpx; text-decoration: none; }
.quick-actions { display: flex; gap: 20rpx; margin-bottom: 30rpx; }
.action-item { flex: 1; background: #fff; border-radius: 16rpx; padding: 30rpx; text-align: center; font-size: 24rpx; color: #666; }
.action-icon { width: 60rpx; height: 60rpx; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-weight: bold; margin-bottom: 12rpx; font-size: 28rpx; }
.action-icon.scan { background: #e3f2fd; color: #2979FF; }
.action-icon.history { background: #f3e5f5; color: #9c27b0; }
.action-icon.profile { background: #e8f5e9; color: #4caf50; }
.recent-section { background: #fff; border-radius: 20rpx; padding: 30rpx; }
.section-title { font-size: 30rpx; font-weight: bold; color: #333; margin-bottom: 20rpx; display: block; }
.report-item { display: flex; justify-content: space-between; align-items: center; padding: 20rpx 0; border-bottom: 1rpx solid #f0f0f0; }
.report-month { font-size: 28rpx; color: #333; }
.report-total { font-size: 24rpx; color: #999; display: block; }
.report-status { font-size: 24rpx; color: #4caf50; }
</style>
