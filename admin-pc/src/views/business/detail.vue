<template>
  <div class="detail-page">
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>{{ biz?.name }} — 台账</span>
          <el-tag :type="biz?.status === 'active' ? 'success' : 'danger'">{{ biz?.status === 'active' ? '正常' : '禁用' }}</el-tag>
        </div>
      </template>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="行业类型">{{ biz?.industry_type }}</el-descriptions-item>
        <el-descriptions-item label="法人">{{ biz?.legal_person }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ biz?.contact_phone }}</el-descriptions-item>
        <el-descriptions-item label="地址" :span="2">{{ biz?.address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ biz?.created_at?.slice(0, 10) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card title="申报台账" style="margin-top: 20px">
      <el-table :data="reports" stripe>
        <el-table-column prop="report_month" label="申报月份" width="120" />
        <el-table-column prop="restaurant_revenue" label="餐饮" width="110">
          <template #default="{ row }">¥{{ row.restaurant_revenue?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="retail_revenue" label="零售" width="110">
          <template #default="{ row }">¥{{ row.retail_revenue?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="accommodation_revenue" label="住宿" width="110">
          <template #default="{ row }">¥{{ row.accommodation_revenue?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="tobacco_alcohol_revenue" label="烟酒" width="110">
          <template #default="{ row }">¥{{ row.tobacco_alcohol_revenue?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="other_revenue" label="其他" width="110">
          <template #default="{ row }">¥{{ row.other_revenue?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="total_revenue" label="合计" width="120">
          <template #default="{ row }"><b>¥{{ row.total_revenue?.toFixed(2) }}</b></template>
        </el-table-column>
        <el-table-column prop="source_type" label="来源" width="80" />
        <el-table-column label="原始照片" width="100">
          <template #default="{ row }">
            <el-button v-if="row.photo_urls" type="primary" link @click="viewPhotos(row)">查看</el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="photoVisible" title="原始扫描照片" width="700px">
      <el-image v-for="url in currentPhotos" :key="url" :src="url" style="width: 100%; margin-bottom: 12px" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getBusinessDetail, getBusinessLedger } from '@/api'

const route = useRoute()
const biz = ref<any>(null)
const reports = ref<any[]>([])
const photoVisible = ref(false)
const currentPhotos = ref<string[]>([])

function getPhotoSrc(url: string): string {
  if (url.startsWith('http')) return url
  return url.startsWith('/') ? url : '/' + url
}

function viewPhotos(row: any) {
  currentPhotos.value = (row.photo_urls || '').split(',').filter(Boolean).map(getPhotoSrc)
  photoVisible.value = true
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    const [bizRes, ledgerRes] = await Promise.all([
      getBusinessDetail(id), getBusinessLedger(id)
    ])
    biz.value = bizRes.data.business || bizRes.data
    reports.value = ledgerRes.data.list || []
  } catch {}
})
</script>

<style scoped>
.detail-page { }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
