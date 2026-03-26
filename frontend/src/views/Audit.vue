<template>
  <div class="audit-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>审计日志</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleExport">
              <el-icon><Download /></el-icon> 导出
            </el-button>
            <el-button 
              v-if="isAuditAdmin" 
              type="danger" 
              :disabled="selectedLogs.length === 0"
              @click="handleBatchDelete"
            >
              <el-icon><Delete /></el-icon> 批量删除
            </el-button>
          </div>
        </div>
      </template>

      <div class="search-form">
        <el-form :inline="true" :model="searchForm">
          <el-form-item label="用户名">
            <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
          </el-form-item>
          <el-form-item label="操作">
            <el-input v-model="searchForm.action" placeholder="请输入操作" clearable />
          </el-form-item>
          <el-form-item label="模块">
            <el-select v-model="searchForm.module" placeholder="请选择模块" clearable>
              <el-option label="认证" value="auth" />
              <el-option label="合同" value="contract" />
              <el-option label="客户" value="customer" />
              <el-option label="审批" value="approval" />
              <el-option label="提醒" value="reminder" />
              <el-option label="用户" value="user" />
              <el-option label="统计" value="statistics" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-form-item>
          <el-form-item label="日期">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table 
        :data="tableData" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="index" label="序号" width="60" align="center" />
        <el-table-column v-if="isAuditAdmin" type="selection" width="55" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="module" label="模块" width="80">
          <template #default="{ row }">
            <el-tag>{{ getModuleText(row.module) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作" min-width="200" />
        <el-table-column prop="method" label="方法" width="80">
          <template #default="{ row }">
            <el-tag :type="getMethodType(row.method)" size="small">{{ row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路径" min-width="180" show-overflow-tooltip />
        <el-table-column prop="ip_address" label="IP地址" width="130" />
        <el-table-column prop="status_code" label="状态码" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status_code >= 400 ? 'danger' : 'success'" size="small">
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column v-if="isAuditAdmin" label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[20, 50, 100, 200]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'
import { getAuditLogs, deleteAuditLog, deleteAuditLogs } from '@/api/audit'
import { Download, Delete } from '@element-plus/icons-vue'

const userStore = useUserStore()
const isAuditAdmin = computed(() => {
  return userStore.userInfo?.role === 'admin' || userStore.userInfo?.role === 'audit_admin'
})

const loading = ref(false)
const tableData = ref([])
const selectedLogs = ref([])
const dateRange = ref([])

const searchForm = reactive({
  username: '',
  action: '',
  module: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const getModuleText = (module) => {
  const map = {
    auth: '认证',
    contract: '合同',
    customer: '客户',
    approval: '审批',
    reminder: '提醒',
    user: '用户',
    statistics: '统计',
    other: '其他'
  }
  return map[module] || module
}

const getMethodType = (method) => {
  const map = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info'
  }
  return map[method] || 'info'
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

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      username: searchForm.username || undefined,
      action: searchForm.action || undefined,
      module: searchForm.module || undefined,
      start_date: dateRange.value?.[0] || undefined,
      end_date: dateRange.value?.[1] || undefined
    }
    const res = await getAuditLogs(params)
    tableData.value = res.logs || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('加载失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  searchForm.username = ''
  searchForm.action = ''
  searchForm.module = ''
  dateRange.value = []
  handleSearch()
}

const handlePageChange = (page) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  loadData()
}

const handleSelectionChange = (selection) => {
  selectedLogs.value = selection
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该日志记录吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteAuditLog(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedLogs.value.length} 条日志记录吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const ids = selectedLogs.value.map(log => log.id)
    await deleteAuditLogs(ids)
    ElMessage.success('批量删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

const handleExport = () => {
  const params = new URLSearchParams({
    username: searchForm.username || '',
    action: searchForm.action || '',
    module: searchForm.module || '',
    start_date: dateRange.value?.[0] || '',
    end_date: dateRange.value?.[1] || ''
  })
  window.open(`/api/audit-logs/export?${params.toString()}`, '_blank')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.audit-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.card-header span {
  font-weight: 600;
  font-size: 15px;
  color: #1E293B;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.header-actions .el-button {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
