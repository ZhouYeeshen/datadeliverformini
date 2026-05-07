// OCR integration for handwritten receipt scanning
// Uploads photo to backend and receives revenue data from OCR processing

import { ocrRecognize } from './api'

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

export async function scanReceipt(photoPath: string): Promise<OCRResult> {
  const result = await ocrRecognize(photoPath)
  return parseOCRResult(result)
}
