<template>
  <div class="customer-page">
    <el-card>
      <el-tabs v-model="activeTab" @tab-change="tabChange">
        <el-tab-pane label="客户管理" name="customers">
          <div class="toolbar">
            <el-form :inline="true" :model="searchForm">
              <el-form-item label="关键词">
                <el-input v-model="searchForm.keyword" placeholder="请输入客户编码或客户名称" clearable />
              </el-form-item>
              <el-form-item label="类型">
                <el-select v-model="searchForm.type" placeholder="请选择类型" clearable>
                  <el-option label="客户" value="customer" />
                  <el-option label="供应商" value="supplier" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleSearch">查询</el-button>
                <el-button @click="handleReset">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon> 新增客户
            </el-button>
          </div>

          <el-table :data="tableData" style="width: 100%" v-loading="loading">
            <el-table-column prop="code" label="客户编码" width="120" />
            <el-table-column prop="name" label="客户名称" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag>{{ row.type === 'customer' ? '客户' : '供应商' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="contact_person" label="联系人" width="100" />
            <el-table-column prop="contact_phone" label="联系电话" width="130" />
            <el-table-column prop="credit_rating" label="信用等级" width="100" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
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
        </el-tab-pane>

        <el-tab-pane label="合同类型" name="contractTypes">
          <div class="toolbar">
            <span></span>
            <el-button type="primary" @click="handleAddType">
              <el-icon><Plus /></el-icon> 新增类型
            </el-button>
          </div>

          <el-table :data="contractTypes" style="width: 100%" v-loading="typesLoading">
            <el-table-column prop="code" label="类型编码" width="150" />
            <el-table-column prop="name" label="类型名称" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="created_at" label="创建时间" width="180" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="warning" link @click="handleEditType(row)">
                    <el-icon><Edit /></el-icon> 编辑
                  </el-button>
                  <el-button type="danger" link @click="handleDeleteType(row)">
                    <el-icon><Delete /></el-icon> 删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" @close="handleDialogClose">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
        <el-form-item label="客户编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入客户编码" />
        </el-form-item>
        <el-form-item label="客户名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入客户名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio label="customer">客户</el-radio>
            <el-radio label="supplier">供应商</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="formData.contact_person" placeholder="请输入联系人" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="formData.contact_phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="formData.contact_email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="formData.address" type="textarea" :rows="2" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="信用等级">
          <el-select v-model="formData.credit_rating" placeholder="请选择信用等级" style="width: 100%">
            <el-option label="AAA" value="AAA" />
            <el-option label="AA" value="AA" />
            <el-option label="A" value="A" />
            <el-option label="B" value="B" />
            <el-option label="C" value="C" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="typeDialogVisible" :title="typeDialogTitle" width="500px" @close="handleTypeDialogClose">
      <el-form ref="typeFormRef" :model="typeForm" :rules="typeFormRules" label-width="100px">
        <el-form-item label="类型编码" prop="code">
          <el-input v-model="typeForm.code" placeholder="请输入类型编码" />
        </el-form-item>
        <el-form-item label="类型名称" prop="name">
          <el-input v-model="typeForm.name" placeholder="请输入类型名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="typeForm.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="typeDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleTypeSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { getCustomerList, createCustomer, updateCustomer, deleteCustomer, getContractTypeList, createContractType, updateContractType, deleteContractType } from '@/api/customer'

const activeTab = ref('customers')
const loading = ref(false)
const typesLoading = ref(false)
const dialogVisible = ref(false)
const typeDialogVisible = ref(false)
const dialogTitle = ref('')
const typeDialogTitle = ref('')
const formRef = ref(null)
const typeFormRef = ref(null)
const tableData = ref([])
const contractTypes = ref([])

const searchForm = reactive({
  keyword: '',
  type: ''
})

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const formData = reactive({
  code: '',
  name: '',
  type: 'customer',
  contact_person: '',
  contact_phone: '',
  contact_email: '',
  address: '',
  credit_rating: ''
})

const typeForm = reactive({
  code: '',
  name: '',
  description: ''
})

const formRules = {
  code: [{ required: true, message: '请输入客户编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入客户名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const typeFormRules = {
  code: [{ required: true, message: '请输入类型编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入类型名称', trigger: 'blur' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      skip: (pagination.page - 1) * pagination.size,
      limit: pagination.size,
      ...searchForm
    }
    const data = await getCustomerList(params)
    tableData.value = data
    pagination.total = data.length
  } finally {
    loading.value = false
  }
}

const loadContractTypes = async () => {
  typesLoading.value = true
  try {
    contractTypes.value = await getContractTypeList({ limit: 1000 })
  } finally {
    typesLoading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增客户'
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑客户'
  Object.assign(formData, row)
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该客户吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  await deleteCustomer(row.id)
  ElMessage.success('删除成功')
  loadData()
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  Object.assign(searchForm, { keyword: '', type: '' })
  handleSearch()
}

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (valid) {
      if (formData.id) {
        await updateCustomer(formData.id, formData)
        ElMessage.success('更新成功')
      } else {
        await createCustomer(formData)
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
    code: '',
    name: '',
    type: 'customer',
    contact_person: '',
    contact_phone: '',
    contact_email: '',
    address: '',
    credit_rating: ''
  })
}

const handleAddType = () => {
  typeDialogTitle.value = '新增合同类型'
  typeDialogVisible.value = true
}

const handleEditType = (row) => {
  typeDialogTitle.value = '编辑合同类型'
  Object.assign(typeForm, row)
  typeDialogVisible.value = true
}

const handleDeleteType = async (row) => {
  await ElMessageBox.confirm('确定要删除该合同类型吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  await deleteContractType(row.id)
  ElMessage.success('删除成功')
  loadContractTypes()
}

const handleTypeSubmit = async () => {
  await typeFormRef.value.validate(async (valid) => {
    if (valid) {
      if (typeForm.id) {
        await updateContractType(typeForm.id, typeForm)
        ElMessage.success('更新成功')
      } else {
        await createContractType(typeForm)
        ElMessage.success('创建成功')
      }
      typeDialogVisible.value = false
      loadContractTypes()
    }
  })
}

const handleTypeDialogClose = () => {
  typeFormRef.value?.resetFields()
  Object.assign(typeForm, { code: '', name: '', description: '' })
}

const tabChange = (tab) => {
  if (tab === 'customers') loadData()
  if (tab === 'contractTypes') loadContractTypes()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.customer-page {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 16px;
}

.toolbar .el-form {
  flex: 1;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 4px;
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