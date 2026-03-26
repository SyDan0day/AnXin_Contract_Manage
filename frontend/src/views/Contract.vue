<template>
  <div class="contract-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>合同管理</span>
          <el-button type="primary" @click="handleAdd" v-if="canCreateContract">
            <el-icon><Plus /></el-icon> 新增合同
          </el-button>
        </div>
      </template>
       <el-form :inline="true" :model="searchForm">
         <el-form-item label="关键词">
           <el-input v-model="searchForm.keyword" placeholder="请输入合同编号或合同标题" clearable @input="handleKeywordInput" />
         </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="请选择状态" clearable @change="handleStatusChange">
              <el-option label="草稿" value="draft" />
              <el-option label="待审批" value="pending" />
              <el-option label="已生效" value="active" />
              <el-option label="已完成" value="completed" />
              <el-option label="已终止" value="terminated" />
              <el-option label="已归档" value="archived" />
            </el-select>
          </el-form-item>
         <el-form-item>
           <el-button type="primary" @click="handleSearch">查询</el-button>
           <el-button @click="handleReset">重置</el-button>
         </el-form-item>
       </el-form>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" align="center" />
        <el-table-column prop="contract_no" label="合同编号" width="150" />
        <el-table-column prop="title" label="合同标题" />
        <el-table-column prop="creator" label="销售负责人" width="120">
          <template #default="{ row }">
            {{ row.creator?.full_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态信息" width="260">
          <template #default="{ row }">
            <!-- 归档状态 -->
            <div v-if="row.status === 'archived'" class="status-info archived-info">
              <div class="status-badge archived-badge">
                <el-icon class="status-icon"><FolderOpened /></el-icon>
                <span>已归档</span>
              </div>
              <div class="status-detail">
                <span class="status-hint">合同已归档保存</span>
              </div>
            </div>
            
            <!-- 审批中状态 -->
            <div v-else-if="row.status === 'pending'" class="status-info pending-info">
              <div class="status-badge pending-badge">
                <el-icon class="status-icon"><Clock /></el-icon>
                <span>审批中</span>
              </div>
              <div class="status-detail">
                <el-progress 
                  :percentage="getApprovalProgress(row)" 
                  :color="['#F59E0B', '#EAB308', '#22C55E']"
                  :stroke-width="6"
                  :show-text="false"
                  style="width: 100%"
                />
                <span class="status-hint">等待{{ row.current_approver || '审批' }}审批中...</span>
              </div>
            </div>
            
            <!-- 已拒绝状态 -->
            <div v-else-if="row.rejection_info" class="status-info rejection-info">
              <div class="status-badge rejection-badge">
                <el-icon class="status-icon"><WarningFilled /></el-icon>
                <span>已拒绝</span>
              </div>
              <div class="status-detail">
                <div class="status-row">
                  <span class="label">拒绝人:</span>
                  <span class="value">{{ row.rejection_info.approver_name || '-' }}</span>
                </div>
                <el-tooltip 
                  :content="row.rejection_info.comment || '无'"
                  placement="top"
                  :show-after="200"
                >
                  <div class="status-row">
                    <span class="label">原因:</span>
                    <span class="value text-ellipsis">{{ row.rejection_info.comment || '-' }}</span>
                  </div>
                </el-tooltip>
              </div>
            </div>
            
            <!-- 其他状态显示横杠 -->
            <span v-else class="text-gray">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="sign_date" label="签约日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.sign_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="end_date" label="到期日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" link @click="handleView(row)">
                <el-icon><View /></el-icon> 详情
              </el-button>
              <el-button type="warning" link @click="handleEdit(row)">
                <el-icon><Edit /></el-icon> 编辑
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)">
                <el-icon><Delete /></el-icon> 删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      @close="handleDialogClose"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
        <el-form-item label="合同标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入合同标题" />
        </el-form-item>
        <el-form-item label="客户" prop="customer_id">
          <el-select v-model="formData.customer_id" placeholder="请选择客户" style="width: 100%">
            <el-option
              v-for="customer in customers"
              :key="customer.id"
              :label="customer.name"
              :value="customer.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="合同类型" prop="contract_type_id">
          <el-select v-model="formData.contract_type_id" placeholder="请选择合同类型" style="width: 100%">
            <el-option
              v-for="type in contractTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="formData.amount" :precision="2" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="签约日期" prop="sign_date">
          <el-date-picker
            v-model="formData.sign_date"
            type="date"
            placeholder="请选择签约日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="开始日期" prop="start_date">
          <el-date-picker
            v-model="formData.start_date"
            type="date"
            placeholder="请选择开始日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束日期" prop="end_date">
          <el-date-picker
            v-model="formData.end_date"
            type="date"
            placeholder="请选择结束日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="付款条件" prop="payment_terms">
          <el-input
            v-model="formData.payment_terms"
            type="textarea"
            :rows="3"
            placeholder="请输入付款条件"
          />
        </el-form-item>
        <el-form-item label="合同内容" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="5"
            placeholder="请输入合同内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, View, Edit, Delete, WarningFilled, FolderOpened, Clock } from '@element-plus/icons-vue'
import { getContractList, getContractDetail, createContract, updateContract, deleteContract } from '@/api/contract'
import { getCustomerList } from '@/api/customer'
import { getContractTypeList } from '@/api/customer'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const userRole = computed(() => userStore.userInfo?.role || 'user')

const canCreateContract = computed(() => {
  return userRole.value === 'user' || userRole.value === 'admin'
})

const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const tableData = ref([])
const customers = ref([])
const contractTypes = ref([])

const searchForm = reactive({
  keyword: '',
  status: ''
})

// 防抖搜索
let searchDebounceTimer = null
const handleDebouncedSearch = () => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
  }
  searchDebounceTimer = setTimeout(() => {
    handleSearch()
  }, 300) // 300ms debounce
}

