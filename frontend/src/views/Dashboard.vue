<template>
  <div class="dashboard">
    <div class="page-header">
      <div class="welcome-section">
        <h1 class="page-title">仪表盘</h1>
        <p class="page-desc">欢迎回来，{{ userStore.userInfo?.username }}！</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" :icon="Plus" @click="createContract" v-if="canCreateContract">新建合同</el-button>
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
                <div class="overview-value">{{ statistics.this_month_contracts }}</div>
                <div class="overview-label">本月新增</div>
              </div>
            </div>
            <div class="overview-item">
              <div class="overview-icon" style="background: linear-gradient(135deg, #10B981, #34D399)">
                <el-icon :size="22"><Money /></el-icon>
              </div>
              <div class="overview-content">
                <div class="overview-value">¥{{ formatAmount(statistics.this_month_amount) }}</div>
                <div class="overview-label">本月金额</div>
              </div>
            </div>
            <div class="overview-item">
              <div class="overview-icon" style="background: linear-gradient(135deg, #F59E0B, #FBBF24)">
                <el-icon :size="22"><CircleCheck /></el-icon>
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ statistics.active_contracts }}</div>
                <div class="overview-label">进行中</div>
              </div>
            </div>
            <div class="overview-item">
              <div class="overview-icon" style="background: linear-gradient(135deg, #EF4444, #F87171)">
                <el-icon :size="22"><Warning /></el-icon>
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ statistics.expiring_soon }}</div>
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
              <el-button type="primary" size="small" @click="$router.push('/reminders')">
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
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
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
const userRole = computed(() => userStore.userInfo?.role || 'user')
const canCreateContract = computed(() => userRole.value === 'user' || userRole.value === 'admin')

const chartRef = ref(null)
const chartType = ref('bar')
const chartInstance = ref(null)
const expiringContracts = ref([])
const statistics = ref({
  total_contracts: 0,
  active_contracts: 0,
  pending_contracts: 0,
  completed_contracts: 0,
  draft_contracts: 0,
  terminated_contracts: 0,
  this_month_contracts: 0,
  this_month_amount: 0,
  expiring_soon: 0,
  total_amount: 0
})

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
    console.log('API返回数据:', data)
    if (data) {
      statistics.value = {
        total_contracts: data.total_contracts ?? 0,
        active_contracts: data.active_contracts ?? 0,
        pending_contracts: data.pending_contracts ?? 0,
        completed_contracts: data.completed_contracts ?? 0,
        draft_contracts: data.draft_contracts ?? 0,
        terminated_contracts: data.terminated_contracts ?? 0,
        this_month_contracts: data.this_month_contracts ?? 0,
        this_month_amount: data.this_month_amount ?? 0,
        expiring_soon: data.expiring_soon ?? 0,
        total_amount: data.total_amount ?? 0
      }
      console.log('统计数据显示:', statistics.value)
      nextTick(() => {
        initChart()
      })
    }
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
  
  if (chartInstance.value) {
    chartInstance.value.dispose()
  }
  
  chartInstance.value = echarts.init(chartRef.value)
  
  let option
  if (chartType.value === 'pie') {
    option = getPieOption()
  } else {
    option = getBarOption()
  }
  
  chartInstance.value.setOption(option)
  
  chartInstance.value.off('click')
  chartInstance.value.on('click', (params) => {
    if (chartType.value === 'pie') {
      const statusMap = {
        '进行中': 'active',
        '待审批': 'pending',
        '已完成': 'completed',
        '草稿': 'draft',
        '已终止': 'terminated'
      }
      const status = statusMap[params.name]
      if (status) {
        router.push({ path: '/contracts', query: { status } })
      }
    } else {
      const barStatusMap = {
        '进行中': 'active',
        '待审批': 'pending',
        '已完成': 'completed',
        '草稿': 'draft',
        '已终止': 'terminated'
      }
      const status = barStatusMap[params.name]
      if (status) {
        router.push({ path: '/contracts', query: { status } })
      }
    }
  })
  
  window.addEventListener('resize', () => {
    if (chartInstance.value) {
      chartInstance.value.resize()
    }
  })
}

watch(chartType, () => {
  if (chartInstance.value) {
    initChart()
  }
})

watch(statistics, () => {
  if (chartInstance.value) {
    initChart()
  }
}, { deep: true })

const getPieOption = () => ({
  tooltip: { 
    trigger: 'item', 
    formatter: '{b}: {c} 个'
  },
  legend: { 
    orient: 'vertical', 
    right: 20, 
    top: 'center',
    textStyle: { color: '#64748B', fontSize: 13 }
  },
  color: ['#6366F1', '#10B981', '#F59E0B', '#EF4444', '#94A3B8', '#EC4899', '#8B5CF6', '#14B8A6'],
  graphic: [{
    type: 'text',
    left: 'center',
    top: '38%',
    z: 100,
    style: {
      text: String(statistics.value.total_contracts || 0),
      textAlign: 'center',
      fill: '#1E293B',
      fontSize: 32,
      fontWeight: 'bold'
    }
  }, {
    type: 'text',
    left: 'center',
    top: '52%',
    z: 100,
    style: {
      text: '合同总数',
      textAlign: 'center',
      fill: '#94A3B8',
      fontSize: 14
    }
  }],
  series: [{
    type: 'pie',
    radius: ['45%', '70%'],
    center: ['50%', '50%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    emphasis: {
      scaleSize: 6,
      label: { show: true, fontSize: 13, fontWeight: 'bold' },
      itemStyle: { shadowBlur: 12, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.15)' }
    },
    data: [
      { value: statistics.value.active_contracts || 0, name: '进行中' },
      { value: statistics.value.pending_contracts || 0, name: '待审批' },
      { value: statistics.value.completed_contracts || 0, name: '已完成' },
      { value: statistics.value.draft_contracts || 0, name: '草稿' },
      { value: statistics.value.terminated_contracts || 0, name: '已终止' }
    ]
  }]
})

const getBarOption = () => ({
  tooltip: { 
    trigger: 'axis',
    formatter: '{b}: {c} 个合同' 
  },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { 
    type: 'category', 
    data: ['进行中', '待审批', '已完成', '草稿', '已终止'],
    axisLine: { lineStyle: { color: '#E2E8F0' } },
    axisLabel: { color: '#64748B', fontSize: 12 }
  },
  yAxis: { 
    type: 'value',
    axisLine: { show: false },
    axisLabel: { color: '#64748B' },
    splitLine: { lineStyle: { color: '#F1F5F9' } }
  },
  series: [{
    data: [
      statistics.value.active_contracts || 0,
      statistics.value.pending_contracts || 0,
      statistics.value.completed_contracts || 0,
      statistics.value.draft_contracts || 0,
      statistics.value.terminated_contracts || 0
    ],
    type: 'bar',
    barWidth: '50%',
    itemStyle: {
      borderRadius: [6, 6, 0, 0],
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: '#6366F1' },
        { offset: 1, color: '#8B5CF6' }
      ])
    },
    emphasis: {
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#4F46E5' },
          { offset: 1, color: '#7C3AED' }
        ])
      }
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

onUnmounted(() => {
  if (chartInstance.value) {
    chartInstance.value.dispose()
    chartInstance.value = null
  }
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
  width: 100%;
}

.card-header span {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #1E293B;
}

.card-header .el-button {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
  padding: 8px 16px;
  border-radius: 8px;
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
