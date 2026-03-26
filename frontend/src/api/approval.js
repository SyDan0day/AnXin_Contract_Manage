/**
 * 审批管理模块 API
 * 封装审批、提醒、通知相关的所有 API 请求
 */

import request from '@/utils/request'

/**
 * 获取合同的审批记录
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回审批记录列表
 */
export const getApprovalRecords = (contractId) => {
  return request({
    url: `/contracts/${contractId}/approvals`,
    method: 'get'
  })
}

/**
 * 获取合同的审批流程详情
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回审批流程详情（包含各级审批状态）
 */
export const getContractWorkflow = (contractId) => {
  return request({
    url: `/workflow/${contractId}`,
    method: 'get'
  })
}

/**
 * 获取待审批合同列表
 * @returns {Promise} 返回待审批的合同列表
 */
export const getPendingApprovals = () => {
  return request({
    url: '/pending-approvals',
    method: 'get'
  })
}

/**
 * 创建审批记录
 * @param {Object} data - 审批数据
 * @param {number} data.contract_id - 合同ID
 * @param {string} [data.status] - 审批状态 (pending/approved/rejected)
 * @param {string} [data.comment] - 审批意见
 * @returns {Promise} 返回创建的审批记录
 */
export const createApproval = (data) => {
  return request({
    url: `/contracts/${data.contract_id}/approvals`,
    method: 'post',
    data
  })
}

/**
 * 撤回审批申请
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回撤回结果
 */
export const withdrawApproval = (contractId) => {
  return request({
    url: `/contracts/${contractId}/approvals/withdraw`,
    method: 'post'
  })
}

/**
 * 更新审批状态
 * @param {number} id - 审批记录ID
 * @param {Object} data - 更新数据
 * @param {string} data.status - 新状态 (approved/rejected)
 * @param {string} [data.comment] - 审批意见
 * @returns {Promise} 返回更新结果
 */
export const updateApproval = (id, data) => {
  return request({
    url: `/approvals/${id}`,
    method: 'put',
    data
  })
}

/**
 * 获取合同的提醒列表
 * @param {number} contractId - 合同ID
 * @returns {Promise} 返回提醒列表
 */
export const getReminders = (contractId) => {
  return request({
    url: `/contracts/${contractId}/reminders`,
    method: 'get'
  })
}

/**
 * 创建合同提醒
 * @param {Object} data - 提醒数据
 * @param {number} data.contract_id - 合同ID
 * @param {string} data.reminder_date - 提醒日期 (YYYY-MM-DD)
 * @param {string} [data.message] - 提醒内容
 * @returns {Promise} 返回创建的提醒
 */
export const createReminder = (data) => {
  return request({
    url: `/contracts/${data.contract_id}/reminders`,
    method: 'post',
    data
  })
}

/**
 * 发送提醒
 * @param {number} id - 提醒ID
 * @returns {Promise} 返回发送结果
 */
export const sendReminder = (id) => {
  return request({
    url: `/reminders/${id}/send`,
    method: 'post'
  })
}

/**
 * 获取即将到期的合同列表
 * @param {number} [days=30] - 提前提醒天数
 * @returns {Promise} 返回即将到期的合同列表
 */
export const getExpiringContracts = (days = 30) => {
  return request({
    url: '/expiring-contracts',
    method: 'get',
    params: { days }
  })
}

/**
 * 获取统计数据
 * 用于仪表盘展示，包括合同数量、金额、状态分布等
 * @returns {Promise} 返回统计数据
 */
export const getStatistics = () => {
  return request({
    url: '/statistics',
    method: 'get'
  })
}

/**
 * 获取通知数量统计
 * 包括待审批数、待处理状态变更数、即将到期数
 * 用于侧边栏菜单红点提示
 * @returns {Promise} 返回各类型通知数量
 */
export const getNotificationCounts = () => {
  return request({
    url: '/notification-counts',
    method: 'get'
  })
}
