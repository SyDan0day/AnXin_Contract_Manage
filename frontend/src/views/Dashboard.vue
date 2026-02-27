<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409eff">
              <el-icon size="30"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.total_contracts }}</div>
              <div class="stat-label">合同总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67c23a">
              <el-icon size="30"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.active_contracts }}</div>
              <div class="stat-label">进行中合同</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e6a23c">
              <el-icon size="30"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.pending_contracts }}</div>
              <div class="stat-label">待审批合同</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f56c6c">
              <el-icon size="30"><Bell /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.expiring_soon }}</div>
              <div class="stat-label">即将到期</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="charts-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>合同金额统计</span>
          </template>
          <div ref="amountChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>本月新增合同</span>
          </template>
          <div class="monthly-stats">
            <div class="stat-item">
              <span class="label">本月合同数：</span>
              <span class="value">{{ statistics.this_month_contracts }}</span>
            </div>
            <div class="stat-item">
              <span class="label">本月合同金额：</span>
              <span class="value">¥{{ statistics.this_month_amount?.toFixed(2) }}</span>
            </div>
          </div>
          <div ref="monthlyChartRef" style="height: 240px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="expiring-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>即将到期合同</span>
          </template>
          <el-table :data="expiringContracts" style="width: 100%">
            <el-table-column prop="contract_no" label="合同编号" width="150" />
            <el-table-column prop="title" label="合同标题" />
            <el-table-column prop="end_date" label="到期日期" width="120" />
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="{ row }">
                ¥{{ row.amount?.toFixed(2) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { getStatistics, getExpiringContracts } from '@/api/approval'

const statistics = ref({})
const expiringContracts = ref([])
const amountChartRef = ref(null)
const monthlyChartRef = ref(null)

const loadStatistics = async () => {
  const data = await getStatistics()
  statistics.value = data
}

const loadExpiringContracts = async () => {
  const data = await getExpiringContracts(30)
  expiringContracts.value = data.contracts
}

const initCharts = () => {
  if (amountChartRef.value) {
    const amountChart = echarts.init(amountChartRef.value)
    amountChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [
        {
          type: 'pie',
          radius: '60%',
          data: [
            { value: statistics.value.total_amount, name: '总金额' }
          ],
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    })
  }

  if (monthlyChartRef.value) {
    const monthlyChart = echarts.init(monthlyChartRef.value)
    monthlyChart.setOption({
      xAxis: { type: 'category', data: ['1月', '2月', '3月', '4月', '5月', '6月'] },
      yAxis: { type: 'value' },
      series: [
        {
          data: [0, 0, 0, 0, 0, statistics.value.this_month_contracts],
          type: 'bar'
        }
      ]
    })
  }
}

onMounted(async () => {
  await loadStatistics()
  await loadExpiringContracts()
  initCharts()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  cursor: pointer;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-right: 16px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.charts-row {
  margin-bottom: 20px;
}

.monthly-stats {
  padding: 20px 0;
}

.stat-item {
  margin-bottom: 10px;
  font-size: 14px;
}

.stat-item .label {
  color: #606266;
}

.stat-item .value {
  font-weight: bold;
  color: #303133;
  margin-left: 10px;
}
</style>