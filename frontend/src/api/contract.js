/**
 * 合同管理模块 API
 * 封装合同相关的所有 API 请求，包括合同 CRUD、执行跟踪、文档管理、生命周期等
 */

import request from '@/utils/request'

/**
 * 获取合同列表
 * @param {Object} [params] - 查询参数
 * @param {string} [params.keyword] - 关键词搜索（合同编号/标题）
 * @param {string} [params.status] - 合同状态筛选
 * @param {number} [params.skip] - 跳过记录数（分页）
 * @param {number} [params.limit] - 返回记录数限制
 * @returns {Promise} 返回合同列表
 */
export const getContractList = (params) => {
  return request({
    url: '/contracts',
    method: 'get',
    params
  })
}

/**
 * 获取合同详情
 * @param {number} id - 合同ID
 * @returns {Promise} 返回合同详细信息
 */
export const getContractDetail = (id) => {
  return request({
    url: `/contracts/${id}`,
    method: 'get'
  })
}

/**
 * 创建合同
 * @param {Object} data - 合同数据
 * @param {string} data.title - 合同标题
 * @param {number} data.customer_id - 客户ID
 * @param {number} data.contract_type_id - 合同类型ID
 * @param {number} data.amount - 合同金额
 * @param {string} [data.sign_date] - 签约日期 (YYYY-MM-DD)
 * @param {string} [data.start_date] - 开始日期 (YYYY-MM-DD)
 * @param {string} [data.end_date] - 结束日期 (YYYY-MM-DD)
 * @param {string} [data.payment_terms] - 付款条件
 * @param {string} [data.content] - 合同内容
 * @returns {Promise} 返回创建的合同信息
 */
export const createContract = (data) => {
  return request({
    url: '/contracts',
    method: 'post',
    data
  })
}

/**
 * 更新合同信息
 * @param {number} id - 合同ID
 * @param {Object} data - 更新的合同数据
 * @returns {Promise} 返回更新结果
 */
export const updateContract = (id, data) => {
  return request({
    url: `/contracts/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除合同
 * @param {number} id - 合同ID
 * @returns {Promise} 返回删除结果
 */
export const deleteContract = (id) => {
  return request({
    url: `/contracts/${id}`,
    method: 'delete'
  })
}

/**
 * 获取合同文档列表
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回文档列表
 */
export const getContractDocuments = (contractId) => {
  return request({
    url: `/contracts/${contractId}/documents`,
    method: 'get'
  })
}

/**
 * 上传合同文档
 * @param {Object} data - 文档数据 (FormData格式)
 * @param {number} data.contract_id - 合同ID
 * @param {File} data.file - 上传的文件
 * @param {string} [data.description] - 文档描述
 * @returns {Promise} 返回上传结果
 */
export const uploadDocument = (data) => {
  return request({
    url: `/contracts/${data.contract_id}/documents`,
    method: 'post',
    data
  })
}

/**
 * 删除文档
 * @param {number} id - 文档ID
 * @returns {Promise} 返回删除结果
 */
export const deleteDocument = (id) => {
  return request({
    url: `/documents/${id}`,
    method: 'delete'
  })
}

/**
 * 获取文档预览令牌
 * @param {number} documentId - 文档ID
 * @returns {Promise} 返回预览令牌
 */
export const getPreviewToken = (documentId) => {
  return request({
    url: `/documents/${documentId}/preview-token`,
    method: 'post'
  })
}

/**
 * 获取合同生命周期记录
 * 记录合同的所有状态变更历史
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回生命周期事件列表
 */
export const getContractLifecycle = (contractId) => {
  return request({
    url: `/contracts/${contractId}/lifecycle`,
    method: 'get'
  })
}

/**
 * 直接更新合同状态（无需审批的状态）
 * @param {number} contractId - 合同ID
 * @param {Object} data - 状态数据
 * @param {string} data.status - 新状态
 * @returns {Promise} 返回更新结果
 */
export const updateContractStatus = (contractId, data) => {
  return request({
    url: `/contracts/${contractId}/status`,
    method: 'put',
    data
  })
}

/**
 * 归档合同
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回归档结果
 */
export const archiveContract = (contractId) => {
  return request({
    url: `/contracts/${contractId}/archive`,
    method: 'post'
  })
}

/**
 * 申请状态变更（需要审批的状态变更）
 * 需要审批的状态：archived(归档)、terminated(终止)、in_progress(执行中)、pending_pay(待付款)
 * @param {number} contractId - 合同ID
 * @param {Object} data - 变更申请数据
 * @param {string} data.to_status - 目标状态
 * @param {string} data.reason - 变更原因
 * @returns {Promise} 返回申请结果，包含是否需要审批
 */
export const requestStatusChange = (contractId, data) => {
  return request({
    url: `/contracts/${contractId}/status-change`,
    method: 'post',
    data
  })
}

/**
 * 获取合同状态变更申请记录
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回状态变更申请记录列表
 */
export const getStatusChangeRequests = (contractId) => {
  return request({
    url: `/contracts/${contractId}/status-change`,
    method: 'get'
  })
}

/**
 * 获取待审批的状态变更列表
 * @returns {Promise} 返回待审批的状态变更列表
 */
export const getPendingStatusChangeApprovals = () => {
  return request({
    url: '/pending-status-changes',
    method: 'get'
  })
}

/**
 * 审批通过状态变更申请
 * @param {number} requestId - 状态变更申请ID
 * @param {Object} [data] - 审批意见
 * @param {string} [data.comment] - 审批备注
 * @returns {Promise} 返回审批结果
 */
export const approveStatusChangeRequest = (requestId, data) => {
  return request({
    url: `/status-change-requests/${requestId}/approve`,
    method: 'post',
    data
  })
}

/**
 * 拒绝状态变更申请
 * @param {number} requestId - 状态变更申请ID
 * @param {Object} [data] - 拒绝原因
 * @param {string} [data.comment] - 拒绝原因说明
 * @returns {Promise} 返回拒绝结果
 */
export const rejectStatusChangeRequest = (requestId, data) => {
  return request({
    url: `/status-change-requests/${requestId}/reject`,
    method: 'post',
    data
  })
}
