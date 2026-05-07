<template>
  <view class="scan-page">
    <!-- Camera view -->
    <view class="camera-area" v-if="!photoPath">
      <camera
        id="scanCamera"
        device-position="back"
        flash="off"
        class="camera"
        @error="onCameraError"
        @ready="onCameraReady"
      >
        <cover-view class="camera-overlay">
          <cover-view class="scan-frame"></cover-view>
          <cover-view class="tips">将单据放入框内拍摄</cover-view>
          <cover-view class="live-text" v-if="liveText">
            <cover-view class="live-text-label">已识别文字（手写+印刷）：</cover-view>
            <cover-view class="live-text-content">{{ liveText }}</cover-view>
          </cover-view>
        </cover-view>
      </camera>
      <view class="capture-bar">
        <button class="capture-btn" @click="takePhoto">拍照</button>
        <button class="album-btn" @click="chooseFromAlbum">从相册选择</button>
      </view>
    </view>

    <!-- Photo preview + actions -->
    <view class="preview-area" v-else-if="!ocrDone">
      <image :src="photoPath" mode="widthFix" class="preview-img" />
      <view class="preview-actions">
        <button class="retake-btn" @click="retake">重拍</button>
        <button class="ocr-btn" @click="doOCR" :loading="recognizing">
          {{ liveText ? '使用已识别手写文字' : '识别单据' }}
        </button>
      </view>

      <view class="ocr-progress" v-if="recognizing">
        <text>正在识别中...</text>
        <progress :percent="ocrProgress" stroke-width="6" activeColor="#2979FF" />
      </view>
    </view>

    <!-- OCR result area -->
    <view class="result-area" v-else>
      <image :src="photoPath" mode="widthFix" class="result-img" />

      <view class="text-card">
        <view class="text-card-header">
          <text class="text-card-title">识别文字内容</text>
          <text class="engine-tag" v-if="engine">{{ engine }}</text>
        </view>
        <view class="raw-text-box">
          <text class="raw-text">{{ rawText || '未识别到文字内容' }}</text>
        </view>
      </view>

      <!-- Editable amounts -->
      <view class="amounts-card" v-if="items.length > 0">
        <view class="amounts-header">
          <text class="amounts-title">金额明细</text>
          <text class="add-row-btn" @click="addRow">+ 添加</text>
        </view>

        <scroll-view scroll-y class="amounts-scroll" :style="{ maxHeight: scrollHeight }">
          <view class="edit-row" v-for="(item, i) in items" :key="i">
            <input
              class="edit-label"
              v-model="item.label"
              placeholder="项目名称"
              @input="onItemChange"
            />
            <view class="edit-value-wrap">
              <text class="yen-sign">¥</text>
              <input
                class="edit-value"
                v-model.number="item.value"
                type="digit"
                placeholder="0.00"
                @input="onItemChange"
              />
            </view>
            <view class="del-btn" @click="delRow(i)" v-if="items.length > 1">×</view>
          </view>
        </scroll-view>

        <view class="amount-total">
          <text class="total-label">合计</text>
          <text class="total-value">¥{{ editableTotal.toFixed(2) }}</text>
        </view>
      </view>

      <view class="amounts-card" v-else>
        <view class="amounts-header">
          <text class="amounts-title">金额明细</text>
          <text class="add-row-btn" @click="addRow">+ 添加</text>
        </view>
        <view class="edit-row">
          <input class="edit-label" v-model="items[0].label" placeholder="项目名称" />
          <view class="edit-value-wrap">
            <text class="yen-sign">¥</text>
            <input class="edit-value" v-model.number="items[0].value" type="digit" placeholder="0.00" />
          </view>
        </view>
        <view class="amount-total">
          <text class="total-label">合计</text>
          <text class="total-value">¥{{ editableTotal.toFixed(2) }}</text>
        </view>
      </view>

      <view class="result-actions">
        <button class="confirm-btn" @click="goReport">确认，去上报</button>
        <button class="retake-btn" @click="retake">重新拍照</button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onUnload } from '@dcloudio/uni-app'
import { SERVER_HOST } from '../../utils/api'
import { scanReceipt, parseText, initDeviceOCR } from '../../utils/ocr'
import { saveOCRResult } from '../../utils/offline-sync'

const photoPath = ref('')
const recognizing = ref(false)
const ocrProgress = ref(0)
const ocrDone = ref(false)
const ocrTimer = ref<ReturnType<typeof setInterval> | null>(null)
const liveText = ref('')

const rawText = ref('')
const engine = ref('')
const photoUrl = ref('')

const scrollHeight = computed(() => {
  const sys = uni.getSystemInfoSync()
  return Math.min(Math.max(items.value.length, 3) * 100, sys.windowHeight * 0.35) + 'px'
})

interface AmountItem {
  label: string
  value: number
}

const items = ref<AmountItem[]>([{ label: '', value: 0 }])