const handleKeywordInput = () => {
  handleDebouncedSearch()
}

const handleStatusChange = () => {
  handleDebouncedSearch()
}

const onBeforeUnmount = () => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
  }
}

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const formData = reactive({
  title: '',
  customer_id: null,
  contract_type_id: null,
  amount: null,
  sign_date: '',
  start_date: '',
  end_date: '',
  payment_terms: '',
  content: ''
})

const formRules = {
  title: [{ required: true, message: '请输入合同标题', trigger: 'blur' }],
  customer_id: [{ required: true, message: '请选择客户', trigger: 'change' }],
  contract_type_id: [{ required: true, message: '请选择合同类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入合同金额', trigger: 'blur' }, { validator: (rule, value) => {
        if (value === '' || value === null) {
          return new Error('请输入合同金额');
        }
        if (isNaN(value) || parseFloat(value) <= 0) {
          return new Error('合同金额必须大于0');
        }
        return true;
      }, trigger: 'blur' }],
  sign_date: [{ required: true, message: '请选择签约日期', trigger: 'change' } ],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' } ],
  end_date: [{ required: true, message: '请选择结束日期', trigger: 'change' }, { validator: (rule, value) => {
        if (!formData.start_date || !value) {
          return true; // 让required规则处理空值
        }
        const startDate = new Date(formData.start_date);
        const endDate = new Date(value);
        if (endDate <= startDate) {
          return new Error('结束日期必须晚于开始日期');
        }
        return true;
      }, trigger: 'blur' }],
  payment_terms: [{ validator: (rule, value) => {
        if (value !== '' && value.length > 500) {
          return new Error('付款条件不能超过500个字符');
        }
        return true;
      }, trigger: 'blur' }],
  content: [{ validator: (rule, value) => {
        if (value !== '' && value.length > 2000) {
          return new Error('合同内容不能超过2000个字符');
        }
        return true;
      }, trigger: 'blur' }]
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

const getStatusType = (status) => {
  const typeMap = {
    draft: 'info',
    pending: 'warning',
    active: 'primary',
    completed: 'success',
    terminated: 'danger',
    archived: 'info'
  }
  return typeMap[status] || ''
}

const getStatusText = (status) => {
  const textMap = {
    draft: '草稿',
    pending: '待审批',
    active: '已生效',
    completed: '已完成',
    terminated: '已终止',
    archived: '已归档'
  }
  return textMap[status] || status
}

const getApprovalProgress = (row) => {
  // 使用后端返回的审批级别数据
  const currentLevel = row.current_approval_level || 1
  const maxLevel = row.max_approval_level || 3
  return Math.round((currentLevel / maxLevel) * 100)
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      skip: (pagination.page - 1) * pagination.size,
      limit: pagination.size,
      ...searchForm
    }
    const res = await getContractList(params)
    // 支持新格式 {items: [], total: 100} 和旧格式 []
    tableData.value = res.items || res
    pagination.total = res.total || (Array.isArray(res) ? res.length : 0)
  } finally {
    loading.value = false
  }
}

