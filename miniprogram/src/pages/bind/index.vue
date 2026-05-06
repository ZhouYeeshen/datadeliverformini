<template>
  <view class="bind-page">
    <view class="header">
      <text class="title">绑定企业信息</text>
      <text class="desc">绑定后即可为企业申报经营数据</text>
    </view>

    <!-- 手机号快速关联已注册企业 -->
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
    <view class="form">
      <view class="form-group">
        <text class="label">企业名称 <text class="required">*</text></text>
        <input class="input" v-model="form.businessName" placeholder="请输入企业名称" />
      </view>

      <view class="form-group">
        <text class="label">行业类型 <text class="required">*</text></text>
        <picker :range="industryOptions" @change="onIndustryChange">
          <view class="picker">
            {{ form.industryType || '请选择行业类型' }}
            <text class="arrow">></text>
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
    </view>

    <button class="submit-btn" @click="handleBind" :loading="loading">确认绑定</button>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAppStore } from '../../store'
import { bindBusiness } from '../../utils/api'

const store = useAppStore()
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
  store.realName = data.real_name || ''

  uni.showToast({ title: '绑定成功', icon: 'success' })
  setTimeout(() => uni.switchTab({ url: '/pages/home/index' }), 800)
}

// 快速关联：仅通过手机号 + 姓名，匹配已有企业
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
    await doBind({
      phone: quickPhone.value,
      real_name: quickName.value
    })
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
    uni.showToast({ title: err.errMsg || '绑定失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.bind-page { min-height: 100vh; background: #f5f5f5; padding: 30rpx; }
.header { text-align: center; padding: 40rpx 0; }
.title { font-size: 36rpx; font-weight: bold; color: #333; }
.desc { font-size: 24rpx; color: #999; display: block; margin-top: 10rpx; }

.quick-bind-card { background: #e3f2fd; border-radius: 16rpx; padding: 30rpx; margin-bottom: 30rpx; border: 2rpx solid #90caf9; }
.card-title { font-size: 26rpx; color: #1565C0; margin-bottom: 20rpx; font-weight: 500; text-align: center; }
.quick-btn {
  width: 100%; height: 80rpx; background: #1565C0; color: #fff;
  border-radius: 10rpx; border: none; font-size: 28rpx; margin-top: 10rpx;
}

.divider { text-align: center; margin: 20rpx 0; position: relative; }
.divider text { font-size: 24rpx; color: #999; background: #f5f5f5; padding: 0 20rpx; position: relative; z-index: 1; }

.form { background: #fff; border-radius: 16rpx; padding: 30rpx; margin-bottom: 40rpx; }
.form-group { margin-bottom: 30rpx; }
.label { font-size: 28rpx; color: #333; margin-bottom: 12rpx; display: block; }
.required { color: #f44336; }
.input { height: 80rpx; border: 1rpx solid #e0e0e0; border-radius: 10rpx; padding: 0 20rpx; font-size: 28rpx; }
.picker { height: 80rpx; border: 1rpx solid #e0e0e0; border-radius: 10rpx; padding: 0 20rpx; display: flex; align-items: center; justify-content: space-between; font-size: 28rpx; color: #333; }
.arrow { color: #999; }
.submit-btn { width: 100%; height: 90rpx; background: #2979FF; color: #fff; border-radius: 50rpx; border: none; font-size: 32rpx; }
</style>
