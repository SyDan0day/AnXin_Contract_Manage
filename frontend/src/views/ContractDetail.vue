<template>
  <div class="contract-detail">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-button text @click="$router.back()">
            <el-icon><ArrowLeft /></el-icon> 返回
          </el-button>
          <span class="title">合同详情</span>
          <el-button type="primary" @click="handleEdit">编辑合同</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" type="border-card" @tab-change="tabChange">
        <el-tab-pane label="基本信息" name="info">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="合同编号">{{ contract.contract_no }}</el-descriptions-item>
            <el-descriptions-item label="合同标题">{{ contract.title }}</el-descriptions-item>
            <el-descriptions-item label="客户名称">{{ contract.customer?.name }}</el-descriptions-item>
            <el-descriptions-item label="合同类型">{{ contract.contract_type?.name }}</el-descriptions-item>
            <el-descriptions-item label="金额">
              <span class="amount">¥{{ contract.amount?.toLocaleString() }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(contract.status)">{{ getStatusText(contract.status) }}</el-tag>
              <el-button type="primary" link size="small" style="margin-left: 8px" @click="showStatusDialog = true">
                变更状态
              </el-button>
              <el-button v-if="contract.status !== 'archived'" type="warning" link size="small" style="margin-left: 8px" @click="handleArchive">
                归档
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="签约日期">{{ formatDate(contract.sign_date) }}</el-descriptions-item>
            <el-descriptions-item label="开始日期">{{ formatDate(contract.start_date) }}</el-descriptions-item>
            <el-descriptions-item label="到期日期">{{ formatDate(contract.end_date) }}</el-descriptions-item>
            <el-descriptions-item label="付款条件" :span="2">{{ contract.payment_terms || '-' }}</el-descriptions-item>
            <el-descriptions-item label="合同内容" :span="2">{{ contract.content || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建人">{{ contract.creator?.full_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDateTime(contract.created_at) }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>

        <el-tab-pane label="生命周期" name="lifecycle">
          <div class="tab-header">
            <span>合同生命周期跟踪</span>
          </div>
          <el-timeline>
            <el-timeline-item
              v-for="(event, index) in lifecycleEvents"
              :key="index"
              :timestamp="formatDateTime(event.created_at)"
              :type="getLifecycleItemType(event.event_type)"
              :hollow="event.event_type === 'progress'"
            >
              <div class="lifecycle-content">
                <div class="lifecycle-title">{{ getLifecycleTitle(event.event_type) }}</div>
                <div class="lifecycle-desc">
                  {{ event.from_status ? `${getStatusText(event.from_status)} → ${getStatusText(event.toStatus)}` : '' }}
                  {{ event.description || '' }}
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-if="lifecycleEvents.length === 0" description="暂无生命周期记录" />
        </el-tab-pane>

        <el-tab-pane label="执行跟踪" name="executions">
          <div class="tab-header">
            <span>执行进度管理</span>
            <el-button type="primary" size="small" @click="showExecutionDialog = true">
              <el-icon><Plus /></el-icon> 添加执行记录
            </el-button>
          </div>
          <el-table :data="executions" v-loading="executionsLoading">
            <el-table-column prop="stage" label="阶段" />
            <el-table-column prop="stage_date" label="阶段日期" width="120">
              <template #default="{ row }">
                {{ formatDate(row.stage_date) }}
              </template>
            </el-table-column>
            <el-table-column prop="progress" label="进度" width="150">
              <template #default="{ row }">
                <el-progress :percentage="row.progress" :color="getProgressColor(row.progress)" />
              </template>
            </el-table-column>
            <el-table-column prop="payment_amount" label="付款金额" width="120">
              <template #default="{ row }">
                ¥{{ row.payment_amount?.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column prop="payment_date" label="付款日期" width="120">
              <template #default="{ row }">
                {{ formatDate(row.payment_date) }}
              </template>
            </el-table-column>
            <el-table-column prop="description" label="说明" />
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button type="danger" link @click="handleDeleteExecution(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="文档管理" name="documents">
          <div class="tab-header">
            <span>合同文档</span>
            <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              :data="uploadData"
              :show-file-list="false"
              :accept="'.doc,.docx,.pdf,.jpg,.jpeg,.png,.gif,.bmp,.webp,.txt,.html,.htm,.xls,.xlsx'"
              :before-upload="handleBeforeUpload"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
            >
              <el-button type="primary" size="small">
                <el-icon><Upload /></el-icon> 上传文档
              </el-button>
              <template #tip>
                <div class="el-upload__tip" style="margin-top: 8px">支持 .doc, .docx, .pdf, 图片, 文本, Excel 格式</div>
              </template>
            </el-upload>
          </div>
          <el-table :data="documents" v-loading="documentsLoading">
            <el-table-column prop="name" label="文档名称" />
            <el-table-column prop="file_type" label="类型" width="100" />
            <el-table-column prop="file_size" label="大小" width="100">
              <template #default="{ row }">
                {{ formatFileSize(row.file_size) }}
              </template>
            </el-table-column>
            <el-table-column prop="version" label="版本" width="80" />
            <el-table-column prop="created_at" label="上传时间" width="180" />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="handlePreview(row)">预览</el-button>
                <el-button type="primary" link @click="handleDownload(row)">下载</el-button>
                <el-button type="danger" link @click="handleDeleteDocument(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="审批记录" name="approvals">
          <div class="tab-header">
            <span>审批历史</span>
            <el-button type="primary" size="small" @click="showApprovalDialog = true" v-if="contract.status === 'draft' || contract.status === 'pending'">
              <el-icon><Plus /></el-icon> 提交审批
            </el-button>
          </div>
          <el-table :data="approvals" v-loading="approvalsLoading">
            <el-table-column prop="approver.full_name" label="审批人" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getApprovalStatusType(row.status)">{{ getApprovalStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="comment" label="审批意见" />
            <el-table-column prop="approved_at" label="审批时间" width="180" />
            <el-table-column prop="created_at" label="提交时间" width="180" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="showExecutionDialog" title="添加执行记录" width="500px">
      <el-form ref="executionFormRef" :model="executionForm" :rules="executionRules" label-width="100px">
        <el-form-item label="阶段名称" prop="stage">
          <el-input v-model="executionForm.stage" placeholder="请输入阶段名称" />
        </el-form-item>
        <el-form-item label="阶段日期" prop="stage_date">
          <el-date-picker v-model="executionForm.stage_date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="进度%" prop="progress">
          <el-slider v-model="executionForm.progress" :marks="{0: '0%', 50: '50%', 100: '100%'}" />
          <div style="font-size: 12px; color: #999; margin-top: 4px">根据付款金额自动计算（合同总金额：¥{{ contractAmount.toLocaleString() }}）</div>
        </el-form-item>
        <el-form-item label="付款金额">
          <el-input-number v-model="executionForm.payment_amount" :precision="2" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="付款日期">
          <el-date-picker v-model="executionForm.payment_date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="executionForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showExecutionDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitExecution">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showApprovalDialog" title="提交审批" width="500px">
      <el-form ref="approvalFormRef" :model="approvalForm" :rules="approvalRules" label-width="100px">
        <el-form-item label="审批意见" prop="comment">
          <el-input v-model="approvalForm.comment" type="textarea" :rows="4" placeholder="请输入审批意见" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showApprovalDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitApproval">提交审批</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showStatusDialog" title="变更合同状态" width="500px">
      <el-form label-width="100px">
        <el-form-item label="当前状态">
          <el-tag :type="getStatusType(contract.status)">{{ getStatusText(contract.status) }}</el-tag>
        </el-form-item>
        <el-form-item label="变更为">
          <el-select v-model="newStatus" placeholder="请选择新状态" style="width: 100%">
            <el-option v-for="opt in statusOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="statusDescription" type="textarea" :rows="3" placeholder="请输入变更说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showStatusDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateStatus">确定</el-button>
      </template>
    </el-dialog>
    
    <!-- 文档预览对话框 -->
    <el-dialog v-model="showPreviewDialog" :title="previewFileName" width="90%" top="5vh" destroy-on-close>
      <div class="preview-container">
        <div v-if="previewLoading" class="preview-loading">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>正在加载预览...</span>
        </div>
        <div v-else-if="previewError" class="preview-error">
          <el-icon><Warning /></el-icon>
          <span>{{ previewError }}</span>
          <el-button type="primary" @click="downloadPreviewFile">下载文件</el-button>
        </div>
        
        <!-- Word 文档或文本文件内容预览 - 表格形式 -->
        <div v-else-if="previewData && previewData.content" class="content-preview">
          <!-- 合同关键信息表格 -->
          <el-table 
            v-if="previewData.fields && previewData.fields.length > 0"
            :data="previewData.fields" 
            border 
            style="width: 100%; margin-bottom: 20px;"
            :header-cell-style="{background: '#f5f7fa', color: '#409eff', fontWeight: 'bold'}"
          >
            <el-table-column prop="label" label="字段名称" width="180" align="center" />
            <el-table-column prop="value" label="提取内容" min-width="300" align="left" />
          </el-table>
          
          <!-- 完整文档内容 -->
          <div class="content-section">
            <h4 style="margin: 0 0 10px 0; color: #409eff;">📄 完整文档内容</h4>
            <div class="document-content">
              <pre>{{ previewData.content }}</pre>
            </div>
          </div>
        </div>
        
        <!-- 其他文件类型使用 iframe 预览 -->
        <iframe 
          v-else-if="previewUrl && previewUrl.startsWith('http')" 
          :src="previewUrl" 
          class="preview-iframe"
          frameborder="0"
        ></iframe>
        
        <!-- 调试信息 -->
        <div v-else style="padding: 20px; background: #f0f0f0;">
          <p><strong>调试信息：</strong></p>
          <p>previewData: {{ previewData }}</p>
          <p>previewUrl: {{ previewUrl }}</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="showPreviewDialog = false">关闭</el-button>
        <el-button type="primary" @click="downloadPreviewFile">下载文件</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Plus, Upload, Loading, Warning } from '@element-plus/icons-vue'
import { getContractDetail, getContractExecutions, createContractExecution, deleteExecution, getContractDocuments, uploadDocument, deleteDocument, getContractLifecycle, updateContractStatus, archiveContract, requestStatusChange } from '@/api/contract'
import { getApprovalRecords, createApproval } from '@/api/approval'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('info')
const contract = ref({})
const executions = ref([])
const documents = ref([])
const approvals = ref([])
const executionsLoading = ref(false)
const documentsLoading = ref(false)
const approvalsLoading = ref(false)
const lifecycleEvents = ref([])

const showExecutionDialog = ref(false)
const showApprovalDialog = ref(false)
const showStatusDialog = ref(false)
const showPreviewDialog = ref(false)
const newStatus = ref('')
const statusDescription = ref('')
const executionFormRef = ref(null)
const approvalFormRef = ref(null)

// 预览相关状态
const previewUrl = ref('')
const previewFileName = ref('')
const previewLoading = ref(false)
const previewError = ref('')
const currentPreviewDocument = ref(null)
const previewData = ref(null) // 存储解析后的文档数据

const statusOptions = [
  { value: 'draft', label: '草稿' },
  { value: 'pending', label: '待审批' },
  { value: 'approved', label: '已批准' },
  { value: 'active', label: '已生效' },
  { value: 'in_progress', label: '执行中' },
  { value: 'pending_pay', label: '待付款' },
  { value: 'completed', label: '已完成' },
  { value: 'terminated', label: '已终止' }
]

const executionForm = reactive({
  stage: '',
  stage_date: '',
  progress: 0,
  payment_amount: 0,
  payment_date: '',
  description: ''
})

const approvalForm = reactive({
  comment: ''
})

const executionRules = {
  stage: [{ required: true, message: '请输入阶段名称', trigger: 'blur' }]
}

const approvalRules = {
  comment: [{ required: true, message: '请输入审批意见', trigger: 'blur' }]
}

const contractId = computed(() => parseInt(route.params.id))

watch(() => route.params.id, () => {
  if (route.params.id) {
    loadContract()
    loadExecutions()
    loadDocuments()
    loadApprovals()
  }
})

const contractAmount = computed(() => contract.value.amount || 0)

watch(() => executionForm.payment_amount, (newVal) => {
  if (contractAmount.value > 0) {
    executionForm.progress = Math.round((newVal / contractAmount.value) * 100)
  }
})

watch(() => executionForm.progress, (newVal) => {
  if (contractAmount.value > 0) {
    executionForm.payment_amount = Math.round((newVal / 100) * contractAmount.value * 100) / 100
  }
})

const uploadUrl = computed(() => `/api/contracts/${contractId.value}/documents`)
const uploadHeaders = computed(() => ({ Authorization: `Bearer ${userStore.token}` }))

const API_BASE = '/api'
const uploadData = computed(() => ({ contract_id: contractId.value }))

const getStatusType = (status) => {
  const map = { 
    draft: 'info', 
    pending: 'warning', 
    approved: 'success', 
    active: 'primary',
    in_progress: 'primary',
    pending_pay: 'warning',
    completed: 'success', 
    terminated: 'danger',
    archived: 'info'
  }
  return map[status] || ''
}

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

const getLifecycleItemType = (eventType) => {
  const map = {
    created: 'primary',
    submitted: 'warning',
    approved: 'success',
    rejected: 'danger',
    activated: 'success',
    progress: 'primary',
    payment: 'warning',
    completed: 'success',
    terminated: 'danger',
    archived: 'info',
    status_changed: 'info'
  }
  return map[eventType] || 'info'
}

const getLifecycleTitle = (eventType) => {
  const map = {
    created: '合同创建',
    submitted: '提交审批',
    approved: '审批通过',
    rejected: '审批拒绝',
    activated: '合同生效',
    progress: '执行进度更新',
    payment: '付款记录',
    completed: '合同完成',
    terminated: '合同终止',
    archived: '合同归档',
    status_changed: '状态变更'
  }
  return map[eventType] || eventType
}

const getApprovalStatusType = (status) => {
  const map = { pending: 'warning', approved: 'success', rejected: 'danger' }
  return map[status] || ''
}

const getApprovalStatusText = (status) => {
  const map = { pending: '待审批', approved: '已批准', rejected: '已拒绝' }
  return map[status] || status
}

const getProgressColor = (progress) => {
  if (progress < 30) return '#EF4444'
  if (progress < 70) return '#F59E0B'
  return '#10B981'
}

const formatFileSize = (bytes) => {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
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

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return dateStr
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  
  return `${year}-${month}-${day}`
}

const loadContract = async () => {
  contract.value = await getContractDetail(contractId.value)
}

const loadExecutions = async () => {
  executionsLoading.value = true
  try {
    executions.value = await getContractExecutions(contractId.value)
  } finally {
    executionsLoading.value = false
  }
}

const loadDocuments = async () => {
  documentsLoading.value = true
  try {
    documents.value = await getContractDocuments(contractId.value)
  } finally {
    documentsLoading.value = false
  }
}

const loadApprovals = async () => {
  approvalsLoading.value = true
  try {
    approvals.value = await getApprovalRecords(contractId.value)
  } finally {
    approvalsLoading.value = false
  }
}

const loadLifecycle = async () => {
  try {
    lifecycleEvents.value = await getContractLifecycle(contractId.value)
  } catch (error) {
    console.error('加载生命周期记录失败:', error)
  }
}

const handleUpdateStatus = async () => {
  try {
    const res = await requestStatusChange(contractId.value, {
      to_status: newStatus.value,
      reason: statusDescription.value
    })
    if (res.direct) {
      ElMessage.success('状态更新成功')
    } else {
      ElMessage.success('状态变更申请已提交，等待管理员审批')
    }
    showStatusDialog.value = false
    newStatus.value = ''
    statusDescription.value = ''
    loadContract()
    loadLifecycle()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

const handleArchive = async () => {
  try {
    await ElMessageBox.confirm('归档操作需要管理员审批通过后生效，是否继续？', '合同归档', {
      confirmButtonText: '确定申请',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await requestStatusChange(contractId.value, {
      to_status: 'archived',
      reason: '申请归档'
    })
    if (res.direct) {
      ElMessage.success('合同归档成功')
    } else {
      ElMessage.success('归档申请已提交，等待管理员审批')
    }
    loadContract()
    loadLifecycle()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '操作失败')
    }
  }
}

const handleEdit = () => {
  router.push(`/contracts?action=edit&id=${contractId.value}`)
}

const handleSubmitExecution = async () => {
  await executionFormRef.value.validate()
  await createContractExecution({ ...executionForm, contract_id: contractId.value })
  ElMessage.success('添加成功')
  showExecutionDialog.value = false
  Object.assign(executionForm, { stage: '', stage_date: '', progress: 0, payment_amount: 0, payment_date: '', description: '' })
  loadExecutions()
}

const handleDeleteExecution = async (row) => {
  await ElMessageBox.confirm('确定删除该执行记录?', '提示', { type: 'warning' })
  await deleteExecution(row.id)
  ElMessage.success('删除成功')
  loadExecutions()
}

const handleBeforeUpload = (file) => {
  const allowedExtensions = ['doc', 'docx', 'pdf', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'txt', 'html', 'htm', 'xls', 'xlsx']
  const extension = file.name.split('.').pop().toLowerCase()
  
  if (!allowedExtensions.includes(extension)) {
    ElMessage.error('不支持的文件格式')
    return false
  }
  
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }
  return true
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  loadDocuments()
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

// 从文本中提取合同关键信息并返回数组
const extractContractFields = (text) => {
  const fields = []
  
  // 定义所有要提取的字段及其正则表达式
  const patterns = [
    { label: '合同编号', pattern: /合同编号[：:]\s*([A-Z0-9\-]+)/ },
    { label: '合同名称', pattern: /合同名称[：:]\s*([^\n]{2,60})/ },
    { label: '甲方（客户）', pattern: /甲方[（(]?客户[）)]?[：:]\s*([^\n]{2,60})/ },
    { label: '乙方', pattern: /乙方[：:]\s*([^\n]{2,60})/ },
    { label: '合同金额', pattern: /合同金额[：:]\s*([\d,]+\.?\d*)\s*(?:元|万)?/ },
    { label: '签订日期', pattern: /签订日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)/ },
    { label: '开始日期', pattern: /开始日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)/ },
    { label: '结束日期', pattern: /结束日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)/ },
    { label: '合同类型', pattern: /合同类型[：:]\s*([^\n]{2,30})/ },
    { label: '联系人', pattern: /联系人[：:]\s*([^\n]{2,20})/ },
    { label: '联系电话', pattern: /(?:联系电话|电话)[：:]\s*([\d\-]{7,20})/ },
    { label: '地址', pattern: /(?:地址)[：:]\s*([^\n]{5,100})/ },
    { label: '开户银行', pattern: /开户银行[：:]\s*([^\n]{2,50})/ },
    { label: '银行账号', pattern: /(?:银行账号|账号)[：:]\s*([\d]{10,30})/ },
  ]
  
  // 提取每个字段
  for (const item of patterns) {
    const match = text.match(item.pattern)
    if (match) {
      fields.push({
        label: item.label,
        value: match[1].trim(),
        found: true
      })
    }
  }
  
  return fields
}

const handlePreview = async (row) => {
  const fileExt = row.name.split('.').pop().toLowerCase()
  const token = localStorage.getItem('token')
  
  // 存储当前预览的文档信息
  currentPreviewDocument.value = row
  previewFileName.value = row.name
  previewError.value = ''
  previewLoading.value = true
  previewData.value = null
  previewUrl.value = ''
  
  try {
    // Word 文档 (.docx) 和文本文件
    if (fileExt === 'docx' || fileExt === 'txt') {
      const response = await axios.get(`/api/documents/${row.id}/preview`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      
      console.log('API响应类型:', typeof response.data)
      console.log('API响应数据:', response.data)
      console.log('content字段类型:', typeof response.data?.content)
      console.log('content字段值:', response.data?.content)
      
      // 提取合同关键信息并生成表格数据
      const contentText = response.data?.content || response.data
      const fields = extractContractFields(contentText)
      
      previewData.value = {
        content: contentText,
        fields: fields
      }
      
      console.log('previewData设置完成:', previewData.value)
      showPreviewDialog.value = true
    } else if (fileExt === 'pdf' || fileExt === 'jpg' || fileExt === 'jpeg' || 
               fileExt === 'png' || fileExt === 'gif' || fileExt === 'bmp' || 
               fileExt === 'webp' || fileExt === 'xls' || fileExt === 'xlsx') {
      // 对于其他文件类型，使用 iframe 预览
      previewUrl.value = `/api/documents/${row.id}/preview?token=${token}`
      previewData.value = null
      showPreviewDialog.value = true
    } else {
      ElMessage.warning(`不支持预览此文件类型 (${fileExt})，请下载查看`)
    }
  } catch (error) {
    console.error('预览错误:', error)
    previewError.value = '预览失败: ' + (error.message || '未知错误')
  } finally {
    previewLoading.value = false
  }
}

const handleDownload = (row) => {
  window.open(row.file_path, '_blank')
}

const downloadPreviewFile = () => {
  if (currentPreviewDocument.value) {
    window.open(currentPreviewDocument.value.file_path, '_blank')
  }
}

const handleDeleteDocument = async (row) => {
  await ElMessageBox.confirm('确定删除该文档?', '提示', { type: 'warning' })
  await deleteDocument(row.id)
  ElMessage.success('删除成功')
  loadDocuments()
}

const handleSubmitApproval = async () => {
  await approvalFormRef.value.validate()
  await createApproval({ contract_id: contractId.value, status: 'pending', comment: approvalForm.comment })
  ElMessage.success('提交成功')
  showApprovalDialog.value = false
  approvalForm.comment = ''
  loadApprovals()
}

const tabChange = (tab) => {
  if (tab === 'executions') loadExecutions()
  if (tab === 'documents') loadDocuments()
  if (tab === 'approvals') loadApprovals()
  if (tab === 'lifecycle') loadLifecycle()
}

onMounted(async () => {
  await loadContract()
  loadApprovals()
})
</script>

<style scoped>
.preview-container {
  width: 100%;
  height: 70vh;
  min-height: 400px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f5f5;
}

.preview-iframe {
  width: 100%;
  height: 100%;
  border: none;
}

.preview-loading,
.preview-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
  gap: 16px;
}

.preview-error {
  color: #f56c6c;
}

.preview-loading .el-icon,
.preview-error .el-icon {
  font-size: 48px;
}

.content-preview {
  height: 100%;
  overflow-y: auto;
  padding: 0;
}

.content-section {
  margin-top: 10px;
}

.document-content {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 4px;
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #e4e7ed;
}

.document-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  line-height: 1.6;
}
</style>
