const fs = require('fs')
const path = require('path')

const appJsonPath = path.join(__dirname, '..', 'dist', 'build', 'mp-weixin', 'app.json')

if (!fs.existsSync(appJsonPath)) {
  console.log('[fix-app-json] app.json not found, skipping')
  process.exit(0)
}

const appJson = JSON.parse(fs.readFileSync(appJsonPath, 'utf-8'))

// Valid values for requiredPrivateInfos in WeChat mini program
const validPrivateInfos = [
  'chooseAddress',
  'chooseLocation',
  'choosePoi',
  'getFuzzyLocation',
  'getLocation',
  'onLocationChange',
  'startLocationUpdate',
  'startLocationUpdateBackground'
]

if (appJson.requiredPrivateInfos) {
  const before = appJson.requiredPrivateInfos.length
  appJson.requiredPrivateInfos = appJson.requiredPrivateInfos.filter(
    (item) => validPrivateInfos.includes(item)
  )
  if (appJson.requiredPrivateInfos.length < before) {
    fs.writeFileSync(appJsonPath, JSON.stringify(appJson, null, 2), 'utf-8')
    console.log(`[fix-app-json] Removed ${before - appJson.requiredPrivateInfos.length} invalid private info entries`)
  } else {
    console.log('[fix-app-json] No invalid entries found')
  }
}
