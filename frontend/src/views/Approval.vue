<template>
  <div class="approval-page">
    <el-card>
      <template #header>
        <div class="header">
          <span>合同审批</span>
          <el-tag type="info">当前角色: {{ userRole }}</el-tag>
        </div>
      </template>
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="contract_no" label="合同编号" width="150" />
        <el-table-column prop="title" label="合同标题" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="当前状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getContractStatusType(row.status)">{{ getContractStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button v-if="row.status === 'draft'" type="warning" link @click="handleSubmit(row)">提交审批</el-button>
            <el-button v-else-if="row.status === 'pending'" type="success" link @click="handleApprove(row)">{{ userRole === 'admin' ? '二级审批' : '一级审批' }}</el-button>
            <el-button v-else-if="row.status === 'active'" type="info" link disabled>已生效</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card style="margin-top: 20px" v-if="userRole === 'admin' || userRole === 'manager'">
      <template #header>
        <div class="header">
          <span>状态变更审批</span>
          <el-tag type="warning">{{ statusChangeData.length }} 条待审批</el-tag>
        </div>
      </template>
      <el-table :data="statusChangeData" style="width: 100%" v-loading="statusChangeLoading">
        <el-table-column prop="contract.contract_no" label="合同编号" width="150" />
        <el-table-column prop="contract.title" label="合同标题" />
        <el-table-column prop="from_status" label="原状态" width="100">
          <template #default="{ row }">
            <el-tag>{{ getStatusText(row.from_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="to_status" label="目标状态" width="100">
          <template #default="{ row }">
            <el-tag type="warning">{{ getStatusText(row.to_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="申请原因" />
        <el-table-column prop="requester.full_name" label="申请人" width="100" />
        <el-table-column prop="created_at" label="申请时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewStatusChange(row)">查看</el-button>
            <el-button type="success" link @click="handleApproveStatusChange(row)">通过</el-button>
            <el-button type="danger" link @click="handleRejectStatusChange(row)">拒绝</el-button>
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
        <el-form-item label="审批级别">
          <el-tag :type="userRole === 'admin' ? 'danger' : 'warning'">
            {{ userRole === 'admin' ? '二级审批 (管理员)' : '一级审批 (经理)' }}
          </el-tag>
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
        <el-button type="primary" @click="handleSubmitApproval">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'
import { getPendingApprovals, createApproval, updateApproval } from '@/api/approval'
import { getPendingStatusChangeApprovals, approveStatusChangeRequest, rejectStatusChangeRequest } from '@/api/contract'

const router = useRouter()
const userStore = useUserStore()
const userRole = userStore.userInfo?.role || 'user'

const loading = ref(false)
const statusChangeLoading = ref(false)
const dialogVisible = ref(false)
const formRef = ref(null)
const tableData = ref([])
const statusChangeData = ref([])
const currentContract = ref({})
const currentApprovalId = ref(null)

const formData = reactive({
  status: 'approved',
  comment: ''
})

const getStatusText = (status) => {
  const map = {
    draft: '草稿',
    pending: '待审批',
    approved: '已批准',
    active: '已生效',
    in_progress: '执行中',
    pending_pay: '待付款',
    completed: '已完成',
    terminated: '已终止',
    archived: '已归档'
  }
  return map[status] || status
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return dateStr
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = date.getHours()
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const ampm = hours < 12 ? '上午' : '下午'
  const hour12 = hours % 12 || 12
  return `${year}-${month}-${day} ${ampm}${hour12}:${minutes}`
}

const getContractStatusType = (status) => {
  const typeMap = {
    draft: 'info',
    pending: 'warning',
    active: 'success',
    approved: 'success',
    rejected: 'danger'
  }
  return typeMap[status] || ''
}

const getContractStatusText = (status) => {
  const textMap = {
    draft: '草稿',
    pending: '待审批',
    active: '进行中',
    approved: '已批准',
    rejected: '已拒绝'
  }
  return textMap[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    const data = await getPendingApprovals()
    tableData.value = data
  } finally {
    loading.value = false
  }
  
  if (userRole === 'admin' || userRole === 'manager') {
    statusChangeLoading.value = true
    try {
      const data = await getPendingStatusChangeApprovals()
      statusChangeData.value = data
    } finally {
      statusChangeLoading.value = false
    }
  }
}

const handleViewStatusChange = (row) => {
  router.push(`/contracts/${row.contract_id}`)
}

const handleApproveStatusChange = async (row) => {
  try {
    await ElMessageBox.confirm(`确定通过将合同 "${row.contract.title}" 状态变更为 "${getStatusText(row.to_status)}" 吗？`, '审批确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'success'
    })
    await approveStatusChangeRequest(row.id, { comment: '同意' })
    ElMessage.success('审批通过')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '操作失败')
    }
  }
}

const handleRejectStatusChange = async (row) => {
  try {
    await ElMessageBox.confirm(`确定拒绝将合同 "${row.contract.title}" 状态变更申请吗？`, '审批确认', {
      confirmButtonText: '确定拒绝',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await rejectStatusChangeRequest(row.id, { comment: '拒绝' })
    ElMessage.success('已拒绝')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '操作失败')
    }
  }
}

const handleView = (row) => {
  router.push(`/contracts/${row.id}`)
}

const handleSubmit = async (row) => {
  const level = userRole === 'admin' ? 2 : 1
  await createApproval({
    contract_id: row.id,
    level: level,
    status: 'pending',
    comment: '提交审批'
  })
  ElMessage.success('已提交审批')
  loadData()
}

const handleApprove = async (row) => {
  currentContract.value = row
  currentApprovalId.value = row.approval_id
  formData.status = 'approved'
  formData.comment = ''
  dialogVisible.value = true
}

const handleSubmitApproval = async () => {
  await updateApproval(currentApprovalId.value, {
    status: formData.status,
    comment: formData.comment
  })
  ElMessage.success(formData.status === 'approved' ? '审批通过' : '审批拒绝')
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

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>