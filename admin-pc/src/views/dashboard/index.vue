<template>
  <div class="dashboard">
    <div class="stat-cards">
      <el-card class="stat-card">
        <div class="stat-num">{{ overview.total_businesses }}</div>
        <div class="stat-label">注册商户总数</div>
      </el-card>
      <el-card class="stat-card">
        <div class="stat-num">{{ overview.current_month_reports }}</div>
        <div class="stat-label">本月已报商户</div>
      </el-card>
      <el-card class="stat-card accent">
        <div class="stat-num">¥{{ fmtMoney(overview.current_month_total) }}</div>
        <div class="stat-label">本月营业总额</div>
      </el-card>
      <el-card class="stat-card">
        <div class="stat-num">{{ unreadCount }}</div>
        <div class="stat-label">未读预警</div>
      </el-card>
    </div>

    <div class="charts-row">
      <el-card class="chart-card" title="月度趋势">
        <v-chart :option="trendOption" style="height: 350px" autoresize />
      </el-card>
      <el-card class="chart-card" title="行业分布">
        <v-chart :option="industryOption" style="height: 350px" autoresize />
      </el-card>
    </div>

    <el-card title="行业营收对比">
      <v-chart :option="industryRevenueOption" style="height: 300px" autoresize />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent } from 'echarts/components'
import { getOverview, getTrend, getIndustryStats, getUnreadAlertCount } from '@/api'

use([CanvasRenderer, BarChart, LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent])

const unreadCount = ref(0)
const overview = reactive({
  total_businesses: 0,
  current_month_reports: 0,
  current_month_total: 0,
  by_industry: {},
  industry_totals: {}
})

const trendData = ref<any[]>([])
const industryData = ref<any>({})

const trendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['营业总额(万)', '上报商户数'] },
  xAxis: { type: 'category', data: trendData.value.map((d: any) => d.month) },
  yAxis: [{ type: 'value', name: '万元' }, { type: 'value', name: '户' }],
  series: [
    {
      name: '营业总额(万)', type: 'line', smooth: true,
      data: trendData.value.map((d: any) => +(d.total / 10000).toFixed(2)),
      itemStyle: { color: '#409EFF' }
    },
    {
      name: '上报商户数', type: 'bar', yAxisIndex: 1,
      data: trendData.value.map((d: any) => d.count),
      itemStyle: { color: '#67C23A' }
    }
  ]
}))

const industryOption = computed(() => {
  const data = overview.by_industry || {}
  return {
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie', radius: ['45%', '75%'],
      data: Object.entries(data).map(([k, v]) => ({ name: k, value: v })),
      label: { formatter: '{b}\n{d}%' }
    }]
  }
})

const industryRevenueOption = computed(() => {
  const data = overview.industry_totals || {}
  return {
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: Object.keys(data) },
    yAxis: { type: 'value', name: '万元' },
    series: [{
      type: 'bar',
      data: Object.values(data).map((v: any) => +(v / 10000).toFixed(2)),
      itemStyle: { color: '#409EFF', borderRadius: [6, 6, 0, 0] }
    }]
  }
})

function fmtMoney(v: number) {
  return v >= 10000 ? (v / 10000).toFixed(2) + '万' : v.toFixed(2)
}

async function loadData() {
  try {
    const [ovRes, trRes, countRes] = await Promise.all([
      getOverview(), getTrend(), getUnreadAlertCount()
    ])
    Object.assign(overview, ovRes.data)
    trendData.value = trRes.data
    unreadCount.value = countRes.data.unread_count
  } catch {}
}

onMounted(() => loadData())
</script>

<style scoped>
.dashboard { display: flex; flex-direction: column; gap: 20px; }

/* Stat cards with gradient backgrounds */
.stat-cards { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }
.stat-card {
  text-align: center;
  padding: 10px;
  position: relative;
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 32px rgba(0, 212, 255, 0.15);
}
.stat-card::before {
  content: '';
  position: absolute;
  top: -50%; right: -50%;
  width: 100%; height: 100%;
  background: radial-gradient(circle, rgba(0,212,255,0.06) 0%, transparent 70%);
  pointer-events: none;
}

/* Gradient text for numbers */
.stat-card .stat-num {
  font-size: 36px;
  font-weight: 800;
  background: linear-gradient(135deg, #00d4ff, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 1px;
  position: relative;
  z-index: 1;
}
.stat-card.accent .stat-num {
  background: linear-gradient(135deg, #00d4ff, #60a5fa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.stat-card .stat-label {
  font-size: 13px;
  color: #64748b;
  margin-top: 8px;
  letter-spacing: 1px;
  text-transform: uppercase;
  position: relative;
  z-index: 1;
}

.charts-row { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.chart-card { }
</style>