const cateLabels: Record<string, string> = {
  restaurant_revenue: '餐饮',
  retail_revenue: '零售',
  accommodation_revenue: '住宿',
  tobacco_alcohol_revenue: '烟酒',
  other_revenue: '其他'
}

const editableTotal = computed(() => {
  const sum = items.value.reduce((s, item) => s + (item.value || 0), 0)
  return Math.round(sum * 100) / 100
})

let vkStop: (() => void) | null = null

function onCameraReady() {
  // @ts-ignore
  const comp = uni.createCameraContext ? uni.createCameraContext('scanCamera') : null
  const { supported, stop } = initDeviceOCR(comp, (text: string) => {
    liveText.value = text.substring(0, 500)
  })
  if (!supported) {
    console.log('[OCR] VKSession not available on this platform, will use backend API fallback')
  }
  vkStop = stop
}

function takePhoto() {
  const ctx = uni.createCameraContext()
  ctx.takePhoto({
    quality: 'high',
    success: (res: any) => {
      photoPath.value = res.tempImagePath
    },
    fail: () => {
      uni.showToast({ title: '拍照失败', icon: 'none' })
    }
  })
}

function chooseFromAlbum() {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album'],
    success: (res) => {
      photoPath.value = res.tempFilePaths[0]
    }
  })
}

function retake() {
  photoPath.value = ''
  ocrDone.value = false
  rawText.value = ''
  liveText.value = ''
  items.value = [{ label: '', value: 0 }]
}

function onCameraError(e: any) {
  uni.showToast({ title: '摄像头不可用', icon: 'none' })
}

function addRow() {
  items.value.push({ label: '', value: 0 })
}

function delRow(i: number) {
  items.value.splice(i, 1)
}

function onItemChange() {
  // Reactive update happens automatically via v-model
}

async function doOCR() {
  if (!photoPath.value) {
    uni.showToast({ title: '请先拍照或选择照片', icon: 'none' })
    return
  }

  recognizing.value = true
  ocrProgress.value = 0

  try {
    ocrTimer.value = setInterval(() => {
      if (ocrProgress.value < 90) ocrProgress.value += 10
    }, 200)

    let result

    // If VKSession detected handwritten/printed text, parse it via backend
    if (liveText.value) {
      ocrProgress.value = 40
      try {
        result = await parseText(liveText.value)
        result.engine = 'device-handwritten'
        result.photo_url = photoPath.value
      } catch {
        // parseText failed, fall through to backend OCR
        result = null
      }
    }

    // Fallback: backend OCR API (printed text only)
    if (!result) {
      try {
        result = await scanReceipt(photoPath.value)
      } catch {
        result = {
          restaurant_revenue: 0, retail_revenue: 0,
          accommodation_revenue: 0, tobacco_alcohol_revenue: 0,
          other_revenue: 0, total_revenue: 0,
          raw_text: '', engine: '', confidence: '', photo_url: ''
        }
        saveOCRResult(result)
        uni.showToast({ title: '网络不可用，已暂存本地', icon: 'none' })
      }
    }

    clearInterval(ocrTimer.value)
    ocrTimer.value = null
    ocrProgress.value = 100

    rawText.value = result.raw_text || ''
    engine.value = result.engine || ''
    let serverPhoto = result.photo_url || ''
    if (serverPhoto && serverPhoto.startsWith('/uploads/')) {
      serverPhoto = SERVER_HOST + serverPhoto
    }
    photoUrl.value = serverPhoto || photoPath.value

    const list: AmountItem[] = []
    const cats: Array<keyof typeof cateLabels> = ['restaurant_revenue', 'retail_revenue', 'accommodation_revenue', 'tobacco_alcohol_revenue', 'other_revenue']
    for (const c of cats) {
      if (result[c] > 0) {
        list.push({ label: cateLabels[c], value: result[c] })
      }
    }
    if (list.length === 0 && result.total_revenue > 0) {
      list.push({ label: '营业额', value: result.total_revenue })
    }
    if (list.length === 0) {
      list.push({ label: '', value: 0 })
    }

    items.value = list
    ocrDone.value = true
  } catch (err: any) {
    uni.showToast({ title: err.message || '识别失败', icon: 'none' })
  } finally {
    recognizing.value = false
  }
}

function goReport() {
  const total = editableTotal.value
  if (total <= 0) {
    uni.showToast({ title: '请输入至少一项金额', icon: 'none' })
    return
  }

  const data = {
    total_revenue: total,
    photo_url: photoUrl.value || photoPath.value,
    photo_local: !photoUrl.value ? photoPath.value : '',
    raw_text: rawText.value,
    items: items.value.filter(i => i.value > 0)
  }
  const params = encodeURIComponent(JSON.stringify(data))
  uni.navigateTo({ url: `/pages/report/index?data=${params}` })
}

onUnload(() => {
  if (ocrTimer.value) {
    clearInterval(ocrTimer.value)
    ocrTimer.value = null
  }
  if (vkStop) {
    vkStop()
    vkStop = null
  }
})
</script>

