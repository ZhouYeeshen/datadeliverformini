<template>
  <div class="stats-page">
    <el-card>
      <div class="header-bar">
        <el-date-picker v-model="selectedMonth" type="month" placeholder="选择月份" value-format="YYYY-MM" @change="loadData" />
        <el-button @click="exportData">导出报表</el-button>
      </div>
    </el-card>

    <el-card title="月度汇总" style="margin-top: 20px">
      <el-descriptions :column="4" border v-if="summary">
        <el-descriptions-item label="月份">{{ summary.month }}</el-descriptions-item>
        <el-descriptions-item label="上报户数">{{ summary.report_count }}</el-descriptions-item>
        <el-descriptions-item label="营业总额">¥{{ summary.total_revenue?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="户均营收">¥{{ (summary.total_revenue / (summary.report_count || 1)).toFixed(2) }}</el-descriptions-item>
      </el-descriptions>

      <v-chart v-if="industryData" :option="industryChartOption" style="height: 350px; margin-top: 20px" autoresize />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { getMonthlySummary } from '@/api'

use([CanvasRenderer, BarChart, GridComponent, TooltipComponent])

const selectedMonth = ref('')
const summary = ref<any>(null)
const industryData = ref<any>({})

const industryChartOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: Object.keys(industryData.value) },
  yAxis: { type: 'value', name: '万元' },
  series: [{
    type: 'bar', name: '营收',
    data: Object.values(industryData.value).map((v: any) => +(v / 10000).toFixed(2)),
    itemStyle: { color: '#409EFF', borderRadius: [6, 6, 0, 0] }
  }]
}))

async function loadData() {
  try {
    const res = await getMonthlySummary(selectedMonth.value)
    summary.value = res.data
    industryData.value = res.data.industry_totals || {}
  } catch {}
}

function exportData() {
  if (!summary.value) return
  const csv = [
    ['月份', '上报户数', '营业总额'],
    [summary.value.month, summary.value.report_count, summary.value.total_revenue]
  ].map(row => row.join(',')).join('\n')
  const blob = new Blob(['﻿' + csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url; a.download = `月报_${summary.value.month}.csv`; a.click()
}

onMounted(() => loadData())
</script>

<style scoped>
.stats-page { display: flex; flex-direction: column; gap: 20px; }
.header-bar { display: flex; gap: 16px; align-items: center; }
</style>
