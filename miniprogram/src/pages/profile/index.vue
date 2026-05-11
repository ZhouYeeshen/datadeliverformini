<template>
  <!-- ========== 未注册视图 ========== -->
  <view class="profile-page" v-if="!store.isBound">
    <view class="register-header">
      <text class="register-icon">👤</text>
      <text class="register-title">请先注册</text>
      <text class="register-desc">绑定微信与企业信息，开始使用经营数据上报</text>
    </view>

    <!-- 手机号快速关联 -->
    <view class="quick-bind-card">
      <view class="card-title">已有企业？手机号快速关联</view>
      <view class="form-group">
        <text class="label">手机号</text>
        <input class="input" v-model="quickPhone" type="number" maxlength="11" placeholder="输入注册时预留的手机号" />
      </view>
      <view class="form-group">
        <text class="label">您的姓名</text>
        <input class="input" v-model="quickName" placeholder="请输入您的姓名" />
      </view>
      <button class="quick-btn" @click="handleQuickBind" :loading="quickLoading">快速关联已有企业</button>
    </view>

    <view class="divider"><text>或注册新企业</text></view>

    <!-- 完整注册表单 -->
    <view class="form-card">
      <view class="form-group">
        <text class="label">企业名称 <text class="required">*</text></text>
        <input class="input" v-model="form.businessName" placeholder="请输入企业名称" />
      </view>
      <view class="form-group">
        <text class="label">行业类型 <text class="required">*</text></text>
        <picker :range="industryOptions" @change="onIndustryChange">
          <view class="picker-row">
            {{ form.industryType || '请选择行业类型' }}
            <text class="arrow">›</text>
          </view>
        </picker>
      </view>
      <view class="form-group">
        <text class="label">法人姓名 <text class="required">*</text></text>
        <input class="input" v-model="form.legalPerson" placeholder="请输入法人姓名" />
      </view>
      <view class="form-group">
        <text class="label">您的姓名 <text class="required">*</text></text>
        <input class="input" v-model="form.realName" placeholder="请输入您的姓名" />
      </view>
      <view class="form-group">
        <text class="label">手机号</text>
        <input class="input" v-model="form.phone" type="number" maxlength="11" placeholder="请输入手机号" />
      </view>
      <view class="form-group">
        <text class="label">营业执照号</text>
        <input class="input" v-model="form.licenseNo" placeholder="选填" />
      </view>
      <button class="submit-btn" @click="handleBind" :loading="loading">确认注册并绑定</button>
    </view>
  </view>

  <!-- ========== 已注册视图 ========== -->
  <view class="profile-page" v-else>
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
        <text class="value status-active">已绑定</text>
      </view>
    </view>

    <view class="menu-section">
      <view class="menu-item" @click="syncOfflineData">
        <text>离线数据同步</text>
        <text class="arrow">›</text>
      </view>
      <view class="menu-item" @click="handleLogout">
        <text class="logout-text">退出登录</text>
        <text class="arrow">›</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAppStore } from '../../store'
import { bindBusiness } from '../../utils/api'
import { getPendingQueue, syncPending } from '../../utils/offline-sync'
import { submitReport } from '../../utils/api'

const store = useAppStore()

// ---- 注册相关 ----
const loading = ref(false)
const quickLoading = ref(false)
const quickPhone = ref('')
const quickName = ref('')
const industryOptions = ['餐饮', '零售', '住宿', '烟酒', '混合']

const form = reactive({
  businessName: '',
  industryType: '',
  legalPerson: '',
  realName: '',
  phone: '',
  licenseNo: ''
})

function onIndustryChange(e: any) {
  form.industryType = industryOptions[e.detail.value]
}

async function doBind(data: any) {
  const codeRes = await uni.login()
  const bindData = { openid: codeRes.code, ...data }
  const res = await bindBusiness(bindData)
  store.setToken(res.token)
  store.setBusiness(res)
  uni.showToast({ title: '注册成功', icon: 'success' })
  setTimeout(() => uni.switchTab({ url: '/pages/home/index' }), 800)
}

