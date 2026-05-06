<template>
  <view class="login-page">
    <view class="logo-area">
      <image class="logo" src="/static/logo.png" mode="aspectFit" />
      <text class="title">经营数据上报系统</text>
      <text class="subtitle">优质个体工商户服务平台</text>
    </view>

    <view class="btn-area">
      <button class="wechat-btn" @click="handleLogin" :loading="loading">
        <text class="wechat-icon">&#xe600;</text>
        微信授权登录
      </button>
      <text class="agreement">登录即同意《用户服务协议》</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAppStore } from '../../store'
import { wechatLogin } from '../../utils/api'

const store = useAppStore()
const loading = ref(false)

const DEV_TEST_CODE = 'DEV_TEST_USER_DEV_TEST_USER_DEV'

async function doLogin(code: string) {
  const res = await wechatLogin(code)
  store.setToken(res.token)

  if (res.is_bound) {
    store.setBusiness(res)
    uni.switchTab({ url: '/pages/home/index' })
  } else {
    uni.navigateTo({ url: '/pages/bind/index' })
  }
}

async function handleLogin() {
  loading.value = true
  try {
    let code = ''
    try {
      const loginRes = await uni.login({ timeout: 5000 })
      code = loginRes.code
    } catch {
      code = 'test_' + Date.now()
    }

    await doLogin(code)
  } catch (err: any) {
    uni.showToast({ title: err.errMsg || '登录失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  // Already authenticated — skip login
  if (store.token && store.isBound) {
    uni.switchTab({ url: '/pages/home/index' })
    return
  }

  if (store.isDev && !store.token) {
    loading.value = true
    try {
      await doLogin(DEV_TEST_CODE)
    } catch {
      // fall through to manual login
    } finally {
      loading.value = false
    }
  }
})
</script>

<style scoped>
.login-page {
  min-height: 100vh; display: flex; flex-direction: column;
  align-items: center; justify-content: center; padding: 60rpx;
  background: linear-gradient(135deg, #2979FF, #1565C0);
}
.logo-area { display: flex; flex-direction: column; align-items: center; margin-bottom: 120rpx; }
.logo { width: 160rpx; height: 160rpx; margin-bottom: 40rpx; }
.title { font-size: 40rpx; color: #fff; font-weight: bold; }
.subtitle { font-size: 26rpx; color: rgba(255,255,255,0.8); margin-top: 12rpx; }
.btn-area { width: 100%; text-align: center; }
.wechat-btn {
  width: 100%; height: 100rpx; background: #fff; border-radius: 50rpx;
  display: flex; align-items: center; justify-content: center;
  font-size: 32rpx; color: #333; font-weight: 500; border: none;
}
.agreement { font-size: 22rpx; color: rgba(255,255,255,0.6); margin-top: 20rpx; }
</style>
