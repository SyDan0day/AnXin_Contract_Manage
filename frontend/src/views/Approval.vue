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
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" link @click="handleView(row)">
                <el-icon><View /></el-icon><span class="btn-text">查看</span>
              </el-button>
              <el-button v-if="row.status === 'draft'" type="warning" link @click="handleSubmit(row)">
                <el-icon><Position /></el-icon><span class="btn-text">提交</span>
              </el-button>
              <el-button v-else-if="row.status === 'pending'" type="primary" link @click="handleApprove(row)">
                <el-icon><Edit /></el-icon><span class="btn-text">审批</span>
              </el-button>
              <el-button v-if="row.status === 'pending' && userRole === 'user'" type="danger" link @click="handleWithdraw(row)">
                <el-icon><Close /></el-icon><span class="btn-text">撤回</span>
              </el-button>
              <el-button v-else-if="row.status === 'active'" type="info" link disabled>
                已生效
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card style="margin-top: 20px" v-if="userRole === 'admin' || userRole === 'sales_manager'">
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
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" link @click="handleViewStatusChange(row)">
                <el-icon><View /></el-icon> 查看
              </el-button>
              <el-button type="success" link @click="handleApproveStatusChange(row)">
                <el-icon><Check /></el-icon> 通过
              </el-button>
              <el-button type="danger" link @click="handleRejectStatusChange(row)">
                <el-icon><Close /></el-icon> 拒绝
              </el-button>
            </div>
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
          <el-tag type="warning">
            {{ getApprovalLevelText(currentApprovalLevel) }}
          </el-tag>
        </el-form-item>
        <el-form-item label="审批结果" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio label="approved">同意</el-radio>
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
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmitApproval">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Position, Check, Close, Edit } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { getPendingApprovals, createApproval, withdrawApproval, updateApproval } from '@/api/approval'
import { getPendingStatusChangeApprovals, approveStatusChangeRequest, rejectStatusChangeRequest } from '@/api/contract'
import axios from 'axios'

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
const currentWorkflowId = ref(null)
const currentApprovalLevel = ref(null)

const formData = reactive({
  status: 'approved',
  comment: ''
})

const getStatusText = (status) => {
  const map = {
    draft: '草稿',
    pending: '待审批',
    active: '已生效',
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
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
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

const getApprovalLevelText = (level) => {
  const levelMap = {
    1: '一级审批 (销售负责人)',
    2: '二级审批 (技术负责人)',
    3: '三级审批 (财务负责人)'
  }
  return levelMap[level] || `第${level}级审批`
}

const loadData = async () => {
  loading.value = true
  try {
    const data = await getPendingApprovals()
    tableData.value = data
  } finally {
    loading.value = false
  }
  
  if (userRole === 'admin' || userRole === 'sales_manager') {
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
  currentWorkflowId.value = row.workflow_id
  currentApprovalLevel.value = row.approval_level
  formData.status = 'approved'
  formData.comment = ''
  dialogVisible.value = true
}

const handleWithdraw = async (row) => {
  try {
    await ElMessageBox.confirm(`确定撤回合同 "${row.title}" 的审批申请吗？`, '撤回确认', {
      confirmButtonText: '确定撤回',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await withdrawApproval(row.id)
    ElMessage.success('已撤回审批申请')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '撤回失败')
    }
  }
}

const handleSubmitApproval = async () => {
  try {
    const response = await axios.put(`/api/approvals/${currentApprovalId.value}`, {
      status: formData.status,
      comment: formData.comment
    }, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data) {
      if (formData.status === 'approved') {
        ElMessage.success('审批通过')
      } else {
        ElMessage.warning('审批已拒绝，合同已退回给销售')
      }
      dialogVisible.value = false
      loadData()
    }
  } catch (error) {
    ElMessage.error('审批失败: ' + (error.response?.data?.error || error.message || '未知错误'))
  }
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
  width: 100%;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-wrap: wrap;
}

.action-buttons .el-button {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 4px 6px;
  min-width: auto;
}

.action-buttons .btn-text {
  display: inline;
}

@media (max-width: 768px) {
  .btn-text {
    display: none;
  }
  .action-buttons {
    gap: 0;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>