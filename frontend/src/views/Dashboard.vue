<template>
  <div class="dashboard">
    <div class="page-header">
      <div class="welcome-section">
        <h1 class="page-title">仪表盘</h1>
        <p class="page-desc">欢迎回来，{{ userStore.userInfo?.username }}！</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" :icon="Plus" @click="createContract">新建合同</el-button>
      </div>
    </div>
    
    <el-row :gutter="24" class="stats-row">
      <el-col :span="6" v-for="(stat, index) in statsCards" :key="index">
        <el-card class="stat-card" :style="{ '--accent-color': stat.color }">
          <div class="stat-content">
            <div class="stat-icon" :style="{ background: stat.gradient }">
              <el-icon :size="24"><component :is="stat.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </div>
          <div class="stat-trend" v-if="stat.trend">
            <span :class="stat.trend > 0 ? 'trend-up' : 'trend-down'">
              <el-icon><CaretTop v-if="stat.trend > 0" /><CaretBottom v-else /></el-icon>
              {{ Math.abs(stat.trend) }}%
            </span>
            <span class="trend-label">较上月</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="24" class="charts-row">
      <el-col :span="14">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><TrendCharts /></el-icon>
                合同统计
              </span>
              <el-radio-group v-model="chartType" size="small">
                <el-radio-button label="pie">饼图</el-radio-button>
                <el-radio-button label="bar">柱状图</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="chartRef" style="height: 320px"></div>
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card class="overview-card">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><DataAnalysis /></el-icon>
                本月概况
              </span>
            </div>
          </template>
          <div class="overview-grid">
            <div class="overview-item">
              <div class="overview-icon" style="background: linear-gradient(135deg, #6366F1, #8B5CF6)">
                <el-icon :size="22"><Document /></el-icon>
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ statistics.this_month_contracts || 0 }}</div>
                <div class="overview-label">新增合同</div>
              </div>
            </div>
            <div class="overview-item">
              <div class="overview-icon" style="background: linear-gradient(135deg, #10B981, #34D399)">
                <el-icon :size="22"><Money /></el-icon>
              </div>
              <div class="overview-content">
                <div class="overview-value">¥{{ formatAmount(statistics.this_month_amount) }}</div>
                <div class="overview-label">合同金额</div>
              </div>
            </div>
            <div class="overview-item">
              <div class="overview-icon" style="background: linear-gradient(135deg, #F59E0B, #FBBF24)">
                <el-icon :size="22"><CircleCheck /></el-icon>
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ statistics.active_contracts || 0 }}</div>
                <div class="overview-label">进行中</div>
              </div>
            </div>
            <div class="overview-item">
              <div class="overview-icon" style="background: linear-gradient(135deg, #EF4444, #F87171)">
                <el-icon :size="22"><Warning /></el-icon>
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ statistics.expiring_soon || 0 }}</div>
                <div class="overview-label">即将到期</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="24">
      <el-col :span="24">
        <el-card class="table-card">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><WarningFilled /></el-icon>
                即将到期合同
              </span>
              <el-button type="primary" link @click="$router.push('/reminders')">
                查看全部 <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
          </template>
          <el-table :data="expiringContracts" style="width: 100%">
            <el-table-column prop="contract_no" label="合同编号" width="160">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ row.contract_no }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="title" label="合同标题" />
            <el-table-column prop="customer_name" label="客户名称" width="150" />
            <el-table-column prop="end_date" label="到期日期" width="120">
              <template #default="{ row }">
                <span :class="{ 'text-danger': isExpiringSoon(row.end_date) }">
                  {{ formatDate(row.end_date) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="140">
              <template #default="{ row }">
                <span class="amount">¥{{ formatAmount(row.amount) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="viewContract(row)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import * as echarts from 'echarts'
import { getStatistics, getExpiringContracts } from '@/api/approval'
import { 
  Document, CircleCheck, Clock, Bell, TrendCharts, DataAnalysis,
  Money, Warning, WarningFilled, ArrowRight, Plus, CaretTop, CaretBottom
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const chartType = ref('bar')
const statistics = ref({})
const expiringContracts = ref([])
const chartRef = ref(null)

const statsCards = computed(() => [
  {
    icon: 'Document',
    label: '合同总数',
    value: statistics.value.total_contracts || 0,
    color: '#6366F1',
    gradient: 'linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%)',
    trend: 12
  },
  {
    icon: 'CircleCheck',
    label: '进行中合同',
    value: statistics.value.active_contracts || 0,
    color: '#10B981',
    gradient: 'linear-gradient(135deg, #10B981 0%, #34D399 100%)',
    trend: 8
  },
  {
    icon: 'Clock',
    label: '待审批',
    value: statistics.value.pending_contracts || 0,
    color: '#F59E0B',
    gradient: 'linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%)',
    trend: -3
  },
  {
    icon: 'Warning',
    label: '即将到期',
    value: statistics.value.expiring_soon || 0,
    color: '#EF4444',
    gradient: 'linear-gradient(135deg, #EF4444 0%, #F87171 100%)',
    trend: 5
  }
])

const formatAmount = (value) => {
  if (!value) return '0.00'
  return Number(value).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return dateStr
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const isExpiringSoon = (date) => {
  if (!date) return false
  const endDate = new Date(date)
  const now = new Date()
  const diffDays = Math.ceil((endDate - now) / (1000 * 60 * 60 * 24))
  return diffDays <= 30
}

const loadStatistics = async () => {
  try {
    const data = await getStatistics()
    statistics.value = data
    initChart()
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

const loadExpiringContracts = async () => {
  try {
    const data = await getExpiringContracts(30)
    expiringContracts.value = data.contracts || []
  } catch (error) {
    console.error('加载到期合同失败:', error)
  }
}

const initChart = () => {
  if (!chartRef.value) return
  
  const chart = echarts.init(chartRef.value)
  
  const option = chartType.value === 'pie' ? getPieOption() : getBarOption()
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

const getPieOption = () => ({
  tooltip: { trigger: 'item', formatter: '¥{c}' },
  legend: { orient: 'vertical', right: 10, top: 'center' },
  color: ['#6366F1', '#10B981', '#F59E0B', '#EF4444', '#94A3B8'],
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    center: ['40%', '50%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: {
      label: { show: true, fontSize: 16, fontWeight: 'bold' },
      itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.3)' }
    },
    data: [
      { value: statistics.value.total_amount || 0, name: '总金额' },
      { value: statistics.value.this_month_amount || 0, name: '本月金额' }
    ]
  }]
})

const getBarOption = () => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { 
    type: 'category', 
    data: ['1月', '2月', '3月', '4月', '5月', '6月'],
    axisLine: { lineStyle: { color: '#E2E8F0' } },
    axisLabel: { color: '#64748B' }
  },
  yAxis: { 
    type: 'value',
    axisLine: { show: false },
    axisLabel: { color: '#64748B' },
    splitLine: { lineStyle: { color: '#F1F5F9' } }
  },
  series: [{
    data: [12, 15, 8, 20, 16, statistics.value.this_month_contracts || 0],
    type: 'bar',
    barWidth: '50%',
    itemStyle: {
      borderRadius: [8, 8, 0, 0],
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: '#6366F1' },
        { offset: 1, color: '#8B5CF6' }
      ])
    }
  }]
})

const viewContract = (row) => {
  router.push(`/contracts/${row.id}`)
}

const createContract = () => {
  router.push('/contracts?action=create')
}

onMounted(async () => {
  await loadStatistics()
  await loadExpiringContracts()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.welcome-section {
  flex: 1;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1E293B;
  margin: 0 0 4px;
}

.page-desc {
  color: #64748B;
  margin: 0;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  border: none;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  overflow: hidden;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1E293B;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #64748B;
  margin-top: 2px;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #F1F5F9;
}

.trend-up, .trend-down {
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 500;
}

.trend-up { color: #10B981; }
.trend-down { color: #EF4444; }

.trend-label {
  color: #94A3B8;
  font-size: 12px;
}

.charts-row {
  margin-bottom: 24px;
}

.chart-card, .table-card, .overview-card {
  border: none;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  font-size: 15px;
}

.card-header span {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #1E293B;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.overview-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: #F8FAFC;
  border-radius: 12px;
  transition: all 0.2s;
}

.overview-item:hover {
  background: #F1F5F9;
}

.overview-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.overview-value {
  font-size: 18px;
  font-weight: 700;
  color: #1E293B;
}

.overview-label {
  font-size: 12px;
  color: #64748B;
  margin-top: 2px;
}

.text-danger {
  color: #EF4444;
  font-weight: 500;
}

.amount {
  font-weight: 600;
  color: #F59E0B;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #F1F5F9;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-table) {
  font-size: 14px;
}

:deep(.el-table th) {
  background: #F8FAFC !important;
  color: #64748B;
  font-weight: 600;
}

:deep(.el-radio-button__inner) {
  border-radius: 8px !important;
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  border: none;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 100%);
}
</style>
