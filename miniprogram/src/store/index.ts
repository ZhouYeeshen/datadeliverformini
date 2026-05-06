import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  const token = ref('')
  const isBound = ref(false)
  const businessId = ref(0)
  const businessName = ref('')
  const realName = ref('')
  const isOnline = ref(true)
  const isDev = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  function setToken(t: string) {
    token.value = t
    uni.setStorageSync('token', t)
  }

  function setBusiness(info: { business_id?: number, business_name?: string, real_name?: string, is_bound?: boolean }) {
    isBound.value = info.is_bound ?? true
    businessId.value = info.business_id ?? 0
    businessName.value = info.business_name ?? ''
    realName.value = info.real_name ?? ''
    uni.setStorageSync('business_id', businessId.value)
    uni.setStorageSync('business_name', businessName.value)
    uni.setStorageSync('real_name', realName.value)
  }

  function restoreSession() {
    // Detect dev mode (WeChat dev tools)
    try {
      const accountInfo = uni.getAccountInfoSync()
      isDev.value = accountInfo.miniProgram.envVersion === 'develop'
    } catch {
      isDev.value = false
    }

    const t = uni.getStorageSync('token')
    const bid = uni.getStorageSync('business_id')
    const bname = uni.getStorageSync('business_name')
    const rname = uni.getStorageSync('real_name')
    if (t) {
      token.value = t
      isBound.value = !!bid
      businessId.value = bid || 0
      businessName.value = bname || ''
      realName.value = rname || ''
    }
  }

  function logout() {
    token.value = ''
    isBound.value = false
    businessId.value = 0
    businessName.value = ''
    realName.value = ''
    uni.removeStorageSync('token')
    uni.removeStorageSync('business_id')
    uni.removeStorageSync('business_name')
    uni.removeStorageSync('real_name')
  }

  function setOnlineStatus(status: boolean) {
    isOnline.value = status
  }

  return {
    token, isBound, businessId, businessName, realName,
    isOnline, isLoggedIn, isDev,
    setToken, setBusiness, restoreSession, logout, setOnlineStatus
  }
})
