// OCR integration — device-side VKSession or backend fallback

import { ocrRecognize, request } from './api'
import { isVKSupported, startVKOCR, stopVKOCR } from './wechat-ocr'

interface OCRResult {
  restaurant_revenue: number
  retail_revenue: number
  accommodation_revenue: number
  tobacco_alcohol_revenue: number
  other_revenue: number
  total_revenue: number
  raw_text: string
  engine: string
  confidence: string
  photo_url: string
}

export function parseOCRResult(rawData: any): OCRResult {
  return {
    restaurant_revenue: rawData?.restaurant_revenue || 0,
    retail_revenue: rawData?.retail_revenue || 0,
    accommodation_revenue: rawData?.accommodation_revenue || 0,
    tobacco_alcohol_revenue: rawData?.tobacco_alcohol_revenue || 0,
    other_revenue: rawData?.other_revenue || 0,
    total_revenue: rawData?.total_revenue || 0,
    raw_text: rawData?.raw_text || '',
    engine: rawData?.engine || '',
    confidence: rawData?.confidence || '',
    photo_url: rawData?.photo_url || ''
  }
}

// Send raw text to backend for category parsing
export async function parseText(rawText: string): Promise<OCRResult> {
  const res = await request('/wx/ocr/parse', {
    method: 'POST',
    data: { raw_text: rawText },
    timeout: 15000
  })
  return parseOCRResult(res)
}

// Backend-based OCR (WeChat OCR API — requires backend API permission)
export async function scanReceipt(photoPath: string): Promise<OCRResult> {
  const result = await ocrRecognize(photoPath)
  return parseOCRResult(result)
}

// Device-side OCR via VKSession
export function initDeviceOCR(
  cameraComp: any,
  onText: (text: string) => void
): { supported: boolean; stop: () => void } {
  if (!isVKSupported()) {
    return { supported: false, stop: () => {} }
  }
  const ok = startVKOCR(cameraComp, onText)
  return { supported: ok, stop: stopVKOCR }
}

export { isVKSupported, startVKOCR, stopVKOCR }