const loadCustomers = async () => {
  customers.value = await getCustomerList({ limit: 1000 })
}

const loadContractTypes = async () => {
  contractTypes.value = await getContractTypeList({ limit: 1000 })
}

const handleAdd = () => {
  dialogTitle.value = '新增合同'
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑合同'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleView = (row) => {
  router.push(`/contracts/${row.id}`)
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该合同吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  await deleteContract(row.id)
  ElMessage.success('删除成功')
  loadData()
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  Object.assign(searchForm, { keyword: '', status: '' })
  handleSearch()
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (valid) {
      if (formData.id) {
        await updateContract(formData.id, formData)
        ElMessage.success('更新成功')
      } else {
        await createContract(formData)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    }
  })
}

const handleDialogClose = () => {
  formRef.value?.resetFields()
  Object.assign(formData, {
    title: '',
    customer_id: null,
    contract_type_id: null,
    amount: null,
    sign_date: '',
    start_date: '',
    end_date: '',
    payment_terms: '',
    content: ''
  })
}

onMounted(async () => {
  if (route.query.status) {
    searchForm.status = route.query.status
  }
  if (route.query.keyword) {
    searchForm.keyword = route.query.keyword
  }
  loadData()
  loadCustomers()
  loadContractTypes()
  
  if (route.query.action === 'create') {
    handleAdd()
    window.history.replaceState({}, '', '/contracts')
  } else if (route.query.action === 'edit' && route.query.id) {
    const id = parseInt(route.query.id)
    const data = await getContractDetail(id)
    handleEdit(data)
    window.history.replaceState({}, '', '/contracts')
  }
})
</script>

<style scoped>
.contract-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-wrap: nowrap;
  justify-content: flex-end;
}

.rejection-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px;
  background: linear-gradient(135deg, #FEF2F2 0%, #FEF2F2 100%);
  border-radius: 8px;
  border: 1px solid #FECACA;
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px;
  border-radius: 8px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  width: fit-content;
}

.status-icon {
  font-size: 12px;
}

.status-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
}

.status-row {
  display: flex;
  gap: 4px;
}

.status-row .label {
  color: #9CA3AF;
  flex-shrink: 0;
}

.status-row .value {
  color: #374151;
  font-weight: 500;
}

.status-hint {
  color: #6B7280;
  font-size: 11px;
}

/* 归档样式 */
.archived-info {
  background: linear-gradient(135deg, #F3F4F6 0%, #F9FAFB 100%);
  border: 1px solid #E5E7EB;
}

.archived-badge {
  background: #6B7280;
  color: white;
}

/* 审批中样式 */
.pending-info {
  background: linear-gradient(135deg, #FFFBEB 0%, #FEF3C7 100%);
  border: 1px solid #FDE68A;
}

.pending-badge {
  background: #F59E0B;
  color: white;
}

/* 拒绝样式 - 重命名避免冲突 */
.rejection-info {
  background: linear-gradient(135deg, #FEF2F2 0%, #FEF2F2 100%);
  border: 1px solid #FECACA;
}

.rejection-badge {
  background: #EF4444;
  color: white;
}

.rejection-badge .status-icon {
  color: #FCA5A5;
}

.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 150px;
}

.text-gray {
  color: #9CA3AF;
}

.action-buttons .el-button {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>