<style scoped>
.scan-page { min-height: 100vh; background: #000; position: relative; }

/* Camera */
.camera { width: 100%; height: 80vh; }
.camera-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.scan-frame { width: 500rpx; height: 500rpx; border: 4rpx solid rgba(41,121,255,0.6); border-radius: 20rpx; }
.tips { color: #fff; margin-top: 40rpx; font-size: 28rpx; opacity: 0.8; }
.capture-bar { position: absolute; bottom: 60rpx; left: 0; right: 0; display: flex; justify-content: center; gap: 40rpx; }
.capture-btn { width: 200rpx; height: 80rpx; background: #2979FF; color: #fff; border-radius: 40rpx; border: none; font-size: 28rpx; }
.album-btn { width: 200rpx; height: 80rpx; background: rgba(255,255,255,0.2); color: #fff; border-radius: 40rpx; border: 1rpx solid rgba(255,255,255,0.4); font-size: 28rpx; }

/* Live text from VKSession */
.live-text { position: absolute; bottom: 140rpx; left: 40rpx; right: 40rpx; background: rgba(0,0,0,0.75); border-radius: 12rpx; padding: 16rpx; max-height: 200rpx; overflow: hidden; }
.live-text-label { color: #4CAF50; font-size: 22rpx; display: block; }
.live-text-content { color: #fff; font-size: 24rpx; word-break: break-all; display: block; margin-top: 8rpx; line-height: 1.5; }

/* Preview */
.preview-area { padding: 30rpx; background: #000; min-height: 100vh; }
.preview-img { width: 100%; border-radius: 12rpx; }
.preview-actions { display: flex; gap: 20rpx; margin-top: 30rpx; padding: 0 20rpx; }
.retake-btn, .ocr-btn { flex: 1; height: 80rpx; border-radius: 40rpx; font-size: 28rpx; }
.ocr-btn { background: #2979FF; color: #fff; border: none; }
.retake-btn { background: #f5f5f5; color: #333; border: none; }
.ocr-progress { color: #fff; text-align: center; padding: 40rpx; }

/* Result area */
.result-area { min-height: 100vh; background: #f5f5f5; padding-bottom: 40rpx; }
.result-img { width: 100%; max-height: 360rpx; object-fit: cover; }

.text-card { background: #fff; margin: 20rpx 24rpx; border-radius: 16rpx; padding: 24rpx; }
.text-card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12rpx; }
.text-card-title { font-size: 28rpx; font-weight: bold; color: #333; }
.engine-tag { font-size: 20rpx; color: #999; background: #f0f0f0; padding: 4rpx 12rpx; border-radius: 8rpx; }
.raw-text-box { background: #f9f9f9; border-radius: 10rpx; padding: 20rpx; max-height: 240rpx; overflow-y: auto; }
.raw-text { font-size: 26rpx; color: #555; line-height: 1.6; white-space: pre-wrap; word-break: break-all; }

/* Editable amounts */
.amounts-card { background: #fff; margin: 0 24rpx 20rpx; border-radius: 16rpx; padding: 24rpx; }
.amounts-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16rpx; }
.amounts-title { font-size: 28rpx; font-weight: bold; color: #333; }
.add-row-btn { font-size: 26rpx; color: #2979FF; padding: 4rpx 12rpx; }

.amounts-scroll { padding-right: 8rpx; }

.edit-row { display: flex; align-items: center; padding: 14rpx 0; border-bottom: 1rpx solid #f0f0f0; gap: 12rpx; }
.edit-label { flex: 1; font-size: 26rpx; color: #333; height: 60rpx; background: #f9f9f9; border-radius: 8rpx; padding: 0 16rpx; }
.edit-value-wrap { display: flex; align-items: center; background: #f9f9f9; border-radius: 8rpx; padding: 0 12rpx; height: 60rpx; }
.yen-sign { font-size: 26rpx; color: #999; margin-right: 4rpx; }
.edit-value { width: 160rpx; font-size: 28rpx; color: #333; text-align: right; height: 60rpx; }
.del-btn { width: 48rpx; height: 48rpx; border-radius: 50%; background: #ff5252; color: #fff; font-size: 32rpx; line-height: 48rpx; text-align: center; flex-shrink: 0; }

.amount-total { display: flex; justify-content: space-between; align-items: center; margin-top: 20rpx; padding-top: 20rpx; border-top: 2rpx solid #2979FF; }
.total-label { font-size: 30rpx; font-weight: bold; color: #333; }
.total-value { font-size: 36rpx; font-weight: bold; color: #2979FF; }

.result-actions { padding: 20rpx 24rpx; display: flex; flex-direction: column; gap: 16rpx; }
.confirm-btn { width: 100%; height: 90rpx; background: #2979FF; color: #fff; border-radius: 50rpx; border: none; font-size: 32rpx; }
</style>
