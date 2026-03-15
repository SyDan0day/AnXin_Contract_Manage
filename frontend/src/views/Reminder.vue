<template>
  <div class="reminder-page">
    <el-card>
      <template #header>
        <span>到期提醒</span>
      </template>
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="提前天数">
          <el-input-number v-model="searchForm.days" :min="1" :max="365" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="contract_no" label="合同编号" width="150" />
        <el-table-column prop="title" label="合同标题" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="end_date" label="到期日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'primary' : 'info'">
              {{ row.status === 'active' ? '进行中' : '已结束' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleCreateReminder(row)">创建提醒</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="创建提醒" width="500px">
      <el-form ref="formRef" :model="formData" label-width="100px">
        <el-form-item label="合同编号">
          <el-input v-model="currentContract.contract_no" disabled />
        </el-form-item>
        <el-form-item label="提醒类型">
          <el-select v-model="formData.type" style="width: 100%">
            <el-option label="合同到期" value="expiry" />
            <el-option label="付款提醒" value="payment" />
          </el-select>
        </el-form-item>
        <el-form-item label="提前天数">
          <el-input-number v-model="formData.days_before" :min="1" :max="365" style="width: 100%" />
        </el-form-item>
        <el-form-item label="提醒日期">
          <el-date-picker
            v-model="formData.reminder_date"
            type="date"
            placeholder="请选择提醒日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getExpiringContracts, createReminder } from '@/api/approval'

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return dateStr
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const loading = ref(false)
const dialogVisible = ref(false)
const formRef = ref(null)
const tableData = ref([])
const currentContract = ref({})

const searchForm = reactive({
  days: 30
})

const formData = reactive({
  contract_id: null,
  type: 'expiry',
  days_before: 7,
  reminder_date: ''
})

const loadData = async () => {
  loading.value = true
  try {
    const data = await getExpiringContracts(searchForm.days)
    tableData.value = data.contracts
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  loadData()
}

const handleCreateReminder = async (row) => {
  currentContract.value = row
  formData.contract_id = row.id
  formData.reminder_date = ''
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await createReminder({
    contract_id: formData.contract_id,
    type: formData.type,
    reminder_date: formData.reminder_date,
    days_before: formData.days_before
  })
  ElMessage.success('创建提醒成功')
  dialogVisible.value = false
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.reminder-page {
  padding: 20px;
}
</style>