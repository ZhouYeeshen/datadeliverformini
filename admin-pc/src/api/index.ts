import axios from 'axios'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

http.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('admin_token')
      window.location.hash = '#/login'
    }
    return Promise.reject(err)
  }
)

// Auth
export function adminLogin(username: string, password: string) {
  return http.post('/auth/admin/login', { username, password })
}

// Business
export function getBusinesses(params: any) {
  return http.get('/admin/businesses', { params })
}
export function getBusinessDetail(id: number) {
  return http.get(`/admin/businesses/${id}`)
}
export function getBusinessLedger(id: number) {
  return http.get(`/admin/businesses/${id}/ledger`)
}
export function createBusiness(data: any) {
  return http.post('/admin/businesses', data)
}
export function updateBusiness(id: number, data: any) {
  return http.put(`/admin/businesses/${id}`, data)
}

// Reports
export function getAdminReports(params: any) {
  return http.get('/admin/reports', { params })
}
export function getAdminReportDetail(id: number) {
  return http.get(`/admin/reports/${id}`)
}
export function updateReport(id: number, data: any) {
  return http.put(`/admin/reports/${id}`, data)
}

// Alerts
export function getAlerts(params: any) {
  return http.get('/admin/alerts', { params })
}
export function getUnreadAlertCount() {
  return http.get('/admin/alerts/unread-count')
}
export function markAlertRead(id: number) {
  return http.put(`/admin/alerts/${id}/read`)
}

// Statistics
export function getOverview() {
  return http.get('/admin/statistics/overview')
}
export function getTrend(months = 12) {
  return http.get('/admin/statistics/trend', { params: { months } })
}
export function getIndustryStats(month?: string) {
  return http.get('/admin/statistics/industry', { params: { month } })
}
export function getMonthlySummary(month?: string) {
  return http.get('/admin/statistics/monthly', { params: { month } })
}

// POS
export function getPosCredentials(params: any) {
  return http.get('/admin/pos/credentials', { params })
}
export function createPosCredential(businessId: number) {
  return http.post('/admin/pos/credentials', { business_id: businessId })
}
export function updatePosCredential(id: number, status: string) {
  return http.put(`/admin/pos/credentials/${id}`, { status })
}
export function getPosLogs(params: any) {
  return http.get('/admin/pos/logs', { params })
}