async function handleQuickBind() {
  if (!quickPhone.value || !quickName.value) {
    uni.showToast({ title: '请填写手机号和姓名', icon: 'none' })
    return
  }
  if (!/^1\d{10}$/.test(quickPhone.value)) {
    uni.showToast({ title: '手机号格式不正确', icon: 'none' })
    return
  }
  quickLoading.value = true
  try {
    await doBind({ phone: quickPhone.value, real_name: quickName.value })
  } catch (err: any) {
    const msg = err.data?.error || err.errMsg || '未找到关联企业'
    uni.showToast({ title: msg + '，请注册新企业', icon: 'none' })
  } finally {
    quickLoading.value = false
  }
}

async function handleBind() {
  if (!form.businessName || !form.industryType || !form.legalPerson || !form.realName) {
    uni.showToast({ title: '请填写必填项', icon: 'none' })
    return
  }
  loading.value = true
  try {
    await doBind({
      business_name: form.businessName,
      industry_type: form.industryType,
      legal_person: form.legalPerson,
      real_name: form.realName,
      phone: form.phone,
      license_no: form.licenseNo
    })
  } catch (err: any) {
    uni.showToast({ title: err.errMsg || '注册失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

// ---- 已注册功能 ----
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

/* ---- 注册视图 ---- */
.register-header {
  background: linear-gradient(135deg, #2979FF, #1565C0);
  padding: 60rpx 30rpx 50rpx;
  text-align: center;
}
.register-icon { font-size: 60rpx; display: block; margin-bottom: 16rpx; }
.register-title { font-size: 36rpx; font-weight: bold; color: #fff; display: block; }
.register-desc { font-size: 24rpx; color: rgba(255,255,255,0.8); margin-top: 12rpx; display: block; }

.quick-bind-card { background: #e3f2fd; border-radius: 16rpx; padding: 30rpx; margin: 20rpx; border: 2rpx solid #90caf9; }
.card-title { font-size: 26rpx; color: #1565C0; margin-bottom: 20rpx; font-weight: 500; text-align: center; }
.quick-btn {
  width: 100%; height: 80rpx; background: #1565C0; color: #fff;
  border-radius: 10rpx; border: none; font-size: 28rpx; margin-top: 10rpx;
}

.divider { text-align: center; margin: 20rpx; position: relative; }
.divider text { font-size: 24rpx; color: #999; background: #f5f5f5; padding: 0 20rpx; position: relative; z-index: 1; }

.form-card { background: #fff; border-radius: 16rpx; padding: 30rpx; margin: 0 20rpx 40rpx; }
.form-group { margin-bottom: 30rpx; }
.label { font-size: 28rpx; color: #333; margin-bottom: 12rpx; display: block; }
.required { color: #f44336; }
.input { height: 80rpx; border: 1rpx solid #e0e0e0; border-radius: 10rpx; padding: 0 20rpx; font-size: 28rpx; }
.picker-row { height: 80rpx; border: 1rpx solid #e0e0e0; border-radius: 10rpx; padding: 0 20rpx; display: flex; align-items: center; justify-content: space-between; font-size: 28rpx; color: #333; }
.arrow { color: #999; }
.submit-btn { width: 100%; height: 90rpx; background: #2979FF; color: #fff; border-radius: 50rpx; border: none; font-size: 32rpx; }

/* ---- 已注册视图 ---- */
.user-card { background: linear-gradient(135deg, #2979FF, #1565C0); padding: 60rpx 30rpx 40rpx; display: flex; flex-direction: column; align-items: center; }
.avatar { width: 120rpx; height: 120rpx; border-radius: 50%; border: 4rpx solid rgba(255,255,255,0.3); margin-bottom: 20rpx; }
.name { font-size: 36rpx; color: #fff; font-weight: bold; }
.biz { font-size: 24rpx; color: rgba(255,255,255,0.8); margin-top: 8rpx; }
.info-section { background: #fff; margin: 20rpx; border-radius: 16rpx; padding: 30rpx; }
.section-title { font-size: 28rpx; color: #999; margin-bottom: 20rpx; }
.info-item { display: flex; justify-content: space-between; padding: 16rpx 0; }
.value { font-size: 28rpx; color: #666; }
.status-active { color: #4caf50; }
.menu-section { background: #fff; margin: 20rpx; border-radius: 16rpx; overflow: hidden; }
.menu-item { display: flex; justify-content: space-between; padding: 30rpx; font-size: 28rpx; color: #333; border-bottom: 1rpx solid #f5f5f5; }
.logout-text { color: #f44336; }
</style>
