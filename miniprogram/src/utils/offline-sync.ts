const QUEUE_KEY = 'offline_report_queue'
const OCR_RESULTS_KEY = 'offline_ocr_results'

interface OfflineReport {
  id: string
  businessId: number
  data: any
  photoUrls: string[]
  createdAt: number
  synced: boolean
}

// Save OCR result for later upload
export function saveOCRResult(result: any): void {
  const list = getOCRResults()
  list.push({ ...result, savedAt: Date.now() })
  uni.setStorageSync(OCR_RESULTS_KEY, JSON.stringify(list))
}

export function getOCRResults(): any[] {
  try {
    return JSON.parse(uni.getStorageSync(OCR_RESULTS_KEY) || '[]')
  } catch { return [] }
}

export function clearOCRResults(): void {
  uni.removeStorageSync(OCR_RESULTS_KEY)
}

// Offline report queue
export function enqueueReport(report: OfflineReport): void {
  const queue = getQueue()
  queue.push(report)
  uni.setStorageSync(QUEUE_KEY, JSON.stringify(queue))
}

export function getQueue(): OfflineReport[] {
  try {
    return JSON.parse(uni.getStorageSync(QUEUE_KEY) || '[]')
  } catch { return [] }
}

export function updateQueueItem(id: string, synced: boolean): void {
  const queue = getQueue()
  const item = queue.find(r => r.id === id)
  if (item) item.synced = synced
  uni.setStorageSync(QUEUE_KEY, JSON.stringify(queue))
}

export function getPendingQueue(): OfflineReport[] {
  return getQueue().filter(r => !r.synced)
}

// Sync all pending reports
export async function syncPending(syncFn: (report: OfflineReport) => Promise<any>): Promise<number> {
  const pending = getPendingQueue()
  let synced = 0

  for (const report of pending) {
    try {
      await syncFn(report)
      updateQueueItem(report.id, true)
      synced++
    } catch (err) {
      console.error('Sync failed for report', report.id, err)
      break // stop on first failure to preserve order
    }
  }

  return synced
}

// Network status monitoring
export function watchNetwork(onOnline: () => void): void {
  uni.onNetworkStatusChange((res) => {
    if (res.isConnected) {
      onOnline()
    }
  })
}
