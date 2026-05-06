// On-device OCR using WeChat VKSession + OCR plugin
// Requires plugin "ocr-plugin" (wx4418e3e031e551be) declared in pages.json

interface OCRAnchor {
  text: string
  confidence: number
  points: Array<{ x: number; y: number }>
}

let vkSession: any = null
let onAnchorUpdate: ((text: string) => void) | null = null

export function isVKSupported(): boolean {
  // @ts-ignore
  return typeof wx !== 'undefined' && typeof wx.createVKSession === 'function'
}

export function startVKOCR(cameraComp: any, onText: (text: string) => void): boolean {
  if (!isVKSupported()) {
    console.warn('[OCR] VKSession not supported')
    return false
  }

  try {
    // @ts-ignore
    vkSession = wx.createVKSession({
      track: {
        OCR: { mode: 1 } // 1 = generic text recognition
      }
    })

    vkSession.on('updateAnchors', (data: { anchors: OCRAnchor[] }) => {
      if (!data.anchors || data.anchors.length === 0) return
      const lines = data.anchors
        .sort((a, b) => a.points[0].y - b.points[0].y)
        .map(a => a.text)
        .filter(Boolean)
      if (lines.length > 0) {
        onText(lines.join('\n'))
      }
    })

    vkSession.start()
    console.log('[OCR] VKSession started')
    return true
  } catch (e) {
    console.error('[OCR] VKSession start failed:', e)
    return false
  }
}

export function stopVKOCR() {
  if (vkSession) {
    try {
      vkSession.stop()
      vkSession.destroy()
    } catch {
      // ignore
    }
    vkSession = null
  }
}
