<template>
  <div class="approval-page">
    <el-card>
      <template #header>
        <span>审批管理</span>
      </template>
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="contract_no" label="合同编号" width="150" />
        <el-table-column prop="title" label="合同标题" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="审批状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.approval_status)">{{ getStatusText(row.approval_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleApprove(row)" :disabled="row.approval_status !== 'pending'">审批</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="合同审批" width="600px">
      <el-form ref="formRef" :model="formData" label-width="100px">
        <el-form-item label="合同编号">
          <el-input v-model="currentContract.contract_no" disabled />
        </el-form-item>
        <el-form-item label="合同标题">
          <el-input v-model="currentContract.title" disabled />
        </el-form-item>
        <el-form-item label="合同金额">
          <el-input :value="'¥' + currentContract.amount?.toFixed(2)" disabled />
        </el-form-item>
        <el-form-item label="审批结果" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审批意见">
          <el-input
            v-model="formData.comment"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见"
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
import { getContractList, getApprovalRecords, createApproval, updateApproval } from '@/api/contract'
import { getApprovalRecords as getApprovals } from '@/api/approval'

const loading = ref(false)
const dialogVisible = ref(false)
const formRef = ref(null)
const tableData = ref([])
const currentContract = ref({})

const formData = reactive({
  status: 'approved',
  comment: ''
})

const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return typeMap[status] || ''
}

const getStatusText = (status) => {
  const textMap = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return textMap[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    const data = await getContractList({ status: 'pending', limit: 100 })
    tableData.value = data.map(contract => ({
      ...contract,
      approval_status: 'pending'
    }))
  } finally {
    loading.value = false
  }
}

const handleApprove = async (row) => {
  currentContract.value = row
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await createApproval({
    contract_id: currentContract.value.id,
    status: formData.status,
    comment: formData.comment
  })
  ElMessage.success('审批成功')
  dialogVisible.value = false
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.approval-page {
  padding: 20px;
}
</style>