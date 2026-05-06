// ============================================================
// 服务器地址配置 —— 切换网络环境时只需修改这里
// ============================================================
// 局域网开发:  http://192.168.1.5:8080
// ngrok 穿透:  https://xxx.ngrok-free.app
// 云服务器:    https://your-domain.com
// ============================================================
const HOST = 'http://127.0.0.1:8080'

export const BASE_URL = HOST + '/api/v1'
export const SERVER_HOST = HOST

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  auth?: boolean
  timeout?: number
}

function getToken(): string {
  return uni.getStorageSync('token') || ''
}

export async function request(path: string, options: RequestOptions = {}) {
  const { method = 'GET', data, auth = true, timeout = 15000 } = options

  const header: any = { 'Content-Type': 'application/json', 'ngrok-skip-browser-warning': '1' }
  if (auth) {
    const token = getToken()
    if (token) header['Authorization'] = `Bearer ${token}`
  }

  try {
    const res = await uni.request({
      url: BASE_URL + path,
      method,
      header,
      data,
      timeout
    })
    return (res.data as any)
  } catch (err) {
    throw err
  }
}

// Auth API
export function wechatLogin(code: string) {
  return request('/auth/wechat/login', { method: 'POST', data: { code }, auth: false })
}

export function bindBusiness(data: any) {
  return request('/auth/wechat/bind', { method: 'POST', data, auth: false })
}

// Report API
export function submitReport(data: any) {
  return request('/wx/reports', { method: 'POST', data })
}

export function getCurrentMonthReport() {
  return request('/wx/reports/current')
}

export function getReportHistory() {
  return request('/wx/reports/history')
}

export function getReportDetail(id: number) {
  return request(`/wx/reports/${id}`)
}

export function getProfile() {
  return request('/wx/profile')
}

// Upload
export function uploadPhoto(filePath: string) {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: BASE_URL + '/wx/upload/photo',
      filePath,
      name: 'photo',
      header: { 'Authorization': `Bearer ${getToken()}`, 'ngrok-skip-browser-warning': '1' },
      success: (res) => resolve(JSON.parse(res.data)),
      fail: reject
    })
  })
}

// OCR
export async function ocrRecognize(filePath: string): Promise<any> {
  // Compress image first to reduce upload size
  let path = filePath
  try {
    const res = await uni.compressImage({ src: filePath, quality: 70, compressedWidth: 1200 })
    if (res.tempFilePath) path = res.tempFilePath
  } catch {
    // Compression may fail in dev tools; use original
  }

  const fs = uni.getFileSystemManager()
  const data = fs.readFileSync(path, 'base64') as string
  return request('/wx/ocr/recognize', {
    method: 'POST',
    data: { image_base64: data },
    timeout: 60000
  })
}
