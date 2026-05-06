<template>
  <div class="alert-page">
    <el-card>
      <div class="search-bar">
        <el-select v-model="alertType" placeholder="预警类型" clearable style="width: 200px" @change="loadData">
          <el-option label="未申报预警" value="MISSING_REPORT" />
          <el-option label="异常波动预警" value="ABNORMAL_CHANGE" />
        </el-select>
        <el-select v-model="isReadFilter" placeholder="读取状态" clearable style="width: 160px" @change="loadData">
          <el-option label="未读" value="false" />
          <el-option label="已读" value="true" />
        </el-select>
        <el-button type="primary" @click="loadData">刷新</el-button>
      </div>

      <el-table :data="list" stripe v-loading="loading" @row-click="handleRowClick">
        <el-table-column prop="alert_type" label="类型" width="130">
          <template #default="{ row }">
            <el-tag :type="row.alert_type === 'MISSING_REPORT' ? 'warning' : 'danger'" size="small">
              {{ row.alert_type === 'MISSING_REPORT' ? '未申报' : '异常波动' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="business_name" label="企业名称" min-width="160" />
        <el-table-column prop="alert_month" label="月份" width="100" />
        <el-table-column label="参考值(前2月均值)" width="140">
          <template #default="{ row }">¥{{ row.reference_value?.toFixed(2) || '-' }}</template>
        </el-table-column>
        <el-table-column label="当月值" width="120">
          <template #default="{ row }">¥{{ row.current_value?.toFixed(2) || '-' }}</template>
        </el-table-column>
        <el-table-column label="变动幅度" width="110">
          <template #default="{ row }">
            <span v-if="row.change_percent" :style="{ color: row.change_percent > 0 ? '#f56c6c' : '#67c23a' }">
              {{ row.change_percent > 0 ? '+' : '' }}{{ row.change_percent?.toFixed(1) }}%
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_read" label="状态" width="80">
          <template #default="{ row }">
            <span :style="{ color: row.is_read ? '#909399' : '#e6a23c', fontWeight: row.is_read ? 'normal' : 'bold' }">
              {{ row.is_read ? '已读' : '未读' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button v-if="!row.is_read" type="primary" link @click.stop="markRead(row.id)">标为已读</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page" :page-size="pageSize"
        :total="total" layout="total, prev, pager, next"
        @current-change="loadData" style="margin-top: 20px; justify-content: center"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onMounted } from 'vue'
import { getAlerts, markAlertRead } from '@/api'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const alertType = ref('')
const isReadFilter = ref('')

async function loadData() {
  loading.value = true
  try {
    const params: any = { page: page.value, page_size: pageSize.value }
    if (alertType.value) params.alert_type = alertType.value
    if (isReadFilter.value) params.is_read = isReadFilter.value

    const res = await getAlerts(params)
    list.value = res.data.list
    total.value = res.data.total
  } finally { loading.value = false }
}

function handleRowClick(row: any) {
  if (!row.is_read) markRead(row.id)
}

async function markRead(id: number) {
  try {
    await markAlertRead(id)
    loadData()
  } catch {}
}

onMounted(() => loadData())
</script>

<style scoped>
.alert-page { }
.search-bar { display: flex; gap: 12px; margin-bottom: 20px; }
</style>
